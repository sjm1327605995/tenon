package antdesign

import (
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntTagColor 定义 Tag 颜色类型。
type AntTagColor string

const (
	AntTagDefault AntTagColor = "default"
	AntTagRed     AntTagColor = "red"
	AntTagOrange  AntTagColor = "orange"
	AntTagGreen   AntTagColor = "green"
	AntTagBlue    AntTagColor = "blue"
)

// AntTag 是 Ant Design 风格的标签组件。
type AntTag struct {
	tenon.BaseWidget
	label    string
	color    AntTagColor
	closable bool
	onClose  func()
}

// NewAntTag 创建标签。
func NewAntTag(label string) *AntTag {
	t := &AntTag{label: label, color: AntTagDefault}
	t.Init(t)
	return t
}

// Render 返回标签 UI。
func (t *AntTag) Render() tenon.Component {
	theme := NewAntTheme()
	bg, textClr, border := t.resolveColors(theme)

	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignCenter).
		SetBackgroundColor(bg).
		SetBorderRadius(4).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(border).
		SetPadding(yoga.EdgeHorizontal, 7).
		SetPadding(yoga.EdgeVertical, 1).
		SetMargin(yoga.EdgeRight, 8)

	root.Add(components.NewText(t.label).
		SetFontSize(theme.FontSizeSM).
		SetColor(textClr))

	if t.closable {
		root.Add(components.NewText(" ×").
			SetFontSize(theme.FontSizeSM).
			SetColor(textClr).
			SetMargin(yoga.EdgeLeft, 4))
		// 点击关闭需要命中检测，简化：整个 Tag 区域可点击关闭
		root.SetPointerEvents(core.PointerEventsAuto)
	}

	return root
}

func (t *AntTag) resolveColors(theme *AntTheme) (bg, text, border color.Color) {
	switch t.color {
	case AntTagRed:
		return theme.TagRedBg, theme.TagRedText, theme.TagRedBg
	case AntTagOrange:
		return theme.TagOrangeBg, theme.TagOrangeText, theme.TagOrangeBg
	case AntTagGreen:
		return theme.TagGreenBg, theme.TagGreenText, theme.TagGreenBg
	case AntTagBlue:
		return theme.TagBlueBg, theme.TagBlueText, theme.TagBlueBg
	default:
		return theme.DisabledBgColor, theme.TextMutedColor, theme.DisabledBgColor
	}
}

// ==================== 链式 API ====================

func (t *AntTag) SetColor(c AntTagColor) *AntTag {
	t.color = c
	return t
}
func (t *AntTag) SetClosable(v bool) *AntTag {
	t.closable = v
	return t
}
func (t *AntTag) SetOnClose(fn func()) *AntTag {
	t.onClose = fn
	return t
}
