package style

import (
	"image/color"

	"github.com/sjm1327605995/tenon/react/common"
	"github.com/sjm1327605995/tenon/react/yoga"
)

type Style struct {
	handleChains []func(element common.Element)
}

func NewStyle() *Style {
	return &Style{}
}
func (s *Style) Apply(element common.Element) {
	for i := range s.handleChains {
		s.handleChains[i](element)
	}
}

func (s *Style) Direction(direction yoga.Direction) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetDirection(direction)
	})
	return s
}

func (s *Style) FlexDirection(flexDirection yoga.FlexDirection) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetFlexDirection(flexDirection)
	})
	return s
}

func (s *Style) JustifyContent(justifyContent yoga.Justify) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetJustifyContent(justifyContent)
	})
	return s
}

func (s *Style) AlignContent(alignContent yoga.Align) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetAlignContent(alignContent)
	})
	return s
}

func (s *Style) AlignItem(alignItems yoga.Align) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetAlignItems(alignItems)
	})
	return s
}

func (s *Style) AlignSelf(alignSelf yoga.Align) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetAlignSelf(alignSelf)
	})
	return s
}

func (s *Style) FlexWrap(flexWrap yoga.Wrap) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetFlexWrap(flexWrap)
	})
	return s
}

func (s *Style) Overflow(overflow yoga.Overflow) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetOverflow(overflow)
	})
	return s
}

func (s *Style) Display(display yoga.Display) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetDisplay(display)
	})
	return s
}

func (s *Style) Flex(flex float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetFlex(flex)
	})
	return s
}

func (s *Style) FlexGrow(flexGrow float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetFlexGrow(flexGrow)
	})
	return s
}

func (s *Style) FlexShrink(flexShrink float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetFlexShrink(flexShrink)
	})
	return s
}

func (s *Style) FlexBasis(flexBasis float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetFlexBasis(flexBasis)
	})
	return s
}

func (s *Style) FlexBasisPercent(flexBasis float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetFlexBasisPercent(flexBasis)
	})
	return s
}

func (s *Style) FlexBasisAuto() *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetFlexBasisAuto()
	})
	return s
}

func (s *Style) Width(points float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetWidth(points)
	})
	return s
}

func (s *Style) WidthPercent(percent float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetWidthPercent(percent)
	})
	return s
}

func (s *Style) WidthAuto() *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetWidthAuto()
	})
	return s
}

func (s *Style) Height(height float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetHeight(height)
	})
	return s
}

func (s *Style) HeightPercent(height float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetHeightPercent(height)
	})
	return s
}

func (s *Style) HeightAuto() *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetHeightAuto()
	})
	return s
}

func (s *Style) PositionType(positionType yoga.PositionType) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetPositionType(positionType)
	})
	return s
}

func (s *Style) Position(edge yoga.Edge, position float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetPosition(edge, position)
	})
	return s
}

func (s *Style) PositionPercent(edge yoga.Edge, position float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetPositionPercent(edge, position)
	})
	return s
}

func (s *Style) Margin(edge yoga.Edge, margin float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetPositionPercent(edge, margin)
	})
	return s
}

func (s *Style) MarginPercent(edge yoga.Edge, margin float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMarginPercent(edge, margin)
	})
	return s
}

func (s *Style) MarginAuto(edge yoga.Edge) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMarginAuto(edge)
	})
	return s
}

func (s *Style) Padding(edge yoga.Edge, padding float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetPadding(edge, padding)
	})
	return s
}

func (s *Style) PaddingPercent(edge yoga.Edge, padding float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetPaddingPercent(edge, padding)
	})
	return s
}

func (s *Style) Border(edge yoga.Edge, border float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetBorder(edge, border)
	})
	return s
}

func (s *Style) Gap(edge yoga.Gutter, gapLength float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetGap(edge, gapLength)
	})
	return s
}

func (s *Style) MinWidth(minWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMinWidth(minWidth)
	})
	return s
}

func (s *Style) MinWidthPercent(minWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMaxWidthPercent(minWidth)
	})
	return s
}
func (s *Style) MinHeight(minHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMinHeight(minHeight)
	})
	return s
}

func (s *Style) MinHeightPercent(minHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMinHeightPercent(minHeight)
	})
	return s
}

func (s *Style) MaxWidth(maxWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMaxWidth(maxWidth)
	})
	return s
}

func (s *Style) MaxWidthPercent(maxWidth float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMaxWidthPercent(maxWidth)
	})
	return s
}

func (s *Style) MaxHeight(maxHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMaxHeight(maxHeight)
	})
	return s
}

func (s *Style) MaxHeightPercent(maxHeight float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetMaxHeightPercent(maxHeight)
	})
	return s
}

func (s *Style) AspectRatio(aspectRatio float32) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.Yoga().StyleSetAspectRatio(aspectRatio)
	})
	return s
}

func (s *Style) BackgroundColor(backgroundColor color.NRGBA) *Style {
	s.handleChains = append(s.handleChains, func(element common.Element) {
		element.SetExtendedStyle(BackgroundColor{
			Color: backgroundColor,
		})
	})
	return s
}
