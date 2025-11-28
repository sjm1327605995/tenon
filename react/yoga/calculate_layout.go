package yoga

import (
	"fmt"
	"math"
	"sync/atomic"
)

var (
	gCurrentGenerationCount uint32 = 0
)

func dimensionWithMargin(node *Node, axis FlexDirection, widthSize float32) float32 {
	return node.getLayout().measuredDimension(dimension(axis)) +
		node.getMarginForAxis(axis, widthSize)
}

func styleDefinesDimension(node *Node, axis FlexDirection, ownerSize float32) bool {
	isDefined := IsDefined(node.getResolvedDimension(dimension(axis)).value)
	resolvedDimension := node.getResolvedDimension(dimension(axis))
	return !(resolvedDimension.unit == UnitAuto ||
		resolvedDimension.unit == UnitUndefined ||
		(resolvedDimension.unit == UnitPoint && isDefined &&
			resolvedDimension.value < 0.0) ||
		(resolvedDimension.unit == UnitPercent && isDefined &&
			(resolvedDimension.value < 0.0 || IsUndefined(ownerSize))))
}

func isLayoutDimensionDefined(node *Node, axis FlexDirection) bool {
	value := node.getLayout().measuredDimension(dimension(axis))
	return IsDefined(value) && value >= 0.0
}

func setChildTrailingPosition(
	node *Node,
	child *Node,
	axis FlexDirection) {
	size := child.getLayout().measuredDimension(dimension(axis))
	child.setLayoutPosition(
		node.getLayout().measuredDimension(dimension(axis))-size-
			child.getLayout().position(flexStartEdge(axis)),
		flexEndEdge(axis))
}

func constrainMaxSizeForMode(node *Node, axis FlexDirection, ownerAxisSize, ownerWidth float32, mode *MeasureMode, size *float32) {
	maxSize := resolveValue(node.getStyle().maxDimension(dimension(axis)).Value(), ownerAxisSize).unwrap() + node.getMarginForAxis(axis, ownerWidth)
	switch *mode {
	case MeasureModeExactly:
	case MeasureModeAtMost:
		*size = If(IsUndefined(maxSize) || maxSize > *size, *size, maxSize)
	case MeasureModeUndefined:
		if IsDefined(maxSize) {
			*mode = MeasureModeAtMost
			*size = maxSize
		}
	}
}

func computeFlexBasisForChild(node *Node, child *Node, width float32, widthMode MeasureMode, height float32, ownerWidth float32, ownerHeight float32, heightMode MeasureMode, direction Direction, layoutMarkerData *LayoutData, depth uint32, generationCount uint32) {
	mainAxis :=
		resolveDirection(node.getStyle().flexDirection(), direction)
	isMainAxisRow := isRow(mainAxis)
	mainAxisSize := If(isMainAxisRow, width, height)
	mainAxisownerSize := If(isMainAxisRow, ownerWidth, ownerHeight)

	var childWidth, childHeight float32
	var childWidthMeasureMode, childHeightMeasureMode MeasureMode

	resolvedFlexBasis := resolveValue(child.resolveFlexBasisPtr(), mainAxisownerSize)
	isRowStyleDimDefined := styleDefinesDimension(child, FlexDirectionRow, ownerWidth)
	isColumnStyleDimDefined := styleDefinesDimension(child, FlexDirectionColumn, ownerHeight)

	if resolvedFlexBasis.isDefined() && IsDefined(mainAxisSize) {
		if child.getLayout().computedFlexBasis.isUndefined() ||
			(child.GetConfig().IsExperimentalFeatureEnabled(
				ExperimentalFeatureWebFlexBasis) &&
				child.getLayout().computedFlexBasisGeneration != generationCount) {
			paddingAndBorder :=
				NewFloatOptional(paddingAndBorderForAxis(child, mainAxis, ownerWidth))
			child.setLayoutComputedFlexBasis(
				FloatOptional{maxOrDefined(resolvedFlexBasis.unwrap(), paddingAndBorder.unwrap())})
		}
	} else if isMainAxisRow && isRowStyleDimDefined {
		// The width is definite, so use that as the flex basis.
		paddingAndBorder := NewFloatOptional(
			paddingAndBorderForAxis(child, FlexDirectionRow, ownerWidth))

		child.setLayoutComputedFlexBasis(FloatOptional{maxOrDefined(
			resolveValue(
				child.getResolvedDimension(DimensionWidth), ownerWidth).unwrap(),
			paddingAndBorder.unwrap())})
	} else if !isMainAxisRow && isColumnStyleDimDefined {
		// The height is definite, so use that as the flex basis.
		paddingAndBorder := NewFloatOptional(
			paddingAndBorderForAxis(child, FlexDirectionColumn, ownerWidth))
		child.setLayoutComputedFlexBasis(FloatOptional{maxOrDefined(
			resolveValue(
				child.getResolvedDimension(DimensionHeight), ownerHeight).unwrap(),
			paddingAndBorder.unwrap())})
	} else {
		// Compute the flex basis and hypothetical main size (i.e. the clamped flex
		// basis).
		childWidth = Undefined
		childHeight = Undefined
		childWidthMeasureMode = MeasureModeUndefined
		childHeightMeasureMode = MeasureModeUndefined

		marginRow := child.getMarginForAxis(FlexDirectionRow, ownerWidth)
		marginColumn :=
			child.getMarginForAxis(FlexDirectionColumn, ownerWidth)

		if isRowStyleDimDefined {
			childWidth =
				resolveValue(
					child.getResolvedDimension(DimensionWidth), ownerWidth).unwrap() + marginRow
			childWidthMeasureMode = MeasureModeExactly
		}
		if isColumnStyleDimDefined {
			childHeight =
				resolveValue(
					child.getResolvedDimension(DimensionHeight), ownerHeight).unwrap() + marginColumn
			childHeightMeasureMode = MeasureModeExactly
		}

		// The W3C spec doesn't say anything about the 'overflow' property, but all
		// major browsers appear to implement the following logic.
		if (!isMainAxisRow && node.getStyle().overflow() == OverflowScroll) ||
			node.getStyle().overflow() != OverflowScroll {
			if IsUndefined(childWidth) && IsDefined(width) {
				childWidth = width
				childWidthMeasureMode = MeasureModeAtMost
			}
		}

		if (isMainAxisRow && node.getStyle().overflow() == OverflowScroll) ||
			node.getStyle().overflow() != OverflowScroll {
			if IsUndefined(childHeight) && IsDefined(height) {
				childHeight = height
				childHeightMeasureMode = MeasureModeAtMost
			}
		}

		childStyle := child.getStyle()
		if childStyle.aspectRatio().isDefined() {
			if !isMainAxisRow && childWidthMeasureMode == MeasureModeExactly {
				childHeight = marginColumn +
					(childWidth-marginRow)/childStyle.aspectRatio().unwrap()
				childHeightMeasureMode = MeasureModeExactly
			} else if isMainAxisRow && childHeightMeasureMode == MeasureModeExactly {
				childWidth = marginRow +
					(childHeight-marginColumn)*childStyle.aspectRatio().unwrap()
				childWidthMeasureMode = MeasureModeExactly
			}
		}

		// If child has no defined size in the cross axis and is set to stretch, set
		// the cross axis to be measured exactly with the available inner width

		hasExactWidth :=
			IsDefined(width) && widthMode == MeasureModeExactly
		childWidthStretch :=
			resolveChildAlignment(node, child) == AlignStretch &&
				childWidthMeasureMode != MeasureModeExactly
		if !isMainAxisRow && !isRowStyleDimDefined && hasExactWidth &&
			childWidthStretch {
			childWidth = width
			childWidthMeasureMode = MeasureModeExactly
			if childStyle.aspectRatio().isDefined() {
				childHeight =
					(childWidth - marginRow) / childStyle.aspectRatio().unwrap()
				childHeightMeasureMode = MeasureModeExactly
			}
		}

		hasExactHeight :=
			IsDefined(height) && heightMode == MeasureModeExactly
		childHeightStretch :=
			resolveChildAlignment(node, child) == AlignStretch &&
				childHeightMeasureMode != MeasureModeExactly
		if isMainAxisRow && !isColumnStyleDimDefined && hasExactHeight &&
			childHeightStretch {
			childHeight = height
			childHeightMeasureMode = MeasureModeExactly

			if childStyle.aspectRatio().isDefined() {
				childWidth =
					(childHeight - marginColumn) * childStyle.aspectRatio().unwrap()
				childWidthMeasureMode = MeasureModeExactly
			}
		}

		constrainMaxSizeForMode(
			child,
			FlexDirectionRow,
			ownerWidth,
			ownerWidth,
			&childWidthMeasureMode,
			&childWidth)
		constrainMaxSizeForMode(
			child,
			FlexDirectionColumn,
			ownerHeight,
			ownerWidth,
			&childHeightMeasureMode,
			&childHeight)

		// Measure the child
		calculateLayoutInternal(
			child,
			childWidth,
			childHeight,
			direction,
			childWidthMeasureMode,
			childHeightMeasureMode,
			ownerWidth,
			ownerHeight,
			false,
			LayoutPassReasonMeasureChild,
			layoutMarkerData,
			depth,
			generationCount)

		child.setLayoutComputedFlexBasis(NewFloatOptional(maxOrDefined(
			child.getLayout().measuredDimension(dimension(mainAxis)),
			paddingAndBorderForAxis(child, mainAxis, ownerWidth))))
	}
	child.setLayoutComputedFlexBasisGeneration(generationCount)

}

func measureNodeWithMeasureFunc(
	node *Node,
	availableWidth float32,
	availableHeight float32,
	widthMeasureMode MeasureMode,
	heightMeasureMode MeasureMode,
	ownerWidth float32,
	ownerHeight float32,
	layoutMarkerData *LayoutData,
	reason LayoutPassReason,
) {
	if !node.HasMeasureFunc() {
		panic("Expected node to have custom measure function")
	}

	if widthMeasureMode == MeasureModeUndefined {
		availableWidth = Undefined
	}
	if heightMeasureMode == MeasureModeUndefined {
		availableHeight = Undefined
	}

	padding := node.getLayout().padding_
	border := node.getLayout().border_
	paddingAndBorderAxisRow := padding[EdgeLeft] +
		padding[EdgeRight] + border[EdgeLeft] + border[EdgeRight]
	paddingAndBorderAxisColumn := padding[EdgeTop] +
		padding[EdgeBottom] + border[EdgeTop] + border[EdgeBottom]

	innerWidth := If(IsUndefined(availableWidth), availableWidth, maxOrDefined(0, availableWidth-paddingAndBorderAxisRow))

	innerHeight := If(IsUndefined(availableHeight), availableHeight, maxOrDefined(0, availableHeight-paddingAndBorderAxisColumn))

	if widthMeasureMode == MeasureModeExactly && heightMeasureMode == MeasureModeExactly {
		node.setLayoutMeasuredDimension(
			boundAxis(
				node, FlexDirectionRow, availableWidth, ownerWidth, ownerWidth),
			DimensionWidth,
		)
		node.setLayoutMeasuredDimension(
			boundAxis(
				node, FlexDirectionColumn, availableHeight, ownerHeight, ownerWidth),
			DimensionHeight,
		)
	} else {
		// Event.Publish(Event.MeasureCallbackStart, node)

		measuredSize := node.measure(innerWidth, widthMeasureMode, innerHeight, heightMeasureMode)

		layoutMarkerData.measureCallbacks += 1
		layoutMarkerData.measureCallbackReasonsCount[reason] += 1

		//Event.Publish(Event.MeasureCallbackEnd, node, MeasureCallbackEndData{
		//     InnerWidth:        innerWidth,
		//     WidthMeasureMode:  unscopedEnum(widthMeasureMode),
		//     InnerHeight:       innerHeight,
		//     HeightMeasureMode: unscopedEnum(heightMeasureMode),
		//     MeasuredWidth:     measuredSize.Width,
		//     MeasuredHeight:    measuredSize.Height,
		//     Reason:            reason,
		// })

		node.setLayoutMeasuredDimension(
			boundAxis(
				node,
				FlexDirectionRow,
				If(widthMeasureMode == MeasureModeUndefined || widthMeasureMode == MeasureModeAtMost, measuredSize.Width+paddingAndBorderAxisRow, availableWidth),
				ownerWidth,
				ownerWidth,
			),
			DimensionWidth,
		)

		node.setLayoutMeasuredDimension(
			boundAxis(
				node,
				FlexDirectionColumn,
				If(heightMeasureMode == MeasureModeUndefined || heightMeasureMode == MeasureModeAtMost, measuredSize.Height+paddingAndBorderAxisColumn, availableHeight),
				ownerHeight,
				ownerWidth,
			),
			DimensionHeight,
		)
	}
}

