package ui

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rivo/uniseg"
)

// tokenize 按 Unicode 换行算法（UAX#14）把一段文字切成“换行单元”：
// 每个单元结束处是一个允许换行的位置。相比手写规则，它能正确处理
// 拉丁词（含尾随空格）、CJK 逐字、连字符后可断、避免在收尾标点（。」！等）前断行、
// 不间断空格不可断等情形。各单元拼接起来仍等于原串（字节偏移保持不变）。
// 调用方已先按 '\n' 拆段，故此处忽略强制换行标志。
func tokenize(s string) []string {
	var toks []string
	state := -1
	for len(s) > 0 {
		var seg string
		seg, s, _, state = uniseg.FirstLineSegmentInString(s, state)
		toks = append(toks, seg)
	}
	return toks
}

func measureW(s string, face *text.GoTextFace, lineH float64) float32 {
	if s == "" {
		return 0
	}
	w, _ := text.Measure(s, face, lineH)
	return float32(w)
}

// wrapForWidth 按可用宽度对文本贪心折行；未约束或本就放得下时返回单行。
// 返回折行结果与实际最大行宽。
func wrapForWidth(s string, face *text.GoTextFace, lineH float64, width float32) ([]string, float32) {
	if face == nil || s == "" {
		return []string{s}, 0
	}
	nat := measureW(s, face, lineH)
	if (width <= 0 || nat <= width+0.5) && !strings.Contains(s, "\n") {
		return []string{s}, nat
	}

	var lines []string
	maxW := float32(0)
	for _, para := range strings.Split(s, "\n") {
		cur := ""
		for _, tk := range tokenize(para) {
			trial := cur + tk
			if cur != "" && measureW(strings.TrimRight(trial, " "), face, lineH) > width {
				line := strings.TrimRight(cur, " ")
				lines = append(lines, line)
				if lw := measureW(line, face, lineH); lw > maxW {
					maxW = lw
				}
				cur = strings.TrimLeft(tk, " ")
			} else {
				cur = trial
			}
		}
		line := strings.TrimRight(cur, " ")
		lines = append(lines, line)
		if lw := measureW(line, face, lineH); lw > maxW {
			maxW = lw
		}
	}
	return lines, maxW
}

// wrapCache 缓存某文本节点上一次的折行结果，命中条件是文本/字体/行高/宽度均未变。
// 目的：按需重绘时整棵树会重画，但内容不变的文本不必每帧重新 measure/折行。
type wrapCache struct {
	text  string
	face  *text.GoTextFace
	lineH float64
	width float32
	lines []string
	maxW  float32
	valid bool
}

// wrapped 返回按 width 折行的结果，优先命中节点级缓存。
func (rn *renderNode) wrapped(width float32) ([]string, float32) {
	c := &rn.wc
	if c.valid && c.text == rn.text && c.face == rn.face && c.lineH == rn.lineH && c.width == width {
		return c.lines, c.maxW
	}
	lines, mw := wrapForWidth(rn.text, rn.face, rn.lineH, width)
	*c = wrapCache{rn.text, rn.face, rn.lineH, width, lines, mw, true}
	return lines, mw
}

// richCache 缓存富文本节点上一次的排版结果，按解析版本 runsRev + 宽度命中。
type richCache struct {
	rev          int
	width        float32
	lines        []richLine
	maxW, totalH float32
	valid        bool
}

// richLayout 返回富文本按 width 的排版结果，优先命中节点级缓存。
func (rn *renderNode) richLayout(width float32) ([]richLine, float32, float32) {
	c := &rn.rc
	if c.valid && c.rev == rn.runsRev && c.width == width {
		return c.lines, c.maxW, c.totalH
	}
	lines, mw, h := layoutRuns(rn.runs, width)
	*c = richCache{rn.runsRev, width, lines, mw, h, true}
	return lines, mw, h
}

// wrapSpan 是一行折行结果，附带其在原始字符串中的字节范围（end 含被折掉的尾部空格）。
type wrapSpan struct {
	text       string
	start, end int
}

// wrapSpans 与 wrapForWidth 折行规则一致，但额外记录每行的字节偏移，用于多行输入的光标/选区映射。
func wrapSpans(s string, face *text.GoTextFace, lineH float64, width float32) []wrapSpan {
	if face == nil || s == "" {
		return []wrapSpan{{s, 0, len(s)}}
	}
	nat := measureW(s, face, lineH)
	if (width <= 0 || nat <= width+0.5) && !strings.Contains(s, "\n") {
		return []wrapSpan{{s, 0, len(s)}}
	}
	var spans []wrapSpan
	base := 0
	for pi, para := range strings.Split(s, "\n") {
		if pi > 0 {
			base++ // 跳过 '\n'
		}
		spans = append(spans, wrapParaSpans(para, base, face, lineH, width)...)
		base += len(para)
	}
	return spans
}

func wrapParaSpans(para string, base int, face *text.GoTextFace, lineH float64, width float32) []wrapSpan {
	var spans []wrapSpan
	lineStart := base
	cur := ""
	pos := base
	for _, tk := range tokenize(para) {
		trial := cur + tk
		if cur != "" && measureW(strings.TrimRight(trial, " "), face, lineH) > width {
			spans = append(spans, wrapSpan{strings.TrimRight(cur, " "), lineStart, pos})
			cur = tk
			lineStart = pos
		} else {
			cur = trial
		}
		pos += len(tk)
	}
	spans = append(spans, wrapSpan{strings.TrimRight(cur, " "), lineStart, pos})
	return spans
}

// xInSpan 返回原始字节偏移 off 在该行内相对行首 tx 的 x 坐标。
func (sp wrapSpan) xInSpan(off int, face *text.GoTextFace, lineH float64, tx float32) float32 {
	rel := off - sp.start
	if rel < 0 {
		rel = 0
	}
	if rel > len(sp.text) {
		rel = len(sp.text)
	}
	return tx + measureW(sp.text[:rel], face, lineH)
}

// offsetInSpan 返回该行内相对行首 rel-x 处最近的原始字节偏移。
func (sp wrapSpan) offsetInSpan(relX float32, face *text.GoTextFace, lineH float64) int {
	if relX <= 0 || sp.text == "" {
		return sp.start
	}
	prevIdx, prevW := 0, float32(0)
	for i := 1; i <= len(sp.text); i++ {
		if i < len(sp.text) && !isRuneStart(sp.text[i]) {
			continue
		}
		w := measureW(sp.text[:i], face, lineH)
		if w >= relX {
			if w-relX < relX-prevW {
				return sp.start + i
			}
			return sp.start + prevIdx
		}
		prevIdx, prevW = i, w
	}
	return sp.start + len(sp.text)
}

func isRuneStart(b byte) bool { return b&0xC0 != 0x80 }
