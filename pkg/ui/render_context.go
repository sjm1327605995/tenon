package ui

import (
	"context"
	"reflect"
	"time"
)

// RenderContext belongs to one mounted object component. UI state operations run
// on the UI thread; Post is safe from background goroutines. It also implements
// context.Context and is canceled when the component or its parent unmounts.
type RenderContext struct {
	fiber    *Fiber
	game     *game
	owner    *Component
	lifetime context.Context
	cancel   context.CancelFunc
}

func (c *RenderContext) Deadline() (time.Time, bool) { return c.lifetime.Deadline() }
func (c *RenderContext) Done() <-chan struct{}       { return c.lifetime.Done() }
func (c *RenderContext) Err() error                  { return c.lifetime.Err() }
func (c *RenderContext) Value(key any) any           { return c.lifetime.Value(key) }

// Mounted reports whether this component's lifetime is still active.
func (c *RenderContext) Mounted() bool { return c != nil && c.lifetime.Err() == nil }

// Update changes component fields and requests a local rerender. Late updates
// after unmount are ignored, including the mutation callback itself.
func (c *RenderContext) Update(update func()) {
	if !c.Mounted() {
		return
	}
	if update != nil {
		update()
	}
	if c.Mounted() {
		c.game.markDirty(c.fiber)
	}
}

// Invalidate requests a local rerender after an external UI-thread mutation.
func (c *RenderContext) Invalidate() { c.Update(nil) }

// Post schedules a state update on the UI thread, guarded by this lifetime.
func (c *RenderContext) Post(update func()) {
	if !c.Mounted() {
		return
	}
	Post(func() { c.Update(update) })
}

// Theme reads and subscribes to the nearest theme provider for this component.
func (c *RenderContext) Theme() Theme { return ContextValueIn(c, themeContext) }

// Viewport reads the owning window's logical bounds and subscribes to resize.
func (c *RenderContext) Viewport() Rect {
	if !c.Mounted() {
		return Rect{}
	}
	if c.game.viewportReaders == nil {
		c.game.viewportReaders = make(map[*Fiber]struct{})
	}
	c.game.viewportReaders[c.fiber] = struct{}{}
	return Rect{W: float32(c.game.w) / uiScale, H: float32(c.game.h) / uiScale}
}

func (c *RenderContext) requireRender() {
	if !c.Mounted() || currentFiber != c.fiber {
		panic("ui: context state and effects must be declared inside their component's Render")
	}
}

// Effect registers an effect and cleanup using the existing Hooks lifecycle.
// Call it in the same order on each Render, just like UseEffect.
func (c *RenderContext) Effect(effect func() Cleanup, deps ...any) {
	c.requireRender()
	UseEffect(effect, deps...)
}

// StateIn is typed state belonging to ctx, with a stable setter that ignores
// updates after unmount. Call it in the same order on each Render.
func StateIn[T any](ctx *RenderContext, initial T) (T, func(T)) {
	ctx.requireRender()
	_, raw := nextHook(ctx.fiber, func() any { return &stateHook{value: initial} })
	h := raw.(*stateHook)
	if h.setter == nil {
		h.setter = func(value T) {
			if !ctx.Mounted() || reflect.DeepEqual(h.value, value) {
				return
			}
			h.value = value
			ctx.game.markDirty(ctx.fiber)
		}
	}
	return valueOrZero[T](h.value), h.setter.(func(T))
}

// ContextValueIn reads a typed provider value using ctx's component scope.
func ContextValueIn[T any](ctx *RenderContext, key *Context[T]) T {
	if ctx.Mounted() {
		for parent := ctx.fiber.parent; parent != nil; parent = parent.parent {
			if parent.typ == typeProvider && parent.ctxID == key.id {
				subscribe(parent, ctx.fiber)
				return parent.ctxValue.(T)
			}
		}
	}
	return key.def
}

var _ context.Context = (*RenderContext)(nil)
