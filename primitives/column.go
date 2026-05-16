package primitives

import (
	"github.com/sjm1327605995/tenon/event"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/layout"
	"github.com/sjm1327605995/tenon/widget"
	"github.com/sjm1327605995/tenon/yoga"
)

// Column lays out children vertically using Yoga.
type Column struct {
	widget.WidgetBase

	Children_ []widget.Widget
	Gap       float32
	Padding   geometry.Insets
	Justify   yoga.Justify
	Align     yoga.Align
}

// NewColumn creates a new Column widget.
func NewColumn(children ...widget.Widget) *Column {
	c := &Column{
		Children_: children,
		Justify:   yoga.JustifyFlexStart,
		Align:     yoga.AlignStretch,
	}
	c.SetVisible(true)
	c.SetEnabled(true)
	return c
}

// Layout calculates column size using Yoga.
func (c *Column) Layout(_ widget.Context, constraints geometry.Constraints) geometry.Size {
	node := layout.NewNode()
	node.SetFlexDirection(yoga.FlexDirectionColumn)
	node.SetJustifyContent(c.Justify)
	node.SetAlignItems(c.Align)
	node.SetPadding(yoga.EdgeAll, c.Padding.Top) // simplified
	node.SetGap(yoga.GutterRow, c.Gap)

	if constraints.HasBoundedWidth() {
		node.SetMaxSize(constraints.MaxWidth, 1e9)
	}
	if constraints.HasBoundedHeight() {
		node.SetMaxSize(1e9, constraints.MaxHeight)
	}

	childNodes := make([]*layout.Node, len(c.Children_))
	for i, child := range c.Children_ {
		cn := layout.NewNode()
		childNodes[i] = cn
		if child != nil {
			cn.SetMeasureFunc(func(width, height float32) geometry.Size {
				return child.Layout(nil, geometry.Constraints{
					MinWidth:  0,
					MaxWidth:  width,
					MinHeight: 0,
					MaxHeight: height,
				})
			})
		}
		node.InsertChild(cn, i)
	}

	node.CalculateLayout(
		constraints.MaxWidth,
		constraints.MaxHeight,
		yoga.DirectionLTR,
	)

	result := node.LayoutSize()
	result = constraints.Constrain(result)
	c.SetBounds(geometry.FromPointSize(geometry.Pt(0, 0), result))

	for i, child := range c.Children_ {
		if child != nil {
			rect := childNodes[i].LayoutRect()
			if sb, ok := child.(interface{ SetBounds(geometry.Rect) }); ok {
				sb.SetBounds(rect)
			}
		}
	}

	return result
}

// Draw renders the column and its children.
func (c *Column) Draw(ctx widget.Context, canvas widget.Canvas) {
	if !c.IsVisible() {
		return
	}
	for _, child := range c.Children_ {
		if child != nil {
			var childBounds geometry.Rect
			if bb, ok := child.(interface{ Bounds() geometry.Rect }); ok {
				childBounds = bb.Bounds()
			}
			canvas.PushTransform(childBounds.Min)
			widget.DrawChild(child, ctx, canvas)
			canvas.PopTransform()
		}
	}
}

// Event handles events.
func (c *Column) Event(_ widget.Context, e event.Event) bool {
	for _, child := range c.Children_ {
		if child != nil && child.Event(nil, e) {
			return true
		}
	}
	return false
}

// Children returns the child widgets.
func (c *Column) Children() []widget.Widget {
	return c.Children_
}
