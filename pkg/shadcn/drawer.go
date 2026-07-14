package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type DrawerProps struct {
	Open     bool
	OnClose  func()
	Height   float32 // 抽屉高度（默认 320）
	children []*ui.Node
}

// Drawer 是从底部滑入的抽屉（vaul 风格）：顶部有抓手，点遮罩/Esc 关闭。
func Drawer(p DrawerProps, children ...*ui.Node) *ui.Node {
	p.children = children
	return ui.Use(drawer, p)
}

func drawer(p DrawerProps) *ui.Node {
	th := ui.UseTheme()
	mounted, prog := ui.UseTransition(p.Open, 240)
	ui.UseEscape(p.Open, func() {
		if p.OnClose != nil {
			p.OnClose()
		}
	})
	if !mounted {
		return nil
	}
	h := p.Height
	if h <= 0 {
		h = 320
	}
	ty := (1 - prog) * h // 从底部滑上来

	handle := ui.Div(ui.Style(ui.Row, ui.JustifyCenter, ui.PaddingXY(0, 6)),
		ui.Div(ui.Style(ui.Width(40), ui.Height(4), ui.Radius(2), ui.Bg(th.Border))))

	panel := []ui.StyleOpt{
		ui.WidthPct(100), ui.Height(h), ui.Column, ui.Gap(12), ui.PaddingXY(20, 12),
		ui.TranslateXY(0, ty), ui.Bg(th.Background), ui.TextColor(th.Foreground),
		ui.Border(1, th.Border), ui.Radius(th.Radius + 6),
	}
	content := append([]*ui.Node{ui.Style(panel...), ui.OnClick(func() {}), handle}, p.children...)

	return ui.Portal(
		ui.TrapFocus(),
		ui.Div(
			ui.Style(ui.Grow(1), ui.Column, ui.JustifyEnd, ui.Bg(ui.Color{A: uint8(140 * prog)})),
			ui.OnClick(func() {
				if p.OnClose != nil {
					p.OnClose()
				}
			}),
			ui.Div(content...),
		),
	)
}
