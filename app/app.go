package app

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sjm1327605995/tenon/animation"
	"github.com/sjm1327605995/tenon/canvas"
	"github.com/sjm1327605995/tenon/event"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/state"
	"github.com/sjm1327605995/tenon/widget"
)

// App implements ebiten.Game and hosts a widget tree.
type App struct {
	root       widget.Widget
	ctx        *widget.ContextImpl
	scheduler  *state.Scheduler
	animCtrl   *animation.Controller
	screenSize geometry.Size
	lastFrame  time.Time

	// Input state
	mousePos       geometry.Point
	lastMousePos   geometry.Point
	mouseButtons   event.ButtonState
	lastMouseBtn   event.ButtonState
	hoveredWidget  widget.Widget
	focusedWidget  widget.Widget
}

// New creates a new App with the given root widget.
func New(root widget.Widget) *App {
	ctx := widget.NewContext()

	sched := state.NewScheduler(func(dirty []widget.Widget) {
		for _, w := range dirty {
			if w != nil {
				if s, ok := w.(interface{ SetNeedsRedraw(bool) }); ok {
					s.SetNeedsRedraw(true)
				}
			}
		}
	})

	a := &App{
		root:      root,
		ctx:       ctx,
		scheduler: sched,
		animCtrl:  animation.NewController(),
		lastFrame: time.Now(),
	}

	ctx.SetScheduler(sched)
	ctx.SetOnScheduleAnimation(func() {
		// Animation requests a frame - ebiten already runs continuously
	})

	if root != nil {
		a.mountTree(root)
	}

	return a
}

// SetRoot sets a new root widget.
func (a *App) SetRoot(root widget.Widget) {
	if a.root != nil {
		a.unmountTree(a.root)
	}
	a.root = root
	if root != nil {
		a.mountTree(root)
	}
}

// Layout implements ebiten.Game.
func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	a.screenSize = geometry.Sz(float32(outsideWidth), float32(outsideHeight))
	a.ctx.SetWindowSize(a.screenSize)
	return outsideWidth, outsideHeight
}

// Update implements ebiten.Game.
func (a *App) Update() error {
	if a.root == nil {
		return nil
	}

	now := time.Now()
	dt := now.Sub(a.lastFrame)
	if dt > 100*time.Millisecond {
		dt = 100 * time.Millisecond
	}
	a.lastFrame = now
	a.ctx.BeginFrame(now)

	// Tick animations
	hasActiveAnim := a.animCtrl.Tick(dt)
	if hasActiveAnim {
		// Ensure continuous rendering while animations are active
	}

	// Flush signal changes
	a.scheduler.Flush()

	// Layout pass
	constraints := geometry.Constraints{
		MinWidth:  0,
		MinHeight: 0,
		MaxWidth:  a.screenSize.Width,
		MaxHeight: a.screenSize.Height,
	}
	a.root.Layout(a.ctx, constraints)
	if sb, ok := a.root.(interface{ SetBounds(geometry.Rect) }); ok {
		sb.SetBounds(geometry.FromPointSize(geometry.Pt(0, 0), a.screenSize))
	}

	// Input handling
	a.handleMouseInput()
	a.handleWheelInput()
	a.handleTouchInput()
	a.handleKeyboardInput()

	return nil
}

// Draw implements ebiten.Game.
func (a *App) Draw(screen *ebiten.Image) {
	if a.root == nil {
		return
	}

	c := canvas.New(screen)
	c.Clear(widget.ColorWhite)
	widget.DrawTree(a.root, a.ctx, c)
}

// Run starts the ebiten game loop.
func (a *App) Run() error {
	ebiten.SetWindowTitle("Tenon")
	ebiten.SetWindowSize(800, 600)
	return ebiten.RunGame(a)
}

// mountTree recursively mounts a widget and its children.
func (a *App) mountTree(w widget.Widget) {
	if w == nil {
		return
	}
	if m, ok := w.(widget.Lifecycle); ok {
		m.Mount(a.ctx)
	}
	for _, child := range w.Children() {
		a.mountTree(child)
	}
}

