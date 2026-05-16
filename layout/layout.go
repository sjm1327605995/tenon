package layout

import (
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/yoga"
)

// Node wraps a yoga.Node with widget identity.
type Node struct {
	Yoga     *yoga.Node
	WidgetID uint64
	Measure  func(width, height float32) geometry.Size
}

// NewNode creates a new layout node.
func NewNode() *Node {
	return &Node{Yoga: yoga.NewNode()}
}

// SetSize sets the node size.
func (n *Node) SetSize(width, height float32) {
	n.Yoga.StyleSetWidth(width)
	n.Yoga.StyleSetHeight(height)
}

// SetMinSize sets the node minimum size.
func (n *Node) SetMinSize(width, height float32) {
	n.Yoga.StyleSetMinWidth(width)
	n.Yoga.StyleSetMinHeight(height)
}

// SetMaxSize sets the node maximum size.
func (n *Node) SetMaxSize(width, height float32) {
	n.Yoga.StyleSetMaxWidth(width)
	n.Yoga.StyleSetMaxHeight(height)
}

// SetSizePercent sets the node size as percentage.
func (n *Node) SetSizePercent(width, height float32) {
	n.Yoga.StyleSetWidthPercent(width)
	n.Yoga.StyleSetHeightPercent(height)
}

// SetFlexDirection sets the flex direction.
func (n *Node) SetFlexDirection(direction yoga.FlexDirection) {
	n.Yoga.StyleSetFlexDirection(direction)
}

// SetJustifyContent sets the justify content.
func (n *Node) SetJustifyContent(justify yoga.Justify) {
	n.Yoga.StyleSetJustifyContent(justify)
}

// SetAlignItems sets the align items.
func (n *Node) SetAlignItems(align yoga.Align) {
	n.Yoga.StyleSetAlignItems(align)
}

// SetAlignSelf sets the align self.
func (n *Node) SetAlignSelf(align yoga.Align) {
	n.Yoga.StyleSetAlignSelf(align)
}

// SetFlexGrow sets the flex grow.
func (n *Node) SetFlexGrow(grow float32) {
	n.Yoga.StyleSetFlexGrow(grow)
}

// SetFlexShrink sets the flex shrink.
func (n *Node) SetFlexShrink(shrink float32) {
	n.Yoga.StyleSetFlexShrink(shrink)
}

// SetFlexBasis sets the flex basis.
func (n *Node) SetFlexBasis(basis float32) {
	n.Yoga.StyleSetFlexBasis(basis)
}

// SetFlexBasisPercent sets the flex basis as percentage.
func (n *Node) SetFlexBasisPercent(percent float32) {
	n.Yoga.StyleSetFlexBasisPercent(percent)
}

// SetFlexBasisAuto sets the flex basis to auto.
func (n *Node) SetFlexBasisAuto() {
	n.Yoga.StyleSetFlexBasisAuto()
}

// SetMargin sets the margin for an edge.
func (n *Node) SetMargin(edge yoga.Edge, margin float32) {
	n.Yoga.StyleSetMargin(edge, margin)
}

// SetMarginPercent sets the margin as percentage.
func (n *Node) SetMarginPercent(edge yoga.Edge, percent float32) {
	n.Yoga.StyleSetMarginPercent(edge, percent)
}

// SetPadding sets the padding for an edge.
func (n *Node) SetPadding(edge yoga.Edge, padding float32) {
	n.Yoga.StyleSetPadding(edge, padding)
}

// SetPaddingPercent sets the padding as percentage.
func (n *Node) SetPaddingPercent(edge yoga.Edge, percent float32) {
	n.Yoga.StyleSetPaddingPercent(edge, percent)
}

// SetBorder sets the border for an edge.
func (n *Node) SetBorder(edge yoga.Edge, border float32) {
	n.Yoga.StyleSetBorder(edge, border)
}

