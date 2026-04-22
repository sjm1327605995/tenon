package scheduler

import (
	"sync"

	"github.com/sjm1327605995/tenon/pkg/core"
)

// Phase 定义调度器的执行阶段
type Phase int

const (
	PhaseNone        Phase = iota // 无状态
	PhaseMounting                 // 挂载阶段
	PhaseUpdating                 // 更新阶段
	PhaseUnmounting               // 卸载阶段
)

// Scheduler 组件生命周期调度器
// 负责管理组件的挂载、更新和卸载流程
type Scheduler struct {
	mu             sync.Mutex
	dirtyComponents map[string]core.Component // 待更新的脏组件
	mountedSet      map[string]bool           // 已挂载组件集合
	phase           Phase                      // 当前执行阶段
}

var instance *Scheduler
var once sync.Once

// GetInstance 获取调度器单例（框架内部使用）
func GetInstance() *Scheduler {
	once.Do(func() {
		instance = &Scheduler{
			dirtyComponents: make(map[string]core.Component),
			mountedSet:      make(map[string]bool),
			phase:           PhaseNone,
		}
	})
	return instance
}

// ----------------------------------------------------------------------------
// 【框架内部方法】调度器核心方法
// ----------------------------------------------------------------------------

// ScheduleUpdate 标记组件为脏组件，需要更新（框架内部调用）
func (s *Scheduler) ScheduleUpdate(component core.Component) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if component == nil {
		return
	}

	s.dirtyComponents[component.ID()] = component
}

// ClearDirty 清除组件的脏标记（框架内部调用）
func (s *Scheduler) ClearDirty(component core.Component) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.dirtyComponents, component.ID())
}

// IsMounted 检查组件是否已挂载（框架内部调用）
func (s *Scheduler) IsMounted(component core.Component) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.mountedSet[component.ID()]
}

// MarkMounted 标记组件为已挂载（框架内部调用）
func (s *Scheduler) MarkMounted(component core.Component) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.mountedSet[component.ID()] = true
}

// MarkUnmounted 标记组件为已卸载（框架内部调用）
func (s *Scheduler) MarkUnmounted(component core.Component) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.mountedSet, component.ID())
}

// ProcessMount 处理组件挂载（框架内部调用）
// 递归调用所有子组件的 ComponentDidMount 方法
func (s *Scheduler) ProcessMount(root core.Component) {
	s.mu.Lock()
	s.phase = PhaseMounting
	s.mu.Unlock()

	processMount(root)

	s.mu.Lock()
	s.phase = PhaseNone
	s.mu.Unlock()
}

func processMount(component core.Component) {
	if component == nil {
		return
	}

	for _, child := range component.GetChildren() {
		processMount(child)
	}

	component.ComponentDidMount()
	GetInstance().MarkMounted(component)
}

// ProcessUpdates 处理组件更新（框架内部调用）
// 遍历所有脏组件，调用生命周期方法
func (s *Scheduler) ProcessUpdates() {
	s.mu.Lock()
	if s.phase != PhaseNone || len(s.dirtyComponents) == 0 {
		s.mu.Unlock()
		return
	}

	s.phase = PhaseUpdating
	dirtyList := make([]core.Component, 0, len(s.dirtyComponents))
	for _, comp := range s.dirtyComponents {
		dirtyList = append(dirtyList, comp)
	}
	s.dirtyComponents = make(map[string]core.Component)
	s.mu.Unlock()

	for _, component := range dirtyList {
		if !GetInstance().IsMounted(component) {
			continue
		}

		prevProps := component.GetProps()
		prevState := component.GetState()

		if !component.ShouldComponentUpdate(nil, nil) {
			continue
		}

		snapshot := component.GetSnapshotBeforeUpdate(prevProps, prevState)
		component.ComponentDidUpdate(prevProps, prevState, snapshot)
	}

	s.mu.Lock()
	s.phase = PhaseNone
	s.mu.Unlock()
}

// ProcessUnmount 处理组件卸载（框架内部调用）
// 递归调用所有子组件的 ComponentWillUnmount 方法
func (s *Scheduler) ProcessUnmount(root core.Component) {
	s.mu.Lock()
	s.phase = PhaseUnmounting
	s.mu.Unlock()

	processUnmount(root)

	s.mu.Lock()
	s.phase = PhaseNone
	s.mu.Unlock()
}

func processUnmount(component core.Component) {
	if component == nil {
		return
	}

	for _, child := range component.GetChildren() {
		processUnmount(child)
	}

	component.ComponentWillUnmount()
	GetInstance().MarkUnmounted(component)
}

// HasDirtyComponents 检查是否有待更新的脏组件（框架内部调用）
func (s *Scheduler) HasDirtyComponents() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.dirtyComponents) > 0
}

// GetPhase 获取当前调度阶段（框架内部调用）
func (s *Scheduler) GetPhase() Phase {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.phase
}

// Reset 重置调度器状态（框架内部调用）
func (s *Scheduler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.dirtyComponents = make(map[string]core.Component)
	s.mountedSet = make(map[string]bool)
	s.phase = PhaseNone
}
