package core

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sjm1327605995/tenon/yoga"
)

// Engine 是 Tenon 核心，负责 Widget → Element 的构建、布局计算、绘制和事件。
type Engine struct {
	rootWidget   Widget
	rootElement  Element

	// 帧调度
	dirtyElements []Element
	buildQueue    []Widget

	// 窗口
	screenWidth   int
	screenHeight  int
	pendingResize bool

	// 输入状态
	lastMouseX      float32
	lastMouseY      float32
	mouseDownTarget Element
	mouseDownX      float32
	mouseDownY      float32
	mouseDownButton ebiten.MouseButton
	hoverTarget     Element
	focusTarget     Element

	// 时间
	lastFrameTime time.Time

	// 动画
	animations []Animation
}

// NewEngine 创建引擎。
func NewEngine(rootWidget Widget, width, height int) *Engine {
	return &Engine{
		rootWidget:    rootWidget,
		screenWidth:   width,
		screenHeight:  height,
		lastFrameTime: time.Now(),
	}
}

// Layout implements ebiten.Game interface.
func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != e.screenWidth || outsideHeight != e.screenHeight {
		e.screenWidth = outsideWidth
		e.screenHeight = outsideHeight
		e.pendingResize = true
	}
	return outsideWidth, outsideHeight
}

// Mount 执行首次挂载，调用 Widget.Build() 并构建 Element 树。
func (e *Engine) Mount() {
	if e.rootWidget == nil {
		return
	}
	e.rootWidget.OnMount(e)
	e.rootElement = e.rootWidget.Build()
	if e.rootElement != nil {
		e.onElementMounted(e.rootElement)
		e.calculateLayout()
	}
}

// onElementMounted 递归给新挂载的 Element 注入 Engine。
func (e *Engine) onElementMounted(el Element) {
	if el == nil {
		return
	}
	el.SetEngine(e)
	applyStyles(el)
	el.OnMount(e)
	for _, child := range el.GetChildren() {
		e.onElementMounted(child)
	}
}

// ==================== 帧循环 ====================

func (e *Engine) Update() error {
	// 0. 窗口大小变化：重新布局并广播事件
	if e.pendingResize {
		e.pendingResize = false
		e.calculateLayout()
		if e.rootElement != nil {
			e.broadcastEvent(e.rootElement, &Event{
				Type:   EventResize,
				Width:  float32(e.screenWidth),
				Height: float32(e.screenHeight),
			})
		}
	}

	now := time.Now()
	deltaTime := float32(now.Sub(e.lastFrameTime).Seconds())
	e.lastFrameTime = now

	// 1. 处理请求重建的 Widget（结构变化）
	e.flushBuildQueue()

	// 2. 事件处理（可能标记 Element 为 dirty）
	e.handleEvents()

	// 3. 批量刷新脏 Element（测量 → 布局 → 清除标记）
	// 放在事件之后，确保同一帧内响应交互带来的布局变化
	e.flushDirtyElements()

	// 4. 更新动画
	e.updateAnimations(deltaTime)

	// 5. 每帧更新 Element（动画等）
	if e.rootElement != nil {
		e.updateElements(e.rootElement)
	}

	return nil
}

// AddAnimation 注册一个动画到引擎。
func (e *Engine) AddAnimation(a Animation) {
	e.animations = append(e.animations, a)
}

// updateAnimations 更新所有动画，移除已完成的。
func (e *Engine) updateAnimations(deltaTime float32) {
	const maxDelta = 1.0 / 30.0
	if deltaTime > maxDelta {
		deltaTime = maxDelta
	}
	alive := e.animations[:0]
	for _, a := range e.animations {
		if a.Update(deltaTime) {
			alive = append(alive, a)
		}
	}
	e.animations = alive
}

func (e *Engine) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 245, G: 245, B: 245, A: 255})
	if e.rootElement != nil {
		e.drawElement(screen, e.rootElement, 0, 0)
	}
}

// ==================== 结构重建 ====================

