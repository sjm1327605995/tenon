package antdesign

import (
	"strconv"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntBadge 是徽标组件，包裹一个子组件并在右上角显示数字/红点。
type AntBadge struct {
	tenon.BaseWidget
	child    tenon.Component // 被包裹的组件
	count    int           // 数字，0 隐藏
	dot      bool          // 纯红点
	overflow string        // 超过时显示，如 "99+"
}

// NewAntBadge 创建徽标。
func NewAntBadge(child tenon.Component) *AntBadge {
	b := &AntBadge{child: child}
	b.Init(b)
	return b
}

// Render 返回徽标 UI。
func (b *AntBadge) Render() tenon.Component {
	theme := NewAntTheme()

	container := components.NewView().
		SetPositionType(yoga.PositionTypeRelative)

	// 被包裹组件
	if b.child != nil {
		container.AddChild(b.child)
	}

	// 徽标
	badge := b.buildBadge(theme)
	if badge != nil {
		badge.SetPosition(yoga.EdgeTop, -6).
			SetPosition(yoga.EdgeRight, -6).
			SetPositionType(yoga.PositionTypeAbsolute)
		container.AddChild(badge)
	}

	return container
}

func (b *AntBadge) buildBadge(theme *AntTheme) *components.View {
	if b.count == 0 && !b.dot {
		return nil
	}

	if b.dot {
		return components.NewView().
			SetWidth(8).SetHeight(8).
			SetBackgroundColor(theme.ErrorColor).
			SetBorderRadius(4)
	}

	label := b.badgeText()
	w := float32(18)
	if len(label) > 2 {
		w = float32(len(label)*8 + 4)
	}

	return components.NewView().
		SetWidth(w).SetHeight(18).
		SetBackgroundColor(theme.ErrorColor).
		SetBorderRadius(9).
		SetJustifyContent(yoga.JustifyCenter).
		SetAlignItems(yoga.AlignCenter).
		Add(components.NewText(label).
			SetFontSize(10).
			SetColor(core.GetTheme().SurfaceColor))
}

func (b *AntBadge) badgeText() string {
	if b.overflow != "" && b.count > 99 {
		return b.overflow
	}
	if b.count > 99 {
		return "99+"
	}
	return strconv.Itoa(b.count)
}

// ==================== 链式 API ====================

func (b *AntBadge) SetCount(n int) *AntBadge {
	b.count = n
	return b
}
func (b *AntBadge) SetDot(v bool) *AntBadge {
	b.dot = v
	return b
}
func (b *AntBadge) SetOverflow(s string) *AntBadge {
	b.overflow = s
	return b
}
