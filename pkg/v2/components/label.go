package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// Label is a text label typically used with form inputs.
type Label struct {
	Text
}

// NewLabel creates a label.
func NewLabel(text string) *Label {
	theme := core.GetTheme()
	l := &Label{}
	l.Init(l)
	l.SetContent(text)
	l.SetColor(theme.TextColor)
	l.SetFontSize(14)
	return l
}

// ElementType returns type identifier.
func (l *Label) ElementType() string { return "Label" }

// SetFor sets a visual association hint (no-op in immediate mode UI).
func (l *Label) SetFor(target string) *Label {
	return l
}
