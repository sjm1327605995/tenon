package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type DialogProps struct {
	Open     bool
	OnClose  func()
	children []*ui.Node
}

// Dialog 是模态对话框：Open 为真时通过 Portal 渲染，带进出场动画；点击遮罩关闭。
func Dialog(p DialogProps, children ...*ui.Node) *ui.Node {
	p.children = children
	return ui.Use(dialog, p)
}

func dialog(p DialogProps) *ui.Node {
	th := ui.UseTheme()
	mounted, prog := ui.UseTransition(p.Open, 180)
	ui.UseEscape(p.Open, func() {
		if p.OnClose != nil {
			p.OnClose()
		}
	})
	if !mounted {
		return nil
	}
	card := append([]*ui.Node{
		ui.Style(ui.Column, ui.Gap(12), ui.Padding(24), ui.MinWidth(320),
			ui.Bg(th.Background), ui.TextColor(th.Foreground), ui.Border(1, th.Border),
			ui.Radius(th.Radius+4), ui.Scale(0.96+0.04*prog)),
		ui.OnClick(func() {}), // 吞掉点击，避免冒泡关闭
	}, p.children...)
	return ui.Portal(
		ui.TrapFocus(), // 模态：键盘焦点限制在对话框内
		ui.Div(
			ui.Style(ui.Grow(1), ui.ItemsCenter, ui.JustifyCenter,
				ui.Bg(ui.Color{R: 0, G: 0, B: 0, A: 140}), ui.Opacity(prog)),
			ui.OnClick(func() {
				if p.OnClose != nil {
					p.OnClose()
				}
			}),
			ui.Div(card...),
		),
	)
}

// DialogTitle / DialogDescription 便捷文本。
func DialogTitle(text string) *ui.Node { return ui.Text(text, ui.FontSize(18), ui.Semibold) }

func DialogDescription(text string) *ui.Node { return ui.Use(mutedText, text) }
