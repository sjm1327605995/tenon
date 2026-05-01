package ui

import "github.com/sjm1327605995/tenon/pkg/v2/render"

// StatefulWidget 是有内部状态的 Widget。
// 对应 Flutter 的 StatefulWidget。
type StatefulWidget interface {
	Widget
	CreateState() State
}

// State 是有状态组件的状态持有者。
// 用户自定义 State 应嵌入 BaseState。
type State interface {
	// Build 构建该状态下的 UI。
	Build(ctx BuildContext) Widget

	// SetState 执行状态变更并标记组件需要 rebuild。
	SetState(fn func())

	// GetWidget 返回当前绑定的 Widget（可能已更新）。
	GetWidget() Widget

	// GetContext 返回 BuildContext。
	GetContext() BuildContext

	// 内部绑定用，用户不应直接调用
	bind(element *StatefulElement, widget Widget)
}

// BaseState 提供 State 的默认实现。
// 用户自定义 State 必须嵌入此结构体，并在 CreateState 中设置 self。
type BaseState struct {
	self    State
	element *StatefulElement
	widget  Widget
}

// Init 初始化 BaseState，应在 CreateState 中调用。
func (b *BaseState) Init(self State) {
	b.self = self
}

func (b *BaseState) bind(element *StatefulElement, widget Widget) {
	b.element = element
	b.widget = widget
}

func (b *BaseState) SetState(fn func()) {
	if fn != nil {
		fn()
	}
	if b.element != nil {
		b.element.markNeedsBuild()
	}
}

func (b *BaseState) GetWidget() Widget {
	return b.widget
}

func (b *BaseState) GetContext() BuildContext {
	if b.element == nil {
		return nil
	}
	return b.element.buildContext
}

// NewStatefulElement 创建 StatefulElement 并初始化。
func NewStatefulElement(widget StatefulWidget) *StatefulElement {
	e := &StatefulElement{}
	e.ComponentElement.BaseElement.Init(e, widget)
	return e
}

// StatefulElement 管理 StatefulWidget 的生命周期和 State 实例。
type StatefulElement struct {
	ComponentElement
	state        State
	child        Element
	buildContext *elementBuildContext
}

func (s *StatefulElement) GetState() State {
	return s.state
}

func (s *StatefulElement) Mount(parent Element, slot int) {
	s.ComponentElement.Mount(parent, slot)

	widget := s.GetWidget().(StatefulWidget)
	s.state = widget.CreateState()

	// 绑定 element 和 widget 到 state
	if bs, ok := s.state.(interface{ bind(*StatefulElement, Widget) }); ok {
		bs.bind(s, widget)
	}

	// 创建 BuildContext
	s.buildContext = &elementBuildContext{element: s}

	// 调用 InitState（如果实现了）
	if init, ok := s.state.(interface{ InitState() }); ok {
		init.InitState()
	}

	// 首次 build
	s.child = UpdateChild(s, nil, s.state.Build(s.buildContext))
}

func (s *StatefulElement) Update(newWidget Widget) {
	oldWidget := s.GetWidget()
	s.ComponentElement.Update(newWidget)

	// 更新 state 持有的 widget 引用
	if bs, ok := s.state.(interface{ bind(*StatefulElement, Widget) }); ok {
		bs.bind(s, newWidget)
	}

	// 调用 DidUpdateWidget（如果实现了）
	if duw, ok := s.state.(interface{ DidUpdateWidget(Widget) }); ok {
		duw.DidUpdateWidget(oldWidget)
	}

	// rebuild
	s.child = UpdateChild(s, s.child, s.state.Build(s.buildContext))
}

func (s *StatefulElement) Unmount() {
	// 调用 Dispose（如果实现了）
	if d, ok := s.state.(interface{ Dispose() }); ok {
		d.Dispose()
	}
	if s.child != nil {
		s.child.Unmount()
		s.child = nil
	}
	// 解除 state 对 element 的引用
	if bs, ok := s.state.(interface{ bind(*StatefulElement, Widget) }); ok {
		bs.bind(nil, nil)
	}
	s.ComponentElement.Unmount()
}

func (s *StatefulElement) GetChildren() []Element {
	if s.child == nil {
		return nil
	}
	return []Element{s.child}
}

func (s *StatefulElement) FindRenderObject() render.RenderObject {
	if s.child != nil {
		return s.child.FindRenderObject()
	}
	return nil
}

// rebuild 被 flushBuild 调用，重新执行 Build 并 diff 子树。
func (s *StatefulElement) rebuild() {
	if s.state == nil {
		return
	}
	s.child = UpdateChild(s, s.child, s.state.Build(s.buildContext))
}

// markNeedsBuild 将自身标记为 dirty，下一帧 flushBuild 时会重新调用 Build。
func (s *StatefulElement) markNeedsBuild() {
	// 通过全局 Engine 调度 rebuild
	if defaultEngine != nil {
		defaultEngine.scheduleBuildFor(s)
	}
}
