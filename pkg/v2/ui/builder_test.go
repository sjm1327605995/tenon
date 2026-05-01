package ui

import (
	"fmt"
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
)

// ========== Builder Tests ==========

func TestBuilder(t *testing.T) {
	called := false
	widget := NewBuilder(func(ctx BuildContext) Widget {
		called = true
		if ctx == nil {
			t.Fatal("BuildContext should not be nil")
		}
		return testTextWidget{content: "builder"}
	})

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	if !called {
		t.Fatal("Builder func should be called during Mount")
	}

	ro := engine.rootElement.FindRenderObject()
	if ro == nil {
		t.Fatal("rootElement should have a RenderObject through child")
	}
}

func TestBuilderRebuild(t *testing.T) {
	buildCount := 0
	widget := NewBuilder(func(ctx BuildContext) Widget {
		buildCount++
		return testTextWidget{content: fmt.Sprintf("build-%d", buildCount)}
	})

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	if buildCount != 1 {
		t.Fatalf("Builder should build once, got %d", buildCount)
	}

	// Trigger global rebuild
	engine.Rebuild()
	engine.flushBuild()

	if buildCount != 2 {
		t.Fatalf("Builder should rebuild on Rebuild, got %d", buildCount)
	}
}

// ========== StatefulBuilder Tests ==========

func TestStatefulBuilder(t *testing.T) {
	widget := NewStatefulBuilder(func(ctx BuildContext, setState func(fn func())) Widget {
		count := 0
		return testTextWidget{content: fmt.Sprintf("count=%d", count)}
	})

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	se, ok := engine.rootElement.(*StatefulElement)
	if !ok {
		t.Fatalf("rootElement should be *StatefulElement, got %T", engine.rootElement)
	}

	if se.state == nil {
		t.Fatal("state should not be nil")
	}
}

func TestStatefulBuilderSetState(t *testing.T) {
	count := 0
	widget := NewStatefulBuilder(func(ctx BuildContext, setState func(fn func())) Widget {
		return testTextWidget{content: fmt.Sprintf("count=%d", count)}
	})

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	se := engine.rootElement.(*StatefulElement)
	state := se.state.(*statefulBuilderState)

	// Before setState
	ro := se.child.FindRenderObject()
	text := ro.(*render.RenderText)
	if text.Content != "count=0" {
		t.Fatalf("initial text should be 'count=0', got %s", text.Content)
	}

	// Trigger setState
	state.SetState(func() {
		count = 5
	})
	engine.flushBuild()

	// After setState: StatefulBuilder rebuilds with captured count
	ro = se.child.FindRenderObject()
	text = ro.(*render.RenderText)
	if text.Content != "count=5" {
		t.Fatalf("after setState text should be 'count=5', got %s", text.Content)
	}
}

func TestStatefulBuilderSetStateInline(t *testing.T) {
	widget := NewStatefulBuilder(func(ctx BuildContext, setState func(fn func())) Widget {
		// Access state via closure; use a static value to verify rebuild
		return NewStatefulBuilder(func(ctx2 BuildContext, setState2 func(fn func())) Widget {
			return testTextWidget{content: "nested"}
		})
	})

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	// Verify nested StatefulBuilder works
	se := engine.rootElement.(*StatefulElement)
	if se.child == nil {
		t.Fatal("outer StatefulBuilder should have a child")
	}

	innerSE, ok := se.child.(*StatefulElement)
	if !ok {
		t.Fatalf("inner should be *StatefulElement, got %T", se.child)
	}

	ro := innerSE.child.FindRenderObject()
	text := ro.(*render.RenderText)
	if text.Content != "nested" {
		t.Fatalf("nested text should be 'nested', got %s", text.Content)
	}
}
