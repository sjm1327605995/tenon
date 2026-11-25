package ui

// isDisabled 按渲染父链计算有效禁用态，使容器禁用也能约束其后代。
func (rn *renderNode) isDisabled() bool {
	for c := rn; c != nil; c = c.parent {
		if c.disabled {
			return true
		}
	}
	return false
}

// interactiveHit 保留禁用节点的命中遮挡，但不把它作为可激活目标返回。
func (rn *renderNode) scrollBy(delta float32) bool {
	if !rn.scroll {
		return false
	}
	maxScroll := max(rn.contentH-rn.bounds.H, float32(0))
	next := min(max(rn.scrollY+delta, float32(0)), maxScroll)
	changed := next != rn.scrollY
	rn.scrollY = next
	return changed
}

func (g *game) interactiveHit(x, y float32) *renderNode {
	n := g.hitTop(x, y)
	if n != nil && n.isDisabled() {
		return nil
	}
	return n
}

// UseInteraction 管理一个元素的悬停/按压态，返回状态和要挂到该元素上的属性。
// 用法：
//
//	hovered, pressed, ia := ui.UseInteraction()
//	return ui.Button(ui.Style(base, ui.StyleIf(hovered, ...)), ia, children...)
//
// 返回的 attrs 用 Spread 展开挂到元素上（也可作为单个 *Node 直接传入变参）。
func UseInteraction() (hovered, pressed bool, attrs *Node) {
	h, setH := UseState(false)
	p, setP := UseState(false)
	return h, p, Attrs(OnHover(setH), OnPress(setP))
}

// Attrs 把多个属性打包成一个 *Node（挂载时透明展开），便于把一组属性作为单值传递。
func Attrs(nodes ...*Node) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) {
		for _, n := range nodes {
			if n != nil && n.typ == typeAttr && n.applyAttr != nil {
				n.applyAttr(hp)
			}
		}
	}}
}
