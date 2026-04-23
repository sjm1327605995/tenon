package core

import (
	"image"
	"image/color"
	"reflect"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sjm1327605995/tenon/yoga"
)

// renderNode 是 Engine 内部的组件树节点，只包含 Host。
// Widget 在挂载时被展开，不会出现在 renderNode 树中。
type renderNode struct {
	host         Host
	sourceWidget Widget   // 若此子树由 Widget 展开，记录来源
	parent       *renderNode
	children     []*renderNode
}

// Engine 是 Tenon 的核心引擎，负责挂载、更新、渲染和事件。
type Engine struct {
	rootHost     Host
	rootNode     *renderNode
	rootWidget   Widget // 根 Widget（当 root 是 Widget 而非 Host 时）

	widgetMounts map[Widget]*renderNode // Widget -> 展开后的根 renderNode
	dirtyWidgets map[Widget]struct{}    // 待更新队列

	overlays []*renderNode // Portal 浮层

	screenWidth  int
	screenHeight int

	focusHost Host // 当前焦点组件

	mouseDownTarget  Host              // 鼠标按下的目标组件（用于拖拽）
	mouseDownButton  ebiten.MouseButton // 按下的鼠标按钮
	lastMouseX       float32           // 上一帧鼠标 X
	lastMouseY       float32           // 上一帧鼠标 Y
	animations     []Animation // 活跃动画列表
	lastFrameTime  time.Time   // 上一帧时间（用于计算 deltaTime）
}

// NewEngine 创建一个引擎，rootHost 必须是宿主组件。
func NewEngine(rootHost Host, width, height int) *Engine {
	return &Engine{
		rootHost:     rootHost,
		widgetMounts: make(map[Widget]*renderNode),
		dirtyWidgets: make(map[Widget]struct{}),
		screenWidth:  width,
		screenHeight: height,
		lastFrameTime: time.Now(),
	}
}

// SetScreenSize 在窗口大小变化时调用。
func (e *Engine) SetScreenSize(w, h int) {
	e.screenWidth = w
	e.screenHeight = h
}

// Mount 执行初始挂载，构建完整的 renderNode 树。
func (e *Engine) Mount() {
	e.rootNode = e.mount(nil, e.rootHost)
	e.injectEngine(e.rootNode)

	// 如果根是 Widget，注册到 widgetMounts，使其 Invalidate 能触发更新
	if e.rootWidget != nil && e.rootNode != nil {
		e.widgetMounts[e.rootWidget] = e.rootNode
		e.rootWidget.SetHostRef(e.rootNode.host)
		e.rootWidget.setEngine(e)
	}

	e.calculateLayout()
	setActiveEngine(e)
}

// injectEngine 给所有 Host 节点注入 Engine 引用。
func (e *Engine) injectEngine(node *renderNode) {
	if node == nil {
		return
	}
	if node.host != nil {
		node.host.SetEngine(e)
	}
	for _, child := range node.children {
		e.injectEngine(child)
	}
}

// ==================== 挂载 / 卸载 ====================

func (e *Engine) mount(parent *renderNode, c Component) *renderNode {
	if c == nil {
		return nil
	}

	// Widget：展开，不创建自己的节点
	if w, ok := c.(Widget); ok {
		w.setEngine(e)
		w.resetHooks()
		rendered := w.Render()
		node := e.mount(parent, rendered)
		if node != nil {
			e.widgetMounts[w] = node
			node.sourceWidget = w
			w.SetHostRef(node.host)
			w.ComponentDidMount()
		}
		return node
	}

	// Host：创建 renderNode
	if h, ok := c.(Host); ok {
		node := &renderNode{host: h, parent: parent}
		h.SetEngine(e)

		// 布局钩子
		bounds := h.GetLayoutBounds()
		if h.LayoutChildren(bounds.Width, bounds.Height) {
			// Host 自行处理子组件布局，跳过 Yoga 插入
		} else {
			for _, child := range h.GetChildren() {
				childNode := e.mount(node, child)
				if childNode != nil && childNode.host != nil {
					node.children = append(node.children, childNode)
					e.insertYogaChild(h, childNode.host)
				}
			}
		}
		return node
	}

	return nil
}

