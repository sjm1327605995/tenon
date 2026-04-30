package widgets

import "github.com/sjm1327605995/tenon/internal/core"


// If is a declarative conditional rendering element.
// It shows the "then" branch when condition is true, otherwise the "else" branch.
type If struct {
	core.BaseElement
	condition *core.State[bool]
	thenFn    func() core.Element
	elseFn    func() core.Element
	current   core.Element
	cleanup   func()
}

// NewIf creates a conditional container.
func NewIf(condition *core.State[bool]) *If {
	i := &If{condition: condition}
	i.Init(i)
	return i
}

// Then sets the branch shown when condition is true.
func (i *If) Then(fn func() core.Element) *If {
	i.thenFn = fn
	return i
}

// Else sets the branch shown when condition is false.
func (i *If) Else(fn func() core.Element) *If {
	i.elseFn = fn
	return i
}

// OnMount mounts and subscribes to condition changes.
func (i *If) OnMount(engine *core.Engine) {
	i.BaseElement.OnMount(engine)
	i.show(i.condition.Get())
	i.cleanup = i.condition.Subscribe(func(v bool) {
		i.show(v)
	})
}

// OnUnmount cleans up.
func (i *If) OnUnmount() {
	if i.cleanup != nil {
		i.cleanup()
		i.cleanup = nil
	}
	if i.current != nil {
		i.RemoveChild(i.current)
		i.current = nil
	}
	i.BaseElement.OnUnmount()
}

func (i *If) show(cond bool) {
	var fn func() core.Element
	if cond {
		fn = i.thenFn
	} else {
		fn = i.elseFn
	}
	if fn == nil {
		if i.current != nil {
			i.RemoveChild(i.current)
			i.current = nil
		}
		return
	}
	if i.current != nil {
		i.RemoveChild(i.current)
		i.current = nil
	}
	newEl := fn()
	if newEl != nil {
		i.current = newEl
		i.AppendChild(newEl)
		i.Mark(core.FlagNeedLayout)
	}
}

// ElementType returns type identifier.
func (i *If) ElementType() string { return "If" }
