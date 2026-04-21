package ui

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type Lane int

const (
	NoLane Lane = 0
	SyncLane Lane = 1 << iota
	InputContinuousLane
	DefaultLane
	TransitionLane1
	TransitionLane2
	IdleLane
)

type SchedulerPriority int

const (
	NoPriority SchedulerPriority = 0
	ImmediatePriority SchedulerPriority = 1
	LowPriority SchedulerPriority = 2
	NormalPriority SchedulerPriority = 3
	HighPriority SchedulerPriority = 4
)

type WorkTag int

const (
	FunctionComponent WorkTag = 0
	ClassComponent WorkTag = 1
	HostRoot WorkTag = 2
	HostComponent WorkTag = 3
	Fragment WorkTag = 4
)

type Phase int

const (
	Render Phase = 0
	Commit Phase = 1
)

type EffectTag int

const (
	NoEffect EffectTag = 0
	Placement EffectTag = 1 << iota
	Update
	PlacementAndUpdate
	Deletion
	Hydrating
)

type Fiber struct {
	tag           WorkTag
	key           string
	stateNode     interface{}
	type_         interface{}
	lanes          Lane
	childLanes     Lane
	alternate      *Fiber
	return_        *Fiber
	child          *Fiber
	sibling        *Fiber
	index          int
	pendingProps   interface{}
	memorizedProps interface{}
	memorizedState interface{}
	ref            interface{}
	effectTag      EffectTag
	nextEffect     *Fiber
	firstEffect    *Fiber
	lastEffect     *Fiber
	phase          Phase
	expirationTime int64
	pendingTime    int64
}

func (f *Fiber) Render() Element {
	if f.type_ == nil {
		return nil
	}
	if fn, ok := f.type_.(func() Element); ok {
		return fn()
	}
	if fn, ok := f.type_.(func(props interface{}) Element); ok {
		if f.pendingProps != nil {
			return fn(f.pendingProps)
		}
		return fn(nil)
	}
	return nil
}

type Hook struct {
	memoizedState interface{}
	baseState     interface{}
	queue         *UpdateQueue
	baseQueue     *UpdateQueue
	next          *Hook
}

type HookUpdate struct {
	action interface{}
	next   *HookUpdate
	_lane  Lane
}

type UpdateQueue struct {
	pending *HookUpdate
}

type SchedulerTask struct {
	callback func() bool
	priority SchedulerPriority
	lanes    Lane
}

type Scheduler struct {
	mu              sync.Mutex
	workLoop        chan func()
	pendingTasks    *list.List
	runningTask     *list.Element
	taskCounter     uint64
	isRunning       Bool
	yieldInterval   int64
	currentTask     *SchedulerTask
	priorityQueue   map[SchedulerPriority]*list.List
}

var defaultScheduler *Scheduler

func init() {
	defaultScheduler = &Scheduler{
		workLoop:      make(chan func(), 100),
		pendingTasks:  list.New(),
		yieldInterval: 5000,
		priorityQueue: map[SchedulerPriority]*list.List{
			ImmediatePriority: list.New(),
			HighPriority:      list.New(),
			NormalPriority:    list.New(),
			LowPriority:       list.New(),
			NoPriority:        list.New(),
		},
	}
	go defaultScheduler.schedulerLoop()
}

func (s *Scheduler) schedulerLoop() {
	for {
		select {
		case task := <-s.workLoop:
			if task != nil {
				task()
			}
		default:
			s.processTask()
		}
	}
}

func (s *Scheduler) processTask() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var task *list.Element
	for priority := ImmediatePriority; priority >= NoPriority; priority-- {
		if queue := s.priorityQueue[priority]; queue != nil && queue.Len() > 0 {
			task = queue.Front()
			queue.Remove(task)
			break
		}
	}

	if task == nil {
		return
	}

	schedulerTask := task.Value.(*SchedulerTask)
	s.currentTask = schedulerTask

	s.mu.Unlock()
	schedulerTask.callback()
	s.mu.Lock()

	s.currentTask = nil
}

func (s *Scheduler) scheduleTask(callback func() bool, priority SchedulerPriority, lanes Lane) uint64 {
	task := &SchedulerTask{
		callback: callback,
		priority: priority,
		lanes:    lanes,
	}

	taskId := atomic.AddUint64(&s.taskCounter, 1)

	s.mu.Lock()
	defer s.mu.Unlock()

	if queue := s.priorityQueue[priority]; queue != nil {
		queue.PushBack(task)
	}

	return taskId
}

func (s *Scheduler) scheduleSyncTask(callback func()) uint64 {
	taskFn := func() bool {
		callback()
		return false
	}
	return s.scheduleTask(taskFn, ImmediatePriority, SyncLane)
}

func (s *Scheduler) scheduleTransitionTask(callback func()) uint64 {
	taskFn := func() bool {
		callback()
		return false
	}
	return s.scheduleTask(taskFn, NormalPriority, TransitionLane1)
}

func (s *Scheduler) scheduleHighPriorityTask(callback func()) uint64 {
	taskFn := func() bool {
		callback()
		return false
	}
	return s.scheduleTask(taskFn, HighPriority, DefaultLane)
}

func (s *Scheduler) getCurrentPriorityLevel() SchedulerPriority {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentTask != nil {
		return s.currentTask.priority
	}
	return NormalPriority
}

func GetCurrentScheduler() *Scheduler {
	return defaultScheduler
}

type ReactContext[T any] struct {
	pendingValue T
	currentValue T
}

func (c *ReactContext[T]) Read() T {
	return c.currentValue
}

func (c *ReactContext[T]) Write(value T) {
	c.pendingValue = value
}