// For nodes with no children, use the available values if they were provided,
// or the minimum size as indicated by the padding and border sizes.
func measureNodeWithoutChildren(
	node *Node,
	availableWidth float32,
	availableHeight float32,
	widthMeasureMode MeasureMode,
	heightMeasureMode MeasureMode,
	ownerWidth float32,
	ownerHeight float32,
) {
	padding := node.getLayout().padding_
	border := node.getLayout().border_

	width := availableWidth
	if widthMeasureMode == MeasureModeUndefined || widthMeasureMode == MeasureModeAtMost {
		width = padding[EdgeLeft] + padding[EdgeRight] + border[EdgeLeft] + border[EdgeRight]
	}
	node.setLayoutMeasuredDimension(
		boundAxis(node, FlexDirectionRow, width, ownerWidth, ownerWidth),
		DimensionWidth,
	)

	height := availableHeight
	if heightMeasureMode == MeasureModeUndefined || heightMeasureMode == MeasureModeAtMost {
		height = padding[EdgeTop] + padding[EdgeBottom] + border[EdgeTop] + border[EdgeBottom]
	}
	node.setLayoutMeasuredDimension(
		boundAxis(node, FlexDirectionColumn, height, ownerHeight, ownerWidth),
		DimensionHeight,
	)
}

func measureNodeWithFixedSize(
	node *Node,
	availableWidth float32,
	availableHeight float32,
	widthMeasureMode MeasureMode,
	heightMeasureMode MeasureMode,
	ownerWidth float32,
	ownerHeight float32,
) bool {
	if (IsDefined(availableWidth) &&
		widthMeasureMode == MeasureModeAtMost && availableWidth <= 0.0) ||
		(IsDefined(availableHeight) &&
			heightMeasureMode == MeasureModeAtMost && availableHeight <= 0.0) ||
		(widthMeasureMode == MeasureModeExactly &&
			heightMeasureMode == MeasureModeExactly) {
		node.setLayoutMeasuredDimension(
			boundAxis(
				node,
				FlexDirectionRow,
				If(IsUndefined(availableWidth) || (widthMeasureMode == MeasureModeAtMost && availableWidth < 0.0), 0.0, availableWidth),
				ownerWidth,
				ownerWidth,
			),
			DimensionWidth,
		)

		node.setLayoutMeasuredDimension(
			boundAxis(
				node,
				FlexDirectionColumn,
				If(IsUndefined(availableHeight) || (heightMeasureMode == MeasureModeAtMost && availableHeight < 0.0), 0.0, availableHeight),
				ownerHeight,
				ownerWidth,
			),
			DimensionHeight,
		)
		return true
	}

	return false
}

func zeroOutLayoutRecursively(node *Node) {
	node.setLayout(LayoutResults{})
	node.setLayoutDimension(0, DimensionWidth)
	node.setLayoutDimension(0, DimensionHeight)
	node.SetHasNewLayout(true)

	node.cloneChildrenIfNeeded()
	for _, child := range node.GetChildren() {
		zeroOutLayoutRecursively(child)
	}
}

func layoutAbsoluteChild(
	node *Node,
	child *Node,
	width float32,
	widthMode MeasureMode,
	height float32,
	direction Direction,
	layoutMarkerData *LayoutData,
	depth uint32,
	generationCount uint32,
) {
	if gDebuging {
		atomic.AddUint32(&gCurrentDebugCount, 1)
		fmt.Printf("layoutAbsoluteChild %d\n", atomic.LoadUint32(&gCurrentDebugCount))
	}
	mainAxis := resolveDirection(node.getStyle().flexDirection(), direction)
	crossAxis := resolveCrossDirection(mainAxis, direction)
	isMainAxisRow := isRow(mainAxis)

	var childWidth, childHeight float32 = Undefined, Undefined
	var childWidthMeasureMode, childHeightMeasureMode MeasureMode = MeasureModeUndefined, MeasureModeUndefined

	marginRow := child.getMarginForAxis(FlexDirectionRow, width)
	marginColumn := child.getMarginForAxis(FlexDirectionColumn, width)

	if styleDefinesDimension(child, FlexDirectionRow, width) {
		childWidth = resolveValue(child.getResolvedDimension(DimensionWidth), width).unwrap() + marginRow
	} else {
		if child.isInlineStartPositionDefined(FlexDirectionRow, direction) &&
			child.isInlineEndPositionDefined(FlexDirectionRow, direction) {
			childWidth = node.getLayout().measuredDimension(DimensionWidth) -
				(node.getInlineStartBorder(FlexDirectionRow, direction) +
					node.getInlineEndBorder(FlexDirectionRow, direction)) -
				(child.getInlineStartPosition(FlexDirectionRow, direction, width) +
					child.getInlineEndPosition(FlexDirectionRow, direction, width))
			childWidth = boundAxis(child, FlexDirectionRow, childWidth, width, width)
		}
	}

	if styleDefinesDimension(child, FlexDirectionColumn, height) {
		childHeight = resolveValue(child.getResolvedDimension(DimensionHeight), height).unwrap() + marginColumn
	} else {
		if child.isInlineStartPositionDefined(FlexDirectionColumn, direction) &&
			child.isInlineEndPositionDefined(FlexDirectionColumn, direction) {
			childHeight = node.getLayout().measuredDimension(DimensionHeight) -
				(node.getInlineStartBorder(FlexDirectionColumn, direction) +
					node.getInlineEndBorder(FlexDirectionColumn, direction)) -
				(child.getInlineStartPosition(FlexDirectionColumn, direction, height) +
					child.getInlineEndPosition(FlexDirectionColumn, direction, height))
			childHeight = boundAxis(child, FlexDirectionColumn, childHeight, height, width)
		}
	}

	childStyle := child.getStyle()
	if (IsUndefined(childWidth) != IsUndefined(childHeight)) && childStyle.aspectRatio().isDefined() {
		if IsUndefined(childWidth) {
			childWidth = marginRow + (childHeight-marginColumn)*childStyle.aspectRatio().unwrap()
		} else if IsUndefined(childHeight) {
			childHeight = marginColumn + (childWidth-marginRow)/childStyle.aspectRatio().unwrap()
		}
	}

	if IsUndefined(childWidth) || IsUndefined(childHeight) {
		if IsUndefined(childWidth) {
			childWidthMeasureMode = MeasureModeUndefined
		} else {
			childWidthMeasureMode = MeasureModeExactly
		}
		if IsUndefined(childHeight) {
			childHeightMeasureMode = MeasureModeUndefined
		} else {
			childHeightMeasureMode = MeasureModeExactly
		}

		if !isMainAxisRow && IsUndefined(childWidth) &&
			widthMode != MeasureModeUndefined && IsDefined(width) && width > 0 {
			childWidth = width
			childWidthMeasureMode = MeasureModeAtMost
		}

		calculateLayoutInternal(
			child,
			childWidth,
			childHeight,
			direction,
			childWidthMeasureMode,
			childHeightMeasureMode,
			childWidth,
			childHeight,
			false,
			LayoutPassReasonAbsMeasureChild,
			layoutMarkerData,
			depth,
			generationCount,
		)
		childWidth = child.getLayout().measuredDimension(DimensionWidth) +
			child.getMarginForAxis(FlexDirectionRow, width)
		childHeight = child.getLayout().measuredDimension(DimensionHeight) +
			child.getMarginForAxis(FlexDirectionColumn, width)
	}

	calculateLayoutInternal(
		child,
		childWidth,
		childHeight,
		direction,
		MeasureModeExactly,
		MeasureModeExactly,
		childWidth,
		childHeight,
		true,
		LayoutPassReasonAbsLayout,
		layoutMarkerData,
		depth,
		generationCount,
	)

	if child.isFlexEndPositionDefined(mainAxis) && !child.isFlexStartPositionDefined(mainAxis) {
		child.setLayoutPosition(
			node.getLayout().measuredDimension(dimension(mainAxis))-
				child.getLayout().measuredDimension(dimension(mainAxis))-
				node.getFlexEndBorder(mainAxis, direction)-
				child.getFlexEndMargin(mainAxis, If(isMainAxisRow, width, height))-
				child.getFlexEndPosition(mainAxis, If(isMainAxisRow, width, height)),
			flexStartEdge(mainAxis),
		)
	} else if !child.isFlexStartPositionDefined(mainAxis) && node.getStyle().justifyContent() == JustifyCenter {
		child.setLayoutPosition(
			(node.getLayout().measuredDimension(dimension(mainAxis))-
				child.getLayout().measuredDimension(dimension(mainAxis)))/2.0,
			flexStartEdge(mainAxis),
		)
	} else if !child.isFlexStartPositionDefined(mainAxis) && node.getStyle().justifyContent() == JustifyFlexEnd {
		child.setLayoutPosition(
			node.getLayout().measuredDimension(dimension(mainAxis))-
				child.getLayout().measuredDimension(dimension(mainAxis)),
			flexStartEdge(mainAxis),
		)
	} else if node.GetConfig().IsExperimentalFeatureEnabled(ExperimentalFeatureAbsolutePercentageAgainstPaddingEdge) &&
		child.isFlexStartPositionDefined(mainAxis) {
		child.setLayoutPosition(
			child.getFlexStartPosition(
				mainAxis,
				node.getLayout().measuredDimension(dimension(mainAxis)),
			)+
				node.getFlexStartBorder(mainAxis, direction)+
				child.getFlexStartMargin(mainAxis, node.getLayout().measuredDimension(dimension(mainAxis))),
			flexStartEdge(mainAxis),
		)
	}

	if child.isFlexEndPositionDefined(crossAxis) && !child.isFlexStartPositionDefined(crossAxis) {
		child.setLayoutPosition(
			node.getLayout().measuredDimension(dimension(crossAxis))-
				child.getLayout().measuredDimension(dimension(crossAxis))-
				node.getFlexEndBorder(crossAxis, direction)-
				child.getFlexEndMargin(crossAxis, If(isMainAxisRow, height, width))-
				child.getFlexEndPosition(crossAxis, If(isMainAxisRow, height, width)),
			flexStartEdge(crossAxis),
		)
	} else if !child.isFlexStartPositionDefined(crossAxis) && resolveChildAlignment(node, child) == AlignCenter {
		child.setLayoutPosition(
			(node.getLayout().measuredDimension(dimension(crossAxis))-
				child.getLayout().measuredDimension(dimension(crossAxis)))/2.0,
			flexStartEdge(crossAxis),
		)
	} else if !child.isFlexStartPositionDefined(crossAxis) &&
		((resolveChildAlignment(node, child) == AlignFlexEnd) !=
			(node.getStyle().flexWrap() == WrapWrapReverse)) {
		child.setLayoutPosition(
			node.getLayout().measuredDimension(dimension(crossAxis))-
				child.getLayout().measuredDimension(dimension(crossAxis)),
			flexStartEdge(crossAxis),
		)
	} else if node.GetConfig().IsExperimentalFeatureEnabled(ExperimentalFeatureAbsolutePercentageAgainstPaddingEdge) &&
		child.isFlexStartPositionDefined(crossAxis) {
		child.setLayoutPosition(
			child.getFlexStartPosition(
				crossAxis,
				node.getLayout().measuredDimension(dimension(crossAxis)),
			)+
				node.getFlexStartBorder(crossAxis, direction)+
				child.getFlexStartMargin(crossAxis, node.getLayout().measuredDimension(dimension(crossAxis))),
			flexStartEdge(crossAxis),
		)
	}
}

