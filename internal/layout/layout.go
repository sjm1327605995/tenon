package layout

import (
	"github.com/sjm1327605995/tenon/yoga"
	"github.com/sjm1327605995/tenon/pkg/types"
)

func CreateElement(element types.Element, children ...types.Element) {
	if element == nil {
		return
	}

	node := yoga.NewNode()
	element.SetNode(node)

	if props := element.GetProps(); props != nil {
		props.ApplyStyle(node)
	}

	for _, child := range element.GetChildren() {
		if child != nil {
			CreateElement(child)
			if childNode := child.GetNode(); childNode != nil {
				node.InsertChild(childNode, node.GetChildCount())
			}
		}
	}
}

func CalculateLayout(element types.Element, width, height float64) {
	if element == nil {
		return
	}

	node := element.GetNode()
	if node == nil {
		return
	}

	node.CalculateLayout(float32(width), float32(height), yoga.DirectionLTR)

	for _, child := range element.GetChildren() {
		CalculateLayout(child, width, height)
	}
}

func GetLayout(element types.Element) types.LayoutRect {
	if element == nil {
		return types.LayoutRect{}
	}

	node := element.GetNode()
	if node == nil {
		return types.LayoutRect{}
	}

	return types.LayoutRect{
		X:      node.LayoutLeft(),
		Y:      node.LayoutTop(),
		Width:  node.LayoutWidth(),
		Height: node.LayoutHeight(),
	}
}

func UpdateElementLayout(element types.Element) {
	if element == nil {
		return
	}

	node := element.GetNode()
	if node != nil {
		element.SetLayout(types.LayoutRect{
			X:      node.LayoutLeft(),
			Y:      node.LayoutTop(),
			Width:  node.LayoutWidth(),
			Height: node.LayoutHeight(),
		})
	}

	for _, child := range element.GetChildren() {
		UpdateElementLayout(child)
	}
}