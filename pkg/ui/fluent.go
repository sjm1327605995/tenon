package ui

import "fmt"

// With 返回附加属性后的新节点；原节点与其子节点切片保持不变。
// 函数组件的属性施加到其返回的单个根节点，不增加布局容器。
// Fragment/Portal/Provider 只有一个子节点时透明传递；多个根时应修饰内部的容器。
func (n *Node) With(attrs ...*Node) *Node {
	if n == nil || len(attrs) == 0 {
		return n
	}
	copyNode := *n
	var applied []*Node
	for _, attr := range attrs {
		if attr == nil {
			continue
		}
		if attr.typ != typeAttr {
			panic("ui.Node.With: only attributes are accepted; pass children to the constructor")
		}
		if attr.key != "" {
			copyNode.key = attr.key
		}
		if attr.trap {
			if n.typ != typePortal {
				panic("ui.Node.With: TrapFocus requires a Portal")
			}
			copyNode.trap = true
		}
		if attr.applyAttr != nil && attr.key == "" {
			applied = append(applied, attr)
		}
	}
	if len(applied) == 0 {
		return &copyNode
	}
	switch n.typ {
	case typeHost, typeText, typeIcon:
		copyNode.attrList = make([]func(*hostProps), len(n.attrList), len(n.attrList)+len(applied))
		copy(copyNode.attrList, n.attrList)
		for _, attr := range applied {
			copyNode.attrList = append(copyNode.attrList, attr.applyAttr)
		}
	case typeComponent:
		// 修饰也是组件输入。稳定复用同一个修饰节点时 Memo 仍可短路；
		// 每次新建修饰时保守重渲染，以免吞掉变化的样式或捕获新值的回调。
		source := n
		if previous, ok := n.props.(*modifiedProps); ok {
			source = previous.source
			merged := make([]*Node, 0, len(previous.attrs)+len(applied))
			merged = append(merged, previous.attrs...)
			applied = append(merged, applied...)
		}
		copyNode.props = &modifiedProps{source: source, attrs: applied}
		copyNode.render = renderModified
	case typeAttr:
		copyNode.applyAttr = func(hp *hostProps) {
			if n.applyAttr != nil {
				n.applyAttr(hp)
			}
			for _, attr := range applied {
				attr.applyAttr(hp)
			}
		}
	case typeFragment, typeProvider, typePortal:
		if len(n.kids) > 1 {
			panic("ui.Node.With: multiple visual roots; modify an inner Div instead")
		}
		if len(n.kids) == 1 {
			copyNode.kids = []*Node{n.kids[0].With(applied...)}
		}
	default:
		panic(fmt.Sprintf("ui.Node.With: node type %d has no single visual root; modify an inner Div instead", n.typ))
	}
	return &copyNode
}

type modifiedProps struct {
	source *Node
	attrs  []*Node
}

func renderModified(props any) *Node {
	p := props.(*modifiedProps)
	return p.source.render(p.source.props).With(p.attrs...)
}

// Style 附加样式，后写的选项覆盖前面的同名选项。
func (n *Node) Style(opts ...StyleOpt) *Node { return n.With(Style(opts...)) }

// Width 设置逻辑像素宽度。
func (n *Node) Width(v float32) *Node { return n.Style(Width(v)) }

// Height 设置逻辑像素高度。
func (n *Node) Height(v float32) *Node { return n.Style(Height(v)) }

// Padding 设置四边内边距。
func (n *Node) Padding(v float32) *Node { return n.Style(Padding(v)) }

// Gap 设置子元素间距。
func (n *Node) Gap(v float32) *Node { return n.Style(Gap(v)) }

// Bg 设置背景色。
func (n *Node) Bg(c Color) *Node { return n.Style(Bg(c)) }

// FontSize 设置字号，容器的字号由后代继承。
func (n *Node) FontSize(v float32) *Node { return n.Style(FontSize(v)) }

// TextColor 设置文字或图标颜色，容器的颜色由后代继承。
func (n *Node) TextColor(c Color) *Node { return n.Style(TextColor(c)) }

// Radius 设置圆角半径。
func (n *Node) Radius(v float32) *Node { return n.Style(Radius(v)) }

// OnClick 设置点击回调，覆盖节点已有的点击回调。
func (n *Node) OnClick(fn func()) *Node { return n.With(OnClick(fn)) }

// Bind 同时设置输入框的受控值与 OnChange 回调。
func (n *Node) Bind(value string, onChange func(string)) *Node {
	return n.With(Value(value), OnChange(onChange))
}

// Placeholder 设置输入框的占位文字。
func (n *Node) Placeholder(text string) *Node { return n.With(Placeholder(text)) }

// OnSubmit 设置单行输入框的回车提交回调。
func (n *Node) OnSubmit(fn func(string)) *Node { return n.With(OnSubmit(fn)) }

// Disabled 设置有效禁用态；外观仍由 Style 或组件自身的 Disabled props 决定。
func (n *Node) Disabled(disabled bool) *Node { return n.With(Disabled(disabled)) }

// Key 返回带稳定身份的新节点，可用于列表重排。
func (n *Node) Key(key string) *Node {
	if n == nil {
		return nil
	}
	copyNode := *n
	copyNode.key = key
	return &copyNode
}
