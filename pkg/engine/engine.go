package engine

import (
	"image/color"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/yoga"
)

// defaultEngine 是当前运行的全局 Engine 实例，供 Rebuild() 等快捷函数使用。
var defaultEngine *Engine

// DefaultEngine 返回当前全局 Engine 实例。
func DefaultEngine() *Engine {
	return defaultEngine
}

// Engine is the core engine of Tenon v2, managing Widget, Element, and RenderObject trees.
type Engine struct {
	buildFunc BuildFunc

	rootWidget       Widget
	rootElement      Element
	rootRenderObject render.RenderObject

	owner *render.PipelineOwner

	screenWidth  int
	screenHeight int

	needsBuild    bool
	dirtyElements []Element
	building      bool

	activeAnimations []*AnimationController
	lastTickTime     time.Time

	lastMouseX      float32
	lastMouseY      float32
	mouseDown       bool
	mouseDownTarget render.RenderObject
	hoverTarget     render.RenderObject
	focusTarget     render.RenderObject

	// 拖放状态
	dragActive      bool
	dragData        any
	dragSource      render.RenderObject
	dragDropTarget  render.RenderObject

	// 焦点管理
	focusManager *FocusManager
	// 快捷键管理
	shortcutManager *ShortcutManager

	// 调试工具
	perfOverlay *PerformanceOverlay
	inspector    *Inspector

	// 可复用缓冲区（避免热路径分配）
	pathBuf       []render.RenderObject // getAbsolutePosition 用
	childrenBuf   []render.RenderObject // drawRenderObject 排序用
}

// NewEngine creates a new Engine.
func NewEngine(buildFunc BuildFunc, width, height int) *Engine {
	e := &Engine{
		buildFunc:    buildFunc,
		screenWidth:  width,
		screenHeight: height,
		owner:        render.NewPipelineOwner(),
		focusManager: NewFocusManager(),
		perfOverlay:  NewPerformanceOverlay(),
	}
	defaultEngine = e
	return e
}

// Rebuild marks the engine to rebuild on the next frame.
func (e *Engine) Rebuild() {
	e.needsBuild = true
}

// RebuildDefault marks the default engine to rebuild on the next frame.
func RebuildDefault() {
	if defaultEngine != nil {
		defaultEngine.Rebuild()
	}
}

// scheduleBuildFor 将指定 Element 加入 dirty 队列，供局部 rebuild 使用。
func (e *Engine) scheduleBuildFor(element Element) {
	if e.building {
		// 如果在 build 过程中触发，直接加入 dirty 队列，下一轮处理
		e.dirtyElements = append(e.dirtyElements, element)
		return
	}
	e.dirtyElements = append(e.dirtyElements, element)
}

func (e *Engine) GetRootRenderObject() render.RenderObject {
	return e.rootRenderObject
}

// GetFocusManager 返回焦点管理器。
func (e *Engine) GetFocusManager() *FocusManager {
	return e.focusManager
}

// GetPerformanceOverlay 返回性能覆盖层。
func (e *Engine) GetPerformanceOverlay() *PerformanceOverlay {
	return e.perfOverlay
}

// GetShortcutManager 返回快捷键管理器。
func (e *Engine) GetShortcutManager() *ShortcutManager {
	if e.shortcutManager == nil {
		e.shortcutManager = NewShortcutManager()
	}
	return e.shortcutManager
}
func (e *Engine) GetInspector() *Inspector {
	if e.inspector == nil {
		e.inspector = NewInspector(e)
	}
	return e.inspector
}

func (e *Engine) GetRootElement() Element {
	return e.rootElement
}

// Mount mounts the root widget and triggers the first build.
func (e *Engine) Mount() {
	e.needsBuild = true
	e.flushBuild()
}

// ==================== ebiten.Game interface ====================

func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != e.screenWidth || outsideHeight != e.screenHeight {
		e.screenWidth = outsideWidth
		e.screenHeight = outsideHeight
		if e.rootRenderObject != nil {
			e.rootRenderObject.MarkNeedsLayout()
		}
	}
	return outsideWidth, outsideHeight
}

