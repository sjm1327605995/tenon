package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
)

// AccordionPane is a single accordion entry.
type AccordionPane struct {
	Title   string
	Content engine.Widget
}

// ShadcnAccordion renders a vertically stacked set of collapsible panels.
// expanded maps item index -> true if open.
// onChange is called with the clicked index; the caller should update expanded
// and (if single mode) collapse other items.
func ShadcnAccordion(items []AccordionPane, expanded map[int]bool, onChange func(int)) engine.Widget {
	if expanded == nil {
		expanded = map[int]bool{}
	}
	t := engine.GetTheme()
	children := make([]engine.Widget, 0, len(items)*2)
	for i, item := range items {
		idx := i
		isOpen := expanded[idx]
		chevron := widgets.IconChevronRight
		if isOpen {
			chevron = widgets.IconChevronDown
		}

		// Header row
		header := widgets.Row(
			widgets.Text(item.Title).
				Color(t.TextColor).
				FontSize(t.FontSizeBase),
			widgets.Container(widgets.Text("")).Grow(1),
			widgets.Icon(chevron).
				Color(t.MutedForegroundColor).
				FontSize(t.FontSizeSM),
		).AlignItems(engine.AlignCenter).JustifyContent(engine.JustifySpaceBetween)

		// Clickable header container (no ghost button overlay needed)
		trigger := widgets.Container(header).
			Pad(engine.EdgeInsetsSymmetric(12, 0)).
			OnTap(func() {
				if onChange != nil {
					onChange(idx)
				}
			})

		children = append(children, trigger)

		if isOpen && item.Content != nil {
			children = append(children, widgets.Container(item.Content).Pad(engine.EdgeInsetsOnly(0, 0, 12, 12)))
		}

		if i < len(items)-1 {
			children = append(children, ShadcnSeparator(SeparatorHorizontal))
		}
	}
	return widgets.Column(children...)
}
