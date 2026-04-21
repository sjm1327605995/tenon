package types

import (
	"github.com/sjm1327605995/tenon/yoga"
)

type UI interface {
	Render() Element
}

type Props interface {
	ApplyStyle(node *yoga.Node)
}

type Element interface {
	GetProps() Props
	GetChildren() []Element
	GetNode() *yoga.Node
	SetNode(node *yoga.Node)
	GetLayout() LayoutRect
	SetLayout(rect LayoutRect)
	Type() string
}

type BaseElement struct {
	props    Props
	children []Element
	node     *yoga.Node
	layout   LayoutRect
}

func (e *BaseElement) GetProps() Props {
	return e.props
}

func (e *BaseElement) GetChildren() []Element {
	return e.children
}

func (e *BaseElement) GetNode() *yoga.Node {
	return e.node
}

func (e *BaseElement) SetNode(node *yoga.Node) {
	e.node = node
}

func (e *BaseElement) GetLayout() LayoutRect {
	return e.layout
}

func (e *BaseElement) SetLayout(rect LayoutRect) {
	e.layout = rect
}

type LayoutRect struct {
	X, Y, Width, Height float32
}

type Value struct {
	Value float32
	Unit  Unit
}

type Unit int

const (
	UnitPx Unit = iota
	UnitPercent
	UnitAuto
)

func Px(value float32) Value {
	return Value{Value: value, Unit: UnitPx}
}

func Percent(value float32) Value {
	return Value{Value: value, Unit: UnitPercent}
}

func Auto() Value {
	return Value{Unit: UnitAuto}
}

func applyDimension(setter func(float32), value Value) {
	switch value.Unit {
	case UnitPx:
		setter(value.Value)
	case UnitPercent:
		setter(value.Value)
	}
}

type FlexStyle interface {
	GetFlexDirection() yoga.FlexDirection
	GetJustifyContent() yoga.Justify
	GetAlignItems() yoga.Align
	GetFlexGrow() float32
}

type SpacingStyle interface {
	GetMargin() Value
	GetPadding() Value
	GetMarginBottom() Value
	GetMarginRight() Value
}

type DimensionStyle interface {
	GetWidth() Value
	GetHeight() Value
	GetBackground() string
}

type TextStyleInterface interface {
	GetFontSize() Value
	GetColor() string
}

type ImageStyleInterface interface {
	GetWidth() Value
	GetHeight() Value
}