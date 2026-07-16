package ui

import "testing"

// 命中测试应与绘制一致：只有设了 Clip 的容器才裁剪后代。
// 溢出到「未裁剪」父容器之外的子元素是看得见的，就必须点得到。
func TestHitOverflowOutsideNonClippingParent(t *testing.T) {
	clicked := false
	app := func(_ struct{}) *Node {
		return Div(Style(Width(400), Height(200)),
			// 父容器 100x100、未设 Clip；子元素绝对定位到父容器之外
			Div(Style(Width(100), Height(100)),
				Div(Style(Absolute, Left(150), Top(10), Width(50), Height(50)),
					OnClick(func() { clicked = true }), Text("out")),
			),
		)
	}
	h := Mount(Use(app, struct{}{}), 400, 200)
	g := h.g

	q := h.Root().ByText("out")
	b := q.Bounds()
	t.Logf("溢出子元素 bounds=%+v（未裁剪，应当可见可点）", b)

	n := g.hitTop(b.X+b.W/2, b.Y+b.H/2)
	if n == nil {
		t.Fatal("溢出子元素完全命中不到")
	}
	for c := n; c != nil; c = c.parent {
		if c.onClick != nil {
			c.onClick()
			break
		}
	}
	if !clicked {
		t.Fatal("溢出到未裁剪父容器之外的元素点不到 —— 命中测试比绘制更严格")
	}
}

// 设了 Clip 的容器：滚出/溢出视口的后代不可见，也不应被点到。
func TestHitRespectsClip(t *testing.T) {
	clicked := false
	app := func(_ struct{}) *Node {
		return Div(Style(Width(400), Height(200)),
			Div(Style(Width(100), Height(100), Clip),
				Div(Style(Absolute, Left(150), Top(10), Width(50), Height(50)),
					OnClick(func() { clicked = true }), Text("hidden")),
			),
		)
	}
	h := Mount(Use(app, struct{}{}), 400, 200)
	g := h.g
	b := h.Root().ByText("hidden").Bounds()

	if n := g.hitTop(b.X+b.W/2, b.Y+b.H/2); n != nil {
		for c := n; c != nil; c = c.parent {
			if c.onClick != nil {
				c.onClick()
				break
			}
		}
	}
	if clicked {
		t.Fatal("被 Clip 裁掉的元素不可见，却被点到了")
	}
}