func (e *Engine) RegisterAnimation(a *AnimationController) {
	for _, existing := range e.activeAnimations {
		if existing == a {
			return
		}
	}
	e.activeAnimations = append(e.activeAnimations, a)
}

func (e *Engine) UnregisterAnimation(a *AnimationController) {
	for i, existing := range e.activeAnimations {
		if existing == a {
			e.activeAnimations = append(e.activeAnimations[:i], e.activeAnimations[i+1:]...)
			return
		}
	}
}

func (e *Engine) Update() error {
	now := time.Now()
	if e.lastTickTime.IsZero() {
		e.lastTickTime = now
	}
	dt := now.Sub(e.lastTickTime)
	e.lastTickTime = now

	// Tick active animations
	if len(e.activeAnimations) > 0 {
		for _, a := range e.activeAnimations {
			a.Tick(dt)
		}
	}

	if e.needsBuild || len(e.dirtyElements) > 0 {
		e.flushBuild()
	}

	// Tick all tickable render objects
	if e.rootRenderObject != nil {
		e.tickRenderObjects(e.rootRenderObject, dt)
	}

	if e.rootRenderObject != nil && e.rootRenderObject.GetYoga() != nil {
		rootYoga := e.rootRenderObject.GetYoga()
		rootYoga.StyleSetWidth(float32(e.screenWidth))
		rootYoga.StyleSetHeight(float32(e.screenHeight))
		rootYoga.CalculateLayout(float32(e.screenWidth), float32(e.screenHeight), yoga.DirectionLTR)
	}

	e.owner.FlushLayout()

	if e.rootRenderObject != nil {
		e.syncBounds(e.rootRenderObject)
	}

	e.owner.FlushPaint()

	e.handleMouseInput()

	// 优先使用 textinput API 处理 IME 输入（中文/日文/韩文输入法）
	textInputHandled := false
	if e.focusTarget != nil {
		if h, ok := tryGetTextInputHandler(e.focusTarget); ok {
			absX, absY := e.getAbsolutePosition(e.focusTarget)
			textInputHandled = h.UpdateTextInput(int(absX), int(absY))
		}
	}

	if !textInputHandled {
		e.handleKeyboardInput()
	}

	e.tickBlink()

	// 记录帧统计
	if e.perfOverlay != nil {
		e.perfOverlay.RecordFrame(0, 0, 0)
	}

	return nil
}

