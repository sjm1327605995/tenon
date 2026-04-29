package core

// DebugEventBus 是调试事件总线接口，用于解耦 Engine 与 Debugger 实现。
type DebugEventBus interface {
	CaptureLayout(trigger string)
	IsEnabled() bool
	AddEventLog(evt interface {
		GetType() string
		GetTarget() string
		GetX() float32
		GetY() float32
		GetDeltaX() float32
		GetDeltaY() float32
	})
	AddStateLog(elementType, elementKey, stateKey string, oldValue, newValue interface{})
}
