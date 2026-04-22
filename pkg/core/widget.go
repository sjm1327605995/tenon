package core

import (
	"fmt"
	"sync/atomic"
)

// ==================== Hook 类型 ====================

type hook interface{}

type stateHook struct {
	state any
}

type effectHook struct {
	effect  func() func()
	cleanup func()
	deps    []any
	hasRun  bool
}

type memoHook struct {
	value any
	deps  []any
}

type refHook struct {
	current any
}

// ==================== BaseWidget ====================

// BaseWidget 是用户自定义组件的基类。
// 用户组件嵌入 BaseWidget 即可自动获得生命周期默认实现和 Hooks。
type BaseWidget struct {
	self    Widget
	engine  *Engine
	hostRef Host

	props any
	state any

	hooks   []hook
	hookIdx int
	dirty   bool
}

// Init 必须在用户组件的构造函数中调用，参数传自己（*MyComponent）。
func (b *BaseWidget) Init(self Widget) {
	b.self = self
}

// === Component 接口默认实现 ===

func (b *BaseWidget) Render() Component { return nil }

func (b *BaseWidget) ComponentDidMount()              {}
func (b *BaseWidget) ComponentWillUnmount()           {}
func (b *BaseWidget) ComponentDidUpdate(_, _ any)     {}
func (b *BaseWidget) ShouldComponentUpdate(_ any) bool { return true }

func (b *BaseWidget) SetHostRef(host Host) { b.hostRef = host }
func (b *BaseWidget) GetHostRef() Host     { return b.hostRef }

// Invalidate 标记组件需要重新 Render，并在下一帧触发更新。
func (b *BaseWidget) Invalidate() {
	b.dirty = true
	if b.engine != nil {
		b.engine.scheduleUpdate(b.self)
	}
}

func (b *BaseWidget) isComponent() {}

func (b *BaseWidget) setEngine(e *Engine) { b.engine = e }
func (b *BaseWidget) getEngine() *Engine  { return b.engine }
func (b *BaseWidget) resetHooks()         { b.hookIdx = 0 }

// ==================== Hooks ====================

// UseState 返回当前状态和一个设置状态的函数。
// 调用设置函数会触发组件重新 Render。
func (b *BaseWidget) UseState(initial any) (any, func(any)) {
	if len(b.hooks) <= b.hookIdx {
		b.hooks = append(b.hooks, &stateHook{state: initial})
	}
	h := b.hooks[b.hookIdx].(*stateHook)
	b.hookIdx++

	setter := func(v any) {
		h.state = v
		b.Invalidate()
	}
	return h.state, setter
}

// UseEffect 在组件挂载后和依赖变化时执行 effect。
// effect 可以返回一个清理函数，在组件卸载或依赖变化前调用。
func (b *BaseWidget) UseEffect(effect func() func(), deps []any) {
	if len(b.hooks) <= b.hookIdx {
		b.hooks = append(b.hooks, &effectHook{
			effect: effect,
			deps:   deps,
		})
	}
	h := b.hooks[b.hookIdx].(*effectHook)
	b.hookIdx++

	shouldRun := !h.hasRun || !depsEqual(h.deps, deps)
	if shouldRun {
		if h.cleanup != nil {
			h.cleanup()
		}
		h.cleanup = effect()
		h.deps = deps
		h.hasRun = true
	}
}

// UseMemo 缓存计算结果，仅在依赖变化时重新计算。
func (b *BaseWidget) UseMemo(fn func() any, deps []any) any {
	if len(b.hooks) <= b.hookIdx {
		b.hooks = append(b.hooks, &memoHook{
			value: fn(),
			deps:  deps,
		})
	}
	h := b.hooks[b.hookIdx].(*memoHook)
	b.hookIdx++

	if !depsEqual(h.deps, deps) {
		h.value = fn()
		h.deps = deps
	}
	return h.value
}

// Ref 是一个可变引用容器。
type Ref struct {
	Current any
}

// UseRef 返回一个持久化的引用对象。
func (b *BaseWidget) UseRef(initial any) *Ref {
	if len(b.hooks) <= b.hookIdx {
		b.hooks = append(b.hooks, &refHook{current: initial})
	}
	h := b.hooks[b.hookIdx].(*refHook)
	b.hookIdx++
	return &Ref{Current: h.current}
}

// UseCallback 是 UseMemo 的语法糖，用于缓存函数引用。
func (b *BaseWidget) UseCallback(fn func(), deps []any) func() {
	return b.UseMemo(func() any { return fn }, deps).(func())
}

var idCounter int64

// UseId 生成一个唯一的 ID 字符串。
func (b *BaseWidget) UseId() string {
	v := atomic.AddInt64(&idCounter, 1)
	return fmt.Sprintf("id-%d", v)
}

// Context 是一个简单的上下文值容器。
type Context struct {
	value any
}

// NewContext 创建一个上下文。
func NewContext(defaultValue any) *Context {
	return &Context{value: defaultValue}
}

// UseContext 读取上下文值。
func (b *BaseWidget) UseContext(ctx *Context) any {
	return ctx.value
}

// UseTransition 返回一个过渡状态标记和一个启动过渡的函数。
// 当前为简化实现，直接同步执行。
func (b *BaseWidget) UseTransition() (bool, func(func())) {
	isPending, _ := b.UseState(false)
	setPending, _ := b.UseState(false)
	_ = setPending

	startTransition := func(fn func()) {
		_ = isPending
		fn()
	}
	return false, startTransition
}

// ==================== 工具函数 ====================

func depsEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
