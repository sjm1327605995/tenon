package ui

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"gioui.org/gpu/headless"
	"gioui.org/op"
)

// 地板贴图必须朝灭点收敛：远端窄、近端宽。
//
// 这是 PlaneImage 存在的全部理由。普通 Img 走仿射，把地板画成平行四边形 —— 对边等长、
// 不收敛，于是卡牌（各自精确投影）和格子对不上。判据就取「顶部一行的宽度 vs 底部一行」。
func TestPlaneImageConvergesToVanishingPoint(t *testing.T) {
	const W, H = 640, 480
	win, err := headless.NewWindow(W, H)
	if err != nil {
		t.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()

	green := image.NewRGBA(image.Rect(0, 0, 400, 300))
	draw.Draw(green, green.Bounds(), &image.Uniform{color.RGBA{0, 200, 0, 255}}, image.Point{}, draw.Src)

	render := func(plane bool) *image.RGBA {
		floor := Img(SrcImage("f", green), Style(Absolute, Left(0), Top(0), Width(400), Height(300)))
		if plane {
			floor = Img(PlaneImage("f", green), Style(Absolute, Left(0), Top(0), Width(400), Height(300)))
		}
		tree := Div(Style(Width(W), Height(H), ItemsCenter, JustifyCenter),
			Div(Style(Width(400), Height(300), Scene3D, Perspective(700), RotateX(50)), floor))
		h := Mount(Use(func(struct{}) *Node { return tree }, struct{}{}), W, H)

		var ops op.Ops
		p := newGioPainter(&ops, W, H)
		paint(p, h.g.rootRN)
		if err := win.Frame(&ops); err != nil {
			t.Fatal(err)
		}
		img := image.NewRGBA(image.Rect(0, 0, W, H))
		if err := win.Screenshot(img); err != nil {
			t.Fatal(err)
		}
		return img
	}

	// 某一行上绿色像素的宽度
	rowWidth := func(img *image.RGBA, y int) int {
		n := 0
		for x := 0; x < W; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if g>>8 > 150 && r>>8 < 100 && b>>8 < 100 {
				n++
			}
		}
		return n
	}
	// 绿色区域的上下边界
	span := func(img *image.RGBA) (top, bot int) {
		top, bot = -1, -1
		for y := 0; y < H; y++ {
			if rowWidth(img, y) > 4 {
				if top < 0 {
					top = y
				}
				bot = y
			}
		}
		return
	}

	for _, tc := range []struct {
		name    string
		plane   bool
		wantCvg bool
	}{
		{"普通 Img（仿射）", false, false},
		{"PlaneImage（单应）", true, true},
	} {
		img := render(tc.plane)
		top, bot := span(img)
		if top < 0 || bot-top < 20 {
			t.Fatalf("%s: 地板没画出来 (top=%d bot=%d)", tc.name, top, bot)
		}
		// 避开抗锯齿边缘，各内缩几像素采样
		wTop := rowWidth(img, top+3)
		wBot := rowWidth(img, bot-3)
		ratio := float64(wBot) / float64(wTop)
		t.Logf("%-20s 顶部宽=%3d 底部宽=%3d 底/顶=%.2f", tc.name, wTop, wBot, ratio)

		converged := ratio > 1.15 // 远端明显更窄才算收敛
		if converged != tc.wantCvg {
			if tc.wantCvg {
				t.Errorf("%s: 底/顶=%.2f，未朝灭点收敛 —— 地板还是平行四边形", tc.name, ratio)
			} else {
				t.Errorf("%s: 底/顶=%.2f 竟然收敛了 —— 前提不成立，本测试无法归因", tc.name, ratio)
			}
		}
	}
}

// 预变形必须只做一次：命中缓存后每帧零成本。场景参数变了才该重算 ——
// 少把一个参数拼进缓存键，窗口一缩放地板就会悄悄用上过期的变形结果、和卡牌错位。
func TestPlaneImageCachesAndInvalidates(t *testing.T) {
	planeCache = map[string]planeCacheEntry{}
	t.Cleanup(func() { planeCache = map[string]planeCacheEntry{} })

	img := image.NewRGBA(image.Rect(0, 0, 64, 48))
	base := layerTransform{cx: 200, cy: 150, w: 400, h: 300, scale: 1, opacity: 1,
		rotateX: 45, perspective: 700}

	if _, _, ok := planeBitmap("k", img, base); !ok {
		t.Fatal("首次预变形失败")
	}
	if n := len(planeCache); n != 1 {
		t.Fatalf("缓存条数=%d, want 1", n)
	}
	// 同样的参数：必须命中，不新增
	planeBitmap("k", img, base)
	if n := len(planeCache); n != 1 {
		t.Errorf("同参数重复调用后缓存条数=%d, want 1 —— 没命中缓存，每帧都在重算", n)
	}

	// 每个影响投影的参数变化都必须使缓存失效
	for _, tc := range []struct {
		name string
		mod  func(*layerTransform)
	}{
		{"倾角", func(t *layerTransform) { t.rotateX = 30 }},
		{"透视", func(t *layerTransform) { t.perspective = 900 }},
		{"尺寸（窗口缩放）", func(t *layerTransform) { t.w, t.h = 500, 375 }},
		{"位置", func(t *layerTransform) { t.cx = 250 }},
		{"绕Y", func(t *layerTransform) { t.rotateY = 10 }},
		{"缩放", func(t *layerTransform) { t.scale = 1.5 }},
	} {
		before := len(planeCache)
		tr := base
		tc.mod(&tr)
		planeBitmap("k", img, tr)
		if len(planeCache) != before+1 {
			t.Errorf("%s 变化后缓存条数仍是 %d —— 该参数没进缓存键，会用上过期的变形结果",
				tc.name, len(planeCache))
		}
	}
}
