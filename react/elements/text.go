package elements

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/widget/material"

	colEmoji "eliasnaur.com/font/noto/emoji/color"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	giotext "gioui.org/text"
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/yoga"
	"golang.org/x/image/math/fixed"
	"image"
)

type Text struct {
	ElementBase
	Theme      *material.Theme
	LabelStyle *material.LabelStyle
}

func (v *Text) Paint(ctx layout.Context) {
	v.LabelStyle.Layout(ctx)
}

// SetStyle applies the given style to this view element.
// This method delegates to the style's Apply method to apply all style properties.
//
// Parameters:
//   - style: The style configuration to apply to this view
func (v *Text) SetStyle(style *styles.Style) {
	style.Apply(v)
}

func (v *Text) GetChildrenCount() int {
	return 0
}

// GetChildAt returns the child element at the specified index.
// Returns nil if the index is out of bounds.
//
// Parameters:
//   - index: The index of the child element to retrieve
//
// Returns:
//   - The child element at the specified index, or nil if the index is invalid
func (v *Text) GetChildAt(index int) api.Element {
	return nil
}

// Render implements the api.Component interface and returns the view itself as a renderable node.
//
// Returns:
//   - The view as a renderable api.Node
func (v *Text) Render() api.Node {
	return v
}

// Style applies the given style to this view and returns the view itself for method chaining.
// This is a convenience method that wraps SetStyle with a return value for fluent API usage.
//
// Parameters:
//   - option: The style configuration to apply to this view
//
// Returns:
//   - The view itself, allowing for method chaining
func (v *Text) Style(option *styles.Style) *Text {
	v.SetStyle(option)
	return v
}
func (v *Text) Content(str string) *Text {
	v.LabelStyle.Text = str
	return v
}
func (v *Text) GetChildren() []api.Element {
	return nil
}

// SetExtendedStyle applies extended style properties to this view element.
// This method handles specialized style types like BackgroundColor and BorderColor.
//
// Parameters:
//   - extendedStyle: The extended style to apply to this view
func (v *Text) SetExtendedStyle(extendedStyle styles.IExtendedStyle) {
	//switch e := extendedStyle.(type) {
	//
	//}
}

func NewText() *Text {
	th := material.NewTheme()
	labelStyle := material.Label(th, 16, "")
	node := yoga.NewNode()
	text := &Text{
		Theme:       th,
		LabelStyle:  &labelStyle,
		ElementBase: ElementBase{Node: node},
	}

	faces, err := opentype.ParseCollection(colEmoji.TTF)
	if err != nil {
		panic(err)
	}

	collection := gofont.Collection()

	th.Shaper = giotext.NewShaper(giotext.WithCollection(append(collection, faces...)))

	node.SetMeasureFunc(func(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
		th.Shaper.LayoutString(giotext.Parameters{
			Font:       text.LabelStyle.Font,
			Alignment:  labelStyle.Alignment,
			PxPerEm:    fixed.I(Metric.Sp(text.LabelStyle.TextSize)),
			MaxLines:   labelStyle.MaxLines,
			WrapPolicy: labelStyle.WrapPolicy,
			MinWidth:   0,
			MaxWidth:   int(width),
			//	Locale:          ctx.Locale,
			LineHeightScale: labelStyle.LineHeightScale,
			LineHeight:      fixed.I(Metric.Sp(labelStyle.LineHeight)),
		}, text.LabelStyle.Text)
		it := textIterator{
			viewport: image.Rectangle{Max: image.Pt(int(width), int(height))},
			maxLines: labelStyle.MaxLines,
		}
		lt := th.Shaper
		var glyphs [32]giotext.Glyph
		line := glyphs[:0]
		for g, ok := lt.NextGlyph(); ok; g, ok = lt.NextGlyph() {
			var ok bool
			if line, ok = it.paintGlyph(g, line); !ok {
				break
			}
		}
		fmt.Println(it.bounds.Max)
		return yoga.Size{
			Width:  float32(it.bounds.Max.X),
			Height: float32(it.bounds.Max.Y),
		}
	})
	return text
}

