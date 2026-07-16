package ui

import (
	"testing"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
)

// 必须注册 key.FocusFilter：gio 只用它把 handler 标成 focusable，而
//   - key.EditEvent（字符输入 / IME 提交文本）只投递给 focusable 的 handler；
//   - keyQueue.Frame 每帧会把焦点从「不 focusable」的 handler 上撤掉。
//
// 少了它的后果是「完全打不了字」+「焦点每帧被撤销又重抢」。key.Filter 不会设置
// focusable，所以两者都得注册 —— 这个测试就是钉住这一点。
func TestGioInputFiltersIncludeFocusFilter(t *testing.T) {
	var hasFocus, hasPointer, hasKey bool
	for _, f := range gioInputFilters() {
		switch f := f.(type) {
		case key.FocusFilter:
			if f.Target == gioTag {
				hasFocus = true
			}
		case pointer.Filter:
			if f.Target == gioTag {
				hasPointer = true
			}
		case key.Filter:
			if f.Focus == gioTag {
				hasKey = true
			}
		}
	}
	if !hasFocus {
		t.Error("缺少 key.FocusFilter：会导致收不到 EditEvent（打不了字）且焦点每帧被撤销")
	}
	if !hasPointer {
		t.Error("缺少 pointer.Filter：收不到鼠标事件")
	}
	if !hasKey {
		t.Error("缺少 key.Filter：收不到导航/快捷键")
	}
}

// EditEvent 的文本应进入本帧的 typedChars，并在下一帧开头被清掉（边沿状态）。
func TestGioEditEventFeedsTypedChars(t *testing.T) {
	g := &gioInput{}
	g.resetFrame()
	if len(g.typedChars()) != 0 {
		t.Fatal("新一帧开头 typedChars 应为空")
	}
	// 模拟 gio 投递的字符输入
	g.typed = append(g.typed, []rune("你好a")...)
	if got := string(g.typedChars()); got != "你好a" {
		t.Fatalf("typedChars=%q want %q", got, "你好a")
	}
	g.resetFrame()
	if got := string(g.typedChars()); got != "" {
		t.Fatalf("resetFrame 后 typedChars 应清空，得到 %q", got)
	}
}
