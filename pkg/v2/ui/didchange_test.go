package ui

import (
	"reflect"
	"testing"
)

// testTheme 是测试用的 InheritedWidget。
type testTheme struct {
	BaseWidget
	Value  string
	child  Widget
}

func (t testTheme) CreateElement() Element {
	return NewInheritedElement(t)
}

func (t testTheme) UpdateShouldNotify(oldWidget InheritedWidget) bool {
	return t.Value != oldWidget.(testTheme).Value
}

func (t testTheme) BuildChild(ctx BuildContext) Widget {
	return t.child
}

// didChangeState 用于测试 didChangeDependencies 的调用。
type didChangeState struct {
	BaseState
	callCount  int
	lastValue  string
}

func (s *didChangeState) InitState() {
	// 不在这里读取 InheritedWidget
}

func (s *didChangeState) DidChangeDependencies() {
	s.callCount++
	ctx := s.GetContext()
	if ctx == nil {
		return
	}
	// 通过 DependOnInheritedWidgetOfExactType 获取 testTheme
	iw, ok := ctx.DependOnInheritedWidgetOfExactType(testThemeType)
	if ok && iw != nil {
		if tw, ok := iw.(testTheme); ok {
			s.lastValue = tw.Value
		}
	}
}

func (s *didChangeState) Build(ctx BuildContext) Widget {
	return nil
}

// didChangeWidget 是对应的 StatefulWidget。
type didChangeWidget struct {
	BaseWidget
}

func (w didChangeWidget) CreateElement() Element {
	return NewStatefulElement(w)
}

func (w didChangeWidget) CreateState() State {
	s := &didChangeState{}
	s.Init(s)
	return s
}

var testThemeType = reflect.TypeOf(testTheme{})

func TestDidChangeDependenciesCalledOnMount(t *testing.T) {
	eng := NewEngine(func() Widget {
		return testTheme{
			Value: "light",
			child: didChangeWidget{},
		}
	}, 800, 600)
	eng.Mount()

	// 找到 didChangeWidget 的 StatefulElement
	root := eng.GetRootElement()
	state := findDidChangeState(root)
	if state == nil {
		t.Fatal("didChangeState not found")
	}
	if state.callCount < 1 {
		t.Errorf("expected DidChangeDependencies called at least once on mount, got %d", state.callCount)
	}
	if state.lastValue != "light" {
		t.Errorf("expected lastValue 'light', got '%s'", state.lastValue)
	}
}

func TestDidChangeDependenciesCalledOnInheritedChange(t *testing.T) {
	var themeValue string
	eng := NewEngine(func() Widget {
		return testTheme{
			Value: themeValue,
			child: didChangeWidget{},
		}
	}, 800, 600)
	eng.Mount()

	root := eng.GetRootElement()
	state := findDidChangeState(root)
	if state == nil {
		t.Fatal("didChangeState not found")
	}
	initialCount := state.callCount

	// 改变 InheritedWidget 的值
	themeValue = "dark"
	eng.Rebuild()
	eng.Update()

	if state.callCount <= initialCount {
		t.Errorf("expected DidChangeDependencies called again after InheritedWidget change, count: %d -> %d", initialCount, state.callCount)
	}
	if state.lastValue != "dark" {
		t.Errorf("expected lastValue 'dark', got '%s'", state.lastValue)
	}
}

func findDidChangeState(el Element) *didChangeState {
	if se, ok := el.(*StatefulElement); ok {
		if ds, ok := se.GetState().(*didChangeState); ok {
			return ds
		}
	}
	for _, child := range el.GetChildren() {
		if s := findDidChangeState(child); s != nil {
			return s
		}
	}
	return nil
}
