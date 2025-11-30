// Package styles provides the style system implementation for the React framework.
// This package contains all style-related interfaces, types, and methods, supporting Flexbox layout and custom style properties.
package styles

import (
	"image/color"

	"github.com/sjm1327605995/tenon/react/yoga"
)

// IExtendedStyle is a marker interface for extended styles.
// Any custom style type should implement this interface to be accepted by StyleElement.
type IExtendedStyle interface {
	// ExtendedStyle is the marker method for the interface, used to identify extended style types.
	ExtendedStyle()
}

// StyleElement is the interface for elements that can have styles applied to them.
// Elements implementing this interface can accept and apply styles defined by Style objects.
type StyleElement interface {
	// Yoga returns the Yoga layout node associated with the element.
	// It is used to access and modify Flexbox layout properties.
	Yoga() *yoga.Node
	// SetExtendedStyle sets the extended style for the element.
	// The style parameter is a custom style object implementing the IExtendedStyle interface.
	SetExtendedStyle(style IExtendedStyle)
}

// Style represents a collection of styles for an element, supporting chained calls to set multiple style properties.
// It uses function chains to delay style application for improved performance.
type Style struct {
	handleChains []func(element StyleElement) // Style processing function chain
	styleCache   map[string]int               // Style cache for optimizing repeated style applications
}

// NewStyle creates and returns a new Style instance.
// It initializes the style processing chain and style cache.
func NewStyle() *Style {
	return &Style{
		styleCache: make(map[string]int),
	}
}

// Apply applies the style to the specified StyleElement.
// It executes all style processing functions added to handleChains in sequence.
// The element parameter is the target element to apply styles to.
func (s *Style) Apply(element StyleElement) {
	for i := range s.handleChains {
		s.handleChains[i](element)
	}
}

// Direction sets the layout direction of the element.
// The direction parameter specifies the flow direction of the element's content, such as left-to-right, right-to-left, etc.
// It returns the Style instance to support method chaining.
func (s *Style) Direction(direction yoga.Direction) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetDirection(direction)
	})
	return s
}

// FlexDirection sets the main axis direction of the Flex container.
// The flexDirection parameter determines the arrangement direction of child elements, such as horizontal (row) or vertical (column).
// It returns the Style instance to support method chaining.
func (s *Style) FlexDirection(flexDirection yoga.FlexDirection) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetFlexDirection(flexDirection)
	})
	return s
}

func (s *Style) JustifyContent(justifyContent yoga.Justify) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetJustifyContent(justifyContent)
	})
	return s
}

func (s *Style) AlignContent(alignContent yoga.Align) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetAlignContent(alignContent)
	})
	return s
}

func (s *Style) AlignItem(alignItems yoga.Align) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetAlignItems(alignItems)
	})
	return s
}

func (s *Style) AlignSelf(alignSelf yoga.Align) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetAlignSelf(alignSelf)
	})
	return s
}

func (s *Style) FlexWrap(flexWrap yoga.Wrap) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetFlexWrap(flexWrap)
	})
	return s
}

func (s *Style) Overflow(overflow yoga.Overflow) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetOverflow(overflow)
	})
	return s
}

func (s *Style) Display(display yoga.Display) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetDisplay(display)
	})
	return s
}

func (s *Style) Flex(flex float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetFlex(flex)
	})
	return s
}

func (s *Style) FlexGrow(flexGrow float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetFlexGrow(flexGrow)
	})
	return s
}

func (s *Style) FlexShrink(flexShrink float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetFlexShrink(flexShrink)
	})
	return s
}

func (s *Style) FlexBasis(flexBasis float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetFlexBasis(flexBasis)
	})
	return s
}

func (s *Style) FlexBasisPercent(flexBasis float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetFlexBasisPercent(flexBasis)
	})
	return s
}

func (s *Style) FlexBasisAuto() *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetFlexBasisAuto()
	})
	return s
}

func (s *Style) Width(points float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetWidth(points)
	})
	return s
}

func (s *Style) WidthPercent(percent float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetWidthPercent(percent)
	})
	return s
}

func (s *Style) WidthAuto() *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetWidthAuto()
	})
	return s
}

func (s *Style) Height(height float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetHeight(height)
	})
	return s
}

func (s *Style) HeightPercent(height float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetHeightPercent(height)
	})
	return s
}

func (s *Style) HeightAuto() *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetHeightAuto()
	})
	return s
}

func (s *Style) PositionType(positionType yoga.PositionType) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetPositionType(positionType)
	})
	return s
}

func (s *Style) Position(edge yoga.Edge, position float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetPosition(edge, position)
	})
	return s
}

func (s *Style) PositionPercent(edge yoga.Edge, position float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetPositionPercent(edge, position)
	})
	return s
}

func (s *Style) Margin(edge yoga.Edge, margin float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetPositionPercent(edge, margin)
	})
	return s
}

func (s *Style) MarginPercent(edge yoga.Edge, margin float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMarginPercent(edge, margin)
	})
	return s
}

func (s *Style) MarginAuto(edge yoga.Edge) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMarginAuto(edge)
	})
	return s
}

func (s *Style) Padding(edge yoga.Edge, padding float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetPadding(edge, padding)
	})
	return s
}

func (s *Style) PaddingPercent(edge yoga.Edge, padding float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetPaddingPercent(edge, padding)
	})
	return s
}

func (s *Style) Border(edge yoga.Edge, border float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetBorder(edge, border)
	})
	return s
}

func (s *Style) Gap(edge yoga.Gutter, gapLength float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetGap(edge, gapLength)
	})
	return s
}

func (s *Style) MinWidth(minWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMinWidth(minWidth)
	})
	return s
}

func (s *Style) MinWidthPercent(minWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMaxWidthPercent(minWidth)
	})
	return s
}
func (s *Style) MinHeight(minHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMinHeight(minHeight)
	})
	return s
}

func (s *Style) MinHeightPercent(minHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMinHeightPercent(minHeight)
	})
	return s
}

func (s *Style) MaxWidth(maxWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMaxWidth(maxWidth)
	})
	return s
}

func (s *Style) MaxWidthPercent(maxWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMaxWidthPercent(maxWidth)
	})
	return s
}

func (s *Style) MaxHeight(maxHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMaxHeight(maxHeight)
	})
	return s
}

func (s *Style) MaxHeightPercent(maxHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetMaxHeightPercent(maxHeight)
	})
	return s
}

func (s *Style) AspectRatio(aspectRatio float32) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.Yoga().StyleSetAspectRatio(aspectRatio)
	})
	return s
}

// BackgroundColor sets the background color of the element.
// The backgroundColor parameter is an NRGBA color value containing red, green, blue, and alpha channels.
// This method uses the extended style mechanism to set the background color, as it is not a standard Flexbox property.
// It returns the Style instance to support method chaining.
func (s *Style) BackgroundColor(backgroundColor color.NRGBA) *Style {
	s.handleChains = append(s.handleChains, func(element StyleElement) {
		element.SetExtendedStyle(BackgroundColor{
			Color: backgroundColor,
		})
	})
	return s
}
