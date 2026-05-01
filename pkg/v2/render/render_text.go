package render

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderText 负责文字的 Yoga 测量和 Ebiten 绘制。
type RenderText struct {
	BaseRenderObject

	Content   string
	FontSize  float32
	TextColor color.Color
	MaxLines  int
	Underline bool
}

func NewRenderText(content string) *RenderText {
	r := &RenderText{
		Content:   content,
		FontSize:  14,
		TextColor: color.Black,
		MaxLines:  0,
	}
	r.BaseRenderObject.Init(r)
	r.yoga = yoga.NewNode()
	r.yoga.SetMeasureFunc(r.measure)
	return r
}

func (r *RenderText) HitTest(x, y float32) bool {
	return false
}

func (r *RenderText) measure(
	node *yoga.Node,
	width float32,
	widthMode yoga.MeasureMode,
	height float32,
	heightMode yoga.MeasureMode,
) yoga.Size {
	face, err := fonts.GetFontFace(fonts.FontDescriptor{
		Family: fonts.FontFamilyDefault,
		Weight: fonts.FontWeightNormal,
		Style:  fonts.FontStyleNormal,
		Size:   r.FontSize,
	})
	if err != nil || face == nil {
		return yoga.Size{
			Width:  float32(len(r.Content)) * r.FontSize * 0.6,
			Height: r.FontSize * 1.5,
		}
	}

	// 使用字体实际 metrics 计算行高，确保 measure 与 Paint 高度一致。
	// 第一行高度 = HAscent + HDescent（文字实际像素高），
	// 后续行行距 = fontSize * 1.2（保持多行间距不变）。
	m := face.Face.Metrics()
	actualLineHeight := m.HAscent + m.HDescent
	lineHeight := float64(r.FontSize) * 1.2

	var maxLineWidth float64
	var totalHeight float64

	lines := r.splitLines(face.Face, float64(width), widthMode)
	for i, line := range lines {
		w, _ := text.Measure(line, face.Face, 0)
		if w > maxLineWidth {
			maxLineWidth = w
		}
		if i == 0 {
			totalHeight += actualLineHeight
		} else {
			totalHeight += lineHeight
		}
	}
	if totalHeight == 0 {
		totalHeight = actualLineHeight
	}

	return yoga.Size{
		Width:  float32(maxLineWidth),
		Height: float32(totalHeight),
	}
}

func (r *RenderText) splitLines(face text.Face, maxWidth float64, widthMode yoga.MeasureMode) []string {
	if r.Content == "" {
		return nil
	}
	// 如果宽度未定义或足够大，单行返回
	if widthMode == yoga.MeasureModeUndefined {
		return []string{r.Content}
	}
	// 简单按字符换行（实际应用应使用更复杂的文本排版）
	var lines []string
	var currentLine string
	for _, ch := range r.Content {
		testLine := currentLine + string(ch)
		w, _ := text.Measure(testLine, face, 0)
		if w > maxWidth && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = string(ch)
			if r.MaxLines > 0 && len(lines) >= r.MaxLines {
				break
			}
		} else {
			currentLine = testLine
		}
	}
	if currentLine != "" && (r.MaxLines == 0 || len(lines) < r.MaxLines) {
		lines = append(lines, currentLine)
	}
	if len(lines) == 0 {
		lines = []string{r.Content}
	}
	return lines
}

func (r *RenderText) Paint(screen *ebiten.Image, offset Offset) {
	if r.Content == "" {
		return
	}
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	face, err := fonts.GetFontFace(fonts.FontDescriptor{
		Family: fonts.FontFamilyDefault,
		Weight: fonts.FontWeightNormal,
		Style:  fonts.FontStyleNormal,
		Size:   r.FontSize,
	})
	if err != nil || face == nil {
		return
	}

	x := float64(offset.X + bounds.X)
	y := float64(offset.Y + bounds.Y)

	lineHeight := float64(r.FontSize) * 1.2

	lines := r.splitLines(face.Face, float64(bounds.Width), yoga.MeasureModeAtMost)
	if r.MaxLines > 0 && len(lines) > r.MaxLines {
		lines = lines[:r.MaxLines]
	}
	if len(lines) == 0 {
		return
	}

	// ebiten text/v2 的 GeoM.Translate 在默认 AlignStart 下定位的是
	// 文字渲染区域的左上角（即文字块的顶部），不是基线。
	// 第一行文字顶部直接对齐 bounds.Y，文字刚好填满 bounds。
	content := strings.Join(lines, "\n")

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(r.TextColor)
	op.GeoM.Translate(x, y)
	op.LineSpacing = lineHeight
	text.Draw(screen, content, face.Face, op)

	// Draw underline
	if r.Underline {
		lastLineY := float32(y) + float32(len(lines)-1)*float32(lineHeight) + float32(lineHeight)*0.85
		vector.DrawFilledRect(screen, float32(x), lastLineY, bounds.Width, 1, r.TextColor, false)
	}
}

func (r *RenderText) SetContent(s string) {
	if r.Content == s {
		return
	}
	r.Content = s
	r.yoga.MarkDirty()
	r.MarkNeedsLayout()
}

func (r *RenderText) SetFontSize(v float32) {
	if r.FontSize == v {
		return
	}
	r.FontSize = v
	r.yoga.MarkDirty()
	r.MarkNeedsLayout()
}

func (r *RenderText) SetColor(c color.Color) {
	if r.TextColor == c {
		return
	}
	r.TextColor = c
	r.MarkNeedsPaint()
}

func (r *RenderText) SetMaxLines(v int) {
	if r.MaxLines == v {
		return
	}
	r.MaxLines = v
	r.yoga.MarkDirty()
	r.MarkNeedsLayout()
}

func (r *RenderText) SetUnderline(v bool) {
	if r.Underline == v {
		return
	}
	r.Underline = v
	r.MarkNeedsPaint()
}
