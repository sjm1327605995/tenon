package ui

import (
	"fmt"
	"testing"
)

func TestUseStoreSwitchesSubscription(t *testing.T) {
	a, b := NewStore(10), NewStore(20)
	var switchStore func(bool)
	renders := 0
	h := MountDefault(View(func() *Node {
		useB, setUseB := UseState(false)
		switchStore = setUseB
		store := a
		if useB {
			store = b
		}
		value := UseStore(store, func(n int) int { return n })
		renders++
		return Text(fmt.Sprint(value))
	}))
	t.Cleanup(h.Close)

	switchStore(true)
	h.Step(0)
	if got := h.Root().AllText(); got != "20" {
		t.Fatalf("after switching store, got %q, want 20", got)
	}
	b.Set(21)
	h.Step(0)
	if got := h.Root().AllText(); got != "21" {
		t.Fatalf("after updating the new store, got %q, want 21", got)
	}
	before := renders
	a.Set(11)
	h.Step(0)
	if renders != before {
		t.Fatal("the old store still triggers renders")
	}
	switchStore(false)
	h.Step(0)
	if got := h.Root().AllText(); got != "11" {
		t.Fatalf("after switching back, got %q, want 11", got)
	}
	h.Close()
	if len(a.subs)+len(b.subs) != 0 {
		t.Fatal("store subscriptions survived unmount")
	}
}

func TestHarnessSchedulesStateAndEffectsOnOwningTree(t *testing.T) {
	var setA func(int)
	var effects []string
	view := func(name string, capture func(func(int))) *Node {
		return View(func() *Node {
			n, set := UseState(0)
			if capture != nil {
				capture(set)
			}
			text := fmt.Sprintf("%s:%d", name, n)
			UseEffect(func() Cleanup {
				effects = append(effects, text)
				return nil
			}, n)
			return Text(text)
		})
	}
	a := MountDefault(view("A", func(set func(int)) { setA = set }))
	t.Cleanup(a.Close)
	b := MountDefault(view("B", nil))
	t.Cleanup(b.Close)

	setA(1)
	b.Step(0)
	if got := a.Root().AllText(); got != "A:0" {
		t.Fatalf("stepping B updated A to %q", got)
	}
	a.Step(0)
	if got := a.Root().AllText(); got != "A:1" {
		t.Fatalf("stepping A got %q, want A:1", got)
	}
	if got := fmt.Sprint(effects); got != "[A:0 B:0 A:1]" {
		t.Fatalf("effects ran in wrong runtime: %s", got)
	}

	b.Close()
	setA(2)
	a.Step(0)
	if got := a.Root().AllText(); got != "A:2" {
		t.Fatalf("closing B stopped A: got %q", got)
	}
}

func TestDuplicateSiblingKeysDoNotLeakMountedChildren(t *testing.T) {
	mounts, cleanups := 0, 0
	child := func() *Node {
		UseEffect(func() Cleanup {
			mounts++
			return func() { cleanups++ }
		})
		return Text("child")
	}
	var setCount func(int)
	h := MountDefault(View(func() *Node {
		count, set := UseState(2)
		setCount = set
		var kids []*Node
		for range count {
			kids = append(kids, Keyed("duplicate", View(child)))
		}
		return Div(kids...)
	}))
	t.Cleanup(h.Close)
	if mounts != 2 {
		t.Fatalf("mounted %d children, want 2", mounts)
	}
	setCount(1)
	h.Step(0)
	if mounts != 2 || cleanups != 1 {
		t.Fatalf("after shrinking list: mounts=%d cleanups=%d, want 2/1", mounts, cleanups)
	}
	setCount(0)
	h.Step(0)
	if cleanups != 2 {
		t.Fatalf("after removing list: cleanups=%d, want 2", cleanups)
	}
}

func TestHooksPreserveNilInterfaceValues(t *testing.T) {
	t.Run("state", func(t *testing.T) {
		var set func(any)
		h := MountDefault(View(func() *Node {
			value, setter := UseState[any](nil)
			set = setter
			return Text(fmt.Sprint(value))
		}))
		t.Cleanup(h.Close)
		if got := h.Root().AllText(); got != "<nil>" {
			t.Fatalf("initial state = %s, want nil", got)
		}
		set("updated")
		h.Step(0)
		if got := h.Root().AllText(); got != "updated" {
			t.Fatalf("updated state = %s, want updated", got)
		}
		set(nil)
		h.Step(0)
		if got := h.Root().AllText(); got != "<nil>" {
			t.Fatalf("cleared state = %s, want nil", got)
		}
	})
	t.Run("reducer", func(t *testing.T) {
		var dispatch func(any)
		h := MountDefault(View(func() *Node {
			value, send := UseReducer(func(_ any, next any) any { return next }, nil)
			dispatch = send
			return Text(fmt.Sprint(value))
		}))
		t.Cleanup(h.Close)
		dispatch(5)
		h.Step(0)
		if got := h.Root().AllText(); got != "5" {
			t.Fatalf("reducer state = %s, want 5", got)
		}
		dispatch(nil)
		h.Step(0)
		if got := h.Root().AllText(); got != "<nil>" {
			t.Fatalf("reducer state = %s, want nil", got)
		}
	})
	t.Run("memo", func(t *testing.T) {
		var set func(int)
		computations := 0
		h := MountDefault(View(func() *Node {
			n, setter := UseState(0)
			set = setter
			value := UseMemo(func() any { computations++; return nil }, n/2)
			return Text(fmt.Sprint(value))
		}))
		t.Cleanup(h.Close)
		set(1)
		h.Step(0)
		if computations != 1 || h.Root().AllText() != "<nil>" {
			t.Fatalf("memo did not cache nil: computations=%d", computations)
		}
		set(2)
		h.Step(0)
		if computations != 2 {
			t.Fatalf("memo did not recompute after deps changed: computations=%d", computations)
		}
	})
}

