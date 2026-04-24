package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntBreadcrumbItem defines a breadcrumb item.
type AntBreadcrumbItem struct {
	Label    string
	Href     string
	Children tenon.Component
}

// AntBreadcrumb is a breadcrumb navigation component.
type AntBreadcrumb struct {
	tenon.BaseWidget
	items     []AntBreadcrumbItem
	separator string
}

// NewAntBreadcrumb creates an AntBreadcrumb.
func NewAntBreadcrumb() *AntBreadcrumb {
	b := &AntBreadcrumb{separator: "/"}
	b.Init(b)
	return b
}

// Render returns the breadcrumb UI.
func (b *AntBreadcrumb) Render() tenon.Component {
	theme := NewAntTheme()

	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignCenter).
		SetFlexWrap(yoga.WrapWrap)

	for i, item := range b.items {
		if i > 0 {
			root.Add(components.NewText(b.separator).
				SetFontSize(theme.FontSizeBase).
				SetColor(theme.TextMutedColor).
				SetMargin(yoga.EdgeHorizontal, 8))
		}

		var content tenon.Component
		if item.Children != nil {
			content = item.Children
		} else {
			labelColor := theme.TextColor
			if item.Href != "" {
				labelColor = theme.LinkColor
			}
			content = components.NewText(item.Label).
				SetFontSize(theme.FontSizeBase).
				SetColor(labelColor)
		}
		root.AddChild(content)
	}

	return root
}

func (b *AntBreadcrumb) SetItems(items []AntBreadcrumbItem) *AntBreadcrumb { b.items = items; return b }
func (b *AntBreadcrumb) SetSeparator(s string) *AntBreadcrumb              { b.separator = s; return b }
