package widget

import (
	"sync"

	"github.com/sjm1327605995/tenon/geometry"
)

// Unbinder is implemented by signal bindings for cleanup.
// It is defined here to avoid importing the state package from widget.
type Unbinder interface {
	Unbind()
}

// Stopper is implemented by effects for cleanup.
// It is defined here to avoid importing the state package from widget.
type Stopper interface {
	Stop()
}

// WidgetBase provides common functionality for widgets.
type WidgetBase struct {
	mu           sync.RWMutex
	bounds       geometry.Rect  // Cached layout bounds
	screenOrigin geometry.Point // Window-space origin, set during Draw pass
	focused      bool           // Whether widget has focus
	visible      bool           // Whether widget is visible
	enabled      bool           // Whether widget accepts input
	needsRedraw  bool           // Whether widget needs re-rendering
	id           string         // Optional ID for debugging
	children     []Widget       // Child widgets
	parent       Widget         // Parent widget (if any)
	bindings     []Unbinder     // Signal bindings (cleaned up on unmount)
	effects      []Stopper      // Effects (stopped on unmount)
	mounted      bool           // Whether widget is currently mounted
}

// NewWidgetBase creates a new WidgetBase with default settings.
func NewWidgetBase() *WidgetBase {
	return &WidgetBase{
		visible:     true,
		enabled:     true,
		needsRedraw: true,
	}
}

// Bounds returns the widget's current bounds (position and size).
func (w *WidgetBase) Bounds() geometry.Rect {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.bounds
}

// SetBounds sets the widget's bounds.
func (w *WidgetBase) SetBounds(bounds geometry.Rect) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.bounds = bounds
}

// Size returns the widget's current size.
func (w *WidgetBase) Size() geometry.Size {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.bounds.Size()
}

// Position returns the widget's top-left position.
func (w *WidgetBase) Position() geometry.Point {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.bounds.Min
}

// IsFocused returns true if the widget currently has focus.
func (w *WidgetBase) IsFocused() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.focused
}

// SetFocused sets the widget's focus state.
func (w *WidgetBase) SetFocused(focused bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.focused = focused
}

// IsVisible returns true if the widget is visible.
func (w *WidgetBase) IsVisible() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.visible
}

// SetVisible sets the widget's visibility.
func (w *WidgetBase) SetVisible(visible bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.visible = visible
}

// IsEnabled returns true if the widget accepts input.
func (w *WidgetBase) IsEnabled() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.enabled
}

// SetEnabled sets whether the widget accepts input.
func (w *WidgetBase) SetEnabled(enabled bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.enabled = enabled
}

// ID returns the widget's ID for debugging purposes.
func (w *WidgetBase) ID() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.id
}

// SetID sets the widget's ID for debugging purposes.
func (w *WidgetBase) SetID(id string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.id = id
}

// Parent returns the widget's parent, or nil if none.
func (w *WidgetBase) Parent() Widget {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.parent
}

// SetParent sets the widget's parent.
func (w *WidgetBase) SetParent(parent Widget) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.parent = parent
}

// Children returns the widget's child widgets.
func (w *WidgetBase) Children() []Widget {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if len(w.children) == 0 {
		return nil
	}
	result := make([]Widget, len(w.children))
	copy(result, w.children)
	return result
}

// ChildCount returns the number of child widgets.
func (w *WidgetBase) ChildCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.children)
}

// ChildAt returns the child at the given index, or nil if out of range.
func (w *WidgetBase) ChildAt(index int) Widget {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if index < 0 || index >= len(w.children) {
		return nil
	}
	return w.children[index]
}

// AddChild adds a child widget.
func (w *WidgetBase) AddChild(child Widget) {
	if child == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.children = append(w.children, child)
}