func TestSharedStoreSchedulesEachOwningHarnessAndPreservesNil(t *testing.T) {
	store := NewStore[any](nil)
	view := func() *Node {
		return Text(fmt.Sprint(UseStore(store, func(v any) any { return v })))
	}
	a := MountDefault(View(view))
	t.Cleanup(a.Close)
	b := MountDefault(View(view))
	t.Cleanup(b.Close)

	store.Set("changed")
	a.Step(0)
	if gotA, gotB := a.Root().AllText(), b.Root().AllText(); gotA != "changed" || gotB != "<nil>" {
		t.Fatalf("after stepping A, got %q/%q, want changed/<nil>", gotA, gotB)
	}
	b.Step(0)
	if got := b.Root().AllText(); got != "changed" {
		t.Fatalf("after stepping B, got %q, want changed", got)
	}
	a.Close()
	store.Set(nil)
	b.Step(0)
	if got := b.Root().AllText(); got != "<nil>" || len(store.subs) != 1 {
		t.Fatalf("after closing A, B=%q subscriptions=%d, want <nil>/1", got, len(store.subs))
	}
}

func TestDiscardedRenderDoesNotRunPendingEffect(t *testing.T) {
	ran := false
	h := MountDefault(ErrorBoundary(func(_ any, _ func()) *Node {
		return Text("fallback")
	}, View(func() *Node {
		UseEffect(func() Cleanup {
			ran = true
			return nil
		})
		panic("failed render")
	})))
	t.Cleanup(h.Close)
	if got := h.Root().AllText(); got != "fallback" || ran {
		t.Fatalf("fallback=%q, discarded effect ran=%t", got, ran)
	}
}

func TestHarnessRetainsFocusViewportAndAnimationAfterMountingAnotherTree(t *testing.T) {
	a := Mount(View(func() *Node {
		return Div(Input(Placeholder("A")), Text(fmt.Sprintf("%.0f/%.1f", Viewport().W, UseElapsed())))
	}), 150, 100)
	t.Cleanup(a.Close)
	a.Root().ByPlaceholder("A").Focus()
	b := Mount(Text("B"), 300, 100)
	t.Cleanup(b.Close)

	a.Resize(180, 100)
	a.Step(100)
	if got := a.Root().AllText(); got != "180/0.1" {
		t.Fatalf("A's viewport and elapsed = %q, want 180/0.1", got)
	}
	if !isFocused(a.Root().ByPlaceholder("A").rn) || b.Focused().Exists() {
		t.Fatal("focus did not stay with A")
	}
}

func TestBackendRuntimesKeepInputAndScaleIsolated(t *testing.T) {
	previousInput, previousScale := input, uiScale
	t.Cleanup(func() { input, uiScale = previousInput, previousScale })
	newRuntime := func(source InputSource, width int) *BackendRuntime {
		return NewBackendRuntime(View(func() *Node {
			hovered, setHovered := UseState(false)
			return Div(
				Style(Width(100), Height(100)),
				OnHover(setHovered),
				Text(fmt.Sprintf("%.0f/%t", Viewport().W, hovered)),
			)
		}), windowConfig{w: width, h: 100}, source)
	}

	sourceA := &harnessInput{x: 10, y: 10}
	a := newRuntime(sourceA, 200)
	a.Resize(200, 100, 2)
	t.Cleanup(a.Close)
	sourceB := &harnessInput{x: 500, y: 500}
	b := newRuntime(sourceB, 300)
	b.Resize(300, 100, 1)
	t.Cleanup(b.Close)

	a.Frame(0)
	b.Frame(0)
	a.Frame(0)
	b.Frame(0)
	if got := (&Harness{g: a.g}).Root().AllText(); got != "100/true" {
		t.Fatalf("runtime A = %q, want 100/true", got)
	}
	if got := (&Harness{g: b.g}).Root().AllText(); got != "300/false" {
		t.Fatalf("runtime B = %q, want 300/false", got)
	}

	a.Resize(400, 100, 2)
	if input != sourceB || uiScale != 1 {
		t.Fatal("resizing inactive runtime A changed runtime B's globals")
	}
	a.Frame(0)
	a.Frame(0)
	if got := (&Harness{g: a.g}).Root().AllText(); got != "200/true" {
		t.Fatalf("resized runtime A = %q, want 200/true", got)
	}
	if input != sourceB || uiScale != 1 {
		t.Fatal("driving runtime A did not restore runtime B's input and scale")
	}
}