func (e *Engine) unmount(node *renderNode) {
	if node == nil {
		return
	}

	// 递归卸载子节点
	for _, child := range node.children {
		e.unmount(child)
	}
	node.children = nil

	// Widget 来源清理
	if node.sourceWidget != nil {
		node.sourceWidget.SetHostRef(nil)
		node.sourceWidget.ComponentWillUnmount()
		delete(e.widgetMounts, node.sourceWidget)
		node.sourceWidget = nil
	}

	// 从 Yoga 树移除
	if node.parent != nil && node.parent.host != nil && node.host != nil {
		e.removeYogaChild(node.parent.host, node.host)
	}
}

// ==================== 更新调度 ====================

func (e *Engine) scheduleUpdate(w Widget) {
	e.dirtyWidgets[w] = struct{}{}
}

// InvalidateAll 标记所有已挂载的 Widget 为 dirty，触发全局重渲染。
// 通常用于主题切换等需要更新整个 UI 树的场景。
func (e *Engine) InvalidateAll() {
	for w := range e.widgetMounts {
		e.dirtyWidgets[w] = struct{}{}
	}
}

func (e *Engine) flushUpdates() {
	for w := range e.dirtyWidgets {
		e.updateWidget(w)
	}
	e.dirtyWidgets = make(map[Widget]struct{})
}

func (e *Engine) updateWidget(w Widget) {
	oldNode, ok := e.widgetMounts[w]
	if !ok || oldNode == nil {
		return
	}

	if !w.ShouldComponentUpdate(nil) {
		return
	}

	// 重新 Render
	w.resetHooks()
	rendered := w.Render()

	// 自动复用同类型 Host：同步属性 → 清空旧 children → 复制新 children
	if oldNode.host != nil {
		if newHost, ok := rendered.(Host); ok {
			if reflect.TypeOf(oldNode.host) == reflect.TypeOf(newHost) {
				syncHost(oldNode.host, newHost)
				newChildren := make([]Component, len(newHost.GetChildren()))
				copy(newChildren, newHost.GetChildren())
				oldNode.host.ClearChildren()
				for _, child := range newChildren {
					oldNode.host.AddChild(child)
				}
				rendered = oldNode.host
			}
		}
	}

	// 优先尝试复用旧的根 Host（现在 rendered 可能是 oldNode.host 自身）
	if oldNode.host != nil {
		if newHost, ok := rendered.(Host); ok {
			if e.replaceHostTree(oldNode, newHost) {
				w.SetHostRef(oldNode.host)
				w.ComponentDidUpdate(nil, nil)
				return
			}
		}
	}

	// 无法复用，全量重建
	parent := oldNode.parent
	e.unmount(oldNode)

	if parent != nil {
		for i, child := range parent.children {
			if child == oldNode {
				parent.children = append(parent.children[:i], parent.children[i+1:]...)
				break
			}
		}
		newNode := e.mount(parent, rendered)
		if newNode != nil {
			parent.children = append(parent.children, newNode)
			if newNode.host != nil && parent.host != nil {
				e.insertYogaChild(parent.host, newNode.host)
			}
			e.widgetMounts[w] = newNode
			newNode.sourceWidget = w
			w.SetHostRef(newNode.host)
		}
	} else {
		// 根节点（parent == nil）
		e.rootNode = nil
		if h, ok := rendered.(Host); ok {
			e.rootHost = h
			newNode := e.mount(nil, rendered)
			e.rootNode = newNode
			if newNode != nil {
				e.widgetMounts[w] = newNode
				newNode.sourceWidget = w
				w.SetHostRef(newNode.host)
			}
		}
	}

	w.ComponentDidUpdate(nil, nil)
}

