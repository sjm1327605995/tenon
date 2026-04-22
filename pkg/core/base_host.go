package core

import "github.com/hajimehoshi/ebiten/v2"

// BaseHost 提供 Host 接口的默认实现。
// 所有内置宿主组件（View / Text / Button / Image）嵌入 BaseHost 即可。
type BaseHost struct {
	self      Host
	element   *Element
	children  []Component
	bounds    LayoutBounds
	focusable bool
	engine    *Engine
}

// Init 必须在宿主组件的构造函数中调用，参数传自己（*MyHost）。
func (b *BaseHost) Init(self Host) {
	b.self = self
	b.element = NewElement()
}

func (b *BaseHost) isComponent() {}

// === 布局 ===

func (b *BaseHost) GetElement() *Element          { return b.element }
func (b *BaseHost) GetLayoutBounds() LayoutBounds { return b.bounds }

func (b *BaseHost) UpdateLayoutBounds(parentX, parentY float32) {
	if b.element == nil || b.element.Yoga == nil {
		return
	}
	b.bounds.X = parentX + b.element.Yoga.LayoutLeft()
	b.bounds.Y = parentY + b.element.Yoga.LayoutTop()
	b.bounds.Width = b.element.Yoga.LayoutWidth()
	b.bounds.Height = b.element.Yoga.LayoutHeight()
}

// LayoutChildren 默认不处理，框架会继续用 Yoga 计算子节点。
func (b *BaseHost) LayoutChildren(parentWidth, parentHeight float32) bool { return false }

// === 绘制 ===

func (b *BaseHost) Draw(screen *ebiten.Image) {}

// DrawChildren 默认返回 false，框架会递归绘制子组件。
func (b *BaseHost) DrawChildren(screen *ebiten.Image) bool { return false }

func (b *BaseHost) ShouldClipChildren() bool { return false }

// === 交互 ===

func (b *BaseHost) HandleEvent(e *Event) bool { return false }
func (b *BaseHost) Update() error             { return nil }

// === 子组件 ===

func (b *BaseHost) GetChildren() []Component { return b.children }

func (b *BaseHost) AddChild(child Component) {
	b.children = append(b.children, child)
}

func (b *BaseHost) RemoveChild(child Component) {
	for i, c := range b.children {
		if c == child {
			b.children = append(b.children[:i], b.children[i+1:]...)
			return
		}
	}
}

func (b *BaseHost) ClearChildren() {
	b.children = nil
}

// GetScrollOffset 默认返回 (0, 0)。
func (b *BaseHost) GetScrollOffset() (x, y float32) {
	return 0, 0
}

// IsFocusable 返回是否可接收焦点。
func (b *BaseHost) IsFocusable() bool {
	return b.focusable
}

// SetFocusable 设置是否可接收焦点。
func (b *BaseHost) SetFocusable(v bool) {
	b.focusable = v
}

// SetEngine 由框架在挂载时注入 Engine 引用。
func (b *BaseHost) SetEngine(e *Engine) {
	b.engine = e
}

// GetEngine 返回关联的 Engine，可能为 nil。
func (b *BaseHost) GetEngine() *Engine {
	return b.engine
}
