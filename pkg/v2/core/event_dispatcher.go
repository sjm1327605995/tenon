package core

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// EventDispatcher 负责事件分发。
type EventDispatcher struct {
	engine *Engine
}

func newEventDispatcher(e *Engine) *EventDispatcher {
	return &EventDispatcher{engine: e}
}

func (ed *EventDispatcher) handleEvents() {
	mx, my := ebiten.CursorPosition()
	fx, fy := float32(mx), float32(my)

	// 1. MouseMove / Hover
	if fx != ed.engine.lastMouseX || fy != ed.engine.lastMouseY {
		ed.handleMouseMove(fx, fy)
	}

	// 2. Mouse wheel
	dx, dy := ebiten.Wheel()
	if math.Abs(dx) >= 0.5 || math.Abs(dy) >= 0.5 {
		ed.handleScroll(fx, fy, float32(dx), float32(dy))
	}

	// 3. Mouse buttons
	for _, btn := range []ebiten.MouseButton{ebiten.MouseButtonLeft, ebiten.MouseButtonRight, ebiten.MouseButtonMiddle} {
		if inpututil.IsMouseButtonJustPressed(btn) {
			ed.handleMouseDown(fx, fy, toCoreMouseButton(btn))
		}
		if inpututil.IsMouseButtonJustReleased(btn) {
			ed.handleMouseUp(fx, fy, toCoreMouseButton(btn))
		}
	}

	// 4. Keyboard
	ed.handleKeyboard()

	ed.engine.lastMouseX = fx
	ed.engine.lastMouseY = fy
}

func (ed *EventDispatcher) handleMouseMove(x, y float32) {
	target := ed.hitTestWithOverlays(x, y)
	if target != ed.engine.hoverTarget {
		ed.dispatchHoverChange(ed.engine.hoverTarget, target, x, y)
		ed.engine.hoverTarget = target
	}
	if target != nil {
		b := target.GetBounds()
		ed.dispatchEvent(target, &Event{
			Type:   EventMouseMove,
			X:      x,
			Y:      y,
			LocalX: x - b.X,
			LocalY: y - b.Y,
			Target: target,
		})
	}
}

func (ed *EventDispatcher) dispatchHoverChange(oldTarget, newTarget Element, x, y float32) {
	if oldTarget == newTarget {
		return
	}
	lca := findLowestCommonAncestor(oldTarget, newTarget)
	for el := oldTarget; el != nil && el != lca; el = el.GetParent() {
		ed.dispatchEvent(el, &Event{Type: EventMouseLeave, X: x, Y: y, Target: el})
	}
	path := buildEventPathTo(newTarget, lca)
	for _, el := range path {
		ed.dispatchEvent(el, &Event{Type: EventMouseEnter, X: x, Y: y, Target: el})
	}
}

func (ed *EventDispatcher) handleMouseDown(x, y float32, btn MouseButton) {
	target := ed.hitTestWithOverlays(x, y)
	if target != nil {
		LogDebug("[Engine] handleMouseDown", "target", target.ElementType(), "bounds", target.GetBounds(), "x", x, "y", y)
		// Debug: log target hierarchy and children
		parentType := "nil"
		if target.GetParent() != nil {
			parentType = target.GetParent().ElementType()
		}
		LogDebug("[Engine] handleMouseDown debug", "target", target.ElementType(), "parent", parentType, "childrenCount", len(target.GetChildren()))
		for i, c := range target.GetChildren() {
			LogDebug("[Engine] handleMouseDown child", "index", i, "type", c.ElementType(), "bounds", c.GetBounds())
		}
		b := target.GetBounds()
		ed.engine.mouseDownTarget = target
		ed.engine.mouseDownX = x
		ed.engine.mouseDownY = y
		ed.engine.mouseDownButton = btn
		// Set focus on click
		if ed.engine.focusTarget != target {
			ed.setFocus(target)
		}
		ed.dispatchEvent(target, &Event{
			Type:   EventMouseDown,
			X:      x,
			Y:      y,
			LocalX: x - b.X,
			LocalY: y - b.Y,
			Button: btn,
			Target: target,
		})
	} else {
		// Click on empty area clears focus
		ed.setFocus(nil)
	}
}