func (e *Engine) scheduleBuild(w Widget) {
	e.buildQueue = append(e.buildQueue, w)
}

func (e *Engine) flushBuildQueue() {
	if len(e.buildQueue) == 0 {
		return
	}
	queue := e.buildQueue
	e.buildQueue = e.buildQueue[:0]

	for _, w := range queue {
		newRoot := w.Build()
		if newRoot == nil {
			continue
		}

		var oldRoot Element
		if w == e.rootWidget {
			oldRoot = e.rootElement
		} else if bw, ok := w.(interface{ GetRootElement() Element }); ok {
			oldRoot = bw.GetRootElement()
		}

		if oldRoot == nil {
			if w == e.rootWidget {
				e.rootElement = newRoot
			}
			if bw, ok := w.(interface{ SetRootElement(Element) }); ok {
				bw.SetRootElement(newRoot)
			}
			e.onElementMounted(newRoot)
		} else {
			finalRoot := e.patchElement(oldRoot, newRoot)
			if w == e.rootWidget {
				e.rootElement = finalRoot
			}
			if bw, ok := w.(interface{ SetRootElement(Element) }); ok {
				bw.SetRootElement(finalRoot)
			}
		}
	}

	// 重建后布局可能变化
	if e.hasLayoutDirty() {
		e.calculateLayout()
	}
}

// patchElement 对新旧 Element 做同级浅对比。
// 同类型复用旧节点并递归同步子节点，不同类型替换根节点。
func (e *Engine) patchElement(oldEl, newEl Element) Element {
	if oldEl.ElementType() == newEl.ElementType() {
		// 同类型：复用旧节点，同步 Yoga 样式
		if oldYoga, newYoga := oldEl.GetYoga(), newEl.GetYoga(); oldYoga != nil && newYoga != nil {
			oldYoga.CopyStyleFrom(newYoga)
		}
		e.patchChildren(oldEl, newEl)
		oldEl.Mark(FlagNeedLayout | FlagNeedDraw)
		return oldEl
	}
	// 类型不同：替换根节点
	if oldEl.GetParent() != nil {
		parent := oldEl.GetParent()
		parent.RemoveChild(oldEl)
		parent.AppendChild(newEl)
	} else {
		e.rootElement = newEl
		e.onElementMounted(newEl)
	}
	return newEl
}

// patchChildren 对两个同类型 Element 的子节点做轻量同级对比。
// 前 minLen 个递归 patch，多余旧节点移除，多余新节点挂载。
func (e *Engine) patchChildren(oldParent, newParent Element) {
	oldChildren := oldParent.GetChildren()
	newChildren := newParent.GetChildren()

	minLen := len(oldChildren)
	if len(newChildren) < minLen {
		minLen = len(newChildren)
	}

	// 1. 同步前 minLen 个子节点（递归复用）
	for i := 0; i < minLen; i++ {
		e.patchElement(oldChildren[i], newChildren[i])
	}

	// 2. 移除多余的旧节点
	for i := len(newChildren); i < len(oldChildren); i++ {
		oldParent.RemoveChild(oldChildren[i])
	}

	// 3. 挂载多余的新节点（先从 newParent 解除 Yoga 关系，再挂到 oldParent）
	if len(newChildren) > len(oldChildren) {
		toAdd := newChildren[len(oldChildren):]
		// 倒序从 newParent 移除，避免索引漂移
		for i := len(toAdd) - 1; i >= 0; i-- {
			newParent.RemoveChild(toAdd[i])
		}
		for _, child := range toAdd {
			oldParent.AppendChild(child)
		}
	}
}

// ==================== 脏标记刷新 ====================

func (e *Engine) markDirty(el Element) {
	if el == nil {
		return
	}
	// 避免重复加入
	for _, d := range e.dirtyElements {
		if d == el {
			return
		}
	}
	e.dirtyElements = append(e.dirtyElements, el)
}

