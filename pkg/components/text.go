package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// Text 是文本宿主组件，支持多行文本和多种换行策略。
type Text struct {
	core.BaseHost
	Content    string
	FontSize   float32
	Color      color.Color
	fontFace   *text.GoTextFace

	whiteSpace WhiteSpace
	wordBreak  WordBreak
	lineHeight float32 // 显式行高，0 表示自动

	// 布局缓存
	cachedLayout   *textLayoutResult
	cachedFace     *text.GoTextFace
	cachedWidth    float32
	cachedContent  string
}

// NewText 创建一个文本组件。
func NewText(content string) *Text {
	theme := core.GetTheme()
	t := &Text{
		Content:    content,
		FontSize:   theme.FontSizeBase + 2,
		Color:      theme.TextColor,
		whiteSpace: WhiteSpaceNormal,
		wordBreak:  WordBreakNormal,
	}
	t.Init(t)
	t.GetElement().Yoga.SetMeasureFunc(t.measure)
	return t
}

func (t *Text) measure(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
	face := t.getFace()
	if face == nil {
		return yoga.Size{
			Width:  float32(len(t.Content)) * t.FontSize * 0.6,
			Height: t.FontSize * 1.5,
		}
	}

	var maxWidth float32
	if widthMode == yoga.MeasureModeExactly {
		maxWidth = width
	} else if widthMode == yoga.MeasureModeAtMost {
		maxWidth = width
	} else {
		maxWidth = 0 // Undefined: 单行模式，不自动换行
	}

	result := t.getLayoutResult(face, maxWidth)
	return yoga.Size{Width: result.width, Height: result.height}
}

func (t *Text) getLayoutResult(face *text.GoTextFace, maxWidth float32) textLayoutResult {
	if t.cachedLayout != nil && t.cachedFace == face && t.cachedWidth == maxWidth && t.cachedContent == t.Content {
		return *t.cachedLayout
	}

	result := computeTextLayout(t.Content, face, maxWidth, t.whiteSpace, t.wordBreak, t.lineHeight)
	t.cachedLayout = &result
	t.cachedFace = face
	t.cachedWidth = maxWidth
	t.cachedContent = t.Content
	return result
}

func (t *Text) invalidateCache() {
	t.cachedLayout = nil
	t.GetElement().Yoga.MarkDirty()
}

func (t *Text) getFace() *text.GoTextFace {
	if t.fontFace != nil {
		return t.fontFace
	}
	fontManager := fonts.GetFontManager()
	face, err := fontManager.GetDefaultFontFace(t.FontSize)
	if err == nil && face != nil && face.Face != nil {
		t.fontFace = face.Face
		return face.Face
	}
	return nil
}

// Draw 绘制文本。
func (t *Text) Draw(screen *ebiten.Image) {
	el := t.GetElement()
	if el == nil || !el.Visible || t.Content == "" {
		return
	}
	bounds := t.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}
	face := t.getFace()
	if face == nil {
		return
	}

	result := t.getLayoutResult(face, bounds.Width)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bounds.X), float64(bounds.Y))
	op.ColorScale.ScaleWithColor(t.Color)
	op.LineSpacing = float64(result.lineHeight)
	text.Draw(screen, result.content, face, op)
}

// ==================== 链式 API ====================

func (t *Text) SetContent(content string) *Text {
	t.Content = content
	t.invalidateCache()
	return t
}
func (t *Text) SetFontSize(size float32) *Text {
	t.FontSize = size
	t.fontFace = nil
	t.invalidateCache()
	return t
}
func (t *Text) SetColor(c color.Color) *Text {
	t.Color = c
	return t
}
func (t *Text) SetFontFamily(family fonts.FontFamily) *Text {
	fontManager := fonts.GetFontManager()
	face, err := fontManager.GetFontFace(fonts.FontDescriptor{
		Family: family, Weight: fonts.FontWeightNormal,
		Style: fonts.FontStyleNormal, Size: t.FontSize,
	})
	if err == nil && face != nil {
		t.fontFace = face.Face
		t.invalidateCache()
	}
	return t
}
func (t *Text) SetFontWeight(weight fonts.FontWeight) *Text {
	fontManager := fonts.GetFontManager()
	face, err := fontManager.GetFontFace(fonts.FontDescriptor{
		Family: fonts.FontFamilySans, Weight: weight,
		Style: fonts.FontStyleNormal, Size: t.FontSize,
	})
	if err == nil && face != nil {
		t.fontFace = face.Face
		t.invalidateCache()
	}
	return t
}
func (t *Text) SetFontStyle(style fonts.FontStyle) *Text {
	fontManager := fonts.GetFontManager()
	face, err := fontManager.GetFontFace(fonts.FontDescriptor{
		Family: fonts.FontFamilySans, Weight: fonts.FontWeightNormal,
		Style: style, Size: t.FontSize,
	})
	if err == nil && face != nil {
		t.fontFace = face.Face
		t.invalidateCache()
	}
	return t
}
func (t *Text) SetMargin(edge yoga.Edge, value float32) *Text {
	t.GetElement().Yoga.StyleSetMargin(edge, value)
	return t
}

// SetWidth 设置文本宽度（影响自动换行时的可用宽度）。
func (t *Text) SetWidth(width float32) *Text {
	t.GetElement().Yoga.StyleSetWidth(width)
	t.invalidateCache()
	return t
}

// SetHeight 设置文本高度。
func (t *Text) SetHeight(height float32) *Text {
	t.GetElement().Yoga.StyleSetHeight(height)
	t.invalidateCache()
	return t
}

// SetWhiteSpace 设置 white-space 换行策略。
func (t *Text) SetWhiteSpace(ws WhiteSpace) *Text {
	t.whiteSpace = ws
	t.invalidateCache()
	return t
}

// SetWordBreak 设置 word-break 断词策略。
func (t *Text) SetWordBreak(wb WordBreak) *Text {
	t.wordBreak = wb
	t.invalidateCache()
	return t
}

// SetLineHeight 设置显式行高（像素），0 表示自动（根据字体计算）。
func (t *Text) SetLineHeight(height float32) *Text {
	t.lineHeight = height
	t.invalidateCache()
	return t
}

// SyncFrom 同步文本属性。
func (t *Text) SyncFrom(other core.Host) {
	if o, ok := other.(*Text); ok {
		t.Content = o.Content
		t.FontSize = o.FontSize
		t.Color = o.Color
		t.whiteSpace = o.whiteSpace
		t.wordBreak = o.wordBreak
		t.lineHeight = o.lineHeight
		// 清除布局缓存，因为文本内容可能已变
		t.cachedLayout = nil
		t.cachedFace = nil
		t.cachedWidth = 0
		t.cachedContent = ""
	}
}
