package ui

import "testing"

// 在同一个元素内部移动光标时，hover 不应反复进出 —— 只要指针没离开元素，
// 就应当只在进入时触发一次。
//
// 注意只能挂一个 OnHover：属性是覆盖语义，挂两个的话后一个会顶掉前一个，
// 于是 setState 不触发、重渲染路径根本没被测到。这里让同一个回调既计数又改状态。
func TestHoverStableWhileMovingInsideElement(t *testing.T) {
	enters, leaves, renders := 0, 0, 0
	comp := func(_ struct{}) *Node {
		renders++
		hov, setHov := UseState(false)
		st := []StyleOpt{Width(200), Height(100)}
		if hov {
			st = append(st, Bg(Hex("#eeeeee")))
		}
		return Box(st,
			OnHover(func(in bool) {
				if in {
					enters++
				} else {
					leaves++
				}
				setHov(in)
			}),
			Text("hi"),
		)
	}
	h := Mount(Use(comp, struct{}{}), 300, 200)
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

	frame(100, 50) // 进入元素中心
	if enters != 1 || leaves != 0 {
		t.Fatalf("进入时 enters=%d leaves=%d want 1/0", enters, leaves)
	}
	rendersAfterEnter := renders

	// 在元素内部小幅移动 40 帧，全程都在框内（元素占 0..200 x 0..100）
	for i := 0; i < 40; i++ {
		frame(100+float32(i%20)*2, 50+float32(i%10))
	}
	t.Logf("内部移动 40 帧后：enters=%d leaves=%d，期间重渲染 %d 次",
		enters, leaves, renders-rendersAfterEnter)
	if enters != 1 || leaves != 0 {
		t.Fatalf("在元素内部移动导致 hover 抖动：enters=%d leaves=%d（期望始终 1/0）", enters, leaves)
	}
}
