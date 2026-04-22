package core

import (
	"image/color"

	"github.com/sjm1327605995/tenon/yoga"
)

type Style struct {
	rules []func(*yoga.Node)
}

func NewStyle() *Style {
	return &Style{
		rules: make([]func(*yoga.Node), 0),
	}
}

func (s *Style) Apply(node *yoga.Node) {
	for _, rule := range s.rules {
		rule(node)
	}
}

func (s *Style) Merge(other *Style) *Style {
	if other == nil {
		return s
	}
	s.rules = append(s.rules, other.rules...)
	return s
}

func (s *Style) Clone() *Style {
	newStyle := &Style{
		rules: make([]func(*yoga.Node), len(s.rules)),
	}
	copy(newStyle.rules, s.rules)
	return newStyle
}

func (s *Style) SetWidth(width float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetWidth(width)
	})
	return s
}

func (s *Style) SetHeight(height float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetHeight(height)
	})
	return s
}

func (s *Style) SetMinWidth(width float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetMinWidth(width)
	})
	return s
}

func (s *Style) SetMinHeight(height float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetMinHeight(height)
	})
	return s
}

func (s *Style) SetMaxWidth(width float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetMaxWidth(width)
	})
	return s
}

func (s *Style) SetMaxHeight(height float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetMaxHeight(height)
	})
	return s
}

func (s *Style) SetFlexDirection(dir yoga.FlexDirection) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetFlexDirection(dir)
	})
	return s
}

func (s *Style) SetJustifyContent(justify yoga.Justify) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetJustifyContent(justify)
	})
	return s
}

func (s *Style) SetAlignItems(align yoga.Align) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetAlignItems(align)
	})
	return s
}

func (s *Style) SetAlignSelf(align yoga.Align) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetAlignSelf(align)
	})
	return s
}

func (s *Style) SetFlexGrow(grow float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetFlexGrow(grow)
	})
	return s
}

func (s *Style) SetFlexShrink(shrink float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetFlexShrink(shrink)
	})
	return s
}

func (s *Style) SetFlexBasis(basis float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetFlexBasis(basis)
	})
	return s
}

func (s *Style) SetFlexWrap(wrap yoga.Wrap) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetFlexWrap(wrap)
	})
	return s
}

func (s *Style) SetPadding(edge yoga.Edge, value float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetPadding(edge, value)
	})
	return s
}

func (s *Style) SetMargin(edge yoga.Edge, value float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetMargin(edge, value)
	})
	return s
}

func (s *Style) SetBorder(edge yoga.Edge, value float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetBorder(edge, value)
	})
	return s
}

func (s *Style) SetPosition(edge yoga.Edge, value float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetPosition(edge, value)
	})
	return s
}

func (s *Style) SetPositionType(positionType yoga.PositionType) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetPositionType(positionType)
	})
	return s
}

func (s *Style) SetAspectRatio(ratio float32) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetAspectRatio(ratio)
	})
	return s
}

func (s *Style) SetDisplay(display yoga.Display) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetDisplay(display)
	})
	return s
}

func (s *Style) SetOverflow(overflow yoga.Overflow) *Style {
	s.rules = append(s.rules, func(node *yoga.Node) {
		node.StyleSetOverflow(overflow)
	})
	return s
}

type VisualStyle struct {
	rules []func(*Element)
}

func NewVisualStyle() *VisualStyle {
	return &VisualStyle{
		rules: make([]func(*Element), 0),
	}
}

func (vs *VisualStyle) Apply(element *Element) {
	for _, rule := range vs.rules {
		rule(element)
	}
}

func (vs *VisualStyle) Merge(other *VisualStyle) *VisualStyle {
	if other == nil {
		return vs
	}
	vs.rules = append(vs.rules, other.rules...)
	return vs
}

func (vs *VisualStyle) Clone() *VisualStyle {
	newStyle := &VisualStyle{
		rules: make([]func(*Element), len(vs.rules)),
	}
	copy(newStyle.rules, vs.rules)
	return newStyle
}

func (vs *VisualStyle) SetBackgroundColor(clr color.Color) *VisualStyle {
	vs.rules = append(vs.rules, func(element *Element) {
		element.BackgroundColor = clr
	})
	return vs
}

func (vs *VisualStyle) SetBorderColor(clr color.Color) *VisualStyle {
	vs.rules = append(vs.rules, func(element *Element) {
		element.BorderColor = clr
	})
	return vs
}

func (vs *VisualStyle) SetBorderRadius(radius float32) *VisualStyle {
	vs.rules = append(vs.rules, func(element *Element) {
		element.BorderRadius = BorderRadius{
			TopLeft:     radius,
			TopRight:    radius,
			BottomRight: radius,
			BottomLeft:  radius,
		}
	})
	return vs
}

func (vs *VisualStyle) SetBorderRadius4(topLeft, topRight, bottomRight, bottomLeft float32) *VisualStyle {
	vs.rules = append(vs.rules, func(element *Element) {
		element.BorderRadius = BorderRadius{
			TopLeft:     topLeft,
			TopRight:    topRight,
			BottomRight: bottomRight,
			BottomLeft:  bottomLeft,
		}
	})
	return vs
}

func (vs *VisualStyle) SetShadow(color color.Color, blur, offsetX, offsetY float32) *VisualStyle {
	vs.rules = append(vs.rules, func(element *Element) {
		element.ShadowColor = color
		element.ShadowBlur = blur
		element.ShadowOffsetX = offsetX
		element.ShadowOffsetY = offsetY
	})
	return vs
}

func (vs *VisualStyle) SetVisible(visible bool) *VisualStyle {
	vs.rules = append(vs.rules, func(element *Element) {
		element.Visible = visible
	})
	return vs
}

func (vs *VisualStyle) SetPointerEvents(pe PointerEvents) *VisualStyle {
	vs.rules = append(vs.rules, func(element *Element) {
		element.PointerEvents = pe
	})
	return vs
}
