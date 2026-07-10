package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type CarouselProps struct {
	Slides        []*ui.Node
	Width, Height float32
}

// Carousel 是横向轮播：一次显示一张，切换带滑动动画，配上一/下一与指示点。
func Carousel(p CarouselProps) *ui.Node { return ui.Use(carousel, p) }

func carouselBtn(th ui.Theme, label string, onClick func()) *ui.Node {
	return ui.Div(
		ui.Style(ui.Width(30), ui.Height(30), ui.Radius(15), ui.Border(1, th.Border),
			ui.ItemsCenter, ui.JustifyCenter, ui.Bg(th.Background)),
		ui.OnClick(onClick),
		ui.Text(label, ui.FontSize(16), ui.TextColor(th.Foreground)),
	)
}

func carousel(p CarouselProps) *ui.Node {
	th := ui.UseTheme()
	idx, setIdx := ui.UseState(0)
	n := len(p.Slides)
	if n == 0 {
		return ui.Div()
	}
	x := ui.UseTween(float32(idx)*p.Width, 300, ui.EaseInOut)

	strip := []*ui.Node{ui.Style(ui.Row, ui.Width(float32(n)*p.Width), ui.Height(p.Height),
		ui.TranslateXY(-x, 0))}
	for _, s := range p.Slides {
		strip = append(strip, ui.Div(
			ui.Style(ui.Width(p.Width), ui.Height(p.Height), ui.ItemsCenter, ui.JustifyCenter), s))
	}
	viewport := ui.Div(
		ui.Style(ui.Width(p.Width), ui.Height(p.Height), ui.Clip, ui.Radius(th.Radius), ui.Border(1, th.Border)),
		ui.Div(strip...),
	)

	dots := []*ui.Node{ui.Style(ui.Row, ui.Gap(6), ui.ItemsCenter)}
	for i := 0; i < n; i++ {
		c := th.Muted
		if i == idx {
			c = th.Primary
		}
		dots = append(dots, ui.Div(ui.Style(ui.Width(8), ui.Height(8), ui.Radius(4), ui.Bg(c))))
	}

	return ui.Div(
		ui.Style(ui.Column, ui.Gap(12), ui.ItemsCenter),
		viewport,
		ui.Div(ui.Style(ui.Row, ui.Gap(14), ui.ItemsCenter),
			carouselBtn(th, "‹", func() {
				if idx > 0 {
					setIdx(idx - 1)
				}
			}),
			ui.Div(dots...),
			carouselBtn(th, "›", func() {
				if idx < n-1 {
					setIdx(idx + 1)
				}
			}),
		),
	)
}