func (ed *EventDispatcher) handleMouseUp(x, y float32, btn MouseButton) {
	if ed.engine.mouseDownTarget != nil {
		b := ed.engine.mouseDownTarget.GetBounds()
		ed.dispatchEvent(ed.engine.mouseDownTarget, &Event{
			Type:   EventMouseUp,
			X:      x,
			Y:      y,
			LocalX: x - b.X,
			LocalY: y - b.Y,
			Button: btn,
			Target: ed.engine.mouseDownTarget,
		})

		// Click = mouseDown + mouseUp on same target or ancestor/descendant
		target := ed.hitTestWithOverlays(x, y)
		clickTarget := ed.resolveClickTarget(ed.engine.mouseDownTarget, target)
		if clickTarget != nil {
			cb := clickTarget.GetBounds()
			ed.dispatchEvent(clickTarget, &Event{
				Type:   EventClick,
				X:      x,
				Y:      y,
				LocalX: x - cb.X,
				LocalY: y - cb.Y,
				Button: btn,
				Target: clickTarget,
			})
		}
	}
	ed.engine.mouseDownTarget = nil
}

func (ed *EventDispatcher) handleScroll(x, y float32, dx, dy float32) {
	target := ed.hitTestWithOverlays(x, y)
	LogDebug("[Engine] handleScroll", "target", target.ElementType())
	if target != nil {
		b := target.GetBounds()
		ed.dispatchEvent(target, &Event{
			Type:   EventScroll,
			X:      x,
			Y:      y,
			LocalX: x - b.X,
			LocalY: y - b.Y,
			DeltaX: dx,
			DeltaY: dy,
			Target: target,
		})
	}
}

func (ed *EventDispatcher) handleKeyboard() {
	// Tab to cycle focus
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		ed.focusNext(ebiten.IsKeyPressed(ebiten.KeyShift))
		return
	}

	if ed.engine.focusTarget == nil {
		return
	}

	// Space/Enter triggers click on focused element
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		b := ed.engine.focusTarget.GetBounds()
		ed.dispatchEvent(ed.engine.focusTarget, &Event{
			Type:   EventClick,
			X:      b.X + b.Width/2,
			Y:      b.Y + b.Height/2,
			LocalX: b.Width / 2,
			LocalY: b.Height / 2,
			Target: ed.engine.focusTarget,
		})
		return
	}

	// KeyDown / KeyUp
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		if inpututil.IsKeyJustPressed(key) {
			ed.dispatchEvent(ed.engine.focusTarget, &Event{Type: EventKeyDown, Key: toCoreKey(key), Target: ed.engine.focusTarget})
		}
		if inpututil.IsKeyJustReleased(key) {
			ed.dispatchEvent(ed.engine.focusTarget, &Event{Type: EventKeyUp, Key: toCoreKey(key), Target: ed.engine.focusTarget})
		}
	}
}

// dispatchEvent sends event to target and bubbles up until consumed.
// 优先使用注册表分发，如果注册表没有消费事件，再调用 HandleEvent。
func (ed *EventDispatcher) dispatchEvent(target Element, event *Event) {
	// Log event to debugger if enabled
	if ed.engine.debugBus != nil && ed.engine.debugBus.IsEnabled() {
		ed.engine.debugBus.AddEventLog(&debugEventWrapper{event: event, targetType: target.ElementType()})
	}

	// 1. 先通过注册表分发（支持冒泡）
	if ed.engine.eventRegistry != nil {
		if ed.engine.eventRegistry.Dispatch(event) {
			return // consumed by registry
		}
	}

	// 2. 回退到 HandleEvent（兼容旧组件）
	el := target
	for el != nil {
		if el.HandleEvent(event) {
			return // consumed
		}
		el = el.GetParent()
	}
}

