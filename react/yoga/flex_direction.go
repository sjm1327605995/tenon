package yoga

func isRow(flexDirection FlexDirection) bool {
	return flexDirection == FlexDirectionRow || flexDirection == FlexDirectionRowReverse
}

func isColumn(flexDirection FlexDirection) bool {
	return flexDirection == FlexDirectionColumn || flexDirection == FlexDirectionColumnReverse
}

func resolveDirection(flexDirection FlexDirection, direction Direction) FlexDirection {
	if direction == DirectionRTL {
		if flexDirection == FlexDirectionRow {
			return FlexDirectionRowReverse
		} else if flexDirection == FlexDirectionRowReverse {
			return FlexDirectionRow
		}
	}
	return flexDirection
}

func resolveCrossDirection(flexDirection FlexDirection, direction Direction) FlexDirection {
	return If(isColumn(flexDirection), resolveDirection(FlexDirectionRow, direction), FlexDirectionColumn)
}

func flexStartEdge(flexDirection FlexDirection) Edge {
	switch flexDirection {
	case FlexDirectionColumn:
		return EdgeTop
	case FlexDirectionColumnReverse:
		return EdgeBottom
	case FlexDirectionRow:
		return EdgeLeft
	case FlexDirectionRowReverse:
		return EdgeRight
	}
	panic("Invalid FlexDirection")
}

func flexEndEdge(flexDirection FlexDirection) Edge {
	switch flexDirection {
	case FlexDirectionColumn:
		return EdgeBottom
	case FlexDirectionColumnReverse:
		return EdgeTop
	case FlexDirectionRow:
		return EdgeRight
	case FlexDirectionRowReverse:
		return EdgeLeft
	}
	panic("Invalid FlexDirection")
}

func inlineStartEdge(flexDirection FlexDirection, direction Direction) Edge {
	if isRow(flexDirection) {
		return If(direction == DirectionRTL, EdgeRight, EdgeLeft)
	}
	return EdgeTop
}

func inlineEndEdge(flexDirection FlexDirection, direction Direction) Edge {
	if isRow(flexDirection) {
		return If(direction == DirectionRTL, EdgeLeft, EdgeRight)
	}
	return EdgeBottom
}

func flexStartRelativeEdge(flexDirection FlexDirection, direction Direction) Edge {
	leadLayoutEdge := inlineStartEdge(flexDirection, direction)
	leadFlexItemEdge := flexStartEdge(flexDirection)
	return If(leadLayoutEdge == leadFlexItemEdge, EdgeStart, EdgeEnd)
}

func flexEndRelativeEdge(flexDirection FlexDirection, direction Direction) Edge {
	trailLayoutEdge := inlineEndEdge(flexDirection, direction)
	trailFlexItemEdge := flexEndEdge(flexDirection)
	return If(trailLayoutEdge == trailFlexItemEdge, EdgeEnd, EdgeStart)
}

func dimension(flexDirection FlexDirection) Dimension {
	switch flexDirection {
	case FlexDirectionColumn, FlexDirectionColumnReverse:
		return DimensionHeight
	case FlexDirectionRow, FlexDirectionRowReverse:
		return DimensionWidth
	}
	panic("Invalid FlexDirection")
}
