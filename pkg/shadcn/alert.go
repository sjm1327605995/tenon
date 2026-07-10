package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type AlertVariant int

const (
	AlertDefault AlertVariant = iota
	AlertDestructive
)

type AlertProps struct {
	Variant  AlertVariant
	children []*ui.Node
}

// Alert 是提示框，配合 AlertTitle / AlertDescription 使用。
func Alert(p AlertProps, children ...*ui.Node) *ui.Node {
	p.children = children
	return ui.Use(alert, p)
}

func alert(p AlertProps) *ui.Node {
	th := ui.UseTheme()
	fg, border := th.Foreground, th.Border
	if p.Variant == AlertDestructive {
		fg, border = th.Destructive, th.Destructive
	}
	base := ui.Style(ui.Column, ui.Gap(6), ui.Padding(16), ui.Radius(th.Radius),
		ui.Border(1, border), ui.Bg(th.Background), ui.TextColor(fg))
	return ui.Div(append([]*ui.Node{base}, p.children...)...)
}

// AlertTitle 继承 Alert 前景色。
func AlertTitle(text string) *ui.Node { return ui.Text(text, ui.FontSize(15)) }

// AlertDescription 使用弱化前景色。
func AlertDescription(text string) *ui.Node { return ui.Use(mutedText, text) }
