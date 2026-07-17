package ui

import (
	"math"
	"sort"
	"time"

	"github.com/sjm1327605995/tenon/yoga"
)

// activeGame 是当前运行的驱动实例，供 hooks 触发局部重渲染。
var activeGame *game

// uiScale 是设备像素比：用户以逻辑像素编写 UI，引擎在物理分辨率下渲染以获得清晰边缘（抗锯齿）。
// 在 syncYoga（尺寸）、文本字号、拖拽/滚轮增量处换算。
var uiScale float32 = 1

type game struct {
	root      *Node
	rootFiber *Fiber
	rootRN    *renderNode
	w, h      int

	dirty          []*Fiber
	pendingEffects []func()
	needsLayout    bool
	focusedFiber   *Fiber
	anims          []*tweenHook
	loops          []*loopHook
	lastFrame      time.Time
	hovered        map[*renderNode]bool

	dragging             *renderNode
	dragLastX, dragLastY float32
	pressedNode          *renderNode
	inputSelecting       bool

	hoverX, hoverY float32 // 上次计算悬停链时的光标位置（用于空闲时跳过重算）

	// 多击检测（双击选词 / 三击选全部）
	lastClickAt            time.Time
	lastClickX, lastClickY float32
	clickCount             int

	portals  []*Fiber
	escStack []*escEntry

	laidOutW, laidOutH int
	boundsDirty        bool
	hasLayoutAnim      bool

	// imeComposing 为真表示正在输入法组字（预编辑）。gio 的 IME 组字后续再接。
	imeComposing bool
}

// FrameSync 预留：控制是否跟随刷新率（gio 循环当前恒重绘，暂未使用）。
var FrameSync = true

// windowConfig 是窗口的启动配置。后端中立：gio 侧（gioRun）把它映射成 app.Option。
// 尺寸都是逻辑像素（dp），由后端按屏幕缩放换算。
type windowConfig struct {
	w, h       int
	title      string
	minW, minH int // 0 = 不限制
	maxW, maxH int // 0 = 不限制
	fullscreen bool
	sync       bool
}

var winCfg = windowConfig{w: 800, h: 600, title: "Tenon UI"}

// WindowSize 设置初始窗口的逻辑尺寸（须在 Run 之前调用）。
func WindowSize(w, h int) {
	if w > 0 && h > 0 {
		winCfg.w, winCfg.h = w, h
	}
}

// WindowTitle 设置窗口标题（须在 Run 之前调用）。
func WindowTitle(t string) { winCfg.title = t }

// WindowMinSize 设置窗口的最小逻辑尺寸，防止被拖到布局撑不开的宽度（0=不限制）。
func WindowMinSize(w, h int) { winCfg.minW, winCfg.minH = w, h }

// WindowMaxSize 设置窗口的最大逻辑尺寸（0=不限制）。
func WindowMaxSize(w, h int) { winCfg.maxW, winCfg.maxH = w, h }

// WindowFullscreen 让窗口以全屏启动。
func WindowFullscreen(on bool) { winCfg.fullscreen = on }

// Run 启动应用；root 通常是一个 Use(...) 组件节点。窗口与事件循环由 gio 后端（backendRun）提供。
// 窗口相关设置（WindowSize/WindowTitle/...）须在此之前调用。
func Run(root *Node) {
	cfg := winCfg
	cfg.sync = FrameSync
	backendRun(root, cfg)
}

// tickLayoutAnim 逐帧推进布局动画：检测位置变化注入残余偏移，并指数衰减到 0。
// tickLayoutAnim 推进 FLIP 布局动画；返回是否仍有节点在移动（供按需重绘判断）。
func (g *game) tickLayoutAnim(dt float32) bool {
	if dt <= 0 || !g.hasLayoutAnim {
		return false
	}
	decay := float32(math.Exp(-float64(dt) / 70))
	active := false
	if g.rootRN != nil {
		active = walkLayoutAnim(g.rootRN, decay) || active
	}
	for _, pf := range g.portals {
		if pf.overlayRoot != nil {
			active = walkLayoutAnim(pf.overlayRoot, decay) || active
		}
	}
	return active
}

