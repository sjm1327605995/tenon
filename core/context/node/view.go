package node

import (
	"image/color"

	"github.com/dhconnelly/rtreego"
	"github.com/millken/yoga"
)

type SearchPoint struct {
	rtreego.Rect
}

type View struct {
	Node
	direction       yoga.Direction
	radius          Radius
	borderColor     color.RGBA
	backgroundColor color.RGBA
}
type Radius struct {
	TopLeft     float32
	BottomLeft  float32
	TopRight    float32
	BottomRight float32
}

func (v *View) OnLayout() {

	yoga.CalculateLayout(v.Yoga(), yoga.Undefined, yoga.Undefined, v.direction)
}

func (v *View) OnDraw(r Renderer, rtree *rtreego.Rtree) {

	//绘制本身
	_ = r.Rectangle(v.Node.X, v.Node.Y, v.Node.Node, v.radius, v.backgroundColor)
	//先画子节点
	v.Node.OnDraw(r, rtree)
}

func NewView() *View {
	return &View{
		Node: Node{
			Node: yoga.NewNode(),
		},
	}
}
func (v *View) SetGap(gutter yoga.Gutter, gapLength float32) *View {
	v.Node.Node.StyleSetGap(gutter, gapLength)
	return v
}
func (v *View) SetPadding(edge yoga.Edge, padding float32) *View {
	v.Node.Node.StyleSetPadding(edge, padding)
	return v
}
func (v *View) SetDirection(direction yoga.Direction) *View {
	v.direction = direction
	v.Node.Node.StyleSetDirection(direction)
	return v
}
func (v *View) SetAlignContent(alignContent yoga.Align) *View {
	v.Node.Node.StyleSetAlignContent(alignContent)
	return v
}
func (v *View) SetAlignItems(alignItems yoga.Align) *View {
	v.Node.Node.StyleSetAlignItems(alignItems)
	return v
}

func (v *View) SetBorder(edge yoga.Edge, border float32) *View {
	v.Node.Node.StyleSetBorder(edge, border)
	return v
}
func (v *View) SetAlignSelf(self yoga.Align) *View {
	v.Node.Node.StyleSetAlignSelf(self)
	return v
}
func (v *View) SetDisplay(display yoga.Display) *View {
	v.Node.Node.StyleSetDisplay(display)
	return v
}
func (v *View) SetFlex(flex float32) *View {
	v.Node.Node.StyleSetFlex(flex)
	return v
}
func (v *View) SetFlexBasis(flexBasis float32) *View {
	v.Node.Node.StyleSetFlexBasis(flexBasis)
	return v
}
func (v *View) SetFlexBasisAuto() *View {
	v.Node.Node.StyleSetFlexBasisAuto()
	return v
}
func (v *View) SetFlexGrow(flexGrow float32) *View {
	v.Node.Node.StyleSetFlexGrow(flexGrow)
	return v
}
func (v *View) SetFlexDirection(flexDirection yoga.FlexDirection) *View {
	v.Node.Node.StyleSetFlexDirection(flexDirection)
	return v
}
func (v *View) SetFlexBasisPercent(flexGrow float32) *View {
	v.Node.Node.StyleSetFlexBasisPercent(flexGrow)
	return v
}
func (v *View) SetFlexShrink(flexShrink float32) *View {
	v.Node.Node.StyleSetFlexShrink(flexShrink)
	return v
}
func (v *View) SetFlexWrap(flexWrap yoga.Wrap) *View {
	v.Node.Node.StyleSetFlexWrap(flexWrap)
	return v
}
func (v *View) SetOverflow(overflow yoga.Overflow) *View {
	v.Node.Node.StyleSetOverflow(overflow)
	return v
}
func (v *View) SetHeightAuto() *View {
	v.Node.Node.StyleSetHeightAuto()
	return v
}
func (v *View) SetHeightPercent(height float32) *View {
	v.Node.Node.StyleSetHeightPercent(height)
	return v
}
func (v *View) SetMargin(edge yoga.Edge, margin float32) *View {
	v.Node.Node.StyleSetMargin(edge, margin)
	return v
}
func (v *View) SetMarginAuto(edge yoga.Edge) *View {
	v.Node.Node.StyleSetMarginAuto(edge)
	return v
}
func (v *View) SetMarginPercent(edge yoga.Edge, margin float32) *View {
	v.Node.Node.StyleSetMarginPercent(edge, margin)
	return v
}
func (v *View) SetMaxHeightPercent(maxHeight float32) *View {
	v.Node.Node.StyleSetMaxHeightPercent(maxHeight)
	return v
}
func (v *View) SetMaxHeight(maxHeight float32) *View {
	v.Node.Node.StyleSetMaxHeight(maxHeight)
	return v
}
func (v *View) SetMinHeightPercent(minHeight float32) *View {
	v.Node.Node.StyleSetMinHeightPercent(minHeight)
	return v
}
func (v *View) SetMinHeight(minHeight float32) *View {
	v.Node.Node.StyleSetMinHeight(minHeight)
	return v
}
func (v *View) SetWidthAuto() *View {
	v.Node.Node.StyleSetWidthAuto()
	return v
}
func (v *View) SetWidthPercent(percent float32) *View {
	v.Node.Node.StyleSetWidthPercent(percent)
	return v
}
func (v *View) SetWidthMaxPercent(percent float32) *View {
	v.Node.Node.StyleSetMaxWidthPercent(percent)
	return v
}
func (v *View) SetWidthMinPercent(percent float32) *View {
	v.Node.Node.StyleSetMinWidthPercent(percent)
	return v
}
func (v *View) SetMaxWidth(maxWidth float32) *View {
	v.Node.Node.StyleSetMaxWidth(maxWidth)
	return v
}
func (v *View) SetMinWidth(minWidth float32) *View {
	v.Node.Node.StyleSetMinWidth(minWidth)
	return v
}
func (v *View) SetJustifyContent(justifyContent yoga.Justify) *View {
	v.Node.Node.StyleSetJustifyContent(justifyContent)
	return v
}

