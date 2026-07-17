package ui

import (
	"fmt"
	"testing"
)

// flex-wrap：yoga 一直支持，只是没透出来。放不下的子元素必须折到下一行。
func TestWrapPutsOverflowOnNextRow(t *testing.T) {
	// 容器 300 宽，子元素 80 宽 + gap 10 -> 一行放得下 3 个（80*3+10*2=260，第 4 个要 350）
	build := func(wrap bool) *Harness {
		st := []StyleOpt{Width(300), Row, Gap(10), ItemsStart}
		if wrap {
			st = append(st, Wrap)
		}
		kids := []*Node{Style(st...)}
		for i := 0; i < 7; i++ {
			kids = append(kids, Div(Style(Width(80), Height(40)),
				Text(fmt.Sprintf("k%d", i))))
		}
		return Mount(Use(func(struct{}) *Node {
			return Div(Style(Width(600), Height(400)), Div(kids...))
		}, struct{}{}), 600, 400)
	}

	nw, w := build(false), build(true)
	rowsOf := func(h *Harness) map[float32][]string {
		rows := map[float32][]string{}
		for i := 0; i < 7; i++ {
			n := fmt.Sprintf("k%d", i)
			b := h.Root().ByText(n).Bounds()
			rows[b.Y] = append(rows[b.Y], n)
		}
		return rows
	}
	nwRows, wRows := rowsOf(nw), rowsOf(w)
	t.Logf("NoWrap（默认）: %d 行 %v", len(nwRows), nwRows)
	t.Logf("Wrap:          %d 行 %v", len(wRows), wRows)

	// A：不换行时全挤在一行（被压缩）—— 前提
	if len(nwRows) != 1 {
		t.Errorf("不换行时排成了 %d 行，want 1 —— 前提不成立", len(nwRows))
	}
	// B：换行时必须折行
	if len(wRows) < 3 {
		t.Errorf("Wrap 后只有 %d 行，want>=3（7 个 x 80px 塞不进 300px 一行）—— 换行没生效", len(wRows))
	}
	// 每行都不能超出容器宽度
	for y, names := range wRows {
		for _, n := range names {
			b := w.Root().ByText(n).Bounds()
			if b.X+b.W > 300.5 {
				t.Errorf("第 %v 行的 %s 右边缘 %.1f 超出容器 300 —— 没真的折行", y, n, b.X+b.W)
			}
		}
	}
}

// 不换行时子元素保持原宽、直接溢出容器（注意：tenon 的 Shrink 默认为 0，
// 与 CSS 的 flex-shrink:1 不同，所以是溢出而不是被压缩）。换行后不得溢出。
func TestWrapStopsOverflow(t *testing.T) {
	build := func(opts ...StyleOpt) (maxRight float32, w0 float32) {
		st := append([]StyleOpt{Width(300), Row, Gap(10)}, opts...)
		kids := []*Node{Style(st...)}
		for i := 0; i < 7; i++ {
			kids = append(kids, Div(Style(Width(80), Height(40)), Text(fmt.Sprintf("k%d", i))))
		}
		h := Mount(Use(func(struct{}) *Node {
			return Div(Style(Width(600), Height(400)), Div(kids...))
		}, struct{}{}), 600, 400)
		for i := 0; i < 7; i++ {
			b := h.Root().ByText(fmt.Sprintf("k%d", i)).Bounds()
			if r := b.X + b.W; r > maxRight {
				maxRight = r
			}
		}
		return maxRight, h.Root().ByText("k0").Bounds().W
	}
	nwRight, nwW := build()
	wRight, wW := build(Wrap)
	t.Logf("不换行：最右边缘=%.0f（容器 300）k0 宽=%.0f", nwRight, nwW)
	t.Logf("换行：  最右边缘=%.0f（容器 300）k0 宽=%.0f", wRight, wW)

	if nwRight <= 300 {
		t.Errorf("不换行时最右边缘 %.0f 没超出容器 300 —— 前提不成立，本测试无法归因", nwRight)
	}
	if wRight > 300.5 {
		t.Errorf("换行后最右边缘 %.0f 仍超出容器 300 —— 换行没挡住溢出", wRight)
	}
	if wW != 80 {
		t.Errorf("换行后 k0 宽 %.0f，want 80 —— 不该被压缩", wW)
	}
}
