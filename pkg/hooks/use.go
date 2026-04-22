package hooks

import (
	"fmt"
	"reflect"
)

// Promise 可等待的 Promise 接口
type Promise[T any] interface {
	Then(onFulfilled func(T), onRejected func(error))
}

// Use React 19 的新 API
// 可以在 render 中读取 Promise 或 Context
// 注意：实际 React 中 use 是可以在条件语句中调用的特殊 Hook

// UsePromise 在 render 中读取 Promise（需要配合 Suspense）
func UsePromise[T any](promise Promise[T]) T {
	hook := getWorkInProgressHook()

	if hook.MemoizedState == nil {
		// 首次调用，设置 pending 状态
		pendingState := &promisePendingState[T]{}
		hook.MemoizedState = pendingState

		promise.Then(
			func(value T) {
				pendingState.value = value
				pendingState.resolved = true
				if currentFiber != nil {
					reconciler.ScheduleUpdateOnFiber(currentFiber, reconciler.DefaultLane)
				}
			},
			func(err error) {
				pendingState.err = err
				pendingState.rejected = true
				if currentFiber != nil {
					reconciler.ScheduleUpdateOnFiber(currentFiber, reconciler.DefaultLane)
				}
			},
		)
	}

	state := hook.MemoizedState.(*promisePendingState[T])
	if state.rejected {
		panic(state.err)
	}
	if !state.resolved {
		// 在真实 React 中会抛出 Promise 让 Suspense 捕获
		panic("Promise not resolved yet")
	}
	return state.value
}

type promisePendingState[T any] struct {
	value    T
	resolved bool
	err      error
	rejected bool
}

// UseResource 读取可消费资源（类似 use，但更通用）
func UseResource[T any](resource interface{}) T {
	// 检查是否是 Context
	if ctx, ok := resource.(Context[T]); ok {
		return UseContext(ctx)
	}

	// 检查是否是 Promise
	if promise, ok := resource.(Promise[T]); ok {
		return UsePromise(promise)
	}

	// 直接返回
	if v, ok := resource.(T); ok {
		return v
	}

	panic(fmt.Sprintf("Unsupported resource type: %v", reflect.TypeOf(resource)))
}
