package render

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

type Clickable struct {
	click       widget.Clickable
	OnClickFunc func()
	Render
}

func (c *Clickable) Layout(ctx layout.Context) layout.Dimensions {
	return c.click.Layout(ctx, c.Render.Layout)
}
func WrapperClickable(render Render) *Clickable {
	return &Clickable{Render: render}
}
