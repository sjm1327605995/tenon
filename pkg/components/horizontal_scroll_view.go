package components

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// HorizontalScrollView 支持水平滚动的容器组件。
type HorizontalScrollView struct {
	core.BaseHost
	content         *View
	scrollX         float32
	maxScrollX      float32
	scrollbarHeight float32
	scrollbarColor  color.Color
	trackColor      color.Color
	dragging        bool
	lastMouseX      float32
	lastMouseY      float32
	velocityX       float32
	flinging        bool
	lastDragX       float32
	lastDragTime    time.Time
}

// NewHorizontalScrollView 创建一个水平滚动视图。
func NewHorizontalScrollView() *HorizontalScrollView {
	theme := core.GetTheme()
	hsv := &HorizontalScrollView{
		scrollbarHeight: 6,
		scrollbarColor:  theme.ScrollbarColor,
		trackColor:      theme.ScrollbarTrackColor,
	}
	hsv.Init(hsv)
	hsv.SetFocusable(true)
	hsv.content = NewView()
	hsv.content.Init(hsv.content)
	hsv.content.GetElement().Yoga.StyleSetHeightPercent(100)
	hsv.AddChild(hsv.content)
	return hsv
}

// Content 返回内容容器，用于添加子组件。
func (hsv *HorizontalScrollView) Content() *View {
	return hsv.content
}

// GetScrollOffset 返回滚动偏移。
func (hsv *HorizontalScrollView) GetScrollOffset() (x, y float32) {
	return hsv.scrollX, 0
}

// ShouldClipChildren 裁剪子组件到自身边界。
func (hsv *HorizontalScrollView) ShouldClipChildren() bool {
	return true
}

// Draw 绘制滚动条轨道和滑块。
func (hsv *HorizontalScrollView) Draw(screen *ebiten.Image) {
	el := hsv.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := hsv.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}

	if hsv.maxScrollX > 0 {
		hsv.drawHorizontalScrollbar(screen, bounds)
	}
}

func (hsv *HorizontalScrollView) drawHorizontalScrollbar(screen *ebiten.Image, bounds core.LayoutBounds) {
	trackX := bounds.X + 2
	trackY := bounds.Y + bounds.Height - hsv.scrollbarHeight - 2
	trackW := bounds.Width - 4

	// 轨道
	vector.FillRect(screen, trackX, trackY, trackW, hsv.scrollbarHeight, hsv.trackColor, false)

	// 滑块
	contentWidth := hsv.content.GetLayoutBounds().Width
	viewportWidth := bounds.Width
	if contentWidth > viewportWidth {
		ratio := viewportWidth / contentWidth
		thumbW := ratio * trackW
		if thumbW < 10 {
			thumbW = 10
		}
		scrollRatio := 0.0
		if hsv.maxScrollX > 0 {
			scrollRatio = float64(-hsv.scrollX / hsv.maxScrollX)
		}
		thumbX := trackX + float32(scrollRatio)*(trackW-thumbW)
		vector.FillRect(screen, thumbX, trackY, thumbW, hsv.scrollbarHeight, hsv.scrollbarColor, false)
	}
}

// HandleEvent 处理滚轮和拖拽事件。
func (hsv *HorizontalScrollView) HandleEvent(e *core.Event) bool {
	bounds := hsv.GetLayoutBounds()

	switch e.Type {
	case core.EventScroll:
		hsv.scrollX += e.DeltaX * 20
		hsv.clampScroll()
		return true

	case core.EventMouseDown:
		hsv.flinging = false
		hsv.velocityX = 0
		if hsv.maxScrollX > 0 {
			trackY := bounds.Y + bounds.Height - hsv.scrollbarHeight - 2
			if e.Y >= trackY && e.Y <= trackY+hsv.scrollbarHeight+4 &&
				e.X >= bounds.X+2 && e.X <= bounds.X+bounds.Width-2 {
				hsv.dragging = true
				hsv.lastMouseX = e.X
				hsv.lastMouseY = e.Y
				hsv.lastDragX = e.X
				hsv.lastDragTime = time.Now()
				return true
			}
		}
		if e.X >= bounds.X && e.X <= bounds.X+bounds.Width &&
			e.Y >= bounds.Y && e.Y <= bounds.Y+bounds.Height {
			hsv.dragging = true
			hsv.lastMouseX = e.X
			hsv.lastMouseY = e.Y
			hsv.lastDragX = e.X
			hsv.lastDragTime = time.Now()
			return true
		}

	case core.EventMouseUp:
		if hsv.dragging {
			hsv.dragging = false
			if math.Abs(float64(hsv.velocityX)) > 50 {
				hsv.flinging = true
			}
		}

	case core.EventMouseMove:
		if hsv.dragging {
			dx := e.X - hsv.lastMouseX
			hsv.scrollX += dx
			hsv.clampScroll()
			now := time.Now()
			dt := float32(now.Sub(hsv.lastDragTime).Seconds())
			if dt > 0 {
				hsv.velocityX = (e.X - hsv.lastDragX) / dt
			}
			hsv.lastDragX = e.X
			hsv.lastDragTime = now
			hsv.lastMouseX = e.X
			hsv.lastMouseY = e.Y
			return true
		}
	}

	return false
}

