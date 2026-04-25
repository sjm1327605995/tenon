package components

import (
	"image"
	"image/color"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// TextInput is a text input field with IME support.
type TextInput struct {
	core.BaseElement
	field            *textinput.Field
	onChange         func(string)
	onSubmit         func(string)
	placeholder      string
	isMultiline      bool

	// Styles
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

	// Cursor blink
	cursorVisible    bool
	blinkTimer       float32

	// Mouse selection
	selecting        bool
	selectAnchor     int
}

// NewTextInput creates a text input.
func NewTextInput() *TextInput {
	theme := core.GetTheme()
	ti := &TextInput{
		field:            &textinput.Field{},
		placeholder:      "",
		isMultiline:      false,
		textColor:        theme.InputTextColor,
		placeholderColor: theme.InputPlaceholderColor,
		cursorColor:      theme.InputTextColor,
		borderColor:      theme.InputBorderColor,
		focusBorderColor: theme.InputFocusBorderColor,
		bgColor:          theme.InputBgColor,
		selectionColor:   theme.InputSelectionColor,
		compositionColor: theme.PrimaryColor,
		fontSize:         theme.FontSizeBase,
		padding:          8,
		cursorVisible:    true,
	}
	ti.Init(ti)
	ti.SetFlag(core.FlagFocusable)
	ti.SetHeight(36)
	return ti
}

// ElementType returns type identifier.
func (ti *TextInput) ElementType() string { return "TextInput" }

// Draw renders the input field.
func (ti *TextInput) Draw(screen *ebiten.Image) {
	if !ti.IsVisible() {
		return
	}
	bounds := ti.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	isFocused := ti.isFocused()

	// Background
	br := core.GetTheme().InputBorderRadius
	drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, core.BorderRadius{TopLeft: br, TopRight: br, BottomRight: br, BottomLeft: br}, ti.bgColor)

	// Border
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

	// Placeholder
	if ti.field.Text() == "" && ti.placeholder != "" && !ti.field.IsFocused() {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(textX), float64(textY))
		op.ColorScale.ScaleWithColor(ti.placeholderColor)
		text.Draw(screen, ti.placeholder, face, op)
		return
	}

	renderText := ti.field.TextForRendering()
	selectionStart, selectionEnd := ti.field.Selection()

	// Composition range
	compStart, compEnd := -1, -1
	if l := ti.field.UncommittedTextLengthInBytes(); l > 0 {
		compStart = selectionStart
		compEnd = selectionStart + l
	}

	// Selection background
	if selectionStart != selectionEnd {
		ti.drawSelection(screen, face, textX, textY, renderText, selectionStart, selectionEnd)
	}

	// Main text
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(textX), float64(textY))
	op.ColorScale.ScaleWithColor(ti.textColor)
	if ti.isMultiline {
		lineH := ti.getLineHeight(face)
		op.LineSpacing = float64(lineH)
	}
	text.Draw(screen, renderText, face, op)

	// Composition underline
	if compStart >= 0 && compEnd > compStart {
		ti.drawCompositionUnderline(screen, face, textX, textY, renderText, compStart, compEnd)
	}

	// Cursor
	if isFocused && ti.cursorVisible {
		cx, cy := ti.getCursorPixelPos(face, textX, textY)
		vector.StrokeLine(screen, cx, cy, cx, cy+float32(face.Size), 1.5, ti.cursorColor, false)
	}
}

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
	return eng.GetFocusTarget() == ti
}

// Update handles input, cursor blink and keyboard events.
func (ti *TextInput) Update() error {
	if !ti.isFocused() {
		if ti.field.IsFocused() {
			ti.field.Blur()
		}
		ti.cursorVisible = false
		ti.blinkTimer = 0
		return nil
	}

	if !ti.field.IsFocused() {
		ti.field.Focus()
	}

	// Handle IME input
	bounds := ti.GetBounds()
	textX := int(bounds.X + ti.padding)
	textY := int(bounds.Y + ti.padding)
	lineHeight := int(ti.getLineHeight(ti.getFace()))
	handled, err := ti.field.HandleInputWithBounds(image.Rect(textX, textY, textX+1, textY+lineHeight))
	if err != nil {
		return err
	}
	if handled {
		ti.cursorVisible = true
		ti.blinkTimer = 0
		if ti.onChange != nil {
			ti.onChange(ti.field.Text())
		}
	}

	// Cursor blink
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
	return nil
}

// HandleEvent processes click, selection and keyboard events.
func (ti *TextInput) HandleEvent(e *core.Event) bool {
	bounds := ti.GetBounds()
	switch e.Type {
	case core.EventFocusIn:
		ti.field.Focus()
		return true
	case core.EventFocusOut:
		ti.field.Blur()
		return true
	case core.EventMouseDown:
		if e.X >= bounds.X && e.X <= bounds.X+bounds.Width &&
			e.Y >= bounds.Y && e.Y <= bounds.Y+bounds.Height {
			ti.selecting = true
			pos := ti.hitTestText(e.X, e.Y)
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				start, end := ti.field.Selection()
				if start == end {
					ti.selectAnchor = start
				} else if absInt(pos-start) < absInt(pos-end) {
					ti.selectAnchor = end
				} else {
					ti.selectAnchor = start
				}
				ti.field.SetSelection(ti.selectAnchor, pos)
			} else {
				ti.selectAnchor = pos
				ti.field.SetSelection(pos, pos)
			}
			ti.cursorVisible = true
			ti.blinkTimer = 0
		}
		return true
	case core.EventMouseMove:
		if ti.selecting {
			pos := ti.hitTestText(e.X, e.Y)
			ti.field.SetSelection(ti.selectAnchor, pos)
		}
		return true
	case core.EventMouseUp:
		if ti.selecting {
			start, end := ti.field.Selection()
			if start == end {
				ti.field.SetSelection(start, start)
			}
			ti.selecting = false
		}
		return true
	case core.EventClick:
		return true
	case core.EventKeyDown:
		return ti.handleKeyDown(e.Key)
	}
	return false
}

