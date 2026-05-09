package engine

import "github.com/sjm1327605995/tenon/internal/render"

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

// BaseStateOf 提供 State 的泛型默认实现，避免 GetWidget() 的类型断言。
type BaseStateOf[W Widget] struct {
	self    State
	element *StatefulElement
	widget  W
}

func (b *BaseStateOf[W]) Init(self State) {
	b.self = self
}

func (b *BaseStateOf[W]) bind(element *StatefulElement, widget Widget) {
	b.element = element
	if widget != nil {
		b.widget = widget.(W)
	}
}

func (b *BaseStateOf[W]) SetState(fn func()) {
	if fn != nil {
		fn()
	}
	if b.element != nil {
		b.element.markNeedsBuild()
	}
}

func (b *BaseStateOf[W]) GetWidget() Widget {
	return b.widget
}

// Widget 返回具体类型的 Widget，无需类型断言。
func (b *BaseStateOf[W]) Widget() W {
	return b.widget
}

func (b *BaseStateOf[W]) GetContext() BuildContext {
	if b.element == nil {
		return nil
	}
	return b.element.buildContext
}

// BaseState 是 BaseStateOf[Widget] 的别名，保持向后兼容。
type BaseState = BaseStateOf[Widget]

// NewStatefulElement 创建 StatefulElement。
func NewStatefulElement(widget StatefulWidget) *StatefulElement {
	e := &StatefulElement{widget: widget}
	e.ComponentElement.BaseElement.Init(e, widget)
	return e
}

// StatefulElement 管理 StatefulWidget 的生命周期和 State 实例。
type StatefulElement struct {
	ComponentElement
	state        State
	child        Element
	buildContext *elementBuildContext
	widget       StatefulWidget
}

func (s *StatefulElement) GetState() State {
	return s.state
}

func (s *StatefulElement) Mount(parent Element, slot int) {
	s.ComponentElement.Mount(parent, slot)

	s.state = s.widget.CreateState()

	if bs, ok := s.state.(stateBinder); ok {
		bs.bind(s, s.widget)
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
	oldWidget := s.widget
	s.ComponentElement.Update(newWidget)
	s.widget = newWidget.(StatefulWidget)

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
