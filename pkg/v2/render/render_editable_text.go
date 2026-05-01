package render

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderEditableText 是可编辑文本的 RenderObject。
// 继承 RenderText，增加光标、焦点和键盘输入处理能力。
type RenderEditableText struct {
	RenderText
	cursorPos   int
	focused     bool
	lastBlink   time.Time
	showCursor  bool
	onChanged   func(string)
	onSubmitted func(string)
}

func NewRenderEditableText() *RenderEditableText {
	r := &RenderEditableText{
		showCursor: true,
	}
	r.BaseRenderObject.Init(r)
	r.yoga = yoga.NewNode()
	r.yoga.SetMeasureFunc(r.measure)
	r.TextColor = color.Black
	r.FontSize = 14
	return r
}

func (r *RenderEditableText) Focus() {
	r.focused = true
	r.showCursor = true
	r.lastBlink = time.Now()
	r.MarkNeedsPaint()
}

func (r *RenderEditableText) Blur() {
	r.focused = false
	r.MarkNeedsPaint()
}

func (r *RenderEditableText) IsFocused() bool {
	return r.focused
}

func (r *RenderEditableText) TickBlink() {
	if !r.focused {
		return
	}
	if time.Since(r.lastBlink) > 530*time.Millisecond {
		r.showCursor = !r.showCursor
		r.lastBlink = time.Now()
		r.MarkNeedsPaint()
	}
}

func (r *RenderEditableText) SetOnChanged(fn func(string)) {
	r.onChanged = fn
}

func (r *RenderEditableText) SetOnSubmitted(fn func(string)) {
	r.onSubmitted = fn
}

// HandleInput 处理键盘输入。
func (r *RenderEditableText) HandleInput(
	chars []rune,
	backspace, enter, left, right, home, end, selectAll bool,
) {
	content := r.Content
	runes := []rune(content)
	changed := false

	for _, ch := range chars {
		if r.cursorPos >= len(runes) {
			runes = append(runes, ch)
		} else {
			runes = append(runes[:r.cursorPos], append([]rune{ch}, runes[r.cursorPos:]...)...)
		}
		r.cursorPos++
		changed = true
	}

	if backspace && r.cursorPos > 0 {
		r.cursorPos--
		runes = append(runes[:r.cursorPos], runes[r.cursorPos+1:]...)
		changed = true
	}

	if left && r.cursorPos > 0 {
		r.cursorPos--
		r.MarkNeedsPaint()
	}
	if right && r.cursorPos < len(runes) {
		r.cursorPos++
		r.MarkNeedsPaint()
	}

	if home {
		r.cursorPos = 0
		r.MarkNeedsPaint()
	}
	if end {
		r.cursorPos = len(runes)
		r.MarkNeedsPaint()
	}

	if selectAll {
		r.cursorPos = len(runes)
		r.MarkNeedsPaint()
	}

	if changed {
		newContent := string(runes)
		if newContent != content {
			r.SetContent(newContent)
			if r.onChanged != nil {
				r.onChanged(newContent)
			}
		}
	}

	if enter && r.onSubmitted != nil {
		r.onSubmitted(string(runes))
	}
}

func (r *RenderEditableText) Paint(screen *ebiten.Image, offset Offset) {
	// 先绘制文本
	r.RenderText.Paint(screen, offset)

	// 绘制光标
	if r.focused && r.showCursor {
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

		runes := []rune(r.Content)
		var prefix string
		if r.cursorPos > 0 && r.cursorPos <= len(runes) {
			prefix = string(runes[:r.cursorPos])
		}

		w, _ := text.Measure(prefix, face.Face, 0)
		lineHeight := float64(r.FontSize) * 1.2

		x := float64(offset.X+bounds.X) + w
		y := float64(offset.Y+bounds.Y)

		cursorH := lineHeight
		if cursorH > float64(bounds.Height) {
			cursorH = float64(bounds.Height)
		}

		vector.StrokeLine(
			screen,
			float32(x), float32(y),
			float32(x), float32(y+cursorH),
			1,
			color.Black,
			false,
		)
	}
}
