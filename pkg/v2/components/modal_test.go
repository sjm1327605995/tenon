package components

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

func TestModal_ChainAPI(t *testing.T) {
	m := NewModal()
	if m.SetTitle("Title") != m {
		t.Error("SetTitle should return *Modal")
	}
	if m.SetOnClose(func() {}) != m {
		t.Error("SetOnClose should return *Modal")
	}
	if m.SetCloseOnMask(false) != m {
		t.Error("SetCloseOnMask should return *Modal")
	}
	if m.SetCloseOnEsc(false) != m {
		t.Error("SetCloseOnEsc should return *Modal")
	}
	if m.SetMaskColor(nil) != m {
		t.Error("SetMaskColor should return *Modal")
	}
	if m.SetPanelSize(300, 200) != m {
		t.Error("SetPanelSize should return *Modal")
	}
	if m.Open() != m {
		t.Error("Open should return *Modal")
	}
}

func TestModal_OpenClose(t *testing.T) {
	m := NewModal()
	if m.IsVisible() {
		t.Fatal("modal should be hidden by default")
	}
	m.Open()
	if !m.IsVisible() {
		t.Fatal("modal should be visible after Open")
	}
	m.Close()
	if m.IsVisible() {
		t.Fatal("modal should be hidden after Close")
	}
}

func TestModal_OnClose(t *testing.T) {
	called := false
	m := NewModal().SetOnClose(func() {
		called = true
	})
	m.Open()
	m.Close()
	if !called {
		t.Fatal("onClose should be called")
	}
}

func TestModal_MaskClickCloses(t *testing.T) {
	m := NewModal()
	m.SetBounds(core.LayoutBounds{X: 0, Y: 0, Width: 800, Height: 600})
	m.Panel().SetBounds(core.LayoutBounds{X: 200, Y: 150, Width: 400, Height: 300})
	m.Open()

	// Simulate click outside panel (top-left corner of screen)
	event := &core.Event{Type: core.EventClick, X: 0, Y: 0}
	m.HandleEvent(event)

	if m.IsVisible() {
		t.Fatal("modal should close on mask click")
	}
}

func TestModal_MaskClickDoesNotCloseWhenDisabled(t *testing.T) {
	m := NewModal().SetCloseOnMask(false)
	m.SetBounds(core.LayoutBounds{X: 0, Y: 0, Width: 800, Height: 600})
	m.Panel().SetBounds(core.LayoutBounds{X: 200, Y: 150, Width: 400, Height: 300})
	m.Open()

	event := &core.Event{Type: core.EventClick, X: 0, Y: 0}
	m.HandleEvent(event)

	if !m.IsVisible() {
		t.Fatal("modal should NOT close when closeOnMask is false")
	}
}

func TestModal_PanelClickDoesNotClose(t *testing.T) {
	m := NewModal()
	m.SetBounds(core.LayoutBounds{X: 0, Y: 0, Width: 800, Height: 600})
	m.Panel().SetBounds(core.LayoutBounds{X: 200, Y: 150, Width: 400, Height: 300})
	m.Open()

	// Click at panel center
	pb := m.Panel().GetBounds()
	event := &core.Event{
		Type: core.EventClick,
		X:    pb.X + pb.Width/2,
		Y:    pb.Y + pb.Height/2,
	}
	m.HandleEvent(event)

	if !m.IsVisible() {
		t.Fatal("modal should NOT close when clicking inside panel")
	}
}

func TestModal_ElementType(t *testing.T) {
	m := NewModal()
	if m.ElementType() != "Modal" {
		t.Fatalf("expected ElementType Modal, got %s", m.ElementType())
	}
}
