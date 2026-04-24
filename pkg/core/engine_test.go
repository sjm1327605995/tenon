package core

import (
	"testing"

	"github.com/sjm1327605995/tenon/yoga"
)

// ==================== 测试用的简单 Host ====================

type testHost struct {
	BaseHost
	name string
}

func newTestHost(name string) *testHost {
	h := &testHost{name: name}
	h.Init(h)
	return h
}

// ==================== 测试用的 Widget ====================

type testWidget struct {
	BaseWidget
	child Component
}

func newTestWidget(child Component) *testWidget {
	w := &testWidget{child: child}
	w.Init(w)
	return w
}

func (w *testWidget) Render() Component {
	return w.child
}

// 可切换 child 的 Widget，用于测试更新
type toggleWidget struct {
	BaseWidget
	childA Component
	childB Component
	useA   bool
}

func newToggleWidget(a, b Component) *toggleWidget {
	w := &toggleWidget{childA: a, childB: b, useA: true}
	w.Init(w)
	return w
}

func (w *toggleWidget) Render() Component {
	if w.useA {
		return w.childA
	}
	return w.childB
}

// ==================== 测试 ====================

func TestEngineMountHost(t *testing.T) {
	root := newTestHost("root")
	root.GetElement().Yoga.StyleSetWidth(100)
	root.GetElement().Yoga.StyleSetHeight(100)

	child := newTestHost("child")
	root.AddChild(child)

	engine := NewEngine(root, 100, 100)
	engine.Mount()

	if engine.rootNode == nil {
		t.Fatal("rootNode should not be nil")
	}
	if len(engine.rootNode.children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(engine.rootNode.children))
	}
	if engine.rootNode.children[0].host != child {
		t.Fatal("child host mismatch")
	}
}

func TestEngineMountWidget(t *testing.T) {
	root := newTestHost("root")
	inner := newTestHost("inner")
	widget := newTestWidget(inner)
	root.AddChild(widget)

	engine := NewEngine(root, 100, 100)
	engine.Mount()

	// Widget 被展开后，root 的子节点应该直接是 inner
	if len(engine.rootNode.children) != 1 {
		t.Fatalf("expected 1 child after widget expansion, got %d", len(engine.rootNode.children))
	}
	if engine.rootNode.children[0].host != inner {
		t.Fatal("widget should expand to inner host")
	}
	// Widget 应该被记录为已挂载
	if _, ok := engine.widgetMounts[widget]; !ok {
		t.Fatal("widget should be in widgetMounts")
	}
}

func TestWidgetInvalidate(t *testing.T) {
	root := newTestHost("root")
	inner := newTestHost("inner")
	widget := newTestWidget(inner)
	root.AddChild(widget)

	engine := NewEngine(root, 100, 100)
	engine.Mount()

	widget.invalidate()
	if len(engine.dirtyWidgets) != 1 {
		t.Fatalf("expected 1 dirty widget, got %d", len(engine.dirtyWidgets))
	}

	// 执行更新
	engine.flushUpdates()
	if len(engine.dirtyWidgets) != 0 {
		t.Fatal("dirtyWidgets should be empty after flush")
	}
}

func TestWidgetUpdatePreservesParentLink(t *testing.T) {
	root := newTestHost("root")
	root.GetElement().Yoga.StyleSetWidth(300)
	root.GetElement().Yoga.StyleSetHeight(200)

	child1 := newTestHost("child1")
	child1.GetElement().Yoga.StyleSetWidth(100)
	child1.GetElement().Yoga.StyleSetHeight(50)

	child2 := newTestHost("child2")
	child2.GetElement().Yoga.StyleSetWidth(150)
	child2.GetElement().Yoga.StyleSetHeight(60)

	toggle := newToggleWidget(child1, child2)
	root.AddChild(toggle)

	engine := NewEngine(root, 300, 200)
	engine.Mount()

	// 初始状态：child1
	if len(engine.rootNode.children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(engine.rootNode.children))
	}
	node1 := engine.rootNode.children[0]
	if node1.host != child1 {
		t.Fatal("initial child should be child1")
	}

	// 切换并更新
	toggle.useA = false
	toggle.invalidate()
	engine.flushUpdates()
	engine.calculateLayout()

	// 更新后：root 下仍然只有 1 个子节点
	if len(engine.rootNode.children) != 1 {
		t.Fatalf("expected 1 child after update, got %d", len(engine.rootNode.children))
	}
	node2 := engine.rootNode.children[0]
	// 由于新旧 Host 类型相同，框架会自动复用旧 Host 并同步属性（类似 React 的 DOM 复用）
	if node2.host != child1 {
		t.Fatal("after toggle, framework should reuse child1 and sync properties from child2")
	}
	// Yoga 样式已被同步，布局结果应该和 child2 一致
	b1 := child1.GetLayoutBounds()
	if b1.Width != 150 {
		t.Fatalf("expected synced width 150, got %f", b1.Width)
	}
	if b1.Height != 60 {
		t.Fatalf("expected synced height 60, got %f", b1.Height)
	}
}

func TestYogaLayout(t *testing.T) {
	root := newTestHost("root")
	root.GetElement().Yoga.StyleSetWidth(800)
	root.GetElement().Yoga.StyleSetHeight(600)

	child := newTestHost("child")
	child.GetElement().Yoga.StyleSetWidth(100)
	child.GetElement().Yoga.StyleSetHeight(50)
	root.AddChild(child)

	engine := NewEngine(root, 800, 600)
	engine.Mount()

	bounds := child.GetLayoutBounds()
	if bounds.Width != 100 {
		t.Fatalf("expected width 100, got %f", bounds.Width)
	}
	if bounds.Height != 50 {
		t.Fatalf("expected height 50, got %f", bounds.Height)
	}
}

