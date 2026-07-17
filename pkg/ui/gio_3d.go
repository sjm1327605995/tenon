package ui

import (
	"math"

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
// 做法是把「形状」和「内容」分开（细节见 drawProjected）：
//
//   - 形状用裁剪：裁到四角真实投影出的四边形，于是远端真的更窄 —— 透视收缩靠它。
//   - 内容用仿射：由其中三角唯一确定，重放整棵子树。文字和图片因此会一起变形
//     （逐图元投影做不到这点：字形轮廓是 shaper 给的不透明 PathSpec，无法重投影）。
//
// Ebiten 那套「离屏纹理 + 带 UV 的贴图网格」在 gio 上行不通：公开 API 既不提供渲染到
// 纹理，也没有逐顶点贴图；而 op.Affine 是仿射，表达不了透视。
//
// # 已知边界：这是近似，不是精确投影
//
// 仿射有 6 个自由度、四角是 8 个约束，所以内容只能匹配三角，第四角必然偏（见 projErr）。
// 由此有两处近似，都实测过、都不显眼：
//
//   - 内容是斜切而非透视。偏差随元素尺寸与倾角增长（80×110 的卡 45° 下约 21px；
//     150×180 在 60° 强透视下约 68px）。它连续平滑，读起来像「图案所在平面略有偏差」。
//   - 可见轮廓是「裁剪梯形 ∩ 仿射平行四边形」，不是纯梯形 —— 内容并不总能填满裁剪，
//     实测覆盖率 80%~99%（见 TestProjectedContentFillsClipQuad）。但两个凸四边形的交集
//     仍是凸四边形，看上去就是一张角度略有不同的卡，覆盖率 0.80 时肉眼也看不出缺角。
//
// 要消掉这两点就得把仿射放大到包住裁剪四边形，代价是卡面内容被放大、边缘被裁掉 ——
// 那个代价比这两处偏差显眼得多，所以不做。
//
// 曾用 n×n 网格逐格近似来消除这个斜切，结果远比它糟（见 git 历史）：每格各用一个仿射，
// 内容跨格对不上会撕裂，且格内内容填不满自己的裁剪格、缺口露出底色 —— 密度越高越像
// 铺了层铁丝网。裁剪外扩、仿射外扩、两者同扩都试过，都盖不住。教训是：形状错了一眼可见，
// 内容斜切几乎看不见，所以该把裁剪的精度用在轮廓上，而不是去细分内容。

// project3D 把元素平面上的一点（相对元素中心的偏移）投影到屏幕坐标。
// 顺序与 CSS 的 rotateX(rx) rotateY(ry) 书写序一致：先绕 Y 再绕 X，然后透视除法，
// 最后叠加 2D 的缩放/绕 Z 旋转/平移。
func project3D(t layerTransform, ox, oy float32) f32.Point {
	sinX, cosX := math.Sincos(float64(t.rotateX) * math.Pi / 180)
	sinY, cosY := math.Sincos(float64(t.rotateY) * math.Pi / 180)
	sinZ, cosZ := math.Sincos(float64(t.rotate) * math.Pi / 180)

	x, y, z := float64(ox), float64(oy), float64(t.transZ)
	// 绕 Y 轴
	x1 := x*cosY + z*sinY
	z1 := -x*sinY + z*cosY
	// 绕 X 轴
	y2 := y*cosX - z1*sinX
	z2 := y*sinX + z1*cosX
	x2 := x1

	if p := float64(t.perspective); p > 0 {
		// 透视除法：Z 朝观察者(正)则放大，远离则缩小。
		denom := 1 - z2/p
		if denom < 0.05 { // 越过视平面时钳制，避免坐标发散/翻转
			denom = 0.05
		}
		x2 /= denom
		y2 /= denom
	}
	if s := float64(t.scale); s != 0 {
		x2 *= s
		y2 *= s
	}
	orgX, orgY := t.origin() // 投影原点：无相机=元素中心，有相机=场景中心
	return f32.Pt(
		orgX+float32(x2*cosZ-y2*sinZ)+t.tx,
		orgY+float32(x2*sinZ+y2*cosZ)+t.ty,
	)
}

// projAffine 求把源矩形 [sx0,sx1]×[sy0,sy1] 映射到 (d0,d1,d2) 的仿射矩阵，
// 其中 d0/d1/d2 分别是左上/右上/左下角的投影点。三点唯一确定一个仿射。
func projAffine(sx0, sy0, sx1, sy1 float32, d0, d1, d2 f32.Point) f32.Affine2D {
	w, h := sx1-sx0, sy1-sy0
	if w == 0 || h == 0 {
		return f32.AffineId()
	}
	sx := (d1.X - d0.X) / w // ∂x/∂sx
	hx := (d2.X - d0.X) / h // ∂x/∂sy
	hy := (d1.Y - d0.Y) / w // ∂y/∂sx
	sy := (d2.Y - d0.Y) / h // ∂y/∂sy
	return f32.NewAffine2D(
		sx, hx, d0.X-sx0*sx-sy0*hx,
		hy, sy, d0.Y-sx0*hy-sy0*sy,
	)
}

// projErr 量化内容的斜切程度：仿射把第四角映到哪、与真透视差多少。轮廓由裁剪保证，
// 所以这不是形状误差，只是卡面内部图案的偏差。元素越大、倾角越陡越大。仅用于测试与诊断。
func projErr(t layerTransform) float32 {
	ox, oy := t.origin()
	x0, y0 := t.cx-t.w/2, t.cy-t.h/2
	x1, y1 := x0+t.w, y0+t.h
	d0 := project3D(t, x0-ox, y0-oy)
	d1 := project3D(t, x1-ox, y0-oy)
	d2 := project3D(t, x0-ox, y1-oy)
	d3 := project3D(t, x1-ox, y1-oy)
	got := projAffine(x0, y0, x1, y1, d0, d1, d2).Transform(f32.Pt(x1, y1))
	return absf(got.X-d3.X) + absf(got.Y-d3.Y)
}

// drawProjected 把录制好的图层按透视贴回：形状和内容分开处理。
//
//   - 形状：裁到四角真实投影出的四边形 —— 透视收缩（远端更窄）由它带来。
//   - 内容：用三角仿射重放 —— 内部略有斜切，但连续、无撕裂。
//
// 内容并不总能填满裁剪四边形（实测覆盖 80%~99%），所以可见轮廓其实是两者的交集。
// 那仍是个凸四边形、看着就是张卡，不值得为此放大内容去硬填满 —— 见文件头「已知边界」。
func (p *gioPainter) drawProjected(call op.CallOp, t layerTransform) {
	if t.opacity <= 0 || t.w <= 0 || t.h <= 0 {
		return
	}
	// 角点相对投影原点的偏移。无相机时原点即元素中心，偏移就是 ±w/2、±h/2；
	// 有相机时原点是场景中心，元素按自身位置投影，因此同场景元素灭点一致。
	ox, oy := t.origin()
	x0, y0 := t.cx-t.w/2, t.cy-t.h/2
	x1, y1 := x0+t.w, y0+t.h
	d0 := project3D(t, x0-ox, y0-oy) // 左上
	d1 := project3D(t, x1-ox, y0-oy) // 右上
	d2 := project3D(t, x0-ox, y1-oy) // 左下
	d3 := project3D(t, x1-ox, y1-oy) // 右下

	var quad clip.Path
	quad.Begin(p.ops)
	quad.MoveTo(d0)
	quad.LineTo(d1)
	quad.LineTo(d3)
	quad.LineTo(d2)
	quad.Close()
	cl := clip.Outline{Path: quad.End()}.Op().Push(p.ops)

	tr := op.Affine(projAffine(x0, y0, x1, y1, d0, d1, d2)).Push(p.ops)
	var os gpaint.OpacityStack
	if t.opacity < 1 {
		os = gpaint.PushOpacity(p.ops, t.opacity)
	}
	call.Add(p.ops)
	if t.opacity < 1 {
		os.Pop()
	}
	tr.Pop()
	cl.Pop()
}
