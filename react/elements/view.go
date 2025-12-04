// Package elements implements the React-like element system for the Tenon framework.
// Elements are the building blocks of the UI that users directly interact with.
package elements

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image/color"
)

// View represents a container element that can hold child elements and apply styles.
// It extends the components.View with additional element-specific functionality.
//
// Fields:
//   - View: The underlying component view that handles layout and basic rendering
//   - Children: The list of child elements contained within this view
type View struct {
	ElementBase
	Background  color.NRGBA // Background is the background color with alpha channel
	BorderColor color.NRGBA // BorderColor is the border color with alpha channel

	CornerRadii styles.CornerRadius
}
type BorderWidth struct {
	Top, Right, Bottom, Left unit.Dp
}

func (v *View) Paint(gtx layout.Context) {

	//x, y := int(node.LayoutLeft()), int(node.LayoutTop())
	width := v.Node.LayoutWidth()
	height := v.Node.LayoutHeight()

	// 2. 转换边框宽度为像素 (float32)
	// 假设 v.BorderWidths 已经包含了 Dp 值。
	top := v.Node.LayoutBorder(yoga.EdgeTop)
	right := v.Node.LayoutBorder(yoga.EdgeRight)
	bottom := v.Node.LayoutBorder(yoga.EdgeBottom)
	left := v.Node.LayoutBorder(yoga.EdgeLeft)

	// 3. 绘制背景 (内部区域)
	if v.Background.A > 0 {
		var p clip.Path
		p.Begin(gtx.Ops)

		// 绘制内圈，顺时针 (reverse=false)
		drawInnerLoop(&p, width, height, top, right, bottom, left, v.CornerRadii, gtx, false)

		paint.FillShape(gtx.Ops, v.Background, clip.Outline{Path: p.End()}.Op())
	}

	// 4. 绘制边框 (环形区域: 外圈 - 内圈)
	if v.BorderColor.A > 0 && (top > 0 || right > 0 || bottom > 0 || left > 0) {
		var p clip.Path
		p.Begin(gtx.Ops)

		// 顺时针绘制外圈 (形成主体)
		drawOuterLoop(&p, width, height, v.CornerRadii, gtx)

		// 逆时针绘制内圈 (挖空中央区域，形成环形边框)
		drawInnerLoop(&p, width, height, top, right, bottom, left, v.CornerRadii, gtx, true)

		paint.FillShape(gtx.Ops, v.BorderColor, clip.Outline{Path: p.End()}.Op())
	}

	// 5. [可选] 如果你需要裁剪子内容 (Overflow: Hidden)，在这里添加裁剪操作
	// clip.RRect{...}.Push(gtx.Ops)
}

const k = 0.55228475 // 贝塞尔曲线系数

func drawOuterLoop(p *clip.Path, w, h float32, r styles.CornerRadius, gtx layout.Context) {
	tl := r.TopLeft
	tr := r.TopRight
	br := r.BottomRight
	bl := r.BottomLeft

	p.MoveTo(f32.Pt(tl, 0))
	p.LineTo(f32.Pt(w-tr, 0))
	if tr > 0 {
		p.CubeTo(f32.Pt(w-tr*(1-k), 0), f32.Pt(w, tr*(1-k)), f32.Pt(w, tr))
	}
	p.LineTo(f32.Pt(w, h-br))
	if br > 0 {
		p.CubeTo(f32.Pt(w, h-br*(1-k)), f32.Pt(w-br*(1-k), h), f32.Pt(w-br, h))
	}
	p.LineTo(f32.Pt(bl, h))
	if bl > 0 {
		p.CubeTo(f32.Pt(bl*(1-k), h), f32.Pt(0, h-bl*(1-k)), f32.Pt(0, h-bl))
	}
	p.LineTo(f32.Pt(0, tl))
	if tl > 0 {
		p.CubeTo(f32.Pt(0, tl*(1-k)), f32.Pt(tl*(1-k), 0), f32.Pt(tl, 0))
	}
	p.Close()
}

