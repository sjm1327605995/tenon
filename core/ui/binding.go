package ui

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/sjm1327605995/tenon/yoga"
)

// ---------- React-like Hooks System ----------

// State 表示一个状态值，使用泛型
type State[T any] struct {
	value    T
	updaters []func()
	mutex    sync.Mutex
}

// StateUpdater 表示状态更新函数类型
type StateUpdater[T any] func(newValue T)

// Hook 表示一个hook实例
type Hook struct {
	State interface{}
}

// 当前组件上下文，用于hooks访问
var currentContext *Context

// 当前组件实例，用于访问组件的hooks列表
var currentComponent *FunctionComponent

// 当前组件的hooks索引
var currentHookIndex int

// UseState 创建一个新的状态值和更新函数，类似React的useState hook
func UseState[T any](initial T) (T, StateUpdater[T]) {
	// 检查当前是否在组件渲染上下文中
	if currentContext == nil || currentComponent == nil {
		// 如果不在上下文中，返回一个简单的状态实现
		state := &State[T]{
			value: initial,
		}
		return state.Value(), func(newValue T) {
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

	// 保存当前上下文引用
	ctx := currentContext

	// 获取当前hook索引
	index := currentHookIndex
	currentHookIndex++

	// 检查是否已有hook，如果没有则创建
	if index >= len(currentComponent.hooks) {
		// 创建新的state
		state := &State[T]{
			value: initial,
		}

		// 订阅状态更新，使用保存的上下文引用
		state.Subscribe(func() {
			ctx.Update()
		})

		// 添加到组件的hooks列表
		currentComponent.hooks = append(currentComponent.hooks, &Hook{State: state})
	}

	// 获取当前hook的state
	state := currentComponent.hooks[index].State.(*State[T])

	// 返回状态值和更新函数
	return state.Value(), func(newValue T) {
		// 更新状态
		state.mutex.Lock()
		state.value = newValue
		updaters := make([]func(), len(state.updaters))
		copy(updaters, state.updaters)
		state.mutex.Unlock()

		// 触发更新
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

// ---------- Component System ----------

// Props 表示组件属性的通用接口
type Props interface{}

// FC 表示函数组件类型，类似React的FunctionComponent
type FC[T Props] func(props T) UI

// FunctionComponent 表示函数组件实例
type FunctionComponent struct {
	RenderFunc FC[Props] // 组件渲染函数
	Props      Props     // 组件属性
	Context    *Context  // 组件上下文
	hooks      []*Hook   // 组件的hooks列表
}

// NewFunctionComponent 创建一个新的函数组件实例
func NewFunctionComponent[T Props](renderFunc FC[T], props T) UI {
	return &FunctionComponent{
		RenderFunc: func(p Props) UI { return renderFunc(p.(T)) },
		Props:      props,
	}
}

// Render 实现UI接口，渲染组件
func (c *FunctionComponent) Render() *Element {
	// 保存之前的状态
	previousContext := currentContext
	previousComponent := currentComponent
	previousHookIndex := currentHookIndex

	// 设置当前组件状态
	currentContext = c.Context
	currentComponent = c
	currentHookIndex = 0

	defer func() {
		// 恢复之前的状态
		currentContext = previousContext
		currentComponent = previousComponent
		currentHookIndex = previousHookIndex
	}()

	// 调用渲染函数
	return c.RenderFunc(c.Props).Render()
}

// SetContext 设置组件上下文
func (c *FunctionComponent) SetContext(ctx *Context) {
	c.Context = ctx
}

// ---------- Context System ----------

// Context 管理组件的渲染上下文，类似React的Context
type Context struct {
	rootUI      UI
	rootEl      *Element
	mutex       sync.Mutex
	needsUpdate bool
}

// NewContext 创建一个新的上下文
func NewContext(root UI) *Context {
	ctx := &Context{
		rootUI: root,
	}

	// 先设置组件上下文，再渲染
	if comp, ok := root.(interface{ SetContext(*Context) }); ok {
		comp.SetContext(ctx)
	}

	// 渲染根元素
	ctx.rootEl = root.Render()

	return ctx
}

// Update 触发组件树更新
func (c *Context) Update() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.needsUpdate = true
	// 标记需要更新，实际更新将在FrameEvent中处理
	// 不要在这里重新创建Element树，因为这会导致组件状态丢失
}

// GetRootElement 获取根元素
func (c *Context) GetRootElement() *Element {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 如果需要更新，重新创建Element树
	if c.needsUpdate {
		c.rootEl = c.rootUI.Render()
	}

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

// ---------- Example Components ----------

// CounterProps 计数器组件属性
type CounterProps struct {
	InitialCount int
}

// Counter 计数器组件，类似React的函数组件
func Counter(props CounterProps) UI {
	// 使用类似React的useState hook
	count, setCount := UseState(props.InitialCount)

	return View(
		Text().Content(fmt.Sprintf("Count: %d", count)),
		View(
			View().Background(color.NRGBA{G: 255, A: 255}).Height(Px(50)).Width(Px(100)).OnClick(func() {
				fmt.Printf("Before increment: %d\n", count)
				setCount(count + 1)
				fmt.Printf("After increment: %d\n", count+1)
			}),
			Text().Content("+"),
		),
		View(
			View().Background(color.NRGBA{R: 255, A: 255}).Height(Px(50)).Width(Px(100)).OnClick(func() {
				fmt.Printf("Before decrement: %d\n", count)
				setCount(count - 1)
				fmt.Printf("After decrement: %d\n", count-1)
			}),
			Text().Content("-"),
		),
	).
		FlexDirection(yoga.FlexDirectionColumn).
		JustifyContent(yoga.JustifyCenter).
		AlignItems(yoga.AlignCenter).
		Width(Percent(100)).
		Height(Percent(100)).
		Background(color.NRGBA{B: 255, A: 128})
}

// CreateCounter 创建一个计数器组件实例
func CreateCounter() UI {
	return NewFunctionComponent(Counter, CounterProps{InitialCount: 0})
}
