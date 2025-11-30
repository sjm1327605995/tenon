package event

import "image"

// Event 是所有事件的基础接口
type Event interface {
	ImplementEvent()
}

// UpdateWindowsSizeEvent 窗口大小更新事件
type UpdateWindowsSizeEvent struct {
	Size image.Point
}

func (u UpdateWindowsSizeEvent) ImplementEvent() {
}

// DestroyEvent 销毁事件
type DestroyEvent struct {
	Err error
}

func (d DestroyEvent) ImplementEvent() {
}
