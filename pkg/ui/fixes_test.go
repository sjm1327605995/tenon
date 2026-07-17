package ui

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// findProvider 在 Fiber 树中找到指定 ctxID 的 provider fiber。
func findProvider(f *Fiber, id int) *Fiber {
	if f.typ == typeProvider && f.ctxID == id {
		return f
	}
	for _, c := range f.children {
		if r := findProvider(c, id); r != nil {
			return r
		}
	}
	return nil
}

// ---- Fix: Context 订阅去重 ----

var (
	dedupCtx    = CreateContext("v")
	dedupSetter func(int)
)

func dedupConsumer(_ struct{}) *Node {
	n, set := UseState(0)
	dedupSetter = set
	_ = UseContext(dedupCtx)
	return Text(itoa(n))
}

func dedupApp(_ struct{}) *Node {
	return dedupCtx.Provider("v", Use(dedupConsumer, struct{}{}))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

// 消费者独立重渲染多次（provider 值不变）不应让 subscribers 无限累积。
func TestContextSubscriberDedup(t *testing.T) {
	g := newGame()
	g.mountRoot(Use(dedupApp, struct{}{}))
	p := findProvider(g.rootFiber, dedupCtx.id)
	if p == nil {
		t.Fatal("provider fiber not found")
	}
	if len(p.subscribers) != 1 {
		t.Fatalf("after mount subscribers=%d, want 1", len(p.subscribers))
	}
	for i := 1; i <= 5; i++ {
		dedupSetter(i)
		g.drain()
	}
	if len(p.subscribers) != 1 {
		t.Fatalf("after 5 self re-renders subscribers=%d, want 1 (dedup)", len(p.subscribers))
	}
}

// ---- Fix: 卸载消费者时退订 ----

var (
	unsubCtx     = CreateContext("v")
	unsubSetShow func(bool)
)

func unsubConsumer(_ struct{}) *Node {
	_ = UseContext(unsubCtx)
	return Text("c")
}

func unsubApp(_ struct{}) *Node {
	show, set := UseState(true)
	unsubSetShow = set
	if show {
		return unsubCtx.Provider("v", Use(unsubConsumer, struct{}{}))
	}
	return unsubCtx.Provider("v")
}

func TestContextUnsubscribeOnUnmount(t *testing.T) {
	g := newGame()
	g.mountRoot(Use(unsubApp, struct{}{}))
	p := findProvider(g.rootFiber, unsubCtx.id)
	if p == nil || len(p.subscribers) != 1 {
		t.Fatalf("after mount provider=%v subscribers=%d, want 1", p != nil, len(p.subscribers))
	}
	unsubSetShow(false)
	g.drain()
	if len(p.subscribers) != 0 {
		t.Fatalf("after consumer unmount subscribers=%d, want 0 (must unsubscribe)", len(p.subscribers))
	}
}

// ---- Fix: 异步图片加载 ----

func findKind(rn *renderNode, k rnKind) *renderNode {
	if rn == nil {
		return nil
	}
	if rn.kind == k {
		return rn
	}
	for _, c := range rn.children {
		if r := findKind(c, k); r != nil {
			return r
		}
	}
	return nil
}

func TestAsyncImageLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.png")
	im := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for i := range im.Pix {
		im.Pix[i] = 255
	}
	im.Set(0, 0, color.RGBA{255, 0, 0, 255})
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, im); err != nil {
		t.Fatal(err)
	}
	f.Close()

	// 清理全局缓存，避免其它测试遗留影响
	delete(imgCache, path)

	g := newGame()
	g.mountRoot(Img(Src(path)))
	rn := findKind(g.rootRN, rnImage)
	if rn == nil {
		t.Fatal("image render node not found")
	}
	if rn.img != nil {
		t.Fatal("image must NOT load synchronously on the render thread")
	}

	deadline := time.Now().Add(3 * time.Second)
	for rn.img == nil && time.Now().Before(deadline) {
		drainPosts() // 后台解码完成后经 Post 回填（模拟每帧 drainPosts）
		g.drain()
		time.Sleep(5 * time.Millisecond)
	}
	if rn.img == nil {
		t.Fatal("image did not load asynchronously within timeout")
	}
	if _, ok := lookupImage(path); !ok {
		t.Fatal("loaded image was not cached")
	}
}

// 失败的加载不得阻塞、不得缓存，且节点保持未加载。
func TestAsyncImageLoadFailure(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist.png")
	delete(imgCache, missing)

	g := newGame()
	g.mountRoot(Img(Src(missing)))
	rn := findKind(g.rootRN, rnImage)
	if rn == nil {
		t.Fatal("image render node not found")
	}

	deadline := time.Now().Add(2 * time.Second)
	for imgLoading[missing] && time.Now().Before(deadline) {
		drainPosts()
		time.Sleep(5 * time.Millisecond)
	}
	drainPosts()
	if rn.img != nil {
		t.Fatal("failed load must leave node unloaded")
	}
	if _, ok := lookupImage(missing); ok {
		t.Fatal("failed load must not be cached (so it can retry later)")
	}
}
