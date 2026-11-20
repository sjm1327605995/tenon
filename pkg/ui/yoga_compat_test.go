package ui

import (
	"testing"

	"github.com/sjm1327605995/tenon/yoga"
)

func TestTenonLayoutDefaultsToBorderBox(t *testing.T) {
	h := MountDefault(Div(Style(Width(100), Height(40), Padding(10), Border(5, Black))))
	defer h.Close()

	if got := h.Root().Bounds(); got.W != 100 || got.H != 40 {
		t.Fatalf("root bounds=%+v, want 100x40 with border-box sizing", got)
	}
	if got := h.g.rootRN.yn.StyleGetBoxSizing(); got != yoga.BoxSizingBorderBox {
		t.Fatalf("Tenon box sizing=%v, want border-box", got)
	}
}
