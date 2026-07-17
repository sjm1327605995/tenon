package ui

import (
	"fmt"
	"testing"
)

// 命中测试必须跟随伪 3D 投影：卡画在哪就得点在哪。
//
// 曾经不跟随（invTransform 只反解了 rotate(Z)/scale/translate），实测 45° 时 25% 点错、
// 60° 时 50%。倾角小的时候「碰巧全中」是假象 —— 位移还没超过卡牌半高，点偏了也落回同一张卡。
//
// 3D 必须放在牌桌上、而不是每张卡自带：元素绕自身中心投影时中心是不动点，
// 采样卡牌中心位移恒为 0，测了等于没测（本测试第一版正是这么写错的）。
func TestHitFollows3DProjection(t *testing.T) {
	const rows, cols = 4, 5
	const cardW, cardH = 80, 110

	for _, tilt := range []float32{15, 30, 45, 60} {
		clicked := ""
		mk := func(r int) *Node {
			cards := []*Node{Style(Row, Gap(10), JustifyCenter)}
			for c := 0; c < cols; c++ {
				n := fmt.Sprintf("r%dc%d", r, c)
				cards = append(cards, Div(Style(Width(cardW), Height(cardH)),
					OnClick(func() { clicked = n }), Text(n)))
			}
			return Div(cards...)
		}
		rowNodes := []*Node{Style(Width(700), Height(560), Column, Gap(10),
			ItemsCenter, JustifyCenter, Perspective(900), RotateX(tilt))}
		for r := 0; r < rows; r++ {
			rowNodes = append(rowNodes, mk(r))
		}
		tree := Div(Style(Width(900), Height(900), ItemsCenter, JustifyCenter), Div(rowNodes...))
		h := Mount(Use(func(struct{}) *Node { return tree }, struct{}{}), 900, 900)

		// 牌桌（带 3D 的那层）的仿射决定卡被画到哪
		var board *renderNode
		var find func(*renderNode)
		find = func(r *renderNode) {
			if r.has3D() {
				board = r
			}
			for _, c := range r.children {
				find(c)
			}
		}
		find(h.g.rootRN)
		if board == nil {
			t.Fatal("没找到带 3D 的牌桌")
		}
		aff := contentAffine(layerOf(board, nil))

		miss, maxShift := 0, float32(0)
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				name := fmt.Sprintf("r%dc%d", r, c)
				rn := h.Root().ByText(name).rn.parent
				// 卡牌中心实际被画到哪：用绘制所用的同一个仿射
				cb := rn.bounds
				ctr := aff.transform(pt{cb.X + cb.W/2, cb.Y + cb.H/2})
				if d := absf(ctr.Y - (cb.Y + cb.H/2)); d > maxShift {
					maxShift = d
				}
				clicked = ""
				n := h.g.hitTop(ctr.X, ctr.Y) // 点「看到的位置」
				for x := n; x != nil; x = x.parent {
					if x.onClick != nil {
						x.onClick()
						break
					}
				}
				if clicked != name {
					miss++
				}
			}
		}
		t.Logf("rotateX=%-3v 最大投影位移=%5.1fpx 点错=%d/%d", tilt, maxShift, miss, rows*cols)
		if miss > 0 {
			t.Errorf("rotateX=%v: %d/%d 张卡点不中 —— 命中没跟随投影", tilt, miss, rows*cols)
		}
	}
}

// Scene3D 下更要紧：相机会把卡挪得更远（远端那张既缩小又位移）。
func TestHitFollowsScene3DCamera(t *testing.T) {
	clicked := ""
	mk := func(i int) *Node {
		n := fmt.Sprintf("c%d", i)
		return Div(Style(Width(90), Height(120)), OnClick(func() { clicked = n }), Text(n))
	}
	cards := []*Node{Style(Width(500), Height(300), Row, JustifyBetween, ItemsCenter,
		PaddingXY(20, 0), Scene3D, Perspective(400), RotateY(40))}
	for i := 0; i < 4; i++ {
		cards = append(cards, mk(i))
	}
	tree := Div(Style(Width(600), Height(400), ItemsCenter, JustifyCenter), Div(cards...))
	h := Mount(Use(func(struct{}) *Node { return tree }, struct{}{}), 600, 400)

	// 找到场景，取相机
	var scene *renderNode
	var find func(*renderNode)
	find = func(r *renderNode) {
		if r.scene3D {
			scene = r
		}
		for _, c := range r.children {
			find(c)
		}
	}
	find(h.g.rootRN)
	if scene == nil {
		t.Fatal("没找到 Scene3D 场景")
	}
	cam := cameraOf(scene)

	for i := 0; i < 4; i++ {
		name := fmt.Sprintf("c%d", i)
		rn := h.Root().ByText(name).rn.parent
		cb := rn.bounds
		ctr := contentAffine(layerOf(rn, cam)).transform(pt{cb.X + cb.W/2, cb.Y + cb.H/2})
		clicked = ""
		n := h.g.hitTop(ctr.X, ctr.Y)
		for x := n; x != nil; x = x.parent {
			if x.onClick != nil {
				x.onClick()
				break
			}
		}
		t.Logf("%s 布局中心=(%.0f,%.0f) 画到=(%.0f,%.0f) 位移=%.0fpx -> 命中=%q",
			name, cb.X+cb.W/2, cb.Y+cb.H/2, ctr.X, ctr.Y,
			absf(ctr.X-(cb.X+cb.W/2)), clicked)
		if clicked != name {
			t.Errorf("点击 %s 画面上的位置却命中 %q —— 命中没跟随相机", name, clicked)
		}
	}
}
