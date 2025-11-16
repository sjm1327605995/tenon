package yoga

type CachedMeasurement struct {
	availableWidth    float32
	availableHeight   float32
	widthMeasureMode  MeasureMode
	heightMeasureMode MeasureMode
	computedWidth     float32
	computedHeight    float32
}
