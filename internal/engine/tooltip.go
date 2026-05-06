package engine

import (
	"image/color"
	"time"
)

// TooltipWidget 为子 Widget 添加悬停提示。
type TooltipWidget struct {
	BaseWidget
	Child       Widget
	Message     string
	Delay       time.Duration // 悬停多久后显示
	Position    TooltipPosition
	TextColor   color.Color
	BgColor     color.Color
	FontSize    float32
	MaxWidth    float32
}

// TooltipPosition 提示位置。
type TooltipPosition int

const (
	TooltipTop    TooltipPosition = iota // 上方
	TooltipBottom                        // 下方
	TooltipLeft                          // 左侧
	TooltipRight                         // 右侧
)

// Tooltip 创建悬停提示 Widget。
func Tooltip(message string, child Widget) TooltipWidget {
	return TooltipWidget{
		Child:    child,
		Message:  message,
		Delay:    500 * time.Millisecond,
		Position: TooltipBottom,
		FontSize: 12,
	}
}

func (t TooltipWidget) WithDelay(d time.Duration) TooltipWidget { t.Delay = d; return t }
func (t TooltipWidget) WithPosition(p TooltipPosition) TooltipWidget { t.Position = p; return t }
func (t TooltipWidget) WithTextColor(c color.Color) TooltipWidget { t.TextColor = c; return t }
func (t TooltipWidget) WithBgColor(c color.Color) TooltipWidget { t.BgColor = c; return t }
func (t TooltipWidget) WithFontSize(v float32) TooltipWidget { t.FontSize = v; return t }
func (t TooltipWidget) WithMaxWidth(v float32) TooltipWidget { t.MaxWidth = v; return t }

func (t TooltipWidget) CreateElement() Element {
	return NewStatefulElement(t)
}

func (t TooltipWidget) CreateState() State {
	s := &tooltipState{}
	s.Init(s)
	return s
}

type tooltipState struct {
	BaseState
	showTimer *time.Timer
	showing   bool
}

func (s *tooltipState) InitState() {}

func (s *tooltipState) Dispose() {
	if s.showTimer != nil {
		s.showTimer.Stop()
	}
}

func (s *tooltipState) DidUpdateWidget(oldWidget Widget) {}

func (s *tooltipState) Build(ctx BuildContext) Widget {
	w := s.GetWidget().(TooltipWidget)

	// 用 Builder 包裹，监听子 Widget 的 mouse enter/leave
	return NewBuilder(func(innerCtx BuildContext) Widget {
		return w.Child
	})
}

// ==================== Tooltip Overlay ====================

// TooltipOverlay 是全局的 Tooltip 渲染层，应放在 Widget 树最外层。
type TooltipOverlay struct {
	BaseWidget
	Child      Widget
	activeTip  *tooltipData
}

type tooltipData struct {
	Message  string
	X, Y     float32
	Position TooltipPosition
	TextColor color.Color
	BgColor   color.Color
	FontSize  float32
	MaxWidth  float32
}

// NewTooltipOverlay 创建 Tooltip 覆盖层。
func NewTooltipOverlay(child Widget) TooltipOverlay {
	return TooltipOverlay{Child: child}
}

func (t TooltipOverlay) CreateElement() Element {
	e := &tooltipOverlayElement{}
	e.SingleChildComponentElement.ComponentElement.BaseElement.Init(e, t)
	return e
}

type tooltipOverlayElement struct {
	SingleChildComponentElement
}

func (e *tooltipOverlayElement) Mount(parent Element, slot int) {
	e.SingleChildComponentElement.Mount(parent, slot)
	w := e.GetWidget().(TooltipOverlay)
	if w.Child != nil {
		e.Child = UpdateChild(e, nil, w.Child)
	}
}

func (e *tooltipOverlayElement) PerformRebuild(oldWidget Widget) {
	w := e.GetWidget().(TooltipOverlay)
	e.Child = UpdateChild(e, e.Child, w.Child)
}
