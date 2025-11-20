package component

import (
	colEmoji "eliasnaur.com/font/noto/emoji/color"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	giotext "gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"golang.org/x/image/math/fixed"
	"image"
)

type Text struct {
	Base[Text]
	fontSize unit.Sp
	content  string
}

func NewText(content string) *Text {
	text := &Text{
		content:  content,
		fontSize: 16,
	}
	text.Base = NewBase(text)
	return text
}
func (t *Text) Update(ctx layout.Context) {
	exceptX, exceptY := t.Node.StyleGetWidth(), t.Node.StyleGetHeight()
	faces, err := opentype.ParseCollection(colEmoji.TTF)
	if err != nil {
		panic(err)
	}
	if exceptX > 0 {
		ctx.Constraints.Max.X = int(exceptX)
	}
	if exceptY > 0 {
		ctx.Constraints.Max.Y = int(exceptY)
	}
	th := material.NewTheme()
	collection := gofont.Collection()

	th.Shaper = giotext.NewShaper(giotext.WithCollection(append(collection, faces...)))

	labelStyle := material.Label(th, t.fontSize, t.content)
	textGio := &TextGio{labelStyle: &labelStyle}
	t.gio = textGio
	th.Shaper.LayoutString(giotext.Parameters{
		Font:            textGio.labelStyle.Font,
		Alignment:       labelStyle.Alignment,
		PxPerEm:         fixed.I(ctx.Sp(t.fontSize)),
		MaxLines:        labelStyle.MaxLines,
		WrapPolicy:      labelStyle.WrapPolicy,
		MinWidth:        0,
		MaxWidth:        ctx.Constraints.Max.X,
		Locale:          ctx.Locale,
		LineHeightScale: labelStyle.LineHeightScale,
		LineHeight:      fixed.I(ctx.Sp(labelStyle.LineHeight)),
	}, t.content)
	var total fixed.Point26_6 // 累计所有字符的宽高（26.6格式）

	for {
		g, ok := th.Shaper.NextGlyph()
		if !ok {
			break
		}

		// 计算单个字符的宽和高（26.6格式，Max是排他的，直接相减）
		glyphSize := g.Bounds.Max.Sub(g.Bounds.Min)
		total = total.Add(glyphSize)
	}
	x, y := total.X.Round(), total.Y.Round()
	textGio.Size = image.Point{X: x, Y: y}
	t.Node.StyleSetWidth(float32(x))
	t.Node.StyleSetHeight(float32(y))

}
func (t *Text) FontSize(size unit.Sp) *Text {
	t.fontSize = size
	return t
}

func (t *Text) Content(content string) *Text {
	t.content = content
	return t
}

type TextGio struct {
	labelStyle *material.LabelStyle
	Size       image.Point
}

func (t *TextGio) Layout(gtx layout.Context) layout.Dimensions {

	gtx.Constraints.Min = t.Size
	gtx.Constraints.Max = t.Size
	return t.labelStyle.Layout(gtx)
}
