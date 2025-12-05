package dom

import (
	"bytes"
	"gioui.org/op"
	"github.com/tdewolff/canvas"
	"image"
	"image/color"
	"math"
	"os"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/millken/yoga"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/tdewolff/canvas/renderers/gio"
)

type Widget interface {
	Layout(gtx layout.Context) layout.Dimensions
}

var Metric unit.Metric

// INode represents a real DOM node that can be rendered.
type INode interface {
	styles.StyleElement
	GetChildren() []INode
	SetChildren(children []INode)
	Paint(gtx layout.Context)
	ApplyProps(vnode *core.VNode)
}

// ElementBase provides a base implementation for an INode.
type ElementBase struct {
	Node     *yoga.Node
	Children []INode
}

func (e *ElementBase) GetYogaNode() *yoga.Node {
	return e.Node
}

func (e *ElementBase) GetChildren() []INode {
	return e.Children
}

func (e *ElementBase) SetChildren(children []INode) {
	e.Children = children
}

// View is a concrete INode for a container.
type View struct {
	ElementBase
	Background  color.NRGBA
	BorderColor color.NRGBA
	CornerRadii styles.CornerRadius
}

func NewView() *View {
	return &View{
		ElementBase: ElementBase{Node: yoga.NewNode()},
		BorderColor: color.NRGBA{A: 255},
	}
}

func (v *View) ApplyProps(vnode *core.VNode) {
	if style, ok := vnode.Props["style"].(*styles.Style); ok {
		style.Apply(v)
	}
}

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

func (v *View) Paint(gtx layout.Context) {
	width := float32(gtx.Dp(unit.Dp(v.Node.LayoutWidth())))
	height := float32(gtx.Dp(unit.Dp(v.Node.LayoutHeight())))
	top := float32(gtx.Dp(unit.Dp(v.Node.LayoutBorder(yoga.EdgeTop))))
	right := float32(gtx.Dp(unit.Dp(v.Node.LayoutBorder(yoga.EdgeRight))))
	bottom := float32(gtx.Dp(unit.Dp(v.Node.LayoutBorder(yoga.EdgeBottom))))
	left := float32(gtx.Dp(unit.Dp(v.Node.LayoutBorder(yoga.EdgeLeft))))

	if v.Background.A > 0 {
		var p clip.Path
		p.Begin(gtx.Ops)
		drawInnerLoop(&p, width, height, top, right, bottom, left, v.CornerRadii, false)
		paint.FillShape(gtx.Ops, v.Background, clip.Outline{Path: p.End()}.Op())
	}

	if v.BorderColor.A > 0 && (top > 0 || right > 0 || bottom > 0 || left > 0) {
		var p clip.Path
		p.Begin(gtx.Ops)
		drawOuterLoop(&p, width, height, v.CornerRadii)
		drawInnerLoop(&p, width, height, top, right, bottom, left, v.CornerRadii, true)
		paint.FillShape(gtx.Ops, v.BorderColor, clip.Outline{Path: p.End()}.Op())
	}
}

// Image is a concrete INode for displaying images.
type Image struct {
	ElementBase
	DefaultSize image.Point
	widget      Widget
}

func NewImage() *Image {
	node := yoga.NewNode()
	img := &Image{
		ElementBase: ElementBase{Node: node},
	}
	return img
}

type Svg struct {
	canvas     *canvas.Canvas
	path       string
	call       op.CallOp
	Record     bool
	dimensions layout.Dimensions
}

const ptPerMm = 72.0 / 25.4