// syncElement 将 src Element 的视觉属性复制到 dst，保留 dst 的 Yoga 节点。
func syncElement(dst, src *Element) {
	if dst == nil || src == nil {
		return
	}
	dst.Visible = src.Visible
	dst.BackgroundColor = src.BackgroundColor
	dst.BorderColor = src.BorderColor
	dst.ShadowColor = src.ShadowColor
	dst.ShadowBlur = src.ShadowBlur
	dst.ShadowOffsetX = src.ShadowOffsetX
	dst.ShadowOffsetY = src.ShadowOffsetY
	dst.PointerEvents = src.PointerEvents
	dst.BorderRadius = src.BorderRadius
	dst.Key = src.Key
}

// syncHost 将 src Host 的属性同步到 dst Host，保留 dst 的 Yoga 节点。
// 如果 dst 实现了 SyncFrom(Host) 方法，还会调用它进行组件特有的属性同步。
func syncHost(dst, src Host) {
	if dst == nil || src == nil {
		return
	}

	// 同步 Element 视觉属性
	syncElement(dst.GetElement(), src.GetElement())

	// 同步 Yoga 样式（保留旧的 Yoga 节点，避免重新创建）
	if dstEl := dst.GetElement(); dstEl != nil && dstEl.Yoga != nil {
		if srcEl := src.GetElement(); srcEl != nil && srcEl.Yoga != nil {
			dstEl.Yoga.CopyStyleFrom(srcEl.Yoga)
		}
	}

	// 组件特有属性同步（可选）
	if s, ok := dst.(interface{ SyncFrom(Host) }); ok {
		s.SyncFrom(src)
	}
}

// replaceHostTree 尝试用新 Host 的内容替换旧 Host 的内容，保留旧 Host 的 Yoga 节点。
// 只有当新旧 Host 是同一个实例（指针相同）时才复用，避免 Yoga 样式不同步导致布局错乱。
// 用户若希望复用，应在 Widget 中把 Host 缓存为字段，在 Render 中返回同一个实例。
func (e *Engine) replaceHostTree(node *renderNode, newHost Host) bool {
	oldHost := node.host
	if oldHost == nil || newHost == nil {
		return false
	}

	// 只有同一个实例才复用，确保 Yoga 样式完全一致
	if oldHost != newHost {
		return false
	}

	// 卸载旧 children
	for _, child := range node.children {
		e.unmount(child)
	}
	node.children = nil

	// 挂载新 children 到旧 Host 下
	for _, child := range newHost.GetChildren() {
		childNode := e.mount(node, child)
		if childNode != nil && childNode.host != nil {
			node.children = append(node.children, childNode)
			e.insertYogaChild(oldHost, childNode.host)
		}
	}

	// 标记 Yoga 需要重新布局（尺寸可能因内容变化）
	if oldHost.GetElement() != nil && oldHost.GetElement().Yoga != nil {
		oldHost.GetElement().Yoga.MarkDirty()
	}

	return true
}

// ==================== 渲染循环 ====================

func (e *Engine) Update() error {
	// 计算 deltaTime
	now := time.Now()
	deltaTime := float32(now.Sub(e.lastFrameTime).Seconds())
	e.lastFrameTime = now

	// 0. 动画更新
	e.updateAnimations(deltaTime)

	// 1. Widget 更新
	e.flushUpdates()

	// 2. 布局计算
	e.calculateLayout()

	// 3. 事件处理
	e.handleEvents()

	// 4. Host 每帧更新
	if e.rootNode != nil {
		e.updateHosts(e.rootNode)
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 245, G: 245, B: 245, A: 255})

	if e.rootNode != nil {
		e.drawNode(screen, e.rootNode)
	}

	for _, overlay := range e.overlays {
		e.drawNode(screen, overlay)
	}
}

// ==================== 布局 ====================

func (e *Engine) calculateLayout() {
	if e.rootHost == nil || e.rootHost.GetElement() == nil || e.rootHost.GetElement().Yoga == nil {
		return
	}

	rootYoga := e.rootHost.GetElement().Yoga
	rootYoga.StyleSetWidth(float32(e.screenWidth))
	rootYoga.StyleSetHeight(float32(e.screenHeight))
	rootYoga.CalculateLayout(float32(e.screenWidth), float32(e.screenHeight), yoga.DirectionLTR)

	if e.rootNode != nil {
		e.updateBounds(e.rootNode, 0, 0)
	}
}

