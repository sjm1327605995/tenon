package event

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon/geometry"
)

// TouchEventType represents the specific type of touch event.
type TouchEventType uint8

const (
	// TouchPress indicates a touch began.
	TouchPress TouchEventType = iota + 1

	// TouchRelease indicates a touch ended.
	TouchRelease

	// TouchMove indicates a touch moved.
	TouchMove
)

// TouchEvent represents a touch input event.
type TouchEvent struct {
	Base

	// TouchType is the specific type of touch event.
	TouchType TouchEventType

	// ID is the unique identifier for this touch point.
	ID int

	// Position is the touch position relative to the widget.
	Position geometry.Point

	// GlobalPosition is the touch position in screen coordinates.
	GlobalPosition geometry.Point
}

// NewTouchEvent creates a new touch event with the current time.
func NewTouchEvent(
	touchType TouchEventType,
	id int,
	position geometry.Point,
	globalPosition geometry.Point,
	mods Modifiers,
) *TouchEvent {
	return &TouchEvent{
		Base:           NewBase(TypeTouch, mods),
		TouchType:      touchType,
		ID:             id,
		Position:       position,
		GlobalPosition: globalPosition,
	}
}

// NewTouchEventWithTime creates a new touch event with a specific timestamp.
func NewTouchEventWithTime(
	touchType TouchEventType,
	id int,
	position geometry.Point,
	globalPosition geometry.Point,
	mods Modifiers,
	t time.Time,
) *TouchEvent {
	return &TouchEvent{
		Base:           NewBaseWithTime(TypeTouch, t, mods),
		TouchType:      touchType,
		ID:             id,
		Position:       position,
		GlobalPosition: globalPosition,
	}
}

// IsPress returns true if this is a touch press event.
func (e *TouchEvent) IsPress() bool {
	return e.TouchType == TouchPress
}

// IsRelease returns true if this is a touch release event.
func (e *TouchEvent) IsRelease() bool {
	return e.TouchType == TouchRelease
}

// IsMove returns true if this is a touch move event.
func (e *TouchEvent) IsMove() bool {
	return e.TouchType == TouchMove
}

// String returns a human-readable representation of the event.
func (e *TouchEvent) String() string {
	return fmt.Sprintf("TouchEvent{Type: %d, ID: %d, Position: %s, Mods: %s}",
		e.TouchType, e.ID, e.Position, e.Modifiers())
}
