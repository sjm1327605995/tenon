package ui

import (
	"math"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/sjm1327605995/tenon/pkg/font"
	"github.com/sjm1327605995/tenon/yoga"
)

// activeGame 是当前运行的驱动实例，供 hooks 触发局部重渲染。
var activeGame *game

// uiScale 是设备像素比：用户以逻辑像素编写 UI，引擎在物理分辨率下渲染以获得清晰边缘（抗锯齿）。
// 在 syncYoga（尺寸）、文本字号、拖拽/滚轮增量处换算。
var uiScale float32 = 1

// initFont 优先加载内置 CJK 字体，回退到 goregular。
func initFont() {
	if len(cjkFont) > 0 {
		if err := font.GetFontManager().ReloadFontFromBytes(font.FontFamilyDefault, cjkFont); err == nil {
			return
		}
	}
	_ = font.InitDefaultFont()
}

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

	portals  []*Fiber
	escStack []*escEntry

	laidOutW, laidOutH int
	boundsDirty        bool
	hasLayoutAnim      bool
}

// Run 启动应用；root 通常是一个 Use(...) 组件节点。
func Run(root *Node) {
	initFont()
	g := &game{root: root, w: 800, h: 600}
	activeGame = g
	ebiten.SetWindowSize(g.w, g.h)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Tenon UI")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

func (g *game) Update() error {
	now := time.Now()
	var dt float32
	if !g.lastFrame.IsZero() {
		dt = float32(now.Sub(g.lastFrame).Milliseconds())
	}
	g.lastFrame = now

	if g.rootFiber == nil {
		g.rootFiber = reconcile(nil, nil, g.root)
		g.needsLayout = true
	} else {
		drainPosts() // 先执行跨 goroutine 排队的更新（在渲染线程上，setState 安全）
		g.handleInput()
		g.tickAnims(dt)
		g.tickLoops(dt)
		// 循环排空脏队列：让 Context / Memo 边界的失效在同一帧内传播完成。
		for guard := 0; len(g.dirty) > 0 && guard < 100; guard++ {
			g.flushDirty()
		}
	}

	if g.needsLayout {
		g.rootRN = rootRenderNode(g.rootFiber)
		g.layout()
		g.needsLayout = false
	}

	g.tickLayoutAnim(dt)
	g.flushEffects()
	return nil
}

// tickLayoutAnim 逐帧推进布局动画：检测位置变化注入残余偏移，并指数衰减到 0。
func (g *game) tickLayoutAnim(dt float32) {
	if dt <= 0 || !g.hasLayoutAnim {
		return
	}
	decay := float32(math.Exp(-float64(dt) / 70))
	if g.rootRN != nil {
		walkLayoutAnim(g.rootRN, decay)
	}
	for _, pf := range g.portals {
		if pf.overlayRoot != nil {
			walkLayoutAnim(pf.overlayRoot, decay)
		}
	}
}

func walkLayoutAnim(rn *renderNode, decay float32) {
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
	}
	for _, c := range rn.children {
		walkLayoutAnim(c, decay)
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(Color{247, 248, 250, 255})
	if g.rootRN != nil {
		paint(screen, g.rootRN)
	}
	// 浮层绘制在主树之上（按树序，靠后者更上层）
	for _, pf := range g.portals {
		if pf.overlayRoot != nil {
			paint(screen, pf.overlayRoot)
		}
	}
}

// SuperSample 是相对物理像素的超采样倍率：在更高分辨率渲染再由 Ebiten 缩放显示，
// 得到抗锯齿的平滑边缘（对圆角/圆形/边框尤其明显）。设为 1 可关闭以换取性能。
var SuperSample float32 = 2

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	ds := float32(ebiten.Monitor().DeviceScaleFactor())
	if ds < 1 {
		ds = 1
	}
	ss := SuperSample
	if ss < 1 {
		ss = 1
	}
	total := ds * ss
	if total > 2.5 { // 限制上限，避免超高分屏下的开销
		total = 2.5
	}
	pw, ph := int(float32(outsideWidth)*total), int(float32(outsideHeight)*total)
	if pw != g.w || ph != g.h || total != uiScale {
		uiScale = total
		g.w, g.h = pw, ph
		g.needsLayout = true
	}
	return pw, ph
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
		if !f.dirty {
			continue
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
	resolveInherited(g.rootRN, Black, false, 16, false)

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
		resolveInherited(root, Black, false, 16, false)
		if root.yn.IsDirty() || windowChanged {
			root.yn.CalculateLayout(float32(g.w), float32(g.h), yoga.DirectionLTR)
			computeBounds(root, 0, 0)
			syncMeasures(root)
		}
	}
}

