package ui

import (
	"reflect"

	values "github.com/sjm1327605995/tenon/pkg/value"
)

type observedValue[T any] struct {
	key  any
	cell *sigSub
	last T
}

// UseValue subscribes the current component to a Value. Use it when layout,
// branches, or custom component logic depend on the value; Text and BindInput
// bind directly and do not need to rerender the component.
func UseValue[T any](source values.Value[T]) T {
	f, g := currentFiber, activeGame
	if f == nil || g == nil {
		panic("ui.UseValue: call inside a component render")
	}
	_, raw := nextHook(f, func() any {
		h := &observedValue[T]{}
		f.addCleanup(func() {
			if h.cell != nil {
				h.cell.cancel()
			}
		})
		return h
	})
	h := raw.(*observedValue[T])
	key := bindingKey(source)
	if h.cell == nil || h.key != key {
		if h.cell != nil {
			h.cell.cancel()
		}
		h.key, h.last = key, source.Get()
		cell := &sigSub{active: true}
		h.cell = cell
		cell.unsub = source.Subscribe(func(next T) {
			if !cell.active || f.unmounted || reflect.DeepEqual(h.last, next) {
				return
			}
			h.last = next
			g.markDirty(f)
		})
	}
	h.last = source.Get()
	return h.last
}