func (e *Engine) updateBounds(node *renderNode, parentX, parentY float32) {
	if node == nil || node.host == nil {
		return
	}

	node.host.UpdateLayoutBounds(parentX, parentY)
	bounds := node.host.GetLayoutBounds()
	scrollX, scrollY := node.host.GetScrollOffset()

	for _, child := range node.children {
		e.updateBounds(child, bounds.X+scrollX, bounds.Y+scrollY)
	}
}

// ==================== 绘制 ====================

func (e *Engine) drawNode(screen *ebiten.Image, node *renderNode) {
	if node == nil || node.host == nil {
		return
	}

	el := node.host.GetElement()
	if el == nil || !el.Visible {
		return
	}

	bounds := node.host.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// Clip
	clip := node.host.ShouldClipChildren()
	if clip {
		rect := image.Rect(int(bounds.X), int(bounds.Y), int(bounds.X+bounds.Width), int(bounds.Y+bounds.Height))
		sub := screen.SubImage(rect).(*ebiten.Image)
		screen = sub
	}

	// 绘制自己
	node.host.Draw(screen)

	// 绘制子组件
	if !node.host.DrawChildren(screen) {
		for _, child := range node.children {
			e.drawNode(screen, child)
		}
	}
}

// ==================== 焦点 ====================

func (e *Engine) setFocusHost(host Host) {
	if e.focusHost == host {
		return
	}
	if e.focusHost != nil {
		e.dispatchEvent(e.focusHost, &Event{Type: EventFocusOut, Target: e.focusHost})
	}
	e.focusHost = host
	if host != nil {
		e.dispatchEvent(host, &Event{Type: EventFocusIn, Target: host})
	}
}

func (e *Engine) collectFocusableHosts() []Host {
	var hosts []Host
	var walk func(*renderNode)
	walk = func(node *renderNode) {
		if node == nil {
			return
		}
		if node.host != nil && node.host.IsFocusable() {
			hosts = append(hosts, node.host)
		}
		for _, child := range node.children {
			walk(child)
		}
	}
	walk(e.rootNode)
	return hosts
}

