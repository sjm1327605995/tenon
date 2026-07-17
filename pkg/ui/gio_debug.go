package ui

import (
	"fmt"
	"os"

	"gioui.org/io/pointer"
)

// 诊断开关：TENON_DEBUG_INPUT=1 时把每帧的输入与悬停决策打到 stdout。
// 用来定位「看得见但复现不出」的交互问题（无头测试里引擎状态与像素都稳定时，
// 就只能看真机上的真实事件流）。默认关闭，零开销。
var debugInput = os.Getenv("TENON_DEBUG_INPUT") != ""

// debugPointer 打印收到的指针事件。
func debugPointer(ev pointer.Event) {
	if !debugInput {
		return
	}
	fmt.Printf("PTR  kind=%-8v pos=%.2f,%.2f btns=%v\n", ev.Kind, ev.Position.X, ev.Position.Y, ev.Buttons)
}

// debugHover 打印一次悬停链变化（进入/离开）。
func debugHover(what string, rn *renderNode, x, y float32) {
	if !debugInput {
		return
	}
	fmt.Printf("HOVER %-5s node=%p kind=%v bounds=%.1f,%.1f %.1fx%.1f cursor=%.2f,%.2f\n",
		what, rn, rn.kind, rn.bounds.X, rn.bounds.Y, rn.bounds.W, rn.bounds.H, x, y)
}

// debugFrame 打印一帧结束时的状态与是否还要下一帧。
func debugFrame(g *game, wantNext bool) {
	if !debugInput {
		return
	}
	x, y := input.cursor()
	fmt.Printf("FRAME cursor=%.2f,%.2f hovered=%d dirty=%d anims=%d loops=%d needsLayout=%v next=%v\n",
		x, y, len(g.hovered), len(g.dirty), len(g.anims), len(g.loops), g.needsLayout, wantNext)
}
