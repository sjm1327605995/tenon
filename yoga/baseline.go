package yoga

func calculateBaseline(node *Node) float32 {
	if node.HasBaselineFunc() {
		//Event.Publish(Event.NodeBaselineStart, node)

		baseline := node.baseline(
			node.getLayout().measuredDimension(DimensionWidth),
			node.getLayout().measuredDimension(DimensionHeight),
		)

		//Event.Publish(Event.NodeBaselineEnd, node)
		if IsNaN(baseline) {
			panic("Expect custom baseline function to not return NaN")
		}
		return baseline
	}

	var baselineChild *Node
	childCount := node.GetChildCount()
	for i := uint32(0); i < childCount; i++ {
		child := node.GetChild(i)
		if child.getLineIndex() > 0 {
			break
		}
		if child.getStyle().positionType() == PositionTypeAbsolute {
			continue
		}
		if resolveChildAlignment(node, child) == AlignBaseline ||
			child.IsReferenceBaseline() {
			baselineChild = child
			break
		}

		if baselineChild == nil {
			baselineChild = child
		}
	}

	if baselineChild == nil {
		return node.getLayout().measuredDimension(DimensionHeight)
	}

	baseline := calculateBaseline(baselineChild)
	return baseline + baselineChild.getLayout().position(EdgeTop)
}

func isBaselineLayout(node *Node) bool {
	if isColumn(node.getStyle().flexDirection()) {
		return false
	}
	if node.getStyle().alignItems() == AlignBaseline {
		return true
	}
	childCount := node.GetChildCount()
	for i := uint32(0); i < childCount; i++ {
		child := node.GetChild(i)
		if child.getStyle().positionType() != PositionTypeAbsolute &&
			child.getStyle().alignSelf() == AlignBaseline {
			return true
		}
	}

	return false
}
