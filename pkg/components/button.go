package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ButtonState 表示按钮的交互状态。
type ButtonState int

const (
	ButtonStateNormal ButtonState = iota
	ButtonStateHover
	ButtonStatePressed
)

// Button 是按钮宿主组件。
type Button struct {
	core.BaseHost
	text          *Text
	leftIcon      core.Component
	rightIcon     core.Component
	state         ButtonState
	onClick       func()
	normalColor   color.Color
	hoverColor    color.Color
	pressedColor  color.Color
	normalBorder  color.Color
	hoverBorder   color.Color
	pressedBorder color.Color
	normalText    color.Color
	hoverText     color.Color
	pressedText   color.Color
	disabled      bool
	dashed        bool
	loading       bool
	loadingAngle  float32
}

// NewButton 创建一个按钮。
func NewButton(label string) *Button {
	theme := core.GetTheme()
	b := &Button{
		state:        ButtonStateNormal,
		normalColor:  theme.ButtonNormalColor,
		hoverColor:   theme.ButtonHoverColor,
		pressedColor: theme.ButtonPressedColor,
		normalText:   theme.ButtonTextColor,
		hoverText:    theme.ButtonTextColor,
		pressedText:  theme.ButtonTextColor,
	}
	b.Init(b)
	b.SetFocusable(true)
	b.SetPadding(yoga.EdgeAll, 12)
	b.SetBorderRadius(theme.ButtonBorderRadius)
	b.SetBackgroundColor(b.normalColor)
	b.SetJustifyContent(yoga.JustifyCenter)
	b.SetAlignItems(yoga.AlignCenter)
	b.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	b.GetElement().Yoga.StyleSetGap(yoga.GutterColumn, 4)

	b.text = NewText(label)
	b.text.SetColor(theme.ButtonTextColor)
	b.text.SetFontSize(theme.FontSizeBase + 2)
	b.text.GetElement().PointerEvents = core.PointerEventsNone
	b.AddChild(b.text)

	return b
}

// Draw 绘制按钮背景、圆角和边框。
func (b *Button) Draw(screen *ebiten.Image) {
	el := b.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := b.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 绘制背景
	if el.BackgroundColor != nil {
		if hasRadius(el.BorderRadius) {
			b.drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BorderRadius, el.BackgroundColor)
		} else {
			vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
		}
	}

	// 绘制边框
	if el.BorderColor != nil {
		b.drawButtonBorder(screen, el, bounds)
	}

	// 绘制 loading 动画（仅在纯图标/无图标模式下中央显示）
	if b.loading && b.leftIcon == nil && b.rightIcon == nil {
		b.drawLoading(screen, bounds)
	}
}

func isBorderDefined(v float32) bool {
	return !math.IsNaN(float64(v)) && v > 0
}

func (b *Button) drawButtonBorder(screen *ebiten.Image, el *core.Element, bounds core.LayoutBounds) {
	yogaNode := el.Yoga
	if yogaNode == nil {
		return
	}
	borderTop := yogaNode.StyleGetBorder(yoga.EdgeTop)
	borderRight := yogaNode.StyleGetBorder(yoga.EdgeRight)
	borderBottom := yogaNode.StyleGetBorder(yoga.EdgeBottom)
	borderLeft := yogaNode.StyleGetBorder(yoga.EdgeLeft)
	if !isBorderDefined(borderTop) && !isBorderDefined(borderRight) && !isBorderDefined(borderBottom) && !isBorderDefined(borderLeft) {
		return
	}

	if hasRadius(el.BorderRadius) {
		maxBorder := max(borderTop, max(borderRight, max(borderBottom, borderLeft)))
		b.drawRoundedRectStroke(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BorderRadius, maxBorder, el.BorderColor)
		return
	}

	if borderTop > 0 {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, borderTop, el.BorderColor, false)
	}
	if borderRight > 0 {
		vector.FillRect(screen, bounds.X+bounds.Width-borderRight, bounds.Y, borderRight, bounds.Height, el.BorderColor, false)
	}
	if borderBottom > 0 {
		vector.FillRect(screen, bounds.X, bounds.Y+bounds.Height-borderBottom, bounds.Width, borderBottom, el.BorderColor, false)
	}
	if borderLeft > 0 {
		vector.FillRect(screen, bounds.X, bounds.Y, borderLeft, bounds.Height, el.BorderColor, false)
	}
}