func (e *Engine) focusNext(reverse bool) {
	focusables := e.collectFocusableHosts()
	if len(focusables) == 0 {
		return
	}
	idx := -1
	for i, h := range focusables {
		if h == e.focusHost {
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
	e.setFocusHost(focusables[idx])
}

// ==================== 事件 ====================

func (e *Engine) handleEvents() {
	mx, my := ebiten.CursorPosition()
	fx, fy := float32(mx), float32(my)

	// 每帧重置光标为默认箭头，组件可通过 MouseMove/Click 覆盖为其他形状
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)

	// Hover 检测：未拖拽时向当前 hover 组件发送 MouseMove，用于光标形状、hover 样式等
	if e.mouseDownTarget == nil {
		hoverTarget := e.hitTest(e.rootNode, fx, fy)
		if hoverTarget != nil {
			bounds := hoverTarget.GetLayoutBounds()
			e.dispatchEvent(hoverTarget, &Event{
				Type:   EventMouseMove,
				X:      fx,
				Y:      fy,
				LocalX: fx - bounds.X,
				LocalY: fy - bounds.Y,
				Target: hoverTarget,
			})
		}
	}

	// 1. 滚轮
	dx, dy := ebiten.Wheel()
	if dx != 0 || dy != 0 {
		target := e.hitTest(e.rootNode, fx, fy)
		if target != nil {
			e.dispatchEvent(target, &Event{
				Type:   EventScroll,
				X:      fx,
				Y:      fy,
				LocalX: fx - target.GetLayoutBounds().X,
				LocalY: fy - target.GetLayoutBounds().Y,
				DeltaX: float32(dx),
				DeltaY: float32(dy),
				Target: target,
			})
		}
	}

	// 2. 鼠标移动 / 拖拽
	if e.mouseDownTarget != nil {
		// 鼠标按下状态下移动：发送 EventMouseMove 给 mouseDownTarget
		if fx != e.lastMouseX || fy != e.lastMouseY {
			bounds := e.mouseDownTarget.GetLayoutBounds()
			e.dispatchEvent(e.mouseDownTarget, &Event{
				Type:   EventMouseMove,
				X:      fx,
				Y:      fy,
				LocalX: fx - bounds.X,
				LocalY: fy - bounds.Y,
				Button: e.mouseDownButton,
				Target: e.mouseDownTarget,
			})
		}
	}

	// 3. 点击（MouseDown / MouseUp / Click）
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		target := e.hitTest(e.rootNode, fx, fy)
		if target != nil {
			// 点击可焦点组件时设置焦点；点击非焦点区域时清除焦点
			if target.IsFocusable() {
				e.setFocusHost(target)
			} else if e.focusHost != nil {
				e.setFocusHost(nil)
			}
			bounds := target.GetLayoutBounds()
			e.mouseDownTarget = target
			e.mouseDownButton = ebiten.MouseButtonLeft
			e.dispatchEvent(target, &Event{
				Type:   EventMouseDown,
				X:      fx,
				Y:      fy,
				LocalX: fx - bounds.X,
				LocalY: fy - bounds.Y,
				Button: ebiten.MouseButtonLeft,
				Target: target,
			})
		} else {
			// 点击空白区域，清除焦点
			if e.focusHost != nil {
				e.setFocusHost(nil)
			}
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if e.mouseDownTarget != nil {
			bounds := e.mouseDownTarget.GetLayoutBounds()
			e.dispatchEvent(e.mouseDownTarget, &Event{
				Type:   EventMouseUp,
				X:      fx,
				Y:      fy,
				LocalX: fx - bounds.X,
				LocalY: fy - bounds.Y,
				Button: ebiten.MouseButtonLeft,
				Target: e.mouseDownTarget,
			})
			// 只有释放时仍在同一组件上才触发 Click
			target := e.hitTest(e.rootNode, fx, fy)
			if target == e.mouseDownTarget {
				e.dispatchEvent(e.mouseDownTarget, &Event{
					Type:   EventClick,
					X:      fx,
					Y:      fy,
					LocalX: fx - bounds.X,
					LocalY: fy - bounds.Y,
					Button: ebiten.MouseButtonLeft,
					Target: e.mouseDownTarget,
				})
			}
			e.mouseDownTarget = nil
		}
	}

	// 4. 键盘：Tab 切换焦点，Space/Enter 触发点击，其他键分发给焦点组件
	e.handleKeyboardEvents()

	// 更新上一帧鼠标位置
	e.lastMouseX = fx
	e.lastMouseY = fy
}

// GetFocusHost 返回当前获得焦点的 Host。
func (e *Engine) GetFocusHost() Host {
	return e.focusHost
}

func (e *Engine) handleKeyboardEvents() {
	// Tab 切换焦点
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		e.focusNext(ebiten.IsKeyPressed(ebiten.KeyShift))
		return
	}

	if e.focusHost == nil {
		return
	}

	// Space / Enter 触发焦点组件的点击
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		bounds := e.focusHost.GetLayoutBounds()
		e.dispatchEvent(e.focusHost, &Event{
			Type:   EventClick,
			X:      bounds.X + bounds.Width/2,
			Y:      bounds.Y + bounds.Height/2,
			LocalX: bounds.Width / 2,
			LocalY: bounds.Height / 2,
			Target: e.focusHost,
		})
		return
	}

	// 其他键盘事件分发给焦点组件
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		if inpututil.IsKeyJustPressed(key) {
			e.dispatchEvent(e.focusHost, &Event{
				Type: EventKeyDown,
				Key:  key,
				Target: e.focusHost,
			})
		}
		if inpututil.IsKeyJustReleased(key) {
			e.dispatchEvent(e.focusHost, &Event{
				Type: EventKeyUp,
				Key:  key,
				Target: e.focusHost,
			})
		}
	}
}

