package components

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// NoticeType 定义通知类型。
type NoticeType int

const (
	NoticeInfo NoticeType = iota
	NoticeSuccess
	NoticeWarning
	NoticeError
)

// SystemNotice 是系统通知组件，自动消失。
type SystemNotice struct {
	core.BaseHost
	message     string
	noticeType  NoticeType
	duration    time.Duration
	showTime    time.Time
	visible     bool
	icon        *View
	label       *Text
	closeBtn    *Button
}

// NewSystemNotice 创建一个系统通知组件。
func NewSystemNotice() *SystemNotice {
	sn := &SystemNotice{
		duration: 5 * time.Second,
	}
	sn.Init(sn)
	sn.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	sn.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
	sn.GetElement().Yoga.StyleSetPadding(yoga.EdgeHorizontal, 12)
	sn.GetElement().Yoga.StyleSetPadding(yoga.EdgeVertical, 8)
	sn.GetElement().Yoga.StyleSetBorderRadius(6)
	sn.GetElement().Yoga.StyleSetMinWidth(200)
	sn.GetElement().Yoga.StyleSetMaxWidth(400)
	sn.GetElement().Yoga.StyleSetPositionType(yoga.PositionTypeAbsolute)
	sn.GetElement().Yoga.StyleSetPosition(yoga.EdgeTop, 16)
	sn.GetElement().Yoga.StyleSetPosition(yoga.EdgeRight, 16)
	sn.GetElement().ShadowColor = color.RGBA{R: 0, G: 0, B: 0, A: 60}
	sn.GetElement().ShadowBlur = 8
	sn.GetElement().ShadowOffsetY = 2
	sn.GetElement().PointerEvents = core.PointerEventsAuto

	// 图标区域
	icon := NewView()
	icon.Init(icon)
	icon.GetElement().Yoga.StyleSetWidth(20)
	icon.GetElement().Yoga.StyleSetHeight(20)
	icon.GetElement().Yoga.StyleSetMargin(yoga.EdgeRight, 8)
	icon.GetElement().Yoga.StyleSetBorderRadius(10)
	icon.GetElement().PointerEvents = core.PointerEventsNone
	sn.icon = icon
	sn.AddChild(icon)

	// 消息文本
	label := NewText("")
	label.SetFontSize(core.GetTheme().FontSizeBase)
	label.GetElement().Yoga.StyleSetFlexGrow(1)
	label.GetElement().PointerEvents = core.PointerEventsNone
	sn.label = label
	sn.AddChild(label)

	// 关闭按钮
	closeBtn := NewButton("×")
	closeBtn.SetType(ButtonTypeText)
	closeBtn.SetOnClick(func() {
		sn.Hide()
	})
	closeBtn.GetElement().Yoga.StyleSetMargin(yoga.EdgeLeft, 8)
	sn.closeBtn = closeBtn
	sn.AddChild(closeBtn)

	sn.Hide()
	return sn
}

// Show 显示通知。
func (sn *SystemNotice) Show(message string, noticeType NoticeType) *SystemNotice {
	sn.message = message
	sn.noticeType = noticeType
	sn.showTime = time.Now()
	sn.visible = true
	sn.GetElement().Visible = true
	sn.label.SetContent(message)
	sn.refreshStyle()
	return sn
}

// Hide 隐藏通知。
func (sn *SystemNotice) Hide() *SystemNotice {
	sn.visible = false
	sn.GetElement().Visible = false
	return sn
}

// SetDuration 设置显示时长。
func (sn *SystemNotice) SetDuration(d time.Duration) *SystemNotice {
	sn.duration = d
	return sn
}

func (sn *SystemNotice) refreshStyle() {
	switch sn.noticeType {
	case NoticeSuccess:
		sn.GetElement().BackgroundColor = color.RGBA{R: 246, G: 255, B: 237, A: 255}
		sn.GetElement().BorderColor = color.RGBA{R: 149, G: 222, B: 100, A: 255}
		sn.icon.GetElement().BackgroundColor = color.RGBA{R: 56, G: 158, B: 13, A: 255}
		sn.label.SetColor(color.RGBA{R: 56, G: 158, B: 13, A: 255})
	case NoticeWarning:
		sn.GetElement().BackgroundColor = color.RGBA{R: 255, G: 251, B: 230, A: 255}
		sn.GetElement().BorderColor = color.RGBA{R: 255, G: 214, B: 102, A: 255}
		sn.icon.GetElement().BackgroundColor = color.RGBA{R: 212, G: 107, B: 8, A: 255}
		sn.label.SetColor(color.RGBA{R: 212, G: 107, B: 8, A: 255})
	case NoticeError:
		sn.GetElement().BackgroundColor = color.RGBA{R: 255, G: 241, B: 240, A: 255}
		sn.GetElement().BorderColor = color.RGBA{R: 255, G: 163, B: 158, A: 255}
		sn.icon.GetElement().BackgroundColor = color.RGBA{R: 207, G: 19, B: 34, A: 255}
		sn.label.SetColor(color.RGBA{R: 207, G: 19, B: 34, A: 255})
	default:
		sn.GetElement().BackgroundColor = color.RGBA{R: 230, G: 244, B: 255, A: 255}
		sn.GetElement().BorderColor = color.RGBA{R: 145, G: 202, B: 255, A: 255}
		sn.icon.GetElement().BackgroundColor = core.GetTheme().PrimaryColor
		sn.label.SetColor(core.GetTheme().PrimaryColor)
	}
}

// Update 每帧检查自动消失。
func (sn *SystemNotice) Update() error {
	if sn.visible && !sn.showTime.IsZero() {
		if time.Since(sn.showTime) > sn.duration {
			sn.Hide()
		}
	}
	return nil
}

// Draw 绘制通知背景和边框。
func (sn *SystemNotice) Draw(screen *ebiten.Image) {
	el := sn.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := sn.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}
	if el.BorderColor != nil {
		vector.StrokeRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, 1, el.BorderColor, false)
	}
}

// HandleEvent 消费所有事件防止穿透。
func (sn *SystemNotice) HandleEvent(e *core.Event) bool {
	return sn.GetElement().Visible
}

// SyncFrom 同步通知属性。
func (sn *SystemNotice) SyncFrom(other core.Host) {
	if o, ok := other.(*SystemNotice); ok {
		sn.message = o.message
		sn.noticeType = o.noticeType
		sn.duration = o.duration
		if sn.visible {
			sn.label.SetContent(sn.message)
			sn.refreshStyle()
		}
	}
}