// broadcastEvent 递归广播事件给元素及其所有后代。
// 优先使用注册表分发，再回退到 HandleEvent。
func (ed *EventDispatcher) broadcastEvent(el Element, event *Event) {
	if el == nil {
		return
	}
	// 先尝试注册表
	if ed.engine.eventRegistry != nil {
		if ed.engine.eventRegistry.DispatchToTarget(el, event) {
			return
		}
	}
	// 回退到 HandleEvent
	el.HandleEvent(event)
	for _, child := range el.GetChildren() {
		ed.broadcastEvent(child, event)
	}
}

// GetFocusTarget returns the currently focused element.
func (ed *EventDispatcher) GetFocusTarget() Element {
	return ed.engine.focusTarget
}

// setFocus changes focused element.
func (ed *EventDispatcher) setFocus(target Element) {
	if ed.engine.focusTarget == target {
		return
	}
	if ed.engine.focusTarget != nil {
		ed.dispatchEvent(ed.engine.focusTarget, &Event{Type: EventFocusOut, Target: ed.engine.focusTarget})
	}
	ed.engine.focusTarget = target
	if target != nil {
		ed.dispatchEvent(target, &Event{Type: EventFocusIn, Target: target})
	}
}

// focusNext cycles focus to next/previous focusable element.
func (ed *EventDispatcher) focusNext(reverse bool) {
	focusables := ed.collectFocusables(ed.engine.rootElement)
	if len(focusables) == 0 {
		return
	}
	idx := -1
	for i, el := range focusables {
		if el == ed.engine.focusTarget {
			idx = i
			break
		}
	}
	if reverse {
		idx--
		if idx < 0 {
			idx = len(focusables) - 1
		}
	} else {
		idx++
		if idx >= len(focusables) {
			idx = 0
		}
	}
	ed.setFocus(focusables[idx])
}

func (ed *EventDispatcher) collectFocusables(el Element) []Element {
	if el == nil {
		return nil
	}
	var result []Element
	if el.HasFlag(FlagFocusable) {
		result = append(result, el)
	}
	for _, child := range el.GetChildren() {
		result = append(result, ed.collectFocusables(child)...)
	}
	return result
}

// GetHoverTarget returns the currently hovered element.
func (ed *EventDispatcher) GetHoverTarget() Element {
	return ed.engine.hoverTarget
}

func (ed *EventDispatcher) resolveClickTarget(downTarget, upTarget Element) Element {
	if downTarget == nil || upTarget == nil {
		return nil
	}
	if downTarget == upTarget {
		return downTarget
	}
	if isAncestorOf(downTarget, upTarget) {
		return downTarget
	}
	if isAncestorOf(upTarget, downTarget) {
		return upTarget
	}
	return nil
}

func (ed *EventDispatcher) hitTest(el Element, x, y float32) Element {
	return ed.hitTestClipped(el, x, y, nil)
}

func (ed *EventDispatcher) hitTestWithOverlays(x, y float32) Element {
	for i := len(ed.engine.overlays) - 1; i >= 0; i-- {
		if ed.engine.overlays[i] != nil && ed.engine.overlays[i].IsVisible() {
			if h := ed.hitTestClipped(ed.engine.overlays[i], x, y, nil); h != nil {
				return h
			}
		}
	}
	return ed.hitTestClipped(ed.engine.rootElement, x, y, nil)
}