func (e *Engine) handleMouseInput() {
	mx, my := ebiten.CursorPosition()
	e.lastMouseX = float32(mx)
	e.lastMouseY = float32(my)

	// Wheel
	wheelX, wheelY := ebiten.Wheel()
	if wheelY != 0 || wheelX != 0 {
		if e.rootRenderObject != nil {
			target := e.hitTest(e.rootRenderObject, e.lastMouseX, e.lastMouseY)
			if target != nil {
				if s, ok := scrollParent(target); ok {
					s.ScrollBy(-float32(wheelX)*30, -float32(wheelY)*30)
				}
			}
		}
	}

	// Hover
	if e.rootRenderObject != nil {
		target := e.hitTest(e.rootRenderObject, e.lastMouseX, e.lastMouseY)
		if target != e.hoverTarget {
			if e.hoverTarget != nil && e.hoverTarget.GetOnMouseLeave() != nil {
				e.hoverTarget.GetOnMouseLeave()()
			}
			if target != nil && target.GetOnMouseEnter() != nil {
				target.GetOnMouseEnter()()
			}
			e.hoverTarget = target
		}
	}

	isPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if isPressed && !e.mouseDown {
		// Mouse down
		e.mouseDown = true
		if e.rootRenderObject != nil {
			target := e.hitTest(e.rootRenderObject, e.lastMouseX, e.lastMouseY)

			e.mouseDownTarget = target
			if target != nil {
				// 启动拖放：如果 target 携带 dragData
				if data := target.GetDragData(); data != nil {
					e.dragActive = true
					e.dragData = data
					e.dragSource = target
				}
				if target.GetOnMouseDown() != nil {
					target.GetOnMouseDown()()
				}
			}
		}
	} else if isPressed && e.mouseDown {
		// Drag
		if e.mouseDownTarget != nil {
			e.mouseDownTarget.HandleDrag(e.lastMouseX, e.lastMouseY)
		}
		if e.dragActive && e.rootRenderObject != nil {
			// 查找当前鼠标下的 DropTarget（排除 dragSource 自身）
			target := e.hitTest(e.rootRenderObject, e.lastMouseX, e.lastMouseY)
			if target == e.dragSource {
				target = nil
			}
			if target != e.dragDropTarget {
				if e.dragDropTarget != nil && e.dragDropTarget.GetOnDragLeave() != nil {
					e.dragDropTarget.GetOnDragLeave()()
				}
				if target != nil && target.GetOnDragEnter() != nil {
					target.GetOnDragEnter()(e.dragData)
				}
				e.dragDropTarget = target
			}
		}
	} else if !isPressed && e.mouseDown {
		// Mouse up / click
		e.mouseDown = false

		// 处理拖放结束
		if e.dragActive {
			if e.dragDropTarget != nil && e.dragDropTarget.GetOnDrop() != nil {
				e.dragDropTarget.GetOnDrop()(e.dragData)
			}
			if e.dragDropTarget != nil && e.dragDropTarget.GetOnDragLeave() != nil {
				e.dragDropTarget.GetOnDragLeave()()
			}
			e.dragActive = false
			e.dragData = nil
			e.dragSource = nil
			e.dragDropTarget = nil
		}

		if e.mouseDownTarget != nil && e.rootRenderObject != nil {
			target := e.hitTest(e.rootRenderObject, e.lastMouseX, e.lastMouseY)
			e.updateFocus(target)
			if target != nil && target == e.mouseDownTarget {
				// 事件冒泡：从 target 向上查找第一个有 OnClick 的祖先
				clickTarget := target
				for clickTarget != nil && clickTarget.GetOnClick() == nil {
					if p := clickTarget.GetParent(); p != nil {
						clickTarget = p
					} else {
						break
					}
				}
				if clickTarget != nil && clickTarget.GetOnClick() != nil {
					clickTarget.GetOnClick()()
				}
			}
			if e.mouseDownTarget.GetOnMouseUp() != nil {
				e.mouseDownTarget.GetOnMouseUp()()
			}
		}
		e.mouseDownTarget = nil
	}
}

func (e *Engine) updateFocus(target render.RenderObject) {
	focuser := e.findFocuser(target)
	if focuser != nil {
		if e.focusTarget != focuser {
			if e.focusTarget != nil {
				if f, ok := e.focusTarget.(Focusabler); ok {
					f.Blur()
				}
			}
			if f, ok := focuser.(Focusabler); ok {
				f.Focus()
			}
			e.focusTarget = focuser
			DismissAllPopups()
		}
	} else {
		if e.focusTarget != nil {
			if f, ok := e.focusTarget.(Focusabler); ok {
				f.Blur()
			}
			e.focusTarget = nil
		}
		DismissAllPopups()
	}
}

func (e *Engine) findFocuser(ro render.RenderObject) render.RenderObject {
	// 从当前节点向上搜索最近的可聚焦祖先
	for ro != nil {
		if f, ok := ro.(Focusabler); ok {
			// 跳过显式标记为不可聚焦的 RenderBox
			if !f.IsFocusable() {
				ro = ro.GetParent()
				continue
			}
			return f
		}
		ro = ro.GetParent()
	}
	return nil
}

