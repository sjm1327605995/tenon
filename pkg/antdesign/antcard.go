package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntCard 是 Ant Design 风格的卡片组件。
type AntCard struct {
	tenon.BaseWidget
	title     string
	extra     tenon.Component // 右上角操作区
	children  []tenon.Component
	shadow    int // 0=无, 1=小, 2=中, 3=大
	hoverable bool
}

// NewAntCard 创建卡片。
func NewAntCard() *AntCard {
	c := &AntCard{shadow: 1}
	c.Init(c)
	return c
}

// Render 返回卡片 UI。
func (c *AntCard) Render() tenon.Component {
	theme := NewAntTheme()

	card := components.NewView().
		SetBackgroundColor(theme.SurfaceColor).
		SetBorderRadius(theme.BorderRadius).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(theme.BorderColor).
		SetOverflow(yoga.OverflowHidden)

	// 阴影
	switch c.shadow {
	case 1:
		card.SetShadow(theme.ShadowColor, 6, 0, 1)
	case 2:
		card.SetShadow(theme.ShadowColor, 12, 0, 4)
	case 3:
		card.SetShadow(theme.ShadowColor, 24, 0, 8)
	}

	// 头部（标题 + extra）
	if c.title != "" || c.extra != nil {
		header := components.NewView().
			SetFlexDirection(yoga.FlexDirectionRow).
			SetJustifyContent(yoga.JustifySpaceBetween).
			SetAlignItems(yoga.AlignCenter).
			SetPadding(yoga.EdgeAll, 16).
			SetBorder(yoga.EdgeBottom, 1).
			SetBorderColor(theme.BorderColor)

		if c.title != "" {
			header.Add(components.NewText(c.title).
				SetFontSize(theme.FontSizeLG).
				SetColor(theme.TextColor).
				SetFontWeight(fonts.FontWeightBold))
		}
		if c.extra != nil {
			header.AddChild(c.extra)
		}
		card.AddChild(header)
	}

	// 内容区
	body := components.NewView().
		SetPadding(yoga.EdgeAll, 16).
		SetFlexDirection(yoga.FlexDirectionColumn)
	for _, child := range c.children {
		body.AddChild(child)
	}
	card.AddChild(body)

	return card
}

// ==================== 链式 API ====================

func (c *AntCard) SetTitle(t string) *AntCard {
	c.title = t
	return c
}
func (c *AntCard) SetExtra(extra tenon.Component) *AntCard {
	c.extra = extra
	return c
}
func (c *AntCard) Add(children ...tenon.Component) *AntCard {
	c.children = append(c.children, children...)
	return c
}
func (c *AntCard) SetShadow(level int) *AntCard {
	c.shadow = level
	return c
}
func (c *AntCard) SetHoverable(v bool) *AntCard {
	c.hoverable = v
	return c
}
