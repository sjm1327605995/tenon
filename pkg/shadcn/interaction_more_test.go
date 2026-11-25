package shadcn

import (
	"fmt"
	"testing"
	"time"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func TestSliderDragUpdatesControlledValue(t *testing.T) {
	app := func(_ struct{}) *ui.Node {
		value, set := ui.UseState(float32(5))
		return ui.Div(
			Slider(SliderProps{Value: value, Min: 0, Max: 10, OnChange: func(v float32) { set(v) }}),
			ui.Text(fmt.Sprint(value)),
		)
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))
	if !h.Root().Child(0).Child(2).Drag(55, 0) {
		t.Fatalf("slider thumb did not accept drag; root children=%d", h.Root().Count())
	}
	if !h.Root().ByText("7.5").Exists() {
		t.Fatalf("slider value texts=%v, want 7.5", h.Root().Texts())
	}
}

func TestToggleAndToggleGroupUpdateControlledState(t *testing.T) {
	app := func(_ struct{}) *ui.Node {
		pressed, setPressed := ui.UseState(false)
		selected, setSelected := ui.UseState("A")
		return ui.Div(
			Toggle(ToggleProps{Pressed: pressed, OnChange: func(v bool) { setPressed(v) }}, ui.Text("Toggle")),
			ToggleGroup(ToggleGroupProps{Value: selected, Options: []string{"A", "B"}, OnChange: func(v string) { setSelected(v) }}),
			ui.Text(fmt.Sprintf("%t/%s", pressed, selected)),
		)
	}
	h := ui.MountDefault(ui.Use(app, struct{}{}))
	if !h.Root().ByText("Toggle").Click() || !h.Root().ByText("true/A").Exists() {
		t.Fatal("toggle did not update controlled pressed state")
	}
	if !h.Root().ByText("B").Click() || !h.Root().ByText("true/B").Exists() {
		t.Fatal("toggle group did not update controlled selection")
	}
}

func TestCarouselNextMovesSlideStrip(t *testing.T) {
	h := ui.MountDefault(Carousel(CarouselProps{
		Slides: []*ui.Node{ui.Text("first"), ui.Text("second")}, Width: 200, Height: 100,
	}))
	defer h.Close()
	if !h.Root().ByText("›").Click() {
		t.Fatal("carousel next button did not handle click")
	}
	initialX := h.Root().Child(0).Child(0).VisualBounds().X
	h.Step(350)
	if got := h.Root().Child(0).Child(0).VisualBounds().X; got != initialX-200 {
		t.Fatalf("slide strip x=%v, want %v", got, initialX-200)
	}
}

func TestResizableDragChangesPanelBounds(t *testing.T) {
	h := ui.MountDefault(Resizable(ResizableProps{
		Width: 300, Height: 100, Left: ui.Text("left"), Right: ui.Text("right"),
	}))
	defer h.Close()
	if got := h.Root().Child(0).Bounds().W; got != 150 {
		t.Fatalf("initial left width=%v, want 150", got)
	}
	if !h.Root().Child(1).Drag(75, 0) {
		t.Fatal("resizable divider did not accept drag")
	}
	if got := h.Root().Child(0).Bounds().W; got != 225 {
		t.Fatalf("dragged left width=%v, want 225", got)
	}
}

func TestToastExpiresAndCleansOverlay(t *testing.T) {
	toastState.items = nil
	toastState.bump = nil
	now := time.Now()
	previousNow := toastNow
	toastNow = func() time.Time { return now }
	t.Cleanup(func() {
		toastState.items = nil
		toastState.bump = nil
		toastNow = previousNow
	})
	h := ui.MountDefault(Toaster())
	defer h.Close()
	Toast("Saved", "complete")
	h.Step(0)
	if len(h.Overlays()) != 1 || !h.Overlays()[0].ByText("Saved").Exists() {
		t.Fatal("toast was not rendered in the overlay")
	}
	now = now.Add(4 * time.Second)
	h.Step(16)
	if len(h.Overlays()) != 0 || len(toastState.items) != 0 {
		t.Fatalf("expired toast was not removed: overlays=%d items=%d", len(h.Overlays()), len(toastState.items))
	}
}
