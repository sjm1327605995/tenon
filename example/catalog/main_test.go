package main

import (
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func TestCatalogControlsAreInteractive(t *testing.T) {
	app := &catalogApp{}
	h := ui.Mount(ui.Render(app), 1000, 800)
	defer h.Close()

	if !h.Root().ByText("Continue").Click() || !h.Root().ByText("Button clicked").Exists() {
		t.Fatal("catalog button did not update interaction status")
	}
	h.Root().ByPlaceholder("Email").SetValue("user@example.com")
	if got := app.text.Get(); got != "user@example.com" {
		t.Fatalf("input binding=%q, want user@example.com", got)
	}
	clickAccessibility(t, h, ui.AccessibilityCheckbox, "")
	if !app.checked {
		t.Fatal("catalog checkbox did not toggle application state")
	}
	clickAccessibility(t, h, ui.AccessibilitySwitch, "")
	if !app.switched {
		t.Fatal("catalog switch did not toggle application state")
	}
	clickAccessibility(t, h, ui.AccessibilityRadioButton, "B")
	if app.radio != "B" {
		t.Fatalf("catalog radio selection=%q, want B", app.radio)
	}

	scrollQueryIntoView(h, h.Root().ByText("Open dialog"))
	if !h.Root().ByText("Open dialog").Click() {
		t.Fatalf("dialog trigger found no click handler; bounds=%+v", h.Root().ByText("Open dialog").Bounds())
	}
	if len(h.Overlays()) == 0 {
		t.Fatalf("dialog trigger did not open the dialog; open=%t", app.dialog)
	}
	h.Escape()
	h.Step(220)
	if app.dialog || len(h.Overlays()) != 0 {
		t.Fatal("Escape did not close the catalog dialog")
	}
}

func clickAccessibility(t *testing.T, h *ui.Harness, role ui.AccessibilityRole, label string) {
	t.Helper()
	for _, node := range h.AccessibilitySnapshot() {
		if node.Role == role && (label == "" || node.Label == label) {
			_, height := h.Size()
			if node.Bounds.Y < 0 || node.Bounds.Y+node.Bounds.H > float32(height) {
				h.Root().ScrollBy(node.Bounds.Y - 100)
				for _, current := range h.AccessibilitySnapshot() {
					if current.Role == role && (label == "" || current.Label == label) {
						node = current
						break
					}
				}
			}
			if !h.ClickAt(node.Bounds.X+node.Bounds.W/2, node.Bounds.Y+node.Bounds.H/2) {
				t.Fatalf("accessible %v %q did not handle click at %+v", role, label, node.Bounds)
			}
			return
		}
	}
	t.Fatalf("accessible %v %q was not found", role, label)
}

func scrollQueryIntoView(h *ui.Harness, query *ui.Query) {
	_, height := h.Size()
	bounds := query.Bounds()
	if bounds.Y < 0 || bounds.Y+bounds.H > float32(height) {
		h.Root().ScrollBy(bounds.Y - 100)
	}
}

func TestCatalogUsesIntrinsicControlWidths(t *testing.T) {
	h := ui.Mount(ui.Render(&catalogApp{}), 1000, 800)
	defer h.Close()
	button := h.Root().Find(func(q *ui.Query) bool {
		return q.Clickable() && q.AllText() == "Continue"
	})
	if !button.Exists() {
		t.Fatal("Continue button not found")
	}
	if width := button.Bounds().W; width <= 0 || width > 200 {
		t.Fatalf("Continue button width = %v, want intrinsic width below 200", width)
	}
	badges := h.Root().FindAll(func(q *ui.Query) bool {
		return q.Kind() == "box" && q.AllText() == "New"
	})
	if len(badges) == 0 {
		t.Fatal("New badge not found")
	}
	width := float32(0)
	for _, badge := range badges {
		if w := badge.Bounds().W; width == 0 || w < width {
			width = w
		}
	}
	if width <= 0 || width > 100 {
		t.Fatalf("New badge width = %v, want intrinsic width below 100", width)
	}
}
