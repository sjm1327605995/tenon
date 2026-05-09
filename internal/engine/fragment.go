package engine

import "github.com/sjm1327605995/tenon/internal/render"

// FragmentWidget 是一个透明容器，本身不产生 RenderObject。
// 用于将多个 Widget 组合在一起，在树中表现为一个节点，但渲染时完全透传子节点。
type FragmentWidget struct {
	BaseWidget
	Children []Widget
}

// Fragment 创建 FragmentWidget。
func Fragment(children ...Widget) FragmentWidget {
	return FragmentWidget{Children: children}
}

func (f FragmentWidget) CreateElement() Element {
	e := &FragmentElement{widget: f}
	e.Init(e, f)
	return e
}

// FragmentElement 是 FragmentWidget 对应的 Element。
// 它本身不产生 RenderObject，子节点的 RenderObject 直接挂载到祖父 RenderObject。
type FragmentElement struct {
	BaseElement
	children []Element
	widget   FragmentWidget
}

func (e *FragmentElement) Mount(parent Element, slot int) {
	e.BaseElement.Mount(parent, slot)
	e.children = UpdateChildren(e, nil, e.widget.Children)
}

func (e *FragmentElement) Update(newWidget Widget) {
	e.BaseElement.Update(newWidget)
	e.widget = newWidget.(FragmentWidget)
	e.children = UpdateChildren(e, e.children, e.widget.Children)
}

func (e *FragmentElement) Unmount() {
	for _, child := range e.children {
		child.Unmount()
	}
	e.children = nil
	e.BaseElement.Unmount()
}

func (e *FragmentElement) GetChildren() []Element {
	return e.children
}

func (e *FragmentElement) FindRenderObject() render.RenderObject {
	for _, child := range e.children {
		if ro := child.FindRenderObject(); ro != nil {
			return ro
		}
	}
	return nil
}
