package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Radio 是单选框组件。
type Radio struct {
	core.BaseHost
	selected     bool
	onChange     func(selected bool)
	boxSize      float32
	borderColor  color.Color
	fillColor    color.Color
	innerColor   color.Color
	text         *Text
}

// NewRadio 创建一个单选框。
func NewRadio(label string) *Radio {
	theme := core.GetTheme()
	r := &Radio{
		selected:    false,
		boxSize:     18,
		borderColor: theme.RadioBorderColor,
		fillColor:   theme.RadioFillColor,
		innerColor:  theme.RadioInnerColor,
	}
	r.Init(r)
	r.SetFocusable(true)
	r.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	r.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)

	if label != "" {
		r.text = NewText(label)
		r.text.SetFontSize(theme.FontSizeBase)
		r.text.SetColor(theme.TextColor)
		r.text.SetMargin(yoga.EdgeLeft, r.boxSize+8)
		r.AddChild(r.text)
	}
	return r
}

// Draw 绘制单选框圆形和选中点。
func (r *Radio) Draw(screen *ebiten.Image) {
	el := r.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := r.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	centerY := bounds.Y + bounds.Height/2
	centerX := bounds.X + r.boxSize/2

	if r.selected {
		// 外圆填充
		drawFilledCirclePath(screen, centerX, centerY, r.boxSize/2, r.fillColor)
		// 内圆白点
		drawFilledCirclePath(screen, centerX, centerY, r.boxSize/4, r.innerColor)
	} else {
		// 外圆白底
		drawFilledCirclePath(screen, centerX, centerY, r.boxSize/2, color.White)
		// 外圆边框
		strokeCirclePath(screen, centerX, centerY, r.boxSize/2, 1.5, r.borderColor)
	}
}

// HandleEvent 处理点击事件。
func (r *Radio) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventMouseDown:
		return true
	case core.EventMouseUp:
		return true
	case core.EventClick:
		r.selected = true
		if r.onChange != nil {
			r.onChange(r.selected)
		}
		return true
	}
	return false
}

// ==================== 链式 API ====================

func (r *Radio) SetSelected(selected bool) *Radio {
	r.selected = selected
	return r
}
func (r *Radio) SetOnChange(fn func(selected bool)) *Radio {
	r.onChange = fn
	return r
}
func (r *Radio) SetBoxSize(size float32) *Radio {
	r.boxSize = size
	return r
}
func (r *Radio) SetBorderColor(clr color.Color) *Radio {
	r.borderColor = clr
	return r
}
func (r *Radio) SetFillColor(clr color.Color) *Radio {
	r.fillColor = clr
	return r
}
func (r *Radio) SetInnerColor(clr color.Color) *Radio {
	r.innerColor = clr
	return r
}
func (r *Radio) SetTextColor(clr color.Color) *Radio {
	if r.text != nil {
		r.text.SetColor(clr)
	}
	return r
}
func (r *Radio) SetFontSize(size float32) *Radio {
	if r.text != nil {
		r.text.SetFontSize(size)
	}
	return r
}
func (r *Radio) SetMargin(edge yoga.Edge, value float32) *Radio {
	r.GetElement().Yoga.StyleSetMargin(edge, value)
	return r
}

// SyncFrom 同步单选框属性。
func (r *Radio) SyncFrom(other core.Host) {
	if o, ok := other.(*Radio); ok {
		r.selected = o.selected
		r.onChange = o.onChange
		r.boxSize = o.boxSize
		r.borderColor = o.borderColor
		r.fillColor = o.fillColor
		r.innerColor = o.innerColor
		if r.text != nil && o.text != nil {
			r.text.Content = o.text.Content
			r.text.cachedLayout = nil
		}
	}
}
