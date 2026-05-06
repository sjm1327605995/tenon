package shadcn

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// ShadcnCalendar renders a month calendar grid.
func ShadcnCalendar(year int, month time.Month, selected time.Time, onSelect func(time.Time)) ui.Widget {
	t := ui.GetTheme()
	title := fmt.Sprintf("%s %d", month.String(), year)

	// Header
	header := widgets.Row(
		widgets.Button(widgets.IconText(widgets.IconArrowLeft)).Variantf(widgets.ButtonGhost).OnTap(func() {
			if onSelect != nil {
				prev := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).AddDate(0, -1, 0)
				onSelect(prev)
			}
		}),
		widgets.Container(widgets.Text(title).Color(render.NewColorFrom(t.TextColor)).FontSize(16)).Grow(1),
		widgets.Button(widgets.IconText(widgets.IconArrowRight)).Variantf(widgets.ButtonGhost).OnTap(func() {
			if onSelect != nil {
				next := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
				onSelect(next)
			}
		}),
	).AlignItems(ui.AlignCenter).JustifyContent(ui.JustifySpaceBetween).Paddingf(ui.EdgeInsetsOnly(0, 4, 0, 4))

	// Weekday labels
	weekdays := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	dayLabels := make([]ui.Widget, 7)
	for i, d := range weekdays {
		dayLabels[i] = widgets.Container(
			widgets.Text(d).Color(render.NewColorFrom(t.MutedForegroundColor)).FontSize(12),
		).W(40).H(36)
	}
	labelRow := widgets.Row(dayLabels...).AlignItems(ui.AlignCenter)

	// Day grid
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startOffset := int(firstDay.Weekday())
	daysInMonth := 32 - time.Date(year, month, 32, 0, 0, 0, 0, time.Local).Day()

	weekRows := make([]ui.Widget, 0, 6)
	for w := 0; w < 6; w++ {
		rowVisible := false
		dayCells := make([]ui.Widget, 7)
		for d := 0; d < 7; d++ {
			cellDay := w*7 + d - startOffset + 1
			if cellDay >= 1 && cellDay <= daysInMonth {
				rowVisible = true
				day := cellDay
				isSelected := cellDay == selected.Day() && month == selected.Month() && year == selected.Year()
				var bg, fg *render.Color
				var borderW float32
				if isSelected {
					bg = newColor(t.PrimaryColor)
					fg = newColor(t.PrimaryForegroundColor)
				} else {
					bg = nil
					fg = newColor(t.TextColor)
					borderW = 0
				}
				cell := widgets.Container(
					widgets.Row(
						widgets.Text(fmt.Sprintf("%d", day)).FontSize(14).Color(fg),
					).JustifyContent(ui.JustifyCenter).AlignItems(ui.AlignCenter),
				).W(40).H(36).JustifyContent(ui.JustifyCenter).AlignItems(ui.AlignCenter).Border(colorToRender(t.BorderColor), borderW).Radius(t.BorderRadius).OnTap(func() {
					if onSelect != nil {
						onSelect(time.Date(year, month, day, 0, 0, 0, 0, time.Local))
					}
				})
				if bg != nil {
					cell = cell.Background(*bg)
				}
				dayCells[d] = cell
			} else {
				dayCells[d] = widgets.Container(widgets.Text("")).W(40).H(36).JustifyContent(ui.JustifyCenter).AlignItems(ui.AlignCenter)
			}
		}
		if rowVisible {
			weekRows = append(weekRows, widgets.Row(dayCells...).AlignItems(ui.AlignCenter).Gapf(2))
		}
	}

	children := make([]ui.Widget, 0, len(weekRows)+2)
	children = append(children, header, labelRow)
	children = append(children, weekRows...)

	return widgets.Container(
		widgets.Column(children...).Gapf(4),
	).W(294)
}
