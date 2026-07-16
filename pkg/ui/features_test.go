package ui

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/sjm1327605995/tenon/pkg/font"
	"github.com/sjm1327605995/tenon/yoga"
)

func TestMultilineSpans(t *testing.T) {
	_ = font.InitDefaultFont()
	ff, err := font.GetDefaultFace(16, 400, false)
	if err != nil {
		t.Fatal(err)
	}
	// 含 \n 时无论宽度都按行分割，且字节偏移正确
	spans := wrapSpans("ab\ncd", ff.Face, 16*1.3, 0)
	if len(spans) != 2 {
		t.Fatalf("spans=%d want 2", len(spans))
	}
	if spans[0].text != "ab" || spans[0].start != 0 || spans[0].end != 2 {
		t.Fatalf("span0=%+v want {ab 0 2}", spans[0])
	}
	if spans[1].text != "cd" || spans[1].start != 3 || spans[1].end != 5 {
		t.Fatalf("span1=%+v want {cd 3 5}", spans[1])
	}
	// 行内偏移映射：行首 x 应回到该行起始字节
	if off := spans[1].offsetInSpan(0, ff.Face, 16*1.3); off != 3 {
		t.Fatalf("offsetInSpan(0)=%d want 3", off)
	}
}

func TestFontWeight(t *testing.T) {
	g := newGame()
	g.rootFiber = reconcile(nil, nil, Div(Style(Bold), // 容器设粗体 -> 后代继承
		Text("inherits-bold"),
		Text("explicit-normal", FontWeight(400)),
		Text("italic", Italic),
	))
	layoutAll(g)
	root := g.rootRN

	if w, f := root.children[0].effWeight, root.children[0].fauxBold; w != 700 || !f {
		t.Fatalf("inherited bold: weight=%d faux=%v want 700/true", w, f)
	}
	if w, f := root.children[1].effWeight, root.children[1].fauxBold; w != 400 || f {
		t.Fatalf("explicit normal: weight=%d faux=%v want 400/false", w, f)
	}
	if !root.children[2].fauxItalic {
		t.Fatal("italic text should be faux-italic (no italic variant registered)")
	}
}

func TestRichTextRuns(t *testing.T) {
	g := newGame()
	// 容器设红色 -> 未显式设色的 run 继承；显式设色/字重的 run 覆盖。
	g.mountRoot(Div(Style(TextColor(Hex("#ff0000"))),
		RichText(
			Text("plain "),
			Text("bold", Bold),
			Text(" big", FontSize(32)),
		),
	))
	rt := g.rootRN.children[0]
	if rt.kind != rnText || len(rt.runs) != 3 {
		t.Fatalf("rich text node: kind=%v runs=%d want text/3", rt.kind, len(rt.runs))
	}
	// 继承颜色
	if rt.runs[0].color != Hex("#ff0000") {
		t.Fatalf("run0 color=%v want inherited red", rt.runs[0].color)
	}
	// 字重覆盖 -> 合成粗体
	if !rt.runs[1].fauxBold {
		t.Fatal("run1 should be faux-bold")
	}
	// 字号覆盖 -> 更大的 lineH 与 ascent
	if rt.runs[2].lineH <= rt.runs[0].lineH {
		t.Fatalf("run2 lineH=%v should exceed run0 lineH=%v", rt.runs[2].lineH, rt.runs[0].lineH)
	}
	// 混排行高取行内最大值，节点高度应等于该行高（单行）
	if rt.bounds.H < float32(rt.runs[2].lineH) {
		t.Fatalf("rich node height=%v should be >= max run lineH=%v", rt.bounds.H, rt.runs[2].lineH)
	}
}

func TestRichTextLayoutBaseline(t *testing.T) {
	_ = font.InitDefaultFont()
	uiScale = 1
	runs := []textRun{
		{text: "A ", style: styleWith(FontSize(16))},
		{text: "B", style: styleWith(FontSize(32))},
	}
	(&renderNode{runs: runs}).resolveRuns(inhText{})
	lines, maxW, h := layoutRuns(runs, 0)
	if len(lines) != 1 {
		t.Fatalf("lines=%d want 1", len(lines))
	}
	if maxW <= 0 || h <= 0 {
		t.Fatalf("maxW=%v h=%v want positive", maxW, h)
	}
	// 行高应为较大字号的 lineH
	if lines[0].height < float32(runs[1].lineH)-0.5 {
		t.Fatalf("line height=%v want >= %v", lines[0].height, runs[1].lineH)
	}
	// 混排基线对齐：两段 draw-y 不同（小字号下移到共同基线）
	base := lines[0].ascent
	if base < runs[1].ascent-0.5 {
		t.Fatalf("line ascent=%v want >= max run ascent %v", base, runs[1].ascent)
	}
}

