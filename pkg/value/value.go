// Package value defines extensible values for UI data binding. Operations and
// notifications are serialized by the caller; UI-bound values run on the UI thread.
package value

import (
	"reflect"
	"slices"
)

// Value is the read and notification protocol. Subscribe does not have to emit
// an initial value; consumers subscribe first and then read Get. The returned
// function ends that subscription and may be called more than once.
type Value[T any] interface {
	Get() T
	Subscribe(func(T)) func()
}

// Bind is a value that also accepts changes. A custom source only needs these
// methods; it does not need to embed or inherit the built-in implementation.
type Bind[T any] interface {
	Value[T]
	Set(T)
}

// Hooks configures optional behavior of a built-in Bind. BeforeSet callbacks run
// in order and may normalize or reject a value. OnChange runs after committing,
// before subscribers. Subscription hooks run once per subscription/cancellation.
type Hooks[T any] struct {
	BeforeSet     func(old, next T) (T, error)
	OnChange      func(old, next T)
	OnError       func(error)
	OnSubscribe   func()
	OnUnsubscribe func()
}

type binding[T any] struct {
	value   T
	hooks   []Hooks[T]
	subs    map[uint64]func(T)
	seq     uint64
	pending []T
	setting bool
}

// New creates a Bind. Get and Set work immediately, even without UI consumers.
func New[T any](initial T, hooks ...Hooks[T]) Bind[T] {
	return &binding[T]{value: initial, hooks: append([]Hooks[T](nil), hooks...)}
}

func (b *binding[T]) Get() T { return b.value }

func (b *binding[T]) Set(next T) {
	b.pending = append(b.pending, next)
	if b.setting {
		return
	}
	b.setting = true
	defer func() { b.setting = false; b.pending = nil }()
	for len(b.pending) > 0 {
		next := b.pending[0]
		var zero T
		b.pending[0] = zero
		b.pending = b.pending[1:]
		b.commit(next)
	}
}

func (b *binding[T]) commit(next T) {
	old := b.value
	for _, hook := range b.hooks {
		if hook.BeforeSet != nil {
			var err error
			next, err = hook.BeforeSet(old, next)
			if err != nil {
				for _, observer := range b.hooks {
					if observer.OnError != nil {
						observer.OnError(err)
					}
				}
				return
			}
		}
	}
	if reflect.DeepEqual(old, next) {
		return
	}
	b.value = next
	ids := make([]uint64, 0, len(b.subs))
	for id := range b.subs {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	for _, hook := range b.hooks {
		if hook.OnChange != nil {
			hook.OnChange(old, next)
		}
	}
	for _, id := range ids {
		if fn, active := b.subs[id]; active {
			fn(next)
		}
	}
}

func (b *binding[T]) Subscribe(fn func(T)) func() {
	if fn == nil {
		return func() {}
	}
	if b.subs == nil {
		b.subs = make(map[uint64]func(T))
	}
	b.seq++
	id := b.seq
	b.subs[id] = fn
	for _, hook := range b.hooks {
		if hook.OnSubscribe != nil {
			hook.OnSubscribe()
		}
	}
	return func() {
		if _, active := b.subs[id]; !active {
			return
		}
		delete(b.subs, id)
		for _, hook := range b.hooks {
			if hook.OnUnsubscribe != nil {
				hook.OnUnsubscribe()
			}
		}
	}
}

type constant[T any] struct{ value T }

func (c constant[T]) Get() T                 { return c.value }
func (constant[T]) Subscribe(func(T)) func() { return func() {} }

// Static wraps an unchanging value in the same protocol as live values.
func Static[T any](v T) Value[T] { return constant[T]{value: v} }

type mapped[A, B any] struct {
	source Value[A]
	mapFn  func(A) B
}

func (m *mapped[A, B]) Get() B { return m.mapFn(m.source.Get()) }
func (m *mapped[A, B]) Subscribe(fn func(B)) func() {
	if fn == nil {
		return func() {}
	}
	return m.source.Subscribe(func(v A) { fn(m.mapFn(v)) })
}

// Map derives a read-only value, for example a formatted label from a Bind[int].
func Map[A, B any](source Value[A], fn func(A) B) Value[B] {
	return &mapped[A, B]{source: source, mapFn: fn}
}
