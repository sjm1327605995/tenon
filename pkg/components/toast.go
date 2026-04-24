package components

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ToastType 定义 Toast 类型。
type ToastType string

const (
	ToastInfo    ToastType = "info"
	ToastSuccess ToastType = "success"
	ToastWarning ToastType = "warning"
	ToastError   ToastType = "error"
)

// ToastPosition 定义 Toast 显示位置。
type ToastPosition string

const (
	ToastTop    ToastPosition = "top"
	ToastCenter ToastPosition = "center"
	ToastBottom ToastPosition = "bottom"
)

// Toast 是轻量提示组件。
type Toast struct {
	core.BaseHost
	message     string
	toastType   ToastType
	position    ToastPosition
	textComp    *Text
	duration    time.Duration
	showTime    time.Time
	autoDismiss bool
}

// NewToast 创建一个 Toast 提示组件。
func NewToast(message string) *Toast {
	t := &Toast{
		message:     message,
		toastType:   ToastInfo,
		position:    ToastTop,
		duration:    3 * time.Second,
		autoDismiss: false,
	}
	t.Init(t)
	t.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	t.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
	t.GetElement().Yoga.StyleSetPadding(yoga.EdgeHorizontal, 16)
	t.GetElement().Yoga.StyleSetPadding(yoga.EdgeVertical, 10)
	t.GetElement().Yoga.StyleSetBorderRadius(8)
	t.GetElement().Yoga.StyleSetPositionType(yoga.PositionTypeAbsolute)
	t.GetElement().Yoga.StyleSetPosition(yoga.EdgeLeft, 0)
	t.GetElement().Yoga.StyleSetPosition(yoga.EdgeRight, 0)
	t.applyStyle()

	t.textComp = NewText(message)
	t.textComp.SetFontSize(core.GetTheme().FontSizeBase)
	t.textComp.GetElement().PointerEvents = core.PointerEventsNone
	t.AddChild(t.textComp)

	return t
}

// Show 显示 Toast，若设置了 autoDismiss 则在 duration 后自动隐藏。
func (t *Toast) Show() *Toast {
	t.GetElement().Visible = true
	if t.autoDismiss {
		t.showTime = time.Now()
	}
	return t
}

// Hide 隐藏 Toast。
func (t *Toast) Hide() *Toast {
	t.GetElement().Visible = false
	return t
}

// SetType 设置 Toast 类型，自动应用对应颜色。
func (t *Toast) SetType(tp ToastType) *Toast {
	t.toastType = tp
	t.applyStyle()
	return t
}

// SetPosition 设置显示位置。
func (t *Toast) SetPosition(pos ToastPosition) *Toast {
	t.position = pos
	t.applyPosition()
	return t
}

// SetDuration 设置自动消失时长。
func (t *Toast) SetDuration(d time.Duration) *Toast {
	t.duration = d
	return t
}

// SetAutoDismiss 设置是否自动消失。
func (t *Toast) SetAutoDismiss(auto bool) *Toast {
	t.autoDismiss = auto
	return t
}

// SetMessage 更新提示文本。
func (t *Toast) SetMessage(msg string) *Toast {
	t.message = msg
	if t.textComp != nil {
		t.textComp.SetContent(msg)
	}
	return t
}