func TestWrapCacheInvalidation(t *testing.T) {
	g := newGame()
	g.mountRoot(Div(Style(Width(120)),
		Text("alpha beta gamma delta epsilon zeta"),
	))
	rn := g.rootRN.children[0]
	if rn.kind != rnText {
		t.Fatalf("expected text node, got kind %v", rn.kind)
	}
	l1, _ := rn.wrapped(120)
	if !rn.wc.valid || rn.wc.width != 120 {
		t.Fatal("cache not populated after first wrap")
	}
	l2, _ := rn.wrapped(120) // 同参再取应命中缓存：返回同一底层切片
	if len(l1) == 0 || &l1[0] != &l2[0] {
		t.Fatal("expected cache hit (same backing slice)")
	}
	rn.wrapped(60) // 宽度变化 -> 失效重算
	if rn.wc.width != 60 {
		t.Fatal("cache not updated for new width")
	}
	rn.text = "xy" // 文本变化 -> 不得返回旧折行
	l3, _ := rn.wrapped(120)
	if len(l3) != 1 || l3[0] != "xy" {
		t.Fatalf("stale wrap after text change: %v", l3)
	}
}

func TestRichCacheInvalidation(t *testing.T) {
	g := newGame()
	g.mountRoot(Div(RichText(Text("hello "), Text("world", Bold))))
	rn := g.rootRN.children[0]
	if len(rn.runs) == 0 {
		t.Fatal("no runs")
	}
	l1, _, _ := rn.richLayout(0)
	l2, _, _ := rn.richLayout(0) // 命中缓存
	if len(l1) == 0 || &l1[0] != &l2[0] {
		t.Fatal("expected rich cache hit")
	}
	rev := rn.runsRev
	// 模拟继承样式变化：改字号后重解析 -> runsRev 自增 -> 缓存失效
	rn.runs[0].style.fontSize, rn.runs[0].style.hasFontSize = 40, true
	rn.resolveRuns(inhText{})
	if rn.runsRev == rev {
		t.Fatal("runsRev should bump after run style change")
	}
	l3, _, _ := rn.richLayout(0)
	if &l3[0] == &l1[0] {
		t.Fatal("expected recompute after runsRev bump (got stale cached slice)")
	}
}

func BenchmarkWrapUncached(b *testing.B) {
	g := newGame()
	g.mountRoot(Div(Style(Width(120)), Text("alpha beta gamma delta epsilon zeta eta theta")))
	rn := g.rootRN.children[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wrapForWidth(rn.text, rn.face, rn.lineH, 120)
	}
}

func BenchmarkWrapCached(b *testing.B) {
	g := newGame()
	g.mountRoot(Div(Style(Width(120)), Text("alpha beta gamma delta epsilon zeta eta theta")))
	rn := g.rootRN.children[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rn.wrapped(120)
	}
}

// 黄金测试：通过录制后端断言"画了什么"，无需 GPU。
func TestPaintRecording(t *testing.T) {
	h := Mount(Use(func(_ struct{}) *Node {
		return Div(Style(Width(120), Height(48), Bg(Hex("#ef4444")), Radius(8), Border(2, Hex("#111111"))),
			Text("Hi"),
		)
	}, struct{}{}), 200, 100)
	ops := h.Paint()

	var fill, stroke, txt bool
	for _, op := range ops {
		switch op.Kind {
		case "rect":
			if op.Color == Hex("#ef4444") && op.Radius == 8 {
				fill = true
			}
		case "stroke":
			if op.Color == Hex("#111111") && op.Width == 2 {
				stroke = true
			}
		case "text":
			if op.Text == "Hi" {
				txt = true
			}
		}
	}
	if !fill || !stroke || !txt {
		t.Fatalf("paint ops missing: fill=%v stroke=%v text=%v\n%+v", fill, stroke, txt, ops)
	}
}

