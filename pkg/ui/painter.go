package ui

// painter 是绘制后端的抽象：paint 遍历只依赖这些原语。当前实现为 gio（gioPainter），
// 测试用录制后端（recordPainter，无需 GPU）。坐标均为物理像素，与 renderNode.bounds 一致。
type painter interface {
	FillRect(x, y, w, h, r float32, c Color)
	FillGradient(x, y, w, h, r float32, from, to Color, angle float32) // 线性渐变填充
	StrokeRect(x, y, w, h, r, width float32, c Color)
	Shadow(x, y, w, h, r, offX, offY, blur, spread float32, c Color)
	Line(x0, y0, x1, y1 float32, c Color)
	DrawText(s string, face fontFace, c Color, x, y float32, fauxBold, fauxItalic bool)
	DrawImage(img bitmap, dst Rect, opacity float32)
	FillPath(path vecPath, x, y float32, c Color)          // 填充路径（实心图标）
	StrokePath(path vecPath, x, y, width float32, c Color) // 描边路径（线性图标，圆角端点）
	PushClip(r Rect, radius float32)                       // radius>0 裁剪到圆角矩形
	PopClip()
	BeginLayer()
	EndLayer(t layerTransform)
}

// layerTransform 是一个图层合回时施加的变换 + 整组透明度（围绕中心 cx,cy）。
type layerTransform struct {
	cx, cy        float32 // 元素中心（变换与透视的锚点）
	w, h          float32 // 元素尺寸（伪 3D 投影按它划分网格）
	scale, rotate float32 // 2D 缩放 + 绕 Z 轴旋转
	tx, ty        float32
	opacity       float32

	// 伪 3D：绕 X/Y 轴旋转 + Z 位移 + 透视距离（0=无透视）。任一非零即走投影路径。
	rotateX, rotateY float32
	transZ           float32
	perspective      float32
}

func (t layerTransform) is3D() bool {
	return t.rotateX != 0 || t.rotateY != 0 || t.transZ != 0
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
func (p *recordPainter) FillGradient(x, y, w, h, r float32, from, to Color, angle float32) {
	p.ops = append(p.ops, PaintOp{Kind: "gradient", Rect: Rect{x, y, w, h}, Radius: r, Color: from})
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
func (p *recordPainter) DrawText(s string, face fontFace, c Color, x, y float32, fb, fi bool) {
	p.ops = append(p.ops, PaintOp{Kind: "text", Text: s, Color: c, X0: x, Y0: y})
}
func (p *recordPainter) DrawImage(img bitmap, d Rect, opacity float32) {
	p.ops = append(p.ops, PaintOp{Kind: "image", Rect: d, Opacity: opacity})
}
func (p *recordPainter) FillPath(path vecPath, x, y float32, c Color) {
	p.ops = append(p.ops, PaintOp{Kind: "path", X0: x, Y0: y, Color: c})
}
func (p *recordPainter) StrokePath(path vecPath, x, y, width float32, c Color) {
	p.ops = append(p.ops, PaintOp{Kind: "strokepath", X0: x, Y0: y, Width: width, Color: c})
}
func (p *recordPainter) PushClip(r Rect, radius float32) {
	p.ops = append(p.ops, PaintOp{Kind: "clip", Rect: r, Radius: radius})
}
func (p *recordPainter) PopClip()    { p.ops = append(p.ops, PaintOp{Kind: "unclip"}) }
func (p *recordPainter) BeginLayer() { p.ops = append(p.ops, PaintOp{Kind: "layer"}) }
func (p *recordPainter) EndLayer(t layerTransform) {
	p.ops = append(p.ops, PaintOp{Kind: "unlayer", Opacity: t.opacity})
}