func (s *Svg) DefaultSize() image.Point {
	return image.Pt(int(ptPerMm*s.canvas.W), int(ptPerMm*s.canvas.H))
}
func (s *Svg) Layout(gtx layout.Context) layout.Dimensions {
	if !s.Record {
		ops := gtx.Ops
		cache := new(op.Ops)
		gtx.Ops = cache
		macro := op.Record(gtx.Ops)
		gtx.Constraints.Min = image.Pt(0, 0)
		c := gio.NewContain(gtx, s.canvas.W, s.canvas.H)
		s.canvas.RenderTo(c)
		s.call = macro.Stop()
		gtx.Ops = ops
		s.Record = true
		s.dimensions = c.Dimensions()
	}
	s.call.Add(gtx.Ops)
	return s.dimensions
}

func NewSvg(path string) *Svg {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fc, err := canvas.ParseSVG(f)
	if err != nil {
		panic(err)
	}
	return &Svg{
		path:   path,
		canvas: fc,
	}
}
func (i *Image) ApplyProps(vnode *core.VNode) {
	if style, ok := vnode.Props["style"].(*styles.Style); ok {
		style.Apply(i)
	}
	if source, ok := vnode.Props["source"].(string); ok {
		i.setSource(source)
	}
	w, h := i.Node.StyleGetWidth(), i.Node.StyleGetHeight()
	if math.IsNaN(float64(w)) {
		w = 0
	}

	if math.IsNaN(float64(h)) {
		h = 0
	}
	// 获取图片原始尺寸
	originalWidth := float32(i.DefaultSize.X)
	originalHeight := float32(i.DefaultSize.Y)

	// 如果原始尺寸无效，使用默认值或保持用户设置
	if originalWidth <= 0 || originalHeight <= 0 {
		// 可以设置一个最小尺寸，或者保持用户设置
		if w == 0 {
			i.Node.StyleSetWidth(originalWidth) // 默认宽度
		}
		if h == 0 {
			i.Node.StyleSetHeight(originalHeight) // 默认高度
		}
		return
	}

	// 计算宽高比
	aspectRatio := originalWidth / originalHeight

	// 处理不同情况
	if w == 0 && h == 0 {
		// 两边都没指定，使用原始尺寸
		i.Node.StyleSetWidth(originalWidth)
		i.Node.StyleSetHeight(originalHeight)
	} else if w == 0 && h > 0 {
		// 只指定了有效高度，宽度按比例缩放
		calculatedWidth := h * aspectRatio
		i.Node.StyleSetWidth(calculatedWidth)
		i.Node.StyleSetHeight(h)
	} else if w > 0 && h == 0 {
		// 只指定了有效宽度，高度按比例缩放
		calculatedHeight := w / aspectRatio
		i.Node.StyleSetWidth(w)
		i.Node.StyleSetHeight(calculatedHeight)
	} else if w <= 0 && h <= 0 {
		// 两边都指定了但无效，使用原始尺寸
		i.Node.StyleSetWidth(originalWidth)
		i.Node.StyleSetHeight(originalHeight)
	}
}

func (i *Image) setSource(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	if bytes.HasPrefix(data, []byte("<svg")) {
		svg := NewSvg(path)
		i.widget = svg
		i.DefaultSize = svg.DefaultSize()
		i.GetYogaNode().MarkDirty()
		return
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return
	}

	imgWidget := widget.Image{
		Src: paint.NewImageOp(img),
		Fit: widget.Contain,
	}
	i.widget = imgWidget
	i.DefaultSize = img.Bounds().Size()
	i.GetYogaNode().MarkDirty()
}

func (i *Image) Paint(ctx layout.Context) {
	w := i.GetYogaNode().LayoutWidth()
	y := i.GetYogaNode().LayoutHeight()
	ctx.Constraints.Max = image.Point{X: int(w), Y: int(y)}
	i.widget.Layout(ctx)
}

func (i *Image) SetExtendedStyle(style styles.IExtendedStyle) {
	// Image currently does not have extended styles.
}

// Text is a concrete INode for displaying text.

const q = 4 * (math.Sqrt2 - 1) / 3
const k = 1 - q

