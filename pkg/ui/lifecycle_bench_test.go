package ui

import (
	"fmt"
	"runtime"
	"testing"
)

func lifecycleBenchRoot(store *Store[int], text *Signal[string], cleanups *int) *Node {
	return View(func() *Node {
		value := UseStore(store, func(value int) int { return value })
		UseEffect(func() Cleanup {
			return func() { *cleanups++ }
		})
		return Div(Text(fmt.Sprint(value)), Text(text))
	})
}

func TestRepeatedMountCloseReleasesSubscriptions(t *testing.T) {
	store := NewStore(1)
	text := NewSignal("value")
	cleanups := 0
	root := lifecycleBenchRoot(store, text, &cleanups)
	for i := 0; i < 50; i++ {
		h := MountDefault(root)
		if len(store.subs) != 1 || len(text.subs) != 1 {
			t.Fatalf("mount %d subscriptions store=%d text=%d, want 1/1", i, len(store.subs), len(text.subs))
		}
		h.Close()
		if len(store.subs) != 0 || len(text.subs) != 0 {
			t.Fatalf("close %d subscriptions store=%d text=%d, want 0/0", i, len(store.subs), len(text.subs))
		}
		h.Close()
	}
	if cleanups != 50 {
		t.Fatalf("effect cleanups=%d, want 50", cleanups)
	}
}

func TestLifecycleMemoryTrend(t *testing.T) {
	root := layoutBenchTree(200)
	for i := 0; i < 20; i++ {
		h := MountDefault(root)
		h.Close()
	}
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	for i := 0; i < 200; i++ {
		h := MountDefault(root)
		h.Close()
	}
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	growth := int64(after.HeapAlloc) - int64(before.HeapAlloc)
	if growth > 4<<20 {
		t.Fatalf("heap grew by %d bytes after 200 complex mount/close cycles", growth)
	}
	t.Logf("heap before=%d after=%d growth=%d", before.HeapAlloc, after.HeapAlloc, growth)
}

func BenchmarkMountClose(b *testing.B) {
	store := NewStore(1)
	text := NewSignal("value")
	cleanups := 0
	root := lifecycleBenchRoot(store, text, &cleanups)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := MountDefault(root)
		h.Close()
	}
}
