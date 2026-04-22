package hooks

import (
	"reflect"

	"github.com/sjm1327605995/tenon/pkg/reconciler"
)

// UseEffect 类似 React 的 useEffect Hook
// 在 commit 阶段的 passive 阶段异步执行
func UseEffect(effect func() func(), deps []interface{}) {
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

	// 检查依赖是否变化
	if hasChanged(prevDeps, deps) {
		// 记录 effect，在 commit 阶段执行
		effectHook := &reconciler.Effect{
			Tag:    reconciler.PassiveEffect,
			Create: effect,
			Destroy: prevDestroy,
			Deps:   deps,
		}

		if hook.Effect == nil {
			hook.Effect = effectHook
		} else {
			// 追加到 effect 链表
			e := hook.Effect
			for e.Next != nil {
				e = e.Next
			}
			e.Next = effectHook
		}

		// 标记 fiber 有 passive effect
		if currentFiber != nil {
			currentFiber.Flags |= reconciler.Passive
			currentFiber.SubtreeFlags |= reconciler.Passive
		}

		hook.MemoizedState = map[string]interface{}{
			"deps":    deps,
			"destroy": nil,
		}
	}
}

// hasChanged 检查依赖数组是否变化
func hasChanged(prevDeps, nextDeps []interface{}) bool {
	if prevDeps == nil {
		return true
	}
	if len(prevDeps) != len(nextDeps) {
		return true
	}
	for i := range prevDeps {
		if !reflect.DeepEqual(prevDeps[i], nextDeps[i]) {
			return true
		}
	}
	return false
}
