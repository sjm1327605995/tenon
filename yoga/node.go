package yoga

import (
	"fmt"
	"sync/atomic"
)

type MeasureFunc func(
	node *Node,
	width float32,
	widthMode MeasureMode,
	height float32,
	heightMode MeasureMode,
) Size

type BaselineFunc func(node *Node, width, height float32) float32

type PrintFunc func(node *Node)

type DirtiedFunc func(node *Node)

type Node struct {
	hasNewLayout_        bool
	isReferenceBaseline_ bool
	isDirty_             bool
	nodeType_            NodeType
	context_             interface{}
	measureFunc_         MeasureFunc
	baselineFunc_        BaselineFunc
	printFunc_           PrintFunc
	dirtiedFunc_         DirtiedFunc
	style_               Style
	layout_              LayoutResults
	lineIndex_           uint32
	owner_               *Node
	children_            []*Node
	config_              *Config
	resolvedDimensions_  [2]Value
}

var (
	nodeDefaults = Node{
		hasNewLayout_:        true,
		isReferenceBaseline_: false,
		isDirty_:             false,
		nodeType_:            NodeTypeDefault,
		context_:             nil,
		measureFunc_:         nil,
		baselineFunc_:        nil,
		printFunc_:           nil,
		dirtiedFunc_:         nil,
		style_:               defaultStyle,
		layout_:              LayoutResults{},
		lineIndex_:           0,
		owner_:               nil,
		children_:            make([]*Node, 0),
		config_:              &defaultConfig,
		resolvedDimensions_:  [2]Value{},
	}
)

// NewNodeWithConfig
func NewNodeWithConfig(config *Config) *Node {
	node := nodeDefaults
	node.config_ = config
	if config.UseWebDefaults() {
		node.useWebDefaults()
	}
	node.StyleSetAlignContent(AlignFlexStart)
	node.StyleSetAlignItems(AlignStretch)
	return &node
}

func NewNode() *Node {
	return NewNodeWithConfig(&defaultConfig)
}
func (node *Node) CalculateLayout(ownerWidth, ownerHeight float32, ownerDirection Direction) {
	CalculateLayout(node, ownerWidth, ownerHeight, ownerDirection)
}

func (node *Node) useWebDefaults() {
	node.StyleSetFlexDirection(FlexDirectionRow)
	node.StyleSetAlignContent(AlignStretch)
}

func (node *Node) GetContext() interface{} {
	return node.context_
}

// GetHasNewLayout
func (node *Node) GetHasNewLayout() bool {
	return node.hasNewLayout_
}

// GetNodeType
func (node *Node) GetNodeType() NodeType {
	return node.nodeType_
}

// HasMeasureFunc
func (node *Node) HasMeasureFunc() bool {
	return node.measureFunc_ != nil
}

func (node *Node) measure(
	width float32,
	widthMode MeasureMode,
	height float32,
	heightMode MeasureMode,
) Size {
	return node.measureFunc_(node, width, widthMode, height, heightMode)
}

// HasBaselineFunc
func (node *Node) HasBaselineFunc() bool {
	return node.baselineFunc_ != nil
}

// baseline
func (node *Node) baseline(width, height float32) float32 {
	return node.baselineFunc_(node, width, height)
}

// hasErrata
func (node *Node) hasErrata(errata Errata) bool {
	return node.config_.HasErrata(errata)
}

// GetDirtiedFunc
func (node *Node) GetDirtiedFunc() DirtiedFunc {
	return node.dirtiedFunc_
}

// getStyle
func (node *Node) getStyle() *Style {
	return &node.style_
}

// getLayout
func (node *Node) getLayout() *LayoutResults {
	return &node.layout_
}

// getLineIndex
func (node *Node) getLineIndex() uint32 {
	return node.lineIndex_
}

// IsReferenceBaseline
func (node *Node) IsReferenceBaseline() bool {
	return node.isReferenceBaseline_
}

// GetOwner
func (node *Node) GetOwner() *Node {
	return node.owner_
}

// GetChildren
func (node *Node) GetChildren() []*Node {
	return node.children_
}

// GetChild
func (node *Node) GetChild(index uint32) *Node {
	if index >= node.GetChildCount() {
		panic("Index out of bounds")
	}
	return node.children_[index]
}

// GetChildCount
func (node *Node) GetChildCount() uint32 {
	return uint32(len(node.children_))
}

// GetConfig
func (node *Node) GetConfig() *Config {
	return node.config_
}

