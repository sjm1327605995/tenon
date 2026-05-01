package ui

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
)

// Element 是 Element 树中的节点，负责 Widget 的生命周期管理和 RenderObject 的创建/更新。
// 对应 Flutter 的 Element。
type Element interface {
	// 生命周期
	Mount(parent Element, slot int)
	Update(newWidget Widget)
	Unmount()

	// 树关系
	GetWidget() Widget
	GetParent() Element
	GetChildren() []Element
	GetSlot() int

	// RenderObject 关联
	FindRenderObject() render.RenderObject
}

// BaseElement 提供 Element 的默认实现。
type BaseElement struct {
	self   Element
	widget Widget
	parent Element
	slot   int
}

func (b *BaseElement) Init(self Element, widget Widget) {
	b.self = self
	b.widget = widget
}

func (b *BaseElement) GetWidget() Widget     { return b.widget }
func (b *BaseElement) GetParent() Element    { return b.parent }
func (b *BaseElement) GetChildren() []Element { return nil }
func (b *BaseElement) GetSlot() int          { return b.slot }

func (b *BaseElement) Mount(parent Element, slot int) {
	b.parent = parent
	b.slot = slot
}

func (b *BaseElement) Update(newWidget Widget) {
	b.widget = newWidget
}

func (b *BaseElement) Unmount() {
	b.parent = nil
}

func (b *BaseElement) FindRenderObject() render.RenderObject {
	// 默认实现：遍历子元素查找 RenderObject
	for _, child := range b.self.GetChildren() {
		if ro := child.FindRenderObject(); ro != nil {
			return ro
		}
	}
	return nil
}

// ComponentElement 管理一个或多个子 Element 的组件元素。
// 对应 Flutter 的 ComponentElement。
type ComponentElement struct {
	BaseElement
}

func NewComponentElement(widget Widget) *ComponentElement {
	e := &ComponentElement{}
	e.BaseElement.Init(e, widget)
	return e
}

func (c *ComponentElement) Mount(parent Element, slot int) {
	c.BaseElement.Mount(parent, slot)
}

func (c *ComponentElement) Update(newWidget Widget) {
	oldWidget := c.widget
	c.BaseElement.Update(newWidget)
	c.PerformRebuild(oldWidget)
}

func (c *ComponentElement) Unmount() {
	c.BaseElement.Unmount()
}

// PerformRebuild 用新 Widget 重建子树。
// 子类可覆盖此方法来控制子树的构建方式。
func (c *ComponentElement) PerformRebuild(oldWidget Widget) {
	// 默认行为：由子类实现
}

// SingleChildComponentElement 管理单个子 Element 的组件元素。
type SingleChildComponentElement struct {
	ComponentElement
	Child Element
}

func NewSingleChildComponentElement(widget Widget) *SingleChildComponentElement {
	e := &SingleChildComponentElement{}
	e.ComponentElement.BaseElement.Init(e, widget)
	return e
}

func (s *SingleChildComponentElement) GetChildren() []Element {
	if s.Child == nil {
		return nil
	}
	return []Element{s.Child}
}

func (s *SingleChildComponentElement) Update(newWidget Widget) {
	oldWidget := s.widget
	s.BaseElement.Update(newWidget)
	s.UpdateChild(oldWidget)
}

// UpdateChild 由子类覆盖，获取新 Widget 的子 Widget 并 diff。
func (s *SingleChildComponentElement) UpdateChild(oldWidget Widget) {
	// 子类实现
}

func (s *SingleChildComponentElement) Unmount() {
	if s.Child != nil {
		s.Child.Unmount()
		s.Child = nil
	}
	s.ComponentElement.Unmount()
}

// RenderObjectElement 管理单个 RenderObject 的 Element。
// 对应 Flutter 的 RenderObjectElement。
type RenderObjectElement struct {
	BaseElement
	RenderObject render.RenderObject
}

