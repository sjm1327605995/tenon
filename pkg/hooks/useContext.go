package hooks

import (
	"sync"
)

// Context 上下文接口
type Context[T any] interface {
	Value() T
	SetValue(T)
	Subscribe(func(T))
	Unsubscribe(func(T))
}

// baseContext 基础上下文实现
type baseContext[T any] struct {
	value     T
	listeners []func(T)
	mu        sync.RWMutex
}

func (c *baseContext[T]) Value() T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func (c *baseContext[T]) SetValue(newValue T) {
	c.mu.Lock()
	c.value = newValue
	c.mu.Unlock()

	for _, listener := range c.listeners {
		listener(newValue)
	}
}

func (c *baseContext[T]) Subscribe(listener func(T)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.listeners = append(c.listeners, listener)
}

func (c *baseContext[T]) Unsubscribe(listener func(T)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, l := range c.listeners {
		if l == listener {
			c.listeners = append(c.listeners[:i], c.listeners[i+1:]...)
			break
		}
	}
}

// CreateContext 创建新的上下文
func CreateContext[T any](initialValue T) Context[T] {
	return &baseContext[T]{
		value:     initialValue,
		listeners: make([]func(T), 0),
	}
}

// UseContext 类似 React 的 useContext Hook
func UseContext[T any](ctx Context[T]) T {
	hook := getWorkInProgressHook()

	var prevValue T
	var listener func(T)
	if hook.MemoizedState != nil {
		state := hook.MemoizedState.(map[string]interface{})
		if v, ok := state["value"].(T); ok {
			prevValue = v
		}
		if l, ok := state["listener"]; ok && l != nil {
			listener = l.(func(T))
			ctx.Unsubscribe(listener)
		}
	}

	// 创建新监听器，当上下文变化时触发更新
	listener = func(newValue T) {
		if currentFiber != nil {
			currentFiber.Lanes = reconciler.MergeLanes(currentFiber.Lanes, reconciler.DefaultLane)
			reconciler.ScheduleUpdateOnFiber(currentFiber, reconciler.DefaultLane)
		}
	}
	ctx.Subscribe(listener)

	currentValue := ctx.Value()
	hook.MemoizedState = map[string]interface{}{
		"value":    currentValue,
		"listener": listener,
	}

	return currentValue
}