// IsDirty
func (node *Node) IsDirty() bool {
	return node.isDirty_
}

// getResolvedDimensions
func (node *Node) getResolvedDimensions() [2]Value {
	return node.resolvedDimensions_
}

// getResolvedDimension
func (node *Node) getResolvedDimension(dimension Dimension) Value {
	return node.resolvedDimensions_[dimension]
}

// computeEdgeValueForColumn
func (node *Node) computeEdgeValueForColumn(
	edges [EdgeCount]CompactValue,
	edge Edge,
) CompactValue {
	if edges[edge].IsDefined() {
		return edges[edge]
	} else if edges[EdgeVertical].IsDefined() {
		return edges[EdgeVertical]
	} else {
		return edges[EdgeAll]
	}
}

// computeEdgeValueForRow
func (node *Node) computeEdgeValueForRow(
	edges [EdgeCount]CompactValue,
	rowEdge Edge,
	edge Edge,
) CompactValue {
	if edges[rowEdge].IsDefined() {
		return edges[rowEdge]
	} else if edges[edge].IsDefined() {
		return edges[edge]
	} else if edges[EdgeHorizontal].IsDefined() {
		return edges[EdgeHorizontal]
	} else {
		return edges[EdgeAll]
	}
}

// getInlineStartEdgeUsingErrata
func (node *Node) getInlineStartEdgeUsingErrata(flexDirection FlexDirection, direction Direction) Edge {
	return If(node.hasErrata(ErrataStartingEndingEdgeFromFlexDirection), flexStartEdge(flexDirection), inlineStartEdge(flexDirection, direction))
}

// getInlineEndEdgeUsingErrata
func (node *Node) getInlineEndEdgeUsingErrata(flexDirection FlexDirection, direction Direction) Edge {
	return If(node.hasErrata(ErrataStartingEndingEdgeFromFlexDirection), flexEndEdge(flexDirection), inlineEndEdge(flexDirection, direction))
}

// isFlexStartPositionDefined
func (node *Node) isFlexStartPositionDefined(axis FlexDirection) bool {
	startEdge := flexStartEdge(axis)
	leadingPosition := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().position_, EdgeStart, startEdge), node.computeEdgeValueForColumn(node.getStyle().position_, startEdge))
	return leadingPosition.IsDefined()
}

// isInlineStartPositionDefined
func (node *Node) isInlineStartPositionDefined(axis FlexDirection, direction Direction) bool {
	startEdge := node.getInlineStartEdgeUsingErrata(axis, direction)
	leadingPosition := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().position_, EdgeStart, startEdge), node.computeEdgeValueForColumn(node.getStyle().position_, startEdge))
	return leadingPosition.IsDefined()
}

// isFlexEndPositionDefined
func (node *Node) isFlexEndPositionDefined(axis FlexDirection) bool {
	endEdge := flexEndEdge(axis)
	trailingPosition := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().position_, EdgeEnd, endEdge), node.computeEdgeValueForColumn(node.getStyle().position_, endEdge))
	return trailingPosition.IsDefined()
}

// isInlineEndPositionDefined
func (node *Node) isInlineEndPositionDefined(axis FlexDirection, direction Direction) bool {
	endEdge := node.getInlineEndEdgeUsingErrata(axis, direction)
	trailingPosition := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().position_, EdgeEnd, endEdge), node.computeEdgeValueForColumn(node.getStyle().position_, endEdge))
	return trailingPosition.IsDefined()
}

// getFlexStartPosition
func (node *Node) getFlexStartPosition(axis FlexDirection, axisSize float32) float32 {
	startEdge := flexStartEdge(axis)
	leadingPosition := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().position_, EdgeStart, startEdge), node.computeEdgeValueForColumn(node.getStyle().position_, startEdge))
	return resolveCompactValue(leadingPosition, axisSize).unwrapOrDefault(0.0)
}

// getInlineStartPosition
func (node *Node) getInlineStartPosition(axis FlexDirection, direction Direction, axisSize float32) float32 {
	startEdge := node.getInlineStartEdgeUsingErrata(axis, direction)
	leadingPosition := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().position_, EdgeStart, startEdge), node.computeEdgeValueForColumn(node.getStyle().position_, startEdge))
	return resolveCompactValue(leadingPosition, axisSize).unwrapOrDefault(0.0)
}

