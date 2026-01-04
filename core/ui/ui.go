package ui

import (
	"image/color"

	"gioui.org/unit"
	"github.com/sjm1327605995/tenon/core/ui/render"
	"github.com/sjm1327605995/tenon/yoga"
)

var Metric unit.Metric

type ElementOption func(element *Element)
type BaseUI[T any] struct {
	This      *T
	PropsFunc []ElementOption
	Clickable bool
}

type Unit uint8

const (
	PxType Unit = 1 + iota
	PercentType
	AutoType
)

type UI interface {
	Render() *Element
}

type Value struct {
	val  float32
	unit Unit
}

func Auto() Value {
	return Value{val: 0.0, unit: AutoType}
}
func Px(val float32) Value {
	return Value{val: val, unit: PxType}
}
func Percent(val float32) Value {
	return Value{val: val, unit: PercentType}
}
func (b *BaseUI[T]) Direction(direction yoga.Direction) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetDirection(direction)
	})
	return b.This
}
func (b *BaseUI[T]) FlexDirection(direction yoga.FlexDirection) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetFlexDirection(direction)
	})
	return b.This
}
func (b *BaseUI[T]) JustifyContent(justifyContent yoga.Justify) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetJustifyContent(justifyContent)
	})
	return b.This
}
func (b *BaseUI[T]) AlignContent(alignContent yoga.Align) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetAlignContent(alignContent)
	})
	return b.This
}
func (b *BaseUI[T]) AlignItems(alignItems yoga.Align) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetAlignItems(alignItems)
	})
	return b.This
}
func (b *BaseUI[T]) AlignSelf(alignSelf yoga.Align) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetAlignItems(alignSelf)
	})
	return b.This
}
func (b *BaseUI[T]) FlexWrap(flexWrap yoga.Wrap) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetFlexWrap(flexWrap)
	})
	return b.This
}
func (b *BaseUI[T]) Overflow(overflow yoga.Overflow) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetOverflow(overflow)
	})
	return b.This
}

func (b *BaseUI[T]) Display(overflow yoga.Display) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetDisplay(overflow)
	})
	return b.This
}
func (b *BaseUI[T]) Flex(flex float32) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetFlex(flex)
	})
	return b.This
}
func (b *BaseUI[T]) FlexGrow(flexGrow float32) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetFlexGrow(flexGrow)
	})
	return b.This
}

func (b *BaseUI[T]) FlexShrink(flexShrink float32) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetFlexShrink(flexShrink)
	})
	return b.This
}
func (b *BaseUI[T]) FlexBasis(flexBasis Value) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		switch flexBasis.unit {
		case PxType:
			element.Yoga.StyleSetFlexBasis(flexBasis.val)
		case PercentType:
			element.Yoga.StyleSetFlexBasisPercent(flexBasis.val)
		case AutoType:
			element.Yoga.StyleSetFlexBasisAuto()
		}
	})
	return b.This
}
func (b *BaseUI[T]) Width(width Value) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {

		switch width.unit {
		case PxType:
			element.Yoga.StyleSetWidth(width.val)
		case PercentType:
			element.Yoga.StyleSetWidthPercent(width.val)
		case AutoType:
			element.Yoga.StyleSetWidthAuto()
		}

	})
	return b.This
}
func (b *BaseUI[T]) Height(height Value) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		switch height.unit {
		case PxType:
			element.Yoga.StyleSetHeight(height.val)
		case PercentType:
			element.Yoga.StyleSetHeightPercent(height.val)
		case AutoType:
			element.Yoga.StyleSetHeightAuto()
		}
	})
	return b.This
}
func (b *BaseUI[T]) PositionType(positionType yoga.PositionType) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetPositionType(positionType)
	})
	return b.This
}

func (b *BaseUI[T]) Position(edge yoga.Edge, position Value) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		switch position.unit {
		case PxType:
			element.Yoga.StyleSetPosition(edge, position.val)
		case PercentType:
			element.Yoga.StyleSetPositionPercent(edge, position.val)
		default:
		}
	})
	return b.This
}

func (b *BaseUI[T]) Margin(edge yoga.Edge, val float32) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetMargin(edge, val)
	})
	return b.This
}
func (b *BaseUI[T]) MarginPercent(edge yoga.Edge, margin float32) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetMarginPercent(edge, margin)
	})
	return b.This
}
func (b *BaseUI[T]) MarginAuto(edge yoga.Edge) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetMarginAuto(edge)
	})
	return b.This
}

