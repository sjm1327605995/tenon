package ui

import (
	"image"
	"testing"

	"gioui.org/f32"
	"gioui.org/gpu/headless"
	"gioui.org/op"
	gpaint "gioui.org/op/paint"
)

// 伪 3D 只作用于绘制：加了透视/旋转后，元素自身与其后兄弟的布局必须一模一样。
func TestPseudo3DDoesNotAffectLayout(t *testing.T) {
	build := func(extra ...StyleOpt) *Node {
		return Box([]StyleOpt{Width(400), Height(600), Padding(20), Column, Gap(16)},
			Box(append([]StyleOpt{Width(200), Height(100)}, extra...), Text("card")),
			Text("below"),
		)
	}
	flat := Mount(build(), 800, 800)
	tilt := Mount(build(Perspective(500), RotateX(35), RotateY(-15), TranslateZ(40)), 800, 800)

	for _, marker := range []string{"card", "below"} {
		fb := flat.Root().ByText(marker).Bounds()
		tb := tilt.Root().ByText(marker).Bounds()
		if fb != tb {
			t.Fatalf("%q 的布局被伪 3D 改变了：flat=%+v tilt=%+v（伪 3D 只应影响绘制）", marker, fb, tb)
		}
	}
}

// 投影的基本性质：无 3D 参数时应恒等；rotateY 让左右两侧到中心的距离不再相等（透视）。
func TestProject3DBasics(t *testing.T) {
	base := layerTransform{cx: 100, cy: 100, w: 80, h: 60, scale: 1, opacity: 1}

	// 没有 3D 时，投影只是把偏移搬到中心
	if got := project3D(base, 10, -5); got != f32.Pt(110, 95) {
		t.Fatalf("无 3D 时 project3D=%v want (110,95)", got)
	}

	// 绕 Y 轴旋转 + 透视：靠近观察者的一侧应被放大（离中心更远）
	p := base
	p.rotateY, p.perspective = 40, 400
	left := project3D(p, -40, 0)
	right := project3D(p, 40, 0)
	dl := base.cx - left.X
	dr := right.X - base.cx
	t.Logf("rotateY=40 perspective=400：左侧距中心 %.1f，右侧距中心 %.1f", dl, dr)
	if dl == dr {
		t.Fatal("绕 Y 轴旋转后左右仍对称 —— 透视没生效")
	}
}

// 投影仿射必须精确落在三个角上（第四角是近似，由裁剪兜底，不校验）。
func TestProjAffineMapsCorners(t *testing.T) {
	d0, d1, d2 := f32.Pt(10, 20), f32.Pt(40, 25), f32.Pt(12, 60)
	a := projAffine(0, 0, 10, 20, d0, d1, d2)

	for _, c := range []struct {
		src  f32.Point
		want f32.Point
		name string
	}{
		{f32.Pt(0, 0), d0, "左上"},
		{f32.Pt(10, 0), d1, "右上"},
		{f32.Pt(0, 20), d2, "左下"},
	} {
		got := a.Transform(c.src)
		if absf(got.X-c.want.X) > 0.01 || absf(got.Y-c.want.Y) > 0.01 {
			t.Errorf("%s角: 映射到 %v want %v", c.name, got, c.want)
		}
	}
}

