package ui

import (
	"image"

	"github.com/sjm1327605995/tenon/core/ui/render"
	"github.com/sjm1327605995/tenon/yoga"
)

type Element struct {
	Type         render.Widget
	Yoga         *yoga.Node
	idx          uint32
	renderObject *render.Node
	parent       *Element
}

func CreateElement(typ render.Widget) *Element {
	element := &Element{
		Type: typ,
		Yoga: yoga.NewNode(),
	}
	// 将Element实例存储到Yoga Node的context中
	element.Yoga.SetContext(element)
	return element
}

func (e *Element) InsertChild(child *Element) {
	e.Yoga.InsertChild(child.Yoga, e.idx)
	child.parent = e
	e.idx++
}

// GetChildren 从Yoga树中获取所有子Element
func (e *Element) GetChildren() []*Element {
	childCount := e.Yoga.GetChildCount()
	children := make([]*Element, 0, childCount)
	for i := uint32(0); i < childCount; i++ {
		childNode := e.Yoga.GetChild(i)
		if childElement, ok := childNode.GetContext().(*Element); ok {
			children = append(children, childElement)
		}
	}
	return children
}

// RemoveChild 从Yoga树中移除子Element
func (e *Element) RemoveChild(child *Element) {
	e.Yoga.RemoveChild(child.Yoga)
	child.parent = nil
}

// GetChildCount 获取子元素数量
func (e *Element) GetChildCount() uint32 {
	return e.Yoga.GetChildCount()
}

// GetChild 通过索引获取子Element
func (e *Element) GetChild(index uint32) *Element {
	childNode := e.Yoga.GetChild(index)
	if childElement, ok := childNode.GetContext().(*Element); ok {
		return childElement
	}
	return nil
}

func (e *Element) Mount() {
	// 先递归挂载所有子元素
	children := e.GetChildren()
	for _, child := range children {
		child.Mount()
	}

	// 再创建renderObject
	e.createRenderObject()
}

func (e *Element) createRenderObject() {
	renderObject := e.Type.ToRender()
	if renderObject.HasDefault() {

	}

	e.renderObject = render.NewNode(renderObject)
	children := e.GetChildren()
	for _, child := range children {
		if child.renderObject != nil {
			e.renderObject.InsertChild(child.renderObject)
		}
	}
}

func (e *Element) UpdateRenderObject() {
	// 更新renderObject的属性
	if e.renderObject != nil {
		width := e.Yoga.LayoutWidth()
		height := e.Yoga.LayoutHeight()
		x := e.Yoga.LayoutLeft()
		y := e.Yoga.LayoutTop()
		e.renderObject.Max = image.Point{X: int(width), Y: int(height)}

		// 计算相对于父元素的偏移量
		if e.parent != nil {
			parentX := e.parent.Yoga.LayoutLeft()
			parentY := e.parent.Yoga.LayoutTop()
			x -= parentX
			y -= parentY
		}

		e.renderObject.Offset = image.Point{X: int(x), Y: int(y)}

		// 递归更新子元素的renderObject
		children := e.GetChildren()
		for _, child := range children {
			child.UpdateRenderObject()
		}
	}
}

func (e *Element) Unmount() {
	// 递归卸载所有子元素
	children := e.GetChildren()
	for _, child := range children {
		child.Unmount()
	}

	// 清理renderObject
	e.renderObject = nil
}

func (e *Element) RenderObject() *render.Node {
	return e.renderObject
}
