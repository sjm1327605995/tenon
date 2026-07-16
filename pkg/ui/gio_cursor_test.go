package ui

import (
	"testing"

	"gioui.org/io/pointer"
)

// 光标落在输入框上应是 I 型，落在可点击元素上应是手型，其余为默认。
func TestCursorShapeFollowsElement(t *testing.T) {
	app := func(_ struct{}) *Node {
		return Div(Style(Row, Fill),
			Div(Style(Width(100), Height(50)), OnClick(func() {}), Text("btn")),
			Input(Style(Width(100), Height(50))),
			Div(Style(Grow(1), Height(50))), // 普通区域
		)
	}
	h := Mount(Use(app, struct{}{}), 400, 100)
	g := h.g

	cases := []struct {
		name   string
		x, y   int
		want   pointer.Cursor
		reason string
	}{
		{"可点击", 50, 25, pointer.CursorPointer, "按钮上应是手型"},
		{"输入框", 150, 25, pointer.CursorText, "文本框上应是 I 型"},
		{"空白", 300, 25, pointer.CursorDefault, "普通区域应是默认指针"},
	}
	for _, c := range cases {
		gioIn.curX, gioIn.curY = c.x, c.y
		if got := gioCursor(g); got != c.want {
			t.Errorf("%s(%d,%d): cursor=%v want %v —— %s", c.name, c.x, c.y, got, c.want, c.reason)
		}
	}
}

// 输入框套在可点击容器里时，最内层语义优先（应显示 I 型而不是手型）。
func TestCursorInnermostWins(t *testing.T) {
	app := func(_ struct{}) *Node {
		return Div(Style(Row, Fill), OnClick(func() {}),
			Input(Style(Width(100), Height(50))),
		)
	}
	h := Mount(Use(app, struct{}{}), 200, 100)
	gioIn.curX, gioIn.curY = 50, 25
	if got := gioCursor(h.g); got != pointer.CursorText {
		t.Fatalf("cursor=%v want CursorText（输入框在可点击容器内，最内层应优先）", got)
	}
}
