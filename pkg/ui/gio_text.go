package ui

import (
	"sync"

	giofont "gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/op"
	"gioui.org/op/clip"
	gpaint "gioui.org/op/paint"
	"gioui.org/text"
	"golang.org/x/image/math/fixed"
)

// ---- gio 后端：字体与文本 ----
//
// 复用引擎内置的 OPPOSans（cjkFont，含 CJK 字形），经 gio 的 opentype 解析后交给
// text.Shaper 整形。gioFont 实现 fontFace：measure 用整形推进量求和，绘制走 Shape→Outline。

var (
	gioShaperOnce sync.Once
	gioShaperInst *text.Shaper
)

func gioShaper() *text.Shaper {
	gioShaperOnce.Do(func() {
		var faces []giofont.FontFace
		if f, err := opentype.Parse(cjkFont); err == nil {
			faces = []giofont.FontFace{{Font: giofont.Font{Typeface: "OPPOSans"}, Face: f}}
		}
		gioShaperInst = text.NewShaper(text.WithCollection(faces))
	})
	return gioShaperInst
}

func f2i(v float32) fixed.Int26_6 { return fixed.Int26_6(v*64 + 0.5) }

// gioFont 是某一字号/字重/斜体下的字体句柄。gio 内部自动合成缺失的粗/斜，故无需伪粗伪斜。
type gioFont struct {
	px     float32
	weight giofont.Weight
	style  giofont.Style
	ascent float32
}

func gioNewFont(px float32, weight int, italic bool) fontFace {
	st := giofont.Regular
	if italic {
		st = giofont.Italic
	}
	// CSS 字重(100..900, 400=常规) -> gio 字重(以 Normal=0 为基准的偏移)。
	f := &gioFont{px: px, weight: giofont.Weight(weight - 400), style: st}
	f.ascent = f.measureAscent()
	return f
}

func (f *gioFont) params() text.Parameters {
	return text.Parameters{
		Font:    giofont.Font{Weight: f.weight, Style: f.style},
		PxPerEm: f2i(f.px),
		// 折行由引擎自己完成，这里给足够大的宽度避免 shaper 把每个字形单独换行（MaxWidth=0
		// 会被当作“可用宽度为 0”，导致文本被逐字竖排、看起来散开不成块）。
		MaxWidth: 1 << 30,
	}
}

func (f *gioFont) Measure(s string, lineH float64) float32 {
	if s == "" {
		return 0
	}
	sh := gioShaper()
	sh.LayoutString(f.params(), s)
	var w fixed.Int26_6
	for {
		g, ok := sh.NextGlyph()
		if !ok {
			break
		}
		w += g.Advance
	}
	return float32(w) / 64
}

func (f *gioFont) measureAscent() float32 {
	sh := gioShaper()
	sh.LayoutString(f.params(), "Ag")
	var asc fixed.Int26_6
	for {
		g, ok := sh.NextGlyph()
		if !ok {
			break
		}
		if g.Ascent > asc {
			asc = g.Ascent
		}
	}
	return float32(asc) / 64
}

func (f *gioFont) Metrics() (float32, bool, bool) { return f.ascent, false, false }

// drawGioText 把一行文本绘制到 ops：(x,y) 为该行左上角，基线落在 y+ascent。
func drawGioText(ops *op.Ops, f *gioFont, s string, c Color, x, y float32) {
	if s == "" {
		return
	}
	sh := gioShaper()
	sh.LayoutString(f.params(), s)
	var glyphs []text.Glyph
	for {
		g, ok := sh.NextGlyph()
		if !ok {
			break
		}
		glyphs = append(glyphs, g)
	}
	if len(glyphs) == 0 {
		return
	}
	tr := op.Affine(gioOffset(x, y+f.ascent)).Push(ops)
	gpaint.ColorOp{Color: nrgba(c)}.Add(ops)
	path := sh.Shape(glyphs)
	cl := clip.Outline{Path: path}.Op().Push(ops)
	gpaint.PaintOp{}.Add(ops)
	cl.Pop()
	if call := sh.Bitmaps(glyphs); call != (op.CallOp{}) {
		call.Add(ops)
	}
	tr.Pop()
}
