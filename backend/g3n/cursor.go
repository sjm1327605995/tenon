package g3n

import (
	"github.com/g3n/engine/window"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

func nativeCursor(role ui.CursorRole) window.Cursor {
	switch role {
	case ui.CursorText:
		return window.IBeamCursor
	case ui.CursorPointer:
		return window.HandCursor
	default:
		return window.ArrowCursor
	}
}
