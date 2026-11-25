package ui

import (
	"strconv"
	"strings"
	"testing"
)

func virtualListBenchmarkRoot() *Node {
	return VirtualList(VirtualListProps{
		Count: 100000, ItemHeight: 20, Height: 200,
		Render: func(i int) *Node { return Text("row-" + strconv.Itoa(i)) },
	})
}

func checkVirtualListBenchmarkWindow(b *testing.B, h *Harness) {
	b.Helper()
	rows := h.Root().FindAll(func(q *Query) bool {
		return q.rn.kind == rnText && strings.HasPrefix(q.Text(), "row-")
	})
	if len(rows) != 17 {
		b.Fatalf("mounted rows = %d, want 17 visible/overscan rows out of 100000", len(rows))
	}
	if h.Root().ByText("row-99999").Exists() {
		b.Fatal("far off-screen row was mounted")
	}
}

func BenchmarkVirtualList100kStartup(b *testing.B) {
	root := virtualListBenchmarkRoot()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := MountDefault(root)
		h.PaintTo(discardPainter{})
		h.Close()
	}
}

func BenchmarkVirtualList100kScroll(b *testing.B) {
	h := MountDefault(virtualListBenchmarkRoot())
	defer h.Close()
	checkVirtualListBenchmarkWindow(b, h)
	scroll := h.Root().ByKind("scroll")
	direction := float32(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if scroll.ScrollOffset() >= 2000 {
			direction = -20
		} else if scroll.ScrollOffset() <= 0 {
			direction = 20
		}
		if !scroll.ScrollBy(direction) {
			b.Fatal("virtual list scroll container is missing")
		}
		h.PaintTo(discardPainter{})
	}
}
