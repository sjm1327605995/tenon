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
	// shadcn v4: rounded-lg border px-4 py-3 text-sm bg-card；destructive 用 text-destructive。
	fg, border := th.CardForeground, th.Border
	if p.Variant == AlertDestructive {
		fg, border = th.Destructive, th.Destructive
	}
	base := ui.Style(ui.Column, ui.Gap(4), ui.PaddingXY(16, 12), ui.Radius(radiusLg(th)),
		ui.Border(1, border), ui.Bg(th.Card), ui.TextColor(fg))
	return ui.Div(append([]*ui.Node{base}, p.children...)...)
}

// AlertTitle 继承 Alert 前景色（text-sm font-medium）。
func AlertTitle(text string) *ui.Node { return ui.Text(text, ui.FontSize(14)) }

// AlertDescription 使用弱化前景色。
func AlertDescription(text string) *ui.Node { return ui.Use(mutedText, text) }
