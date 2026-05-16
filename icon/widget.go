package icon

import (
	"github.com/sjm1327605995/tenon/event"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/state"
	"github.com/sjm1327605995/tenon/widget"
)

// defaultIconSize is the default display size in logical pixels.
const defaultIconSize float32 = 24

// IconOption configures an [IconWidget].
type IconOption func(*iconConfig)

// iconConfig holds configuration for constructing an IconWidget.
type iconConfig struct {
	size        float32
	color       widget.Color
	colorSignal state.ReadonlySignal[widget.Color]
	label       string
}

// Size sets the display size of the icon in logical pixels.
func Size(s float32) IconOption {
	return func(c *iconConfig) { c.size = s }
}

// Color sets the icon stroke color.
func Color(color widget.Color) IconOption {
	return func(c *iconConfig) { c.color = color }
}

// ColorSignal binds the icon color to a reactive signal.
func ColorSignal(sig state.ReadonlySignal[widget.Color]) IconOption {
	return func(c *iconConfig) { c.colorSignal = sig }
}

// Label sets the accessibility label for the icon.
func Label(text string) IconOption {
	return func(c *iconConfig) { c.label = text }
}

// IconWidget is a display widget that renders a vector icon.
type IconWidget struct {
	widget.WidgetBase

	data        IconData
	size        float32
	color       widget.Color
	colorSignal state.ReadonlySignal[widget.Color]
	label       string
}

// NewIcon creates a new icon widget displaying the given icon data.
func NewIcon(data IconData, opts ...IconOption) *IconWidget {
	cfg := iconConfig{
		size:  defaultIconSize,
		color: widget.ColorBlack,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	w := &IconWidget{
		data:        data,
		size:        cfg.size,
		color:       cfg.color,
		colorSignal: cfg.colorSignal,
		label:       cfg.label,
	}
	w.SetVisible(true)
	w.SetEnabled(true)
	return w
}

// Data returns the icon data being displayed.
func (w *IconWidget) Data() IconData {
	return w.data
}

// IconSize returns the display size in logical pixels.
func (w *IconWidget) IconSize() float32 {
	return w.size
}

// IconColor returns the resolved icon color.
func (w *IconWidget) IconColor() widget.Color {
	if w.colorSignal != nil {
		return w.colorSignal.Get()
	}
	return w.color
}

// IconLabel returns the accessibility label.
func (w *IconWidget) IconLabel() string {
	return w.label
}

// Layout returns the icon's preferred size.
func (w *IconWidget) Layout(_ widget.Context, constraints geometry.Constraints) geometry.Size {
	preferred := geometry.Sz(w.size, w.size)
	resultSize := constraints.Constrain(preferred)
	w.SetBounds(geometry.FromPointSize(w.Position(), resultSize))
	return resultSize
}

// Draw renders the icon to the canvas.
func (w *IconWidget) Draw(_ widget.Context, canvas widget.Canvas) {
	if !w.IsVisible() {
		return
	}
	bounds := w.Bounds()
	if bounds.IsEmpty() {
		return
	}
	Draw(canvas, w.data, bounds, w.IconColor())
}

// Event returns false. Icon widgets do not consume events.
func (w *IconWidget) Event(_ widget.Context, _ event.Event) bool {
	return false
}

// Children returns nil. Icon is a leaf widget.
func (w *IconWidget) Children() []widget.Widget {
	return nil
}

// Mount creates signal bindings for push-based invalidation.
func (w *IconWidget) Mount(ctx widget.Context) {
	sched := ctx.Scheduler()
	if sched == nil {
		return
	}
	if w.colorSignal != nil {
		b := state.BindToScheduler(w.colorSignal, w, sched)
		w.AddBinding(b)
	}
}

// Unmount is called when the icon widget is removed from the widget tree.
func (w *IconWidget) Unmount() {
	// Bindings are cleaned up automatically by WidgetBase.CleanupBindings().
}

// Compile-time interface checks.
var (
	_ widget.Widget    = (*IconWidget)(nil)
	_ widget.Lifecycle = (*IconWidget)(nil)
)
