package render

import (
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/widget"
)

type Clickable struct {
	click widget.Clickable
	Render
}

func (c *Clickable) Layout(ctx layout.Context) layout.Dimensions {
	return c.clickable(ctx, &c.click, func(gtx layout.Context) layout.Dimensions {
		return c.Render.Layout(gtx)
	})

}
func WrapperClickable(render Render) *Clickable {
	return &Clickable{Render: render}
}
func (c *Clickable) clickable(gtx layout.Context, button *widget.Clickable, w layout.Widget) layout.Dimensions {
	return button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Button.Add(gtx.Ops)
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
				if button.Pressed() {
					c.Render.OnClick()
				}
				//for _, c := range button.History() {
				//	drawInk(gtx, c)
				//}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			w,
		)
	})
}
