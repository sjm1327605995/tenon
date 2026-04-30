package components

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

type testWidget struct {
	core.BaseWidget
	sv *ScrollView
}

func (w *testWidget) Render() core.Element {
	root := native.NewView()
	root.SetWidthPercent(100)
	root.SetHeightPercent(100)

	w.sv = NewScrollView()
	w.sv.SetHeight(300)
	root.AppendChild(w.sv)

	Content := w.sv.Content()
	Content.SetFlexDirection(yoga.FlexDirectionColumn)
	Content.SetGap(yoga.GutterAll, 8)
	Content.SetPadding(yoga.EdgeAll, 16)

	for i := 0; i < 30; i++ {
		row := native.NewView().
			SetFlexDirection(yoga.FlexDirectionRow).
			SetAlignItems(yoga.AlignCenter).
			SetGap(yoga.GutterAll, 8)
		row.SetHeight(40)
		row.Add(
			native.NewText(fmt.Sprintf("Item %d", i+1)),
			native.NewView().SetFlexGrow(1),
			NewButton("Action"),
		)
		Content.Add(row)
	}

	return root
}

func TestScrollViewScroll(t *testing.T) {
	fonts.InitDefaultFont()

	w := &testWidget{}
	w.Init(w)

	eng := core.NewEngine(w, 900, 600)
	eng.Mount()

	dummyImg := ebiten.NewImage(1, 1)
	eng.Draw(dummyImg)

	sv := w.sv

	// Check initial state
	fmt.Printf("Initial: contentH=%.2f viewportH=%.2f maxScrollY=%.2f scrollY=%.2f engine=%v\n",
		sv.Content().GetBounds().Height, sv.GetBounds().Height, sv.maxScrollY, sv.scrollY, sv.GetEngine() != nil)

	// Simulate scroll event
	event := &core.Event{
		Type:   core.EventScroll,
		DeltaY: -1, // scroll down
		X:      sv.GetBounds().X + 10,
		Y:      sv.GetBounds().Y + 10,
	}
	consumed := sv.HandleEvent(event)
	fmt.Printf("After scroll: consumed=%v scrollY=%.2f maxScrollY=%.2f\n", consumed, sv.scrollY, sv.maxScrollY)

	// Manually trigger engine update cycle to apply layout changes
	eng.Update()
	eng.Draw(dummyImg)

	fmt.Printf("After layout: scrollY=%.2f contentY=%.2f contentH=%.2f\n",
		sv.scrollY, sv.Content().GetBounds().Y, sv.Content().GetBounds().Height)

	// Simulate mouse down and drag
	sv.HandleEvent(&core.Event{
		Type: core.EventMouseDown,
		X:    sv.GetBounds().X + 10,
		Y:    sv.GetBounds().Y + 10,
	})

	// drag down 50px
	sv.HandleEvent(&core.Event{
		Type: core.EventMouseMove,
		X:    sv.GetBounds().X + 10,
		Y:    sv.GetBounds().Y + 60,
	})
	fmt.Printf("After drag: scrollY=%.2f maxScrollY=%.2f\n", sv.scrollY, sv.maxScrollY)

	eng.Update()
	eng.Draw(dummyImg)

	fmt.Printf("Final: scrollY=%.2f contentY=%.2f\n", sv.scrollY, sv.Content().GetBounds().Y)
}

func TestScrollViewButtonClickAfterScroll(t *testing.T) {
	fonts.InitDefaultFont()

	w := &testWidget{}
	w.Init(w)

	eng := core.NewEngine(w, 900, 600)
	eng.Mount()

	dummyImg := ebiten.NewImage(1, 1)
	eng.Draw(dummyImg)

	sv := w.sv
	Content := sv.Content()

	// Get first row and its button
	row0 := Content.GetChildren()[0]
	button0 := row0.GetChildren()[2] // native.Text, native.View, Button

	fmt.Printf("Before scroll: row0Y=%.2f button0Y=%.2f button0Bounds=%+v\n",
		row0.GetBounds().Y, button0.GetBounds().Y, button0.GetBounds())

	// Simulate scroll down 100px
	sv.HandleEvent(&core.Event{
		Type:   core.EventScroll,
		DeltaY: -5,
		X:      sv.GetBounds().X + 10,
		Y:      sv.GetBounds().Y + 10,
	})

	// Trigger engine update to apply scroll offset
	eng.Update()
	eng.Draw(dummyImg)

	fmt.Printf("After scroll: scrollY=%.2f row0Y=%.2f button0Y=%.2f button0Bounds=%+v\n",
		sv.scrollY, row0.GetBounds().Y, button0.GetBounds().Y, button0.GetBounds())

	// Now we manually test if the button bounds are correct by dispatching a click event
	b := button0.GetBounds()
	clickX := b.X + b.Width/2
	clickY := b.Y + b.Height/2

	fmt.Printf("Click at: x=%.2f y=%.2f\n", clickX, clickY)

	// Check if click is inside button bounds
	inside := clickX >= b.X && clickX < b.X+b.Width && clickY >= b.Y && clickY < b.Y+b.Height
	fmt.Printf("Click inside button: %v\n", inside)

	// Also check row bounds
	rowBounds := row0.GetBounds()
	insideRow := clickX >= rowBounds.X && clickX < rowBounds.X+rowBounds.Width &&
		clickY >= rowBounds.Y && clickY < rowBounds.Y+rowBounds.Height
	fmt.Printf("Click inside row: %v\n", insideRow)
}

