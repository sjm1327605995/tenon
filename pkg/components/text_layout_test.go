package components

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sjm1327605995/tenon/pkg/fonts"
)

func init() {
	_ = fonts.InitDefaultFont()
}

func getTestFace(t *testing.T) *text.GoTextFace {
	fm := fonts.GetFontManager()
	face, err := fm.GetDefaultFontFace(16)
	if err != nil || face == nil || face.Face == nil {
		t.Skip("default font not available")
	}
	return face.Face
}

func TestWhiteSpaceNoWrap(t *testing.T) {
	face := getTestFace(t)
	result := computeTextLayout("Hello   World", face, 100, WhiteSpaceNoWrap, WordBreakNormal, 0)
	if result.lineCount != 1 {
		t.Fatalf("expected 1 line for nowrap, got %d", result.lineCount)
	}
	if result.content != "Hello World" {
		t.Fatalf("expected collapsed text, got %q", result.content)
	}
}

func TestWhiteSpacePre(t *testing.T) {
	face := getTestFace(t)
	result := computeTextLayout("Hello   World", face, 100, WhiteSpacePre, WordBreakNormal, 0)
	if result.lineCount != 1 {
		t.Fatalf("expected 1 line for pre, got %d", result.lineCount)
	}
	if result.content != "Hello   World" {
		t.Fatalf("expected preserved spaces, got %q", result.content)
	}
}

func TestWhiteSpacePreLine(t *testing.T) {
	face := getTestFace(t)
	result := computeTextLayout("Hello   World\n  Foo  Bar", face, 1000, WhiteSpacePreLine, WordBreakNormal, 0)
	if result.lineCount != 2 {
		t.Fatalf("expected 2 lines for pre-line, got %d", result.lineCount)
	}
	lines := splitLines(result.content)
	if lines[0] != "Hello World" || lines[1] != "Foo Bar" {
		t.Fatalf("expected collapsed spaces with preserved lines, got %q", result.content)
	}
}

func TestAutoWrap(t *testing.T) {
	face := getTestFace(t)
	// 一个很长的英文句子，maxWidth 很小，应该被换行
	result := computeTextLayout("The quick brown fox jumps over the lazy dog", face, 80, WhiteSpaceNormal, WordBreakNormal, 0)
	if result.lineCount <= 1 {
		t.Fatalf("expected multiple lines for auto wrap, got %d", result.lineCount)
	}
}

func TestWordBreakBreakAll(t *testing.T) {
	face := getTestFace(t)
	result := computeTextLayout("abcdefghij", face, 50, WhiteSpaceNormal, WordBreakBreakAll, 0)
	// break-all 下每个字符都可断，应该产生多行
	if result.lineCount <= 1 {
		t.Fatalf("expected multiple lines for break-all, got %d", result.lineCount)
	}
}

func TestCJKWrap(t *testing.T) {
	face := getTestFace(t)
	result := computeTextLayout("这是一个很长的中文文本用于测试换行功能", face, 100, WhiteSpaceNormal, WordBreakNormal, 0)
	if result.lineCount <= 1 {
		t.Fatalf("expected multiple lines for CJK text, got %d", result.lineCount)
	}
}

func TestPreWrap(t *testing.T) {
	face := getTestFace(t)
	result := computeTextLayout("Line1\nLine2\nLine3", face, 1000, WhiteSpacePreWrap, WordBreakNormal, 0)
	if result.lineCount != 3 {
		t.Fatalf("expected 3 lines for pre-wrap, got %d", result.lineCount)
	}
}

func TestLineHeight(t *testing.T) {
	face := getTestFace(t)
	result1 := computeTextLayout("A\nB", face, 1000, WhiteSpacePreWrap, WordBreakNormal, 0)
	result2 := computeTextLayout("A\nB", face, 1000, WhiteSpacePreWrap, WordBreakNormal, 40)
	if result2.height <= result1.height {
		t.Fatalf("explicit line height should increase total height")
	}
}


