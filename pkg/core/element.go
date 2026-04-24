package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// Element 是持久化的渲染节点接口。
// View、Text、Button 等 Native 组件直接实现此接口。
type Element interface {
	// === 树关系 ===
	GetParent() Element
	SetParent(p Element)
	GetChildren() []Element
	AppendChild(child Element)
	RemoveChild(child Element)
	ClearChildren()

	// === Yoga / 布局 ===
	GetYoga() *yoga.Node
	GetBounds() LayoutBounds
	SetBounds(b LayoutBounds)

	// === 绘制 / 事件 / 更新 ===
	Draw(screen *ebiten.Image)
	HandleEvent(e *Event) bool
	Update() error

	// === 生命周期 ===
	OnMount(engine *Engine)
	OnUnmount()

	// === 脏标记 ===
	Mark(flags ElementFlags)
	GetFlags() ElementFlags
	ClearDirty()

	// === 类型与标识 ===
	ElementType() string
	SetKey(key string)
	GetKey() string

	// === 样式标签 ===
	SetTag(tag string)
	GetTag() string
	SetClass(classes ...string)
	GetClass() []string

	// === Engine ===
	SetEngine(engine *Engine)
	GetEngine() *Engine
}

// ElementFlags 使用 uint64 bitmap 打包所有状态标志。
// 低 32 位用于持久状态（可见性、焦点等），高 32 位用于脏标记（每帧清除）。
type ElementFlags uint64

const (
	// === 持久状态（低 32 位，不会被 ClearDirty 清除）===
	FlagVisible ElementFlags = 1 << iota // 是否可见
	FlagFocusable                        // 是否可聚焦
	FlagClipChildren                     // 是否裁剪子节点到自身边界

	_ = iota + 29 // 占位，确保脏标记从第 32 位开始

	// === 脏标记（高 32 位，flushDirtyElements 后清除）===
	FlagNeedMeasure ElementFlags = 1 << 32 // 需要重新测量（如文字排版）
	FlagNeedLayout  ElementFlags = 1 << 33 // 布局属性变了（width/margin/flex...）
	FlagNeedDraw    ElementFlags = 1 << 34 // 视觉属性变了（color/text/radius...）
)

// LayoutBounds 描述组件在屏幕上的位置和尺寸。
type LayoutBounds struct {
	X, Y, Width, Height float32
}

// BorderRadius 描述四个角的圆角半径。
type BorderRadius struct {
	TopLeft, TopRight, BottomRight, BottomLeft float32
}

// PointerEvents 控制组件是否响应指针事件。
type PointerEvents int

const (
	PointerEventsAuto PointerEvents = iota
	PointerEventsNone
)

// ==================== BaseElement ====================

// BaseElement 提供 Element 接口的默认实现。
// 所有 Native 组件（View、Text、Button...）内嵌 BaseElement 即可。
type BaseElement struct {
	self      Element
	engine    *Engine
	yoga      *yoga.Node
	parent    Element
	children  []Element
	bounds    LayoutBounds
	flags     ElementFlags // uint64 bitmap：低32位持久状态，高32位脏标记
	key       string
	tag       string
	classes   []string
}

// Init 初始化 BaseElement，必须在子类构造函数中调用。
func (b *BaseElement) Init(self Element) {
	b.self = self
	b.yoga = yoga.NewNode()
	b.flags = FlagVisible // 默认可见
}

// === 树关系 ===

func (b *BaseElement) GetParent() Element     { return b.parent }
func (b *BaseElement) SetParent(p Element)    { b.parent = p }
func (b *BaseElement) GetChildren() []Element { return b.children }

func (b *BaseElement) AppendChild(child Element) {
	if child == nil {
		return
	}
	b.children = append(b.children, child)
	child.SetParent(b.self)
	if b.yoga != nil && child.GetYoga() != nil {
		b.yoga.InsertChild(child.GetYoga(), b.yoga.GetChildCount())
	}
	if b.engine != nil {
		b.engine.onElementMounted(child)
	}
}

func (b *BaseElement) RemoveChild(child Element) {
	for i, c := range b.children {
		if c == child {
			b.children = append(b.children[:i], b.children[i+1:]...)
			child.SetParent(nil)
			if b.yoga != nil && child.GetYoga() != nil {
				b.yoga.RemoveChild(child.GetYoga())
			}
			child.OnUnmount()
			return
		}
	}
}

func (b *BaseElement) ClearChildren() {
	for _, c := range b.children {
		c.SetParent(nil)
		c.OnUnmount()
	}
	b.children = b.children[:0]
	if b.yoga != nil {
		b.yoga.RemoveAllChildren()
	}
}

// === Yoga / 布局 ===

func (b *BaseElement) GetYoga() *yoga.Node      { return b.yoga }
func (b *BaseElement) GetBounds() LayoutBounds   { return b.bounds }
func (b *BaseElement) SetBounds(lb LayoutBounds) { b.bounds = lb }

// === 绘制 / 事件 / 更新（默认空实现）===

func (b *BaseElement) Draw(screen *ebiten.Image)     {}
func (b *BaseElement) HandleEvent(e *Event) bool     { return false }
func (b *BaseElement) Update() error                 { return nil }
func (b *BaseElement) OnMount(engine *Engine)        { b.engine = engine }
func (b *BaseElement) OnUnmount()                    {}

