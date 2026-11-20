package ui

import (
	"fmt"
	"testing"
)

func TestViewportReaderUpdatesOnResizeAndUnsubscribes(t *testing.T) {
	renders := 0
	reader := func(_ struct{}) *Node {
		renders++
		vp := Viewport()
		return Text(fmt.Sprintf("%.0fx%.0f", vp.W, vp.H))
	}
	var setVisible func(bool)
	h := Mount(Use(func(_ struct{}) *Node {
		visible, set := UseState(true)
		setVisible = set
		return Div(Style(Fill), If(visible, Use(reader, struct{}{})))
	}, struct{}{}), 320, 240)
	h.Resize(200, 160)
	if !h.Root().ByText("200x160").Exists() {
		t.Fatalf("viewport reader did not update: %v", h.Root().Texts())
	}
	before := renders
	h.Step(0)
	if renders != before {
		t.Fatal("unchanged viewport triggered a rerender")
	}
	setVisible(false)
	h.Step(0)
	if len(h.g.viewportReaders) != 0 {
		t.Fatal("unmounted viewport reader retained a subscription")
	}
	h.Resize(400, 300)
	if renders != before {
		t.Fatal("unmounted viewport reader rendered after resize")
	}
}