func (e *Engine) flushDirtyElements() {
	if len(e.dirtyElements) == 0 {
		return
	}

	// 1. 先测量（文字排版等）
	for _, el := range e.dirtyElements {
		if el.GetFlags()&FlagNeedMeasure != 0 {
			if m, ok := el.(interface{ FlushMeasure() }); ok {
				m.FlushMeasure()
			}
		}
	}

	// 2. 再布局（只要有一个节点标记了 NeedLayout，就从根重算）
	if e.hasLayoutDirty() {
		e.calculateLayout()
	}

	// 3. 清除所有标记（只清脏标记，保留持久状态）
	for _, el := range e.dirtyElements {
		el.ClearDirty()
	}
	e.dirtyElements = e.dirtyElements[:0]
}

func (e *Engine) hasLayoutDirty() bool {
	for _, el := range e.dirtyElements {
		if el.GetFlags()&FlagNeedLayout != 0 {
			return true
		}
	}
	return false
}

// ==================== 布局 ====================

func (e *Engine) calculateLayout() {
	if e.rootElement == nil || e.rootElement.GetYoga() == nil {
		return
	}
	rootYoga := e.rootElement.GetYoga()
	rootYoga.StyleSetWidth(float32(e.screenWidth))
	rootYoga.StyleSetHeight(float32(e.screenHeight))
	rootYoga.CalculateLayout(float32(e.screenWidth), float32(e.screenHeight), yoga.DirectionLTR)

	if e.rootElement != nil {
		e.updateBounds(e.rootElement, 0, 0)
	}
}

func (e *Engine) updateBounds(el Element, parentX, parentY float32) {
	if el == nil || el.GetYoga() == nil {
		return
	}
	y := el.GetYoga()
	b := LayoutBounds{
		X:      parentX + y.LayoutLeft(),
		Y:      parentY + y.LayoutTop(),
		Width:  y.LayoutWidth(),
		Height: y.LayoutHeight(),
	}
	el.SetBounds(b)
	for _, child := range el.GetChildren() {
		e.updateBounds(child, el.GetBounds().X, el.GetBounds().Y)
	}
}

// ==================== 绘制 ====================

func (e *Engine) drawElement(screen *ebiten.Image, el Element, offsetX, offsetY float32) {
	if el == nil || !el.IsVisible() {
		return
	}
	bounds := el.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 若标记了 ClipChildren，子节点绘制在裁剪后的子图上
	var childScreen *ebiten.Image = screen
	childOffsetX := offsetX
	childOffsetY := offsetY
	if el.HasFlag(FlagClipChildren) {
		sub := screen.SubImage(image.Rect(
			int(bounds.X), int(bounds.Y),
			int(bounds.X+bounds.Width), int(bounds.Y+bounds.Height),
		))
		if subImg, ok := sub.(*ebiten.Image); ok {
			childScreen = subImg
			// SubImage uses the same coordinate system as the original image;
			// it only acts as a clip mask. Children should use screen coordinates.
			childOffsetX = 0
			childOffsetY = 0
		}
	}

	// Use screen coordinates for drawing. SubImage will clip automatically.
	relBounds := LayoutBounds{
		X:      bounds.X - childOffsetX,
		Y:      bounds.Y - childOffsetY,
		Width:  bounds.Width,
		Height: bounds.Height,
	}
	el.SetBounds(relBounds)
	el.Draw(childScreen)
	el.SetBounds(bounds)

	for _, child := range el.GetChildren() {
		e.drawElement(childScreen, child, childOffsetX, childOffsetY)
	}
}

// ==================== Element 每帧更新 ====================

func (e *Engine) updateElements(el Element) {
	if el == nil {
		return
	}
	_ = el.Update()
	for _, child := range el.GetChildren() {
		e.updateElements(child)
	}
}

// ==================== 事件系统 ====================

type eventState struct {
	mouseX, mouseY     float32
	mouseDownTarget    Element
	mouseDownX         float32
	mouseDownY         float32
	hoverTarget        Element
	focusTarget        Element
	mouseDownButton    ebiten.MouseButton
}

