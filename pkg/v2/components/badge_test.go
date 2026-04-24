package components

import (
	"testing"
)

func TestBadge_ChainAPI(t *testing.T) {
	b := NewBadge("5")
	if b.SetCount(10) != b {
		t.Error("SetCount should return *Badge")
	}
	if b.SetMaxCount(9) != b {
		t.Error("SetMaxCount should return *Badge")
	}
	if b.SetDotMode(true) != b {
		t.Error("SetDotMode should return *Badge")
	}
	if b.SetTextColor(nil) != b {
		t.Error("SetTextColor should return *Badge")
	}
}

func TestBadge_Overflow(t *testing.T) {
	b := NewBadge("").SetCount(150).SetMaxCount(99)
	if b.textEl == nil {
		t.Skip("textEl nil in dot mode")
	}
	if b.textEl.content != "99+" {
		t.Fatalf("expected '99+', got %q", b.textEl.content)
	}
}

func TestBadge_DotMode(t *testing.T) {
	b := NewBadge("")
	if !b.dotMode {
		t.Fatal("empty badge should be dot mode")
	}
	b.SetDotMode(false)
	if b.dotMode {
		t.Fatal("SetDotMode(false) should disable dot mode")
	}
}

func TestBadge_ElementType(t *testing.T) {
	b := NewBadge("1")
	if b.ElementType() != "Badge" {
		t.Fatalf("expected ElementType Badge, got %s", b.ElementType())
	}
}
