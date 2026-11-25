package g3n

import (
	"testing"
	"time"

	"github.com/g3n/engine/core"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

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

type pointerTestComponent struct{ events []PointerEvent }

func (c *pointerTestComponent) Update(time.Duration)         {}
func (c *pointerTestComponent) Close()                       {}
func (c *pointerTestComponent) HandlePointer(e PointerEvent) { c.events = append(c.events, e) }

func TestHandlePointerDispatchesToMountedSceneComponents(t *testing.T) {
	component := &pointerTestComponent{}
	scene := GameSceneRuntime(ui.NewScene3DRuntime([]ui.Game3DComponent{component}))
	scene.root = core.NewNode()
	game := &Game{scenes: []*Scene{scene}}
	game.HandlePointer(PointerEvent{X: 0.25, Y: -0.5, Inside: true, LeftClick: true})
	if len(component.events) != 1 || component.events[0].X != 0.25 || !component.events[0].LeftClick {
		t.Fatalf("component pointer events = %+v", component.events)
	}
	scene.root.SetVisible(false)
	game.HandlePointer(PointerEvent{Inside: true})
	if len(component.events) != 1 {
		t.Fatal("hidden scene still received pointer input")
	}
}