type RadiusType uint8

const (
	RadiusAll RadiusType = iota
	RadiusTopLeft
	RadiusTopRight
	RadiusBottomLeft
	RadiusBottomRight
)

func (v *View) SetRadius(edge RadiusType, radius float32) *View {
	switch edge {
	case RadiusTopLeft:
		v.radius.TopLeft = radius
	case RadiusBottomLeft:
		v.radius.BottomLeft = radius
	case RadiusTopRight:
		v.radius.TopRight = radius
	case RadiusBottomRight:
		v.radius.BottomRight = radius
	default:
		v.radius.TopLeft = radius
		v.radius.TopRight = radius
		v.radius.BottomLeft = radius
		v.radius.BottomRight = radius
	}
	return v
}
func (v *View) SetBackgroundColor(color color.RGBA) *View {
	v.backgroundColor = color
	return v
}
func (v *View) SetWidth(width float32) *View {
	v.Node.Node.StyleSetWidth(width)
	return v
}
func (v *View) SetBorderColor(color color.RGBA) *View {
	v.borderColor = color
	return v
}
func (v *View) SetHeight(height float32) *View {
	v.Node.Node.StyleSetHeight(height)
	return v
}
func (v *View) SetPaddingPercent(edge yoga.Edge, padding float32) *View {
	v.Node.Node.StyleSetPaddingPercent(edge, padding)
	return v
}
func (v *View) SetPosition(edge yoga.Edge, position float32) *View {
	v.Node.Node.StyleSetPosition(edge, position)
	return v
}
func (v *View) SetPositionPercent(edge yoga.Edge, position float32) *View {
	v.Node.Node.StyleSetPositionPercent(edge, position)
	return v
}
func (v *View) SetPositionType(positionType yoga.PositionType) *View {
	v.Node.Node.StyleSetPositionType(positionType)
	return v
}
func (v *View) AddChild(children ...INode) *View {
	v.children = append(v.children, children...)
	for i := range children {
		v.Node.Node.InsertChild(children[i].Yoga(), uint32(i))
	}
	return v
}
func (v *View) SetOnClick(onClick func(v *View)) *View {
	v.Node.Click = func() {
		onClick(v)
	}
	v.IsRtreeNode = 1
	return v
}