func (ed *EventDispatcher) hitTestClipped(el Element, x, y float32, clipBounds *LayoutBounds) Element {
	if el == nil || !el.IsVisible() {
		return nil
	}
	b := el.GetBounds()
	if b.Width <= 0 || b.Height <= 0 {
		return nil
	}
	if clipBounds != nil {
		if b.X+b.Width <= clipBounds.X || b.X >= clipBounds.X+clipBounds.Width ||
			b.Y+b.Height <= clipBounds.Y || b.Y >= clipBounds.Y+clipBounds.Height {
			return nil
		}
	}
	children := el.GetChildren()
	var childClip *LayoutBounds
	if el.HasFlag(FlagClipChildren) {
		childClip = &b
	} else {
		childClip = clipBounds
	}
	for i := len(children) - 1; i >= 0; i-- {
		if ed.engine.IsOverlay(children[i]) {
			continue
		}
		if h := ed.hitTestClipped(children[i], x, y, childClip); h != nil {
			return h
		}
	}
	if el.GetPointerEvents() == PointerEventsNone {
		return nil
	}
	localX, localY := x, y
	t := el.GetTransform()
	if !t.IsIdentity() {
		localX, localY = inverseTransformPoint(b, t, x, y)
	}
	if b.Contains(localX, localY) {
		if clipBounds != nil {
			if clipBounds.Contains(x, y) {
				return el
			}
			return nil
		}
		return el
	}
	return nil
}

func (ed *EventDispatcher) applyInverseTransform(el Element, x, y float32) (float32, float32) {
	b := el.GetBounds()
	t := el.GetTransform()
	if t.IsIdentity() {
		return x, y
	}
	g := BuildTransformGeoM(b, t)
	if !g.IsInvertible() {
		return x, y
	}
	g.Invert()
	ix, iy := g.Apply(float64(x), float64(y))
	return float32(ix), float32(iy)
}

func findLowestCommonAncestor(a, b Element) Element {
	if a == nil || b == nil {
		return nil
	}
	ancestors := make(map[Element]bool)
	for el := a; el != nil; el = el.GetParent() {
		ancestors[el] = true
	}
	for el := b; el != nil; el = el.GetParent() {
		if ancestors[el] {
			return el
		}
	}
	return nil
}

