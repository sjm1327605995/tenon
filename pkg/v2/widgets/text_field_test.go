package widgets

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/yoga"
)

func TestTextFieldLayout(t *testing.T) {
	// Create a TextField widget
	tf := TextField("hello").
		W(200).
		H(40).
		Placeholder("Enter text").
		Pad(ui.EdgeInsetsAll(8))

	// Create element and mount
	el := tf.CreateElement().(*TextFieldElement)
	el.Mount(nil, 0)

	// Get render objects
	box := el.GetRenderObject().(*render.RenderBox)
	children := box.GetChildren()
	if len(children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(children))
	}
	editable := children[0].(*render.RenderEditableText)

	// Manually run yoga layout to get bounds
	box.GetYoga().StyleSetWidth(200)
	box.GetYoga().StyleSetHeight(40)
	box.GetYoga().StyleSetPadding(yoga.EdgeTop, 8)
	box.GetYoga().StyleSetPadding(yoga.EdgeRight, 8)
	box.GetYoga().StyleSetPadding(yoga.EdgeBottom, 8)
	box.GetYoga().StyleSetPadding(yoga.EdgeLeft, 8)

	editable.GetYoga().StyleSetWidthAuto()
	editable.GetYoga().StyleSetHeightAuto()

	box.GetYoga().CalculateLayout(200, 40, yoga.DirectionLTR)

	// Sync bounds
	var syncBounds func(ro render.RenderObject)
	syncBounds = func(ro render.RenderObject) {
		yn := ro.GetYoga()
		if yn != nil {
			left := yn.LayoutLeft()
			top := yn.LayoutTop()
			w := yn.LayoutWidth()
			h := yn.LayoutHeight()
			ro.SetBounds(render.Bounds{X: left, Y: top, Width: w, Height: h})
		}
		for _, child := range ro.GetChildren() {
			syncBounds(child)
		}
	}
	syncBounds(box)

	t.Logf("RenderBox bounds: %+v", box.GetBounds())
	t.Logf("RenderEditableText bounds: %+v", editable.GetBounds())

	// Hit test at center of TextField
	boxBounds := box.GetBounds()
	centerX := boxBounds.X + boxBounds.Width/2
	centerY := boxBounds.Y + boxBounds.Height/2

	t.Logf("HitTest center (%v, %v): box=%v editable=%v",
		centerX, centerY,
		box.HitTest(centerX, centerY),
		editable.HitTest(centerX-boxBounds.X, centerY-boxBounds.Y))

	// Hit test at empty content
	editable.SetContent("")
	editable.GetYoga().CalculateLayout(184, 40, yoga.DirectionLTR)
	syncBounds(editable)
	t.Logf("After empty content, editable bounds: %+v", editable.GetBounds())
}
