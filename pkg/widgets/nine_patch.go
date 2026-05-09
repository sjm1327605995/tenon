package widgets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// NinePatchWidget 使用九宫格图片渲染可任意缩放的 UI 面板背景。
// 适用于游戏 GUI 中的对话框、按钮背景、血条边框等场景。
type NinePatchWidget struct {
	engine.BaseWidget
	Source    *ebiten.Image
	Slice     render.BorderSlice
	Width     float32
	Height    float32
	TintColor color.Color
}

// NinePatch 创建九宫格图片 Widget。
// slice 定义从源图四边到不拉伸区域的距离（像素）。
func NinePatch(src *ebiten.Image, slice render.BorderSlice) NinePatchWidget {
	return NinePatchWidget{
		Source: src,
		Slice:  slice,
	}
}

func (n NinePatchWidget) W(v float32) NinePatchWidget {
	n.Width = v
	return n
}

func (n NinePatchWidget) H(v float32) NinePatchWidget {
	n.Height = v
	return n
}

func (n NinePatchWidget) Tint(c color.Color) NinePatchWidget {
	n.TintColor = c
	return n
}

func (n NinePatchWidget) CreateElement() engine.Element {
	return engine.NewRenderObjectElement(n)
}

func (n NinePatchWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	r := render.NewRenderNinePatch()
	r.Source = n.Source
	r.Slice = n.Slice
	r.TintColor = n.TintColor
	return r
}

func (n NinePatchWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	r := ro.(*render.RenderNinePatch)
	old, _ := oldWidget.(NinePatchWidget)
	if old.Source != n.Source {
		r.Source = n.Source
	}
	if old.Slice != n.Slice {
		r.Slice = n.Slice
	}
	if !render.ColorEquals(old.TintColor, n.TintColor) {
		r.TintColor = n.TintColor
	}
}
