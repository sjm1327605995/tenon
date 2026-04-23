package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntEmpty is an empty state placeholder.
type AntEmpty struct {
	tenon.BaseWidget
	description string
	image       string
	children    []tenon.Component
}

// NewAntEmpty creates an AntEmpty.
func NewAntEmpty() *AntEmpty {
	e := &AntEmpty{description: "No Data"}
	e.Init(e)
	return e
}

// Render returns the empty state UI.
func (e *AntEmpty) Render() tenon.Component {
	theme := NewAntTheme()

	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetAlignItems(yoga.AlignCenter).
		SetJustifyContent(yoga.JustifyCenter).
		SetPadding(yoga.EdgeAll, 32)

	if e.image != "" {
		root.AddChild(components.NewImage().
			SetSource(e.image).
			SetWidth(120).
			SetHeight(120))
	} else {
		// Simple placeholder box
		root.Add(components.NewView().
			SetWidth(120).
			SetHeight(120).
			SetBackgroundColor(theme.BackgroundColor).
			SetBorderRadius(theme.BorderRadius))
	}

	root.Add(components.NewText(e.description).
		SetFontSize(theme.FontSizeBase).
		SetColor(theme.TextMutedColor).
		SetMargin(yoga.EdgeTop, 16))

	for _, child := range e.children {
		root.AddChild(child)
	}

	return root
}

func (e *AntEmpty) SetDescription(d string) *AntEmpty { e.description = d; return e }
func (e *AntEmpty) SetImage(img string) *AntEmpty     { e.image = img; return e }
func (e *AntEmpty) Add(children ...tenon.Component) *AntEmpty {
	e.children = append(e.children, children...)
	return e
}