// getFlexEndPosition
func (node *Node) getFlexEndPosition(axis FlexDirection, axisSize float32) float32 {
	endEdge := flexEndEdge(axis)
	trailingPosition := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().position_, EdgeEnd, endEdge), node.computeEdgeValueForColumn(node.getStyle().position_, endEdge))
	return resolveCompactValue(trailingPosition, axisSize).unwrapOrDefault(0.0)
}

// getInlineEndPosition
func (node *Node) getInlineEndPosition(axis FlexDirection, direction Direction, axisSize float32) float32 {
	endEdge := node.getInlineEndEdgeUsingErrata(axis, direction)
	trailingPosition := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().position_, EdgeEnd, endEdge), node.computeEdgeValueForColumn(node.getStyle().position_, endEdge))
	return resolveCompactValue(trailingPosition, axisSize).unwrapOrDefault(0.0)
}

// getFlexStartMargin
func (node *Node) getFlexStartMargin(axis FlexDirection, widthSize float32) float32 {
	startEdge := flexStartEdge(axis)
	leadingMargin := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().margin_, EdgeStart, startEdge), node.computeEdgeValueForColumn(node.getStyle().margin_, startEdge))
	return resolveCompactValue(leadingMargin, widthSize).unwrapOrDefault(0.0)
}

// getInlineStartMargin
func (node *Node) getInlineStartMargin(axis FlexDirection, direction Direction, widthSize float32) float32 {
	startEdge := node.getInlineStartEdgeUsingErrata(axis, direction)
	leadingMargin := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().margin_, EdgeStart, startEdge), node.computeEdgeValueForColumn(node.getStyle().margin_, startEdge))
	return resolveCompactValue(leadingMargin, widthSize).unwrapOrDefault(0.0)
}

// getFlexEndMargin
func (node *Node) getFlexEndMargin(axis FlexDirection, widthSize float32) float32 {
	endEdge := flexEndEdge(axis)
	trailingMargin := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().margin_, EdgeEnd, endEdge), node.computeEdgeValueForColumn(node.getStyle().margin_, endEdge))
	return resolveCompactValue(trailingMargin, widthSize).unwrapOrDefault(0.0)
}

// getInlineEndMargin
func (node *Node) getInlineEndMargin(axis FlexDirection, direction Direction, widthSize float32) float32 {
	endEdge := node.getInlineEndEdgeUsingErrata(axis, direction)
	trailingMargin := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().margin_, EdgeEnd, endEdge), node.computeEdgeValueForColumn(node.getStyle().margin_, endEdge))
	return resolveCompactValue(trailingMargin, widthSize).unwrapOrDefault(0.0)
}

// getFlexStartBorder
func (node *Node) getFlexStartBorder(axis FlexDirection, direction Direction) float32 {
	leadRelativeFlexItemEdge := flexStartRelativeEdge(axis, direction)
	startEdge := flexStartEdge(axis)
	leadingBorder := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().border_, leadRelativeFlexItemEdge, startEdge), node.computeEdgeValueForColumn(node.getStyle().border_, startEdge))
	return maxOrDefined(leadingBorder.Value().value, 0)
}

// getInlineStartBorder
func (node *Node) getInlineStartBorder(axis FlexDirection, direction Direction) float32 {
	startEdge := node.getInlineStartEdgeUsingErrata(axis, direction)
	leadingBorder := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().border_, EdgeStart, startEdge), node.computeEdgeValueForColumn(node.getStyle().border_, startEdge))
	return maxOrDefined(leadingBorder.Value().value, 0)
}

// getFlexEndBorder
func (node *Node) getFlexEndBorder(axis FlexDirection, direction Direction) float32 {
	trailRelativeFlexItemEdge := flexEndRelativeEdge(axis, direction)
	trailingBorder := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().border_, trailRelativeFlexItemEdge, flexEndEdge(axis)), node.computeEdgeValueForColumn(node.getStyle().border_, flexEndEdge(axis)))
	return maxOrDefined(trailingBorder.Value().value, 0)
}

// getInlineEndBorder
func (node *Node) getInlineEndBorder(axis FlexDirection, direction Direction) float32 {
	endEdge := node.getInlineEndEdgeUsingErrata(axis, direction)
	trailingBorder := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().border_, EdgeEnd, endEdge), node.computeEdgeValueForColumn(node.getStyle().border_, endEdge))
	return maxOrDefined(trailingBorder.Value().value, 0)
}