func calculateAvailableInnerDimension(
	node *Node,
	dimension Dimension,
	availableDim float32,
	paddingAndBorder float32,
	ownerDim float32,
) float32 {
	availableInnerDim := availableDim - paddingAndBorder
	// Max dimension overrides predefined dimension value; Min dimension in turn
	// overrides both of the above
	if IsDefined(availableInnerDim) {
		// We want to make sure our available height does not violate min and max
		// constraints
		minDimensionOptional := resolveValue(node.getStyle().minDimension(dimension).Value(), ownerDim)
		minInnerDim := If(minDimensionOptional.isUndefined(), 0.0, minDimensionOptional.unwrap()-paddingAndBorder)

		maxDimensionOptional := resolveValue(node.getStyle().maxDimension(dimension).Value(), ownerDim)

		maxInnerDim := If(maxDimensionOptional.isUndefined(), math.MaxFloat32, maxDimensionOptional.unwrap()-paddingAndBorder)

		availableInnerDim = maxOrDefined(minOrDefined(availableInnerDim, maxInnerDim), minInnerDim)
	}

	return availableInnerDim
}

func computeFlexBasisForChildren(
	node *Node,
	availableInnerWidth float32,
	availableInnerHeight float32,
	widthMeasureMode MeasureMode,
	heightMeasureMode MeasureMode,
	direction Direction,
	mainAxis FlexDirection,
	performLayout bool,
	layoutMarkerData *LayoutData,
	depth uint32,
	generationCount uint32,
) float32 {
	totalOuterFlexBasis := float32(0.0)
	var singleFlexChild *Node
	children := node.GetChildren()
	measureModeMainDim := If(isRow(mainAxis), widthMeasureMode, heightMeasureMode)

	if measureModeMainDim == MeasureModeExactly {
		for _, child := range children {
			if child.isNodeFlexible() {
				if singleFlexChild != nil ||
					inexactEqual(child.resolveFlexGrow(), 0.0) ||
					inexactEqual(child.resolveFlexShrink(), 0.0) {
					singleFlexChild = nil
					break
				} else {
					singleFlexChild = child
				}
			}
		}
	}

	for _, child := range children {
		child.resolveDimension()
		if child.getStyle().display() == DisplayNone {
			zeroOutLayoutRecursively(child)
			child.SetHasNewLayout(true)
			child.setDirty(false)
			continue
		}
		if performLayout {
			childDirection := child.resolveDirection(direction)
			mainDim := If(isRow(mainAxis), availableInnerWidth, availableInnerHeight)
			crossDim := If(isRow(mainAxis), availableInnerHeight, availableInnerWidth)
			child.setPosition(childDirection, mainDim, crossDim, availableInnerWidth)
		}

		if child.getStyle().positionType() == PositionTypeAbsolute {
			continue
		}
		if child == singleFlexChild {
			child.setLayoutComputedFlexBasisGeneration(generationCount)
			child.setLayoutComputedFlexBasis(NewFloatOptional(0))
		} else {
			computeFlexBasisForChild(
				node,
				child,
				availableInnerWidth,
				widthMeasureMode,
				availableInnerHeight,
				availableInnerWidth,
				availableInnerHeight,
				heightMeasureMode,
				direction,
				layoutMarkerData,
				depth,
				generationCount,
			)
		}

		totalOuterFlexBasis += (child.getLayout().computedFlexBasis.unwrap() +
			child.getMarginForAxis(mainAxis, availableInnerWidth))
	}

	return totalOuterFlexBasis
}

// It distributes the free space to the flexible items and ensures that the size
// of the flex items abide the min and max constraints. At the end of this
// function the child nodes would have proper size. Prior using this function
// please ensure that distributeFreeSpaceFirstPass is called.
func distributeFreeSpaceSecondPass(
	flexLine *FlexLine,
	node *Node,
	mainAxis FlexDirection,
	crossAxis FlexDirection,
	mainAxisOwnerSize float32,
	availableInnerMainDim float32,
	availableInnerCrossDim float32,
	availableInnerWidth float32,
	availableInnerHeight float32,
	mainAxisOverflows bool,
	measureModeCrossDim MeasureMode,
	performLayout bool,
	layoutMarkerData *LayoutData,
	depth uint32,
	generationCount uint32,
) float32 {
	var childFlexBasis float32
	var flexShrinkScaledFactor float32
	var flexGrowFactor float32
	var deltaFreeSpace float32
	isMainAxisRow := isRow(mainAxis)
	isNodeFlexWrap := node.getStyle().flexWrap() != WrapNoWrap

	for _, currentLineChild := range flexLine.itemsInFlow {
		childFlexBasis = boundAxisWithinMinAndMax(
			currentLineChild,
			mainAxis,
			currentLineChild.getLayout().computedFlexBasis,
			mainAxisOwnerSize,
		).unwrap()
		updatedMainSize := childFlexBasis

		if IsDefined(flexLine.layout.remainingFreeSpace) && flexLine.layout.remainingFreeSpace < 0 {
			flexShrinkScaledFactor = -currentLineChild.resolveFlexShrink() * childFlexBasis
			if flexShrinkScaledFactor != 0 {
				var childSize float32

				if IsDefined(flexLine.layout.totalFlexShrinkScaledFactors) && flexLine.layout.totalFlexShrinkScaledFactors == 0 {
					childSize = childFlexBasis + flexShrinkScaledFactor
				} else {
					childSize = childFlexBasis + (flexLine.layout.remainingFreeSpace/flexLine.layout.totalFlexShrinkScaledFactors)*flexShrinkScaledFactor
				}

				updatedMainSize = boundAxis(
					currentLineChild,
					mainAxis,
					childSize,
					availableInnerMainDim,
					availableInnerWidth,
				)
			}
		} else if IsDefined(flexLine.layout.remainingFreeSpace) && flexLine.layout.remainingFreeSpace > 0 {
			flexGrowFactor = currentLineChild.resolveFlexGrow()
			if !IsNaN(flexGrowFactor) && flexGrowFactor != 0 {
				updatedMainSize = boundAxis(
					currentLineChild,
					mainAxis,
					childFlexBasis+(flexLine.layout.remainingFreeSpace/flexLine.layout.totalFlexGrowFactors)*flexGrowFactor,
					availableInnerMainDim,
					availableInnerWidth,
				)
			}
		}

		deltaFreeSpace += updatedMainSize - childFlexBasis

		marginMain := currentLineChild.getMarginForAxis(mainAxis, availableInnerWidth)
		marginCross := currentLineChild.getMarginForAxis(crossAxis, availableInnerWidth)

		var childCrossSize float32
		childMainSize := updatedMainSize + marginMain
		var childCrossMeasureMode MeasureMode
		childMainMeasureMode := MeasureModeExactly

		childStyle := currentLineChild.getStyle()
		if childStyle.aspectRatio().isDefined() {
			if isMainAxisRow {
				childCrossSize = (childMainSize - marginMain) / childStyle.aspectRatio().unwrap()
			} else {
				childCrossSize = (childMainSize - marginMain) * childStyle.aspectRatio().unwrap()
			}
			childCrossMeasureMode = MeasureModeExactly

			childCrossSize += marginCross
		} else if !IsNaN(availableInnerCrossDim) &&
			!styleDefinesDimension(currentLineChild, crossAxis, availableInnerCrossDim) &&
			measureModeCrossDim == MeasureModeExactly &&
			!(isNodeFlexWrap && mainAxisOverflows) &&
			resolveChildAlignment(node, currentLineChild) == AlignStretch &&
			currentLineChild.getFlexStartMarginValue(crossAxis).unit != UnitAuto &&
			currentLineChild.marginTrailingValue(crossAxis).unit != UnitAuto {
			childCrossSize = availableInnerCrossDim
			childCrossMeasureMode = MeasureModeExactly
		} else if !styleDefinesDimension(currentLineChild, crossAxis, availableInnerCrossDim) {
			childCrossSize = availableInnerCrossDim
			childCrossMeasureMode = If(IsUndefined(childCrossSize), MeasureModeUndefined, MeasureModeAtMost)
		} else {
			childCrossSize = resolveValue(
				currentLineChild.getResolvedDimension(dimension(crossAxis)),
				availableInnerCrossDim,
			).unwrap() + marginCross
			isLoosePercentageMeasurement :=
				currentLineChild.getResolvedDimension(dimension(crossAxis)).unit == UnitPercent &&
					measureModeCrossDim != MeasureModeExactly
			if IsUndefined(childCrossSize) || isLoosePercentageMeasurement {
				childCrossMeasureMode = MeasureModeUndefined
			} else {
				childCrossMeasureMode = MeasureModeExactly
			}
		}

		constrainMaxSizeForMode(
			currentLineChild,
			mainAxis,
			availableInnerMainDim,
			availableInnerWidth,
			&childMainMeasureMode,
			&childMainSize,
		)
		constrainMaxSizeForMode(
			currentLineChild,
			crossAxis,
			availableInnerCrossDim,
			availableInnerWidth,
			&childCrossMeasureMode,
			&childCrossSize,
		)

		requiresStretchLayout :=
			!styleDefinesDimension(currentLineChild, crossAxis, availableInnerCrossDim) &&
				resolveChildAlignment(node, currentLineChild) == AlignStretch &&
				currentLineChild.getFlexStartMarginValue(crossAxis).unit != UnitAuto &&
				currentLineChild.marginTrailingValue(crossAxis).unit != UnitAuto

		childWidth := childMainSize
		childHeight := childCrossSize
		childWidthMeasureMode := childMainMeasureMode
		childHeightMeasureMode := childCrossMeasureMode
		if !isMainAxisRow {
			childWidth = childCrossSize
			childHeight = childMainSize
			childWidthMeasureMode = childCrossMeasureMode
			childHeightMeasureMode = childMainMeasureMode
		}

		isLayoutPass := performLayout && !requiresStretchLayout
		calculateLayoutInternal(
			currentLineChild,
			childWidth,
			childHeight,
			node.getLayout().direction(),
			childWidthMeasureMode,
			childHeightMeasureMode,
			availableInnerWidth,
			availableInnerHeight,
			isLayoutPass,
			If(isLayoutPass, LayoutPassReasonFlexLayout, LayoutPassReasonFlexMeasure),
			layoutMarkerData,
			depth,
			generationCount,
		)
		node.setLayoutHadOverflow(
			node.getLayout().hadOverflow() ||
				currentLineChild.getLayout().hadOverflow(),
		)
	}
	return deltaFreeSpace
}

