package ui

import "time"

// ---- 后端中立的输入抽象 ----
//
// 引擎的输入分发（handleInput 及其子过程）只经 InputSource 读取指针/键盘状态，不直接依赖
// Gio。每个后端提供一个实现；Gio 由事件循环把每帧事件累积成可轮询的状态。

// MouseButton 是后端中立的鼠标按键。
type MouseButton int

const (
	MouseLeft MouseButton = iota
	MouseRight
)

type mouseBtn = MouseButton

const (
	btnLeft  = MouseLeft
	btnRight = MouseRight
)

// InputKey 是引擎用到的按键的中立枚举（各后端映射到自己的键码/键名）。
type InputKey int

const (
	KeyF12 InputKey = iota
	KeyTab
	KeyShift
	KeyCtrl
	KeyMeta
	KeyDown
	KeyUp
	KeyLeft
	KeyRight
	KeyEscape
	KeyEnter
	KeySpace
	KeyA
	KeyC
	KeyX
	KeyV
	KeyBackspace
	KeyDelete
	KeyHome
	KeyEnd
	keyCount
)

type inKey = InputKey

const (
	keyF12       = KeyF12
	keyTab       = KeyTab
	keyShift     = KeyShift
	keyCtrl      = KeyCtrl
	keyMeta      = KeyMeta
	keyDown      = KeyDown
	keyUp        = KeyUp
	keyLeft      = KeyLeft
	keyRight     = KeyRight
	keyEscape    = KeyEscape
	keyEnter     = KeyEnter
	keySpace     = KeySpace
	keyA         = KeyA
	keyC         = KeyC
	keyX         = KeyX
	keyV         = KeyV
	keyBackspace = KeyBackspace
	keyDelete    = KeyDelete
	keyHome      = KeyHome
	keyEnd       = KeyEnd
	inKeyCount   = keyCount
)

// InputSource 提供一帧内的输入快照。边沿状态（JustPressed、字符、滚轮）由后端在每帧重置。
type InputSource interface {
	// Cursor 返回光标位置：物理像素，保留亚像素精度。
	Cursor() (float32, float32)
	Wheel() (float32, float32)         // 本帧滚轮增量，物理像素，正值向右/向下
	MousePressed(MouseButton) bool     // 按钮当前是否按住
	MouseJustPressed(MouseButton) bool // 按钮本帧是否刚按下
	KeyPressed(InputKey) bool          // 键当前是否按住
	KeyJustPressed(InputKey) bool      // 键本帧是否刚按下
	TypedChars() []rune                // 本帧输入的字符（不含控制键）
}

type inputSource = InputSource

var input inputSource

func setGlobalInput(source InputSource) InputSource {
	previous := input
	input = source
	return previous
}

// 长按重复用墙钟计时（与刷新率脱钩）：首触发后延迟 repeatDelay，再每 repeatInterval 重复一次。
const (
	repeatDelay    = 450 * time.Millisecond
	repeatInterval = 33 * time.Millisecond // ~30 次/秒
)

// keyNextRepeat 记录每个键下次允许重复触发的时刻。
var keyNextRepeat = map[inKey]time.Time{}

// repeatKey 在按下瞬间触发一次，长按后按固定时间间隔重复（不依赖帧率）。
func repeatKey(k inKey) bool {
	if input.KeyJustPressed(k) {
		keyNextRepeat[k] = time.Now().Add(repeatDelay)
		return true
	}
	if !input.KeyPressed(k) {
		delete(keyNextRepeat, k)
		return false
	}
	next, ok := keyNextRepeat[k]
	if !ok { // 聚焦时键已被按住：先建立计时，不立即触发
		keyNextRepeat[k] = time.Now().Add(repeatDelay)
		return false
	}
	now := time.Now()
	if now.Before(next) {
		return false
	}
	keyNextRepeat[k] = now.Add(repeatInterval)
	return true
}
