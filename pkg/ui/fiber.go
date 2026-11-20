package ui

import "reflect"

// Fiber 是可变的身份/状态层（相当于 React Fiber / Flutter Element）。
// 跨重建复用，持有 hooks、子 Fiber、以及 host/text 的 renderNode。
type Fiber struct {
	typ nodeType
	key string

	parent   *Fiber
	children []*Fiber
	level    int

	// component
	fnPtr         uintptr
	props         any
	render        func(any) *Node
	propsEqual    func(a, b any) bool
	memo          bool
	hooks         []any
	instanceType  reflect.Type
	renderContext *RenderContext

	// provider
	ctxID        int
	ctxValue     any
	subscribers  map[*Fiber]struct{} // 消费者集合（去重）；值变化时统一标脏
	providerSubs []*Fiber            // 本 fiber 订阅的 provider，卸载时据此退订

	// host / text
	tag       string
	rnode     *renderNode
	leafAttrs bool // Text/Icon 曾应用链式属性，移除修饰时也须清理旧属性

	// portal（独立浮层布局根）
	overlayRoot *renderNode
	portalTrap  bool // TrapFocus()：该浮层限制键盘焦点在内部（模态）

	dirty     bool
	queued    bool
	unmounted bool

	// cleanups 是卸载时要执行的清理（Store/Signal 订阅退订等），由 addCleanup 登记。
	cleanups []func()

	// 错误边界（ErrorBoundary）：捕获后代 render 期的 panic
	errBoundary bool
	caughtErr   any
}

// addCleanup 登记卸载清理（订阅退订等）。多次登记的函数按登记顺序执行。
func (f *Fiber) addCleanup(fn func()) { f.cleanups = append(f.cleanups, fn) }

// boundaryCount 是活动 ErrorBoundary 数量；为 0 时 renderComponent 走无 defer 的快路径。
var boundaryCount int

// nearestBoundary 返回 f 最近的、当前未处于错误态的祖先错误边界。
func nearestBoundary(f *Fiber) *Fiber {
	for a := f.parent; a != nil; a = a.parent {
		if a.errBoundary && a.caughtErr == nil {
			return a
		}
	}
	return nil
}

func sameType(f *Fiber, n *Node) bool {
	if f.typ != n.typ {
		return false
	}
	switch n.typ {
	case typeComponent:
		return f.fnPtr == n.fnPtr && f.instanceType == n.instanceType && f.key == n.key
	case typeHost:
		return f.tag == n.tag && f.key == n.key
	case typeProvider:
		return f.ctxID == n.ctxID && f.key == n.key
	case typeFragment, typePortal, typeText, typeIcon:
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
	if parent != nil {
		f.level = parent.level + 1
	}
	switch n.typ {
	case typeComponent:
		f.fnPtr, f.props, f.render, f.propsEqual, f.memo = n.fnPtr, n.props, n.render, n.propsEqual, n.memo
		f.instanceType = n.instanceType
		renderComponent(f)
	case typeProvider:
		f.ctxID, f.ctxValue = n.ctxID, n.ctxValue
		f.children = mountList(f, n.kids)
	case typeFragment:
		f.children = mountList(f, n.kids)
	case typePortal:
		f.overlayRoot = newBoxRenderNode()
		f.portalTrap = n.trap
		f.children = mountList(f, n.kids)
	case typeHost:
		f.tag = n.tag
		f.rnode = newHostRenderNode(n.tag)
		f.rnode.owner = f
		hp := buildHostProps(n)
		applyHostProps(f, f.rnode, hp)
		releaseHostProps(hp)
		f.children = mountList(f, n.kids)
	case typeText:
		hp := buildHostProps(n)
		f.rnode = newTextRenderNode(n.text, hp.style, cloneRuns(n.runs))
		f.rnode.owner = f
		applyLeafAttrs(f, n, hp)
		attachTextBinds(f, n.binds, f.rnode)
		releaseHostProps(hp)
	case typeIcon:
		hp := buildHostProps(n)
		f.rnode = newIconRenderNode(n.iconPath, n.iconSize, n.iconStroke, n.iconRaw, n.iconW, n.iconH, hp.style)
		f.rnode.owner = f
		applyLeafAttrs(f, n, hp)
		releaseHostProps(hp)
	}
	return f
}

// 普通文字/图标沿用轻量路径，只有附加或移除链式属性时才同步完整的布局和事件属性。
func applyLeafAttrs(f *Fiber, n *Node, hp *hostProps) {
	hasAttrs := len(n.attrList) > 0
	if hasAttrs || f.leafAttrs {
		applyHostProps(f, f.rnode, hp)
	}
	f.leafAttrs = hasAttrs
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
			for s := range subs {
				if !s.unmounted && activeGame != nil {
					activeGame.markDirty(s)
				}
			}
		}
		f.children = reconcileList(f, f.children, n.kids)
	case typeFragment, typePortal:
		f.portalTrap = n.trap
		f.children = reconcileList(f, f.children, n.kids)
	case typeHost:
		hp := buildHostProps(n)
		applyHostProps(f, f.rnode, hp)
		releaseHostProps(hp)
		f.children = reconcileList(f, f.children, n.kids)
	case typeText:
		hp := buildHostProps(n)
		attachTextBinds(f, n.binds, f.rnode)
		if len(n.binds) > 0 {
			// 绑定节点：文本由 Signal 直写，样式仍随 render 更新；占位 text 不覆盖活值。
			f.rnode.applyTextStyle(hp.style)
		} else {
			f.rnode.setText(n.text, hp.style, cloneRuns(n.runs))
		}
		applyLeafAttrs(f, n, hp)
		releaseHostProps(hp)
	case typeIcon:
		hp := buildHostProps(n)
		f.rnode.setIcon(n.iconPath, n.iconSize, n.iconStroke, n.iconRaw, n.iconW, n.iconH, hp.style)
		applyLeafAttrs(f, n, hp)
		releaseHostProps(hp)
	}
}

