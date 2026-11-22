package g3n

import "testing"

func TestPointerEventUsesSceneViewportCoordinates(t *testing.T) {
	event := pointerEventAt(790, 366.5, 300, 99, 980, 535, true, false)
	if !event.Inside || event.X != 0 || event.Y != 0 || !event.LeftClick {
		t.Fatalf("viewport center = %+v", event)
	}
	event = pointerEventAt(300, 99, 300, 99, 980, 535, false, false)
	if !event.Inside || event.X != -1 || event.Y != 1 {
		t.Fatalf("viewport top-left = %+v", event)
	}
	event = pointerEventAt(100, 300, 300, 99, 980, 535, true, false)
	if event.Inside || !event.LeftClick {
		t.Fatalf("outside click = %+v", event)
	}
}