// textIterator computes the bounding box of and paints text.
type textIterator struct {
	// viewport is the rectangle of document coordinates that the iterator is
	// trying to fill with text.
	viewport image.Rectangle
	// maxLines is the maximum number of text lines that should be displayed.
	maxLines int

	// truncated tracks the count of truncated runes in the text.
	truncated int
	// linesSeen tracks the quantity of line endings this iterator has seen.
	linesSeen int
	// lineOff tracks the origin for the glyphs in the current line.
	lineOff f32.Point
	// padding is the space needed outside of the bounds of the text to ensure no
	// part of a glyph is clipped.
	padding image.Rectangle
	// bounds is the logical bounding box of the text.
	bounds image.Rectangle
	// visible tracks whether the most recently iterated glyph is visible within
	// the viewport.
	visible bool
	// first tracks whether the iterator has processed a glyph yet.
	first bool
	// baseline tracks the location of the first line of text's baseline.
	baseline int
}

// processGlyph checks whether the glyph is visible within the iterator's configured
// viewport and (if so) updates the iterator's text dimensions to include the glyph.
func (it *textIterator) processGlyph(g giotext.Glyph, ok bool) (visibleOrBefore bool) {
	if it.maxLines > 0 {
		if g.Flags&giotext.FlagTruncator != 0 && g.Flags&giotext.FlagClusterBreak != 0 {
			// A glyph carrying both of these flags provides the count of truncated runes.
			it.truncated = int(g.Runes)
		}
		if g.Flags&giotext.FlagLineBreak != 0 {
			it.linesSeen++
		}
		if it.linesSeen == it.maxLines && g.Flags&giotext.FlagParagraphBreak != 0 {
			return false
		}
	}
	// Compute the maximum extent to which glyphs overhang on the horizontal
	// axis.
	if d := g.Bounds.Min.X.Floor(); d < it.padding.Min.X {
		// If the distance between the dot and the left edge of this glyph is
		// less than the current padding, increase the left padding.
		it.padding.Min.X = d
	}
	if d := (g.Bounds.Max.X - g.Advance).Ceil(); d > it.padding.Max.X {
		// If the distance between the dot and the right edge of this glyph
		// minus the logical advance of this glyph is greater than the current
		// padding, increase the right padding.
		it.padding.Max.X = d
	}
	if d := (g.Bounds.Min.Y + g.Ascent).Floor(); d < it.padding.Min.Y {
		// If the distance between the dot and the top of this glyph is greater
		// than the ascent of the glyph, increase the top padding.
		it.padding.Min.Y = d
	}
	if d := (g.Bounds.Max.Y - g.Descent).Ceil(); d > it.padding.Max.Y {
		// If the distance between the dot and the bottom of this glyph is greater
		// than the descent of the glyph, increase the bottom padding.
		it.padding.Max.Y = d
	}
	logicalBounds := image.Rectangle{
		Min: image.Pt(g.X.Floor(), int(g.Y)-g.Ascent.Ceil()),
		Max: image.Pt((g.X + g.Advance).Ceil(), int(g.Y)+g.Descent.Ceil()),
	}
	if !it.first {
		it.first = true
		it.baseline = int(g.Y)
		it.bounds = logicalBounds
	}

	above := logicalBounds.Max.Y < it.viewport.Min.Y
	below := logicalBounds.Min.Y > it.viewport.Max.Y
	left := logicalBounds.Max.X < it.viewport.Min.X
	right := logicalBounds.Min.X > it.viewport.Max.X
	it.visible = !above && !below && !left && !right
	if it.visible {
		it.bounds.Min.X = min(it.bounds.Min.X, logicalBounds.Min.X)
		it.bounds.Min.Y = min(it.bounds.Min.Y, logicalBounds.Min.Y)
		it.bounds.Max.X = max(it.bounds.Max.X, logicalBounds.Max.X)
		it.bounds.Max.Y = max(it.bounds.Max.Y, logicalBounds.Max.Y)
	}
	return ok && !below
}

func fixedToFloat(i fixed.Int26_6) float32 {
	return float32(i) / 64.0
}

// paintGlyph buffers up and paints text glyphs. It should be invoked iteratively upon each glyph
// until it returns false. The line parameter should be a slice with
// a backing array of sufficient size to buffer multiple glyphs.
// A modified slice will be returned with each invocation, and is
// expected to be passed back in on the following invocation.
// This design is awkward, but prevents the line slice from escaping
// to the heap.
func (it *textIterator) paintGlyph(glyph giotext.Glyph, line []giotext.Glyph) ([]giotext.Glyph, bool) {
	visibleOrBefore := it.processGlyph(glyph, true)
	if it.visible {
		if len(line) == 0 {
			it.lineOff = f32.Point{X: fixedToFloat(glyph.X), Y: float32(glyph.Y)}.Sub(layout.FPt(it.viewport.Min))
		}
		line = append(line, glyph)
	}
	return line, visibleOrBefore
}