func (ti *TextInput) handleKeyDown(key ebiten.Key) bool {
	switch key {
	case ebiten.KeyBackspace:
		ti.handleBackspace()
		return true
	case ebiten.KeyDelete:
		ti.handleDelete()
		return true
	case ebiten.KeyLeft:
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			ti.extendSelection(-1)
		} else {
			ti.clearSelection()
			ti.moveCursor(-1)
		}
		return true
	case ebiten.KeyRight:
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			ti.extendSelection(1)
		} else {
			ti.clearSelection()
			ti.moveCursor(1)
		}
		return true
	case ebiten.KeyEnter:
		if ti.isMultiline {
			ti.insertRune('\n')
		} else if ti.onSubmit != nil {
			ti.onSubmit(ti.field.Text())
		}
		return true
	case ebiten.KeyHome:
		ti.clearSelection()
		ti.field.SetSelection(0, 0)
		return true
	case ebiten.KeyEnd:
		text := ti.field.Text()
		ti.clearSelection()
		ti.field.SetSelection(len(text), len(text))
		return true
	}
	if key == ebiten.KeyA && ebiten.IsKeyPressed(ebiten.KeyControl) {
		text := ti.field.Text()
		ti.field.SetSelection(0, len(text))
		return true
	}
	// Character input handled by textinput.Field via HandleInputWithBounds in Update
	return false
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

func (ti *TextInput) hitTestText(mouseX, mouseY float32) int {
	face := ti.getFace()
	if face == nil {
		return 0
	}
	bounds := ti.GetBounds()
	textX := bounds.X + ti.padding
	textY := bounds.Y + ti.padding

	renderText := ti.field.TextForRendering()
	lineH := ti.getLineHeight(face)
	lines := splitLines(renderText)

	relY := mouseY - textY
	lineIdx := int(relY / lineH)
	if lineIdx < 0 {
		lineIdx = 0
	}

	byteOffset := 0
	for i := 0; i < lineIdx && i < len(lines); i++ {
		byteOffset += len(lines[i]) + 1
	}

	if lineIdx >= len(lines) {
		return len(renderText)
	}

	line := lines[lineIdx]
	relX := mouseX - textX
	if relX <= 0 {
		return byteOffset
	}

	bestIdx := len(line)
	bestDist := relX
	low, high := 0, len(line)
	for low < high {
		mid := (low + high) / 2
		w, _ := text.Measure(line[:mid], face, 0)
		w32 := float32(w)
		dist := absFloat32(w32 - relX)
		if dist < bestDist {
			bestDist = dist
			bestIdx = mid
		}
		if w32 < relX {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return byteOffset + bestIdx
}

// ==================== Chain API ====================

func (ti *TextInput) SetText(text string) *TextInput {
	ti.field.SetTextAndSelection(text, len(text), len(text))
	ti.Mark(core.FlagNeedDraw)
	return ti
}

func (ti *TextInput) SetPlaceholder(p string) *TextInput {
	ti.placeholder = p
	ti.Mark(core.FlagNeedDraw)
	return ti
}

func (ti *TextInput) SetMultiline(v bool) *TextInput {
	ti.isMultiline = v
	ti.Mark(core.FlagNeedDraw)
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
	ti.Mark(core.FlagNeedMeasure | core.FlagNeedDraw)
	return ti
}

func (ti *TextInput) SetTextColor(clr color.Color) *TextInput {
	ti.textColor = clr
	ti.Mark(core.FlagNeedDraw)
	return ti
}

func (ti *TextInput) SetPlaceholderColor(clr color.Color) *TextInput {
	ti.placeholderColor = clr
	ti.Mark(core.FlagNeedDraw)
	return ti
}

func (ti *TextInput) SetBorderColor(clr color.Color) *TextInput {
	ti.borderColor = clr
	ti.Mark(core.FlagNeedDraw)
	return ti
}

func (ti *TextInput) SetFocusBorderColor(clr color.Color) *TextInput {
	ti.focusBorderColor = clr
	ti.Mark(core.FlagNeedDraw)
	return ti
}

func (ti *TextInput) SetBackgroundColor(clr color.Color) *TextInput {
	ti.bgColor = clr
	ti.Mark(core.FlagNeedDraw)
	return ti
}

func (ti *TextInput) SetInnerPadding(padding float32) *TextInput {
	ti.padding = padding
	ti.Mark(core.FlagNeedDraw)
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

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func absFloat32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}
