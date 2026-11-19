package ui

import (
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	gpaint "gioui.org/op/paint"
)

// ---- gio 后端：伪 3D 投影 ----
//
// 把已录制的图层（op.Record 宏）按透视贴回屏幕，实现 CSS 的
// transform: perspective(p) rotateX(rx) rotateY(ry) translateZ(tz)。内容仍是平面，
// 只是被贴到一个倾斜的四边形上 —— 所谓「伪 3D」。
//
// 做法是把「形状」和「内容」分开（细节见 replayProjected）：
//
//   - 形状用裁剪：裁到四角真实投影出的四边形，于是远端真的更窄 —— 透视收缩靠它。
//   - 内容用仿射：每个元素贴自己四角投影的前三角，重放该元素的宏。文字和图片因此
//     会一起变形（逐图元投影做不到这点：字形轮廓是 shaper 给的不透明 PathSpec，
//     无法重投影）。
//
// 内容按「元素」拆分（paintProjectedLayer）：一个投影层 = 层自身 + 每个直接子元素
// 各录一个宏，回放时各用自己的三角仿射。仿射第四角的误差随元素尺寸增长，拆到元素
// 粒度后，层里一行 150×20 的文字斜切只有亚像素级。
//
// 「元素自身就是一整块图」的纯图片（铺满层的大图、倾斜的照片）是剩余误差的集中地，
// 由 projExactImage 自动识别、走精确单应贴图（ProjectedImage，与 PlaneImage 同一份
// unitToQuad 数学）—— 无需调用方换用特殊组件。仍带可见斜切的只剩拒绝自动拦截的
// 情形：object-fit≠Fill、带阴影/变换的复合元素。
//
// Ebiten 那套「离屏纹理 + 带 UV 的贴图网格」在 gio 上行不通：公开 API 既不提供渲染到
// 纹理，也没有逐顶点贴图；而 op.Affine 是仿射，表达不了透视。
//
// # 已知边界：元素级近似，不是精确投影
//
// 仿射有 6 个自由度、四角是 8 个约束，所以单个元素的内容只能匹配三角，第四角必然偏
//（误差随元素尺寸与倾角增长，见 projErr）。逐元素拆分把它从「整层尺寸」压到「元素
// 尺寸」：80×110 的卡 45° 下曾斜 21px，卡内一行文字现在只剩零点几 px。它连续平滑，
// 读起来像「图案所在平面略有偏差」。
//
// 可见轮廓仍是「裁剪梯形 ∩ 仿射平行四边形」的交集，不是纯梯形 —— 内容并不总能填满
// 裁剪（实测覆盖率 80%~99%），但两个凸四边形的交集仍是凸四边形，看上去就是一张角度
// 略有不同的卡，不值得为此放大内容去硬填满。
//
// 曾用 n×n 网格逐格近似来消除这个斜切，结果远比它糟（见 git 历史）：每格各用一个仿射，
// 内容跨格对不上会撕裂，且格内内容填不满自己的裁剪格、缺口露出底色 —— 密度越高越像
// 铺了层铁丝网。逐元素拆分没有这个病：元素的边界正是内容自身的边界，不存在跨格内容。
// 教训是：形状错了一眼可见，内容斜切几乎看不见，所以该把裁剪的精度用在轮廓上，
// 而不是去细分内容。

// gioPt / gioAffine 是中立投影数学（project3d.go）到 gio 类型的适配 —— 边界只在这里。
func gioPt(p pt) f32.Point { return f32.Pt(p.X, p.Y) }

func gioAffine(a affine2D) f32.Affine2D {
	return f32.NewAffine2D(a.sx, a.hx, a.ox, a.hy, a.sy, a.oy)
}

// drawProjected 把录制好的图层按透视贴回：形状和内容分开处理。
//
//   - 形状：裁到四角真实投影出的四边形 —— 透视收缩（远端更窄）由它带来。
//   - 内容：用三角仿射重放 —— 内部略有斜切，但连续、无撕裂。
//
// 内容并不总能填满裁剪四边形（实测覆盖 80%~99%），所以可见轮廓其实是两者的交集。
// 那仍是个凸四边形、看着就是张卡，不值得为此放大内容去硬填满 —— 见文件头「已知边界」。
//
// 本函数现在只服务于直接调 EndLayer(3D) 的路径（场景底背景等单元素层）；
// 带子元素的投影层走 paintProjectedLayer 的逐元素回放（见文件头）。
func (p *gioPainter) drawProjected(call op.CallOp, t layerTransform) {
	if t.opacity <= 0 || t.w <= 0 || t.h <= 0 {
		return
	}
	// 单元素层：元素矩形即层矩形，退化为整层单仿射。
	r := Rect{X: t.cx - t.w/2, Y: t.cy - t.h/2, W: t.w, H: t.h}
	p.replayProjected(projLayerRec{t: t, elements: []projElement{{call: call, r: r}}})
}

