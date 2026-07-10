package ui

import (
	"fmt"
	"testing"

	"github.com/sjm1327605995/tenon/pkg/font"
	"github.com/sjm1327605995/tenon/yoga"
)

func newGame() *game {
	_ = font.InitDefaultFont()
	g := &game{w: 400, h: 400}
	activeGame = g
	return g
}

func (g *game) mountRoot(n *Node) {
	g.rootFiber = reconcile(nil, nil, n)
	layoutAll(g)
}

func (g *game) drain() {
	for guard := 0; len(g.dirty) > 0 && guard < 100; guard++ {
		g.flushDirty()
	}
	layoutAll(g)
}

func firstText(g *game) string {
	var ts []string
	collectTexts(g.rootRN, &ts)
	if len(ts) > 0 {
		return ts[0]
	}
	return ""
}

// ---- Context 跨 Memo 边界传播 ----

var (
	themeCtx        = CreateContext("light")
	ctxChildRenders int
	themeSetter     func(string)
)

func ctxChild(_ struct{}) *Node {
	ctxChildRenders++
	return Text(UseContext(themeCtx))
}

func ctxParent(_ struct{}) *Node {
	theme, set := UseState("light")
	themeSetter = set
	return themeCtx.Provider(theme, Memo(ctxChild, struct{}{}))
}

func TestContextThroughMemo(t *testing.T) {
	g := newGame()
	ctxChildRenders = 0
	g.mountRoot(Use(ctxParent, struct{}{}))

	if firstText(g) != "light" || ctxChildRenders != 1 {
		t.Fatalf("initial: text=%q renders=%d", firstText(g), ctxChildRenders)
	}

	themeSetter("dark")
	g.drain()

	if firstText(g) != "dark" {
		t.Fatalf("after provider change, text=%q want dark", firstText(g))
	}
	if ctxChildRenders != 2 {
		t.Fatalf("memo child renders=%d, want 2 (context change must pierce memo)", ctxChildRenders)
	}
}

// ---- useReducer ----

func TestUseReducer(t *testing.T) {
	g := newGame()
	var dispatch func(int)
	comp := func(_ struct{}) *Node {
		count, d := UseReducer(func(s, a int) int { return s + a }, 0)
		dispatch = d
		return Text(fmt.Sprintf("%d", count))
	}
	g.mountRoot(Use(comp, struct{}{}))
	if firstText(g) != "0" {
		t.Fatalf("initial=%q", firstText(g))
	}
	dispatch(5)
	g.drain()
	if firstText(g) != "5" {
		t.Fatalf("after dispatch(5)=%q want 5", firstText(g))
	}
	dispatch(-2)
	g.drain()
	if firstText(g) != "3" {
		t.Fatalf("after dispatch(-2)=%q want 3", firstText(g))
	}
}

// ---- Memo 短路 ----

var (
	memoRenders  int
	plainRenders int
	memoTickSet  func(int)
)

func memoChild(p struct{ X int }) *Node  { memoRenders++; return Text(fmt.Sprintf("m%d", p.X)) }
func plainChild(p struct{ X int }) *Node { plainRenders++; return Text(fmt.Sprintf("p%d", p.X)) }

func memoParent(_ struct{}) *Node {
	_, set := UseState(0)
	memoTickSet = set
	return Div(
		Memo(memoChild, struct{ X int }{X: 1}),
		Use(plainChild, struct{ X int }{X: 1}),
	)
}

func TestMemoBailout(t *testing.T) {
	g := newGame()
	memoRenders, plainRenders = 0, 0
	g.mountRoot(Use(memoParent, struct{}{}))
	if memoRenders != 1 || plainRenders != 1 {
		t.Fatalf("initial memo=%d plain=%d", memoRenders, plainRenders)
	}
	memoTickSet(1) // parent re-renders, children props unchanged
	g.drain()
	if memoRenders != 1 {
		t.Fatalf("memo child re-rendered (=%d); should bail on equal props", memoRenders)
	}
	if plainRenders != 2 {
		t.Fatalf("plain child renders=%d, want 2 (default cascade)", plainRenders)
	}
}

// ---- Memo + 稳定回调 props：带 func 的 props 仍能短路 ----

type cbProps struct {
	V int
	F func(int)
}

var (
	cbChildRenders int
	cbTickSet      func(int)
)

func cbChild(p cbProps) *Node { cbChildRenders++; return Text(fmt.Sprintf("%d", p.V)) }

