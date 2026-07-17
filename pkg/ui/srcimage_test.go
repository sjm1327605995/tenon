package ui

import (
	"image"
	"image/color"
	"testing"
)

func solidImage(w, h int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

// SrcImage 应当同步装好位图：既无 IO 也无解码，不该像 Src 那样等一次 Post 才可见，
// 否则首帧会闪一下空白。
func TestSrcImageInstallsSynchronously(t *testing.T) {
	resetImageCache(t)
	img := solidImage(8, 12, color.RGBA{255, 0, 0, 255})

	h := Mount(Img(SrcImage("mem:red", img)), 100, 100)

	q := h.Root().ByKind("image")
	if !q.Exists() {
		t.Fatal("树里没有 img 节点")
	}
	if q.rn.img == nil {
		t.Fatal("SrcImage 未同步安装位图：首帧会闪空白")
	}
	if w, hh := q.rn.img.Size(); w != 8 || hh != 12 {
		t.Fatalf("位图尺寸 = %dx%d, 期望 8x12", w, hh)
	}
}

// 同一 key 的多个节点必须共用一张位图。否则每个节点各建一份，缓存预算形同虚设 ——
// 决斗盘上同名卡、卡背会同时出现很多次。
func TestSrcImageSharesBitmapByKey(t *testing.T) {
	resetImageCache(t)
	img := solidImage(4, 4, color.RGBA{0, 255, 0, 255})

	h := Mount(Div(
		Img(SrcImage("mem:green", img)),
		Img(SrcImage("mem:green", img)),
	), 100, 100)

	imgs := h.Root().FindAll(func(q *Query) bool { return q.rn.kind == rnImage && q.rn.img != nil })
	if len(imgs) != 2 {
		t.Fatalf("找到 %d 个已装图的 img 节点, 期望 2", len(imgs))
	}
	if imgs[0].rn.img != imgs[1].rn.img {
		t.Error("同一 key 的两个节点建了两张位图，未共用缓存")
	}
	if n, _, _ := ImageCacheStats(); n != 1 {
		t.Errorf("缓存条数 = %d, 期望 1", n)
	}
}

// 内容变了就换 key —— 同一个节点的 key 变了，必须真的换图。
//
// 必须让同一棵树重渲染来改 key，不能 Mount 两个独立 harness 各用一个 key：那样两个节点
// 都是全新的（rn.imgSrc 为空），走的是首次安装、而不是「key 变化」这条路径。实测把重建
// 逻辑整个打断，两个独立 Mount 的写法照样全过 —— 那种测试对它宣称的行为毫无约束力。
func TestSrcImageKeyChangeSwapsBitmap(t *testing.T) {
	resetImageCache(t)
	small := solidImage(4, 4, color.RGBA{0, 0, 255, 255})
	big := solidImage(16, 16, color.RGBA{0, 0, 255, 255})

	app := func(_ struct{}) *Node {
		big2, setBig := UseState(false)
		key, img := "mem:v1", small
		if big2 {
			key, img = "mem:v2", big
		}
		return Div(Style(Width(100), Height(100)),
			Img(SrcImage(key, img), Style(Width(50), Height(50))),
			Div(Style(Width(20), Height(20)), OnClick(func() { setBig(true) }), Text("go")),
		)
	}
	h := Mount(Use(app, struct{}{}), 100, 100)

	imgW := func() int {
		q := h.Root().ByKind("image")
		if !q.Exists() || q.rn.img == nil {
			t.Fatal("img 节点没装上图")
		}
		w, _ := q.rn.img.Size()
		return w
	}
	if w := imgW(); w != 4 {
		t.Fatalf("初始位图宽 = %d, 期望 4", w)
	}

	b := h.Root().ByText("go").Bounds()
	h.ClickAt(b.X+b.W/2, b.Y+b.H/2) // 触发重渲染，key: v1 -> v2

	if w := imgW(); w != 16 {
		t.Errorf("key 从 v1 换到 v2 后位图宽 = %d, 期望 16 —— 同一节点没有换图", w)
	}
}
