package ui

import "gioui.org/io/event"

// AccessibilityRole identifies the platform-facing role of a UI element.
type AccessibilityRole uint8

const (
	AccessibilityStaticText AccessibilityRole = iota + 1
	AccessibilityButton
	AccessibilityTextInput
	AccessibilityCheckbox
	AccessibilitySwitch
	AccessibilityRadioButton
)

// AccessibilityNode is a backend-neutral semantic description. Parent is the
// node's index in the snapshot, or -1 for a root. Bounds are physical pixels.
type AccessibilityNode struct {
	Role        AccessibilityRole
	Label       string
	Value       string
	Enabled     bool
	Focused     bool
	Selected    bool
	HasSelected bool
	Bounds      Rect
	Parent      int
	nativeTag   event.Tag
}

// AccessibilityPainter is an optional Painter extension for semantic output.
// Backends that do not implement it retain their existing behavior.
type AccessibilityPainter interface {
	PushAccessibility(AccessibilityNode)
	PopAccessibility()
}

// AccessibilityLabel overrides the accessible name of a text, button or input node.
func AccessibilityLabel(label string) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.accessibilityLabel = label }}
}

// AccessibilityRoleOf overrides the inferred interactive role. It also keeps
// disabled composite controls present when they have no click handler.
func AccessibilityRoleOf(role AccessibilityRole) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.accessibilityRole = role }}
}

// AccessibilitySelected reports a checkbox, switch or radio button's state.
func AccessibilitySelected(selected bool) *Node {
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) {
		hp.accessibilitySelected = selected
		hp.accessibilitySelectedSet = true
	}}
}

func accessibilityFor(rn *renderNode) (AccessibilityNode, bool) {
	node := AccessibilityNode{Enabled: !rn.isDisabled(), Focused: isFocused(rn), Bounds: rn.bounds}
	switch rn.kind {
	case rnText:
		node.Role = AccessibilityStaticText
		node.Label = rn.text
	case rnInput:
		node.Role = AccessibilityTextInput
		node.Label = rn.placeholder
		node.Value = rn.value
	default:
		node.Role = rn.accessibilityRole
		if node.Role == 0 {
			if rn.onClick == nil {
				return AccessibilityNode{}, false
			}
			node.Role = AccessibilityButton
		}
		node.Label = accessibilityText(rn)
	}
	node.Selected = rn.accessibilitySelected
	node.HasSelected = rn.accessibilitySelectedSet
	if rn.accessibilityLabel != "" {
		node.Label = rn.accessibilityLabel
	}
	if node.Role != AccessibilityStaticText {
		if rn.accessibilityTag == nil {
			rn.accessibilityTag = new(int)
		}
		node.nativeTag = rn.accessibilityTag
	}
	if node.Role == AccessibilityStaticText && (node.Label == "" || hasInteractiveAccessibilityAncestor(rn)) {
		return AccessibilityNode{}, false
	}
	return node, true
}

func hasInteractiveAccessibilityAncestor(rn *renderNode) bool {
	for parent := rn.parent; parent != nil; parent = parent.parent {
		if parent.kind == rnInput || parent.onClick != nil || parent.accessibilityRole != 0 {
			return true
		}
	}
	return false
}

func accessibilityText(rn *renderNode) string {
	if rn.kind == rnText {
		return rn.text
	}
	label := ""
	for _, child := range rn.children {
		text := accessibilityText(child)
		if text == "" {
			continue
		}
		if label != "" {
			label += " "
		}
		label += text
	}
	return label
}

func collectAccessibility(rn *renderNode, parent int, nodes *[]AccessibilityNode) {
	if rn == nil {
		return
	}
	if node, ok := accessibilityFor(rn); ok {
		node.Parent = parent
		*nodes = append(*nodes, node)
		parent = len(*nodes) - 1
	}
	for _, child := range rn.children {
		collectAccessibility(child, parent, nodes)
	}
}

// AccessibilitySnapshot returns the current semantic tree in depth-first order.
func (h *Harness) AccessibilitySnapshot() []AccessibilityNode {
	if h == nil || h.g == nil {
		return nil
	}
	return accessibilitySnapshot(h.g)
}

// AccessibilitySnapshot returns the current semantic tree in depth-first order.
func (r *BackendRuntime) AccessibilitySnapshot() []AccessibilityNode {
	if r == nil || r.closed {
		return nil
	}
	return accessibilitySnapshot(r.g)
}

func accessibilitySnapshot(g *game) []AccessibilityNode {
	if g == nil {
		return nil
	}
	var nodes []AccessibilityNode
	collectAccessibility(g.rootRN, -1, &nodes)
	for _, portal := range g.portals {
		if portal.overlayRoot != nil {
			collectAccessibility(portal.overlayRoot, -1, &nodes)
		}
	}
	return nodes
}
