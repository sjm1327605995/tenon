package widgets

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// TextWidget 配置文字显示。
type TextWidget struct {
	engine.BaseWidget
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

func (t TextWidget) CreateElement() engine.Element {
	return engine.NewRenderObjectElement(t)
}

// CreateRenderObject implements RenderObjectFactory.
func (t TextWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	r := render.NewRenderText(t.content)
	r.SetFontSize(t.fontSize)
	r.SetColor(t.textColor)
	r.SetMaxLines(t.maxLines)
	r.SetUnderline(t.underline)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (t TextWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	r := ro.(*render.RenderText)
	old := oldWidget.(TextWidget)
	if old.content != t.content {
		r.SetContent(t.content)
	}
	if old.fontSize != t.fontSize {
		r.SetFontSize(t.fontSize)
	}
	if !old.textColor.Equals(t.textColor) {
		r.SetColor(t.textColor)
	}
	if old.maxLines != t.maxLines {
		r.SetMaxLines(t.maxLines)
	}
	if old.underline != t.underline {
		r.SetUnderline(t.underline)
	}
}
