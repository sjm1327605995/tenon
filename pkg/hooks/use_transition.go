package hooks

import (
	"time"

	"github.com/sjm1327605995/tenon/pkg/reconciler"
)

// TransitionState 过渡状态
type TransitionState struct {
	IsPending bool
}

// UseTransition 类似 React 的 useTransition Hook
// 返回 [isPending, startTransition]
func UseTransition() (bool, func(func())) {
	hook := getWorkInProgressHook()

	var state *TransitionState
	if hook.MemoizedState == nil {
		state = &TransitionState{IsPending: false}
		hook.MemoizedState = state
	} else {
		state = hook.MemoizedState.(*TransitionState)
	}

	startTransition := func(callback func()) {
		state.IsPending = true

		// 使用 TransitionLane 调度更新
		if currentFiber != nil {
			update := &reconciler.Update{
				Action: func(prevState, props interface{}) interface{} {
					callback()
					return prevState
				},
				Lane: reconciler.TransitionLane1,
			}
			reconciler.EnqueueUpdate(currentFiber, update)
			reconciler.ScheduleUpdateOnFiber(currentFiber, reconciler.TransitionLane1)
		}

		// 模拟异步完成（实际应在 commit 完成后设置）
		go func() {
			time.Sleep(16 * time.Millisecond)
			state.IsPending = false
		}()
	}

	return state.IsPending, startTransition
}

// useDeferredValue 延迟更新值（降低优先级）
func UseDeferredValue[T any](value T) T {
	deferredValue, _ := UseState(value)

	UseEffect(func() func() {
		// 延迟更新
		return nil
	}, []interface{}{value})

	return deferredValue
}
