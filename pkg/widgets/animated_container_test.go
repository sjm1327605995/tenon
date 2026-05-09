package widgets

import (
	"testing"
	"time"

	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

func TestAnimatedContainerInitialValues(t *testing.T) {
	eng := engine.NewEngine(func() engine.Widget {
		return NewAnimatedContainer().
			WithSize(100, 50).
			WithBackground(*render.NewColor(255, 0, 0, 255)).
			WithRadius(8)
	}, 400, 300)
	eng.Mount()

	statefulEl := eng.GetRootElement().(*engine.StatefulElement)
	state := statefulEl.GetState().(*animatedContainerState)

	if state.width != 100 {
		t.Errorf("expected width=100, got %v", state.width)
	}
	if state.height != 50 {
		t.Errorf("expected height=50, got %v", state.height)
	}
	if state.borderRadius != 8 {
		t.Errorf("expected borderRadius=8, got %v", state.borderRadius)
	}
}

func TestAnimatedContainerTransition(t *testing.T) {
	var container engine.Widget
	eng := engine.NewEngine(func() engine.Widget {
		return container
	}, 400, 300)
	eng.Mount()

	// Initial: 100x50 red
	container = NewAnimatedContainer().
		WithSize(100, 50).
		WithBackground(*render.NewColor(255, 0, 0, 255)).
		WithDuration(100 * time.Millisecond)
	eng.Rebuild()
	eng.Update()

	statefulEl := eng.GetRootElement().(*engine.StatefulElement)
	state := statefulEl.GetState().(*animatedContainerState)

	if state.width != 100 {
		t.Fatalf("initial width should be 100, got %v", state.width)
	}

	// Change to 200x100 blue
	container = NewAnimatedContainer().
		WithSize(200, 100).
		WithBackground(*render.NewColor(0, 0, 255, 255)).
		WithDuration(100 * time.Millisecond)
	eng.Rebuild()
	eng.Update()

	// After rebuild, state should have recorded start/target values and started animation
	if !state.animating {
		t.Fatal("expected animation to be started")
	}
	if state.startWidth != 100 {
		t.Errorf("expected startWidth=100, got %v", state.startWidth)
	}
	if state.targetWidth != 200 {
		t.Errorf("expected targetWidth=200, got %v", state.targetWidth)
	}

	// Simulate halfway through animation (progress = 0.5)
	state.ctrl.Value = 0.5
	state.onAnimationTick()

	midWidth := state.width
	if midWidth <= 100 || midWidth >= 200 {
		t.Errorf("expected mid width between 100 and 200, got %v", midWidth)
	}

	// Complete animation
	state.ctrl.Value = 1.0
	state.onAnimationTick()

	if state.width != 200 {
		t.Errorf("expected final width=200, got %v", state.width)
	}
	if state.height != 100 {
		t.Errorf("expected final height=100, got %v", state.height)
	}
}

func TestAnimatedContainerNoAnimationWhenUnchanged(t *testing.T) {
	var container engine.Widget
	eng := engine.NewEngine(func() engine.Widget {
		return container
	}, 400, 300)
	eng.Mount()

	container = NewAnimatedContainer().
		WithSize(100, 50).
		WithDuration(100 * time.Millisecond)
	eng.Rebuild()
	eng.Update()

	statefulEl := eng.GetRootElement().(*engine.StatefulElement)
	state := statefulEl.GetState().(*animatedContainerState)

	// Rebuild with same values
	container = NewAnimatedContainer().
		WithSize(100, 50).
		WithDuration(100 * time.Millisecond)
	eng.Rebuild()
	eng.Update()

	if state.animating {
		t.Error("expected no animation when properties unchanged")
	}
}

func TestLerpFloat32(t *testing.T) {
	if v := lerpFloat32(0, 100, 0.5); v != 50 {
		t.Errorf("expected 50, got %v", v)
	}
	if v := lerpFloat32(0, 100, 0); v != 0 {
		t.Errorf("expected 0, got %v", v)
	}
	if v := lerpFloat32(0, 100, 1); v != 100 {
		t.Errorf("expected 100, got %v", v)
	}
}

func TestLerpColor(t *testing.T) {
	a := render.NewColor(0, 0, 0, 255)
	b := render.NewColor(255, 255, 255, 255)
	mid := lerpColor(a, b, 0.5)
	if mid == nil {
		t.Fatal("expected non-nil color")
	}
	r, g, bb, _ := mid.RGBA()
	// RGBA returns 0-65535
	if r < 32000 || r > 33600 {
		t.Errorf("expected mid red ~32768, got %v", r)
	}
	if g < 32000 || g > 33600 {
		t.Errorf("expected mid green ~32768, got %v", g)
	}
	if bb < 32000 || bb > 33600 {
		t.Errorf("expected mid blue ~32768, got %v", bb)
	}
}