func (b *Button) drawRoundedRectFill(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

func (b *Button) drawRoundedRectStroke(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, stroke float32, clr color.Color) {
	if b.dashed {
		b.drawDashedRoundedRectStroke(screen, x, y, w, h, r, stroke, clr)
		return
	}

	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}

func (b *Button) drawDashedRoundedRectStroke(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, stroke float32, clr color.Color) {
	dash := float32(4)
	gap := float32(4)
	pattern := dash + gap

	// 绘制四条边的虚线（包含圆角部分）
	b.drawDashedSegment(screen, x, y+r.TopLeft, x, y+h-r.BottomLeft, stroke, dash, gap, clr)         // 左
	b.drawDashedSegment(screen, x+w, y+r.TopRight, x+w, y+h-r.BottomRight, stroke, dash, gap, clr)   // 右
	b.drawDashedSegment(screen, x+r.TopLeft, y, x+w-r.TopRight, y, stroke, dash, gap, clr)           // 上
	b.drawDashedSegment(screen, x+r.BottomLeft, y+h, x+w-r.BottomRight, y+h, stroke, dash, gap, clr) // 下

	// 四个圆角用实线弧线衔接
	var path vector.Path
	if r.TopLeft > 0 {
		path.MoveTo(x, y+r.TopLeft)
		path.Arc(x+r.TopLeft, y+r.TopLeft, r.TopLeft, math.Pi, 3*math.Pi/2, vector.Clockwise)
	}
	if r.TopRight > 0 {
		path.MoveTo(x+w-r.TopRight, y)
		path.Arc(x+w-r.TopRight, y+r.TopRight, r.TopRight, -math.Pi/2, 0, vector.Clockwise)
	}
	if r.BottomRight > 0 {
		path.MoveTo(x+w, y+h-r.BottomRight)
		path.Arc(x+w-r.BottomRight, y+h-r.BottomRight, r.BottomRight, 0, math.Pi/2, vector.Clockwise)
	}
	if r.BottomLeft > 0 {
		path.MoveTo(x+r.BottomLeft, y+h)
		path.Arc(x+r.BottomLeft, y+h-r.BottomLeft, r.BottomLeft, math.Pi/2, math.Pi, vector.Clockwise)
	}
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)

	_ = pattern
}

