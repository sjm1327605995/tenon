package ui

import "github.com/sjm1327605995/tenon/core/ui/render"

type ImageUI struct {
	*BaseUI[ImageUI]
	style render.ImageStyle
}

func Image() *ImageUI {
	return &ImageUI{}
}

func (v *ImageUI) Render() *Element {
	element := CreateElement(v.style)
	for i := range v.PropsFunc {
		v.PropsFunc[i](element)
	}

	return element
}
