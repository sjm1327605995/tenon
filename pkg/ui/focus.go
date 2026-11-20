package ui

// focusNode 是键盘与程序化聚焦的共同入口，保持选区和滚动位置一致。
func (g *game) focusNode(n *renderNode) {
	g.visualVersion++
	if n == nil || !n.focusable || n.isDisabled() {
		g.focusedFiber = nil
		return
	}
	g.focusedFiber = n.owner
	if n.kind == rnInput {
		n.caretPos, n.selAnchor = len(n.value), len(n.value)
	}
	g.revealFocus(n)
}

// focusClicked 聚焦命中节点或最近的可聚焦祖先（例如按钮里的文字/图标）。
func (g *game) focusClicked(n *renderNode) {
	for c := n; c != nil; c = c.parent {
		if c.focusable {
			g.focusNode(c)
			return
		}
	}
	g.focusedFiber = nil
}

// revealFocus 逐层滚动，确保 Tab/方向键移到的控件处于可视范围。
// 这里只改变滚动偏移和绝对 bounds，不使 Yoga 布局失效。
func (g *game) revealFocus(n *renderNode) {
	target := n.bounds
	for c := n.parent; c != nil; c = c.parent {
		if !c.scroll {
			continue
		}
		var delta float32
		if target.Y < c.bounds.Y || target.H > c.bounds.H {
			delta = target.Y - c.bounds.Y
		} else if bottom := target.Y + target.H; bottom > c.bounds.Y+c.bounds.H {
			delta = bottom - c.bounds.Y - c.bounds.H
		}
		next := min(max(c.scrollY+delta, 0), max(c.contentH-c.bounds.H, 0))
		if next == c.scrollY {
			continue
		}
		target.Y -= next - c.scrollY
		c.scrollY = next
		g.boundsDirty, g.needsLayout = true, true
	}
}

// clearDisabledInteractions 在控件或祖先被禁用时清理旧焦点与指针捕获。
func (g *game) clearDisabledInteractions() {
	if f := g.focusedFiber; f != nil && f.rnode != nil && f.rnode.isDisabled() {
		f.rnode.composeLo, f.rnode.composeHi = 0, 0
		f.rnode.selAnchor = f.rnode.caretPos
		g.focusedFiber = nil
		g.inputSelecting, g.imeComposing = false, false
	}
	if rn := g.pressedNode; rn != nil && rn.isDisabled() {
		g.pressedNode = nil
		if rn.onPress != nil {
			rn.onPress(false)
		}
	}
	if rn := g.dragging; rn != nil && rn.isDisabled() {
		g.dragging = nil
	}
	if g.rootRN != nil {
		for level := g.maxLevel; ; level-- {
			g.rootRN.clearDisabledHoverAtLevel(level)
			if level == 0 {
				break
			}
		}
	}
	for _, portal := range g.portals {
		if portal.overlayRoot != nil {
			for level := g.maxLevel; ; level-- {
				portal.overlayRoot.clearDisabledHoverAtLevel(level)
				if level == 0 {
					break
				}
			}
		}
	}
}
