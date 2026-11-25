package ui

import (
	"fmt"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
)

type gioWindowRuntime struct {
	runtime *BackendRuntime
	win     *app.Window
	input   *gioInput
	ime     gioIMEState
	clip    *gioClipboardState
	group   *gioWindowGroup
	ops     op.Ops
	last    time.Time
}

type gioWindowGroup struct {
	windows []*gioWindowRuntime
}

func (g *gioWindowGroup) close() {
	for _, window := range g.windows {
		window.win.Perform(system.ActionClose)
	}
}

var gioWindowFrameMu sync.Mutex
var gioWindowWakes = struct {
	sync.Mutex
	items map[*game]func()
}{items: make(map[*game]func())}

func registerGioWindowWake(g *game, wake func()) {
	gioWindowWakes.Lock()
	gioWindowWakes.items[g] = wake
	gioWindowWakes.Unlock()
	setBackendWake(wakeGioWindows)
}

func unregisterGioWindowWake(g *game) {
	gioWindowWakes.Lock()
	delete(gioWindowWakes.items, g)
	gioWindowWakes.Unlock()
}

func wakeGioWindows() {
	gioWindowWakes.Lock()
	wakes := make([]func(), 0, len(gioWindowWakes.items))
	for _, wake := range gioWindowWakes.items {
		wakes = append(wakes, wake)
	}
	gioWindowWakes.Unlock()
	for _, wake := range wakes {
		wake()
	}
}

func newGioWindowRuntime(root *Node, cfg windowConfig) *gioWindowRuntime {
	win := new(app.Window)
	win.Option(gioWindowOptions(cfg)...)
	tag := new(int)
	clipboard := &gioClipboardState{tag: tag}
	input := &gioInput{tag: tag, clip: clipboard}
	runtime := NewBackendRuntime(root, cfg, input)
	window := &gioWindowRuntime{
		runtime: runtime,
		win:     win,
		input:   input,
		ime:     gioIMEState{tag: tag},
		clip:    clipboard,
	}
	registerGioClipboard(runtime.g, clipboard)
	registerGioWindowWake(runtime.g, win.Invalidate)
	return window
}

func (w *gioWindowRuntime) run() {
	defer func() {
		unregisterGioWindowWake(w.runtime.g)
		unregisterGioClipboard(w.runtime.g)
		w.runtime.Close()
	}()
	for {
		switch e := w.win.Event().(type) {
		case app.DestroyEvent:
			return
		case app.FrameEvent:
			w.frame(e)
		}
	}
}

func (w *gioWindowRuntime) frame(e app.FrameEvent) {
	gioWindowFrameMu.Lock()
	defer gioWindowFrameMu.Unlock()
	restore := w.runtime.activate()
	defer restore()

	scale := e.Metric.PxPerDp
	if scale < 1 {
		scale = 1
	}
	w.runtime.Resize(e.Size.X, e.Size.Y, scale)
	g := w.runtime.g
	w.input.resetFrame()
	w.input.accessibilityTags = w.runtime.g.inputTags
	if w.input.process(e.Source) {
		g.visualVersion++
	}
	w.ime.applyEdits(g, w.input)
	w.clip.flush(e.Source, g)

	var dt float32
	now := e.Now
	if !w.last.IsZero() {
		dt = float32(now.Sub(w.last).Seconds() * 1000)
	}
	w.last = now
	layoutMoving := w.runtime.Frame(dt)

	w.ops.Reset()
	g.frameCache.Add(g, &w.ops, g.framePaintVersion(now))
	area := clip.Rect{Max: e.Size}.Push(&w.ops)
	event.Op(&w.ops, w.input.tag)
	key.InputHintOp{Tag: w.input.tag, Hint: key.HintAny}.Add(&w.ops)
	gioCursor(g, w.input).Add(&w.ops)
	area.Pop()
	if !w.input.focused {
		e.Source.Execute(key.FocusCmd{Tag: w.input.tag})
	}
	if w.input.snippetReq != nil {
		w.ime.handleSnippetReq(e.Source, g, *w.input.snippetReq)
	}
	w.ime.sync(e.Source, g)

	e.Frame(&w.ops)
	if captureGioFrame(g) {
		if w.group != nil {
			w.group.close()
		} else {
			w.win.Perform(system.ActionClose)
		}
	}
	scheduleNextFrame(e, g, layoutMoving)
	debugFrame(g, wantsNextFrame(g, layoutMoving))
}

// RunWindows opens one Gio window per root and blocks until all windows close.
// Frame submission is serialized because component rendering is single-threaded.
func RunWindows(roots ...*Node) error {
	if len(roots) == 0 {
		return nil
	}
	if _, ok := activeBackend.(gioBackend); !ok {
		return fmt.Errorf("ui: RunWindows requires the Gio backend")
	}
	windows := make([]*gioWindowRuntime, 0, len(roots))
	for i, root := range roots {
		cfg := winCfg
		cfg.sync = FrameSync
		if len(roots) > 1 {
			cfg.title = fmt.Sprintf("%s — %d", cfg.title, i+1)
		}
		windows = append(windows, newGioWindowRuntime(root, cfg))
	}
	group := &gioWindowGroup{windows: windows}
	for _, window := range windows {
		window.group = group
	}
	var wg sync.WaitGroup
	wg.Add(len(windows))
	for _, window := range windows {
		go func() {
			defer wg.Done()
			window.run()
		}()
	}
	wg.Wait()
	return nil
}
