package ui

import "reflect"

type nodeType int

const (
	typeHost nodeType = iota
	typeText
	typeComponent
	typeAttr
	typeFragment
	typeProvider
	typePortal
	typeIcon
)

// Node 是不可变的 UI 描述（相当于 React Element / VNode）。
// host 元素、文本、组件、属性统一用 Node 表示，可在变参里自由穿插。
type Node struct {
	typ nodeType

	// host
	tag      string
	attrList []func(*hostProps)
	kids     []*Node
	key      string

	// text
	text      string
	textStyle StyleProps
	runs      []textRun // 富文本：多段混排样式（RichText）

	// icon（SVG 图标）
	iconPath   string
	iconSize   float32
	iconStroke float32 // >0 描边宽度（viewBox 单位），0 为填充

	// component
	fnPtr      uintptr
	props      any
	render     func(any) *Node
	propsEqual func(a, b any) bool
	memo       bool

	// provider
	ctxID    int
	ctxValue any

	// attr
	applyAttr func(*hostProps)
	trap      bool // TrapFocus()：作为 Portal 属性时标记该浮层为焦点陷阱（模态）
}

// hostProps 是把某个 host 元素的所有属性归拢后的结果。
type hostProps struct {
	style   StyleProps
	onClick func()
	class   string
	id      string
	key     string

	// input
	value       string
	hasValue    bool
	onChange    func(string)
	placeholder string
	multiline   bool
	onHover     func(bool)
	onPress     func(bool)
	onDrag      func(dx, dy float32)
	measure     *measureHook
	scrollRef   *scrollHook // UseScroll：把该 ScrollView 的滚动状态写回

	// 方向键导航组（ArrowNav）：组内可聚焦项用方向键移动焦点
	navGroup  bool
	navOrient NavOrient

	// img
	src       string
	objectFit ObjectFit
}

func buildHostProps(n *Node) hostProps {
	hp := hostProps{style: newStyleProps()}
	for _, ap := range n.attrList {
		ap(&hp)
	}
	return hp
}

// el 构造 host 元素：属性 Node 进 attrList，其余进 kids。
func el(tag string, args []*Node) *Node {
	n := &Node{typ: typeHost, tag: tag}
	for _, a := range args {
		if a == nil {
			continue
		}
		if a.typ == typeAttr {
			if a.applyAttr != nil {
				n.attrList = append(n.attrList, a.applyAttr)
			}
			if a.key != "" {
				n.key = a.key
			}
			continue
		}
		n.kids = append(n.kids, a)
	}
	return n
}

// ---- host 元素构造器（HTML 风格）----

func Div(args ...*Node) *Node    { return el("div", args) }
func Span(args ...*Node) *Node   { return el("span", args) }
func Button(args ...*Node) *Node { return el("button", args) }
func Input(args ...*Node) *Node  { return el("input", args) }
func Img(args ...*Node) *Node    { return el("img", args) }

// ScrollView 是可垂直滚动的容器（内容超出时裁剪 + 滚轮滚动 + 滚动条）。
func ScrollView(args ...*Node) *Node { return el("scroll", args) }

// Fragment 分组多个子节点而不产生 host 元素（相当于 React.Fragment）。
func Fragment(args ...*Node) *Node { return grouped(typeFragment, args) }

// Portal 把子树渲染到顶层浮层，脱离父级裁剪与层叠顺序（相当于 React 的 createPortal）。
// 子内容在全屏尺寸的独立布局根中排布，绘制/命中都在主树之上。
func Portal(args ...*Node) *Node { return grouped(typePortal, args) }

func grouped(t nodeType, args []*Node) *Node {
	n := &Node{typ: t}
	for _, a := range args {
		if a == nil {
			continue
		}
		if a.typ == typeAttr {
			if a.key != "" {
				n.key = a.key
			}
			if a.trap {
				n.trap = true
			}
			continue
		}
		n.kids = append(n.kids, a)
	}
	return n
}

// Text 是文本节点，可选自带样式（颜色/字号）。
func Text(s string, opts ...StyleOpt) *Node {
	st := newStyleProps()
	for _, o := range opts {
		o(&st)
	}
	return &Node{typ: typeText, text: s, textStyle: st}
}

// RichText 在单个文本节点内混排多段不同样式的文字（富文本 span）。
// 每个子节点用 Text(...) 描述一段，其文字与样式（颜色/字号/字重/斜体）被收集为一段 run；
// 整体作为一个段落统一折行、按基线对齐，未显式设置的样式继承自容器。
//
//	ui.RichText(
//	    ui.Text("Hello "),
//	    ui.Text("world", ui.Bold, ui.TextColor(ui.Hex("#ef4444"))),
//	    ui.Text(" 你好", ui.FontSize(20)),
//	)
func RichText(spans ...*Node) *Node {
	rt := &Node{typ: typeText}
	for _, s := range spans {
		if s == nil || s.typ != typeText {
			continue
		}
		if len(s.runs) > 0 { // 允许嵌套 RichText：展平其 runs
			rt.runs = append(rt.runs, s.runs...)
			continue
		}
		rt.runs = append(rt.runs, textRun{text: s.text, style: s.textStyle})
	}
	return rt
}

