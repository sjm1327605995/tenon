package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// BreadcrumbItem represents a single breadcrumb entry.
type BreadcrumbItem struct {
	Label    string
	Href     string // placeholder for navigation reference
	IsActive bool
}

// Breadcrumb is a navigation hierarchy indicator.
type Breadcrumb struct {
	core.BaseElement
	items   []BreadcrumbItem
	sepText string
}

// NewBreadcrumb creates a breadcrumb.
func NewBreadcrumb(items []BreadcrumbItem) *Breadcrumb {
	b := &Breadcrumb{
		items:   items,
		sepText: "/",
	}
	b.Init(b)
	b.SetFlexDirection(yoga.FlexDirectionRow)
	b.SetAlignItems(yoga.AlignCenter)
	b.SetGap(yoga.GutterAll, 8)
	b.buildItems()
	return b
}

// ElementType returns type identifier.
func (b *Breadcrumb) ElementType() string { return "Breadcrumb" }

func (b *Breadcrumb) buildItems() {
	for i, item := range b.items {
		if i > 0 {
			sep := NewText(b.sepText).SetColor(core.GetTheme().MutedForegroundColor)
			sep.SetFontSize(14)
			b.AppendChild(sep)
		}
		var txt *Text
		if item.IsActive {
			txt = NewText(item.Label).SetColor(core.GetTheme().TextColor)
		} else {
			txt = NewText(item.Label).SetColor(core.GetTheme().MutedForegroundColor)
		}
		txt.SetFontSize(14)
		b.AppendChild(txt)
	}
}

// SetItems rebuilds with new items.
func (b *Breadcrumb) SetItems(items []BreadcrumbItem) *Breadcrumb {
	b.items = items
	b.ClearChildren()
	b.buildItems()
	b.Mark(core.FlagNeedLayout)
	return b
}

// SetSeparator sets the separator character.
func (b *Breadcrumb) SetSeparator(sep string) *Breadcrumb {
	b.sepText = sep
	b.ClearChildren()
	b.buildItems()
	b.Mark(core.FlagNeedLayout)
	return b
}
