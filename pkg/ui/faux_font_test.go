package ui

import (
	"image"
	"testing"

	"gioui.org/gpu/headless"
	"gioui.org/op"
	gpaint "gioui.org/op/paint"
)

// 内置字体集只有一张 OPPOSans Medium，gio 不会自己合成粗体/斜体 —— 它只挑最接近的那张。
// 所以必须由我们合成，否则 ui.Bold 完全没有效果（曾经如此：Bold 与 Regular 的字形 ID
// 与宽度完全相同，painter 接口的 fauxBold 参数被 gio 后端整个忽略）。
func TestFauxBoldAndItalicAreRequested(t *testing.T) {
	for _, tc := range []struct {
		weight           int
		italic           bool
		wantBold, wantIt bool
	}{
		{400, false, false, false},
		{500, false, false, false}, // Medium：还不到合成阈值
		{600, false, true, false},  // Semibold 起
		{700, false, true, false},
		{400, true, false, true},
		{700, true, true, true},
	} {
		f := gioNewFont(20, tc.weight, tc.italic)
		_, fb, fi := f.Metrics()
		if fb != tc.wantBold || fi != tc.wantIt {
			t.Errorf("weight=%d italic=%v -> fauxBold=%v fauxItalic=%v, want %v/%v",
				tc.weight, tc.italic, fb, fi, tc.wantBold, tc.wantIt)
		}
	}
}

// 合成必须真的改变像素。只断言 Metrics 返回 true 是不够的 ——
// 那只说明「请求了」，不说明「画出来了」（这正是之前 fauxBold 参数被忽略的样子）。
func TestFauxBoldAndItalicChangePixels(t *testing.T) {
	const W, H = 240, 80
	win, err := headless.NewWindow(W, H)
	if err != nil {
		t.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()

	render := func(weight int, italic bool) *image.RGBA {
		var ops op.Ops
		gpaint.ColorOp{Color: nrgba(Color{255, 255, 255, 255})}.Add(&ops)
		gpaint.PaintOp{}.Add(&ops)
		p := newGioPainter(&ops, W, H)
		f := gioNewFont(28, weight, italic)
		_, fb, fi := f.Metrics()
		p.DrawText("Hamburg", f, Color{0, 0, 0, 255}, 10, 20, fb, fi)
		if err := win.Frame(&ops); err != nil {
			t.Fatal(err)
		}
		img := image.NewRGBA(image.Rect(0, 0, W, H))
		if err := win.Screenshot(img); err != nil {
			t.Fatal(err)
		}
		return img
	}
	ink := func(img *image.RGBA) int {
		n := 0
		for i := 0; i < len(img.Pix); i += 4 {
			if img.Pix[i] < 128 {
				n++
			}
		}
		return n
	}
	// 最左侧墨迹的 x：斜体会让上半部右移
	leftmostAt := func(img *image.RGBA, y int) int {
		for x := 0; x < W; x++ {
			if r, _, _, _ := img.At(x, y).RGBA(); r>>8 < 128 {
				return x
			}
		}
		return -1
	}

	reg, bold, ital := render(400, false), render(700, false), render(400, true)
	rInk, bInk := ink(reg), ink(bold)
	t.Logf("常规墨迹=%d 粗体墨迹=%d (%.2fx)", rInk, bInk, float64(bInk)/float64(rInk))

	if rInk == 0 {
		t.Fatal("常规字一个墨点都没有")
	}
	if bInk <= rInk {
		t.Errorf("粗体墨迹 %d 不多于常规 %d —— 合成粗体没生效，Bold 与 Regular 长得一样", bInk, rInk)
	}
	if float64(bInk)/float64(rInk) > 2.0 {
		t.Errorf("粗体墨迹是常规的 %.2f 倍 —— 描边太宽，糊成一团", float64(bInk)/float64(rInk))
	}

	// 斜体绕基线剪切：基线以上向右倾。所以要在「字形上部」采样 ——
	// 采样行必须动态找（写死行号很容易落在空白上，本测试第一版就是这么错的）。
	firstInkRow := func(img *image.RGBA) int {
		for y := 0; y < H; y++ {
			if leftmostAt(img, y) >= 0 {
				return y
			}
		}
		return -1
	}
	row := firstInkRow(reg) + 2 // 字形顶部略下一点，避开抗锯齿
	topReg, topIt := leftmostAt(reg, row), leftmostAt(ital, row)
	t.Logf("字形上部 y=%d 处最左墨迹：常规 x=%d 斜体 x=%d", row, topReg, topIt)
	if topReg < 0 || topIt < 0 {
		t.Fatalf("采样行 y=%d 上没有墨迹，前提不成立", row)
	}
	if topIt <= topReg {
		t.Errorf("斜体在 y=%d 的最左墨迹 x=%d 未右移（常规 x=%d）—— 合成斜体没生效",
			row, topIt, topReg)
	}
}
