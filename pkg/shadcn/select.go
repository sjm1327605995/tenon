package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type SelectProps struct {
	Value       string
	Options     []string
	OnChange    func(string)
	Placeholder string
}

// Select 是下拉选择框：点击展开选项，选中后回填。
func Select(p SelectProps) *ui.Node { return ui.Use(selectC, p) }

func selectC(p SelectProps) *ui.Node {
	th := ui.UseTheme()
	open, setOpen := ui.UseState(false)
	ref, rect := ui.UseMeasure()
	ui.UseEscape(open, func() { setOpen(false) })

	label, labelColor := p.Value, th.Foreground
	if label == "" {
		label, labelColor = p.Placeholder, th.MutedForeground
	}

	trigger := ui.Div(ref,
		ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.Gap(8), ui.Height(38),
			ui.PaddingXY(12, 0), ui.Radius(th.Radius), ui.Border(1, th.Input),
			ui.Bg(th.Background), ui.MinWidth(180)),
		ui.OnClick(func() { setOpen(!open) }),
		ui.Text(label, ui.FontSize(14), ui.TextColor(labelColor)),
		ui.Text("▾", ui.FontSize(12), ui.TextColor(th.MutedForeground)),
	)

	rows := make([]*ui.Node, 0, len(p.Options))
	for _, opt := range p.Options {
		o := opt
		rows = append(rows, menuRow(o, func() {
			if p.OnChange != nil {
				p.OnChange(o)
			}
			setOpen(false)
		}))
	}

	return ui.Fragment(
		trigger,
		ui.If(open, floatPanel(th, rect, func() { setOpen(false) },
			[]ui.StyleOpt{ui.Column, ui.Padding(4), ui.MinWidth(rect.W)}, rows...)),
	)
}
