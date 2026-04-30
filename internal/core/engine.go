package core

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// Engine 是 Tenon 核心，负责 Widget → Element 的构建、布局计算、绘制和事件。
type Engine struct {
	rootWidget  Widget
	rootElement Element

	// 帧调度
	dirtyElements []Element
	buildQueue    []Widget

	// 窗口
	screenWidth   int
	screenHeight  int
	pendingResize bool

	// 弹出层
	overlays []Element

	// 输入状态
	lastMouseX      float32
	lastMouseY      float32
	mouseDownTarget Element
	mouseDownX      float32
	mouseDownY      float32
	mouseDownButton MouseButton
	hoverTarget     Element
	focusTarget     Element

	// 时间
	lastFrameTime time.Time

	// 性能指标
	perf PerfMetrics

	// 生命周期追踪
	lifecycleLogs   []LifecycleLog
	maxLifecycleLog int

	// 渲染依赖追踪
	renderTracker renderTracker

	// 动画
	animations []Animation

	// 子引擎
	builder         *BuildEngine
	renderer        *RenderEngine
	eventDispatcher *EventDispatcher
	animator        *AnimationEngine

	// 事件注册表
	eventRegistry *EventRegistry

	// 样式注册表
	styleRegistry *StyleRegistry

	// 脏标记事件总线：收集元素属性变更事件，合并后统一刷新
	dirtyBus *DirtyEventBus

	// 调试器
	debugBus DebugEventBus
}

