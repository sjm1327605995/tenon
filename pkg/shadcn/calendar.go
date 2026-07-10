package shadcn

import (
	"fmt"
	"time"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

type CalendarProps struct {
	Value    time.Time // 选中日期（零值表示未选）
	OnChange func(time.Time)
}

// Calendar 是月历日期选择器。
func Calendar(p CalendarProps) *ui.Node { return ui.Use(calendar, p) }

func sameDay(a, b time.Time) bool {
	return !b.IsZero() && a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func calendar(p CalendarProps) *ui.Node {
	th := ui.UseTheme()
	base := p.Value
	if base.IsZero() {
		base = time.Now()
	}
	month, setMonth := ui.UseState(time.Date(base.Year(), base.Month(), 1, 0, 0, 0, 0, time.Local))

	navBtn := func(label string, onClick func()) *ui.Node {
		return ui.Div(ui.Style(ui.Width(28), ui.Height(28), ui.Radius(th.Radius-2),
			ui.ItemsCenter, ui.JustifyCenter), ui.OnClick(onClick),
			ui.Text(label, ui.FontSize(16), ui.TextColor(th.MutedForeground)))
	}
	header := ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.PaddingXY(4, 6)),
		navBtn("‹", func() { setMonth(month.AddDate(0, -1, 0)) }),
		ui.Text(fmt.Sprintf("%d 年 %d 月", month.Year(), int(month.Month())), ui.FontSize(14)),
		navBtn("›", func() { setMonth(month.AddDate(0, 1, 0)) }),
	)

	// 星期表头
	wkCells := []*ui.Node{ui.Style(ui.Row)}
	for _, wd := range []string{"日", "一", "二", "三", "四", "五", "六"} {
		wkCells = append(wkCells, ui.Div(
			ui.Style(ui.Width(34), ui.Height(26), ui.ItemsCenter, ui.JustifyCenter),
			ui.Text(wd, ui.FontSize(12), ui.TextColor(th.MutedForeground))))
	}

	blank := func() *ui.Node { return ui.Div(ui.Style(ui.Width(34), ui.Height(34))) }
	dayCell := func(day int, date time.Time) *ui.Node {
		st := []ui.StyleOpt{ui.Width(34), ui.Height(34), ui.Radius(th.Radius - 2),
			ui.ItemsCenter, ui.JustifyCenter}
		fg := th.Foreground
		if sameDay(date, p.Value) {
			st = append(st, ui.Bg(th.Primary))
			fg = th.PrimaryForeground
		}
		return ui.Div(ui.Style(st...),
			ui.OnClick(func() {
				if p.OnChange != nil {
					p.OnChange(date)
				}
			}),
			ui.Text(fmt.Sprintf("%d", day), ui.FontSize(13), ui.TextColor(fg)))
	}

	startOffset := int(month.Weekday()) // 周日=0
	daysIn := time.Date(month.Year(), month.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()

	var cells []*ui.Node
	for i := 0; i < startOffset; i++ {
		cells = append(cells, blank())
	}
	for d := 1; d <= daysIn; d++ {
		date := time.Date(month.Year(), month.Month(), d, 0, 0, 0, 0, time.Local)
		cells = append(cells, dayCell(d, date))
	}
	for len(cells)%7 != 0 {
		cells = append(cells, blank())
	}

	grid := []*ui.Node{ui.Style(ui.Column), header, ui.Div(wkCells...)}
	for i := 0; i < len(cells); i += 7 {
		row := append([]*ui.Node{ui.Style(ui.Row)}, cells[i:i+7]...)
		grid = append(grid, ui.Div(row...))
	}

	return ui.Div(append([]*ui.Node{
		ui.Style(ui.Column, ui.Padding(12), ui.Border(1, th.Border), ui.Radius(th.Radius),
			ui.Bg(th.Background), ui.TextColor(th.Foreground)),
	}, grid...)...)
}