func (e *Engine) handleEvents() {
	mx, my := ebiten.CursorPosition()
	fx, fy := float32(mx), float32(my)

	// 1. MouseMove / Hover
	if fx != e.lastMouseX || fy != e.lastMouseY {
		e.handleMouseMove(fx, fy)
	}

	// 2. Mouse wheel
	dx, dy := ebiten.Wheel()
	if math.Abs(dx) >= 0.5 || math.Abs(dy) >= 0.5 {
		e.handleScroll(fx, fy, float32(dx), float32(dy))
	}

	// 3. Mouse buttons
	for _, btn := range []ebiten.MouseButton{ebiten.MouseButtonLeft, ebiten.MouseButtonRight, ebiten.MouseButtonMiddle} {
		if inpututil.IsMouseButtonJustPressed(btn) {
			e.handleMouseDown(fx, fy, btn)
		}
		if inpututil.IsMouseButtonJustReleased(btn) {
			e.handleMouseUp(fx, fy, btn)
		}
	}

	// 4. Keyboard
	e.handleKeyboard()

	e.lastMouseX = fx
	e.lastMouseY = fy
}

func (e *Engine) handleMouseMove(x, y float32) {
	target := e.hitTest(e.rootElement, x, y)
	if target != e.hoverTarget {
		// MouseLeave old target
		if e.hoverTarget != nil {
			e.dispatchEvent(e.hoverTarget, &Event{Type: EventMouseMove, X: x, Y: y, Target: e.hoverTarget})
		}
		e.hoverTarget = target
	}
	if target != nil {
		b := target.GetBounds()
		e.dispatchEvent(target, &Event{
			Type:   EventMouseMove,
			X:      x,
			Y:      y,
			LocalX: x - b.X,
			LocalY: y - b.Y,
			Target: target,
		})
	}
}

func (e *Engine) handleMouseDown(x, y float32, btn ebiten.MouseButton) {
	target := e.hitTest(e.rootElement, x, y)
	if target != nil {
		LogDebug("[Engine] handleMouseDown", "target", target.ElementType(), "bounds", target.GetBounds(), "x", x, "y", y)
		// Debug: log target hierarchy and children
		parentType := "nil"
		if target.GetParent() != nil {
			parentType = target.GetParent().ElementType()
		}
		LogDebug("[Engine] handleMouseDown debug", "target", target.ElementType(), "parent", parentType, "childrenCount", len(target.GetChildren()))
		for i, c := range target.GetChildren() {
			LogDebug("[Engine] handleMouseDown child", "index", i, "type", c.ElementType(), "bounds", c.GetBounds())
		}
		b := target.GetBounds()
		e.mouseDownTarget = target
		e.mouseDownX = x
		e.mouseDownY = y
		e.mouseDownButton = btn
		// Set focus on click
		if e.focusTarget != target {
			e.setFocus(target)
		}
		e.dispatchEvent(target, &Event{
			Type:   EventMouseDown,
			X:      x,
			Y:      y,
			LocalX: x - b.X,
			LocalY: y - b.Y,
			Button: btn,
			Target: target,
		})
	} else {
		// Click on empty area clears focus
		e.setFocus(nil)
	}
}

func (e *Engine) handleMouseUp(x, y float32, btn ebiten.MouseButton) {
	if e.mouseDownTarget != nil {
		b := e.mouseDownTarget.GetBounds()
		e.dispatchEvent(e.mouseDownTarget, &Event{
			Type:   EventMouseUp,
			X:      x,
			Y:      y,
			LocalX: x - b.X,
			LocalY: y - b.Y,
			Button: btn,
			Target: e.mouseDownTarget,
		})

		// Click = mouseDown + mouseUp on same target
		target := e.hitTest(e.rootElement, x, y)
		if target == e.mouseDownTarget {
			e.dispatchEvent(target, &Event{
				Type:   EventClick,
				X:      x,
				Y:      y,
				LocalX: x - b.X,
				LocalY: y - b.Y,
				Button: btn,
				Target: target,
			})
		}
	}
	e.mouseDownTarget = nil
}