// ---- 逐元素投影（paintProjectedLayer 的录制与回放）----
//
// 一个投影层不再录成一个宏、回放时套整层单仿射，而是拆成「层自身 + 每个直接子元素」
// 各自成宏，回放时每个元素贴自己四角投影的前三角（elementContentAffine）。
// 仿射的第四角误差随元素尺寸增长，拆到元素粒度后，层里 150×20 的一行文字斜切从
// 「整层尺寸决定的几十 px」降到亚像素级；只剩元素自身那么大的一整块图（铺满层的大图）
// 仍有可见斜切 —— 那正是 PlaneImage 的场景。

func (p *gioPainter) BeginProjected(t layerTransform) {
	p.projs = append(p.projs, projLayerRec{t: t})
}

func (p *gioPainter) BeginProjectedElement(r Rect) {
	p.projEl = append(p.projEl, projElementRec{r: r, macro: op.Record(p.ops)})
}

func (p *gioPainter) EndProjectedElement() {
	n := len(p.projEl) - 1
	e := p.projEl[n]
	p.projEl = p.projEl[:n]
	call := e.macro.Stop()
	top := &p.projs[len(p.projs)-1]
	top.elements = append(top.elements, projElement{call: call, r: e.r})
}

func (p *gioPainter) ProjectedImage(img bitmap, quad [4]pt, opacity float32, cacheKey string, clip *Rect, clipR float32) {
	top := &p.projs[len(p.projs)-1]
	e := projElement{img: img, quad: quad, key: cacheKey, opacity: opacity, clipR: clipR}
	if clip != nil {
		c := *clip
		e.clip = &c
	}
	top.elements = append(top.elements, e)
}

func (p *gioPainter) EndProjected() {
	n := len(p.projs) - 1
	rec := p.projs[n]
	p.projs = p.projs[:n]
	p.replayProjected(rec)
}

// replayProjected 回放一个投影层：整层裁到真实投影四边形（形状精确），
// 每个元素用自己的三角仿射重放（斜切误差 = 元素尺寸级）。
func (p *gioPainter) replayProjected(rec projLayerRec) {
	t := rec.t
	if t.opacity <= 0 || len(rec.elements) == 0 {
		return
	}
	// 整层裁剪：四角真实投影出的四边形 —— 透视收缩（远端更窄）由它带来。
	d0, d1, d2, d3 := projCorners(t)
	var quad clip.Path
	quad.Begin(p.ops)
	quad.MoveTo(gioPt(d0))
	quad.LineTo(gioPt(d1))
	quad.LineTo(gioPt(d3))
	quad.LineTo(gioPt(d2))
	quad.Close()
	cl := clip.Outline{Path: quad.End()}.Op().Push(p.ops)

	var os gpaint.OpacityStack
	if t.opacity < 1 {
		os = gpaint.PushOpacity(p.ops, t.opacity)
	}
	for _, e := range rec.elements {
		// 精确单应贴图元素：整块图预变形到四角的真实投影，绕开仿射斜切。
		if e.img != nil {
			gi, ok := e.img.(*gioImage)
			if !ok {
				continue
			}
			bmp, box, ok := quadBitmap(e.key, gi.src, e.quad[0], e.quad[1], e.quad[2], e.quad[3])
			if !ok {
				continue
			}
			if e.clip != nil {
				// 层裁剪：圆角矩形路径在元素仿射下构造（变形方式与常规元素路径一致）。
				tr := op.Affine(gioAffine(elementContentAffine(t, e.r))).Push(p.ops)
				st := p.rrectOp(e.clip.X, e.clip.Y, e.clip.W, e.clip.H, e.clipR).Push(p.ops)
				p.DrawImage(bmp, box, e.opacity)
				st.Pop()
				tr.Pop()
			} else {
				p.DrawImage(bmp, box, e.opacity)
			}
			continue
		}
		if e.r.W <= 0 || e.r.H <= 0 {
			continue // 零尺寸元素：其仿射退化为恒等，内容也必然为空
		}
		tr := op.Affine(gioAffine(elementContentAffine(t, e.r))).Push(p.ops)
		e.call.Add(p.ops)
		tr.Pop()
	}
	if t.opacity < 1 {
		os.Pop()
	}
	cl.Pop()
}
