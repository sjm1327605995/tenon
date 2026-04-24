package antdesign

import (
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
)

// AntLink is a hyperlink-styled text component.
type AntLink struct {
	tenon.BaseWidget
	content string
	href    string
	clr     color.Color
}

// NewAntLink creates an AntLink.
func NewAntLink(content string) *AntLink {
	l := &AntLink{content: content}
	l.Init(l)
	return l
}

// Render returns the link UI.
func (l *AntLink) Render() tenon.Component {
	theme := NewAntTheme()
	clr := theme.LinkColor
	if l.clr != nil {
		clr = l.clr
	}
	return components.NewText(l.content).
		SetColor(clr).
		SetFontSize(theme.FontSizeBase)
}

func (l *AntLink) SetHref(h string) *AntLink       { l.href = h; return l }
func (l *AntLink) SetColor(c color.Color) *AntLink { l.clr = c; return l }
