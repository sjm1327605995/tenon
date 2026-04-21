package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

type View struct {
	*core.BaseComponent
}

func NewView() *View {
	return &View{
		BaseComponent: core.NewBaseComponent(),
	}
}

func (v *View) GetBase() *core.BaseComponent {
	return v.BaseComponent
}

func (v *View) Draw(screen *ebiten.Image) {
	element := v.Render()
	if element == nil || !element.Visible {
		return
	}

	bounds := v.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	v.drawShadow(screen, element)

	if element.BackgroundColor != nil {
		v.drawBackground(screen, element, bounds)
	}

	v.drawBorder(screen, element)

	for _, child := range v.GetChildren() {
		child.Draw(screen)
	}
}

func (v *View) drawBackground(screen *ebiten.Image, element *core.Element, bounds core.LayoutBounds) {
	if element.BorderRadius.TopLeft > 0 || element.BorderRadius.TopRight > 0 ||
		element.BorderRadius.BottomRight > 0 || element.BorderRadius.BottomLeft > 0 {
		v.drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height,
			element.BorderRadius, element.BackgroundColor)
	} else {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, element.BackgroundColor, false)
	}
}

func (v *View) drawRoundedRectFill(screen *ebiten.Image, x, y, width, height float32, radius core.BorderRadius, clr color.Color) {
	var path vector.Path

	path.MoveTo(x+radius.TopLeft, y)
	path.LineTo(x+width-radius.TopRight, y)

	if radius.TopRight > 0 {
		path.Arc(x+width-radius.TopRight, y+radius.TopRight, radius.TopRight, -math.Pi/2, 0, vector.Clockwise)
	} else {
		path.LineTo(x+width, y)
	}

	path.LineTo(x+width, y+height-radius.BottomRight)

	if radius.BottomRight > 0 {
		path.Arc(x+width-radius.BottomRight, y+height-radius.BottomRight, radius.BottomRight, 0, math.Pi/2, vector.Clockwise)
	} else {
		path.LineTo(x+width, y+height)
	}

	path.LineTo(x+radius.BottomLeft, y+height)

	if radius.BottomLeft > 0 {
		path.Arc(x+radius.BottomLeft, y+height-radius.BottomLeft, radius.BottomLeft, math.Pi/2, math.Pi, vector.Clockwise)
	} else {
		path.LineTo(x, y+height)
	}

	path.LineTo(x, y+radius.TopLeft)

	if radius.TopLeft > 0 {
		path.Arc(x+radius.TopLeft, y+radius.TopLeft, radius.TopLeft, math.Pi, 3*math.Pi/2, vector.Clockwise)
	} else {
		path.LineTo(x, y)
	}

	path.Close()

	fillOp := &vector.FillOptions{}
	drawOp := &vector.DrawPathOptions{}
	drawOp.ColorScale.ScaleWithColor(clr)
	drawOp.AntiAlias = true
	vector.FillPath(screen, &path, fillOp, drawOp)
}

func (v *View) AddChild(child core.Component) error {
	return v.BaseComponent.AddChild(child)
}

func (v *View) Add(children ...core.Component) *View {
	for _, child := range children {
		_ = v.BaseComponent.AddChild(child)
	}
	return v
}

func (v *View) drawBorder(screen *ebiten.Image, element *core.Element) {
	if element.BorderColor == nil {
		return
	}

	bounds := v.GetLayoutBounds()
	yogaNode := element.Yoga
	borderTop := yogaNode.StyleGetBorder(yoga.EdgeTop)
	borderRight := yogaNode.StyleGetBorder(yoga.EdgeRight)
	borderBottom := yogaNode.StyleGetBorder(yoga.EdgeBottom)
	borderLeft := yogaNode.StyleGetBorder(yoga.EdgeLeft)

	// 如果有圆角，使用圆角边框
	if element.BorderRadius.TopLeft > 0 || element.BorderRadius.TopRight > 0 ||
		element.BorderRadius.BottomRight > 0 || element.BorderRadius.BottomLeft > 0 {
		// 使用 Path 绘制圆角边框
		v.drawRoundedRectStroke(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height,
			element.BorderRadius, borderTop, element.BorderColor)
		return
	}

	// 否则使用普通矩形边框
	if borderTop > 0 {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, borderTop, element.BorderColor, false)
	}
	if borderRight > 0 {
		vector.FillRect(screen, bounds.X+bounds.Width-borderRight, bounds.Y, borderRight, bounds.Height, element.BorderColor, false)
	}
	if borderBottom > 0 {
		vector.FillRect(screen, bounds.X, bounds.Y+bounds.Height-borderBottom, bounds.Width, borderBottom, element.BorderColor, false)
	}
	if borderLeft > 0 {
		vector.FillRect(screen, bounds.X, bounds.Y, borderLeft, bounds.Height, element.BorderColor, false)
	}
}

