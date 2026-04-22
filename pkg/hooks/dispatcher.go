package hooks

import (
	"github.com/sjm1327605995/tenon/pkg/reconciler"
)

var (
	currentFiber *reconciler.Fiber
	hookIndex    int
)

// SetCurrentFiber 设置当前正在渲染的 Fiber（reconciler 调用）
func SetCurrentFiber(fiber *reconciler.Fiber) {
	currentFiber = fiber
	hookIndex = 0
}

// GetCurrentFiber 获取当前 Fiber
func GetCurrentFiber() *reconciler.Fiber {
	return currentFiber
}

// getWorkInProgressHook 获取或创建当前 hook
func getWorkInProgressHook() *reconciler.HookState {
	if currentFiber == nil {
		panic("Hooks can only be called inside a function component's render function")
	}

	// 获取 hooks 链表
	head := currentFiber.MemoizedState

	// 按索引查找/创建
	var currentHook *reconciler.HookState
	if head == nil {
		// 第一个 hook
		currentHook = &reconciler.HookState{}
		currentFiber.MemoizedState = currentHook
	} else {
		currentHook = head
		for i := 0; i < hookIndex; i++ {
			if currentHook.Next == nil {
				currentHook.Next = &reconciler.HookState{}
			}
			currentHook = currentHook.Next
		}
	}

	hookIndex++
	return currentHook
}

// dispatchSetState 创建 setState dispatch 函数
func dispatchSetState(fiber *reconciler.Fiber, hook *reconciler.HookState) func(interface{}) {
	return func(action interface{}) {
		// 创建更新
		update := &reconciler.Update{
			Action: action,
			Lane:   reconciler.DefaultLane,
		}

		// 检查是否是函数更新
		if fn, ok := action.(func(interface{}, interface{}) interface{}); ok {
			// 如果是函数，直接计算新状态
			newState := fn(hook.BaseState, fiber.RenderProps)
			update.HasEagerState = true
			update.EagerState = newState
			hook.BaseState = newState
			hook.MemoizedState = newState
		}

		reconciler.EnqueueUpdate(fiber, update)
		reconciler.ScheduleUpdateOnFiber(fiber, reconciler.DefaultLane)
	}
}

// dispatchOptimisticSetState 乐观更新 dispatch
func dispatchOptimisticSetState(fiber *reconciler.Fiber, hook *reconciler.HookState) func(interface{}) {
	return func(action interface{}) {
		update := &reconciler.Update{
			Action: action,
			Lane:   reconciler.TransitionLane1,
		}
		reconciler.EnqueueUpdate(fiber, update)
		reconciler.ScheduleUpdateOnFiber(fiber, reconciler.TransitionLane1)
	}
}
