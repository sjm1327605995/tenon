package component

import (
	colEmoji "eliasnaur.com/font/noto/emoji/color"
	"fmt"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	giotext "gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/sjm1327605995/tenon/react/yoga"
	"golang.org/x/image/math/fixed"
	"image"
)

type Text struct {
	Base[Text]
	fontSize unit.Sp
	content  string
}

func NewText(content string) *Text {
	text := &Text{
		content:  content,
		fontSize: 16,
	}
	text.Base = NewBase(text)
	return text
}
func (t *Text) Update(ctx layout.Context) {
	exceptX, exceptY := t.Node.StyleGetWidth(), t.Node.StyleGetHeight()
	faces, err := opentype.ParseCollection(colEmoji.TTF)
	if err != nil {
		panic(err)
	}
	if exceptX > 0 {
		ctx.Constraints.Max.X = int(exceptX)
	}
	if exceptY > 0 {
		ctx.Constraints.Max.Y = int(exceptY)
	}
	th := material.NewTheme()
	collection := gofont.Collection()

	th.Shaper = giotext.NewShaper(giotext.WithCollection(append(collection, faces...)))

	labelStyle := material.Label(th, t.fontSize, t.content)
	textGio := &TextGio{labelStyle: &labelStyle}
	t.gio = textGio
	th.Shaper.LayoutString(giotext.Parameters{
		Font:            textGio.labelStyle.Font,
		Alignment:       labelStyle.Alignment,
		PxPerEm:         fixed.I(ctx.Sp(t.fontSize)),
		MaxLines:        labelStyle.MaxLines,
		WrapPolicy:      labelStyle.WrapPolicy,
		MinWidth:        0,
		MaxWidth:        ctx.Constraints.Max.X,
		Locale:          ctx.Locale,
		LineHeightScale: labelStyle.LineHeightScale,
		LineHeight:      fixed.I(ctx.Sp(labelStyle.LineHeight)),
	}, t.content)
	it := textIterator{
		viewport: image.Rectangle{Max: ctx.Constraints.Max},
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
	t.Node.SetMeasureFunc(func(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
		fmt.Println("Measure")
		return yoga.Size{}
	})

	textGio.Size = image.Point{X: it.bounds.Max.X, Y: it.bounds.Max.Y}
	t.Node.StyleSetWidth(float32(it.bounds.Max.X))
	t.Node.StyleSetHeight(float32(it.bounds.Max.Y))

}
func (t *Text) FontSize(size unit.Sp) *Text {
	t.fontSize = size
	return t
}

func (t *Text) Content(content string) *Text {
	t.content = content
	return t
}

type TextGio struct {
	labelStyle *material.LabelStyle
	Size       image.Point
}

func (t *TextGio) Layout(gtx layout.Context) layout.Dimensions {

	gtx.Constraints.Min = t.Size
	gtx.Constraints.Max = t.Size
	return t.labelStyle.Layout(gtx)
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
