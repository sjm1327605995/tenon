package yoga

import "math"

func RoundValueToPixelGrid(value float64, pointScaleFactor float64, forceCeil bool, forceFloor bool) float32 {
	scaledValue := value * pointScaleFactor
	fractial := math.Mod(scaledValue, 1.0)
	if fractial < 0.00001 {
		fractial++
	}
	if inexactEqual(fractial, 0) {
		scaledValue = scaledValue - fractial
	} else if inexactEqual(fractial, 1) {
		scaledValue = scaledValue - fractial + 1
	} else if forceCeil {
		scaledValue = scaledValue - fractial + 1
	} else if forceFloor {
		scaledValue = scaledValue - fractial
	} else {
		scaledValue = scaledValue - fractial +
			If(!math.IsNaN(fractial) &&
				(fractial > 0.5 || inexactEqual(fractial, 0.5)), 1.0, 0.0)
	}
	return If(math.IsNaN(scaledValue) || math.IsNaN(pointScaleFactor), Undefined, float32(scaledValue/pointScaleFactor))
}

func roundLayoutResultsToPixelGrid(node *Node, absoluteLeft float64, absoluteTop float64) {
	pointScaleFactor := float64(node.GetConfig().GetPointScaleFactor())
	nodeLeft := float64(node.getLayout().position(EdgeLeft))
	nodeTop := float64(node.getLayout().position(EdgeTop))

	nodeWidth := float64(node.getLayout().dimension(DimensionWidth))
	nodeHeight := float64(node.getLayout().dimension(DimensionHeight))

	absoluteNodeLeft := absoluteLeft + nodeLeft
	absoluteNodeTop := absoluteTop + nodeTop

	absoluteNodeRight := absoluteNodeLeft + nodeWidth
	absoluteNodeBottom := absoluteNodeTop + nodeHeight

	if pointScaleFactor != 0.0 {
		// If a node has a custom measure function we never want to round down its
		// size as this could lead to unwanted text truncation.
		textRounding := node.GetNodeType() == NodeTypeText

		node.setLayoutPosition(RoundValueToPixelGrid(nodeLeft, pointScaleFactor, false, textRounding), EdgeLeft)

		node.setLayoutPosition(RoundValueToPixelGrid(nodeTop, pointScaleFactor, false, textRounding), EdgeTop)

		// We multiply dimension by scale factor and if the result is close to the
		// whole number, we don't have any fraction To verify if the result is close
		// to whole number we want to check both floor and ceil numbers
		hasFractionalWidth := !inexactEqual(math.Mod(nodeWidth*pointScaleFactor, 1.0), 0) &&
			!inexactEqual(math.Mod(nodeWidth*pointScaleFactor, 1.0), 1.0)
		hasFractionalHeight := !inexactEqual(math.Mod(nodeHeight*pointScaleFactor, 1.0), 0) &&
			!inexactEqual(math.Mod(nodeHeight*pointScaleFactor, 1.0), 1.0)

		node.setLayoutDimension(
			RoundValueToPixelGrid(
				absoluteNodeRight,
				pointScaleFactor,
				(textRounding && hasFractionalWidth),
				(textRounding && !hasFractionalWidth))-
				RoundValueToPixelGrid(
					absoluteNodeLeft, pointScaleFactor, false, textRounding), DimensionWidth)

		node.setLayoutDimension(
			RoundValueToPixelGrid(
				absoluteNodeBottom,
				pointScaleFactor,
				(textRounding && hasFractionalHeight),
				(textRounding && !hasFractionalHeight))-
				RoundValueToPixelGrid(
					absoluteNodeTop, pointScaleFactor, false, textRounding),
			DimensionHeight)
	}

	for _, child := range node.GetChildren() {
		roundLayoutResultsToPixelGrid(child, absoluteNodeLeft, absoluteNodeTop)
	}
}
