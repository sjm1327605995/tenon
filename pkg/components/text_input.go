package components

import (
	"image"
	"image/color"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// TextInput 是文本输入框组件，基于 Ebiten exp/textinput 实现，支持 IME。
type TextInput struct {
	core.BaseHost
	field            *textinput.Field
	onChange         func(string)
	onSubmit         func(string)
	placeholder      string
	isMultiline      bool

	// 样式
	textColor        color.Color
	placeholderColor color.Color
	cursorColor      color.Color
	borderColor      color.Color
	focusBorderColor color.Color
	bgColor          color.Color
	selectionColor   color.Color
	compositionColor color.Color
	fontSize         float32
	fontFace         *text.GoTextFace
	padding          float32

	// 光标闪烁
	cursorVisible    bool
	blinkTimer       float32
}

// NewTextInput 创建一个文本输入框。
func NewTextInput() *TextInput {
	ti := &TextInput{
		field:            &textinput.Field{},
		placeholder:      "",
		isMultiline:      false,
		textColor:        color.RGBA{R: 33, G: 37, B: 41, A: 255},
		placeholderColor: color.RGBA{R: 170, G: 170, B: 170, A: 255},
		cursorColor:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
		borderColor:      color.RGBA{R: 200, G: 200, B: 200, A: 255},
		focusBorderColor: color.RGBA{R: 0, G: 123, B: 255, A: 255},
		bgColor:          color.White,
		selectionColor:   color.RGBA{R: 0, G: 123, B: 255, A: 100},
		compositionColor: color.RGBA{R: 0, G: 123, B: 255, A: 255},
		fontSize:         14,
		padding:          8,
		cursorVisible:    true,
	}
	ti.Init(ti)
	ti.SetFocusable(true)
	ti.GetElement().Yoga.StyleSetHeight(36)
	return ti
}

// Draw 绘制输入框。
func (ti *TextInput) Draw(screen *ebiten.Image) {
	el := ti.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := ti.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	isFocused := ti.isFocused()

	// 背景
	ti.drawRoundedRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, 4, ti.bgColor)

	// 边框
	borderClr := ti.borderColor
	if isFocused {
		borderClr = ti.focusBorderColor
	}
	vector.StrokeLine(screen, bounds.X, bounds.Y+bounds.Height-1, bounds.X+bounds.Width, bounds.Y+bounds.Height-1, 1.5, borderClr, false)

	face := ti.getFace()
	if face == nil {
		return
	}

	textX := bounds.X + ti.padding
	textY := bounds.Y + ti.padding

	// 占位文本
	if ti.field.Text() == "" && ti.placeholder != "" && !ti.field.IsFocused() {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(textX), float64(textY))
		op.ColorScale.ScaleWithColor(ti.placeholderColor)
		text.Draw(screen, ti.placeholder, face, op)
		return
	}

	// 获取包含合成文本的渲染文本
	renderText := ti.field.TextForRendering()
	selectionStart, selectionEnd := ti.field.Selection()

	// 计算合成文本范围（字节索引）
	compStart, compEnd := -1, -1
	if l := ti.field.UncommittedTextLengthInBytes(); l > 0 {
		compStart = selectionStart
		compEnd = selectionStart + l
	}

	// 选区背景
	if selectionStart != selectionEnd {
		ti.drawSelection(screen, face, textX, textY, renderText, selectionStart, selectionEnd)
	}

	// 主文本
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(textX), float64(textY))
	op.ColorScale.ScaleWithColor(ti.textColor)
	if ti.isMultiline {
		lineH := ti.getLineHeight(face)
		op.LineSpacing = float64(lineH)
	}
	text.Draw(screen, renderText, face, op)

	// 合成文本下划线
	if compStart >= 0 && compEnd > compStart {
		ti.drawCompositionUnderline(screen, face, textX, textY, renderText, compStart, compEnd)
	}

	// 光标
	if isFocused && ti.cursorVisible {
		cx, cy := ti.getCursorPixelPos(face, textX, textY)
		vector.StrokeLine(screen, cx, cy, cx, cy+float32(face.Size), 1.5, ti.cursorColor, false)
	}
}

func (ti *TextInput) drawRoundedRect(screen *ebiten.Image, x, y, w, h, r float32, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, core.BorderRadius{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r})
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

