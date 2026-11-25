package g3n

// PointerEvent is a scene-space pointer snapshot. X and Y are normalized device
// coordinates in [-1,1] when Inside is true. Click fields are frame edges.
type PointerEvent struct {
	X, Y                  float32
	Inside                bool
	LeftClick, RightClick bool
}

// PointerHandler is implemented by scene components that accept pointer input.
type PointerHandler interface {
	HandlePointer(PointerEvent)
}

func pointerEventAt(cursorX, cursorY float32, x, y, width, height int, leftClick, rightClick bool) PointerEvent {
	event := PointerEvent{Inside: width > 0 && height > 0, LeftClick: leftClick, RightClick: rightClick}
	if !event.Inside {
		return event
	}
	event.X = (cursorX-float32(x))/float32(width)*2 - 1
	event.Y = (float32(y+height)-cursorY)/float32(height)*2 - 1
	event.Inside = cursorX >= float32(x) && cursorX < float32(x+width) && cursorY >= float32(y) && cursorY < float32(y+height)
	return event
}

// HandlePointer dispatches input to visible game scenes in mounting order.
func (g *Game) HandlePointer(event PointerEvent) {
	for _, scene := range g.scenes {
		if scene.root == nil || !scene.root.Visible() {
			continue
		}
		for _, child := range scene.Components() {
			if handler, ok := child.(PointerHandler); ok {
				handler.HandlePointer(event)
			}
		}
	}
}
