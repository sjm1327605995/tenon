package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type BadgeVariant int

const (
	BadgeDefault BadgeVariant = iota
	BadgeSecondary
	BadgeDestructive
	BadgeOutline
)

type BadgeProps struct {
	Variant  BadgeVariant
	children []*ui.Node
}

// Badge 是小圆角标签。
func Badge(p BadgeProps, children ...*ui.Node) *ui.Node {
	p.children = children
	return ui.Use(badge, p)
}

func badge(p BadgeProps) *ui.Node {
	th := ui.UseTheme()
	bg, fg, border := th.Primary, th.PrimaryForeground, ui.Transparent
	bordered := false
	switch p.Variant {
	case BadgeSecondary:
		bg, fg = th.Secondary, th.SecondaryForeground
	case BadgeDestructive:
		bg, fg = th.Destructive, th.DestructiveForeground
	case BadgeOutline:
		bg, fg, border, bordered = ui.Transparent, th.Foreground, th.Border, true
	}
	st := []ui.StyleOpt{
		ui.Row, ui.ItemsCenter, ui.JustifyCenter, ui.PaddingXY(10, 2), ui.Radius(999),
		ui.Bg(bg), ui.TextColor(fg), ui.FontSize(12),
	}
	if bordered {
		st = append(st, ui.Border(1, border))
	}
	return ui.Div(append([]*ui.Node{ui.Style(st...)}, p.children...)...)
}
