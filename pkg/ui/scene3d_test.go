package ui

import (
	"image"
	"testing"

	"gioui.org/gpu/headless"
	"gioui.org/op"
	gpaint "gioui.org/op/paint"
)

// 共享相机的全部意义在于「灭点一致」，而它可被观测的后果是：
// 同一场景里，元素的投影形状取决于它在场景中的位置 —— 远端的卡更小。
//
// 没有相机时每张卡绕自己中心投影，摆在哪里形状都一样（这正是摆一桌不像同一张桌子的原因）。
// 所以本测试的判据是「左右两张卡的高度必须不同」，并用 A/B（不加 Scene3D）证明
// 这个差异确实由相机带来。
func TestScene3DSharesVanishingPoint(t *testing.T) {
	win, err := headless.NewWindow(600, 400)
	if err != nil {
		t.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()

	// scene=true: 桌子是 Scene3D，卡是直接子元素，不自带 3D
	// scene=false: 桌子不倾斜，每张卡自己带同样的 3D（旧模型）
	build := func(scene bool) *renderNode {
		mk := func() *Node {
			st := []StyleOpt{Width(90), Height(120), Bg(Hex("#000000"))}
			if !scene {
				st = append(st, Perspective(400), RotateY(40))
			}
			return Div(Style(st...))
		}
		boardSt := []StyleOpt{Width(500), Height(300), Row, JustifyBetween, ItemsCenter,
			PaddingXY(20, 0)}
		if scene {
			boardSt = append(boardSt, Scene3D, Perspective(400), RotateY(40))
		}
		tree := Div(Style(Width(600), Height(400), ItemsCenter, JustifyCenter),
			Div(Style(boardSt...), mk(), mk()))
		h := Mount(Use(func(struct{}) *Node { return tree }, struct{}{}), 600, 400)
		return h.g.rootRN
	}

	// 某一列上黑色像素的高度
	measure := func(rn *renderNode) (leftH, rightH int) {
		var ops op.Ops
		gpaint.ColorOp{Color: nrgba(Color{255, 255, 255, 255})}.Add(&ops)
		gpaint.PaintOp{}.Add(&ops)
		p := newGioPainter(&ops, 600, 400)
		paint(p, rn)
		if err := win.Frame(&ops); err != nil {
			t.Fatal(err)
		}
		img := image.NewRGBA(image.Rect(0, 0, 600, 400))
		if err := win.Screenshot(img); err != nil {
			t.Fatal(err)
		}
		colH := func(x int) int {
			n := 0
			for y := 0; y < 400; y++ {
				if r, _, _, _ := img.At(x, y).RGBA(); r>>8 < 128 {
					n++
				}
			}
			return n
		}
		// 取两张卡各自最高的那一列（避开投影后的水平偏移）
		best := func(lo, hi int) int {
			m := 0
			for x := lo; x < hi; x++ {
				if h := colH(x); h > m {
					m = h
				}
			}
			return m
		}
		return best(60, 290), best(310, 540)
	}

	sl, sr := measure(build(true))
	nl, nr := measure(build(false))
	t.Logf("Scene3D（共享相机）: 左卡高=%d 右卡高=%d  差=%d", sl, sr, sl-sr)
	t.Logf("各自投影（无相机）: 左卡高=%d 右卡高=%d  差=%d", nl, nr, nl-nr)

	if sl == 0 || sr == 0 {
		t.Fatalf("场景里的卡没画出来：left=%d right=%d", sl, sr)
	}
	// A：无相机时两张卡形状必须一样 —— 位置不影响投影
	if nl != nr {
		t.Errorf("无相机时左右卡高 %d/%d 竟然不同 —— 前提不成立，本测试无法归因", nl, nr)
	}
	// B：有相机时必须不同 —— 位置影响投影，才谈得上灭点一致
	if sl == sr {
		t.Fatalf("Scene3D 下左右卡高都是 %d —— 位置没影响投影，相机没生效", sl)
	}
	// rotateY=40：右侧远离观察者，应更矮
	if sl <= sr {
		t.Errorf("Scene3D 下左卡 %d 未高于右卡 %d —— 远近关系反了", sl, sr)
	}
}

// 相机只影响绘制，不能动布局。
func TestScene3DDoesNotAffectLayout(t *testing.T) {
	build := func(scene bool) *Harness {
		st := []StyleOpt{Width(500), Height(300), Row, Gap(20), ItemsCenter, JustifyCenter}
		if scene {
			st = append(st, Scene3D, Perspective(500), RotateX(45))
		}
		return Mount(Use(func(struct{}) *Node {
			return Div(Style(Width(600), Height(400), ItemsCenter, JustifyCenter),
				Div(Style(st...),
					Div(Style(Width(90), Height(120)), Text("a")),
					Div(Style(Width(90), Height(120)), Text("b"))))
		}, struct{}{}), 600, 400)
	}
	flat, sc := build(false), build(true)
	for _, m := range []string{"a", "b"} {
		fb, sb := flat.Root().ByText(m).Bounds(), sc.Root().ByText(m).Bounds()
		if fb != sb {
			t.Errorf("%q 的布局被 Scene3D 改变了：flat=%+v scene=%+v", m, fb, sb)
		}
	}
}