// It distributes the free space to the flexible items.For those flexible items
// whose min and max constraints are triggered, those flex item's clamped size
// is removed from the remaingfreespace.
func distributeFreeSpaceFirstPass(
	flexLine *FlexLine,
	mainAxis FlexDirection,
	mainAxisOwnerSize float32,
	availableInnerMainDim float32,
	availableInnerWidth float32,
) {
	flexShrinkScaledFactor := float32(0)
	flexGrowFactor := float32(0)
	baseMainSize := float32(0)
	boundMainSize := float32(0)
	deltaFreeSpace := float32(0)

	for _, currentLineChild := range flexLine.itemsInFlow {
		childFlexBasis := boundAxisWithinMinAndMax(
			currentLineChild,
			mainAxis,
			currentLineChild.getLayout().computedFlexBasis,
			mainAxisOwnerSize,
		).unwrap()

		if flexLine.layout.remainingFreeSpace < 0 {
			flexShrinkScaledFactor =
				-currentLineChild.resolveFlexShrink() * childFlexBasis

			// Is this child able to shrink?
			if IsDefined(flexShrinkScaledFactor) &&
				flexShrinkScaledFactor != 0 {
				baseMainSize = childFlexBasis +
					flexLine.layout.remainingFreeSpace/
						flexLine.layout.totalFlexShrinkScaledFactors*
						flexShrinkScaledFactor
				boundMainSize = boundAxis(
					currentLineChild,
					mainAxis,
					baseMainSize,
					availableInnerMainDim,
					availableInnerWidth,
				)
				if IsDefined(baseMainSize) && IsDefined(boundMainSize) &&
					baseMainSize != boundMainSize {
					// By excluding this item's size and flex factor from remaining, this
					// item's min/max constraints should also trigger in the second pass
					// resulting in the item's size calculation being identical in the
					// first and second passes.
					deltaFreeSpace += boundMainSize - childFlexBasis
					flexLine.layout.totalFlexShrinkScaledFactors -=
						(-currentLineChild.resolveFlexShrink() *
							currentLineChild.getLayout().computedFlexBasis.unwrap())
				}
			}
		} else if IsDefined(flexLine.layout.remainingFreeSpace) &&
			flexLine.layout.remainingFreeSpace > 0 {
			flexGrowFactor = currentLineChild.resolveFlexGrow()

			// Is this child able to grow?
			if IsDefined(flexGrowFactor) && flexGrowFactor != 0 {
				baseMainSize = childFlexBasis +
					flexLine.layout.remainingFreeSpace/
						flexLine.layout.totalFlexGrowFactors*flexGrowFactor
				boundMainSize = boundAxis(
					currentLineChild,
					mainAxis,
					baseMainSize,
					availableInnerMainDim,
					availableInnerWidth,
				)

				if IsDefined(baseMainSize) && IsDefined(boundMainSize) &&
					baseMainSize != boundMainSize {
					// By excluding this item's size and flex factor from remaining, this
					// item's min/max constraints should also trigger in the second pass
					// resulting in the item's size calculation being identical in the
					// first and second passes.
					deltaFreeSpace += boundMainSize - childFlexBasis
					flexLine.layout.totalFlexGrowFactors -= flexGrowFactor
				}
			}
		}
	}
	flexLine.layout.remainingFreeSpace -= deltaFreeSpace
}

// Do two passes over the flex items to figure out how to distribute the
// remaining space.
//
// The first pass finds the items whose min/max constraints trigger, freezes
// them at those sizes, and excludes those sizes from the remaining space.
//
// The second pass sets the size of each flexible item. It distributes the
// remaining space amongst the items whose min/max constraints didn't trigger in
// the first pass. For the other items, it sets their sizes by forcing their
// min/max constraints to trigger again.
//
// This two pass approach for resolving min/max constraints deviates from the
// spec. The spec
// (https://www.w3.org/TR/CSS-flexbox-1/#resolve-flexible-lengths) describes a
// process that needs to be repeated a variable number of times. The algorithm
// implemented here won't handle all cases but it was simpler to implement and
// it mitigates performance concerns because we know exactly how many passes
// it'll do.
//
// At the end of this function the child nodes would have the proper size
// assigned to them.
func resolveFlexibleLength(
	node *Node,
	flexLine *FlexLine,
	mainAxis FlexDirection,
	crossAxis FlexDirection,
	mainAxisOwnerSize float32,
	availableInnerMainDim float32,
	availableInnerCrossDim float32,
	availableInnerWidth float32,
	availableInnerHeight float32,
	mainAxisOverflows bool,
	measureModeCrossDim MeasureMode,
	performLayout bool,
	layoutMarkerData *LayoutData,
	depth uint32,
	generationCount uint32,
) {
	originalFreeSpace := flexLine.layout.remainingFreeSpace
	// First pass: detect the flex items whose min/max constraints trigger
	distributeFreeSpaceFirstPass(
		flexLine,
		mainAxis,
		mainAxisOwnerSize,
		availableInnerMainDim,
		availableInnerWidth,
	)

	// Second pass: resolve the sizes of the flexible items
	distributedFreeSpace := distributeFreeSpaceSecondPass(
		flexLine,
		node,
		mainAxis,
		crossAxis,
		mainAxisOwnerSize,
		availableInnerMainDim,
		availableInnerCrossDim,
		availableInnerWidth,
		availableInnerHeight,
		mainAxisOverflows,
		measureModeCrossDim,
		performLayout,
		layoutMarkerData,
		depth,
		generationCount,
	)

	flexLine.layout.remainingFreeSpace = originalFreeSpace - distributedFreeSpace
}

func justifyMainAxis(
	node *Node,
	flexLine *FlexLine,
	startOfLineIndex uint32,
	mainAxis FlexDirection,
	crossAxis FlexDirection,
	direction Direction,
	measureModeMainDim MeasureMode,
	measureModeCrossDim MeasureMode,
	mainAxisOwnerSize float32,
	ownerWidth float32,
	availableInnerMainDim float32,
	availableInnerCrossDim float32,
	availableInnerWidth float32,
	performLayout bool) {

	style := node.getStyle()
	leadingPaddingAndBorderMain := If(node.hasErrata(ErrataStartingEndingEdgeFromFlexDirection), node.getInlineStartPaddingAndBorder(mainAxis, direction, ownerWidth), node.getFlexStartPaddingAndBorder(mainAxis, direction, ownerWidth))
	trailingPaddingAndBorderMain := If(node.hasErrata(ErrataStartingEndingEdgeFromFlexDirection), node.getInlineEndPaddingAndBorder(mainAxis, direction, ownerWidth), node.getFlexEndPaddingAndBorder(mainAxis, direction, ownerWidth))

	gap := node.getGapForAxis(mainAxis)

	if measureModeMainDim == MeasureModeAtMost && flexLine.layout.remainingFreeSpace > 0 {
		if style.minDimension(dimension(mainAxis)).IsDefined() &&
			resolveValue(style.minDimension(dimension(mainAxis)).Value(), mainAxisOwnerSize).isDefined() {
			// This condition makes sure that if the size of main dimension(after
			// considering child nodes main dim, leading and trailing padding etc)
			// falls below min dimension, then the remainingFreeSpace is reassigned
			// considering the min dimension

			// `minAvailableMainDim` denotes minimum available space in which child
			// can be laid out, it will exclude space consumed by padding and border.
			minAvailableMainDim := resolveValue(style.minDimension(dimension(mainAxis)).Value(), mainAxisOwnerSize).unwrap() -
				leadingPaddingAndBorderMain - trailingPaddingAndBorderMain
			occupiedSpaceByChildNodes := availableInnerMainDim - flexLine.layout.remainingFreeSpace
			flexLine.layout.remainingFreeSpace = maxOrDefined(0.0, minAvailableMainDim-occupiedSpaceByChildNodes)
		} else {
			flexLine.layout.remainingFreeSpace = 0
		}
	}

	numberOfAutoMarginsOnCurrentLine := 0
	for i := startOfLineIndex; i < flexLine.endOfLineIndex; i++ {
		child := node.GetChild(i)
		if child.getStyle().positionType() != PositionTypeAbsolute {
			if child.getFlexStartMarginValue(mainAxis).unit == UnitAuto {
				numberOfAutoMarginsOnCurrentLine++
			}
			if child.marginTrailingValue(mainAxis).unit == UnitAuto {
				numberOfAutoMarginsOnCurrentLine++
			}
		}
	}
	// In order to position the elements in the main axis, we have two controls.
	// The space between the beginning and the first element and the space between
	// each two elements.
	leadingMainDim := float32(0)
	betweenMainDim := gap
	justifyContent := style.justifyContent()

	if numberOfAutoMarginsOnCurrentLine == 0 {
		switch justifyContent {
		case JustifyCenter:
			leadingMainDim = flexLine.layout.remainingFreeSpace / 2
		case JustifyFlexEnd:
			leadingMainDim = flexLine.layout.remainingFreeSpace
		case JustifySpaceBetween:
			if len(flexLine.itemsInFlow) > 1 {
				betweenMainDim += maxOrDefined(flexLine.layout.remainingFreeSpace, 0.0) /
					float32(len(flexLine.itemsInFlow)-1)
			}
		case JustifySpaceEvenly:
			leadingMainDim = flexLine.layout.remainingFreeSpace /
				float32(len(flexLine.itemsInFlow)+1)
			betweenMainDim += leadingMainDim
		case JustifySpaceAround:
			leadingMainDim = 0.5 * flexLine.layout.remainingFreeSpace /
				float32(len(flexLine.itemsInFlow))
			betweenMainDim += leadingMainDim * 2
		case JustifyFlexStart:
			break
		}
	}

	flexLine.layout.mainDim = leadingPaddingAndBorderMain + leadingMainDim
	flexLine.layout.crossDim = 0

	maxAscentForCurrentLine := float32(0)
	maxDescentForCurrentLine := float32(0)
	isNodeBaselineLayout := isBaselineLayout(node)

	for i := startOfLineIndex; i < flexLine.endOfLineIndex; i++ {
		child := node.GetChild(i)
		childStyle := child.getStyle()
		childLayout := child.getLayout()

		if childStyle.display() == DisplayNone {
			continue
		}

		if childStyle.positionType() == PositionTypeAbsolute &&
			child.isInlineStartPositionDefined(mainAxis, direction) {
			if performLayout {
				child.setLayoutPosition(
					child.getInlineStartPosition(mainAxis, direction, availableInnerMainDim)+
						node.getInlineStartBorder(mainAxis, direction)+
						child.getInlineStartMargin(mainAxis, direction, availableInnerWidth),
					flexStartEdge(mainAxis))
			}
		} else {
			if childStyle.positionType() != PositionTypeAbsolute {
				if child.getFlexStartMarginValue(mainAxis).unit == UnitAuto {
					flexLine.layout.mainDim += flexLine.layout.remainingFreeSpace /
						float32(numberOfAutoMarginsOnCurrentLine)
				}

				if performLayout {
					child.setLayoutPosition(
						childLayout.position(flexStartEdge(mainAxis))+
							flexLine.layout.mainDim,
						flexStartEdge(mainAxis))
				}

				if child != flexLine.itemsInFlow[len(flexLine.itemsInFlow)-1] {
					flexLine.layout.mainDim += betweenMainDim
				}

				if child.marginTrailingValue(mainAxis).unit == UnitAuto {
					flexLine.layout.mainDim += flexLine.layout.remainingFreeSpace /
						float32(numberOfAutoMarginsOnCurrentLine)
				}

				canSkipFlex := !performLayout && measureModeCrossDim == MeasureModeExactly
				if canSkipFlex {
					flexLine.layout.mainDim +=
						child.getMarginForAxis(mainAxis, availableInnerWidth) +
							childLayout.computedFlexBasis.unwrap()
					flexLine.layout.crossDim = availableInnerCrossDim
				} else {
					flexLine.layout.mainDim +=
						dimensionWithMargin(child, mainAxis, availableInnerWidth)

					if isNodeBaselineLayout {
						ascent := calculateBaseline(child) +
							child.getInlineStartMargin(FlexDirectionColumn, direction, availableInnerWidth)
						descent :=
							childLayout.measuredDimension(DimensionHeight) +
								child.getMarginForAxis(FlexDirectionColumn, availableInnerWidth) -
								ascent

						maxAscentForCurrentLine = maxOrDefined(maxAscentForCurrentLine, ascent)
						maxDescentForCurrentLine = maxOrDefined(maxDescentForCurrentLine, descent)
					} else {
						flexLine.layout.crossDim = maxOrDefined(
							flexLine.layout.crossDim,
							dimensionWithMargin(child, crossAxis, availableInnerWidth))
					}
				}
			} else if performLayout {
				child.setLayoutPosition(
					childLayout.position(flexStartEdge(mainAxis))+
						node.getInlineStartBorder(mainAxis, direction)+
						leadingMainDim,
					flexStartEdge(mainAxis))
			}
		}
	}

	flexLine.layout.mainDim += trailingPaddingAndBorderMain

	if isNodeBaselineLayout {
		flexLine.layout.crossDim =
			maxAscentForCurrentLine + maxDescentForCurrentLine
	}
}