// NewEngine 创建引擎。
func NewEngine(rootWidget Widget, width, height int) *Engine {
	e := &Engine{
		rootWidget:      rootWidget,
		screenWidth:     width,
		screenHeight:    height,
		lastFrameTime:   time.Now(),
		eventRegistry:   NewEventRegistry(),
		styleRegistry:   NewStyleRegistry(),
		dirtyBus:        &DirtyEventBus{},
		maxLifecycleLog: 200,
	}
	e.builder = newBuildEngine(e)
	e.renderer = newRenderEngine(e)
	e.eventDispatcher = newEventDispatcher(e)
	e.animator = newAnimationEngine(e)
	return e
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

// SetDebugger 设置调试事件总线。
func (e *Engine) SetDebugger(d DebugEventBus) {
	e.debugBus = d
	SetStateDebugHook(func(elementType, key string, oldValue, newValue interface{}) {
		d.AddStateLog(elementType, key, "", oldValue, newValue)
	})
}

func (e *Engine) GetRootElement() Element {
	return e.rootElement
}

func (e *Engine) SetStyleRegistry(r *StyleRegistry) {
	e.styleRegistry = r
}

func (e *Engine) RegisterStyle(tag string, apply StyleFunc) {
	if e.styleRegistry == nil {
		e.styleRegistry = NewStyleRegistry()
	}
	e.styleRegistry.RegisterStyle(tag, apply)
}

func (e *Engine) RegisterStyleForType(elemType, tag string, apply StyleFunc) {
	if e.styleRegistry == nil {
		e.styleRegistry = NewStyleRegistry()
	}
	e.styleRegistry.RegisterStyleForType(elemType, tag, apply)
}

func (e *Engine) beginRenderTracking(widget Widget) {
	e.renderTracker.widget = widget
	e.renderTracker.active = true
	e.renderTracker.states = nil
	activeRenderTracker = &e.renderTracker
}

func (e *Engine) endRenderTracking() []stateNotifier {
	e.renderTracker.active = false
	states := e.renderTracker.states
	e.renderTracker.states = nil
	activeRenderTracker = nil
	return states
}

func (e *Engine) Mount() {
	if e.rootWidget == nil {
		return
	}
	setThemeChangeNotifier(e.RequestRedrawAll)
	if bw, ok := e.rootWidget.(interface{ SetEngine(*Engine) }); ok {
		bw.SetEngine(e)
	}
	e.rootWidget.OnMount()
	e.rootElement = e.rootWidget.Render()
	if e.rootElement != nil {
		e.onElementMounted(e.rootElement)
		e.initOverlayLayer()
		e.calculateLayout()
	}
}

func (e *Engine) initOverlayLayer() {
}

func (e *Engine) AddOverlay(el Element) {
	for _, o := range e.overlays {
		if o == el {
			return
		}
	}
	e.overlays = append(e.overlays, el)
}

func (e *Engine) RemoveOverlay(el Element) {
	for i, o := range e.overlays {
		if o == el {
			e.overlays = append(e.overlays[:i], e.overlays[i+1:]...)
			return
		}
	}
}

func (e *Engine) IsOverlay(el Element) bool {
	for _, o := range e.overlays {
		if o == el {
			return true
		}
	}
	return false
}

func (e *Engine) RequestRedrawAll() {
	if e.rootWidget != nil {
		e.builder.scheduleBuild(e.rootWidget)
	}
}

// onElementMounted 递归给新挂载的 Element 注入 Engine。
func (e *Engine) onElementMounted(el Element) {
	if el == nil {
		return
	}
	el.SetEngine(e)
	if e.styleRegistry != nil {
		e.styleRegistry.ApplyStyles(el)
	}
	el.OnMount(e)
	el.FlushDelayedListeners()
	e.recordLifecycle("mount", el)
	// 补发挂载前已存在的脏标记（创建时 engine 为 nil，导致 Mark 只设了 flag 没进 dirtyBus）
	if el.GetFlags()&FlagDirtyMask != 0 {
		e.dirtyBus.Post(el)
	}
	for _, child := range el.GetChildren() {
		if child.GetEngine() == nil {
			e.onElementMounted(child)
		}
	}
}

// ==================== 帧循环 ====================

func (e *Engine) Update() error {
	frameStart := time.Now()

	if e.pendingResize {
		e.pendingResize = false
		e.calculateLayout()
		if e.rootElement != nil {
			e.eventDispatcher.broadcastEvent(e.rootElement, &Event{
				Type:   EventResize,
				Width:  float32(e.screenWidth),
				Height: float32(e.screenHeight),
			})
		}
	}

	now := time.Now()
	deltaTime := float32(now.Sub(e.lastFrameTime).Seconds())
	e.lastFrameTime = now

	e.builder.flushBuildQueue()

	eventsStart := time.Now()
	e.eventDispatcher.handleEvents()
	e.perf.LastEventTime = time.Since(eventsStart)

	layoutStart := time.Now()
	e.flushDirtyElements()
	e.perf.LastLayoutTime = time.Since(layoutStart)

	e.animator.updateAnimations(deltaTime)

	if e.rootElement != nil {
		e.updateElements(e.rootElement)
	}
	for _, overlay := range e.overlays {
		if overlay != nil && overlay.IsVisible() {
			e.updateElements(overlay)
		}
	}

	e.perf.LastFrameTime = time.Since(frameStart)
	e.perf.FrameCount++
	e.perf.LastFrameTimestamp = time.Now()

	return nil
}

// AddAnimation 注册一个动画到引擎。
func (e *Engine) AddAnimation(a Animation) {
	e.animator.AddAnimation(a)
}

func (e *Engine) Draw(screen *ebiten.Image) {
	e.renderer.drawScreen(screen)
}

// ==================== 结构重建 ====================

func (e *Engine) scheduleBuild(w Widget) {
	e.builder.scheduleBuild(w)
}

// ==================== 脏标记事件总线 ====================

// DirtyEventBus 负责收集元素的脏标记事件，合并后统一刷新。
// 同一元素的多次属性变更在 pending 中只保留一个，flag 合并由 Element 自身的 bitmap 完成。
type DirtyEventBus struct {
	pending []Element
}

// Post 投递一个元素的脏标记事件。已在 pending 中的元素会被忽略（去重）。
func (b *DirtyEventBus) Post(el Element) {
	for _, p := range b.pending {
		if p == el {
			return
		}
	}
	b.pending = append(b.pending, el)
}

// Flush 取出所有待刷新元素并清空队列。
func (b *DirtyEventBus) Flush() []Element {
	batch := b.pending
	b.pending = nil
	return batch
}

// Len 返回当前待处理事件数。
func (b *DirtyEventBus) Len() int {
	return len(b.pending)
}

// ==================== 脏标记刷新 ====================

func (e *Engine) flushDirtyElements() {
	// 从事件总线取出合并后的元素批次
	if e.dirtyBus.Len() == 0 {
		return
	}
	e.dirtyElements = e.dirtyBus.Flush()
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
	// 同时检查已刷入 dirtyElements 的和事件总线中 pending 的
	for _, el := range e.dirtyElements {
		if el.GetFlags()&FlagNeedLayout != 0 {
			return true
		}
	}
	for _, el := range e.dirtyBus.pending {
		if el.GetFlags()&FlagNeedLayout != 0 {
			return true
		}
	}
	return false
}

// ==================== 布局 ====================

func (e *Engine) calculateLayout() {
	if e.rootElement == nil {
		return
	}
	layoutEl, ok := e.rootElement.(LayoutElement)
	if !ok {
		return
	}
	rootYoga := layoutEl.GetYoga()
	if rootYoga == nil {
		return
	}
	rootYoga.StyleSetWidth(float32(e.screenWidth))
	rootYoga.StyleSetHeight(float32(e.screenHeight))
	rootYoga.CalculateLayout(float32(e.screenWidth), float32(e.screenHeight), yoga.DirectionLTR)

	if e.rootElement != nil {
		e.updateBounds(e.rootElement, 0, 0)
	}

	if e.debugBus != nil && e.debugBus.IsEnabled() {
		e.debugBus.CaptureLayout("calculateLayout")
	}
}

func (e *Engine) updateBounds(el Element, parentX, parentY float32) {
	if el == nil {
		return
	}
	layoutEl, ok := el.(LayoutElement)
	if !ok {
		return
	}
	y := layoutEl.GetYoga()
	if y == nil {
		return
	}
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

// ==================== Element 每帧更新 ====================

func (e *Engine) updateElements(el Element) {
	if el == nil {
		return
	}
	_ = el.Update()
	for _, child := range el.GetChildren() {
		if e.IsOverlay(child) {
			continue
		}
		e.updateElements(child)
	}
}

// ==================== 事件系统 ====================

// GetFocusTarget returns the currently focused element.
func (e *Engine) GetFocusTarget() Element {
	return e.eventDispatcher.GetFocusTarget()
}

func (e *Engine) GetHoverTarget() Element {
	return e.eventDispatcher.GetHoverTarget()
}

// ==================== 性能指标 ====================

type PerfMetrics struct {
	FrameCount         int           `json:"frameCount"`
	LastFrameTime      time.Duration `json:"lastFrameTime"`
	LastFrameTimestamp time.Time     `json:"lastFrameTimestamp"`
	LastDrawTime       time.Duration `json:"lastDrawTime"`
	LastLayoutTime     time.Duration `json:"lastLayoutTime"`
	LastEventTime      time.Duration `json:"lastEventTime"`
	ElementCount       int           `json:"elementCount"`
	AnimationCount     int           `json:"animationCount"`
	DirtyElementCount  int           `json:"dirtyElementCount"`
	EventListenerCount int           `json:"eventListenerCount"`
}

func (e *Engine) GetPerfMetrics() PerfMetrics {
	e.perf.ElementCount = e.countElements(e.rootElement)
	e.perf.AnimationCount = len(e.animations)
	e.perf.DirtyElementCount = len(e.dirtyElements)
	if e.eventRegistry != nil {
		e.perf.EventListenerCount = e.countEventListeners()
	}
	return e.perf
}

func (e *Engine) countElements(el Element) int {
	if el == nil {
		return 0
	}
	count := 1
	for _, child := range el.GetChildren() {
		count += e.countElements(child)
	}
	return count
}

func (e *Engine) countEventListeners() int {
	if e.eventRegistry == nil {
		return 0
	}
	info := e.eventRegistry.DebugInfo()
	total := 0
	for _, l := range info {
		total += l.Count
	}
	return total
}

// ==================== 生命周期追踪 ====================

type LifecycleLog struct {
	Timestamp   time.Time `json:"timestamp"`
	Action      string    `json:"action"`
	ElementType string    `json:"elementType"`
	ElementKey  string    `json:"elementKey,omitempty"`
}

func (e *Engine) recordLifecycle(action string, el Element) {
	if e.maxLifecycleLog == 0 {
		return
	}
	entry := LifecycleLog{
		Timestamp:   time.Now(),
		Action:      action,
		ElementType: el.ElementType(),
		ElementKey:  el.GetKey(),
	}
	e.lifecycleLogs = append(e.lifecycleLogs, entry)
	if len(e.lifecycleLogs) > e.maxLifecycleLog {
		e.lifecycleLogs = e.lifecycleLogs[len(e.lifecycleLogs)-e.maxLifecycleLog:]
	}
}

func (e *Engine) GetLifecycleLogs() []LifecycleLog {
	return e.lifecycleLogs
}

// GetEventRegistryDebugInfo 返回事件注册表的调试信息。
func (e *Engine) GetEventRegistryDebugInfo() []DebugListenerInfo {
	if e.eventRegistry == nil {
		return nil
	}
	return e.eventRegistry.DebugInfo()
}

// GetEventListenersForElement 返回指定元素的事件监听器信息。
func (e *Engine) GetEventListenersForElement(el Element) []DebugListenerInfo {
	if e.eventRegistry == nil || el == nil {
		return nil
	}
	return e.eventRegistry.DebugListenersForElement(el)
}
