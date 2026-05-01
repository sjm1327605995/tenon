package widgets

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// EditableTextWidget 是可编辑文本 Widget。
type EditableTextWidget struct {
	ui.BaseWidget
	content     string
	fontSize    float32
	textColor   color.Color
	onChanged   func(string)
	onSubmitted func(string)
}

// EditableText 创建可编辑文本 Widget。
func EditableText(content string) EditableTextWidget {
	return EditableTextWidget{
		content:   content,
		fontSize:  14,
		textColor: color.Black,
	}
}

func (e EditableTextWidget) Size(size float32) EditableTextWidget {
	e.fontSize = size
	return e
}

func (e EditableTextWidget) Color(c color.Color) EditableTextWidget {
	e.textColor = c
	return e
}

func (e EditableTextWidget) OnChange(fn func(string)) EditableTextWidget {
	e.onChanged = fn
	return e
}

func (e EditableTextWidget) OnSubmit(fn func(string)) EditableTextWidget {
	e.onSubmitted = fn
	return e
}

func (e EditableTextWidget) CreateElement() ui.Element {
	el := &EditableTextElement{}
	el.RenderObjectElement.BaseElement.Init(el, e)
	return el
}

// EditableTextElement 是 EditableTextWidget 对应的 Element。
type EditableTextElement struct {
	ui.RenderObjectElement
}

func (e *EditableTextElement) Mount(parent ui.Element, slot int) {
	e.RenderObject = e.CreateRenderObject()
	e.RenderObjectElement.Mount(parent, slot)
}

func (e *EditableTextElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(EditableTextWidget)
	r := render.NewRenderEditableText()
	r.SetContent(w.content)
	r.SetFontSize(w.fontSize)
	r.SetColor(w.textColor)
	r.SetOnChanged(w.onChanged)
	r.SetOnSubmitted(w.onSubmitted)
	return r
}

func (e *EditableTextElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(EditableTextWidget)
	r := e.GetRenderObject().(*render.RenderEditableText)
	old := oldWidget.(EditableTextWidget)

	if old.fontSize != w.fontSize {
		r.SetFontSize(w.fontSize)
	}
	if !render.ColorEquals(old.textColor, w.textColor) {
		r.SetColor(w.textColor)
	}
	// onChanged / onSubmitted are function pointers; update unconditionally
	r.SetOnChanged(w.onChanged)
	r.SetOnSubmitted(w.onSubmitted)
	// Content is not updated here; mutated directly by keyboard input via SetContent
}

// ========== TextFieldWidget ==========

// TextFieldWidget 是文本输入框，包含背景、边框、内边距和一个可编辑文本。
type TextFieldWidget struct {
	ui.BaseWidget
	editable     EditableTextWidget
	width        float32
	height       float32
	background   *render.Color
	borderColor  *render.Color
	borderWidth  float32
	borderRadius float32
	padding      ui.EdgeInsets
	flexGrow     float32
	flexShrink   float32
}

// TextField 创建文本输入框。
func TextField(initial string) TextFieldWidget {
	return TextFieldWidget{
		editable:     EditableText(initial),
		width:        200,
		height:       40,
		background:   render.NewColor(255, 255, 255, 255),
		borderColor:  render.NewColor(200, 200, 200, 255),
		borderWidth:  1,
		borderRadius: 4,
		padding:      ui.EdgeInsets{Left: 8, Right: 8, Top: 8, Bottom: 8},
	}
}

func (t TextFieldWidget) W(v float32) TextFieldWidget {
	t.width = v
	return t
}

func (t TextFieldWidget) H(v float32) TextFieldWidget {
	t.height = v
	return t
}

func (t TextFieldWidget) Bg(c render.Color) TextFieldWidget {
	t.background = &c
	return t
}

func (t TextFieldWidget) Border(c render.Color, width float32) TextFieldWidget {
	t.borderColor = &c
	t.borderWidth = width
	return t
}

func (t TextFieldWidget) Radius(v float32) TextFieldWidget {
	t.borderRadius = v
	return t
}

func (t TextFieldWidget) Pad(insets ui.EdgeInsets) TextFieldWidget {
	t.padding = insets
	return t
}

func (t TextFieldWidget) Grow(v float32) TextFieldWidget {
	t.flexGrow = v
	return t
}

func (t TextFieldWidget) Shrink(v float32) TextFieldWidget {
	t.flexShrink = v
	return t
}

func (t TextFieldWidget) OnChange(fn func(string)) TextFieldWidget {
	t.editable = t.editable.OnChange(fn)
	return t
}

func (t TextFieldWidget) OnSubmit(fn func(string)) TextFieldWidget {
	t.editable = t.editable.OnSubmit(fn)
	return t
}

func (t TextFieldWidget) CreateElement() ui.Element {
	e := &TextFieldElement{}
	e.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, t)
	return e
}

// TextFieldElement 是 TextFieldWidget 对应的 Element。
type TextFieldElement struct {
	ui.SingleChildRenderObjectElement
}

func (e *TextFieldElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderBox()
	w := e.GetWidget().(TextFieldWidget)
	applyTextFieldProps(r, TextFieldWidget{}, w)
	return r
}

func (e *TextFieldElement) UpdateRenderObject(oldWidget ui.Widget) {
	r := e.GetRenderObject().(*render.RenderBox)
	old := oldWidget.(TextFieldWidget)
	w := e.GetWidget().(TextFieldWidget)
	applyTextFieldProps(r, old, w)
}

func (e *TextFieldElement) UpdateChild(oldWidget ui.Widget) {
	w := e.GetWidget().(TextFieldWidget)
	e.Child = ui.UpdateChild(e, e.Child, w.editable)
}

func (e *TextFieldElement) Mount(parent ui.Element, slot int) {
	e.RenderObject = e.CreateRenderObject()
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(TextFieldWidget)
	e.Child = ui.UpdateChild(e, nil, w.editable)
}

func applyTextFieldProps(r *render.RenderBox, old, w TextFieldWidget) {
	if !render.ColorPtrEquals(old.background, w.background) {
		r.SetBackgroundColor(w.background)
	}
	if !render.ColorPtrEquals(old.borderColor, w.borderColor) {
		r.SetBorderColor(w.borderColor)
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
	if old.flexGrow != w.flexGrow {
		r.StyleSetFlexGrow(w.flexGrow)
	}
	if old.flexShrink != w.flexShrink {
		r.StyleSetFlexShrink(w.flexShrink)
	}
	if old.padding != w.padding {
		r.StyleSetPadding(ui.EdgeTop, w.padding.Top)
		r.StyleSetPadding(ui.EdgeRight, w.padding.Right)
		r.StyleSetPadding(ui.EdgeBottom, w.padding.Bottom)
		r.StyleSetPadding(ui.EdgeLeft, w.padding.Left)
	}
}
