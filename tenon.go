package tenon

import (
	"github.com/sjm1327605995/tenon/internal/reconciler"
	"github.com/sjm1327605995/tenon/internal/renderer"
	"github.com/sjm1327605995/tenon/pkg/components"
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

const (
	UnitPx      = types.UnitPx
	UnitPercent = types.UnitPercent
	UnitAuto    = types.UnitAuto
)

var Px = types.Px
var Percent = types.Percent
var Auto = types.Auto

var View = components.View
var Text = components.Text
var Image = components.Image

type Reconciler = reconciler.Reconciler

func NewReconciler(rootUI UI) *Reconciler {
	return reconciler.NewReconciler(rootUI)
}

func RenderToHTML(element Element) string {
	return renderer.RenderToHTML(element)
}

func SaveHTML(element Element, filename string) error {
	return renderer.SaveHTML(element, filename)
}