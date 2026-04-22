package reconciler

import (
	"github.com/sjm1327605995/tenon/pkg/core"
)

// commitRoot 提交根节点
func commitRoot(root *Fiber) {
	if root == nil {
		return
	}

	// 阶段 1: BeforeMutation（获取快照等）
	commitBeforeMutationEffects(root)

	// 阶段 2: Mutation（对组件树进行变更）
	commitMutationEffects(root)

	// 交换 current 和 workInProgress
	if root.Alternate != nil {
		currentFiberTree = root
	}

	// 阶段 3: Layout（useLayoutEffect, 挂载回调）
	commitLayoutEffects(root)

	// 阶段 4: Passive（useEffect，异步调度）
	schedulePassiveEffects(root)
}

// commitBeforeMutationEffects 提交前副作用
func commitBeforeMutationEffects(finishedWork *Fiber) {
	// 遍历 subtree，执行 getSnapshotBeforeUpdate 等
	// 简化版：目前主要用于类组件，函数组件不需要
}

// commitMutationEffects 变更副作用
func commitMutationEffects(finishedWork *Fiber) {
	traverseFiberTree(finishedWork, func(fiber *Fiber) {
		if fiber.Flags&Placement != 0 {
			commitPlacement(fiber)
			fiber.Flags &^= Placement
		}
		if fiber.Flags&Update != 0 {
			commitUpdate(fiber)
			fiber.Flags &^= Update
		}
		if fiber.Flags&ChildDeletion != 0 {
			// 处理子节点删除
			fiber.Flags &^= ChildDeletion
		}
	})
}

// commitPlacement 处理挂载
func commitPlacement(fiber *Fiber) {
	// 不需要 DOM 操作，组件已经在树中
	// 标记组件需要重新计算布局
	if fiber.StateNode != nil {
		fiber.StateNode.MarkDirty()
	}
}

// commitUpdate 处理更新
func commitUpdate(fiber *Fiber) {
	if fiber.StateNode != nil {
		fiber.StateNode.MarkDirty()
	}
}

// commitLayoutEffects 提交 Layout 副作用（同步执行）
func commitLayoutEffects(finishedWork *Fiber) {
	traverseFiberTree(finishedWork, func(fiber *Fiber) {
		if fiber.Tag == FunctionComponent {
			commitHookLayoutEffects(fiber, LayoutEffect)
		}
		// 类组件的 componentDidMount/Update
		if fiber.StateNode != nil {
			// 通知组件更新完成
		}
	})
}

// commitHookLayoutEffects 提交 hooks 的 layout effects
func commitHookLayoutEffects(finishedWork *Fiber, effectTag EffectTag) {
	effect := finishedWork.MemoizedState
	for effect != nil && effect.Effect != nil {
		if effect.Effect.Tag&effectTag != 0 {
			if effect.Effect.Destroy != nil {
				effect.Effect.Destroy()
			}
			if effect.Effect.Create != nil {
				cleanup := effect.Effect.Create()
				if cleanup != nil {
					effect.Effect.Destroy = cleanup
				}
			}
		}
		effect = effect.Next
	}
}

// schedulePassiveEffects 调度 Passive 副作用（useEffect，异步）
var passiveEffectsQueue []*Fiber

func schedulePassiveEffects(finishedWork *Fiber) {
	// 收集所有有 Passive 标记的 fiber
	traverseFiberTree(finishedWork, func(fiber *Fiber) {
		if fiber.SubtreeFlags&Passive != 0 || fiber.Flags&Passive != 0 {
			passiveEffectsQueue = append(passiveEffectsQueue, fiber)
		}
	})
}

// FlushPassiveEffects 执行 Passive 副作用（由 renderer 调用）
func FlushPassiveEffects() {
	if len(passiveEffectsQueue) == 0 {
		return
	}
	queue := passiveEffectsQueue
	passiveEffectsQueue = nil

	for _, fiber := range queue {
		commitHookPassiveEffects(fiber)
	}
}

// commitHookPassiveEffects 提交 useEffect
func commitHookPassiveEffects(finishedWork *Fiber) {
	effect := finishedWork.MemoizedState
	for effect != nil && effect.Effect != nil {
		if effect.Effect.Tag&PassiveEffect != 0 {
			if effect.Effect.Destroy != nil {
				effect.Effect.Destroy()
			}
			if effect.Effect.Create != nil {
				cleanup := effect.Effect.Create()
				if cleanup != nil {
					effect.Effect.Destroy = cleanup
				}
			}
		}
		effect = effect.Next
	}
}

// traverseFiberTree 遍历 Fiber 树
func traverseFiberTree(root *Fiber, callback func(*Fiber)) {
	if root == nil {
		return
	}
	callback(root)

	// 先序遍历
	child := root.Child
	for child != nil {
		traverseFiberTree(child, callback)
		child = child.Sibling
	}
}

// CleanupFiberEffects 清理 Fiber 的所有副作用
func CleanupFiberEffects(fiber *Fiber) {
	if fiber == nil {
		return
	}

	// 清理 hooks 副作用
	effect := fiber.MemoizedState
	for effect != nil && effect.Effect != nil {
		if effect.Effect.Destroy != nil {
			effect.Effect.Destroy()
		}
		effect = effect.Next
	}

	// 递归清理子节点
	child := fiber.Child
	for child != nil {
		CleanupFiberEffects(child)
		child = child.Sibling
	}
}

// BuildFiberTreeFromComponent 从组件根构建 Fiber 树（初始化用）
func BuildFiberTreeFromComponent(root core.Component, tag FiberTag, typeName string) *Fiber {
	rootFiber := CreateFiber(tag, typeName, "")
	rootFiber.StateNode = root
	rootFiber.Depth = 0

	buildChildrenFibers(rootFiber, root)
	return rootFiber
}

func buildChildrenFibers(parentFiber *Fiber, parentComponent core.Component) {
	if parentComponent == nil {
		return
	}

	children := parentComponent.GetChildren()
	var firstChild *Fiber
	var prevSibling *Fiber

	for _, child := range children {
		childFiber := CreateFiber(HostComponent, getComponentTypeName(child), "")
		childFiber.StateNode = child
		childFiber.Return = parentFiber
		childFiber.Depth = parentFiber.Depth + 1

		if firstChild == nil {
			firstChild = childFiber
			parentFiber.Child = childFiber
		} else {
			prevSibling.Sibling = childFiber
		}
		prevSibling = childFiber

		// 递归构建孙节点
		buildChildrenFibers(childFiber, child)
	}
}

func getComponentTypeName(c core.Component) string {
	if c == nil {
		return ""
	}
	// 简单方案：通过类型断言判断
	switch c.(type) {
	case interface{ SetContent(string) *core.BaseComponent }:
		return "text"
	case interface{ SetOnClick(func()) *core.BaseComponent }:
		return "button"
	default:
		return "view"
	}
}
