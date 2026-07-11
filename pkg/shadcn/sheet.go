package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type SheetSide int

const (
	SheetRight SheetSide = iota
	SheetLeft
)

type SheetProps struct {
	Open     bool
	OnClose  func()
	Side     SheetSide
	Width    float32
	children []*ui.Node
}

// Sheet 是从屏幕边缘滑入的抽屉面板。
func Sheet(p SheetProps, children ...*ui.Node) *ui.Node {
	p.children = children
	return ui.Use(sheet, p)
}

func sheet(p SheetProps) *ui.Node {
	th := ui.UseTheme()
	mounted, prog := ui.UseTransition(p.Open, 220)
	ui.UseEscape(p.Open, func() {
		if p.OnClose != nil {
			p.OnClose()
		}
	})
	if !mounted {
		return nil
	}
	w := p.Width
	if w <= 0 {
		w = 360
	}
	tx := (1 - prog) * w
	justify := ui.JustifyEnd
	if p.Side == SheetLeft {
		justify = ui.JustifyStart
		tx = -tx
	}

	panel := append([]ui.StyleOpt{
		ui.Width(w), ui.Column, ui.Gap(12), ui.Padding(24), ui.TranslateXY(tx, 0),
		ui.Bg(th.Background), ui.TextColor(th.Foreground), ui.Border(1, th.Border),
	}, nil...)

	return ui.Portal(
		ui.TrapFocus(), // 模态：键盘焦点限制在抽屉内
		ui.Div(
			ui.Style(ui.Grow(1), ui.Row, justify, ui.Bg(ui.Color{R: 0, G: 0, B: 0, A: uint8(140 * prog)})),
			ui.OnClick(func() {
				if p.OnClose != nil {
					p.OnClose()
				}
			}),
			ui.Div(append([]*ui.Node{ui.Style(panel...), ui.OnClick(func() {})}, p.children...)...),
		),
	)
}

// SheetTitle / SheetDescription 便捷文本。
func SheetTitle(text string) *ui.Node { return ui.Text(text, ui.FontSize(18)) }

func SheetDescription(text string) *ui.Node { return ui.Use(mutedText, text) }
