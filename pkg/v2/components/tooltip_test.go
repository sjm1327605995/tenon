package components

import (
	"testing"
)

func TestTooltip_ChainAPI(t *testing.T) {
	tt := NewTooltip("hint")
	if tt.SetContent("new hint") != tt {
		t.Error("SetContent should return *Tooltip")
	}
	if tt.SetTextColor(nil) != tt {
		t.Error("SetTextColor should return *Tooltip")
	}
	if tt.Show() != tt {
		t.Error("Show should return *Tooltip")
	}
	if tt.Hide() != tt {
		t.Error("Hide should return *Tooltip")
	}
}

func TestTooltip_Visibility(t *testing.T) {
	tt := NewTooltip("hint")
	if tt.IsVisible() {
		t.Fatal("tooltip should be hidden by default")
	}
	tt.Show()
	if !tt.IsVisible() {
		t.Fatal("tooltip should be visible after Show")
	}
	tt.Hide()
	if tt.IsVisible() {
		t.Fatal("tooltip should be hidden after Hide")
	}
}

func TestTooltip_ElementType(t *testing.T) {
	tt := NewTooltip("hint")
	if tt.ElementType() != "Tooltip" {
		t.Fatalf("expected ElementType Tooltip, got %s", tt.ElementType())
	}
}
