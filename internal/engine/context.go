package engine

import "reflect"

// BuildContext 是 Widget 构建时的上下文。
// 提供跨层数据查找和祖先访问能力。
type BuildContext interface {
	// GetWidget 返回当前上下文关联的 Widget。
	GetWidget() Widget

	// FindAncestorWidgetOfExactType 向上查找最近的指定类型的祖先 Widget。
	// 返回找到的 Widget 和是否成功。不会注册依赖关系。
	FindAncestorWidgetOfExactType(t reflect.Type) (Widget, bool)

	// FindAncestorElementOfExactType 向上查找最近的指定类型的祖先 Element。
	// 供 InheritedElement 内部使用。不会注册依赖关系。
	FindAncestorElementOfExactType(t reflect.Type) (Element, bool)

	// DependOnInheritedWidgetOfExactType 向上查找最近的指定类型的 InheritedWidget，
	// 并将当前 Element 注册为依赖者。当 InheritedWidget 数据变化时，
	// 当前 Element 会自动 rebuild。
	DependOnInheritedWidgetOfExactType(t reflect.Type) (Widget, bool)
}

// elementBuildContext 是 BuildContext 的实现，依附于某个 Element。
type elementBuildContext struct {
	element Element
}

func (c *elementBuildContext) GetWidget() Widget {
	if c.element == nil {
		return nil
	}
	return c.element.GetWidget()
}

func (c *elementBuildContext) FindAncestorWidgetOfExactType(t reflect.Type) (Widget, bool) {
	el, ok := c.FindAncestorElementOfExactType(t)
	if !ok || el == nil {
		return nil, false
	}
	return el.GetWidget(), true
}

func (c *elementBuildContext) FindAncestorElementOfExactType(t reflect.Type) (Element, bool) {
	for p := c.element.GetParent(); p != nil; p = p.GetParent() {
		if reflect.TypeOf(p) == t {
			return p, true
		}
	}
	return nil, false
}

func (c *elementBuildContext) DependOnInheritedWidgetOfExactType(t reflect.Type) (Widget, bool) {
	iw, ok := getInheritedWidgetOfExactType(c.element, t)
	return iw, ok
}
