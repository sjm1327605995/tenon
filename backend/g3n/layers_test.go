package g3n

import (
	"math"
	"testing"

	"github.com/g3n/engine/core"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

func TestLayerScaleTranslateAndOpacity(t *testing.T) {
	root := core.NewNode()
	painter := newPainter(root, nil, 200, 200)
	painter.panel(10, 20, 20, 10, ui.White)
	painter.applyLayer(0, ui.Layer{Scale: 2, ScaleY: 3, TX: 4, TY: 5, Opacity: 0.5})
	panel := root.Children()[0].(layerPanel)
	assertPanelGeometry(t, panel, 24, 65, 40, 30)
	if alpha := panel.Color4().A; alpha != 0.5 {
		t.Fatalf("layer alpha = %v, want 0.5", alpha)
	}
}

func TestLayerRotationUsesSharedCenter(t *testing.T) {
	root := core.NewNode()
	painter := newPainter(root, nil, 200, 200)
	painter.panel(10, 10, 10, 10, ui.White)
	painter.panel(40, 10, 10, 10, ui.White)
	painter.applyLayer(0, ui.Layer{CX: 30, CY: 15, Scale: 1, ScaleY: 1, Rotate: 90, TX: 5, TY: -2, Opacity: 1})
	first := root.Children()[0].(layerPanel)
	second := root.Children()[1].(layerPanel)
	assertPanelGeometry(t, first, 30, -7, 10, 10)
	assertPanelGeometry(t, second, 30, 23, 10, 10)
	for _, panel := range []layerPanel{first, second} {
		if angle := panel.GetPanel().Rotation().Z; !closeFloat(angle, -math.Pi/2) {
			t.Fatalf("panel rotation = %v, want -pi/2", angle)
		}
	}
}

func assertPanelGeometry(t *testing.T, panel layerPanel, x, y, width, height float32) {
	t.Helper()
	position := panel.Position()
	if !closeFloat(position.X, x) || !closeFloat(position.Y, y) || !closeFloat(panel.Width(), width) || !closeFloat(panel.Height(), height) {
		t.Fatalf("panel geometry = (%v,%v %vx%v), want (%v,%v %vx%v)", position.X, position.Y, panel.Width(), panel.Height(), x, y, width, height)
	}
}

func closeFloat(a, b float32) bool {
	if a > b {
		return a-b < 0.0001
	}
	return b-a < 0.0001
}
