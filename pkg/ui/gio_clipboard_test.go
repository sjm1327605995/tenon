package ui

import "testing"

func focusedInputGame(value string, caret, anchor int, multiline bool) (*game, *renderNode) {
	rn := &renderNode{kind: rnInput, value: value, caretPos: caret, selAnchor: anchor, multiline: multiline}
	f := &Fiber{rnode: rn}
	rn.owner = f
	return &game{focusedFiber: f}, rn
}

// 有选区时粘贴应替换选区，无选区时插到光标处。
func TestPasteReplacesSelection(t *testing.T) {
	g, rn := focusedInputGame("abcd", 3, 1, false) // 选中 "bc"
	var got string
	rn.onChange = func(s string) { got = s }

	pasteIntoFocused(g, "XY")

	if rn.value != "aXYd" {
		t.Fatalf("value=%q want %q", rn.value, "aXYd")
	}
	if got != "aXYd" {
		t.Fatalf("onChange=%q want %q", got, "aXYd")
	}
	if rn.caretPos != 3 || rn.selAnchor != 3 {
		t.Fatalf("粘贴后光标应折叠到插入末尾: caret=%d anchor=%d want 3/3", rn.caretPos, rn.selAnchor)
	}
}

func TestPasteInsertsAtCaret(t *testing.T) {
	g, rn := focusedInputGame("ad", 1, 1, false)
	pasteIntoFocused(g, "bc")
	if rn.value != "abcd" {
		t.Fatalf("value=%q want %q", rn.value, "abcd")
	}
}

// 单行输入框不接受换行：从别处复制的多行文本会把布局撑坏，应折成空格。
func TestPasteStripsNewlinesInSingleLine(t *testing.T) {
	g, rn := focusedInputGame("", 0, 0, false)
	pasteIntoFocused(g, "a\r\nb\nc")
	if rn.value != "a b c" {
		t.Fatalf("单行粘贴 value=%q want %q", rn.value, "a b c")
	}
}

// 多行输入框应保留换行。
func TestPasteKeepsNewlinesInMultiline(t *testing.T) {
	g, rn := focusedInputGame("", 0, 0, true)
	pasteIntoFocused(g, "a\nb")
	if rn.value != "a\nb" {
		t.Fatalf("多行粘贴 value=%q want %q", rn.value, "a\nb")
	}
}

// 没有聚焦输入框时粘贴应无副作用。
func TestPasteWithoutFocusIsNoop(t *testing.T) {
	g := &game{}
	pasteIntoFocused(g, "x") // 不应 panic
}

// 复制/剪切要排队写入系统剪贴板，同时更新应用内缓存。
func TestSetClipboardQueuesSystemWrite(t *testing.T) {
	gioClip.pendingWrite = nil
	setClipboard("hello")
	if clipboardText != "hello" {
		t.Fatalf("应用内缓存=%q want hello", clipboardText)
	}
	if gioClip.pendingWrite == nil || *gioClip.pendingWrite != "hello" {
		t.Fatal("没有排队写入系统剪贴板 —— 复制的内容不会出现在其它程序里")
	}
	gioClip.pendingWrite = nil
}
