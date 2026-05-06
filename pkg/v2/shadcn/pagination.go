package shadcn

import (
	"fmt"

	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// ShadcnPagination renders a row of page buttons.
// current and total are 1-based.
func ShadcnPagination(current, total int, onChange func(int)) ui.Widget {
	if total <= 0 {
		return widgets.Row()
	}
	if current < 1 {
		current = 1
	}
	if current > total {
		current = total
	}
	maxVisible := 5
	if total < maxVisible {
		maxVisible = total
	}
	start, end := visibleRange(current, total, maxVisible)

	children := make([]ui.Widget, 0, maxVisible+2)

	prev := widgets.Button(widgets.IconText(widgets.IconArrowLeft)).Variantf(widgets.ButtonOutline).OnTap(func() {
		if current > 1 && onChange != nil {
			onChange(current - 1)
		}
	}).SetDisabled(current <= 1)
	children = append(children, prev)

	for page := start; page <= end; page++ {
		p := page
		var btn ui.Widget
		if p == current {
			btn = widgets.Button(fmt.Sprintf("%d", p)).
				Variantf(widgets.ButtonDefault).
				OnTap(func() {
					if onChange != nil {
						onChange(p)
					}
				})
		} else {
			btn = widgets.Button(fmt.Sprintf("%d", p)).
				Variantf(widgets.ButtonOutline).
				OnTap(func() {
					if onChange != nil {
						onChange(p)
					}
				})
		}
		children = append(children, btn)
	}

	next := widgets.Button(widgets.IconText(widgets.IconArrowRight)).Variantf(widgets.ButtonOutline).OnTap(func() {
		if current < total && onChange != nil {
			onChange(current + 1)
		}
	}).SetDisabled(current >= total)
	children = append(children, next)

	return widgets.Row(children...).AlignItems(ui.AlignCenter).Gapf(4)
}

func visibleRange(current, total, maxVisible int) (int, int) {
	start := current - 2
	if start < 1 {
		start = 1
	}
	end := start + maxVisible - 1
	if end > total {
		end = total
		start = end - maxVisible + 1
		if start < 1 {
			start = 1
		}
	}
	return start, end
}
