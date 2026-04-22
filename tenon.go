package tenon

import (
	"github.com/sjm1327605995/tenon/internal/renderer"
	"github.com/sjm1327605995/tenon/pkg/core"
)

type Component = core.Component
type LayoutBounds = core.LayoutBounds
type Element = core.Element
type Style = core.Style
type VisualStyle = core.VisualStyle

var NewStyle = core.NewStyle
var NewVisualStyle = core.NewVisualStyle

func Run(root core.Component, width, height int) {
	renderer.Run(root, width, height)
}
