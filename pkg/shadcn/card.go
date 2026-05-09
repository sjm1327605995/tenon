package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
)

// ShadcnCard renders a content container with optional header/title/description/footer.
func ShadcnCard(title, desc string, children, footer []engine.Widget) engine.Widget {
	t := engine.GetTheme()
	var content []engine.Widget
	if title != "" || desc != "" {
		headerChildren := []engine.Widget{}
		if title != "" {
			headerChildren = append(headerChildren, widgets.Text(title).FontSize(16).Color(render.NewColorFrom(t.CardForegroundColor)))
		}
		if desc != "" {
			headerChildren = append(headerChildren, widgets.Text(desc).FontSize(14).Color(render.NewColorFrom(t.TextMutedColor)))
		}
		content = append(content, widgets.Column(headerChildren...).Gapf(4).Paddingf(engine.EdgeInsets{Top: 24, Right: 24, Bottom: 0, Left: 24}))
	}
	if len(children) > 0 {
		content = append(content, widgets.Column(children...).Gapf(8).Paddingf(engine.EdgeInsetsAll(24)))
	}
	if len(footer) > 0 {
		content = append(content, widgets.Row(footer...).Gapf(8).JustifyContent(engine.JustifyFlexEnd).Paddingf(engine.EdgeInsets{Top: 0, Right: 24, Bottom: 24, Left: 24}))
	}
	return widgets.Container(
		widgets.Column(content...),
	).Background(colorToRender(t.CardColor)).Border(colorToRender(t.BorderColor), 1).Radius(t.BorderRadius)
}
