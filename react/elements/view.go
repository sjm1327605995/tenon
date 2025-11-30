package elements

import (
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/components"
)

type View struct {
	*components.View
	Children []api.Element
}

func (v *View) SetStyle(style *styles.Style) {
	style.Apply(v)
}

// Rendering 实现渲染功能，使用api.Renderer接口
func (v *View) Rendering(renderer api.Renderer) {
	renderer.DrawView(v.View)
}

// GetChildrenCount 返回View的子元素数量
func (v *View) GetChildrenCount() int {
	return len(v.Yoga().GetChildren())
}

// GetChildAt 返回指定索引的子元素
func (v *View) GetChildAt(index int) api.Element {
	if index < 0 || index >= len(v.Children) {
		return nil
	}
	return v.Children[index]
}

func (v *View) Render() api.Node {
	return v
}
func (v *View) Style(option *styles.Style) *View {
	v.SetStyle(option)
	return v
}
func (v *View) Child(nodes ...api.Component) *View {
	for i := range nodes {
		element := nodes[i].Render()
		v.Yoga().InsertChild(element.Yoga(), uint32(i))
		v.Children = append(v.Children, element)
	}
	return v
}
func (v *View) GetChildren() []api.Element {
	return v.Children
}

// GetView 返回内部的component.View实例，用于样式系统设置属性
func (v *View) GetView() *components.View {
	return v.View
}

func (v *View) SetExtendedStyle(extendedStyle styles.IExtendedStyle) {
	switch e := extendedStyle.(type) {
	case styles.BackgroundColor:
		v.View.Background = e.Color
	case styles.BorderColor:
		v.View.BorderColor = e.Color
	}
}

func NewView() *View {
	return &View{
		View: components.NewView(),
	}
}
