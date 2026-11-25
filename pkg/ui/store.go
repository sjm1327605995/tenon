package ui

import "reflect"

// Store 是组件树之外的应用级状态容器（对标 zustand 的精简版）。
//
// 与 UseState 的分工：组件内部状态用 UseState；跨组件、页面级共享的业务状态放 Store。
// 与 Signal 的分工：Store 变化重渲染订阅组件（粒度=组件，经 reconcile）；Signal 直连
// 叶子 renderNode（粒度=节点），适合每帧级高频值。
//
// 线程约定与 setState 一致：Set/Update 须在 UI 线程（事件回调、render、ui.Post 内）调用；
// 后台 goroutine 用 ui.Post 包一层。
type Store[T any] struct {
	value T
	subs  map[uint64]*storeSub
	seq   uint64
}

// storeSub 是一次组件订阅：selector 把状态投影成组件关心的切片，投影不变则跳过重渲染。
type storeSub struct {
	fiber  *Fiber
	sel    func(any) any // 类型擦除的 selector：func(T) S
	last   any           // 上次渲染所用的投影值
	equals func(any, any) bool
}

// NewStore 创建应用级状态容器。
func NewStore[T any](initial T) *Store[T] { return &Store[T]{value: initial} }

// Get 返回当前状态（只读快照；修改请走 Set/Update）。
func (s *Store[T]) Get() T { return s.value }

// Set 整体替换状态并通知订阅者：selector 投影不变的组件不重渲染（同 React 的
// useSyncExternalStore + selector 语义）。新值与旧值深度相等时不做任何事。
func (s *Store[T]) Set(v T) {
	if reflect.DeepEqual(s.value, v) {
		return
	}
	s.value = v
	for _, sub := range s.subs {
		next := sub.sel(v)
		if sub.equals(sub.last, next) {
			continue
		}
		sub.last = next
		if !sub.fiber.unmounted && sub.fiber.game != nil {
			sub.fiber.game.markDirty(sub.fiber)
		}
	}
}

// Update 以当前状态为入参的函数式更新。
func (s *Store[T]) Update(fn func(T) T) { s.Set(fn(s.value)) }

type storeHook struct {
	sub         *storeSub
	store       any
	unsubscribe func()
}

// UseStore 订阅 store：selector 选中的切片变化时本组件重渲染，未变化则跳过；
// 组件卸载时自动退订。多个组件可各选不同的切片，互不重渲染。
//
//	count := ui.UseStore(cartStore, func(s Cart) int { return s.Count })
func UseStore[T any, S any](s *Store[T], selector func(T) S) S {
	return useStore(s, selector, func(a, b S) bool { return reflect.DeepEqual(a, b) })
}

// UseStoreEq 与 UseStore 相同，但用自定义相等函数判定选中切片是否变化——
// 热路径上避免 reflect.DeepEqual 的开销（如切片很大或含不可比字段时）。
func UseStoreEq[T any, S any](s *Store[T], selector func(T) S, equals func(S, S) bool) S {
	return useStore(s, selector, equals)
}

func useStore[T any, S any](s *Store[T], selector func(T) S, equals func(S, S) bool) S {
	f := currentFiber
	_, raw := nextHook(f, func() any { return &storeHook{} })
	h := raw.(*storeHook)
	if h.store != s {
		if h.unsubscribe != nil {
			h.unsubscribe()
		}
		sub := &storeSub{
			fiber:  f,
			equals: func(a, b any) bool { return equals(valueOrZero[S](a), valueOrZero[S](b)) },
		}
		sub.sel = func(v any) any { return selector(valueOrZero[T](v)) }
		sub.last = sub.sel(s.value)
		s.seq++
		id := s.seq
		if s.subs == nil {
			s.subs = map[uint64]*storeSub{}
		}
		s.subs[id] = sub
		h.sub = sub
		h.store = s
		h.unsubscribe = func() { delete(s.subs, id) }
		f.addCleanup(func() {
			if h.unsubscribe != nil {
				h.unsubscribe()
				h.unsubscribe = nil
			}
		})
	}
	// 每轮 render 刷新 selector 闭包：selector 可能捕获本轮作用域的变量。
	h.sub.sel = func(v any) any { return selector(valueOrZero[T](v)) }
	// getSnapshot 式复核：投影值在渲染期被直接改过（非经 Set）时对齐，防止显示回退。
	h.sub.last = h.sub.sel(s.value)
	return valueOrZero[S](h.sub.last)
}
