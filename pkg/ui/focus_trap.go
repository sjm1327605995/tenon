package ui

// 焦点陷阱（模态无障碍）：把 TrapFocus() 作为 Portal 的属性，该浮层存在期间
// Tab/Shift+Tab 只在其内部循环，不会跑到背景内容上。配合 UseEscape（Esc 关闭）
// 与遮罩（拦截点击）构成完整模态。多个模态时，最上层（绘制在最上）的生效。
//
//	ui.Portal(ui.TrapFocus(),
//	    ui.Div(...dialog content with focusable elements...),
//	)
func TrapFocus() *Node {
	return &Node{typ: typeAttr, trap: true}
}

// trapScope 返回当前生效的焦点陷阱范围（最上层带 TrapFocus 的浮层根），无则返回 nil。
func (g *game) trapScope() *renderNode {
	for i := len(g.portals) - 1; i >= 0; i-- {
		if pf := g.portals[i]; pf.portalTrap && pf.overlayRoot != nil {
			return pf.overlayRoot
		}
	}
	return nil
}