func (t *Toast) applyStyle() {
	theme := core.GetTheme()
	switch t.toastType {
	case ToastSuccess:
		t.GetElement().BackgroundColor = color.RGBA{R: 246, G: 255, B: 237, A: 255}
		t.GetElement().BorderColor = color.RGBA{R: 149, G: 222, B: 100, A: 255}
		if t.textComp != nil {
			t.textComp.SetColor(color.RGBA{R: 56, G: 158, B: 13, A: 255})
		}
	case ToastWarning:
		t.GetElement().BackgroundColor = color.RGBA{R: 255, G: 251, B: 230, A: 255}
		t.GetElement().BorderColor = color.RGBA{R: 255, G: 214, B: 102, A: 255}
		if t.textComp != nil {
			t.textComp.SetColor(color.RGBA{R: 212, G: 107, B: 8, A: 255})
		}
	case ToastError:
		t.GetElement().BackgroundColor = color.RGBA{R: 255, G: 241, B: 240, A: 255}
		t.GetElement().BorderColor = color.RGBA{R: 255, G: 163, B: 158, A: 255}
		if t.textComp != nil {
			t.textComp.SetColor(color.RGBA{R: 207, G: 19, B: 34, A: 255})
		}
	default: // info
		t.GetElement().BackgroundColor = color.RGBA{R: 230, G: 244, B: 255, A: 255}
		t.GetElement().BorderColor = color.RGBA{R: 145, G: 202, B: 255, A: 255}
		if t.textComp != nil {
			t.textComp.SetColor(theme.PrimaryColor)
		}
	}
}

func (t *Toast) applyPosition() {
	yogaNode := t.GetElement().Yoga
	switch t.position {
	case ToastTop:
		yogaNode.StyleSetPosition(yoga.EdgeTop, 20)
		yogaNode.StyleSetPosition(yoga.EdgeBottom, -1)
	case ToastCenter:
		yogaNode.StyleSetPosition(yoga.EdgeTop, -1)
		yogaNode.StyleSetPosition(yoga.EdgeBottom, -1)
		yogaNode.StyleSetAlignSelf(yoga.AlignCenter)
	case ToastBottom:
		yogaNode.StyleSetPosition(yoga.EdgeTop, -1)
		yogaNode.StyleSetPosition(yoga.EdgeBottom, 20)
	}
}

// Update 每帧检查自动消失。
func (t *Toast) Update() error {
	if t.autoDismiss && t.GetElement().Visible && !t.showTime.IsZero() {
		if time.Since(t.showTime) > t.duration {
			t.Hide()
		}
	}
	return nil
}

// Draw 绘制 Toast 背景和边框。
func (t *Toast) Draw(screen *ebiten.Image) {
	el := t.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := t.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 背景
	if el.BackgroundColor != nil {
		if hasRadius(el.BorderRadius) {
			t.drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BorderRadius, el.BackgroundColor)
		} else {
			vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
		}
	}

	// 边框
	if el.BorderColor != nil {
		if hasRadius(el.BorderRadius) {
			t.drawRoundedRectStroke(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BorderRadius, 1, el.BorderColor)
		} else {
			vector.StrokeRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, 1, el.BorderColor, false)
		}
	}
}

func (t *Toast) drawRoundedRectFill(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

func (t *Toast) drawRoundedRectStroke(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, stroke float32, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}

// HandleEvent 消费所有事件，防止穿透到下层组件。
func (t *Toast) HandleEvent(e *core.Event) bool {
	return t.GetElement().Visible
}

// ==================== 链式 API ====================

func (t *Toast) SetWidth(width float32) *Toast {
	t.GetElement().Yoga.StyleSetWidth(width)
	return t
}
func (t *Toast) SetMargin(edge yoga.Edge, value float32) *Toast {
	t.GetElement().Yoga.StyleSetMargin(edge, value)
	return t
}
func (t *Toast) SetPadding(edge yoga.Edge, value float32) *Toast {
	t.GetElement().Yoga.StyleSetPadding(edge, value)
	return t
}

// SyncFrom 同步 Toast 属性。
func (t *Toast) SyncFrom(other core.Host) {
	if o, ok := other.(*Toast); ok {
		t.message = o.message
		t.toastType = o.toastType
		t.position = o.position
		t.duration = o.duration
		t.autoDismiss = o.autoDismiss
		if t.textComp != nil && o.textComp != nil {
			t.textComp.Content = o.textComp.Content
			t.textComp.cachedLayout = nil
		}
		t.applyStyle()
	}
}
