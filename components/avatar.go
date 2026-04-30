package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// Avatar is a circular user profile image with fallback text.
type Avatar struct {
	core.BaseElement
	size       float32
	fallback   string
	bgColor    color.Color
	textColor  color.Color
	fallbackEl *native.Text
}

// NewAvatar creates an avatar.
func NewAvatar(fallback string) *Avatar {
	theme := core.GetTheme()
	a := &Avatar{
		size:      40,
		fallback:  fallback,
		bgColor:   theme.SecondaryColor,
		textColor: theme.SecondaryForegroundColor,
	}
	a.Init(a)
	a.SetWidth(a.size)
	a.SetHeight(a.size)
	a.SetJustifyContent(yoga.JustifyCenter)
	a.SetAlignItems(yoga.AlignCenter)

	a.fallbackEl = native.NewText(fallback).SetColor(a.textColor)
	a.fallbackEl.SetFontSize(14)
	a.AppendChild(a.fallbackEl)
	return a
}

// ElementType returns type identifier.
func (a *Avatar) ElementType() string { return "Avatar" }

// Draw renders the circular avatar background.
func (a *Avatar) Draw(screen *ebiten.Image) {
	bounds := a.GetBounds()
	r := bounds.Width / 2
	cx := bounds.X + r
	cy := bounds.Y + r
	native.DrawFilledCirclePath(screen, cx, cy, r, a.bgColor)
}

// SetSize sets the avatar diameter.
func (a *Avatar) SetSize(size float32) *Avatar {
	a.size = size
	a.SetWidth(size)
	a.SetHeight(size)
	a.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return a
}

// SetFallback sets the fallback text.
func (a *Avatar) SetFallback(text string) *Avatar {
	a.fallback = text
	a.fallbackEl.SetContent(text)
	return a
}

// SetBgColor sets the avatar background color.
func (a *Avatar) SetBgColor(c color.Color) *Avatar {
	a.bgColor = c
	a.Mark(core.FlagNeedDraw)
	return a
}

// SetTextColor sets the fallback text color.
func (a *Avatar) SetTextColor(c color.Color) *Avatar {
	a.textColor = c
	a.fallbackEl.SetColor(c)
	return a
}
