package ui

import (
	"fmt"
	"testing"
)

// ---------- Store ----------

type cartState struct {
	Count int
	Note  string
}

var (
	cartStore = NewStore(cartState{Count: 1, Note: "a"})

	cartRendersCount, cartRendersNote int
)

func cartCountView(_ struct{}) *Node {
	cartRendersCount++
	n := UseStore(cartStore, func(s cartState) int { return s.Count })
	return Text(fmt.Sprintf("count=%d", n))
}

func cartNoteView(_ struct{}) *Node {
	cartRendersNote++
	n := UseStore(cartStore, func(s cartState) string { return s.Note })
	return Text(fmt.Sprintf("note=%s", n))
}

func TestUseStoreRerendersSelectedOnly(t *testing.T) {
	cartStore.Set(cartState{Count: 1, Note: "a"}) // 复位（同值时 DeepEqual 直接跳过）
	cartRendersCount, cartRendersNote = 0, 0

	g := &game{w: 400, h: 400}
	activeGame = g
	g.rootFiber = reconcile(nil, nil, Use(func(_ struct{}) *Node {
		return Div(Style(Column), Use(cartCountView, struct{}{}), Use(cartNoteView, struct{}{}))
	}, struct{}{}))
	layoutAll(g)

	if cartRendersCount != 1 || cartRendersNote != 1 {
		t.Fatalf("initial renders: count=%d note=%d, want 1/1", cartRendersCount, cartRendersNote)
	}

	// 只改 Note：仅 note 组件入脏并重渲染。
	cartStore.Set(cartState{Count: 1, Note: "b"})
	if len(g.dirty) != 1 {
		t.Fatalf("dirty queue = %d, want 1 (only note view)", len(g.dirty))
	}
	g.flushDirty()
	layoutAll(g)
	if cartRendersCount != 1 || cartRendersNote != 2 {
		t.Fatalf("after note change: count=%d note=%d, want 1/2", cartRendersCount, cartRendersNote)
	}
	var texts []string
	collectTexts(g.rootRN, &texts)
	if !contains(texts, "note=b") || !contains(texts, "count=1") {
		t.Fatalf("texts = %v", texts)
	}

	// 只改 Count：count 组件读到新值。
	cartStore.Set(cartState{Count: 7, Note: "b"})
	if len(g.dirty) != 1 {
		t.Fatalf("dirty queue = %d, want 1 (only count view)", len(g.dirty))
	}
	g.flushDirty()
	layoutAll(g)
	if cartRendersCount != 2 || cartRendersNote != 2 {
		t.Fatalf("after count change: count=%d note=%d, want 2/2", cartRendersCount, cartRendersNote)
	}
	texts = nil
	collectTexts(g.rootRN, &texts)
	if !contains(texts, "count=7") {
		t.Fatalf("texts = %v, want count=7", texts)
	}

	// 同值 Set：全不触发。
	g.flushDirty()
	cartStore.Set(cartState{Count: 7, Note: "b"})
	if len(g.dirty) != 0 {
		t.Fatalf("same-value Set queued %d dirty, want 0", len(g.dirty))
	}
}

func TestUseStoreUnsubscribesOnUnmount(t *testing.T) {
	cartStore.Set(cartState{Count: 1, Note: "a"})
	before := len(cartStore.subs)

	g := &game{w: 400, h: 400}
	activeGame = g
	g.rootFiber = reconcile(nil, nil, Use(cartCountView, struct{}{}))
	if got := len(cartStore.subs); got != before+1 {
		t.Fatalf("subs after mount = %d, want %d", got, before+1)
	}

	reconcile(nil, g.rootFiber, nil) // 整体卸载
	if got := len(cartStore.subs); got != before {
		t.Fatalf("subs after unmount = %d, want %d（退订泄漏）", got, before)
	}

	g.dirty = nil
	cartStore.Set(cartState{Count: 99, Note: "a"}) // 已卸载的订阅不应再标脏
	if len(g.dirty) != 0 {
		t.Fatalf("set after unmount queued %d dirty, want 0", len(g.dirty))
	}
}

// ---------- Signal ----------

