package core

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

type LayoutBounds struct {
	X, Y, Width, Height float32
}

type Element struct {
	Visible         bool
	Yoga            *yoga.Node
	BackgroundColor color.Color
	BorderColor     color.Color
	ShadowColor     color.Color
	ShadowBlur      float32
	ShadowOffsetX   float32
	ShadowOffsetY   float32
	PointerEvents   PointerEvents
	BorderRadius    BorderRadius
	Key             string // 用于 reconciler 匹配
}

type BorderRadius struct {
	TopLeft     float32
	TopRight    float32
	BottomRight float32
	BottomLeft  float32
}

type PointerEvents int

const (
	PointerEventsAuto PointerEvents = iota
	PointerEventsNone
)

// ============================================================================
// Component 接口（精简版，移除 Lifecycle）
// ============================================================================
type Component interface {
	// 渲染相关
	Draw(screen *ebiten.Image)
	DrawOverlay(screen *ebiten.Image)
	Update() error

	// 输入处理
	HandleInput() bool

	// 子组件管理
	GetChildren() []Component
	AddChild(child Component) error

	// 布局相关
	GetLayoutBounds() LayoutBounds
	GetElement() *Element
	Render() *Element
	UpdateLayoutBounds(parentX, parentY float32)

	// 脏检查机制
	ID() string
	MarkDirty()
	IsDirty() bool
	ClearDirty()
}

// ============================================================================
// BaseComponent - 基础组件实现
// ============================================================================
type BaseComponent struct {
	children []Component
	element  *Element
	yogaNode *yoga.Node
	bounds   LayoutBounds

	id    string
	dirty bool
	mu    sync.RWMutex

	self Component
}

func NewBaseComponent() BaseComponent {
	node := yoga.NewNode()
	return BaseComponent{
		children: make([]Component, 0),
		element: &Element{
			Visible:       true,
			Yoga:          node,
			PointerEvents: PointerEventsAuto,
		},
		yogaNode:  node,
		id:        generateComponentID(),
		dirty:     false,
	}
}

func (b *BaseComponent) Init(self Component) {
	b.self = self
}

func (b *BaseComponent) Self() Component {
	return b.self
}

// InternalComponent 接口实现

func (b *BaseComponent) GetChildren() []Component {
	return b.children
}

func (b *BaseComponent) AddChild(child Component) error {
	childYoga := child.GetElement().Yoga
	if childYoga != nil && b.yogaNode != nil {
		if childYoga.GetOwner() != nil {
			childYoga.GetOwner().RemoveChild(childYoga)
		}
		b.yogaNode.InsertChild(childYoga, b.yogaNode.GetChildCount())
	}
	b.children = append(b.children, child)
	return nil
}

func (b *BaseComponent) GetLayoutBounds() LayoutBounds {
	return b.bounds
}

func (b *BaseComponent) GetElement() *Element {
	return b.element
}

func (b *BaseComponent) Render() *Element {
	return b.element
}

func (b *BaseComponent) UpdateLayoutBounds(parentX, parentY float32) {
	if b.yogaNode == nil {
		return
	}

	relX := b.yogaNode.LayoutLeft()
	relY := b.yogaNode.LayoutTop()

	b.bounds = LayoutBounds{
		X:      parentX + relX,
		Y:      parentY + relY,
		Width:  b.yogaNode.LayoutWidth(),
		Height: b.yogaNode.LayoutHeight(),
	}
}

