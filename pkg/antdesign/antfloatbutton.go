package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
)

// AntFloatButton 是 Ant Design 风格的浮动按钮。
type AntFloatButton struct {
	tenon.BaseWidget
	label   string
	onClick func()
}

// NewAntFloatButton 创建一个浮动按钮。
func NewAntFloatButton(label string) *AntFloatButton {
	f := &AntFloatButton{label: label}
	f.Init(f)
	return f
}

// SetOnClick 设置点击回调。
func (f *AntFloatButton) SetOnClick(fn func()) *AntFloatButton {
	f.onClick = fn
	return f
}

func (f *AntFloatButton) Render() tenon.Component {
	theme := tenon.GetTheme()
	fb := components.NewFloatButton(f.label)
	fb.SetBackgroundColor(theme.PrimaryColor)
	fb.SetTextColor(theme.ButtonTextColor)
	if f.onClick != nil {
		fb.SetOnClick(f.onClick)
	}
	return fb
}