// drawInnerLoop 绘制内圈
// reverse = false: 顺时针 (用于背景)
// reverse = true:  逆时针 (用于边框挖孔)
// drawInnerLoop 绘制内圈 (修正了控制点和半径计算)
// reverse = false: 顺时针 (用于背景)
// reverse = true:  逆时针 (用于边框挖孔)
func drawInnerLoop(p *clip.Path, w, h, top, right, bottom, left float32, r styles.CornerRadius, gtx layout.Context, reverse bool) {
	// 确保 max 函数可用

	// 1. 计算内圈有效半径 R_x' 和 R_y'
	// R_inner = max(0, R_outer - BorderWidth)
	tlRx := max(0, r.TopLeft-left)
	tlRy := max(0, r.TopLeft-top)

	trRx := max(0, r.TopRight-right)
	trRy := max(0, r.TopRight-top)

	brRx := max(0, r.BottomRight-right)
	brRy := max(0, r.BottomRight-bottom)

	blRx := max(0, r.BottomLeft-left)
	blRy := max(0, r.BottomLeft-bottom)

	// 2. 定义关键点 (内圈的四个角点)
	ptTlStart := f32.Pt(left, top+tlRy) // TL角 Y轴起点
	ptTlEnd := f32.Pt(left+tlRx, top)   // TL角 X轴终点 (也是 Top Edge 的起点)

	ptTrStart := f32.Pt(w-right-trRx, top) // TR角 X轴起点 (Top Edge 的终点)
	ptTrEnd := f32.Pt(w-right, top+trRy)   // TR角 Y轴终点 (也是 Right Edge 的起点)

	ptBrStart := f32.Pt(w-right, h-bottom-brRy) // BR角 Y轴起点 (Right Edge 的终点)
	ptBrEnd := f32.Pt(w-right-brRx, h-bottom)   // BR角 X轴终点 (也是 Bottom Edge 的起点)

	ptBlStart := f32.Pt(left+blRx, h-bottom) // BL角 X轴起点 (Bottom Edge 的终点)
	ptBlEnd := f32.Pt(left, h-bottom-blRy)   // BL角 Y轴终点 (也是 Left Edge 的起点)

	// 3. 路径绘制
	if !reverse {
		// --- 顺时针 (用于背景) ---
		p.MoveTo(ptTlEnd) // 始于 Top Edge

		// Top Line & Corner (ptTlEnd -> ptTrStart -> ptTrEnd)
		p.LineTo(ptTrStart)
		if trRx > 0 || trRy > 0 {
			// P1: 向左偏移 trRx*k, P2: 向下偏移 trRy*k
			cp1 := f32.Pt(w-right-trRx+trRx*k, top)
			cp2 := f32.Pt(w-right, top+trRy-trRy*k)
			p.CubeTo(cp1, cp2, ptTrEnd)
		} else {
			p.LineTo(ptTrEnd)
		}

		// Right Line & Corner (ptTrEnd -> ptBrStart -> ptBrEnd)
		p.LineTo(ptBrStart)
		if brRx > 0 || brRy > 0 {
			// P1: 向下偏移 brRy*k, P2: 向左偏移 brRx*k
			cp1 := f32.Pt(w-right, h-bottom-brRy+brRy*k)
			cp2 := f32.Pt(w-right-brRx+brRx*k, h-bottom)
			p.CubeTo(cp1, cp2, ptBrEnd)
		} else {
			p.LineTo(ptBrEnd)
		}

		// Bottom Line & Corner (ptBrEnd -> ptBlStart -> ptBlEnd)
		p.LineTo(ptBlStart)
		if blRx > 0 || blRy > 0 {
			// P1: 向右偏移 blRx*k, P2: 向上偏移 blRy*k
			cp1 := f32.Pt(left+blRx-blRx*k, h-bottom)
			cp2 := f32.Pt(left, h-bottom-blRy+blRy*k)
			p.CubeTo(cp1, cp2, ptBlEnd)
		} else {
			p.LineTo(ptBlEnd)
		}

		// Left Line & Corner (ptBlEnd -> ptTlStart -> ptTlEnd)
		p.LineTo(ptTlStart)
		if tlRx > 0 || tlRy > 0 {
			// P1: 向上偏移 tlRy*k, P2: 向右偏移 tlRx*k
			cp1 := f32.Pt(left, top+tlRy-tlRy*k)
			cp2 := f32.Pt(left+tlRx-tlRx*k, top)
			p.CubeTo(cp1, cp2, ptTlEnd)
		} else {
			p.LineTo(ptTlEnd)
		}
		p.Close()

	} else {
		// --- 逆时针 (用于边框挖孔) ---
		// 路径与顺时针完全相反，从 ptTlEnd 开始，逆时针绕行
		p.MoveTo(ptTlEnd)

		// 1. Top-Left Corner (ptTlEnd -> ptTlStart)
		if tlRx > 0 || tlRy > 0 {
			// P1: 向右偏移 tlRx*k, P2: 向上偏移 tlRy*k
			cp1 := f32.Pt(left+tlRx-tlRx*k, top)
			cp2 := f32.Pt(left, top+tlRy-tlRy*k)
			p.CubeTo(cp1, cp2, ptTlStart)
		} else {
			p.LineTo(ptTlStart)
		}

		// 2. Left Line (ptTlStart -> ptBlEnd)
		p.LineTo(ptBlEnd)

		// 3. Bottom-Left Corner (ptBlEnd -> ptBlStart)
		if blRx > 0 || blRy > 0 {
			// P1: 向上偏移 blRy*k, P2: 向右偏移 blRx*k
			cp1 := f32.Pt(left, h-bottom-blRy+blRy*k)
			cp2 := f32.Pt(left+blRx-blRx*k, h-bottom)
			p.CubeTo(cp1, cp2, ptBlStart)
		} else {
			p.LineTo(ptBlStart)
		}

		// 4. Bottom Line (ptBlStart -> ptBrEnd)
		p.LineTo(ptBrEnd)

		// 5. Bottom-Right Corner (ptBrEnd -> ptBrStart)
		if brRx > 0 || brRy > 0 {
			// P1: 向左偏移 brRx*k, P2: 向下偏移 brRy*k
			cp1 := f32.Pt(w-right-brRx+brRx*k, h-bottom)
			cp2 := f32.Pt(w-right, h-bottom-brRy+brRy*k)
			p.CubeTo(cp1, cp2, ptBrStart)
		} else {
			p.LineTo(ptBrStart)
		}

		// 6. Right Line (ptBrStart -> ptTrEnd)
		p.LineTo(ptTrEnd)

		// 7. Top-Right Corner (ptTrEnd -> ptTrStart)
		if trRx > 0 || trRy > 0 {
			// P1: 向下偏移 trRy*k, P2: 向左偏移 trRx*k
			cp1 := f32.Pt(w-right, top+trRy-trRy*k)
			cp2 := f32.Pt(w-right-trRx+trRx*k, top)
			p.CubeTo(cp1, cp2, ptTrStart)
		} else {
			p.LineTo(ptTrStart)
		}

		// 8. Top Line (ptTrStart -> ptTlEnd)
		p.LineTo(ptTlEnd)

		p.Close()
	}
}

