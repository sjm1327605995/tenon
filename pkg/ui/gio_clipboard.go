package ui

import (
	"io"
	"strings"

	"gioui.org/io/clipboard"
	gioinput "gioui.org/io/input"
)

// ---- gio 后端：系统剪贴板 ----
//
// 引擎默认的剪贴板只在应用内有效，跟别的程序不通。这里把它接到系统剪贴板。
//
// 难点在于读写不对称：gio 的写是同步命令，读却是异步的（ReadCmd -> 之后才来
// transfer.DataEvent），而引擎的 getClipboard() 是同步的。若让引擎同步读缓存，
// 在别的程序里复制后过来粘贴就会粘到过期内容。所以粘贴由后端接管：
//
//	Ctrl+V -> 压住引擎本帧的 keyV，发 ReadCmd
//	         -> DataEvent 回来后直接落到聚焦输入框（替换选区）
//
// 复制/剪切仍走引擎的 setClipboard，只是在这里排队成 WriteCmd。

// gioClipMime 是文本传输的 MIME 类型（与 gio 的 widget.Editor 一致）。
const gioClipMime = "application/text"

type gioClipboardState struct {
	pendingWrite   *string // 待写入系统剪贴板的文本（复制/剪切）
	pasteRequested bool    // 本帧按下了 Ctrl+V，待发 ReadCmd
	pastedText     *string // 系统剪贴板回来的文本，待落到输入框
}

var gioClip gioClipboardState

func init() {
	// 复制/剪切：登记到默认 provider 的后端钩子上，排队写入系统剪贴板。
	backendWriteClipboard = func(s string) {
		t := s
		gioClip.pendingWrite = &t
		if backendWake != nil {
			backendWake() // 可能发生在静止帧里，需要叫醒循环去执行 WriteCmd
		}
	}
}

// flush 在帧内执行排队的剪贴板命令，并把读回的文本落到聚焦输入框。
func (c *gioClipboardState) flush(src gioinput.Source, g *game) {
	if w := c.pendingWrite; w != nil {
		c.pendingWrite = nil
		src.Execute(clipboard.WriteCmd{
			Type: gioClipMime,
			Data: io.NopCloser(strings.NewReader(*w)),
		})
	}
	if c.pasteRequested {
		c.pasteRequested = false
		src.Execute(clipboard.ReadCmd{Tag: gioTag})
	}
	if t := c.pastedText; t != nil {
		c.pastedText = nil
		clipboardText = *t
		pasteIntoFocused(g, *t)
	}
}

// pasteIntoFocused 把文本粘到聚焦输入框：有选区则替换，否则插到光标处。
func pasteIntoFocused(g *game, text string) {
	rn := focusedInput(g)
	if rn == nil || text == "" {
		return
	}
	val := rn.value
	lo := clampi(min(rn.selAnchor, rn.caretPos), 0, len(val))
	hi := clampi(max(rn.selAnchor, rn.caretPos), 0, len(val))
	if !rn.multiline { // 单行输入框不接受换行
		text = strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", " "), "\n", " ")
	}
	nv := val[:lo] + text + val[hi:]
	caret := lo + len(text)
	rn.value = nv
	rn.caretPos, rn.selAnchor = caret, caret
	if rn.onChange != nil {
		rn.onChange(nv)
	}
}
