package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Kbd renders keyboard key text in a styled box.
type Kbd struct {
	View
	textEl *Text
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

	k.textEl = NewText(text).SetColor(theme.MutedForegroundColor)
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
