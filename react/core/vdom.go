package core

// VNode represents a virtual DOM node.
type VNode struct {
	Type     string
	Props    map[string]interface{}
	Children []*VNode
}

// NewVNode creates a new virtual DOM node.
func NewVNode(nodeType string, props map[string]interface{}, children ...*VNode) *VNode {
	if props == nil {
		props = make(map[string]interface{})
	}
	return &VNode{
		Type:     nodeType,
		Props:    props,
		Children: children,
	}
}