func (b *BaseComponent) Update() error {
	for _, child := range b.children {
		if err := child.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (b *BaseComponent) Draw(screen *ebiten.Image) {
}

func (b *BaseComponent) DrawOverlay(screen *ebiten.Image) {
}

func (b *BaseComponent) HandleInput() bool {
	return false
}

// CalculateLayout 计算布局（框架内部调用）
func CalculateLayout(root Component, width, height float32) {
	element := root.GetElement()
	if element == nil || element.Yoga == nil {
		return
	}

	element.Yoga.StyleSetWidth(width)
	element.Yoga.StyleSetHeight(height)
	element.Yoga.CalculateLayout(width, height, yoga.DirectionLTR)
	updateAllBounds(root, 0, 0)
}

func updateAllBounds(component Component, parentX, parentY float32) {
	if component == nil {
		return
	}

	component.UpdateLayoutBounds(parentX, parentY)
	bounds := component.GetLayoutBounds()

	for _, child := range component.GetChildren() {
		updateAllBounds(child, bounds.X, bounds.Y)
	}
}

var (
	componentCounter int64
	componentMutex   sync.Mutex
)

func generateComponentID() string {
	componentMutex.Lock()
	defer componentMutex.Unlock()

	componentCounter++
	return fmt.Sprintf("component_%d_%d", time.Now().UnixNano(), componentCounter)
}

func (b *BaseComponent) ID() string {
	return b.id
}

func (b *BaseComponent) MarkDirty() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dirty = true
}

func (b *BaseComponent) IsDirty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dirty
}

func (b *BaseComponent) ClearDirty() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dirty = false
}

// ============================================================================
// 样式设置方法（链式调用）
// ============================================================================

func (b *BaseComponent) ApplyStyle(style *Style) {
	if style != nil && b.yogaNode != nil {
		style.Apply(b.yogaNode)
	}
}

func (b *BaseComponent) ApplyVisualStyle(style *VisualStyle) {
	if style != nil && b.element != nil {
		style.Apply(b.element)
	}
}

func (b *BaseComponent) SetWidth(width float32) Component {
	b.element.Yoga.StyleSetWidth(width)
	return b.self
}

func (b *BaseComponent) SetHeight(height float32) Component {
	b.element.Yoga.StyleSetHeight(height)
	return b.self
}

func (b *BaseComponent) SetMinWidth(width float32) Component {
	b.element.Yoga.StyleSetMinWidth(width)
	return b.self
}

func (b *BaseComponent) SetMinHeight(height float32) Component {
	b.element.Yoga.StyleSetMinHeight(height)
	return b.self
}

func (b *BaseComponent) SetMaxWidth(width float32) Component {
	b.element.Yoga.StyleSetMaxWidth(width)
	return b.self
}

func (b *BaseComponent) SetMaxHeight(height float32) Component {
	b.element.Yoga.StyleSetMaxHeight(height)
	return b.self
}

func (b *BaseComponent) SetFlexDirection(dir yoga.FlexDirection) Component {
	b.element.Yoga.StyleSetFlexDirection(dir)
	return b.self
}

func (b *BaseComponent) SetJustifyContent(justify yoga.Justify) Component {
	b.element.Yoga.StyleSetJustifyContent(justify)
	return b.self
}

func (b *BaseComponent) SetAlignItems(align yoga.Align) Component {
	b.element.Yoga.StyleSetAlignItems(align)
	return b.self
}

func (b *BaseComponent) SetAlignSelf(align yoga.Align) Component {
	b.element.Yoga.StyleSetAlignSelf(align)
	return b.self
}

func (b *BaseComponent) SetFlexGrow(grow float32) Component {
	b.element.Yoga.StyleSetFlexGrow(grow)
	return b.self
}

func (b *BaseComponent) SetFlexShrink(shrink float32) Component {
	b.element.Yoga.StyleSetFlexShrink(shrink)
	return b.self
}

func (b *BaseComponent) SetFlexBasis(basis float32) Component {
	b.element.Yoga.StyleSetFlexBasis(basis)
	return b.self
}

func (b *BaseComponent) SetFlexWrap(wrap yoga.Wrap) Component {
	b.element.Yoga.StyleSetFlexWrap(wrap)
	return b.self
}

func (b *BaseComponent) SetPadding(edge yoga.Edge, value float32) Component {
	b.element.Yoga.StyleSetPadding(edge, value)
	return b.self
}

