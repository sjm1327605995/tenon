package ui

import (
	"testing"

	"gioui.org/op"
)

func BenchmarkGioFrameCacheStatic(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	var ops op.Ops
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ops.Reset()
		h.g.frameCache.Add(h.g, &ops, 1)
	}
}
