package types

import (
	"github.com/sjm1327605995/tenon/yoga"
)

type ImageStyle struct {
	Width  Value
	Height Value
}

func (s *ImageStyle) GetWidth() Value {
	return s.Width
}

func (s *ImageStyle) GetHeight() Value {
	return s.Height
}

type ImageProps struct {
	Style  *ImageStyle
	Source string
}

func (p *ImageProps) ApplyStyle(node *yoga.Node) {
	if p.Style == nil {
		return
	}
	if p.Style.Width.Unit != UnitAuto {
		applyDimension(node.StyleSetWidth, p.Style.Width)
	}
	if p.Style.Height.Unit != UnitAuto {
		applyDimension(node.StyleSetHeight, p.Style.Height)
	}
}

type ImageElement struct {
	BaseElement
	Props *ImageProps
}

func NewImageElement(props *ImageProps) *ImageElement {
	return &ImageElement{
		Props: props,
	}
}

func (e *ImageElement) GetProps() Props {
	return e.Props
}

func (e *ImageElement) GetChildren() []Element {
	return nil
}

func (e *ImageElement) Type() string {
	return "image"
}