func (b *BaseComponent) SetMargin(edge yoga.Edge, value float32) Component {
	b.element.Yoga.StyleSetMargin(edge, value)
	return b.self
}

func (b *BaseComponent) SetBorder(edge yoga.Edge, value float32) Component {
	b.element.Yoga.StyleSetBorder(edge, value)
	return b.self
}

func (b *BaseComponent) SetPosition(edge yoga.Edge, value float32) Component {
	b.element.Yoga.StyleSetPosition(edge, value)
	return b.self
}

func (b *BaseComponent) SetPositionType(positionType yoga.PositionType) Component {
	b.element.Yoga.StyleSetPositionType(positionType)
	return b.self
}

func (b *BaseComponent) SetAspectRatio(ratio float32) Component {
	b.element.Yoga.StyleSetAspectRatio(ratio)
	return b.self
}

func (b *BaseComponent) SetDisplay(display yoga.Display) Component {
	b.element.Yoga.StyleSetDisplay(display)
	return b.self
}

func (b *BaseComponent) SetOverflow(overflow yoga.Overflow) Component {
	b.element.Yoga.StyleSetOverflow(overflow)
	return b.self
}

func (b *BaseComponent) SetBackgroundColor(clr color.Color) Component {
	b.element.BackgroundColor = clr
	return b.self
}

func (b *BaseComponent) SetBorderColor(clr color.Color) Component {
	b.element.BorderColor = clr
	return b.self
}

func (b *BaseComponent) SetBorderRadius(radius float32) Component {
	b.element.BorderRadius = BorderRadius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomRight: radius,
		BottomLeft:  radius,
	}
	return b.self
}

func (b *BaseComponent) SetBorderRadius4(topLeft, topRight, bottomRight, bottomLeft float32) Component {
	b.element.BorderRadius = BorderRadius{
		TopLeft:     topLeft,
		TopRight:    topRight,
		BottomRight: bottomRight,
		BottomLeft:  bottomLeft,
	}
	return b.self
}

func (b *BaseComponent) SetShadow(color color.Color, blur, offsetX, offsetY float32) Component {
	b.element.ShadowColor = color
	b.element.ShadowBlur = blur
	b.element.ShadowOffsetX = offsetX
	b.element.ShadowOffsetY = offsetY
	return b.self
}

func (b *BaseComponent) SetVisible(visible bool) Component {
	b.element.Visible = visible
	return b.self
}

func (b *BaseComponent) SetPointerEvents(pe PointerEvents) Component {
	b.element.PointerEvents = pe
	return b.self
}

func (b *BaseComponent) Add(children ...Component) Component {
	for _, child := range children {
		_ = b.AddChild(child)
	}
	return b.self
}

func (b *BaseComponent) SetStyle(style *Style) Component {
	b.ApplyStyle(style)
	return b.self
}

func (b *BaseComponent) SetVisualStyle(style *VisualStyle) Component {
	b.ApplyVisualStyle(style)
	return b.self
}

func (b *BaseComponent) SetKey(key string) Component {
	b.element.Key = key
	return b.self
}

// RemoveChild 移除子组件
func (b *BaseComponent) RemoveChild(child Component) bool {
	for i, c := range b.children {
		if c.ID() == child.ID() {
			// 从 yoga 树中移除
			if childYoga := child.GetElement().Yoga; childYoga != nil && b.yogaNode != nil {
				b.yogaNode.RemoveChild(childYoga)
			}
			b.children = append(b.children[:i], b.children[i+1:]...)
			return true
		}
	}
	return false
}

// ClearChildren 清空所有子组件
func (b *BaseComponent) ClearChildren() {
	for _, child := range b.children {
		if childYoga := child.GetElement().Yoga; childYoga != nil && b.yogaNode != nil {
			b.yogaNode.RemoveChild(childYoga)
		}
	}
	b.children = b.children[:0]
}
