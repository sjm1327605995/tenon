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
// 把已录制的图层（op.Record 宏）按透视投影贴回屏幕，实现 CSS 的
// transform: perspective(p) rotateX(rx) rotateY(ry) translateZ(tz)。内容仍是平面，
// 只是被投影到一个四边形上 —— 所谓「伪 3D」。
//
// 做法：Ebiten 那套（离屏纹理 + 带 UV 的贴图网格）在 gio 上行不通，gio 的公开 API
// 既不提供渲染到纹理、也没有逐顶点贴图。但 gio 的宏（op.CallOp）可以被反复重放，
// 于是改成：把元素划成 n×n 网格，每个小格「裁剪到投影后的四边形 + 施加该格的仿射」
// 再重放一次整棵子树。单个仿射无法表达透视，但每个小格内的透视近似为仿射，网格够密
// 就足够像。文字与图片也因此会跟着变形（若逐图元投影则做不到，字形轮廓无法重投影）。
//
// 代价是子树被重放 n×n 次，所以网格取得克制。

// gio3DGrid 是每个方向上的网格数。密度越高越接近真实透视，代价是子树重放次数按平方增长。
const gio3DGrid = 8

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
	return f32.Pt(
		t.cx+float32(x2*cosZ-y2*sinZ)+t.tx,
		t.cy+float32(x2*sinZ+y2*cosZ)+t.ty,
	)
}

// cellAffine 求把源矩形 [sx0,sx1]×[sy0,sy1] 映射到 (d0,d1,d2) 的仿射矩阵，
// 其中 d0/d1/d2 分别是该格左上/右上/左下角的投影点。三点唯一确定一个仿射；
// 第四角的偏差就是本格内以仿射近似透视的误差，网格越密越小。
func cellAffine(sx0, sy0, sx1, sy1 float32, d0, d1, d2 f32.Point) f32.Affine2D {
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

// drawProjected 用网格重放把录制好的图层按透视贴回。
func (p *gioPainter) drawProjected(call op.CallOp, t layerTransform) {
	if t.opacity <= 0 || t.w <= 0 || t.h <= 0 {
		return
	}
	const n = gio3DGrid
	x0, y0 := t.cx-t.w/2, t.cy-t.h/2

	// 先算好 (n+1)² 个网格顶点的投影，供相邻格共用（保证格与格严丝合缝）。
	pts := make([]f32.Point, (n+1)*(n+1))
	for j := 0; j <= n; j++ {
		for i := 0; i <= n; i++ {
			ox := -t.w/2 + t.w*float32(i)/n
			oy := -t.h/2 + t.h*float32(j)/n
			pts[j*(n+1)+i] = project3D(t, ox, oy)
		}
	}

	var os gpaint.OpacityStack
	if t.opacity < 1 {
		os = gpaint.PushOpacity(p.ops, t.opacity)
	}
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			d0 := pts[j*(n+1)+i]       // 左上
			d1 := pts[j*(n+1)+i+1]     // 右上
			d2 := pts[(j+1)*(n+1)+i]   // 左下
			d3 := pts[(j+1)*(n+1)+i+1] // 右下

			sx0 := x0 + t.w*float32(i)/n
			sy0 := y0 + t.h*float32(j)/n
			sx1 := x0 + t.w*float32(i+1)/n
			sy1 := y0 + t.h*float32(j+1)/n

			// 裁到投影后的四边形：相邻格共用顶点，因此彼此贴合、不重不漏。
			var quad clip.Path
			quad.Begin(p.ops)
			quad.MoveTo(d0)
			quad.LineTo(d1)
			quad.LineTo(d3)
			quad.LineTo(d2)
			quad.Close()
			cl := clip.Outline{Path: quad.End()}.Op().Push(p.ops)

			tr := op.Affine(cellAffine(sx0, sy0, sx1, sy1, d0, d1, d2)).Push(p.ops)
			call.Add(p.ops) // 重放整棵子树，只有本格内的部分能透过裁剪显出来
			tr.Pop()
			cl.Pop()
		}
	}
	if t.opacity < 1 {
		os.Pop()
	}
}