func (b *BaseUI[T]) Padding(edge yoga.Edge, padding float32) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetPadding(edge, padding)
	})
	return b.This
}
func (b *BaseUI[T]) PaddingPercent(edge yoga.Edge, padding float32) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetPaddingPercent(edge, padding)
	})
	return b.This
}

func (b *BaseUI[T]) Gap(gutter yoga.Gutter, gapLength float32) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetGap(gutter, gapLength)
	})
	return b.This
}

func (b *BaseUI[T]) MinWidth(value Value) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		switch value.unit {
		case PxType:
			element.Yoga.StyleSetMinWidth(value.val)
		case PercentType:
			element.Yoga.StyleSetMinWidthPercent(value.val)
		default:
		}
	})
	return b.This
}
func (b *BaseUI[T]) MinHeight(value Value) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		switch value.unit {
		case PxType:
			element.Yoga.StyleSetMinHeight(value.val)
		case PercentType:
			element.Yoga.StyleSetMinHeightPercent(value.val)
		default:
		}
	})
	return b.This
}
func (b *BaseUI[T]) MaxWidth(value Value) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		switch value.unit {
		case PxType:
			element.Yoga.StyleSetMaxWidth(value.val)
		case PercentType:
			element.Yoga.StyleSetMaxWidthPercent(value.val)
		default:
		}
	})
	return b.This
}
func (b *BaseUI[T]) MaxHeight(value Value) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		switch value.unit {
		case PxType:
			element.Yoga.StyleSetMaxHeight(value.val)
		case PercentType:
			element.Yoga.StyleSetMaxHeightPercent(value.val)
		default:
		}
	})
	return b.This
}
func (b *BaseUI[T]) AspectRatio(aspectRatio float32) *T {
	b.PropsFunc = append(b.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetAspectRatio(aspectRatio)
	})
	return b.This
}

func (b *BaseUI[T]) Render() *Element {
	element := CreateElement(nil)
	return element
}
func NewBaseUI[T any](this *T) *BaseUI[T] {
	return &BaseUI[T]{This: this}
}

type ViewUI struct {
	*BaseUI[ViewUI]
	onClick  func()
	Children []UI
	style    render.ViewStyle
}

func View(children ...UI) *ViewUI {
	v := &ViewUI{
		style: render.ViewStyle{BorderColor: color.NRGBA{A: 255}},
	}
	v.BaseUI = NewBaseUI[ViewUI](v)
	v.Children = children
	return v
}

func (v *ViewUI) Background(nrgba color.NRGBA) *ViewUI {
	v.style.Background = nrgba
	return v
}
func (v *ViewUI) OnClick(onClick func()) *ViewUI {
	v.onClick = onClick
	return v
}
func (v *ViewUI) Border(edge yoga.Edge, border float32) *ViewUI {
	v.PropsFunc = append(v.PropsFunc, func(element *Element) {
		element.Yoga.StyleSetBorder(edge, border)
	})
	switch edge {
	case yoga.EdgeBottom:
		v.style.Bottom = border
	case yoga.EdgeTop:
		v.style.Top = border
	case yoga.EdgeRight:
		v.style.Right = border
	case yoga.EdgeLeft:
		v.style.Left = border
	default:
		v.style.Bottom = border
		v.style.Top = border
		v.style.Right = border
		v.style.Left = border
	}
	return v
}
func (v *ViewUI) BorderRadius(radius ...float32) *ViewUI {
	switch len(radius) {
	case 4:
		v.style.CornerRadii.TopLeft = radius[0]
		v.style.CornerRadii.TopRight = radius[1]
		v.style.CornerRadii.BottomLeft = radius[2]
		v.style.CornerRadii.BottomRight = radius[3]
	case 3:
		v.style.CornerRadii.TopLeft = radius[0]
		v.style.CornerRadii.TopRight = radius[1]
		v.style.CornerRadii.BottomLeft = radius[1]
		v.style.CornerRadii.BottomRight = radius[2]
	case 2:
		v.style.CornerRadii.TopLeft = radius[0]
		v.style.CornerRadii.BottomRight = radius[0]
		v.style.CornerRadii.TopRight = radius[1]
		v.style.CornerRadii.BottomLeft = radius[1]
	case 1:
		v.style.CornerRadii.TopLeft = radius[0]
		v.style.CornerRadii.TopRight = radius[0]
		v.style.CornerRadii.BottomRight = radius[0]
		v.style.CornerRadii.BottomLeft = radius[0]
	default:

	}
	return v
}
func (v *ViewUI) Render() *Element {
	element := CreateElement(v.style)
	for i := range v.PropsFunc {
		v.PropsFunc[i](element)
	}
	for _, child := range v.Children {
		element.InsertChild(child.Render())
	}
	return element
}
