package yoga

func paddingAndBorderForAxis(node *Node, axis FlexDirection, widthSize float32) float32 {
	// The total padding/border for a given axis does not depend on the direction
	// so hardcoding LTR here to avoid piping direction to this function
	return node.getInlineStartPaddingAndBorder(axis, DirectionLTR, widthSize) +
		node.getInlineEndPaddingAndBorder(axis, DirectionLTR, widthSize)
}

func boundAxisWithinMinAndMax(node *Node, axis FlexDirection, value FloatOptional, axisSize float32) FloatOptional {
	var min, max FloatOptional

	if isColumn(axis) {
		min = resolveValue(
			node.getStyle().minDimension(DimensionHeight).Value(), axisSize)
		max = resolveValue(
			node.getStyle().maxDimension(DimensionHeight).Value(), axisSize)
	} else if isRow(axis) {
		min = resolveValue(
			node.getStyle().minDimension(DimensionWidth).Value(), axisSize)
		max = resolveValue(
			node.getStyle().maxDimension(DimensionWidth).Value(), axisSize)
	}

	if max.unwrap() >= 0 && value.unwrap() > max.unwrap() {
		return max
	}

	if min.unwrap() >= 0 && value.unwrap() < min.unwrap() {
		return min
	}

	return value
}

// Like boundAxisWithinMinAndMax but also ensures that the value doesn't
// go below the padding and border amount.
func boundAxis(
	node *Node,
	axis FlexDirection,
	value float32,
	axisSize float32,
	widthSize float32,
) float32 {
	return maxOrDefined(
		boundAxisWithinMinAndMax(node, axis, FloatOptional{value}, axisSize).unwrap(),
		paddingAndBorderForAxis(node, axis, widthSize),
	)
}