// 裁剪容器应在子节点绘制前后成对产生 clip / unclip 指令。
func TestPaintClipBalanced(t *testing.T) {
	h := Mount(Use(func(_ struct{}) *Node {
		return Div(Style(Width(100), Height(50), Clip), Text("x"))
	}, struct{}{}), 200, 100)
	ops := h.Paint()

	var clip, unclip, depth, maxDepth int
	for _, op := range ops {
		switch op.Kind {
		case "clip":
			clip++
			depth++
			if depth > maxDepth {
				maxDepth = depth
			}
		case "unclip":
			unclip++
			depth--
		}
	}
	if clip != 1 || unclip != 1 || depth != 0 || maxDepth != 1 {
		t.Fatalf("clip=%d unclip=%d endDepth=%d maxDepth=%d; want 1/1/0/1", clip, unclip, depth, maxDepth)
	}
}

// ErrorBoundary 捕获初次挂载时子组件的 panic，显示 fallback，不崩溃。
func TestErrorBoundaryCatchesMount(t *testing.T) {
	boom := func(_ struct{}) *Node { panic("boom") }
	app := func(_ struct{}) *Node {
		return ErrorBoundary(
			func(err any, _ func()) *Node { return Text("caught: " + ErrText(err)) },
			Use(boom, struct{}{}),
		)
	}
	h := Mount(Use(app, struct{}{}), 200, 100)
	if !h.Root().ByText("caught: boom").Exists() {
		t.Fatalf("fallback not rendered; texts=%v", h.Root().Texts())
	}
}

// 无边界时 panic 应照常抛出（不被静默吞掉）。
func TestNoBoundaryPanicsPropagate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic to propagate without a boundary")
		}
	}()
	boom := func(_ struct{}) *Node { panic("x") }
	Mount(Use(boom, struct{}{}), 100, 100)
}

// 触发再重试：出错 -> fallback -> retry 清错 -> 子树恢复正常。
func TestErrorBoundaryRetry(t *testing.T) {
	shouldBoom := true
	child := func(_ struct{}) *Node {
		if shouldBoom {
			panic("transient")
		}
		return Text("ok")
	}
	var retryFn func()
	app := func(_ struct{}) *Node {
		return ErrorBoundary(
			func(err any, retry func()) *Node { retryFn = retry; return Text("failed") },
			Use(child, struct{}{}),
		)
	}
	h := Mount(Use(app, struct{}{}), 200, 100)
	if !h.Root().ByText("failed").Exists() {
		t.Fatalf("expected fallback; texts=%v", h.Root().Texts())
	}
	shouldBoom = false // 修复"外部条件"
	retryFn()          // 重试
	h.Step(0)          // 让脏队列刷新
	if !h.Root().ByText("ok").Exists() {
		t.Fatalf("expected recovery after retry; texts=%v", h.Root().Texts())
	}
}

// 线性渐变背景应记录一条 gradient 填充指令。
func TestGradientFill(t *testing.T) {
	h := Mount(Use(func(_ struct{}) *Node {
		return Div(Style(Width(100), Height(40), LinearGradient(Hex("#ff0000"), Hex("#0000ff"), 90)))
	}, struct{}{}), 200, 100)
	ops := h.Paint()
	found := false
	for _, op := range ops {
		if op.Kind == "gradient" && op.Color == Hex("#ff0000") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected gradient fill op; ops=%+v", ops)
	}
}

// object-fit 的几何计算（不依赖真实图片）。
func TestObjectFitRect(t *testing.T) {
	box := Rect{X: 0, Y: 0, W: 200, H: 100}
	// 图片 100×100（1:1），放进 200×100 的框
	if dr, clip := fitRect(100, 100, box, FitFill); dr != box || clip {
		t.Fatalf("fill: %v clip=%v want full box no clip", dr, clip)
	}
	// contain：缩到高度 100，宽 100，水平居中 -> x=50
	if dr, clip := fitRect(100, 100, box, FitContain); dr != (Rect{50, 0, 100, 100}) || clip {
		t.Fatalf("contain: %v clip=%v want {50 0 100 100} no clip", dr, clip)
	}
	// cover：放大到宽度 200，高 200，垂直居中 -> y=-50，需裁剪
	if dr, clip := fitRect(100, 100, box, FitCover); dr != (Rect{0, -50, 200, 200}) || !clip {
		t.Fatalf("cover: %v clip=%v want {0 -50 200 200} clip", dr, clip)
	}
}

