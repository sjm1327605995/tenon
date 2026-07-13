package ui

import "strconv"

// VirtualListProps 配置一个虚拟滚动列表。
type VirtualListProps struct {
	Count      int               // 总行数
	ItemHeight float32           // 每行固定高度（逻辑 px）
	Height     float32           // 视口高度（逻辑 px；列表可滚动区域）
	Render     func(i int) *Node // 渲染第 i 行的内容
	Overscan   int               // 视口上下额外多渲染的行数（默认 3），减少快速滚动时的空白
}

// VirtualList 是定高行的虚拟滚动列表：无论总行数多大，都只渲染视口附近的少量行，
// 用上/下占位撑出正确的总高度与滚动条。适合上千/上万行的长列表。
//
//	ui.VirtualList(ui.VirtualListProps{
//	    Count: 10000, ItemHeight: 28, Height: 320,
//	    Render: func(i int) *ui.Node { return ui.Text(fmt.Sprintf("行 %d", i)) },
//	})
func VirtualList(p VirtualListProps) *Node { return Use(virtualList, p) }

func virtualList(p VirtualListProps) *Node {
	ref, info := UseScroll()

	ih := p.ItemHeight
	if ih <= 0 {
		ih = 1
	}
	over := p.Overscan
	if over <= 0 {
		over = 3
	}
	viewport := info.Viewport
	if viewport <= 0 { // 首帧尚未测得视口，先用配置高度估算
		viewport = p.Height
	}

	first := int(info.Offset/ih) - over
	if first < 0 {
		first = 0
	}
	visible := int(viewport/ih) + over*2 + 1
	last := first + visible
	if last > p.Count {
		last = p.Count
	}
	if first > last {
		first = last
	}

	// 直接作为 ScrollView 的子节点：上占位 + 可视行（keyed）+ 下占位，
	// 三段高度之和恒为 Count*ItemHeight，保证滚动位置与滚动条稳定。
	kids := make([]*Node, 0, (last-first)+3)
	kids = append(kids, ref, Style(Height(p.Height)))
	kids = append(kids, Div(Style(Height(float32(first)*ih)))) // 上占位
	for i := first; i < last; i++ {
		kids = append(kids, Keyed(strconv.Itoa(i),
			Div(Style(Height(ih), WidthPct(100)), p.Render(i))))
	}
	kids = append(kids, Div(Style(Height(float32(p.Count-last)*ih)))) // 下占位
	return ScrollView(kids...)
}
