package core

// Switch 是声明式路由/页面切换组件。
// 它根据 State 的值自动切换显示的子 Element，用户不需要写 switch 或手动管理子节点。
// Switch 本身是 Element，可以像 View、Text 一样被挂载到树中。
type Switch[T comparable] struct {
	BaseElement
	state      *State[T]
	cases      map[T]func() Element
	defaultFn  func() Element
	current    Element
	cleanup    func()
}

// NewSwitch 创建声明式切换容器。
func NewSwitch[T comparable](state *State[T]) *Switch[T] {
	s := &Switch[T]{
		state: state,
		cases: make(map[T]func() Element),
	}
	s.Init(s)
	return s
}

// Case 注册一个分支。fn 是延迟创建函数，只在切换到该分支时调用。
func (s *Switch[T]) Case(key T, fn func() Element) *Switch[T] {
	s.cases[key] = fn
	return s
}

// Default 注册默认分支，当没有匹配的 Case 时显示。
func (s *Switch[T]) Default(fn func() Element) *Switch[T] {
	s.defaultFn = fn
	return s
}

// OnMount 挂载时自动订阅 State 并显示当前分支。
func (s *Switch[T]) OnMount(engine *Engine) {
	s.BaseElement.OnMount(engine)

	// 立即显示当前值对应的分支
	s.show(s.state.Get())

	// 订阅后续变化
	s.cleanup = s.state.Subscribe(func(v T) {
		s.show(v)
	})
}

// OnUnmount 卸载时清理订阅和子节点。
func (s *Switch[T]) OnUnmount() {
	if s.cleanup != nil {
		s.cleanup()
		s.cleanup = nil
	}
	if s.current != nil {
		s.RemoveChild(s.current)
		s.current = nil
	}
	s.BaseElement.OnUnmount()
}

// show 切换到指定分支。
func (s *Switch[T]) show(key T) {
	// 查找对应的创建函数
	fn, ok := s.cases[key]
	if !ok && s.defaultFn != nil {
		fn = s.defaultFn
	}
	if fn == nil {
		// 没有匹配，清空
		if s.current != nil {
			s.RemoveChild(s.current)
			s.current = nil
		}
		return
	}

	// 如果当前显示的已经是这个分支，跳过
	// TODO: 如果需要支持 Key 复用，可以在这里扩展

	// 卸载旧的
	if s.current != nil {
		s.RemoveChild(s.current)
		s.current = nil
	}

	// 创建并挂载新的
	newEl := fn()
	if newEl != nil {
		s.current = newEl
		s.AppendChild(newEl)
		// 新子节点加入后需要重新布局
		s.Mark(FlagNeedLayout)
	}
}

// ElementType 返回类型标识。
func (s *Switch[T]) ElementType() string { return "Switch" }