// drawSelection 绘制选区背景。
func (ti *TextInput) drawSelection(screen *ebiten.Image, face *text.GoTextFace, textX, textY float32, renderText string, selStart, selEnd int) {
	if selStart > selEnd {
		selStart, selEnd = selEnd, selStart
	}
	if selStart < 0 {
		selStart = 0
	}
	if selEnd > len(renderText) {
		selEnd = len(renderText)
	}
	if selStart >= selEnd {
		return
	}

	lineH := ti.getLineHeight(face)
	lines := splitLines(renderText)
	var byteOffset int
	for lineIdx, line := range lines {
		lineLen := len(line)
		lineStart := byteOffset
		lineEnd := byteOffset + lineLen

		if selEnd <= lineStart || selStart >= lineEnd {
			byteOffset = lineEnd
			continue
		}

		sStart := selStart
		if sStart < lineStart {
			sStart = lineStart
		}
		sEnd := selEnd
		if sEnd > lineEnd {
			sEnd = lineEnd
		}

		beforeSel := renderText[lineStart:sStart]
		selText := renderText[sStart:sEnd]

		startW, _ := text.Measure(beforeSel, face, 0)
		selW, _ := text.Measure(selText, face, 0)

		vector.FillRect(screen, textX+float32(startW), textY+float32(lineIdx)*lineH, float32(selW), lineH, ti.selectionColor, false)

		byteOffset = lineEnd
	}
}

// drawCompositionUnderline 绘制合成文本下划线。
func (ti *TextInput) drawCompositionUnderline(screen *ebiten.Image, face *text.GoTextFace, textX, textY float32, renderText string, compStart, compEnd int) {
	if compStart > compEnd {
		compStart, compEnd = compEnd, compStart
	}
	if compStart < 0 {
		compStart = 0
	}
	if compEnd > len(renderText) {
		compEnd = len(renderText)
	}

	lines := splitLines(renderText)
	var byteOffset int
	for lineIdx, line := range lines {
		lineLen := len(line)
		lineStart := byteOffset
		lineEnd := byteOffset + lineLen

		if compEnd <= lineStart || compStart >= lineEnd {
			byteOffset = lineEnd
			continue
		}

		cStart := compStart
		if cStart < lineStart {
			cStart = lineStart
		}
		cEnd := compEnd
		if cEnd > lineEnd {
			cEnd = lineEnd
		}

		beforeComp := renderText[lineStart:cStart]
		compText := renderText[cStart:cEnd]

		startW, _ := text.Measure(beforeComp, face, 0)
		compW, _ := text.Measure(compText, face, 0)

		lineH := ti.getLineHeight(face)
		y := textY + float32(lineIdx)*lineH + lineH - 2
		vector.StrokeLine(screen, textX+float32(startW), y, textX+float32(startW)+float32(compW), y, 1.5, ti.compositionColor, false)

		byteOffset = lineEnd
	}
}

// getCursorPixelPos 计算光标像素位置。
func (ti *TextInput) getCursorPixelPos(face *text.GoTextFace, textX, textY float32) (float32, float32) {
	renderText := ti.field.TextForRendering()
	selectionStart, _ := ti.field.Selection()
	if selectionStart > len(renderText) {
		selectionStart = len(renderText)
	}
	if selectionStart < 0 {
		selectionStart = 0
	}

	beforeCursor := renderText[:selectionStart]
	lineH := ti.getLineHeight(face)
	lines := splitLines(beforeCursor)
	lastLine := ""
	if len(lines) > 0 {
		lastLine = lines[len(lines)-1]
	}

	w, _ := text.Measure(lastLine, face, 0)
	cursorX := textX + float32(w)
	cursorY := textY + float32(len(lines)-1)*lineH
	return cursorX, cursorY
}

func (ti *TextInput) getLineHeight(face *text.GoTextFace) float32 {
	if face == nil {
		return ti.fontSize * 1.2
	}
	m := face.Metrics()
	h := float32(m.HAscent + m.HDescent)
	if h <= 0 {
		return float32(face.Size) * 1.2
	}
	return h
}

func (ti *TextInput) getFace() *text.GoTextFace {
	if ti.fontFace != nil {
		return ti.fontFace
	}
	fontManager := fonts.GetFontManager()
	face, err := fontManager.GetDefaultFontFace(ti.fontSize)
	if err == nil && face != nil && face.Face != nil {
		ti.fontFace = face.Face
		return face.Face
	}
	return nil
}

func (ti *TextInput) isFocused() bool {
	eng := ti.GetEngine()
	if eng == nil {
		return false
	}
	return eng.GetFocusHost() == ti
}