// hitTest mirrors Engine.hitTest for testing.
func hitTest(el core.Element, x, y float32) core.Element {
	return hitTestClipped(el, x, y, nil)
}

func hitTestClipped(el core.Element, x, y float32, clipBounds *core.LayoutBounds) core.Element {
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
	var childClip *core.LayoutBounds
	if el.HasFlag(core.FlagClipChildren) {
		childClip = &b
	} else {
		childClip = clipBounds
	}
	for i := len(children) - 1; i >= 0; i-- {
		if hit := hitTestClipped(children[i], x, y, childClip); hit != nil {
			return hit
		}
	}
	if el.GetPointerEvents() == core.PointerEventsNone {
		return nil
	}
	if x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height {
		if clipBounds != nil {
			if x >= clipBounds.X && x < clipBounds.X+clipBounds.Width &&
				y >= clipBounds.Y && y < clipBounds.Y+clipBounds.Height {
				return el
			}
			return nil
		}
		return el
	}
	return nil
}

// TestScrollViewHitTestAfterScroll verifies that hitTest finds the Button
// immediately after a scroll event, before the next Engine.Update() cycle.
func TestScrollViewHitTestAfterScroll(t *testing.T) {
	fonts.InitDefaultFont()

	w := &testWidget{}
	w.Init(w)

	eng := core.NewEngine(w, 900, 600)
	eng.Mount()

	dummyImg := ebiten.NewImage(1, 1)
	eng.Draw(dummyImg)

	sv := w.sv
	Content := sv.Content()
	row0 := Content.GetChildren()[0]
	button0 := row0.GetChildren()[2]

	// Before scroll: hitTest from root on button center should return button or its text child
	root := w.sv.GetParent()
	b0 := button0.GetBounds()
	clickX := b0.X + b0.Width/2
	clickY := b0.Y + b0.Height/2
	before := hitTest(root, clickX, clickY)
	if before != button0 && !isChildOf(before, button0) {
		t.Fatalf("before scroll: expected hitTest to return Button or its child, got %v", before.ElementType())
	}

	// Scroll down 20px via HandleEvent (simulates what Engine.handleEvents does)
	sv.HandleEvent(&core.Event{
		Type:   core.EventScroll,
		DeltaY: -1,
		X:      sv.GetBounds().X + 10,
		Y:      sv.GetBounds().Y + 10,
	})

	// After scroll but BEFORE eng.Update():
	// With the fix, applyScrollToContent is called inside HandleEvent,
	// so bounds should already reflect the new scrollY.
	b1 := button0.GetBounds()
	after := hitTest(root, clickX, b1.Y+b1.Height/2)
	if after == nil {
		t.Fatalf("after scroll (before Update): hitTest returned nil")
	}
	if after != button0 && !isChildOf(after, button0) {
		t.Fatalf("after scroll (before Update): expected hitTest to return Button or its child, got %v (bounds=%+v)",
			after.ElementType(), after.GetBounds())
	}

	// Also verify the button is no longer at its original position
	if b1.Y == b0.Y {
		t.Fatalf("expected button bounds to change after scroll, got same Y=%.2f", b1.Y)
	}
}

func isChildOf(child, parent core.Element) bool {
	for _, c := range parent.GetChildren() {
		if c == child {
			return true
		}
		if isChildOf(child, c) {
			return true
		}
	}
	return false
}
