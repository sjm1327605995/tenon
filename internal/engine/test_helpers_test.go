package engine

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
)

// mockTextWidget 是 ui 包内的简单文本 Widget，用于测试 TestEnvironment。
type mockTextWidget struct {
	BaseWidget
	content string
}

func (m mockTextWidget) CreateElement() Element {
	return NewRenderObjectElement(m)
}

func (m mockTextWidget) CreateRenderObject(element Element) render.RenderObject {
	return render.NewRenderText(m.content)
}

func (m mockTextWidget) UpdateRenderObject(ro render.RenderObject, oldWidget Widget) {
	r := ro.(*render.RenderText)
	r.Content = m.content
}

func TestTestEnvironmentFindText(t *testing.T) {
	env := TestWidget(t, mockTextWidget{content: "Hello"}, 400, 300)

	rt := env.FindText("Hello")
	if rt == nil {
		t.Fatal("expected to find RenderText with 'Hello'")
	}
	if rt.Content != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", rt.Content)
	}
}

func TestTestEnvironmentFindTextNotFound(t *testing.T) {
	env := TestWidget(t, mockTextWidget{content: "Hello"}, 400, 300)

	rt := env.FindText("World")
	if rt != nil {
		t.Error("expected nil for non-existent text")
	}
}

func TestTestEnvironmentAssertBounds(t *testing.T) {
	env := TestWidget(t, mockTextWidget{content: "Test"}, 400, 300)
	if env.RootElement == nil {
		t.Fatal("nil root render object")
	}
	// 根 RenderObject 应该有非零尺寸
	b := env.RootElement.GetBounds()
	if b.Width <= 0 || b.Height <= 0 {
		t.Errorf("expected positive bounds, got %.1fx%.1f", b.Width, b.Height)
	}
}

func TestTestEnvironmentRebuild(t *testing.T) {
	var content string
	eng := NewEngine(func() Widget {
		return mockTextWidget{content: content}
	}, 400, 300)
	eng.Mount()

	env := &TestEnvironment{
		t:      t,
		Engine: eng,
		Root:   eng.GetRootElement(),
	}
	env.RootElement = eng.GetRootRenderObject()

	content = "First"
	env.Rebuild()
	rt := env.FindText("First")
	if rt == nil {
		t.Fatal("expected 'First' after rebuild")
	}

	content = "Second"
	env.Rebuild()
	rt = env.FindText("Second")
	if rt == nil {
		t.Fatal("expected 'Second' after rebuild")
	}
}

func TestTestEnvironmentWithNilWidget(t *testing.T) {
	env := TestWidget(t, nil, 400, 300)
	// nil widget should not panic
	if env.Root != nil {
		// 如果 root 不是 nil，FindText 应该安全返回 nil
		rt := env.FindText("anything")
		if rt != nil {
			t.Error("expected nil for nil widget tree")
		}
	}
}