func buildEventPathTo(target, stopAt Element) []Element {
	var path []Element
	for el := target; el != nil && el != stopAt; el = el.GetParent() {
		path = append(path, el)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func isAncestorOf(ancestor, descendant Element) bool {
	for p := descendant.GetParent(); p != nil; p = p.GetParent() {
		if p == ancestor {
			return true
		}
	}
	return false
}

func inverseTransformPoint(bounds LayoutBounds, t Transform, x, y float32) (float32, float32) {
	g := BuildTransformGeoM(bounds, t)
	if !g.IsInvertible() {
		return x, y
	}
	g.Invert()
	ix, iy := g.Apply(float64(x), float64(y))
	return float32(ix), float32(iy)
}

// debugEventWrapper adapts *Event to the debugger interface.
type debugEventWrapper struct {
	event      *Event
	targetType string
}

func (w *debugEventWrapper) GetType() string {
	switch w.event.Type {
	case EventMouseMove:
		return "MouseMove"
	case EventMouseDown:
		return "MouseDown"
	case EventMouseUp:
		return "MouseUp"
	case EventClick:
		return "Click"
	case EventScroll:
		return "Scroll"
	case EventKeyDown:
		return "KeyDown"
	case EventKeyUp:
		return "KeyUp"
	case EventFocusIn:
		return "FocusIn"
	case EventFocusOut:
		return "FocusOut"
	case EventResize:
		return "Resize"
	default:
		return "Unknown"
	}
}

func (w *debugEventWrapper) GetTarget() string  { return w.targetType }
func (w *debugEventWrapper) GetX() float32      { return w.event.X }
func (w *debugEventWrapper) GetY() float32      { return w.event.Y }
func (w *debugEventWrapper) GetDeltaX() float32 { return w.event.DeltaX }
func (w *debugEventWrapper) GetDeltaY() float32 { return w.event.DeltaY }

// toCoreMouseButton 将 ebiten 鼠标按键映射为 core 抽象类型。
func toCoreMouseButton(btn ebiten.MouseButton) MouseButton {
	switch btn {
	case ebiten.MouseButtonLeft:
		return MouseButtonLeft
	case ebiten.MouseButtonRight:
		return MouseButtonRight
	case ebiten.MouseButtonMiddle:
		return MouseButtonMiddle
	default:
		return MouseButton(btn)
	}
}

// toCoreKey 将 ebiten 键盘按键映射为 core 抽象类型。
func toCoreKey(key ebiten.Key) Key {
	switch key {
	case ebiten.KeyA:
		return KeyA
	case ebiten.KeyB:
		return KeyB
	case ebiten.KeyC:
		return KeyC
	case ebiten.KeyD:
		return KeyD
	case ebiten.KeyE:
		return KeyE
	case ebiten.KeyF:
		return KeyF
	case ebiten.KeyG:
		return KeyG
	case ebiten.KeyH:
		return KeyH
	case ebiten.KeyI:
		return KeyI
	case ebiten.KeyJ:
		return KeyJ
	case ebiten.KeyK:
		return KeyK
	case ebiten.KeyL:
		return KeyL
	case ebiten.KeyM:
		return KeyM
	case ebiten.KeyN:
		return KeyN
	case ebiten.KeyO:
		return KeyO
	case ebiten.KeyP:
		return KeyP
	case ebiten.KeyQ:
		return KeyQ
	case ebiten.KeyR:
		return KeyR
	case ebiten.KeyS:
		return KeyS
	case ebiten.KeyT:
		return KeyT
	case ebiten.KeyU:
		return KeyU
	case ebiten.KeyV:
		return KeyV
	case ebiten.KeyW:
		return KeyW
	case ebiten.KeyX:
		return KeyX
	case ebiten.KeyY:
		return KeyY
	case ebiten.KeyZ:
		return KeyZ
	case ebiten.KeyDigit0:
		return Key0
	case ebiten.KeyDigit1:
		return Key1
	case ebiten.KeyDigit2:
		return Key2
	case ebiten.KeyDigit3:
		return Key3
	case ebiten.KeyDigit4:
		return Key4
	case ebiten.KeyDigit5:
		return Key5
	case ebiten.KeyDigit6:
		return Key6
	case ebiten.KeyDigit7:
		return Key7
	case ebiten.KeyDigit8:
		return Key8
	case ebiten.KeyDigit9:
		return Key9
	case ebiten.KeySpace:
		return KeySpace
	case ebiten.KeyEnter:
		return KeyEnter
	case ebiten.KeyTab:
		return KeyTab
	case ebiten.KeyEscape:
		return KeyEscape
	case ebiten.KeyBackspace:
		return KeyBackspace
	case ebiten.KeyDelete:
		return KeyDelete
	case ebiten.KeyArrowLeft:
		return KeyLeft
	case ebiten.KeyArrowRight:
		return KeyRight
	case ebiten.KeyArrowUp:
		return KeyUp
	case ebiten.KeyArrowDown:
		return KeyDown
	case ebiten.KeyHome:
		return KeyHome
	case ebiten.KeyEnd:
		return KeyEnd
	case ebiten.KeyShift:
		return KeyShift
	case ebiten.KeyControl:
		return KeyCtrl
	case ebiten.KeyAlt:
		return KeyAlt
	case ebiten.KeyF1:
		return KeyF1
	case ebiten.KeyF2:
		return KeyF2
	case ebiten.KeyF3:
		return KeyF3
	case ebiten.KeyF4:
		return KeyF4
	case ebiten.KeyF5:
		return KeyF5
	case ebiten.KeyF6:
		return KeyF6
	case ebiten.KeyF7:
		return KeyF7
	case ebiten.KeyF8:
		return KeyF8
	case ebiten.KeyF9:
		return KeyF9
	case ebiten.KeyF10:
		return KeyF10
	case ebiten.KeyF11:
		return KeyF11
	case ebiten.KeyF12:
		return KeyF12
	default:
		return Key(key)
	}
}