func cbParent(_ struct{}) *Node {
	_, setTick := UseState(0)
	cbTickSet = setTick
	_, dispatch := UseReducer(func(s, a int) int { return s + a }, 0)
	return Div(Memo(cbChild, cbProps{V: 1, F: dispatch}))
}

func TestMemoStableCallback(t *testing.T) {
	g := newGame()
	cbChildRenders = 0
	g.mountRoot(Use(cbParent, struct{}{}))
	if cbChildRenders != 1 {
		t.Fatalf("initial cbChildRenders=%d", cbChildRenders)
	}
	cbTickSet(1) // parent re-renders; child props: V unchanged, F is stable dispatch
	g.drain()
	if cbChildRenders != 1 {
		t.Fatalf("memo child re-rendered (=%d); stable callback prop should still bail", cbChildRenders)
	}
}

// ---- UseTween：数值补间动画 ----

func TestUseTween(t *testing.T) {
	g := newGame()
	var setTarget func(float32)
	var last float32
	comp := func(_ struct{}) *Node {
		target, set := UseState(float32(0))
		setTarget = set
		last = UseTween(target, 100, Linear) // 100ms 线性
		return Text(fmt.Sprintf("%.0f", last))
	}
	g.mountRoot(Use(comp, struct{}{}))
	if last != 0 {
		t.Fatalf("initial=%v", last)
	}

	setTarget(100)
	g.drain()
	if last != 0 || len(g.anims) != 1 {
		t.Fatalf("after target change last=%v anims=%d", last, len(g.anims))
	}

	g.tickAnims(50)
	g.drain()
	if last < 45 || last > 55 {
		t.Fatalf("mid last=%v want ~50", last)
	}

	g.tickAnims(60) // 越过终点
	g.drain()
	if last != 100 {
		t.Fatalf("end last=%v want 100", last)
	}
	if len(g.anims) != 0 {
		t.Fatalf("anims not cleared: %d", len(g.anims))
	}
}

// ---- Transform / 图层：needsLayer 判定 ----

func TestLayerFlags(t *testing.T) {
	g := newGame()
	g.rootFiber = reconcile(nil, nil, Div(Style(Scale(1.2)),
		Div(Style(Opacity(0.5)), Text("x")),
		Div(Style(Width(10), Height(10)))))
	layoutAll(g)

	outer := g.rootRN
	if !outer.hasTransform() || !outer.needsLayer() {
		t.Fatalf("outer: transform=%v layer=%v", outer.hasTransform(), outer.needsLayer())
	}
	withOpacity := outer.children[0]
	if withOpacity.hasTransform() {
		t.Fatal("opacity child should not have transform")
	}
	if !withOpacity.needsLayer() {
		t.Fatal("opacity<1 with children should need a layer (group opacity)")
	}
	plain := outer.children[1]
	if plain.needsLayer() {
		t.Fatal("plain box should not need a layer")
	}
}

// ---- Portal：脱离父树、绘制/命中在主树之上 ----

func inSubtree(root, n *renderNode) bool {
	if root == n {
		return true
	}
	for _, c := range root.children {
		if inSubtree(c, n) {
			return true
		}
	}
	return false
}

func TestPortalOverlay(t *testing.T) {
	g := newGame()
	comp := func(_ struct{}) *Node {
		return Div(Style(Width(400), Height(400)),
			Div(Style(Width(400), Height(400)), OnClick(func() {})), // 主树背景
			Portal(
				Div(Style(Grow(1), ItemsCenter, JustifyCenter, Bg(Color{0, 0, 0, 120})),
					Div(Style(Width(200), Height(120), Bg(White)), OnClick(func() {})), // 居中对话框
				),
			),
		)
	}
	g.mountRoot(Use(comp, struct{}{}))

	if len(g.portals) != 1 {
		t.Fatalf("portals=%d want 1", len(g.portals))
	}
	// 屏幕中心命中的应是浮层里的对话框，而非主树背景
	top := g.hitTop(200, 200)
	if top == nil {
		t.Fatal("hitTop returned nil")
	}
	if !inSubtree(g.portals[0].overlayRoot, top) {
		t.Fatal("center hit should land in the portal overlay")
	}
	if inSubtree(g.rootRN, top) {
		t.Fatal("portal content must not live in the main render tree")
	}
}

// ---- 键盘焦点导航 ----

