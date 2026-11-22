package ui

import (
	"testing"
	"time"
)

func TestChildBitmapGrowsAndReusesWords(t *testing.T) {
	var bitmap childBitmap
	for _, index := range []uint32{0, 63, 64, 255, 1000} {
		bitmap.Set(index)
		if !bitmap.Has(index) {
			t.Fatalf("index %d is not set", index)
		}
	}
	if len(bitmap.words) != 16 {
		t.Fatalf("bitmap words=%d, want power-of-two capacity 16", len(bitmap.words))
	}
	bitmap.Delete(64)
	if bitmap.Has(64) {
		t.Fatal("index 64 was not deleted")
	}
	bitmap.Reset()
	if bitmap.Any() || len(bitmap.words) != 16 {
		t.Fatal("Reset did not retain cleared storage")
	}
}

func TestHoverBitmapRetainsCurrentState(t *testing.T) {
	enters := 0
	h := MountDefault(Div(Style(Width(200), Height(100)), OnHover(func(in bool) {
		if in {
			enters++
		}
	})))
	defer h.Close()
	old := input
	input = &harnessInput{x: 100, y: 50}
	defer func() { input = old }()
	h.g.updateHover()
	rn := h.g.rootRN
	input = &harnessInput{x: 101, y: 50}
	h.g.updateHover()
	if rn.hover.self != [2]bool{true, true} {
		t.Fatalf("hover state did not retain current and previous bits: %v", rn.hover.self)
	}
	if enters != 1 {
		t.Fatalf("enters=%d, want 1 after moving inside", enters)
	}
}

func TestHoverLeaveRestoresAnimatedCard(t *testing.T) {
	enters, leaves := 0, 0
	card := func(_ struct{}) *Node {
		hovered, setHovered := UseState(false)
		lift := UseTween(boolf(hovered), 110, EaseOutQuad)
		size := float32(60) * (1 + 0.24*lift)
		return posTestCard(100-size/2, 100-size/2, size,
			OnHover(func(value bool) {
				if value {
					enters++
				} else {
					leaves++
				}
				setHovered(value)
			}))
	}
	parent := func(_ struct{}) *Node {
		UseElapsed()
		return Div(Style(Fill), Div(Style(Fill), Use(card, struct{}{})))
	}
	h := Mount(Use(parent, struct{}{}), 300, 200)
	defer h.Close()
	var hoverNode *renderNode
	h.g.rootRN.forEachAtLevel(h.g.maxLevel, func(n *renderNode) {
		if n.onHover != nil {
			hoverNode = n
		}
	})
	h.MoveTo(100, 100)
	h.Step(110)
	h.MoveTo(100, 100)
	if got := hoverNode.bounds.W; got <= 60 {
		t.Fatalf("hover did not enlarge card: %v", got)
	}
	h.MoveTo(10, 10)
	if enters != 1 || leaves != 1 {
		t.Fatalf("hover events after leave: enters=%d leaves=%d", enters, leaves)
	}
	h.Step(110)
	if got := hoverNode.bounds.W; got != 60 {
		t.Fatalf("card size after leave=%v, want 60", got)
	}
}

func TestBackendRuntimeHoverLeaveDuringTween(t *testing.T) {
	enters, leaves := 0, 0
	card := func(_ struct{}) *Node {
		hovered, setHovered := UseState(false)
		lift := UseTween(boolf(hovered), 110, EaseOutQuad)
		size := float32(60) * (1 + 0.24*lift)
		return posTestCard(100-size/2, 100-size/2, size, OnHover(func(value bool) {
			if value {
				enters++
			} else {
				leaves++
			}
			setHovered(value)
		}))
	}
	parent := func(_ struct{}) *Node {
		return Div(Style(Fill), Div(Style(Fill), Use(card, struct{}{})))
	}
	source := &harnessInput{x: 10, y: 10}
	runtime := NewBackendRuntime(Use(parent, struct{}{}), windowConfig{w: 300, h: 200}, source)
	defer runtime.Close()
	runtime.Frame(16)
	var hoverNode *renderNode
	runtime.g.rootRN.forEachAtLevel(runtime.g.maxLevel, func(n *renderNode) {
		if n.onHover != nil {
			hoverNode = n
		}
	})

	source.x, source.y = 100, 100
	runtime.Frame(16)
	runtime.Frame(16)
	if enters != 1 || hoverNode.bounds.W <= 60 {
		t.Fatalf("runtime enter failed: enters=%d size=%v", enters, hoverNode.bounds.W)
	}
	if !runtime.PointerInteractive(100, 100) {
		t.Fatal("runtime did not identify the hovered control")
	}
	if runtime.PointerInteractive(10, 10) {
		t.Fatal("plain game viewport was captured by a non-interactive container")
	}
	source.x, source.y = 10, 10
	runtime.Frame(16)
	for range 8 {
		runtime.Frame(16)
	}
	if leaves != 1 {
		t.Fatalf("runtime leave events=%d, want 1", leaves)
	}
	if got := hoverNode.bounds.W; got != 60 {
		t.Fatalf("runtime card size after leave=%v, want 60", got)
	}
}

func TestBackendRuntimeCloseRestoresInputAndIsIdempotent(t *testing.T) {
	previous := input
	source := &harnessInput{}
	runtime := NewBackendRuntime(Text("runtime"), windowConfig{w: 100, h: 100}, source)
	if input != source {
		t.Fatal("runtime input was not installed")
	}
	runtime.Close()
	if input != previous {
		t.Fatal("runtime close did not restore the previous input")
	}
	runtime.Close()
	if runtime.Frame(16) || runtime.WantsNextFrame(false) || runtime.PointerInteractive(0, 0) || runtime.PaintVersion(time.Now()) != 0 {
		t.Fatal("closed runtime still reports activity")
	}
	if w, h := runtime.Size(); w != 0 || h != 0 {
		t.Fatalf("closed runtime size = %dx%d", w, h)
	}
}

func TestBackendRuntimeRejectsTypedNilInput(t *testing.T) {
	var source *harnessInput
	defer func() {
		if recover() == nil {
			t.Fatal("typed-nil input did not panic")
		}
	}()
	NewBackendRuntime(Text("runtime"), windowConfig{}, source)
}

func posTestCard(x, y, size float32, attrs ...*Node) *Node {
	attrs = append(attrs, Style(Absolute, Left(x), Top(y), Width(size), Height(size)))
	return Div(attrs...)
}

func TestRenderNodeLevelsAndLevelLookup(t *testing.T) {
	h := MountDefault(Div(Style(Width(100), Height(100)), Div(Text("leaf"))))
	defer h.Close()
	root := h.g.rootRN
	child := root.children[0]
	grandchild := child.children[0]
	if root.level != 0 || child.level != 1 || grandchild.level != 2 || h.g.maxLevel != 2 {
		t.Fatalf("unexpected levels: root=%d child=%d grandchild=%d max=%d",
			root.level, child.level, grandchild.level, h.g.maxLevel)
	}
	var atLevelOne []*renderNode
	root.forEachAtLevel(1, func(n *renderNode) { atLevelOne = append(atLevelOne, n) })
	if len(atLevelOne) != 1 || atLevelOne[0] != child {
		t.Fatal("level lookup did not return the direct child")
	}
}
