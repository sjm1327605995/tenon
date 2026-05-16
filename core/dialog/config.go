package dialog

import (
	"github.com/sjm1327605995/tenon/state"
	"github.com/sjm1327605995/tenon/widget"
)

// config holds the dialog's configuration, set at construction time via options.
type config struct {
	title               string
	titleFn             func() string
	titleSignal         state.Signal[string]
	readonlyTitleSignal state.ReadonlySignal[string]
	content             widget.Widget
	actions             []Action
	dismissible         bool
	escToClose          bool
	onClose             func()
	maxWidth            float32
	maxHeight           float32
	painter             Painter
}

// ResolvedTitle returns the current display title.
// Priority: ReadonlySignal > Signal > Fn > Static.
func (c *config) ResolvedTitle() string {
	if c.readonlyTitleSignal != nil {
		return c.readonlyTitleSignal.Get()
	}
	if c.titleSignal != nil {
		return c.titleSignal.Get()
	}
	if c.titleFn != nil {
		return c.titleFn()
	}
	return c.title
}