// 带圆角的裁剪容器应记录一条 radius>0 的 clip 指令（走圆角遮罩路径）。
func TestRoundedClipRecordsRadius(t *testing.T) {
	h := Mount(Use(func(_ struct{}) *Node {
		return Div(Style(Width(100), Height(50), Radius(12), Clip), Text("x"))
	}, struct{}{}), 200, 100)
	ops := h.Paint()
	found := false
	for _, op := range ops {
		if op.Kind == "clip" && op.Radius > 0 {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a rounded clip op (radius>0); ops=%+v", ops)
	}
}

// tokenize 现在遵循 UAX#14 换行规则。
func TestUAX14Tokenize(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"test-case", []string{"test-", "case"}}, // 连字符后可断
		{"foo bar", []string{"foo ", "bar"}},     // 拉丁词含尾随空格
		{"你好。世界", []string{"你", "好。", "世", "界"}}, // 收尾标点不落行首
		{"hello，world", []string{"hello，", "world"}},
	}
	for _, c := range cases {
		got := tokenize(c.in)
		if strings.Join(got, "") != c.in {
			t.Fatalf("tokenize(%q)=%q must concatenate back to input", c.in, got)
		}
		if len(got) != len(c.want) {
			t.Fatalf("tokenize(%q)=%q want %q", c.in, got, c.want)
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Fatalf("tokenize(%q)=%q want %q", c.in, got, c.want)
			}
		}
	}
}

// Icon：按尺寸测量，并通过录制后端确认真的描边/填充了一条路径（继承颜色）。
func TestIconRendersPath(t *testing.T) {
	h := Mount(Use(func(_ struct{}) *Node {
		return Div(Style(TextColor(Hex("#3b82f6"))), // 容器设色 -> 图标继承
			Icon(IconCheck, 20),           // 描边
			IconFill(IconChevronDown, 16), // 填充
		)
	}, struct{}{}), 200, 100)

	icons := h.Root().FindAll(func(q *Query) bool { return q.rn.kind == rnIcon })
	if len(icons) != 2 {
		t.Fatalf("icon nodes = %d want 2", len(icons))
	}
	// 尺寸测量：20×20 与 16×16（uiScale=1）
	if b := icons[0].Bounds(); b.W != 20 || b.H != 20 {
		t.Fatalf("icon0 bounds = %vx%v want 20x20", b.W, b.H)
	}
	// 绘制：一条描边路径 + 一条填充路径，颜色继承自容器
	ops := h.Paint()
	var stroke, fill int
	for _, op := range ops {
		if op.Kind == "strokepath" && op.Color == Hex("#3b82f6") {
			stroke++
		}
		if op.Kind == "path" && op.Color == Hex("#3b82f6") {
			fill++
		}
	}
	if stroke != 1 || fill != 1 {
		t.Fatalf("icon paint: stroke=%d fill=%d want 1/1\n%+v", stroke, fill, ops)
	}
}

// VirtualList 只渲染视口附近的少量行，滚动后窗口跟随移动。
func TestVirtualListWindowing(t *testing.T) {
	render := func(i int) *Node { return Text("row-" + strconv.Itoa(i)) }
	h := Mount(Use(func(_ struct{}) *Node {
		return VirtualList(VirtualListProps{Count: 1000, ItemHeight: 20, Height: 100, Render: render})
	}, struct{}{}), 300, 140)

	countRows := func() int {
		return len(h.Root().FindAll(func(q *Query) bool {
			return q.rn.kind == rnText && strings.HasPrefix(q.Text(), "row-")
		}))
	}
	// 视口 100 / 行高 20 = 5 行；加 overscan，应是十几行，绝非 1000
	if n := countRows(); n == 0 || n > 40 {
		t.Fatalf("rendered rows = %d, want a small window (not ~1000)", n)
	}
	if !h.Root().ByText("row-0").Exists() {
		t.Fatal("row-0 should render at top")
	}
	if h.Root().ByText("row-500").Exists() {
		t.Fatal("row-500 should NOT render (far off-screen)")
	}
	// 向下滚动到约第 25 行处
	if !h.Root().ByKind("scroll").ScrollBy(500) {
		t.Fatal("no scroll container found")
	}
	if h.Root().ByText("row-0").Exists() {
		t.Fatal("row-0 should be virtualized out after scrolling")
	}
	if !h.Root().ByText("row-25").Exists() {
		t.Fatalf("row-25 should be in window after scroll; texts=%v", h.Root().Texts())
	}
	if n := countRows(); n > 40 {
		t.Fatalf("window grew unexpectedly after scroll: %d rows", n)
	}
}

