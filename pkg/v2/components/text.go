package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Text renders text with font support.
type Text struct {
	core.BaseElement
	content    string
	fontSize   float64
	color      color.Color
	faceSource *text.GoTextFaceSource
}

// NewText creates a Text element.
func NewText(content string) *Text {
	t := &Text{
		content:  content,
		fontSize: 16,
		color:    color.Black,
	}
	t.Init(t)
	// Auto-load default font face source if available
	if face, err := fonts.GetDefaultFontFace(float32(t.fontSize)); err == nil {
		t.faceSource = face.Face.Source
	}
	// Text nodes measure their own size
	t.GetYoga().SetMeasureFunc(t.measure)
	return t
}

// ElementType returns type identifier.
func (t *Text) ElementType() string { return "Text" }

// Draw renders the text.
func (t *Text) Draw(screen *ebiten.Image) {
	if !t.IsVisible() || t.content == "" {
		return
	}
	bounds := t.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	clr := t.color
	if clr == nil {
		clr = color.Black
	}

	// Simple text rendering using Ebiten text package
	face := &text.GoTextFace{
		Source: t.faceSource,
		Size:   t.fontSize,
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bounds.X), float64(bounds.Y))
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, t.content, face, op)
}

// measure implements yoga.MeasureFunc.
func (t *Text) measure(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
	face := &text.GoTextFace{
		Source: t.faceSource,
		Size:   t.fontSize,
	}
	w, h := text.Measure(t.content, face, 0)
	return yoga.Size{Width: float32(w), Height: float32(h)}
}

// FlushMeasure recalculates text layout when content changes.
func (t *Text) FlushMeasure() {
	// Trigger yoga remeasure
	t.GetYoga().MarkDirty()
	t.Mark(core.FlagNeedLayout)
}

// Chain API

func (t *Text) SetContent(content string) *Text {
	if t.content != content {
		t.content = content
		t.Mark(core.FlagNeedMeasure | core.FlagNeedDraw)
	}
	return t
}

func (t *Text) SetFontSize(size float64) *Text {
	t.fontSize = size
	if face, err := fonts.GetDefaultFontFace(float32(size)); err == nil {
		t.faceSource = face.Face.Source
	}
	t.Mark(core.FlagNeedMeasure | core.FlagNeedDraw)
	return t
}

func (t *Text) SetColor(c color.Color) *Text {
	t.color = c
	t.Mark(core.FlagNeedDraw)
	return t
}

func (t *Text) SetFaceSource(src *text.GoTextFaceSource) *Text {
	t.faceSource = src
	t.Mark(core.FlagNeedMeasure | core.FlagNeedDraw)
	return t
}
