package render

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/yoga"
)

// ButtonState 描述按钮的交互状态。
type ButtonState int

const (
	ButtonStateNormal ButtonState = iota
	ButtonStateHover
	ButtonStatePressed
	ButtonStateDisabled
)

// RenderButton 是按钮的 RenderObject，负责根据状态绘制背景和文字。
type RenderButton struct {
	RenderBox
	NormalColor   color.Color
	HoverColor    color.Color
	PressedColor  color.Color
	DisabledColor  color.Color
	TextColor      color.Color
	HoverTextColor color.Color
	BorderColor    color.Color
	State          ButtonState
	IsOutline      bool
	IsGhost        bool
	IsLink         bool
	Loading        bool
}

func NewRenderButton() *RenderButton {
	r := &RenderButton{}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	r.yoga.StyleSetJustifyContent(yoga.JustifyCenter)
	r.yoga.StyleSetAlignItems(yoga.AlignCenter)
	r.yoga.StyleSetPadding(yoga.EdgeHorizontal, 16)
	r.yoga.StyleSetPadding(yoga.EdgeVertical, 10)
	r.State = ButtonStateNormal
	return r
}

func (r *RenderButton) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	x := offset.X + bounds.X
	y := offset.Y + bounds.Y
	w := bounds.Width
	h := bounds.Height

	// 根据状态确定背景色
	bgColor := r.resolveBgColor()

	// Disabled 状态不绘制阴影
	if r.State != ButtonStateDisabled {
		PaintBoxShadow(screen, x, y, w, h, r.BorderRadius, r.ShadowColor, r.ShadowBlur, r.ShadowOffsetX, r.ShadowOffsetY)
	}
	PaintBackground(screen, x, y, w, h, r.BorderRadius, bgColor)

	// Link 变体不绘制边框
	if !r.IsLink {
		PaintBorder(screen, x, y, w, h, r.BorderRadius, r.BorderWidth, r.BorderColor)
	}

	// 绘制 loading spinner
	if r.Loading {
		r.drawLoading(screen, x+w/2, y+h/2)
	}
}

func (r *RenderButton) resolveBgColor() color.Color {
	switch r.State {
	case ButtonStateHover:
		if r.HoverColor != nil {
			return r.HoverColor
		}
		if r.NormalColor != nil {
			return Lighten(r.NormalColor, 20)
		}
	case ButtonStatePressed:
		if r.PressedColor != nil {
			return r.PressedColor
		}
		if r.NormalColor != nil {
			return Darken(r.NormalColor, 20)
		}
	case ButtonStateDisabled:
		if r.DisabledColor != nil {
			return r.DisabledColor
		}
	}
	if r.IsOutline || r.IsGhost || r.IsLink {
		return nil
	}
	return r.NormalColor
}

func (r *RenderButton) drawLoading(screen *ebiten.Image, cx, cy float32) {
	radius := float32(6)
	stroke := float32(2)
	// 简单的弧形 spinner
	path := &vector.Path{}
	path.MoveTo(cx+radius, cy)
	path.Arc(cx, cy, radius, 0, float32(math.Pi), vector.Clockwise)
	vector.StrokePath(screen, path, &vector.StrokeOptions{Width: stroke}, &vector.DrawPathOptions{
		ColorScale: toColorScale(r.TextColor),
	})
}

func (r *RenderButton) SetState(s ButtonState) {
	if r.State == s {
		return
	}
	r.State = s
	r.MarkNeedsPaint()

	// Link variant: auto-update child text color on state change
	if r.IsLink {
		var textColor color.Color
		if s == ButtonStateHover || s == ButtonStatePressed {
			textColor = r.HoverTextColor
		} else {
			textColor = r.TextColor
		}
		if textColor != nil {
			r.updateChildTextColor(textColor)
		}
	}
}

func (r *RenderButton) updateChildTextColor(c color.Color) {
	for _, child := range r.GetChildren() {
		if text, ok := child.(*RenderText); ok {
			text.SetColor(c)
		}
	}
}
