package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// TabItem defines a single tab.
type TabItem struct {
	Label   string
	Content ui.Widget
}

// ShadcnTabs renders a tabbed interface.
func ShadcnTabs(items []TabItem, activeIndex int, onChange func(int)) ui.Widget {
	var tabButtons []ui.Widget
	for i, tab := range items {
		idx := i
		variant := widgets.ButtonGhost
		if i == activeIndex {
			variant = widgets.ButtonDefault
		}
		btn := widgets.Button(tab.Label).Variantf(variant).OnTap(func() {
			if onChange != nil {
				onChange(idx)
			}
		})
		tabButtons = append(tabButtons, btn)
	}
	tabRow := widgets.Row(tabButtons...).Gapf(4).Paddingf(ui.EdgeInsets{Bottom: 12})

	var children []ui.Widget
	children = append(children, tabRow)
	if activeIndex >= 0 && activeIndex < len(items) {
		if items[activeIndex].Content != nil {
			children = append(children, items[activeIndex].Content)
		}
	}
	return widgets.Container(widgets.Column(children...)).WPct(100)
}
