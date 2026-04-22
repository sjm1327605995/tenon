package reconciler

import (
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/hooks"
)

// beginWork 开始处理一个 Fiber
func beginWork(workInProgress *Fiber) *Fiber {
	switch workInProgress.Tag {
	case HostRoot:
		return updateHostRoot(workInProgress)
	case HostComponent:
		return updateHostComponent(workInProgress)
	case FunctionComponent:
		return updateFunctionComponent(workInProgress)
	default:
		return nil
	}
}

// updateHostRoot 更新根节点
func updateHostRoot(current *Fiber) *Fiber {
	// 处理更新队列
	ProcessUpdateQueue(current, nil)

	// 协调子节点
	if current.StateNode != nil {
		return reconcileChildren(current, current.StateNode.GetChildren())
	}
	return nil
}

// updateHostComponent 更新宿主组件
func updateHostComponent(current *Fiber) *Fiber {
	// 宿主组件自身不需要执行函数，只需协调子节点
	if current.StateNode != nil {
		return reconcileChildren(current, current.StateNode.GetChildren())
	}
	return nil
}

// updateFunctionComponent 更新函数组件
func updateFunctionComponent(current *Fiber) *Fiber {
	if current.RenderFn == nil {
		return nil
	}

	// 设置当前 Fiber，让 hooks 能绑定
	hooks.SetCurrentFiber(current)
	SetWorkInProgress(current)
	ResetWorkInProgressHook()

	// 处理 props 更新
	props := current.RenderProps
	if current.PendingProps != nil {
		props = current.PendingProps.(map[string]interface{})
		current.RenderProps = props
		current.PendingProps = nil
	}

	// 执行函数组件
	result := current.RenderFn(props)

	// 重置 hooks 上下文
	hooks.SetCurrentFiber(nil)
	ResetWorkInProgressHook()

	// 如果函数返回一个组件，把它作为子节点
	if result != nil {
		return reconcileSingleChild(current, result)
	}

	return nil
}

// reconcileChildren 协调子节点列表
func reconcileChildren(returnFiber *Fiber, children []core.Component) *Fiber {
	if len(children) == 0 {
		returnFiber.Child = nil
		return nil
	}

	var firstChild *Fiber
	var prevSibling *Fiber
	existingChild := returnFiber.Child

	for i, child := range children {
		key := ""
		if child.GetElement() != nil && child.GetElement().Key != "" {
			key = child.GetElement().Key
		} else {
			key = getComponentTypeName(child) + "_" + string(rune(i))
		}

		// 尝试复用已有 Fiber
		var childFiber *Fiber
		if existingChild != nil && existingChild.Key == key {
			childFiber = CreateWorkInProgress(existingChild)
			childFiber.StateNode = child
			existingChild = existingChild.Sibling
		} else {
			childFiber = CreateFiber(HostComponent, getComponentTypeName(child), key)
			childFiber.StateNode = child
			childFiber.Flags |= Placement
		}

		childFiber.Return = returnFiber
		childFiber.Depth = returnFiber.Depth + 1

		if firstChild == nil {
			firstChild = childFiber
			returnFiber.Child = childFiber
		} else {
			prevSibling.Sibling = childFiber
		}
		prevSibling = childFiber
	}

	// 标记多余子节点为删除
	for existingChild != nil {
		existingChild.Flags |= Deletion
		existingChild = existingChild.Sibling
	}

	return firstChild
}

// reconcileSingleChild 协调单个子节点
func reconcileSingleChild(returnFiber *Fiber, child core.Component) *Fiber {
	existingChild := returnFiber.Child
	key := ""
	if child.GetElement() != nil && child.GetElement().Key != "" {
		key = child.GetElement().Key
	}

	var childFiber *Fiber
	if existingChild != nil && existingChild.Key == key {
		childFiber = CreateWorkInProgress(existingChild)
		childFiber.StateNode = child
	} else {
		childFiber = CreateFiber(HostComponent, getComponentTypeName(child), key)
		childFiber.StateNode = child
		childFiber.Flags |= Placement
	}

	childFiber.Return = returnFiber
	childFiber.Depth = returnFiber.Depth + 1
	returnFiber.Child = childFiber
	return childFiber
}