// This is the main routine that implements a subset of the flexbox layout
// algorithm described in the W3C CSS documentation:
// https://www.w3.org/TR/CSS3-flexbox/.
//
// Limitations of this algorithm, compared to the full standard:
//   - Display property is always assumed to be 'flex' except for Text nodes,
//     which are assumed to be 'inline-flex'.
//   - The 'zIndex' property (or any form of z ordering) is not supported. Nodes
//     are stacked in document order.
//   - The 'order' property is not supported. The order of flex items is always
//     defined by document order.
//   - The 'visibility' property is always assumed to be 'visible'. Values of
//     'collapse' and 'hidden' are not supported.
//   - There is no support for forced breaks.
//   - It does not support vertical inline directions (top-to-bottom or
//     bottom-to-top text).
//
// Deviations from standard:
//   - Section 4.5 of the spec indicates that all flex items have a default
//     minimum main size. For text blocks, for example, this is the width of the
//     widest word. Calculating the minimum width is expensive, so we forego it
//     and assume a default minimum main size of 0.
//   - Min/Max sizes in the main axis are not honored when resolving flexible
//     lengths.
//   - The spec indicates that the default value for 'flexDirection' is 'row',
//     but the algorithm below assumes a default of 'column'.
//
// Input parameters:
//   - node: current node to be sized and laid out
//   - availableWidth & availableHeight: available size to be used for sizing
//     the node or YGUndefined if the size is not available; interpretation
//     depends on layout flags
//   - ownerDirection: the inline (text) direction within the owner
//     (left-to-right or right-to-left)
//   - widthMeasureMode: indicates the sizing rules for the width (see below
//     for explanation)
//   - heightMeasureMode: indicates the sizing rules for the height (see below
//     for explanation)
//   - performLayout: specifies whether the caller is interested in just the
//     dimensions of the node or it requires the entire node and its subtree to
//     be laid out (with final positions)
//
// Details:
//
//	This routine is called recursively to lay out subtrees of flexbox
//	elements. It uses the information in node.style, which is treated as a
//	read-only input. It is responsible for setting the layout.direction and
//	layout.measuredDimensions fields for the input node as well as the
//	layout.position and layout.lineIndex fields for its child nodes. The
//	layout.measuredDimensions field includes any border or padding for the
//	node but does not include margins.
//
//	The spec describes four different layout modes: "fill available", "max
//	content", "min content", and "fit content". Of these, we don't use "min
//	content" because we don't support default minimum main sizes (see above
//	for details). Each of our measure modes maps to a layout mode from the
//	spec (https://www.w3.org/TR/CSS3-sizing/#terms):
//	  - MeasureMode::Undefined: max content
//	  - MeasureMode::Exactly: fill available
//	  - MeasureMode::AtMost: fit content
//
//	When calling calculateLayoutImpl and calculateLayoutInternal, if the
//	caller passes an available size of undefined then it must also pass a
//	measure mode of MeasureMode::Undefined in that dimension.
func calculateLayoutImpl(
	node *Node,
	availableWidth float32,
	availableHeight float32,
	ownerDirection Direction,
	widthMeasureMode MeasureMode,
	heightMeasureMode MeasureMode,
	ownerWidth float32,
	ownerHeight float32,
	performLayout bool,
	layoutMarkerData *LayoutData,
	depth uint32,
	generationCount uint32,
	reason LayoutPassReason,
) {
	if gDebuging {
		atomic.AddUint32(&gCurrentDebugCount, 1)
		fmt.Printf("calculateLayoutImpl  %d\n", atomic.LoadUint32(&gCurrentDebugCount))
	}
	if !If(IsUndefined(availableWidth), widthMeasureMode == MeasureModeUndefined, true) {
		panic("availableWidth is indefinite so widthMeasureMode must be MeasureModeUndefined")
	}
	if !If(IsUndefined(availableHeight), heightMeasureMode == MeasureModeUndefined, true) {
		panic("availableHeight is indefinite so heightMeasureMode must be MeasureModeUndefined")
	}

	if performLayout {
		layoutMarkerData.layouts += 1
	} else {
		layoutMarkerData.measures += 1
	}

	direction := node.resolveDirection(ownerDirection)
	node.setLayoutDirection(direction)

	flexRowDirection := resolveDirection(FlexDirectionRow, direction)
	flexColumnDirection := resolveDirection(FlexDirectionColumn, direction)

	startEdge := EdgeLeft
	endEdge := EdgeRight
	if direction == DirectionRTL {
		startEdge = EdgeRight
		endEdge = EdgeLeft
	}

	marginRowLeading := node.getInlineStartMargin(flexRowDirection, direction, ownerWidth)
	node.setLayoutMargin(marginRowLeading, startEdge)
	marginRowTrailing := node.getInlineEndMargin(flexRowDirection, direction, ownerWidth)
	node.setLayoutMargin(marginRowTrailing, endEdge)
	marginColumnLeading := node.getInlineStartMargin(flexColumnDirection, direction, ownerWidth)
	node.setLayoutMargin(marginColumnLeading, EdgeTop)
	marginColumnTrailing := node.getInlineEndMargin(flexColumnDirection, direction, ownerWidth)
	node.setLayoutMargin(marginColumnTrailing, EdgeBottom)

	marginAxisRow := marginRowLeading + marginRowTrailing
	marginAxisColumn := marginColumnLeading + marginColumnTrailing

	node.setLayoutBorder(node.getInlineStartBorder(flexRowDirection, direction), startEdge)
	node.setLayoutBorder(node.getInlineEndBorder(flexRowDirection, direction), endEdge)
	node.setLayoutBorder(node.getInlineStartBorder(flexColumnDirection, direction), EdgeTop)
	node.setLayoutBorder(node.getInlineEndBorder(flexColumnDirection, direction), EdgeBottom)

	node.setLayoutPadding(node.getInlineStartPadding(flexRowDirection, direction, ownerWidth), startEdge)
	node.setLayoutPadding(node.getInlineEndPadding(flexRowDirection, direction, ownerWidth), endEdge)
	node.setLayoutPadding(node.getInlineStartPadding(flexColumnDirection, direction, ownerWidth), EdgeTop)
	node.setLayoutPadding(node.getInlineEndPadding(flexColumnDirection, direction, ownerWidth), EdgeBottom)

	if node.HasMeasureFunc() {
		measureNodeWithMeasureFunc(
			node,
			availableWidth-marginAxisRow,
			availableHeight-marginAxisColumn,
			widthMeasureMode,
			heightMeasureMode,
			ownerWidth,
			ownerHeight,
			layoutMarkerData,
			reason,
		)
		return
	}

	childCount := node.GetChildCount()
	if childCount == 0 {
		measureNodeWithoutChildren(
			node,
			availableWidth-marginAxisRow,
			availableHeight-marginAxisColumn,
			widthMeasureMode,
			heightMeasureMode,
			ownerWidth,
			ownerHeight,
		)
		return
	}

	if !performLayout && measureNodeWithFixedSize(
		node,
		availableWidth-marginAxisRow,
		availableHeight-marginAxisColumn,
		widthMeasureMode,
		heightMeasureMode,
		ownerWidth,
		ownerHeight,
	) {
		return
	}

	node.cloneChildrenIfNeeded()
	node.setLayoutHadOverflow(false)

	// STEP 1: CALCULATE VALUES FOR REMAINDER OF ALGORITHM
	mainAxis := resolveDirection(node.getStyle().flexDirection(), direction)
	crossAxis := resolveCrossDirection(mainAxis, direction)
	isMainAxisRow := isRow(mainAxis)
	isNodeFlexWrap := node.getStyle().flexWrap() != WrapNoWrap

	mainAxisownerSize := ownerHeight
	crossAxisownerSize := ownerWidth
	if isMainAxisRow {
		mainAxisownerSize = ownerWidth
		crossAxisownerSize = ownerHeight
	}

	paddingAndBorderAxisMain := paddingAndBorderForAxis(node, mainAxis, ownerWidth)
	paddingAndBorderAxisCross := paddingAndBorderForAxis(node, crossAxis, ownerWidth)
	leadingPaddingAndBorderCross := node.getInlineStartPaddingAndBorder(crossAxis, direction, ownerWidth)

	measureModeMainDim := heightMeasureMode
	measureModeCrossDim := widthMeasureMode
	if isMainAxisRow {
		measureModeMainDim = widthMeasureMode
		measureModeCrossDim = heightMeasureMode
	}

	paddingAndBorderAxisRow := paddingAndBorderAxisCross
	paddingAndBorderAxisColumn := paddingAndBorderAxisMain
	if isMainAxisRow {
		paddingAndBorderAxisRow = paddingAndBorderAxisMain
		paddingAndBorderAxisColumn = paddingAndBorderAxisCross
	}
	// STEP 2: DETERMINE AVAILABLE SIZE IN MAIN AND CROSS DIRECTIONS
	availableInnerWidth := calculateAvailableInnerDimension(
		node,
		DimensionWidth,
		availableWidth-marginAxisRow,
		paddingAndBorderAxisRow,
		ownerWidth,
	)
	availableInnerHeight := calculateAvailableInnerDimension(
		node,
		DimensionHeight,
		availableHeight-marginAxisColumn,
		paddingAndBorderAxisColumn,
		ownerHeight,
	)

	availableInnerMainDim := availableInnerHeight
	availableInnerCrossDim := availableInnerWidth
	if isMainAxisRow {
		availableInnerMainDim = availableInnerWidth
		availableInnerCrossDim = availableInnerHeight
	}

	// STEP 3: DETERMINE FLEX BASIS FOR EACH ITEM
	totalMainDim := float32(0)
	totalMainDim += computeFlexBasisForChildren(
		node,
		availableInnerWidth,
		availableInnerHeight,
		widthMeasureMode,
		heightMeasureMode,
		direction,
		mainAxis,
		performLayout,
		layoutMarkerData,
		depth,
		generationCount,
	)

	if childCount > 1 {
		totalMainDim += node.getGapForAxis(mainAxis) * float32(childCount-1)
	}

	mainAxisOverflows := (measureModeMainDim != MeasureModeUndefined) && totalMainDim > availableInnerMainDim

	if isNodeFlexWrap && mainAxisOverflows && measureModeMainDim == MeasureModeAtMost {
		measureModeMainDim = MeasureModeExactly
	}

	// STEP 4: COLLECT FLEX ITEMS INTO FLEX LINES
	startOfLineIndex := uint32(0)
	endOfLineIndex := uint32(0)
	lineCount := uint32(0)
	totalLineCrossDim := float32(0)
	crossAxisGap := node.getGapForAxis(crossAxis)
	maxLineMainDim := float32(0)

	for endOfLineIndex < childCount {
		flexLine := calculateFlexLine(
			node,
			ownerDirection,
			mainAxisownerSize,
			availableInnerWidth,
			availableInnerMainDim,
			startOfLineIndex,
			lineCount,
		)

		endOfLineIndex = flexLine.endOfLineIndex
		// If we don't need to measure the cross axis, we can skip the entire flex
		// step.
		canSkipFlex := !performLayout && measureModeCrossDim == MeasureModeExactly

		// STEP 5: RESOLVING FLEXIBLE LENGTHS ON MAIN AXIS
		// Calculate the remaining available space that needs to be allocated. If
		// the main dimension size isn't known, it is computed based on the line
		// length, so there's no more space left to distribute.
		sizeBasedOnContent := false
		if measureModeMainDim != MeasureModeExactly {
			style := node.getStyle()
			minInnerWidth := resolveValue(style.minDimension(DimensionWidth).Value(), ownerWidth).unwrap() - paddingAndBorderAxisRow
			maxInnerWidth := resolveValue(style.maxDimension(DimensionWidth).Value(), ownerWidth).unwrap() - paddingAndBorderAxisRow
			minInnerHeight := resolveValue(style.minDimension(DimensionHeight).Value(), ownerHeight).unwrap() - paddingAndBorderAxisColumn
			maxInnerHeight := resolveValue(style.maxDimension(DimensionHeight).Value(), ownerHeight).unwrap() - paddingAndBorderAxisColumn

			minInnerMainDim := minInnerHeight
			maxInnerMainDim := maxInnerHeight
			if isMainAxisRow {
				minInnerMainDim = minInnerWidth
				maxInnerMainDim = maxInnerWidth
			}

			if IsDefined(minInnerMainDim) && flexLine.sizeConsumed < minInnerMainDim {
				availableInnerMainDim = minInnerMainDim
			} else if IsDefined(maxInnerMainDim) && flexLine.sizeConsumed > maxInnerMainDim {
				availableInnerMainDim = maxInnerMainDim
			} else {
				useLegacyStretchBehaviour := node.hasErrata(ErrataStretchFlexBasis)

				if !useLegacyStretchBehaviour && ((IsDefined(flexLine.layout.totalFlexGrowFactors) && flexLine.layout.totalFlexGrowFactors == 0) || (IsDefined(node.resolveFlexGrow()) && node.resolveFlexGrow() == 0)) {
					availableInnerMainDim = flexLine.sizeConsumed
				}

				sizeBasedOnContent = !useLegacyStretchBehaviour
			}
		}

		if !sizeBasedOnContent && IsDefined(availableInnerMainDim) {
			flexLine.layout.remainingFreeSpace = availableInnerMainDim - flexLine.sizeConsumed
		} else if flexLine.sizeConsumed < 0 {
			flexLine.layout.remainingFreeSpace = -flexLine.sizeConsumed
		}

		if !canSkipFlex {
			resolveFlexibleLength(
				node,
				&flexLine,
				mainAxis,
				crossAxis,
				mainAxisownerSize,
				availableInnerMainDim,
				availableInnerCrossDim,
				availableInnerWidth,
				availableInnerHeight,
				mainAxisOverflows,
				measureModeCrossDim,
				performLayout,
				layoutMarkerData,
				depth,
				generationCount,
			)
		}

		node.setLayoutHadOverflow(node.getLayout().hadOverflow() || (flexLine.layout.remainingFreeSpace < 0))

		// STEP 6: MAIN-AXIS JUSTIFICATION & CROSS-AXIS SIZE DETERMINATION

		// At this point, all the children have their dimensions set in the main
		// axis. Their dimensions are also set in the cross axis with the exception
		// of items that are aligned "stretch". We need to compute these stretch
		// values and set the final positions.
		justifyMainAxis(
			node,
			&flexLine,
			startOfLineIndex,
			mainAxis,
			crossAxis,
			direction,
			measureModeMainDim,
			measureModeCrossDim,
			mainAxisownerSize,
			ownerWidth,
			availableInnerMainDim,
			availableInnerCrossDim,
			availableInnerWidth,
			performLayout,
		)

		containerCrossAxis := availableInnerCrossDim
		if measureModeCrossDim == MeasureModeUndefined || measureModeCrossDim == MeasureModeAtMost {
			containerCrossAxis = boundAxis(
				node,
				crossAxis,
				flexLine.layout.crossDim+paddingAndBorderAxisCross,
				crossAxisownerSize,
				ownerWidth,
			) - paddingAndBorderAxisCross
		}

		if !isNodeFlexWrap && measureModeCrossDim == MeasureModeExactly {
			flexLine.layout.crossDim = availableInnerCrossDim
		}

		flexLine.layout.crossDim = boundAxis(
			node,
			crossAxis,
			flexLine.layout.crossDim+paddingAndBorderAxisCross,
			crossAxisownerSize,
			ownerWidth,
		) - paddingAndBorderAxisCross
		// STEP 7: CROSS-AXIS ALIGNMENT
		// We can skip child alignment if we're just measuring the container.
		if performLayout {
			for i := startOfLineIndex; i < endOfLineIndex; i++ {
				child := node.GetChild(i)
				if child.getStyle().display() == DisplayNone {
					continue
				}
				if child.getStyle().positionType() == PositionTypeAbsolute {
					// If the child is absolutely positioned and has a top/left/bottom/right set, override all the previously computed positions to set it correctly.
					isChildLeadingPosDefined := child.isInlineStartPositionDefined(crossAxis, direction)
					if isChildLeadingPosDefined {
						child.setLayoutPosition(
							child.getInlineStartPosition(crossAxis, direction, availableInnerCrossDim)+
								node.getInlineStartBorder(crossAxis, direction)+
								child.getInlineStartMargin(crossAxis, direction, availableInnerWidth),
							flexStartEdge(crossAxis),
						)
					}
					// If leading position is not defined or calculations result in NaN, default to border + margin
					if !isChildLeadingPosDefined || IsUndefined(child.getLayout().position(flexStartEdge(crossAxis))) {
						child.setLayoutPosition(
							node.getInlineStartBorder(crossAxis, direction)+
								child.getInlineStartMargin(crossAxis, direction, availableInnerWidth),
							flexStartEdge(crossAxis),
						)
					}
				} else {
					leadingCrossDim := leadingPaddingAndBorderCross

					// For a relative children, we're either using alignItems (owner) or alignSelf (child) in order to determine the position in the cross axis
					alignItem := resolveChildAlignment(node, child)

					// If the child uses align stretch, we need to lay it out one more time, this time forcing the cross-axis size to be the computed cross size for the current line.
					if alignItem == AlignStretch &&
						child.getFlexStartMarginValue(crossAxis).unit != UnitAuto &&
						child.marginTrailingValue(crossAxis).unit != UnitAuto {
						// If the child defines a definite size for its cross axis, there's no need to stretch.
						if !styleDefinesDimension(child, crossAxis, availableInnerCrossDim) {
							childMainSize := child.getLayout().measuredDimension(dimension(mainAxis))
							childStyle := child.getStyle()
							childCrossSize := If(childStyle.aspectRatio().isDefined(), child.getMarginForAxis(crossAxis, availableInnerWidth)+If(isMainAxisRow, childMainSize/childStyle.aspectRatio().unwrap(), childMainSize*childStyle.aspectRatio().unwrap()), flexLine.layout.crossDim)

							childMainSize += child.getMarginForAxis(mainAxis, availableInnerWidth)

							childMainMeasureMode := MeasureModeExactly
							childCrossMeasureMode := MeasureModeExactly
							constrainMaxSizeForMode(
								child,
								mainAxis,
								availableInnerMainDim,
								availableInnerWidth,
								&childMainMeasureMode,
								&childMainSize,
							)
							constrainMaxSizeForMode(
								child,
								crossAxis,
								availableInnerCrossDim,
								availableInnerWidth,
								&childCrossMeasureMode,
								&childCrossSize,
							)

							childWidth := childCrossSize
							childHeight := childMainSize
							if isMainAxisRow {
								childWidth, childHeight = childMainSize, childCrossSize
							}

							alignContent := node.getStyle().alignContent()
							crossAxisDoesNotGrow := isNodeFlexWrap && alignContent != AlignStretch

							childWidthMeasureMode := If(IsUndefined(childWidth) || (!isMainAxisRow && crossAxisDoesNotGrow), MeasureModeUndefined, MeasureModeExactly)
							childHeightMeasureMode := If(IsUndefined(childHeight) || (isMainAxisRow && crossAxisDoesNotGrow), MeasureModeUndefined, MeasureModeExactly)

							calculateLayoutInternal(
								child,
								childWidth,
								childHeight,
								direction,
								childWidthMeasureMode,
								childHeightMeasureMode,
								availableInnerWidth,
								availableInnerHeight,
								true,
								LayoutPassReasonStretch,
								layoutMarkerData,
								depth,
								generationCount,
							)
						}
					} else {
						remainingCrossDim := containerCrossAxis - dimensionWithMargin(child, crossAxis, availableInnerWidth)

						if child.getFlexStartMarginValue(crossAxis).unit == UnitAuto &&
							child.marginTrailingValue(crossAxis).unit == UnitAuto {
							leadingCrossDim += maxOrDefined(0.0, remainingCrossDim/2)
						} else if child.marginTrailingValue(crossAxis).unit == UnitAuto {
							// No-Op
						} else if child.getFlexStartMarginValue(crossAxis).unit == UnitAuto {
							leadingCrossDim += maxOrDefined(0.0, remainingCrossDim)
						} else if alignItem == AlignFlexStart {
							// No-Op
						} else if alignItem == AlignCenter {
							leadingCrossDim += remainingCrossDim / 2
						} else {
							leadingCrossDim += remainingCrossDim
						}
					}
					// And we apply the position
					child.setLayoutPosition(
						child.getLayout().position(flexStartEdge(crossAxis))+totalLineCrossDim+leadingCrossDim,
						flexStartEdge(crossAxis),
					)
				}
			}
		}

		appliedCrossGap := crossAxisGap
		if lineCount == 0 {
			appliedCrossGap = 0.0
		}
		totalLineCrossDim += flexLine.layout.crossDim + appliedCrossGap
		maxLineMainDim = maxOrDefined(maxLineMainDim, flexLine.layout.mainDim)
		lineCount++
		startOfLineIndex = endOfLineIndex
	}

	// STEP 8: MULTI-LINE CONTENT ALIGNMENT
	// currentLead stores the size of the cross dim
	if performLayout && (isNodeFlexWrap || isBaselineLayout(node)) {
		crossDimLead := float32(0)
		currentLead := leadingPaddingAndBorderCross
		if IsDefined(availableInnerCrossDim) {
			remainingAlignContentDim := availableInnerCrossDim - totalLineCrossDim
			switch node.getStyle().alignContent() {
			case AlignFlexEnd:
				currentLead += remainingAlignContentDim
			case AlignCenter:
				currentLead += remainingAlignContentDim / 2
			case AlignStretch:
				if availableInnerCrossDim > totalLineCrossDim {
					crossDimLead = remainingAlignContentDim / float32(lineCount)
				}
			case AlignSpaceAround:
				if availableInnerCrossDim > totalLineCrossDim {
					currentLead += remainingAlignContentDim / (2 * float32(lineCount))
					if lineCount > 1 {
						crossDimLead = remainingAlignContentDim / float32(lineCount)
					}
				} else {
					currentLead += remainingAlignContentDim / 2
				}
			case AlignSpaceEvenly:
				if availableInnerCrossDim > totalLineCrossDim {
					currentLead += remainingAlignContentDim / float32(lineCount+1)
					if lineCount > 1 {
						crossDimLead = remainingAlignContentDim / float32(lineCount+1)
					}
				} else {
					currentLead += remainingAlignContentDim / 2
				}
			case AlignSpaceBetween:
				if availableInnerCrossDim > totalLineCrossDim && lineCount > 1 {
					crossDimLead = remainingAlignContentDim / float32(lineCount-1)
				}
			}
		}
		endIndex := uint32(0)
		for i := uint32(0); i < lineCount; i++ {
			startIndex := endIndex
			var ii uint32
			lineHeight := float32(0)
			maxAscentForCurrentLine := float32(0)
			maxDescentForCurrentLine := float32(0)
			for ii = startIndex; ii < childCount; ii++ {
				child := node.GetChild(ii)
				if child.getStyle().display() == DisplayNone {
					continue
				}
				if child.getStyle().positionType() != PositionTypeAbsolute {
					if child.getLineIndex() != i {
						break
					}
					if isLayoutDimensionDefined(child, crossAxis) {
						lineHeight = maxOrDefined(
							lineHeight,
							child.getLayout().measuredDimension(dimension(crossAxis))+
								child.getMarginForAxis(crossAxis, availableInnerWidth),
						)
					}
					if resolveChildAlignment(node, child) == AlignBaseline {
						ascent := calculateBaseline(child) +
							child.getInlineStartMargin(FlexDirectionColumn, direction, availableInnerWidth)
						descent := child.getLayout().measuredDimension(DimensionHeight) +
							child.getMarginForAxis(FlexDirectionColumn, availableInnerWidth) - ascent
						maxAscentForCurrentLine = maxOrDefined(maxAscentForCurrentLine, ascent)
						maxDescentForCurrentLine = maxOrDefined(maxDescentForCurrentLine, descent)
						lineHeight = maxOrDefined(lineHeight, maxAscentForCurrentLine+maxDescentForCurrentLine)
					}
				}
			}
			endIndex = ii
			lineHeight += crossDimLead
			currentLead += If(i != 0, crossAxisGap, 0.0)
			if performLayout {
				for ii = startIndex; ii < endIndex; ii++ {
					child := node.GetChild(ii)
					if child.getStyle().display() == DisplayNone {
						continue
					}
					if child.getStyle().positionType() != PositionTypeAbsolute {
						switch resolveChildAlignment(node, child) {
						case AlignFlexStart:
							child.setLayoutPosition(
								currentLead+
									child.getInlineStartMargin(crossAxis, direction, availableInnerWidth),
								flexStartEdge(crossAxis),
							)
						case AlignFlexEnd:
							child.setLayoutPosition(
								currentLead+lineHeight-
									child.getInlineEndMargin(crossAxis, direction, availableInnerWidth)-
									child.getLayout().measuredDimension(dimension(crossAxis)),
								flexStartEdge(crossAxis),
							)
						case AlignCenter:
							childHeight := child.getLayout().measuredDimension(dimension(crossAxis))
							child.setLayoutPosition(
								currentLead+(lineHeight-childHeight)/2,
								flexStartEdge(crossAxis),
							)
						case AlignStretch:
							child.setLayoutPosition(
								currentLead+
									child.getInlineStartMargin(crossAxis, direction, availableInnerWidth),
								flexStartEdge(crossAxis),
							)

							// Remeasure child with the line height as it as been only
							// measured with the owners height yet.
							if !styleDefinesDimension(child, crossAxis, availableInnerCrossDim) {
								childWidth := float32(0)
								childHeight := float32(0)
								if isMainAxisRow {
									childWidth = child.getLayout().measuredDimension(DimensionWidth) +
										child.getMarginForAxis(mainAxis, availableInnerWidth)
									childHeight = lineHeight
								} else {
									childWidth = lineHeight
									childHeight = child.getLayout().measuredDimension(DimensionHeight) +
										child.getMarginForAxis(crossAxis, availableInnerWidth)
								}
								if !(inexactEqual(childWidth, child.getLayout().measuredDimension(DimensionWidth)) &&
									inexactEqual(childHeight, child.getLayout().measuredDimension(DimensionHeight))) {
									calculateLayoutInternal(
										child,
										childWidth,
										childHeight,
										direction,
										MeasureModeExactly,
										MeasureModeExactly,
										availableInnerWidth,
										availableInnerHeight,
										true,
										LayoutPassReasonMultilineStretch,
										layoutMarkerData,
										depth,
										generationCount,
									)
								}
							}
						case AlignBaseline:
							child.setLayoutPosition(
								currentLead+maxAscentForCurrentLine-
									calculateBaseline(child)+
									child.getInlineStartPosition(FlexDirectionColumn, direction, availableInnerCrossDim),
								EdgeTop,
							)
						}
					}
				}
			}
			currentLead += lineHeight
		}
	}
	// STEP 9: COMPUTING FINAL DIMENSIONS

	node.setLayoutMeasuredDimension(
		boundAxis(
			node,
			FlexDirectionRow,
			availableWidth-marginAxisRow,
			ownerWidth,
			ownerWidth),
		DimensionWidth,
	)

	node.setLayoutMeasuredDimension(
		boundAxis(
			node,
			FlexDirectionColumn,
			availableHeight-marginAxisColumn,
			ownerHeight,
			ownerWidth),
		DimensionHeight,
	)

	// If the user didn't specify a width or height for the node, set the dimensions based on the children.
	if measureModeMainDim == MeasureModeUndefined ||
		(node.getStyle().overflow() != OverflowScroll &&
			measureModeMainDim == MeasureModeAtMost) {
		// Clamp the size to the min/max size, if specified, and make sure it doesn't go below the padding and border amount.
		node.setLayoutMeasuredDimension(
			boundAxis(
				node,
				mainAxis,
				maxLineMainDim,
				mainAxisownerSize,
				ownerWidth),
			dimension(mainAxis),
		)
	} else if measureModeMainDim == MeasureModeAtMost &&
		node.getStyle().overflow() == OverflowScroll {
		node.setLayoutMeasuredDimension(
			maxOrDefined(
				minOrDefined(
					availableInnerMainDim+paddingAndBorderAxisMain,
					boundAxisWithinMinAndMax(
						node,
						mainAxis,
						FloatOptional{maxLineMainDim},
						mainAxisownerSize,
					).unwrap(),
				),
				paddingAndBorderAxisMain,
			),
			dimension(mainAxis),
		)
	}

	if measureModeCrossDim == MeasureModeUndefined ||
		(node.getStyle().overflow() != OverflowScroll &&
			measureModeCrossDim == MeasureModeAtMost) {
		// Clamp the size to the min/max size, if specified, and make sure it doesn't go below the padding and border amount.
		node.setLayoutMeasuredDimension(
			boundAxis(
				node,
				crossAxis,
				totalLineCrossDim+paddingAndBorderAxisCross,
				crossAxisownerSize,
				ownerWidth,
			),
			dimension(crossAxis),
		)
	} else if measureModeCrossDim == MeasureModeAtMost &&
		node.getStyle().overflow() == OverflowScroll {
		node.setLayoutMeasuredDimension(
			maxOrDefined(
				minOrDefined(
					availableInnerCrossDim+paddingAndBorderAxisCross,
					boundAxisWithinMinAndMax(
						node,
						crossAxis,
						FloatOptional{totalLineCrossDim + paddingAndBorderAxisCross},
						crossAxisownerSize,
					).unwrap(),
				),
				paddingAndBorderAxisCross,
			),
			dimension(crossAxis),
		)
	}

	// As we only wrapped in normal direction yet, we need to reverse the positions on wrap-reverse.
	if performLayout && node.getStyle().flexWrap() == WrapWrapReverse {
		for i := uint32(0); i < childCount; i++ {
			child := node.GetChild(i)
			if child.getStyle().positionType() != PositionTypeAbsolute {
				child.setLayoutPosition(
					node.getLayout().measuredDimension(dimension(crossAxis))-
						child.getLayout().position(flexStartEdge(crossAxis))-
						child.getLayout().measuredDimension(dimension(crossAxis)),
					flexStartEdge(crossAxis),
				)
			}
		}
	}

	if performLayout {
		// STEP 10: SIZING AND POSITIONING ABSOLUTE CHILDREN
		for _, child := range node.GetChildren() {
			if child.getStyle().display() == DisplayNone ||
				child.getStyle().positionType() != PositionTypeAbsolute {
				continue
			}
			absolutePercentageAgainstPaddingEdge :=
				node.GetConfig().IsExperimentalFeatureEnabled(
					ExperimentalFeatureAbsolutePercentageAgainstPaddingEdge)

			layoutAbsoluteChild(
				node,
				child,
				If(absolutePercentageAgainstPaddingEdge, node.getLayout().measuredDimension(DimensionWidth), availableInnerWidth),
				If(isMainAxisRow, measureModeMainDim, measureModeCrossDim),
				If(absolutePercentageAgainstPaddingEdge, node.getLayout().measuredDimension(DimensionHeight), availableInnerHeight),
				direction,
				layoutMarkerData,
				depth,
				generationCount)
		}

		// STEP 11: SETTING TRAILING POSITIONS FOR CHILDREN
		needsMainTrailingPos := mainAxis == FlexDirectionRowReverse ||
			mainAxis == FlexDirectionColumnReverse
		needsCrossTrailingPos := crossAxis == FlexDirectionRowReverse ||
			crossAxis == FlexDirectionColumnReverse
		// Set trailing position if necessary.
		if needsMainTrailingPos || needsCrossTrailingPos {
			for i := uint32(0); i < childCount; i++ {
				child := node.GetChild(i)
				if child.getStyle().display() == DisplayNone {
					continue
				}
				if needsMainTrailingPos {
					setChildTrailingPosition(node, child, mainAxis)
				}

				if needsCrossTrailingPos {
					setChildTrailingPosition(node, child, crossAxis)
				}
			}
		}
	}
}

