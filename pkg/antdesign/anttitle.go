package antdesign

import (
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/fonts"
)

// AntTitleLevel defines heading levels.
type AntTitleLevel int

const (
	AntTitleH1 AntTitleLevel = 1
	AntTitleH2 AntTitleLevel = 2
	AntTitleH3 AntTitleLevel = 3
	AntTitleH4 AntTitleLevel = 4
)

// AntTitle is a heading component.
type AntTitle struct {
	tenon.BaseWidget
	content string
	level   AntTitleLevel
	clr     color.Color
}

// NewAntTitle creates an AntTitle.
func NewAntTitle(content string) *AntTitle {
	t := &AntTitle{content: content, level: AntTitleH1}
	t.Init(t)
	return t
}

// Render returns the title UI.
func (t *AntTitle) Render() tenon.Component {
	theme := NewAntTheme()
	fontSize, fontWeight := t.resolveStyle(theme)

	clr := theme.TextColor
	if t.clr != nil {
		clr = t.clr
	}

	return components.NewText(t.content).
		SetFontSize(fontSize).
		SetFontWeight(fontWeight).
		SetColor(clr)
}

func (t *AntTitle) resolveStyle(theme *AntTheme) (fontSize float32, fontWeight fonts.FontWeight) {
	switch t.level {
	case AntTitleH1:
		return theme.FontSizeBase + 24, fonts.FontWeightBold
	case AntTitleH2:
		return theme.FontSizeBase + 16, fonts.FontWeightBold
	case AntTitleH3:
		return theme.FontSizeBase + 8, fonts.FontWeightBold
	case AntTitleH4:
		return theme.FontSizeBase + 4, fonts.FontWeightBold
	default:
		return theme.FontSizeBase + 24, fonts.FontWeightBold
	}
}

// SetLevel sets the heading level.
func (t *AntTitle) SetLevel(l AntTitleLevel) *AntTitle {
	t.level = l
	return t
}

// SetColor sets the text color.
func (t *AntTitle) SetColor(c color.Color) *AntTitle {
	t.clr = c
	return t
}
