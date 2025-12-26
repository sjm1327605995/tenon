package render

import (
	"image"

	"gioui.org/layout"
)

type Render interface {
	Layout(ctx layout.Context) layout.Dimensions
	DefaultSize() image.Point
	HasDefault() bool
}
type Widget interface {
	ToRender() Render
}
