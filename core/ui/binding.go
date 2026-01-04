package ui

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/sjm1327605995/tenon/yoga"
)

// State 表示一个状态值，使用泛型
type State[T any] struct {
	value    T
	updaters []func()
	mutex    sync.Mutex
}

// StateUpdater 表示状态更新函数类型
type StateUpdater[T any] func(newValue T)

// UseState 创建一个新的状态值和更新函数
func UseState[T any](initial T) (State[T], StateUpdater[T]) {
	state := State[T]{
		value: initial,
	}
	return state, func(newValue T) {
		state.mutex.Lock()
		state.value = newValue
		updaters := make([]func(), len(state.updaters))
		copy(updaters, state.updaters)
		state.mutex.Unlock()
		for _, updater := range updaters {
			updater()
		}
	}
}

// Subscribe 添加状态更新监听器
func (s *State[T]) Subscribe(updater func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.updaters = append(s.updaters, updater)
}

// Value 获取当前状态值
func (s *State[T]) Value() T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.value
}

// Props 表示组件的属性，使用泛型
type Props[T any] struct {
	Value T
}

// Component 表示一个通用组件，使用泛型
type Component[T any] struct {
	Props    Props[T]
	renderFn func(props Props[T]) UI
	updater  func()
}

// ComponentFn 表示组件渲染函数类型
type ComponentFn[T any] func(props Props[T]) UI

// NewComponent 创建一个新的组件
func NewComponent[T any](renderFn ComponentFn[T]) *Component[T] {
	return &Component[T]{
		renderFn: renderFn,
	}
}

// SetProps 设置组件属性
func (c *Component[T]) SetProps(props Props[T]) {
	c.Props = props
	if c.updater != nil {
		c.updater()
	}
}

// Render 实现UI接口
func (c *Component[T]) Render() *Element {
	return c.renderFn(c.Props).Render()
}

// SetUpdater 设置更新函数
func (c *Component[T]) SetUpdater(updater func()) {
	c.updater = updater
}

// Context 管理组件的渲染上下文
type Context struct {
	rootUI      UI
	rootEl      *Element
	mutex       sync.Mutex
	needsUpdate bool
}

// NewContext 创建一个新的上下文
func NewContext(root UI) *Context {
	return &Context{
		rootUI: root,
		rootEl: root.Render(),
	}
}

// Update 触发组件树更新
func (c *Context) Update() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.needsUpdate = true
	// 重新创建Element树
	c.rootEl = c.rootUI.Render()
}

// GetRootElement 获取根元素
func (c *Context) GetRootElement() *Element {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.rootEl
}

// NeedsUpdate 检查是否需要更新
func (c *Context) NeedsUpdate() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.needsUpdate
}

// ClearNeedsUpdate 清除更新标志
func (c *Context) ClearNeedsUpdate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.needsUpdate = false
}

// StatefulComponent 表示一个有状态组件的基类
type StatefulComponent struct {
	context *Context
}

// SetContext 设置组件上下文
func (c *StatefulComponent) SetContext(ctx *Context) {
	c.context = ctx
}

// UpdateContext 更新上下文
func (c *StatefulComponent) UpdateContext() {
	if c.context != nil {
		c.context.Update()
	}
}

// CounterProps 计数器组件属性
type CounterProps struct {
	InitialCount int
}

// CounterComponent 计数器组件
type CounterComponent struct {
	count     State[int]
	increment func()
	decrement func()
}

// NewCounterComponent 创建一个新的计数器组件
func NewCounterComponent() UI {
	c := &CounterComponent{}
	count, setCount := UseState(0)
	c.count = count
	c.increment = func() {
		setCount(c.count.Value() + 1)
	}
	c.decrement = func() {
		setCount(c.count.Value() - 1)
	}
	return c
}

// Render 实现UI接口
func (c *CounterComponent) Render() *Element {
	// 这里应该返回一个UI组件，而不是直接创建Element
	// 但为了保持与现有代码的兼容性，我们暂时返回nil
	// 后续需要重新设计组件系统
	return nil
}

// CounterUI 计数器UI组件，返回ui.UI接口
type CounterUI struct {
	count     State[int]
	increment func()
	decrement func()
}

// NewCounterUI 创建一个新的计数器UI组件
func NewCounterUI() UI {
	c := &CounterUI{}
	count, setCount := UseState(0)
	c.count = count
	c.increment = func() {
		setCount(c.count.Value() + 1)
	}
	c.decrement = func() {
		setCount(c.count.Value() - 1)
	}
	return c
}

// Render 实现UI接口，返回ui.UI组件
func (c *CounterUI) Render() *Element {
	// 这里应该返回一个UI组件，而不是直接创建Element
	// 但为了保持与现有代码的兼容性，我们暂时返回nil
	// 后续需要重新设计组件系统
	return nil
}

// 创建一个新的组件类型，用于包装有状态组件
type StatelessComponent struct {
	renderFn func() UI
}

// NewStatelessComponent 创建一个新的无状态组件
func NewStatelessComponent(renderFn func() UI) UI {
	return &StatelessComponent{
		renderFn: renderFn,
	}
}

// Render 实现UI接口
func (c *StatelessComponent) Render() *Element {
	return c.renderFn().Render()
}

// CreateCounter 创建一个计数器组件，返回ui.UI接口
func CreateCounter() UI {
	count, _ := UseState(0)

	return NewStatelessComponent(func() UI {
		return View(
			Text().Content(fmt.Sprintf("Count: %d", count.Value())),
			View(
				View().Background(color.NRGBA{G: 255, A: 255}).Height(Px(50)).Width(Px(100)),
				Text().Content("+"),
			),
			View(
				View().Background(color.NRGBA{R: 255, A: 255}).Height(Px(50)).Width(Px(100)),
				Text().Content("-"),
			),
		).
			FlexDirection(yoga.FlexDirectionColumn).
			JustifyContent(yoga.JustifyCenter).
			AlignItems(yoga.AlignCenter).
			Width(Percent(100)).
			Height(Percent(100)).
			Background(color.NRGBA{B: 255, A: 128})
	})
}
