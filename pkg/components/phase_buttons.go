package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// DuelPhase 定义决斗阶段。
type DuelPhase string

const (
	PhaseDP DuelPhase = "DP" // Draw Phase
	PhaseSP DuelPhase = "SP" // Standby Phase
	PhaseM1 DuelPhase = "M1" // Main Phase 1
	PhaseBP DuelPhase = "BP" // Battle Phase
	PhaseBS DuelPhase = "BS" // Battle Step
	PhaseM2 DuelPhase = "M2" // Main Phase 2
	PhaseEP DuelPhase = "EP" // End Phase
)

// PhaseButton 是单个阶段按钮。
type PhaseButton struct {
	core.BaseHost
	phase      DuelPhase
	active     bool
	completed  bool
	onClick    func(phase DuelPhase)
	label      *Text
}

// NewPhaseButton 创建一个阶段按钮。
func NewPhaseButton(phase DuelPhase) *PhaseButton {
	pb := &PhaseButton{phase: phase}
	pb.Init(pb)
	pb.GetElement().Yoga.StyleSetPadding(yoga.EdgeHorizontal, 12)
	pb.GetElement().Yoga.StyleSetPadding(yoga.EdgeVertical, 6)
	pb.GetElement().Yoga.StyleSetBorderRadius(4)
	pb.GetElement().Yoga.StyleSetMinWidth(40)
	pb.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
	pb.GetElement().PointerEvents = core.PointerEventsAuto

	label := NewText(string(phase))
	label.SetFontSize(core.GetTheme().FontSizeSM)
	label.GetElement().PointerEvents = core.PointerEventsNone
	pb.label = label
	pb.AddChild(label)

	pb.refreshStyle()
	return pb
}

// SetActive 设置是否为当前阶段。
func (pb *PhaseButton) SetActive(active bool) *PhaseButton {
	pb.active = active
	pb.refreshStyle()
	return pb
}

// SetCompleted 设置是否已完成。
func (pb *PhaseButton) SetCompleted(completed bool) *PhaseButton {
	pb.completed = completed
	pb.refreshStyle()
	return pb
}

// SetOnClick 设置点击回调。
func (pb *PhaseButton) SetOnClick(fn func(phase DuelPhase)) *PhaseButton {
	pb.onClick = fn
	return pb
}

func (pb *PhaseButton) refreshStyle() {
	theme := core.GetTheme()
	if pb.active {
		pb.GetElement().BackgroundColor = theme.PrimaryColor
		pb.label.SetColor(theme.TextInverseColor)
	} else if pb.completed {
		pb.GetElement().BackgroundColor = theme.SuccessColor
		pb.label.SetColor(theme.TextInverseColor)
	} else {
		pb.GetElement().BackgroundColor = theme.SurfaceColor
		pb.label.SetColor(theme.TextMutedColor)
	}
}

// HandleEvent 处理点击事件。
func (pb *PhaseButton) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventMouseEnter:
		if !pb.active {
			pb.GetElement().BackgroundColor = core.GetTheme().PrimaryLightColor
			pb.label.SetColor(core.GetTheme().PrimaryColor)
		}
		return true
	case core.EventMouseLeave:
		pb.refreshStyle()
		return true
	case core.EventMouseDown:
		return true
	case core.EventMouseUp:
		return true
	case core.EventClick:
		if pb.onClick != nil {
			pb.onClick(pb.phase)
		}
		return true
	}
	return false
}

// PhaseButtons 是阶段按钮组，DP→SP→M1→BP→BS→M2→EP。
type PhaseButtons struct {
	core.BaseHost
	phases     []DuelPhase
	current    DuelPhase
	onChange   func(phase DuelPhase)
	buttons    map[DuelPhase]*PhaseButton
	separator  *Text
}

// NewPhaseButtons 创建阶段按钮组。
func NewPhaseButtons() *PhaseButtons {
	pb := &PhaseButtons{
		phases: []DuelPhase{PhaseDP, PhaseSP, PhaseM1, PhaseBP, PhaseBS, PhaseM2, PhaseEP},
		buttons: make(map[DuelPhase]*PhaseButton),
	}
	pb.Init(pb)
	pb.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	pb.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
	pb.GetElement().Yoga.StyleSetGap(yoga.GutterColumn, 4)

	for i, phase := range pb.phases {
		if i > 0 {
			sep := NewText("→")
			sep.SetFontSize(core.GetTheme().FontSizeSM)
			sep.SetColor(core.GetTheme().TextMutedColor)
			sep.GetElement().PointerEvents = core.PointerEventsNone
			pb.AddChild(sep)
		}
		btn := NewPhaseButton(phase)
		btn.SetOnClick(func(p DuelPhase) {
			pb.SetCurrent(p)
			if pb.onChange != nil {
				pb.onChange(p)
			}
		})
		pb.buttons[phase] = btn
		pb.AddChild(btn)
	}

	return pb
}

// SetCurrent 设置当前阶段。
func (pb *PhaseButtons) SetCurrent(phase DuelPhase) *PhaseButtons {
	pb.current = phase
	found := false
	for _, p := range pb.phases {
		btn, ok := pb.buttons[p]
		if !ok {
			continue
		}
		if p == phase {
			btn.SetActive(true).SetCompleted(false)
			found = true
		} else if found {
			btn.SetActive(false).SetCompleted(false)
		} else {
			btn.SetActive(false).SetCompleted(true)
		}
	}
	return pb
}

// SetOnChange 设置阶段切换回调。
func (pb *PhaseButtons) SetOnChange(fn func(phase DuelPhase)) *PhaseButtons {
	pb.onChange = fn
	return pb
}

// ==================== 链式 API ====================

func (pb *PhaseButtons) SetWidth(width float32) *PhaseButtons {
	pb.GetElement().Yoga.StyleSetWidth(width)
	return pb
}
func (pb *PhaseButtons) SetMargin(edge yoga.Edge, value float32) *PhaseButtons {
	pb.GetElement().Yoga.StyleSetMargin(edge, value)
	return pb
}

// SyncFrom 同步阶段按钮组属性。
func (pb *PhaseButtons) SyncFrom(other core.Host) {
	if o, ok := other.(*PhaseButtons); ok {
		pb.current = o.current
		pb.onChange = o.onChange
		pb.SetCurrent(pb.current)
	}
}
