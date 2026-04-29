package components

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// Switch is a toggle switch component.
type Switch struct {
	core.BaseElement
	checked       bool
	onChange      func(checked bool)
	trackWidth    float32
	trackHeight   float32
	offColor      color.Color
	onColor       color.Color
	thumbColor    color.Color
	thumbProgress float32 // 0=off, 1=on, for animation interpolation
	thumbAnim     *core.Tween
}

// NewSwitch creates a switch.
func NewSwitch() *Switch {
	theme := core.GetTheme()
	s := &Switch{
		checked:     false,
		trackWidth:  44,
		trackHeight: 24,
		offColor:    theme.SwitchOffColor,
		onColor:     theme.SwitchOnColor,
		thumbColor:  theme.SwitchThumbColor,
	}
	s.Init(s)
	s.SetFlag(core.FlagFocusable)
	s.SetWidth(s.trackWidth)
	s.SetHeight(s.trackHeight)
	return s
}

// ElementType returns type identifier.
func (s *Switch) ElementType() string { return "Switch" }

// Draw renders the switch track and thumb.
func (s *Switch) Draw(screen *ebiten.Image) {
	bounds := s.GetBounds()

	trackColor := s.offColor
	if s.checked {
		trackColor = s.onColor
	}

	// Draw rounded track
	r := bounds.Height / 2
	drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, core.BorderRadius{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r}, trackColor)

	// Draw thumb with animation interpolation
	thumbRadius := bounds.Height/2 - 2
	thumbY := bounds.Y + bounds.Height/2
	leftX := bounds.X + thumbRadius + 2
	rightX := bounds.X + bounds.Width - thumbRadius - 2
	thumbX := core.LerpFloat32(leftX, rightX, s.thumbProgress)
	drawFilledCirclePath(screen, thumbX, thumbY, thumbRadius, s.thumbColor)
}

// HandleEvent processes click events.
func (s *Switch) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventMouseDown:
		return true
	case core.EventMouseUp:
		return true
	case core.EventClick:
		s.checked = !s.checked
		if s.onChange != nil {
			s.onChange(s.checked)
		}
		s.startThumbAnimation()
		return true
	}
	return false
}

func (s *Switch) startThumbAnimation() {
	target := float32(0)
	if s.checked {
		target = 1
	}
	if s.thumbAnim != nil {
		s.thumbAnim.Stop()
	}
	startProgress := s.thumbProgress
	s.thumbAnim = core.NewTween(200*time.Millisecond, core.EaseInOutQuad).
		OnUpdate(func(progress float32) {
			s.thumbProgress = core.LerpFloat32(startProgress, target, progress)
			s.Mark(core.FlagNeedDraw)
		})
	s.thumbAnim.Start()
	if engine := s.GetEngine(); engine != nil {
		engine.AddAnimation(s.thumbAnim)
	}
}

// SetChecked sets the switch state.
func (s *Switch) SetChecked(checked bool) *Switch {
	s.checked = checked
	s.thumbProgress = 0
	if checked {
		s.thumbProgress = 1
	}
	s.Mark(core.FlagNeedDraw)
	return s
}

// DebugProps returns visual properties for debugger preview.
func (s *Switch) DebugProps() map[string]interface{} {
	props := make(map[string]interface{})
	trackColor := s.offColor
	if s.checked {
		trackColor = s.onColor
	}
	if trackColor != nil {
		props["backgroundColor"] = colorToCSS(trackColor)
	}
	if s.thumbColor != nil {
		props["thumbColor"] = colorToCSS(s.thumbColor)
	}
	props["checked"] = s.checked
	props["trackWidth"] = s.trackWidth
	props["trackHeight"] = s.trackHeight
	// Switch track is always fully rounded
	if s.trackHeight > 0 {
		props["borderRadius"] = s.trackHeight / 2
	}
	return props
}

// SyncFrom 同步新 Switch 的属性到当前 Element（声明式重建）。
func (s *Switch) SyncFrom(src core.Element) {
	other, ok := src.(*Switch)
	if !ok {
		return
	}
	sb := &SyncBuilder{}
	if s.checked != other.checked {
		s.checked = other.checked
		s.thumbProgress = 0
		if s.checked {
			s.thumbProgress = 1
		}
		sb.NeedDraw = true
	}
	if s.trackWidth != other.trackWidth || s.trackHeight != other.trackHeight {
		s.trackWidth = other.trackWidth
		s.trackHeight = other.trackHeight
		sb.NeedDraw = true
	}
	syncColor(sb, &s.offColor, other.offColor)
	syncColor(sb, &s.onColor, other.onColor)
	syncColor(sb, &s.thumbColor, other.thumbColor)
	sb.MarkDraw(s)
}

// SetOnChange sets the change callback.
func (s *Switch) SetOnChange(fn func(checked bool)) *Switch {
	s.onChange = fn
	return s
}

// SetOffColor sets the track off-color.
func (s *Switch) SetOffColor(clr color.Color) *Switch {
	s.offColor = clr
	s.Mark(core.FlagNeedDraw)
	return s
}

// SetOnColor sets the track on-color.
func (s *Switch) SetOnColor(clr color.Color) *Switch {
	s.onColor = clr
	s.Mark(core.FlagNeedDraw)
	return s
}

// SetThumbColor sets the thumb color.
func (s *Switch) SetThumbColor(clr color.Color) *Switch {
	s.thumbColor = clr
	s.Mark(core.FlagNeedDraw)
	return s
}

// SetTrackSize sets the track width and height.
func (s *Switch) SetTrackSize(width, height float32) *Switch {
	s.trackWidth = width
	s.trackHeight = height
	s.SetWidth(width)
	s.SetHeight(height)
	s.Mark(core.FlagNeedDraw)
	return s
}

