package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// ShadcnTextarea renders a multi-line text input with vertical scrolling.
func ShadcnTextarea(initial string, onChange func(string)) ui.Widget {
	t := ui.GetTheme()
	return widgets.Textarea(initial).
		Placeholder("Type something...").
		Bg(render.Color{R: 255, G: 255, B: 255, A: 255}).
		Border(colorToRender(t.InputBorderColor), 1).
		FocusBorder(colorToRender(t.InputFocusBorderColor)).
		Radius(t.InputBorderRadius).
		Pad(ui.EdgeInsetsAll(12)).
		H(80).
		OnChange(onChange)
}
