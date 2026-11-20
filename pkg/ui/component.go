package ui

import (
	"context"
	"reflect"
	"sync/atomic"
)

// IComponent is an object component. Embed Component in a struct and implement
// Render; the runtime invokes Render through this interface.
type IComponent interface {
	NodeSource
	Render() NodeSource
	Context() *RenderContext
	ComponentBase() *Component
}

// Component supplies lifecycle-bound context to an embedding object. Render is
// deliberately not defined here: embedded Go methods do not dispatch to methods
// on the enclosing struct. The runtime dispatches through IComponent instead.
type Component struct {
	context atomic.Pointer[RenderContext]
}

// Context returns the current component context, or nil before mounting.
func (c *Component) Context() *RenderContext {
	if c == nil {
		return nil
	}
	return c.context.Load()
}

// ComponentBase lets the runtime bind context to the actual embedded base.
func (c *Component) ComponentBase() *Component { return c }

// SetState changes fields and schedules this component's render on the UI thread.
func (c *Component) SetState(update func()) { c.Context().Update(update) }

// Build is the declaration protocol. Object components are resolved through
// IComponent by the runtime before this fallback can be reached.
func (*Component) Build() *Node {
	panic("ui: embed Component and implement Render() ui.NodeSource; the runtime renders it through IComponent")
}

type componentInstanceProps struct{ component IComponent }

// Render describes an object component without invoking it before mounting.
// A component's concrete type and key determine its Fiber identity.
func Render(component IComponent) *Node {
	if component == nil {
		return nil
	}
	v := reflect.ValueOf(component)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}
	n := Use(renderComponentInstance, componentInstanceProps{component: component})
	n.instanceType = reflect.TypeOf(component)
	return n
}

func renderComponentInstance(p componentInstanceProps) *Node {
	f := currentFiber
	base := p.component.ComponentBase()
	if base == nil {
		panic("ui: IComponent.ComponentBase returned nil")
	}
	if previous := base.context.Load(); previous != nil && previous.Mounted() && previous.fiber != f {
		panic("ui: the same component instance cannot be mounted in two places")
	}
	ctx := f.renderContext
	if ctx == nil {
		parent := context.Background()
		for ancestor := f.parent; ancestor != nil; ancestor = ancestor.parent {
			if ancestor.renderContext != nil {
				parent = ancestor.renderContext
				break
			}
		}
		lifetime, cancel := context.WithCancel(parent)
		ctx = &RenderContext{fiber: f, game: activeGame, lifetime: lifetime, cancel: cancel}
		f.renderContext = ctx
	}
	if ctx.owner != nil && ctx.owner != base && ctx.owner.context.Load() == ctx {
		ctx.owner.context.Store(nil)
	}
	ctx.owner = base
	base.context.Store(ctx)
	return nodeFromSource(p.component.Render())
}
