package core

// Widget 是轻量的结构描述对象。
// 每个 Widget 只实现 Build()，产出初始 Element 树。
// 框架在挂载时调用一次 Build()，之后不再主动调用，除非用户显式请求。
type Widget interface {
	// Render 构建并返回 Element 树。只会在以下时机被框架调用：
	// 1. 首次挂载
	// 2. 用户调用 RequestBuild() 后，下一帧刷新前
	Render() Element

	// RequestBuild 请求框架在下一帧重新调用 Build()。
	// 仅用于结构变化（条件渲染、列表增删、页面切换）。
	// 属性变化不应调用此方法，直接通过 Element Setter 更新即可。
	RequestBuild()

	// OnMount / OnUnmount 生命周期
	OnMount()
	OnUnmount()
}

// BaseWidget 提供 Widget 的默认实现。
type BaseWidget struct {
	self        Widget
	engine      *Engine
	needBuild   bool
	rootElement Element

	// 自动依赖追踪：State 订阅清理函数
	stateCleanups []func()
}

// Init 初始化 BaseWidget，必须在子类构造函数中调用。
func (b *BaseWidget) Init(self Widget) {
	b.self = self
}

// Render 默认返回 nil，子类必须覆盖。
func (b *BaseWidget) Render() Element { return nil }

// RenderWithTracking 执行 Render() 并自动追踪 State 依赖。
// 框架在 flushBuildQueue 中调用此方法，用户无需关心。
func (b *BaseWidget) RenderWithTracking() Element {
	if b.engine != nil {
		b.engine.beginRenderTracking(b.self)
	}

	el := b.self.Render()

	var states []stateNotifier
	if b.engine != nil {
		states = b.engine.endRenderTracking()
	}

	// 清理上一帧的订阅
	for _, cleanup := range b.stateCleanups {
		cleanup()
	}
	b.stateCleanups = nil

	// 为本次 Render 访问的所有 State 建立订阅
	// State 变化时自动触发 Widget 重建
	for _, s := range states {
		cleanup := s._subscribeRebuild(func() {
			b.self.RequestBuild()
		})
		b.stateCleanups = append(b.stateCleanups, cleanup)
	}

	return el
}

// RequestBuild 标记需要在下一帧重建。
func (b *BaseWidget) RequestBuild() {
	b.needBuild = true
	if b.engine != nil {
		b.engine.scheduleBuild(b)
	}
}

func (b *BaseWidget) SetEngine(engine *Engine) { b.engine = engine }
func (b *BaseWidget) OnMount()                           {}
func (b *BaseWidget) GetEngine() *Engine                { return b.engine }
func (b *BaseWidget) OnUnmount() {
	// 清理所有 State 订阅
	for _, cleanup := range b.stateCleanups {
		cleanup()
	}
	b.stateCleanups = nil
}

// GetRootElement 返回该 Widget 已挂载的根 Element。
func (b *BaseWidget) GetRootElement() Element { return b.rootElement }

// SetRootElement 设置该 Widget 的根 Element，框架内部使用。
func (b *BaseWidget) SetRootElement(el Element) { b.rootElement = el }
