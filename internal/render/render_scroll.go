package render

import (
	"image/color"

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
	contentDirty  bool
}

func NewRenderScroll() *RenderScroll {
	r := &RenderScroll{}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.clipChildren = true
	r.contentDirty = true
	return r
}

func (r *RenderScroll) AddChild(child RenderObject) {
	r.RenderBox.AddChild(child)
	r.contentDirty = true
}

func (r *RenderScroll) RemoveChild(child RenderObject) {
	r.RenderBox.RemoveChild(child)
	r.contentDirty = true
}

// GetScrollOffset 返回当前的滚动偏移量。
// 正值表示内容向上/向左滚动（视口向下/向右移动）。
func (r *RenderScroll) GetScrollOffset() Offset {
	return Offset{X: r.scrollOffsetX, Y: r.scrollOffsetY}
}

func (r *RenderScroll) MarkNeedsLayout() {
	r.RenderBox.MarkNeedsLayout()
	r.contentDirty = true
}

func (r *RenderScroll) computeContentSize() {
	if !r.contentDirty {
		return
	}
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
	r.contentDirty = false
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
func (r *RenderScroll) Paint(screen *ebiten.Image, offset Offset) {
	r.computeContentSize()
	r.clampScrollOffset()

	// 绘制背景
	r.RenderBox.Paint(screen, offset)

	// 绘制滚动条指示器
	r.drawScrollBar(screen, offset)
}

// drawScrollBar 在内容超出视口时绘制圆角胶囊滚动条，自动避开圆角边界。
func (r *RenderScroll) drawScrollBar(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 计算安全边距：至少 4px，且不小于最大圆角半径，防止 thumb 超出圆角
	br := r.RenderBox.BorderRadius
	maxRadius := br.TopLeft
	if br.TopRight > maxRadius {
		maxRadius = br.TopRight
	}
	if br.BottomLeft > maxRadius {
		maxRadius = br.BottomLeft
	}
	if br.BottomRight > maxRadius {
		maxRadius = br.BottomRight
	}
	margin := float32(4)
	if maxRadius > margin {
		margin = maxRadius
	}

	thumbColor := color.RGBA{R: 150, G: 150, B: 150, A: 180}
	thumbWidth := float32(3)
	thumbRadius := UniformBorderRadius(thumbWidth / 2)

	// 垂直滚动条：不绘制 track，只绘制圆角胶囊 thumb
	if r.contentHeight > bounds.Height {
		minThumbHeight := float32(16)

		thumbX := offset.X + bounds.X + bounds.Width - margin - thumbWidth
		trackY := offset.Y + bounds.Y + margin
		trackH := bounds.Height - margin*2

		// thumb 高度按视口/内容比例，但不少于最小值
		thumbHeight := bounds.Height / r.contentHeight * trackH
		if thumbHeight < minThumbHeight {
			thumbHeight = minThumbHeight
		}
		if thumbHeight > trackH {
			thumbHeight = trackH
		}

		maxScroll := r.contentHeight - bounds.Height
		var thumbY float32
		if maxScroll > 0 {
			thumbY = trackY + r.scrollOffsetY/maxScroll*(trackH-thumbHeight)
		} else {
			thumbY = trackY
		}

		DrawRoundedRectFill(screen, thumbX, thumbY, thumbWidth, thumbHeight, thumbRadius, thumbColor)
	}

	// 水平滚动条
	if r.contentWidth > bounds.Width {
		minThumbWidth := float32(16)

		trackX := offset.X + bounds.X + margin
		trackW := bounds.Width - margin*2
		thumbY := offset.Y + bounds.Y + bounds.Height - margin - thumbWidth

		thumbWidthH := bounds.Width / r.contentWidth * trackW
		if thumbWidthH < minThumbWidth {
			thumbWidthH = minThumbWidth
		}
		if thumbWidthH > trackW {
			thumbWidthH = trackW
		}

		maxScroll := r.contentWidth - bounds.Width
		var thumbX float32
		if maxScroll > 0 {
			thumbX = trackX + r.scrollOffsetX/maxScroll*(trackW-thumbWidthH)
		} else {
			thumbX = trackX
		}

		DrawRoundedRectFill(screen, thumbX, thumbY, thumbWidthH, thumbWidth, thumbRadius, thumbColor)
	}
}
