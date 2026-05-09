package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
)

// BreadcrumbItem is a single entry in a breadcrumb trail.
type BreadcrumbItem struct {
	Label    string
	Href     string // navigation reference (unused in UI)
	IsActive bool
	OnTap    func()
}

// ShadcnBreadcrumb renders a horizontal breadcrumb trail.
// Use engine.GetTheme().TextColor for active, MutedForegroundColor for inactive.
func ShadcnBreadcrumb(items []BreadcrumbItem, separator string) engine.Widget {
	if separator == "" {
		separator = "/"
	}
	t := engine.GetTheme()
	children := make([]engine.Widget, 0, len(items)*2-1)
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
	return widgets.Row(children...).AlignItems(engine.AlignCenter).Gapf(6)
}
