package ui

import (
	"fmt"
	"testing"
	"time"
)

func TestUIRuntimeLifecycleConformance(t *testing.T) {
	tests := []struct {
		name string
		new  func(*testing.T, *Node) (drive func(), close func(), text func() string, paint func() []PaintOp)
	}{
		{
			name: "harness",
			new: func(t *testing.T, root *Node) (func(), func(), func() string, func() []PaintOp) {
				h := MountDefault(root)
				return func() { h.Step(0) }, h.Close, func() string { return h.Root().AllText() }, h.Paint
			},
		},
		{
			name: "backend",
			new: func(t *testing.T, root *Node) (func(), func(), func() string, func() []PaintOp) {
				source := &harnessInput{}
				runtime := NewBackendRuntime(root, windowConfig{w: 800, h: 600}, source)
				runtime.Frame(0)
				return func() { runtime.Frame(0) }, runtime.Close,
					func() string { return (&Harness{g: runtime.g}).Root().AllText() },
					func() []PaintOp {
						painter := &recordPainter{}
						runtime.Paint(painter)
						return painter.ops
					}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var set func(int)
			mounts, cleanups := 0, 0
			root := View(func() *Node {
				value, setter := UseState(0)
				set = setter
				UseEffect(func() Cleanup {
					mounts++
					return func() { cleanups++ }
				})
				return Text(fmt.Sprint(value))
			})
			drive, close, text, paint := test.new(t, root)
			if mounts != 1 || text() != "0" {
				t.Fatalf("mount: mounts=%d text=%q, want 1/0", mounts, text())
			}
			set(1)
			drive()
			if text() != "1" {
				t.Fatalf("updated text=%q, want 1", text())
			}
			if len(paint()) == 0 {
				t.Fatal("mounted runtime produced no paint operations")
			}
			close()
			if cleanups != 1 || len(paint()) != 0 {
				t.Fatalf("close: cleanups=%d paint ops=%d, want 1/0", cleanups, len(paint()))
			}
			close()
			if cleanups != 1 {
				t.Fatalf("repeated close ran %d cleanups, want 1", cleanups)
			}
		})
	}
}

func TestBackendRuntimesCloseOutOfOrderRestoresExternalInput(t *testing.T) {
	originalInput, originalScale := input, uiScale
	external := &harnessInput{}
	input, uiScale = external, float32(1.25)
	t.Cleanup(func() { input, uiScale = originalInput, originalScale })

	a := NewBackendRuntime(Text("A"), windowConfig{w: 100, h: 100}, &harnessInput{})
	a.Resize(100, 100, 2)
	b := NewBackendRuntime(Text("B"), windowConfig{w: 100, h: 100}, &harnessInput{})
	a.Close()
	b.Close()

	if input != external || uiScale != 1.25 {
		t.Fatalf("out-of-order close restored input=%T scale=%v, want external/1.25", input, uiScale)
	}
}

func TestBackendRuntimeCloseIsTerminal(t *testing.T) {
	mounts, cleanups := 0, 0
	var set func(int)
	runtime := NewBackendRuntime(View(func() *Node {
		value, setter := UseState(0)
		set = setter
		UseEffect(func() Cleanup {
			mounts++
			return func() { cleanups++ }
		})
		return Text(fmt.Sprint(value))
	}), windowConfig{w: 100, h: 100}, &harnessInput{})
	runtime.Frame(0)
	runtime.Close()
	defer runtime.Close()

	set(1)
	runtime.Resize(200, 200, 2)
	if runtime.Frame(0) || runtime.WantsNextFrame(false) {
		t.Fatal("closed runtime requested another frame")
	}
	if w, h := runtime.Size(); w != 0 || h != 0 {
		t.Fatalf("closed runtime size = %dx%d, want 0x0", w, h)
	}
	if version := runtime.PaintVersion(time.Now()); version != 0 {
		t.Fatalf("closed runtime paint version = %d, want 0", version)
	}
	if runtime.PointerInteractive(10, 10) || len(runtime.AccessibilitySnapshot()) != 0 {
		t.Fatal("closed runtime still exposes interaction or semantics")
	}
	painter := &recordPainter{}
	runtime.Paint(painter)
	if len(painter.ops) != 0 || mounts != 1 || cleanups != 1 {
		t.Fatalf("closed runtime became active: paint=%d mounts=%d cleanups=%d", len(painter.ops), mounts, cleanups)
	}

	replacement := NewBackendRuntime(Text("replacement"), windowConfig{w: 100, h: 100}, &harnessInput{})
	defer replacement.Close()
	replacement.Frame(0)
	painter = &recordPainter{}
	replacement.Paint(painter)
	if len(painter.ops) == 0 {
		t.Fatal("new runtime could not mount after the old runtime closed")
	}
}