func (e *Engine) hitTest(node *renderNode, x, y float32) Host {
	if node == nil || node.host == nil {
		return nil
	}

	// 后序遍历：子组件优先
	for i := len(node.children) - 1; i >= 0; i-- {
		if h := e.hitTest(node.children[i], x, y); h != nil {
			return h
		}
	}

	el := node.host.GetElement()
	if el != nil && el.PointerEvents == PointerEventsNone {
		return nil
	}

	bounds := node.host.GetLayoutBounds()
	if x >= bounds.X && x <= bounds.X+bounds.Width &&
		y >= bounds.Y && y <= bounds.Y+bounds.Height {
		return node.host
	}

	return nil
}

func (e *Engine) dispatchEvent(target Host, event *Event) {
	// 从 target 向上冒泡，直到有 Host 消费事件
	node := e.findNodeByHost(e.rootNode, target)
	for node != nil && node.host != nil {
		if node.host.HandleEvent(event) {
			return
		}
		node = node.parent
	}
}

func (e *Engine) findNodeByHost(node *renderNode, host Host) *renderNode {
	if node == nil {
		return nil
	}
	if node.host == host {
		return node
	}
	for _, child := range node.children {
		if found := e.findNodeByHost(child, host); found != nil {
			return found
		}
	}
	return nil
}

// ==================== Host 更新 ====================

func (e *Engine) updateHosts(node *renderNode) {
	if node == nil || node.host == nil {
		return
	}
	_ = node.host.Update()
	for _, child := range node.children {
		e.updateHosts(child)
	}
}

// ==================== 动画管理 ====================

func (e *Engine) AddAnimation(anim Animation) {
	if anim == nil {
		return
	}
	e.animations = append(e.animations, anim)
}

func (e *Engine) updateAnimations(deltaTime float32) {
	if len(e.animations) == 0 {
		return
	}
	// 防止 deltaTime 过大（如窗口失去焦点后恢复）
	if deltaTime > 0.1 {
		deltaTime = 1.0 / 60.0
	}
	active := e.animations[:0]
	for _, anim := range e.animations {
		if anim.Update(deltaTime) {
			active = append(active, anim)
		}
	}
	e.animations = active
}

// ==================== Yoga 工具 ====================

func (e *Engine) insertYogaChild(parent, child Host) {
	if parent == nil || child == nil {
		return
	}
	pel := parent.GetElement()
	cel := child.GetElement()
	if pel == nil || pel.Yoga == nil || cel == nil || cel.Yoga == nil {
		return
	}
	pel.Yoga.InsertChild(cel.Yoga, pel.Yoga.GetChildCount())
}

func (e *Engine) removeYogaChild(parent, child Host) {
	if parent == nil || child == nil {
		return
	}
	pel := parent.GetElement()
	cel := child.GetElement()
	if pel == nil || pel.Yoga == nil || cel == nil || cel.Yoga == nil {
		return
	}
	pel.Yoga.RemoveChild(cel.Yoga)
}

// SetRootWidget 设置根 Widget，使 Engine 能正确处理根级别的 Invalidate。
func (e *Engine) SetRootWidget(w Widget) {
	e.rootWidget = w
}

// ==================== Portal ====================

// MountOverlay 将组件挂载到 overlay 层，不参与主树布局。
func (e *Engine) MountOverlay(c Component, anchor Host) *renderNode {
	node := e.mount(nil, c)
	if node != nil {
		e.overlays = append(e.overlays, node)
	}
	return node
}

// UnmountOverlay 从 overlay 层卸载。
func (e *Engine) UnmountOverlay(node *renderNode) {
	if node == nil {
		return
	}
	e.unmount(node)
	for i, o := range e.overlays {
		if o == node {
			e.overlays = append(e.overlays[:i], e.overlays[i+1:]...)
			break
		}
	}
}
