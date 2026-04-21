package types

import (
	"github.com/sjm1327605995/tenon/yoga"
)

type TextStyle struct {
	FontSize Value
	Color    string
}

func (s *TextStyle) GetFontSize() Value {
	return s.FontSize
}

func (s *TextStyle) GetColor() string {
	return s.Color
}

type TextProps struct {
	Style   *TextStyle
	Content string
}

func (p *TextProps) ApplyStyle(node *yoga.Node) {
	if p.Style == nil {
		return
	}
	if p.Style.FontSize.Unit != UnitAuto {
		content := p.Content
		if content == "" {
			content = "Text"
		}
		node.StyleSetWidth(float32(len(content)) * p.Style.FontSize.Value * 0.6)
		node.StyleSetHeight(p.Style.FontSize.Value * 1.2)
	}
}

type TextElement struct {
	BaseElement
	Props *TextProps
}

func NewTextElement(props *TextProps) *TextElement {
	return &TextElement{
		Props: props,
	}
}

func (e *TextElement) GetProps() Props {
	return e.Props
}

func (e *TextElement) GetChildren() []Element {
	return nil
}

func (e *TextElement) Type() string {
	return "text"
}

func TextMeasureFunc(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
	var content string
	if props, ok := node.GetContext().(*TextProps); ok {
		content = props.Content
	} else {
		content = "Text"
	}

	fontSize := float32(16)
	if props, ok := node.GetContext().(*TextProps); ok {
		if props.Style != nil && props.Style.FontSize.Unit != UnitAuto {
			fontSize = props.Style.FontSize.Value
		}
	}

	contentWidth := float32(len(content)) * fontSize * 0.6
	contentHeight := fontSize * 1.2

	return yoga.Size{
		Width:  contentWidth,
		Height: contentHeight,
	}
}