package ui

import (
	"sync"

	"gioui.org/f32"
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
	gioShaperOnce.Do(func() { gioShaperInst = newGioShaper() })
	return gioShaperInst
}

// newGioShaper 用内置字体建 shaper。
//
// 解析失败必须炸：cjkFont 是 //go:embed 进来的资源，解析不了只可能是资源本身坏了或
// 换错了文件 —— 属于构建期问题，不是运行期可恢复的状况。而如果在这里静默吞掉错误，
// text.NewShaper 会回落到系统字体：开发机上通常装着中文字体、看起来一切正常，
// 换到没有 CJK 字体的机器上就是满屏豆腐，且没有任何线索指向真正的原因。
func newGioShaper() *text.Shaper {
	face, err := opentype.Parse(cjkFont)
	if err != nil {
		panic("ui: 内置字体解析失败（assets/OPPOSans-Medium.ttf 可能损坏或被替换）: " + err.Error())
	}
	return text.NewShaper(text.WithCollection(
		[]giofont.FontFace{{Font: giofont.Font{Typeface: "OPPOSans"}, Face: face}},
	))
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

// Metrics 报告基线位置，以及是否需要合成粗体/斜体。
//
// 内置字体集只有一张 OPPOSans Medium（Regular 字形、常规字重）。gio 不会自己合成 ——
// 它只会在字体集里挑最接近的那张，于是 Bold 和 Regular 会渲染得一模一样（实测字形 ID
// 与宽度完全相同）。所以字重/斜体只能由我们合成，painter 接口早就为此留了这两个参数。
func (f *gioFont) Metrics() (float32, bool, bool) {
	return f.ascent, f.weight >= giofont.Weight(fauxBoldMinWeight-400), f.style == giofont.Italic
}

// fauxBoldMinWeight 是启用合成粗体的 CSS 字重阈值（Semibold 起）。
const fauxBoldMinWeight = 600

// fauxBoldWidth 是合成粗体的描边宽度：在填充之外再描一圈，让字形整体变粗。
// 取字号的比例，这样各字号下的观感一致。
func fauxBoldWidth(px float32) float32 { return px / 22 }

// fauxItalicShear 是合成斜体的倾角（弧度，约 12 度），与常见字体的斜角接近。
// 注意 gio 的 Shear 收的是弧度（内部取 tan），不是斜率。
const fauxItalicShear = 0.21

// drawGioText 把一行文本绘制到 ops：(x,y) 为该行左上角，基线落在 y+ascent。
// fauxBold/fauxItalic 由调用方从 Metrics 取得：内置字体只有一张常规face，粗体与斜体
// 必须在这里合成，否则 ui.Bold 会完全没有效果。
func drawGioText(ops *op.Ops, f *gioFont, s string, c Color, x, y float32, fauxBold, fauxItalic bool) {
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
	aff := gioOffset(x, y+f.ascent)
	if fauxItalic {
		// 必须「先在字形坐标里剪切、再平移」：Mul 是矩阵乘（A.Mul(B) = 先 B 后 A），
		// 写反了剪切就作用在屏幕坐标上，整行字会被 y 推着横向平移而不是倾斜。
		// 字形在基线以上是负 y，故用负角度让上半部右倾。
		aff = aff.Mul(f32.Affine2D{}.Shear(f32.Pt(0, 0), -fauxItalicShear, 0))
	}
	tr := op.Affine(aff).Push(ops)
	gpaint.ColorOp{Color: nrgba(c)}.Add(ops)
	path := sh.Shape(glyphs)
	cl := clip.Outline{Path: path}.Op().Push(ops)
	gpaint.PaintOp{}.Add(ops)
	cl.Pop()
	if fauxBold {
		// 在填充之外再描一圈轮廓：字形整体向外扩，即合成粗体
		st := clip.Stroke{Path: path, Width: fauxBoldWidth(f.px)}.Op().Push(ops)
		gpaint.PaintOp{}.Add(ops)
		st.Pop()
	}
	if call := sh.Bitmaps(glyphs); call != (op.CallOp{}) {
		call.Add(ops)
	}
	tr.Pop()
}
