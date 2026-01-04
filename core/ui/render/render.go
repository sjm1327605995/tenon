package render

import (
	"image"

	"gioui.org/layout"
)

type Render interface {
	Layout(ctx layout.Context) layout.Dimensions
	DefaultSize() image.Point
	HasDefault() bool
	Clickable() bool
}
type Widget interface {
	ToRender() Render
}
type BaseRender struct{}

func (b BaseRender) Layout(ctx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

func (b BaseRender) DefaultSize() image.Point {
	return image.Point{}
}

func (b BaseRender) HasDefault() bool {
	return false
}

func (b BaseRender) Clickable() bool {
	return false
}
