package core

import "github.com/hajimehoshi/ebiten/v2"

// EventType 定义事件类型。
type EventType int

const (
	EventMouseMove EventType = iota
	EventMouseDown
	EventMouseUp
	EventClick
	EventScroll
	EventKeyDown
	EventKeyUp
	EventFocusIn
	EventFocusOut
	EventMouseEnter
	EventMouseLeave
)

// Event 是统一的 UI 事件结构。
type Event struct {
	Type     EventType
	X, Y     float32 // 窗口坐标
	LocalX   float32 // 目标组件本地坐标
	LocalY   float32 // 目标组件本地坐标
	DeltaX   float32 // 滚轮水平增量
	DeltaY   float32 // 滚轮垂直增量
	Key      ebiten.Key
	Button   ebiten.MouseButton
	Target   Host
}
