package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// ---- Slider ----

type SliderProps struct {
	Value, Min, Max float32
	OnChange        func(float32)
}

func Slider(p SliderProps) *ui.Node { return ui.Use(slider, p) }

func slider(p SliderProps) *ui.Node {
	th := ui.UseTheme()
	w := float32(220)
	span := p.Max - p.Min
	if span <= 0 {
		span = 1
	}
	fill := clampf((p.Value-p.Min)/span, 0, 1) * w
	return ui.Div(
		ui.Style(ui.Width(w), ui.Height(20), ui.JustifyStart, ui.ItemsCenter),
		ui.Div(ui.Style(ui.Width(w), ui.Height(6), ui.Radius(3), ui.Bg(th.Secondary), ui.Absolute, ui.Top(7))),
		ui.Div(ui.Style(ui.Width(fill), ui.Height(6), ui.Radius(3), ui.Bg(th.Primary), ui.Absolute, ui.Top(7))),
		ui.Div(
			ui.Style(ui.Width(18), ui.Height(18), ui.Radius(9), ui.Bg(th.Background),
				ui.Border(2, th.Primary), ui.Absolute, ui.Top(1), ui.Left(fill-9)),
			ui.OnDrag(func(dx, _ float32) {
				if p.OnChange != nil {
					p.OnChange(clampf(p.Value+dx/w*span, p.Min, p.Max))
				}
			}),
		),
	)
}

// ---- Progress ----

func Progress(value float32) *ui.Node { return ui.Use(progress, value) }

func progress(value float32) *ui.Node {
	th := ui.UseTheme()
	value = clampf(value, 0, 1)
	w := float32(240)
	return ui.Div(
		ui.Style(ui.Width(w), ui.Height(8), ui.Radius(4), ui.Bg(th.Secondary)),
		ui.Div(ui.Style(ui.Width(w*value), ui.Height(8), ui.Radius(4), ui.Bg(th.Primary))),
	)
}
