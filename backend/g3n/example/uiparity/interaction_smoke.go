package main

import (
	"fmt"
	"time"

	"github.com/g3n/engine/window"
	"github.com/sjm1327605995/tenon/backend/g3n"
	"github.com/sjm1327605995/tenon/example/uiparity/shared"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func scheduleNativeInteraction(game *g3n.Game, state *shared.State) {
	post := func(delay time.Duration, fn func()) {
		time.AfterFunc(delay, func() { ui.Post(fn) })
	}
	click := func(x, y float32) func() {
		return func() {
			game.Window.Dispatch(window.OnCursor, &window.CursorEvent{Xpos: x, Ypos: y})
			game.Window.Dispatch(window.OnMouseDown, &window.MouseEvent{Xpos: x, Ypos: y, Button: window.MouseButtonLeft})
			game.Window.Dispatch(window.OnMouseUp, &window.MouseEvent{Xpos: x, Ypos: y, Button: window.MouseButtonLeft})
		}
	}
	key := func(key window.Key, mods window.ModifierKey) {
		game.Window.Dispatch(window.OnKeyDown, &window.KeyEvent{Key: key, Mods: mods})
		game.Window.Dispatch(window.OnKeyUp, &window.KeyEvent{Key: key, Mods: mods})
	}
	post(200*time.Millisecond, click(240, 228))
	post(400*time.Millisecond, func() { key(window.KeyTab, 0) })
	post(500*time.Millisecond, func() {
		for _, char := range "你好" {
			game.Window.Dispatch(window.OnChar, &window.CharEvent{Char: char})
		}
	})
	post(600*time.Millisecond, func() {
		game.Window.Dispatch(window.OnKeyDown, &window.KeyEvent{Key: window.KeyLeftControl})
		key(window.KeyA, window.ModControl)
	})
	post(650*time.Millisecond, func() {
		game.Window.Dispatch(window.OnKeyUp, &window.KeyEvent{Key: window.KeyLeftControl})
		key(window.KeyBackspace, 0)
	})
	post(700*time.Millisecond, func() {
		game.Window.Dispatch(window.OnChar, &window.CharEvent{Char: '好'})
	})
	post(800*time.Millisecond, func() {
		if state.Status.Get() != "按钮已点击" || state.Text.Get() != "好" {
			panic(fmt.Sprintf("native interaction failed: status=%q text=%q", state.Status.Get(), state.Text.Get()))
		}
		fmt.Println("native button, Tab, selection and Chinese input smoke passed")
	})
	post(1500*time.Millisecond, g3n.Exit)
}
