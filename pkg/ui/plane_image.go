package ui

import (
	"fmt"
	"image"
	"image/draw"
	"math"
)

// ---- PlaneImage：Scene3D 的地板贴图 ----
//
// 问题：Scene3D 用仿射画内容，偏差随元素尺寸增长（见 gio_3d.go 的「已知边界」）。
// 卡牌够小、不显眼，但地板不行 —— 它是参照系（卡要落在它的格子里），而且它是场上最大的
// 元素：640x420 的桌面在 50° 下第四角偏 211px，画出来是斜切的平行四边形，格线不朝灭点收敛。
//
// 关键观察：地板是静态的。一局里图和投影参数都不变，所以不必在每帧的绘制路径里对付透视 ——
// 按精确单应把图 CPU 预变形一次，得到的就是一张普通位图，正着贴上去即可，每帧零额外成本。
//
// 变形用的单应取自 projCorners，与卡牌共用同一份 project3D，因此格子和卡不会错位 ——
// 这正是这件事必须由 tenon 做、而不是由调用方自己拼的原因。
//
// 只适用于「铺满场景平面的静态图」。倾斜的视频、逐帧变的图不适用：每帧重新变形会很贵。

// planeCacheEntry 是一次预变形的结果。
type planeCacheEntry struct {
	bmp    bitmap
	origin Rect // 变形结果在屏幕上的位置（投影四边形的包围盒）
}

// planeCache 按 (key + 场景参数 + 尺寸) 缓存预变形结果。场景参数变了（窗口缩放、改倾角）
// 就得重算，所以它们必须进缓存键 —— 否则窗口一缩放地板就和卡牌错位。
var planeCache = map[string]planeCacheEntry{}

// planeCacheKey 把 key 与所有影响投影结果的量拼成缓存键。
// 少拼一个量，就会在那个量变化时悄悄用上过期的变形结果。
func planeCacheKey(key string, t layerTransform) string {
	return fmt.Sprintf("%s|%.2f,%.2f|%.2f,%.2f|%.3f,%.3f,%.3f,%.3f|%.3f,%.3f|%.3f,%.3f",
		key, t.cx, t.cy, t.w, t.h,
		t.rotateX, t.rotateY, t.transZ, t.perspective,
		t.scale, t.rotate, t.tx, t.ty)
}

// planeBBox 返回投影四边形的屏幕包围盒（向外取整到整像素）。
func planeBBox(t layerTransform) Rect {
	d0, d1, d2, d3 := projCorners(t)
	minX := minf(minf(d0.X, d1.X), minf(d2.X, d3.X))
	maxX := maxf(maxf(d0.X, d1.X), maxf(d2.X, d3.X))
	minY := minf(minf(d0.Y, d1.Y), minf(d2.Y, d3.Y))
	maxY := maxf(maxf(d0.Y, d1.Y), maxf(d2.Y, d3.Y))
	x0 := float32(math.Floor(float64(minX)))
	y0 := float32(math.Floor(float64(minY)))
	x1 := float32(math.Ceil(float64(maxX)))
	y1 := float32(math.Ceil(float64(maxY)))
	return Rect{x0, y0, x1 - x0, y1 - y0}
}

func maxf(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// warpPlane 把 src 按 t 的精确投影预变形，返回一张覆盖投影四边形包围盒的 RGBA。
// 与 warpQuad 同一实现 —— unitToQuad 的逆映射直接给出单位坐标，源图铺满单位正方形。
func warpPlane(src image.Image, t layerTransform, box Rect) *image.RGBA {
	d0, d1, d2, d3 := projCorners(t)
	return warpQuad(src, d0, d1, d2, d3, box)
}

// warpLoop 是逆映射变形的公共像素循环：对 box 内每个目标像素，inv 反解出源图单位
// 坐标；落在消失线上或单位正方形之外的像素留透明（包围盒四角不会糊上边缘色）。
func warpLoop(src image.Image, inv homography, box Rect) *image.RGBA {
	w, h := int(box.W), int(box.H)
	if w <= 0 || h <= 0 {
		return nil
	}
	out := image.NewRGBA(image.Rect(0, 0, w, h))

	sb := src.Bounds()
	sw, sh := float32(sb.Dx()), float32(sb.Dy())

	// 转成 RGBA 直接索引，避免逐像素走 At() 的接口开销与颜色模型转换
	var rgba *image.RGBA
	if r, ok := src.(*image.RGBA); ok {
		rgba = r
	} else {
		rgba = image.NewRGBA(sb)
		draw.Draw(rgba, sb, src, sb.Min, draw.Src)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// 取目标像素中心反解，避免整体偏半个像素
			p, ok := inv.transform(pt{box.X + float32(x) + 0.5, box.Y + float32(y) + 0.5})
			if !ok {
				continue // 落在消失线上：该像素不属于这个平面
			}
			su, sv := p.X*sw, p.Y*sh
			if su < 0 || sv < 0 || su >= sw || sv >= sh {
				continue // 四边形之外：留透明
			}
			si := rgba.PixOffset(sb.Min.X+int(su), sb.Min.Y+int(sv))
			di := out.PixOffset(x, y)
			copy(out.Pix[di:di+4], rgba.Pix[si:si+4])
		}
	}
	return out
}

