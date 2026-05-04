package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// TextareaWidget 是多行文本输入框，内部使用 RenderScroll 实现垂直滚动。
type TextareaWidget struct {
	ui.BaseWidget
	content          string
	placeholder      string
	fontSize         float32
	width            float32
	height           float32
	background       *render.Color
	borderColor      *render.Color
	focusBorderColor *render.Color
	borderWidth      float32
	borderRadius     float32
	padding          ui.EdgeInsets
	onChanged        func(string)
}

// Textarea 创建多行文本输入框。
func Textarea(initial string) TextareaWidget {
	return TextareaWidget{
		content:      initial,
		fontSize:     14,
		width:        200,
		height:       80,
		background:   render.NewColor(255, 255, 255, 255),
		borderColor:  render.NewColor(200, 200, 200, 255),
		borderWidth:  1,
		borderRadius: 4,
		padding:      ui.EdgeInsets{Left: 8, Right: 8, Top: 8, Bottom: 8},
	}
}

func (t TextareaWidget) Placeholder(text string) TextareaWidget {
	t.placeholder = text
	return t
}

func (t TextareaWidget) Size(size float32) TextareaWidget {
	t.fontSize = size
	return t
}

func (t TextareaWidget) W(v float32) TextareaWidget {
	t.width = v
	return t
}

func (t TextareaWidget) H(v float32) TextareaWidget {
	t.height = v
	return t
}

func (t TextareaWidget) Bg(c render.Color) TextareaWidget {
	t.background = &c
	return t
}

func (t TextareaWidget) Border(c render.Color, width float32) TextareaWidget {
	t.borderColor = &c
	t.borderWidth = width
	return t
}

func (t TextareaWidget) FocusBorder(c render.Color) TextareaWidget {
	t.focusBorderColor = &c
	return t
}

func (t TextareaWidget) Radius(v float32) TextareaWidget {
	t.borderRadius = v
	return t
}

func (t TextareaWidget) Pad(insets ui.EdgeInsets) TextareaWidget {
	t.padding = insets
	return t
}

func (t TextareaWidget) OnChange(fn func(string)) TextareaWidget {
	t.onChanged = fn
	return t
}

func (t TextareaWidget) buildEditable() ui.Widget {
	e := EditableText(t.content).
		Multiline(true).
		Size(t.fontSize)
	if t.placeholder != "" {
		e = e.Placeholder(t.placeholder)
	}
	if t.onChanged != nil {
		e = e.OnChange(t.onChanged)
	}
	return e
}

func (t TextareaWidget) CreateElement() ui.Element {
	el := &TextareaElement{}
	el.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(el, t)
	return el
}

// TextareaElement 是 TextareaWidget 对应的 Element。
type TextareaElement struct {
	ui.SingleChildRenderObjectElement
}

func (e *TextareaElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderScroll()
	w := e.GetWidget().(TextareaWidget)
	applyTextareaProps(r, TextareaWidget{}, w)
	return r
}

func (e *TextareaElement) UpdateRenderObject(oldWidget ui.Widget) {
	r := e.GetRenderObject().(*render.RenderScroll)
	old := oldWidget.(TextareaWidget)
	w := e.GetWidget().(TextareaWidget)
	applyTextareaProps(r, old, w)
}

func (e *TextareaElement) UpdateChild(oldWidget ui.Widget) {
	w := e.GetWidget().(TextareaWidget)
	e.Child = ui.UpdateChild(e, e.Child, w.buildEditable())
}

func (e *TextareaElement) Mount(parent ui.Element, slot int) {
	e.RenderObject = e.CreateRenderObject()
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(TextareaWidget)
	e.Child = ui.UpdateChild(e, nil, w.buildEditable())
}

func applyTextareaProps(r *render.RenderScroll, old, w TextareaWidget) {
	if !render.ColorPtrEquals(old.background, w.background) {
		r.SetBackgroundColor(w.background)
	}
	if !render.ColorPtrEquals(old.borderColor, w.borderColor) {
		r.SetBorderColor(w.borderColor)
	}
	if !render.ColorPtrEquals(old.focusBorderColor, w.focusBorderColor) {
		r.SetFocusedBorderColor(w.focusBorderColor)
	}
	if old.borderWidth != w.borderWidth {
		r.SetBorderWidth(w.borderWidth)
	}
	if old.borderRadius != w.borderRadius {
		r.SetBorderRadius(w.borderRadius)
	}
	r.SetClipChildren(true)

	if old.width != w.width {
		if w.width > 0 {
			r.StyleSetWidth(w.width)
		} else {
			r.StyleSetWidthAuto()
		}
	}
	if old.height != w.height {
		if w.height > 0 {
			r.StyleSetHeight(w.height)
		} else {
			r.StyleSetHeightAuto()
		}
	}
	if old.padding != w.padding {
		r.StyleSetPadding(ui.EdgeTop, w.padding.Top)
		r.StyleSetPadding(ui.EdgeRight, w.padding.Right)
		r.StyleSetPadding(ui.EdgeBottom, w.padding.Bottom)
		r.StyleSetPadding(ui.EdgeLeft, w.padding.Left)
	}
}
