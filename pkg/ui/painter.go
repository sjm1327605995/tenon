package ui

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// painter 是绘制后端的抽象：paint 遍历只依赖这些原语，从而可替换后端 ——
// 当前 ebiten（ebitenPainter），测试用录制后端（recordPainter），将来可加 Skia 等。
// 坐标均为物理像素，与 renderNode.bounds 一致。
type painter interface {
	FillRect(x, y, w, h, r float32, c Color)
	StrokeRect(x, y, w, h, r, width float32, c Color)
	Shadow(x, y, w, h, r, offX, offY, blur, spread float32, c Color)
	Line(x0, y0, x1, y1 float32, c Color)
	DrawText(s string, face *text.GoTextFace, c Color, x, y float32, fauxBold, fauxItalic bool)
	DrawImage(img *ebiten.Image, dst Rect, opacity float32)
	PushClip(r Rect)
	PopClip()
	BeginLayer()
	EndLayer(t layerTransform)
}

// layerTransform 是一个离屏图层合回时施加的变换 + 整组透明度（围绕中心）。
type layerTransform struct {
	cx, cy        float32
	scale, rotate float32
	tx, ty        float32
	opacity       float32
}

// ---- ebiten 后端 ----

// ebitenPainter 把原语画到 *ebiten.Image。stack 顶为当前绘制目标，
// 裁剪（SubImage）与离屏图层都以压栈形式表达，保证嵌套平衡。
type ebitenPainter struct {
	w, h  int
	stack []*ebiten.Image
}

func newEbitenPainter(dst *ebiten.Image, w, h int) *ebitenPainter {
	return &ebitenPainter{w: w, h: h, stack: []*ebiten.Image{dst}}
}

func (p *ebitenPainter) top() *ebiten.Image { return p.stack[len(p.stack)-1] }

func (p *ebitenPainter) FillRect(x, y, w, h, r float32, c Color) {
	fillRoundRect(p.top(), x, y, w, h, r, c)
}
func (p *ebitenPainter) StrokeRect(x, y, w, h, r, width float32, c Color) {
	strokeRoundRect(p.top(), x, y, w, h, r, width, c)
}
func (p *ebitenPainter) Shadow(x, y, w, h, r, offX, offY, blur, spread float32, c Color) {
	fillShadow(p.top(), x, y, w, h, r, offX, offY, blur, spread, c)
}
func (p *ebitenPainter) Line(x0, y0, x1, y1 float32, c Color) {
	ebitenDrawLine(p.top(), x0, y0, x1, y1, c)
}
func (p *ebitenPainter) DrawText(s string, face *text.GoTextFace, c Color, x, y float32, fb, fi bool) {
	drawText(p.top(), s, face, c, x, y, fb, fi)
}

func (p *ebitenPainter) DrawImage(img *ebiten.Image, d Rect, opacity float32) {
	ib := img.Bounds()
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(float64(d.W)/float64(ib.Dx()), float64(d.H)/float64(ib.Dy()))
	op.GeoM.Translate(float64(d.X), float64(d.Y))
	op.ColorScale.ScaleAlpha(opacity)
	p.top().DrawImage(img, op)
}

func (p *ebitenPainter) PushClip(r Rect) {
	dst := p.top()
	ir := image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H))
	if ir.Dx() > 0 && ir.Dy() > 0 {
		p.stack = append(p.stack, dst.SubImage(ir).(*ebiten.Image))
	} else {
		p.stack = append(p.stack, dst) // 空裁剪：压入同一目标，保持栈平衡
	}
}
func (p *ebitenPainter) PopClip() { p.stack = p.stack[:len(p.stack)-1] }

func (p *ebitenPainter) BeginLayer() {
	p.stack = append(p.stack, acquireLayer(p.w, p.h))
}
func (p *ebitenPainter) EndLayer(t layerTransform) {
	layer := p.top()
	p.stack = p.stack[:len(p.stack)-1]
	dst := p.top()
	cx, cy := float64(t.cx), float64(t.cy)
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-cx, -cy)
	op.GeoM.Scale(float64(t.scale), float64(t.scale))
	op.GeoM.Rotate(float64(t.rotate) * math.Pi / 180)
	op.GeoM.Translate(cx+float64(t.tx), cy+float64(t.ty))
	op.ColorScale.ScaleAlpha(t.opacity)
	dst.DrawImage(layer, op)
	releaseLayer(layer)
}

// ---- 录制后端（测试用，无需 GPU）----

// PaintOp 是一条绘制指令记录。Harness.Paint 返回其有序列表，用于无头断言"画了什么"
// （黄金测试）：颜色、位置、文本、裁剪/图层边界都可检查，无需真实渲染。
type PaintOp struct {
	Kind           string // rect / stroke / shadow / line / text / image / clip / unclip / layer / unlayer
	Rect           Rect
	Radius, Width  float32
	Color          Color
	Text           string
	X0, Y0, X1, Y1 float32
	Opacity        float32
}

// recordPainter 把每个原语调用记成一条 PaintOp。
type recordPainter struct{ ops []PaintOp }

func (p *recordPainter) FillRect(x, y, w, h, r float32, c Color) {
	p.ops = append(p.ops, PaintOp{Kind: "rect", Rect: Rect{x, y, w, h}, Radius: r, Color: c})
}
func (p *recordPainter) StrokeRect(x, y, w, h, r, width float32, c Color) {
	p.ops = append(p.ops, PaintOp{Kind: "stroke", Rect: Rect{x, y, w, h}, Radius: r, Width: width, Color: c})
}
func (p *recordPainter) Shadow(x, y, w, h, r, offX, offY, blur, spread float32, c Color) {
	p.ops = append(p.ops, PaintOp{Kind: "shadow", Rect: Rect{x, y, w, h}, Radius: r, Color: c})
}
func (p *recordPainter) Line(x0, y0, x1, y1 float32, c Color) {
	p.ops = append(p.ops, PaintOp{Kind: "line", X0: x0, Y0: y0, X1: x1, Y1: y1, Color: c})
}
func (p *recordPainter) DrawText(s string, face *text.GoTextFace, c Color, x, y float32, fb, fi bool) {
	p.ops = append(p.ops, PaintOp{Kind: "text", Text: s, Color: c, X0: x, Y0: y})
}
func (p *recordPainter) DrawImage(img *ebiten.Image, d Rect, opacity float32) {
	p.ops = append(p.ops, PaintOp{Kind: "image", Rect: d, Opacity: opacity})
}
func (p *recordPainter) PushClip(r Rect) { p.ops = append(p.ops, PaintOp{Kind: "clip", Rect: r}) }
func (p *recordPainter) PopClip()        { p.ops = append(p.ops, PaintOp{Kind: "unclip"}) }
func (p *recordPainter) BeginLayer()     { p.ops = append(p.ops, PaintOp{Kind: "layer"}) }
func (p *recordPainter) EndLayer(t layerTransform) {
	p.ops = append(p.ops, PaintOp{Kind: "unlayer", Opacity: t.opacity})
}
