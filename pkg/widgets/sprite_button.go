package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// SpriteButtonWidget 是使用图片作为背景的按钮，支持多状态切换。
// 可配合 NinePatch 实现可缩放的按钮皮肤。
type SpriteButtonWidget struct {
	engine.BaseWidget
	normal   *ebiten.Image
	hover    *ebiten.Image
	pressed  *ebiten.Image
	disabled *ebiten.Image
	slice    render.BorderSlice
	label    string
	onClick  func()
	width    float32
	height   float32
	disabledState bool
	loading       bool
}

// SpriteButton 创建图片按钮。
// normal 为必传，hover/pressed/disabled 为可选；未提供时回退到 normal。
func SpriteButton(normal *ebiten.Image) SpriteButtonWidget {
	return SpriteButtonWidget{normal: normal}
}

func (b SpriteButtonWidget) Hover(img *ebiten.Image) SpriteButtonWidget {
	b.hover = img
	return b
}

func (b SpriteButtonWidget) Pressed(img *ebiten.Image) SpriteButtonWidget {
	b.pressed = img
	return b
}

func (b SpriteButtonWidget) Disabled(img *ebiten.Image) SpriteButtonWidget {
	b.disabled = img
	return b
}

// Slice 设置 NinePatch 切图参数，使按钮背景可缩放。
func (b SpriteButtonWidget) Slice(s render.BorderSlice) SpriteButtonWidget {
	b.slice = s
	return b
}

func (b SpriteButtonWidget) Label(text string) SpriteButtonWidget {
	b.label = text
	return b
}

func (b SpriteButtonWidget) OnTap(fn func()) SpriteButtonWidget {
	b.onClick = fn
	return b
}

func (b SpriteButtonWidget) W(v float32) SpriteButtonWidget {
	b.width = v
	return b
}

func (b SpriteButtonWidget) H(v float32) SpriteButtonWidget {
	b.height = v
	return b
}

func (b SpriteButtonWidget) SetDisabled(v bool) SpriteButtonWidget { return b.Disable(v) }
func (b SpriteButtonWidget) Disable(v bool) SpriteButtonWidget {
	b.disabledState = v
	return b
}

func (b SpriteButtonWidget) SetLoading(v bool) SpriteButtonWidget { return b.Loading(v) }
func (b SpriteButtonWidget) Loading(v bool) SpriteButtonWidget {
	b.loading = v
	return b
}

func (b SpriteButtonWidget) CreateElement() engine.Element {
	return engine.NewSingleChildRenderObjectElement(b)
}

func (b SpriteButtonWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	r := render.NewRenderBox()
	r.SetBackgroundImage(b.normal, b.slice)
	r.StyleSetJustifyContent(engine.JustifyCenter)
	r.StyleSetAlignItems(engine.AlignCenter)
	if b.width > 0 {
		r.StyleSetWidth(b.width)
	}
	if b.height > 0 {
		r.StyleSetHeight(b.height)
	}
	if b.label == "" {
		r.StyleSetWidthPercent(100)
	}

	setState := func(s render.ButtonState) {
		if b.disabledState || b.loading {
			if b.disabled != nil {
				r.SetBackgroundImage(b.disabled, b.slice)
			}
			return
		}
		switch s {
		case render.ButtonStateHover:
			if b.hover != nil {
				r.SetBackgroundImage(b.hover, b.slice)
			}
		case render.ButtonStatePressed:
			if b.pressed != nil {
				r.SetBackgroundImage(b.pressed, b.slice)
			}
		default:
			r.SetBackgroundImage(b.normal, b.slice)
		}
	}

	r.SetOnMouseEnter(func() { setState(render.ButtonStateHover) })
	r.SetOnMouseLeave(func() { setState(render.ButtonStateNormal) })
	r.SetOnMouseDown(func() { setState(render.ButtonStatePressed) })
	r.SetOnMouseUp(func() { setState(render.ButtonStateHover) })
	r.SetOnClick(b.onClick)

	if b.disabledState {
		setState(render.ButtonStateDisabled)
	}

	return r
}

func (b SpriteButtonWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	r := ro.(*render.RenderBox)
	old := oldWidget.(SpriteButtonWidget)

	if old.normal != b.normal || old.hover != b.hover || old.pressed != b.pressed || old.disabled != b.disabled || old.slice != b.slice {
		// 重新设置当前状态对应的图片
		if b.disabledState && b.disabled != nil {
			r.SetBackgroundImage(b.disabled, b.slice)
		} else {
			r.SetBackgroundImage(b.normal, b.slice)
		}
	}
	if old.width != b.width {
		if b.width > 0 {
			r.StyleSetWidth(b.width)
		} else {
			r.StyleSetWidthAuto()
		}
	}
	if old.height != b.height {
		if b.height > 0 {
			r.StyleSetHeight(b.height)
		} else {
			r.StyleSetHeightAuto()
		}
	}
	r.SetOnClick(b.onClick)
}

func (b SpriteButtonWidget) GetChildWidget() engine.Widget {
	if b.label == "" {
		return nil
	}
	return Text(b.label).FontSize(engine.GetTheme().FontSizeBase)
}
