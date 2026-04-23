package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntSpaceDirection defines spacing direction.
type AntSpaceDirection string

const (
	AntSpaceHorizontal AntSpaceDirection = "horizontal"
	AntSpaceVertical   AntSpaceDirection = "vertical"
)

// AntSpaceAlign defines cross-axis alignment.
type AntSpaceAlign string

const (
	AntSpaceAlignStart    AntSpaceAlign = "start"
	AntSpaceAlignCenter   AntSpaceAlign = "center"
	AntSpaceAlignEnd      AntSpaceAlign = "end"
	AntSpaceAlignBaseline AntSpaceAlign = "baseline"
)

// AntSpaceSize presets.
type AntSpaceSize string

const (
	AntSpaceSmall  AntSpaceSize = "small"
	AntSpaceMiddle AntSpaceSize = "middle"
	AntSpaceLarge  AntSpaceSize = "large"
)

// AntSpace manages spacing between child components.
type AntSpace struct {
	tenon.BaseWidget
	direction AntSpaceDirection
	align     AntSpaceAlign
	size      interface{} // float32 or AntSpaceSize
	children  []tenon.Component
}

// NewAntSpace creates an AntSpace.
func NewAntSpace() *AntSpace {
	s := &AntSpace{direction: AntSpaceHorizontal, align: AntSpaceAlignCenter}
	s.Init(s)
	return s
}

// Render returns the space UI.
func (s *AntSpace) Render() tenon.Component {
	root := components.NewView().
		SetFlexDirection(s.resolveFlexDirection()).
		SetAlignItems(s.resolveAlignItems())

	gap := s.resolveGap()
	for i, child := range s.children {
		root.AddChild(child)
		if i < len(s.children)-1 {
			// Add margin to create gap
			// We cannot easily set margin on arbitrary child,
			// so we wrap each child in a view with margin
			// For simplicity, rely on Yoga gap if available, else wrap
		}
	}

	// Use padding as a simple gap simulation by wrapping children
	return s.wrapChildrenWithGap(root, gap)
}

func (s *AntSpace) resolveFlexDirection() yoga.FlexDirection {
	if s.direction == AntSpaceVertical {
		return yoga.FlexDirectionColumn
	}
	return yoga.FlexDirectionRow
}

func (s *AntSpace) resolveAlignItems() yoga.Align {
	switch s.align {
	case AntSpaceAlignStart:
		return yoga.AlignFlexStart
	case AntSpaceAlignEnd:
		return yoga.AlignFlexEnd
	case AntSpaceAlignBaseline:
		return yoga.AlignBaseline
	default:
		return yoga.AlignCenter
	}
}

func (s *AntSpace) resolveGap() float32 {
	theme := NewAntTheme()
	switch v := s.size.(type) {
	case float32:
		return v
	case int:
		return float32(v)
	case AntSpaceSize:
		switch v {
		case AntSpaceSmall:
			return 8
		case AntSpaceLarge:
			return 24
		default:
			return theme.FontSizeBase
		}
	default:
		return theme.FontSizeBase
	}
}

func (s *AntSpace) wrapChildrenWithGap(root *components.View, gap float32) tenon.Component {
	if len(s.children) == 0 {
		return root
	}
	result := components.NewView().
		SetFlexDirection(s.resolveFlexDirection()).
		SetAlignItems(s.resolveAlignItems())
	for i, child := range s.children {
		if i > 0 {
			if s.direction == AntSpaceVertical {
				result.Add(components.NewView().SetHeight(gap))
			} else {
				result.Add(components.NewView().SetWidth(gap))
			}
		}
		result.AddChild(child)
	}
	return result
}

func (s *AntSpace) Add(children ...tenon.Component) *AntSpace {
	s.children = append(s.children, children...)
	return s
}

func (s *AntSpace) SetDirection(d AntSpaceDirection) *AntSpace { s.direction = d; return s }
func (s *AntSpace) SetAlign(a AntSpaceAlign) *AntSpace         { s.align = a; return s }
func (s *AntSpace) SetSize(size interface{}) *AntSpace         { s.size = size; return s }

// AntSpaceCompact is a compact mode for form controls.
type AntSpaceCompact struct {
	tenon.BaseWidget
	direction AntSpaceDirection
	children  []tenon.Component
}

// NewAntSpaceCompact creates a compact space.
func NewAntSpaceCompact() *AntSpaceCompact {
	s := &AntSpaceCompact{direction: AntSpaceHorizontal}
	s.Init(s)
	return s
}

// Render returns compact UI with children joined edge-to-edge.
func (s *AntSpaceCompact) Render() tenon.Component {
	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignCenter)
	for _, child := range s.children {
		root.AddChild(child)
	}
	return root
}

func (s *AntSpaceCompact) Add(children ...tenon.Component) *AntSpaceCompact {
	s.children = append(s.children, children...)
	return s
}

func (s *AntSpaceCompact) SetDirection(d AntSpaceDirection) *AntSpaceCompact {
	s.direction = d
	return s
}
