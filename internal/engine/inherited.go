package engine

import (
	"reflect"

	"github.com/sjm1327605995/tenon/internal/render"
)

// InheritedWidget 是一种特殊 Widget，用于在树中向下传递数据。
// 当数据变化时，自动 rebuild 所有依赖它的子节点。
type InheritedWidget interface {
	Widget
	// UpdateShouldNotify 判断新数据是否发生了变化。
	// 返回 true 时，框架会自动通知所有依赖该 Widget 的子节点 rebuild。
	UpdateShouldNotify(oldWidget InheritedWidget) bool
	// BuildChild 返回该 InheritedWidget 包裹的子 Widget。
	BuildChild(ctx BuildContext) Widget
}

// InheritedElement 管理 InheritedWidget 的依赖关系。
type InheritedElement struct {
	ComponentElement
	Child        Element
	dependents   map[Element]struct{}
	buildContext *elementBuildContext
	widget       InheritedWidget
}

func NewInheritedElement(widget Widget) *InheritedElement {
	e := &InheritedElement{
		dependents: make(map[Element]struct{}),
		widget:     widget.(InheritedWidget),
	}
	e.BaseElement.Init(e, widget)
	return e
}

func (i *InheritedElement) Mount(parent Element, slot int) {
	i.ComponentElement.Mount(parent, slot)
	i.buildContext = &elementBuildContext{element: i}
	i.Child = UpdateChild(i, nil, i.widget.BuildChild(i.buildContext))
}

func (i *InheritedElement) Update(newWidget Widget) {
	oldWidget := i.widget
	i.BaseElement.Update(newWidget)
	i.widget = newWidget.(InheritedWidget)

	if oldWidget.UpdateShouldNotify(i.widget) {
		i.notifyDependents()
	}

	// 更新子树
	i.Child = UpdateChild(i, i.Child, i.widget.BuildChild(i.buildContext))
}

func (i *InheritedElement) Unmount() {
	if i.Child != nil {
		i.Child.Unmount()
		i.Child = nil
	}
	i.ComponentElement.Unmount()
}

func (i *InheritedElement) GetChildren() []Element {
	if i.Child == nil {
		return nil
	}
	return []Element{i.Child}
}

func (i *InheritedElement) FindRenderObject() render.RenderObject {
	if i.Child != nil {
		return i.Child.FindRenderObject()
	}
	return nil
}

// addDependent 注册一个依赖该 InheritedWidget 的子 Element。
func (i *InheritedElement) addDependent(dependent Element) {
	if i.dependents == nil {
		i.dependents = make(map[Element]struct{})
	}
	i.dependents[dependent] = struct{}{}
}

// removeDependent 移除依赖。
func (i *InheritedElement) removeDependent(dependent Element) {
	delete(i.dependents, dependent)
}

// notifyDependents 通知所有依赖者 rebuild。
func (i *InheritedElement) notifyDependents() {
	for dependent := range i.dependents {
		if se, ok := dependent.(*StatefulElement); ok {
			se.didChangeDependencies()
		}
	}
}

// getInheritedWidgetOfExactType 从给定 Element 向上查找指定类型的 InheritedElement。
// 如果找到，将调用者 Element 注册为依赖者。
func getInheritedWidgetOfExactType(from Element, t reflect.Type) (InheritedWidget, bool) {
	for p := from.GetParent(); p != nil; p = p.GetParent() {
		if ie, ok := p.(*InheritedElement); ok {
			if reflect.TypeOf(ie.widget) == t {
				ie.addDependent(from)
				return ie.widget, true
			}
		}
	}
	return nil, false
}
