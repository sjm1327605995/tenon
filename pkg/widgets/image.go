package widgets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// ImageWidget 配置图片显示。
type ImageWidget struct {
	engine.BaseWidget
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

func (i ImageWidget) CreateElement() engine.Element {
	return engine.NewRenderObjectElement(i)
}

// CreateRenderObject implements RenderObjectFactory.
func (i ImageWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	r := render.NewRenderImage()
	applyImageProps(r, i)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (i ImageWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	r := ro.(*render.RenderImage)
	old, _ := oldWidget.(ImageWidget)

	if old.source != i.source {
		r.SetSource(i.source)
	}
	if old.objectFit != i.objectFit {
		r.SetObjectFit(i.objectFit)
	}
	if old.borderRadius != i.borderRadius {
		r.SetBorderRadius(i.borderRadius)
	}
	if !render.ColorEquals(old.tintColor, i.tintColor) {
		r.SetTintColor(i.tintColor)
	}
	if old.width != i.width {
		r.SetWidth(i.width)
	}
	if old.height != i.height {
		r.SetHeight(i.height)
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
