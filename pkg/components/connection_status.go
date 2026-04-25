package components

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ConnectionState 定义连接状态。
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateReconnecting
	StateError
)

// ConnectionStatus 是网络连接状态指示器组件。
type ConnectionStatus struct {
	core.BaseHost
	state       ConnectionState
	label       *Text
	pulseTime   float64
	lastChange  time.Time
}

// NewConnectionStatus 创建一个连接状态指示器。
func NewConnectionStatus() *ConnectionStatus {
	cs := &ConnectionStatus{
		state: StateDisconnected,
	}
	cs.Init(cs)
	cs.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	cs.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
	cs.GetElement().Yoga.StyleSetPadding(yoga.EdgeHorizontal, 12)
	cs.GetElement().Yoga.StyleSetPadding(yoga.EdgeVertical, 6)
	cs.GetElement().Yoga.StyleSetBorderRadius(4)
	cs.GetElement().Yoga.StyleSetGap(yoga.GutterColumn, 6)

	label := NewText("未连接")
	label.SetFontSize(core.GetTheme().FontSizeSM)
	label.GetElement().PointerEvents = core.PointerEventsNone
	cs.label = label
	cs.AddChild(label)

	cs.refreshStyle()
	return cs
}

// SetState 设置连接状态。
func (cs *ConnectionStatus) SetState(state ConnectionState) *ConnectionStatus {
	if cs.state != state {
		cs.state = state
		cs.lastChange = time.Now()
		cs.refreshStyle()
	}
	return cs
}

// GetState 获取当前连接状态。
func (cs *ConnectionStatus) GetState() ConnectionState {
	return cs.state
}

func (cs *ConnectionStatus) refreshStyle() {
	theme := core.GetTheme()
	switch cs.state {
	case StateConnected:
		cs.label.SetContent("已连接")
		cs.label.SetColor(color.RGBA{R: 56, G: 158, B: 13, A: 255})
		cs.GetElement().BackgroundColor = color.RGBA{R: 246, G: 255, B: 237, A: 255}
	case StateConnecting:
		cs.label.SetContent("连接中...")
		cs.label.SetColor(color.RGBA{R: 212, G: 107, B: 8, A: 255})
		cs.GetElement().BackgroundColor = color.RGBA{R: 255, G: 251, B: 230, A: 255}
	case StateReconnecting:
		cs.label.SetContent("重连中...")
		cs.label.SetColor(color.RGBA{R: 207, G: 19, B: 34, A: 255})
		cs.GetElement().BackgroundColor = color.RGBA{R: 255, G: 241, B: 240, A: 255}
	case StateError:
		cs.label.SetContent("连接失败")
		cs.label.SetColor(color.RGBA{R: 207, G: 19, B: 34, A: 255})
		cs.GetElement().BackgroundColor = color.RGBA{R: 255, G: 241, B: 240, A: 255}
	default:
		cs.label.SetContent("未连接")
		cs.label.SetColor(theme.TextMutedColor)
		cs.GetElement().BackgroundColor = theme.SurfaceColor
	}
}

// Update 每帧更新脉冲动画。
func (cs *ConnectionStatus) Update() error {
	if cs.state == StateConnecting || cs.state == StateReconnecting {
		cs.pulseTime += 0.05
		if cs.pulseTime > 2*3.14159 {
			cs.pulseTime -= 2 * 3.14159
		}
	}
	return nil
}

// Draw 绘制状态指示灯。
func (cs *ConnectionStatus) Draw(screen *ebiten.Image) {
	el := cs.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := cs.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 背景
	if el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}

	// 状态指示灯
	indicatorX := bounds.X + 8
	indicatorY := bounds.Y + bounds.Height/2
	radius := float32(4)

	var indicatorColor color.Color
	switch cs.state {
	case StateConnected:
		indicatorColor = color.RGBA{R: 56, G: 158, B: 13, A: 255}
	case StateConnecting, StateReconnecting:
		// 脉冲效果
		alpha := uint8(128 + 127*float32((1+sin(cs.pulseTime))/2))
		indicatorColor = color.RGBA{R: 212, G: 107, B: 8, A: alpha}
	case StateError:
		indicatorColor = color.RGBA{R: 207, G: 19, B: 34, A: 255}
	default:
		indicatorColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	}

	vector.FillRect(screen, indicatorX-radius, indicatorY-radius, radius*2, radius*2, indicatorColor, false)
}

func sin(x float64) float64 {
	// 简单的sin近似
	x = x - float64(int(x/(2*3.14159)))*2*3.14159
	if x < 0 {
		x += 2 * 3.14159
	}
	if x > 3.14159 {
		x = 2*3.14159 - x
	}
	// 使用泰勒展开近似 sin(x) for x in [0, pi]
	x2 := x * x
	return x * (1 - x2/6*(1-x2/20*(1-x2/42)))
}

// ==================== 链式 API ====================

func (cs *ConnectionStatus) SetWidth(width float32) *ConnectionStatus {
	cs.GetElement().Yoga.StyleSetWidth(width)
	return cs
}
func (cs *ConnectionStatus) SetMargin(edge yoga.Edge, value float32) *ConnectionStatus {
	cs.GetElement().Yoga.StyleSetMargin(edge, value)
	return cs
}

// SyncFrom 同步连接状态属性。
func (cs *ConnectionStatus) SyncFrom(other core.Host) {
	if o, ok := other.(*ConnectionStatus); ok {
		cs.SetState(o.state)
	}
}
