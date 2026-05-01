package ui

import "github.com/sjm1327605995/tenon/pkg/v2/render"

// BuilderFunc 是 Builder widget 的构建函数签名。
type BuilderFunc func(ctx BuildContext) Widget

// NewBuilder 创建 Builder widget。
func NewBuilder(builder BuilderFunc) Builder {
	return Builder{Builder: builder}
}

// Builder 是一个简单的 widget，其 UI 由构建函数决定。
// 适用于需要在 BuildContext 上下文中构建子树的场景。
type Builder struct {
	BaseWidget
	Builder BuilderFunc
}

func (b Builder) CreateElement() Element {
	e := &statelessElement{}
	e.ComponentElement.BaseElement.Init(e, b)
	return e
}

// statelessElement 是 Builder 对应的 Element。
type statelessElement struct {
	ComponentElement
	child        Element
	buildContext *elementBuildContext
}

func (e *statelessElement) Mount(parent Element, slot int) {
	e.ComponentElement.Mount(parent, slot)
	e.buildContext = &elementBuildContext{element: e}
	w := e.GetWidget().(Builder)
	e.child = UpdateChild(e, nil, w.Builder(e.buildContext))
}

func (e *statelessElement) Update(newWidget Widget) {
	oldWidget := e.widget
	e.BaseElement.Update(newWidget)
	e.performRebuild(oldWidget)
}

func (e *statelessElement) performRebuild(oldWidget Widget) {
	w := e.GetWidget().(Builder)
	e.child = UpdateChild(e, e.child, w.Builder(e.buildContext))
}

func (e *statelessElement) GetChildren() []Element {
	if e.child == nil {
		return nil
	}
	return []Element{e.child}
}

func (e *statelessElement) FindRenderObject() render.RenderObject {
	if e.child != nil {
		return e.child.FindRenderObject()
	}
	return nil
}

// StatefulBuilderFunc 是 StatefulBuilder 的构建函数签名。
// setState 参数等价于 State.SetState。
type StatefulBuilderFunc func(ctx BuildContext, setState func(fn func())) Widget

// NewStatefulBuilder 创建 StatefulBuilder widget。
func NewStatefulBuilder(builder StatefulBuilderFunc) StatefulBuilder {
	return StatefulBuilder{Builder: builder}
}

// StatefulBuilder 是有内部状态的内联 widget。
// 适用于简单交互场景，不需要单独定义 StatefulWidget + State 类型。
//
// 示例：
//
//	StatefulBuilder{
//		Builder: func(ctx BuildContext, setState func(func())) Widget {
//			return Button("Count: " + strconv.Itoa(count)).OnTap(func() {
//				setState(func() { count++ })
//			})
//		},
//	}
type StatefulBuilder struct {
	BaseWidget
	Builder StatefulBuilderFunc
}

func (s StatefulBuilder) CreateElement() Element {
	return NewStatefulElement(s)
}

func (s StatefulBuilder) CreateState() State {
	st := &statefulBuilderState{}
	st.Init(st)
	return st
}

// statefulBuilderState 是 StatefulBuilder 的内部 State 实现。
type statefulBuilderState struct {
	BaseState
}

func (s *statefulBuilderState) Build(ctx BuildContext) Widget {
	w := s.GetWidget().(StatefulBuilder)
	return w.Builder(ctx, s.SetState)
}