// iconViewBox 是内置图标假定的 SVG viewBox 边长（lucide/heroicons 等均为 24）。
const iconViewBox = 24

// Icon 渲染一个线性（描边）SVG 图标：d 是 24×24 viewBox 下的 path data，size 是目标边长（逻辑 px）。
// 颜色继承自文本颜色（类似 SVG 的 currentColor），可用 TextColor 覆盖。匹配 lucide/feather 等图标集。
func Icon(d string, size float32, opts ...StyleOpt) *Node { return iconNode(d, size, 2, opts) }

// IconFill 渲染填充（实心）SVG 图标（Material/Heroicons-solid 等单一填充路径）。
func IconFill(d string, size float32, opts ...StyleOpt) *Node { return iconNode(d, size, 0, opts) }

func iconNode(d string, size, stroke float32, opts []StyleOpt) *Node {
	st := newStyleProps()
	for _, o := range opts {
		o(&st)
	}
	return &Node{typ: typeIcon, iconPath: d, iconSize: size, iconStroke: stroke, textStyle: st}
}

// ---- 属性构造器 ----

func Class(v string) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.class = v }}
}
func Id(v string) *Node { return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.id = v }} }
func OnClick(fn func()) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.onClick = fn }}
}
func Key(v string) *Node {
	return &Node{typ: typeAttr, key: v, applyAttr: func(hp *hostProps) { hp.key = v }}
}
func Style(opts ...StyleOpt) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) {
		for _, o := range opts {
			o(&hp.style)
		}
	}}
}

// Value 设置受控输入的当前值。
func Value(v string) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.value = v; hp.hasValue = true }}
}

// OnChange 在输入内容变化时回调（受控组件）。
func OnChange(fn func(string)) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.onChange = fn }}
}

// OnHover 在指针进入/离开元素时回调（true=进入，false=离开）。
func OnHover(fn func(bool)) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.onHover = fn }}
}

// OnPress 在元素上按下/松开左键时回调（true=按下，false=松开），用于按压态样式。
func OnPress(fn func(bool)) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.onPress = fn }}
}

// OnDrag 在元素上按住左键拖动时逐帧回调（dx,dy 为本帧屏幕位移）。
func OnDrag(fn func(dx, dy float32)) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.onDrag = fn }}
}

// Placeholder 设置输入占位符。
func Placeholder(v string) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.placeholder = v }}
}

// Multiline 让 Input 成为多行文本域：自动折行、按 Enter 换行、高度随内容增长。
func Multiline() *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.multiline = true }}
}

// Src 设置图片来源（文件路径）。
func Src(v string) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.src = v }}
}

// ObjectFit 决定图片如何适配其框：FitFill 拉伸（默认）/ FitContain 完整放入留白 / FitCover 填满并裁剪。
type ObjectFit int

const (
	FitFill    ObjectFit = iota // 拉伸填满（默认，可能变形）
	FitContain                  // 等比缩放完整放入，多余处留白
	FitCover                    // 等比缩放填满，超出部分裁掉
)

// Fit 设置 Img 的适配方式（object-fit）。
func Fit(mode ObjectFit) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.objectFit = mode }}
}

// Use 挂载一个带 typed props 的函数组件。
func Use[P any](fn func(P) *Node, props P) *Node {
	return component(fn, props, false)
}

// Memo 与 Use 相同，但当 props 深度相等时跳过重渲染（React.memo 语义）。
func Memo[P any](fn func(P) *Node, props P) *Node {
	return component(fn, props, true)
}

func component[P any](fn func(P) *Node, props P, memo bool) *Node {
	return &Node{
		typ:        typeComponent,
		fnPtr:      reflect.ValueOf(fn).Pointer(),
		props:      props,
		memo:       memo,
		render:     func(p any) *Node { return fn(p.(P)) },
		propsEqual: shallowEqual,
	}
}

// If 条件渲染：cond 为真返回 n，否则返回 nil（nil 子节点会被忽略）。
func If(cond bool, n *Node) *Node {
	if cond {
		return n
	}
	return nil
}

// Keyed 给任意节点（含组件）打上 key，用于列表稳定身份。
func Keyed(k string, n *Node) *Node {
	if n != nil {
		n.key = k
	}
	return n
}
