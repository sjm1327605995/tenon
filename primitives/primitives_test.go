package primitives

import (
	"testing"

	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/widget"
)

func TestColumnLayout(t *testing.T) {
	c := NewColumn(
		NewText("A"),
		NewText("B"),
	)
	c.Children_[0].(*Text).FontSize = 10
	c.Children_[1].(*Text).FontSize = 10

	constraints := geometry.Constraints{
		MinWidth:  0,
		MaxWidth:  100,
		MinHeight: 0,
		MaxHeight: 100,
	}
	size := c.Layout(nil, constraints)
	if size.Width <= 0 || size.Height <= 0 {
		t.Fatalf("expected non-zero size, got %v", size)
	}

	// Each text is ~6px wide and 12px tall (with scaling).
	// Column should stack them vertically.
	if c.Children_[0].(*Text).Bounds().Min.Y != 0 {
		t.Errorf("first child should start at y=0, got %v", c.Children_[0].(*Text).Bounds().Min.Y)
	}
}

func TestBoxLayout(t *testing.T) {
	child := NewText("Hello")
	child.FontSize = 10
	b := NewBox(child)
	b.Padding = geometry.Insets{Left: 10, Top: 10, Right: 10, Bottom: 10}

	constraints := geometry.Constraints{
		MinWidth:  0,
		MaxWidth:  200,
		MinHeight: 0,
		MaxHeight: 200,
	}
	size := b.Layout(nil, constraints)
	if size.Width <= 0 || size.Height <= 0 {
		t.Fatalf("expected non-zero size, got %v", size)
	}

	// Child should be offset by padding.
	childBounds := child.Bounds()
	if childBounds.Min.X != 10 {
		t.Errorf("expected child x=10, got %v", childBounds.Min.X)
	}
	if childBounds.Min.Y != 10 {
		t.Errorf("expected child y=10, got %v", childBounds.Min.Y)
	}
}

func TestTextMeasure(t *testing.T) {
	txt := NewText("ABCD")
	txt.FontSize = 10
	w := txt.Layout(nil, geometry.Constraints{
		MaxWidth:  100,
		MaxHeight: 100,
	})
	if w.Width <= 0 {
		t.Fatalf("expected positive width, got %v", w)
	}
}

func TestRowLayout(t *testing.T) {
	r := NewRow(
		NewText("Left"),
		NewText("Right"),
	)
	r.Children_[0].(*Text).FontSize = 10
	r.Children_[1].(*Text).FontSize = 10

	constraints := geometry.Constraints{
		MinWidth:  0,
		MaxWidth:  200,
		MinHeight: 0,
		MaxHeight: 100,
	}
	size := r.Layout(nil, constraints)
	if size.Width <= 0 || size.Height <= 0 {
		t.Fatalf("expected non-zero size, got %v", size)
	}

	// Row should place children side by side.
	if r.Children_[0].(*Text).Bounds().Min.X != 0 {
		t.Errorf("first child should start at x=0, got %v", r.Children_[0].(*Text).Bounds().Min.X)
	}
}

func TestWidgetInterfaces(t *testing.T) {
	var _ widget.Widget = (*Box)(nil)
	var _ widget.Widget = (*Text)(nil)
	var _ widget.Widget = (*Row)(nil)
	var _ widget.Widget = (*Column)(nil)
}
