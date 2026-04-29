package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Radio is a single-select radio button with optional label.
type Radio struct {
	core.BaseElement
	selected    bool
	onChange    func(selected bool)
	boxSize     float32
	borderColor color.Color
	fillColor   color.Color
	innerColor  color.Color
	labelEl     *Text
}

// NewRadio creates a radio button.
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
	r.SetFlag(core.FlagFocusable)
	r.SetFlexDirection(yoga.FlexDirectionRow)
	r.SetAlignItems(yoga.AlignCenter)

	if label != "" {
		r.labelEl = NewText(label).SetColor(theme.TextColor)
		r.labelEl.SetMargin(yoga.EdgeLeft, r.boxSize+8)
		r.AppendChild(r.labelEl)
	}
	return r
}

// ElementType returns type identifier.
func (r *Radio) ElementType() string { return "Radio" }

// Draw renders the radio circle and inner dot.
func (r *Radio) Draw(screen *ebiten.Image) {
	bounds := r.GetBounds()

	centerY := bounds.Y + bounds.Height/2
	centerX := bounds.X + r.boxSize/2

	bgClr := core.GetTheme().BackgroundColor
	if r.selected {
		drawFilledCirclePath(screen, centerX, centerY, r.boxSize/2, r.fillColor)
		drawFilledCirclePath(screen, centerX, centerY, r.boxSize/4, r.innerColor)
	} else {
		drawFilledCirclePath(screen, centerX, centerY, r.boxSize/2, bgClr)
		strokeCirclePath(screen, centerX, centerY, r.boxSize/2, 1.5, r.borderColor)
	}
}

// HandleEvent processes click events.
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
		r.Mark(core.FlagNeedDraw)
		return true
	}
	return false
}

// SetSelected sets the selected state.
func (r *Radio) SetSelected(selected bool) *Radio {
	r.selected = selected
	r.Mark(core.FlagNeedDraw)
	return r
}

// DebugProps returns visual properties for debugger preview.
func (r *Radio) DebugProps() map[string]interface{} {
	props := make(map[string]interface{})
	props["selected"] = r.selected
	props["boxSize"] = r.boxSize
	if r.borderColor != nil {
		props["borderColor"] = colorToCSS(r.borderColor)
	}
	if r.selected && r.fillColor != nil {
		props["backgroundColor"] = colorToCSS(r.fillColor)
	}
	if r.innerColor != nil {
		props["innerColor"] = colorToCSS(r.innerColor)
	}
	return props
}

// SyncFrom 同步新 Radio 的属性到当前 Element（声明式重建）。
func (r *Radio) SyncFrom(src core.Element) {
	other, ok := src.(*Radio)
	if !ok {
		return
	}
	sb := &SyncBuilder{}
	syncField(sb, &r.selected, other.selected)
	syncField(sb, &r.boxSize, other.boxSize)
	syncColor(sb, &r.borderColor, other.borderColor)
	syncColor(sb, &r.fillColor, other.fillColor)
	syncColor(sb, &r.innerColor, other.innerColor)
	sb.MarkDraw(r)
}

// SetOnChange sets the change callback.
func (r *Radio) SetOnChange(fn func(selected bool)) *Radio {
	r.onChange = fn
	return r
}

// SetBoxSize sets the radio size.
func (r *Radio) SetBoxSize(size float32) *Radio {
	r.boxSize = size
	r.Mark(core.FlagNeedDraw)
	return r
}

// SetBorderColor sets the border color.
func (r *Radio) SetBorderColor(clr color.Color) *Radio {
	r.borderColor = clr
	r.Mark(core.FlagNeedDraw)
	return r
}

// SetFillColor sets the fill color.
func (r *Radio) SetFillColor(clr color.Color) *Radio {
	r.fillColor = clr
	r.Mark(core.FlagNeedDraw)
	return r
}

// SetInnerColor sets the inner dot color.
func (r *Radio) SetInnerColor(clr color.Color) *Radio {
	r.innerColor = clr
	r.Mark(core.FlagNeedDraw)
	return r
}

// SetTextColor sets the label text color.
func (r *Radio) SetTextColor(clr color.Color) *Radio {
	if r.labelEl != nil {
		r.labelEl.SetColor(clr)
	}
	return r
}

// SetFontSize sets the label font size.
func (r *Radio) SetFontSize(size float64) *Radio {
	if r.labelEl != nil {
		r.labelEl.SetFontSize(size)
	}
	return r
}