// getAbsolutePosition 计算 RenderObject 在屏幕上的绝对视觉位置（考虑 scroll offset）。
func (e *Engine) getAbsolutePosition(ro render.RenderObject) (float32, float32) {
	var x, y float32
	// 收集从根节点到当前节点的路径（复用缓冲区）
	e.pathBuf = e.pathBuf[:0]
	for ro != nil {
		e.pathBuf = append(e.pathBuf, ro)
		ro = ro.GetParent()
	}
	// 从根到当前节点累加
	for i := len(e.pathBuf) - 1; i >= 0; i-- {
		node := e.pathBuf[i]
		b := node.GetBounds()
		x += b.X
		y += b.Y
		if s, ok := tryGetScroll(node); ok {
			so := s.GetScrollOffset()
			x -= so.X
			y -= so.Y
		}
	}
	return x, y
}

func (e *Engine) handleKeyboardInput() {
	// F12 切换检查器
	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		e.GetInspector().Toggle()
		return
	}
	// F11 切换性能覆盖层
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		e.perfOverlay.Toggle()
		return
	}

	// 快捷键
	if e.shortcutManager != nil && e.shortcutManager.Handle() {
		return
	}

	// Tab 键切换焦点
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			e.focusManager.PreviousFocus()
		} else {
			e.focusManager.NextFocus()
		}
		return
	}

	if e.focusTarget == nil {
		return
	}
	f, ok := tryGetInputHandler(e.focusTarget)
	if !ok {
		return
	}

	chars := ebiten.AppendInputChars(nil)

	backspace := inpututil.IsKeyJustPressed(ebiten.KeyBackspace)
	if inpututil.KeyPressDuration(ebiten.KeyBackspace) > 30 && inpututil.KeyPressDuration(ebiten.KeyBackspace)%5 == 0 {
		backspace = true
	}

	enter := inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	left := inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	right := inpututil.IsKeyJustPressed(ebiten.KeyRight)
	if inpututil.KeyPressDuration(ebiten.KeyLeft) > 30 && inpututil.KeyPressDuration(ebiten.KeyLeft)%5 == 0 {
		left = true
	}
	if inpututil.KeyPressDuration(ebiten.KeyRight) > 30 && inpututil.KeyPressDuration(ebiten.KeyRight)%5 == 0 {
		right = true
	}

	home := inpututil.IsKeyJustPressed(ebiten.KeyHome)
	end := inpututil.IsKeyJustPressed(ebiten.KeyEnd)
	selectAll := inpututil.IsKeyJustPressed(ebiten.KeyA) && ebiten.IsKeyPressed(ebiten.KeyControl)

	// 剪贴板操作
	ctrl := ebiten.IsKeyPressed(ebiten.KeyControl)
	if ctrl {
		if ch, ok := tryGetClipboardHandler(e.focusTarget); ok {
			clipboard := GetClipboard()
			if inpututil.IsKeyJustPressed(ebiten.KeyC) {
				ch.Copy(clipboard.Write)
				return
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyX) {
				ch.Cut(clipboard.Write)
				return
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyV) {
				ch.Paste(clipboard.Read)
				return
			}
		}
	}

	if len(chars) > 0 || backspace || enter || left || right || home || end || selectAll {
		f.HandleInput(chars, backspace, enter, left, right, home, end, selectAll)
	}
}

func (e *Engine) tickBlink() {
	if e.focusTarget == nil {
		return
	}
	if b, ok := tryGetBlinker(e.focusTarget); ok {
		b.TickBlink()
	}
}

