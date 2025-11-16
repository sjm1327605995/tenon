package yoga

type FlexLineRunningLayout struct {
	// Total flex grow factors of flex items which are to be laid in the current
	// line. This is decremented as free space is distributed.
	totalFlexGrowFactors float32

	// Total flex shrink factors of flex items which are to be laid in the current
	// line. This is decremented as free space is distributed.
	totalFlexShrinkScaledFactors float32

	// The amount of available space within inner dimensions of the line which may
	// still be distributed.
	remainingFreeSpace float32

	// The size of the mainDim for the row after considering size, padding, margin
	// and border of flex items. This is used to calculate maxLineDim after going
	// through all the rows to decide on the main axis size of owner.
	mainDim float32

	// The size of the crossDim for the row after considering size, padding,
	// margin and border of flex items. Used for calculating containers crossSize.
	crossDim float32
}

type FlexLine struct {
	// List of children which are part of the line flow. This means they are not
	// positioned absolutely, or with `display: "none"`, and do not overflow the
	// available dimensions.
	itemsInFlow []*Node

	// Accumulation of the dimensions and margin of all the children on the
	// current line. This will be used in order to either set the dimensions of
	// the node if none already exist or to compute the remaining space left for
	// the flexible children.
	sizeConsumed float32

	// The index of the first item beyond the current line.
	endOfLineIndex uint32

	// Layout information about the line computed in steps after line-breaking
	layout FlexLineRunningLayout
}

// Calculates where a line starting at a given index should break, returning
// information about the collective children on the line.
//
// This function assumes that all the children of node have their
// computedFlexBasis properly computed (To do this, use the
// computeFlexBasisForChildren function).
func calculateFlexLine(
	node *Node,
	ownerDirection Direction,
	mainAxisOwnerSize float32,
	availableInnerWidth float32,
	availableInnerMainDim float32,
	startOfLineIndex uint32,
	lineCount uint32,
) FlexLine {
	itemsInFlow := make([]*Node, 0, len(node.GetChildren()))

	sizeConsumed := float32(0)
	totalFlexGrowFactors := float32(0)
	totalFlexShrinkScaledFactors := float32(0)
	endOfLineIndex := startOfLineIndex
	firstElementInLineIndex := startOfLineIndex

	sizeConsumedIncludingMinConstraint := float32(0)
	mainAxis := resolveDirection(node.getStyle().flexDirection(), node.resolveDirection(ownerDirection))
	isNodeFlexWrap := node.getStyle().flexWrap() != WrapNoWrap
	gap := node.getGapForAxis(mainAxis)

	for ; endOfLineIndex < uint32(len(node.GetChildren())); endOfLineIndex++ {
		child := node.GetChild(endOfLineIndex)
		if child.getStyle().display() == DisplayNone || child.getStyle().positionType() == PositionTypeAbsolute {
			if firstElementInLineIndex == endOfLineIndex {
				firstElementInLineIndex++
			}
			continue
		}

		isFirstElementInLine := (endOfLineIndex - firstElementInLineIndex) == 0

		child.setLineIndex(lineCount)
		childMarginMainAxis := child.getMarginForAxis(mainAxis, availableInnerWidth)
		childLeadingGapMainAxis := float32(0)
		if !isFirstElementInLine {
			childLeadingGapMainAxis = gap
		}
		flexBasisWithMinAndMaxConstraints := boundAxisWithinMinAndMax(
			child,
			mainAxis,
			child.getLayout().computedFlexBasis,
			mainAxisOwnerSize,
		).unwrap()

		if sizeConsumedIncludingMinConstraint+flexBasisWithMinAndMaxConstraints+childMarginMainAxis+childLeadingGapMainAxis > availableInnerMainDim &&
			isNodeFlexWrap && len(itemsInFlow) > 0 {
			break
		}

		sizeConsumedIncludingMinConstraint += flexBasisWithMinAndMaxConstraints + childMarginMainAxis + childLeadingGapMainAxis
		sizeConsumed += flexBasisWithMinAndMaxConstraints + childMarginMainAxis + childLeadingGapMainAxis

		if child.isNodeFlexible() {
			totalFlexGrowFactors += child.resolveFlexGrow()
			totalFlexShrinkScaledFactors += -child.resolveFlexShrink() * child.getLayout().computedFlexBasis.unwrap()
		}

		itemsInFlow = append(itemsInFlow, child)
	}

	if totalFlexGrowFactors > 0 && totalFlexGrowFactors < 1 {
		totalFlexGrowFactors = 1
	}

	if totalFlexShrinkScaledFactors > 0 && totalFlexShrinkScaledFactors < 1 {
		totalFlexShrinkScaledFactors = 1
	}

	return FlexLine{
		itemsInFlow:    itemsInFlow,
		sizeConsumed:   sizeConsumed,
		endOfLineIndex: endOfLineIndex,
		layout: FlexLineRunningLayout{
			totalFlexGrowFactors:         totalFlexGrowFactors,
			totalFlexShrinkScaledFactors: totalFlexShrinkScaledFactors,
		},
	}
}