func (v *View) drawRoundedRectStroke(screen *ebiten.Image, x, y, width, height float32, radius core.BorderRadius, strokeWidth float32, clr color.Color) {
	// 使用更精确的数学计算来绘制圆角边框
	var path vector.Path

	// 绘制圆角矩形路径，使用精确的弧度值
	path.MoveTo(x+radius.TopLeft, y)

	// 上边
	path.LineTo(x+width-radius.TopRight, y)

	// 右上角 - 使用精确的弧度值
	if radius.TopRight > 0 {
		path.Arc(x+width-radius.TopRight, y+radius.TopRight, radius.TopRight, -math.Pi/2, 0, vector.Clockwise)
	} else {
		path.LineTo(x+width, y)
	}

	// 右边
	path.LineTo(x+width, y+height-radius.BottomRight)

	// 右下角
	if radius.BottomRight > 0 {
		path.Arc(x+width-radius.BottomRight, y+height-radius.BottomRight, radius.BottomRight, 0, math.Pi/2, vector.Clockwise)
	} else {
		path.LineTo(x+width, y+height)
	}

	// 下边
	path.LineTo(x+radius.BottomLeft, y+height)

	// 左下角
	if radius.BottomLeft > 0 {
		path.Arc(x+radius.BottomLeft, y+height-radius.BottomLeft, radius.BottomLeft, math.Pi/2, math.Pi, vector.Clockwise)
	} else {
		path.LineTo(x, y+height)
	}

	// 左边
	path.LineTo(x, y+radius.TopLeft)

	// 左上角
	if radius.TopLeft > 0 {
		path.Arc(x+radius.TopLeft, y+radius.TopLeft, radius.TopLeft, math.Pi, 3*math.Pi/2, vector.Clockwise)
	} else {
		path.LineTo(x, y)
	}

	path.Close()

	// 描边路径 - 启用抗锯齿
	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	strokeOp.MiterLimit = 10

	drawOp := &vector.DrawPathOptions{}
	drawOp.ColorScale.ScaleWithColor(clr)
	drawOp.AntiAlias = true // 启用抗锯齿

	vector.StrokePath(screen, &path, strokeOp, drawOp)
}

func (v *View) drawShadow(screen *ebiten.Image, element *core.Element) {
	if element.ShadowColor == nil || element.ShadowBlur <= 0 {
		return
	}

	bounds := v.GetLayoutBounds()
	shadowWidth := bounds.Width + element.ShadowBlur*2
	shadowHeight := bounds.Height + element.ShadowBlur*2
	shadowX := bounds.X - element.ShadowBlur + element.ShadowOffsetX
	shadowY := bounds.Y - element.ShadowBlur + element.ShadowOffsetY
	vector.FillRect(screen, shadowX, shadowY, shadowWidth, shadowHeight, element.ShadowColor, false)
}

