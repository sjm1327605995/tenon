package core

// State 是可观察的信号容器。
// 调用 Set() 时只通知订阅者，不触发任何重绘或布局。
// 是否刷新由订阅者（通常是 Element Setter）自行决定。
//
// 自动依赖追踪：Widget.Render() 期间调用 Get() 会自动记录访问，
// Render 完成后框架自动订阅这些 State，变化时自动触发 Widget 重建。
type State[T any] struct {
	value     T
	subs      []func(T)
	onSetHook func(oldVal, newVal T)
}

// renderTracker 在 Widget.Render() 期间收集 State 访问。
type renderTracker struct {
	widget Widget
	active bool
	states []stateNotifier
}

// activeRenderTracker 在 Render 期间指向当前 Engine 的 tracker。
// 这是同步上下文指针，不构成持久全局状态。
var activeRenderTracker *renderTracker

var stateDebugHook func(elementType, key string, oldValue, newValue interface{})

func SetStateDebugHook(hook func(elementType, key string, oldValue, newValue interface{})) {
	stateDebugHook = hook
}

// NewState 创建状态容器。
func NewState[T any](v T) *State[T] {
	return &State[T]{value: v}
}

// Get 读取当前值。
// 如果在 Widget.Render() 执行期间调用，会自动记录依赖关系，
// 框架后续会自动订阅该 State，变化时触发 Widget 重建。
func (s *State[T]) Get() T {
	if activeRenderTracker != nil && activeRenderTracker.active && activeRenderTracker.widget != nil {
		recordRenderState(s)
	}
	return s.value
}

// Set 修改值并通知所有订阅者。
func (s *State[T]) Set(v T) {
	old := s.value
	s.value = v
	if s.onSetHook != nil {
		s.onSetHook(old, v)
	}
	if stateDebugHook != nil {
		stateDebugHook("", "", old, v)
	}
	for _, fn := range s.subs {
		if fn != nil {
			fn(v)
		}
	}
}

// Subscribe 注册变化回调，返回取消订阅函数。
func (s *State[T]) Subscribe(fn func(T)) func() {
	idx := len(s.subs)
	s.subs = append(s.subs, fn)
	return func() {
		if idx < len(s.subs) && s.subs[idx] != nil {
			s.subs[idx] = nil
		}
	}
}

// Bind 是 Subscribe 的语法糖，适合链式场景。
// 返回取消订阅函数，由 Element 在卸载时调用。
func (s *State[T]) Bind(fn func(T)) func() {
	fn(s.value)
	return s.Subscribe(fn)
}

// _subscribeRebuild 内部接口，用于 Widget 自动重建。
func (s *State[T]) _subscribeRebuild(fn func()) func() {
	return s.Subscribe(func(T) { fn() })
}

// _subscribeRebuild 内部接口，用于 Widget 自动重建。
type stateNotifier interface {
	_subscribeRebuild(fn func()) func()
}

// recordRenderState 记录一个 State 被当前 Widget 访问。
func recordRenderState(s stateNotifier) {
	if activeRenderTracker == nil || !activeRenderTracker.active {
		return
	}
	for _, existing := range activeRenderTracker.states {
		if existing == s {
			return
		}
	}
	activeRenderTracker.states = append(activeRenderTracker.states, s)
}

// ==================== 兼容旧 any 类型 State ====================

// AnyState 是 State[any] 的别名，兼容旧代码。
type AnyState = State[any]

// NewAnyState 创建 State[any]。
func NewAnyState(v any) *AnyState {
	return NewState[any](v)
}
