package g3n

import (
	"testing"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

func TestVectorPanelUsesSignedDistanceMaterial(t *testing.T) {
	panel := gradientPanel(10, 20, 240, 80, 16, ui.Red, ui.Blue, 45)
	if got := panel.material.Shader(); got != "tenon_vector" {
		t.Fatalf("shader = %q, want tenon_vector", got)
	}
	if panel.material.shape.X != 240 || panel.material.shape.Y != 80 || panel.material.shape.Z != 16 {
		t.Fatalf("shape = %+v", panel.material.shape)
	}
	if panel.material.direction.X == 0 || panel.material.direction.Y == 0 {
		t.Fatal("gradient direction was not initialized")
	}
}

func TestApplyLayerPreservesGradientMaterial(t *testing.T) {
	root := core.NewNode()
	panel := gradientPanel(100, 100, 200, 50, 12, ui.White, ui.Blue, 90)
	root.Add(panel)
	painter := &painter{root: root}
	painter.applyLayer(0, ui.Layer{CX: 100, CY: 100, Scale: 2, Opacity: 0.5})
	if panel.Width() != 400 || panel.Height() != 100 {
		t.Fatalf("scaled size = %v×%v, want 400×100", panel.Width(), panel.Height())
	}
	if got := panel.material.opacity; got != 0.5 {
		t.Fatalf("material opacity = %v, want 0.5", got)
	}
}

func TestPanelClipFollowsLayerTransform(t *testing.T) {
	root := core.NewNode()
	painter := &painter{root: root}
	painter.PushClip(ui.Rect{X: 10, Y: 20, W: 30, H: 40}, 0)
	panel := painter.panel(0, 0, 100, 40, ui.White)
	painter.PopClip()
	painter.applyLayer(0, ui.Layer{Scale: 2, ScaleY: 0.5, Opacity: 1})
	if len(panel.clips) != 1 {
		t.Fatalf("clip count = %d, want 1", len(panel.clips))
	}
	clip := panel.clips[0]
	if clip.X != 20 || clip.Y != 10 || clip.W != 60 || clip.H != 20 {
		t.Fatalf("transformed clip = %+v", clip)
	}
}

func TestApplyLayerSupportsIndependentVerticalScale(t *testing.T) {
	root := core.NewNode()
	painter := &painter{root: root}
	painter.panel(0, 0, 100, 40, ui.White)
	painter.applyLayer(0, ui.Layer{Scale: 2, ScaleY: 0.5, Opacity: 1})
	panel := root.ChildAt(0).(layerPanel)
	if panel.Width() != 200 || panel.Height() != 20 {
		t.Fatalf("scaled size = %v×%v, want 200×20", panel.Width(), panel.Height())
	}
}

func TestApplyLayerScalesAndFadesLabels(t *testing.T) {
	root := core.NewNode()
	painter := &painter{root: root}
	face := (&Backend{fonts: make(map[fontKey]*font)}).NewFont(16, 400, false)
	painter.DrawText("DUELIST", face, ui.White, 100, 100, false, false)
	painter.applyLayer(0, ui.Layer{CX: 100, CY: 100, Scale: 2, Opacity: 0.5})
	panel, ok := root.ChildAt(0).(layerPanel)
	if !ok {
		t.Fatal("label does not participate in layer transforms")
	}
	if got := panel.Color4().A; got != 0.5 {
		t.Fatalf("label opacity = %v, want 0.5", got)
	}
	if got := panel.Position(); got.X != 100 || got.Y != 100 {
		t.Fatalf("label position = %+v, want (100,100)", got)
	}
}

func TestApplyLayerScalesAndFadesPanels(t *testing.T) {
	root := core.NewNode()
	painter := &painter{root: root}
	painter.panel(100, 100, 200, 50, ui.White)
	painter.applyLayer(0, ui.Layer{CX: 100, CY: 100, Scale: 2, Opacity: 0.5})
	panel := root.ChildAt(0).(layerPanel)
	if got := panel.Position(); got.X != 100 || got.Y != 100 {
		t.Fatalf("scaled position = %+v, want (100,100)", got)
	}
	if panel.Width() != 400 || panel.Height() != 100 {
		t.Fatalf("scaled size = %v×%v, want 400×100", panel.Width(), panel.Height())
	}
	if got := panel.Color4().A; got != 0.5 {
		t.Fatalf("opacity = %v, want 0.5", got)
	}
}

func TestWithTransformsOwnsGraphicRenderSetup(t *testing.T) {
	panel := withTransforms(gui.NewPanel(10, 10))
	materials := panel.GetGraphic().Materials()
	if len(materials) == 0 {
		t.Fatal("transformed panel has no graphic material")
	}
	if owner := materials[0].IGraphic(); owner != panel {
		t.Fatal("graphic material does not point to the transformed panel")
	}
}
