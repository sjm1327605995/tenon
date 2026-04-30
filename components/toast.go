package components

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// ToastType defines the visual style of a toast.
type ToastType int

const (
	ToastInfo ToastType = iota
	ToastSuccess
	ToastWarning
	ToastError
)

// Toast is a self-dismissing notification bubble.
// Add it to your view tree and call Show() to display messages.
type Toast struct {
	core.BaseElement
	textEl      *native.Text
	toastType   ToastType
	duration    time.Duration
	showTime    time.Time
	visible     bool
	bgColor     color.Color
	borderColor color.Color
	iconColor   color.Color
	radius      float32
	padding     float32
}

// NewToast creates a toast notification.
func NewToast() *Toast {
	_ = core.GetTheme()
	t := &Toast{
		duration: 3 * time.Second,
		radius:   8,
		padding:  12,
	}
	t.Init(t)
	t.SetPositionType(yoga.PositionTypeAbsolute)
	t.SetVisible(false)
	t.SetPadding(yoga.EdgeAll, t.padding)
	t.SetFlexDirection(yoga.FlexDirectionRow)
	t.SetAlignItems(yoga.AlignCenter)
	t.SetGap(yoga.GutterAll, 8)

	t.textEl = native.NewText("").SetColor(color.White)
	t.textEl.SetFontSize(14)
	t.AppendChild(t.textEl)

	t.setTypeColors(ToastInfo)
	return t
}

// ElementType returns type identifier.
func (t *Toast) ElementType() string { return "Toast" }

// Draw renders toast background and optional icon indicator.
func (t *Toast) Draw(screen *ebiten.Image) {
	bounds := t.GetBounds()

	// Background
	native.DrawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height,
		core.BorderRadius{TopLeft: t.radius, TopRight: t.radius, BottomRight: t.radius, BottomLeft: t.radius},
		t.bgColor)

	// Border stroke
	if t.borderColor != nil {
		var path vector.Path
		native.BuildRoundedRectPath(&path, bounds.X, bounds.Y, bounds.Width, bounds.Height,
			core.BorderRadius{TopLeft: t.radius, TopRight: t.radius, BottomRight: t.radius, BottomLeft: t.radius})
		strokeOp := &vector.StrokeOptions{Width: 1, MiterLimit: 10}
		op := &vector.DrawPathOptions{}
		op.ColorScale.ScaleWithColor(t.borderColor)
		op.AntiAlias = true
		vector.StrokePath(screen, &path, strokeOp, op)
	}

	// Icon dot
	iconSize := float32(6)
	iconX := bounds.X + t.padding + iconSize
	iconY := bounds.Y + bounds.Height/2
	native.DrawFilledCirclePath(screen, iconX, iconY, iconSize, t.iconColor)
}

// Update handles auto-dismiss.
func (t *Toast) Update() error {
	if !t.visible {
		return nil
	}
	if time.Since(t.showTime) > t.duration {
		t.Hide()
	}
	return nil
}

// Show displays the toast with a message.
func (t *Toast) Show(message string) *Toast {
	return t.ShowTyped(message, ToastInfo)
}

// ShowTyped displays the toast with a specific type.
func (t *Toast) ShowTyped(message string, toastType ToastType) *Toast {
	t.textEl.SetContent(message)
	t.toastType = toastType
	t.setTypeColors(toastType)
	t.visible = true
	t.showTime = time.Now()
	t.SetVisible(true)
	t.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return t
}

// Hide immediately dismisses the toast.
func (t *Toast) Hide() *Toast {
	t.visible = false
	t.SetVisible(false)
	t.Mark(core.FlagNeedDraw)
	return t
}

// SetDuration sets the auto-dismiss duration.
func (t *Toast) SetDuration(d time.Duration) *Toast {
	t.duration = d
	return t
}

// SyncFrom 同步新 Toast 的属性到当前 Element（声明式重建）。
func (t *Toast) SyncFrom(src core.Element) {
	other, ok := src.(*Toast)
	if !ok {
		return
	}
	sb := &core.SyncBuilder{}
	if t.visible != other.visible {
		t.visible = other.visible
		t.SetVisible(other.visible)
		sb.NeedDraw = true
	}
	if t.toastType != other.toastType {
		t.toastType = other.toastType
		t.setTypeColors(t.toastType)
		sb.NeedDraw = true
	}
	core.SyncField(sb, &t.duration, other.duration)
	if t.radius != other.radius || t.padding != other.padding {
		t.radius = other.radius
		t.padding = other.padding
		sb.NeedDraw = true
	}
	sb.MarkDraw(t)
}

func (t *Toast) setTypeColors(tt ToastType) {
	theme := core.GetTheme()
	switch tt {
	case ToastSuccess:
		t.bgColor = color.RGBA{R: 82, G: 196, B: 26, A: 230}
		t.iconColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		t.borderColor = color.RGBA{R: 82, G: 196, B: 26, A: 255}
	case ToastWarning:
		t.bgColor = color.RGBA{R: 250, G: 173, B: 20, A: 230}
		t.iconColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		t.borderColor = color.RGBA{R: 250, G: 173, B: 20, A: 255}
	case ToastError:
		t.bgColor = theme.DestructiveColor
		if c, ok := t.bgColor.(color.RGBA); ok {
			t.bgColor = color.RGBA{R: c.R, G: c.G, B: c.B, A: 230}
		}
		t.iconColor = theme.DestructiveForegroundColor
		t.borderColor = theme.DestructiveColor
	default: // Info
		t.bgColor = theme.PrimaryColor
		if c, ok := t.bgColor.(color.RGBA); ok {
			t.bgColor = color.RGBA{R: c.R, G: c.G, B: c.B, A: 230}
		}
		t.iconColor = theme.PrimaryForegroundColor
		t.borderColor = theme.PrimaryColor
	}
	if t.textEl != nil {
		t.textEl.SetColor(theme.PrimaryForegroundColor)
	}
	t.Mark(core.FlagNeedDraw)
}

// ToastManager is a container that can show multiple toasts stacked.
type ToastManager struct {
	core.BaseElement
	toasts []*Toast
	gap    float32
}

// NewToastManager creates a toast manager anchored to a corner.
func NewToastManager() *ToastManager {
	tm := &ToastManager{
		gap: 8,
	}
	tm.Init(tm)
	tm.SetPositionType(yoga.PositionTypeAbsolute)
	tm.SetFlexDirection(yoga.FlexDirectionColumn)
	tm.SetAlignItems(yoga.AlignFlexEnd)
	tm.SetGap(yoga.GutterAll, tm.gap)
	return tm
}

// ElementType returns type identifier.
func (tm *ToastManager) ElementType() string { return "ToastManager" }

// Show adds a new toast to the stack.
func (tm *ToastManager) Show(message string, duration time.Duration) *Toast {
	t := NewToast().SetDuration(duration)
	t.Show(message)
	tm.toasts = append(tm.toasts, t)
	tm.AppendChild(t)
	tm.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return t
}

// Update cleans up hidden toasts.
func (tm *ToastManager) Update() error {
	alive := tm.toasts[:0]
	for _, t := range tm.toasts {
		_ = t.Update()
		if t.visible {
			alive = append(alive, t)
		} else {
			tm.RemoveChild(t)
		}
	}
	if len(alive) != len(tm.toasts) {
		tm.toasts = alive
		tm.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	}
	return nil
}

