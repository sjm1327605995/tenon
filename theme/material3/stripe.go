package material3

import (
	"github.com/sjm1327605995/tenon/core/stripe"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/icon"
	"github.com/sjm1327605995/tenon/widget"
)

// StripePainter renders tool window stripes using Material 3 design tokens.
//
// If Theme is nil, StripePainter falls back to the default M3 purple palette.
type StripePainter struct {
	Theme *Theme
}

// resolveColors returns M3-derived colors for stripe painting.
func (p StripePainter) resolveColors() stripeColors {
	if p.Theme == nil {
		return m3DefaultStripeColors
	}
	cs := p.Theme.Colors
	return stripeColors{
		Background:  cs.SurfaceContainer,
		IconColor:   cs.OnSurface,
		ActiveFg:    cs.Primary,
		HoverBg:     cs.OnSurface.WithAlpha(0.04),
		PressedBg:   cs.OnSurface.WithAlpha(0.08),
		ActiveBg:    cs.Primary.WithAlpha(0.08),
	}
}

// PaintBackground renders the stripe background with M3 surface-container color.
func (p StripePainter) PaintBackground(canvas widget.Canvas, bounds geometry.Rect) {
	if bounds.IsEmpty() {
		return
	}
	colors := p.resolveColors()
	canvas.DrawRect(bounds, colors.Background)
}

// PaintButton renders a tool window button with icon and optional label.
func (p StripePainter) PaintButton(canvas widget.Canvas, state stripe.ButtonPaintState) {
	if state.Bounds.IsEmpty() {
		return
	}
	colors := p.resolveColors()

	bg := widget.ColorTransparent
	switch {
	case state.Pressed:
		bg = colors.PressedBg
	case state.Hovered:
		bg = colors.HoverBg
	case state.Active:
		bg = colors.ActiveBg
	}
	if !bg.IsTransparent() {
		canvas.DrawRoundRect(state.Bounds, bg, m3StripeBtnRadius)
	}

	fg := colors.IconColor
	if state.Active {
		fg = colors.ActiveFg
	}

	iconBounds := m3StripeIconBounds(state.Bounds, state.ShowLabel)
	icon.Draw(canvas, state.Icon, iconBounds, fg)

	if state.ShowLabel && state.Label != "" {
		textBounds := m3StripeTextBounds(state.Bounds, iconBounds)
		canvas.DrawText(state.Label, textBounds, m3StripeLabelFontSize, fg, false, widget.TextAlignCenter)
	}
}

func m3StripeIconBounds(btnBounds geometry.Rect, showLabel bool) geometry.Rect {
	if showLabel {
		centerX := btnBounds.Min.X + btnBounds.Width()/2
		y := btnBounds.Min.Y + m3StripeIconPaddingLabel
		return geometry.NewRect(centerX-m3StripeIconSize/2, y, m3StripeIconSize, m3StripeIconSize)
	}
	centerX := btnBounds.Min.X + btnBounds.Width()/2
	centerY := btnBounds.Min.Y + btnBounds.Height()/2
	return geometry.NewRect(centerX-m3StripeIconSize/2, centerY-m3StripeIconSize/2, m3StripeIconSize, m3StripeIconSize)
}

func m3StripeTextBounds(btnBounds, iconRect geometry.Rect) geometry.Rect {
	y := iconRect.Max.Y + m3StripeIconTextGap
	return geometry.NewRect(
		btnBounds.Min.X+m3StripeTextPaddingH,
		y,
		btnBounds.Width()-m3StripeTextPaddingH*2,
		btnBounds.Max.Y-y,
	)
}

type stripeColors struct {
	Background widget.Color
	IconColor  widget.Color
	ActiveFg   widget.Color
	HoverBg    widget.Color
	PressedBg  widget.Color
	ActiveBg   widget.Color
}

var m3DefaultStripeColors = stripeColors{
	Background: widget.Hex(0xECE6F0),
	IconColor:  widget.Hex(0x1C1B1F),
	ActiveFg:   widget.Hex(0x6750A4),
	HoverBg:    widget.RGBA(0.12, 0.12, 0.13, 0.04),
	PressedBg:  widget.RGBA(0.12, 0.12, 0.13, 0.08),
	ActiveBg:   widget.Hex(0x6750A4).WithAlpha(0.08),
}

const (
	m3StripeIconSize         float32 = 20
	m3StripeIconPaddingLabel float32 = 4
	m3StripeIconTextGap      float32 = 3
	m3StripeLabelFontSize    float32 = 11
	m3StripeTextPaddingH     float32 = 6
	m3StripeBtnRadius        float32 = 12
)

var _ stripe.Painter = StripePainter{}
