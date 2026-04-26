package core

// State 是可观察的信号容器。
// 调用 Set() 时只通知订阅者，不触发任何重绘或布局。
// 是否刷新由订阅者（通常是 Element Setter）自行决定。
type State[T any] struct {
	value T
	subs  []func(T)
}

// NewState 创建状态容器。
func NewState[T any](v T) *State[T] {
	return &State[T]{value: v}
}

// Get 读取当前值。
func (s *State[T]) Get() T {
	return s.value
}

// Set 修改值并通知所有订阅者。
// 使用 comparable 约束时可以做去重，这里保持通用，由调用方保证。
func (s *State[T]) Set(v T) {
	s.value = v
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
			s.subs[idx] = nil // mark as removed
		}
	}
}

// Bind 是 Subscribe 的语法糖，适合链式场景。
// 返回取消订阅函数，由 Element 在卸载时调用。
func (s *State[T]) Bind(fn func(T)) func() {
	fn(s.value) // 立即执行一次，同步初始值
	return s.Subscribe(fn)
}

// ==================== 兼容旧 any 类型 State ====================

// AnyState 是 State[any] 的别名，兼容旧代码。
type AnyState = State[any]

// NewAnyState 创建 State[any]。
func NewAnyState(v any) *AnyState {
	return NewState[any](v)
}