// 端到端：带 rotateY + perspective 的卡片必须真的被投影成梯形 ——
// 靠近观察者的一侧更高、另一侧更矮；且卡片里的文字也要跟着一起变形。
func TestGio3DCardIsProjected(t *testing.T) {
	win, err := headless.NewWindow(300, 240)
	if err != nil {
		t.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()

	render := func(rotY float32) *image.RGBA {
		var ops op.Ops
		gpaint.ColorOp{Color: nrgba(Color{255, 255, 255, 255})}.Add(&ops)
		gpaint.PaintOp{}.Add(&ops)
		p := newGioPainter(&ops, 300, 240)

		p.BeginLayer()
		p.FillRect(90, 60, 120, 120, 0, Color{0, 0, 0, 255})
		p.DrawText("Hi", gioNewFont(20, 400, false), Color{255, 255, 255, 255}, 110, 100, false, false)
		p.EndLayer(layerTransform{
			cx: 150, cy: 120, w: 120, h: 120,
			scale: 1, opacity: 1,
			rotateY: rotY, perspective: 300,
		})
		if err := win.Frame(&ops); err != nil {
			t.Fatalf("frame: %v", err)
		}
		img := image.NewRGBA(image.Rect(0, 0, 300, 240))
		if err := win.Screenshot(img); err != nil {
			t.Fatalf("screenshot: %v", err)
		}
		return img
	}

	// 某一列上黑色像素的高度
	colHeight := func(img *image.RGBA, x int) int {
		n := 0
		for y := 0; y < 240; y++ {
			r, _, _, _ := img.At(x, y).RGBA()
			if r>>8 < 128 {
				n++
			}
		}
		return n
	}

	flat := render(0)
	if h := colHeight(flat, 150); h < 100 {
		t.Fatalf("未旋转时卡片中列高=%d，want>=100（前提不成立）", h)
	}

	// rotateY=45 时卡片横向收缩（约 x=101..187），采样必须落在投影后的图形内部。
	// 靠近观察者的一侧（左）应更高，远离的一侧（右）更矮 —— 这正是透视。
	tilt := render(45)
	left, right := colHeight(tilt, 110), colHeight(tilt, 180)
	t.Logf("rotateY=45：左侧(x=110)列高=%d 右侧(x=180)列高=%d", left, right)
	if left == 0 || right == 0 {
		t.Fatalf("投影后卡片缺失：left=%d right=%d", left, right)
	}
	if left <= right {
		t.Fatalf("左侧列高 %d 未高于右侧 %d —— 没有透视收缩，投影没生效", left, right)
	}
}

// 卡片里的文字必须跟着一起投影变形。这正是「重放整棵子树」相对「逐图元投影」的价值所在：
// 字形轮廓由 shaper 给出（clip.PathSpec），拿不到顶点、无法自己重投影，
// 但整棵子树可以被套上投影仿射重放一遍。
func TestGio3DWarpsTextToo(t *testing.T) {
	win, err := headless.NewWindow(300, 240)
	if err != nil {
		t.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()

	// 只画白底+黑字（不画卡片底色），统计文字像素的水平跨度
	render := func(rotY float32) *image.RGBA {
		var ops op.Ops
		gpaint.ColorOp{Color: nrgba(Color{255, 255, 255, 255})}.Add(&ops)
		gpaint.PaintOp{}.Add(&ops)
		p := newGioPainter(&ops, 300, 240)
		p.BeginLayer()
		p.DrawText("MMMM", gioNewFont(28, 700, false), Color{0, 0, 0, 255}, 100, 100, false, false)
		p.EndLayer(layerTransform{
			cx: 150, cy: 120, w: 120, h: 60,
			scale: 1, opacity: 1,
			rotateY: rotY, perspective: 300,
		})
		if err := win.Frame(&ops); err != nil {
			t.Fatalf("frame: %v", err)
		}
		img := image.NewRGBA(image.Rect(0, 0, 300, 240))
		if err := win.Screenshot(img); err != nil {
			t.Fatalf("screenshot: %v", err)
		}
		return img
	}
	span := func(img *image.RGBA) int {
		lo, hi := 300, -1
		for x := 0; x < 300; x++ {
			for y := 0; y < 240; y++ {
				if r, _, _, _ := img.At(x, y).RGBA(); r>>8 < 128 {
					if x < lo {
						lo = x
					}
					if x > hi {
						hi = x
					}
					break
				}
			}
		}
		if hi < 0 {
			return 0
		}
		return hi - lo
	}

	flat, tilt := span(render(0)), span(render(50))
	t.Logf("文字水平跨度：未旋转=%d，rotateY=50=%d", flat, tilt)
	if flat == 0 || tilt == 0 {
		t.Fatalf("文字缺失：flat=%d tilt=%d", flat, tilt)
	}
	if tilt >= flat {
		t.Fatalf("旋转后文字跨度 %d 未收窄（未旋转 %d）—— 文字没有跟着投影", tilt, flat)
	}
}
