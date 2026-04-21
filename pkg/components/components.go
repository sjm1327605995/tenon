package components

import (
	"github.com/sjm1327605995/tenon/pkg/types"
)

type ViewComponent struct {
	Props    *types.ViewProps
	Children []types.UI
}

func ViewFunc(props *types.ViewProps, children ...types.UI) *ViewComponent {
	return &ViewComponent{
		Props:    props,
		Children: children,
	}
}

func (c *ViewComponent) Render() types.Element {
	children := make([]types.Element, len(c.Children))
	for i, child := range c.Children {
		children[i] = child.Render()
	}
	return types.NewViewElement(c.Props, children...)
}

type TextComponent struct {
	Props *types.TextProps
}

func TextFunc(props *types.TextProps) *TextComponent {
	return &TextComponent{Props: props}
}

func (c *TextComponent) Render() types.Element {
	return types.NewTextElement(c.Props)
}

type ImageComponent struct {
	Props *types.ImageProps
}

func ImageFunc(props *types.ImageProps) *ImageComponent {
	return &ImageComponent{Props: props}
}

func (c *ImageComponent) Render() types.Element {
	return types.NewImageElement(c.Props)
}
