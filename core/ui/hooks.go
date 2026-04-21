package ui

import (
	"fmt"
	"sync"
)

type Dispatch func(action interface{})

type HookDispatcher struct {
	mu            sync.RWMutex
	dispatchers   map[*Fiber]*HookDispatcherState
	currentFiber  *Fiber
	currentHook   *Hook
	hookIndex     int
}

type HookDispatcherState struct {
	hooks      []*Hook
	workInProgressHook *Hook
}

var currentDispatcher *HookDispatcher
var currentRenderingFiber *Fiber

func init() {
	currentDispatcher = &HookDispatcher{
		dispatchers: make(map[*Fiber]*HookDispatcherState),
	}
}

func (d *HookDispatcher) setCurrentFiber(fiber *Fiber) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.currentFiber = fiber
	d.hookIndex = 0
}

func (d *HookDispatcher) getCurrentFiber() *Fiber {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.currentFiber
}

func (d *HookDispatcher) createHook(memoizedState interface{}) *Hook {
	return &Hook{
		memoizedState: memoizedState,
		baseState:     memoizedState,
		queue:         nil,
		baseQueue:     nil,
		next:          nil,
	}
}

func (d *HookDispatcher) updateHook() (*Hook, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.currentFiber == nil {
		return nil, fmt.Errorf("no current fiber")
	}

	state, exists := d.dispatchers[d.currentFiber]
	if !exists {
		return nil, fmt.Errorf("no hooks state for fiber")
	}

	if d.hookIndex >= len(state.hooks) {
		newHook := d.createHook(nil)
		state.hooks = append(state.hooks, newHook)
		state.workInProgressHook = newHook
	} else {
		state.workInProgressHook = state.hooks[d.hookIndex]
	}

	hook := state.workInProgressHook
	d.hookIndex++
	return hook, nil
}

func (d *HookDispatcher) mountHook(memoizedState interface{}) *Hook {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.currentFiber == nil {
		return nil
	}

	state, exists := d.dispatchers[d.currentFiber]
	if !exists {
		state = &HookDispatcherState{
			hooks: make([]*Hook, 0),
		}
		d.dispatchers[d.currentFiber] = state
	}

	hook := d.createHook(memoizedState)
	state.hooks = append(state.hooks, hook)
	return hook
}

type BasicState interface{}
type Reducer func(oldState BasicState, action interface{}) BasicState

func UseState(initialState interface{}) (interface{}, Dispatch) {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(initialState)
	}

	queue := &UpdateQueue{}
	hook.queue = queue

	dispatch := func(action interface{}) {
		fiber := dispatcher.getCurrentFiber()
		if fiber == nil {
			return
		}

		update := &HookUpdate{
			action: action,
			next:   nil,
		}

		fiberQueue := hook.queue
		if fiberQueue.pending == nil {
			update.next = update
		} else {
			update.next = fiberQueue.pending.next
			fiberQueue.pending.next = update
		}
		hook.queue.pending = update

		scheduler := GetCurrentScheduler()
		scheduler.scheduleSyncTask(func() {
			dispatcher.setCurrentFiber(fiber)
			processUpdateQueue(hook)
		})
	}

	processUpdateQueue(hook)

	return hook.memoizedState, dispatch
}

func processUpdateQueue(hook *Hook) {
	queue := hook.queue
	if queue == nil || queue.pending == nil {
		return
	}

	first := queue.pending.next
	current := first
	previousUpdate := first

	newState := hook.memoizedState

	for current != nil {
		if current == previousUpdate && current != first {
			break
		}

		action := current.action
		if action != nil {
			if fn, ok := action.(func(interface{}) interface{}); ok {
				newState = fn(newState)
			} else if fn, ok := action.(func() interface{}); ok {
				newState = fn()
			} else {
				newState = action
			}
		}

		if current.next == first {
			break
		}
		current = current.next
	}

	hook.memoizedState = newState
	hook.baseState = newState
}

