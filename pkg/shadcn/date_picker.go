package shadcn

import (
	"time"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

type DatePickerProps struct {
	Value       time.Time
	OnChange    func(time.Time)
	Placeholder string
}

// DatePicker 是「触发按钮 + 日历浮层」的日期选择器：点击展开日历，选中回填并关闭。
func DatePicker(p DatePickerProps) *ui.Node { return ui.Use(datePicker, p) }

func datePicker(p DatePickerProps) *ui.Node {
	th := ui.UseTheme()
	open, setOpen := ui.UseState(false)
	ref, rect := ui.UseMeasure()
	ui.UseEscape(open, func() { setOpen(false) })

	label, color := p.Placeholder, th.MutedForeground
	if label == "" {
		label = "选择日期"
	}
	if !p.Value.IsZero() {
		label, color = p.Value.Format("2006-01-02"), th.Foreground
	}

	trigger := ui.Div(ref,
		ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.Gap(8), ui.Height(36),
			ui.PaddingXY(12, 0), ui.Radius(radiusMd(th)), ui.Border(1, th.Input),
			ui.Bg(th.Background), ui.MinWidth(200)),
		ui.OnClick(func() { setOpen(!open) }),
		ui.Text(label, ui.FontSize(14), ui.TextColor(color)),
		ui.Icon(ui.IconChevronDown, 16, ui.TextColor(th.MutedForeground)),
	)

	return ui.Fragment(trigger,
		ui.If(open, floatPanel(th, rect, func() { setOpen(false) }, []ui.StyleOpt{ui.Padding(8)},
			Calendar(CalendarProps{Value: p.Value, OnChange: func(t time.Time) {
				if p.OnChange != nil {
					p.OnChange(t)
				}
				setOpen(false)
			}}),
		)),
	)
}
