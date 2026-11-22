package ui

// Painter 是绘制后端的抽象：paint 遍历只依赖这些原语。默认实现为 Gio，
// 测试用录制后端（recordPainter，无需 GPU）。坐标均为物理像素，与 renderNode.bounds 一致。
type Painter interface {
	FillRect(x, y, w, h, r float32, c Color)
	FillGradient(x, y, w, h, r float32, from, to Color, angle float32) // 线性渐变填充
	StrokeRect(x, y, w, h, r, width float32, c Color)
	Shadow(x, y, w, h, r, offX, offY, blur, spread float32, c Color)
	Line(x0, y0, x1, y1 float32, c Color)
	DrawText(s string, face FontFace, c Color, x, y float32, fauxBold, fauxItalic bool)
	DrawImage(img Bitmap, dst Rect, opacity float32)
	FillPath(path VecPath, x, y float32, c Color)          // 填充路径（实心图标）
	StrokePath(path VecPath, x, y, width float32, c Color) // 描边路径（线性图标，圆角端点）
	PushClip(r Rect, radius float32)                       // radius>0 裁剪到圆角矩形
	PopClip()
	BeginLayer()
	EndLayer(t LayerTransform)
	// 投影层（伪 3D）：逐元素录制、逐元素仿射回放。BeginProjected/EndProjected 之间
	// 用 BeginProjectedElement/EndProjectedElement 成对包裹每个元素的绘制（第一个元素
	// 应是层自身，rect 传层的 bounds）；回放时每个元素贴自己四角投影的前三角 ——
	// 斜切误差从「整层尺寸」降到「元素尺寸」。命中端用 elementContentAffine 的逆，
	// 与绘制端共用同一份 projAffine，画哪点哪。
	BeginProjected(t LayerTransform)
	BeginProjectedElement(r Rect)
	EndProjectedElement()
	// ProjectedImage 在投影层里登记一个「精确单应贴图」元素：整块图按 quad（元素四角的
	// 真实投影，顺序 左上/右上/左下/右下）预变形后正着贴，绕开仿射斜切 —— 用于「元素
	// 自身就是一整块图」的情形。clip 非 nil 时按层裁剪矩形+圆角裁剪（裁剪路径随元素
	// 仿射变形，与常规元素路径一致）。cacheKey 标识图像内容。
	ProjectedImage(img Bitmap, quad [4]Point, opacity float32, cacheKey string, clip *Rect, clipR float32)
	EndProjected()
}

type painter = Painter

// LayerTransform 是一个图层合回时施加的变换和透明度。后端通过 Snapshot 读取。
type LayerTransform = layerTransform

// Point 是投影四边形中的二维点。
type Point = pt

// Layer 是 LayerTransform 的只读快照，供外部渲染后端使用。
type Layer struct {
	CX, CY              float32
	W, H                float32
	Scale, ScaleY       float32
	Rotate              float32
	TX, TY              float32
	Opacity             float32
	RotateX, RotateY    float32
	TransZ, Perspective float32
	CameraX, CameraY    float32
	HasCamera           bool
}

// Snapshot 返回图层的全部绘制变换。
func (t layerTransform) Snapshot() Layer {
	return Layer{
		CX: t.cx, CY: t.cy, W: t.w, H: t.h,
		Scale: t.scale, ScaleY: t.scaleY, Rotate: t.rotate, TX: t.tx, TY: t.ty, Opacity: t.opacity,
		RotateX: t.rotateX, RotateY: t.rotateY, TransZ: t.transZ, Perspective: t.perspective,
		CameraX: t.camX, CameraY: t.camY, HasCamera: t.hasCam,
	}
}

// layerTransform 是一个图层合回时施加的变换 + 整组透明度（围绕中心 cx,cy）。
type layerTransform struct {
	cx, cy        float32 // 元素中心（变换与透视的锚点）
	w, h          float32 // 元素尺寸（伪 3D 按四角投影求仿射时要用）
	scale, scaleY float32 // 2D 横向/纵向缩放
	rotate        float32 // 绕 Z 轴旋转
	tx, ty        float32
	opacity       float32

	// 伪 3D：绕 X/Y 轴旋转 + Z 位移 + 透视距离（0=无透视）。任一非零即走投影路径。
	rotateX, rotateY float32
	transZ           float32
	perspective      float32

	// 共享相机（Scene3D）：投影原点改用 camX,camY，元素按自身在场景中的位置投影。
	// 于是同一场景里所有元素的灭点一致 —— 摆一桌卡才像同一张桌子。
	// hasCam 为假时原点就是元素自身中心（等价于 CSS 的 transform: perspective(...)）。
	camX, camY float32
	hasCam     bool
}

// origin 是投影原点：有相机用相机原点，否则用元素中心。
func (t layerTransform) origin() (float32, float32) {
	if t.hasCam {
		return t.camX, t.camY
	}
	return t.cx, t.cy
}

func (t layerTransform) is3D() bool {
	return t.rotateX != 0 || t.rotateY != 0 || t.transZ != 0
}

func (t layerTransform) xyScale() (float32, float32) {
	y := t.scaleY
	if y == 0 {
		y = t.scale
	}
	return t.scale, y
}

// ---- 录制后端（测试用，无需 GPU）----

// PaintOp 是一条绘制指令记录。Harness.Paint 返回其有序列表，用于无头断言"画了什么"
// （黄金测试）：颜色、位置、文本、裁剪/图层边界都可检查，无需真实渲染。
type PaintOp struct {
	Kind           string // rect / stroke / shadow / line / text / image / clip / unclip / layer / unlayer / projlayer / unprojlayer / projelement / unprojelement
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
func (p *recordPainter) BeginProjected(t layerTransform) {
	p.ops = append(p.ops, PaintOp{Kind: "projlayer",
		Rect: Rect{X: t.cx - t.w/2, Y: t.cy - t.h/2, W: t.w, H: t.h}, Opacity: t.opacity})
}
func (p *recordPainter) BeginProjectedElement(r Rect) {
	p.ops = append(p.ops, PaintOp{Kind: "projelement", Rect: r})
}
func (p *recordPainter) EndProjectedElement() {
	p.ops = append(p.ops, PaintOp{Kind: "unprojelement"})
}
func (p *recordPainter) ProjectedImage(img bitmap, quad [4]pt, opacity float32, cacheKey string, clip *Rect, clipR float32) {
	r := quadBBox(quad[0], quad[1], quad[2], quad[3])
	kind := "projimage"
	if clip != nil {
		kind = "projimageclip"
	}
	p.ops = append(p.ops, PaintOp{Kind: kind, Rect: r, Radius: clipR, Opacity: opacity})
}
func (p *recordPainter) EndProjected() {
	p.ops = append(p.ops, PaintOp{Kind: "unprojlayer"})
}
