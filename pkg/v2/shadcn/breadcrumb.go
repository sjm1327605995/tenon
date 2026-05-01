package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// BreadcrumbItem is a single entry in a breadcrumb trail.
type BreadcrumbItem struct {
	Label    string
	Href     string // navigation reference (unused in UI)
	IsActive bool
	OnTap    func()
}

// ShadcnBreadcrumb renders a horizontal breadcrumb trail.
// Use ui.GetTheme().TextColor for active, MutedForegroundColor for inactive.
func ShadcnBreadcrumb(items []BreadcrumbItem, separator string) ui.Widget {
	if separator == "" {
		separator = "/"
	}
	t := ui.GetTheme()
	children := make([]ui.Widget, 0, len(items)*2-1)
	for i, item := range items {
		if i > 0 {
			children = append(children,
				widgets.Text(separator).
					Color(render.NewColorFrom(t.MutedForegroundColor)).
					FontSize(t.FontSizeSM),
				)
		}
		lbl := item.Label
		if item.IsActive && item.OnTap == nil {
			// Active item without tap handler: plain text
			children = append(children,
				widgets.Text(lbl).
					Color(render.NewColorFrom(t.TextColor)).
					FontSize(t.FontSizeSM),
				)
		} else if item.OnTap != nil {
			// Clickable item (active or not): Ghost Button
			children = append(children,
				widgets.Button(lbl).
					Variantf(widgets.ButtonGhost).
					OnTap(item.OnTap).
					SetDisabled(false),
			)
		} else {
			// Inactive item without tap handler: muted text
			children = append(children,
				widgets.Text(lbl).
					Color(render.NewColorFrom(t.MutedForegroundColor)).
					FontSize(t.FontSizeSM),
			)
		}
	}
	return widgets.Row(children...).AlignItems(ui.AlignCenter).Gapf(6)
}
