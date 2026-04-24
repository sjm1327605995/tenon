package components

import (
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// WhiteSpace 定义 CSS white-space 属性的行为。
type WhiteSpace int

const (
	WhiteSpaceNormal WhiteSpace = iota   // 合并空白，自动换行
	WhiteSpaceNoWrap                     // 合并空白，不换行
	WhiteSpacePre                        // 保留空白，不换行
	WhiteSpacePreWrap                    // 保留空白，自动换行
	WhiteSpacePreLine                    // 合并空白（保留换行），自动换行
)

// WordBreak 定义 CSS word-break 属性的行为。
type WordBreak int

const (
	WordBreakNormal WordBreak = iota     // 默认：CJK 可断，英文按词
	WordBreakBreakAll                    // 所有字符可断
	WordBreakKeepAll                     // 保持词完整
)

// textLayoutResult 保存文本布局结果。
type textLayoutResult struct {
	content    string  // 带 \n 的文本
	width      float32
	height     float32
	lineCount  int
	lineHeight float32
}

// computeTextLayout 根据策略计算文本布局。
func computeTextLayout(rawContent string, face *text.GoTextFace, maxWidth float32, ws WhiteSpace, wb WordBreak, explicitLineHeight float32) textLayoutResult {
	// 计算行高
	lineHeight := explicitLineHeight
	if lineHeight <= 0 {
		m := face.Metrics()
		lineHeight = float32(m.HAscent + m.HDescent)
		if lineHeight <= 0 {
			lineHeight = float32(face.Size) * 1.2
		}
	}

	// 处理 white-space
	paragraphs := processWhiteSpace(rawContent, ws)

	// 是否自动换行
	shouldWrap := ws != WhiteSpaceNoWrap && ws != WhiteSpacePre && maxWidth > 0

	var wrappedLines []string
	for _, para := range paragraphs {
		if shouldWrap {
			wrapped := wrapParagraph(para, maxWidth, wb, face)
			wrappedLines = append(wrappedLines, wrapped...)
		} else {
			wrappedLines = append(wrappedLines, para)
		}
	}

	content := strings.Join(wrappedLines, "\n")
	if content == "" {
		return textLayoutResult{lineHeight: lineHeight}
	}

	w, h := text.Measure(content, face, float64(lineHeight))
	return textLayoutResult{
		content:    content,
		width:      float32(w),
		height:     float32(h),
		lineCount:  len(wrappedLines),
		lineHeight: lineHeight,
	}
}

func processWhiteSpace(content string, ws WhiteSpace) []string {
	switch ws {
	case WhiteSpacePre, WhiteSpacePreWrap:
		normalized := strings.ReplaceAll(content, "\r\n", "\n")
		normalized = strings.ReplaceAll(normalized, "\r", "\n")
		return strings.Split(normalized, "\n")
	case WhiteSpacePreLine:
		rawLines := strings.Split(content, "\n")
		result := make([]string, len(rawLines))
		for i, line := range rawLines {
			result[i] = collapseWhiteSpace(line)
		}
		return result
	case WhiteSpaceNormal, WhiteSpaceNoWrap:
		collapsed := collapseWhiteSpace(content)
		return []string{collapsed}
	default:
		return []string{content}
	}
}

func collapseWhiteSpace(s string) string {
	var b strings.Builder
	inSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			inSpace = true
		} else {
			if inSpace && b.Len() > 0 {
				b.WriteByte(' ')
			}
			inSpace = false
			b.WriteRune(r)
		}
	}
	return b.String()
}

// token 表示一个文本片段。
type token struct {
	text         string
	isBreakable  bool // 该 token 内部是否可进一步拆分
	isWhitespace bool
}

func tokenizeParagraph(para string, wb WordBreak) []token {
	if wb == WordBreakBreakAll {
		var tokens []token
		for _, r := range para {
			tokens = append(tokens, token{text: string(r), isBreakable: true})
		}
		return tokens
	}

	if wb == WordBreakKeepAll {
		return tokenizeBySpace(para, false)
	}

	return tokenizeNormal(para)
}