// SetPosition sets the position for an edge.
func (n *Node) SetPosition(edge yoga.Edge, position float32) {
	n.Yoga.StyleSetPosition(edge, position)
}

// SetPositionPercent sets the position as percentage.
func (n *Node) SetPositionPercent(edge yoga.Edge, percent float32) {
	n.Yoga.StyleSetPositionPercent(edge, percent)
}

// SetPositionType sets the position type.
func (n *Node) SetPositionType(positionType yoga.PositionType) {
	n.Yoga.StyleSetPositionType(positionType)
}

// SetDisplay sets the display.
func (n *Node) SetDisplay(display yoga.Display) {
	n.Yoga.StyleSetDisplay(display)
}

// SetGap sets the gap for a gutter.
func (n *Node) SetGap(gutter yoga.Gutter, gapLength float32) {
	n.Yoga.StyleSetGap(gutter, gapLength)
}

// SetWrap sets the flex wrap.
func (n *Node) SetWrap(wrap yoga.Wrap) {
	n.Yoga.StyleSetFlexWrap(wrap)
}

// InsertChild inserts a child node at the given index.
func (n *Node) InsertChild(child *Node, index int) {
	n.Yoga.InsertChild(child.Yoga, uint32(index))
}

// RemoveChild removes a child node.
func (n *Node) RemoveChild(child *Node) {
	n.Yoga.RemoveChild(child.Yoga)
}

// ChildCount returns the number of children.
func (n *Node) ChildCount() uint32 {
	return n.Yoga.GetChildCount()
}

// CalculateLayout calculates the layout for the node tree.
func (n *Node) CalculateLayout(ownerWidth, ownerHeight float32, ownerDirection yoga.Direction) {
	n.Yoga.CalculateLayout(ownerWidth, ownerHeight, ownerDirection)
}

// LayoutLeft returns the computed left position.
func (n *Node) LayoutLeft() float32 {
	return n.Yoga.LayoutLeft()
}

// LayoutTop returns the computed top position.
func (n *Node) LayoutTop() float32 {
	return n.Yoga.LayoutTop()
}

// LayoutWidth returns the computed width.
func (n *Node) LayoutWidth() float32 {
	return n.Yoga.LayoutWidth()
}

// LayoutHeight returns the computed height.
func (n *Node) LayoutHeight() float32 {
	return n.Yoga.LayoutHeight()
}

// LayoutMargin returns the computed margin for an edge.
func (n *Node) LayoutMargin(edge yoga.Edge) float32 {
	return n.Yoga.LayoutMargin(edge)
}

// LayoutPadding returns the computed padding for an edge.
func (n *Node) LayoutPadding(edge yoga.Edge) float32 {
	return n.Yoga.LayoutPadding(edge)
}

// LayoutBorder returns the computed border for an edge.
func (n *Node) LayoutBorder(edge yoga.Edge) float32 {
	return n.Yoga.LayoutBorder(edge)
}

// LayoutRect returns the computed layout as a geometry.Rect.
func (n *Node) LayoutRect() geometry.Rect {
	return geometry.Rect{
		Min: geometry.Pt(n.LayoutLeft(), n.LayoutTop()),
		Max: geometry.Pt(
			n.LayoutLeft()+n.LayoutWidth(),
			n.LayoutTop()+n.LayoutHeight(),
		),
	}
}

// LayoutSize returns the computed layout size.
func (n *Node) LayoutSize() geometry.Size {
	return geometry.Sz(n.LayoutWidth(), n.LayoutHeight())
}

// SetMeasureFunc sets a custom measure function for leaf nodes.
func (n *Node) SetMeasureFunc(measure func(width, height float32) geometry.Size) {
	n.Measure = measure
	if measure != nil {
		n.Yoga.SetMeasureFunc(func(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
			size := measure(width, height)
			return yoga.Size{Width: size.Width, Height: size.Height}
		})
	} else {
		n.Yoga.SetMeasureFunc(nil)
	}
}
