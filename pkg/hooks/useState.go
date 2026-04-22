package hooks

import (
	"reflect"
)

// UseState 类似 React 的 useState Hook
// 返回当前值和 setter 函数
func UseState[T any](initialValue T) (T, func(T)) {
	hook := getWorkInProgressHook()

	// 初始化状态
	if hook.MemoizedState == nil {
		hook.BaseState = initialValue
		hook.MemoizedState = initialValue
	}

	currentValue, ok := hook.MemoizedState.(T)
	if !ok {
		currentValue = initialValue
	}

	setter := func(newValue T) {
		if !reflect.DeepEqual(hook.MemoizedState, newValue) {
			hook.BaseState = newValue
			hook.MemoizedState = newValue

			if currentFiber != nil {
				update := &reconciler.Update{
					Action: newValue,
					Lane:   reconciler.DefaultLane,
				}
				reconciler.EnqueueUpdate(currentFiber, update)
				reconciler.ScheduleUpdateOnFiber(currentFiber, reconciler.DefaultLane)
			}
		}
	}

	return currentValue, setter
}

// UseReducer 类似 React 的 useReducer
func UseReducer[T any, A any](reducer func(T, A) T, initialValue T) (T, func(A)) {
	hook := getWorkInProgressHook()

	if hook.MemoizedState == nil {
		hook.BaseState = initialValue
		hook.MemoizedState = initialValue
	}

	currentValue, ok := hook.MemoizedState.(T)
	if !ok {
		currentValue = initialValue
	}

	dispatch := func(action A) {
		newValue := reducer(currentValue, action)
		hook.BaseState = newValue
		hook.MemoizedState = newValue

		if currentFiber != nil {
			update := &reconciler.Update{
				Action: newValue,
				Lane:   reconciler.DefaultLane,
			}
			reconciler.EnqueueUpdate(currentFiber, update)
			reconciler.ScheduleUpdateOnFiber(currentFiber, reconciler.DefaultLane)
		}
	}

	return currentValue, dispatch
}
