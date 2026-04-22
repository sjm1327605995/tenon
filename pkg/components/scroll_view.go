package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ScrollView 是支持垂直滚动的容器组件。
type ScrollView struct {
	core.BaseHost
	content        *View
	scrollX        float32
	scrollY        float32
	maxScrollX     float32
	maxScrollY     float32
	scrollbarWidth float32
	scrollbarColor color.Color
	trackColor     color.Color
	dragging       bool
	lastMouseX     float32
	lastMouseY     float32
}

// NewScrollView 创建一个滚动视图。
func NewScrollView() *ScrollView {
	sv := &ScrollView{
		scrollbarWidth: 6,
		scrollbarColor: color.RGBA{R: 150, G: 150, B: 150, A: 200},
		trackColor:     color.RGBA{R: 230, G: 230, B: 230, A: 100},
	}
	sv.Init(sv)
	sv.SetFocusable(true)
	sv.content = NewView()
	sv.content.Init(sv.content)
	sv.content.GetElement().Yoga.StyleSetWidthPercent(100)
	sv.AddChild(sv.content)
	return sv
}

// Content 返回内容容器，用于添加子组件。
func (sv *ScrollView) Content() *View {
	return sv.content
}

// GetScrollOffset 返回滚动偏移。
func (sv *ScrollView) GetScrollOffset() (x, y float32) {
	return sv.scrollX, sv.scrollY
}

// ShouldClipChildren 裁剪子组件到自身边界。
func (sv *ScrollView) ShouldClipChildren() bool {
	return true
}

// Draw 绘制滚动条轨道和滑块。
func (sv *ScrollView) Draw(screen *ebiten.Image) {
	el := sv.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := sv.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 绘制背景
	if el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}

	// 垂直滚动条
	if sv.maxScrollY > 0 {
		sv.drawVerticalScrollbar(screen, bounds)
	}
}

func (sv *ScrollView) drawVerticalScrollbar(screen *ebiten.Image, bounds core.LayoutBounds) {
	trackX := bounds.X + bounds.Width - sv.scrollbarWidth - 2
	trackY := bounds.Y + 2
	trackH := bounds.Height - 4

	// 轨道
	vector.FillRect(screen, trackX, trackY, sv.scrollbarWidth, trackH, sv.trackColor, false)

	// 滑块
	contentHeight := sv.content.GetLayoutBounds().Height
	viewportHeight := bounds.Height
	if contentHeight > viewportHeight {
		ratio := viewportHeight / contentHeight
		thumbH := ratio * trackH
		if thumbH < 10 {
			thumbH = 10
		}
		scrollRatio := -sv.scrollY / sv.maxScrollY
		thumbY := trackY + scrollRatio*(trackH-thumbH)
		vector.FillRect(screen, trackX, thumbY, sv.scrollbarWidth, thumbH, sv.scrollbarColor, false)
	}
}

// HandleEvent 处理滚轮和拖拽事件。
func (sv *ScrollView) HandleEvent(e *core.Event) bool {
	bounds := sv.GetLayoutBounds()

	switch e.Type {
	case core.EventScroll:
		sv.scrollY += e.DeltaY * 20
		sv.clampScroll()
		return true

	case core.EventMouseDown:
		if sv.maxScrollY > 0 {
			// 检查是否点在滚动条上
			trackX := bounds.X + bounds.Width - sv.scrollbarWidth - 2
			trackY := bounds.Y + 2
			trackH := bounds.Height - 4
			if e.X >= trackX && e.X <= trackX+sv.scrollbarWidth+4 &&
				e.Y >= trackY && e.Y <= trackY+trackH {
				sv.dragging = true
				sv.lastMouseX = e.X
				sv.lastMouseY = e.Y
				return true
			}
		}
		// 点在内容区域也捕获拖拽
		if e.X >= bounds.X && e.X <= bounds.X+bounds.Width &&
			e.Y >= bounds.Y && e.Y <= bounds.Y+bounds.Height {
			sv.dragging = true
			sv.lastMouseX = e.X
			sv.lastMouseY = e.Y
			return true
		}

	case core.EventMouseUp:
		sv.dragging = false

	case core.EventMouseMove:
		if sv.dragging {
			dy := e.Y - sv.lastMouseY
			sv.scrollY += dy
			sv.clampScroll()
			sv.lastMouseX = e.X
			sv.lastMouseY = e.Y
			return true
		}
	}

	return false
}

func (sv *ScrollView) clampScroll() {
	if sv.scrollY > 0 {
		sv.scrollY = 0
	}
	if sv.scrollY < -sv.maxScrollY {
		sv.scrollY = -sv.maxScrollY
	}
}

// Update 每帧更新最大滚动范围。
func (sv *ScrollView) Update() error {
	contentH := sv.content.GetLayoutBounds().Height
	viewportH := sv.GetLayoutBounds().Height
	if contentH > viewportH {
		sv.maxScrollY = contentH - viewportH
	} else {
		sv.maxScrollY = 0
		sv.scrollY = 0
	}
	sv.clampScroll()
	return nil
}

// ==================== 链式 API ====================

func (sv *ScrollView) SetWidth(width float32) *ScrollView {
	sv.GetElement().Yoga.StyleSetWidth(width)
	return sv
}
func (sv *ScrollView) SetHeight(height float32) *ScrollView {
	sv.GetElement().Yoga.StyleSetHeight(height)
	return sv
}
func (sv *ScrollView) SetFlexDirection(dir yoga.FlexDirection) *ScrollView {
	sv.GetElement().Yoga.StyleSetFlexDirection(dir)
	return sv
}
func (sv *ScrollView) SetPadding(edge yoga.Edge, value float32) *ScrollView {
	sv.GetElement().Yoga.StyleSetPadding(edge, value)
	return sv
}
func (sv *ScrollView) SetMargin(edge yoga.Edge, value float32) *ScrollView {
	sv.GetElement().Yoga.StyleSetMargin(edge, value)
	return sv
}
func (sv *ScrollView) SetBackgroundColor(clr color.Color) *ScrollView {
	sv.GetElement().BackgroundColor = clr
	return sv
}
func (sv *ScrollView) SetScrollbarColor(clr color.Color) *ScrollView {
	sv.scrollbarColor = clr
	return sv
}
func (sv *ScrollView) SetScrollbarWidth(width float32) *ScrollView {
	sv.scrollbarWidth = width
	return sv
}
func (sv *ScrollView) SetBorderRadius(radius float32) *ScrollView {
	sv.GetElement().BorderRadius = core.BorderRadius{
		TopLeft: radius, TopRight: radius,
		BottomRight: radius, BottomLeft: radius,
	}
	return sv
}
func (sv *ScrollView) SetBorder(edge yoga.Edge, value float32) *ScrollView {
	sv.GetElement().Yoga.StyleSetBorder(edge, value)
	return sv
}
func (sv *ScrollView) SetBorderColor(clr color.Color) *ScrollView {
	sv.GetElement().BorderColor = clr
	return sv
}

// SyncFrom 同步滚动视图属性（保留滚动位置等交互状态）。
func (sv *ScrollView) SyncFrom(other core.Host) {
	if o, ok := other.(*ScrollView); ok {
		sv.content = o.content
		sv.scrollbarWidth = o.scrollbarWidth
		sv.scrollbarColor = o.scrollbarColor
		sv.trackColor = o.trackColor
	}
}