func (e *Engine) hitTest(ro render.RenderObject, x, y float32) render.RenderObject {
	if ro == nil || !ro.IsVisible() {
		return nil
	}

	scrollX, scrollY := float32(0), float32(0)
	if s, ok := tryGetScroll(ro); ok {
		so := s.GetScrollOffset()
		scrollX, scrollY = so.X, so.Y
	}

	bounds := ro.GetBounds()

	// 子节点按 z-index 降序排列（z-index 相同时后绘制的优先），
	// 确保视觉上层的元素优先接收事件。
	children := ro.GetChildren()
	if len(children) > 1 {
		type childInfo struct {
			ro  render.RenderObject
			idx int
			z   int
		}
		infos := make([]childInfo, len(children))
		for i, c := range children {
			infos[i] = childInfo{ro: c, idx: i, z: e.getSubtreeMaxZIndex(c)}
		}
		sort.Slice(infos, func(i, j int) bool {
			if infos[i].z != infos[j].z {
				return infos[i].z > infos[j].z
			}
			return infos[i].idx > infos[j].idx
		})
		children = make([]render.RenderObject, len(infos))
		for i, info := range infos {
			children[i] = info.ro
		}
	}
	for _, child := range children {
		if result := e.hitTest(child, x-bounds.X+scrollX, y-bounds.Y+scrollY); result != nil {
			return result
		}
	}

	// 然后检查当前节点本身
	if ro.HitTest(x, y) {
		return ro
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	if e.rootRenderObject != nil {
		e.drawRenderObject(screen, e.rootRenderObject, render.Offset{})
	}
	// 性能覆盖层
	if e.perfOverlay != nil {
		e.perfOverlay.Paint(screen, render.Offset{})
	}
	// 检查器覆盖层
	if e.inspector != nil {
		e.inspector.Update()
		e.inspector.Paint(screen, render.Offset{})
	}
}

// ==================== Widget diff & build ====================

func (e *Engine) flushBuild() {
	if e.building {
		return
	}
	e.building = true
	defer func() { e.building = false }()

	// 循环处理，因为 rebuild 可能产生新的 dirty elements
	// 限制循环次数防止无限递归
	for loop := 0; loop < 10; loop++ {
		hasWork := len(e.dirtyElements) > 0 || e.needsBuild
		if !hasWork {
			break
		}

		// 1. 处理局部 dirty elements（StatefulElement / fcElement 的 setState）
		dirtyList := e.dirtyElements
		e.dirtyElements = nil
		for _, el := range dirtyList {
			switch r := el.(type) {
			case *StatefulElement:
				r.rebuild()
			case *fcElement:
				r.PerformRebuild(nil)
			}
		}

		// 2. 处理全局 rebuild（buildFunc + Rebuild）
		if e.needsBuild {
			e.needsBuild = false
			e.flushGlobalBuild()
		}
	}

	// 3. 同步 rootRenderObject：局部 rebuild 可能替换/重建 render object 树，
	//    但 RenderObjectElement.Mount 只能挂到最近的 RenderObjectElement 祖先。
	//    当 Navigator transition 插入 Opacity/SlideOffset 等无 RenderObject 的节点时，
	//    新的 root render object 不会被自动 attach，需要在这里同步。
	if e.rootElement != nil {
		newRootRO := e.rootElement.FindRenderObject()
		if newRootRO != e.rootRenderObject {
			if e.rootRenderObject != nil {
				e.rootRenderObject.Detach()
			}
			e.rootRenderObject = newRootRO
			if e.rootRenderObject != nil {
				e.rootRenderObject.Attach(e.owner)
				e.rootRenderObject.MarkNeedsLayout()
			}
		}
	}
}

func (e *Engine) flushGlobalBuild() {
	newRootWidget := e.buildFunc()
	if newRootWidget == nil {
		return
	}

	if e.rootElement == nil {
		// First mount
		e.rootWidget = newRootWidget
		e.rootElement = newRootWidget.CreateElement()
		e.rootElement.Mount(nil, 0)
		if e.rootElement != nil {
			e.rootRenderObject = e.rootElement.FindRenderObject()
			if e.rootRenderObject != nil {
				e.rootRenderObject.Attach(e.owner)
				e.rootRenderObject.MarkNeedsLayout()
			}
		}
	} else {
		// Diff update
		if CanUpdate(e.rootWidget, newRootWidget) {
			e.rootWidget = newRootWidget
			e.rootElement.Update(newRootWidget)
		} else {
			// Type changed, remount
			e.rootElement.Unmount()
			e.rootWidget = newRootWidget
			e.rootElement = newRootWidget.CreateElement()
			e.rootElement.Mount(nil, 0)
			if e.rootElement != nil {
				e.rootRenderObject = e.rootElement.FindRenderObject()
				if e.rootRenderObject != nil {
					e.rootRenderObject.Attach(e.owner)
					e.rootRenderObject.MarkNeedsLayout()
				}
			}
		}
	}
}

// ==================== Sync bounds ====================

func (e *Engine) tickRenderObjects(ro render.RenderObject, dt time.Duration) {
	if ro == nil {
		return
	}
	if t, ok := ro.(render.Tickable); ok {
		if t.Tick(dt) {
			ro.MarkNeedsPaint()
		}
	}
	for _, child := range ro.GetChildren() {
		if child.GetParent() == ro {
			e.tickRenderObjects(child, dt)
		}
	}
}

func (e *Engine) syncBounds(ro render.RenderObject) {
	if ro == nil {
		return
	}
	y := ro.GetYoga()
	if y == nil {
		return
	}
	bounds := render.Bounds{
		X:      y.LayoutLeft(),
		Y:      y.LayoutTop(),
		Width:  y.LayoutWidth(),
		Height: y.LayoutHeight(),
	}
	ro.SetBounds(bounds)
	for _, child := range ro.GetChildren() {
		e.syncBounds(child)
	}

	// 对于 Stack，Yoga 在计算父节点尺寸时会忽略 absolute 子节点，
	// 导致 Stack 的 bounds 可能为 0，进而导致整个 Stack 及其子树被跳过绘制。
	// 这里基于子元素的实际 bounds 扩展 Stack 的 bounds，并调整无约束的 center 子元素位置。
	if stack, ok := ro.(*render.RenderStack); ok {
		e.expandStackBounds(stack)
	}
}

// expandStackBounds 基于子元素 bounds 扩展 Stack 的 bounds，并手动居中无约束的 absolute 子元素。
func (e *Engine) expandStackBounds(stack *render.RenderStack) {
	var maxW, maxH float32
	for _, child := range stack.GetChildren() {
		b := child.GetBounds()
		childW := b.X + b.Width
		childH := b.Y + b.Height
		if childW > maxW {
			maxW = childW
		}
		if childH > maxH {
			maxH = childH
		}
	}
	bounds := stack.GetBounds()
	changed := false
	if maxW > bounds.Width {
		bounds.Width = maxW
		changed = true
	}
	if maxH > bounds.Height {
		bounds.Height = maxH
		changed = true
	}
	if changed {
		stack.SetBounds(bounds)
	}

	// 对于没有方向约束的 absolute 子元素（Center 模式），手动居中
	for _, child := range stack.GetChildren() {
		childYoga := child.GetYoga()
		if childYoga == nil || childYoga.StyleGetPositionType() != yoga.PositionTypeAbsolute {
			continue
		}
		left := childYoga.StyleGetPosition(yoga.EdgeLeft)
		right := childYoga.StyleGetPosition(yoga.EdgeRight)
		top := childYoga.StyleGetPosition(yoga.EdgeTop)
		bottom := childYoga.StyleGetPosition(yoga.EdgeBottom)
		if !left.IsUndefined() || !right.IsUndefined() || !top.IsUndefined() || !bottom.IsUndefined() {
			continue
		}
		cb := child.GetBounds()
		cb.X = (bounds.Width - cb.Width) / 2
		cb.Y = (bounds.Height - cb.Height) / 2
		child.SetBounds(cb)
	}
}


func (e *Engine) drawRenderObject(screen *ebiten.Image, ro render.RenderObject, parentOffset render.Offset) {
	if ro == nil || !ro.IsVisible() {
		return
	}
	bounds := ro.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	t := ro.GetTransform()
	if !t.IsIdentity() {
		e.drawRenderObjectWithTransform(screen, ro, parentOffset, t)
		return
	}

	e._drawRenderObject(screen, ro, parentOffset)
}

func (e *Engine) _drawRenderObject(screen *ebiten.Image, ro render.RenderObject, parentOffset render.Offset) {
	bounds := ro.GetBounds()
	ro.Paint(screen, parentOffset)

	absX := parentOffset.X + bounds.X
	absY := parentOffset.Y + bounds.Y

	scrollX, scrollY := float32(0), float32(0)
	if s, ok := tryGetScroll(ro); ok {
		so := s.GetScrollOffset()
		scrollX, scrollY = so.X, so.Y
	}

	childScreen := screen
	childOffset := render.Offset{X: absX - scrollX, Y: absY - scrollY}
	if c, ok := tryGetClipper(ro); ok && c.ClipChildren() {
		sub := render.SubImage(screen,
			int(absX), int(absY),
			int(bounds.Width), int(bounds.Height))
		if sub != nil {
			childScreen = sub
			// Ebiten SubImage 使用全局坐标系，子节点仍需用全局坐标绘制
			childOffset = render.Offset{X: absX - scrollX, Y: absY - scrollY}
		}
	}

	children := ro.GetChildren()
	if len(children) > 1 {
		e.childrenBuf = append(e.childrenBuf[:0], children...)
		sort.SliceStable(e.childrenBuf, func(i, j int) bool {
			return e.getSubtreeMaxZIndex(e.childrenBuf[i]) < e.getSubtreeMaxZIndex(e.childrenBuf[j])
		})
			// Copy to a new slice so recursive calls can't overwrite our iteration array.
		children = make([]render.RenderObject, len(e.childrenBuf))
		copy(children, e.childrenBuf)
	}
	for _, child := range children {
		e.drawRenderObject(childScreen, child, childOffset)
	}
}

// getSubtreeMaxZIndex 递归计算 RenderObject 子树中的最大 zIndex。
func (e *Engine) getSubtreeMaxZIndex(ro render.RenderObject) int {
	maxZ := ro.GetZIndex()
	for _, child := range ro.GetChildren() {
		if z := e.getSubtreeMaxZIndex(child); z > maxZ {
			maxZ = z
		}
	}
	return maxZ
}

func (e *Engine) drawRenderObjectWithTransform(screen *ebiten.Image, ro render.RenderObject, parentOffset render.Offset, t render.Transform) {
	bounds := ro.GetBounds()
	iw, ih := int(bounds.Width), int(bounds.Height)
	if iw <= 0 || ih <= 0 {
		return
	}

	temp := ebiten.NewImage(iw, ih)
	defer temp.Dispose()

	// 将子树绘制到临时 surface，节点的 (0,0) 对齐到 temp 的 (0,0)
	e._drawRenderObject(temp, ro, render.Offset{X: -bounds.X, Y: -bounds.Y})

	// 3D 变换使用 DrawTriangles 实现透视投影
	if t.Has3D() {
		e.drawWith3DTransform(screen, temp, bounds, parentOffset, t)
		return
	}

	geoM := render.BuildTransformGeoM(render.Bounds{X: 0, Y: 0, Width: bounds.Width, Height: bounds.Height}, t)
	geoM.Translate(float64(parentOffset.X+bounds.X), float64(parentOffset.Y+bounds.Y))

	op := &ebiten.DrawImageOptions{GeoM: geoM}
	if t.Alpha < 1 {
		op.ColorScale.Scale(float32(t.Alpha), float32(t.Alpha), float32(t.Alpha), float32(t.Alpha))
	}
	screen.DrawImage(temp, op)
}

func (e *Engine) drawWith3DTransform(screen *ebiten.Image, temp *ebiten.Image, bounds render.Bounds, parentOffset render.Offset, t render.Transform) {
	b := render.Bounds{X: parentOffset.X + bounds.X, Y: parentOffset.Y + bounds.Y, Width: bounds.Width, Height: bounds.Height}
	verts := render.Build3DVertices(b, t)

	w, h := float32(temp.Bounds().Dx()), float32(temp.Bounds().Dy())
	vertices := []ebiten.Vertex{
		{DstX: float32(verts[0][0]), DstY: float32(verts[0][1]), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: t.Alpha},
		{DstX: float32(verts[1][0]), DstY: float32(verts[1][1]), SrcX: w, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: t.Alpha},
		{DstX: float32(verts[2][0]), DstY: float32(verts[2][1]), SrcX: w, SrcY: h, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: t.Alpha},
		{DstX: float32(verts[3][0]), DstY: float32(verts[3][1]), SrcX: 0, SrcY: h, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: t.Alpha},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}

	op := &ebiten.DrawTrianglesOptions{Filter: ebiten.FilterLinear}
	screen.DrawTriangles(vertices, indices, temp, op)
}

// ==================== Run ====================

// Run starts the ebiten game loop.
func (e *Engine) Run() {
	e.Mount()
	ebiten.SetWindowSize(e.screenWidth, e.screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(e); err != nil {
		panic(err)
	}
}

// ==================== EdgeInsets ====================

type EdgeInsets struct {
	Top, Right, Bottom, Left float32
}

func (e EdgeInsets) All() float32 {
	return e.Top + e.Right + e.Bottom + e.Left
}

func EdgeInsetsAll(v float32) EdgeInsets {
	return EdgeInsets{Top: v, Right: v, Bottom: v, Left: v}
}

func EdgeInsetsSymmetric(horizontal, vertical float32) EdgeInsets {
	return EdgeInsets{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

func EdgeInsetsOnly(top, right, bottom, left float32) EdgeInsets {
	return EdgeInsets{Top: top, Right: right, Bottom: bottom, Left: left}
}

// Yoga type aliases
type (
	FlexDirection = yoga.FlexDirection
	Justify       = yoga.Justify
	Align         = yoga.Align
	Wrap          = yoga.Wrap
	PositionType  = yoga.PositionType
	Display       = yoga.Display
	Overflow      = yoga.Overflow
	Edge          = yoga.Edge
	Gutter        = yoga.Gutter
	Direction     = yoga.Direction
)

const (
	FlexDirectionColumn    = yoga.FlexDirectionColumn
	FlexDirectionRow       = yoga.FlexDirectionRow
	JustifyFlexStart       = yoga.JustifyFlexStart
	JustifyCenter          = yoga.JustifyCenter
	JustifyFlexEnd         = yoga.JustifyFlexEnd
	JustifySpaceBetween    = yoga.JustifySpaceBetween
	JustifySpaceAround     = yoga.JustifySpaceAround
	JustifySpaceEvenly     = yoga.JustifySpaceEvenly
	AlignAuto              = yoga.AlignAuto
	AlignFlexStart         = yoga.AlignFlexStart
	AlignCenter            = yoga.AlignCenter
	AlignFlexEnd           = yoga.AlignFlexEnd
	AlignStretch           = yoga.AlignStretch
	AlignBaseline          = yoga.AlignBaseline
	WrapNoWrap             = yoga.WrapNoWrap
	WrapWrap               = yoga.WrapWrap
	PositionTypeRelative   = yoga.PositionTypeRelative
	PositionTypeAbsolute   = yoga.PositionTypeAbsolute
	DisplayFlex            = yoga.DisplayFlex
	DisplayNone            = yoga.DisplayNone
	OverflowVisible        = yoga.OverflowVisible
	OverflowHidden         = yoga.OverflowHidden
	OverflowScroll         = yoga.OverflowScroll
	EdgeLeft               = yoga.EdgeLeft
	EdgeTop                = yoga.EdgeTop
	EdgeRight              = yoga.EdgeRight
	EdgeBottom             = yoga.EdgeBottom
	EdgeStart              = yoga.EdgeStart
	EdgeEnd                = yoga.EdgeEnd
	EdgeHorizontal         = yoga.EdgeHorizontal
	EdgeVertical           = yoga.EdgeVertical
	EdgeAll                = yoga.EdgeAll
	GutterColumn           = yoga.GutterColumn
	GutterRow              = yoga.GutterRow
	GutterAll              = yoga.GutterAll
	DirectionLTR           = yoga.DirectionLTR
	DirectionRTL           = yoga.DirectionRTL
)
