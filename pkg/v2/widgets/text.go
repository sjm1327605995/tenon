package widgets

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// TextWidget 配置文字显示。
type TextWidget struct {
	ui.BaseWidget
	content   string
	fontSize  float32
	textColor *render.Color
	maxLines  int
	underline bool
}

// Text 创建文字 Widget。
func Text(content string) TextWidget {
	return TextWidget{
		content:   content,
		fontSize:  14,
		textColor: render.NewColor(0, 0, 0, 255),
		maxLines:  0,
	}
}

// FontSize 设置字体大小，返回新实例。
func (t TextWidget) FontSize(size float32) TextWidget {
	t.fontSize = size
	return t
}

// Color 设置文字颜色，返回新实例。
func (t TextWidget) Color(c color.Color) TextWidget {
	t.textColor = render.NewColorFrom(c)
	return t
}

// MaxLines 设置最大行数，返回新实例。
func (t TextWidget) MaxLines(n int) TextWidget {
	t.maxLines = n
	return t
}

// Underline 设置是否显示下划线，返回新实例。
func (t TextWidget) Underline() TextWidget {
	t.underline = true
	return t
}

func (t TextWidget) CreateElement() ui.Element {
	e := &TextElement{}
	e.RenderObjectElement.BaseElement.Init(e, t)
	return e
}

// TextElement 是 TextWidget 对应的 Element。
type TextElement struct {
	ui.RenderObjectElement
	ro *render.RenderText
}

func (e *TextElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderText)
	e.RenderObject = e.ro
	e.RenderObjectElement.Mount(parent, slot)
}

func (e *TextElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(TextWidget)
	r := render.NewRenderText(w.content)
	r.SetFontSize(w.fontSize)
	r.SetColor(w.textColor)
	r.SetMaxLines(w.maxLines)
	r.SetUnderline(w.underline)
	return r
}

func (e *TextElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(TextWidget)
	if old, ok := oldWidget.(TextWidget); !ok || old.content != w.content {
		e.ro.SetContent(w.content)
	}
	if old, ok := oldWidget.(TextWidget); !ok || old.fontSize != w.fontSize {
		e.ro.SetFontSize(w.fontSize)
	}
	if old, ok := oldWidget.(TextWidget); !ok || !old.textColor.Equals(w.textColor) {
		e.ro.SetColor(w.textColor)
	}
	if old, ok := oldWidget.(TextWidget); !ok || old.maxLines != w.maxLines {
		e.ro.SetMaxLines(w.maxLines)
	}
	if old, ok := oldWidget.(TextWidget); !ok || old.underline != w.underline {
		e.ro.SetUnderline(w.underline)
	}
}
