package ui

import "github.com/sjm1327605995/tenon/pkg/v2/render"

// ==================== State 生命周期接口 ====================
// 将 State 的可选生命周期方法定义为命名接口，
// 避免 StatefulElement 中散落的 interface{} 断言。

// stateBinder 是 State 内部绑定接口。
type stateBinder interface {
	bind(element *StatefulElement, widget Widget)
}

// stateInitStater 实现 InitState 生命周期。
type stateInitStater interface {
	InitState()
}

// stateDidChangeDepender 实现 DidChangeDependencies 生命周期。
type stateDidChangeDepender interface {
	DidChangeDependencies()
}

// stateDidUpdater 实现 DidUpdateWidget 生命周期。
type stateDidUpdater interface {
	DidUpdateWidget(oldWidget Widget)
}

// stateDisposer 实现 Dispose 生命周期。
type stateDisposer interface {
	Dispose()
}

// StatefulWidget 是有内部状态的 Widget。
type StatefulWidget interface {
	Widget
	CreateState() State
}

// State 是有状态组件的状态持有者。
type State interface {
	Build(ctx BuildContext) Widget
	SetState(fn func())
	GetWidget() Widget
	GetContext() BuildContext
}

// BaseState 提供 State 的默认实现。
type BaseState struct {
	self    State
	element *StatefulElement
	widget  Widget
}

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

// NewStatefulElement 创建 StatefulElement。
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

	if bs, ok := s.state.(stateBinder); ok {
		bs.bind(s, widget)
	}

	s.buildContext = &elementBuildContext{element: s}

	if init, ok := s.state.(stateInitStater); ok {
		init.InitState()
	}

	if ddc, ok := s.state.(stateDidChangeDepender); ok {
		ddc.DidChangeDependencies()
	}

	s.child = UpdateChild(s, nil, s.state.Build(s.buildContext))
}

func (s *StatefulElement) Update(newWidget Widget) {
	oldWidget := s.GetWidget()
	s.ComponentElement.Update(newWidget)

	if bs, ok := s.state.(stateBinder); ok {
		bs.bind(s, newWidget)
	}

	if duw, ok := s.state.(stateDidUpdater); ok {
		duw.DidUpdateWidget(oldWidget)
	}

	s.child = UpdateChild(s, s.child, s.state.Build(s.buildContext))
}

func (s *StatefulElement) Unmount() {
	if d, ok := s.state.(stateDisposer); ok {
		d.Dispose()
	}
	if s.child != nil {
		s.child.Unmount()
		s.child = nil
	}
	if bs, ok := s.state.(stateBinder); ok {
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

func (s *StatefulElement) rebuild() {
	if s.state == nil {
		return
	}
	s.child = UpdateChild(s, s.child, s.state.Build(s.buildContext))
}

func (s *StatefulElement) markNeedsBuild() {
	if defaultEngine != nil {
		defaultEngine.scheduleBuildFor(s)
	}
}

func (s *StatefulElement) didChangeDependencies() {
	if ddc, ok := s.state.(stateDidChangeDepender); ok {
		ddc.DidChangeDependencies()
	}
	s.markNeedsBuild()
}