// Update 处理输入、光标闪烁和键盘事件。
func (ti *TextInput) Update() error {
	if !ti.isFocused() {
		if ti.field.IsFocused() {
			ti.field.Blur()
		}
		ti.cursorVisible = false
		ti.blinkTimer = 0
		return nil
	}

	// 聚焦 Field（启动 IME）
	if !ti.field.IsFocused() {
		ti.field.Focus()
	}

	// 处理文本输入（含 IME）
	bounds := ti.GetLayoutBounds()
	textX := int(bounds.X + ti.padding)
	textY := int(bounds.Y + ti.padding)
	lineHeight := int(ti.getLineHeight(ti.getFace()))

	handled, err := ti.field.HandleInputWithBounds(image.Rect(textX, textY, textX+1, textY+lineHeight))
	if err != nil {
		return err
	}

	// 如果 IME 处理了输入，重置光标闪烁并触发 onChange
	if handled {
		ti.cursorVisible = true
		ti.blinkTimer = 0
		if ti.onChange != nil {
			ti.onChange(ti.field.Text())
		}
	}

	// 光标闪烁（IME 合成时不闪烁，保持可见）
	if ti.field.UncommittedTextLengthInBytes() > 0 {
		ti.cursorVisible = true
		ti.blinkTimer = 0
	} else {
		ti.blinkTimer += 1.0 / 60.0
		if ti.blinkTimer >= 0.53 {
			ti.cursorVisible = !ti.cursorVisible
			ti.blinkTimer = 0
		}
	}

	// IME 未处理时，手动处理键盘事件
	if !handled {
		ti.handleKeyboardInput()
	}

	return nil
}

func (ti *TextInput) handleKeyboardInput() {
	// Backspace（支持长按重复）
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) || inpututil.KeyPressDuration(ebiten.KeyBackspace) > 30 {
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) || inpututil.KeyPressDuration(ebiten.KeyBackspace)%6 == 0 {
			ti.handleBackspace()
		}
	}

	// Delete（支持长按重复）
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) || inpututil.KeyPressDuration(ebiten.KeyDelete) > 30 {
		if inpututil.IsKeyJustPressed(ebiten.KeyDelete) || inpututil.KeyPressDuration(ebiten.KeyDelete)%6 == 0 {
			ti.handleDelete()
		}
	}

	// 左方向键
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.KeyPressDuration(ebiten.KeyLeft) > 30 {
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.KeyPressDuration(ebiten.KeyLeft)%6 == 0 {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				ti.extendSelection(-1)
			} else {
				ti.clearSelection()
				ti.moveCursor(-1)
			}
		}
	}

	// 右方向键
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.KeyPressDuration(ebiten.KeyRight) > 30 {
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.KeyPressDuration(ebiten.KeyRight)%6 == 0 {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				ti.extendSelection(1)
			} else {
				ti.clearSelection()
				ti.moveCursor(1)
			}
		}
	}

	// Enter
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if ti.isMultiline {
			ti.insertRune('\n')
		} else if ti.onSubmit != nil {
			ti.onSubmit(ti.field.Text())
		}
	}

	// Home / End
	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		ti.clearSelection()
		ti.field.SetSelection(0, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		text := ti.field.Text()
		ti.clearSelection()
		ti.field.SetSelection(len(text), len(text))
	}

	// Ctrl+A 全选
	if inpututil.IsKeyJustPressed(ebiten.KeyA) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		text := ti.field.Text()
		ti.field.SetSelection(0, len(text))
	}
}

func (ti *TextInput) insertRune(ch rune) {
	text := ti.field.Text()
	start, end := ti.field.Selection()
	if start != end {
		text = text[:start] + text[end:]
		end = start
	}
	text = text[:end] + string(ch) + text[end:]
	start += utf8.RuneLen(ch)
	ti.field.SetTextAndSelection(text, start, start)
	ti.cursorVisible = true
	ti.blinkTimer = 0
	if ti.onChange != nil {
		ti.onChange(text)
	}
}

func (ti *TextInput) handleBackspace() {
	text := ti.field.Text()
	start, end := ti.field.Selection()
	if start != end {
		text = text[:start] + text[end:]
		ti.field.SetTextAndSelection(text, start, start)
	} else if start > 0 {
		_, l := utf8.DecodeLastRuneInString(text[:start])
		text = text[:start-l] + text[end:]
		start -= l
		ti.field.SetTextAndSelection(text, start, start)
	}
	ti.cursorVisible = true
	ti.blinkTimer = 0
	if ti.onChange != nil {
		ti.onChange(text)
	}
}

func (ti *TextInput) handleDelete() {
	text := ti.field.Text()
	start, end := ti.field.Selection()
	if start != end {
		text = text[:start] + text[end:]
		ti.field.SetTextAndSelection(text, start, start)
	} else if end < len(text) {
		_, l := utf8.DecodeRuneInString(text[end:])
		text = text[:end] + text[end+l:]
		ti.field.SetTextAndSelection(text, end, end)
	}
	ti.cursorVisible = true
	ti.blinkTimer = 0
	if ti.onChange != nil {
		ti.onChange(text)
	}
}

func (ti *TextInput) moveCursor(delta int) {
	text := ti.field.Text()
	start, _ := ti.field.Selection()
	if delta < 0 {
		if start > 0 {
			_, l := utf8.DecodeLastRuneInString(text[:start])
			start -= l
		}
	} else {
		if start < len(text) {
			_, l := utf8.DecodeRuneInString(text[start:])
			start += l
		}
	}
	ti.field.SetSelection(start, start)
	ti.cursorVisible = true
	ti.blinkTimer = 0
}

