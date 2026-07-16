package ui

import "time"

// ---- 后端中立的输入抽象 ----
//
// 引擎的输入分发（handleInput 及其子过程）只经 input 读取指针/键盘状态，不直接依赖
// ebiten / gio。每个后端提供一个 inputSource 实现：ebiten 为轮询式（inpututil 差分），
// gio 由事件循环把每帧事件累积成同样可轮询的状态。

type mouseBtn int

const (
	btnLeft mouseBtn = iota
	btnRight
)

// inKey 是引擎用到的按键的中立枚举（各后端映射到自己的键码/键名）。
type inKey int

const (
	keyF12 inKey = iota
	keyTab
	keyShift
	keyCtrl
	keyMeta
	keyDown
	keyUp
	keyLeft
	keyRight
	keyEscape
	keyEnter
	keySpace
	keyA
	keyC
	keyX
	keyV
	keyBackspace
	keyDelete
	keyHome
	keyEnd
	inKeyCount
)

type inputSource interface {
	cursor() (int, int)             // 光标位置（物理像素）
	wheel() (float32, float32)      // 本帧滚轮增量（x,y）
	mousePressed(mouseBtn) bool     // 按钮当前是否按住
	mouseJustPressed(mouseBtn) bool // 按钮本帧是否刚按下
	keyPressed(inKey) bool          // 键当前是否按住
	keyJustPressed(inKey) bool      // 键本帧是否刚按下
	typedChars() []rune             // 本帧输入的字符（不含控制键）
}

// input 是当前后端的输入源，由 gio_run.go 的 init 设定。
var input inputSource

// 长按重复用墙钟计时（与刷新率脱钩）：首触发后延迟 repeatDelay，再每 repeatInterval 重复一次。
const (
	repeatDelay    = 450 * time.Millisecond
	repeatInterval = 33 * time.Millisecond // ~30 次/秒
)

// keyNextRepeat 记录每个键下次允许重复触发的时刻。
var keyNextRepeat = map[inKey]time.Time{}

// repeatKey 在按下瞬间触发一次，长按后按固定时间间隔重复（不依赖帧率）。
func repeatKey(k inKey) bool {
	if input.keyJustPressed(k) {
		keyNextRepeat[k] = time.Now().Add(repeatDelay)
		return true
	}
	if !input.keyPressed(k) {
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