func TestFocusNav(t *testing.T) {
	g := newGame()
	g.rootFiber = reconcile(nil, nil, Div(Style(Column),
		Button(Style(Width(50), Height(20)), OnClick(func() {}), Text("a")),
		Input(Style(Width(50), Height(20))),
		Div(Style(Width(50), Height(20))), // 不可聚焦
		Button(Style(Width(50), Height(20)), OnClick(func() {}), Text("b")),
	))
	layoutAll(g)

	var list []*renderNode
	collectFocusables(g.rootRN, &list)
	if len(list) != 3 {
		t.Fatalf("focusables=%d want 3 (2 buttons + 1 input)", len(list))
	}
	if nextFocus(list, nil, true) != list[0] {
		t.Fatal("nil forward should focus first")
	}
	if nextFocus(list, list[0], true) != list[1] {
		t.Fatal("forward from 0 should be 1")
	}
	if nextFocus(list, list[2], true) != list[0] {
		t.Fatal("forward should wrap to 0")
	}
	if nextFocus(list, list[0], false) != list[2] {
		t.Fatal("backward from 0 should wrap to last")
	}
}

// ---- 文本折行 ----

func TestTextWrap(t *testing.T) {
	g := newGame()
	long := "The quick brown fox jumps over the lazy dog again and again and again"
	g.rootFiber = reconcile(nil, nil, Div(Style(Width(140)), Text(long, FontSize(16))))
	layoutAll(g)

	txt := g.rootRN.children[0]
	if txt.bounds.H <= float32(txt.lineH)+1 {
		t.Fatalf("text height %v not wrapped (single lineH %v)", txt.bounds.H, txt.lineH)
	}
	if txt.bounds.W > 141 {
		t.Fatalf("wrapped text width %v exceeds container 140", txt.bounds.W)
	}
}

// ---- 扩展库基座：主题 + 交互 ----

func TestUseTheme(t *testing.T) {
	g := newGame()
	var got Theme
	child := func(_ struct{}) *Node { got = UseTheme(); return Text("x") }
	app := func(_ struct{}) *Node { return ThemeProvider(DarkTheme, Use(child, struct{}{})) }
	g.mountRoot(Use(app, struct{}{}))
	if got.Background != DarkTheme.Background {
		t.Fatalf("UseTheme got Background=%v want dark %v", got.Background, DarkTheme.Background)
	}

	// 无 Provider 时回退 LightTheme
	var def Theme
	solo := func(_ struct{}) *Node { def = UseTheme(); return Text("y") }
	g2 := newGame()
	g2.mountRoot(Use(solo, struct{}{}))
	if def.Background != LightTheme.Background {
		t.Fatal("UseTheme without provider should be LightTheme")
	}
}

func TestOnPressWiring(t *testing.T) {
	g := newGame()
	last := false
	g.rootFiber = reconcile(nil, nil, Div(Style(Width(40), Height(40)),
		OnPress(func(b bool) { last = b })))
	layoutAll(g)
	if g.rootRN.onPress == nil {
		t.Fatal("onPress not wired to render node")
	}
	g.rootRN.onPress(true)
	if !last {
		t.Fatal("press-down not delivered")
	}
	g.rootRN.onPress(false)
	if last {
		t.Fatal("press-up not delivered")
	}
}

// ---- 增量 relink：纯绘制变更不脏化 yoga（避免整树重算）----

func TestIncrementalRelink(t *testing.T) {
	g := newGame()
	var setBg func(Color)
	comp := func(_ struct{}) *Node {
		bg, set := UseState(White)
		setBg = set
		return Div(Style(Width(100), Height(100), Bg(bg)),
			Div(Style(Width(50), Height(50))))
	}
	g.mountRoot(Use(comp, struct{}{})) // 初次布局后 yoga 干净

	setBg(Red)          // 仅改背景色（纯绘制）
	g.flushDirty()      // 重渲染 + 协调，不布局
	relink(g.rootFiber) // 增量：结构未变，不应脏化 yoga
	if g.rootRN.yn.IsDirty() {
		t.Fatal("paint-only change dirtied yoga; CalculateLayout would redo full solve")
	}
}

// ---- 输入光标定位 & 剪贴板 ----

func TestInputCaretFromX(t *testing.T) {
	g := newGame()
	g.rootFiber = reconcile(nil, nil, Input(Style(Width(200), Height(30), FontSize(16)), Value("hello")))
	layoutAll(g)
	rn := g.rootRN
	if rn.kind != rnInput {
		t.Fatalf("root kind=%v want input", rn.kind)
	}
	if c := rn.caretFromX(rn.bounds.X - 10); c != 0 {
		t.Fatalf("far-left caret=%d want 0", c)
	}
	if c := rn.caretFromX(rn.bounds.X + 1000); c != len("hello") {
		t.Fatalf("far-right caret=%d want %d", c, len("hello"))
	}
}

