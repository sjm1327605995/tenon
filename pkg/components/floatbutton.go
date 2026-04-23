package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// FloatButton 是一个固定在屏幕角落的浮动按钮。
type FloatButton struct {
	View
	btn     *Button
	onClick func()
}

// NewFloatButton 创建一个浮动按钮。
func NewFloatButton(label string) *FloatButton {
	fb := &FloatButton{}
	fb.Init(fb)
	fb.SetPositionType(yoga.PositionTypeAbsolute).
		SetPosition(yoga.EdgeRight, 24).
		SetPosition(yoga.EdgeBottom, 24).
		SetWidth(56).
		SetHeight(56).
		SetBorderRadius(28).
		SetBackgroundColor(color.RGBA{A: 0}).
		SetJustifyContent(yoga.JustifyCenter).
		SetAlignItems(yoga.AlignCenter).
		SetPointerEvents(core.PointerEventsAuto)

	fb.btn = NewButton(label)
	fb.btn.SetWidth(56).SetHeight(56).SetBorderRadius(28)
	fb.btn.GetElement().PointerEvents = core.PointerEventsNone
	fb.AddChild(fb.btn)

	return fb
}

// SetOnClick 设置点击回调。
func (fb *FloatButton) SetOnClick(fn func()) *FloatButton {
	fb.onClick = fn
	return fb
}

// SetPosition 设置浮动按钮的位置（right/bottom 等）。
func (fb *FloatButton) SetPosition(edge yoga.Edge, value float32) *FloatButton {
	fb.GetElement().Yoga.StyleSetPosition(edge, value)
	return fb
}

// SetBackgroundColor 设置按钮背景色。
func (fb *FloatButton) SetBackgroundColor(clr color.Color) *FloatButton {
	fb.btn.SetBackgroundColors(clr, clr, clr)
	return fb
}

// SetTextColor 设置按钮文字颜色。
func (fb *FloatButton) SetTextColor(clr color.Color) *FloatButton {
	fb.btn.SetTextColor(clr)
	return fb
}

// HandleEvent 处理点击事件。
func (fb *FloatButton) HandleEvent(e *core.Event) bool {
	if e.Type == core.EventClick && fb.onClick != nil {
		fb.onClick()
		return true
	}
	return false
}

// Draw 绘制浮动按钮。
func (fb *FloatButton) Draw(screen *ebiten.Image) {
	fb.View.Draw(screen)
}
