package ui

import "testing"

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
