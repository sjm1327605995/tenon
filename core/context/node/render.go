package node

import (
	"image"
	"image/color"

	"github.com/millken/yoga"
)

type Renderer interface {
	Rectangle(startX, startY float32, currentNode *yoga.Node, radius Radius, color color.RGBA) error
	Image(startX, startY float32, node *yoga.Node, image image.Image) error
}
