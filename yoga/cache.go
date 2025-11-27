package yoga

func sizeIsExactAndMatchesOldMeasuredSize(sizeMode MeasureMode, size, lastComputedSize float32) bool {
	return sizeMode == MeasureModeExactly && inexactEqual(size, lastComputedSize)
}

func oldSizeIsUnspecifiedAndStillFits(sizeMode MeasureMode, size float32, lastSizeMode MeasureMode, lastComputedSize float32) bool {
	return sizeMode == MeasureModeAtMost && lastSizeMode == MeasureModeUndefined &&
		(size >= lastComputedSize || inexactEqual(size, lastComputedSize))
}

func newMeasureSizeIsStricterAndStillValid(sizeMode MeasureMode, size float32, lastSizeMode MeasureMode, lastSize, lastComputedSize float32) bool {
	return lastSizeMode == MeasureModeAtMost && sizeMode == MeasureModeAtMost && !IsNaN(lastSize) &&
		!IsNaN(size) && !IsNaN(lastComputedSize) && lastSize > size &&
		(lastComputedSize <= size || inexactEqual(size, lastComputedSize))
}

func canUseCachedMeasurement(widthMode MeasureMode, availableWidth float32, heightMode MeasureMode, availableHeight float32,
	lastWidthMode MeasureMode, lastAvailableWidth float32, lastHeightMode MeasureMode, lastAvailableHeight float32,
	lastComputedWidth, lastComputedHeight, marginRow, marginColumn float32, config *Config) bool {
	if (!IsNaN(lastComputedHeight) && lastComputedHeight < 0) ||
		(!IsNaN(lastComputedWidth) && lastComputedWidth < 0) {
		return false
	}

	pointScaleFactor := config.GetPointScaleFactor()

	useRoundedComparison := config != nil && pointScaleFactor != 0
	effectiveWidth := availableWidth
	effectiveHeight := availableHeight
	effectiveLastWidth := lastAvailableWidth
	effectiveLastHeight := lastAvailableHeight

	if useRoundedComparison {
		effectiveWidth = RoundValueToPixelGrid(float64(availableWidth), float64(pointScaleFactor), false, false)
		effectiveHeight = RoundValueToPixelGrid(float64(availableHeight), float64(pointScaleFactor), false, false)
		effectiveLastWidth = RoundValueToPixelGrid(float64(lastAvailableWidth), float64(pointScaleFactor), false, false)
		effectiveLastHeight = RoundValueToPixelGrid(float64(lastAvailableHeight), float64(pointScaleFactor), false, false)
	}

	hasSameWidthSpec := lastWidthMode == widthMode && inexactEqual(effectiveLastWidth, effectiveWidth)
	hasSameHeightSpec := lastHeightMode == heightMode && inexactEqual(effectiveLastHeight, effectiveHeight)

	widthIsCompatible :=
		hasSameWidthSpec ||
			sizeIsExactAndMatchesOldMeasuredSize(widthMode, availableWidth-marginRow, lastComputedWidth) ||
			oldSizeIsUnspecifiedAndStillFits(
				widthMode,
				availableWidth-marginRow,
				lastWidthMode,
				lastComputedWidth) ||
			newMeasureSizeIsStricterAndStillValid(
				widthMode,
				availableWidth-marginRow,
				lastWidthMode,
				lastAvailableWidth,
				lastComputedWidth)

	heightIsCompatible :=
		hasSameHeightSpec ||
			sizeIsExactAndMatchesOldMeasuredSize(
				heightMode, availableHeight-marginColumn, lastComputedHeight) ||
			oldSizeIsUnspecifiedAndStillFits(
				heightMode,
				availableHeight-marginColumn,
				lastHeightMode,
				lastComputedHeight) ||
			newMeasureSizeIsStricterAndStillValid(
				heightMode,
				availableHeight-marginColumn,
				lastHeightMode,
				lastAvailableHeight,
				lastComputedHeight)

	return widthIsCompatible && heightIsCompatible
}
