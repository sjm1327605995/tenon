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
