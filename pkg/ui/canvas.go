package ui

import (
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// whiteImage 用作 DrawTriangles 的纹理源（1x1 白色）。
// 延迟创建：避免在图形上下文就绪前（如单元测试导入包时）分配 GPU 资源。
var (
	whiteOnce sync.Once
	whiteImg  *ebiten.Image
)

func whiteImage() *ebiten.Image {
	whiteOnce.Do(func() {
		whiteImg = ebiten.NewImage(1, 1)
		whiteImg.Fill(White)
	})
	return whiteImg
}

// paint 是单线程的，以下 scratch 缓冲在各绘制原语间顺序复用，消除 paint 热路径每帧的
// 顶点/索引/路径分配。每次生成的顶点都会被 DrawTriangles 立即拷贝走，且各绘制函数内
// 不会在消费完当前顶点前再次占用同一缓冲，故顺序复用是安全的。
var (
	scratchVerts []ebiten.Vertex
	scratchIdx   []uint16
	scratchPath  vector.Path
)

// fillVerts / strokeVerts 用复用缓冲把 path 三角化，返回的切片仅在下一次绘制前有效。
func fillVerts(path *vector.Path) ([]ebiten.Vertex, []uint16) {
	scratchVerts, scratchIdx = path.AppendVerticesAndIndicesForFilling(scratchVerts[:0], scratchIdx[:0])
	return scratchVerts, scratchIdx
}
func strokeVerts(path *vector.Path, sop *vector.StrokeOptions) ([]ebiten.Vertex, []uint16) {
	scratchVerts, scratchIdx = path.AppendVerticesAndIndicesForStroke(scratchVerts[:0], scratchIdx[:0], sop)
	return scratchVerts, scratchIdx
}

// roundRectPath 构造圆角矩形路径（复用包级 scratchPath，调用方须在下次调用前消费完）。
func roundRectPath(x, y, w, h, r float32) *vector.Path {
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}
	p := &scratchPath
	p.Reset()
	if r <= 0 {
		p.MoveTo(x, y)
		p.LineTo(x+w, y)
		p.LineTo(x+w, y+h)
		p.LineTo(x, y+h)
		p.Close()
		return p
	}
	p.MoveTo(x+r, y)
	p.LineTo(x+w-r, y)
	p.ArcTo(x+w, y, x+w, y+r, r)
	p.LineTo(x+w, y+h-r)
	p.ArcTo(x+w, y+h, x+w-r, y+h, r)
	p.LineTo(x+r, y+h)
	p.ArcTo(x, y+h, x, y+h-r, r)
	p.LineTo(x, y+r)
	p.ArcTo(x, y, x+r, y, r)
	p.Close()
	return p
}

func colorScale(c Color) (float32, float32, float32, float32) {
	return float32(c.R) / 255, float32(c.G) / 255, float32(c.B) / 255, float32(c.A) / 255
}

// fillRoundRect 填充圆角矩形。
func fillRoundRect(dst *ebiten.Image, x, y, w, h, r float32, c Color) {
	if c.transparent() || w <= 0 || h <= 0 {
		return
	}
	if r <= 0 {
		vector.DrawFilledRect(dst, x, y, w, h, c, true)
		return
	}
	p := roundRectPath(x, y, w, h, r)
	vs, is := fillVerts(p)
	cr, cg, cb, ca := colorScale(c)
	for i := range vs {
		vs[i].SrcX, vs[i].SrcY = 0, 0
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = cr, cg, cb, ca
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true}
	dst.DrawTriangles(vs, is, whiteImage(), op)
}

// fillGradientRoundRect 用线性渐变填充圆角矩形：按顶点在渐变方向上的投影插值颜色。
func fillGradientRoundRect(dst *ebiten.Image, x, y, w, h, r float32, from, to Color, angleDeg float32) {
	if w <= 0 || h <= 0 {
		return
	}
	path := roundRectPath(x, y, w, h, r)
	vs, is := fillVerts(path)
	if len(is) == 0 {
		return
	}
	rad := float64(angleDeg) * math.Pi / 180
	dx, dy := float32(math.Cos(rad)), float32(math.Sin(rad))
	// 四角在方向上的投影范围，作为 t∈[0,1] 的归一化区间
	minP, maxP := float32(math.MaxFloat32), float32(-math.MaxFloat32)
	for _, c := range [4][2]float32{{x, y}, {x + w, y}, {x, y + h}, {x + w, y + h}} {
		pr := c[0]*dx + c[1]*dy
		if pr < minP {
			minP = pr
		}
		if pr > maxP {
			maxP = pr
		}
	}
	span := maxP - minP
	if span == 0 {
		span = 1
	}
	for i := range vs {
		t := (vs[i].DstX*dx + vs[i].DstY*dy - minP) / span
		cr, cg, cb, ca := colorScale(Mix(from, to, t))
		vs[i].SrcX, vs[i].SrcY = 0, 0
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = cr, cg, cb, ca
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true}
	dst.DrawTriangles(vs, is, whiteImage(), op)
}