// SetStyle applies the given style to this view element.
// This method delegates to the style's Apply method to apply all style properties.
//
// Parameters:
//   - style: The style configuration to apply to this view
func (v *View) SetStyle(style *styles.Style) {
	style.Apply(v)
}

// Render implements the api.Component interface and returns the view itself as a renderable node.
//
// Returns:
//   - The view as a renderable api.Node
func (v *View) Render() api.Node {
	return v
}

// Style applies the given style to this view and returns the view itself for method chaining.
// This is a convenience method that wraps SetStyle with a return value for fluent API usage.
//
// Parameters:
//   - option: The style configuration to apply to this view
//
// Returns:
//   - The view itself, allowing for method chaining
func (v *View) Style(option *styles.Style) *View {
	v.SetStyle(option)
	return v
}

// Child adds child components to this view and returns the view itself for method chaining.
// Each component is rendered and added to both the Yoga layout tree and the view's children list.
//
// Parameters:
//   - nodes: The components to add as children to this view
//
// Returns:
//   - The view itself, allowing for method chaining
func (v *View) Child(nodes ...api.Component) *View {
	for i := range nodes {
		element := nodes[i].Render()
		v.Yoga().InsertChild(element.Yoga(), uint32(i))
		v.Children = append(v.Children, element)
	}
	return v
}

// SetExtendedStyle applies extended style properties to this view element.
// This method handles specialized style types like BackgroundColor and BorderColor.
//
// Parameters:
//   - extendedStyle: The extended style to apply to this view
func (v *View) SetExtendedStyle(extendedStyle styles.IExtendedStyle) {
	switch e := extendedStyle.(type) {
	case styles.BackgroundColor:
		v.Background = e.Color
	case styles.BorderColor:
		v.BorderColor = e.Color
	case styles.CornerRadius:
		v.CornerRadii = e
	}
}

// NewView creates a new View element with default properties.
//
// Returns:
//   - A pointer to a newly created View element
func NewView() *View {
	return &View{
		ElementBase: ElementBase{
			Node: yoga.NewNode(),
		},
		BorderColor: color.NRGBA{A: 255},
	}
}
