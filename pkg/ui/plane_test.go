package ui

import "testing"

// 单应必须精确落在四个角上 —— 这正是仿射做不到的（6 个自由度只能匹配三角）。
// 如果这条不成立，地板的边就对不上卡牌的落位。
func TestPlaneHomographyHitsAllFourCorners(t *testing.T) {
	for _, tc := range []struct{ rx, ry, pers float32 }{
		{45, 0, 900}, {45, 20, 900}, {35, 20, 700}, {60, 15, 400}, {0, 50, 500}, {0, 0, 0},
	} {
		tr := layerTransform{
			cx: 400, cy: 300, w: 640, h: 420, scale: 1, opacity: 1,
			rotateX: tc.rx, rotateY: tc.ry, perspective: tc.pers,
		}
		d0, d1, d2, d3 := projCorners(tr)
		m := planeHomography(tr)
		x0, y0 := tr.cx-tr.w/2, tr.cy-tr.h/2
		x1, y1 := x0+tr.w, y0+tr.h

		for _, c := range []struct {
			src  pt
			want pt
			name string
		}{
			{pt{x0, y0}, d0, "左上"},
			{pt{x1, y0}, d1, "右上"},
			{pt{x0, y1}, d2, "左下"},
			{pt{x1, y1}, d3, "右下"}, // 仿射对不上的就是这一角
		} {
			got, ok := m.transform(c.src)
			if !ok {
				t.Fatalf("rx=%v ry=%v: %s角落在消失线上", tc.rx, tc.ry, c.name)
			}
			if absf(got.X-c.want.X) > 0.01 || absf(got.Y-c.want.Y) > 0.01 {
				t.Errorf("rx=%v ry=%v: %s角 -> %v, want %v", tc.rx, tc.ry, c.name, got, c.want)
			}
		}
	}
}

// 对比仿射：同样的第四角，仿射会偏出几十上百像素，单应必须是 0。
// 这是「为什么需要 PlaneImage」的量化依据。
func TestHomographyBeatsAffineOnFourthCorner(t *testing.T) {
	tr := layerTransform{
		cx: 400, cy: 300, w: 640, h: 420, scale: 1, opacity: 1,
		rotateX: 50, perspective: 1000,
	}
	_, _, _, d3 := projCorners(tr)
	x0, y0 := tr.cx-tr.w/2, tr.cy-tr.h/2
	corner := pt{x0 + tr.w, y0 + tr.h}

	aff := contentAffine(tr).transform(corner)
	hom, _ := planeHomography(tr).transform(corner)
	affErr := absf(aff.X-d3.X) + absf(aff.Y-d3.Y)
	homErr := absf(hom.X-d3.X) + absf(hom.Y-d3.Y)
	t.Logf("640x420 牌桌 @50°：仿射第四角偏 %.1fpx，单应偏 %.3fpx", affErr, homErr)

	if affErr < 50 {
		t.Fatalf("仿射误差只有 %.1fpx —— 前提不成立，本测试证明不了单应的价值", affErr)
	}
	if homErr > 0.01 {
		t.Errorf("单应第四角偏了 %.3fpx，应当精确", homErr)
	}
}

// 逆变换必须真的是逆：屏幕点反解回平面点，再正着投影回去，应当回到原处。
// 预变形靠的就是这条（对每个目标像素反解出源像素）。
func TestPlaneHomographyInverseRoundTrips(t *testing.T) {
	tr := layerTransform{
		cx: 400, cy: 300, w: 640, h: 420, scale: 1, opacity: 1,
		rotateX: 45, rotateY: 15, perspective: 800,
	}
	m := planeHomography(tr)
	inv := m.invert()
	x0, y0 := tr.cx-tr.w/2, tr.cy-tr.h/2

	maxErr := float32(0)
	for _, u := range []float32{0, 0.13, 0.5, 0.87, 1} {
		for _, v := range []float32{0, 0.29, 0.5, 0.71, 1} {
			src := pt{x0 + tr.w*u, y0 + tr.h*v}
			scr, ok := m.transform(src)
			if !ok {
				t.Fatalf("正变换失败 u=%v v=%v", u, v)
			}
			back, ok := inv.transform(scr)
			if !ok {
				t.Fatalf("逆变换失败 u=%v v=%v", u, v)
			}
			if e := absf(back.X-src.X) + absf(back.Y-src.Y); e > maxErr {
				maxErr = e
			}
		}
	}
	t.Logf("25 个采样点往返最大误差 = %.4fpx", maxErr)
	if maxErr > 0.05 {
		t.Errorf("往返误差 %.4fpx 过大 —— 逆变换不准，预变形会整体错位", maxErr)
	}
}
