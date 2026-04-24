package components

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

func TestWindow_ChainAPI(t *testing.T) {
	w := NewWindow("Test", 400, 300)
	if w.SetTitle("New Title") != w {
		t.Error("SetTitle should return *Window")
	}
	if w.SetOnClose(func() {}) != w {
		t.Error("SetOnClose should return *Window")
	}
	if w.Show() != w {
		t.Error("Show should return *Window")
	}
}

func TestWindow_ElementType(t *testing.T) {
	w := NewWindow("Test", 400, 300)
	if w.ElementType() != "Window" {
		t.Fatalf("expected ElementType Window, got %s", w.ElementType())
	}
}

func TestWindow_ShowClose(t *testing.T) {
	w := NewWindow("Test", 400, 300)
	if !w.IsVisible() {
		t.Fatal("window should be visible by default")
	}
	w.Close()
	if w.IsVisible() {
		t.Fatal("window should be hidden after Close")
	}
}

func TestWindow_OnClose(t *testing.T) {
	called := false
	w := NewWindow("Test", 400, 300).SetOnClose(func() {
		called = true
	})
	w.Close()
	if !called {
		t.Fatal("onClose should be called")
	}
}

func TestWindow_DragMovesBounds(t *testing.T) {
	w := NewWindow("Test", 400, 300)
	w.SetBounds(core.LayoutBounds{X: 100, Y: 100, Width: 400, Height: 300})
	w.titleBar.SetBounds(core.LayoutBounds{X: 100, Y: 100, Width: 400, Height: 32})
	w.closeBtn.SetBounds(core.LayoutBounds{X: 480, Y: 108, Width: 16, Height: 16})
	w.Content().SetBounds(core.LayoutBounds{X: 112, Y: 144, Width: 376, Height: 244})

	// Simulate mouse down on title bar
	w.HandleEvent(&core.Event{Type: core.EventMouseDown, X: 150, Y: 110})
	if !w.dragging {
		t.Fatal("dragging should be true after MouseDown on title bar")
	}

	// Simulate drag: dx=30, dy=50
	w.dragOffsetX = 50
	w.dragOffsetY = 10
	w.moveTo(130, 150)

	b := w.GetBounds()
	if b.X != 130 || b.Y != 150 {
		t.Fatalf("expected window at (130,150), got (%.0f,%.0f)", b.X, b.Y)
	}

	// Content panel should have shifted by the same delta
	cpb := w.Content().GetBounds()
	if cpb.X != 142 || cpb.Y != 194 { // 112+30=142, 144+50=194
		t.Fatalf("expected content at (142,194), got (%.0f,%.0f)", cpb.X, cpb.Y)
	}

	// MouseUp ends drag
	w.HandleEvent(&core.Event{Type: core.EventMouseUp})
	if w.dragging {
		t.Fatal("dragging should be false after MouseUp")
	}
}

func TestWindow_CloseButtonClick(t *testing.T) {
	w := NewWindow("Test", 400, 300)
	w.SetBounds(core.LayoutBounds{X: 100, Y: 100, Width: 400, Height: 300})
	w.closeBtn.SetBounds(core.LayoutBounds{X: 480, Y: 108, Width: 16, Height: 16})

	w.HandleEvent(&core.Event{Type: core.EventClick, X: 488, Y: 116})

	if w.IsVisible() {
		t.Fatal("window should close when clicking close button")
	}
}

func TestWindow_ClickOnContentDoesNotClose(t *testing.T) {
	w := NewWindow("Test", 400, 300)
	w.SetBounds(core.LayoutBounds{X: 100, Y: 100, Width: 400, Height: 300})
	w.Content().SetBounds(core.LayoutBounds{X: 100, Y: 132, Width: 400, Height: 268})

	w.HandleEvent(&core.Event{Type: core.EventClick, X: 200, Y: 200})

	if !w.IsVisible() {
		t.Fatal("window should NOT close when clicking content area")
	}
}

func TestWindow_ContentPanel(t *testing.T) {
	w := NewWindow("Test", 400, 300)
	w.Content().AppendChild(NewText("hello"))
	children := w.Content().GetChildren()
	if len(children) != 1 {
		t.Fatalf("expected 1 child in content panel, got %d", len(children))
	}
}
