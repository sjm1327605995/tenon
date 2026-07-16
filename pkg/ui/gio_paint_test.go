package ui

import (
	"image"
	"testing"

	"gioui.org/gpu/headless"
	"gioui.org/op"
)

// renderOps 走真实 gio 绘制路径把一组原语渲染到离屏图，用于像素级断言。
func renderOps(t *testing.T, w, h int, draw func(p *gioPainter)) *image.RGBA {
	t.Helper()
	win, err := headless.NewWindow(w, h)
	if err != nil {
		t.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()
	var ops op.Ops
	p := newGioPainter(&ops, w, h)
	draw(p)
	if err := win.Frame(&ops); err != nil {
		t.Fatalf("frame: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if err := win.Screenshot(img); err != nil {
		t.Fatalf("screenshot: %v", err)
	}
	return img
}

func darkPixels(img *image.RGBA, x0, x1, y0, y1 int) int {
	n := 0
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r>>8 < 200 || g>>8 < 200 || b>>8 < 200 {
				n++
			}
		}
	}
	return n
}

// 回归：Radius 远大于尺寸（Radius(9999) 即 shadcn 的“全圆角/胶囊”写法）不得破坏整帧绘制。
// gio 的 clip.RRect 不钳制半径，半径超过短边一半会让路径退化，曾导致同一帧内其它内容全部消失
// （现象：侧边栏里有 Switch 时整个侧边栏空白）。rrectOp 必须把半径钳到 min(w,h)/2。
func TestGioHugeRadiusDoesNotBreakFrame(t *testing.T) {
	black := Color{0, 0, 0, 255}
	img := renderOps(t, 200, 120, func(p *gioPainter) {
		p.FillRect(0, 0, 200, 120, 0, Color{255, 255, 255, 255})
		p.DrawText("Before", gioNewFont(16, 400, false), black, 10, 10, false, false)
		// 全圆角胶囊（半径 >> 尺寸）
		p.FillRect(10, 50, 32, 18, 9999, Color{20, 20, 20, 255})
		p.DrawText("After", gioNewFont(16, 400, false), black, 10, 85, false, false)
	})
	if n := darkPixels(img, 0, 200, 10, 45); n == 0 {
		t.Error("胶囊之前绘制的文字消失了（圆角半径未钳制会破坏整帧）")
	}
	if n := darkPixels(img, 0, 200, 50, 70); n == 0 {
		t.Error("全圆角胶囊自身没画出来")
	}
	if n := darkPixels(img, 0, 200, 85, 118); n == 0 {
		t.Error("胶囊之后绘制的文字消失了")
	}
}

// 圆角必须真的画成圆角：填充与描边在角上都应留白，直边中点则必须有色。
func TestGioRoundedCorners(t *testing.T) {
	black := Color{0, 0, 0, 255}
	// 填充：60x60、半径 20 的黑块画在 (20,20)
	fill := renderOps(t, 100, 100, func(p *gioPainter) {
		p.FillRect(0, 0, 100, 100, 0, Color{255, 255, 255, 255})
		p.FillRect(20, 20, 60, 60, 20, black)
	})
	if n := darkPixels(fill, 21, 25, 21, 25); n != 0 {
		t.Errorf("填充左上角应被圆角切掉，却有 %d 个深色像素（圆角没生效）", n)
	}
	if n := darkPixels(fill, 45, 55, 21, 25); n == 0 {
		t.Error("填充上边中点应有颜色")
	}
	// 描边：同样的圆角矩形，2px 边框
	stroke := renderOps(t, 100, 100, func(p *gioPainter) {
		p.FillRect(0, 0, 100, 100, 0, Color{255, 255, 255, 255})
		p.StrokeRect(20, 20, 60, 60, 20, 2, black)
	})
	if n := darkPixels(stroke, 21, 25, 21, 25); n != 0 {
		t.Errorf("描边左上角应被圆角切掉，却有 %d 个深色像素（边框圆角没生效）", n)
	}
	if n := darkPixels(stroke, 45, 55, 19, 23); n == 0 {
		t.Error("描边上边中点应有边框")
	}
}

// clip 内绘制多段文字（侧边栏的典型形态）：每段都必须可见。
func TestGioMultipleTextsInsideClip(t *testing.T) {
	black := Color{0, 0, 0, 255}
	img := renderOps(t, 300, 200, func(p *gioPainter) {
		p.FillRect(0, 0, 300, 200, 0, Color{255, 255, 255, 255})
		p.PushClip(Rect{X: 0, Y: 0, W: 240, H: 200}, 0)
		p.DrawText("First", gioNewFont(16, 400, false), black, 10, 10, false, false)
		p.DrawText("Second", gioNewFont(15, 600, false), black, 10, 50, false, false)
		p.DrawText("Third", gioNewFont(12, 500, false), black, 10, 90, false, false)
		p.DrawText("Fourth", gioNewFont(14, 400, false), black, 10, 130, false, false)
		p.PopClip()
	})
	for i, band := range [][2]int{{10, 45}, {50, 85}, {90, 120}, {130, 165}} {
		if n := darkPixels(img, 0, 300, band[0], band[1]); n == 0 {
			t.Errorf("clip 内第 %d 段文字没渲染出来", i+1)
		}
	}
}

// 图层（BeginLayer/EndLayer 用 op.Record 宏）内的绘制、以及图层之后的绘制都必须可见。
func TestGioLayerRenders(t *testing.T) {
	img := renderOps(t, 100, 100, func(p *gioPainter) {
		p.FillRect(0, 0, 100, 100, 0, Color{255, 255, 255, 255})
		p.BeginLayer()
		p.FillRect(10, 10, 30, 30, 0, Color{0, 0, 0, 255})
		p.EndLayer(layerTransform{cx: 25, cy: 25, scale: 1, opacity: 1})
		p.FillRect(60, 60, 30, 30, 0, Color{0, 0, 0, 255}) // 图层之后
	})
	if n := darkPixels(img, 10, 40, 10, 40); n == 0 {
		t.Error("图层内容没画出来")
	}
	if n := darkPixels(img, 60, 90, 60, 90); n == 0 {
		t.Error("图层之后的绘制丢失")
	}
}
