package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Text renders text with font support and multiline layout.
type Text struct {
	core.BaseElement
	content    string
	fontSize   float64
	color      color.Color
	faceSource *text.GoTextFaceSource

	// Layout strategy
	whiteSpace WhiteSpace
	wordBreak  WordBreak
	lineHeight float32 // explicit line height, 0 = auto

	// Layout cache
	cachedLayout  *textLayoutResult
	cachedFace    *text.GoTextFace
	cachedWidth   float32
	cachedContent string
}

// NewText creates a Text element.
func NewText(content string) *Text {
	t := &Text{
		content:    content,
		fontSize:   16,
		color:      color.Black,
		whiteSpace: WhiteSpaceNormal,
		wordBreak:  WordBreakNormal,
	}
	t.Init(t)
	if face, err := fonts.GetDefaultFontFace(float32(t.fontSize)); err == nil {
		t.faceSource = face.Face.Source
	}
	t.GetYoga().SetMeasureFunc(t.measure)
	return t
}

// ElementType returns type identifier.
func (t *Text) ElementType() string { return "Text" }

// IsNative returns true as Text is a native rendering component.
func (t *Text) IsNative() bool { return true }

// SyncFrom 同步新 Text 的属性到当前 Element（声明式重建）。
func (t *Text) SyncFrom(src core.Element) {
	other, ok := src.(*Text)
	if !ok {
		return
	}
	// 同步内容（触发重新排版）
	if t.content != other.content {
		t.content = other.content
		t.invalidateCache()
	}
	// 同步字体大小
	if t.fontSize != other.fontSize {
		t.fontSize = other.fontSize
		t.invalidateCache()
	}
	// 同步颜色
	if !colorsEqual(t.color, other.color) {
		t.color = other.color
		t.Mark(core.FlagNeedDraw)
	}
	// 同步 faceSource
	if t.faceSource != other.faceSource {
		t.faceSource = other.faceSource
		t.invalidateCache()
	}
	// 同步布局策略
	if t.whiteSpace != other.whiteSpace || t.wordBreak != other.wordBreak || t.lineHeight != other.lineHeight {
		t.whiteSpace = other.whiteSpace
		t.wordBreak = other.wordBreak
		t.lineHeight = other.lineHeight
		t.invalidateCache()
	}
}

func (t *Text) getFace() *text.GoTextFace {
	if t.faceSource == nil {
		return nil
	}
	return &text.GoTextFace{
		Source: t.faceSource,
		Size:   t.fontSize,
	}
}

func (t *Text) getLayoutResult(face *text.GoTextFace, maxWidth float32) textLayoutResult {
	if t.cachedLayout != nil && t.cachedFace != nil &&
		t.cachedFace.Source == face.Source && t.cachedFace.Size == face.Size &&
		t.cachedWidth == maxWidth && t.cachedContent == t.content {
		return *t.cachedLayout
	}
	result := computeTextLayout(t.content, face, maxWidth, t.whiteSpace, t.wordBreak, t.lineHeight)
	t.cachedLayout = &result
	t.cachedFace = face
	t.cachedWidth = maxWidth
	t.cachedContent = t.content
	return result
}

func (t *Text) invalidateCache() {
	t.cachedLayout = nil
	t.cachedFace = nil
	t.GetYoga().MarkDirty()
	t.Mark(core.FlagNeedLayout)
}

// measure implements yoga.MeasureFunc.
func (t *Text) measure(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
	face := t.getFace()
	if face == nil {
		return yoga.Size{
			Width:  float32(len(t.content)) * float32(t.fontSize) * 0.6,
			Height: float32(t.fontSize) * 1.5,
		}
	}

	var maxWidth float32
	if widthMode == yoga.MeasureModeExactly || widthMode == yoga.MeasureModeAtMost {
		maxWidth = width
	}
	// Undefined mode: single line, no wrapping

	result := t.getLayoutResult(face, maxWidth)
	return yoga.Size{Width: result.width, Height: result.height}
}

