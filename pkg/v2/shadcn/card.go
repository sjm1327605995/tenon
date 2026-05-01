package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// ShadcnCard renders a content container with optional header/title/description/footer.
func ShadcnCard(title, desc string, children, footer []ui.Widget) ui.Widget {
	t := ui.GetTheme()
	var content []ui.Widget
	if title != "" || desc != "" {
		headerChildren := []ui.Widget{}
		if title != "" {
			headerChildren = append(headerChildren, widgets.Text(title).FontSize(16).Color(render.NewColorFrom(t.CardForegroundColor)))
		}
		if desc != "" {
			headerChildren = append(headerChildren, widgets.Text(desc).FontSize(14).Color(render.NewColorFrom(t.TextMutedColor)))
		}
		content = append(content, widgets.Column(headerChildren...).Gapf(4).Paddingf(ui.EdgeInsets{Top: 24, Right: 24, Bottom: 0, Left: 24}))
	}
	if len(children) > 0 {
		content = append(content, widgets.Column(children...).Gapf(8).Paddingf(ui.EdgeInsetsAll(24)))
	}
	if len(footer) > 0 {
		content = append(content, widgets.Row(footer...).Gapf(8).JustifyContent(ui.JustifyFlexEnd).Paddingf(ui.EdgeInsets{Top: 0, Right: 24, Bottom: 24, Left: 24}))
	}
	return widgets.Container(
		widgets.Column(content...),
	).Background(colorToRender(t.CardColor)).Border(colorToRender(t.BorderColor), 1).Radius(t.BorderRadius)
}
