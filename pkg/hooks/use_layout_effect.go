package hooks

import (
	"github.com/sjm1327605995/tenon/pkg/reconciler"
)

// UseLayoutEffect 类似 React 的 useLayoutEffect Hook
// 在 commit 阶段的 layout 阶段同步执行（在绘制前）
func UseLayoutEffect(effect func() func(), deps []interface{}) {
	hook := getWorkInProgressHook()

	var prevDeps []interface{}
	var prevDestroy func()
	if hook.MemoizedState != nil {
		state := hook.MemoizedState.(map[string]interface{})
		prevDeps = state["deps"].([]interface{})
		if d, ok := state["destroy"]; ok && d != nil {
			prevDestroy = d.(func())
		}
	}

	if hasChanged(prevDeps, deps) {
		// 先执行上一次的 cleanup
		if prevDestroy != nil {
			prevDestroy()
		}

		// 立即执行 effect（layout effect 是同步的）
		destroy := effect()

		effectHook := &reconciler.Effect{
			Tag:     reconciler.LayoutEffect,
			Create:  effect,
			Destroy: destroy,
			Deps:    deps,
		}

		if hook.Effect == nil {
			hook.Effect = effectHook
		} else {
			e := hook.Effect
			for e.Next != nil {
				e = e.Next
			}
			e.Next = effectHook
		}

		if currentFiber != nil {
			currentFiber.Flags |= reconciler.LayoutMask
			currentFiber.SubtreeFlags |= reconciler.LayoutMask
		}

		hook.MemoizedState = map[string]interface{}{
			"deps":    deps,
			"destroy": destroy,
		}
	}
}

// UseInsertionEffect 类似 React 18+ 的 useInsertionEffect
// 在所有 DOM 变更之前同步执行，用于注入样式
func UseInsertionEffect(effect func() func(), deps []interface{}) {
	hook := getWorkInProgressHook()

	var prevDeps []interface{}
	var prevDestroy func()
	if hook.MemoizedState != nil {
		state := hook.MemoizedState.(map[string]interface{})
		prevDeps = state["deps"].([]interface{})
		if d, ok := state["destroy"]; ok && d != nil {
			prevDestroy = d.(func())
		}
	}

	if hasChanged(prevDeps, deps) {
		if prevDestroy != nil {
			prevDestroy()
		}

		destroy := effect()

		effectHook := &reconciler.Effect{
			Tag:     reconciler.InsertionEffect,
			Create:  effect,
			Destroy: destroy,
			Deps:    deps,
		}

		if hook.Effect == nil {
			hook.Effect = effectHook
		} else {
			e := hook.Effect
			for e.Next != nil {
				e = e.Next
			}
			e.Next = effectHook
		}

		if currentFiber != nil {
			currentFiber.Flags |= reconciler.InsertionMask
			currentFiber.SubtreeFlags |= reconciler.InsertionMask
		}

		hook.MemoizedState = map[string]interface{}{
			"deps":    deps,
			"destroy": destroy,
		}
	}
}