// getFlexStartPadding
func (node *Node) getFlexStartPadding(axis FlexDirection, direction Direction, widthSize float32) float32 {
	leadRelativeFlexItemEdge := flexStartRelativeEdge(axis, direction)
	leadingPadding := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().padding_, leadRelativeFlexItemEdge, flexStartEdge(axis)), node.computeEdgeValueForColumn(node.getStyle().padding_, flexStartEdge(axis)))
	return maxOrDefined(resolveCompactValue(leadingPadding, widthSize).unwrap(), 0)
}

// getInlineStartPadding
func (node *Node) getInlineStartPadding(axis FlexDirection, direction Direction, widthSize float32) float32 {
	startEdge := node.getInlineStartEdgeUsingErrata(axis, direction)
	leadingPadding := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().padding_, EdgeStart, startEdge), node.computeEdgeValueForColumn(node.getStyle().padding_, startEdge))
	return maxOrDefined(resolveCompactValue(leadingPadding, widthSize).unwrap(), 0)
}

// getFlexEndPadding
func (node *Node) getFlexEndPadding(axis FlexDirection, direction Direction, widthSize float32) float32 {
	trailRelativeFlexItemEdge := flexEndRelativeEdge(axis, direction)
	trailingPadding := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().padding_, trailRelativeFlexItemEdge, flexEndEdge(axis)), node.computeEdgeValueForColumn(node.getStyle().padding_, flexEndEdge(axis)))
	return maxOrDefined(resolveCompactValue(trailingPadding, widthSize).unwrap(), 0)
}

// getInlineEndPadding
func (node *Node) getInlineEndPadding(axis FlexDirection, direction Direction, widthSize float32) float32 {
	endEdge := node.getInlineEndEdgeUsingErrata(axis, direction)
	trailingPadding := If(isRow(axis), node.computeEdgeValueForRow(node.getStyle().padding_, EdgeEnd, endEdge), node.computeEdgeValueForColumn(node.getStyle().padding_, endEdge))
	return maxOrDefined(resolveCompactValue(trailingPadding, widthSize).unwrap(), 0)
}

// getFlexStartPaddingAndBorder
func (node *Node) getFlexStartPaddingAndBorder(axis FlexDirection, direction Direction, widthSize float32) float32 {
	return node.getFlexStartPadding(axis, direction, widthSize) + node.getFlexStartBorder(axis, direction)
}

// getInlineStartPaddingAndBorder
func (node *Node) getInlineStartPaddingAndBorder(axis FlexDirection, direction Direction, widthSize float32) float32 {
	return node.getInlineStartPadding(axis, direction, widthSize) + node.getInlineStartBorder(axis, direction)
}

// getFlexEndPaddingAndBorder
func (node *Node) getFlexEndPaddingAndBorder(axis FlexDirection, direction Direction, widthSize float32) float32 {
	return node.getFlexEndPadding(axis, direction, widthSize) + node.getFlexEndBorder(axis, direction)
}

// getInlineEndPaddingAndBorder
func (node *Node) getInlineEndPaddingAndBorder(axis FlexDirection, direction Direction, widthSize float32) float32 {
	return node.getInlineEndPadding(axis, direction, widthSize) + node.getInlineEndBorder(axis, direction)
}

// getMarginForAxis
func (node *Node) getMarginForAxis(axis FlexDirection, widthSize float32) float32 {
	return node.getInlineStartMargin(axis, DirectionLTR, widthSize) + node.getInlineEndMargin(axis, DirectionLTR, widthSize)
}

// getGapForAxis
func (node *Node) getGapForAxis(axis FlexDirection) float32 {
	gap := If(isRow(axis), node.getStyle().resolveColumnGap(), node.getStyle().resolveRowGap())
	return maxOrDefined(resolveCompactValue(gap, 0).unwrap(), 0)
}

// SetContext
func (node *Node) SetContext(context interface{}) {
	node.context_ = context
}

// SetPrintFunc
func (node *Node) SetPrintFunc(printFunc PrintFunc) {
	node.printFunc_ = printFunc
}

// SetHasNewLayout
func (node *Node) SetHasNewLayout(hasNewLayout bool) {
	node.hasNewLayout_ = hasNewLayout
}

// SetNodeType
func (node *Node) SetNodeType(nodeType NodeType) {
	node.nodeType_ = nodeType
}

