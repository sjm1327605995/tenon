package engine

import (
	"fmt"
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
)

// TestLocalRebuild verifies that SetState only rebuilds the specific StatefulWidget subtree.
func TestLocalRebuild(t *testing.T) {
	outerBuildCount := 0

	widget := NewStatefulBuilder(func(ctx BuildContext, setState func(fn func())) Widget {
		count := 0
		return NewBuilder(func(ctx BuildContext) Widget {
			outerBuildCount++
			return NewStatefulBuilder(func(ctx2 BuildContext, setState2 func(fn func())) Widget {
				return testTextWidget{content: fmt.Sprintf("inner-%d", count)}
			})
		})
	})

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	if outerBuildCount != 1 {
		t.Fatalf("outer should build once, got %d", outerBuildCount)
	}

	// The inner StatefulBuilder's state setState should only trigger inner rebuild
	innerSE := engine.rootElement.(*StatefulElement).child.(*statelessElement).child.(*StatefulElement)
	innerState := innerSE.state.(*statefulBuilderState)
	innerState.SetState(func() {
		// state change
	})
	engine.flushBuild()

	// outer Builder should NOT rebuild on inner SetState
	if outerBuildCount != 1 {
		t.Fatalf("outer should NOT rebuild on inner SetState, got %d", outerBuildCount)
	}
}

// TestLocalRebuildWithRenderObject verifies that SetState correctly updates RenderObject.
func TestLocalRebuildWithRenderObject(t *testing.T) {
	widget := NewStatefulBuilder(func(ctx BuildContext, setState func(fn func())) Widget {
		content := "before"
		return NewBuilder(func(ctx BuildContext) Widget {
			return testTextWidget{content: content}
		})
	})

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	// Verify initial text
	ro := engine.rootElement.FindRenderObject()
	text := ro.(*render.RenderText)
	if text.Content != "before" {
		t.Fatalf("initial text should be 'before', got %s", text.Content)
	}

	// Trigger local rebuild by creating a new StatefulBuilder with different state
	// In practice, users would call SetState. Here we simulate by updating the widget.
	newWidget := NewStatefulBuilder(func(ctx BuildContext, setState func(fn func())) Widget {
		content := "after"
		return NewBuilder(func(ctx BuildContext) Widget {
			return testTextWidget{content: content}
		})
	})
	engine.rootElement.Update(newWidget)
	engine.flushBuild()

	// Verify updated text
	ro = engine.rootElement.FindRenderObject()
	text = ro.(*render.RenderText)
	if text.Content != "after" {
		t.Fatalf("after update text should be 'after', got %s", text.Content)
	}
}

// TestSetStateTriggersFlushBuild verifies that scheduleBuildFor + flushBuild works end-to-end.
func TestSetStateTriggersFlushBuild(t *testing.T) {
	buildCount := 0
	widget := NewStatefulBuilder(func(ctx BuildContext, setState func(fn func())) Widget {
		buildCount++
		return testTextWidget{content: fmt.Sprintf("build-%d", buildCount)}
	})

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	if buildCount != 1 {
		t.Fatalf("initial build count should be 1, got %d", buildCount)
	}

	se := engine.rootElement.(*StatefulElement)
	se.state.SetState(func() {
		// state change
	})

	// Before flushBuild, dirtyElements should contain the element
	if len(engine.dirtyElements) != 1 {
		t.Fatalf("dirtyElements should have 1 element, got %d", len(engine.dirtyElements))
	}

	engine.flushBuild()

	if buildCount != 2 {
		t.Fatalf("after setState build count should be 2, got %d", buildCount)
	}
}
