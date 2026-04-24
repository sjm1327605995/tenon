package core

// Widget 是轻量的结构描述对象。
// 每个 Widget 只实现 Build()，产出初始 Element 树。
// 框架在挂载时调用一次 Build()，之后不再主动调用，除非用户显式请求。
type Widget interface {
	// Build 构建并返回 Element 树。只会在以下时机被框架调用：
	// 1. 首次挂载
	// 2. 用户调用 RequestBuild() 后，下一帧刷新前
	Build() Element

	// RequestBuild 请求框架在下一帧重新调用 Build()。
	// 仅用于结构变化（条件渲染、列表增删、页面切换）。
	// 属性变化不应调用此方法，直接通过 Element Setter 更新即可。
	RequestBuild()

	// OnMount / OnUnmount 生命周期
	OnMount(engine *Engine)
	OnUnmount()
}

// BaseWidget 提供 Widget 的默认实现。
type BaseWidget struct {
	self         Widget
	engine       *Engine
	needBuild    bool
	rootElement  Element
}

// Init 初始化 BaseWidget，必须在子类构造函数中调用。
func (b *BaseWidget) Init(self Widget) {
	b.self = self
}

// Build 默认返回 nil，子类必须覆盖。
func (b *BaseWidget) Build() Element { return nil }

// RequestBuild 标记需要在下一帧重建。
func (b *BaseWidget) RequestBuild() {
	b.needBuild = true
	if b.engine != nil {
		b.engine.scheduleBuild(b)
	}
}

func (b *BaseWidget) OnMount(engine *Engine) { b.engine = engine }
func (b *BaseWidget) OnUnmount()              {}

// GetRootElement 返回该 Widget 已挂载的根 Element。
func (b *BaseWidget) GetRootElement() Element { return b.rootElement }

// SetRootElement 设置该 Widget 的根 Element，框架内部使用。
func (b *BaseWidget) SetRootElement(el Element) { b.rootElement = el }
