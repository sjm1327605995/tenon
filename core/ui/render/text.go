package render

import (
	"image"
	_ "image/jpeg"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

type TextStyle struct {
	material.LabelStyle
}

func (i *TextStyle) ToRender() Render {
	return i
}

type Text struct {
	LabelStyle material.LabelStyle
}

func (i *TextStyle) DefaultSize() image.Point {
	return image.Point{}
}

func (i *TextStyle) HasDefault() bool {
	return true
}

func (i *TextStyle) Layout(ctx layout.Context) layout.Dimensions {
	return i.LabelStyle.Layout(ctx)
}
