package ui

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	gpaint "gioui.org/op/paint"
)

// ---- gio 后端：painter 实现 ----
//
// 把 painter 原语翻译成 gio 的 op.Ops。坐标为物理像素，与 gio 的 op 像素空间 1:1。
// 裁剪/图层用 clip 栈、op.Record 宏 + op.Affine + gpaint.PushOpacity 表达。
// 注：伪 3D（rotateX/Y/透视）在 gio 上不支持，EndLayer 只施加 2D 仿射 + 整组透明度。

func nrgba(c Color) color.NRGBA           { return color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A} }
func gioOffset(x, y float32) f32.Affine2D { return f32.Affine2D{}.Offset(f32.Pt(x, y)) }

func absf(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

// gioImage / gioPath 是 gio 版位图与矢量句柄。
type gioImage struct{ src image.Image }

func (i *gioImage) Size() (int, int) { b := i.src.Bounds(); return b.Dx(), b.Dy() }

// gioPath 存 SVG path 的 d 与缩放；实际的 clip.Path 在绘制时（有 op.Ops）惰性构建。
type gioPath struct {
	d     string
	scale float32
}

func (*gioPath) isVecPath() {}

// buildOutline 把 SVG 路径解析进 clip.Path（整体平移 ox,oy），返回其 PathSpec。
func (p *gioPainter) buildPath(gp *gioPath, ox, oy float32) clip.PathSpec {
	var pth clip.Path
	pth.Begin(p.ops)
	parseSVGInto(gp.d, gp.scale, &gioPathSink{p: &pth, ox: ox, oy: oy})
	return pth.End()
}

type gioPainter struct {
	ops    *op.Ops
	w, h   int
	clips  []func()     // PushClip/PopClip 的出栈回调（LIFO）
	layers []op.MacroOp // BeginLayer/EndLayer 的宏栈
}

func newGioPainter(ops *op.Ops, w, h int) *gioPainter { return &gioPainter{ops: ops, w: w, h: h} }

// rrectOp 生成（可含圆角的）矩形裁剪 Op；r<=0 为直角矩形。
//
// 圆角必须钳制到短边的一半：调用方常用 Radius(9999) 表示「全圆角/胶囊」（shadcn 的
// radiusFull），而 gio 的 clip.RRect 不做钳制，半径超过尺寸会让路径退化并破坏整帧绘制。
// clampRadius 把圆角钳到短边的一半。调用方常用 Radius(9999) 表示「全圆角/胶囊」
// （shadcn 的 radiusFull），而 gio 不做钳制，半径超过尺寸会让路径退化并破坏整帧绘制。
func clampRadius(w, h, r float32) float32 {
	if lim := minf(w, h) / 2; r > lim {
		r = lim
	}
	if r < 0 {
		r = 0
	}
	return r
}

// rrectPath 按 float 坐标构造圆角矩形路径。
//
// 不用 clip.RRect：它的 Rect 是 image.Rectangle（整数），把布局的小数坐标取整后，
// 在非整数缩放（如 150%）下会出现 1px 的接缝与抖动。圆弧用与 gio 相同的三次贝塞尔近似。
func (p *gioPainter) rrectPath(x, y, w, h, r float32) clip.PathSpec {
	const q = 4 * (math.Sqrt2 - 1) / 3
	const iq = 1 - q
	r = clampRadius(w, h, r)
	west, north, east, south := x, y, x+w, y+h

	var pth clip.Path
	pth.Begin(p.ops)
	pth.MoveTo(f32.Pt(west+r, north))
	pth.LineTo(f32.Pt(east-r, north))
	pth.CubeTo(f32.Pt(east-r*iq, north), f32.Pt(east, north+r*iq), f32.Pt(east, north+r))
	pth.LineTo(f32.Pt(east, south-r))
	pth.CubeTo(f32.Pt(east, south-r*iq), f32.Pt(east-r*iq, south), f32.Pt(east-r, south))
	pth.LineTo(f32.Pt(west+r, south))
	pth.CubeTo(f32.Pt(west+r*iq, south), f32.Pt(west, south-r*iq), f32.Pt(west, south-r))
	pth.LineTo(f32.Pt(west, north+r))
	pth.CubeTo(f32.Pt(west, north+r*iq), f32.Pt(west+r*iq, north), f32.Pt(west+r, north))
	pth.Close()
	return pth.End()
}

func (p *gioPainter) rrectOp(x, y, w, h, r float32) clip.Op {
	if clampRadius(w, h, r) <= 0 {
		// 直角矩形：clip.Rect 是整数，但直角边缘本就落在像素格上，取整无损且更省。
		return clip.Rect(image.Rect(int(x+0.5), int(y+0.5), int(x+w+0.5), int(y+h+0.5))).Op()
	}
	return clip.Outline{Path: p.rrectPath(x, y, w, h, r)}.Op()
}

func minf(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func (p *gioPainter) fillShape(shape clip.Op, c Color) {
	if c.A == 0 {
		return
	}
	st := shape.Push(p.ops)
	gpaint.ColorOp{Color: nrgba(c)}.Add(p.ops)
	gpaint.PaintOp{}.Add(p.ops)
	st.Pop()
}

func (p *gioPainter) FillRect(x, y, w, h, r float32, c Color) {
	if w <= 0 || h <= 0 {
		return
	}
	p.fillShape(p.rrectOp(x, y, w, h, r), c)
}

func (p *gioPainter) FillGradient(x, y, w, h, r float32, from, to Color, angle float32) {
	if w <= 0 || h <= 0 {
		return
	}
	st := p.rrectOp(x, y, w, h, r).Push(p.ops)
	rad := float64(angle) * math.Pi / 180
	dx, dy := float32(math.Cos(rad)), float32(math.Sin(rad))
	cx, cy := x+w/2, y+h/2
	half := (absf(w*dx) + absf(h*dy)) / 2 // 沿渐变方向的近似半跨度
	gpaint.LinearGradientOp{
		Stop1:  f32.Pt(cx-dx*half, cy-dy*half),
		Color1: nrgba(from),
		Stop2:  f32.Pt(cx+dx*half, cy+dy*half),
		Color2: nrgba(to),
	}.Add(p.ops)
	gpaint.PaintOp{}.Add(p.ops)
	st.Pop()
}

// StrokeRect 描边（可含圆角的）矩形。与填充共用同一条 float 圆角路径，两者始终一致。
func (p *gioPainter) StrokeRect(x, y, w, h, r, width float32, c Color) {
	if c.A == 0 || width <= 0 || w <= 0 || h <= 0 {
		return
	}
	st := clip.Stroke{Path: p.rrectPath(x, y, w, h, r), Width: width}.Op().Push(p.ops)
	gpaint.ColorOp{Color: nrgba(c)}.Add(p.ops)
	gpaint.PaintOp{}.Add(p.ops)
	st.Pop()
}

func (p *gioPainter) Shadow(x, y, w, h, r, offX, offY, blur, spread float32, c Color) {
	if c.A == 0 || w <= 0 || h <= 0 {
		return
	}
	bx, by := x+offX-spread, y+offY-spread
	bw, bh := w+spread*2, h+spread*2
	br := r + spread
	if br < 0 {
		br = 0
	}
	if blur <= 0 {
		p.FillRect(bx, by, bw, bh, br, c)
		return
	}
	const layers = 5
	la := c.Alpha(1.0 / float32(layers)) // 各层等 alpha，向中心累积、向外柔和衰减
	for i := layers; i >= 1; i-- {
		g := blur * float32(i) / float32(layers)
		p.FillRect(bx-g, by-g, bw+g*2, bh+g*2, br+g, la)
	}
	p.FillRect(bx, by, bw, bh, br, la)
}

func (p *gioPainter) Line(x0, y0, x1, y1 float32, c Color) {
	if c.A == 0 {
		return
	}
	var pth clip.Path
	pth.Begin(p.ops)
	pth.MoveTo(f32.Pt(x0, y0))
	pth.LineTo(f32.Pt(x1, y1))
	st := clip.Stroke{Path: pth.End(), Width: 1}.Op().Push(p.ops)
	gpaint.ColorOp{Color: nrgba(c)}.Add(p.ops)
	gpaint.PaintOp{}.Add(p.ops)
	st.Pop()
}

func (p *gioPainter) DrawText(s string, face fontFace, c Color, x, y float32, fauxBold, fauxItalic bool) {
	drawGioText(p.ops, face.(*gioFont), s, c, x, y)
}

func (p *gioPainter) DrawImage(img bitmap, d Rect, opacity float32) {
	gi := img.(*gioImage)
	iw, ih := gi.Size()
	if iw <= 0 || ih <= 0 {
		return
	}
	aff := f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(d.W/float32(iw), d.H/float32(ih))).Offset(f32.Pt(d.X, d.Y))
	tr := op.Affine(aff).Push(p.ops)
	var os gpaint.OpacityStack
	if opacity < 1 {
		os = gpaint.PushOpacity(p.ops, opacity)
	}
	gpaint.NewImageOp(gi.src).Add(p.ops)
	gpaint.PaintOp{}.Add(p.ops)
	if opacity < 1 {
		os.Pop()
	}
	tr.Pop()
}

func (p *gioPainter) FillPath(path vecPath, x, y float32, c Color) {
	if c.A == 0 {
		return
	}
	st := clip.Outline{Path: p.buildPath(path.(*gioPath), x, y)}.Op().Push(p.ops)
	gpaint.ColorOp{Color: nrgba(c)}.Add(p.ops)
	gpaint.PaintOp{}.Add(p.ops)
	st.Pop()
}

func (p *gioPainter) StrokePath(path vecPath, x, y, width float32, c Color) {
	if c.A == 0 || width <= 0 {
		return
	}
	st := clip.Stroke{Path: p.buildPath(path.(*gioPath), x, y), Width: width}.Op().Push(p.ops)
	gpaint.ColorOp{Color: nrgba(c)}.Add(p.ops)
	gpaint.PaintOp{}.Add(p.ops)
	st.Pop()
}

func (p *gioPainter) PushClip(r Rect, radius float32) {
	st := p.rrectOp(r.X, r.Y, r.W, r.H, radius).Push(p.ops)
	p.clips = append(p.clips, st.Pop)
}

func (p *gioPainter) PopClip() {
	n := len(p.clips) - 1
	pop := p.clips[n]
	p.clips = p.clips[:n]
	pop()
}

func (p *gioPainter) BeginLayer() {
	p.layers = append(p.layers, op.Record(p.ops))
}

func (p *gioPainter) EndLayer(t layerTransform) {
	n := len(p.layers) - 1
	call := p.layers[n].Stop()
	p.layers = p.layers[:n]

	if t.is3D() { // 伪 3D：仿射表达不了透视，改走网格投影重放
		p.drawProjected(call, t)
		return
	}

	// 纯 2D：绕中心缩放/旋转 + 平移。
	aff := f32.Affine2D{}
	ctr := f32.Pt(t.cx, t.cy)
	if t.scale != 1 && t.scale != 0 {
		aff = aff.Scale(ctr, f32.Pt(t.scale, t.scale))
	}
	if t.rotate != 0 {
		aff = aff.Rotate(ctr, float32(float64(t.rotate)*math.Pi/180))
	}
	if t.tx != 0 || t.ty != 0 {
		aff = aff.Offset(f32.Pt(t.tx, t.ty))
	}
	tr := op.Affine(aff).Push(p.ops)
	var os gpaint.OpacityStack
	if t.opacity < 1 {
		os = gpaint.PushOpacity(p.ops, t.opacity)
	}
	call.Add(p.ops)
	if t.opacity < 1 {
		os.Pop()
	}
	tr.Pop()
}