func (b *Button) drawDashedSegment(screen *ebiten.Image, x1, y1, x2, y2, stroke, dash, gap float32, clr color.Color) {
	if x1 == x2 && y1 == y2 {
		return
	}
	isVertical := x1 == x2
	length := float32(0)
	if isVertical {
		length = y2 - y1
	} else {
		length = x2 - x1
	}
	if length < 0 {
		length = -length
	}
	if length <= 0 {
		return
	}

	pos := float32(0)
	drawing := true
	segment := dash

	for pos < length {
		if drawing {
			if pos+segment > length {
				segment = length - pos
			}
			if isVertical {
				// 垂直线，居中于 x
				sy := min(y1, y2) + pos
				vector.FillRect(screen, x1-stroke/2, sy, stroke, segment, clr, false)
			} else {
				// 水平线，居中于 y
				sx := min(x1, x2) + pos
				vector.FillRect(screen, sx, y1-stroke/2, segment, stroke, clr, false)
			}
		}
		pos += segment
		drawing = !drawing
		if drawing {
			segment = dash
		} else {
			segment = gap
		}
	}
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// Update 每帧更新 loading 动画。
func (b *Button) Update() error {
	if b.loading {
		b.loadingAngle += 0.15
		if b.loadingAngle > 2*math.Pi {
			b.loadingAngle -= 2 * math.Pi
		}
	}
	return nil
}

// HandleEvent 处理按钮交互事件。
func (b *Button) HandleEvent(e *core.Event) bool {
	if b.disabled {
		return false
	}

	switch e.Type {
	case core.EventMouseEnter:
		if b.state != ButtonStatePressed {
			b.state = ButtonStateHover
			b.refreshColor()
		}
		return true
	case core.EventMouseLeave:
		if b.state != ButtonStatePressed {
			b.state = ButtonStateNormal
			b.refreshColor()
		}
		return true
	case core.EventMouseDown:
		b.state = ButtonStatePressed
		b.refreshColor()
		return true
	case core.EventMouseUp:
		b.state = ButtonStateNormal
		b.refreshColor()
		return true
	case core.EventClick:
		if b.onClick != nil {
			b.onClick()
		}
		return true
	}
	return false
}

func (b *Button) refreshColor() {
	switch b.state {
	case ButtonStateHover:
		b.SetBackgroundColor(b.hoverColor)
		b.SetBorderColor(b.hoverBorder)
		b.SetTextColor(b.hoverText)
	case ButtonStatePressed:
		b.SetBackgroundColor(b.pressedColor)
		b.SetBorderColor(b.pressedBorder)
		b.SetTextColor(b.pressedText)
	default:
		b.SetBackgroundColor(b.normalColor)
		b.SetBorderColor(b.normalBorder)
		b.SetTextColor(b.normalText)
	}
}

// ==================== 链式 API ====================

func (b *Button) SetOnClick(fn func()) *Button {
	b.onClick = fn
	return b
}
func (b *Button) SetDisabled(disabled bool) *Button {
	b.disabled = disabled
	if disabled {
		b.SetBackgroundColor(color.RGBA{R: 108, G: 117, B: 125, A: 255})
	} else {
		b.SetBackgroundColor(b.normalColor)
	}
	return b
}
func (b *Button) SetText(text string) *Button {
	b.text.SetContent(text)
	return b
}
func (b *Button) SetTextColor(clr color.Color) *Button {
	b.text.SetColor(clr)
	return b
}
func (b *Button) SetBackgroundColors(normal, hover, pressed color.Color) *Button {
	b.normalColor = normal
	b.hoverColor = hover
	b.pressedColor = pressed
	b.SetBackgroundColor(b.normalColor)
	return b
}
func (b *Button) SetBorderColors(normal, hover, pressed color.Color) *Button {
	b.normalBorder = normal
	b.hoverBorder = hover
	b.pressedBorder = pressed
	b.SetBorderColor(b.normalBorder)
	return b
}
func (b *Button) SetTextColors(normal, hover, pressed color.Color) *Button {
	b.normalText = normal
	b.hoverText = hover
	b.pressedText = pressed
	b.SetTextColor(b.normalText)
	return b
}
func (b *Button) SetDashed(dashed bool) *Button {
	b.dashed = dashed
	return b
}
func (b *Button) SetLoading(loading bool) *Button {
	b.loading = loading
	return b
}
func (b *Button) SetLeftIcon(icon core.Component) *Button {
	if b.leftIcon != nil {
		b.RemoveChild(b.leftIcon)
	}
	b.leftIcon = icon
	if icon != nil {
		// 将图标插入到 text 之前（如果 text 有内容）
		b.RemoveChild(b.text)
		b.AddChild(icon)
		if b.text != nil && b.text.Content != "" {
			b.AddChild(b.text)
		}
		if b.rightIcon != nil {
			b.RemoveChild(b.rightIcon)
			b.AddChild(b.rightIcon)
		}
	}
	return b
}
func (b *Button) SetRightIcon(icon core.Component) *Button {
	if b.rightIcon != nil {
		b.RemoveChild(b.rightIcon)
	}
	b.rightIcon = icon
	if icon != nil {
		b.AddChild(icon)
	}
	return b
}

// drawLoading 在按钮中央绘制旋转的 loading 圆弧。
func (b *Button) drawLoading(screen *ebiten.Image, bounds core.LayoutBounds) {
	cx := bounds.X + bounds.Width/2
	cy := bounds.Y + bounds.Height/2
	r := float32(6)
	stroke := float32(2)
	clr := b.text.Color
	if clr == nil {
		clr = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}

	var path vector.Path
	startAngle := b.loadingAngle
	endAngle := b.loadingAngle + 1.2
	path.MoveTo(cx+r*float32(math.Cos(float64(startAngle))), cy+r*float32(math.Sin(float64(startAngle))))
	path.Arc(cx, cy, r, startAngle, endAngle, vector.Clockwise)
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	vector.StrokePath(screen, &path, strokeOp, op)
}
func (b *Button) SetWidth(width float32) *Button {
	b.GetElement().Yoga.StyleSetWidth(width)
	return b
}
func (b *Button) SetWidthPercent(percent float32) *Button {
	b.GetElement().Yoga.StyleSetWidthPercent(percent)
	return b
}
func (b *Button) SetHeight(height float32) *Button {
	b.GetElement().Yoga.StyleSetHeight(height)
	return b
}
func (b *Button) SetMargin(edge yoga.Edge, value float32) *Button {
	b.GetElement().Yoga.StyleSetMargin(edge, value)
	return b
}
func (b *Button) SetPadding(edge yoga.Edge, value float32) *Button {
	b.GetElement().Yoga.StyleSetPadding(edge, value)
	return b
}
func (b *Button) SetJustifyContent(justify yoga.Justify) *Button {
	b.GetElement().Yoga.StyleSetJustifyContent(justify)
	return b
}
func (b *Button) SetAlignItems(align yoga.Align) *Button {
	b.GetElement().Yoga.StyleSetAlignItems(align)
	return b
}
func (b *Button) SetBackgroundColor(clr color.Color) *Button {
	b.GetElement().BackgroundColor = clr
	return b
}
func (b *Button) SetBorderColor(clr color.Color) *Button {
	b.GetElement().BorderColor = clr
	return b
}
func (b *Button) SetBorderRadius(radius float32) *Button {
	b.GetElement().BorderRadius = core.BorderRadius{
		TopLeft: radius, TopRight: radius,
		BottomRight: radius, BottomLeft: radius,
	}
	return b
}

// SyncFrom 同步按钮属性。
func (b *Button) SyncFrom(other core.Host) {
	if o, ok := other.(*Button); ok {
		b.normalColor = o.normalColor
		b.hoverColor = o.hoverColor
		b.pressedColor = o.pressedColor
		b.normalBorder = o.normalBorder
		b.hoverBorder = o.hoverBorder
		b.pressedBorder = o.pressedBorder
		b.normalText = o.normalText
		b.hoverText = o.hoverText
		b.pressedText = o.pressedText
		b.dashed = o.dashed
		b.disabled = o.disabled
		b.loading = o.loading
		b.onClick = o.onClick
		b.leftIcon = o.leftIcon
		b.rightIcon = o.rightIcon
		if b.text != nil && o.text != nil {
			b.text.Content = o.text.Content
			b.text.cachedLayout = nil
		}
		// 同步当前状态的颜色到 Element，避免 Widget 更新后颜色丢失
		b.refreshColor()
	}
}
