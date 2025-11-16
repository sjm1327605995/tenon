package yoga

// LayoutLeft returns left
func (n *Node) LayoutLeft() float32 {
	return n.getLayout().position(EdgeLeft)
}

// LayoutTop
func (n *Node) LayoutTop() float32 {
	return n.getLayout().position(EdgeTop)
}

// LayoutRight
func (n *Node) LayoutRight() float32 {
	return n.getLayout().position(EdgeRight)
}

// LayoutBottom
func (n *Node) LayoutBottom() float32 {
	return n.getLayout().position(EdgeBottom)
}

// LayoutWidth returns width
func (n *Node) LayoutWidth() float32 {
	return n.getLayout().dimension(DimensionWidth)
}

// LayoutHeight returns height
func (n *Node) LayoutHeight() float32 {
	return n.getLayout().dimension(DimensionHeight)
}

// LayoutDirection
func (n *Node) LayoutDirection() Direction {
	return n.getLayout().direction()
}

// LayoutHadOverflow
func (n *Node) LayoutHadOverflow() bool {
	return n.getLayout().hadOverflow()
}

// LayoutMargin
func (n *Node) LayoutMargin(edge Edge) float32 {
	member := n.getLayout().margin_
	if edge > EdgeEnd {
		panic("Cannot get layout properties of multi-edge shorthands")
	}
	if edge == EdgeStart {
		if n.getLayout().direction() == DirectionRTL {
			return member[EdgeRight]
		} else {
			return member[EdgeLeft]
		}
	}

	if edge == EdgeEnd {
		if n.getLayout().direction() == DirectionRTL {
			return member[EdgeLeft]
		} else {
			return member[EdgeRight]
		}
	}

	return member[edge]
}

// LayoutBorder
func (n *Node) LayoutBorder(edge Edge) float32 {
	member := n.getLayout().border_
	if edge > EdgeEnd {
		panic("Cannot get layout properties of multi-edge shorthands")
	}
	if edge == EdgeStart {
		if n.getLayout().direction() == DirectionRTL {
			return member[EdgeRight]
		} else {
			return member[EdgeLeft]
		}
	}

	if edge == EdgeEnd {
		if n.getLayout().direction() == DirectionRTL {
			return member[EdgeLeft]
		} else {
			return member[EdgeRight]
		}
	}

	return member[edge]
}

// LayoutPadding
func (n *Node) LayoutPadding(edge Edge) float32 {
	member := n.getLayout().padding_
	if edge > EdgeEnd {
		panic("Cannot get layout properties of multi-edge shorthands")
	}
	if edge == EdgeStart {
		if n.getLayout().direction() == DirectionRTL {
			return member[EdgeRight]
		} else {
			return member[EdgeLeft]
		}
	}

	if edge == EdgeEnd {
		if n.getLayout().direction() == DirectionRTL {
			return member[EdgeLeft]
		} else {
			return member[EdgeRight]
		}
	}

	return member[edge]
}
