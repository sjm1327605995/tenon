package ui

import (
	"unicode"

	"github.com/rivo/uniseg"
)

// 文本切分（基于 Unicode 分段，rivo/uniseg）：光标/退格按“字素簇”移动，
// Ctrl+方向/双击按“词”移动，正确处理 emoji、组合字、ZWJ 等多码点单元。

// prevGraphemeBoundary 返回 pos 之前一个字素簇的起始字节偏移（退格/左移一个字符的目标）。
func prevGraphemeBoundary(s string, pos int) int {
	if pos <= 0 {
		return 0
	}
	if pos > len(s) {
		pos = len(s)
	}
	off, last, state := 0, 0, -1
	for off < pos {
		var cl string
		cl, _, _, state = uniseg.FirstGraphemeClusterInString(s[off:], state)
		if len(cl) == 0 {
			break
		}
		last = off
		off += len(cl)
	}
	return last
}

// nextGraphemeBoundary 返回 pos 之后一个字素簇的结束字节偏移（删除/右移一个字符的目标）。
func nextGraphemeBoundary(s string, pos int) int {
	if pos < 0 {
		pos = 0
	}
	if pos >= len(s) {
		return len(s)
	}
	cl, _, _, _ := uniseg.FirstGraphemeClusterInString(s[pos:], -1)
	return pos + len(cl)
}

func isSpaceWord(w string) bool {
	if w == "" {
		return false
	}
	for _, r := range w {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// nextWordBoundary 返回 pos 之后下一个词的结束偏移（跳过空白，落在词末）。
func nextWordBoundary(s string, pos int) int {
	off, state := 0, -1
	for off < len(s) {
		var w string
		w, _, state = uniseg.FirstWordInString(s[off:], state)
		if len(w) == 0 {
			break
		}
		end := off + len(w)
		if end > pos && !isSpaceWord(w) {
			return end
		}
		off = end
	}
	return len(s)
}

// prevWordBoundary 返回 pos 之前上一个词的起始偏移（跳过空白，落在词首）。
func prevWordBoundary(s string, pos int) int {
	off, last, state := 0, 0, -1
	for off < len(s) {
		var w string
		w, _, state = uniseg.FirstWordInString(s[off:], state)
		if len(w) == 0 || off >= pos {
			break
		}
		if !isSpaceWord(w) {
			last = off
		}
		off += len(w)
	}
	return last
}

// wordAt 返回包含（或紧邻）pos 的词的字节区间 [start,end)，用于双击选词。
func wordAt(s string, pos int) (int, int) {
	off, state := 0, -1
	for off < len(s) {
		var w string
		w, _, state = uniseg.FirstWordInString(s[off:], state)
		if len(w) == 0 {
			break
		}
		end := off + len(w)
		if pos < end {
			return off, end
		}
		off = end
	}
	return len(s), len(s)
}
