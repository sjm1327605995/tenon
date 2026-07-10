package ui

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func isCJK(r rune) bool {
	return r >= 0x3000 || (r >= 0x2E80 && r < 0x3000)
}

// tokenize 把字符串切成换行单元：拉丁词（连同尾随空格）为一个单元，CJK 字符各自成单元。
func tokenize(s string) []string {
	var toks []string
	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			toks = append(toks, string(cur))
			cur = cur[:0]
		}
	}
	for _, r := range s {
		switch {
		case r == ' ':
			cur = append(cur, r)
			flush()
		case isCJK(r):
			flush()
			toks = append(toks, string(r))
		default:
			cur = append(cur, r)
		}
	}
	flush()
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
	if width <= 0 || nat <= width+0.5 {
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
