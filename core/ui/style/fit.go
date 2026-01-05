package style

type Fit uint8

const (
	// Unscaled does not alter the scale of a widget.
	Unscaled Fit = iota
	// Contain scales widget as large as possible without cropping
	// and it preserves aspect-ratio.
	Contain
	// Cover scales the widget to cover the constraint area and
	// preserves aspect-ratio.
	Cover
	// ScaleDown scales the widget smaller without cropping,
	// when it exceeds the constraint area.
	// It preserves aspect-ratio.
	ScaleDown
	// Fill stretches the widget to the constraints and does not
	// preserve aspect-ratio.
	Fill
)