// SetMeasureFunc
func (node *Node) SetMeasureFunc(measureFunc MeasureFunc) {
	if measureFunc == nil {
		node.SetNodeType(NodeTypeDefault)
	} else {
		if node.GetChildCount() != 0 {
			panic("Cannot set measure function: Nodes with measure functions cannot have children.")
		}
		node.SetNodeType(NodeTypeText)
	}
	node.measureFunc_ = measureFunc
}

// SetBaselineFunc
func (node *Node) SetBaselineFunc(baselineFunc BaselineFunc) {
	node.baselineFunc_ = baselineFunc
}

// SetDirtiedFunc
func (node *Node) SetDirtiedFunc(dirtiedFunc DirtiedFunc) {
	node.dirtiedFunc_ = dirtiedFunc
}

// setStyle
func (node *Node) setStyle(style Style) {
	node.style_ = style
}

// setLayout
func (node *Node) setLayout(layout LayoutResults) {
	node.layout_ = layout
}

// setLineIndex
func (node *Node) setLineIndex(lineIndex uint32) {
	node.lineIndex_ = lineIndex
}

// SetIsReferenceBaseline
func (node *Node) SetIsReferenceBaseline(isReferenceBaseline bool) {
	if node.IsReferenceBaseline() == isReferenceBaseline {
		node.isReferenceBaseline_ = isReferenceBaseline
		node.markDirtyAndPropagate()
	}
}

// setOwner
func (node *Node) setOwner(owner *Node) {
	node.owner_ = owner
}

// setChildren
func (node *Node) setChildren(children []*Node) {
	node.children_ = children
}

// SetConfig
func (node *Node) SetConfig(config *Config) {
	if config == nil {
		panic("config cannot be nil")
	}
	if config.UseWebDefaults() != node.config_.UseWebDefaults() {
		panic("UseWebDefaults may not be changed after constructing a Node")
	}
	if configUpdateInvalidatesLayout(node.config_, config) {
		node.markDirtyAndPropagate()
	}
	node.config_ = config
}

// setDirty
func (node *Node) setDirty(isDirty bool) {
	if isDirty == node.isDirty_ {
		return
	}
	node.isDirty_ = isDirty
	if isDirty && node.dirtiedFunc_ != nil {
		node.dirtiedFunc_(node)
	}
}

// setLayoutLastOwnerDirection
func (node *Node) setLayoutLastOwnerDirection(direction Direction) {
	node.getLayout().lastOwnerDirection = direction
}

// setLayoutComputedFlexBasis
func (node *Node) setLayoutComputedFlexBasis(computedFlexBasis FloatOptional) {
	node.getLayout().computedFlexBasis = computedFlexBasis
}

// setLayoutComputedFlexBasisGeneration
func (node *Node) setLayoutComputedFlexBasisGeneration(computedFlexBasisGeneration uint32) {
	node.getLayout().computedFlexBasisGeneration = computedFlexBasisGeneration
}

// setLayoutMeasuredDimension
func (node *Node) setLayoutMeasuredDimension(
	measuredDimension float32,
	dimension Dimension,
) {
	node.getLayout().setMeasuredDimension(dimension, measuredDimension)
}

// setLayoutHadOverflow
func (node *Node) setLayoutHadOverflow(hadOverflow bool) {
	node.getLayout().setHadOverflow(hadOverflow)
}

// setLayoutDimension
func (node *Node) setLayoutDimension(dimensionValue float32, dimension Dimension) {
	node.getLayout().setDimension(dimension, dimensionValue)
}

// setLayoutDirection
func (node *Node) setLayoutDirection(direction Direction) {
	node.getLayout().setDirection(direction)
}

// setLayoutMargin
func (node *Node) setLayoutMargin(margin float32, edge Edge) {
	node.getLayout().setMargin(edge, margin)
}

// setLayoutBorder
func (node *Node) setLayoutBorder(border float32, edge Edge) {
	node.getLayout().setBorder(edge, border)
}

// setLayoutPadding
func (node *Node) setLayoutPadding(padding float32, edge Edge) {
	node.getLayout().setPadding(edge, padding)
}

// setLayoutPosition
func (node *Node) setLayoutPosition(position float32, edge Edge) {
	if gDebuging {
		atomic.AddUint32(&gCurrentDebugCount, 1)
		fmt.Printf("setLayoutPosition: %d %s=%f\n", atomic.LoadUint32(&gCurrentDebugCount),
			edge.String(), position)
	}
	node.getLayout().setPosition(edge, position)
}

