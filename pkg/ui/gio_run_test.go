package ui

import "testing"

// hasLayoutAnim 只表示「树里存在 Animated 节点」，一旦置位就再也不复位，因此不能拿它
// 当「是否还要出下一帧」的判据 —— 否则只要用过一次 Animated，界面就永远醒着重绘，
// 按需重绘等于白做。要用 tickLayoutAnim 的返回值（是否真的还在移动）。
func TestLayoutAnimFlagIsNotAnActivitySignal(t *testing.T) {
	g := &game{}
	rn := &renderNode{animatedLayout: true, bounds: Rect{X: 10, Y: 10}}
	g.rootRN = rn
	g.hasLayoutAnim = true

	// 第一帧：记录基准位置，还没有位移
	g.tickLayoutAnim(16)
	// 位置没再变化 -> 残余偏移衰减到 0 -> 不应再声称在动
	moving := false
	for i := 0; i < 60; i++ {
		moving = g.tickLayoutAnim(16)
	}
	if moving {
		t.Fatal("位置不再变化后 tickLayoutAnim 仍报告在移动，界面将永远无法进入静止")
	}
	if !g.hasLayoutAnim {
		t.Fatal("前提失效：hasLayoutAnim 本就不会复位，所以它不能作为活动信号")
	}
}

// 光标闪烁半周期必须与 caretVisible() 的判据一致，否则定时唤醒会和实际闪烁错拍。
func TestCaretBlinkPeriodMatchesCaretVisible(t *testing.T) {
	if caretBlinkPeriod.Milliseconds() != 500 {
		t.Fatalf("caretBlinkPeriod=%v，但 caretVisible() 是按 500ms 判定的", caretBlinkPeriod)
	}
}
