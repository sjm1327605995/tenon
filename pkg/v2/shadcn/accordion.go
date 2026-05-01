package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// AccordionPane is a single accordion entry.
type AccordionPane struct {
	Title   string
	Content ui.Widget
}

// ShadcnAccordion renders a vertically stacked set of collapsible panels.
// expanded maps item index -> true if open.
// onChange is called with the clicked index; the caller should update expanded
// and (if single mode) collapse other items.
func ShadcnAccordion(items []AccordionPane, expanded map[int]bool, onChange func(int)) ui.Widget {
	if expanded == nil {
		expanded = map[int]bool{}
	}
	t := ui.GetTheme()
	children := make([]ui.Widget, 0, len(items)*2)
	for i, item := range items {
		idx := i
		isOpen := expanded[idx]
		chevron := "▶"
		if isOpen {
			chevron = "▼"
		}

		// Header row
		header := widgets.Row(
			widgets.Text(item.Title).
				Color(render.NewColorFrom(t.TextColor)).
				FontSize(14),
			widgets.Container(widgets.Text("")).Grow(1),
			widgets.Text(chevron).
				Color(render.NewColorFrom(t.MutedForegroundColor)).
				FontSize(12),
		).AlignItems(ui.AlignCenter).JustifyContent(ui.JustifySpaceBetween)

		// Full-width clickable ghost button as trigger
		trigger := widgets.Button("").
			Variantf(widgets.ButtonGhost).
			OnTap(func() {
				if onChange != nil {
					onChange(idx)
				}
			})

		// Stack: ghost button fills area, header text overlays on top
		triggerStack := widgets.Stack(
			widgets.Container(trigger).WPct(100).HPct(100),
			widgets.Container(header).Pad(ui.EdgeInsetsSymmetric(12, 0)).WPct(100),
		)

		children = append(children, triggerStack)

		if isOpen && item.Content != nil {
			children = append(children, widgets.Container(item.Content).Pad(ui.EdgeInsetsOnly(0, 0, 12, 12)))
		}

		if i < len(items)-1 {
			children = append(children, ShadcnSeparator(SeparatorHorizontal))
		}
	}
	return widgets.Column(children...)
}
