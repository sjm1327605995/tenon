package g3n

import (
	"testing"

	"github.com/g3n/engine/window"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

func TestNativeCursorMapping(t *testing.T) {
	for role, want := range map[ui.CursorRole]window.Cursor{
		ui.CursorDefault: window.ArrowCursor,
		ui.CursorPointer: window.HandCursor,
		ui.CursorText:    window.IBeamCursor,
	} {
		if got := nativeCursor(role); got != want {
			t.Errorf("nativeCursor(%v) = %v, want %v", role, got, want)
		}
	}
}
