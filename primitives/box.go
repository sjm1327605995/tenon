package primitives

import (
	"github.com/sjm1327605995/tenon/event"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/widget"
)

// Box is a container widget with optional background, border, and padding.
type Box struct {
	widget.WidgetBase

	Background widget.Color
	BorderColor widget.Color
	BorderWidth float32
	Padding     geometry.Insets
	Radius      float32
	Child       widget.Widget
}

// NewBox creates a new Box widget.
func NewBox(child widget.Widget) *Box {
	b := &Box{Child: child}
	b.SetVisible(true)
	b.SetEnabled(true)
	return b
}

// Layout calculates the box size.
func (b *Box) Layout(_ widget.Context, constraints geometry.Constraints) geometry.Size {
	innerConstraints := constraints.Deflate(b.Padding)
	var childSize geometry.Size
	if b.Child != nil {
		childSize = b.Child.Layout(nil, innerConstraints)
	}
	size := geometry.Sz(
		childSize.Width+b.Padding.Horizontal(),
		childSize.Height+b.Padding.Vertical(),
	)
	result := constraints.Constrain(size)
	b.SetBounds(geometry.FromPointSize(geometry.Pt(0, 0), result))
	if b.Child != nil {
		childPos := geometry.Pt(b.Padding.Left, b.Padding.Top)
		if sb, ok := b.Child.(interface{ SetBounds(geometry.Rect) }); ok {
			sb.SetBounds(geometry.FromPointSize(childPos, childSize))
		}
	}
	return result
}

// Draw renders the box and its child.
func (b *Box) Draw(ctx widget.Context, canvas widget.Canvas) {
	if !b.IsVisible() {
		return
	}
	bounds := b.Bounds()
	if !b.Background.IsTransparent() {
		if b.Radius > 0 {
			canvas.DrawRoundRect(bounds, b.Background, b.Radius)
		} else {
			canvas.DrawRect(bounds, b.Background)
		}
	}
	if b.BorderWidth > 0 && !b.BorderColor.IsTransparent() {
		if b.Radius > 0 {
			canvas.StrokeRoundRect(bounds, b.BorderColor, b.Radius, b.BorderWidth)
		} else {
			canvas.StrokeRect(bounds, b.BorderColor, b.BorderWidth)
		}
	}
	if b.Child != nil {
		canvas.PushTransform(b.Bounds().Min)
		widget.DrawChild(b.Child, ctx, canvas)
		canvas.PopTransform()
	}
}

// Event handles events.
func (b *Box) Event(_ widget.Context, e event.Event) bool {
	if b.Child != nil {
		return b.Child.Event(nil, e)
	}
	return false
}

// Children returns the child widgets.
func (b *Box) Children() []widget.Widget {
	if b.Child != nil {
		return []widget.Widget{b.Child}
	}
	return nil
}
