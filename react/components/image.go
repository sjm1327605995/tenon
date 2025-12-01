package components

import (
	"github.com/sjm1327605995/tenon/react/yoga"
	"image"
)

type Image struct {
	yoga   *yoga.Node
	Origin image.Image
	Path   string
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
