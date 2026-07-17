package ui

import "testing"

// Portal 的浮层根节点是全屏且不裁剪的容器，它自己不是可命中的目标：
// 浮层空白处的点必须穿透到主树。
//
// 否则会形成自激：悬停触发器 -> 弹出浮层(Tooltip/下拉) -> 全屏浮层根吞掉命中 ->
// 触发器 unhover -> 浮层收起 -> 又命中到触发器 -> 再弹出……每帧一次，表现为 hover 闪烁。
func TestPortalOverlayDoesNotSwallowHits(t *testing.T) {
	enters, leaves := 0, 0
	trigger := func(_ struct{}) *Node {
		hov, setHov := UseState(false)
		kids := []*Node{
			Style(Width(80), Height(36)),
			OnHover(func(in bool) {
				if in {
					enters++
				} else {
					leaves++
				}
				setHov(in)
			}),
			Text("悬停我"),
		}
		if hov { // 悬停时弹出浮层（Tooltip 的做法）：浮在触发器下方，不遮挡触发器本身
			kids = append(kids, Portal(
				Div(Style(Absolute, Left(0), Top(100), Width(120), Height(30)), Text("提示")),
			))
		}
		return Div(kids...)
	}
	h := Mount(Use(trigger, struct{}{}), 400, 300)
	g := h.g

	frame := func(x, y float32) {
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
	}

	frame(40, 18) // 悬停到触发器上（触发器占 0..80 x 0..36）
	if enters != 1 {
		t.Fatalf("首次悬停 enters=%d want 1", enters)
	}

	// 浮层弹出后，光标一直停在触发器内部并小幅移动：不应反复进出
	for i := 0; i < 20; i++ {
		frame(40+float32(i%10), 18)
	}
	t.Logf("浮层弹出后在触发器内移动 20 帧：enters=%d leaves=%d", enters, leaves)
	if enters != 1 || leaves != 0 {
		t.Fatalf("浮层弹出后 hover 自激振荡：enters=%d leaves=%d（期望 1/0）—— "+
			"全屏浮层根吞掉了主树的命中", enters, leaves)
	}
}
