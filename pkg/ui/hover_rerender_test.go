package ui

import "testing"

// 组件持续重渲染（真实应用里常有动画在跑）时，在元素内部移动光标不应让 hover 抖动。
//
// g.hovered 用 *renderNode 指针做 key，一旦重渲染会重建 renderNode，指针就失效，
// 于是每次重算悬停链都变成「leave 旧节点 + enter 新节点」—— 表现为一动就闪。
// 光标静止时因为 updateHover 会提前返回而看不出来，所以必须「持续重渲染 + 移动」一起测。
func TestHoverStableWhileRerenderingAndMoving(t *testing.T) {
	enters, leaves := 0, 0
	var lastNode *renderNode
	nodeChanges := 0

	comp := func(_ struct{}) *Node {
		UseElapsed() // 持续动画：每帧都把本组件标脏 -> 每帧重渲染
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
		drainPosts()
		g.handleInput()
		g.tickAnims(16)
		g.tickLoops(16) // 推进 UseElapsed -> 标脏
		for guard := 0; len(g.dirty) > 0 && guard < 100; guard++ {
			g.flushDirty()
		}
		if g.needsLayout {
			g.rootRN = rootRenderNode(g.rootFiber)
			g.layout()
			g.needsLayout = false
		}
		// 记录承载 onHover 的那个 renderNode 是否被换掉了
		if n := g.hitTop(x, y); n != nil {
			for c := n; c != nil; c = c.parent {
				if c.onHover != nil {
					if lastNode != nil && c != lastNode {
						nodeChanges++
					}
					lastNode = c
					break
				}
			}
		}
	}

	frame(100, 50)
	if enters != 1 || leaves != 0 {
		t.Fatalf("进入时 enters=%d leaves=%d want 1/0", enters, leaves)
	}

	// 持续重渲染的同时，在元素内部移动 30 帧
	for i := 0; i < 30; i++ {
		frame(100+float32(i%20), 50+float32(i%8))
	}
	t.Logf("持续重渲染 + 内部移动 30 帧：enters=%d leaves=%d，onHover 节点被换掉 %d 次",
		enters, leaves, nodeChanges)
	if enters != 1 || leaves != 0 {
		t.Fatalf("重渲染期间在元素内移动导致 hover 抖动：enters=%d leaves=%d（期望 1/0）", enters, leaves)
	}
}
