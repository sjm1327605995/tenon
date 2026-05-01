package ui

import (
	"testing"
)

// ========== GlobalKey Tests ==========

func TestGlobalKeyRegistration(t *testing.T) {
	gk := NewGlobalKey()

	widget := testStatefulWidget{initial: 42}
	widget.SetKey(gk)

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	// Verify GlobalKey can find the element
	el := getGlobalKeyElement(gk)
	if el == nil {
		t.Fatal("GlobalKey should be registered after Mount")
	}

	if el.GetWidget() != widget {
		t.Fatal("GlobalKey should point to the correct element")
	}
}

func TestGlobalKeyCurrentWidget(t *testing.T) {
	gk := NewGlobalKey()

	widget := testStatefulWidget{initial: 10}
	widget.SetKey(gk)

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	w := gk.CurrentWidget()
	if w == nil {
		t.Fatal("CurrentWidget should not be nil")
	}

	if w != widget {
		t.Fatal("CurrentWidget should return the mounted widget")
	}
}

func TestGlobalKeyCurrentState(t *testing.T) {
	gk := NewGlobalKey()

	widget := testStatefulWidget{initial: 99}
	widget.SetKey(gk)

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	state := gk.CurrentState()
	if state == nil {
		t.Fatal("CurrentState should not be nil for StatefulWidget")
	}

	ts, ok := state.(*testState)
	if !ok {
		t.Fatalf("CurrentState should be *testState, got %T", state)
	}

	if ts.count != 99 {
		t.Fatalf("state count should be 99, got %d", ts.count)
	}
}

func TestGlobalKeyCurrentContext(t *testing.T) {
	gk := NewGlobalKey()

	widget := testStatefulWidget{initial: 1}
	widget.SetKey(gk)

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	ctx := gk.CurrentContext()
	if ctx == nil {
		t.Fatal("CurrentContext should not be nil")
	}

	if ctx.GetWidget() != widget {
		t.Fatal("CurrentContext should return the widget")
	}
}

func TestGlobalKeyUnregistration(t *testing.T) {
	gk := NewGlobalKey()

	widget := testStatefulWidget{initial: 0}
	widget.SetKey(gk)

	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	if getGlobalKeyElement(gk) == nil {
		t.Fatal("GlobalKey should be registered")
	}

	engine.rootElement.Unmount()

	if getGlobalKeyElement(gk) != nil {
		t.Fatal("GlobalKey should be unregistered after Unmount")
	}
}

func TestGlobalKeyUnique(t *testing.T) {
	k1 := NewGlobalKey()
	k2 := NewGlobalKey()

	if k1.Equals(k2) {
		t.Fatal("Two NewGlobalKey() should not be equal")
	}
}

// ========== ValueKey Tests ==========

func TestValueKey(t *testing.T) {
	k1 := NewValueKey("hello")
	k2 := NewValueKey("hello")
	k3 := NewValueKey("world")

	if !k1.Equals(k2) {
		t.Fatal("ValueKey with same value should be equal")
	}

	if k1.Equals(k3) {
		t.Fatal("ValueKey with different value should not be equal")
	}
}

// ========== CanUpdate with Key Tests ==========

func TestCanUpdateWithKey(t *testing.T) {
	w1 := testStatefulWidget{initial: 1}
	w1.SetKey(NewValueKey("a"))

	w2 := testStatefulWidget{initial: 2}
	w2.SetKey(NewValueKey("a"))

	w3 := testStatefulWidget{initial: 3}
	w3.SetKey(NewValueKey("b"))

	if !CanUpdate(w1, w2) {
		t.Fatal("Same type and same ValueKey should be updatable")
	}

	if CanUpdate(w1, w3) {
		t.Fatal("Same type but different ValueKey should not be updatable")
	}
}
