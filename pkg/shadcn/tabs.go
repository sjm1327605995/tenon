package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// ---- Tabs ----

type TabsProps struct {
	Tabs     []string
	Active   int
	OnChange func(int)
}

// Tabs 渲染分段式标签栏（内容区由外部按 Active 渲染）。
func Tabs(p TabsProps) *ui.Node { return ui.Use(tabs, p) }

func tabs(p TabsProps) *ui.Node {
	th := ui.UseTheme()
	kids := []*ui.Node{ui.Style(ui.Row, ui.Gap(4), ui.Padding(4), ui.Radius(th.Radius), ui.Bg(th.Muted))}
	for i, label := range p.Tabs {
		active := i == p.Active
		idx := i
		st := []ui.StyleOpt{ui.PaddingXY(14, 6), ui.Radius(th.Radius - 2), ui.ItemsCenter, ui.JustifyCenter}
		fg := th.MutedForeground
		if active {
			st = append(st, ui.Bg(th.Background))
			fg = th.Foreground
		}
		kids = append(kids, ui.Button(
			ui.Style(st...),
			ui.OnClick(func() {
				if p.OnChange != nil {
					p.OnChange(idx)
				}
			}),
			ui.Text(label, ui.FontSize(13), ui.TextColor(fg)),
		))
	}
	return ui.Div(kids...)
}

// ---- Toggle ----

type ToggleProps struct {
	Pressed  bool
	OnChange func(bool)
	children []*ui.Node
}

func Toggle(p ToggleProps, children ...*ui.Node) *ui.Node {
	p.children = children
	return ui.Use(toggle, p)
}

func toggle(p ToggleProps) *ui.Node {
	th := ui.UseTheme()
	hovered, _, ia := ui.UseInteraction()
	st := []ui.StyleOpt{ui.Height(36), ui.PaddingXY(12, 0), ui.Radius(th.Radius),
		ui.ItemsCenter, ui.JustifyCenter, ui.TextColor(th.Foreground)}
	switch {
	case p.Pressed:
		st = append(st, ui.Bg(th.Accent), ui.TextColor(th.AccentForeground))
	case hovered:
		st = append(st, ui.Bg(ui.Mix(th.Background, th.Accent, 0.5)))
	}
	attrs := []*ui.Node{
		ui.Style(st...),
		ui.OnClick(func() {
			if p.OnChange != nil {
				p.OnChange(!p.Pressed)
			}
		}),
		ia,
	}
	return ui.Button(append(attrs, p.children...)...)
}

// ---- ToggleGroup ----

type ToggleGroupProps struct {
	Value    string
	Options  []string
	OnChange func(string)
}

// ToggleGroup 是单选的分段切换组。
func ToggleGroup(p ToggleGroupProps) *ui.Node { return ui.Use(toggleGroup, p) }

func toggleGroup(p ToggleGroupProps) *ui.Node {
	th := ui.UseTheme()
	kids := []*ui.Node{ui.Style(ui.Row, ui.Gap(2), ui.Padding(2), ui.Radius(th.Radius), ui.Bg(th.Muted))}
	for _, opt := range p.Options {
		o := opt
		active := opt == p.Value
		st := []ui.StyleOpt{ui.PaddingXY(14, 6), ui.Radius(th.Radius - 2), ui.ItemsCenter, ui.JustifyCenter}
		if active {
			st = append(st, ui.Bg(th.Background))
		}
		kids = append(kids, ui.Button(
			ui.Style(st...),
			ui.OnClick(func() {
				if p.OnChange != nil {
					p.OnChange(o)
				}
			}),
			ui.Text(opt, ui.FontSize(13), ui.TextColor(th.Foreground)),
		))
	}
	return ui.Div(kids...)
}