func (hsv *HorizontalScrollView) clampScroll() {
	if hsv.scrollX > 0 {
		hsv.scrollX = 0
	}
	if hsv.scrollX < -hsv.maxScrollX {
		hsv.scrollX = -hsv.maxScrollX
	}
}

// Update 每帧更新最大滚动范围和 fling 状态。
func (hsv *HorizontalScrollView) Update() error {
	contentW := hsv.content.GetLayoutBounds().Width
	viewportW := hsv.GetLayoutBounds().Width
	if contentW > viewportW {
		hsv.maxScrollX = contentW - viewportW
	} else {
		hsv.maxScrollX = 0
		hsv.scrollX = 0
	}

	if hsv.flinging {
		hsv.scrollX += hsv.velocityX * (1.0 / 60.0)
		hsv.velocityX *= 0.9
		if math.Abs(float64(hsv.velocityX)) < 10 {
			hsv.flinging = false
			hsv.velocityX = 0
		}
	}

	hsv.clampScroll()
	return nil
}

// ==================== 链式 API ====================

func (hsv *HorizontalScrollView) SetWidth(width float32) *HorizontalScrollView {
	hsv.GetElement().Yoga.StyleSetWidth(width)
	return hsv
}
func (hsv *HorizontalScrollView) SetWidthPercent(percent float32) *HorizontalScrollView {
	hsv.GetElement().Yoga.StyleSetWidthPercent(percent)
	return hsv
}
func (hsv *HorizontalScrollView) SetHeight(height float32) *HorizontalScrollView {
	hsv.GetElement().Yoga.StyleSetHeight(height)
	return hsv
}
func (hsv *HorizontalScrollView) SetHeightPercent(percent float32) *HorizontalScrollView {
	hsv.GetElement().Yoga.StyleSetHeightPercent(percent)
	return hsv
}
func (hsv *HorizontalScrollView) SetFlexDirection(dir yoga.FlexDirection) *HorizontalScrollView {
	hsv.GetElement().Yoga.StyleSetFlexDirection(dir)
	return hsv
}
func (hsv *HorizontalScrollView) SetPadding(edge yoga.Edge, value float32) *HorizontalScrollView {
	hsv.GetElement().Yoga.StyleSetPadding(edge, value)
	return hsv
}
func (hsv *HorizontalScrollView) SetMargin(edge yoga.Edge, value float32) *HorizontalScrollView {
	hsv.GetElement().Yoga.StyleSetMargin(edge, value)
	return hsv
}
func (hsv *HorizontalScrollView) SetBackgroundColor(clr color.Color) *HorizontalScrollView {
	hsv.GetElement().BackgroundColor = clr
	return hsv
}
func (hsv *HorizontalScrollView) SetScrollbarColor(clr color.Color) *HorizontalScrollView {
	hsv.scrollbarColor = clr
	return hsv
}
func (hsv *HorizontalScrollView) SetScrollbarHeight(height float32) *HorizontalScrollView {
	hsv.scrollbarHeight = height
	return hsv
}
func (hsv *HorizontalScrollView) SetBorderRadius(radius float32) *HorizontalScrollView {
	hsv.GetElement().BorderRadius = core.BorderRadius{
		TopLeft: radius, TopRight: radius,
		BottomRight: radius, BottomLeft: radius,
	}
	return hsv
}
func (hsv *HorizontalScrollView) SetBorder(edge yoga.Edge, value float32) *HorizontalScrollView {
	hsv.GetElement().Yoga.StyleSetBorder(edge, value)
	return hsv
}
func (hsv *HorizontalScrollView) SetBorderColor(clr color.Color) *HorizontalScrollView {
	hsv.GetElement().BorderColor = clr
	return hsv
}

// SyncFrom 同步水平滚动视图属性。
func (hsv *HorizontalScrollView) SyncFrom(other core.Host) {
	if o, ok := other.(*HorizontalScrollView); ok {
		hsv.content = o.content
		hsv.scrollbarHeight = o.scrollbarHeight
		hsv.scrollbarColor = o.scrollbarColor
		hsv.trackColor = o.trackColor
	}
}
