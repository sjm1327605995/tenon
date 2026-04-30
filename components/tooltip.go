package components

import (
	"image/color"

	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// Tooltip is a floating hint label, typically shown on hover.
// Use absolute positioning to place it relative to a target element.
type Tooltip struct {
	native.View
	textEl *native.Text
}

// NewTooltip creates a tooltip with the given Content.
func NewTooltip(Content string) *Tooltip {
	theme := core.GetTheme()
	tt := &Tooltip{}
	tt.Init(tt)
	tt.SetPadding(yoga.EdgeHorizontal, 10)
	tt.SetPadding(yoga.EdgeVertical, 6)
	tt.SetBackgroundColor(theme.TextColor)
	tt.SetBorderRadius(theme.BorderRadius / 2)
	tt.SetPositionType(yoga.PositionTypeAbsolute)
	tt.SetVisible(false)

	tt.textEl = native.NewText(Content).SetColor(theme.BackgroundColor)
	tt.textEl.SetFontSize(12)
	tt.AppendChild(tt.textEl)
	return tt
}

// ElementType returns type identifier.
func (tt *Tooltip) ElementType() string { return "Tooltip" }

// SetContent updates the tooltip text.
func (tt *Tooltip) SetContent(Content string) *Tooltip {
	tt.textEl.SetContent(Content)
	return tt
}

// Show makes the tooltip visible.
func (tt *Tooltip) Show() *Tooltip {
	tt.SetVisible(true)
	if eng := tt.GetEngine(); eng != nil {
		eng.AddOverlay(tt)
	}
	return tt
}

// Hide makes the tooltip invisible.
func (tt *Tooltip) Hide() *Tooltip {
	tt.SetVisible(false)
	if eng := tt.GetEngine(); eng != nil {
		eng.RemoveOverlay(tt)
	}
	return tt
}

// SetTextColor sets the tooltip text color.
func (tt *Tooltip) SetTextColor(clr color.Color) *Tooltip {
	if tt.textEl != nil {
		tt.textEl.SetColor(clr)
	}
	return tt
}
