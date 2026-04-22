package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

type Text struct {
	core.BaseComponent
	Content  string
	FontSize float32
	Color    color.Color
	FontFace *fonts.FontFace
}

func NewText(content string) *Text {
	t := &Text{
		BaseComponent: core.NewBaseComponent(),
		Content:       content,
		FontSize:      16,
		Color:         color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
	t.Init(t)
	t.GetElement().Yoga.SetMeasureFunc(t.measure)
	return t
}

func (t *Text) measure(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
	face := t.getFace()
	if face != nil {
		w, h := text.Measure(t.Content, face, 0)
		return yoga.Size{Width: float32(w), Height: float32(h)}
	}

	contentWidth := float32(len(t.Content)) * t.FontSize * 0.6
	contentHeight := t.FontSize * 1.5
	return yoga.Size{Width: contentWidth, Height: contentHeight}
}

func (t *Text) getFace() *text.GoTextFace {
	if t.FontFace != nil && t.FontFace.Face != nil {
		return t.FontFace.Face
	}

	fontManager := fonts.GetFontManager()
	face, err := fontManager.GetDefaultFontFace(t.FontSize)
	if err == nil && face != nil {
		t.FontFace = face
		return face.Face
	}

	return nil
}

func (t *Text) Draw(screen *ebiten.Image) {
	element := t.Render()
	if element == nil || !element.Visible {
		return
	}

	bounds := t.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if t.Content == "" {
		return
	}

	face := t.getFace()
	if face == nil {
		return
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bounds.X), float64(bounds.Y))
	op.ColorScale.ScaleWithColor(t.Color)
	text.Draw(screen, t.Content, face, op)
}

func (t *Text) Update() error {
	return nil
}

func (t *Text) DrawOverlay(screen *ebiten.Image) {
}

func (t *Text) HandleInput() bool {
	return false
}

func (t *Text) SetContent(content string) *Text {
	t.Content = content
	t.GetElement().Yoga.MarkDirty()
	return t
}

func (t *Text) SetFontSize(size float32) *Text {
	t.FontSize = size
	t.FontFace = nil
	t.GetElement().Yoga.MarkDirty()
	return t
}

func (t *Text) SetColor(c color.Color) *Text {
	t.Color = c
	return t
}

func (t *Text) SetFontFace(face *fonts.FontFace) *Text {
	t.FontFace = face
	t.GetElement().Yoga.MarkDirty()
	return t
}

func (t *Text) SetFontFamily(family fonts.FontFamily) *Text {
	fontManager := fonts.GetFontManager()
	face, err := fontManager.GetFontFace(fonts.FontDescriptor{
		Family: family,
		Weight: fonts.FontWeightNormal,
		Style:  fonts.FontStyleNormal,
		Size:   t.FontSize,
	})
	if err == nil {
		t.FontFace = face
		t.GetElement().Yoga.MarkDirty()
	}
	return t
}

func (t *Text) SetFontWeight(weight fonts.FontWeight) *Text {
	if t.FontFace != nil {
		fontManager := fonts.GetFontManager()
		face, err := fontManager.GetFontFace(fonts.FontDescriptor{
			Family: t.FontFace.Descriptor.Family,
			Weight: weight,
			Style:  t.FontFace.Descriptor.Style,
			Size:   t.FontSize,
		})
		if err == nil {
			t.FontFace = face
			t.GetElement().Yoga.MarkDirty()
		}
	}
	return t
}

func (t *Text) SetFontStyle(style fonts.FontStyle) *Text {
	if t.FontFace != nil {
		fontManager := fonts.GetFontManager()
		face, err := fontManager.GetFontFace(fonts.FontDescriptor{
			Family: t.FontFace.Descriptor.Family,
			Weight: t.FontFace.Descriptor.Weight,
			Style:  style,
			Size:   t.FontSize,
		})
		if err == nil {
			t.FontFace = face
			t.GetElement().Yoga.MarkDirty()
		}
	}
	return t
}
