package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// Kbd renders keyboard key text in a styled box.
type Kbd struct {
	native.View
	textEl *native.Text
}

// NewKbd creates a keyboard key element.
func NewKbd(text string) *Kbd {
	theme := core.GetTheme()
	k := &Kbd{}
	k.Init(k)
	k.SetPadding(yoga.EdgeHorizontal, 8)
	k.SetPadding(yoga.EdgeVertical, 4)
	k.SetBackgroundColor(theme.MutedColor)
	k.SetBorderRadius(theme.BorderRadius / 2)
	k.SetBorderColor(theme.BorderColor)

	k.textEl = native.NewText(text).SetColor(theme.MutedForegroundColor)
	k.textEl.SetFontSize(12)
	k.AppendChild(k.textEl)
	return k
}

// ElementType returns type identifier.
func (k *Kbd) ElementType() string { return "Kbd" }

// SetContent updates the key text.
func (k *Kbd) SetContent(text string) *Kbd {
	k.textEl.SetContent(text)
	return k
}
