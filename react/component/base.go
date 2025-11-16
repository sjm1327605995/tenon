package component

import (
	"image"
	"image/color"

	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/yoga"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Base[T any] struct {
	Node            *yoga.Node
	Idx             uint32
	This            *T
	BackgroundColor color.NRGBA
}

func (b *Base[T]) Yoga() *yoga.Node {
	return b.Node
}

func (b *Base[T]) Layout(gtx layout.Context) layout.Dimensions {
	w := b.Node.StyleGetWidth()
	h := b.Node.StyleGetHeight()
	size := image.Pt(int(w), int(h))
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: b.BackgroundColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: image.Pt(int(w), int(h))}
}

func NewBase[T any](this *T) Base[T] {
	return Base[T]{Node: yoga.NewNode(), This: this}
}

func (b *Base[T]) Body(children ...core.Node) *T {
	for i := range children {
		children[i].Yoga().SetContext(children[i])
		b.Node.InsertChild(children[i].Yoga(), b.Idx)
		b.Idx += uint32(i)
	}
	return b.This
}
func (b *Base[T]) Direction(direction yoga.Direction) *T {
	b.Node.StyleSetDirection(direction)
	return b.This
}
func (b *Base[T]) FlexDirection(flexDirection yoga.FlexDirection) *T {
	b.Node.StyleSetFlexDirection(flexDirection)
	return b.This
}

func (b *Base[T]) JustifyContent(justifyContent yoga.Justify) *T {
	b.Node.StyleSetJustifyContent(justifyContent)
	return b.This
}

func (b *Base[T]) AlignContent(alignContent yoga.Align) *T {
	b.Node.StyleSetAlignContent(alignContent)
	return b.This
}

func (b *Base[T]) AlignItems(alignItems yoga.Align) *T {
	b.Node.StyleSetAlignItems(alignItems)
	return b.This
}

func (b *Base[T]) AlignSelf(alignSelf yoga.Align) *T {
	b.Node.StyleSetAlignSelf(alignSelf)
	return b.This
}

func (b *Base[T]) FlexWrap(flexWrap yoga.Wrap) *T {
	b.Node.StyleSetFlexWrap(flexWrap)
	return b.This
}

func (b *Base[T]) Overflow(overflow yoga.Overflow) *T {
	b.Node.StyleSetOverflow(overflow)
	return b.This
}

func (b *Base[T]) Display(display yoga.Display) *T {
	b.Node.StyleSetDisplay(display)
	return b.This
}

func (b *Base[T]) Flex(flex float32) *T {
	b.Node.StyleSetFlex(flex)
	return b.This
}

func (b *Base[T]) FlexGrow(flexGrow float32) *T {
	b.Node.StyleSetFlexGrow(flexGrow)
	return b.This
}

func (b *Base[T]) FlexShrink(flexShrink float32) *T {
	b.Node.StyleSetFlexShrink(flexShrink)
	return b.This
}

func (b *Base[T]) FlexBasis(flexBasis float32) *T {
	b.Node.StyleSetFlexBasis(flexBasis)
	return b.This
}

func (b *Base[T]) FlexBasisPercent(flexBasis float32) *T {
	b.Node.StyleSetFlexBasisPercent(flexBasis)
	return b.This
}

func (b *Base[T]) FlexBasisAuto() *T {
	b.Node.StyleSetFlexBasisAuto()
	return b.This
}

func (b *Base[T]) Width(points float32) *T {
	b.Node.StyleSetWidth(points)
	return b.This
}

func (b *Base[T]) WidthPercent(percent float32) *T {
	b.Node.StyleSetWidthPercent(percent)
	return b.This
}

func (b *Base[T]) WidthAuto() *T {
	b.Node.StyleSetWidthAuto()
	return b.This
}

func (b *Base[T]) Height(height float32) *T {
	b.Node.StyleSetHeight(height)
	return b.This
}

func (b *Base[T]) HeightPercent(height float32) *T {
	b.Node.StyleSetHeightPercent(height)
	return b.This
}

func (b *Base[T]) HeightAuto() *T {
	b.Node.StyleSetHeightAuto()
	return b.This
}

func (b *Base[T]) PositionType(positionType yoga.PositionType) *T {
	b.Node.StyleSetPositionType(positionType)
	return b.This
}

func (b *Base[T]) StyleSetPosition(edge yoga.Edge, position float32) *T {
	b.Node.StyleSetPosition(edge, position)
	return b.This
}

func (b *Base[T]) PositionPercent(edge yoga.Edge, position float32) *T {
	b.Node.StyleSetPositionPercent(edge, position)
	return b.This
}

func (b *Base[T]) Margin(edge yoga.Edge, margin float32) *T {
	b.Node.StyleSetMargin(edge, margin)
	return b.This
}

func (b *Base[T]) MarginPercent(edge yoga.Edge, margin float32) *T {
	b.Node.StyleSetMarginPercent(edge, margin)
	return b.This
}

func (b *Base[T]) MarginAuto(edge yoga.Edge) *T {
	b.Node.StyleSetMarginAuto(edge)
	return b.This
}

func (b *Base[T]) Padding(edge yoga.Edge, padding float32) *T {
	b.Node.StyleSetPadding(edge, padding)
	return b.This
}

func (b *Base[T]) PaddingPercent(edge yoga.Edge, padding float32) *T {
	b.Node.StyleSetPaddingPercent(edge, padding)
	return b.This
}

func (b *Base[T]) BorderAll(border float32) *T {
	b.Node.StyleSetBorder(yoga.EdgeAll, border)
	return b.This
}

//func (b *Base[T]) Border(edge yoga.Edge, border float32) *T {
//	b.Node.StyleSetBorder(edge, border)
//	return b.This
//}

func (b *Base[T]) Gap(gutter yoga.Gutter, gapLength float32) *T {
	b.Node.StyleSetGap(gutter, gapLength)
	return b.This
}

func (b *Base[T]) MinWidth(minWidth float32) *T {
	b.Node.StyleSetMinWidth(minWidth)
	return b.This
}

func (b *Base[T]) MinWidthPercent(minWidth float32) *T {
	b.Node.StyleSetMinWidthPercent(minWidth)
	return b.This
}

func (b *Base[T]) MinHeight(minHeight float32) *T {
	b.Node.StyleSetMinHeight(minHeight)
	return b.This
}

func (b *Base[T]) MinHeightPercent(minHeight float32) *T {
	b.Node.StyleSetMinHeightPercent(minHeight)
	return b.This
}

func (b *Base[T]) MaxWidth(maxWidth float32) *T {
	b.Node.StyleSetMaxWidth(maxWidth)
	return b.This
}

func (b *Base[T]) MaxWidthPercent(maxWidth float32) *T {
	b.Node.StyleSetMaxWidthPercent(maxWidth)
	return b.This
}

func (b *Base[T]) MaxHeigh(maxHeight float32) *T {
	b.Node.StyleSetMaxHeight(maxHeight)
	return b.This
}

func (b *Base[T]) MaxHeightPercent(maxHeight float32) *T {
	b.Node.StyleSetMaxHeightPercent(maxHeight)
	return b.This
}

func (b *Base[T]) AspectRatio(aspectRatio float32) *T {
	b.Node.StyleSetAspectRatio(aspectRatio)
	return b.This
}

func (b *Base[T]) Background(color color.NRGBA) *T {
	b.BackgroundColor = color
	return b.This
}
