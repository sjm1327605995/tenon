// Package styles provides the style system implementation for the React framework.
// This package contains all style-related interfaces, types, and methods, supporting Flexbox layout and custom style properties.
package styles

import (
	"image/color"

	"github.com/sjm1327605995/tenon/react/yoga"
)

// IExtendedStyle is a marker interface for extended styles.
type IExtendedStyle interface {
	ExtendedStyle()
}

// StyleElement is the interface for elements that can have styles applied to them.
type StyleElement interface {
	GetYogaNode() *yoga.Node
	SetExtendedStyle(style IExtendedStyle)
}

// Style represents a collection of styles for an element.
type Style struct {
	handleChains []func(element StyleElement)
}

// NewStyle creates and returns a new Style instance.
func NewStyle() *Style {
	return &Style{}
}

// Apply applies the style to the specified StyleElement.
func (s *Style) Apply(element StyleElement) {
	for _, f := range s.handleChains {
		f(element)
	}
}

// --- Yoga Style Methods ---

func (s *Style) Direction(direction yoga.Direction) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetDirection(direction)
	})
	return s
}

func (s *Style) FlexDirection(flexDirection yoga.FlexDirection) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetFlexDirection(flexDirection)
	})
	return s
}

func (s *Style) JustifyContent(justifyContent yoga.Justify) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetJustifyContent(justifyContent)
	})
	return s
}

func (s *Style) AlignContent(alignContent yoga.Align) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetAlignContent(alignContent)
	})
	return s
}

func (s *Style) AlignItem(alignItems yoga.Align) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetAlignItems(alignItems)
	})
	return s
}

func (s *Style) AlignSelf(alignSelf yoga.Align) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetAlignSelf(alignSelf)
	})
	return s
}

func (s *Style) FlexWrap(flexWrap yoga.Wrap) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetFlexWrap(flexWrap)
	})
	return s
}

func (s *Style) Overflow(overflow yoga.Overflow) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetOverflow(overflow)
	})
	return s
}

func (s *Style) Display(display yoga.Display) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetDisplay(display)
	})
	return s
}

func (s *Style) Flex(flex float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetFlex(flex)
	})
	return s
}

func (s *Style) FlexGrow(flexGrow float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetFlexGrow(flexGrow)
	})
	return s
}

func (s *Style) FlexShrink(flexShrink float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetFlexShrink(flexShrink)
	})
	return s
}

func (s *Style) FlexBasis(flexBasis float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetFlexBasis(flexBasis)
	})
	return s
}

func (s *Style) FlexBasisPercent(flexBasis float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetFlexBasisPercent(flexBasis)
	})
	return s
}

func (s *Style) FlexBasisAuto() *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetFlexBasisAuto()
	})
	return s
}

func (s *Style) Width(points float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetWidth(points)
	})
	return s
}

func (s *Style) WidthPercent(percent float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetWidthPercent(percent)
	})
	return s
}

func (s *Style) WidthAuto() *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetWidthAuto()
	})
	return s
}

func (s *Style) Height(height float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetHeight(height)
	})
	return s
}

func (s *Style) HeightPercent(height float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetHeightPercent(height)
	})
	return s
}

func (s *Style) HeightAuto() *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetHeightAuto()
	})
	return s
}

func (s *Style) PositionType(positionType yoga.PositionType) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetPositionType(positionType)
	})
	return s
}

func (s *Style) Position(edge yoga.Edge, position float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetPosition(edge, position)
	})
	return s
}

func (s *Style) PositionPercent(edge yoga.Edge, position float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetPositionPercent(edge, position)
	})
	return s
}

func (s *Style) Margin(edge yoga.Edge, margin float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMargin(edge, margin)
	})
	return s
}

func (s *Style) MarginPercent(edge yoga.Edge, margin float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMarginPercent(edge, margin)
	})
	return s
}

func (s *Style) MarginAuto(edge yoga.Edge) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMarginAuto(edge)
	})
	return s
}

func (s *Style) Padding(edge yoga.Edge, padding float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetPadding(edge, padding)
	})
	return s
}

func (s *Style) PaddingPercent(edge yoga.Edge, padding float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetPaddingPercent(edge, padding)
	})
	return s
}

func (s *Style) Border(edge yoga.Edge, border float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetBorder(edge, border)
	})
	return s
}

func (s *Style) Gap(edge yoga.Gutter, gapLength float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetGap(edge, gapLength)
	})
	return s
}

func (s *Style) MinWidth(minWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMinWidth(minWidth)
	})
	return s
}

func (s *Style) MinWidthPercent(minWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMinWidthPercent(minWidth)
	})
	return s
}
func (s *Style) MinHeight(minHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMinHeight(minHeight)
	})
	return s
}

func (s *Style) MinHeightPercent(minHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMinHeightPercent(minHeight)
	})
	return s
}

func (s *Style) MaxWidth(maxWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMaxWidth(maxWidth)
	})
	return s
}

func (s *Style) MaxWidthPercent(maxWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMaxWidthPercent(maxWidth)
	})
	return s
}

func (s *Style) MaxHeight(maxHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMaxHeight(maxHeight)
	})
	return s
}

func (s *Style) MaxHeightPercent(maxHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetMaxHeightPercent(maxHeight)
	})
	return s
}

func (s *Style) AspectRatio(aspectRatio float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.GetYogaNode().StyleSetAspectRatio(aspectRatio)
	})
	return s
}

// --- Extended Style Methods ---

func (s *Style) BackgroundColor(backgroundColor color.NRGBA) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.SetExtendedStyle(BackgroundColor{
			Color: backgroundColor,
		})
	})
	return s
}

func (s *Style) BorderColor(borderColor color.NRGBA) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.SetExtendedStyle(BorderColor{
			Color: borderColor,
		})
	})
	return s
}

func (s *Style) CornerRadius(radius CornerRadius) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.SetExtendedStyle(radius)
	})
	return s
}
