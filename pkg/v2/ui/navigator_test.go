package ui

import (
	"testing"
)

// getTestNavigator 从 StatefulElement 获取 NavigatorState。
func getTestNavigator(t *testing.T, eng *Engine) NavigatorState {
	t.Helper()
	se, ok := eng.GetRootElement().(*StatefulElement)
	if !ok {
		t.Fatal("expected root element to be StatefulElement")
	}
	nav, ok := se.GetState().(NavigatorState)
	if !ok {
		t.Fatal("expected state to implement NavigatorState")
	}
	return nav
}

func TestNavigatorInitialState(t *testing.T) {
	routes := map[string]RouteBuilder{
		"home": func(ctx BuildContext, params RouteParams) Widget {
			return nil
		},
	}
	eng := NewEngine(func() Widget {
		return Navigator(routes, "home")
	}, 800, 600)
	eng.Mount()

	nav := getTestNavigator(t, eng)
	if nav.CurrentPage() != "home" {
		t.Errorf("expected current page 'home', got '%s'", nav.CurrentPage())
	}
	if nav.PageCount() != 1 {
		t.Errorf("expected page count 1, got %d", nav.PageCount())
	}
}

func TestNavigatorPush(t *testing.T) {
	routes := map[string]RouteBuilder{
		"home": func(ctx BuildContext, params RouteParams) Widget {
			return nil
		},
		"settings": func(ctx BuildContext, params RouteParams) Widget {
			return nil
		},
	}
	eng := NewEngine(func() Widget {
		return Navigator(routes, "home")
	}, 800, 600)
	eng.Mount()

	nav := getTestNavigator(t, eng)
	nav.Push("settings")

	if nav.CurrentPage() != "settings" {
		t.Errorf("expected 'settings', got '%s'", nav.CurrentPage())
	}
	if nav.PageCount() != 2 {
		t.Errorf("expected 2, got %d", nav.PageCount())
	}
}

func TestNavigatorPop(t *testing.T) {
	routes := map[string]RouteBuilder{
		"home":     func(ctx BuildContext, p RouteParams) Widget { return nil },
		"settings": func(ctx BuildContext, p RouteParams) Widget { return nil },
	}
	eng := NewEngine(func() Widget {
		return Navigator(routes, "home")
	}, 800, 600)
	eng.Mount()

	nav := getTestNavigator(t, eng)
	nav.Push("settings")
	nav.Pop()

	if nav.CurrentPage() != "home" {
		t.Errorf("expected 'home', got '%s'", nav.CurrentPage())
	}
	if nav.PageCount() != 1 {
		t.Errorf("expected 1, got %d", nav.PageCount())
	}
}

func TestNavigatorPopToRoot(t *testing.T) {
	routes := map[string]RouteBuilder{
		"a": func(ctx BuildContext, p RouteParams) Widget { return nil },
		"b": func(ctx BuildContext, p RouteParams) Widget { return nil },
		"c": func(ctx BuildContext, p RouteParams) Widget { return nil },
	}
	eng := NewEngine(func() Widget {
		return Navigator(routes, "a")
	}, 800, 600)
	eng.Mount()

	nav := getTestNavigator(t, eng)
	nav.Push("b")
	nav.Push("c")
	if nav.PageCount() != 3 {
		t.Fatalf("expected 3, got %d", nav.PageCount())
	}

	nav.PopToRoot()
	if nav.CurrentPage() != "a" {
		t.Errorf("expected 'a', got '%s'", nav.CurrentPage())
	}
	if nav.PageCount() != 1 {
		t.Errorf("expected 1, got %d", nav.PageCount())
	}
}

func TestNavigatorPushReplacement(t *testing.T) {
	routes := map[string]RouteBuilder{
		"a": func(ctx BuildContext, p RouteParams) Widget { return nil },
		"b": func(ctx BuildContext, p RouteParams) Widget { return nil },
		"c": func(ctx BuildContext, p RouteParams) Widget { return nil },
	}
	eng := NewEngine(func() Widget {
		return Navigator(routes, "a")
	}, 800, 600)
	eng.Mount()

	nav := getTestNavigator(t, eng)
	nav.Push("b")
	nav.PushReplacement("c")

	if nav.CurrentPage() != "c" {
		t.Errorf("expected 'c', got '%s'", nav.CurrentPage())
	}
	if nav.PageCount() != 2 {
		t.Errorf("expected 2, got %d", nav.PageCount())
	}
}

func TestNavigatorPopCannotPopRoot(t *testing.T) {
	routes := map[string]RouteBuilder{
		"home": func(ctx BuildContext, p RouteParams) Widget { return nil },
	}
	eng := NewEngine(func() Widget {
		return Navigator(routes, "home")
	}, 800, 600)
	eng.Mount()

	nav := getTestNavigator(t, eng)
	nav.Pop() // 不应 panic
	if nav.CurrentPage() != "home" {
		t.Errorf("expected 'home', got '%s'", nav.CurrentPage())
	}
}

func TestNavigatorWithParams(t *testing.T) {
	routes := map[string]RouteBuilder{
		"detail": func(ctx BuildContext, params RouteParams) Widget { return nil },
	}
	eng := NewEngine(func() Widget {
		return NavigatorWithParams(routes, "detail", RouteParams{"id": 42})
	}, 800, 600)
	eng.Mount()

	nav := getTestNavigator(t, eng)
	if nav.CurrentPage() != "detail" {
		t.Errorf("expected 'detail', got '%s'", nav.CurrentPage())
	}
}

func TestNavigatorUnknownRoute(t *testing.T) {
	routes := map[string]RouteBuilder{
		"home": func(ctx BuildContext, p RouteParams) Widget { return nil },
	}
	eng := NewEngine(func() Widget {
		return Navigator(routes, "home")
	}, 800, 600)
	eng.Mount()

	nav := getTestNavigator(t, eng)
	nav.Push("nonexistent") // 不应 panic
	if nav.CurrentPage() != "home" {
		t.Errorf("expected 'home', got '%s'", nav.CurrentPage())
	}
	if nav.PageCount() != 1 {
		t.Errorf("expected 1, got %d", nav.PageCount())
	}
}

func TestRouteParamsGet(t *testing.T) {
	params := RouteParams{"name": "test", "count": 5}
	if params.Get("name") != "test" {
		t.Errorf("expected 'test', got '%s'", params.Get("name"))
	}
	if params.GetInt("count") != 5 {
		t.Errorf("expected 5, got %d", params.GetInt("count"))
	}
	if params.Get("missing") != "" {
		t.Errorf("expected empty, got '%s'", params.Get("missing"))
	}
}

// TestGetNavigatorFromChild 验证子页面通过 BuildContext 获取 Navigator。
func TestGetNavigatorFromChild(t *testing.T) {
	var capturedNav NavigatorState
	routes := map[string]RouteBuilder{
		"home": func(ctx BuildContext, p RouteParams) Widget {
			capturedNav = GetNavigator(ctx)
			return nil
		},
	}
	eng := NewEngine(func() Widget {
		return Navigator(routes, "home")
	}, 800, 600)
	eng.Mount()

	if capturedNav == nil {
		t.Fatal("expected child to capture navigator via GetNavigator(ctx)")
	}
	if capturedNav.CurrentPage() != "home" {
		t.Errorf("expected 'home', got '%s'", capturedNav.CurrentPage())
	}
}
