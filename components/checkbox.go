package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// Checkbox is a checkable box with optional label.
type Checkbox struct {
	core.BaseElement
	checked     bool
	onChange    func(checked bool)
	boxSize     float32
	borderColor color.Color
	fillColor   color.Color
	checkColor  color.Color
	labelEl     *native.Text
}

// NewCheckbox creates a checkbox.
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
	cb.SetFlag(core.FlagFocusable)
	cb.SetFlexDirection(yoga.FlexDirectionRow)
	cb.SetAlignItems(yoga.AlignCenter)

	if label != "" {
		cb.labelEl = native.NewText(label).SetColor(theme.TextColor)
		cb.labelEl.SetMargin(yoga.EdgeLeft, cb.boxSize+8)
		cb.AppendChild(cb.labelEl)
	}
	return cb
}

// ElementType returns type identifier.
func (cb *Checkbox) ElementType() string { return "Checkbox" }

// Draw renders the checkbox box and checkmark.
func (cb *Checkbox) Draw(screen *ebiten.Image) {
	bounds := cb.GetBounds()

	boxY := bounds.Y + (bounds.Height-cb.boxSize)/2
	boxX := bounds.X

	bgClr := core.GetTheme().BackgroundColor
	if cb.checked {
		cb.drawRoundedBox(screen, boxX, boxY, cb.boxSize, cb.boxSize, cb.fillColor)
		cb.drawCheckmark(screen, boxX, boxY, cb.boxSize)
	} else {
		cb.drawRoundedBox(screen, boxX, boxY, cb.boxSize, cb.boxSize, bgClr)
		cb.drawBoxBorder(screen, boxX, boxY, cb.boxSize, cb.boxSize, cb.borderColor)
	}
}

func (cb *Checkbox) drawRoundedBox(screen *ebiten.Image, x, y, w, h float32, clr color.Color) {
	r := float32(3)
	var path vector.Path
	native.BuildRoundedRectPath(&path, x, y, w, h, core.BorderRadius{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r})
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

func (cb *Checkbox) drawBoxBorder(screen *ebiten.Image, x, y, w, h float32, clr color.Color) {
	r := float32(3)
	var path vector.Path
	native.BuildRoundedRectPath(&path, x, y, w, h, core.BorderRadius{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r})
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

// HandleEvent processes click events.
func (cb *Checkbox) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventMouseDown:
		return true
	case core.EventMouseUp:
		return true
	case core.EventClick:
		cb.checked = !cb.checked
		if cb.onChange != nil {
			cb.onChange(cb.checked)
		}
		cb.Mark(core.FlagNeedDraw)
		return true
	}
	return false
}

// SetChecked sets the checked state.
func (cb *Checkbox) SetChecked(checked bool) *Checkbox {
	cb.checked = checked
	cb.Mark(core.FlagNeedDraw)
	return cb
}

// DebugProps returns visual properties for debugger preview.
func (cb *Checkbox) DebugProps() map[string]interface{} {
	props := make(map[string]interface{})
	props["checked"] = cb.checked
	props["boxSize"] = cb.boxSize
	if cb.borderColor != nil {
		props["borderColor"] = native.ColorToCSS(cb.borderColor)
	}
	if cb.checked && cb.fillColor != nil {
		props["backgroundColor"] = native.ColorToCSS(cb.fillColor)
	}
	if cb.checkColor != nil {
		props["checkColor"] = native.ColorToCSS(cb.checkColor)
	}
	return props
}

// SyncFrom 同步新 Checkbox 的属性到当前 Element（声明式重建）。
func (cb *Checkbox) SyncFrom(src core.Element) {
	other, ok := src.(*Checkbox)
	if !ok {
		return
	}
	sb := &core.SyncBuilder{}
	core.SyncField(sb, &cb.checked, other.checked)
	core.SyncField(sb, &cb.boxSize, other.boxSize)
	core.SyncColor(sb, &cb.borderColor, other.borderColor)
	core.SyncColor(sb, &cb.fillColor, other.fillColor)
	core.SyncColor(sb, &cb.checkColor, other.checkColor)
	sb.MarkDraw(cb)
}

// SetOnChange sets the change callback.
func (cb *Checkbox) SetOnChange(fn func(checked bool)) *Checkbox {
	cb.onChange = fn
	return cb
}

// SetBoxSize sets the box size.
func (cb *Checkbox) SetBoxSize(size float32) *Checkbox {
	cb.boxSize = size
	cb.Mark(core.FlagNeedDraw)
	return cb
}

// SetBorderColor sets the border color.
func (cb *Checkbox) SetBorderColor(clr color.Color) *Checkbox {
	cb.borderColor = clr
	cb.Mark(core.FlagNeedDraw)
	return cb
}

// SetFillColor sets the fill color.
func (cb *Checkbox) SetFillColor(clr color.Color) *Checkbox {
	cb.fillColor = clr
	cb.Mark(core.FlagNeedDraw)
	return cb
}

// SetCheckColor sets the checkmark color.
func (cb *Checkbox) SetCheckColor(clr color.Color) *Checkbox {
	cb.checkColor = clr
	cb.Mark(core.FlagNeedDraw)
	return cb
}

// SetTextColor sets the label text color.
func (cb *Checkbox) SetTextColor(clr color.Color) *Checkbox {
	if cb.labelEl != nil {
		cb.labelEl.SetColor(clr)
	}
	return cb
}

// SetFontSize sets the label font size.
func (cb *Checkbox) SetFontSize(size float64) *Checkbox {
	if cb.labelEl != nil {
		cb.labelEl.SetFontSize(size)
	}
	return cb
}

