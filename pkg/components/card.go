package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// CardFace 表示卡片的朝向。
type CardFace int

const (
	CardFaceFront CardFace = iota
	CardFaceBack
)

// Card 是游戏王卡片组件，支持正反面显示、高亮、ATK/DEF 指示等。
type Card struct {
	core.BaseHost
	face        CardFace
	highlighted bool
	selected    bool
	atk         int
	def         int
	level       int
	overlayNum  int // XYZ 叠放数量
	linkArrows  [8]bool
	onClick     func()
	onHover     func(hovered bool)
	cardImage   *Image
	backImage   *Image
	contentView *View
}

// NewCard 创建一个卡片组件。
func NewCard() *Card {
	c := &Card{
		face: CardFaceFront,
	}
	c.Init(c)
	c.SetFocusable(true)
	c.GetElement().Yoga.StyleSetAspectRatio(0.686) // 游戏王卡片比例 59:86
	c.GetElement().Yoga.StyleSetWidth(80)
	c.GetElement().BorderRadius = core.BorderRadius{TopLeft: 4, TopRight: 4, BottomRight: 4, BottomLeft: 4}

	c.contentView = NewView()
	c.contentView.Init(c.contentView)
	c.contentView.GetElement().Yoga.StyleSetWidthPercent(100)
	c.contentView.GetElement().Yoga.StyleSetHeightPercent(100)
	c.contentView.GetElement().PointerEvents = core.PointerEventsNone
	c.AddChild(c.contentView)

	return c
}

// SetFace 设置卡片朝向。
func (c *Card) SetFace(face CardFace) *Card {
	c.face = face
	return c
}

// SetHighlighted 设置高亮状态。
func (c *Card) SetHighlighted(v bool) *Card {
	c.highlighted = v
	return c
}

// SetSelected 设置选中状态。
func (c *Card) SetSelected(v bool) *Card {
	c.selected = v
	return c
}

// SetATK 设置攻击力。
func (c *Card) SetATK(atk int) *Card {
	c.atk = atk
	return c
}

// SetDEF 设置守备力。
func (c *Card) SetDEF(def int) *Card {
	c.def = def
	return c
}

// SetLevel 设置等级。
func (c *Card) SetLevel(level int) *Card {
	c.level = level
	return c
}

// SetOverlayNum 设置 XYZ 叠放数量。
func (c *Card) SetOverlayNum(n int) *Card {
	c.overlayNum = n
	return c
}

// SetLinkArrows 设置 Link 箭头方向（8 个方向，顺时针从上方开始）。
func (c *Card) SetLinkArrows(arrows [8]bool) *Card {
	c.linkArrows = arrows
	return c
}

// SetCardImage 设置正面卡图。
func (c *Card) SetCardImage(img *Image) *Card {
	if c.cardImage != nil {
		c.contentView.RemoveChild(c.cardImage)
	}
	c.cardImage = img
	if img != nil {
		img.GetElement().PointerEvents = core.PointerEventsNone
		c.contentView.AddChild(img)
	}
	return c
}

// SetBackImage 设置背面卡图。
func (c *Card) SetBackImage(img *Image) *Card {
	c.backImage = img
	return c
}

// SetOnClick 设置点击回调。
func (c *Card) SetOnClick(fn func()) *Card {
	c.onClick = fn
	return c
}

// SetOnHover 设置悬停回调。
func (c *Card) SetOnHover(fn func(hovered bool)) *Card {
	c.onHover = fn
	return c
}

