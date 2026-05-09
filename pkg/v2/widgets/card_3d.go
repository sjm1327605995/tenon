package widgets

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// Card3DWidget 是一个 3D 卡牌组件，支持双面纹理和透视旋转。
// 通过 RotateX / RotateY 实现倾斜和拉伸效果。
type Card3DWidget struct {
	ui.BaseWidget
	front       *ebiten.Image
	back        *ebiten.Image
	rotateX     float32
	rotateY     float32
	perspective float32
	width       float32
	height      float32
}

// Card3D 创建 3D 卡牌；front 为正面纹理，back 为反面纹理（nil 时用正面代替）。
func Card3D(front, back *ebiten.Image) Card3DWidget {
	return Card3DWidget{front: front, back: back}
}

func (c Card3DWidget) RotateX(deg float32) Card3DWidget {
	c.rotateX = deg
	return c
}

func (c Card3DWidget) RotateY(deg float32) Card3DWidget {
	c.rotateY = deg
	return c
}

func (c Card3DWidget) Perspective(dist float32) Card3DWidget {
	c.perspective = dist
	return c
}

func (c Card3DWidget) W(v float32) Card3DWidget {
	c.width = v
	return c
}

func (c Card3DWidget) H(v float32) Card3DWidget {
	c.height = v
	return c
}

func (c Card3DWidget) CreateElement() ui.Element {
	return ui.NewSingleChildRenderObjectElement(c)
}

func (c Card3DWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	img := c.visibleFace()
	r := render.NewRenderImage()
	r.SetSource(img)
	if c.width > 0 {
		r.SetWidth(c.width)
	}
	if c.height > 0 {
		r.SetHeight(c.height)
	}
	t := render.DefaultTransform()
	t.RotateX = c.rotateX
	t.RotateY = c.rotateY
	t.Perspective = c.perspective
	if t.Perspective == 0 {
		t.Perspective = 800 // 默认透视距离
	}
	r.SetTransform(t)
	return r
}

func (c Card3DWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	r := ro.(*render.RenderImage)
	old := oldWidget.(Card3DWidget)
	img := c.visibleFace()
	r.SetSource(img)
	if old.width != c.width && c.width > 0 {
		r.SetWidth(c.width)
	}
	if old.height != c.height && c.height > 0 {
		r.SetHeight(c.height)
	}
	t := r.GetTransform()
	if old.rotateX != c.rotateX || old.rotateY != c.rotateY || old.perspective != c.perspective {
		t.RotateX = c.rotateX
		t.RotateY = c.rotateY
		t.Perspective = c.perspective
		if t.Perspective == 0 {
			t.Perspective = 800
		}
		r.SetTransform(t)
	}
}

func (c Card3DWidget) GetChildWidget() ui.Widget {
	return nil
}

// visibleFace 根据 RotateY 判断当前应显示正面还是反面。
func (c Card3DWidget) visibleFace() *ebiten.Image {
	// 将角度规范化到 0~360
	ry := math.Mod(float64(c.rotateY), 360)
	if ry < 0 {
		ry += 360
	}
	// 90~270 度显示反面
	if ry > 90 && ry < 270 {
		if c.back != nil {
			return c.back
		}
	}
	return c.front
}