// Draw renders the text.
func (t *Text) Draw(screen *ebiten.Image) {
	if t.content == "" {
		return
	}
	bounds := t.GetBounds()

	face := t.getFace()
	if face == nil {
		return
	}

	result := t.getLayoutResult(face, bounds.Width)

	clr := t.color
	if clr == nil {
		clr = color.Black
	}

	op := &text.DrawOptions{}
	tr := t.GetTransform()
	op.GeoM.Concat(core.BuildTransformGeoM(bounds, tr))
	core.ApplyColorScaleAlpha(&op.ColorScale, tr.Alpha)
	op.ColorScale.ScaleWithColor(clr)
	op.LineSpacing = float64(result.lineHeight)
	text.Draw(screen, result.content, face, op)
}

// FlushMeasure recalculates text layout when content changes.
func (t *Text) FlushMeasure() {
	t.invalidateCache()
}

// Chain API

func (t *Text) SetContent(content string) *Text {
	if t.content != content {
		t.content = content
		t.invalidateCache()
	}
	return t
}

func (t *Text) SetFontSize(size float64) *Text {
	t.fontSize = size
	if face, err := fonts.GetDefaultFontFace(float32(size)); err == nil {
		t.faceSource = face.Face.Source
	}
	t.invalidateCache()
	return t
}

func (t *Text) SetColor(c color.Color) *Text {
	t.color = c
	t.Mark(core.FlagNeedDraw)
	return t
}

func (t *Text) SetFaceSource(src *text.GoTextFaceSource) *Text {
	t.faceSource = src
	t.invalidateCache()
	return t
}

// SetFontFamily sets the font family.
func (t *Text) SetFontFamily(family fonts.FontFamily) *Text {
	face, err := fonts.GetFontFace(fonts.FontDescriptor{
		Family: family, Weight: fonts.FontWeightNormal,
		Style: fonts.FontStyleNormal, Size: float32(t.fontSize),
	})
	if err == nil && face != nil {
		t.faceSource = face.Face.Source
		t.invalidateCache()
	}
	return t
}

// SetFontWeight sets the font weight.
func (t *Text) SetFontWeight(weight fonts.FontWeight) *Text {
	face, err := fonts.GetFontFace(fonts.FontDescriptor{
		Family: fonts.FontFamilySans, Weight: weight,
		Style: fonts.FontStyleNormal, Size: float32(t.fontSize),
	})
	if err == nil && face != nil {
		t.faceSource = face.Face.Source
		t.invalidateCache()
	}
	return t
}

// SetFontStyle sets the font style.
func (t *Text) SetFontStyle(style fonts.FontStyle) *Text {
	face, err := fonts.GetFontFace(fonts.FontDescriptor{
		Family: fonts.FontFamilySans, Weight: fonts.FontWeightNormal,
		Style: style, Size: float32(t.fontSize),
	})
	if err == nil && face != nil {
		t.faceSource = face.Face.Source
		t.invalidateCache()
	}
	return t
}

// SetWhiteSpace sets the white-space wrapping strategy.
func (t *Text) SetWhiteSpace(ws WhiteSpace) *Text {
	t.whiteSpace = ws
	t.invalidateCache()
	return t
}

// SetWordBreak sets the word-break strategy.
func (t *Text) SetWordBreak(wb WordBreak) *Text {
	t.wordBreak = wb
	t.invalidateCache()
	return t
}

// SetLineHeight sets explicit line height in pixels, 0 means auto.
func (t *Text) SetLineHeight(height float32) *Text {
	t.lineHeight = height
	t.invalidateCache()
	return t
}

func (t *Text) DebugProps() map[string]interface{} {
	props := make(map[string]interface{})
	props["content"] = t.content
	props["fontSize"] = t.fontSize
	if t.color != nil {
		props["color"] = colorToCSS(t.color)
	}
	return props
}
