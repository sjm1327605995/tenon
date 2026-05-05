package render

import (
	"image"
	"image/color"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderEditableText 是可编辑文本的 RenderObject。
// 继承 RenderText，增加光标、焦点和键盘输入处理能力。
// 内部自动管理滚动偏移，确保光标始终可见（类似 <input> / <textarea> 的溢出滚动行为）。
type RenderEditableText struct {
	RenderText
	cursorPos        int
	focused          bool
	lastBlink        time.Time
	showCursor       bool
	placeholder      string
	placeholderColor color.Color
	multiline        bool
	onChanged        func(string)
	onSubmitted      func(string)

	// 选区
	selectionStart int // 选区起始位置（-1 表示无选区）
	selectionEnd   int // 选区结束位置

	// textinput.Field 提供 IME 支持（中文/日文/韩文输入法）
	field *textinput.Field

	// 内部滚动偏移（仅 RenderEditableText 自己使用，不依赖 RenderScroll）
	scrollOffsetX float32
	scrollOffsetY float32
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
	r.MaxLines = 1 // 默认可编辑文本为单行模式
	return r
}

func (r *RenderEditableText) measure(
	node *yoga.Node,
	width float32,
	widthMode yoga.MeasureMode,
	height float32,
	heightMode yoga.MeasureMode,
) yoga.Size {
	size := r.RenderText.measure(node, width, widthMode, height, heightMode)
	// EditableText 不应被 measure 为 0 尺寸，否则无法接收点击和绘制光标
	minW := float32(r.FontSize) * 0.6
	if size.Width <= 0 {
		if widthMode == yoga.MeasureModeExactly && width > 0 {
			size.Width = width
		} else if widthMode == yoga.MeasureModeAtMost && width > 0 {
			if minW < width {
				size.Width = minW
			} else {
				size.Width = width
			}
		} else {
			size.Width = minW
		}
	}
	if size.Height <= 0 {
		minH := r.FontSize * 1.2
		if heightMode == yoga.MeasureModeExactly && height > 0 {
			size.Height = height
		} else if heightMode == yoga.MeasureModeAtMost && height > 0 {
			if minH < height {
				size.Height = minH
			} else {
				size.Height = height
			}
		} else {
			size.Height = minH
		}
	}
	return size
}

func (r *RenderEditableText) HitTest(x, y float32) bool {
	if !r.visible {
		return false
	}
	return r.bounds.Contains(x, y)
}

func (r *RenderEditableText) Focus() {
	r.focused = true
	r.showCursor = true
	r.lastBlink = time.Now()
	r.MarkNeedsPaint()
	if p := r.GetParent(); p != nil {
		p.MarkNeedsPaint()
	}

	// 启动 textinput 会话以支持 IME
	if r.field == nil {
		r.field = &textinput.Field{}
	}
	r.field.Focus()
	r.syncToField()
	r.ensureCursorVisible()
}

func (r *RenderEditableText) Blur() {
	r.focused = false
	r.MarkNeedsPaint()
	if p := r.GetParent(); p != nil {
		p.MarkNeedsPaint()
	}
	if r.field != nil {
		r.field.Blur()
	}
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

func (r *RenderEditableText) GetPlaceholder() string {
	return r.placeholder
}

func (r *RenderEditableText) SetPlaceholder(text string) {
	if r.placeholder == text {
		return
	}
	r.placeholder = text
	r.MarkNeedsPaint()
}

func (r *RenderEditableText) GetPlaceholderColor() color.Color {
	return r.placeholderColor
}

func (r *RenderEditableText) SetPlaceholderColor(c color.Color) {
	if r.placeholderColor == c {
		return
	}
	r.placeholderColor = c
	r.MarkNeedsPaint()
}

func (r *RenderEditableText) SetMultiline(v bool) {
	if r.multiline == v {
		return
	}
	r.multiline = v
	if v {
		r.RenderText.MaxLines = 0
	} else {
		r.RenderText.MaxLines = 1
	}
	r.yoga.MarkDirty()
	r.MarkNeedsLayout()
}

func (r *RenderEditableText) GetField() *textinput.Field {
	return r.field
}

func (r *RenderEditableText) GetContent() string {
	return r.Content
}

func (r *RenderEditableText) GetCursorPos() int {
	return r.cursorPos
}

func (r *RenderEditableText) SetCursorPos(pos int) {
	r.cursorPos = pos
	if r.field != nil {
		startBytes := runeIndexToByteIndex(r.Content, r.cursorPos)
		r.field.SetTextAndSelection(r.Content, startBytes, startBytes)
	}
	r.ensureCursorVisible()
	r.MarkNeedsPaint()
}

func (r *RenderEditableText) SetContent(s string) {
	r.RenderText.SetContent(s)
}

func (r *RenderEditableText) SetOnChanged(fn func(string)) {
	r.onChanged = fn
}

func (r *RenderEditableText) SetOnSubmitted(fn func(string)) {
	r.onSubmitted = fn
}

// syncToField 将 RenderEditableText 的当前状态同步到 textinput.Field。
func (r *RenderEditableText) syncToField() {
	if r.field == nil {
		return
	}
	startBytes := runeIndexToByteIndex(r.Content, r.cursorPos)
	r.field.SetTextAndSelection(r.Content, startBytes, startBytes)
}

// UpdateTextInput 每帧调用，处理 IME 输入。
// absX, absY 是 RenderEditableText 在屏幕上的绝对坐标。
// 返回 true 表示 textinput 已处理输入，引擎应跳过传统键盘输入。
func (r *RenderEditableText) UpdateTextInput(absX, absY int) bool {
	if r.field == nil || !r.field.IsFocused() {
		return false
	}

	bounds := r.bounds
	rect := image.Rect(
		int(absX), int(absY),
		int(absX)+int(bounds.Width), int(absY)+int(bounds.Height),
	)

	handled, err := r.field.HandleInputWithBounds(rect)
	if err != nil {
		return false
	}
	if !handled {
		return false
	}

	// 同步 field 的文本回 RenderEditableText
	// 注意：必须先读取 field 的 selection，再调用 SetContent，
	// 避免 SetContent 内部的行为覆盖 field 状态。
	newText := r.field.Text()
	startBytes, _ := r.field.Selection()
	if newText != r.Content {
		r.SetContent(newText)
		if r.onChanged != nil {
			r.onChanged(newText)
		}
	}

	// 同步光标位置（字节 -> rune 索引）
	r.cursorPos = byteIndexToRuneIndex(newText, startBytes)
	r.ensureCursorVisible()
	r.MarkNeedsPaint()
	return true
}

// ==================== 选区操作 ====================

// HasSelection 是否有选区。
func (r *RenderEditableText) HasSelection() bool {
	return r.selectionStart >= 0 && r.selectionStart != r.selectionEnd
}

// GetSelection 返回选区的起止位置（start, end），无选区返回 (-1, -1)。
func (r *RenderEditableText) GetSelection() (int, int) {
	if !r.HasSelection() {
		return -1, -1
	}
	if r.selectionStart <= r.selectionEnd {
		return r.selectionStart, r.selectionEnd
	}
	return r.selectionEnd, r.selectionStart
}

// SetSelection 设置选区。
func (r *RenderEditableText) SetSelection(start, end int) {
	r.selectionStart = start
	r.selectionEnd = end
}

// ClearSelection 清除选区。
func (r *RenderEditableText) ClearSelection() {
	r.selectionStart = -1
	r.selectionEnd = -1
}

// SelectAll 全选。
func (r *RenderEditableText) SelectAll() {
	r.selectionStart = 0
	r.selectionEnd = len([]rune(r.Content))
}

// SelectedText 获取选中的文本。
func (r *RenderEditableText) SelectedText() string {
	if !r.HasSelection() {
		return ""
	}
	runes := []rune(r.Content)
	s, e := r.GetSelection()
	if s < 0 || e > len(runes) {
		return ""
	}
	return string(runes[s:e])
}

// DeleteSelection 删除选区内容。
func (r *RenderEditableText) DeleteSelection() bool {
	if !r.HasSelection() {
		return false
	}
	runes := []rune(r.Content)
	s, e := r.GetSelection()
	r.Content = string(append(runes[:s], runes[e:]...))
	r.cursorPos = s
	r.ClearSelection()
	if r.onChanged != nil {
		r.onChanged(r.Content)
	}
	r.MarkNeedsPaint()
	return true
}

// ==================== 剪贴板操作 ====================

// Copy 将选中文本复制到剪贴板（返回选中的文本）。
func (r *RenderEditableText) Copy(clipboardWrite func(string)) string {
	selected := r.SelectedText()
	if selected != "" && clipboardWrite != nil {
		clipboardWrite(selected)
	}
	return selected
}

// Cut 剪切选中文本到剪贴板。
func (r *RenderEditableText) Cut(clipboardWrite func(string)) string {
	selected := r.SelectedText()
	if selected != "" && clipboardWrite != nil {
		clipboardWrite(selected)
		r.DeleteSelection()
	}
	return selected
}

// Paste 从剪贴板粘贴文本。
func (r *RenderEditableText) Paste(clipboardRead func() string) {
	if clipboardRead == nil {
		return
	}
	text := clipboardRead()
	if text == "" {
		return
	}
	// 先删除选区
	r.DeleteSelection()
	// 插入文本
	runes := []rune(r.Content)
	insertRunes := []rune(text)
	if r.cursorPos >= len(runes) {
		runes = append(runes, insertRunes...)
	} else {
		newRunes := make([]rune, 0, len(runes)+len(insertRunes))
		newRunes = append(newRunes, runes[:r.cursorPos]...)
		newRunes = append(newRunes, insertRunes...)
		newRunes = append(newRunes, runes[r.cursorPos:]...)
		runes = newRunes
	}
	r.cursorPos += len(insertRunes)
	r.Content = string(runes)
	if r.onChanged != nil {
		r.onChanged(r.Content)
	}
	r.MarkNeedsPaint()
}

// HandleInput 处理键盘输入（传统方式，不支持 IME。
// 当 textinput 可用时，引擎会优先调用 UpdateTextInput，而不是 HandleInput。
func (r *RenderEditableText) HandleInput(
	chars []rune,
	backspace, enter, left, right, home, end, selectAll bool,
) {
	content := r.Content
	runes := []rune(content)
	changed := false
	cursorMoved := false

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
		cursorMoved = true
	}
	if right && r.cursorPos < len(runes) {
		r.cursorPos++
		cursorMoved = true
	}
	if home {
		r.cursorPos = 0
		cursorMoved = true
	}
	if end {
		r.cursorPos = len(runes)
		cursorMoved = true
	}
	if selectAll {
		r.cursorPos = len(runes)
		cursorMoved = true
	}

	if enter {
		if r.multiline {
			if r.cursorPos >= len(runes) {
				runes = append(runes, '\n')
			} else {
				runes = append(runes[:r.cursorPos], append([]rune{'\n'}, runes[r.cursorPos:]...)...)
			}
			r.cursorPos++
			changed = true
		} else if r.onSubmitted != nil {
			r.onSubmitted(string(runes))
		}
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

	if changed || cursorMoved || enter {
		if r.field != nil {
			r.syncToField()
		}
		r.ensureCursorVisible()
		r.MarkNeedsPaint()
	}
}

// ensureCursorVisible 根据光标位置自动调整滚动偏移，确保光标始终在可视区域内。
func (r *RenderEditableText) ensureCursorVisible() {
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

	displayText := r.Content
	if r.field != nil && r.field.IsFocused() {
		displayText = r.field.TextForRendering()
	}

	cursorText := displayText
	if r.field != nil && r.field.IsFocused() {
		startBytes, _ := r.field.Selection()
		uncommittedLen := r.field.UncommittedTextLengthInBytes()
		if uncommittedLen > 0 {
			cursorBytes := startBytes + uncommittedLen
			runeIdx := byteIndexToRuneIndex(cursorText, cursorBytes)
			cursorText = string([]rune(cursorText)[:runeIdx])
		} else {
			// 无合成文本时，优先使用 r.cursorPos（避免 field.Selection 同步延迟）
			cursorText = string([]rune(cursorText)[:r.cursorPos])
		}
	} else {
		cursorText = string([]rune(cursorText)[:r.cursorPos])
	}

	cursorW, _ := text.Measure(cursorText, face.Face, 0)
	lineHeight := float64(r.FontSize) * 1.2
	margin := float64(r.FontSize) * 0.5

	if !r.multiline {
		// 单行模式：水平滚动
		visibleWidth := float64(bounds.Width)

		if cursorW > visibleWidth+float64(r.scrollOffsetX) {
			// 光标超出右边界
			r.scrollOffsetX = float32(cursorW - visibleWidth + margin)
		} else if cursorW < float64(r.scrollOffsetX)+margin {
			// 光标在滚动区域左侧
			newOffset := float32(cursorW - margin)
			if newOffset < 0 {
				newOffset = 0
			}
			r.scrollOffsetX = newOffset
		}
		// 清除垂直滚动
		r.scrollOffsetY = 0
	} else {
		// 多行模式：垂直滚动
		cursorLine := r.cursorDisplayLine(face.Face)
		cursorY := float64(cursorLine) * lineHeight

		// 优先委托给可滚动父级（如 RenderScroll）
		delegated := false
		if p := r.GetParent(); p != nil {
			if scroller, ok := p.(interface{ ScrollTo(x, y float32) }); ok {
				visibleHeight := float64(bounds.Height)
				var targetY float32
				if cursorY+lineHeight > visibleHeight+margin {
					targetY = float32(cursorY + lineHeight - visibleHeight + margin)
				} else if cursorY < margin {
					targetY = float32(cursorY - margin)
					if targetY < 0 {
						targetY = 0
					}
				}
				scroller.ScrollTo(0, targetY)
				delegated = true
			}
		}

		if !delegated {
			// 无可滚动父级，回退到内部 scrollOffsetY
			visibleHeight := float64(bounds.Height)
			if cursorY+lineHeight > visibleHeight+float64(r.scrollOffsetY) {
				r.scrollOffsetY = float32(cursorY + lineHeight - visibleHeight + margin)
			} else if cursorY < float64(r.scrollOffsetY)+margin {
				newOffset := float32(cursorY - margin)
				if newOffset < 0 {
					newOffset = 0
				}
				r.scrollOffsetY = newOffset
			}
		}
		// 清除水平滚动
		r.scrollOffsetX = 0
	}
}

// cursorDisplayLine 计算光标所在的实际显示行号（0-based），同时考虑 SplitLines 产生的软换行
// 和用户输入的硬换行 \n。
func (r *RenderEditableText) cursorDisplayLine(face text.Face) int {
	contentRunes := []rune(r.Content)
	if r.cursorPos < 0 {
		r.cursorPos = 0
	}
	if r.cursorPos > len(contentRunes) {
		r.cursorPos = len(contentRunes)
	}
	if r.bounds.Width <= 0 {
		return strings.Count(string(contentRunes[:r.cursorPos]), "\n")
	}

	lines := r.RenderText.SplitLines(face, float64(r.bounds.Width), yoga.MeasureModeAtMost)
	cursorRunes := contentRunes[:r.cursorPos]

	displayLine := 0
	processedRunes := 0

	for _, line := range lines {
		lineRunes := []rune(line)
		if processedRunes+len(lineRunes) >= len(cursorRunes) {
			// 光标在当前 line 中，计算该行内部 \n 数量
			remainingRunes := len(cursorRunes) - processedRunes
			prefixRunes := lineRunes[:remainingRunes]
			displayLine += strings.Count(string(prefixRunes), "\n")
			return displayLine
		}
		// 当前 line 完全在光标前
		displayLine += strings.Count(line, "\n") + 1
		processedRunes += len(lineRunes)
	}

	return displayLine
}

func (r *RenderEditableText) GetScrollOffset() Offset {
	if r.multiline {
		// 多行模式下滚动由外层 RenderScroll 管理，不暴露内部偏移
		return Offset{X: 0, Y: 0}
	}
	return Offset{X: r.scrollOffsetX, Y: r.scrollOffsetY}
}

func (r *RenderEditableText) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	// 确定要显示的文本（包含 IME 合成中的文本）
	displayText := r.Content
	if r.field != nil && r.field.IsFocused() {
		displayText = r.field.TextForRendering()
	}

	face, err := fonts.GetFontFace(fonts.FontDescriptor{
		Family: fonts.FontFamilyDefault,
		Weight: fonts.FontWeightNormal,
		Style:  fonts.FontStyleNormal,
		Size:   r.FontSize,
	})

	lineHeight := float64(r.FontSize) * 1.2

	// 计算垂直居中偏移（仅单行模式）
	var centerOffset float32
	if !r.multiline && face != nil {
		m := face.Face.Metrics()
		textBlockHeight := float32(m.HAscent + m.HDescent)
		centerOffset = (bounds.Height - textBlockHeight) / 2
		if centerOffset < 0 {
			centerOffset = 0
		}
	}

	// 应用滚动偏移和垂直居中后的绘制偏移
	textOffset := offset
	if !r.multiline {
		// 单行模式使用内部滚动偏移；多行模式由外层 RenderScroll 管理滚动
		textOffset.X -= r.scrollOffsetX
		textOffset.Y -= r.scrollOffsetY
	}
	textOffset.Y += centerOffset

	// 绘制文本或 placeholder
	if displayText != "" {
		origContent := r.RenderText.Content
		r.RenderText.Content = displayText
		r.RenderText.Paint(screen, textOffset)
		r.RenderText.Content = origContent
	} else if r.placeholder != "" {
		origContent := r.RenderText.Content
		origColor := r.RenderText.TextColor
		r.RenderText.Content = r.placeholder
		r.RenderText.TextColor = r.placeholderColor
		r.RenderText.Paint(screen, textOffset)
		r.RenderText.Content = origContent
		r.RenderText.TextColor = origColor
	}

	// 绘制光标
	if r.focused && r.showCursor {
		if bounds.Width <= 0 || bounds.Height <= 0 {
			return
		}

		if err != nil || face == nil {
			return
		}

		// 光标位置基于当前显示文本的长度
		// 优先使用 r.cursorPos（避免 field.Selection 同步延迟）
		cursorText := displayText
		if r.field != nil && r.field.IsFocused() {
			startBytes, _ := r.field.Selection()
			// IME 合成阶段，未提交的文本长度
			uncommittedLen := r.field.UncommittedTextLengthInBytes()
			if uncommittedLen > 0 {
				// 有合成文本时，光标放在合成文本末尾
				cursorBytes := startBytes + uncommittedLen
				runeIdx := byteIndexToRuneIndex(cursorText, cursorBytes)
				cursorText = string([]rune(cursorText)[:runeIdx])
			} else {
				// 无合成文本时，直接使用 r.cursorPos
				cursorText = string([]rune(cursorText)[:r.cursorPos])
			}
		} else {
			cursorText = string([]rune(cursorText)[:r.cursorPos])
		}

		w, _ := text.Measure(cursorText, face.Face, 0)

		// 光标坐标使用与文本相同的偏移体系
		textX := float64(textOffset.X + bounds.X)
		textY := float64(textOffset.Y + bounds.Y)

		x := textX + w
		y := textY

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

// runeIndexToByteIndex 将 rune 索引转换为 UTF-8 字节索引。
func runeIndexToByteIndex(s string, runeIdx int) int {
	if runeIdx <= 0 {
		return 0
	}
	byteIdx := 0
	for i := 0; i < runeIdx && byteIdx < len(s); i++ {
		_, size := utf8.DecodeRuneInString(s[byteIdx:])
		byteIdx += size
	}
	return byteIdx
}

// byteIndexToRuneIndex 将 UTF-8 字节索引转换为 rune 索引。
func byteIndexToRuneIndex(s string, byteIdx int) int {
	if byteIdx <= 0 {
		return 0
	}
	runeIdx := 0
	b := 0
	for b < byteIdx && b < len(s) {
		_, size := utf8.DecodeRuneInString(s[b:])
		b += size
		runeIdx++
	}
	return runeIdx
}
