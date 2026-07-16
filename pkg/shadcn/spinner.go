package shadcn

import (
	"math"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

type SpinnerProps struct {
	Size  float32  // 边长（逻辑 px，默认 16）
	Color ui.Color // 颜色（默认 MutedForeground）
}

// Spinner 是加载指示器（持续旋转）。
func Spinner(p SpinnerProps) *ui.Node { return ui.Use(spinner, p) }

func spinner(p SpinnerProps) *ui.Node {
	th := ui.UseTheme()
	sz := p.Size
	if sz <= 0 {
		sz = 16
	}
	c := p.Color
	if (c == ui.Color{}) {
		c = th.MutedForeground
	}

	deg := ui.UseElapsed() * 300 // ~0.83 圈/秒
	dot := sz / 5
	r := sz/2 - dot/2
	kids := []*ui.Node{ui.Style(ui.Width(sz), ui.Height(sz), ui.Rotate(deg))}
	const n = 8
	for i := 0; i < n; i++ {
		a := float64(i) / float64(n) * 2 * math.Pi
		cx := sz/2 + float32(math.Cos(a))*r - dot/2
		cy := sz/2 + float32(math.Sin(a))*r - dot/2
		kids = append(kids, ui.Div(ui.Style(
			ui.Width(dot), ui.Height(dot), ui.Radius(dot/2), ui.Bg(c),
			ui.Opacity(float32(i+1)/float32(n)), ui.Absolute, ui.Left(cx), ui.Top(cy))))
	}
	return ui.Div(kids...)
}
