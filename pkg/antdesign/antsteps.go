package antdesign

import (
	"fmt"
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntStepItem defines a single step.
type AntStepItem struct {
	Title       string
	Description string
	Icon        string
	Status      string // wait/process/finish/error
	Disabled    bool
}

// AntSteps is a step component.
type AntSteps struct {
	tenon.BaseWidget
	current   int
	direction string // horizontal/vertical
	size      string // default/small
	items     []AntStepItem
	onChange  func(current int)
}

// NewAntSteps creates an AntSteps.
func NewAntSteps() *AntSteps {
	s := &AntSteps{current: 0, direction: "horizontal", size: "default"}
	s.Init(s)
	return s
}

// Render returns the steps UI.
func (s *AntSteps) Render() tenon.Component {
	theme := NewAntTheme()

	isVertical := s.direction == "vertical"
	root := components.NewView()
	if isVertical {
		root.SetFlexDirection(yoga.FlexDirectionColumn)
	} else {
		root.SetFlexDirection(yoga.FlexDirectionRow)
	}

	for i, item := range s.items {
		status := item.Status
		if status == "" {
			if i < s.current {
				status = "finish"
			} else if i == s.current {
				status = "process"
			} else {
				status = "wait"
			}
		}

		step := s.renderStep(item, status, i, theme)
		root.AddChild(step)

		// Connector line between steps
		if i < len(s.items)-1 {
			connectorColor := theme.BorderColor
			if i < s.current {
				connectorColor = theme.PrimaryColor
			}
			if isVertical {
				root.Add(components.NewView().
					SetWidth(1).
					SetHeight(24).
					SetBackgroundColor(connectorColor).
					SetMargin(yoga.EdgeLeft, 15))
			} else {
				root.Add(components.NewView().
					SetWidth(24).
					SetHeight(1).
					SetBackgroundColor(connectorColor).
					SetMargin(yoga.EdgeTop, 15))
			}
		}
	}

	return root
}

func (s *AntSteps) renderStep(item AntStepItem, status string, index int, theme *AntTheme) tenon.Component {
	var dotColor, textColor color.Color
	var iconText string
	switch status {
	case "finish":
		dotColor = theme.PrimaryColor
		textColor = theme.TextColor
		iconText = "\u2713"
	case "process":
		dotColor = theme.PrimaryColor
		textColor = theme.PrimaryColor
		iconText = fmt.Sprintf("%d", index+1)
	case "error":
		dotColor = theme.ErrorColor
		textColor = theme.ErrorColor
		iconText = "\u2717"
	default: // wait
		dotColor = theme.BorderColor
		textColor = theme.TextMutedColor
		iconText = fmt.Sprintf("%d", index+1)
	}

	if item.Icon != "" {
		iconText = item.Icon
	}

	isVertical := s.direction == "vertical"
	step := components.NewView()
	if isVertical {
		step.SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignFlexStart)
	} else {
		step.SetFlexDirection(yoga.FlexDirectionColumn).SetAlignItems(yoga.AlignCenter)
	}

	// Icon circle
	iconSize := float32(32)
	if s.size == "small" {
		iconSize = 24
	}
	iconCircle := components.NewView().
		SetWidth(iconSize).
		SetHeight(iconSize).
		SetBackgroundColor(dotColor).
		SetBorderRadius(iconSize / 2).
		SetJustifyContent(yoga.JustifyCenter).
		SetAlignItems(yoga.AlignCenter)
	iconCircle.Add(components.NewText(iconText).
		SetFontSize(theme.FontSizeSM).
		SetColor(theme.SurfaceColor))
	step.AddChild(iconCircle)

	// Text area
	textArea := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn)
	if isVertical {
		textArea.SetMargin(yoga.EdgeLeft, 12)
	} else {
		textArea.SetMargin(yoga.EdgeTop, 8)
	}

	textArea.Add(components.NewText(item.Title).
		SetFontSize(theme.FontSizeBase).
		SetColor(textColor).
		SetFontWeight(fonts.FontWeightBold))

	if item.Description != "" {
		textArea.Add(components.NewText(item.Description).
			SetFontSize(theme.FontSizeSM).
			SetColor(theme.TextMutedColor).
			SetMargin(yoga.EdgeTop, 4))
	}

	step.AddChild(textArea)
	return step
}

func (s *AntSteps) SetCurrent(c int) *AntSteps             { s.current = c; return s }
func (s *AntSteps) SetDirection(d string) *AntSteps        { s.direction = d; return s }
func (s *AntSteps) SetSize(sz string) *AntSteps            { s.size = sz; return s }
func (s *AntSteps) SetItems(items []AntStepItem) *AntSteps { s.items = items; return s }
func (s *AntSteps) SetOnChange(fn func(int)) *AntSteps     { s.onChange = fn; return s }