func TestClipboard(t *testing.T) {
	defer func() {
		getClipboard = func() string { return clipboardText }
		setClipboard = func(s string) { clipboardText = s }
	}()
	SetClipboardText("hi")
	if Clipboard() != "hi" {
		t.Fatalf("clipboard = %q want hi", Clipboard())
	}
	var store string
	SetClipboardProvider(func() string { return store }, func(s string) { store = s })
	SetClipboardText("world")
	if Clipboard() != "world" || store != "world" {
		t.Fatal("clipboard provider override failed")
	}
}

// ---- Post：跨 goroutine 安全更新 ----

func TestPost(t *testing.T) {
	g := newGame()
	got := 0
	var setV func(int)
	comp := func(_ struct{}) *Node {
		v, set := UseState(0)
		setV = set
		got = v
		return Text(fmt.Sprintf("%d", v))
	}
	g.mountRoot(Use(comp, struct{}{}))

	done := make(chan struct{})
	go func() { // 模拟后台任务
		Post(func() { setV(42) })
		close(done)
	}()
	<-done

	drainPosts() // 渲染线程排空 -> setV 在此执行
	g.drain()
	if got != 42 {
		t.Fatalf("state after Post = %d, want 42", got)
	}
}

// ---- 布局门控：干净帧跳过 computeBounds ----

func TestLayoutGating(t *testing.T) {
	g := newGame()
	g.rootFiber = reconcile(nil, nil, Div(Style(Width(100), Height(100)),
		Div(Style(Width(50), Height(50)))))
	g.rootRN = rootRenderNode(g.rootFiber)
	g.layout() // 初次：计算 bounds

	child := g.rootRN.children[0]
	child.bounds = Rect{X: -999, Y: -999} // 人为破坏，检测 computeBounds 是否重跑

	g.layout() // 无变化的第二次布局：应跳过 computeBounds
	if child.bounds.X != -999 {
		t.Fatal("computeBounds ran on a clean layout; gating failed")
	}

	// 窗口变化应触发重算，恢复正确 bounds
	g.w += 10
	g.layout()
	if child.bounds.X == -999 {
		t.Fatal("computeBounds skipped after resize; should recompute")
	}
}

// ---- UseEscape：栈式 Esc 处理器 ----

func TestUseEscape(t *testing.T) {
	g := newGame()
	fired := 0
	var setActive func(bool)
	comp := func(_ struct{}) *Node {
		active, set := UseState(true)
		setActive = set
		UseEscape(active, func() { fired++ })
		return Text("x")
	}
	g.mountRoot(Use(comp, struct{}{}))
	g.flushEffects()
	if len(g.escStack) != 1 {
		t.Fatalf("escStack=%d want 1", len(g.escStack))
	}
	(*g.escStack[0].fn)()
	if fired != 1 {
		t.Fatal("escape handler not fired")
	}

	setActive(false)
	g.drain()
	g.flushEffects()
	if len(g.escStack) != 0 {
		t.Fatalf("escStack=%d want 0 after deactivate", len(g.escStack))
	}
}

// ---- UseMeasure：读回元素屏幕矩形（锚定基座）----

func TestUseMeasure(t *testing.T) {
	g := newGame()
	var got Rect
	comp := func(_ struct{}) *Node {
		ref, rect := UseMeasure()
		got = rect
		return Div(ref, Style(Width(120), Height(40), Absolute, Left(20), Top(10)))
	}
	g.mountRoot(Use(comp, struct{}{}))
	g.drain() // 测量写回后重渲染，读到真实矩形

	if got.X != 20 || got.Y != 10 || got.W != 120 || got.H != 40 {
		t.Fatalf("measured rect = %+v want {20 10 120 40}", got)
	}
}

// ---- 基础组件 ----

func findDraggable(rn *renderNode) *renderNode {
	if rn.onDrag != nil {
		return rn
	}
	for _, c := range rn.children {
		if r := findDraggable(c); r != nil {
			return r
		}
	}
	return nil
}

func TestCheckbox(t *testing.T) {
	g := newGame()
	got := false
	g.rootFiber = reconcile(nil, nil, Checkbox(false, func(v bool) { got = v }))
	layoutAll(g)
	if g.rootRN.onClick == nil {
		t.Fatal("checkbox should be clickable")
	}
	g.rootRN.onClick()
	if !got {
		t.Fatal("checkbox onChange should receive true")
	}
}