// walkLayoutAnim 返回子树内是否还有非零残余偏移（仍在动画中）。
func walkLayoutAnim(rn *renderNode, decay float32) bool {
	active := false
	if rn.animatedLayout {
		if rn.hasPrevPos {
			rn.offX += rn.prevPosX - rn.bounds.X
			rn.offY += rn.prevPosY - rn.bounds.Y
		}
		rn.offX *= decay
		rn.offY *= decay
		if rn.offX < 0.5 && rn.offX > -0.5 {
			rn.offX = 0
		}
		if rn.offY < 0.5 && rn.offY > -0.5 {
			rn.offY = 0
		}
		rn.prevPosX, rn.prevPosY = rn.bounds.X, rn.bounds.Y
		rn.hasPrevPos = true
		if rn.offX != 0 || rn.offY != 0 {
			active = true
		}
	}
	for _, c := range rn.children {
		active = walkLayoutAnim(c, decay) || active
	}
	return active
}

func (g *game) markDirty(f *Fiber) {
	if f.queued || f.unmounted {
		return
	}
	f.queued = true
	f.dirty = true
	g.dirty = append(g.dirty, f)
}

func (g *game) flushDirty() {
	q := g.dirty
	g.dirty = nil
	// 浅层优先：祖先先重渲染，其协调会顺带处理已脏的后代，避免重复渲染。
	sort.SliceStable(q, func(i, j int) bool { return depth(q[i]) < depth(q[j]) })
	for _, f := range q {
		f.queued = false
		if !f.dirty || f.unmounted {
			continue // 跳过本帧内已被卸载的 fiber，避免在死组件上重跑 render/effect
		}
		renderComponent(f)
	}
	g.needsLayout = true
}

// tickAnims 推进所有活动补间动画，并把其所属组件标记为需重渲染。
func (g *game) tickAnims(dt float32) {
	if len(g.anims) == 0 || dt <= 0 {
		return
	}
	live := g.anims[:0]
	for _, h := range g.anims {
		if h.fiber.unmounted || !h.active {
			continue
		}
		if h.advance(dt) {
			live = append(live, h)
		}
		g.markDirty(h.fiber)
	}
	g.anims = live
}

// tickLoops 推进持续动画（UseElapsed），每帧累加时间并标记重渲染。
func (g *game) tickLoops(dt float32) {
	if len(g.loops) == 0 || dt <= 0 {
		return
	}
	live := g.loops[:0]
	for _, h := range g.loops {
		if h.fiber.unmounted {
			h.active = false
			continue
		}
		h.elapsed += dt / 1000
		g.markDirty(h.fiber)
		live = append(live, h)
	}
	g.loops = live
}

func (g *game) flushEffects() {
	if len(g.pendingEffects) == 0 {
		return
	}
	effs := g.pendingEffects
	g.pendingEffects = nil
	for _, e := range effs {
		e()
	}
}

func (g *game) layout() {
	if g.rootRN == nil {
		return
	}
	relink(g.rootFiber) // 增量：仅结构变化时改动 yoga 链接
	resolveInherited(g.rootRN, inhText{})

	windowChanged := g.w != g.laidOutW || g.h != g.laidOutH
	if g.rootRN.yn.IsDirty() || windowChanged {
		g.rootRN.yn.CalculateLayout(float32(g.w), float32(g.h), yoga.DirectionLTR)
		g.laidOutW, g.laidOutH = g.w, g.h
		g.boundsDirty = true
	}
	if g.boundsDirty {
		computeBounds(g.rootRN, 0, 0)
		syncMeasures(g.rootRN)
		g.boundsDirty = false
	}
	g.layoutPortals(windowChanged)
}

// layoutPortals 为每个 Portal 建立全屏独立布局根并计算其 bounds（同样按需重算）。
func (g *game) layoutPortals(windowChanged bool) {
	g.portals = g.portals[:0]
	collectPortals(g.rootFiber, &g.portals)
	for _, pf := range g.portals {
		root := pf.overlayRoot
		var kids []*renderNode
		collectChildRenderNodes(pf, &kids)
		if !renderNodesEqual(root.children, kids) {
			root.yn.RemoveAllChildren()
			root.children = kids
			for i, k := range kids {
				k.parent = root
				root.yn.InsertChild(k.yn, uint32(i))
			}
		}
		root.yn.StyleSetWidth(float32(g.w))
		root.yn.StyleSetHeight(float32(g.h))
		root.yn.StyleSetFlexDirection(yoga.FlexDirectionColumn)
		resolveInherited(root, inhText{})
		if root.yn.IsDirty() || windowChanged {
			root.yn.CalculateLayout(float32(g.w), float32(g.h), yoga.DirectionLTR)
			computeBounds(root, 0, 0)
			syncMeasures(root)
		}
	}
}

