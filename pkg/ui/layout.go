package ui

// 布局便捷构造：都是对 Div + Style 的薄封装，用来消除最常见的样板
//
//	Div(Style(Row, ItemsCenter, Gap(8)), a, b)
//
// 需要额外样式时，在子节点里再传一个 Style(...)——它会与预置样式合并（后者覆盖前者），
// 例如 HStack(8, Style(JustifyBetween), a, b)。

// HStack 是横向 flex 容器（行、交叉轴居中、子项间距 gap）。
func HStack(gap float32, kids ...*Node) *Node {
	return Div(append([]*Node{Style(Row, ItemsCenter, Gap(gap))}, kids...)...)
}

// VStack 是纵向 flex 容器（列、子项间距 gap）。
func VStack(gap float32, kids ...*Node) *Node {
	return Div(append([]*Node{Style(Column, Gap(gap))}, kids...)...)
}

// Center 让子节点在主轴与交叉轴上都居中。
func Center(kids ...*Node) *Node {
	return Div(append([]*Node{Style(ItemsCenter, JustifyCenter)}, kids...)...)
}

// Spacer 是可伸缩占位（flex-grow:1），把两侧的兄弟节点推向容器两端。
func Spacer() *Node { return Div(Style(Grow(1))) }
