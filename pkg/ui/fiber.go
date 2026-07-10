package ui

import "reflect"

// Fiber 是可变的身份/状态层（相当于 React Fiber / Flutter Element）。
// 跨重建复用，持有 hooks、子 Fiber、以及 host/text 的 renderNode。
type Fiber struct {
	typ nodeType
	key string

	parent   *Fiber
	children []*Fiber

	// component
	fnPtr      uintptr
	props      any
	render     func(any) *Node
	propsEqual func(a, b any) bool
	memo       bool
	hooks      []any

	// provider
	ctxID       int
	ctxValue    any
	subscribers []*Fiber

	// host / text
	tag   string
	rnode *renderNode

	// portal（独立浮层布局根）
	overlayRoot *renderNode

	dirty     bool
	queued    bool
	unmounted bool
}

func sameType(f *Fiber, n *Node) bool {
	if f.typ != n.typ {
		return false
	}
	switch n.typ {
	case typeComponent:
		return f.fnPtr == n.fnPtr && f.key == n.key
	case typeHost:
		return f.tag == n.tag && f.key == n.key
	case typeProvider:
		return f.ctxID == n.ctxID && f.key == n.key
	case typeFragment, typePortal, typeText:
		return f.key == n.key
	}
	return false
}

// reconcile 协调单个子节点：可复用则原地更新，否则卸载旧的、挂载新的。
func reconcile(parent, old *Fiber, n *Node) *Fiber {
	if n == nil {
		if old != nil {
			unmount(old)
		}
		return nil
	}
	if old != nil && sameType(old, n) {
		updateFiber(old, n)
		return old
	}
	if old != nil {
		unmount(old)
	}
	return mountFiber(parent, n)
}

func mountFiber(parent *Fiber, n *Node) *Fiber {
	f := &Fiber{typ: n.typ, key: n.key, parent: parent}
	switch n.typ {
	case typeComponent:
		f.fnPtr, f.props, f.render, f.propsEqual, f.memo = n.fnPtr, n.props, n.render, n.propsEqual, n.memo
		renderComponent(f)
	case typeProvider:
		f.ctxID, f.ctxValue = n.ctxID, n.ctxValue
		f.children = mountList(f, n.kids)
	case typeFragment:
		f.children = mountList(f, n.kids)
	case typePortal:
		f.overlayRoot = newBoxRenderNode()
		f.children = mountList(f, n.kids)
	case typeHost:
		f.tag = n.tag
		f.rnode = newHostRenderNode(n.tag)
		f.rnode.owner = f
		applyHostProps(f.rnode, buildHostProps(n))
		f.children = mountList(f, n.kids)
	case typeText:
		f.rnode = newTextRenderNode(n.text, n.textStyle)
		f.rnode.owner = f
	}
	return f
}

func mountList(parent *Fiber, kids []*Node) []*Fiber {
	var out []*Fiber
	for _, kn := range kids {
		if cf := reconcile(parent, nil, kn); cf != nil {
			out = append(out, cf)
		}
	}
	return out
}

func updateFiber(f *Fiber, n *Node) {
	switch n.typ {
	case typeComponent:
		same := f.propsEqual(f.props, n.props)
		f.props = n.props
		f.render = n.render
		f.memo = n.memo
		if f.memo && same && !f.dirty {
			return // memo 短路
		}
		renderComponent(f)
	case typeProvider:
		changed := !reflect.DeepEqual(f.ctxValue, n.ctxValue)
		f.ctxValue = n.ctxValue
		if changed {
			subs := f.subscribers
			f.subscribers = nil
			for _, s := range subs {
				if !s.unmounted && activeGame != nil {
					activeGame.markDirty(s)
				}
			}
		}
		f.children = reconcileList(f, f.children, n.kids)
	case typeFragment, typePortal:
		f.children = reconcileList(f, f.children, n.kids)
	case typeHost:
		applyHostProps(f.rnode, buildHostProps(n))
		f.children = reconcileList(f, f.children, n.kids)
	case typeText:
		f.rnode.setText(n.text, n.textStyle)
	}
}

