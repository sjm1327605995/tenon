package engine

import "github.com/sjm1327605995/tenon/internal/render"

// ==================== RenderObject 能力检查辅助函数 ====================
// 将 engine.go 中散落的 interface{...} 断言集中管理。

// Scrollabler 是可滚动的 RenderObject。
type Scrollabler interface {
	render.RenderObject
	GetScrollOffset() render.Offset
	ScrollBy(dx, dy float32)
}

// tryGetScroll 尝试将 RenderObject 转为 Scrollabler。
func tryGetScroll(ro render.RenderObject) (Scrollabler, bool) {
	s, ok := ro.(Scrollabler)
	return s, ok
}

// Clipper 是可裁剪子节点的 RenderObject。
type Clipper interface {
	render.RenderObject
	ClipChildren() bool
}

// tryGetClipper 尝试将 RenderObject 转为 Clipper。
func tryGetClipper(ro render.RenderObject) (Clipper, bool) {
	c, ok := ro.(Clipper)
	return c, ok
}

// Focusabler 是可聚焦的 RenderObject。
type Focusabler interface {
	render.RenderObject
	Focus()
	Blur()
	IsFocused() bool
}

// tryGetFocuser 从 ro 向上搜索最近的 Focusabler 祖先。
func tryGetFocuser(ro render.RenderObject) (Focusabler, bool) {
	for ro != nil {
		if f, ok := ro.(Focusabler); ok {
			return f, true
		}
		ro = ro.GetParent()
	}
	return nil, false
}

// InputHandler 是可处理键盘输入的 RenderObject。
type InputHandler interface {
	render.RenderObject
	HandleInput(chars []rune, backspace, enter, left, right, home, end, selectAll bool)
}

// tryGetInputHandler 尝试将 RenderObject 转为 InputHandler。
func tryGetInputHandler(ro render.RenderObject) (InputHandler, bool) {
	h, ok := ro.(InputHandler)
	return h, ok
}

// TextInputHandler 是可处理 IME 输入的 RenderObject。
type TextInputHandler interface {
	render.RenderObject
	UpdateTextInput(absX, absY int) bool
}

// tryGetTextInputHandler 尝试将 RenderObject 转为 TextInputHandler。
func tryGetTextInputHandler(ro render.RenderObject) (TextInputHandler, bool) {
	h, ok := ro.(TextInputHandler)
	return h, ok
}

// Blinker 是有光标闪烁的 RenderObject。
type Blinker interface {
	render.RenderObject
	TickBlink()
}

// tryGetBlinker 尝试将 RenderObject 转为 Blinker。
func tryGetBlinker(ro render.RenderObject) (Blinker, bool) {
	b, ok := ro.(Blinker)
	return b, ok
}

// ClipboardHandler 是支持剪贴板操作的 RenderObject。
type ClipboardHandler interface {
	render.RenderObject
	Copy(clipboardWrite func(string)) string
	Cut(clipboardWrite func(string)) string
	Paste(clipboardRead func() string)
	HasSelection() bool
	SelectAll()
}

// tryGetClipboardHandler 尝试将 RenderObject 转为 ClipboardHandler。
func tryGetClipboardHandler(ro render.RenderObject) (ClipboardHandler, bool) {
	c, ok := ro.(ClipboardHandler)
	return c, ok
}

// scrollParent 在祖先中查找可滚动的 RenderObject。
func scrollParent(ro render.RenderObject) (Scrollabler, bool) {
	for ro != nil {
		if s, ok := ro.(Scrollabler); ok {
			return s, true
		}
		ro = ro.GetParent()
	}
	return nil, false
}
