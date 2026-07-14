package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

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
	return ui.Spinner(sz, c)
}