// hitTop 自顶向下命中：先浮层（逆序），再主树。
//
// 浮层根是覆盖全窗的容器，它本身不是可命中的目标 —— 只有它的内容才是。hitNode 在没命中
// 任何子节点时会返回容器自己，这里必须把「只命中到浮层根」当作没命中，让点穿透到主树。
// 否则任一浮层存在就会吞掉整个界面的命中，并形成自激：悬停触发器 -> 弹出 Tooltip 浮层 ->
// 浮层根吞掉命中 -> 触发器 unhover -> 浮层收起 -> 又命中到触发器 -> 再弹出，每帧一次，
// 看起来就是 hover 框在闪。
func (g *game) hitTop(x, y float32) *renderNode {
	for i := len(g.portals) - 1; i >= 0; i-- {
		if r := g.portals[i].overlayRoot; r != nil {
			if h := hitNode(r, x, y); h != nil && h != r {
				return h
			}
		}
	}
	if g.rootRN != nil {
		return hitNode(g.rootRN, x, y)
	}
	return nil
}

func (g *game) handleInput() {
	if g.rootRN == nil {
		return
	}
	g.updateHover()

	// 滚轮：把滚动施加到光标下最近的可滚动祖先。
	// wheel() 与 bounds 同为物理像素，直接相加即可 —— 不要再乘 uiScale，那会在高 DPI 上
	// 把滚动速度按缩放倍数放大（150% 缩放下快 1.5 倍）。
	if _, wy := input.wheel(); wy != 0 {
		x, y := input.cursor()
		for c := g.hitTop(x, y); c != nil; c = c.parent {
			if c.scroll {
				c.scrollY += wy
				g.needsLayout = true
				g.boundsDirty = true // 滚动改变绝对 bounds 但不脏化 yoga
				break
			}
		}
	}

	if input.mouseJustPressed(btnLeft) {
		x, y := input.cursor()
		n := g.hitTop(x, y)
		// 聚焦：命中 input 则聚焦；单击定位光标并开始拖选，双击选词，三击选全部
		if n != nil && n.kind == rnInput {
			g.focusedFiber = n.owner
			c := n.caretAt(x, y)
			now := time.Now()
			if absf(x-g.lastClickX) < 4 && absf(y-g.lastClickY) < 4 && now.Sub(g.lastClickAt) < 400*time.Millisecond {
				g.clickCount++
			} else {
				g.clickCount = 1
			}
			if g.clickCount > 3 {
				g.clickCount = 1
			}
			g.lastClickAt, g.lastClickX, g.lastClickY = now, x, y
			switch g.clickCount {
			case 2: // 双击选词
				n.selAnchor, n.caretPos = wordAt(n.value, c)
				g.inputSelecting = false
			case 3: // 三击选全部
				n.selAnchor, n.caretPos = 0, len(n.value)
				g.inputSelecting = false
			default:
				n.caretPos, n.selAnchor = c, c
				g.inputSelecting = true
			}
		} else {
			g.focusedFiber = nil
			g.clickCount = 0
		}
		// 按压开始：向上找第一个 onPress
		for c := n; c != nil; c = c.parent {
			if c.onPress != nil {
				g.pressedNode = c
				c.onPress(true)
				break
			}
		}
		// 拖拽开始：向上找第一个 onDrag
		for c := n; c != nil; c = c.parent {
			if c.onDrag != nil {
				g.dragging = c
				g.dragLastX, g.dragLastY = float32(x), float32(y)
				break
			}
		}
		// 点击冒泡：从命中节点向上找第一个 onClick
		for c := n; c != nil; c = c.parent {
			if c.onClick != nil {
				c.onClick()
				break
			}
		}
	}
	// 右键：向上冒泡找第一个 onContextMenu，回调光标逻辑坐标
	if input.mouseJustPressed(btnRight) {
		x, y := input.cursor()
		for c := g.hitTop(x, y); c != nil; c = c.parent {
			if c.onContextMenu != nil {
				c.onContextMenu(float32(x)/uiScale, float32(y)/uiScale)
				break
			}
		}
	}
	g.updatePress()
	g.updateInputSelection()
	g.handleKeyboardNav()
	g.updateDrag()
	g.editFocusedInput()
}

// updateInputSelection 在聚焦输入框上拖动鼠标时扩展选区（单行）。
func (g *game) updateInputSelection() {
	if !g.inputSelecting {
		return
	}
	if !input.mousePressed(btnLeft) {
		g.inputSelecting = false
		return
	}
	f := g.focusedFiber
	if f == nil || f.unmounted || f.rnode == nil || f.rnode.kind != rnInput {
		return
	}
	x, y := input.cursor()
	f.rnode.caretPos = f.rnode.caretAt(float32(x), float32(y)) // anchor 不动 -> 形成选区
}

