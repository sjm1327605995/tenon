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
	// shadcn v4: list bg-muted rounded-lg p-[3px]；trigger rounded-md, active bg-background shadow-sm。
	kids := []*ui.Node{ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(4), ui.Padding(3),
		ui.Radius(radiusLg(th)), ui.Bg(th.Muted)),
		ui.ArrowNav(ui.NavHorizontal)} // 左右方向键在标签间移动焦点（WAI-ARIA tabs）
	for i, label := range p.Tabs {
		active := i == p.Active
		idx := i
		st := []ui.StyleOpt{ui.PaddingXY(10, 6), ui.Radius(radiusMd(th)), ui.ItemsCenter, ui.JustifyCenter}
		fg := th.MutedForeground
		if active {
			st = append(st, ui.Bg(th.Background), shadowSm())
			fg = th.Foreground
		}
		kids = append(kids, ui.Button(
			ui.Style(st...),
			ui.OnClick(func() {
				if p.OnChange != nil {
					p.OnChange(idx)
				}
			}),
			ui.Text(label, ui.FontSize(14), ui.TextColor(fg)),
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
	// shadcn v4: h-9 min-w-9 rounded-md px-2；hover bg-muted；on bg-accent text-accent-foreground。
	st := []ui.StyleOpt{ui.Height(36), ui.MinWidth(36), ui.PaddingXY(8, 0), ui.Radius(radiusMd(th)),
		ui.ItemsCenter, ui.JustifyCenter, ui.TextColor(th.Foreground)}
	switch {
	case p.Pressed:
		st = append(st, ui.Bg(th.Accent), ui.TextColor(th.AccentForeground))
	case hovered:
		st = append(st, ui.Bg(th.Muted))
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
	// shadcn v4 outline: 连成一体的带边框分组，选中项 bg-accent。
	kids := []*ui.Node{ui.Style(ui.Row, ui.ItemsCenter, ui.Radius(radiusMd(th)),
		ui.Border(1, th.Border), ui.Clip)}
	for _, opt := range p.Options {
		o := opt
		active := opt == p.Value
		st := []ui.StyleOpt{ui.Height(36), ui.PaddingXY(12, 0), ui.ItemsCenter, ui.JustifyCenter,
			ui.TextColor(th.Foreground)}
		if active {
			st = append(st, ui.Bg(th.Accent), ui.TextColor(th.AccentForeground))
		}
		kids = append(kids, ui.Button(
			ui.Style(st...),
			ui.OnClick(func() {
				if p.OnChange != nil {
					p.OnChange(o)
				}
			}),
			ui.Text(opt, ui.FontSize(14)),
		))
	}
	return ui.Div(kids...)
}
