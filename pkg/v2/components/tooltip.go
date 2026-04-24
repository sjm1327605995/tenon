package components

import (
	"image/color"

	"github.com/sjm1327605995/tenon/yoga"
)

// Tooltip is a floating hint label, typically shown on hover.
// Use absolute positioning to place it relative to a target element.
type Tooltip struct {
	View
	textEl *Text
}

// NewTooltip creates a tooltip with the given content.
func NewTooltip(content string) *Tooltip {
	tt := &Tooltip{}
	tt.Init(tt)
	tt.SetPadding(yoga.EdgeAll, 6)
	tt.SetBackgroundColor(color.RGBA{R: 50, G: 50, B: 50, A: 230})
	tt.SetBorderRadius(4)
	tt.SetPositionType(yoga.PositionTypeAbsolute)
	tt.SetVisible(false)

	tt.textEl = NewText(content).SetColor(color.White)
	tt.textEl.SetFontSize(10)
	tt.AppendChild(tt.textEl)
	return tt
}

// ElementType returns type identifier.
func (tt *Tooltip) ElementType() string { return "Tooltip" }

// SetContent updates the tooltip text.
func (tt *Tooltip) SetContent(content string) *Tooltip {
	tt.textEl.SetContent(content)
	return tt
}

// Show makes the tooltip visible.
func (tt *Tooltip) Show() *Tooltip {
	tt.SetVisible(true)
	return tt
}

// Hide makes the tooltip invisible.
func (tt *Tooltip) Hide() *Tooltip {
	tt.SetVisible(false)
	return tt
}

// SetTextColor sets the tooltip text color.
func (tt *Tooltip) SetTextColor(clr color.Color) *Tooltip {
	tt.textEl.SetColor(clr)
	return tt
}