// updatePress 在左键松开时结束按压态并回调 onPress(false)。
func (g *game) updatePress() {
	if g.pressedNode == nil {
		return
	}
	if !input.mousePressed(btnLeft) {
		if g.pressedNode.onPress != nil {
			g.pressedNode.onPress(false)
		}
		g.pressedNode = nil
	}
}

// handleKeyboardNav 处理 Tab 焦点切换、Enter/Space 激活、Esc 失焦。
// 决策逻辑抽到 focusNext/fireEscape/activateFocused，供无窗口的测试驱动复用。
func (g *game) handleKeyboardNav() {
	if input.keyJustPressed(keyTab) {
		g.focusNext(!input.keyPressed(keyShift))
	}
	// 方向键：在导航组内移动焦点（输入框聚焦时 moveFocusInGroup 自动放行给光标）
	switch {
	case input.keyJustPressed(keyDown):
		g.moveFocusInGroup(true, NavVertical)
	case input.keyJustPressed(keyUp):
		g.moveFocusInGroup(false, NavVertical)
	case input.keyJustPressed(keyRight):
		g.moveFocusInGroup(true, NavHorizontal)
	case input.keyJustPressed(keyLeft):
		g.moveFocusInGroup(false, NavHorizontal)
	}
	if input.keyJustPressed(keyEscape) {
		g.fireEscape()
	}
	// Enter/Space 激活聚焦的可点击元素（输入框不拦截，交给文本编辑）
	if input.keyJustPressed(keyEnter) || input.keyJustPressed(keySpace) {
		g.activateFocused()
	}
}

// focusNext 把焦点移到 Tab 顺序中的下一个（forward）或上一个可聚焦元素（含浮层），
// 环形回绕；返回新的焦点节点，无可聚焦元素时返回 nil。
func (g *game) focusNext(forward bool) *renderNode {
	var list []*renderNode
	if scope := g.trapScope(); scope != nil {
		collectFocusables(scope, &list) // 模态：焦点只在浮层内循环
	} else {
		collectFocusables(g.rootRN, &list)
		for _, pf := range g.portals {
			if pf.overlayRoot != nil {
				collectFocusables(pf.overlayRoot, &list)
			}
		}
	}
	var cur *renderNode
	if g.focusedFiber != nil {
		cur = g.focusedFiber.rnode
	}
	n := nextFocus(list, cur, forward)
	if n != nil {
		g.focusedFiber = n.owner
		if n.kind == rnInput {
			n.caretPos = len(n.value)
		}
	}
	return n
}

// fireEscape 触发 Esc：最上层浮层优先响应，无浮层时取消聚焦。
func (g *game) fireEscape() {
	if n := len(g.escStack); n > 0 {
		(*g.escStack[n-1].fn)()
	} else {
		g.focusedFiber = nil
	}
}

// activateFocused 用 Enter/Space 激活当前聚焦的可点击元素（输入框除外，交给文本编辑）。
// 返回是否触发了 onClick。
func (g *game) activateFocused() bool {
	if g.focusedFiber == nil || g.focusedFiber.unmounted {
		return false
	}
	rn := g.focusedFiber.rnode
	if rn != nil && rn.kind != rnInput && rn.onClick != nil {
		rn.onClick()
		return true
	}
	return false
}

// updateDrag 按住左键时逐帧把屏幕位移交给拖拽目标，松开则结束。
func (g *game) updateDrag() {
	if g.dragging == nil {
		return
	}
	if !input.mousePressed(btnLeft) {
		g.dragging = nil
		return
	}
	x, y := input.cursor()
	dx, dy := float32(x)-g.dragLastX, float32(y)-g.dragLastY
	if dx != 0 || dy != 0 {
		g.dragLastX, g.dragLastY = float32(x), float32(y)
		if g.dragging.onDrag != nil {
			// 光标是物理像素，回调给用户逻辑像素
			g.dragging.onDrag(dx/uiScale, dy/uiScale)
		}
	}
}