const spacer = "                                                            "

func spacerWithLength(level uint32) string {
	spacerLen := uint32(len(spacer))
	if level > spacerLen {
		return spacer
	} else {
		return spacer[spacerLen-level:]
	}
}

func measureModeName(mode MeasureMode, performLayout bool) string {
	switch mode {
	case MeasureModeUndefined:
		if performLayout {
			return "LAY_UNDEFINED"
		} else {
			return "UNDEFINED"
		}
	case MeasureModeExactly:
		if performLayout {
			return "LAY_EXACTLY"
		} else {
			return "EXACTLY"
		}
	case MeasureModeAtMost:
		if performLayout {
			return "LAY_AT_MOST"
		} else {
			return "AT_MOST"
		}
	}
	return ""
}

// This is a wrapper around the calculateLayoutImpl function. It determines
// whether the layout request is redundant and can be skipped.
//
// Parameters:
//
//	Input parameters are the same as calculateLayoutImpl (see above)
//	Return parameter is true if layout was performed, false if skipped
func calculateLayoutInternal(
	node *Node,
	availableWidth float32,
	availableHeight float32,
	ownerDirection Direction,
	widthMeasureMode MeasureMode,
	heightMeasureMode MeasureMode,
	ownerWidth float32,
	ownerHeight float32,
	performLayout bool,
	reason LayoutPassReason,
	layoutMarkerData *LayoutData,
	depth uint32,
	generationCount uint32,
) bool {
	if gDebuging {
		atomic.AddUint32(&gCurrentDebugCount, 1)
		fmt.Printf("calculateLayoutInternal  %d\n", atomic.LoadUint32(&gCurrentDebugCount))
	}
	layout := node.getLayout()

	depth++

	needToVisitNode :=
		(node.IsDirty() && layout.generationCount != generationCount) ||
			layout.lastOwnerDirection != ownerDirection

	if needToVisitNode {
		layout.nextCachedMeasurementsIndex = 0
		layout.cachedLayout.availableWidth = -1
		layout.cachedLayout.availableHeight = -1
		layout.cachedLayout.widthMeasureMode = MeasureModeUndefined
		layout.cachedLayout.heightMeasureMode = MeasureModeUndefined
		layout.cachedLayout.computedWidth = -1
		layout.cachedLayout.computedHeight = -1
	}

	var cachedResults *CachedMeasurement

	// Determine whether the results are already cached. We maintain a separate
	// cache for layouts and measurements. A layout operation modifies the
	// positions and dimensions for nodes in the subtree. The algorithm assumes
	// that each node gets laid out a maximum of one time per tree layout, but
	// multiple measurements may be required to resolve all of the flex
	// dimensions. We handle nodes with measure functions specially here because
	// they are the most expensive to measure, so it's worth avoiding redundant
	// measurements if at all possible.
	if node.HasMeasureFunc() {
		marginAxisRow := node.getMarginForAxis(FlexDirectionRow, ownerWidth)
		marginAxisColumn := node.getMarginForAxis(FlexDirectionColumn, ownerWidth)

		if canUseCachedMeasurement(
			widthMeasureMode,
			availableWidth,
			heightMeasureMode,
			availableHeight,
			layout.cachedLayout.widthMeasureMode,
			layout.cachedLayout.availableWidth,
			layout.cachedLayout.heightMeasureMode,
			layout.cachedLayout.availableHeight,
			layout.cachedLayout.computedWidth,
			layout.cachedLayout.computedHeight,
			marginAxisRow,
			marginAxisColumn,
			node.GetConfig(),
		) {
			cachedResults = &layout.cachedLayout
		} else {
			for i := uint32(0); i < layout.nextCachedMeasurementsIndex; i++ {
				if canUseCachedMeasurement(
					widthMeasureMode,
					availableWidth,
					heightMeasureMode,
					availableHeight,
					layout.cachedMeasurements[i].widthMeasureMode,
					layout.cachedMeasurements[i].availableWidth,
					layout.cachedMeasurements[i].heightMeasureMode,
					layout.cachedMeasurements[i].availableHeight,
					layout.cachedMeasurements[i].computedWidth,
					layout.cachedMeasurements[i].computedHeight,
					marginAxisRow,
					marginAxisColumn,
					node.GetConfig(),
				) {
					cachedResults = &layout.cachedMeasurements[i]
					break
				}
			}
		}
	} else if performLayout {
		if inexactEqual(layout.cachedLayout.availableWidth, availableWidth) &&
			inexactEqual(layout.cachedLayout.availableHeight, availableHeight) &&
			layout.cachedLayout.widthMeasureMode == widthMeasureMode &&
			layout.cachedLayout.heightMeasureMode == heightMeasureMode {
			cachedResults = &layout.cachedLayout
		}
	} else {
		for i := uint32(0); i < layout.nextCachedMeasurementsIndex; i++ {
			if inexactEqual(layout.cachedMeasurements[i].availableWidth, availableWidth) &&
				inexactEqual(layout.cachedMeasurements[i].availableHeight, availableHeight) &&
				layout.cachedMeasurements[i].widthMeasureMode == widthMeasureMode &&
				layout.cachedMeasurements[i].heightMeasureMode == heightMeasureMode {
				cachedResults = &layout.cachedMeasurements[i]
				break
			}
		}
	}

	if !needToVisitNode && cachedResults != nil {
		layout.setMeasuredDimension(DimensionWidth, cachedResults.computedWidth)
		layout.setMeasuredDimension(DimensionHeight, cachedResults.computedHeight)

		if performLayout {
			layoutMarkerData.cachedLayouts += 1
		} else {
			layoutMarkerData.cachedMeasures += 1
		}
		if gPrintChanges && gPrintSkips {
			vlog(nil, node, LogLevelVerbose, "%s%d.{[skipped] ", spacerWithLength(depth), depth)
			node.print()
			vlog(nil, node, LogLevelVerbose, "wm: %s, hm: %s, aw: %f ah: %f => d: (%f, %f) %s\n",
				measureModeName(widthMeasureMode, performLayout),
				measureModeName(heightMeasureMode, performLayout),
				availableWidth,
				availableHeight,
				cachedResults.computedWidth,
				cachedResults.computedHeight,
				reason)
		}
	} else {
		if gPrintChanges {
			vlog(nil,
				node,
				LogLevelVerbose,
				"%s%d.{%s",
				spacerWithLength(depth),
				depth,
				If(needToVisitNode, "*", ""))
			node.print()
			vlog(nil,
				node,
				LogLevelVerbose,
				"wm: %s, hm: %s, aw: %f ah: %f %s\n",
				measureModeName(widthMeasureMode, performLayout),
				measureModeName(heightMeasureMode, performLayout),
				availableWidth,
				availableHeight,
				reason)
		}
		calculateLayoutImpl(
			node,
			availableWidth,
			availableHeight,
			ownerDirection,
			widthMeasureMode,
			heightMeasureMode,
			ownerWidth,
			ownerHeight,
			performLayout,
			layoutMarkerData,
			depth,
			generationCount,
			reason,
		)
		if gPrintChanges {
			vlog(nil,
				node,
				LogLevelVerbose,
				"%s%d.}%s",
				spacerWithLength(depth),
				depth,
				If(needToVisitNode, "*", ""))
			node.print()
			vlog(nil,
				node,
				LogLevelVerbose,
				"wm: %s, hm: %s, d: (%f, %f) %s\n",
				measureModeName(widthMeasureMode, performLayout),
				measureModeName(heightMeasureMode, performLayout),
				layout.measuredDimension(DimensionWidth),
				layout.measuredDimension(DimensionHeight),
				reason)
		}
		layout.lastOwnerDirection = ownerDirection

		if cachedResults == nil {
			layoutMarkerData.maxMeasureCache = max(layoutMarkerData.maxMeasureCache, layout.nextCachedMeasurementsIndex+1)

			if layout.nextCachedMeasurementsIndex == uint32(MaxCachedMeasurements) {
				if gPrintChanges {
					vlog(nil, node, LogLevelVerbose, "Out of cache entries!\n")
				}
				layout.nextCachedMeasurementsIndex = 0
			}

			var newCacheEntry *CachedMeasurement
			if performLayout {
				newCacheEntry = &layout.cachedLayout
			} else {
				newCacheEntry = &layout.cachedMeasurements[layout.nextCachedMeasurementsIndex]
				layout.nextCachedMeasurementsIndex++
			}

			newCacheEntry.availableWidth = availableWidth
			newCacheEntry.availableHeight = availableHeight
			newCacheEntry.widthMeasureMode = widthMeasureMode
			newCacheEntry.heightMeasureMode = heightMeasureMode
			newCacheEntry.computedWidth = layout.measuredDimension(DimensionWidth)
			newCacheEntry.computedHeight = layout.measuredDimension(DimensionHeight)
		}
	}

	if performLayout {
		node.setLayoutDimension(layout.measuredDimension(DimensionWidth), DimensionWidth)
		node.setLayoutDimension(layout.measuredDimension(DimensionHeight), DimensionHeight)

		node.SetHasNewLayout(true)
		node.setDirty(false)
	}

	layout.generationCount = generationCount

	var layoutType LayoutType
	if performLayout {
		if !needToVisitNode && cachedResults == &layout.cachedLayout {
			layoutType = LayoutTypeCachedLayout
		} else {
			layoutType = LayoutTypeLayout
		}
	} else {
		if cachedResults != nil {
			layoutType = LayoutTypeCachedMeasure
		} else {
			layoutType = LayoutTypeMeasure
		}
	}
	_ = layoutType //TODO: use layoutType
	// EventPublishNodeLayout(node, layoutType)

	return needToVisitNode || cachedResults == nil
}

