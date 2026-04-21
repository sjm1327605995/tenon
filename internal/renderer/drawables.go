package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/types"
)

func RegisterDrawables() {
	RegisterDrawable("view", newViewDrawable)
	RegisterDrawable("text", newTextDrawable)
	RegisterDrawable("image", newImageDrawable)
}

type ViewDrawable struct {
	element types.Element
}

func newViewDrawable(element types.Element) Drawable {
	return &ViewDrawable{element: element}
}

func (d *ViewDrawable) Draw(screen *ebiten.Image, x, y int) {
	viewElement, ok := d.element.(*types.ViewElement)
	if !ok {
		return
	}

	layout := viewElement.GetLayout()
	props := viewElement.Props

	width := int(layout.Width)
	height := int(layout.Height)

	if width <= 0 || height <= 0 {
		return
	}

	if props.Style != nil && props.Style.Background != "" {
		bgColor := ParseColor(props.Style.Background)
		screen.Fill(bgColor)
	}
}

type TextDrawable struct {
	element types.Element
}

func newTextDrawable(element types.Element) Drawable {
	return &TextDrawable{element: element}
}

func (d *TextDrawable) Draw(screen *ebiten.Image, x, y int) {
}

type ImageDrawable struct {
	element types.Element
}

func newImageDrawable(element types.Element) Drawable {
	return &ImageDrawable{element: element}
}

func (d *ImageDrawable) Draw(screen *ebiten.Image, x, y int) {
	imgElement, ok := d.element.(*types.ImageElement)
	if !ok {
		return
	}

	props := imgElement.Props
	if props.Style == nil {
		return
	}

	width := int(props.Style.Width.Value)
	height := int(props.Style.Height.Value)

	if width <= 0 || height <= 0 {
		return
	}

	placeholderColor := ParseColor("#e0e0e0")
	screen.Fill(placeholderColor)
}