// Draw 绘制卡片背景、边框和高亮效果。
func (c *Card) Draw(screen *ebiten.Image) {
	el := c.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := c.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 背景（卡背或卡面底色）
	bgColor := color.RGBA{R: 30, G: 30, B: 30, A: 255}
	if c.face == CardFaceFront {
		bgColor = color.RGBA{R: 20, G: 20, B: 25, A: 255}
	}
	c.drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BorderRadius, bgColor)

	// 高亮边框
	if c.highlighted || c.selected {
		borderWidth := float32(2)
		borderColor := color.RGBA{R: 255, G: 215, B: 0, A: 255} // 金色高亮
		if c.selected {
			borderColor = color.RGBA{R: 0, G: 150, B: 255, A: 255} // 蓝色选中
		}
		c.drawRoundedRectStroke(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BorderRadius, borderWidth, borderColor)
	}

	// XYZ 叠放指示
	if c.overlayNum > 0 && c.face == CardFaceFront {
		c.drawOverlayIndicator(screen, bounds)
	}

	// Link 箭头
	if c.linkArrows != [8]bool{} && c.face == CardFaceFront {
		c.drawLinkArrows(screen, bounds)
	}
}

func (c *Card) drawOverlayIndicator(screen *ebiten.Image, bounds core.LayoutBounds) {
	badgeW := float32(20)
	badgeH := float32(16)
	badgeX := bounds.X + bounds.Width - badgeW - 4
	badgeY := bounds.Y + 4
	vector.FillRect(screen, badgeX, badgeY, badgeW, badgeH, color.RGBA{R: 128, G: 0, B: 128, A: 255}, false)
}

func (c *Card) drawLinkArrows(screen *ebiten.Image, bounds core.LayoutBounds) {
	cx := bounds.X + bounds.Width/2
	cy := bounds.Y + bounds.Height/2
	r := float32(4)
	arrowColor := color.RGBA{R: 0, G: 200, B: 255, A: 255}

	// 8 个方向的位置偏移
	offsets := [8][2]float32{
		{0, -1},     // 上
		{0.7, -0.7}, // 右上
		{1, 0},      // 右
		{0.7, 0.7},  // 右下
		{0, 1},      // 下
		{-0.7, 0.7}, // 左下
		{-1, 0},     // 左
		{-0.7, -0.7},// 左上
	}

	for i, active := range c.linkArrows {
		if !active {
			continue
		}
		dx := offsets[i][0] * (bounds.Width/2 - 6)
		dy := offsets[i][1] * (bounds.Height/2 - 6)
		vector.FillRect(screen, cx+dx-r, cy+dy-r, r*2, r*2, arrowColor, false)
	}
}

// HandleEvent 处理点击和悬停事件。
func (c *Card) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventMouseEnter:
		if c.onHover != nil {
			c.onHover(true)
		}
		return true
	case core.EventMouseLeave:
		if c.onHover != nil {
			c.onHover(false)
		}
		return true
	case core.EventMouseDown:
		return true
	case core.EventMouseUp:
		return true
	case core.EventClick:
		if c.onClick != nil {
			c.onClick()
		}
		return true
	}
	return false
}

// ==================== 绘制辅助 ====================

func (c *Card) drawRoundedRectFill(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

func (c *Card) drawRoundedRectStroke(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, stroke float32, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}

// ==================== 链式 API ====================

func (c *Card) SetWidth(width float32) *Card {
	c.GetElement().Yoga.StyleSetWidth(width)
	return c
}
func (c *Card) SetHeight(height float32) *Card {
	c.GetElement().Yoga.StyleSetHeight(height)
	return c
}
func (c *Card) SetMargin(edge yoga.Edge, value float32) *Card {
	c.GetElement().Yoga.StyleSetMargin(edge, value)
	return c
}

// SyncFrom 同步卡片属性。
func (c *Card) SyncFrom(other core.Host) {
	if o, ok := other.(*Card); ok {
		c.face = o.face
		c.highlighted = o.highlighted
		c.selected = o.selected
		c.atk = o.atk
		c.def = o.def
		c.level = o.level
		c.overlayNum = o.overlayNum
		c.linkArrows = o.linkArrows
		c.onClick = o.onClick
		c.onHover = o.onHover
		c.cardImage = o.cardImage
		c.backImage = o.backImage
	}
}