func NewRenderObjectElement(widget Widget) *RenderObjectElement {
	e := &RenderObjectElement{}
	e.BaseElement.Init(e, widget)
	return e
}

func (r *RenderObjectElement) GetRenderObject() render.RenderObject { return r.RenderObject }
func (r *RenderObjectElement) FindRenderObject() render.RenderObject { return r.RenderObject }

func (r *RenderObjectElement) Mount(parent Element, slot int) {
	r.BaseElement.Mount(parent, slot)
	if r.RenderObject == nil {
		if roe, ok := r.self.(interface {
			CreateRenderObject() render.RenderObject
		}); ok {
			r.RenderObject = roe.CreateRenderObject()
		}
	}
	if r.RenderObject != nil {
		// 将 RenderObject 挂载到父 RenderObject 树
		if parent != nil {
			if parentRO := parent.FindRenderObject(); parentRO != nil {
				parentRO.AddChild(r.RenderObject)
			}
		}
	}
}

func (r *RenderObjectElement) Update(newWidget Widget) {
	oldWidget := r.widget
	r.BaseElement.Update(newWidget)
	if roe, ok := r.self.(interface {
		UpdateRenderObject(oldWidget Widget)
	}); ok {
		roe.UpdateRenderObject(oldWidget)
	}
}

func (r *RenderObjectElement) Unmount() {
	if r.RenderObject != nil {
		// 从父 RenderObject 移除
		if r.parent != nil {
			if parentRO := r.parent.FindRenderObject(); parentRO != nil {
				parentRO.RemoveChild(r.RenderObject)
			}
		}
		r.RenderObject.Detach()
		r.RenderObject = nil
	}
	r.BaseElement.Unmount()
}

// CreateRenderObject 由子类覆盖，创建对应的 RenderObject。
func (r *RenderObjectElement) CreateRenderObject() render.RenderObject {
	return nil
}

// UpdateRenderObject 由子类覆盖，将新 Widget 的属性同步到 RenderObject。
func (r *RenderObjectElement) UpdateRenderObject(oldWidget Widget) {
	// 子类实现
}

// SingleChildRenderObjectElement 管理单个子 Element 的 RenderObjectElement。
type SingleChildRenderObjectElement struct {
	RenderObjectElement
	Child Element
}

func NewSingleChildRenderObjectElement(widget Widget) *SingleChildRenderObjectElement {
	e := &SingleChildRenderObjectElement{}
	e.RenderObjectElement.BaseElement.Init(e, widget)
	return e
}

func (s *SingleChildRenderObjectElement) GetChildren() []Element {
	if s.Child == nil {
		return nil
	}
	return []Element{s.Child}
}

func (s *SingleChildRenderObjectElement) Mount(parent Element, slot int) {
	s.RenderObjectElement.Mount(parent, slot)
	// 子元素在子类的 Mount 或 Update 中创建
}

func (s *SingleChildRenderObjectElement) Update(newWidget Widget) {
	oldWidget := s.widget
	s.RenderObjectElement.Update(newWidget)
	if sce, ok := s.self.(interface {
		UpdateChild(oldWidget Widget)
	}); ok {
		sce.UpdateChild(oldWidget)
	}
}

// UpdateChild 由子类覆盖。
func (s *SingleChildRenderObjectElement) UpdateChild(oldWidget Widget) {
	// 子类实现
}

func (s *SingleChildRenderObjectElement) Unmount() {
	if s.Child != nil {
		s.Child.Unmount()
		s.Child = nil
	}
	s.RenderObjectElement.Unmount()
}

// MultiChildRenderObjectElement 管理多个子 Element 的 RenderObjectElement。
type MultiChildRenderObjectElement struct {
	RenderObjectElement
	Children []Element
}

func NewMultiChildRenderObjectElement(widget Widget) *MultiChildRenderObjectElement {
	e := &MultiChildRenderObjectElement{}
	e.RenderObjectElement.BaseElement.Init(e, widget)
	return e
}