// renderComponent 执行组件 render 并协调其唯一子树。
func renderComponent(f *Fiber) {
	f.dirty = false
	childNode := renderWithHooks(f)
	var old *Fiber
	if len(f.children) > 0 {
		old = f.children[0]
	}
	if cf := reconcile(f, old, childNode); cf != nil {
		f.children = []*Fiber{cf}
	} else {
		f.children = nil
	}
}

func unmount(f *Fiber) {
	f.unmounted = true
	if activeGame != nil && activeGame.focusedFiber == f {
		activeGame.focusedFiber = nil
	}
	for _, h := range f.hooks {
		if e, ok := h.(*effectHook); ok && e.cleanup != nil {
			e.cleanup()
			e.cleanup = nil
		}
	}
	for _, c := range f.children {
		unmount(c)
	}
	// 从 yoga 父节点摘除；不 Reset（节点即将被丢弃，Reset 对仍挂在父上的节点会 panic）。
	if f.rnode != nil && f.rnode.parent != nil {
		f.rnode.parent.yn.RemoveChild(f.rnode.yn)
		f.rnode.parent = nil
	}
}

// reconcileList 协调子节点列表（带 Key 的简化 diff）。
func reconcileList(parent *Fiber, old []*Fiber, nodes []*Node) []*Fiber {
	keyed := make(map[string]*Fiber)
	var unkeyed []*Fiber
	for _, c := range old {
		if c.key != "" {
			keyed[c.key] = c
		} else {
			unkeyed = append(unkeyed, c)
		}
	}

	out := make([]*Fiber, 0, len(nodes))
	uidx := 0
	for _, nd := range nodes {
		if nd == nil {
			continue
		}
		var match *Fiber
		if nd.key != "" {
			if m, ok := keyed[nd.key]; ok && sameType(m, nd) {
				match = m
				delete(keyed, nd.key)
			}
		} else if uidx < len(unkeyed) {
			if sameType(unkeyed[uidx], nd) {
				match = unkeyed[uidx]
			}
			uidx++
		}

		if match != nil {
			updateFiber(match, nd)
			out = append(out, match)
		} else {
			out = append(out, mountFiber(parent, nd))
		}
	}

	for _, c := range unkeyed[min(uidx, len(unkeyed)):] {
		unmount(c)
	}
	for _, c := range keyed {
		unmount(c)
	}
	return out
}

// ---- yoga 树链接 ----

// rootRenderNode 从 Fiber 向下找到第一个 renderNode（跳过组件 Fiber）。
func rootRenderNode(f *Fiber) *renderNode {
	if f.rnode != nil {
		return f.rnode
	}
	for _, c := range f.children {
		if r := rootRenderNode(c); r != nil {
			return r
		}
	}
	return nil
}

func collectChildRenderNodes(f *Fiber, out *[]*renderNode) {
	for _, c := range f.children {
		if c.typ == typePortal {
			continue // 浮层内容不进入主 yoga 树
		}
		if c.rnode != nil {
			*out = append(*out, c.rnode)
		} else {
			collectChildRenderNodes(c, out)
		}
	}
}

// collectPortals 收集整棵 Fiber 树中的所有 Portal（按树序，靠后者绘制在更上层）。
func collectPortals(f *Fiber, out *[]*Fiber) {
	for _, c := range f.children {
		if c.typ == typePortal {
			*out = append(*out, c)
		}
		collectPortals(c, out)
	}
}

// relink 依据 Fiber 树同步 yoga 父子关系。
// 仅在子节点集合真正变化时重建 yoga 链接——否则跳过，保持 yoga 布局缓存有效，
// 让纯绘制变更（颜色/hover 等）不触发整棵树重新计算。
func relink(f *Fiber) {
	if f.rnode != nil && f.rnode.container() {
		var kids []*renderNode
		collectChildRenderNodes(f, &kids)
		if !renderNodesEqual(f.rnode.children, kids) {
			f.rnode.yn.RemoveAllChildren()
			f.rnode.children = kids
			for i, k := range kids {
				k.parent = f.rnode
				f.rnode.yn.InsertChild(k.yn, uint32(i))
			}
		}
	}
	for _, c := range f.children {
		relink(c)
	}
}

func renderNodesEqual(a, b []*renderNode) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func depth(f *Fiber) int {
	d := 0
	for p := f.parent; p != nil; p = p.parent {
		d++
	}
	return d
}
