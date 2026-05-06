package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// AlertVariant defines the visual style of an alert.
type AlertVariant int

const (
	AlertDefault AlertVariant = iota
	AlertDestructive
)

// ShadcnAlert renders a callout box for important messages.
func ShadcnAlert(title, desc string, variant AlertVariant) ui.Widget {
	t := ui.GetTheme()
	iconColor := t.PrimaryColor
	titleColor := t.TextColor
	if variant == AlertDestructive {
		iconColor = t.DestructiveColor
		titleColor = t.DestructiveColor
	}

	content := []ui.Widget{
		widgets.Text(title).FontSize(14).Color(render.NewColorFrom(titleColor)),
	}
	if desc != "" {
		content = append(content, widgets.Text(desc).FontSize(14).Color(render.NewColorFrom(t.TextMutedColor)))
	}
	body := widgets.Column(content...).Gapf(4)

	return widgets.Container(
		widgets.Row(
			widgets.Icon(widgets.IconInfo).FontSize(16).Color(render.NewColorFrom(iconColor)),
			body,
		).AlignItems(ui.AlignCenter).Gapf(12),
	).Background(colorToRender(t.CardColor)).Border(colorToRender(t.BorderColor), 1).Radius(t.BorderRadius).Pad(ui.EdgeInsetsAll(16)).WPct(100)
}