// hitTop 自顶向下命中：先浮层（逆序），再主树。
func (g *game) hitTop(x, y float32) *renderNode {
	for i := len(g.portals) - 1; i >= 0; i-- {
		if r := g.portals[i].overlayRoot; r != nil {
			if h := hitNode(r, x, y); h != nil {
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
	if _, wy := ebiten.Wheel(); wy != 0 {
		x, y := ebiten.CursorPosition()
		for c := g.hitTop(float32(x), float32(y)); c != nil; c = c.parent {
			if c.scroll {
				c.scrollY -= float32(wy) * 24 * uiScale
				g.needsLayout = true
				g.boundsDirty = true // 滚动改变绝对 bounds 但不脏化 yoga
				break
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		n := g.hitTop(float32(x), float32(y))
		// 聚焦：命中 input 则聚焦并把光标定位到点击处、开始选区拖拽
		if n != nil && n.kind == rnInput {
			g.focusedFiber = n.owner
			c := n.caretFromX(float32(x))
			n.caretPos, n.selAnchor = c, c
			g.inputSelecting = true
		} else {
			g.focusedFiber = nil
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
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.inputSelecting = false
		return
	}
	f := g.focusedFiber
	if f == nil || f.unmounted || f.rnode == nil || f.rnode.kind != rnInput || f.rnode.multiline {
		return
	}
	x, _ := ebiten.CursorPosition()
	f.rnode.caretPos = f.rnode.caretFromX(float32(x)) // anchor 不动 -> 形成选区
}

// updatePress 在左键松开时结束按压态并回调 onPress(false)。
func (g *game) updatePress() {
	if g.pressedNode == nil {
		return
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.pressedNode.onPress != nil {
			g.pressedNode.onPress(false)
		}
		g.pressedNode = nil
	}
}

// handleKeyboardNav 处理 Tab 焦点切换、Enter/Space 激活、Esc 失焦。
// 决策逻辑抽到 focusNext/fireEscape/activateFocused，供无窗口的测试驱动复用。
func (g *game) handleKeyboardNav() {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.focusNext(!ebiten.IsKeyPressed(ebiten.KeyShift))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.fireEscape()
	}
	// Enter/Space 激活聚焦的可点击元素（输入框不拦截，交给文本编辑）
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.activateFocused()
	}
}

// focusNext 把焦点移到 Tab 顺序中的下一个（forward）或上一个可聚焦元素（含浮层），
// 环形回绕；返回新的焦点节点，无可聚焦元素时返回 nil。
func (g *game) focusNext(forward bool) *renderNode {
	var list []*renderNode
	collectFocusables(g.rootRN, &list)
	for _, pf := range g.portals {
		if pf.overlayRoot != nil {
			collectFocusables(pf.overlayRoot, &list)
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
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.dragging = nil
		return
	}
	x, y := ebiten.CursorPosition()
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
func (g *game) editFocusedInput() {
	f := g.focusedFiber
	if f == nil || f.unmounted || f.rnode == nil || f.rnode.kind != rnInput {
		return
	}
	rn := f.rnode
	val := rn.value
	caret := clampi(rn.caretPos, 0, len(val))
	anchor := clampi(rn.selAnchor, 0, len(val))

	shift := ebiten.IsKeyPressed(ebiten.KeyShift)
	ctrl := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)

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
		case inpututil.IsKeyJustPressed(ebiten.KeyA):
			anchor, caret = 0, len(val)
		case inpututil.IsKeyJustPressed(ebiten.KeyC):
			if hasSel() {
				setClipboard(val[selLo():selHi()])
			}
		case inpututil.IsKeyJustPressed(ebiten.KeyX):
			if hasSel() {
				setClipboard(val[selLo():selHi()])
				delSel()
			}
		case inpututil.IsKeyJustPressed(ebiten.KeyV):
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
		for _, r := range ebiten.AppendInputChars(nil) {
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

	if repeatKey(ebiten.KeyBackspace) {
		if hasSel() {
			delSel()
		} else if caret > 0 {
			_, sz := utf8.DecodeLastRuneInString(val[:caret])
			val = val[:caret-sz] + val[caret:]
			caret -= sz
			anchor = caret
		}
	}
	if repeatKey(ebiten.KeyDelete) {
		if hasSel() {
			delSel()
		} else if caret < len(val) {
			_, sz := utf8.DecodeRuneInString(val[caret:])
			val = val[:caret] + val[caret+sz:]
			anchor = caret
		}
	}
	if repeatKey(ebiten.KeyLeft) {
		if hasSel() && !shift {
			caret = selLo()
		} else if caret > 0 {
			_, sz := utf8.DecodeLastRuneInString(val[:caret])
			caret -= sz
		}
		afterMove()
	}
	if repeatKey(ebiten.KeyRight) {
		if hasSel() && !shift {
			caret = selHi()
		} else if caret < len(val) {
			_, sz := utf8.DecodeRuneInString(val[caret:])
			caret += sz
		}
		afterMove()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		caret = 0
		afterMove()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		caret = len(val)
		afterMove()
	}
	if rn.multiline && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
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
	x, y := ebiten.CursorPosition()
	now := map[*renderNode]bool{}
	for c := g.hitTop(float32(x), float32(y)); c != nil; c = c.parent {
		if c.onHover != nil {
			now[c] = true
		}
	}
	for rn := range now {
		if !g.hovered[rn] {
			rn.onHover(true)
		}
	}
	for rn := range g.hovered {
		if !now[rn] {
			rn.onHover(false)
		}
	}
	g.hovered = now
}

// repeatKey 在按下瞬间触发，长按后按固定间隔重复。
func repeatKey(k ebiten.Key) bool {
	d := inpututil.KeyPressDuration(k)
	return d == 1 || (d >= 30 && (d-30)%3 == 0)
}