func (v *View) Update() error {
	for _, child := range v.GetChildren() {
		if err := child.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (v *View) DrawOverlay(screen *ebiten.Image) {
	for _, child := range v.GetChildren() {
		child.DrawOverlay(screen)
	}
}

func (v *View) HandleInput() bool {
	children := v.GetChildren()
	for i := len(children) - 1; i >= 0; i-- {
		if children[i].HandleInput() {
			return true
		}
	}

	element := v.Render()
	if element != nil && element.PointerEvents == core.PointerEventsNone {
		return false
	}

	bounds := v.GetLayoutBounds()
	mouseX, mouseY := ebiten.CursorPosition()
	if float32(mouseX) >= bounds.X && float32(mouseX) <= bounds.X+bounds.Width &&
		float32(mouseY) >= bounds.Y && float32(mouseY) <= bounds.Y+bounds.Height {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return true
		}
	}
	return false
}

func (v *View) SetBackgroundColor(c color.Color) *View {
	v.GetElement().BackgroundColor = c
	return v
}

func (v *View) SetBorderColor(c color.Color) *View {
	v.GetElement().BorderColor = c
	return v
}

func (v *View) SetWidth(width float32) *View {
	v.GetElement().Yoga.StyleSetWidth(width)
	return v
}

func (v *View) SetHeight(height float32) *View {
	v.GetElement().Yoga.StyleSetHeight(height)
	return v
}

func (v *View) SetFlexDirection(dir yoga.FlexDirection) *View {
	v.GetElement().Yoga.StyleSetFlexDirection(dir)
	return v
}

func (v *View) SetJustifyContent(justify yoga.Justify) *View {
	v.GetElement().Yoga.StyleSetJustifyContent(justify)
	return v
}

func (v *View) SetAlignItems(align yoga.Align) *View {
	v.GetElement().Yoga.StyleSetAlignItems(align)
	return v
}

func (v *View) SetFlexGrow(grow float32) *View {
	v.GetElement().Yoga.StyleSetFlexGrow(grow)
	return v
}

func (v *View) SetPadding(edge yoga.Edge, value float32) *View {
	v.GetElement().Yoga.StyleSetPadding(edge, value)
	return v
}

func (v *View) SetMargin(edge yoga.Edge, value float32) *View {
	v.GetElement().Yoga.StyleSetMargin(edge, value)
	return v
}

func (v *View) SetBorder(edge yoga.Edge, value float32) *View {
	v.GetElement().Yoga.StyleSetBorder(edge, value)
	return v
}

// SetBorderRadius 设置所有角的圆角半径
func (v *View) SetBorderRadius(radius float32) *View {
	elem := v.GetElement()
	elem.BorderRadius = core.BorderRadius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomRight: radius,
		BottomLeft:  radius,
	}
	return v
}

// SetBorderRadiusTL 设置左上角圆角半径
func (v *View) SetBorderRadiusTL(radius float32) *View {
	elem := v.GetElement()
	elem.BorderRadius.TopLeft = radius
	return v
}

// SetBorderRadiusTR 设置右上角圆角半径
func (v *View) SetBorderRadiusTR(radius float32) *View {
	elem := v.GetElement()
	elem.BorderRadius.TopRight = radius
	return v
}

// SetBorderRadiusBR 设置右下角圆角半径
func (v *View) SetBorderRadiusBR(radius float32) *View {
	elem := v.GetElement()
	elem.BorderRadius.BottomRight = radius
	return v
}

// SetBorderRadiusBL 设置左下角圆角半径
func (v *View) SetBorderRadiusBL(radius float32) *View {
	elem := v.GetElement()
	elem.BorderRadius.BottomLeft = radius
	return v
}

// SetBorderRadius4 设置四个角的圆角半径（CSS 样式）
func (v *View) SetBorderRadius4(topLeft, topRight, bottomRight, bottomLeft float32) *View {
	elem := v.GetElement()
	elem.BorderRadius = core.BorderRadius{
		TopLeft:     topLeft,
		TopRight:    topRight,
		BottomRight: bottomRight,
		BottomLeft:  bottomLeft,
	}
	return v
}

func (v *View) SetShadow(color color.Color, blur, offsetX, offsetY float32) *View {
	elem := v.GetElement()
	elem.ShadowColor = color
	elem.ShadowBlur = blur
	elem.ShadowOffsetX = offsetX
	elem.ShadowOffsetY = offsetY
	return v
}