// unmountTree recursively unmounts a widget and its children.
func (a *App) unmountTree(w widget.Widget) {
	if w == nil {
		return
	}
	for _, child := range w.Children() {
		a.unmountTree(child)
	}
	if m, ok := w.(widget.Lifecycle); ok {
		m.Unmount()
	}
	if c, ok := w.(interface{ CleanupBindings() }); ok {
		c.CleanupBindings()
	}
}

// handleMouseInput processes mouse events.
func (a *App) handleMouseInput() {
	x, y := ebiten.CursorPosition()
	a.mousePos = geometry.Pt(float32(x), float32(y))

	// Update button state
	btnState := event.ButtonState(0)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		btnState |= event.ButtonStateLeft
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		btnState |= event.ButtonStateRight
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		btnState |= event.ButtonStateMiddle
	}

	mods := a.modifiers()
	globalPos := a.mousePos

	// Detect button presses/releases
	for _, mb := range []struct {
		ebitenBtn ebiten.MouseButton
		eventBtn  event.Button
		stateBit  event.ButtonState
	}{
		{ebiten.MouseButtonLeft, event.ButtonLeft, event.ButtonStateLeft},
		{ebiten.MouseButtonRight, event.ButtonRight, event.ButtonStateRight},
		{ebiten.MouseButtonMiddle, event.ButtonMiddle, event.ButtonStateMiddle},
	} {
		if inpututil.IsMouseButtonJustPressed(mb.ebitenBtn) {
			target := a.hitTest(a.root, globalPos)
			if target != nil {
				ev := event.NewMouseEvent(event.MousePress, mb.eventBtn, btnState,
					a.toLocal(target, globalPos), globalPos, mods)
				if target.Event(a.ctx, ev) {
					a.focusedWidget = target
				}
			}
		}
		if inpututil.IsMouseButtonJustReleased(mb.ebitenBtn) {
			target := a.hitTest(a.root, globalPos)
			if target != nil {
				ev := event.NewMouseEvent(event.MouseRelease, mb.eventBtn, btnState,
					a.toLocal(target, globalPos), globalPos, mods)
				target.Event(a.ctx, ev)
			}
		}
	}

	// Detect mouse move / enter / leave
	if a.mousePos != a.lastMousePos {
		target := a.hitTest(a.root, globalPos)

		// Leave previous hovered widget
		if a.hoveredWidget != nil && a.hoveredWidget != target {
			ev := event.NewMouseEvent(event.MouseLeave, event.ButtonNone, btnState,
				a.toLocal(a.hoveredWidget, globalPos), globalPos, mods)
			a.hoveredWidget.Event(a.ctx, ev)
		}

		// Enter new widget
		if target != nil && target != a.hoveredWidget {
			ev := event.NewMouseEvent(event.MouseEnter, event.ButtonNone, btnState,
				a.toLocal(target, globalPos), globalPos, mods)
			target.Event(a.ctx, ev)
		}

		// Move event
		if target != nil {
			mouseType := event.MouseMove
			if btnState.AnyPressed() {
				mouseType = event.MouseDrag
			}
			ev := event.NewMouseEvent(mouseType, event.ButtonNone, btnState,
				a.toLocal(target, globalPos), globalPos, mods)
			target.Event(a.ctx, ev)
		}

		a.hoveredWidget = target
	}

	a.lastMousePos = a.mousePos
	a.lastMouseBtn = btnState
}

// handleWheelInput processes mouse wheel scroll events.
func (a *App) handleWheelInput() {
	dx, dy := ebiten.Wheel()
	if dx == 0 && dy == 0 {
		return
	}

	target := a.hoveredWidget
	if target == nil {
		target = a.hitTest(a.root, a.mousePos)
	}
	if target == nil {
		return
	}

	mods := a.modifiers()
	globalPos := a.mousePos
	ev := event.NewWheelEvent(
		geometry.Pt(float32(dx), float32(dy)),
		a.toLocal(target, globalPos),
		globalPos,
		mods,
	)
	target.Event(a.ctx, ev)
}