// fillShadow 近似 box-shadow：以偏移+扩散为基准，向外叠加若干半透明圆角矩形，
// 得到中心接近目标不透明度、向外柔和衰减的投影。offY 正值向下。
func fillShadow(dst *ebiten.Image, x, y, w, h, r, offX, offY, blur, spread float32, c Color) {
	if c.transparent() || w <= 0 || h <= 0 {
		return
	}
	bx, by := x+offX-spread, y+offY-spread
	bw, bh := w+spread*2, h+spread*2
	br := r + spread
	if br < 0 {
		br = 0
	}
	if blur <= 0 {
		fillRoundRect(dst, bx, by, bw, bh, br, c)
		return
	}
	const layers = 5
	la := c.Alpha(1.0 / float32(layers)) // 各层等 alpha，重叠处向中心累积、向外衰减
	for i := layers; i >= 1; i-- {
		g := blur * float32(i) / float32(layers)
		fillRoundRect(dst, bx-g, by-g, bw+g*2, bh+g*2, br+g, la)
	}
	fillRoundRect(dst, bx, by, bw, bh, br, la)
}

// fillPathAt 以 (x,y) 为偏移填充一条矢量路径（实心图标）。
func fillPathAt(dst *ebiten.Image, path *vector.Path, x, y float32, c Color) {
	if path == nil || c.transparent() {
		return
	}
	vs, is := fillVerts(path)
	if len(is) == 0 {
		return
	}
	cr, cg, cb, ca := colorScale(c)
	for i := range vs {
		vs[i].DstX += x
		vs[i].DstY += y
		vs[i].SrcX, vs[i].SrcY = 0, 0
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = cr, cg, cb, ca
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true, FillRule: ebiten.FillRuleNonZero}
	dst.DrawTriangles(vs, is, whiteImage(), op)
}

// strokePathAt 以 (x,y) 为偏移描边一条矢量路径（线性图标，圆角端点/拐角）。
func strokePathAt(dst *ebiten.Image, path *vector.Path, x, y, width float32, c Color) {
	if path == nil || c.transparent() || width <= 0 {
		return
	}
	sop := &vector.StrokeOptions{Width: width, LineCap: vector.LineCapRound, LineJoin: vector.LineJoinRound}
	vs, is := strokeVerts(path, sop)
	if len(is) == 0 {
		return
	}
	cr, cg, cb, ca := colorScale(c)
	for i := range vs {
		vs[i].DstX += x
		vs[i].DstY += y
		vs[i].SrcX, vs[i].SrcY = 0, 0
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = cr, cg, cb, ca
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true}
	dst.DrawTriangles(vs, is, whiteImage(), op)
}

// maskRoundRect 用一个白色圆角矩形以 DestinationIn 混合裁剪 dst：
// 保留圆角内的像素、清除圆角外（含四个尖角）的像素。用于圆角裁剪。
func maskRoundRect(dst *ebiten.Image, r Rect, radius float32) {
	path := roundRectPath(r.X, r.Y, r.W, r.H, radius)
	vs, is := fillVerts(path)
	for i := range vs {
		vs[i].SrcX, vs[i].SrcY = 0, 0
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = 1, 1, 1, 1
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true, Blend: ebiten.BlendDestinationIn}
	dst.DrawTriangles(vs, is, whiteImage(), op)
}

// ebitenDrawLine 画一条直线（用于输入光标）。
func ebitenDrawLine(dst *ebiten.Image, x0, y0, x1, y1 float32, c Color) {
	vector.StrokeLine(dst, x0, y0, x1, y1, 1, c, false)
}

// strokeRoundRect 描边圆角矩形（边框绘制在内侧对齐）。
func strokeRoundRect(dst *ebiten.Image, x, y, w, h, r, width float32, c Color) {
	if c.transparent() || width <= 0 {
		return
	}
	inset := width / 2
	p := roundRectPath(x+inset, y+inset, w-width, h-width, r-inset)
	vs, is := strokeVerts(p, &vector.StrokeOptions{Width: width})
	cr, cg, cb, ca := colorScale(c)
	for i := range vs {
		vs[i].SrcX, vs[i].SrcY = 0, 0
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = cr, cg, cb, ca
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true}
	dst.DrawTriangles(vs, is, whiteImage(), op)
}