func UseReducer(reducer Reducer, initialState interface{}) (interface{}, Dispatch) {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(initialState)
	}

	if hook.queue == nil {
		queue := &UpdateQueue{}
		hook.queue = queue
	}

	dispatch := func(action interface{}) {
		fiber := dispatcher.getCurrentFiber()
		if fiber == nil {
			return
		}

		update := &HookUpdate{
			action: action,
			next:   nil,
		}

		fiberQueue := hook.queue
		if fiberQueue.pending == nil {
			update.next = update
		} else {
			update.next = fiberQueue.pending.next
			fiberQueue.pending.next = update
		}
		hook.queue.pending = update

		scheduler := GetCurrentScheduler()
		scheduler.scheduleSyncTask(func() {
			dispatcher.setCurrentFiber(fiber)
			processReducerUpdateQueue(hook, reducer)
		})
	}

	if hook.memoizedState == nil {
		processReducerUpdateQueue(hook, reducer)
	}

	return hook.memoizedState, dispatch
}

func processReducerUpdateQueue(hook *Hook, reducer Reducer) {
	queue := hook.queue
	if queue == nil || queue.pending == nil {
		return
	}

	newState := hook.memoizedState
	if newState == nil {
		newState = hook.baseState
	}

	first := queue.pending.next
	current := first
	previousUpdate := first

	for current != nil {
		if current == previousUpdate && current != first {
			break
		}

		action := current.action
		if action != nil {
			newState = reducer(newState.(BasicState), action)
		}

		if current.next == first {
			break
		}
		current = current.next
	}

	hook.memoizedState = newState
}

type EffectCallback func()
type EffectDestroyCallback func()

type Effect struct {
	tag      EffectTag
	create  EffectCallback
	destroy  EffectDestroyCallback
	callback interface{}
	dependencies []interface{}
	next    *Effect
}

type EffectList struct {
	firstEffect *Effect
	lastEffect  *Effect
}

func (d *HookDispatcher) mountEffect(create EffectCallback, callback interface{}, dependencies []interface{}) *Effect {
	effect := &Effect{
		tag:         Update,
		create:      create,
		destroy:     nil,
		callback:    callback,
		dependencies: dependencies,
		next:        nil,
	}
	return effect
}

func UseEffect(create EffectCallback, callback interface{}, dependencies []interface{}) {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(nil)
	}

	effect := &Effect{
		tag:         Update,
		create:      create,
		destroy:     nil,
		callback:    callback,
		dependencies: dependencies,
		next:        nil,
	}

	hook.memoizedState = effect
}

func UseLayoutEffect(create EffectCallback, callback interface{}, dependencies []interface{}) {
	UseEffect(create, callback, dependencies)
}

func UseRef[T any](initialValue T) *Ref[T] {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(&Ref[T]{current: initialValue})
	}

	if hook.memoizedState == nil {
		hook.memoizedState = &Ref[T]{current: initialValue}
	}

	return hook.memoizedState.(*Ref[T])
}

type Ref[T any] struct {
	current T
}

func (r *Ref[T]) Get() T {
	return r.current
}

func (r *Ref[T]) Set(value T) {
	r.current = value
}

func UseCallback[T any](callback T, dependencies []interface{}) T {
	return callback
}

func UseMemo[T any](factory func() T, dependencies []interface{}) T {
	dispatcher := currentDispatcher

	hook, err := dispatcher.updateHook()
	if err != nil {
		hook = dispatcher.mountHook(factory())
		return hook.memoizedState.(T)
	}

	if hook.memoizedState == nil {
		hook.memoizedState = factory()
	}

	return hook.memoizedState.(T)
}

func UseContext[T any](context *ReactContext[T]) T {
	return context.currentValue
}

func CreateContext[T any](defaultValue T) *ReactContext[T] {
	return &ReactContext[T]{
		pendingValue: defaultValue,
		currentValue: defaultValue,
	}
}


