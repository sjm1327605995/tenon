package engine

import (
	"reflect"
	"testing"

	"github.com/sjm1327605995/tenon/internal/render"
)

// TestEnvironment 是 Widget 测试环境，无需 GUI 即可测试组件。
type TestEnvironment struct {
	t          *testing.T
	Engine     *Engine
	Root       Element
	RootElement render.RenderObject
}

// TestWidget 构建并布局一个 Widget，返回测试环境。
// width/height 指定测试画布尺寸。
func TestWidget(t *testing.T, widget Widget, width, height int) *TestEnvironment {
	t.Helper()
	eng := NewEngine(func() Widget {
		return widget
	}, width, height)
	eng.Mount()
	// 触发一次 Update 计算 Yoga 布局
	eng.Update()

	env := &TestEnvironment{
		t:      t,
		Engine: eng,
		Root:   eng.GetRootElement(),
	}
	env.RootElement = eng.GetRootRenderObject()
	return env
}

// Rebuild 触发全局 rebuild 并刷新布局。
func (env *TestEnvironment) Rebuild() {
	env.Engine.Rebuild()
	env.Engine.Update()
	env.RootElement = env.Engine.GetRootRenderObject()
}

// UpdateWidget 用新的 buildFunc 替换构建函数并 rebuild。
func (env *TestEnvironment) UpdateWidget(buildFunc BuildFunc) {
	env.Engine.buildFunc = buildFunc
	env.Engine.Rebuild()
	env.Engine.Update()
	env.RootElement = env.Engine.GetRootRenderObject()
}

// FindText 在 RenderObject 树中查找包含指定文本的 RenderText。
func (env *TestEnvironment) FindText(text string) *render.RenderText {
	env.t.Helper()
	return findRenderText(env.RootElement, text)
}

// FindAllText 返回所有包含指定文本的 RenderText。
func (env *TestEnvironment) FindAllText(text string) []*render.RenderText {
	var results []*render.RenderText
	findAllRenderText(env.RootElement, text, &results)
	return results
}

// TapAt 在指定坐标模拟点击。
func (env *TestEnvironment) TapAt(x, y float32) {
	env.t.Helper()
	if env.RootElement == nil {
		env.t.Fatal("TapAt: no root render object")
	}
	target := env.Engine.hitTest(env.RootElement, x, y)
	if target != nil && target.GetOnClick() != nil {
		target.GetOnClick()()
	}
}

// AssertBounds 断言 RenderObject 的尺寸。
func (env *TestEnvironment) AssertBounds(ro render.RenderObject, expectedWidth, expectedHeight float32) {
	env.t.Helper()
	if ro == nil {
		env.t.Fatal("AssertBounds: nil RenderObject")
	}
	b := ro.GetBounds()
	if abs(b.Width-expectedWidth) > 1 {
		env.t.Errorf("expected width %.1f, got %.1f", expectedWidth, b.Width)
	}
	if abs(b.Height-expectedHeight) > 1 {
		env.t.Errorf("expected height %.1f, got %.1f", expectedHeight, b.Height)
	}
}

// AssertPosition 断言 RenderObject 的位置。
func (env *TestEnvironment) AssertPosition(ro render.RenderObject, expectedX, expectedY float32) {
	env.t.Helper()
	if ro == nil {
		env.t.Fatal("AssertPosition: nil RenderObject")
	}
	b := ro.GetBounds()
	if abs(b.X-expectedX) > 1 {
		env.t.Errorf("expected X %.1f, got %.1f", expectedX, b.X)
	}
	if abs(b.Y-expectedY) > 1 {
		env.t.Errorf("expected Y %.1f, got %.1f", expectedY, b.Y)
	}
}

// AssertVisible 断言 RenderObject 可见。
func (env *TestEnvironment) AssertVisible(ro render.RenderObject) {
	env.t.Helper()
	if ro == nil {
		env.t.Fatal("AssertVisible: nil RenderObject")
	}
	if !ro.IsVisible() {
		env.t.Error("expected RenderObject to be visible")
	}
}

// AssertHidden 断言 RenderObject 不可见。
func (env *TestEnvironment) AssertHidden(ro render.RenderObject) {
	env.t.Helper()
	if ro == nil {
		env.t.Fatal("AssertHidden: nil RenderObject")
	}
	if ro.IsVisible() {
		env.t.Error("expected RenderObject to be hidden")
	}
}

// FindByType 在 Element 树中查找指定类型的 Element。
func (env *TestEnvironment) FindByType(target interface{}) Element {
	env.t.Helper()
	return findByType(env.Root, target)
}

// GetRootElement 返回根 Element。
func (env *TestEnvironment) GetRootElement() Element {
	return env.Root
}

// GetRootRenderObject 返回根 RenderObject。
func (env *TestEnvironment) GetRootRenderObject() render.RenderObject {
	return env.RootElement
}

// ---- 内部辅助 ----

func findRenderText(ro render.RenderObject, text string) *render.RenderText {
	if ro == nil {
		return nil
	}
	if rt, ok := ro.(*render.RenderText); ok {
		if rt.Content == text {
			return rt
		}
	}
	for _, child := range ro.GetChildren() {
		if result := findRenderText(child, text); result != nil {
			return result
		}
	}
	return nil
}

func findAllRenderText(ro render.RenderObject, text string, results *[]*render.RenderText) {
	if ro == nil {
		return
	}
	if rt, ok := ro.(*render.RenderText); ok {
		if rt.Content == text {
			*results = append(*results, rt)
		}
	}
	for _, child := range ro.GetChildren() {
		findAllRenderText(child, text, results)
	}
}

func findByType(el Element, target interface{}) Element {
	if el == nil {
		return nil
	}
	// 简单的类型匹配通过 reflect
	if reflect.TypeOf(el) == reflect.TypeOf(target) {
		return el
	}
	for _, child := range el.GetChildren() {
		if result := findByType(child, target); result != nil {
			return result
		}
	}
	return nil
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
