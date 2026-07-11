package ui

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/sjm1327605995/tenon/pkg/font"
)

// textRun 是富文本中一段样式一致的文本。一个 rnText 节点若持有 runs，即按多段混排绘制。
type textRun struct {
	text  string
	style StyleProps

	// resolveRuns 解析后的生效值（在物理像素下取字体，逻辑同 setEffectiveText）
	face       *text.GoTextFace
	color      Color
	lineH      float64
	ascent     float32
	fauxBold   bool
	fauxItalic bool

	// 解析缓存键，避免每帧重复取字体
	rSize   float32
	rWeight int
	rItalic bool
	rScale  float32
	rValid  bool
}

// richSeg 是排版后落在某一行上的一段（同一 run 的连续 token）。
type richSeg struct {
	run  int
	text string
	x    float32 // 相对行首的 x（物理像素）
}

// richLine 是排版后的一行。
type richLine struct {
	segs   []richSeg
	width  float32
	ascent float32 // 行内最大 ascent，用于混排基线对齐
	height float32 // 行高 = 行内最大 lineH
}

// runStyleEqual 只比较文本相关字段（RichText 只从 Text 节点取这些）。
func runStyleEqual(a, b StyleProps) bool {
	return a.hasColor == b.hasColor && a.color == b.color &&
		a.hasFontSize == b.hasFontSize && a.fontSize == b.fontSize &&
		a.hasWeight == b.hasWeight && a.weight == b.weight &&
		a.hasItalic == b.hasItalic && a.italic == b.italic
}

// cloneRuns 复制文字与样式，丢弃解析缓存（renderNode 拥有独立缓存，不污染不可变的 Node）。
func cloneRuns(src []textRun) []textRun {
	if len(src) == 0 {
		return nil
	}
	out := make([]textRun, len(src))
	for i := range src {
		out[i] = textRun{text: src[i].text, style: src[i].style}
	}
	return out
}

func runsEqual(a, b []textRun) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].text != b[i].text || !runStyleEqual(a[i].style, b[i].style) {
			return false
		}
	}
	return true
}

// resolveRuns 解析每段 run 的生效颜色/字号/字重/斜体（未显式设置的字段继承自 ctx）。
// 须在测量（CalculateLayout）之前调用，与 resolveInherited 同一时机。
func (rn *renderNode) resolveRuns(ctx inhText) {
	for i := range rn.runs {
		r := &rn.runs[i]
		st := r.style
		c := Black
		if st.hasColor {
			c = st.color
		} else if ctx.hasColor {
			c = ctx.color
		}
		s := float32(16)
		if st.hasFontSize {
			s = st.fontSize
		} else if ctx.hasSize {
			s = ctx.size
		}
		w := 400
		if st.hasWeight {
			w = st.weight
		} else if ctx.hasWeight {
			w = ctx.weight
		}
		it := false
		if st.hasItalic {
			it = st.italic
		} else if ctx.hasItal {
			it = ctx.italic
		}
		if s <= 0 {
			s = 16
		}
		if w <= 0 {
			w = 400
		}
		r.color = c
		if r.rValid && r.rSize == s && r.rWeight == w && r.rItalic == it && r.rScale == uiScale && r.face != nil {
			continue
		}
		px := s * uiScale
		r.lineH = float64(px) * 1.3
		if ff, err := font.GetDefaultFace(px, font.FontWeight(w), it); err == nil {
			r.face = ff.Face
			r.fauxBold = ff.FauxBold
			r.fauxItalic = ff.FauxItalic
			r.ascent = float32(ff.Face.Metrics().HAscent)
		}
		r.rSize, r.rWeight, r.rItalic, r.rScale, r.rValid = s, w, it, uiScale, true
		rn.runsRev++ // 字体/字号改变 → 令排版缓存失效
	}
}

// layoutRuns 把多段 runs 排成若干行（贪心折行，规则与 wrapForWidth 一致，跨 run 边界允许换行）。
// 返回行、最大行宽、总高度（均为物理像素）。width<=0 表示不约束（不折行，仅按 '\n' 断行）。
func layoutRuns(runs []textRun, width float32) ([]richLine, float32, float32) {
	var lines []richLine
	cur := richLine{}
	curX := float32(0)
	maxW := float32(0)

	grow := func(r *textRun) {
		if r.ascent > cur.ascent {
			cur.ascent = r.ascent
		}
		if float32(r.lineH) > cur.height {
			cur.height = float32(r.lineH)
		}
	}
	pushLine := func() {
		cur.width = curX
		if curX > maxW {
			maxW = curX
		}
		lines = append(lines, cur)
		cur = richLine{}
		curX = 0
	}

	for ri := range runs {
		r := &runs[ri]
		if r.face == nil {
			continue
		}
		for pi, para := range strings.Split(r.text, "\n") {
			if pi > 0 {
				grow(r) // 让（可能为空的）当前行有合理行高
				pushLine()
			}
			for _, tk := range tokenize(para) {
				trimmed := strings.TrimRight(tk, " ")
				if width > 0 && len(cur.segs) > 0 && curX+measureW(trimmed, r.face, r.lineH) > width {
					pushLine()
					tk = strings.TrimLeft(tk, " ")
					if tk == "" {
						continue
					}
				}
				cur.segs = append(cur.segs, richSeg{ri, tk, curX})
				curX += measureW(tk, r.face, r.lineH)
				grow(r)
			}
		}
	}
	pushLine()

	var totalH float32
	for i := range lines {
		if lines[i].height == 0 && len(runs) > 0 {
			lines[i].height = float32(runs[0].lineH)
			lines[i].ascent = runs[0].ascent
		}
		totalH += lines[i].height
	}
	return lines, maxW, totalH
}

// paintRichText 逐行、逐段绘制富文本，混排时按基线对齐。
func paintRichText(p painter, rn *renderNode, o float32) {
	b := rn.bounds
	lines, _, _ := rn.richLayout(b.W)
	y := b.Y
	for _, ln := range lines {
		baseline := y + ln.ascent
		for _, sg := range ln.segs {
			r := &rn.runs[sg.run]
			ty := baseline - r.ascent
			p.DrawText(sg.text, r.face, r.color.Alpha(o), b.X+sg.x, ty, r.fauxBold, r.fauxItalic)
		}
		y += ln.height
	}
}
