package event

import "image"

type Event interface {
	ImplementEvent()
}
type UpdateWindowsSizeEvent struct {
	Size image.Point
}

func (u UpdateWindowsSizeEvent) ImplementEvent() {

}
