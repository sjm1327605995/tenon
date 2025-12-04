package dom

import (
	"bytes"
	"image"
	"image/color"
	"os"

	colEmoji "eliasnaur.com/font/noto/emoji/color"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	giotext "gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/inkeliz/giosvg"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/yoga"
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
	widget Widget
}

func NewImage() *Image {
	node := yoga.NewNode()
	img := &Image{
		ElementBase: ElementBase{Node: node},
	}
	return img
}

func (i *Image) ApplyProps(vnode *core.VNode) {
	if style, ok := vnode.Props["style"].(*styles.Style); ok {
		style.Apply(i)
	}
	if source, ok := vnode.Props["source"].(string); ok {
		i.setSource(source)
	}
}

func (i *Image) setSource(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	if bytes.HasPrefix(data, []byte("<svg")) {
		vector, err := giosvg.NewVector(data)
		if err != nil {
			panic(err)
		}
		iconRuntime := giosvg.NewIcon(vector)
		i.widget = iconRuntime
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
type Text struct {
	ElementBase
	LabelStyle material.LabelStyle
}

func NewText() *Text {
	th := material.NewTheme()
	faces, _ := opentype.ParseCollection(colEmoji.TTF)
	collection := gofont.Collection()
	th.Shaper = giotext.NewShaper(giotext.WithCollection(append(collection, faces...)))

	return &Text{
		ElementBase: ElementBase{Node: yoga.NewNode()},
		LabelStyle:  material.Label(th, unit.Sp(16), ""),
	}
}

func (t *Text) ApplyProps(vnode *core.VNode) {
	if style, ok := vnode.Props["style"].(*styles.Style); ok {
		style.Apply(t)
	}
	if content, ok := vnode.Props["content"].(string); ok {
		t.LabelStyle.Text = content
		t.GetYogaNode().MarkDirty()
	}
}

func (t *Text) Paint(ctx layout.Context) {
	if t.LabelStyle.Text == "" {
		return
	}
	t.LabelStyle.Layout(ctx)
}

func (t *Text) SetExtendedStyle(style styles.IExtendedStyle) {
	// Text currently does not have extended styles.
}

// --- Drawing helpers ---
const k = 0.55228475 // Bezier curve coefficient

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
