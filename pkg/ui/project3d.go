package ui

import "math"

// ---- 伪 3D 的投影数学 ----
//
// 后端中立：这里只算坐标，不认识任何 gio 类型（见 backend.go 的边界规则）。
// 绘制端（gio_3d.go）把结果转成 gio 的仿射去画，命中端（render.go）取同一个仿射的逆 ——
// 两边共用这里的一套数学，才不会出现「画在一处、点在另一处」。

// pt 是一个二维点。
type pt struct{ X, Y float32 }

// affine2D 是 2x3 仿射：(x,y) -> (sx*x + hx*y + ox, hy*x + sy*y + oy)。
type affine2D struct{ sx, hx, ox, hy, sy, oy float32 }

func (a affine2D) transform(p pt) pt {
	return pt{a.sx*p.X + a.hx*p.Y + a.ox, a.hy*p.X + a.sy*p.Y + a.oy}
}

// invert 返回逆仿射；退化（行列式为 0）时返回恒等，调用方会因此落回未变换的坐标。
func (a affine2D) invert() affine2D {
	det := a.sx*a.sy - a.hx*a.hy
	if det == 0 {
		return affine2D{1, 0, 0, 0, 1, 0}
	}
	inv := 1 / det
	sx, hx := a.sy*inv, -a.hx*inv
	hy, sy := -a.hy*inv, a.sx*inv
	return affine2D{
		sx, hx, -(sx*a.ox + hx*a.oy),
		hy, sy, -(hy*a.ox + sy*a.oy),
	}
}

// project3D 把元素平面上的一点（相对投影原点的偏移）投影到屏幕坐标。
// 顺序与 CSS 的 rotateX(rx) rotateY(ry) 书写序一致：先绕 Y 再绕 X，然后透视除法，
// 最后叠加 2D 的缩放/绕 Z 旋转/平移。
func project3D(t layerTransform, ox, oy float32) pt {
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
	orgX, orgY := t.origin() // 无相机=元素中心，有相机=场景中心
	return pt{
		orgX + float32(x2*cosZ-y2*sinZ) + t.tx,
		orgY + float32(x2*sinZ+y2*cosZ) + t.ty,
	}
}

// projAffine 求把源矩形 [sx0,sx1]x[sy0,sy1] 映射到 (d0,d1,d2) 的仿射，
// 其中 d0/d1/d2 分别是左上/右上/左下角的投影点。三点唯一确定一个仿射。
func projAffine(sx0, sy0, sx1, sy1 float32, d0, d1, d2 pt) affine2D {
	w, h := sx1-sx0, sy1-sy0
	if w == 0 || h == 0 {
		return affine2D{1, 0, 0, 0, 1, 0}
	}
	sx := (d1.X - d0.X) / w // dx/dsx
	hx := (d2.X - d0.X) / h // dx/dsy
	hy := (d1.Y - d0.Y) / w // dy/dsx
	sy := (d2.Y - d0.Y) / h // dy/dsy
	return affine2D{
		sx, hx, d0.X - sx0*sx - sy0*hx,
		hy, sy, d0.Y - sx0*hy - sy0*sy,
	}
}

// projCorners 返回元素四角的真实投影：左上/右上/左下/右下。
//
// 角点取的是「相对投影原点的偏移」。无相机时原点即元素中心，偏移就是 ±w/2、±h/2；
// 有相机时原点是场景中心，元素按自身位置投影，因此同场景元素灭点一致。
func projCorners(t layerTransform) (d0, d1, d2, d3 pt) {
	orgX, orgY := t.origin()
	x0, y0 := t.cx-t.w/2, t.cy-t.h/2
	x1, y1 := x0+t.w, y0+t.h
	return project3D(t, x0-orgX, y0-orgY), project3D(t, x1-orgX, y0-orgY),
		project3D(t, x0-orgX, y1-orgY), project3D(t, x1-orgX, y1-orgY)
}

// contentAffine 是内容实际被施加的那个仿射（由前三角确定）。
// 绘制用它，命中测试用它的逆 —— 同一个矩阵，两边不可能对不上。
func contentAffine(t layerTransform) affine2D {
	d0, d1, d2, _ := projCorners(t)
	x0, y0 := t.cx-t.w/2, t.cy-t.h/2
	return projAffine(x0, y0, x0+t.w, y0+t.h, d0, d1, d2)
}

// projErr 量化内容的斜切程度：仿射把第四角映到哪、与真透视差多少。轮廓由裁剪保证，
// 所以这不是形状误差，只是元素内部图案的偏差。元素越大、倾角越陡越大。仅用于测试与诊断。
func projErr(t layerTransform) float32 {
	_, _, _, d3 := projCorners(t)
	x0, y0 := t.cx-t.w/2, t.cy-t.h/2
	got := contentAffine(t).transform(pt{x0 + t.w, y0 + t.h})
	return absf(got.X-d3.X) + absf(got.Y-d3.Y)
}

