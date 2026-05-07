package widgets

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// EditableTextWidget 是可编辑文本 Widget。
type EditableTextWidget struct {
	ui.BaseWidget
	content          string
	placeholder      string
	fontSize         float32
	textColor        color.Color
	placeholderColor color.Color
	multiline        bool
	onChanged        func(string)
	onSubmitted      func(string)
}

// EditableText 创建可编辑文本 Widget。
func EditableText(content string) EditableTextWidget {
	return EditableTextWidget{
		content:          content,
		fontSize:         14,
		textColor:        color.Black,
		placeholderColor: ui.GetTheme().InputPlaceholderColor,
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

func (e EditableTextWidget) Placeholder(text string) EditableTextWidget {
	e.placeholder = text
	return e
}

func (e EditableTextWidget) PlaceholderColor(c color.Color) EditableTextWidget {
	e.placeholderColor = c
	return e
}

func (e EditableTextWidget) Multiline(v bool) EditableTextWidget {
	e.multiline = v
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
	return ui.NewRenderObjectElement(e)
}

// CreateRenderObject implements RenderObjectFactory.
func (e EditableTextWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	r := render.NewRenderEditableText()
	r.SetContent(e.content)
	r.SetPlaceholder(e.placeholder)
	r.SetPlaceholderColor(e.placeholderColor)
	r.SetMultiline(e.multiline)
	r.SetFontSize(e.fontSize)
	r.SetColor(e.textColor)
	r.SetOnChanged(e.onChanged)
	r.SetOnSubmitted(e.onSubmitted)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (e EditableTextWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	r := ro.(*render.RenderEditableText)
	old := oldWidget.(EditableTextWidget)

	if old.content != e.content {
		r.SetContent(e.content)
	}
	if old.fontSize != e.fontSize {
		r.SetFontSize(e.fontSize)
	}
	if !render.ColorEquals(old.textColor, e.textColor) {
		r.SetColor(e.textColor)
	}
	if old.placeholder != e.placeholder {
		r.SetPlaceholder(e.placeholder)
	}
	if !render.ColorEquals(old.placeholderColor, e.placeholderColor) {
		r.SetPlaceholderColor(e.placeholderColor)
	}
	if old.multiline != e.multiline {
		r.SetMultiline(e.multiline)
	}
	r.SetOnChanged(e.onChanged)
	r.SetOnSubmitted(e.onSubmitted)
}

// ========== TextFieldWidget ==========

// TextFieldWidget 是文本输入框，包含背景、边框、内边距和一个可编辑文本。
type TextFieldWidget struct {
	ui.BaseWidget
	editable         EditableTextWidget
	placeholder      string
	width            float32
	height           float32
	background       *render.Color
	borderColor      *render.Color
	focusBorderColor *render.Color
	borderWidth      float32
	borderRadius     float32
	padding          ui.EdgeInsets
	flexGrow         float32
	flexShrink       float32
	multiline        bool
}

// TextField 创建文本输入框。
func TextField(initial string) TextFieldWidget {
	return TextFieldWidget{
		editable:         EditableText(initial),
		width:            200,
		height:           40,
		background:       render.NewColor(255, 255, 255, 255),
		borderColor:      render.NewColor(200, 200, 200, 255),
		focusBorderColor: render.NewColor(23, 23, 23, 255),
		borderWidth:      1,
		borderRadius:     4,
		padding:          ui.EdgeInsets{Left: 8, Right: 8, Top: 8, Bottom: 8},
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

func (t TextFieldWidget) Placeholder(text string) TextFieldWidget {
	t.placeholder = text
	return t
}

func (t TextFieldWidget) Multiline(v bool) TextFieldWidget {
	t.multiline = v
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

func (t TextFieldWidget) FocusBorder(c render.Color) TextFieldWidget {
	t.focusBorderColor = &c
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

func (t TextFieldWidget) buildEditable() EditableTextWidget {
	e := t.editable
	if t.placeholder != "" && e.placeholder == "" {
		e = e.Placeholder(t.placeholder)
	}
	if t.multiline {
		e = e.Multiline(true)
	}
	return e
}

func (t TextFieldWidget) CreateElement() ui.Element {
	return ui.NewSingleChildRenderObjectElement(t)
}

// CreateRenderObject implements RenderObjectFactory.
func (t TextFieldWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	r := render.NewRenderBox()
	applyTextFieldProps(r, TextFieldWidget{}, t)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (t TextFieldWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	old := oldWidget.(TextFieldWidget)
	applyTextFieldProps(ro.(*render.RenderBox), old, t)
}

// GetChildWidget implements SingleChildProvider.
func (t TextFieldWidget) GetChildWidget() ui.Widget {
	return t.buildEditable()
}

func applyTextFieldProps(r *render.RenderBox, old, w TextFieldWidget) {
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
