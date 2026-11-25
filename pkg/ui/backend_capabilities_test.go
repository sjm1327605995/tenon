package ui

import (
	"image"
	"testing"
)

type plainTestBackend struct{}

func (plainTestBackend) NewFont(float32, int, bool) FontFace { return nil }
func (plainTestBackend) NewBitmap(image.Image) Bitmap        { return nil }
func (plainTestBackend) NewVecPath(string, float32) VecPath  { return nil }
func (plainTestBackend) Run(*Node, WindowConfig)             {}

type capabilityTestBackend struct {
	plainTestBackend
	capabilities BackendCapabilities
}

func (b capabilityTestBackend) Capabilities() BackendCapabilities { return b.capabilities }

func TestRunWindowsRequiresGioBackend(t *testing.T) {
	previous := activeBackend
	t.Cleanup(func() { SetBackend(previous) })
	SetBackend(plainTestBackend{})
	if err := RunWindows(Text("window")); err == nil {
		t.Fatal("RunWindows accepted a non-Gio backend")
	}
	if err := RunWindows(); err != nil {
		t.Fatalf("empty RunWindows returned %v", err)
	}
}

func TestBackendCapabilities(t *testing.T) {
	previous := activeBackend
	t.Cleanup(func() { SetBackend(previous) })

	if got := CurrentBackendCapabilities(); !got.SVGPaths || !got.RoundedClipping || !got.LayerTransforms || !got.ExactProjection || !got.Accessibility {
		t.Fatalf("Gio capabilities=%+v, want all advanced features", got)
	}
	SetBackend(plainTestBackend{})
	if got := CurrentBackendCapabilities(); got != (BackendCapabilities{}) {
		t.Fatalf("undeclared backend capabilities=%+v, want conservative zero value", got)
	}
	want := BackendCapabilities{SVGPaths: true, LayerTransforms: true}
	SetBackend(capabilityTestBackend{capabilities: want})
	if got := CurrentBackendCapabilities(); got != want {
		t.Fatalf("declared backend capabilities=%+v, want %+v", got, want)
	}
}
