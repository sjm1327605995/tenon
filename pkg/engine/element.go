package engine

import (
	"sync"

	"github.com/sjm1327605995/tenon/pkg/render"
)

// globalKeyRegistry 维护 GlobalKey 到 Element 的映射。
var (
	globalKeyRegistry   = make(map[*GlobalKey]Element)
	globalKeyRegistryMu sync.RWMutex
)

func registerGlobalKey(key *GlobalKey, element Element) {
	globalKeyRegistryMu.Lock()
	globalKeyRegistry[key] = element
	globalKeyRegistryMu.Unlock()
}

func unregisterGlobalKey(key *GlobalKey) {
	globalKeyRegistryMu.Lock()
	delete(globalKeyRegistry, key)
	globalKeyRegistryMu.Unlock()
}

func getGlobalKeyElement(key *GlobalKey) Element {
	globalKeyRegistryMu.RLock()
	el := globalKeyRegistry[key]
	globalKeyRegistryMu.RUnlock()
	return el
}

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
	GetRenderObject() render.RenderObject
	GetBuildContext() BuildContext
	GetState() State

	// InheritedWidget 关联
	GetInheritedWidget() InheritedWidget
	AddDependent(dependent Element)
	DidChangeDependencies()
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
func (b *BaseElement) GetRenderObject() render.RenderObject { return nil }
func (b *BaseElement) GetBuildContext() BuildContext        { return nil }
func (b *BaseElement) GetState() State                      { return nil }
func (b *BaseElement) GetInheritedWidget() InheritedWidget  { return nil }
func (b *BaseElement) AddDependent(dependent Element)        {}
func (b *BaseElement) DidChangeDependencies()                {}

func (b *BaseElement) Mount(parent Element, slot int) {
	b.parent = parent
	b.slot = slot
	if key := b.widget.GetKey(); key != nil {
		if gk := key.AsGlobalKey(); gk != nil {
			registerGlobalKey(gk, b.self)
		}
	}
}

func (b *BaseElement) Update(newWidget Widget) {
	b.widget = newWidget
}

func (b *BaseElement) Unmount() {
	if key := b.widget.GetKey(); key != nil {
		if gk := key.AsGlobalKey(); gk != nil {
			unregisterGlobalKey(gk)
		}
	}
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
		if factory, ok := r.widget.(RenderObjectFactory); ok {
			r.RenderObject = factory.CreateRenderObject(r.self)
		}
	}
	if r.RenderObject != nil {
		// 将 RenderObject 挂载到最近的 RenderObjectElement 祖先
		// 跳过 ComponentElement（如 StatefulElement），因为它们没有自己的 RenderObject
		for p := parent; p != nil; p = p.GetParent() {
			if pro := p.GetRenderObject(); pro != nil {
				pro.AddChild(r.RenderObject)
				break
			}
		}
	}
}

func (r *RenderObjectElement) Update(newWidget Widget) {
	oldWidget := r.widget
	r.BaseElement.Update(newWidget)
	if updater, ok := r.widget.(RenderObjectUpdater); ok && r.RenderObject != nil {
		updater.UpdateRenderObject(r.RenderObject, oldWidget)
	}
}

func (r *RenderObjectElement) Unmount() {
	if r.RenderObject != nil {
		// 从最近的 RenderObjectElement 祖先移除
		// 跳过 ComponentElement（如 StatefulElement），因为它们没有自己的 RenderObject
		for p := r.parent; p != nil; p = p.GetParent() {
			if pro := p.GetRenderObject(); pro != nil {
				pro.RemoveChild(r.RenderObject)
				break
			}
		}
		r.RenderObject.Detach()
		r.RenderObject = nil
	}
	r.BaseElement.Unmount()
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
	// 委托给 Widget 管理子元素（新接口）
	if provider, ok := s.widget.(SingleChildProvider); ok {
		s.Child = UpdateChild(s.self, s.Child, provider.GetChildWidget())
	}
	// 否则子元素在旧 Element 子类的 Mount 或 Update 中创建
}

func (s *SingleChildRenderObjectElement) Update(newWidget Widget) {
	s.RenderObjectElement.Update(newWidget)
	if provider, ok := s.widget.(SingleChildProvider); ok {
		s.Child = UpdateChild(s.self, s.Child, provider.GetChildWidget())
	}
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
	// 委托给 Widget 管理子元素（新接口）
	if provider, ok := m.widget.(MultiChildProvider); ok {
		m.Children = UpdateChildren(m.self, m.Children, provider.GetChildrenWidgets())
	}
}

func (m *MultiChildRenderObjectElement) Update(newWidget Widget) {
	m.RenderObjectElement.Update(newWidget)
	if provider, ok := m.widget.(MultiChildProvider); ok {
		m.Children = UpdateChildren(m.self, m.Children, provider.GetChildrenWidgets())
	}
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

// UpdateChildren 对新旧两组 Widget 做同级对比，按位置复用旧 Element。
// 同位置且 CanUpdate 匹配则复用，否则销毁重建。
func UpdateChildren(parent Element, oldChildren []Element, newWidgets []Widget) []Element {
	newChildren := make([]Element, 0, len(newWidgets))

	for i, newWidget := range newWidgets {
		var oldChild Element
		if i < len(oldChildren) {
			oldChild = oldChildren[i]
		}

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

	// 卸载多余的旧子节点
	for i := len(newWidgets); i < len(oldChildren); i++ {
		if oldChildren[i] != nil {
			oldChildren[i].Unmount()
		}
	}

	return newChildren
}
