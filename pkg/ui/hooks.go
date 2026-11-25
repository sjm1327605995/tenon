package ui

import "reflect"

// 渲染期调度器：渲染在 Gio 单协程内同步进行，包级指针即可，无需 GLS。
var (
	currentFiber *Fiber
	hookIndex    int
)

// Cleanup 是 UseEffect 的清理函数（可为 nil）。
type Cleanup func()

type stateHook struct {
	value    any
	setter   any // 缓存以保证 setter 跨渲染引用稳定（React 语义）
	dispatch any
	reducer  any
}

type effectHook struct {
	deps    []any
	cleanup Cleanup
	hasRun  bool
}

type refHook struct {
	value any
}

// renderWithHooks 在 hook 调度上下文中执行组件的 render。
func renderWithHooks(f *Fiber) *Node {
	prevF, prevI := currentFiber, hookIndex
	currentFiber, hookIndex = f, 0
	defer func() { currentFiber, hookIndex = prevF, prevI }()
	return f.render(f.props)
}

func nextHook(f *Fiber, make func() any) (int, any) {
	i := hookIndex
	hookIndex++
	if i == len(f.hooks) {
		f.hooks = append(f.hooks, make())
	}
	return i, f.hooks[i]
}

func valueOrZero[T any](value any) T {
	if value == nil {
		var zero T
		return zero
	}
	return value.(T)
}

// UseState 声明组件本地状态，返回当前值和 setter。
// setter 捕获持久的 Fiber，调用即自动把该组件排入局部重渲染队列。
func UseState[T any](initial T) (T, func(T)) {
	f := currentFiber
	_, raw := nextHook(f, func() any { return &stateHook{value: initial} })
	h := raw.(*stateHook)
	if h.setter == nil {
		h.setter = func(nv T) {
			if f.unmounted || reflect.DeepEqual(h.value, nv) {
				return
			}
			h.value = nv
			if f.game != nil {
				f.game.markDirty(f)
			}
		}
	}
	return valueOrZero[T](h.value), h.setter.(func(T))
}

// UseEffect 在 commit 后执行副作用；deps 变化时重跑，返回的 Cleanup 在下次重跑前/卸载时调用。
func UseEffect(fn func() Cleanup, deps ...any) {
	f := currentFiber
	_, raw := nextHook(f, func() any { return &effectHook{} })
	h := raw.(*effectHook)
	if h.hasRun && depsEqual(h.deps, deps) {
		return
	}
	h.deps = deps
	h.hasRun = true
	if f.game != nil {
		f.game.pendingEffects = append(f.game.pendingEffects, func() {
			if f.unmounted {
				return
			}
			if h.cleanup != nil {
				h.cleanup()
			}
			h.cleanup = fn()
		})
	}
}

// UseReducer 以 reducer 管理状态，返回当前状态和稳定的 dispatch。
func UseReducer[S any, A any](reducer func(S, A) S, initial S) (S, func(A)) {
	f := currentFiber
	_, raw := nextHook(f, func() any { return &stateHook{value: initial} })
	h := raw.(*stateHook)
	h.reducer = reducer // 每次渲染刷新 reducer，dispatch 引用保持稳定
	if h.dispatch == nil {
		h.dispatch = func(a A) {
			if f.unmounted {
				return
			}
			next := h.reducer.(func(S, A) S)(valueOrZero[S](h.value), a)
			if reflect.DeepEqual(h.value, next) {
				return
			}
			h.value = next
			if f.game != nil {
				f.game.markDirty(f)
			}
		}
	}
	return valueOrZero[S](h.value), h.dispatch.(func(A))
}

type memoHook struct {
	deps        []any
	value       any
	initialized bool
}

// UseMemo 缓存计算结果，deps 变化时才重算。
func UseMemo[T any](create func() T, deps ...any) T {
	f := currentFiber
	_, raw := nextHook(f, func() any { return &memoHook{deps: nil, value: nil} })
	h := raw.(*memoHook)
	if !h.initialized || !depsEqual(h.deps, deps) {
		h.value = create()
		h.deps = deps
		h.initialized = true
	}
	return valueOrZero[T](h.value)
}

// UseCallback 缓存回调函数，deps 变化时才更新（引用稳定）。
func UseCallback[F any](fn F, deps ...any) F {
	return UseMemo(func() F { return fn }, deps...)
}

// UseRef 返回一个跨重渲染稳定的可变引用。
func UseRef[T any](initial T) *T {
	f := currentFiber
	_, raw := nextHook(f, func() any { v := initial; return &refHook{value: &v} })
	return raw.(*refHook).value.(*T)
}

func depsEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
