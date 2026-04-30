package native

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ShadowText renders text with a shadow/outline effect by drawing twice:
// first the shadow offset, then the main text.
type ShadowText struct {
	core.BaseElement
	Content       string
	fontSize      float64
	color         color.Color
	shadowColor   color.Color
	shadowOffsetX float32
	shadowOffsetY float32
	hAlign        text.Align
	vAlign        text.Align
	faceSource    *text.GoTextFaceSource
}

// NewShadowText creates a ShadowText element.
func NewShadowText(Content string) *ShadowText {
	st := &ShadowText{
		Content:       Content,
		fontSize:      16,
		color:         color.Black,
		shadowColor:   color.RGBA{0, 0, 0, 128},
		shadowOffsetX: 1,
		shadowOffsetY: 1,
		hAlign:        text.AlignStart,
		vAlign:        text.AlignStart,
	}
	st.Init(st)
	if face, err := fonts.GetDefaultFontFace(float32(st.fontSize)); err == nil {
		st.faceSource = face.Face.Source
	}
	st.GetYoga().SetMeasureFunc(st.measure)
	return st
}

// ElementType returns type identifier.
func (st *ShadowText) ElementType() string { return "ShadowText" }

func (st *ShadowText) getFace() *text.GoTextFace {
	if st.faceSource == nil {
		return nil
	}
	return &text.GoTextFace{Source: st.faceSource, Size: st.fontSize}
}

func (st *ShadowText) measure(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
	face := st.getFace()
	if face == nil {
		return yoga.Size{
			Width:  float32(len(st.Content)) * float32(st.fontSize) * 0.6,
			Height: float32(st.fontSize) * 1.5,
		}
	}
	w, h := text.Measure(st.Content, face, 0)
	return yoga.Size{Width: float32(w), Height: float32(h)}
}

// Draw renders shadow text.
func (st *ShadowText) Draw(screen *ebiten.Image) {
	if st.Content == "" {
		return
	}
	bounds := st.GetBounds()

	face := st.getFace()
	if face == nil {
		return
	}

	clr := st.color
	if clr == nil {
		clr = color.Black
	}
	sclr := st.shadowColor
	if sclr == nil {
		sclr = color.RGBA{0, 0, 0, 128}
	}

	tr := st.GetTransform()

	// Draw shadow
	sop := &text.DrawOptions{}
	sop.GeoM.Concat(core.BuildTransformGeoM(bounds, tr))
	core.ApplyColorScaleAlpha(&sop.ColorScale, tr.Alpha)
	sop.ColorScale.ScaleWithColor(sclr)
	sop.GeoM.Translate(float64(st.shadowOffsetX), float64(st.shadowOffsetY))
	sop.PrimaryAlign = st.hAlign
	sop.SecondaryAlign = st.vAlign
	text.Draw(screen, st.Content, face, sop)

	// Draw main text
	op := &text.DrawOptions{}
	op.GeoM.Concat(core.BuildTransformGeoM(bounds, tr))
	core.ApplyColorScaleAlpha(&op.ColorScale, tr.Alpha)
	op.ColorScale.ScaleWithColor(clr)
	op.PrimaryAlign = st.hAlign
	op.SecondaryAlign = st.vAlign
	text.Draw(screen, st.Content, face, op)
}

// Chain API

func (st *ShadowText) SetContent(Content string) *ShadowText {
	if st.Content != Content {
		st.Content = Content
		st.GetYoga().MarkDirty()
		st.Mark(core.FlagNeedLayout)
	}
	return st
}

func (st *ShadowText) SetFontSize(size float64) *ShadowText {
	st.fontSize = size
	if face, err := fonts.GetDefaultFontFace(float32(size)); err == nil {
		st.faceSource = face.Face.Source
	}
	st.GetYoga().MarkDirty()
	st.Mark(core.FlagNeedLayout)
	return st
}

func (st *ShadowText) SetColor(c color.Color) *ShadowText {
	st.color = c
	st.Mark(core.FlagNeedDraw)
	return st
}

func (st *ShadowText) SetShadowColor(c color.Color) *ShadowText {
	st.shadowColor = c
	st.Mark(core.FlagNeedDraw)
	return st
}

func (st *ShadowText) SetShadowOffset(dx, dy float32) *ShadowText {
	st.shadowOffsetX = dx
	st.shadowOffsetY = dy
	st.Mark(core.FlagNeedDraw)
	return st
}

func (st *ShadowText) SetAlign(h, v text.Align) *ShadowText {
	st.hAlign = h
	st.vAlign = v
	st.Mark(core.FlagNeedDraw)
	return st
}

func (st *ShadowText) SetFaceSource(src *text.GoTextFaceSource) *ShadowText {
	st.faceSource = src
	st.GetYoga().MarkDirty()
	st.Mark(core.FlagNeedLayout)
	return st
}

func (st *ShadowText) SetFontFamily(family fonts.FontFamily) *ShadowText {
	face, err := fonts.GetFontFace(fonts.FontDescriptor{
		Family: family, Weight: fonts.FontWeightNormal,
		Style: fonts.FontStyleNormal, Size: float32(st.fontSize),
	})
	if err == nil && face != nil {
		st.faceSource = face.Face.Source
		st.GetYoga().MarkDirty()
		st.Mark(core.FlagNeedLayout)
	}
	return st
}