func CalculateLayout(node *Node, ownerWidth, ownerHeight float32, ownerDirection Direction) {
	//EventPublishLayoutPassStart(node)
	markerData := defaultLayoutData

	// Increment the generation count. This will force the recursive routine to
	// visit all dirty nodes at least once. Subsequent visits will be skipped if
	// the input parameters don't change.
	atomic.AddUint32(&gCurrentGenerationCount, 1)
	if gDebuging {
		atomic.AddUint32(&gCurrentDebugCount, 1)
		fmt.Printf("CalculateLayout  %d\n", atomic.LoadUint32(&gCurrentDebugCount))
	}
	node.resolveDimension()
	width := Undefined
	widthMeasureMode := MeasureModeUndefined
	style := node.getStyle()
	if styleDefinesDimension(node, FlexDirectionRow, ownerWidth) {
		width = (resolveValue(node.getResolvedDimension(dimension(FlexDirectionRow)), ownerWidth).unwrap() +
			node.getMarginForAxis(FlexDirectionRow, ownerWidth))
		widthMeasureMode = MeasureModeExactly
	} else if resolveValue(style.maxDimension(DimensionWidth).Value(), ownerWidth).isDefined() {
		width = resolveValue(style.maxDimension(DimensionWidth).Value(), ownerWidth).unwrap()
		widthMeasureMode = MeasureModeAtMost
	} else {
		width = ownerWidth
		if IsUndefined(width) {
			widthMeasureMode = MeasureModeUndefined
		} else {
			widthMeasureMode = MeasureModeExactly
		}
	}

	height := Undefined
	heightMeasureMode := MeasureModeUndefined
	if styleDefinesDimension(node, FlexDirectionColumn, ownerHeight) {
		height = (resolveValue(node.getResolvedDimension(dimension(FlexDirectionColumn)), ownerHeight).unwrap() +
			node.getMarginForAxis(FlexDirectionColumn, ownerWidth))
		heightMeasureMode = MeasureModeExactly
	} else if resolveValue(style.maxDimension(DimensionHeight).Value(), ownerHeight).isDefined() {
		height = resolveValue(style.maxDimension(DimensionHeight).Value(), ownerHeight).unwrap()
		heightMeasureMode = MeasureModeAtMost
	} else {
		height = ownerHeight
		if IsUndefined(height) {
			heightMeasureMode = MeasureModeUndefined
		} else {
			heightMeasureMode = MeasureModeExactly
		}
	}

	if calculateLayoutInternal(node, width, height, ownerDirection, widthMeasureMode, heightMeasureMode, ownerWidth, ownerHeight, true, LayoutPassReasonInitial, &markerData, 0, atomic.LoadUint32(&gCurrentGenerationCount)) {
		node.setPosition(node.getLayout().direction(), ownerWidth, ownerHeight, ownerWidth)
		roundLayoutResultsToPixelGrid(node, 0.0, 0.0)

		// #ifdef DEBUG
		if node.GetConfig().ShouldPrintTree() {
			vprint(node, PrintOptionsLayout|PrintOptionsChildren|PrintOptionsStyle)
		}
		// #endif
	}

	//EventPublishLayoutPassEnd(node, markerData)
}