// handleTouchInput processes touch screen events.
func (a *App) handleTouchInput() {
	// Just pressed touches
	pressedIDs := inpututil.AppendJustPressedTouchIDs(nil)
	for _, id := range pressedIDs {
		x, y := ebiten.TouchPosition(id)
		globalPos := geometry.Pt(float32(x), float32(y))
		target := a.hitTest(a.root, globalPos)
		if target != nil {
			mods := a.modifiers()
			ev := event.NewTouchEvent(event.TouchPress, int(id),
				a.toLocal(target, globalPos), globalPos, mods)
			if target.Event(a.ctx, ev) {
				a.focusedWidget = target
			}
		}
	}

	// Active touches (move)
	justPressedSet := make(map[ebiten.TouchID]bool)
	for _, id := range pressedIDs {
		justPressedSet[id] = true
	}
	for _, id := range ebiten.TouchIDs() {
		if justPressedSet[id] {
			continue
		}
		x, y := ebiten.TouchPosition(id)
		if x >= 0 && y >= 0 {
			globalPos := geometry.Pt(float32(x), float32(y))
			target := a.hitTest(a.root, globalPos)
			if target != nil {
				mods := a.modifiers()
				ev := event.NewTouchEvent(event.TouchMove, int(id),
					a.toLocal(target, globalPos), globalPos, mods)
				target.Event(a.ctx, ev)
			}
		}
	}

	// Just released touches
	releasedIDs := inpututil.AppendJustReleasedTouchIDs(nil)
	for _, id := range releasedIDs {
		globalPos := a.mousePos // fallback
		target := a.hitTest(a.root, globalPos)
		if target != nil {
			mods := a.modifiers()
			ev := event.NewTouchEvent(event.TouchRelease, int(id),
				a.toLocal(target, globalPos), globalPos, mods)
			target.Event(a.ctx, ev)
		}
	}
}

// handleKeyboardInput processes keyboard events.
func (a *App) handleKeyboardInput() {
	mods := a.modifiers()

	// Just pressed keys
	pressedKeys := inpututil.AppendJustPressedKeys(nil)
	for _, key := range pressedKeys {
		if a.focusedWidget != nil {
			k := a.mapKey(key)
			ev := event.NewKeyEvent(event.KeyPress, k, 0, mods)
			a.focusedWidget.Event(a.ctx, ev)
		}
	}

	// Just released keys
	releasedKeys := inpututil.AppendJustReleasedKeys(nil)
	for _, key := range releasedKeys {
		if a.focusedWidget != nil {
			k := a.mapKey(key)
			ev := event.NewKeyEvent(event.KeyRelease, k, 0, mods)
			a.focusedWidget.Event(a.ctx, ev)
		}
	}
}

// hitTest finds the deepest widget containing the given point.
func (a *App) hitTest(w widget.Widget, pt geometry.Point) widget.Widget {
	if w == nil {
		return nil
	}

	// Check visibility
	if vis, ok := w.(interface{ IsVisible() bool }); ok && !vis.IsVisible() {
		return nil
	}

	// Check bounds
	var bounds geometry.Rect
	if bb, ok := w.(interface{ Bounds() geometry.Rect }); ok {
		bounds = bb.Bounds()
	} else {
		return nil
	}

	if !bounds.Contains(pt) {
		return nil
	}

	// Check children (topmost first)
	children := w.Children()
	for i := len(children) - 1; i >= 0; i-- {
		child := children[i]
		if hit := a.hitTest(child, pt); hit != nil {
			return hit
		}
	}

	return w
}

// toLocal converts a global point to widget-local coordinates.
func (a *App) toLocal(w widget.Widget, globalPt geometry.Point) geometry.Point {
	var bounds geometry.Rect
	if bb, ok := w.(interface{ Bounds() geometry.Rect }); ok {
		bounds = bb.Bounds()
	}
	return geometry.Pt(globalPt.X-bounds.Min.X, globalPt.Y-bounds.Min.Y)
}

