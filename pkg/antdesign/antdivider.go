package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntDivider 是 Ant Design 风格的分隔线，支持文字。
type AntDivider struct {
	tenon.BaseWidget
	text     string
	vertical bool // true=竖向
	plain    bool // 文字是否使用普通样式
	align    string // left | center | right，仅水平时有效
}

// NewAntDivider 创建分隔线。
func NewAntDivider() *AntDivider {
	d := &AntDivider{align: "center"}
	d.Init(d)
	return d
}

// Render 返回分隔线 UI。
func (d *AntDivider) Render() tenon.Component {
	theme := NewAntTheme()

	if d.vertical {
		return components.NewView().
			SetWidth(1).
			SetHeightPercent(100).
			SetBackgroundColor(theme.DividerColor)
	}

	// 水平分隔线 + 文字
	if d.text == "" {
		return components.NewView().
			SetHeight(1).
			SetWidthPercent(100).
			SetBackgroundColor(theme.DividerColor).
			SetMargin(yoga.EdgeVertical, 12)
	}

	// 带文字的水平分隔线：左线 + 文字 + 右线
	textComp := components.NewView().
		SetMargin(yoga.EdgeHorizontal, 16).
		Add(components.NewText(d.text).
			SetFontSize(theme.FontSizeSM).
			SetColor(theme.TextMutedColor))

	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignCenter).
		SetMargin(yoga.EdgeVertical, 12)

	newLine := func() *components.View {
		return components.NewView().
			SetHeight(1).
			SetBackgroundColor(theme.DividerColor)
	}

	switch d.align {
	case "left":
		root.Add(newLine().SetWidth(24))
		root.Add(textComp)
		root.Add(newLine().SetFlexGrow(1))
	case "right":
		root.Add(newLine().SetFlexGrow(1))
		root.Add(textComp)
		root.Add(newLine().SetWidth(24))
	default: // center
		root.Add(newLine().SetFlexGrow(1))
		root.Add(textComp)
		root.Add(newLine().SetFlexGrow(1))
	}
	return root
}

// ==================== 链式 API ====================

func (d *AntDivider) SetText(t string) *AntDivider {
	d.text = t
	return d
}
func (d *AntDivider) SetVertical(v bool) *AntDivider {
	d.vertical = v
	return d
}
func (d *AntDivider) SetPlain(v bool) *AntDivider {
	d.plain = v
	return d
}
func (d *AntDivider) SetAlign(a string) *AntDivider {
	d.align = a
	return d
}
