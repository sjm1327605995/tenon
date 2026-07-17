package ui

// ---- 无头测试用的输入源 ----
//
// Harness 过去靠直接调 g.focusNext()/g.activateFocused() 来「按键」，那是与生产平行的
// 另一条路：真实帧循环走的是 handleKeyboardNav + editFocusedInput，读的是 inputSource。
// 两条路一分岔，测试就会给出与用户实际操作不同的答案 —— OnSubmit 的双重提交就是这么
// 溜过去的（harness 只跑了 activateFocused，没跑真实按键路径，全绿）。
//
// harnessInput 让测试可以真的「按下一个键」，然后跑生产同款的处理顺序。

type harnessInput struct {
	x, y    float32
	just    [inKeyCount]bool
	held    [inKeyCount]bool
	chars   []rune
	wheelX  float32
	wheelY  float32
	btnHeld [3]bool
	btnJust [3]bool
}

func (i *harnessInput) cursor() (float32, float32) { return i.x, i.y }
func (i *harnessInput) wheel() (float32, float32)  { return i.wheelX, i.wheelY }
func (i *harnessInput) mousePressed(b mouseBtn) bool {
	return int(b) < len(i.btnHeld) && i.btnHeld[b]
}
func (i *harnessInput) mouseJustPressed(b mouseBtn) bool {
	return int(b) < len(i.btnJust) && i.btnJust[b]
}
func (i *harnessInput) keyPressed(k inKey) bool     { return int(k) < len(i.held) && i.held[k] }
func (i *harnessInput) keyJustPressed(k inKey) bool { return int(k) < len(i.just) && i.just[k] }
func (i *harnessInput) typedChars() []rune          { return i.chars }

// pressKey 模拟按下一个键并跑一遍生产帧循环里处理按键的那两步，顺序与 run.go 的
// updateInput 一致：handleKeyboardNav（Tab/Esc/Enter 激活）在前，editFocusedInput
// （文本编辑、单行回车提交）在后。顺序必须与生产一致，否则测不出两者的相互作用。
func (h *Harness) pressKey(k inKey) {
	old := input
	hi := &harnessInput{}
	hi.just[k], hi.held[k] = true, true
	input = hi
	defer func() { input = old }()

	h.g.handleKeyboardNav()
	h.g.editFocusedInput()
	h.settle()
}
