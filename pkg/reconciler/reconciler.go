package reconciler

import (
	"github.com/sjm1327605995/tenon/pkg/core"
)

var (
	rootFiber *Fiber
	isRendering bool
)

// RenderRoot 渲染根节点（主入口）
func RenderRoot(root core.Component) {
	if root == nil {
		return
	}

	isRendering = true
	defer func() { isRendering = false }()

	if rootFiber == nil {
		// 首次渲染
		rootFiber = BuildFiberTreeFromComponent(root, HostRoot, "root")
	} else {
		// 更新渲染：创建 workInProgress 树
		rootFiber = CreateWorkInProgress(rootFiber)
		rootFiber.StateNode = root
	}

	// 准备新的工作栈
	prepareFreshStack(rootFiber, DefaultUpdateLanes)

	// 执行工作循环
	workLoop()

	// 提交阶段
	commitRoot(rootFiber)
}

// ScheduleUpdateOnFiber 调度 Fiber 更新（由 hooks 调用）
func ScheduleUpdateOnFiber(fiber *Fiber, lane Lane) {
	if fiber == nil {
		return
	}

	// 标记车道
	fiber.Lanes = MergeLanes(fiber.Lanes, lane)

	// 冒泡到根节点
	parent := fiber.Return
	for parent != nil {
		parent.ChildLanes = MergeLanes(parent.ChildLanes, lane)
		parent = parent.Return
	}

	// 根节点收到更新标记
	if rootFiber != nil {
		rootFiber.Lanes = MergeLanes(rootFiber.Lanes, lane)
	}
}

// IsRendering 检查是否处于渲染阶段
func IsRendering() bool {
	return isRendering
}

// GetRootFiber 获取根 Fiber
func GetRootFiber() *Fiber {
	return rootFiber
}

// RebuildFiberTree 重建 Fiber 树（当组件树结构变化时）
func RebuildFiberTree(root core.Component) {
	if root == nil {
		return
	}

	// 清理旧的 effects
	if rootFiber != nil {
		CleanupFiberEffects(rootFiber)
	}

	rootFiber = BuildFiberTreeFromComponent(root, HostRoot, "root")
	RenderRoot(root)
}

// UpdateFiberFromComponent 根据组件更新 Fiber
func UpdateFiberFromComponent(fiber *Fiber, component core.Component) {
	if fiber == nil || component == nil {
		return
	}
	fiber.StateNode = component
	// 子节点需要重新协调
	fiber.Child = nil
	buildChildrenFibers(fiber, component)
}
