package engine

import (
	"fmt"
	"testing"

	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/pkg/v2/render"
)

func init() {
	fonts.InitDefaultFont()
	SetTheme(DefaultLightTheme())
}

// ========== Test Widgets ==========

// testTextWidget is a minimal text widget for testing (avoids import cycle with widgets package).
type testTextWidget struct {
	BaseWidget
	content string
}

func (t testTextWidget) CreateElement() Element {
	return NewRenderObjectElement(t)
}

func (t testTextWidget) CreateRenderObject(element Element) render.RenderObject {
	r := render.NewRenderText(t.content)
	r.SetFontSize(14)
	r.SetColor(render.NewColor(0, 0, 0, 255))
	return r
}

func (t testTextWidget) UpdateRenderObject(ro render.RenderObject, oldWidget Widget) {
	r := ro.(*render.RenderText)
	if old, ok := oldWidget.(testTextWidget); !ok || old.content != t.content {
		r.SetContent(t.content)
	}
}

type testStatefulWidget struct {
	BaseWidget
	initial int
}

func (t testStatefulWidget) CreateElement() Element {
	return NewStatefulElement(t)
}

func (t testStatefulWidget) CreateState() State {
	s := &testState{}
	s.Init(s)
	return s
}

type testState struct {
	BaseState
	count int
}

func (s *testState) InitState() {
	s.count = s.GetWidget().(testStatefulWidget).initial
}

func (s *testState) Build(ctx BuildContext) Widget {
	return testTextWidget{content: fmt.Sprintf("count=%d", s.count)}
}

// ========== Tests ==========

func TestStatefulElementMount(t *testing.T) {
	widget := testStatefulWidget{initial: 42}
	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	if engine.rootElement == nil {
		t.Fatal("rootElement should not be nil")
	}

	se, ok := engine.rootElement.(*StatefulElement)
	if !ok {
		t.Fatalf("rootElement should be *StatefulElement, got %T", engine.rootElement)
	}

	if se.state == nil {
		t.Fatal("state should not be nil after Mount")
	}

	state := se.state.(*testState)
	if state.count != 42 {
		t.Fatalf("InitState should set count to 42, got %d", state.count)
	}

	if se.child == nil {
		t.Fatal("child should not be nil after Mount")
	}
}

func TestStatefulElementSetState(t *testing.T) {
	widget := testStatefulWidget{initial: 0}
	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	se := engine.rootElement.(*StatefulElement)
	state := se.state.(*testState)

	// Before setState
	ro := se.child.FindRenderObject()
	if ro == nil {
		t.Fatal("child RenderObject should not be nil")
	}
	text, ok := ro.(*render.RenderText)
	if !ok {
		t.Fatalf("child should be RenderText, got %T", ro)
	}
	if text.Content != "count=0" {
		t.Fatalf("initial text should be 'count=0', got %s", text.Content)
	}

	// Trigger setState
	state.SetState(func() {
		state.count = 5
	})

	// flushBuild should process dirtyElements
	engine.flushBuild()

	// After setState
	ro = se.child.FindRenderObject()
	text = ro.(*render.RenderText)
	if text.Content != "count=5" {
		t.Fatalf("after setState text should be 'count=5', got %s", text.Content)
	}
}

func TestStatefulElementUnmount(t *testing.T) {
	widget := testStatefulWidget{initial: 0}
	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	se := engine.rootElement.(*StatefulElement)
	state := se.state.(*testState)

	// We can't easily test Dispose without adding a callback,
	// but we can verify Unmount doesn't panic and child is cleaned up
	engine.rootElement.Unmount()

	if se.child != nil {
		t.Fatal("child should be nil after Unmount")
	}

	// Accessing state after Unmount is allowed, but element is nil
	if state.element != nil {
		t.Fatal("state.element should be nil after Unmount")
	}
}

func TestStatefulElementDidUpdateWidget(t *testing.T) {
	widget := testStatefulWidget{initial: 10}
	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	se := engine.rootElement.(*StatefulElement)
	state := se.state.(*testState)

	if state.count != 10 {
		t.Fatalf("initial count should be 10, got %d", state.count)
	}

	// Update with new widget (same type, different initial)
	newWidget := testStatefulWidget{initial: 20}
	engine.rootElement.Update(newWidget)
	engine.flushBuild()

	// State should persist (count stays 10), widget should update
	if state.count != 10 {
		t.Fatalf("state should persist across updates, count should be 10, got %d", state.count)
	}

	if state.GetWidget().(testStatefulWidget).initial != 20 {
		t.Fatal("widget should be updated to new widget")
	}
}

func TestBuildContextFindAncestor(t *testing.T) {
	widget := testStatefulWidget{initial: 0}
	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	se := engine.rootElement.(*StatefulElement)
	ctx := se.buildContext

	if ctx == nil {
		t.Fatal("buildContext should not be nil")
	}

	if ctx.GetWidget() == nil {
		t.Fatal("GetWidget should return the StatefulWidget")
	}
}
