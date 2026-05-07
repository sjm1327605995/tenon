package tenon

import "sync"

// State 是 go-tui 风格的泛型状态对象。
type State[T any] struct {
	mu    sync.RWMutex
	value T
	app   *App
}

// NewState 创建一个具有初始值的 State。
func NewState[T any](initial T) *State[T] {
	return &State[T]{value: initial}
}

// Get 返回当前值。线程安全。
func (s *State[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Set 更新值并标记需要重渲染。
func (s *State[T]) Set(v T) {
	s.mu.Lock()
	s.value = v
	s.mu.Unlock()

	if s.app != nil {
		s.app.Rebuild()
	} else {
		Rebuild()
	}
}

// Update 对当前值应用函数并设置结果。
func (s *State[T]) Update(fn func(T) T) {
	s.Set(fn(s.Get()))
}

// BindApp 将此 State 绑定到指定 App。
func (s *State[T]) BindApp(app *App) {
	s.mu.Lock()
	s.app = app
	s.mu.Unlock()
}