// === 标志位操作（uint64 bitmap）===

// Mark 设置脏标记，并通知引擎加入刷新队列。
func (b *BaseElement) Mark(flags ElementFlags) {
	if b.flags&FlagDirtyMask == 0 && b.engine != nil {
		b.engine.markDirty(b.self)
	}
	b.flags |= flags
}

// HasFlag 检查是否包含指定标志。
func (b *BaseElement) HasFlag(f ElementFlags) bool { return b.flags&f != 0 }

// SetFlag 设置标志位（持久状态或脏标记）。
func (b *BaseElement) SetFlag(f ElementFlags) { b.flags |= f }

// ClearFlag 清除指定标志位。
func (b *BaseElement) ClearFlag(f ElementFlags) { b.flags &^= f }

// GetFlags 返回完整 bitmap。
func (b *BaseElement) GetFlags() ElementFlags { return b.flags }

// ClearDirty 只清除高 32 位的脏标记，保留低 32 位持久状态。
func (b *BaseElement) ClearDirty() { b.flags &^= FlagDirtyMask }

// FlagDirtyMask 用于隔离高 32 位脏标记。
const FlagDirtyMask ElementFlags = 0xFFFFFFFF00000000

// === 类型与标识 ===

func (b *BaseElement) ElementType() string   { return "BaseElement" }
func (b *BaseElement) SetKey(key string)     { b.key = key }
func (b *BaseElement) GetKey() string        { return b.key }
func (b *BaseElement) SetTag(tag string)     { b.tag = tag }
func (b *BaseElement) GetTag() string        { return b.tag }
func (b *BaseElement) SetClass(c ...string)  { b.classes = append(b.classes[:0], c...) }
func (b *BaseElement) GetClass() []string    { return b.classes }

// === Engine ===

func (b *BaseElement) SetEngine(e *Engine) { b.engine = e }
func (b *BaseElement) GetEngine() *Engine  { return b.engine }

// ==================== 链式 API（布局）====================

func (b *BaseElement) SetWidth(v float32) Element {
	b.yoga.StyleSetWidth(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetWidthPercent(v float32) Element {
	b.yoga.StyleSetWidthPercent(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetHeight(v float32) Element {
	b.yoga.StyleSetHeight(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetHeightPercent(v float32) Element {
	b.yoga.StyleSetHeightPercent(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetMinWidth(v float32) Element {
	b.yoga.StyleSetMinWidth(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetMinHeight(v float32) Element {
	b.yoga.StyleSetMinHeight(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetMaxWidth(v float32) Element {
	b.yoga.StyleSetMaxWidth(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetMaxHeight(v float32) Element {
	b.yoga.StyleSetMaxHeight(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetFlexDirection(dir yoga.FlexDirection) Element {
	b.yoga.StyleSetFlexDirection(dir)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetJustifyContent(v yoga.Justify) Element {
	b.yoga.StyleSetJustifyContent(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetAlignItems(v yoga.Align) Element {
	b.yoga.StyleSetAlignItems(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetAlignSelf(v yoga.Align) Element {
	b.yoga.StyleSetAlignSelf(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetFlexGrow(v float32) Element {
	b.yoga.StyleSetFlexGrow(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetFlexShrink(v float32) Element {
	b.yoga.StyleSetFlexShrink(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetFlexBasis(v float32) Element {
	b.yoga.StyleSetFlexBasis(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetFlexWrap(v yoga.Wrap) Element {
	b.yoga.StyleSetFlexWrap(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetPadding(edge yoga.Edge, v float32) Element {
	b.yoga.StyleSetPadding(edge, v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetMargin(edge yoga.Edge, v float32) Element {
	b.yoga.StyleSetMargin(edge, v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetBorder(edge yoga.Edge, v float32) Element {
	b.yoga.StyleSetBorder(edge, v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetPosition(edge yoga.Edge, v float32) Element {
	b.yoga.StyleSetPosition(edge, v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetPositionType(v yoga.PositionType) Element {
	b.yoga.StyleSetPositionType(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetAspectRatio(v float32) Element {
	b.yoga.StyleSetAspectRatio(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetDisplay(v yoga.Display) Element {
	b.yoga.StyleSetDisplay(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetOverflow(v yoga.Overflow) Element {
	b.yoga.StyleSetOverflow(v)
	b.Mark(FlagNeedLayout)
	return b.self
}

func (b *BaseElement) SetGap(gutter yoga.Gutter, v float32) Element {
	b.yoga.StyleSetGap(gutter, v)
	b.Mark(FlagNeedLayout)
	return b.self
}

// Add 是 AppendChild 的链式语法糖。
func (b *BaseElement) Add(children ...Element) Element {
	for _, c := range children {
		b.AppendChild(c)
	}
	return b.self
}

// SetVisible 控制是否可见（视觉属性，不重排）。
func (b *BaseElement) SetVisible(v bool) Element {
	if v {
		b.flags |= FlagVisible
	} else {
		b.flags &^= FlagVisible
	}
	b.Mark(FlagNeedDraw)
	return b.self
}

// IsVisible 检查可见性。
func (b *BaseElement) IsVisible() bool { return b.flags&FlagVisible != 0 }
