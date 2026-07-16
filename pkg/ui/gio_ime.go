package ui

import (
	"unicode/utf8"

	"gioui.org/f32"
	gioinput "gioui.org/io/input"
	"gioui.org/io/key"
)

// ---- gio 后端：输入法（IME）桥 ----
//
// gio 生态没有可独立复用的 IME 组件：实现埋在 widget.Editor 里，而它同时拥有文本状态、
// 布局、绘制和自己的指针/键盘处理，没法只借用其中的 IME 部分（ebiten 那边有独立的
// exp/textinput.Field，gio 没有对应物）。因此这里按 widget.Editor 的做法讲同一套
// io/key 协议：
//
//	发布：SelectionCmd（光标/选区 + 光标屏幕位置）、SnippetCmd（光标周围的上下文文本）
//	消费：SnippetEvent（输入法索要某段文本）→ 回 SnippetCmd
//	      EditEvent（替换某 rune 区间）→ 经 typedChars 走引擎既有的插入逻辑
//
// 不发布这些，输入法就不知道往哪弹候选窗、也拿不到上下文，于是只能直接走 WM_CHAR ——
// 表现就是「只能输入英文，切不了输入法」。
//
// 注意：引擎内部一律用字节偏移，gio 协议用 rune 偏移，边界处需转换。

// gioIMEState 记录上次发布给输入法的状态，避免每帧重复 Execute。
type gioIMEState struct {
	snippet key.Snippet
	sel     key.Range
	caret   key.Caret
	valid   bool
}

var gioIME gioIMEState

// focusedInput 返回当前聚焦的输入框渲染节点（没有则 nil）。
func focusedInput(g *game) *renderNode {
	f := g.focusedFiber
	if f == nil || f.unmounted || f.rnode == nil || f.rnode.kind != rnInput {
		return nil
	}
	return f.rnode
}

// runeIdx 把字节偏移换成 rune 偏移。
func runeIdx(s string, byteOff int) int {
	if byteOff <= 0 {
		return 0
	}
	if byteOff > len(s) {
		byteOff = len(s)
	}
	return utf8.RuneCountInString(s[:byteOff])
}

// byteIdx 把 rune 偏移换成字节偏移（runeIdx 的逆）。
func byteIdx(s string, runeOff int) int {
	if runeOff <= 0 {
		return 0
	}
	i := 0
	for off := range s {
		if i == runeOff {
			return off
		}
		i++
	}
	return len(s)
}

// applyEdits 把本帧的输入法编辑落到聚焦输入框上。
//
// EditEvent 语义是「把 rng 这段 rune 替换成 text」，不是「在光标处插入」：组字时输入法会
// 用新串反复替换上一次的区间（"s" 之后用 "shi" 替换 [0,1)）。按插入处理会打出 "sshi"。
func (m *gioIMEState) applyEdits(g *game, in *gioInput) {
	rn := focusedInput(g)
	if rn == nil {
		return
	}
	for _, e := range in.edits {
		val := rn.value
		s, en := byteIdx(val, e.rng.Start), byteIdx(val, e.rng.End)
		if s > en {
			s, en = en, s
		}
		nv := val[:s] + e.text + val[en:]
		caret := s + len(e.text)
		rn.value = nv
		rn.caretPos, rn.selAnchor = caret, caret
		if rn.onChange != nil {
			rn.onChange(nv)
		}
	}
	if r := in.selEvt; r != nil {
		val := rn.value
		rn.selAnchor = clampi(byteIdx(val, r.Start), 0, len(val))
		rn.caretPos = clampi(byteIdx(val, r.End), 0, len(val))
	}
	// 组字区间 -> 预编辑下划线。gio 的模型里组字文本已经在 value 内，这里只标范围。
	if r := in.composeEvt; r != nil {
		val := rn.value
		lo := clampi(byteIdx(val, r.Start), 0, len(val))
		hi := clampi(byteIdx(val, r.End), 0, len(val))
		rn.composeLo, rn.composeHi = lo, hi
		g.imeComposing = hi > lo
	}
}

// sync 把聚焦输入框的选区/光标位置与上下文文本发布给输入法。无聚焦输入框时清空状态。
func (m *gioIMEState) sync(src gioinput.Source, g *game) {
	rn := focusedInput(g)
	if rn == nil {
		m.valid = false
		m.snippet = key.Snippet{}
		return
	}
	val := rn.value
	caretB := clampi(rn.caretPos, 0, len(val))
	anchorB := clampi(rn.selAnchor, 0, len(val))
	sel := key.Range{Start: runeIdx(val, anchorB), End: runeIdx(val, caretB)}

	cr := rn.caretRect(caretB) // image.Rectangle，物理像素
	caret := key.Caret{
		Pos:    f32.Pt(float32(cr.Min.X), float32(cr.Max.Y)),
		Ascent: float32(cr.Dy()),
	}
	if !m.valid || sel != m.sel || caret != m.caret {
		m.sel, m.caret, m.valid = sel, caret, true
		src.Execute(key.SelectionCmd{Tag: gioTag, Range: sel, Caret: caret})
	}
	m.sendSnippet(src, val, key.Range{Start: 0, End: utf8.RuneCountInString(val)})
}

// sendSnippet 把 [rng) 范围内的文本发给输入法（仅在与上次不同时发）。
func (m *gioIMEState) sendSnippet(src gioinput.Source, val string, rng key.Range) {
	if rng.Start > rng.End {
		rng.Start, rng.End = rng.End, rng.Start
	}
	n := utf8.RuneCountInString(val)
	if rng.Start > n {
		rng.Start = n
	}
	if rng.End > n {
		rng.End = n
	}
	sn := key.Snippet{Range: rng, Text: substrRunes(val, rng.Start, rng.End)}
	if sn == m.snippet {
		return
	}
	m.snippet = sn
	src.Execute(key.SnippetCmd{Tag: gioTag, Snippet: sn})
}

// substrRunes 按 rune 下标取子串。
func substrRunes(s string, start, end int) string {
	if start >= end {
		return ""
	}
	i, si, ei := 0, -1, len(s)
	for off := range s {
		if i == start {
			si = off
		}
		if i == end {
			ei = off
			break
		}
		i++
	}
	if si < 0 {
		return ""
	}
	return s[si:ei]
}

// handleSnippetReq 响应输入法对某段文本的索要（SnippetEvent）。
func (m *gioIMEState) handleSnippetReq(src gioinput.Source, g *game, rng key.Range) {
	rn := focusedInput(g)
	if rn == nil {
		return
	}
	m.sendSnippet(src, rn.value, rng)
}
