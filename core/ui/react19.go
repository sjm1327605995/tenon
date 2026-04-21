package ui

import (
	"sync/atomic"
)

var isBatching = atomic.Value{}
var batchDepth int32

type Transition struct {
	_lanes Lane
	_status int32
}

const (
	TransitionPending int32 = 0
	TransitionResolved int32 = 1
	TransitionDiscarded int32 = 2
)

func StartTransition(fiber *Fiber, callback func()) *Transition {
	transition := &Transition{
		_lanes: TransitionLane1,
		_status: TransitionPending,
	}

	scheduler := GetCurrentScheduler()
	scheduler.scheduleTransitionTask(func() {
		callback()
		atomic.StoreInt32(&transition._status, TransitionResolved)
	})

	return transition
}

func UseTransition() (bool, func(func())) {
	isPending, dispatch := UseState(false)

	startTransition := func(callback func()) {
		dispatch(true)

		fiber := currentDispatcher.getCurrentFiber()
		if fiber != nil {
			transition := StartTransition(fiber, func() {
				callback()
				dispatch(false)
			})
			_ = transition
		} else {
			callback()
			dispatch(false)
		}
	}

	return isPending.(bool), startTransition
}

func UseDeferredValue[T any](value T, comparator func(old, new T) bool) T {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(value)
		return value
	}

	if hook.memoizedState == nil {
		hook.memoizedState = value
		return value
	}

	oldValue := hook.memoizedState.(T)
	if comparator != nil && comparator(oldValue, value) {
		return oldValue
	}

	scheduler := GetCurrentScheduler()
	scheduler.scheduleTransitionTask(func() {
		hook.memoizedState = value
	})

	return oldValue
}

func UseInsertionEffect(create EffectCallback, callback interface{}, dependencies []interface{}) {
	UseEffect(create, callback, dependencies)
}

func UseSyncExternalStore(subscribe func(func()) func(), getSnapshot func() interface{}, getServerSnapshot func() interface{}) interface{} {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(getSnapshot())
	}

	if hook.memoizedState == nil {
		hook.memoizedState = getSnapshot()
	}

	_ = subscribe

	return hook.memoizedState
}

func UseId() string {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(generateId())
	} else if hook.memoizedState == nil {
		hook.memoizedState = generateId()
	}

	return hook.memoizedState.(string)
}

var idCounter int64

func generateId() string {
	return "react-id-" + string(rune(atomic.AddInt64(&idCounter, 1)))
}

type ActionState int

const (
	ActionStateIdle ActionState = 0
	ActionStatePending ActionState = 1
	ActionStateResolved ActionState = 2
	ActionStateRejected ActionState = 3
)

type ActionResult[T any] struct {
	status ActionState
	value  T
	error  error
}

func UseActionState(reducer func(state interface{}, action interface{}) interface{}, initialState interface{}, action func() error) (interface{}, func()) {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(&ActionResult[interface{}]{
			status: ActionStateIdle,
			value:  initialState,
		})
	}

	state := hook.memoizedState.(*ActionResult[interface{}])

	executeAction := func() {
		dispatcher.setCurrentFiber(dispatcher.getCurrentFiber())
		err := action()
		if err != nil {
			state.status = ActionStateRejected
			state.error = err
		} else {
			state.status = ActionStateResolved
		}
	}

	return state.value, executeAction
}

func UseOptimistic(value interface{}) (interface{}, func(interface{})) {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(value)
	}

	dispatch := func(newValue interface{}) {
		hook.memoizedState = newValue
	}

	return hook.memoizedState, dispatch
}

func FlushSync(callback func()) {
	scheduler := GetCurrentScheduler()
	scheduler.scheduleSyncTask(callback)
}

func BatchedUpdates(callback func()) {
	isBatching.Store(int32(1))
	defer func() {
		isBatching.Store(int32(0))
	}()
	callback()
}

func UnbatchedUpdates(callback func()) {
	callback()
}

var (
	renderRoot *Fiber
	rootHostContext interface{}
)

func CreateRoot() *FiberRoot {
	return &FiberRoot{
		current: nil,
		pendingLanes: NoLane,
	}
}

type FiberRoot struct {
	current       *Fiber
	pendingLanes  Lane
	finishedWork  *Fiber
	callbackQueue []func()
}

func (root *FiberRoot) Render() {
	scheduler := GetCurrentScheduler()
	scheduler.scheduleSyncTask(func() {
		workLoop(root)
	})
}

func workLoop(root *FiberRoot) {
	for root.current != nil {
		performUnitOfWork(root.current)
		root.current = root.current.alternate
	}
}

func performUnitOfWork(fiber *Fiber) {
	if fiber.type_ != nil {
		if fn, ok := fiber.type_.(func() Element); ok {
			element := fn()
			if element != nil {
				currentDispatcher.setCurrentFiber(fiber)
			}
		}
	}

	if fiber.child != nil {
		performUnitOfWork(fiber.child)
	}
}