func (e *Engine) handleScroll(x, y float32, dx, dy float32) {
	target := e.hitTest(e.rootElement, x, y)
	LogDebug("[Engine] handleScroll", "target", target.ElementType())
	if target != nil {
		b := target.GetBounds()
		e.dispatchEvent(target, &Event{
			Type:   EventScroll,
			X:      x,
			Y:      y,
			LocalX: x - b.X,
			LocalY: y - b.Y,
			DeltaX: dx,
			DeltaY: dy,
			Target: target,
		})
	}
}

func (e *Engine) handleKeyboard() {
	// Tab to cycle focus
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		e.focusNext(ebiten.IsKeyPressed(ebiten.KeyShift))
		return
	}

	if e.focusTarget == nil {
		return
	}

	// Space/Enter triggers click on focused element
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		b := e.focusTarget.GetBounds()
		e.dispatchEvent(e.focusTarget, &Event{
			Type:   EventClick,
			X:      b.X + b.Width/2,
			Y:      b.Y + b.Height/2,
			LocalX: b.Width / 2,
			LocalY: b.Height / 2,
			Target: e.focusTarget,
		})
		return
	}

	// KeyDown / KeyUp
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		if inpututil.IsKeyJustPressed(key) {
			e.dispatchEvent(e.focusTarget, &Event{Type: EventKeyDown, Key: key, Target: e.focusTarget})
		}
		if inpututil.IsKeyJustReleased(key) {
			e.dispatchEvent(e.focusTarget, &Event{Type: EventKeyUp, Key: key, Target: e.focusTarget})
		}
	}
}

// dispatchEvent sends event to target and bubbles up until consumed.
func (e *Engine) dispatchEvent(target Element, event *Event) {
	el := target
	for el != nil {
		if el.HandleEvent(event) {
			return // consumed
		}
		el = el.GetParent()
	}
}

// broadcastEvent 递归广播事件给元素及其所有后代。
func (e *Engine) broadcastEvent(el Element, event *Event) {
	if el == nil {
		return
	}
	el.HandleEvent(event)
	for _, child := range el.GetChildren() {
		e.broadcastEvent(child, event)
	}
}

// GetFocusTarget returns the currently focused element.
func (e *Engine) GetFocusTarget() Element {
	return e.focusTarget
}

// setFocus changes focused element.
func (e *Engine) setFocus(target Element) {
	if e.focusTarget == target {
		return
	}
	if e.focusTarget != nil {
		e.focusTarget.HandleEvent(&Event{Type: EventFocusOut, Target: e.focusTarget})
	}
	e.focusTarget = target
	if target != nil {
		target.HandleEvent(&Event{Type: EventFocusIn, Target: target})
	}
}

// focusNext cycles focus to next/previous focusable element.
func (e *Engine) focusNext(reverse bool) {
	focusables := e.collectFocusables(e.rootElement)
	if len(focusables) == 0 {
		return
	}
	idx := -1
	for i, el := range focusables {
		if el == e.focusTarget {
			idx = i
			break
		}
	}
	if reverse {
		idx--
		if idx < 0 {
			idx = len(focusables) - 1
		}
	} else {
		idx++
		if idx >= len(focusables) {
			idx = 0
		}
	}
	e.setFocus(focusables[idx])
}

func (e *Engine) collectFocusables(el Element) []Element {
	if el == nil {
		return nil
	}
	var result []Element
	if el.HasFlag(FlagFocusable) {
		result = append(result, el)
	}
	for _, child := range el.GetChildren() {
		result = append(result, e.collectFocusables(child)...)
	}
	return result
}

func (e *Engine) hitTest(el Element, x, y float32) Element {
	if el == nil || !el.IsVisible() {
		return nil
	}
	// 后序遍历：子节点优先
	children := el.GetChildren()
	for i := len(children) - 1; i >= 0; i-- {
		if h := e.hitTest(children[i], x, y); h != nil {
			return h
		}
	}
	b := el.GetBounds()
	if x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height {
		return el
	}
	return nil
}
