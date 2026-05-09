package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
)

// ShadcnTextarea renders a multi-line text input with vertical scrolling.
func ShadcnTextarea(initial string, onChange func(string)) engine.Widget {
	t := engine.GetTheme()
	return widgets.Textarea(initial).
		Placeholder("Type something...").
		Bg(render.Color{R: 255, G: 255, B: 255, A: 255}).
		Border(colorToRender(t.InputBorderColor), 1).
		FocusBorder(colorToRender(t.InputFocusBorderColor)).
		Radius(t.InputBorderRadius).
		Pad(engine.EdgeInsetsAll(12)).
		H(80).
		OnChange(onChange)
}