// editFocusedInput 处理聚焦输入框的键盘编辑：文本输入、选区、剪切/复制/粘贴/全选（受控回流）。
// 注：CJK 输入法组字（IME 预编辑）尚未在 gio 后端接入，此处只做手动编辑。
func (g *game) editFocusedInput() {
	f := g.focusedFiber
	if f == nil || f.unmounted || f.rnode == nil || f.rnode.kind != rnInput {
		g.imeComposing = false
		return
	}
	rn := f.rnode

	val := rn.value
	caret := clampi(rn.caretPos, 0, len(val))
	anchor := clampi(rn.selAnchor, 0, len(val))

	shift := input.keyPressed(keyShift)
	ctrl := input.keyPressed(keyCtrl) || input.keyPressed(keyMeta)

	selLo := func() int { return min(anchor, caret) }
	selHi := func() int { return max(anchor, caret) }
	hasSel := func() bool { return anchor != caret }
	delSel := func() {
		lo, hi := selLo(), selHi()
		val = val[:lo] + val[hi:]
		caret, anchor = lo, lo
	}
	// 移动光标后按是否 shift 决定保留/折叠选区
	afterMove := func() {
		if !shift {
			anchor = caret
		}
	}

	// 快捷键
	if ctrl {
		switch {
		case input.keyJustPressed(keyA):
			anchor, caret = 0, len(val)
		case input.keyJustPressed(keyC):
			if hasSel() {
				setClipboard(val[selLo():selHi()])
			}
		case input.keyJustPressed(keyX):
			if hasSel() {
				setClipboard(val[selLo():selHi()])
				delSel()
			}
		case input.keyJustPressed(keyV):
			if hasSel() {
				delSel()
			}
			p := getClipboard()
			val = val[:caret] + p + val[caret:]
			caret += len(p)
			anchor = caret
		}
	} else {
		// 文本输入（有选区时先替换）
		for _, r := range input.typedChars() {
			if r >= 0x20 && r != 0x7f {
				if hasSel() {
					delSel()
				}
				s := string(r)
				val = val[:caret] + s + val[caret:]
				caret += len(s)
				anchor = caret
			}
		}
	}

	if repeatKey(keyBackspace) {
		if hasSel() {
			delSel()
		} else if caret > 0 {
			prev := prevGraphemeBoundary(val, caret) // Ctrl 删整词，否则删一个字素簇
			if ctrl {
				prev = prevWordBoundary(val, caret)
			}
			val = val[:prev] + val[caret:]
			caret, anchor = prev, prev
		}
	}
	if repeatKey(keyDelete) {
		if hasSel() {
			delSel()
		} else if caret < len(val) {
			next := nextGraphemeBoundary(val, caret)
			if ctrl {
				next = nextWordBoundary(val, caret)
			}
			val = val[:caret] + val[next:]
			anchor = caret
		}
	}
	if repeatKey(keyLeft) {
		if hasSel() && !shift {
			caret = selLo()
		} else if caret > 0 {
			if ctrl {
				caret = prevWordBoundary(val, caret)
			} else {
				caret = prevGraphemeBoundary(val, caret)
			}
		}
		afterMove()
	}
	if repeatKey(keyRight) {
		if hasSel() && !shift {
			caret = selHi()
		} else if caret < len(val) {
			if ctrl {
				caret = nextWordBoundary(val, caret)
			} else {
				caret = nextGraphemeBoundary(val, caret)
			}
		}
		afterMove()
	}
	if input.keyJustPressed(keyHome) {
		caret = 0
		afterMove()
	}
	if input.keyJustPressed(keyEnd) {
		caret = len(val)
		afterMove()
	}
	if rn.multiline && input.keyJustPressed(keyEnter) {
		if hasSel() {
			delSel()
		}
		val = val[:caret] + "\n" + val[caret:]
		caret++
		anchor = caret
	}

	rn.caretPos = clampi(caret, 0, len(val))
	rn.selAnchor = clampi(anchor, 0, len(val))
	if val != rn.value && rn.onChange != nil {
		rn.onChange(val)
	}
}

// updateHover 计算光标下的悬停链，触发 enter/leave 回调（回调内一般 setState 驱动重渲染）。
func (g *game) updateHover() {
	x, y := input.cursor()
	// 空闲时（光标未移动、无待处理重渲染、无动画在跑）悬停链不会变，
	// 跳过整树命中遍历与每帧 map 分配，避免静态界面空转。
	animating := len(g.anims) > 0 || len(g.loops) > 0 || g.hasLayoutAnim
	if !animating && g.hovered != nil && x == g.hoverX && y == g.hoverY &&
		!g.needsLayout && len(g.dirty) == 0 {
		return
	}
	g.hoverX, g.hoverY = x, y
	now := map[*renderNode]bool{}
	for c := g.hitTop(x, y); c != nil; c = c.parent {
		if c.onHover != nil {
			now[c] = true
		}
	}
	for rn := range now {
		if !g.hovered[rn] {
			debugHover("enter", rn, x, y)
			rn.onHover(true)
		}
	}
	for rn := range g.hovered {
		if !now[rn] {
			debugHover("leave", rn, x, y)
			rn.onHover(false)
		}
	}
	g.hovered = now
}

// 长按重复计时（repeatKey / keyNextRepeat / repeatDelay）已移至 input.go，与后端无关。