// ---- 平面单应（PlaneImage 用）----
//
// 一个平面被透视投影后的像，数学上就是一个单应变换（homography）：3x3 齐次矩阵，8 个
// 自由度，正好能把矩形映成任意凸四边形 —— 这是仿射（6 个自由度）做不到的，也正是
// drawProjected 只能近似的原因。
//
// 但如果内容是一张静态图，就不必在每帧的绘制路径里对付它：可以离线按精确单应把图预变形
// 一次，结果是一张普通位图，正着贴上去即可。这里提供那份数学。
//
// 关键在于它必须与 Scene3D 给子元素用的投影完全一致，否则地板的格子和卡牌会对不上。
// 所以四角一律取自 projCorners —— 与卡牌同一份 project3D，不另算一套。

// homography 是 3x3 齐次变换：(x,y,1) -> (a*x+b*y+c, d*x+e*y+f, g*x+h*y+i)，
// 结果再除以第三个分量。
type homography struct{ a, b, c, d, e, f, g, h, i float32 }

// transform 施加变换；返回的 ok 为假表示该点落在消失线上（第三分量为 0），结果无意义。
func (m homography) transform(p pt) (pt, bool) {
	w := m.g*p.X + m.h*p.Y + m.i
	if w == 0 {
		return pt{}, false
	}
	return pt{(m.a*p.X + m.b*p.Y + m.c) / w, (m.d*p.X + m.e*p.Y + m.f) / w}, true
}

// unitToQuad 求把单位正方形 (0,0),(1,0),(1,1),(0,1) 映到给定四边形的单应。
// 经典构造（Heckbert），四点顺序须为 左上/右上/右下/左下。
func unitToQuad(d0, d1, d3, d2 pt) homography {
	// d0=左上(0,0) d1=右上(1,0) d3=右下(1,1) d2=左下(0,1)
	sx := d0.X - d1.X + d3.X - d2.X
	sy := d0.Y - d1.Y + d3.Y - d2.Y
	if sx == 0 && sy == 0 { // 退化为仿射（平行四边形），无透视
		return homography{
			d1.X - d0.X, d2.X - d0.X, d0.X,
			d1.Y - d0.Y, d2.Y - d0.Y, d0.Y,
			0, 0, 1,
		}
	}
	dx1, dx2 := d1.X-d3.X, d2.X-d3.X
	dy1, dy2 := d1.Y-d3.Y, d2.Y-d3.Y
	den := dx1*dy2 - dx2*dy1
	if den == 0 {
		return homography{1, 0, 0, 0, 1, 0, 0, 0, 1} // 四点共线：无解，返回恒等
	}
	g := (sx*dy2 - dx2*sy) / den
	h := (dx1*sy - sx*dy1) / den
	return homography{
		d1.X - d0.X + g*d1.X, d2.X - d0.X + h*d2.X, d0.X,
		d1.Y - d0.Y + g*d1.Y, d2.Y - d0.Y + h*d2.Y, d0.Y,
		g, h, 1,
	}
}

// invert 求逆（伴随矩阵；齐次矩阵差一个常数因子不影响结果，故不必除行列式）。
func (m homography) invert() homography {
	return homography{
		m.e*m.i - m.f*m.h, m.c*m.h - m.b*m.i, m.b*m.f - m.c*m.e,
		m.f*m.g - m.d*m.i, m.a*m.i - m.c*m.g, m.c*m.d - m.a*m.f,
		m.d*m.h - m.e*m.g, m.b*m.g - m.a*m.h, m.a*m.e - m.b*m.d,
	}
}

// planeHomography 返回「元素自身的矩形 -> 它在场景中被投影出的四边形」的单应。
// 四角取自 projCorners，因此与同场景的卡牌共用同一份投影，格子和卡不会错位。
func planeHomography(t layerTransform) homography {
	d0, d1, d2, d3 := projCorners(t)
	// 先把元素矩形归一化到单位正方形，再映到投影四边形
	x0, y0 := t.cx-t.w/2, t.cy-t.h/2
	toUnit := homography{
		1 / t.w, 0, -x0 / t.w,
		0, 1 / t.h, -y0 / t.h,
		0, 0, 1,
	}
	return unitToQuad(d0, d1, d3, d2).mul(toUnit)
}

// mul 返回 m∘n（先 n 后 m）。
func (m homography) mul(n homography) homography {
	return homography{
		m.a*n.a + m.b*n.d + m.c*n.g, m.a*n.b + m.b*n.e + m.c*n.h, m.a*n.c + m.b*n.f + m.c*n.i,
		m.d*n.a + m.e*n.d + m.f*n.g, m.d*n.b + m.e*n.e + m.f*n.h, m.d*n.c + m.e*n.f + m.f*n.i,
		m.g*n.a + m.h*n.d + m.i*n.g, m.g*n.b + m.h*n.e + m.i*n.h, m.g*n.c + m.h*n.f + m.i*n.i,
	}
}
