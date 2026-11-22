package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type SelectOption struct {
	Label string
	Value string
}

type SelectProps struct {
	Value       string
	Options     []string
	Items       []SelectOption
	OnChange    func(string)
	Placeholder string
	Style       []ui.StyleOpt
}

// Select 是下拉选择框：点击展开选项，选中后回填。
func Select(p SelectProps) *ui.Node { return ui.Use(selectC, p) }

func selectC(p SelectProps) *ui.Node {
	th := ui.UseTheme()
	open, setOpen := ui.UseState(false)
	value, setValue := ui.UseState(p.Value)
	ui.UseEffect(func() ui.Cleanup {
		setValue(p.Value)
		return nil
	}, p.Value)
	ref, rect := ui.UseMeasure()
	ui.UseEscape(open, func() { setOpen(false) })

	items := append([]SelectOption(nil), p.Items...)
	if len(items) == 0 {
		for _, option := range p.Options {
			items = append(items, SelectOption{Label: option, Value: option})
		}
	}
	label := value
	for _, option := range items {
		if option.Value == value {
			label = option.Label
			break
		}
	}
	labelColor := th.Foreground
	if label == "" {
		label, labelColor = p.Placeholder, th.MutedForeground
	}

	style := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.Gap(8), ui.Height(38),
		ui.PaddingXY(12, 0), ui.Radius(th.Radius), ui.Border(1, th.Input),
		ui.Bg(th.Background), ui.MinWidth(180)}
	style = append(style, p.Style...)
	trigger := ui.Div(ref,
		ui.Style(style...),
		ui.OnClick(func() { setOpen(!open) }),
		ui.Text(label, ui.FontSize(14), ui.TextColor(labelColor)),
		ui.Icon(ui.IconChevronDown, 16, ui.TextColor(th.MutedForeground)),
	)

	rows := make([]*ui.Node, 0, len(items))
	for _, option := range items {
		option := option
		rows = append(rows, comboRow(option.Label, option.Value == value, func() {
			setValue(option.Value)
			if p.OnChange != nil {
				p.OnChange(option.Value)
			}
			setOpen(false)
		}))
	}

	return ui.Fragment(
		trigger,
		ui.If(open, floatPanel(th, rect, func() { setOpen(false) },
			rect.W, []ui.StyleOpt{ui.Column, ui.Padding(4)}, rows...)),
	)
}
