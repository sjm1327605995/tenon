package hooks

import (
	"reflect"

	"github.com/sjm1327605995/tenon/pkg/reconciler"
)

// optimisticState 乐观更新状态
type optimisticState[T any] struct {
	Value        T
	OptimisticValue T
	IsOptimistic bool
}

// UseOptimistic React 19 的 useOptimistic Hook
// 在执行异步操作时立即更新 UI，异步完成后再确认或回滚
func UseOptimistic[T any](
	state T,
	updateFn func(currentState T, optimisticValue interface{}) T,
) (T, func(interface{})) {
	hook := getWorkInProgressHook()

	var optState *optimisticState[T]
	if hook.MemoizedState == nil {
		optState = &optimisticState[T]{
			Value:        state,
			OptimisticValue: state,
			IsOptimistic: false,
		}
		hook.MemoizedState = optState
	} else {
		optState = hook.MemoizedState.(*optimisticState[T])
		// 如果外部 state 变化，重置乐观状态
		if !reflect.DeepEqual(optState.Value, state) {
			optState.Value = state
			optState.OptimisticValue = state
			optState.IsOptimistic = false
		}
	}

	addOptimistic := func(optimisticValue interface{}) {
		optState.IsOptimistic = true
		optState.OptimisticValue = updateFn(optState.Value, optimisticValue)

		// 触发低优先级更新
		if currentFiber != nil {
			update := &reconciler.Update{
				Action: optState.OptimisticValue,
				Lane:   reconciler.TransitionLane2,
			}
			reconciler.EnqueueUpdate(currentFiber, update)
			reconciler.ScheduleUpdateOnFiber(currentFiber, reconciler.TransitionLane2)
		}
	}

	if optState.IsOptimistic {
		return optState.OptimisticValue, addOptimistic
	}
	return optState.Value, addOptimistic
}
