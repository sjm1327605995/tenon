package yoga

func calculateBaseline(node *Node) float32 {
	if node.HasBaselineFunc() {
		publishEvent(Event{Type: EventNodeBaselineStart, Node: node, Config: node.GetConfig()})
		baseline := node.baseline(
			node.getLayout().measuredDimension(DimensionWidth),
			node.getLayout().measuredDimension(DimensionHeight),
		)

		if IsNaN(baseline) {
			panic("Expect custom baseline function to not return NaN")
		}
		publishEvent(Event{Type: EventNodeBaselineEnd, Node: node, Config: node.GetConfig(), Baseline: baseline})
		return baseline
	}

	var baselineChild *Node
	children := node.getLayoutChildren()
	for _, child := range children {
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
	for _, child := range node.getLayoutChildren() {
		if child.getStyle().positionType() != PositionTypeAbsolute &&
			child.getStyle().alignSelf() == AlignBaseline {
			return true
		}
	}

	return false
}
