package hooks

import (
	"fmt"
	"reflect"
	"sync"
)

type HookType int

const (
	HookTypeState HookType = iota
	HookTypeEffect
	HookTypeContext
)

type Hook interface {
	Type() HookType
	Update(value interface{})
	Cleanup()
}

type ComponentHooks struct {
	hooks     []Hook
	index     int
	component Component
}

type Component interface {
	ID() string
	MarkDirty()
}

var (
	currentComponent Component
	currentHooks     *ComponentHooks
	componentMutex   sync.RWMutex
	components       = make(map[string]*ComponentHooks)
)

func NewComponentHooks(component Component) *ComponentHooks {
	return &ComponentHooks{
		hooks:     make([]Hook, 0),
		component: component,
	}
}

func WithHooks(component Component, renderFunc func()) {
	componentMutex.Lock()
	defer componentMutex.Unlock()

	hooks, exists := components[component.ID()]
	if !exists {
		hooks = NewComponentHooks(component)
		components[component.ID()] = hooks
	}

	prevComponent := currentComponent
	prevHooks := currentHooks

	currentComponent = component
	currentHooks = hooks

	hooks.index = 0
	renderFunc()

	currentComponent = prevComponent
	currentHooks = prevHooks

	if len(hooks.hooks) > hooks.index {
		for i := hooks.index; i < len(hooks.hooks); i++ {
			hooks.hooks[i].Cleanup()
		}
		hooks.hooks = hooks.hooks[:hooks.index]
	}
}

func getCurrentHooks() *ComponentHooks {
	if currentHooks == nil {
		panic("Hooks can only be called within a component's render function")
	}
	return currentHooks
}

func nextHook() Hook {
	hooks := getCurrentHooks()
	if hooks.index < len(hooks.hooks) {
		hook := hooks.hooks[hooks.index]
		hooks.index++
		return hook
	}
	panic("Hooks must be called in the same order every render")
}

func addHook(hook Hook) {
	hooks := getCurrentHooks()
	hooks.hooks = append(hooks.hooks, hook)
	hooks.index++
}

func CleanupComponentHooks(componentID string) {
	componentMutex.Lock()
	defer componentMutex.Unlock()

	if hooks, exists := components[componentID]; exists {
		for _, hook := range hooks.hooks {
			hook.Cleanup()
		}
		delete(components, componentID)
	}
}

func DebugHooks() {
	componentMutex.RLock()
	defer componentMutex.RUnlock()

	fmt.Printf("Total components with hooks: %d\n", len(components))
	for id, hooks := range components {
		fmt.Printf("Component %s: %d hooks\n", id, len(hooks.hooks))
		for i, hook := range hooks.hooks {
			fmt.Printf("  Hook %d: Type=%v\n", i, hook.Type())
		}
	}
}

type HookValue struct {
	value     interface{}
	setter    func(interface{})
	listeners []func(interface{})
}

func NewHookValue(initialValue interface{}) *HookValue {
	return &HookValue{
		value:     initialValue,
		listeners: make([]func(interface{}), 0),
	}
}

func (hv *HookValue) Get() interface{} {
	return hv.value
}

func (hv *HookValue) Set(newValue interface{}) {
	if !reflect.DeepEqual(hv.value, newValue) {
		hv.value = newValue
		for _, listener := range hv.listeners {
			listener(newValue)
		}
	}
}

func (hv *HookValue) Subscribe(listener func(interface{})) {
	hv.listeners = append(hv.listeners, listener)
}
