package shadcn

import (
	"fmt"

	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
)

// ShadcnTable renders a simple data table.
func ShadcnTable(headers []string, rows [][]string) engine.Widget {
	t := engine.GetTheme()

	headerCells := make([]engine.Widget, len(headers))
	for i, h := range headers {
		headerCells[i] = widgets.Container(
			widgets.Text(h).Color(render.NewColorFrom(t.TextColor)).FontSize(12),
		).Pad(engine.EdgeInsetsAll(12)).WPct(100.0 / float32(len(headers)))
	}
	headerRow := widgets.Container(
		widgets.Row(headerCells...).AlignItems(engine.AlignStretch),
	).Background(colorToRender(t.MutedColor)).Border(colorToRender(t.BorderColor), 1).WPct(100)

	dataRows := make([]engine.Widget, len(rows))
	for r, row := range rows {
		cells := make([]engine.Widget, len(row))
		for c, text := range row {
			cells[c] = widgets.Container(
				widgets.Text(text).Color(render.NewColorFrom(t.TextColor)).FontSize(14),
			).Pad(engine.EdgeInsetsAll(12)).WPct(100.0 / float32(len(headers)))
		}
		dataRows[r] = widgets.Container(
			widgets.Row(cells...).AlignItems(engine.AlignStretch),
		).Border(colorToRender(t.BorderColor), 1).WPct(100)
	}

	children := make([]engine.Widget, 0, len(dataRows)+1)
	children = append(children, headerRow)
	for i, row := range dataRows {
		if cw, ok := row.(interface{ SetKey(k engine.Key) }); ok {
			cw.SetKey(engine.NewStringKey(fmt.Sprintf("table-row-%d", i)))
		}
		children = append(children, row)
	}

	column := widgets.Column(children...)
	column.SetKey(engine.NewStringKey("shadcn-table"))
	return column
}
