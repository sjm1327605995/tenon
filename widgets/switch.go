package widgets

import "github.com/sjm1327605995/tenon/internal/core"


// Switch 是声明式路由/页面切换组件。
// 它根据 core.State 的值自动切换显示的子 core.Element，用户不需要写 switch 或手动管理子节点。
// Switch 本身是 core.Element，可以像 View、Text 一样被挂载到树中。
type Switch[T comparable] struct {
	core.BaseWidget
	state     *core.State[T]
	cases     map[T]func() core.Element
	defaultFn func() core.Element
	current   core.Element
	cleanup   func()
}

// NewSwitch 创建声明式切换容器。
func NewSwitch[T comparable](state *core.State[T]) *Switch[T] {
	s := &Switch[T]{
		state: state,
		cases: make(map[T]func() core.Element),
	}
	s.Init(s)
	return s
}

// Render builds the switched element tree.
func (s *Switch[T]) Render() core.Element {
	return s.show(s.state.Get())
}

// Case 注册一个分支。fn 是延迟创建函数，只在切换到该分支时调用。
func (s *Switch[T]) Case(key T, fn func() core.Element) *Switch[T] {
	s.cases[key] = fn
	return s
}

// Default 注册默认分支，当没有匹配的 Case 时显示。
func (s *Switch[T]) Default(fn func() core.Element) *Switch[T] {
	s.defaultFn = fn
	return s
}

// OnMount 挂载时自动订阅 core.State 并显示当前分支。
func (s *Switch[T]) OnMount() {
	// 订阅后续变化
	s.cleanup = s.state.Subscribe(func(v T) {
		s.RequestBuild()
	})
}

// OnUnmount 卸载时清理订阅和子节点。
func (s *Switch[T]) OnUnmount() {
	if s.cleanup != nil {
		s.cleanup()
		s.cleanup = nil
	}
}

func (s *Switch[T]) show(key T) core.Element {
	fn, ok := s.cases[key]
	if !ok && s.defaultFn != nil {
		fn = s.defaultFn
	}
	if fn == nil {
		return nil
	}
	return fn()
}
