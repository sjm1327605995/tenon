package yoga

func resolveChildAlignment(node *Node, child *Node) Align {
	align := If(child.getStyle().alignSelf() == AlignAuto, node.getStyle().alignItems(), child.getStyle().alignSelf())
	if align == AlignBaseline && isColumn(node.getStyle().flexDirection()) {
		return AlignFlexStart
	}
	return align
}