func drawOuterLoop(p *clip.Path, w, h float32, r styles.CornerRadius) {
	tl, tr, br, bl := r.TopLeft, r.TopRight, r.BottomRight, r.BottomLeft
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

func drawInnerLoop(p *clip.Path, w, h, top, right, bottom, left float32, r styles.CornerRadius, reverse bool) {
	tlRx, tlRy := max(0, r.TopLeft-left), max(0, r.TopLeft-top)
	trRx, trRy := max(0, r.TopRight-right), max(0, r.TopRight-top)
	brRx, brRy := max(0, r.BottomRight-right), max(0, r.BottomRight-bottom)
	blRx, blRy := max(0, r.BottomLeft-left), max(0, r.BottomLeft-bottom)

	ptTlStart := f32.Pt(left, top+tlRy)
	ptTlEnd := f32.Pt(left+tlRx, top)
	ptTrStart := f32.Pt(w-right-trRx, top)
	ptTrEnd := f32.Pt(w-right, top+trRy)
	ptBrStart := f32.Pt(w-right, h-bottom-brRy)
	ptBrEnd := f32.Pt(w-right-brRx, h-bottom)
	ptBlStart := f32.Pt(left+blRx, h-bottom)
	ptBlEnd := f32.Pt(left, h-bottom-blRy)

	if !reverse {
		p.MoveTo(ptTlEnd)
		p.LineTo(ptTrStart)
		if trRx > 0 || trRy > 0 {
			p.CubeTo(f32.Pt(w-right-trRx+trRx*k, top), f32.Pt(w-right, top+trRy-trRy*k), ptTrEnd)
		} else {
			p.LineTo(ptTrEnd)
		}
		p.LineTo(ptBrStart)
		if brRx > 0 || brRy > 0 {
			p.CubeTo(f32.Pt(w-right, h-bottom-brRy+brRy*k), f32.Pt(w-right-brRx+brRx*k, h-bottom), ptBrEnd)
		} else {
			p.LineTo(ptBrEnd)
		}
		p.LineTo(ptBlStart)
		if blRx > 0 || blRy > 0 {
			p.CubeTo(f32.Pt(left+blRx-blRx*k, h-bottom), f32.Pt(left, h-bottom-blRy+blRy*k), ptBlEnd)
		} else {
			p.LineTo(ptBlEnd)
		}
		p.LineTo(ptTlStart)
		if tlRx > 0 || tlRy > 0 {
			p.CubeTo(f32.Pt(left, top+tlRy-tlRy*k), f32.Pt(left+tlRx-tlRx*k, top), ptTlEnd)
		} else {
			p.LineTo(ptTlEnd)
		}
		p.Close()
	} else {
		p.MoveTo(ptTlEnd)
		if tlRx > 0 || tlRy > 0 {
			p.CubeTo(f32.Pt(left+tlRx-tlRx*k, top), f32.Pt(left, top+tlRy-tlRy*k), ptTlStart)
		} else {
			p.LineTo(ptTlStart)
		}
		p.LineTo(ptBlEnd)
		if blRx > 0 || blRy > 0 {
			p.CubeTo(f32.Pt(left, h-bottom-blRy+blRy*k), f32.Pt(left+blRx-blRx*k, h-bottom), ptBlStart)
		} else {
			p.LineTo(ptBlStart)
		}
		p.LineTo(ptBrEnd)
		if brRx > 0 || brRy > 0 {
			p.CubeTo(f32.Pt(w-right-brRx+brRx*k, h-bottom), f32.Pt(w-right, h-bottom-brRy+brRy*k), ptBrStart)
		} else {
			p.LineTo(ptBrStart)
		}
		p.LineTo(ptTrEnd)
		if trRx > 0 || trRy > 0 {
			p.CubeTo(f32.Pt(w-right, top+trRy-trRy*k), f32.Pt(w-right-trRx+trRx*k, top), ptTrStart)
		} else {
			p.LineTo(ptTrStart)
		}
		p.LineTo(ptTlEnd)
		p.Close()
	}
}
