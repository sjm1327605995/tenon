package ui

import "testing"

// ZIndex 必须同时决定「画在上面」和「点得到」—— 两者分岔就是最坏的 bug：
// 看得见的那个点不到。所以这里一并断言绘制顺序与命中结果。
func TestZIndexOrdersPaintAndHitTogether(t *testing.T) {
	// 三个完全重叠的方块，兄弟顺序 a,b,c（默认 c 在最上）
	build := func(za, zb, zc int) *Harness {
		mk := func(name string, z int) *Node {
			return Div(Style(Absolute, Left(0), Top(0), Width(100), Height(100), ZIndex(z)),
				Text(name))
		}
		return Mount(Use(func(struct{}) *Node {
			return Div(Style(Width(200), Height(200)),
				Div(Style(Absolute, Left(0), Top(0), Width(100), Height(100)),
					mk("a", za), mk("b", zb), mk("c", zc)))
		}, struct{}{}), 200, 200)
	}

	// 谁最后被画（=在最上面）
	topPainted := func(h *Harness) string {
		last := ""
		for _, op := range h.Paint() {
			if op.Kind == "text" {
				last = op.Text
			}
		}
		return last
	}
	// 点重叠区域命中谁：命中最深的节点，往上找到带文字的那个方块
	topHit := func(h *Harness) string {
		for c := h.g.hitTop(50, 50); c != nil; c = c.parent {
			if c.kind == rnText {
				return c.text
			}
			for _, ch := range c.children {
				if ch.kind == rnText {
					return ch.text
				}
			}
		}
		return ""
	}

	for _, tc := range []struct {
		name       string
		za, zb, zc int
		want       string
	}{
		{"全为 0：兄弟顺序决定，最后的在上", 0, 0, 0, "c"},
		{"a 抬高", 5, 0, 0, "a"},
		{"b 抬高", 0, 5, 0, "b"},
		{"负值把 c 压到最下", 0, 0, -1, "b"},
		{"同值时仍按兄弟顺序", 3, 3, 0, "b"},
	} {
		h := build(tc.za, tc.zb, tc.zc)
		painted, hit := topPainted(h), topHit(h)
		t.Logf("%-24s z=(%d,%d,%d) 画在最上=%q 点中=%q", tc.name, tc.za, tc.zb, tc.zc, painted, hit)
		if painted != tc.want {
			t.Errorf("%s: 最上层画的是 %q，want %q", tc.name, painted, tc.want)
		}
		if hit != painted {
			t.Errorf("%s: 画在最上的是 %q 却点中了 %q —— 绘制与命中顺序分岔", tc.name, painted, hit)
		}
	}
}
