package primitives

import (
	"github.com/sjm1327605995/tenon/event"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/layout"
	"github.com/sjm1327605995/tenon/widget"
	"github.com/sjm1327605995/tenon/yoga"
)

// Row lays out children horizontally using Yoga.
type Row struct {
	widget.WidgetBase

	Children_ []widget.Widget
	Gap       float32
	Padding   geometry.Insets
	Justify   yoga.Justify
	Align     yoga.Align
}

// NewRow creates a new Row widget.
func NewRow(children ...widget.Widget) *Row {
	r := &Row{
		Children_: children,
		Justify:   yoga.JustifyFlexStart,
		Align:     yoga.AlignStretch,
	}
	r.SetVisible(true)
	r.SetEnabled(true)
	return r
}

// Layout calculates row size using Yoga.
func (r *Row) Layout(_ widget.Context, constraints geometry.Constraints) geometry.Size {
	node := layout.NewNode()
	node.SetFlexDirection(yoga.FlexDirectionRow)
	node.SetJustifyContent(r.Justify)
	node.SetAlignItems(r.Align)
	node.SetPadding(yoga.EdgeAll, r.Padding.Top) // simplified
	node.SetGap(yoga.GutterColumn, r.Gap)

	if constraints.HasBoundedWidth() {
		node.SetMaxSize(constraints.MaxWidth, 1e9)
	}
	if constraints.HasBoundedHeight() {
		node.SetMaxSize(1e9, constraints.MaxHeight)
	}

	childNodes := make([]*layout.Node, len(r.Children_))
	for i, child := range r.Children_ {
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
	r.SetBounds(geometry.FromPointSize(geometry.Pt(0, 0), result))

	for i, child := range r.Children_ {
		if child != nil {
			rect := childNodes[i].LayoutRect()
			child.SetBounds(rect)
		}
	}

	return result
}

// Draw renders the row and its children.
func (r *Row) Draw(ctx widget.Context, canvas widget.Canvas) {
	if !r.IsVisible() {
		return
	}
	for _, child := range r.Children_ {
		if child != nil {
			childBounds := child.Bounds()
			canvas.PushTransform(childBounds.Min)
			widget.DrawChild(child, ctx, canvas)
			canvas.PopTransform()
		}
	}
}

// Event handles events.
func (r *Row) Event(_ widget.Context, e event.Event) bool {
	for _, child := range r.Children_ {
		if child != nil && child.Event(nil, e) {
			return true
		}
	}
	return false
}

// Children returns the child widgets.
func (r *Row) Children() []widget.Widget {
	return r.Children_
}
