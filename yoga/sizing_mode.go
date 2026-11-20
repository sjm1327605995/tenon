package yoga

// SizingMode describes how a box is sized along one axis. MeasureMode remains
// the callback-facing API; the two representations are losslessly convertible.
type SizingMode uint8

const (
	SizingModeStretchFit SizingMode = iota
	SizingModeMaxContent
	SizingModeFitContent
)

// MeasureMode returns the equivalent callback measurement constraint.
func (m SizingMode) MeasureMode() MeasureMode {
	switch m {
	case SizingModeStretchFit:
		return MeasureModeExactly
	case SizingModeMaxContent:
		return MeasureModeUndefined
	case SizingModeFitContent:
		return MeasureModeAtMost
	default:
		panic("invalid sizing mode")
	}
}

// SizingModeFromMeasureMode converts a callback measurement constraint.
func SizingModeFromMeasureMode(mode MeasureMode) SizingMode {
	switch mode {
	case MeasureModeExactly:
		return SizingModeStretchFit
	case MeasureModeUndefined:
		return SizingModeMaxContent
	case MeasureModeAtMost:
		return SizingModeFitContent
	default:
		panic("invalid measure mode")
	}
}