// modifiers returns current modifier key state.
func (a *App) modifiers() event.Modifiers {
	var mods event.Modifiers
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		mods |= event.ModShift
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		mods |= event.ModCtrl
	}
	if ebiten.IsKeyPressed(ebiten.KeyAlt) {
		mods |= event.ModAlt
	}
	if ebiten.IsKeyPressed(ebiten.KeyMeta) {
		mods |= event.ModSuper
	}
	return mods
}

// mapKey converts an ebiten.Key to an event.Key.
func (a *App) mapKey(key ebiten.Key) event.Key {
	switch key {
	case ebiten.KeyA:
		return event.KeyA
	case ebiten.KeyB:
		return event.KeyB
	case ebiten.KeyC:
		return event.KeyC
	case ebiten.KeyD:
		return event.KeyD
	case ebiten.KeyE:
		return event.KeyE
	case ebiten.KeyF:
		return event.KeyF
	case ebiten.KeyG:
		return event.KeyG
	case ebiten.KeyH:
		return event.KeyH
	case ebiten.KeyI:
		return event.KeyI
	case ebiten.KeyJ:
		return event.KeyJ
	case ebiten.KeyK:
		return event.KeyK
	case ebiten.KeyL:
		return event.KeyL
	case ebiten.KeyM:
		return event.KeyM
	case ebiten.KeyN:
		return event.KeyN
	case ebiten.KeyO:
		return event.KeyO
	case ebiten.KeyP:
		return event.KeyP
	case ebiten.KeyQ:
		return event.KeyQ
	case ebiten.KeyR:
		return event.KeyR
	case ebiten.KeyS:
		return event.KeyS
	case ebiten.KeyT:
		return event.KeyT
	case ebiten.KeyU:
		return event.KeyU
	case ebiten.KeyV:
		return event.KeyV
	case ebiten.KeyW:
		return event.KeyW
	case ebiten.KeyX:
		return event.KeyX
	case ebiten.KeyY:
		return event.KeyY
	case ebiten.KeyZ:
		return event.KeyZ
	case ebiten.Key0, ebiten.KeyNumpad0:
		return event.Key0
	case ebiten.Key1, ebiten.KeyNumpad1:
		return event.Key1
	case ebiten.Key2, ebiten.KeyNumpad2:
		return event.Key2
	case ebiten.Key3, ebiten.KeyNumpad3:
		return event.Key3
	case ebiten.Key4, ebiten.KeyNumpad4:
		return event.Key4
	case ebiten.Key5, ebiten.KeyNumpad5:
		return event.Key5
	case ebiten.Key6, ebiten.KeyNumpad6:
		return event.Key6
	case ebiten.Key7, ebiten.KeyNumpad7:
		return event.Key7
	case ebiten.Key8, ebiten.KeyNumpad8:
		return event.Key8
	case ebiten.Key9, ebiten.KeyNumpad9:
		return event.Key9
	case ebiten.KeyEnter, ebiten.KeyNumpadEnter:
		return event.KeyEnter
	case ebiten.KeySpace:
		return event.KeySpace
	case ebiten.KeyEscape:
		return event.KeyEscape
	case ebiten.KeyBackspace:
		return event.KeyBackspace
	case ebiten.KeyTab:
		return event.KeyTab
	case ebiten.KeyArrowUp:
		return event.KeyUp
	case ebiten.KeyArrowDown:
		return event.KeyDown
	case ebiten.KeyArrowLeft:
		return event.KeyLeft
	case ebiten.KeyArrowRight:
		return event.KeyRight
	case ebiten.KeyShift:
		return event.KeyLeftShift
	case ebiten.KeyControl:
		return event.KeyLeftCtrl
	case ebiten.KeyAlt:
		return event.KeyLeftAlt
	default:
		return event.KeyUnknown
	}
}
