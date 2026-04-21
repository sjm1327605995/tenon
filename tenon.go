package tenon

import (
	"github.com/sjm1327605995/tenon/internal/renderer"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/types"
)

type UI = types.UI
type Element = types.Element
type Props = types.Props
type LayoutRect = types.LayoutRect
type Value = types.Value
type Unit = types.Unit

type ViewProps = types.ViewProps
type ViewStyle = types.ViewStyle
type TextProps = types.TextProps
type TextStyle = types.TextStyle
type ImageProps = types.ImageProps
type ImageStyle = types.ImageStyle

type Component = core.Component
type LayoutBounds = core.LayoutBounds

const (
	UnitPx      = types.UnitPx
	UnitPercent = types.UnitPercent
	UnitAuto    = types.UnitAuto
)

var Px = types.Px
var Percent = types.Percent
var Auto = types.Auto

var ViewFunc = components.ViewFunc
var TextFunc = components.TextFunc
var ImageFunc = components.ImageFunc

var NewView = components.NewView
var NewText = components.NewText
var NewImage = components.NewImage

func Run(root core.Component, width, height int) {
	renderer.Run(root, width, height)
}

func RenderToHTML(element Element) string {
	return renderer.RenderToHTML(element)
}

func SaveHTML(element Element, filename string) error {
	return renderer.SaveHTML(element, filename)
}
