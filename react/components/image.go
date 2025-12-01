package components

import (
	"image"

	"github.com/sjm1327605995/tenon/react/yoga"
)

type Image struct {
	yoga   *yoga.Node
	Origin image.Image
	Scale  float32
}

func (i *Image) Yoga() *yoga.Node {
	return i.yoga
}
func NewImage() *Image {
	return &Image{
		yoga: yoga.NewNode(),
	}
}