// renderComponent 执行组件 render 并协调其唯一子树。
// 存在错误边界时用 recover 捕获本组件（或其后代）render 期的 panic，交给最近的边界降级显示。
func renderComponent(f *Fiber) {
	f.dirty = false
	if boundaryCount > 0 {
		defer func() {
			if r := recover(); r != nil {
				b := nearestBoundary(f)
				if b == nil {
					panic(r) // 无边界可兜 -> 继续抛出
				}
				b.caughtErr = r
				if activeGame != nil {
					activeGame.markDirty(b)
				}
			}
		}()
	}
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
	if f.renderContext != nil {
		f.renderContext.cancel()
	}
	if f.errBoundary {
		boundaryCount--
	}
	// 退订本 fiber 消费过的所有 provider，避免 provider 悬挂已卸载消费者。
	for _, p := range f.providerSubs {
		delete(p.subscribers, f)
	}
	f.providerSubs = nil
	if g := activeGame; g != nil {
		delete(g.viewportReaders, f)
		if g.focusedFiber == f {
			g.focusedFiber = nil
		}
		// 清理指向本节点的悬空交互指针，防止在已卸载节点上回调 onPress/onDrag/onHover。
		if rn := f.rnode; rn != nil {
			if g.pressedNode == rn {
				g.pressedNode = nil
			}
			if g.dragging == rn {
				g.dragging = nil
			}
			if rn.parent != nil {
				rn.parent.hover.children[0].Delete(rn.index)
				rn.parent.hover.children[1].Delete(rn.index)
			}
			rn.hover.self = [2]bool{}
		}
	}
	for _, h := range f.hooks {
		if e, ok := h.(*effectHook); ok && e.cleanup != nil {
			e.cleanup()
			e.cleanup = nil
		}
	}
	// 通用清理（Store/Signal 订阅退订等）。
	for _, c := range f.cleanups {
		c()
	}
	f.cleanups = nil
	for _, c := range f.children {
		unmount(c)
	}
	// 从 yoga 父节点摘除；不 Reset（节点即将被丢弃，Reset 对仍挂在父上的节点会 panic）。
	if f.rnode != nil {
		if f.rnode.parent != nil {
			f.rnode.parent.yn.RemoveChild(f.rnode.yn)
			f.rnode.parent = nil
		}
		f.rnode.yn.Finalize()
	}
	if f.overlayRoot != nil {
		f.overlayRoot.yn.Finalize()
	}
}

// reconcileList 协调子节点列表（带 Key 的简化 diff）。
func reconcileList(parent *Fiber, old []*Fiber, nodes []*Node) []*Fiber {
	unkeyedOnly := true
	for _, child := range old {
		if child.key != "" {
			unkeyedOnly = false
			break
		}
	}
	if unkeyedOnly {
		for _, node := range nodes {
			if node != nil && node.key != "" {
				unkeyedOnly = false
				break
			}
		}
	}
	if unkeyedOnly {
		out := old[:0]
		for _, node := range nodes {
			if node == nil {
				continue
			}
			var previous *Fiber
			if len(out) < len(old) {
				previous = old[len(out)]
			}
			out = append(out, reconcile(parent, previous, node))
		}
		if len(out) < len(old) {
			for _, child := range old[len(out):] {
				unmount(child)
			}
		}
		return out
	}

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
			} else {
				// 该位置已消费，但旧节点并未复用，必须立刻卸载。
				// 否则会遗漏 effects、订阅与组件 Context 的生命周期清理。
				unmount(unkeyed[uidx])
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
	if f == nil {
		return nil
	}
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

func attachRenderChild(parent, child *renderNode, index uint32) {
	child.parent = parent
	child.index = index
	child.level = parent.level + 1
	if activeGame != nil && child.level > activeGame.maxLevel {
		activeGame.maxLevel = child.level
	}
	parent.yn.InsertChild(child.yn, index)
}

// relink 依据 Fiber 树同步 yoga 父子关系。
// 仅在子节点集合真正变化时重建 yoga 链接——否则跳过，保持 yoga 布局缓存有效，
// 让纯绘制变更（颜色/hover 等）不触发整棵树重新计算。
func relink(f *Fiber) {
	if f.rnode != nil && f.rnode.container() {
		if !renderChildrenEqual(f, f.rnode.children) {
			var kids []*renderNode
			collectChildRenderNodes(f, &kids)
			f.rnode.yn.RemoveAllChildren()
			f.rnode.children = kids
			for i, k := range kids {
				attachRenderChild(f.rnode, k, uint32(i))
			}
		}
	}
	for _, c := range f.children {
		relink(c)
	}
}

func renderChildrenEqual(f *Fiber, existing []*renderNode) bool {
	index := 0
	equal := true
	var visit func(*Fiber)
	visit = func(parent *Fiber) {
		for _, child := range parent.children {
			if child.typ == typePortal {
				continue
			}
			if child.rnode != nil {
				if index >= len(existing) || existing[index] != child.rnode {
					equal = false
				}
				index++
			} else {
				visit(child)
			}
		}
	}
	visit(f)
	return equal && index == len(existing)
}
