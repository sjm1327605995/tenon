package ui

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
)

// ========== Test InheritedWidget ==========

// testInheritedWidget 是一个测试用的 InheritedWidget，向下传递一个整数值。
type testInheritedWidget struct {
	BaseWidget
	value int
	child Widget
}

func (t testInheritedWidget) CreateElement() Element {
	return NewInheritedElement(t)
}

func (t testInheritedWidget) UpdateShouldNotify(oldWidget InheritedWidget) bool {
	return t.value != oldWidget.(testInheritedWidget).value
}

func (t testInheritedWidget) BuildChild(ctx BuildContext) Widget {
	return t.child
}

// testInheritedConsumer 消费 testInheritedWidget 的值。
type testInheritedConsumer struct {
	BaseWidget
}

func (c testInheritedConsumer) CreateElement() Element {
	return NewStatefulElement(c)
}

func (c testInheritedConsumer) CreateState() State {
	s := &testInheritedConsumerState{}
	s.Init(s)
	return s
}

type testInheritedConsumerState struct {
	BaseState
	value int
}

func (s *testInheritedConsumerState) Build(ctx BuildContext) Widget {
	if iw, ok := ctx.DependOnInheritedWidgetOfExactType(reflect.TypeOf(testInheritedWidget{})); ok {
		s.value = iw.(testInheritedWidget).value
	}
	return testTextWidget{content: "consumer"}
}

// ========== Tests ==========

func TestInheritedWidgetMount(t *testing.T) {
	inner := testTextWidget{content: "inner"}
	widget := testInheritedWidget{value: 42, child: inner}
	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	ie, ok := engine.rootElement.(*InheritedElement)
	if !ok {
		t.Fatalf("rootElement should be *InheritedElement, got %T", engine.rootElement)
	}

	if ie.Child == nil {
		t.Fatal("InheritedElement should have a child")
	}

	ro := ie.Child.FindRenderObject()
	text := ro.(*render.RenderText)
	if text.Content != "inner" {
		t.Fatalf("child should be 'inner', got %s", text.Content)
	}
}

func TestInheritedWidgetDependOn(t *testing.T) {
	consumer := testInheritedConsumer{}
	widget := testInheritedWidget{value: 100, child: consumer}
	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	ie := engine.rootElement.(*InheritedElement)
	if len(ie.dependents) != 1 {
		t.Fatalf("should have 1 dependent, got %d", len(ie.dependents))
	}

	// Verify the consumer state got the value
	consumerEl := ie.Child.(*StatefulElement)
	state := consumerEl.state.(*testInheritedConsumerState)
	if state.value != 100 {
		t.Fatalf("consumer should read value 100, got %d", state.value)
	}
}

func TestInheritedWidgetNotifyDependents(t *testing.T) {
	consumer := testInheritedConsumer{}
	widget := testInheritedWidget{value: 1, child: consumer}
	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	ie := engine.rootElement.(*InheritedElement)
	consumerEl := ie.Child.(*StatefulElement)
	state := consumerEl.state.(*testInheritedConsumerState)

	if state.value != 1 {
		t.Fatalf("initial value should be 1, got %d", state.value)
	}

	// Update with new value (should notify dependents)
	newWidget := testInheritedWidget{value: 99, child: consumer}
	engine.rootElement.Update(newWidget)
	engine.flushBuild()

	if state.value != 99 {
		t.Fatalf("after update value should be 99, got %d", state.value)
	}
}

func TestInheritedWidgetNoNotifyWhenEqual(t *testing.T) {
	// A deeper consumer that depends on the inherited widget
	deepNotifyCount := 0
	deepConsumer := NewStatefulBuilder(func(ctx BuildContext, setState func(fn func())) Widget {
		if _, ok := ctx.DependOnInheritedWidgetOfExactType(reflect.TypeOf(testInheritedWidget{})); ok {
			// depends on inherited widget
		}
		return testTextWidget{content: fmt.Sprintf("notify-%d", deepNotifyCount)}
	})

	// Wrap deepConsumer in a Builder so it's not the direct child of InheritedWidget
	child := NewBuilder(func(ctx BuildContext) Widget {
		return deepConsumer
	})

	widget := testInheritedWidget{value: 5, child: child}
	engine := NewEngine(func() Widget { return widget }, 400, 300)
	engine.Mount()

	// Get the InheritedElement and its dependent
	ie := engine.rootElement.(*InheritedElement)
	if len(ie.dependents) != 1 {
		t.Fatalf("should have 1 dependent, got %d", len(ie.dependents))
	}

	// Clear dirtyElements before test
	engine.dirtyElements = nil

	// Update with same value (should NOT add dependents to dirtyElements)
	newWidget := testInheritedWidget{value: 5, child: child}
	engine.rootElement.Update(newWidget)
	engine.flushBuild()

	if len(engine.dirtyElements) != 0 {
		t.Fatalf("dirtyElements should be empty when value unchanged, got %d", len(engine.dirtyElements))
	}

	// Update with different value (should notify dependents → dirtyElements)
	newWidget = testInheritedWidget{value: 10, child: child}
	engine.rootElement.Update(newWidget)

	if len(engine.dirtyElements) != 1 {
		t.Fatalf("dirtyElements should have 1 element after value change, got %d", len(engine.dirtyElements))
	}
}