// relativePosition
func (node *Node) relativePosition(axis FlexDirection, direction Direction, axisSize float32) float32 {
	if node.isInlineStartPositionDefined(axis, direction) {
		return node.getInlineStartPosition(axis, direction, axisSize)
	}
	return -1 * node.getInlineEndPosition(axis, direction, axisSize)
}

// setPosition
func (node *Node) setPosition(
	direction Direction,
	mainSize float32,
	crossSize float32,
	ownerWidth float32) {
	if gDebuging {
		atomic.AddUint32(&gCurrentDebugCount, 1)
		fmt.Printf("setPosition: %d (%s,%f,%f,%f)\n", atomic.LoadUint32(&gCurrentDebugCount),
			direction.String(), mainSize, crossSize, ownerWidth)
	}
	/* Root nodes should be always layouted as LTR, so we don't return negative
	 * values. */
	directionRespectingRoot := If(node.owner_ != nil, direction, DirectionLTR)
	mainAxis := resolveDirection(node.getStyle().flexDirection(), directionRespectingRoot)
	crossAxis := resolveCrossDirection(mainAxis, directionRespectingRoot)
	// Here we should check for `PositionType::Static` and in this case zero inset
	// properties (left, right, top, bottom, begin, end).
	// https://www.w3.org/TR/css-position-3/#valdef-position-static
	relativePositionMain := node.relativePosition(mainAxis, directionRespectingRoot, mainSize)
	relativePositionCross := node.relativePosition(crossAxis, directionRespectingRoot, crossSize)

	mainAxisLeadingEdge := node.getInlineStartEdgeUsingErrata(mainAxis, direction)
	mainAxisTrailingEdge := node.getInlineEndEdgeUsingErrata(mainAxis, direction)
	crossAxisLeadingEdge := node.getInlineStartEdgeUsingErrata(crossAxis, direction)
	crossAxisTrailingEdge := node.getInlineEndEdgeUsingErrata(crossAxis, direction)

	node.setLayoutPosition((node.getInlineStartMargin(mainAxis, direction, ownerWidth) + relativePositionMain), mainAxisLeadingEdge)
	node.setLayoutPosition((node.getInlineEndMargin(mainAxis, direction, ownerWidth) + relativePositionMain), mainAxisTrailingEdge)
	node.setLayoutPosition((node.getInlineStartMargin(crossAxis, direction, ownerWidth) + relativePositionCross), crossAxisLeadingEdge)
	node.setLayoutPosition((node.getInlineEndMargin(crossAxis, direction, ownerWidth) + relativePositionCross), crossAxisTrailingEdge)
}

// getFlexStartMarginValue
func (node *Node) getFlexStartMarginValue(axis FlexDirection) Value {
	if isRow(axis) && node.getStyle().margin_[EdgeStart].IsDefined() {
		return node.getStyle().margin(EdgeStart).Value()
	} else {
		return node.getStyle().margin(flexStartEdge(axis)).Value()
	}
}

// marginTrailingValue
func (node *Node) marginTrailingValue(axis FlexDirection) Value {
	if isRow(axis) && node.getStyle().margin_[EdgeEnd].IsDefined() {
		return node.getStyle().margin(EdgeEnd).Value()
	} else {
		return node.getStyle().margin(flexEndEdge(axis)).Value()
	}
}

// resolveFlexBasisPtr
func (node *Node) resolveFlexBasisPtr() Value {
	flexBasis := node.getStyle().flexBasis().Value()
	if flexBasis.unit != UnitAuto && flexBasis.unit != UnitUndefined {
		return flexBasis
	}
	if node.getStyle().flex().isDefined() && node.getStyle().flex().unwrap() > 0 {
		return If(node.GetConfig().UseWebDefaults(), ValueAuto, ValueZero)
	}
	return ValueAuto
}

// resolveDimension
func (node *Node) resolveDimension() {
	style := node.getStyle()
	for dim := DimensionWidth; dim < DimensionCount; dim++ {
		if style.maxDimension(dim).IsDefined() &&
			inexactEquals(style.maxDimension(dim), style.minDimension(dim)) {
			node.resolvedDimensions_[dim] = style.maxDimension(dim).Value()
		} else {
			node.resolvedDimensions_[dim] = style.dimension(dim).Value()
		}
	}
}

