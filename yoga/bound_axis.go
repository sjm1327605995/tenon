package yoga

func paddingAndBorderForAxis(node *Node, axis FlexDirection, widthSize float32) float32 {
	// The total padding/border for an axis does not depend on direction.
	return node.getInlineStartPaddingAndBorder(axis, DirectionLTR, widthSize) +
		node.getInlineEndPaddingAndBorder(axis, DirectionLTR, widthSize)
}

func resolveSizeWithBoxSizing(node *Node, value Value, referenceLength, ownerWidth float32, dimension Dimension) FloatOptional {
	resolved := resolveValue(value, referenceLength)
	if node.getStyle().boxSizing() == BoxSizingBorderBox || resolved.isUndefined() {
		return resolved
	}
	axis := FlexDirectionRow
	if dimension == DimensionHeight {
		axis = FlexDirectionColumn
	}
	return FloatOptional{resolved.unwrap() + paddingAndBorderForAxis(node, axis, ownerWidth)}
}

func resolveDimensionSize(node *Node, dimension Dimension, referenceLength, ownerWidth float32) FloatOptional {
	return resolveSizeWithBoxSizing(node, node.getResolvedDimension(dimension), referenceLength, ownerWidth, dimension)
}

func boundAxisWithinMinAndMax(node *Node, axis FlexDirection, value FloatOptional, axisSize, widthSize float32) FloatOptional {
	dimension := dimension(axis)
	min := resolveSizeWithBoxSizing(node, node.getStyle().minDimension(dimension).Value(), axisSize, widthSize, dimension)
	max := resolveSizeWithBoxSizing(node, node.getStyle().maxDimension(dimension).Value(), axisSize, widthSize, dimension)

	if max.unwrap() >= 0 && value.unwrap() > max.unwrap() {
		return max
	}
	if min.unwrap() >= 0 && value.unwrap() < min.unwrap() {
		return min
	}
	return value
}

// boundAxis also ensures that the value doesn't go below padding and border.
func boundAxis(
	node *Node,
	axis FlexDirection,
	value float32,
	axisSize float32,
	widthSize float32,
) float32 {
	return maxOrDefined(
		boundAxisWithinMinAndMax(node, axis, FloatOptional{value}, axisSize, widthSize).unwrap(),
		paddingAndBorderForAxis(node, axis, widthSize),
	)
}
