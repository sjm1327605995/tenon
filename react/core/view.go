package core

import (
	"github.com/sjm1327605995/tenon/react/common"
	"github.com/sjm1327605995/tenon/react/common/component"
	"github.com/sjm1327605995/tenon/react/style"
)

type View struct {
	*component.View
	Children []common.Element
}

// Rendering 实现渲染功能，使用common.Renderer接口
func (v *View) Rendering(renderer common.Renderer) {
	renderer.DrawView(v.View)
}

// GetChildrenCount 返回View的子元素数量
func (v *View) GetChildrenCount() int {
	return len(v.Yoga().GetChildren())
}

// GetChildAt 返回指定索引的子元素
func (v *View) GetChildAt(index int) common.Element {
	if index < 0 || index >= len(v.Children) {
		return nil
	}
	return v.Children[index]
}

func (v *View) Render() common.Node {
	return v
}
func (v *View) Style(option ...style.Option) *View {
	for i := range option {
		option[i].Apply(v)
	}
	return v
}
func (v *View) Child(nodes ...common.Component) *View {
	for i := range nodes {
		element := nodes[i].Render()
		v.Yoga().InsertChild(element.Yoga(), uint32(i))
		v.Children = append(v.Children, element)
	}
	return v
}
func (v *View) GetChildren() []common.Element {
	return v.Children
}

func NewView() *View {
	return &View{
		View: component.NewView(),
	}
}
