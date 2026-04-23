package core

import "github.com/hajimehoshi/ebiten/v2"

// Component 是 Widget 和 Host 的统一标记类型。
// 外部包不能随意实现，必须通过 BaseWidget 或内置 Host。
type Component interface {
	isComponent()
}

// Widget 是用户自定义组件接口。
// 用户组件只负责描述 UI（Render）和逻辑（Hooks、生命周期），
// 不参与直接绘制、布局和事件处理。
type Widget interface {
	Component

	// Render 返回 UI 描述。返回 Host 或另一个 Widget，nil 表示不渲染。
	Render() Component

	// === 生命周期（BaseWidget 提供空默认实现）===
	ComponentDidMount()
	ComponentWillUnmount()
	ComponentDidUpdate(prevProps, prevState any)
	ShouldComponentUpdate(nextProps any) bool

	// === 命令式引用（框架在挂载后注入）===
	SetHostRef(host Host)
	GetHostRef() Host

	// === 内部（框架调用）===
	setEngine(e *Engine)
	getEngine() *Engine
	resetHooks()
}

// Host 是宿主组件接口（View / Text / Button / Image / ScrollView 等）。
// 它们是真实 UI 元素，直接参与布局、绘制、交互，没有生命周期概念。
type Host interface {
	Component

	// === 布局 ===
	GetElement() *Element
	GetLayoutBounds() LayoutBounds
	UpdateLayoutBounds(parentX, parentY float32)

	// LayoutChildren 布局钩子。
	// 返回 true = 自行处理完毕，框架不再用 Yoga 计算子节点。
	LayoutChildren(parentWidth, parentHeight float32) bool

	// === 绘制 ===
	Draw(screen *ebiten.Image)

	// DrawChildren 绘制钩子。
	// 返回 true = 自行处理完毕，框架不再递归绘制子组件。
	DrawChildren(screen *ebiten.Image) bool

	// ShouldClipChildren 是否将子组件裁剪到自身边界。
	ShouldClipChildren() bool

	// === 交互 ===
	// HandleEvent 处理事件。返回 true = 事件被消费，停止冒泡。
	HandleEvent(e *Event) bool

	// Update 每帧更新（动画等）。
	Update() error

	// === 子组件 ===
	GetChildren() []Component
	AddChild(child Component)
	RemoveChild(child Component)
	ClearChildren()

	// === 滚动偏移 ===
	// GetScrollOffset 返回该 Host 对其子组件的滚动偏移量。
	// 默认实现返回 (0, 0)。
	GetScrollOffset() (x, y float32)

	// === 焦点 ===
	// IsFocusable 返回该 Host 是否可接收焦点。
	IsFocusable() bool
	// SetFocusable 设置该 Host 是否可接收焦点。
	SetFocusable(bool)

	// === Engine 引用 ===
	// SetEngine 由框架注入 Engine 引用。
	SetEngine(e *Engine)
	// GetEngine 返回关联的 Engine。
	GetEngine() *Engine
}
