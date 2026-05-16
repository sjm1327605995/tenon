package primitives

import (
	"github.com/sjm1327605995/tenon/event"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/widget"
)

// Text is a leaf widget that renders a string.
type Text struct {
	widget.WidgetBase

	Content  string
	FontSize float32
	Color    widget.Color
	Bold     bool
	Align    widget.TextAlign
}

// NewText creates a new Text widget.
func NewText(content string) *Text {
	t := &Text{
		Content:  content,
		FontSize: 14,
		Color:    widget.ColorBlack,
		Align:    widget.TextAlignLeft,
	}
	t.SetVisible(true)
	t.SetEnabled(true)
	return t
}

// Layout calculates the text size.
func (t *Text) Layout(_ widget.Context, constraints geometry.Constraints) geometry.Size {
	width := constraints.ConstrainWidth(float32(len(t.Content)) * t.FontSize * 0.6)
	height := constraints.ConstrainHeight(t.FontSize * 1.2)
	result := geometry.Sz(width, height)
	t.SetBounds(geometry.FromPointSize(geometry.Pt(0, 0), result))
	return result
}

// Draw renders the text.
func (t *Text) Draw(_ widget.Context, canvas widget.Canvas) {
	if !t.IsVisible() {
		return
	}
	canvas.DrawText(t.Content, t.Bounds(), t.FontSize, t.Color, t.Bold, t.Align)
}

// Event handles events (text is passive).
func (t *Text) Event(_ widget.Context, _ event.Event) bool {
	return false
}

// Children returns nil (leaf widget).
func (t *Text) Children() []widget.Widget {
	return nil
}
