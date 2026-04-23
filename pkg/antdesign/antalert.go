package antdesign

import (
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntAlertType 定义 Alert 类型。
type AntAlertType string

const (
	AntAlertError   AntAlertType = "error"
	AntAlertWarning AntAlertType = "warning"
	AntAlertSuccess AntAlertType = "success"
	AntAlertInfo    AntAlertType = "info"
)

// AntAlert 是警告提示组件。
type AntAlert struct {
	tenon.BaseWidget
	msg      string
	desc     string
	alertType AntAlertType
	closable bool
	showIcon bool
	banner   bool // 顶部通栏
}

// NewAntAlert 创建警告提示。
func NewAntAlert(msg string) *AntAlert {
	a := &AntAlert{msg: msg, alertType: AntAlertInfo, showIcon: true}
	a.Init(a)
	return a
}

// Render 返回 Alert UI。
func (a *AntAlert) Render() tenon.Component {
	theme := NewAntTheme()
	bg, textClr, icon := a.resolveStyle(theme)

	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignFlexStart).
		SetBackgroundColor(bg).
		SetBorderRadius(theme.BorderRadius).
		SetPadding(yoga.EdgeAll, 12).
		SetMargin(yoga.EdgeBottom, 12)

	if a.banner {
		root.SetBorderRadius(0)
		root.SetMargin(yoga.EdgeBottom, 0)
	}

	// 图标
	if a.showIcon {
		root.Add(components.NewText(icon).
			SetFontSize(theme.FontSizeLG).
			SetMargin(yoga.EdgeRight, 8))
	}

	// 内容区
	body := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetFlexGrow(1)

	body.Add(components.NewText(a.msg).
		SetFontSize(theme.FontSizeBase).
		SetColor(textClr))

	if a.desc != "" {
		body.Add(components.NewText(a.desc).
			SetFontSize(theme.FontSizeSM).
			SetColor(theme.TextMutedColor).
			SetMargin(yoga.EdgeTop, 4))
	}
	root.AddChild(body)

	// 关闭按钮
	if a.closable {
		root.Add(components.NewText("×").
			SetFontSize(theme.FontSizeLG).
			SetColor(theme.TextMutedColor).
			SetMargin(yoga.EdgeLeft, 8))
	}

	return root
}

func (a *AntAlert) resolveStyle(theme *AntTheme) (bg, text color.Color, icon string) {
	switch a.alertType {
	case AntAlertError:
		return theme.AlertErrorBg, theme.ErrorColor, "✖"
	case AntAlertWarning:
		return theme.AlertWarningBg, theme.WarningColor, "⚠"
	case AntAlertSuccess:
		return theme.AlertSuccessBg, theme.SuccessColor, "✓"
	default: // info
		return theme.AlertInfoBg, theme.InfoColor, "ℹ"
	}
}

// ==================== 链式 API ====================

func (a *AntAlert) SetDesc(d string) *AntAlert {
	a.desc = d
	return a
}
func (a *AntAlert) SetType(t AntAlertType) *AntAlert {
	a.alertType = t
	return a
}
func (a *AntAlert) SetClosable(v bool) *AntAlert {
	a.closable = v
	return a
}
func (a *AntAlert) SetShowIcon(v bool) *AntAlert {
	a.showIcon = v
	return a
}
func (a *AntAlert) SetBanner(v bool) *AntAlert {
	a.banner = v
	return a
}