func (ti *TextInput) extendSelection(delta int) {
	start, end := ti.field.Selection()
	if start == end {
		start = end
	}
	text := ti.field.Text()
	if delta < 0 {
		if end > 0 {
			_, l := utf8.DecodeLastRuneInString(text[:end])
			end -= l
		}
	} else {
		if end < len(text) {
			_, l := utf8.DecodeRuneInString(text[end:])
			end += l
		}
	}
	ti.field.SetSelection(start, end)
}

func (ti *TextInput) clearSelection() {
	start, end := ti.field.Selection()
	if start != end {
		ti.field.SetSelection(start, start)
	}
}

// HandleEvent 处理点击事件，聚焦/失焦并设置鼠标光标。
func (ti *TextInput) HandleEvent(e *core.Event) bool {
	bounds := ti.GetLayoutBounds()
	switch e.Type {
	case core.EventMouseDown:
		// 设置鼠标光标为文本形状
		ebiten.SetCursorShape(ebiten.CursorShapeText)
		return true
	case core.EventClick:
		if e.X >= bounds.X && e.X <= bounds.X+bounds.Width &&
			e.Y >= bounds.Y && e.Y <= bounds.Y+bounds.Height {
			// 点击输入框内部聚焦已在 Engine 中处理
		}
		return true
	case core.EventMouseMove:
		// 鼠标悬停时设置文本光标
		if e.X >= bounds.X && e.X <= bounds.X+bounds.Width &&
			e.Y >= bounds.Y && e.Y <= bounds.Y+bounds.Height {
			ebiten.SetCursorShape(ebiten.CursorShapeText)
		}
		return true
	}
	return false
}

// ==================== 链式 API ====================

func (ti *TextInput) SetText(text string) *TextInput {
	ti.field.SetTextAndSelection(text, len(text), len(text))
	return ti
}
func (ti *TextInput) SetPlaceholder(p string) *TextInput {
	ti.placeholder = p
	return ti
}
func (ti *TextInput) SetMultiline(v bool) *TextInput {
	ti.isMultiline = v
	return ti
}
func (ti *TextInput) SetOnChange(fn func(string)) *TextInput {
	ti.onChange = fn
	return ti
}
func (ti *TextInput) SetOnSubmit(fn func(string)) *TextInput {
	ti.onSubmit = fn
	return ti
}
func (ti *TextInput) SetFontSize(size float32) *TextInput {
	ti.fontSize = size
	ti.fontFace = nil
	return ti
}
func (ti *TextInput) SetTextColor(clr color.Color) *TextInput {
	ti.textColor = clr
	return ti
}
func (ti *TextInput) SetPlaceholderColor(clr color.Color) *TextInput {
	ti.placeholderColor = clr
	return ti
}
func (ti *TextInput) SetBorderColor(clr color.Color) *TextInput {
	ti.borderColor = clr
	return ti
}
func (ti *TextInput) SetFocusBorderColor(clr color.Color) *TextInput {
	ti.focusBorderColor = clr
	return ti
}
func (ti *TextInput) SetBackgroundColor(clr color.Color) *TextInput {
	ti.bgColor = clr
	return ti
}
func (ti *TextInput) SetWidth(width float32) *TextInput {
	ti.GetElement().Yoga.StyleSetWidth(width)
	return ti
}
func (ti *TextInput) SetHeight(height float32) *TextInput {
	ti.GetElement().Yoga.StyleSetHeight(height)
	return ti
}
func (ti *TextInput) SetPadding(padding float32) *TextInput {
	ti.padding = padding
	return ti
}
func (ti *TextInput) SetMargin(edge yoga.Edge, value float32) *TextInput {
	ti.GetElement().Yoga.StyleSetMargin(edge, value)
	return ti
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, r := range s {
		if r == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	lines = append(lines, s[start:])
	return lines
}

// SyncFrom 同步文本输入框样式（保留输入状态）。
func (ti *TextInput) SyncFrom(other core.Host) {
	if o, ok := other.(*TextInput); ok {
		ti.onChange = o.onChange
		ti.onSubmit = o.onSubmit
		ti.placeholder = o.placeholder
		ti.isMultiline = o.isMultiline
		ti.textColor = o.textColor
		ti.placeholderColor = o.placeholderColor
		ti.cursorColor = o.cursorColor
		ti.borderColor = o.borderColor
		ti.focusBorderColor = o.focusBorderColor
		ti.bgColor = o.bgColor
		ti.selectionColor = o.selectionColor
		ti.compositionColor = o.compositionColor
		ti.fontSize = o.fontSize
		ti.fontFace = o.fontFace
		ti.padding = o.padding
	}
}