func TestWidgetRenderReturnsNil(t *testing.T) {
	root := newTestHost("root")
	widget := &testWidget{}
	widget.Init(widget)
	root.AddChild(widget)

	engine := NewEngine(root, 100, 100)
	engine.Mount()

	// Widget Render 返回 nil，不应产生子节点
	if len(engine.rootNode.children) != 0 {
		t.Fatalf("expected 0 children for nil-render widget, got %d", len(engine.rootNode.children))
	}
}

func TestNestedWidget(t *testing.T) {
	root := newTestHost("root")
	leaf := newTestHost("leaf")

	innerWidget := newTestWidget(leaf)
	outerWidget := newTestWidget(innerWidget)
	root.AddChild(outerWidget)

	engine := NewEngine(root, 100, 100)
	engine.Mount()

	// outerWidget -> innerWidget -> leaf
	// 最终 root 下应该只有 leaf
	if len(engine.rootNode.children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(engine.rootNode.children))
	}
	if engine.rootNode.children[0].host != leaf {
		t.Fatal("nested widget expansion failed")
	}
}

func TestHostAdd(t *testing.T) {
	root := newTestHost("root")
	c1 := newTestHost("c1")
	c2 := newTestHost("c2")

	root.AddChild(c1)
	root.AddChild(c2)

	engine := NewEngine(root, 100, 100)
	engine.Mount()

	if len(engine.rootNode.children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(engine.rootNode.children))
	}
}

func TestYogaFlexLayout(t *testing.T) {
	root := newTestHost("root")
	root.GetElement().Yoga.StyleSetWidth(300)
	root.GetElement().Yoga.StyleSetHeight(100)
	root.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)

	c1 := newTestHost("c1")
	c1.GetElement().Yoga.StyleSetFlexGrow(1)
	c1.GetElement().Yoga.StyleSetHeight(100)

	c2 := newTestHost("c2")
	c2.GetElement().Yoga.StyleSetFlexGrow(2)
	c2.GetElement().Yoga.StyleSetHeight(100)

	root.AddChild(c1)
	root.AddChild(c2)

	engine := NewEngine(root, 300, 100)
	engine.Mount()

	b1 := c1.GetLayoutBounds()
	b2 := c2.GetLayoutBounds()

	if b1.Width != 100 {
		t.Fatalf("expected c1 width 100, got %f", b1.Width)
	}
	if b2.Width != 200 {
		t.Fatalf("expected c2 width 200, got %f", b2.Width)
	}
}

// ==================== 焦点系统测试 ====================

func TestFocusSystem(t *testing.T) {
	root := newTestHost("root")
	root.GetElement().Yoga.StyleSetWidth(300)
	root.GetElement().Yoga.StyleSetHeight(200)

	btn1 := newTestHost("btn1")
	btn1.SetFocusable(true)
	btn1.GetElement().Yoga.StyleSetWidth(100)
	btn1.GetElement().Yoga.StyleSetHeight(50)

	btn2 := newTestHost("btn2")
	btn2.SetFocusable(true)
	btn2.GetElement().Yoga.StyleSetWidth(100)
	btn2.GetElement().Yoga.StyleSetHeight(50)

	root.AddChild(btn1)
	root.AddChild(btn2)

	engine := NewEngine(root, 300, 200)
	engine.Mount()

	// 初始无焦点
	if engine.focusHost != nil {
		t.Fatal("initial focus should be nil")
	}

	// 收集可焦点组件
	focusables := engine.collectFocusableHosts()
	if len(focusables) != 2 {
		t.Fatalf("expected 2 focusables, got %d", len(focusables))
	}

	// 设置焦点
	engine.setFocusHost(btn1)
	if engine.focusHost != btn1 {
		t.Fatal("focus should be btn1")
	}

	// 切换下一个
	engine.focusNext(false)
	if engine.focusHost != btn2 {
		t.Fatal("focus should move to btn2")
	}

	// 循环回到开头
	engine.focusNext(false)
	if engine.focusHost != btn1 {
		t.Fatal("focus should wrap to btn1")
	}

	// 反向切换
	engine.focusNext(true)
	if engine.focusHost != btn2 {
		t.Fatal("focus should move back to btn2")
	}
}

func TestFocusableHostInterface(t *testing.T) {
	h := newTestHost("h")
	if h.IsFocusable() {
		t.Fatal("default focusable should be false")
	}
	h.SetFocusable(true)
	if !h.IsFocusable() {
		t.Fatal("focusable should be true after SetFocusable(true)")
	}
}

func TestScrollOffset(t *testing.T) {
	root := newTestHost("root")
	root.GetElement().Yoga.StyleSetWidth(200)
	root.GetElement().Yoga.StyleSetHeight(200)

	child := newTestHost("child")
	child.GetElement().Yoga.StyleSetWidth(50)
	child.GetElement().Yoga.StyleSetHeight(50)
	root.AddChild(child)

	engine := NewEngine(root, 200, 200)
	engine.Mount()

	// 默认无偏移
	x, y := root.GetScrollOffset()
	if x != 0 || y != 0 {
		t.Fatalf("default scroll offset should be (0,0), got (%f,%f)", x, y)
	}

	// 布局后 child 应该在 (0,0)
	b := child.GetLayoutBounds()
	if b.X != 0 || b.Y != 0 {
		t.Fatalf("expected child at (0,0), got (%f,%f)", b.X, b.Y)
	}
}