// RemoveChild removes a child widget.
func (w *WidgetBase) RemoveChild(child Widget) bool {
	if child == nil {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	for i, c := range w.children {
		if c != child {
			continue
		}
		lastIdx := len(w.children) - 1
		w.children[i] = w.children[lastIdx]
		w.children[lastIdx] = nil
		w.children = w.children[:lastIdx]
		return true
	}
	return false
}

// RemoveChildAt removes the child at the given index.
func (w *WidgetBase) RemoveChildAt(index int) Widget {
	w.mu.Lock()
	defer w.mu.Unlock()
	if index < 0 || index >= len(w.children) {
		return nil
	}
	child := w.children[index]
	copy(w.children[index:], w.children[index+1:])
	w.children[len(w.children)-1] = nil
	w.children = w.children[:len(w.children)-1]
	return child
}

// ClearChildren removes all child widgets.
func (w *WidgetBase) ClearChildren() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := range w.children {
		w.children[i] = nil
	}
	w.children = w.children[:0]
}

// InsertChild inserts a child widget at the given index.
func (w *WidgetBase) InsertChild(index int, child Widget) {
	if child == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if index < 0 {
		index = 0
	}
	if index >= len(w.children) {
		w.children = append(w.children, child)
		return
	}
	w.children = append(w.children, nil)
	copy(w.children[index+1:], w.children[index:])
	w.children[index] = child
}

// HasChildren returns true if the widget has any children.
func (w *WidgetBase) HasChildren() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.children) > 0
}

// ContainsPoint returns true if the point is within the widget's bounds.
func (w *WidgetBase) ContainsPoint(p geometry.Point) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.bounds.Contains(p)
}

// ScreenOrigin returns the widget's top-left corner in window (screen) coordinates.
func (w *WidgetBase) ScreenOrigin() geometry.Point {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.screenOrigin
}

// SetScreenOrigin records the widget's window-space origin.
func (w *WidgetBase) SetScreenOrigin(origin geometry.Point) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.screenOrigin = origin
}

// ScreenBounds returns the widget's bounds in window (screen) coordinates.
func (w *WidgetBase) ScreenBounds() geometry.Rect {
	w.mu.RLock()
	defer w.mu.RUnlock()
	size := w.bounds.Size()
	return geometry.FromPointSize(w.screenOrigin, size)
}

// LocalToGlobal converts a point from local coordinates to global (window) coordinates.
func (w *WidgetBase) LocalToGlobal(p geometry.Point) geometry.Point {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return p.Add(w.screenOrigin)
}

// GlobalToLocal converts a point from global (window) coordinates to local coordinates.
func (w *WidgetBase) GlobalToLocal(p geometry.Point) geometry.Point {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return p.Sub(w.screenOrigin)
}

// IsMounted reports whether the widget is currently in the mounted tree.
func (w *WidgetBase) IsMounted() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.mounted
}

// SetMounted sets the widget's mounted state.
func (w *WidgetBase) SetMounted(m bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.mounted = m
}

// NeedsRedraw reports whether the widget needs re-rendering.
func (w *WidgetBase) NeedsRedraw() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.needsRedraw
}

// SetNeedsRedraw marks the widget as needing re-rendering.
func (w *WidgetBase) SetNeedsRedraw(v bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.needsRedraw = v
}

// ClearRedraw clears the widget's needsRedraw flag after a successful draw.
func (w *WidgetBase) ClearRedraw() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.needsRedraw = false
}

// AddBinding registers a signal binding for automatic cleanup on unmount.
func (w *WidgetBase) AddBinding(b Unbinder) {
	if b == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.bindings = append(w.bindings, b)
}

// AddEffect registers an effect for automatic cleanup on unmount.
func (w *WidgetBase) AddEffect(e Stopper) {
	if e == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.effects = append(w.effects, e)
}

// CleanupBindings unbinds all signal bindings and stops all effects.
func (w *WidgetBase) CleanupBindings() {
	w.mu.Lock()
	bindings := w.bindings
	effects := w.effects
	w.bindings = nil
	w.effects = nil
	w.mu.Unlock()

	for _, b := range bindings {
		b.Unbind()
	}
	for _, e := range effects {
		e.Stop()
	}
}
