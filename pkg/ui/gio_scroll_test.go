package ui

import "testing"

// 滚动量的单位契约：inputSource.wheel() 与 renderNode.bounds 同为物理像素，引擎直接相加。
//
// 曾经的做法是后端把 gio 的像素滚动量除以 24 伪造成「档位」，引擎再乘回 24*uiScale ——
// 一来一回不但多余，还因为 gio 的 Scroll 本就是物理像素而多乘了一次 uiScale，
// 导致 150% 缩放的屏幕上滚动快 1.5 倍（100% 屏上恰好正确，所以一直没被发现）。
func TestScrollAppliedInPhysicalPixelsRegardlessOfDPI(t *testing.T) {
	app := func(_ struct{}) *Node {
		return ScrollView(nil, Style(Width(100), Height(100)),
			Div(Style(Width(100), Height(1000))), // 内容比视口高，可滚动
		)
	}
	for _, scale := range []float32{1, 1.5, 2} {
		h := Mount(Use(app, struct{}{}), 200, 200)
		g := h.g
		uiScale = scale // 模拟不同 DPI

		sc := findScrollNode(g.rootRN)
		if sc == nil {
			t.Fatal("没找到可滚动节点")
		}
		before := sc.scrollY

		gioIn.setCursor(50, 50)
		gioIn.wheelY = 120 // 物理像素
		g.handleInput()
		gioIn.wheelY = 0

		if got := sc.scrollY - before; got != 120 {
			t.Errorf("uiScale=%v 时滚动了 %v 像素，want 120（滚动量不应随 DPI 缩放）", scale, got)
		}
	}
	uiScale = 1
}

func findScrollNode(rn *renderNode) *renderNode {
	if rn == nil {
		return nil
	}
	if rn.scroll {
		return rn
	}
	for _, c := range rn.children {
		if r := findScrollNode(c); r != nil {
			return r
		}
	}
	return nil
}
