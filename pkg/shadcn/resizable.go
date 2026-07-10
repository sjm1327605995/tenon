package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type ResizableProps struct {
	Width, Height float32
	Left, Right   *ui.Node
}

// Resizable 是左右两栏 + 可拖动分隔条，拖动调整两栏占比。
func Resizable(p ResizableProps) *ui.Node { return ui.Use(resizable, p) }

func resizable(p ResizableProps) *ui.Node {
	th := ui.UseTheme()
	ratio, setRatio := ui.UseState(float32(0.5))
	leftW := clampf(ratio, 0.15, 0.85) * p.Width

	return ui.Div(
		ui.Style(ui.Row, ui.Width(p.Width), ui.Height(p.Height),
			ui.Border(1, th.Border), ui.Radius(th.Radius), ui.Clip),
		ui.Div(ui.Style(ui.Width(leftW), ui.Height(p.Height), ui.Clip), p.Left),
		ui.Div(
			ui.Style(ui.Width(6), ui.Height(p.Height), ui.Bg(th.Border)),
			ui.OnDrag(func(dx, _ float32) {
				setRatio(clampf(ratio+dx/p.Width, 0.15, 0.85))
			}),
		),
		ui.Div(ui.Style(ui.Grow(1), ui.Height(p.Height), ui.Clip), p.Right),
	)
}
