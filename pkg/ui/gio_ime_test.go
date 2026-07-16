package ui

import (
	"testing"

	"gioui.org/io/key"
)

// 回归：输入法组字发的是「替换某 rune 区间」，不是「在光标处追加」。
// 打 "shi" 时输入法先给 "s"，再用 "shi" 替换 [0,1)。当成追加会写出 "sshi"。
func TestIMEEditEventReplacesRange(t *testing.T) {
	var got string
	rn := &renderNode{kind: rnInput}
	rn.onChange = func(s string) { got = s }
	f := &Fiber{rnode: rn}
	rn.owner = f
	g := &game{focusedFiber: f}

	in := &gioInput{}
	in.edits = []gioEdit{
		{rng: key.Range{Start: 0, End: 0}, text: "s"},   // 组字开始
		{rng: key.Range{Start: 0, End: 1}, text: "shi"}, // 用新串替换上一次的区间
	}
	gioIME.applyEdits(g, in)

	if rn.value != "shi" {
		t.Fatalf("value=%q want %q（追加而非替换会得到 \"sshi\"）", rn.value, "shi")
	}
	if got != "shi" {
		t.Fatalf("onChange 收到 %q want %q", got, "shi")
	}
	if rn.caretPos != len("shi") {
		t.Fatalf("caretPos=%d want %d", rn.caretPos, len("shi"))
	}
}

// 组字提交中文：区间替换要按 rune 边界算，不能按字节。
func TestIMECommitCJKUsesRuneRange(t *testing.T) {
	rn := &renderNode{kind: rnInput, value: "shi"}
	f := &Fiber{rnode: rn}
	rn.owner = f
	g := &game{focusedFiber: f}

	in := &gioInput{}
	// 把 3 个 rune 的拼音替换成 1 个汉字
	in.edits = []gioEdit{{rng: key.Range{Start: 0, End: 3}, text: "是"}}
	gioIME.applyEdits(g, in)

	if rn.value != "是" {
		t.Fatalf("value=%q want %q", rn.value, "是")
	}
	if rn.caretPos != len("是") {
		t.Fatalf("caretPos=%d want %d（应为字节偏移）", rn.caretPos, len("是"))
	}
}

// 引擎用字节偏移、gio 的 IME 协议用 rune 偏移，转换错了会让输入法把候选插到错位置。
func TestRuneIdxFromByteOffset(t *testing.T) {
	const s = "a你好b" // 字节: a=1, 你=3, 好=3, b=1
	cases := []struct{ byteOff, want int }{
		{0, 0},
		{1, 1},  // "a"
		{4, 2},  // "a你"
		{7, 3},  // "a你好"
		{8, 4},  // "a你好b"
		{99, 4}, // 越界钳制
		{-1, 0},
	}
	for _, c := range cases {
		if got := runeIdx(s, c.byteOff); got != c.want {
			t.Errorf("runeIdx(%q, %d)=%d want %d", s, c.byteOff, got, c.want)
		}
	}
}

func TestSubstrRunes(t *testing.T) {
	const s = "a你好b"
	cases := []struct {
		start, end int
		want       string
	}{
		{0, 4, "a你好b"},
		{1, 3, "你好"},
		{0, 1, "a"},
		{3, 4, "b"},
		{2, 2, ""},
		{3, 1, ""}, // start>=end
	}
	for _, c := range cases {
		if got := substrRunes(s, c.start, c.end); got != c.want {
			t.Errorf("substrRunes(%q,%d,%d)=%q want %q", s, c.start, c.end, got, c.want)
		}
	}
}
