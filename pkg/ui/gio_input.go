package ui

import (
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
)

// ---- gio 后端：输入 ----
//
// gio 是事件式输入；这里把每帧收到的 pointer/key/edit 事件累积成与 ebiten 一致的可轮询状态，
// 从而复用引擎的 handleInput。持久状态（光标/按住/修饰键）跨帧保留，边沿状态（刚按下/滚轮/
// 输入字符）每帧开头清零、由本帧事件填充。

// gioTag 是整窗输入区域与键盘焦点的标识。
var gioTag = new(int)

type gioInput struct {
	curX, curY     int
	wheelX, wheelY float32
	btnDown        [2]bool // btnLeft/btnRight 当前按住
	btnJust        [2]bool // 本帧刚按下
	keyDown        [inKeyCount]bool
	keyJust        [inKeyCount]bool
	mods           key.Modifiers
	typed          []rune
	focused        bool       // gio 是否已把键盘焦点给了 gioTag
	snippetReq     *key.Range // 本帧输入法索要的上下文范围（SnippetEvent）
}

var gioIn = &gioInput{}

func (g *gioInput) cursor() (int, int)               { return g.curX, g.curY }
func (g *gioInput) wheel() (float32, float32)        { return g.wheelX, g.wheelY }
func (g *gioInput) mousePressed(b mouseBtn) bool     { return g.btnDown[b] }
func (g *gioInput) mouseJustPressed(b mouseBtn) bool { return g.btnJust[b] }
func (g *gioInput) typedChars() []rune               { return g.typed }

func (g *gioInput) keyPressed(k inKey) bool {
	switch k {
	case keyShift:
		return g.mods&key.ModShift != 0
	case keyCtrl:
		return g.mods&key.ModCtrl != 0
	case keyMeta:
		return g.mods&(key.ModCommand|key.ModSuper) != 0
	}
	return g.keyDown[k]
}

func (g *gioInput) keyJustPressed(k inKey) bool { return g.keyJust[k] }

// resetFrame 每帧开头清零边沿状态（保留光标/按住/修饰键等持久状态）。
func (g *gioInput) resetFrame() {
	g.wheelX, g.wheelY = 0, 0
	g.btnJust = [2]bool{}
	g.keyJust = [inKeyCount]bool{}
	g.typed = g.typed[:0]
	g.snippetReq = nil
}

// gioKeyName 把 gio 键名映射到中立键（仅特殊/导航/快捷键；普通字符走 EditEvent）。
var gioKeyName = map[key.Name]inKey{
	key.NameTab: keyTab, key.NameEscape: keyEscape,
	key.NameReturn: keyEnter, key.NameEnter: keyEnter,
	key.NameLeftArrow: keyLeft, key.NameRightArrow: keyRight,
	key.NameUpArrow: keyUp, key.NameDownArrow: keyDown,
	key.NameDeleteBackward: keyBackspace, key.NameDeleteForward: keyDelete,
	key.NameHome: keyHome, key.NameEnd: keyEnd, key.NameSpace: keySpace,
	"A": keyA, "C": keyC, "X": keyX, "V": keyV, "F12": keyF12,
}

// gioInputFilters 是每帧向 Source 注册的事件过滤器。
//
// key.FocusFilter 必不可少：gio 只用它把 handler 标成 focusable，而
//   - key.EditEvent（字符/IME 提交文本）只投递给 focusable 的 handler；
//   - keyQueue.Frame 每帧会把焦点从「不 focusable」的 handler 上撤掉。
//
// 少了它，就会既打不了字，又陷入「每帧请求焦点→每帧被撤销」的抖动。
// 注意 key.Filter 不会设置 focusable，两者都要注册。
func gioInputFilters() []event.Filter {
	sc := pointer.ScrollRange{Min: -100000, Max: 100000}
	navOpt := key.ModShift | key.ModShortcut | key.ModShortcutAlt
	fs := []event.Filter{
		pointer.Filter{Target: gioTag, Kinds: pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Scroll, ScrollX: sc, ScrollY: sc},
		key.FocusFilter{Target: gioTag},
	}
	for _, n := range []key.Name{
		key.NameTab, key.NameEscape, key.NameReturn, key.NameEnter,
		key.NameLeftArrow, key.NameRightArrow, key.NameUpArrow, key.NameDownArrow,
		key.NameDeleteBackward, key.NameDeleteForward, key.NameHome, key.NameEnd, key.NameSpace, "F12",
	} {
		fs = append(fs, key.Filter{Focus: gioTag, Name: n, Optional: navOpt})
	}
	for _, n := range []key.Name{"A", "C", "X", "V"} { // 选中/复制/剪切/粘贴（需快捷修饰键）
		fs = append(fs, key.Filter{Focus: gioTag, Name: n, Required: key.ModShortcut, Optional: key.ModShift})
	}
	return fs
}

// gioProcessInput 排空本帧队列并更新 gioIn。返回是否收到过任何事件（用于判断是否需重绘）。
func (g *gioInput) process(src interface {
	Event(...event.Filter) (event.Event, bool)
}) {
	filters := gioInputFilters()
	for {
		ev, ok := src.Event(filters...)
		if !ok {
			break
		}
		switch ev := ev.(type) {
		case pointer.Event:
			g.curX, g.curY = int(ev.Position.X), int(ev.Position.Y)
			g.mods = ev.Modifiers
			switch ev.Kind {
			case pointer.Scroll:
				g.wheelX += -ev.Scroll.X / 24 // 引擎再乘 24*uiScale，得到 ≈原始像素
				g.wheelY += -ev.Scroll.Y / 24
			case pointer.Press:
				if ev.Buttons&pointer.ButtonPrimary != 0 && !g.btnDown[btnLeft] {
					g.btnJust[btnLeft] = true
				}
				if ev.Buttons&pointer.ButtonSecondary != 0 && !g.btnDown[btnRight] {
					g.btnJust[btnRight] = true
				}
				g.setButtons(ev.Buttons)
			case pointer.Release, pointer.Drag:
				g.setButtons(ev.Buttons)
			case pointer.Move:
				g.setButtons(ev.Buttons)
			}
		case key.Event:
			g.mods = ev.Modifiers
			if k, ok := gioKeyName[ev.Name]; ok {
				if ev.State == key.Press {
					if !g.keyDown[k] {
						g.keyJust[k] = true
					}
					g.keyDown[k] = true
				} else {
					g.keyDown[k] = false
				}
			}
		case key.EditEvent: // 普通字符输入（含 IME 提交的文本）
			g.typed = append(g.typed, []rune(ev.Text)...)
		case key.FocusEvent:
			g.focused = ev.Focus
		case key.SnippetEvent: // 输入法索要上下文文本，本帧稍后回它
			r := key.Range(ev)
			g.snippetReq = &r
		}
	}
}

func (g *gioInput) setButtons(b pointer.Buttons) {
	g.btnDown[btnLeft] = b&pointer.ButtonPrimary != 0
	g.btnDown[btnRight] = b&pointer.ButtonSecondary != 0
}
