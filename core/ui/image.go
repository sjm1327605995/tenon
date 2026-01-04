package ui

import (
	"github.com/sjm1327605995/tenon/core/ui/render"
	"github.com/sjm1327605995/tenon/core/ui/style"
)

type ImageUI struct {
	*BaseUI[ImageUI]
	style render.ImageStyle
}

func Image() *ImageUI {
	img := new(ImageUI)
	img.BaseUI = NewBaseUI[ImageUI](img)
	return img
}
func (v *ImageUI) Source(src string) *ImageUI {
	v.style.Src = src
	return v
}
func (v *ImageUI) Fit(f style.Fit) *ImageUI {
	v.style.Fit = f
	return v
}
func (v *ImageUI) Render() *Element {
	element := CreateElement(v.style)
	for i := range v.PropsFunc {
		v.PropsFunc[i](element)
	}

	return element
}