// planeBitmap 取（或建）预变形结果。命中缓存时零成本。
func planeBitmap(key string, src image.Image, t layerTransform) (bitmap, Rect, bool) {
	box := planeBBox(t)
	if box.W <= 0 || box.H <= 0 {
		return nil, Rect{}, false
	}
	ck := planeCacheKey(key, t)
	if e, ok := planeCache[ck]; ok {
		return e.bmp, e.origin, true
	}
	img := warpPlane(src, t, box)
	if img == nil {
		return nil, Rect{}, false
	}
	bmp := backendNewBitmap(img)
	planeCache[ck] = planeCacheEntry{bmp: bmp, origin: box}
	return bmp, box, true
}

// paintPlaneImage 把地板按精确单应预变形后正着贴上，绕开 drawProjected 的仿射近似。
// 返回是否已处理；退化情形（无场景参数、尺寸为 0）返回 false，交回常规路径。
func paintPlaneImage(p painter, rn *renderNode, cam *camera3D) bool {
	t := layerOf(rn, cam)
	if !t.is3D() { // 场景没倾斜：普通贴图即可，不必变形
		return false
	}
	bmp, box, ok := planeBitmap(rn.imgSrc, rn.planeImg, t)
	if !ok {
		return false
	}
	p.DrawImage(bmp, box, rn.opacity)
	return true
}

// ---- 任意四边形的精确单应贴图（投影层内普通图片用）----
//
// 逐元素投影把斜切压到元素尺寸级（见 gio_3d.go），但「元素自身就是一整块图」时误差
// 仍随图的大小增长 —— 倾斜卡片上铺满的照片是最显眼的例子。这样的元素走这里：按四角
// 真实投影预变形后正着贴，与 PlaneImage 同一份数学（unitToQuad），只是四边形直接
// 来自 projCorners，不必凑一个 layerTransform。

// quadCache 缓存「源图 + 目标四边形」的预变形结果。动画期间四边形每帧都变、每帧
// 都会新增条目，所以设上限：满了整体清空（动画是突发性的，清空后热点条目会迅速重建）。
var quadCache = map[string]planeCacheEntry{}

const quadCacheLimit = 64

func quadCacheKey(key string, d0, d1, d2, d3 pt) string {
	return fmt.Sprintf("%s|%.2f,%.2f|%.2f,%.2f|%.2f,%.2f|%.2f,%.2f",
		key, d0.X, d0.Y, d1.X, d1.Y, d2.X, d2.Y, d3.X, d3.Y)
}

// quadBBox 返回四边形的屏幕包围盒（向外取整到整像素）。
func quadBBox(d0, d1, d2, d3 pt) Rect {
	minX := minf(minf(d0.X, d1.X), minf(d2.X, d3.X))
	maxX := maxf(maxf(d0.X, d1.X), maxf(d2.X, d3.X))
	minY := minf(minf(d0.Y, d1.Y), minf(d2.Y, d3.Y))
	maxY := maxf(maxf(d0.Y, d1.Y), maxf(d2.Y, d3.Y))
	x0 := float32(math.Floor(float64(minX)))
	y0 := float32(math.Floor(float64(minY)))
	x1 := float32(math.Ceil(float64(maxX)))
	y1 := float32(math.Ceil(float64(maxY)))
	return Rect{x0, y0, x1 - x0, y1 - y0}
}

// warpQuad 把 src（铺满单位正方形）按逆映射映到四边形，返回覆盖包围盒的 RGBA。
// 与 warpPlane 同一套 warpLoop，单应由四角直接构造。
func warpQuad(src image.Image, d0, d1, d2, d3 pt, box Rect) *image.RGBA {
	return warpLoop(src, unitToQuad(d0, d1, d3, d2).invert(), box)
}

// quadBitmap 取（或建）「源图 -> 四边形」的预变形结果；命中缓存时零成本。
func quadBitmap(key string, src image.Image, d0, d1, d2, d3 pt) (bitmap, Rect, bool) {
	box := quadBBox(d0, d1, d2, d3)
	if box.W <= 0 || box.H <= 0 {
		return nil, Rect{}, false
	}
	ck := quadCacheKey(key, d0, d1, d2, d3)
	if e, ok := quadCache[ck]; ok {
		return e.bmp, e.origin, true
	}
	img := warpQuad(src, d0, d1, d2, d3, box)
	if img == nil {
		return nil, Rect{}, false
	}
	if len(quadCache) >= quadCacheLimit {
		quadCache = map[string]planeCacheEntry{} // 见 quadCache 注释：动画突发清空
	}
	bmp := backendNewBitmap(img)
	quadCache[ck] = planeCacheEntry{bmp: bmp, origin: box}
	return bmp, box, true
}
