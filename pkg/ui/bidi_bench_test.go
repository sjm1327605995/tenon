package ui

import (
	"strings"
	"testing"
)

var bidiBenchText = strings.Repeat("English אבג Arabic أبجد and English. ", 32)

func BenchmarkBidiDisplay(b *testing.B) {
	face := bidiTestFace{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bidiDisplay(bidiBenchText, face, 13)
	}
}

func BenchmarkBidiDisplayLTR(b *testing.B) {
	face := bidiTestFace{}
	text := strings.Repeat("Plain English 中文 123. ", 32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bidiDisplay(text, face, 13)
	}
}

func BenchmarkBidiCaretX(b *testing.B) {
	face := bidiTestFace{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bidiCaretX(bidiBenchText, len(bidiBenchText)/2, face, 13, 0)
	}
}

func BenchmarkBidiNodeDisplay(b *testing.B) {
	rn := &renderNode{text: bidiBenchText, face: bidiTestFace{}, lineH: 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rn.displayText(bidiBenchText)
	}
}
