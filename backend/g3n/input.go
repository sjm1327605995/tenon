package g3n

import (
	"github.com/g3n/engine/window"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

type input struct {
	win       window.IWindow
	x, y      float32
	wheelX    float32
	wheelY    float32
	mouseDown map[ui.MouseButton]bool
	mouseJust map[ui.MouseButton]bool
	keyDown   map[ui.InputKey]bool
	keyJust   map[ui.InputKey]bool
	typed     []rune
}

func newInput(win window.IWindow) *input {
	i := &input{
		win:       win,
		mouseDown: make(map[ui.MouseButton]bool),
		mouseJust: make(map[ui.MouseButton]bool),
		keyDown:   make(map[ui.InputKey]bool),
		keyJust:   make(map[ui.InputKey]bool),
	}
	win.Subscribe(window.OnCursor, func(_ string, event interface{}) {
		e := event.(*window.CursorEvent)
		i.x, i.y = e.Xpos, e.Ypos
	})
	win.Subscribe(window.OnScroll, func(_ string, event interface{}) {
		e := event.(*window.ScrollEvent)
		i.wheelX -= e.Xoffset * 40
		i.wheelY -= e.Yoffset * 40
	})
	win.Subscribe(window.OnMouseDown, func(_ string, event interface{}) {
		e := event.(*window.MouseEvent)
		button := mouseButton(e.Button)
		if !i.mouseDown[button] {
			i.mouseJust[button] = true
		}
		i.mouseDown[button] = true
	})
	win.Subscribe(window.OnMouseUp, func(_ string, event interface{}) {
		e := event.(*window.MouseEvent)
		i.mouseDown[mouseButton(e.Button)] = false
	})
	win.Subscribe(window.OnKeyDown, func(_ string, event interface{}) {
		e := event.(*window.KeyEvent)
		for _, key := range inputKeys(e.Key) {
			if !i.keyDown[key] {
				i.keyJust[key] = true
			}
			i.keyDown[key] = true
		}
	})
	win.Subscribe(window.OnKeyUp, func(_ string, event interface{}) {
		e := event.(*window.KeyEvent)
		for _, key := range inputKeys(e.Key) {
			i.keyDown[key] = false
		}
	})
	win.Subscribe(window.OnChar, func(_ string, event interface{}) {
		e := event.(*window.CharEvent)
		i.typed = append(i.typed, e.Char)
	})
	if native, ok := win.(*window.GlfwWindow); ok {
		ui.SetClipboardProvider(native.GetClipboardString, native.SetClipboardString)
	}
	return i
}

func mouseButton(button window.MouseButton) ui.MouseButton {
	if button == window.MouseButtonRight {
		return ui.MouseRight
	}
	return ui.MouseLeft
}

func inputKeys(key window.Key) []ui.InputKey {
	if mapped, ok := keyMap[key]; ok {
		return []ui.InputKey{mapped}
	}
	return nil
}

var keyMap = map[window.Key]ui.InputKey{
	window.KeyF12:          ui.KeyF12,
	window.KeyTab:          ui.KeyTab,
	window.KeyLeftShift:    ui.KeyShift,
	window.KeyRightShift:   ui.KeyShift,
	window.KeyLeftControl:  ui.KeyCtrl,
	window.KeyRightControl: ui.KeyCtrl,
	window.KeyLeftSuper:    ui.KeyMeta,
	window.KeyRightSuper:   ui.KeyMeta,
	window.KeyDown:         ui.KeyDown,
	window.KeyUp:           ui.KeyUp,
	window.KeyLeft:         ui.KeyLeft,
	window.KeyRight:        ui.KeyRight,
	window.KeyEscape:       ui.KeyEscape,
	window.KeyEnter:        ui.KeyEnter,
	window.KeyKPEnter:      ui.KeyEnter,
	window.KeySpace:        ui.KeySpace,
	window.KeyA:            ui.KeyA,
	window.KeyC:            ui.KeyC,
	window.KeyX:            ui.KeyX,
	window.KeyV:            ui.KeyV,
	window.KeyBackspace:    ui.KeyBackspace,
	window.KeyDelete:       ui.KeyDelete,
	window.KeyHome:         ui.KeyHome,
	window.KeyEnd:          ui.KeyEnd,
}

func (i *input) endFrame() {
	clear(i.mouseJust)
	clear(i.keyJust)
	i.wheelX, i.wheelY = 0, 0
	i.typed = i.typed[:0]
}

func (i *input) Cursor() (float32, float32) {
	sx, sy := i.win.GetScale()
	return i.x * float32(sx), i.y * float32(sy)
}

func (i *input) Wheel() (float32, float32) {
	sx, sy := i.win.GetScale()
	return i.wheelX * float32(sx), i.wheelY * float32(sy)
}

func (i *input) MousePressed(button ui.MouseButton) bool     { return i.mouseDown[button] }
func (i *input) MouseJustPressed(button ui.MouseButton) bool { return i.mouseJust[button] }
func (i *input) KeyPressed(key ui.InputKey) bool             { return i.keyDown[key] }
func (i *input) KeyJustPressed(key ui.InputKey) bool         { return i.keyJust[key] }
func (i *input) TypedChars() []rune                          { return i.typed }
