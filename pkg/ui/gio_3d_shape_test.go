package ui

import (
	"image"
	"testing"

	"gioui.org/gpu/headless"
	"gioui.org/op"
	gpaint "gioui.org/op/paint"
)

// 伪 3D 的轮廓必须是真梯形：远端更窄。这是「形状用裁剪」的核心保证 ——
// 单靠仿射画出来的是平行四边形，对边等长，透视收缩会整个消失（曾经如此，见 git 历史）。
func TestProjectedOutlineIsTrapezoidNotParallelogram(t *testing.T) {
	win, err := headless.NewWindow(300, 240)
	if err != nil {
		t.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()

	// 纯色方块 + 绕 Y 轴旋转：左侧靠近观察者应更高，右侧更矮
	var ops op.Ops
	gpaint.ColorOp{Color: nrgba(Color{255, 255, 255, 255})}.Add(&ops)
	gpaint.PaintOp{}.Add(&ops)
	p := newGioPainter(&ops, 300, 240)
	p.BeginLayer()
	p.FillRect(90, 60, 120, 120, 0, Color{0, 0, 0, 255})
	p.EndLayer(layerTransform{
		cx: 150, cy: 120, w: 120, h: 120, scale: 1, opacity: 1,
		rotateY: 45, perspective: 300,
	})
	if err := win.Frame(&ops); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 300, 240))
	if err := win.Screenshot(img); err != nil {
		t.Fatal(err)
	}

	colHeight := func(x int) int {
		n := 0
		for y := 0; y < 240; y++ {
			if r, _, _, _ := img.At(x, y).RGBA(); r>>8 < 128 {
				n++
			}
		}
		return n
	}
	left, right := colHeight(110), colHeight(180)
	t.Logf("rotateY=45：左侧(x=110)列高=%d 右侧(x=180)列高=%d", left, right)
	if left == 0 || right == 0 {
		t.Fatalf("投影后图形缺失：left=%d right=%d", left, right)
	}
	// 平行四边形的对边等长；梯形才有收缩。留 8px 余量以容忍抗锯齿。
	if left <= right+8 {
		t.Errorf("左侧列高 %d 未明显高于右侧 %d —— 轮廓退化成平行四边形，透视收缩没了",
			left, right)
	}
}

// 记录并守住「内容对裁剪四边形的覆盖率」。
//
// 这里不要求覆盖到 1.0：仿射只能匹配三角，内容填不满裁剪是既定的近似，实测 80%~99%，
// 而且肉眼看不出（可见轮廓是两个凸四边形的交集，仍像一张卡）。本测试是防止它劣化 ——
// 覆盖率若明显跌破这个区间，说明仿射或投影算错了，卡会缺角。
//
// 判据必须是面积。曾用「黑色区间内部有没有白点」来查，那是没牙的：内容填不满时图形是
// 边缘被切短，不是中间破洞 —— 故意把内容缩小 3% 都照样通过。
func TestProjectedContentFillsClipQuad(t *testing.T) {
	win, err := headless.NewWindow(300, 300)
	if err != nil {
		t.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()

	for _, tc := range []struct{ rx, ry, pers float32 }{
		{45, 0, 900}, {45, 20, 900}, {35, 20, 700}, {60, 15, 400}, {0, 50, 500}, {30, -40, 600},
	} {
		tr := layerTransform{
			cx: 150, cy: 150, w: 150, h: 180, scale: 1, opacity: 1,
			rotateX: tc.rx, rotateY: tc.ry, perspective: tc.pers,
		}
		var ops op.Ops
		gpaint.ColorOp{Color: nrgba(Color{255, 255, 255, 255})}.Add(&ops) // 白底
		gpaint.PaintOp{}.Add(&ops)
		p := newGioPainter(&ops, 300, 300)
		p.BeginLayer()
		p.FillRect(75, 60, 150, 180, 0, Color{0, 0, 0, 255}) // 纯黑实心卡，恰好铺满元素
		p.EndLayer(tr)
		if err := win.Frame(&ops); err != nil {
			t.Fatal(err)
		}
		img := image.NewRGBA(image.Rect(0, 0, 300, 300))
		if err := win.Screenshot(img); err != nil {
			t.Fatal(err)
		}

		painted := 0
		for y := 0; y < 300; y++ {
			for x := 0; x < 300; x++ {
				if r, _, _, _ := img.At(x, y).RGBA(); r>>8 < 128 {
					painted++
				}
			}
		}

		// 裁剪梯形的解析面积（鞋带公式）
		d0 := project3D(tr, -tr.w/2, -tr.h/2)
		d1 := project3D(tr, tr.w/2, -tr.h/2)
		d3 := project3D(tr, tr.w/2, tr.h/2)
		d2 := project3D(tr, -tr.w/2, tr.h/2)
		poly := []struct{ X, Y float32 }{{d0.X, d0.Y}, {d1.X, d1.Y}, {d3.X, d3.Y}, {d2.X, d2.Y}}
		var a2 float32
		for i := range poly {
			j := (i + 1) % len(poly)
			a2 += poly[i].X*poly[j].Y - poly[j].X*poly[i].Y
		}
		want := absf(a2) / 2

		ratio := float32(painted) / want
		t.Logf("rotateX=%-4v rotateY=%-4v persp=%-4v 斜切=%5.1fpx 画出=%5d 梯形面积=%6.0f 覆盖率=%.3f",
			tc.rx, tc.ry, tc.pers, projErr(tr), painted, want, ratio)
		if ratio < 0.75 {
			t.Errorf("rotateX=%v rotateY=%v: 只覆盖了裁剪梯形的 %.1f%%（既有水平 80%%~99%%）"+
				" —— 明显劣化，卡会缺角", tc.rx, tc.ry, ratio*100)
		}
		if ratio > 1.02 {
			t.Errorf("rotateX=%v rotateY=%v: 覆盖率 %.1f%% 超过裁剪四边形 —— 裁剪没生效",
				tc.rx, tc.ry, ratio*100)
		}
	}
}
