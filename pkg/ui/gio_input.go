package ui

import (
	"io"

	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/transfer"
)

// ---- gio 后端：输入 ----
//
// gio 是事件式输入；这里把每帧收到的 pointer/key/edit 事件累积成与 ebiten 一致的可轮询状态，
// 从而复用引擎的 handleInput。持久状态（光标/按住/修饰键）跨帧保留，边沿状态（刚按下/滚轮/
// 输入字符）每帧开头清零、由本帧事件填充。

// gioTag 是整窗输入区域与键盘焦点的标识。
var gioTag = new(int)

type gioInput struct {
	curX, curY     float32 // 光标位置，原样保留 gio 给的亚像素精度
	wheelX, wheelY float32
	btnDown        [2]bool // btnLeft/btnRight 当前按住
	btnJust        [2]bool // 本帧刚按下
	keyDown        [inKeyCount]bool
	keyJust        [inKeyCount]bool
	mods           key.Modifiers
	typed          []rune
	focused        bool       // gio 是否已把键盘焦点给了 gioTag
	snippetReq     *key.Range // 本帧输入法索要的上下文范围（SnippetEvent）

	// 输入法编辑：EditEvent 是「替换某 rune 区间」，不是「在光标处追加」。
	edits      []gioEdit  // 本帧的替换（组字更新与提交都走这里）
	selEvt     *key.Range // 输入法移动了选区
	composeEvt *key.Range // 组字中的区间（用于预编辑下划线）；空区间表示组字结束
}

// gioEdit 是一次「把 rng 这段 rune 替换成 text」的编辑。
type gioEdit struct {
	rng  key.Range
	text string
}

var gioIn = &gioInput{}

func (g *gioInput) cursor() (float32, float32) { return g.curX, g.curY }

// setCursor 记录光标位置（物理像素，保留亚像素精度）。
func (g *gioInput) setCursor(x, y float32)           { g.curX, g.curY = x, y }
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
	g.edits = g.edits[:0]
	g.selEvt, g.composeEvt = nil, nil
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
	// 系统剪贴板读回的数据（Ctrl+V 发出 ReadCmd 后经此送达）。
	fs = append(fs, transfer.TargetFilter{Target: gioTag, Type: gioClipMime})
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
			debugPointer(ev)
			g.setCursor(ev.Position.X, ev.Position.Y) // 原样透传，不取整
			g.mods = ev.Modifiers
			switch ev.Kind {
			case pointer.Scroll:
				// gio 的 Scroll 已经是物理像素，原样透传（inputSource 的契约即为物理像素）。
				g.wheelX += ev.Scroll.X
				g.wheelY += ev.Scroll.Y
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
			k, ok := gioKeyName[ev.Name]
			switch {
			case !ok:
			case k == keyV && ev.State == key.Press && ev.Modifiers.Contain(key.ModShortcut):
				// 粘贴由后端接管：系统剪贴板是异步读的，让引擎同步读缓存会粘到过期内容。
				gioClip.pasteRequested = true
			case ev.State == key.Press:
				if !g.keyDown[k] {
					g.keyJust[k] = true
				}
				g.keyDown[k] = true
			default:
				g.keyDown[k] = false
			}
		case key.EditEvent:
			// 必须按 Range 替换：组字过程中输入法反复用新串替换上一次的区间
			// （"s" -> 用 "shi" 替换 [0,1)）。当成追加会写出 "sshi"。
			g.edits = append(g.edits, gioEdit{rng: ev.Range, text: ev.Text})
		case key.FocusEvent:
			g.focused = ev.Focus
		case key.SnippetEvent: // 输入法索要上下文文本，本帧稍后回它
			r := key.Range(ev)
			g.snippetReq = &r
		case key.SelectionEvent:
			r := key.Range(ev)
			g.selEvt = &r
		case key.CompositionEvent: // 组字区间（预编辑），用于下划线显示
			r := key.Range(ev)
			g.composeEvt = &r
		case transfer.DataEvent: // 系统剪贴板读回的文本
			if ev.Type == gioClipMime {
				if rc := ev.Open(); rc != nil {
					b, err := io.ReadAll(rc)
					rc.Close()
					if err == nil {
						s := string(b)
						gioClip.pastedText = &s
					}
				}
			}
		}
	}
}

// gioCursor 按光标下的元素挑选系统指针形状：文本框用 I 型、可点击元素用手型。
// 从命中节点向上冒泡，最内层的语义优先（输入框套在可点击容器里时应显示 I 型）。
func gioCursor(g *game) pointer.Cursor {
	if g.rootRN == nil {
		return pointer.CursorDefault
	}
	x, y := gioIn.cursor()
	for c := g.hitTop(x, y); c != nil; c = c.parent {
		switch {
		case c.kind == rnInput:
			return pointer.CursorText
		case c.onClick != nil:
			return pointer.CursorPointer
		}
	}
	return pointer.CursorDefault
}

func (g *gioInput) setButtons(b pointer.Buttons) {
	g.btnDown[btnLeft] = b&pointer.ButtonPrimary != 0
	g.btnDown[btnRight] = b&pointer.ButtonSecondary != 0
}