func TestSliderClamp(t *testing.T) {
	g := newGame()
	var got float32 = -1
	g.rootFiber = reconcile(nil, nil, Slider(90, 0, 100, func(v float32) { got = v }))
	layoutAll(g)
	d := findDraggable(g.rootRN)
	if d == nil {
		t.Fatal("slider thumb should be draggable")
	}
	d.onDrag(10000, 0) // 远超范围，应钳制到 max
	if got != 100 {
		t.Fatalf("slider clamp got=%v want 100", got)
	}
}

// ---- 文本样式继承 ----

func TestTextInheritance(t *testing.T) {
	g := newGame()
	red := Hex("#ff0000")
	blue := Hex("#0000ff")
	g.rootFiber = reconcile(nil, nil, Div(Style(TextColor(red), FontSize(24)),
		Text("inherits"), // 继承 -> red/24
		Text("override", TextColor(blue), FontSize(12)), // 覆盖 -> blue/12
		Div(Style(), // 透传
			Text("nested"), // -> red/24
		),
	))
	layoutAll(g)
	root := g.rootRN

	if c, s := root.children[0].color, root.children[0].effSize; c != red || s != 24 {
		t.Fatalf("inherit: color=%v size=%v want red/24", c, s)
	}
	if c, s := root.children[1].color, root.children[1].effSize; c != blue || s != 12 {
		t.Fatalf("override: color=%v size=%v want blue/12", c, s)
	}
	nested := root.children[2].children[0]
	if c, s := nested.color, nested.effSize; c != red || s != 24 {
		t.Fatalf("nested inherit: color=%v size=%v want red/24", c, s)
	}
}

// ---- 布局动画（FLIP）：位置变化注入偏移并衰减 ----

func TestFlipOffset(t *testing.T) {
	rn := &renderNode{animatedLayout: true}

	walkLayoutAnim(rn, 0.8) // 首帧：记录 prevPos，无偏移
	if rn.offY != 0 {
		t.Fatalf("initial offY=%v want 0", rn.offY)
	}

	rn.bounds.Y = 100 // 位置跳到 100
	walkLayoutAnim(rn, 0.8)
	if rn.offY != -80 { // (0-100) 再 *0.8
		t.Fatalf("after move offY=%v want -80", rn.offY)
	}

	walkLayoutAnim(rn, 0.8) // 位置不变，继续衰减
	if rn.offY != -64 {
		t.Fatalf("decay offY=%v want -64", rn.offY)
	}

	for i := 0; i < 60; i++ {
		walkLayoutAnim(rn, 0.8)
	}
	if rn.offY != 0 {
		t.Fatalf("did not settle: offY=%v", rn.offY)
	}
}

// ---- 变换感知的命中测试 ----

func TestTransformHitTest(t *testing.T) {
	g := newGame()
	// 100x100 的盒子在 (0,0)，以中心 (50,50) 放大 2 倍 -> 视觉横跨 [-50,150]。
	g.rootFiber = reconcile(nil, nil, Div(Style(Width(100), Height(100), Scale(2))))
	layoutAll(g)
	rn := g.rootRN

	// (140,50) 在放大后的视觉内，但在未变换 bounds(0..100) 之外 —— 应命中。
	if hitNode(rn, 140, 50) != rn {
		t.Fatal("scaled element should be hit at (140,50)")
	}
	// (160,50) 超出视觉右边界(150) —— 应落空。
	if hitNode(rn, 160, 50) != nil {
		t.Fatal("(160,50) is outside scaled visual, should miss")
	}
}

// ---- 拖拽回调接线 ----

func TestDragWiring(t *testing.T) {
	g := newGame()
	var total float32
	g.rootFiber = reconcile(nil, nil,
		Div(Style(Width(50), Height(50)), OnDrag(func(dx, _ float32) { total += dx })))
	layoutAll(g)
	rn := g.rootRN
	if rn.onDrag == nil {
		t.Fatal("onDrag not wired to render node")
	}
	rn.onDrag(5, 0)
	rn.onDrag(3, 0)
	if total != 8 {
		t.Fatalf("accumulated drag=%v want 8", total)
	}
}

// ---- UseTransition：退场动画结束后才卸载 ----