var (
	testSig        = NewSignal(0)
	sigBindRenders int
)

func sigBindView(_ struct{}) *Node {
	sigBindRenders++
	return Div(Style(Column),
		Text("static"),
		Text("", Bind(testSig, func(v int) string { return fmt.Sprintf("v=%d", v) })))
}

func TestSignalBindDirectWrite(t *testing.T) {
	testSig.Set(0) // 复位
	sigBindRenders = 0

	g := &game{w: 400, h: 400}
	activeGame = g
	g.rootFiber = reconcile(nil, nil, Use(sigBindView, struct{}{}))
	layoutAll(g)

	if sigBindRenders != 1 {
		t.Fatalf("renders = %d, want 1", sigBindRenders)
	}
	var texts []string
	collectTexts(g.rootRN, &texts)
	if !contains(texts, "v=0") {
		t.Fatalf("initial texts = %v, want v=0（挂载时应回填当前值）", texts)
	}

	// 信号变化：直写 renderNode——组件不重渲染、不入脏队列，只请求重排。
	g.needsLayout = false
	testSig.Set(42)
	if sigBindRenders != 1 {
		t.Fatalf("bind 触发了组件重渲染（renders=%d），应直写不走 reconcile", sigBindRenders)
	}
	if len(g.dirty) != 0 {
		t.Fatalf("bind 入了 %d 个脏 fiber，应完全不经过 reconcile", len(g.dirty))
	}
	if !g.needsLayout {
		t.Fatal("bind 后应请求重排（needsLayout）")
	}
	texts = nil
	collectTexts(g.rootRN, &texts)
	if !contains(texts, "v=42") {
		t.Fatalf("texts = %v, want v=42", texts)
	}

	// 同值 bailout：不标脏。
	g.needsLayout = false
	testSig.Set(42)
	if g.needsLayout {
		t.Fatal("同值 Set 不应请求重排")
	}

	// 卸载后退订：Set 不再写死节点也不标脏。
	reconcile(nil, g.rootFiber, nil)
	g.needsLayout = false
	testSig.Set(43)
	if g.needsLayout {
		t.Fatal("卸载后 Set 不应再请求重排（退订泄漏）")
	}
}

// ---------- BindInput ----------

var testNameSig = NewSignal("")

func findInputRN(rn *renderNode) *renderNode {
	if rn.kind == rnInput {
		return rn
	}
	for _, c := range rn.children {
		if r := findInputRN(c); r != nil {
			return r
		}
	}
	return nil
}

func TestBindInputTwoWay(t *testing.T) {
	testNameSig.Set("") // 复位

	g := &game{w: 400, h: 400}
	activeGame = g
	g.rootFiber = reconcile(nil, nil, Input(Style(Width(200)), BindInput(testNameSig)))
	layoutAll(g)

	rn := findInputRN(g.rootRN)
	if rn == nil {
		t.Fatal("no input render node")
	}
	if rn.value != "" {
		t.Fatalf("initial value = %q, want empty", rn.value)
	}

	// 外部写信号：直写输入框（不走组件重渲染）。
	testNameSig.Set("hello")
	if rn.value != "hello" {
		t.Fatalf("外部 Set 后 rn.value = %q, want hello", rn.value)
	}
	if len(g.dirty) != 0 {
		t.Fatalf("外部 Set 入了 %d 个脏 fiber，应直写", len(g.dirty))
	}

	// 打字路径：引擎改 rn.value 后回调 onChange（镜像 harness.SetValue），信号被回写。
	rn.value = "hello world"
	if rn.onChange != nil {
		rn.onChange(rn.value)
	}
	if got := testNameSig.Get(); got != "hello world" {
		t.Fatalf("signal after typing = %q, want hello world", got)
	}

	// 回写不自激：打字后信号回调发现 rn.value 已是最新，bailout（上面若循环早已栈溢出）。

	// 卸载后退订。
	reconcile(nil, g.rootFiber, nil)
	g.needsLayout = false
	testNameSig.Set("bye")
	if g.needsLayout {
		t.Fatal("卸载后 Set 不应再写输入框（退订泄漏）")
	}
}
