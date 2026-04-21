package core

import (
	"image/color"

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
	BorderRadius    BorderRadius // CSS 样式四个角的圆角半径
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

type Component interface {
	Draw(screen *ebiten.Image)
	Update() error
	DrawOverlay(screen *ebiten.Image)
	HandleInput() bool
	GetChildren() []Component
	AddChild(child Component) error
	GetLayoutBounds() LayoutBounds
	GetElement() *Element
	Render() *Element
	UpdateLayoutBounds(parentX, parentY float32)
}

type BaseComponent struct {
	children []Component
	element  *Element
	yogaNode *yoga.Node
	bounds   LayoutBounds
}

func NewBaseComponent() *BaseComponent {
	node := yoga.NewNode()
	return &BaseComponent{
		children: make([]Component, 0),
		element: &Element{
			Visible:       true,
			Yoga:          node,
			PointerEvents: PointerEventsAuto,
		},
		yogaNode: node,
	}
}

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

	// 获取相对于父节点的位置
	relX := b.yogaNode.LayoutLeft()
	relY := b.yogaNode.LayoutTop()

	// 转换为绝对位置
	abX := parentX + relX
	abY := parentY + relY

	b.bounds = LayoutBounds{
		X:      abX,
		Y:      abY,
		Width:  b.yogaNode.LayoutWidth(),
		Height: b.yogaNode.LayoutHeight(),
	}
}

func CalculateLayout(root Component, width, height float32) {
	element := root.GetElement()
	if element == nil || element.Yoga == nil {
		return
	}

	// 设置根节点的宽度和高度
	element.Yoga.StyleSetWidth(width)
	element.Yoga.StyleSetHeight(height)

	// 计算布局
	element.Yoga.CalculateLayout(width, height, yoga.DirectionLTR)
	updateAllBounds(root, 0, 0) // 从(0,0)开始计算绝对位置
}

func updateAllBounds(component Component, parentX, parentY float32) {
	if component == nil {
		return
	}

	// 更新当前组件的绝对位置
	component.UpdateLayoutBounds(parentX, parentY)
	bounds := component.GetLayoutBounds()

	// 递归更新子组件的绝对位置
	for _, child := range component.GetChildren() {
		updateAllBounds(child, bounds.X, bounds.Y)
	}
}
