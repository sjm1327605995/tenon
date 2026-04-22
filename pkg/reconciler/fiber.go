package reconciler

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sjm1327605995/tenon/pkg/core"
)

// FiberTag 标识 Fiber 类型
type FiberTag int

const (
	HostRoot FiberTag = iota
	HostComponent
	FunctionComponent
)

// SideEffectFlags 副作用标记
type SideEffectFlags int

const (
	NoFlags           SideEffectFlags = 0
	Placement                         = 1 << 0
	Update                            = 1 << 1
	ChildDeletion                     = 1 << 2
	Passive                           = 1 << 3 // useEffect
	LayoutMask                        = 1 << 4 // useLayoutEffect
	InsertionMask                     = 1 << 5 // useInsertionEffect
)

// Fiber 是调度的基本单元，每个组件对应一个 Fiber
// 不需要 VDOM，Fiber 直接关联到组件实例
type Fiber struct {
	Tag  FiberTag
	Type string
	Key  string

	// 关联的组件实例（HostComponent 时有效）
	StateNode core.Component

	// 函数组件专用
	RenderFn   func(props map[string]interface{}) core.Component
	RenderProps map[string]interface{}

	// 树结构
	Return  *Fiber // parent
	Child   *Fiber // first child
	Sibling *Fiber // next sibling

	// Hooks 状态链表头
	MemoizedState *HookState

	// 更新队列
	UpdateQueue *UpdateQueue

	// 优先级车道
	Lanes      Lane
	ChildLanes Lane

	// 副作用标记（冒泡）
	Flags        SideEffectFlags
	SubtreeFlags SideEffectFlags

	// 双缓冲
	Alternate *Fiber

	// 深度
	Depth int

	// ID
	ID string
}

// HookState hooks 链表节点
type HookState struct {
	MemoizedState interface{}
	BaseState     interface{}
	Queue         *UpdateQueue
	Next          *HookState
	Effect        *Effect
}

// Effect 副作用存储
type Effect struct {
	Tag       EffectTag
	Create    func() func()
	Destroy   func()
	Deps      []interface{}
	Next      *Effect
}

// EffectTag 副作用类型
type EffectTag int

const (
	PassiveEffect EffectTag = 1 << 0
	LayoutEffect  EffectTag = 1 << 1
	InsertionEffect EffectTag = 1 << 2
)

// UpdateQueue 更新队列
type UpdateQueue struct {
	BaseState   interface{}
	FirstUpdate *Update
	LastUpdate  *Update
}

// Update 更新节点
type Update struct {
	Action        interface{}
	Next          *Update
	Lane          Lane
	HasEagerState bool
	EagerState    interface{}
}

var (
	fiberIDCounter int64
	fiberIDMutex   sync.Mutex
)

func generateFiberID() string {
	fiberIDMutex.Lock()
	defer fiberIDMutex.Unlock()
	fiberIDCounter++
	return fmt.Sprintf("fiber_%d_%d", time.Now().UnixNano(), fiberIDCounter)
}

// CreateFiber 创建新 Fiber
func CreateFiber(tag FiberTag, typeName string, key string) *Fiber {
	return &Fiber{
		Tag:   tag,
		Type:  typeName,
		Key:   key,
		ID:    generateFiberID(),
		Lanes: NoLane,
	}
}

// CreateWorkInProgress 基于 current 创建工作中的 Fiber
func CreateWorkInProgress(current *Fiber) *Fiber {
	workInProgress := &Fiber{
		Tag:          current.Tag,
		Type:         current.Type,
		Key:          current.Key,
		StateNode:    current.StateNode,
		RenderFn:     current.RenderFn,
		RenderProps:  current.RenderProps,
		Return:       current.Return,
		Depth:        current.Depth,
		ID:           current.ID,
		Lanes:        current.Lanes,
		ChildLanes:   current.ChildLanes,
		MemoizedState: copyHookState(current.MemoizedState),
		Alternate:    current,
	}
	current.Alternate = workInProgress
	return workInProgress
}

func copyHookState(head *HookState) *HookState {
	if head == nil {
		return nil
	}
	newHead := &HookState{
		MemoizedState: head.MemoizedState,
		BaseState:     head.BaseState,
	}
	current := newHead
	for src := head.Next; src != nil; src = src.Next {
		newHook := &HookState{
			MemoizedState: src.MemoizedState,
			BaseState:     src.BaseState,
		}
		current.Next = newHook
		current = newHook
	}
	return newHead
}

// ResetFiber 重置 Fiber 状态（用于复用）
func ResetFiber(fiber *Fiber) {
	fiber.Flags = NoFlags
	fiber.SubtreeFlags = NoFlags
	fiber.UpdateQueue = nil
}

// EnqueueUpdate 向 Fiber 添加更新
func EnqueueUpdate(fiber *Fiber, update *Update) {
	queue := fiber.UpdateQueue
	if queue == nil {
		queue = &UpdateQueue{}
		fiber.UpdateQueue = queue
	}

	if queue.LastUpdate == nil {
		queue.FirstUpdate = update
		queue.LastUpdate = update
	} else {
		queue.LastUpdate.Next = update
		queue.LastUpdate = update
	}
}

// ProcessUpdateQueue 处理更新队列
func ProcessUpdateQueue(workInProgress *Fiber, props interface{}) {
	queue := workInProgress.UpdateQueue
	if queue == nil || queue.FirstUpdate == nil {
		return
	}

	newState := queue.BaseState
	update := queue.FirstUpdate
	for update != nil {
		if actionFn, ok := update.Action.(func(interface{}, interface{}) interface{}); ok {
			newState = actionFn(newState, props)
		} else {
			newState = update.Action
		}
		update = update.Next
	}

	queue.BaseState = newState
	if workInProgress.MemoizedState == nil {
		workInProgress.MemoizedState = &HookState{}
	}
	workInProgress.MemoizedState.BaseState = newState
	workInProgress.MemoizedState.MemoizedState = newState
	queue.FirstUpdate = nil
	queue.LastUpdate = nil
}

// ResetHookStateIndex 用于 hooks 遍历
func ResetHookStateIndex(fiber *Fiber) {
	// hooks 通过 currentHook 指针遍历，不需要 index
}
