package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntParagraph is a paragraph component.
type AntParagraph struct {
	tenon.BaseWidget
	content  string
	ellipsis bool
	copyable bool
	editable bool
	rows     int
	onChange func(string)
	onCopy   func(string)
}

// NewAntParagraph creates an AntParagraph.
func NewAntParagraph(content string) *AntParagraph {
	p := &AntParagraph{content: content, rows: 3}
	p.Init(p)
	return p
}

// Render returns the paragraph UI.
func (p *AntParagraph) Render() tenon.Component {
	theme := NewAntTheme()

	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn)

	textComp := components.NewText(p.content).
		SetFontSize(theme.FontSizeBase).
		SetColor(theme.TextColor)

	if p.ellipsis {
		textComp.SetWhiteSpace(components.WhiteSpaceNormal)
		textComp.SetWordBreak(components.WordBreakBreakAll)
		// Limit to approximate rows via max height
		lineHeight := theme.FontSizeBase * 1.5
		textComp.SetHeight(float32(p.rows) * lineHeight)
		root.SetOverflow(yoga.OverflowHidden)
	}

	root.AddChild(textComp)

	// Copy button row
	if p.copyable {
		copyBtn := components.NewButton("Copy").
			SetOnClick(func() {
				if p.onCopy != nil {
					p.onCopy(p.content)
				}
			})
		copyBtn.SetMargin(yoga.EdgeTop, 4)
		root.AddChild(copyBtn)
	}

	return root
}

func (p *AntParagraph) SetEllipsis(v bool) *AntParagraph          { p.ellipsis = v; return p }
func (p *AntParagraph) SetCopyable(v bool) *AntParagraph          { p.copyable = v; return p }
func (p *AntParagraph) SetEditable(v bool) *AntParagraph          { p.editable = v; return p }
func (p *AntParagraph) SetRows(n int) *AntParagraph               { p.rows = n; return p }
func (p *AntParagraph) SetOnChange(fn func(string)) *AntParagraph { p.onChange = fn; return p }
func (p *AntParagraph) SetOnCopy(fn func(string)) *AntParagraph   { p.onCopy = fn; return p }