// resolveDirection
func (node *Node) resolveDirection(ownerDirection Direction) Direction {
	if node.getStyle().direction() == DirectionInherit {
		return If(ownerDirection != DirectionInherit, ownerDirection, DirectionLTR)
	} else {
		return node.getStyle().direction()
	}
}

// clearChildren
func (node *Node) clearChildren() {
	node.children_ = make([]*Node, 0)
}

// replaceChild
func (node *Node) replaceChild(oldChild, newChild *Node) {
	for i, child := range node.children_ {
		if child == oldChild {
			node.children_[i] = newChild
			break
		}
	}
}

func (node *Node) replaceChildIdx(child *Node, index uint32) {
	node.children_[index] = child
}

// insertChild
func (node *Node) InsertChild(child *Node, index uint32) {
	if child.GetOwner() != nil {
		panic("Child already has a owner, it must be removed first.")
	}
	if node.HasMeasureFunc() {
		panic("Cannot add child: Nodes with measure functions cannot have children.")
	}
	node.children_ = append(node.children_, nil)
	copy(node.children_[index+1:], node.children_[index:])
	node.children_[index] = child
	child.setOwner(node)
	node.markDirtyAndPropagate()
}

// SwapChild
func (node *Node) SwapChild(child *Node, index uint32) {
	node.replaceChildIdx(child, index)
	child.setOwner(node)
}

func (node *Node) RemoveChild(child *Node) {
	for i, c := range node.children_ {
		if c == child {
			node.removeChildIdx(uint32(i))
			child.setOwner(nil)
			break
		}
	}
}

func (node *Node) removeChildIdx(index uint32) {
	copy(node.children_[index:], node.children_[index+1:])
	node.children_[len(node.children_)-1] = nil
	node.children_ = node.children_[:len(node.children_)-1]
}

// RemoveAllChildren
func (node *Node) RemoveAllChildren() {
	firstChild := node.GetChild(0)
	if firstChild.GetOwner() == node {
		for _, child := range node.children_ {
			child.setLayout(LayoutResults{})
			child.setOwner(nil)
		}
	}
	node.clearChildren()
	node.markDirtyAndPropagate()
}

// cloneChildrenIfNeeded
func (node *Node) cloneChildrenIfNeeded() {
	for i, child := range node.children_ {
		if child.GetOwner() != node {
			child := node.config_.cloneNode(child, node, uint32(i))
			child.setOwner(node)
		}
	}
}

// MarkDirty
func (node *Node) MarkDirty() {
	node.markDirtyAndPropagate()
}

// markDirtyAndPropagate
func (node *Node) markDirtyAndPropagate() {
	if !node.isDirty_ {
		node.setDirty(true)
		node.setLayoutComputedFlexBasis(undefinedFloatOptional)
		if node.GetOwner() != nil {
			node.GetOwner().markDirtyAndPropagate()
		}
	}
}

// resolveFlexGrow
func (node *Node) resolveFlexGrow() float32 {
	if node.GetOwner() == nil {
		return 0
	}
	if node.getStyle().flexGrow().isDefined() {
		return node.getStyle().flexGrow().unwrap()
	}
	if node.getStyle().flex().isDefined() && node.getStyle().flex().unwrap() > 0 {
		return node.getStyle().flex().unwrap()
	}
	return DefaultFlexGrow
}

// resolveFlexShrink
func (node *Node) resolveFlexShrink() float32 {
	if node.GetOwner() == nil {
		return 0
	}
	if node.getStyle().flexShrink().isDefined() {
		return node.getStyle().flexShrink().unwrap()
	}
	if !node.GetConfig().UseWebDefaults() && node.getStyle().flex().isDefined() && node.getStyle().flex().unwrap() < 0 {
		return -node.getStyle().flex().unwrap()
	}
	return If(node.GetConfig().UseWebDefaults(), WebDefaultFlexShrink, DefaultFlexShrink)
}

// isNodeFlexible
func (node *Node) isNodeFlexible() bool {
	return (node.getStyle().positionType() != PositionTypeAbsolute) && (node.resolveFlexGrow() != 0 || node.resolveFlexShrink() != 0)
}

// print
func (node *Node) print() {
	if node.printFunc_ != nil {
		node.printFunc_(node)
	}
}

// Reset
func (node *Node) Reset() {
	if node.GetChildCount() != 0 {
		panic("Cannot reset a node which still has children attached")
	}
	if node.GetOwner() != nil {
		panic("Cannot reset a node still attached to a owner")
	}
	node = NewNodeWithConfig(node.GetConfig())
}