func TestUseTransition(t *testing.T) {
	g := newGame()
	var setVisible func(bool)
	var mounted bool
	var prog float32
	comp := func(_ struct{}) *Node {
		vis, set := UseState(false)
		setVisible = set
		mounted, prog = UseTransition(vis, 100)
		return If(mounted, Text("x"))
	}
	g.mountRoot(Use(comp, struct{}{}))
	if mounted {
		t.Fatal("initial should not be mounted")
	}

	setVisible(true)
	g.drain()
	if !mounted || prog != 0 {
		t.Fatalf("on show: mounted=%v prog=%v want true/0", mounted, prog)
	}

	g.tickAnims(200) // 进场结束
	g.drain()
	if !mounted || prog != 1 {
		t.Fatalf("after enter: mounted=%v prog=%v want true/1", mounted, prog)
	}

	setVisible(false)
	g.drain()
	if !mounted {
		t.Fatal("should stay mounted while exit animates")
	}

	g.tickAnims(200) // 退场结束
	g.drain()
	if mounted || prog != 0 {
		t.Fatalf("after exit: mounted=%v prog=%v want false/0", mounted, prog)
	}
}

// ---- ScrollView：滚动偏移与边界钳制 ----

func TestScrollOffset(t *testing.T) {
	g := newGame()
	kids := []*Node{Style(Height(100), Width(120), Column)}
	for i := 0; i < 5; i++ {
		kids = append(kids, Div(Style(Height(40), Width(100)), Text(fmt.Sprintf("row%d", i))))
	}
	g.rootFiber = reconcile(nil, nil, ScrollView(kids...))
	layoutAll(g)

	sc := g.rootRN
	if sc.kind != rnScroll {
		t.Fatalf("root kind=%v want scroll", sc.kind)
	}
	if sc.contentH != 200 {
		t.Fatalf("contentH=%v want 200", sc.contentH)
	}
	y0 := sc.children[0].bounds.Y

	sc.scrollY = 50
	computeBounds(sc, 0, 0)
	if got := sc.children[0].bounds.Y; got != y0-50 {
		t.Fatalf("after scroll child0.Y=%v want %v", got, y0-50)
	}

	// 过度滚动应钳制到 maxScroll = content(200) - viewport(100) = 100
	sc.scrollY = 999
	computeBounds(sc, 0, 0)
	if sc.scrollY != 100 {
		t.Fatalf("scrollY=%v want clamped 100", sc.scrollY)
	}
}

// ---- keyed list：重排复用 Fiber，不重新挂载 ----

var itemMounts map[string]int

func keyedItem(p struct{ Label string }) *Node {
	UseEffect(func() Cleanup {
		itemMounts[p.Label]++
		return nil
	}) // 无 deps 但空切片：仅挂载时运行一次
	return Text(p.Label)
}

func keyedList(order []string) *Node {
	kids := make([]*Node, len(order))
	for i, lbl := range order {
		kids[i] = Fragment(Key(lbl), Use(keyedItem, struct{ Label string }{Label: lbl}))
	}
	return Div(kids...)
}

func TestKeyedListReuse(t *testing.T) {
	g := newGame()
	itemMounts = map[string]int{}

	g.rootFiber = reconcile(nil, nil, keyedList([]string{"A", "B", "C"}))
	layoutAll(g)
	g.flushEffects()
	if itemMounts["A"] != 1 || itemMounts["B"] != 1 || itemMounts["C"] != 1 {
		t.Fatalf("initial mounts = %v", itemMounts)
	}

	// 重排 —— 相同 key，应复用而非重挂载
	updateFiber(g.rootFiber, keyedList([]string{"C", "A", "B"}))
	layoutAll(g)
	g.flushEffects()
	if itemMounts["A"] != 1 || itemMounts["B"] != 1 || itemMounts["C"] != 1 {
		t.Fatalf("after reorder mounts = %v (should be unchanged)", itemMounts)
	}

	// 顺序应为 C,A,B
	var ts []string
	collectTexts(g.rootRN, &ts)
	if fmt.Sprint(ts) != "[C A B]" {
		t.Fatalf("order after reorder = %v want [C A B]", ts)
	}

	// 删除 B：其 Fragment 应卸载
	updateFiber(g.rootFiber, keyedList([]string{"C", "A"}))
	layoutAll(g)
	ts = nil
	collectTexts(g.rootRN, &ts)
	if fmt.Sprint(ts) != "[C A]" {
		t.Fatalf("after delete = %v want [C A]", ts)
	}
	_ = yoga.DirectionLTR
}
