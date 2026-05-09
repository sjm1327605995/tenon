package engine

import "github.com/sjm1327605995/tenon/internal/render"

// FCFunc 是函数组件的渲染函数签名。
type FCFunc func(h *Hooks) Widget

// FCWidget 包装函数组件。
type FCWidget struct {
	BaseWidget
	Render FCFunc
}

// FC 创建函数组件 Widget。
func FC(render FCFunc) Widget {
	return FCWidget{Render: render}
}

func (f FCWidget) CreateElement() Element {
	return newFCElement(f)
}

// ==================== fcElement ====================

// fcElement 管理函数组件的生命周期和 hooks。
type fcElement struct {
	ComponentElement
	child        Element
	hooks        *Hooks
	buildContext *elementBuildContext
}

func newFCElement(widget FCWidget) *fcElement {
	e := &fcElement{}
	e.ComponentElement.BaseElement.Init(e, widget)
	e.hooks = &Hooks{element: e}
	return e
}

func (e *fcElement) Mount(parent Element, slot int) {
	e.ComponentElement.Mount(parent, slot)
	e.buildContext = &elementBuildContext{element: e}
	e.PerformRebuild(nil)
}

func (e *fcElement) PerformRebuild(oldWidget Widget) {
	e.hooks.reset()
	w := e.GetWidget().(FCWidget)
	e.child = UpdateChild(e, e.child, w.Render(e.hooks))
}

func (e *fcElement) GetChildren() []Element {
	if e.child == nil {
		return nil
	}
	return []Element{e.child}
}

func (e *fcElement) FindRenderObject() render.RenderObject {
	if e.child != nil {
		return e.child.FindRenderObject()
	}
	return nil
}

func (e *fcElement) markNeedsBuild() {
	if defaultEngine != nil {
		defaultEngine.scheduleBuildFor(e)
	}
}

func (e *fcElement) Unmount() {
	for _, ef := range e.hooks.effects {
		if ef.cleanup != nil {
			ef.cleanup()
		}
	}
	if e.child != nil {
		e.child.Unmount()
		e.child = nil
	}
	e.ComponentElement.Unmount()
}