// ArrowNav：方向键在导航组内环形移动焦点；组外方向不动作。
func TestArrowNavRovingFocus(t *testing.T) {
	comp := func(_ struct{}) *Node {
		return Div(
			Button(OnClick(func() {}), Text("outside")),
			Div(Style(Column), ArrowNav(NavVertical),
				Button(OnClick(func() {}), Text("i0")),
				Button(OnClick(func() {}), Text("i1")),
				Button(OnClick(func() {}), Text("i2")),
			),
		)
	}
	h := Mount(Use(comp, struct{}{}), 300, 300)

	// 先聚焦到组内第一个项（按钮，而非其文本子节点）
	byLabel := func(s string) *Query {
		return h.Root().Find(func(q *Query) bool { return q.Clickable() && q.AllText() == s })
	}
	byLabel("i0").Focus()
	if got := h.Arrow(NavVertical, true).AllText(); got != "i1" {
		t.Fatalf("Down -> %q want i1", got)
	}
	if got := h.Arrow(NavVertical, true).AllText(); got != "i2" {
		t.Fatalf("Down -> %q want i2", got)
	}
	if got := h.Arrow(NavVertical, true).AllText(); got != "i0" { // 环形回绕
		t.Fatalf("Down wrap -> %q want i0", got)
	}
	if got := h.Arrow(NavVertical, false).AllText(); got != "i2" { // 反向回绕
		t.Fatalf("Up wrap -> %q want i2", got)
	}
	// 水平方向在纵向组里不动作
	if got := h.Arrow(NavHorizontal, true).AllText(); got != "i2" {
		t.Fatalf("Horizontal in vertical group moved focus -> %q want i2", got)
	}
	// 焦点在组外时方向键不动作
	byLabel("outside").Focus()
	if got := h.Arrow(NavVertical, true).AllText(); got != "outside" {
		t.Fatalf("arrow outside group moved focus -> %q want outside", got)
	}
}

// 模态（Portal + TrapFocus）打开时 Tab 只在浮层内循环，不逃到背景。
func TestFocusTrapInModal(t *testing.T) {
	comp := func(_ struct{}) *Node {
		return Div(
			Button(OnClick(func() {}), Text("bg1")),
			Button(OnClick(func() {}), Text("bg2")),
			Portal(TrapFocus(),
				Button(OnClick(func() {}), Text("m1")),
				Button(OnClick(func() {}), Text("m2")),
			),
		)
	}
	h := Mount(Use(comp, struct{}{}), 400, 300)

	seen := map[string]bool{}
	for i := 0; i < 4; i++ {
		seen[h.Tab().AllText()] = true
	}
	if seen["bg1"] || seen["bg2"] {
		t.Fatalf("focus escaped modal into background: %v", seen)
	}
	if !seen["m1"] || !seen["m2"] {
		t.Fatalf("modal items should be reachable: %v", seen)
	}
}

// 无 TrapFocus 的普通浮层（如下拉/提示）不应限制焦点。
func TestNonModalPortalDoesNotTrap(t *testing.T) {
	comp := func(_ struct{}) *Node {
		return Div(
			Button(OnClick(func() {}), Text("bg1")),
			Portal(Button(OnClick(func() {}), Text("p1"))),
		)
	}
	h := Mount(Use(comp, struct{}{}), 400, 300)
	seen := map[string]bool{}
	for i := 0; i < 3; i++ {
		seen[h.Tab().AllText()] = true
	}
	if !seen["bg1"] || !seen["p1"] {
		t.Fatalf("non-modal portal should allow focus everywhere: %v", seen)
	}
}

