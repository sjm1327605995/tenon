package g3n

import (
	"math"
	"testing"

	"github.com/g3n/engine/core"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

func TestLineUsesOneContinuousPanel(t *testing.T) {
	root := core.NewNode()
	painter := newPainter(root, nil, 100, 100)
	painter.Line(10, 20, 30, 20, ui.White)
	if len(root.Children()) != 1 {
		t.Fatalf("horizontal line produced %d panels, want 1", len(root.Children()))
	}
	panel := root.Children()[0].(layerPanel)
	assertPanelGeometry(t, panel, 10, 20, 20, 1)
}

func TestDiagonalLinePreservesLengthAndAngle(t *testing.T) {
	root := core.NewNode()
	painter := newPainter(root, nil, 100, 100)
	painter.Line(10, 10, 13, 14, ui.White)
	panel := root.Children()[0].(layerPanel)
	assertPanelGeometry(t, panel, 9, 12, 5, 1)
	if angle := panel.GetPanel().Rotation().Z; !closeFloat(angle, -float32(math.Atan2(4, 3))) {
		t.Fatalf("line angle = %v, want -atan2(4,3)", angle)
	}
}

func TestDegenerateLineDoesNotDraw(t *testing.T) {
	root := core.NewNode()
	painter := newPainter(root, nil, 100, 100)
	painter.Line(10, 10, 10, 10, ui.White)
	if len(root.Children()) != 0 {
		t.Fatal("zero-length line produced a panel")
	}
}
