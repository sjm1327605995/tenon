package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// ShadcnTable renders a simple data table.
func ShadcnTable(headers []string, rows [][]string) ui.Widget {
	t := ui.GetTheme()

	headerCells := make([]ui.Widget, len(headers))
	for i, h := range headers {
		headerCells[i] = widgets.Container(
			widgets.Text(h).Color(render.NewColorFrom(t.TextColor)).FontSize(12),
		).Pad(ui.EdgeInsetsAll(12)).WPct(100.0 / float32(len(headers)))
	}
	headerRow := widgets.Container(
		widgets.Row(headerCells...).AlignItems(ui.AlignStretch),
	).Background(colorToRender(t.MutedColor)).Border(colorToRender(t.BorderColor), 1).WPct(100)

	dataRows := make([]ui.Widget, len(rows))
	for r, row := range rows {
		cells := make([]ui.Widget, len(row))
		for c, text := range row {
			cells[c] = widgets.Container(
				widgets.Text(text).Color(render.NewColorFrom(t.TextColor)).FontSize(14),
			).Pad(ui.EdgeInsetsAll(12)).WPct(100.0 / float32(len(headers)))
		}
		dataRows[r] = widgets.Container(
			widgets.Row(cells...).AlignItems(ui.AlignStretch),
		).Border(colorToRender(t.BorderColor), 1).WPct(100)
	}

	children := make([]ui.Widget, 0, len(dataRows)+1)
	children = append(children, headerRow)
	children = append(children, dataRows...)

	return widgets.Column(children...)
}
