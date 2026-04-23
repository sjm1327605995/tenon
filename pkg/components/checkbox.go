package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Checkbox 是复选框组件。
type Checkbox struct {
	core.BaseHost
	checked      bool
	onChange     func(checked bool)
	boxSize      float32
	borderColor  color.Color
	fillColor    color.Color
	checkColor   color.Color
	text         *Text
}

// NewCheckbox 创建一个复选框。
func NewCheckbox(label string) *Checkbox {
	theme := core.GetTheme()
	cb := &Checkbox{
		checked:     false,
		boxSize:     18,
		borderColor: theme.CheckboxBorderColor,
		fillColor:   theme.CheckboxFillColor,
		checkColor:  theme.CheckboxCheckColor,
	}
	cb.Init(cb)
	cb.SetFocusable(true)
	// 使用 Row 布局让文字在方框右侧；不固定高度，由文字自适应
	cb.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	cb.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)

	if label != "" {
		cb.text = NewText(label)
		cb.text.SetFontSize(theme.FontSizeBase)
		cb.text.SetColor(theme.TextColor)
		// 左边距 = 方框宽度 + 文字与方框间距
		cb.text.SetMargin(yoga.EdgeLeft, cb.boxSize+8)
		cb.AddChild(cb.text)
	}
	return cb
}

// Draw 绘制复选框方框和对勾。
func (cb *Checkbox) Draw(screen *ebiten.Image) {
	el := cb.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := cb.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	boxY := bounds.Y + (bounds.Height-cb.boxSize)/2
	boxX := bounds.X

	// 绘制方框背景
	if cb.checked {
		cb.drawRoundedBox(screen, boxX, boxY, cb.boxSize, cb.boxSize, cb.fillColor)
		// 绘制对勾
		cb.drawCheckmark(screen, boxX, boxY, cb.boxSize)
	} else {
		cb.drawRoundedBox(screen, boxX, boxY, cb.boxSize, cb.boxSize, color.White)
		// 绘制边框
		cb.drawBoxBorder(screen, boxX, boxY, cb.boxSize, cb.boxSize, cb.borderColor)
	}
}

func (cb *Checkbox) drawRoundedBox(screen *ebiten.Image, x, y, w, h float32, clr color.Color) {
	r := float32(3)
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, core.BorderRadius{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r})
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

func (cb *Checkbox) drawBoxBorder(screen *ebiten.Image, x, y, w, h float32, clr color.Color) {
	r := float32(3)
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, core.BorderRadius{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r})
	strokeOp := &vector.StrokeOptions{Width: 1.5, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}

func (cb *Checkbox) drawCheckmark(screen *ebiten.Image, x, y, size float32) {
	padding := size * 0.2
	x1 := x + padding
	y1 := y + size/2
	x2 := x + size*0.4
	y2 := y + size - padding
	x3 := x + size - padding
	y3 := y + padding

	var path vector.Path
	path.MoveTo(x1, y1)
	path.LineTo(x2, y2)
	path.LineTo(x3, y3)

	strokeOp := &vector.StrokeOptions{Width: 2, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(cb.checkColor)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}

// HandleEvent 处理点击事件。
func (cb *Checkbox) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventClick:
		cb.checked = !cb.checked
		if cb.onChange != nil {
			cb.onChange(cb.checked)
		}
		return true
	}
	return false
}

// ==================== 链式 API ====================

func (cb *Checkbox) SetChecked(checked bool) *Checkbox {
	cb.checked = checked
	return cb
}
func (cb *Checkbox) SetOnChange(fn func(checked bool)) *Checkbox {
	cb.onChange = fn
	return cb
}
func (cb *Checkbox) SetBoxSize(size float32) *Checkbox {
	cb.boxSize = size
	cb.GetElement().Yoga.StyleSetHeight(size)
	return cb
}
func (cb *Checkbox) SetBorderColor(clr color.Color) *Checkbox {
	cb.borderColor = clr
	return cb
}
func (cb *Checkbox) SetFillColor(clr color.Color) *Checkbox {
	cb.fillColor = clr
	return cb
}
func (cb *Checkbox) SetCheckColor(clr color.Color) *Checkbox {
	cb.checkColor = clr
	return cb
}
func (cb *Checkbox) SetTextColor(clr color.Color) *Checkbox {
	if cb.text != nil {
		cb.text.SetColor(clr)
	}
	return cb
}
func (cb *Checkbox) SetFontSize(size float32) *Checkbox {
	if cb.text != nil {
		cb.text.SetFontSize(size)
	}
	return cb
}
func (cb *Checkbox) SetMargin(edge yoga.Edge, value float32) *Checkbox {
	cb.GetElement().Yoga.StyleSetMargin(edge, value)
	return cb
}
func (cb *Checkbox) SetWidth(width float32) *Checkbox {
	cb.GetElement().Yoga.StyleSetWidth(width)
	return cb
}
func (cb *Checkbox) SetHeight(height float32) *Checkbox {
	cb.GetElement().Yoga.StyleSetHeight(height)
	return cb
}

// SyncFrom 同步复选框属性。
func (cb *Checkbox) SyncFrom(other core.Host) {
	if o, ok := other.(*Checkbox); ok {
		cb.checked = o.checked
		cb.onChange = o.onChange
		cb.boxSize = o.boxSize
		cb.borderColor = o.borderColor
		cb.fillColor = o.fillColor
		cb.checkColor = o.checkColor
		if cb.text != nil && o.text != nil {
			cb.text.Content = o.text.Content
			cb.text.cachedLayout = nil
		}
	}
}