func (m *MultiChildRenderObjectElement) GetChildren() []Element { return m.Children }

func (m *MultiChildRenderObjectElement) Mount(parent Element, slot int) {
	m.RenderObjectElement.Mount(parent, slot)
}

func (m *MultiChildRenderObjectElement) Update(newWidget Widget) {
	oldWidget := m.widget
	m.RenderObjectElement.Update(newWidget)
	if mce, ok := m.self.(interface {
		UpdateChildren(oldWidget Widget)
	}); ok {
		mce.UpdateChildren(oldWidget)
	}
}

// UpdateChildren 由子类覆盖。
func (m *MultiChildRenderObjectElement) UpdateChildren(oldWidget Widget) {
	// 子类实现
}

func (m *MultiChildRenderObjectElement) Unmount() {
	for _, child := range m.Children {
		child.Unmount()
	}
	m.Children = nil
	m.RenderObjectElement.Unmount()
}

// UpdateChild 是对单个 child 进行 diff 的便捷方法。
func UpdateChild(parent Element, child Element, newWidget Widget) Element {
	if newWidget == nil {
		if child != nil {
			child.Unmount()
		}
		return nil
	}
	if child != nil {
		if CanUpdate(child.GetWidget(), newWidget) {
			child.Update(newWidget)
			return child
		}
		child.Unmount()
	}
	newChild := newWidget.CreateElement()
	newChild.Mount(parent, 0)
	return newChild
}

// UpdateChildren 对新旧两组 Widget 做同级对比，尽可能复用旧 Element。
// 这是 ComponentElement / MultiChildRenderObjectElement 的核心 diff 方法。
func UpdateChildren(parent Element, oldChildren []Element, newWidgets []Widget) []Element {
	newChildren := make([]Element, 0, len(newWidgets))

	// 用于旧 Element 复用的查找表（按 Key）
	oldKeyed := make(map[string]Element)
	for _, old := range oldChildren {
		if old.GetWidget() != nil && !IsNilKey(old.GetWidget().GetKey()) {
			oldKeyed[old.GetWidget().GetKey().String()] = old
		}
	}

	var oldIndex int
	for i, newWidget := range newWidgets {
		var oldChild Element

		// 1. 尝试按 Key 查找可复用的旧 Element
		if !IsNilKey(newWidget.GetKey()) {
			if keyed, ok := oldKeyed[newWidget.GetKey().String()]; ok {
				oldChild = keyed
				delete(oldKeyed, newWidget.GetKey().String())
			}
		}

		// 2. 无 Key 匹配时，按位置取旧 Element
		if oldChild == nil && oldIndex < len(oldChildren) {
			for oldIndex < len(oldChildren) {
				candidate := oldChildren[oldIndex]
				if candidate.GetWidget() != nil && !IsNilKey(candidate.GetWidget().GetKey()) {
					if _, stillInOld := oldKeyed[candidate.GetWidget().GetKey().String()]; stillInOld {
						oldIndex++
						continue
					}
				}
				oldChild = candidate
				oldIndex++
				break
			}
		}

		// 3. 尝试复用或创建新 Element
		if oldChild != nil && CanUpdate(oldChild.GetWidget(), newWidget) {
			oldChild.Update(newWidget)
			newChildren = append(newChildren, oldChild)
		} else {
			if oldChild != nil {
				oldChild.Unmount()
			}
			newChild := newWidget.CreateElement()
			newChild.Mount(parent, i)
			newChildren = append(newChildren, newChild)
		}
	}

	// 4. 卸载剩余未被复用的旧 Element
	for oldIndex < len(oldChildren) {
		oldChildren[oldIndex].Unmount()
		oldIndex++
	}
	for _, leftover := range oldKeyed {
		leftover.Unmount()
	}

	return newChildren
}
