package render

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderScroll 是滚动容器的 RenderObject。
// 内容可以超出视口，通过 scrollOffset 控制显示区域。
type RenderScroll struct {
	RenderBox
	scrollOffsetX float32
	scrollOffsetY float32
	contentWidth  float32
	contentHeight float32
}

func NewRenderScroll() *RenderScroll {
	r := &RenderScroll{}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.clipChildren = true
	return r
}

// GetScrollOffset 返回当前的滚动偏移量。
// 正值表示内容向上/向左滚动（视口向下/向右移动）。
func (r *RenderScroll) GetScrollOffset() Offset {
	return Offset{X: r.scrollOffsetX, Y: r.scrollOffsetY}
}

func (r *RenderScroll) computeContentSize() {
	var maxRight, maxBottom float32
	for _, child := range r.GetChildren() {
		b := child.GetBounds()
		right := b.X + b.Width
		bottom := b.Y + b.Height
		if right > maxRight {
			maxRight = right
		}
		if bottom > maxBottom {
			maxBottom = bottom
		}
	}
	r.contentWidth = maxRight
	r.contentHeight = maxBottom
}

// ScrollBy 按增量滚动。
func (r *RenderScroll) ScrollBy(dx, dy float32) {
	r.computeContentSize()
	r.scrollOffsetX += dx
	r.scrollOffsetY += dy
	r.clampScrollOffset()
	r.MarkNeedsPaint()
}

// ScrollTo 滚动到指定位置。
func (r *RenderScroll) ScrollTo(x, y float32) {
	r.computeContentSize()
	r.scrollOffsetX = x
	r.scrollOffsetY = y
	r.clampScrollOffset()
	r.MarkNeedsPaint()
}

func (r *RenderScroll) clampScrollOffset() {
	maxX := r.contentWidth - r.bounds.Width
	maxY := r.contentHeight - r.bounds.Height
	if maxX < 0 {
		maxX = 0
	}
	if maxY < 0 {
		maxY = 0
	}
	if r.scrollOffsetX < 0 {
		r.scrollOffsetX = 0
	}
	if r.scrollOffsetX > maxX {
		r.scrollOffsetX = maxX
	}
	if r.scrollOffsetY < 0 {
		r.scrollOffsetY = 0
	}
	if r.scrollOffsetY > maxY {
		r.scrollOffsetY = maxY
	}
}

// Paint 绘制滚动容器。
// 在绘制时计算内容尺寸并钳制滚动偏移。
func (r *RenderScroll) Paint(screen *ebiten.Image, offset Offset) {
	// 计算内容尺寸（此时 bounds 已同步为最新 Yoga 结果）
	var maxRight, maxBottom float32
	for _, child := range r.GetChildren() {
		b := child.GetBounds()
		right := b.X + b.Width
		bottom := b.Y + b.Height
		if right > maxRight {
			maxRight = right
		}
		if bottom > maxBottom {
			maxBottom = bottom
		}
	}
	r.contentWidth = maxRight
	r.contentHeight = maxBottom
	r.clampScrollOffset()

	// 绘制背景（RenderBox.Paint 只绘制背景和边框，不绘制子节点）
	r.RenderBox.Paint(screen, offset)
}
