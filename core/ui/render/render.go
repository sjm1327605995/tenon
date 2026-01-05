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
	OnClick()
}
type Widget interface {
	ToRender() Render
}
type BaseRender struct {
	onClickFunc func()
}

func (b *BaseRender) Layout(ctx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

func (b *BaseRender) DefaultSize() image.Point {
	return image.Point{}
}

func (b *BaseRender) HasDefault() bool {
	return false
}

func (b *BaseRender) Clickable() bool {
	return b.onClickFunc != nil
}
func (b *BaseRender) SetOnClick(onClick func()) {
	b.onClickFunc = onClick
}
func (b *BaseRender) OnClick() {
	b.onClickFunc()
}
