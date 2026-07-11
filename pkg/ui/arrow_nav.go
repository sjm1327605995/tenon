package ui

// 方向键导航（roving focus，无障碍）：把 ArrowNav(...) 放到容器上，
// 组内可聚焦项即可用方向键移动焦点（菜单/列表用上下，标签页/工具条用左右）。
// Tab 仍在整树内切换；方向键只在组内循环。输入框获得焦点时方向键交给光标，不参与导航。
//
//	ui.Div(ui.ArrowNav(ui.NavVertical),
//	    MenuItem(...), MenuItem(...), MenuItem(...),
//	)

// NavOrient 是方向键导航的朝向。
type NavOrient int

const (
	NavVertical   NavOrient = iota // 上/下
	NavHorizontal                  // 左/右
	NavBoth                        // 上下左右都可
)

// ArrowNav 把容器标记为方向键导航组。
func ArrowNav(o NavOrient) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.navGroup = true; hp.navOrient = o }}
}

func orientMatch(group, want NavOrient) bool {
	return group == want || group == NavBoth
}

// moveFocusInGroup 在最近的、朝向匹配的导航组内把焦点移到上/下一个可聚焦项（环形）。
// 返回是否发生了移动（供调用方判断是否消费该按键）。输入框聚焦时不参与（方向键归光标）。
func (g *game) moveFocusInGroup(forward bool, want NavOrient) bool {
	if g.focusedFiber == nil || g.focusedFiber.rnode == nil {
		return false
	}
	cur := g.focusedFiber.rnode
	if cur.kind == rnInput {
		// 左右永远归光标；上下只有多行输入才归光标（单行输入上下可离框导航，如搜索框→结果）
		if want == NavHorizontal || (want == NavVertical && cur.multiline) {
			return false
		}
	}
	var group *renderNode
	for c := cur; c != nil; c = c.parent {
		if c.navGroup && orientMatch(c.navOrient, want) {
			group = c
			break
		}
	}
	if group == nil {
		return false
	}
	var list []*renderNode
	collectFocusables(group, &list)
	n := nextFocus(list, cur, forward)
	if n == nil || n == cur {
		return false
	}
	g.focusedFiber = n.owner
	if n.kind == rnInput {
		n.caretPos, n.selAnchor = len(n.value), len(n.value)
	}
	return true
}