func tokenizeBySpace(para string, cjkSeparate bool) []token {
	var tokens []token
	var word strings.Builder

	flushWord := func() {
		if word.Len() > 0 {
			w := word.String()
			tokens = append(tokens, token{text: w, isBreakable: !cjkSeparate || !isAllCJK(w)})
			word.Reset()
		}
	}

	for _, r := range para {
		if unicode.IsSpace(r) {
			flushWord()
			tokens = append(tokens, token{text: string(r), isWhitespace: true, isBreakable: false})
		} else if cjkSeparate && isCJK(r) {
			flushWord()
			tokens = append(tokens, token{text: string(r), isBreakable: true})
		} else {
			word.WriteRune(r)
		}
	}
	flushWord()
	return tokens
}

func tokenizeNormal(para string) []token {
	var tokens []token
	var word strings.Builder
	inCJK := false

	flushWord := func() {
		if word.Len() > 0 {
			tokens = append(tokens, token{text: word.String(), isBreakable: true})
			word.Reset()
		}
	}

	for _, r := range para {
		if unicode.IsSpace(r) {
			flushWord()
			tokens = append(tokens, token{text: string(r), isWhitespace: true, isBreakable: false})
			inCJK = false
		} else if isCJK(r) {
			flushWord()
			inCJK = true
			// CJK 字符每个都是独立的可断行 token
			tokens = append(tokens, token{text: string(r), isBreakable: true})
		} else {
			if inCJK {
				flushWord()
			}
			inCJK = false
			word.WriteRune(r)
		}
	}
	flushWord()
	return tokens
}

func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Hiragana, r) ||
		unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Hangul, r)
}

func isAllCJK(s string) bool {
	for _, r := range s {
		if !isCJK(r) {
			return false
		}
	}
	return len(s) > 0
}

func wrapParagraph(para string, maxWidth float32, wb WordBreak, face *text.GoTextFace) []string {
	if para == "" {
		return []string{""}
	}

	tokens := tokenizeParagraph(para, wb)
	var lines []string
	var line strings.Builder
	var lineWidth float32

	flushLine := func() {
		s := line.String()
		s = strings.TrimRightFunc(s, unicode.IsSpace)
		lines = append(lines, s)
		line.Reset()
		lineWidth = 0
	}

	for _, tok := range tokens {
		if tok.isWhitespace {
			line.WriteString(tok.text)
			lineWidth += measureString(tok.text, face)
			continue
		}

		tokWidth := measureString(tok.text, face)

		if line.Len() == 0 {
			if tokWidth > maxWidth && tok.isBreakable {
				parts := breakWord(tok.text, maxWidth, face)
				for j, part := range parts {
					if j < len(parts)-1 {
						lines = append(lines, part)
					} else {
						line.WriteString(part)
						lineWidth = measureString(part, face)
					}
				}
			} else {
				line.WriteString(tok.text)
				lineWidth = tokWidth
			}
		} else {
			testLine := line.String() + tok.text
			testWidth := measureString(testLine, face)

			if testWidth <= maxWidth {
				line.WriteString(tok.text)
				lineWidth = testWidth
			} else {
				flushLine()

				if tokWidth > maxWidth && tok.isBreakable {
					parts := breakWord(tok.text, maxWidth, face)
					for j, part := range parts {
						if j < len(parts)-1 {
							lines = append(lines, part)
						} else {
							line.WriteString(part)
							lineWidth = measureString(part, face)
						}
					}
				} else {
					line.WriteString(tok.text)
					lineWidth = tokWidth
				}
			}
		}
	}

	if line.Len() > 0 {
		s := line.String()
		s = strings.TrimRightFunc(s, unicode.IsSpace)
		lines = append(lines, s)
	}

	return lines
}

// breakWord 将一个过长的词按字符拆分成多行。
func breakWord(word string, maxWidth float32, face *text.GoTextFace) []string {
	var parts []string
	var current strings.Builder

	for _, r := range word {
		test := current.String() + string(r)
		if measureString(test, face) > maxWidth && current.Len() > 0 {
			parts = append(parts, current.String())
			current.Reset()
		}
		current.WriteRune(r)
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func measureString(s string, face *text.GoTextFace) float32 {
	if s == "" {
		return 0
	}
	w, _ := text.Measure(s, face, 0)
	return float32(w)
}
