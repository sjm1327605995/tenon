package ui

// CursorRole is the backend-neutral pointer shape for the current hit target.
type CursorRole uint8

const (
	CursorDefault CursorRole = iota
	CursorPointer
	CursorText
)

// PointerCursor reports the pointer shape for physical pixel coordinates.
// Disabled targets use CursorDefault and do not fall through to underlying controls.
func (r *BackendRuntime) PointerCursor(x, y float32) CursorRole {
	if r == nil || r.closed {
		return CursorDefault
	}
	for n := r.g.hitTop(x, y); n != nil; n = n.parent {
		if n.isDisabled() {
			return CursorDefault
		}
		if n.kind == rnInput {
			return CursorText
		}
		if n.onClick != nil {
			return CursorPointer
		}
	}
	return CursorDefault
}
