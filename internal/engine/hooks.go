package engine

import (
	"reflect"
)

// ==================== Hooks 核心 ====================

// Hooks 为函数组件提供状态管理能力。
// 每次 render 时 hooks 按调用顺序匹配状态槽位。
type Hooks struct {
	element     *fcElement
	states      []any
	effects     []effectHook
	stateIndex  int
	effectIndex int
	renderCount int
}

// reset 在每次 render 前重置索引。
func (h *Hooks) reset() {
	h.stateIndex = 0
	h.effectIndex = 0
	h.renderCount++
}

// ==================== useState ====================

// UseState 返回当前值的 getter 和 setter。
// 使用 getter 而非直接返回值，可确保闭包内始终读到最新状态。
func UseState[T any](h *Hooks, initial T) (func() T, func(T)) {
	idx := h.stateIndex
	h.stateIndex++

	if idx >= len(h.states) {
		h.states = append(h.states, initial)
	}

	get := func() T {
		return h.states[idx].(T)
	}

	set := func(v T) {
		if reflect.DeepEqual(h.states[idx], v) {
			return
		}
		h.states[idx] = v
		if h.element != nil {
			h.element.markNeedsBuild()
		}
	}

	return get, set
}

// ==================== useEffect ====================

type effectHook struct {
	cleanup func()
	deps    []any
}

// UseEffect 在依赖变化时执行副作用。
// effect 返回的 cleanup 函数在依赖变化或组件卸载时调用。
func UseEffect(h *Hooks, effect func() func(), deps []any) {
	idx := h.effectIndex
	h.effectIndex++

	if idx >= len(h.effects) {
		// 首次渲染
		h.effects = append(h.effects, effectHook{
			deps: append([]any(nil), deps...),
		})
		if cleanup := effect(); cleanup != nil {
			h.effects[idx].cleanup = cleanup
		}
		return
	}

	old := h.effects[idx]
	if !depsEqual(old.deps, deps) {
		if old.cleanup != nil {
			old.cleanup()
		}
		h.effects[idx].deps = append([]any(nil), deps...)
		h.effects[idx].cleanup = nil
		if cleanup := effect(); cleanup != nil {
			h.effects[idx].cleanup = cleanup
		}
	}
}

// ==================== useMemo ====================

type memoSlot struct {
	value any
	deps  []any
}

// UseMemo 缓存 factory 的计算结果，仅在依赖变化时重新计算。
func UseMemo[T any](h *Hooks, factory func() T, deps []any) T {
	idx := h.stateIndex
	h.stateIndex++

	if idx >= len(h.states) {
		val := factory()
		h.states = append(h.states, memoSlot{value: val, deps: append([]any(nil), deps...)})
		return val
	}

	slot := h.states[idx].(memoSlot)
	if !depsEqual(slot.deps, deps) {
		slot.value = factory()
		slot.deps = append([]any(nil), deps...)
		h.states[idx] = slot
	}

	return slot.value.(T)
}

// ==================== useCallback ====================

// UseCallback 缓存函数引用，仅在依赖变化时更新。
// 等价于 UseMemo(h, func() T { return fn }, deps)。
func UseCallback[T any](h *Hooks, fn T, deps []any) T {
	return UseMemo(h, func() T { return fn }, deps)
}

// ==================== 辅助函数 ====================

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
