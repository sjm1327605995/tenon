package ui

import "testing"

// elemErr 量化逐元素投影的斜切：元素四角投影的第四角，与「元素自己的三角仿射映过去」差多少。
// 与 projErr（整层）同算法，只是把 contentAffine 换成 elementContentAffine。
func elemErr(t layerTransform, r Rect) float32 {
	et := t
	et.cx = r.X + r.W/2
	et.cy = r.Y + r.H/2
	et.w, et.h = r.W, r.H
	_, _, _, d3 := projCorners(et)
	got := elementContentAffine(t, r).transform(pt{r.X + r.W, r.Y + r.H})
	return absf(got.X-d3.X) + absf(got.Y-d3.Y)
}

// 逐元素投影的价值：层内小元素（一行文字、一张小图）的斜切误差必须从「整层尺寸级」
// 降到「元素尺寸级」。640x420 的桌面 50° 曾整层斜 211px（见 node.go 的 PlaneImage 文档），
// 桌面上 150x20 的一行文字必须只剩亚像素级。
func TestElementAffineShrinksSkew(t *testing.T) {
	base := layerTransform{cx: 320, cy: 210, w: 640, h: 420, perspective: 800, rotateX: 50}
	layerErr := projErr(base)
	if layerErr < 50 {
		t.Fatalf("整层斜切只有 %.1fpx —— 前提不成立，本测试证明不了逐元素的价值", layerErr)
	}
	t.Logf("整层（640x420 @50°）斜切 %.1fpx", layerErr)

	for _, r := range []Rect{
		{X: 100, Y: 100, W: 150, H: 20}, // 一行文字
		{X: 300, Y: 250, W: 80, H: 110}, // 一张卡
		{X: 20, Y: 30, W: 40, H: 14},    // 一个小标签
	} {
		e := elemErr(base, r)
		t.Logf("  元素 %gx%g：斜切 %.2fpx", r.W, r.H, e)
		// 判据：比整层单仿射小一个数量级，且绝对值在肉眼不可辨的范围内。
		// 元素尺寸级误差是仿射近似的固有下限（铺满层的大图请用 PlaneImage）。
		if e > layerErr/10 || e > 12 {
			t.Errorf("元素 %gx%g 斜切 %.2fpx，应为整层的 1/10 以下且 <12px", r.W, r.H, e)
		}
	}
}

// 结构黄金测试：带 3D 变换的层必须按「元素」录制（层自身 + 每个直接子元素各一个宏），
// 且元素矩形就是子元素的 bounds —— 回放时各自套用三角仿射的前提。
func TestPaintProjectedRecordsElements(t *testing.T) {
	tree := Div(Style(Width(400), Height(300), Perspective(600), RotateX(40), Column, Gap(8)),
		Div(Style(Width(200), Height(40), Bg(Hex("#111111")))),
		Div(Style(Width(200), Height(40), Bg(Hex("#222222")))),
		Text("label"))
	h := Mount(Use(func(struct{}) *Node { return tree }, struct{}{}), 500, 500)
	ops := h.Paint()

	var kinds []string
	var elements []Rect
	inProj := false
	for _, op := range ops {
		kinds = append(kinds, op.Kind)
		switch op.Kind {
		case "projlayer":
			inProj = true
		case "unprojlayer":
			inProj = false
		case "projelement":
			if inProj {
				elements = append(elements, op.Rect)
			}
		}
	}
	if !contains(kinds, "projlayer") || !contains(kinds, "unprojlayer") {
		t.Fatalf("缺少投影层录制：kinds=%v", kinds)
	}
	// 期望 5 个元素：层自身 + 两个子 Div + 文本 + 收尾元素；层自身与收尾元素同为层 bounds。
	if len(elements) != 5 {
		t.Fatalf("投影元素数=%d，want 5（层自身 + 3 子元素 + 收尾元素）：%v", len(elements), elements)
	}
	if elements[0] != elements[4] {
		t.Errorf("元素0（层自身）=%v 与末元素=%v 应同为层 bounds", elements[0], elements[4])
	}
	if elements[1].W != 200 || elements[1].H != 40 {
		t.Errorf("子元素矩形=%v，want 200x40（子元素录制时用其 bounds）", elements[1])
	}
	// 纯 2D 层不走投影路径。
	flat := Mount(Use(func(struct{}) *Node {
		return Div(Style(Width(400), Height(300), Opacity(0.5)), Div(Style(Width(10), Height(10))))
	}, struct{}{}), 500, 500)
	for _, op := range flat.Paint() {
		if op.Kind == "projlayer" {
			t.Fatal("纯 2D 层不应走投影录制路径")
		}
	}
}

// 逐元素命中：点击子元素被投影后的画面位置，必须命中它自己；而「整层单仿射」会算出的
// 旧位置在大倾角下与子元素真实画面位置差几十 px —— 点旧位置不应命中（证明命中确实
// 用的是逐元素仿射，而不是碰巧位移小全中）。
func TestHitUsesElementAffine(t *testing.T) {
	clicked := ""
	tree := Div(Style(Width(900), Height(900), ItemsCenter, JustifyCenter),
		Div(Style(Width(600), Height(400), Column, Gap(12), ItemsCenter, JustifyCenter,
			Perspective(700), RotateX(50)),
			Div(Style(Width(150), Height(30)), OnClick(func() { clicked = "row" }), Text("row")),
			Div(Style(Width(500), Height(30)), OnClick(func() { clicked = "wide" }), Text("wide"))))
	h := Mount(Use(func(struct{}) *Node { return tree }, struct{}{}), 900, 900)

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
	lt := layerOf(board, nil)

	// 点每一行被逐元素仿射映到的中心：必须命中。
	for _, name := range []string{"row", "wide"} {
		rn := h.Root().ByText(name).rn.parent
		cb := rn.bounds
		ctr := elementContentAffine(lt, cb).transform(pt{cb.X + cb.W/2, cb.Y + cb.H/2})
		clicked = ""
		if n := hitNode(h.g.rootRN, ctr.X, ctr.Y); n != nil {
			for x := n; x != nil; x = x.parent {
				if x.onClick != nil {
					x.onClick()
					break
				}
			}
		}
		if clicked != name {
			t.Errorf("点 %s 的投影中心却命中 %q —— 逐元素命中没生效", name, clicked)
		}
	}

	// 对照：整层单仿射给「row」算出的画面中心，与逐元素仿射差多少？
	rn := h.Root().ByText("row").rn.parent
	cb := rn.bounds
	cElem := elementContentAffine(lt, cb).transform(pt{cb.X + cb.W/2, cb.Y + cb.H/2})
	cLayer := contentAffine(lt).transform(pt{cb.X + cb.W/2, cb.Y + cb.H/2})
	d := absf(cElem.X-cLayer.X) + absf(cElem.Y-cLayer.Y)
	t.Logf("两仿射对同一中心差 %.1fpx", d)
	if d < 2 {
		t.Fatalf("对照位移只有 %.1fpx —— 证明不了逐元素命中在用不同的矩阵", d)
	}
	// 点「整层单仿射」的旧位置：不应命中 row（它被画在逐元素位置）。
	clicked = ""
	if n := hitNode(h.g.rootRN, cLayer.X, cLayer.Y); n != nil {
		for x := n; x != nil; x = x.parent {
			if x.onClick != nil {
				x.onClick()
				break
			}
		}
	}
	if clicked == "row" {
		t.Error("点整层单仿射的旧位置竟命中 row —— 命中与绘制可能用了不同的矩阵")
	}
}
