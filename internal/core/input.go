package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Key 表示键盘按键的抽象标识符。
// 该类型与具体输入后端（如 ebiten）解耦，
// 底层映射由 EventDispatcher 在事件分发时完成。
type Key int

// MouseButton 表示鼠标按键的抽象标识符。
type MouseButton int

const (
	MouseButtonLeft   MouseButton = 0
	MouseButtonRight  MouseButton = 1
	MouseButtonMiddle MouseButton = 2
)

// 常用键盘按键常量。
// 值保持与主流输入系统兼容，但代码不应依赖具体数值。
const (
	KeyA Key = iota + 1
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySpace
	KeyEnter
	KeyTab
	KeyEscape
	KeyBackspace
	KeyDelete
	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyHome
	KeyEnd
	KeyShift
	KeyCtrl
	KeyAlt
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)

// IsMouseButtonPressed 检查指定鼠标按键是否被按下。
func IsMouseButtonPressed(btn MouseButton) bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButton(btn))
}

// IsKeyPressed 检查指定键盘按键是否被按下。
func IsKeyPressed(key Key) bool {
	return ebiten.IsKeyPressed(ebiten.Key(key))
}

// IsKeyJustPressed 检查指定键盘按键是否在本帧被按下。
func IsKeyJustPressed(key Key) bool {
	return inpututil.IsKeyJustPressed(ebiten.Key(key))
}
