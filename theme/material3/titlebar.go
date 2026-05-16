package material3

import (
	"github.com/sjm1327605995/tenon/core/titlebar"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/widget"
)

// TitleBarPainter renders title bars using Material 3 design tokens.
//
// If Theme is nil, TitleBarPainter falls back to the default M3 purple palette.
type TitleBarPainter struct {
	Theme *Theme
}

// resolveColors returns M3-derived colors for title bar painting.
func (p TitleBarPainter) resolveColors() titlebarColors {
	if p.Theme == nil {
		return m3DefaultTitleBarColors
	}
	cs := p.Theme.Colors
	return titlebarColors{
		Background:   cs.SurfaceContainer,
		ControlFg:    cs.OnSurface,
		HoverBg:      cs.OnSurface.WithAlpha(0.08),
		PressedBg:    cs.OnSurface.WithAlpha(0.12),
		CloseHoverBg: widget.Hex(0xC42B1C),
		ClosePressBg: widget.Hex(0xB22A1A),
	}
}

// DrawBackground renders the title bar background with M3 surface color.
func (p TitleBarPainter) DrawBackground(canvas widget.Canvas, bounds geometry.Rect, state titlebar.BackgroundState) {
	if bounds.IsEmpty() {
		return
	}
	colors := p.resolveColors()
	bg := colors.Background
	if !state.Focused {
		bg = bg.WithAlpha(0.7)
	}
	canvas.DrawRect(bounds, bg)
}

// DrawControlButton renders a window control button using M3 colors.
func (p TitleBarPainter) DrawControlButton(canvas widget.Canvas, bounds geometry.Rect, control titlebar.ControlType, state titlebar.ControlState) {
	if bounds.IsEmpty() {
		return
	}
	colors := p.resolveColors()

	bg := widget.ColorTransparent
	if control == titlebar.ControlClose {
		if state.Pressed {
			bg = colors.ClosePressBg
		} else if state.Hovered {
			bg = colors.CloseHoverBg
		}
	} else {
		if state.Pressed {
			bg = colors.PressedBg
		} else if state.Hovered {
			bg = colors.HoverBg
		}
	}
	if !bg.IsTransparent() {
		canvas.DrawRect(bounds, bg)
	}

	fg := colors.ControlFg
	if control == titlebar.ControlClose && state.Hovered {
		fg = widget.ColorWhite
	}

	cx := bounds.Min.X + bounds.Width()/2
	cy := bounds.Min.Y + bounds.Height()/2

	switch control {
	case titlebar.ControlMinimize:
		canvas.DrawLine(
			geometry.Pt(cx-m3TitleIconHalf, cy),
			geometry.Pt(cx+m3TitleIconHalf, cy),
			fg, m3TitleIconStroke,
		)
	case titlebar.ControlMaximize:
		r := geometry.NewRect(cx-m3TitleIconHalf, cy-m3TitleIconHalf, m3TitleIconSize, m3TitleIconSize)
		canvas.StrokeRect(r, fg, m3TitleIconStroke)
	case titlebar.ControlRestore:
		offset := float32(2)
		half := m3TitleIconHalf - 1
		back := geometry.NewRect(cx-half+offset, cy-half-offset, half*2, half*2)
		canvas.StrokeRect(back, fg, m3TitleIconStroke)
		front := geometry.NewRect(cx-half, cy-half, half*2, half*2)
		canvas.DrawRect(front, colors.Background)
		canvas.StrokeRect(front, fg, m3TitleIconStroke)
	case titlebar.ControlClose:
		canvas.DrawLine(
			geometry.Pt(cx-m3TitleIconHalf, cy-m3TitleIconHalf),
			geometry.Pt(cx+m3TitleIconHalf, cy+m3TitleIconHalf),
			fg, m3TitleCloseStroke,
		)
		canvas.DrawLine(
			geometry.Pt(cx+m3TitleIconHalf, cy-m3TitleIconHalf),
			geometry.Pt(cx-m3TitleIconHalf, cy+m3TitleIconHalf),
			fg, m3TitleCloseStroke,
		)
	}
}

type titlebarColors struct {
	Background   widget.Color
	ControlFg    widget.Color
	HoverBg      widget.Color
	PressedBg    widget.Color
	CloseHoverBg widget.Color
	ClosePressBg widget.Color
}

var m3DefaultTitleBarColors = titlebarColors{
	Background:   widget.Hex(0xECE6F0),
	ControlFg:    widget.Hex(0x1C1B1F),
	HoverBg:      widget.RGBA(0.12, 0.12, 0.13, 0.08),
	PressedBg:    widget.RGBA(0.12, 0.12, 0.13, 0.12),
	CloseHoverBg: widget.Hex(0xC42B1C),
	ClosePressBg: widget.Hex(0xB22A1A),
}

const (
	m3TitleIconSize    float32 = 10
	m3TitleIconHalf    float32 = 5
	m3TitleIconStroke  float32 = 1
	m3TitleCloseStroke float32 = 1.5
)

var _ titlebar.Painter = TitleBarPainter{}