func TestGraphemeBoundaries(t *testing.T) {
	// "a👍b"：👍 是 4 字节单个字素簇，退格/方向应把它当一个字符
	s := "a👍b"
	if got := prevGraphemeBoundary(s, 5); got != 1 { // 光标在 👍 之后(byte5) 退一格 -> 1
		t.Fatalf("prevGrapheme(%q,5)=%d want 1", s, got)
	}
	if got := nextGraphemeBoundary(s, 1); got != 5 { // 从 👍 起始前进一格 -> 越过整个 emoji
		t.Fatalf("nextGrapheme(%q,1)=%d want 5", s, got)
	}
	// 组合字符 e + U+0301 (é 分解形式) 应作为一个字素簇
	combining := "éx"
	if got := prevGraphemeBoundary(combining, 3); got != 0 {
		t.Fatalf("prevGrapheme(combining,3)=%d want 0", got)
	}
}

func TestWordBoundaries(t *testing.T) {
	s := "foo bar baz"
	if got := nextWordBoundary(s, 0); got != 3 { // "foo" 末
		t.Fatalf("nextWord(0)=%d want 3", got)
	}
	if got := nextWordBoundary(s, 3); got != 7 { // 跳过空格 -> "bar" 末
		t.Fatalf("nextWord(3)=%d want 7", got)
	}
	if got := prevWordBoundary(s, 11); got != 8 { // "baz" 首
		t.Fatalf("prevWord(11)=%d want 8", got)
	}
	if got := prevWordBoundary(s, 7); got != 4 { // "bar" 首
		t.Fatalf("prevWord(7)=%d want 4", got)
	}
	a, b := wordAt(s, 5) // 点在 "bar" 内
	if a != 4 || b != 7 {
		t.Fatalf("wordAt(5)=(%d,%d) want (4,7)", a, b)
	}
}

// 退格应按字素簇删除整个 emoji，而不是拆成半个码点。
func TestBackspaceDeletesGrapheme(t *testing.T) {
	val := "a👍"
	caret := len(val) // 6
	// 复刻 editFocusedInput 的退格逻辑（非 ctrl）
	prev := prevGraphemeBoundary(val, caret)
	val = val[:prev] + val[caret:]
	if val != "a" {
		t.Fatalf("after backspace = %q want %q", val, "a")
	}
}

func styleWith(opts ...StyleOpt) StyleProps {
	st := newStyleProps()
	for _, o := range opts {
		o(&st)
	}
	return st
}

// 预编辑（IME 组字串）以虚拟方式插入到 caret，不改变受控 value。
func TestPreeditRenderInsertion(t *testing.T) {
	rn := &renderNode{kind: rnInput, value: "abcd", caretPos: 2, selAnchor: 2,
		preedit: "XY", preeditAt: 2, preeditCaret: 2}
	// 复现 paintInput 中的显示串推导
	val := rn.value
	at := rn.preeditAt
	disp := val[:at] + rn.preedit + val[at:]
	if disp != "abXYcd" {
		t.Fatalf("display=%q want abXYcd", disp)
	}
	if rn.value != "abcd" {
		t.Fatalf("controlled value mutated to %q", rn.value)
	}
	caret := at + rn.preeditCaret
	if caret != 4 {
		t.Fatalf("preedit caret=%d want 4", caret)
	}
}

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

func TestOnClickWiring(t *testing.T) {
	g := newGame()
	got := false
	g.rootFiber = reconcile(nil, nil, Div(OnClick(func() { got = true })))
	layoutAll(g)
	if g.rootRN.onClick == nil {
		t.Fatal("node with OnClick should be clickable")
	}
	g.rootRN.onClick()
	if !got {
		t.Fatal("OnClick handler should fire")
	}
}

func TestOnDragWiring(t *testing.T) {
	g := newGame()
	var got float32 = -1
	g.rootFiber = reconcile(nil, nil, Div(OnDrag(func(dx, dy float32) { got = dx })))
	layoutAll(g)
	d := findDraggable(g.rootRN)
	if d == nil {
		t.Fatal("node with OnDrag should be draggable")
	}
	d.onDrag(42, 0)
	if got != 42 {
		t.Fatalf("OnDrag handler got=%v want 42", got)
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
