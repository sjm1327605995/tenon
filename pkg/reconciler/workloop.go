package reconciler

import (
	"time"
)

// 时间切片配置
const (
	frameYieldMs = 5 // 每帧留给 JS/go 的时间（Ebiten 是 60fps，约 16ms）
)

var (
	workInProgress       *Fiber
	workInProgressRoot   *Fiber
	currentFiberTree     *Fiber
	renderLanes          Lanes
	shouldYield          bool
	workInProgressHook   *HookState
	currentHook          *HookState
)

// SetWorkInProgress 设置当前工作的 Fiber
func SetWorkInProgress(fiber *Fiber) {
	workInProgress = fiber
}

// GetWorkInProgress 获取当前工作的 Fiber
func GetWorkInProgress() *Fiber {
	return workInProgress
}

// GetCurrentFiberTree 获取当前的 Fiber 树
func GetCurrentFiberTree() *Fiber {
	return currentFiberTree
}

// SetCurrentHook 设置当前 hooks 指针
func SetCurrentHook(hook *HookState) {
	currentHook = hook
}

// GetCurrentHook 获取当前 hooks 指针
func GetCurrentHook() *HookState {
	return currentHook
}

// GetWorkInProgressHook 获取工作进度中的 hook
func GetWorkInProgressHook() *HookState {
	return workInProgressHook
}

// ResetWorkInProgressHook 重置 hook 指针
func ResetWorkInProgressHook() {
	workInProgressHook = nil
	currentHook = nil
}

// workLoop 执行工作循环（支持时间切片）
func workLoop() {
	startTime := time.Now()
	for workInProgress != nil && !shouldYieldToHost(startTime) {
		performUnitOfWork(workInProgress)
	}
}

// performUnitOfWork 执行单个工作单元
func performUnitOfWork(unitOfWork *Fiber) {
	next := beginWork(unitOfWork)
	if next == nil {
		completeUnitOfWork(unitOfWork)
	} else {
		workInProgress = next
	}
}

// completeUnitOfWork 完成工作单元，冒泡 subtree flags
func completeUnitOfWork(unitOfWork *Fiber) {
	node := unitOfWork
	for {
		returnFiber := node.Return
		if returnFiber != nil {
			returnFiber.SubtreeFlags |= node.SubtreeFlags | node.Flags
		}

		if returnFiber == nil {
			workInProgress = nil
			return
		}

		sibling := node.Sibling
		if sibling != nil {
			workInProgress = sibling
			return
		}
		node = returnFiber
	}
}

// shouldYieldToHost 检查是否应该让出执行权
func shouldYieldToHost(startTime time.Time) bool {
	elapsed := time.Since(startTime)
	return elapsed > frameYieldMs*time.Millisecond
}

// prepareFreshStack 准备新的工作栈
func prepareFreshStack(root *Fiber, lanes Lanes) {
	workInProgressRoot = root
	workInProgress = root
	renderLanes = lanes
	shouldYield = false
	ResetWorkInProgressHook()
}
