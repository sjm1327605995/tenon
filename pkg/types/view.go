package types

import (
	"github.com/sjm1327605995/tenon/yoga"
)

type ViewStyle struct {
	Width          Value
	Height         Value
	Background     string
	Padding        Value
	Margin         Value
	FlexDirection  yoga.FlexDirection
	JustifyContent yoga.Justify
	AlignItems     yoga.Align
	FlexGrow       float32
	MarginBottom   Value
	MarginRight    Value
}

func (s *ViewStyle) GetFlexDirection() yoga.FlexDirection {
	return s.FlexDirection
}

func (s *ViewStyle) GetJustifyContent() yoga.Justify {
	return s.JustifyContent
}

func (s *ViewStyle) GetAlignItems() yoga.Align {
	return s.AlignItems
}

func (s *ViewStyle) GetFlexGrow() float32 {
	return s.FlexGrow
}

func (s *ViewStyle) GetMargin() Value {
	return s.Margin
}

func (s *ViewStyle) GetPadding() Value {
	return s.Padding
}

func (s *ViewStyle) GetMarginBottom() Value {
	return s.MarginBottom
}

func (s *ViewStyle) GetMarginRight() Value {
	return s.MarginRight
}

func (s *ViewStyle) GetWidth() Value {
	return s.Width
}

func (s *ViewStyle) GetHeight() Value {
	return s.Height
}

func (s *ViewStyle) GetBackground() string {
	return s.Background
}

type ViewProps struct {
	Style   *ViewStyle
	OnClick func()
}

func (p *ViewProps) ApplyStyle(node *yoga.Node) {
	if p.Style == nil {
		return
	}
	s := p.Style
	if s.Width.Unit != UnitAuto {
		applyDimension(node.StyleSetWidth, s.Width)
	}
	if s.Height.Unit != UnitAuto {
		applyDimension(node.StyleSetHeight, s.Height)
	}
	if s.FlexDirection != 0 {
		node.StyleSetFlexDirection(s.FlexDirection)
	}
	if s.JustifyContent != 0 {
		node.StyleSetJustifyContent(s.JustifyContent)
	}
	if s.AlignItems != 0 {
		node.StyleSetAlignItems(s.AlignItems)
	}
	if s.FlexGrow > 0 {
		node.StyleSetFlexGrow(s.FlexGrow)
	}
	if s.MarginBottom.Unit != UnitAuto {
		applyDimension(func(v float32) { node.StyleSetMargin(yoga.EdgeBottom, v) }, s.MarginBottom)
	}
	if s.MarginRight.Unit != UnitAuto {
		applyDimension(func(v float32) { node.StyleSetMargin(yoga.EdgeRight, v) }, s.MarginRight)
	}
}

type ViewElement struct {
	BaseElement
	Props *ViewProps
}

func NewViewElement(props *ViewProps, children ...Element) *ViewElement {
	e := &ViewElement{
		Props: props,
	}
	e.children = children
	return e
}

func (e *ViewElement) GetProps() Props {
	return e.Props
}

func (e *ViewElement) GetChildren() []Element {
	return e.children
}

func (e *ViewElement) Type() string {
	return "view"
}