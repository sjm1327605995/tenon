package widgets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// ImageWidget 配置图片显示。
type ImageWidget struct {
	ui.BaseWidget
	source       *ebiten.Image
	objectFit    render.ObjectFit
	width        float32
	height       float32
	borderRadius float32
	tintColor    color.Color
}

// Image 创建图片 Widget。
func Image(src *ebiten.Image) ImageWidget {
	return ImageWidget{
		source:    src,
		objectFit: render.ObjectFitCover,
	}
}

// Fit 设置 ObjectFit 模式。
func (i ImageWidget) Fit(fit render.ObjectFit) ImageWidget {
	i.objectFit = fit
	return i
}

// W 设置固定宽度。
func (i ImageWidget) W(v float32) ImageWidget {
	i.width = v
	return i
}

// H 设置固定高度。
func (i ImageWidget) H(v float32) ImageWidget {
	i.height = v
	return i
}

// Radius 设置圆角半径。
func (i ImageWidget) Radius(v float32) ImageWidget {
	i.borderRadius = v
	return i
}

// Tint 设置着色颜色。
func (i ImageWidget) Tint(c color.Color) ImageWidget {
	i.tintColor = c
	return i
}

func (i ImageWidget) CreateElement() ui.Element {
	e := &ImageElement{}
	e.RenderObjectElement.BaseElement.Init(e, i)
	return e
}

// ImageElement 是 ImageWidget 对应的 Element。
type ImageElement struct {
	ui.RenderObjectElement
	ro *render.RenderImage
}

func (e *ImageElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderImage)
	e.RenderObject = e.ro
	e.RenderObjectElement.Mount(parent, slot)
}

func (e *ImageElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderImage()
	applyImageProps(r, e.GetWidget().(ImageWidget))
	return r
}

func (e *ImageElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(ImageWidget)
	old, _ := oldWidget.(ImageWidget)

	if old.source != w.source {
		e.ro.SetSource(w.source)
	}
	if old.objectFit != w.objectFit {
		e.ro.SetObjectFit(w.objectFit)
	}
	if old.borderRadius != w.borderRadius {
		e.ro.SetBorderRadius(w.borderRadius)
	}
	if !render.ColorEquals(old.tintColor, w.tintColor) {
		e.ro.SetTintColor(w.tintColor)
	}
	if old.width != w.width {
		e.ro.SetWidth(w.width)
	}
	if old.height != w.height {
		e.ro.SetHeight(w.height)
	}
}

func applyImageProps(r *render.RenderImage, w ImageWidget) {
	r.SetSource(w.source)
	r.SetObjectFit(w.objectFit)
	r.SetBorderRadius(w.borderRadius)
	r.SetTintColor(w.tintColor)
	r.SetWidth(w.width)
	r.SetHeight(w.height)
}
