package ui

import (
	"image"
	"testing"

	"gioui.org/gpu/headless"
	"gioui.org/op"
	gpaint "gioui.org/op/paint"
)

// 光标在同一元素内部移动时，连续帧渲染出的画面必须完全一致（hover 高亮不应忽有忽无）。
// 这直接复现「在边框内移动会闪烁」的现象：引擎层的 hover 状态是稳的，所以要看真实像素。
func TestHoverNoVisualFlickerWhileMovingInside(t *testing.T) {
	win, err := headless.NewWindow(300, 200)
	if err != nil {
		t.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()

	comp := func(_ struct{}) *Node {
		hov, setHov := UseState(false)
		st := []StyleOpt{Width(200), Height(100), Bg(Hex("#ffffff"))}
		if hov {
			st = append(st, Bg(Hex("#333333")))
		}
		return Box(st, OnHover(setHov), Text("hi"))
	}
	h := Mount(Use(comp, struct{}{}), 300, 200)
	g := h.g

	// 驱动一帧并渲染出像素
	shot := func(x, y float32) *image.RGBA {
		t.Helper()
		gioIn.setCursor(x, y)
		g.handleInput()
		for guard := 0; len(g.dirty) > 0 && guard < 100; guard++ {
			g.flushDirty()
		}
		if g.needsLayout {
			g.rootRN = rootRenderNode(g.rootFiber)
			g.layout()
			g.needsLayout = false
		}
		var ops op.Ops
		gpaint.ColorOp{Color: nrgba(Color{247, 248, 250, 255})}.Add(&ops)
		gpaint.PaintOp{}.Add(&ops)
		p := newGioPainter(&ops, g.w, g.h)
		paint(p, g.rootRN)
		if err := win.Frame(&ops); err != nil {
			t.Fatalf("frame: %v", err)
		}
		img := image.NewRGBA(image.Rect(0, 0, 300, 200))
		if err := win.Screenshot(img); err != nil {
			t.Fatalf("screenshot: %v", err)
		}
		return img
	}

	base := shot(100, 50) // 进入元素，hover 高亮出现
	if darkPixels(base, 20, 180, 20, 80) == 0 {
		t.Fatal("hover 高亮没出现，测试前提不成立")
	}

	// 在元素内部移动，每帧都应画出与上一帧完全相同的画面
	prev := base
	for i, pos := range [][2]float32{{102, 52}, {104, 48}, {98, 55}, {110, 50}, {100, 50}} {
		cur := shot(pos[0], pos[1])
		if diff := diffPixels(prev, cur); diff != 0 {
			t.Errorf("第 %d 次内部移动到 (%v,%v)：与上一帧有 %d 个像素不同 —— 画面在闪",
				i+1, pos[0], pos[1], diff)
		}
		prev = cur
	}
}

func diffPixels(a, b *image.RGBA) int {
	n := 0
	for i := range a.Pix {
		if a.Pix[i] != b.Pix[i] {
			n++
		}
	}
	return n
}
