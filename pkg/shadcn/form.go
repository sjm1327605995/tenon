package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// ---- Checkbox ----

type CheckboxProps struct {
	Checked  bool
	OnChange func(bool)
	Disabled bool
}

func Checkbox(p CheckboxProps) *ui.Node { return ui.Use(checkbox, p) }

func checkbox(p CheckboxProps) *ui.Node {
	th := ui.UseTheme()
	st := []ui.StyleOpt{ui.Width(18), ui.Height(18), ui.Radius(4), ui.ItemsCenter,
		ui.JustifyCenter, ui.Border(1, th.Primary)}
	if p.Checked {
		st = append(st, ui.Bg(th.Primary))
	}
	if p.Disabled {
		st = append(st, ui.Opacity(0.5))
	}
	kids := []*ui.Node{ui.Style(st...)}
	if !p.Disabled {
		kids = append(kids, ui.OnClick(func() {
			if p.OnChange != nil {
				p.OnChange(!p.Checked)
			}
		}))
	}
	if p.Checked {
		kids = append(kids, ui.Text("✓", ui.FontSize(13), ui.TextColor(th.PrimaryForeground)))
	}
	return ui.Div(kids...)
}

// ---- Switch ----

type SwitchProps struct {
	Checked  bool
	OnChange func(bool)
	Disabled bool
}

func Switch(p SwitchProps) *ui.Node { return ui.Use(switchC, p) }

func switchC(p SwitchProps) *ui.Node {
	th := ui.UseTheme()
	target := float32(0)
	if p.Checked {
		target = 1
	}
	x := ui.UseTween(target, 140, ui.EaseOut)
	st := []ui.StyleOpt{ui.Width(40), ui.Height(22), ui.Radius(11), ui.JustifyStart,
		ui.ItemsCenter, ui.Bg(ui.Mix(th.Input, th.Primary, x))}
	if p.Disabled {
		st = append(st, ui.Opacity(0.5))
	}
	attrs := []*ui.Node{ui.Style(st...)}
	if !p.Disabled {
		attrs = append(attrs, ui.OnClick(func() {
			if p.OnChange != nil {
				p.OnChange(!p.Checked)
			}
		}))
	}
	thumb := ui.Div(ui.Style(ui.Width(16), ui.Height(16), ui.Radius(8), ui.Bg(th.Background),
		ui.Absolute, ui.Top(3), ui.Left(3+x*18)))
	return ui.Div(append(attrs, thumb)...)
}

// ---- RadioGroup ----

type RadioGroupProps struct {
	Value    string
	Options  []string
	OnChange func(string)
}

func RadioGroup(p RadioGroupProps) *ui.Node { return ui.Use(radioGroup, p) }

func radioGroup(p RadioGroupProps) *ui.Node {
	th := ui.UseTheme()
	kids := []*ui.Node{ui.Style(ui.Column, ui.Gap(10))}
	for _, opt := range p.Options {
		o := opt
		selected := opt == p.Value
		var dot *ui.Node
		if selected {
			dot = ui.Div(ui.Style(ui.Width(9), ui.Height(9), ui.Radius(5), ui.Bg(th.Primary)))
		}
		kids = append(kids, ui.Div(
			ui.Style(ui.Row, ui.Gap(8), ui.ItemsCenter),
			ui.OnClick(func() {
				if p.OnChange != nil {
					p.OnChange(o)
				}
			}),
			ui.Div(ui.Style(ui.Width(18), ui.Height(18), ui.Radius(9), ui.Border(1, th.Primary),
				ui.ItemsCenter, ui.JustifyCenter), dot),
			ui.Text(opt, ui.FontSize(14)),
		))
	}
	return ui.Div(kids...)
}
