package hooks

import (
	"github.com/sjm1327605995/tenon/pkg/reconciler"
)

// ActionState 动作状态
type ActionState[T any] struct {
	Value    T
	Pending  bool
	Error    error
}

// UseActionState React 19 的 useActionState Hook
// 用于处理表单动作和异步状态
func UseActionState[T any](
	action func(prevState T, formData map[string]string) (T, error),
	initialValue T,
) (T, func(map[string]string), bool, error) {
	hook := getWorkInProgressHook()

	var state *ActionState[T]
	if hook.MemoizedState == nil {
		state = &ActionState[T]{
			Value:   initialValue,
			Pending: false,
		}
		hook.MemoizedState = state
	} else {
		state = hook.MemoizedState.(*ActionState[T])
	}

	dispatch := func(formData map[string]string) {
		state.Pending = true

		// 异步执行 action
		go func() {
			newValue, err := action(state.Value, formData)

			// 使用 TransitionLane 调度结果更新
			if currentFiber != nil {
				update := &reconciler.Update{
					Action: func(prevState, props interface{}) interface{} {
						return &ActionState[T]{
							Value:   newValue,
							Pending: false,
							Error:   err,
						}
					},
					Lane: reconciler.TransitionLane1,
				}
				reconciler.EnqueueUpdate(currentFiber, update)
				reconciler.ScheduleUpdateOnFiber(currentFiber, reconciler.TransitionLane1)
			}
		}()
	}

	return state.Value, dispatch, state.Pending, state.Error
}
