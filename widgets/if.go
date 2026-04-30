package widgets

import "github.com/sjm1327605995/tenon/internal/core"


// If is a declarative conditional rendering element.
// It shows the "then" branch when condition is true, otherwise the "else" branch.
type If struct {
	core.BaseWidget
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

// Render builds the conditional element tree.
func (i *If) Render() core.Element {
	return i.show(i.condition.Get())
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
func (i *If) OnMount() {
	i.cleanup = i.condition.Subscribe(func(v bool) {
		i.RequestBuild()
	})
}

// OnUnmount cleans up.
func (i *If) OnUnmount() {
	if i.cleanup != nil {
		i.cleanup()
		i.cleanup = nil
	}
}

func (i *If) show(cond bool) core.Element {
	var fn func() core.Element
	if cond {
		fn = i.thenFn
	} else {
		fn = i.elseFn
	}
	if fn == nil {
		return nil
	}
	return fn()
}
