// Package event defines the event system for the Tenon React framework.
// This package contains interfaces and structs for handling various UI events in the application.
package event

import "image"

// Event is the base interface for all event types in the framework.
// Any type that implements this interface can be used as an event in the event system.
// The ImplementEvent method serves as a marker to identify event types.

type Event interface {
	ImplementEvent()
}

// UpdateWindowsSizeEvent represents an event that occurs when the window size changes.
// This event contains the new size of the window.
//
// Fields:
//   - Size: The new size of the window as an image.Point (Width, Height)
type UpdateWindowsSizeEvent struct {
	Size image.Point
}

// ImplementEvent implements the Event interface for UpdateWindowsSizeEvent.
// This is a marker method required by the Event interface.
func (u UpdateWindowsSizeEvent) ImplementEvent() {
	// No specific implementation needed, this is a marker method
}
