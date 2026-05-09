package render

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// BorderSlice 定义九宫格切片的边距。
// Left/Top/Right/Bottom 分别表示从源图四边到中心可拉伸区域的距离（像素）。
type BorderSlice struct {
	Left, Top, Right, Bottom int
}

// IsZero 判断切片是否为零值（即不进行九宫格处理）。
func (bs BorderSlice) IsZero() bool {
	return bs.Left == 0 && bs.Top == 0 && bs.Right == 0 && bs.Bottom == 0
}

// RenderNinePatch 绘制九宫格图片，用于游戏 UI 中需要任意缩放的面板背景、
// 按钮底图、血条边框等场景。
type RenderNinePatch struct {
	BaseRenderObject

	Source    *ebiten.Image
	Slice     BorderSlice
	TintColor color.Color
}

func NewRenderNinePatch() *RenderNinePatch {
	r := &RenderNinePatch{}
	r.BaseRenderObject.Init(r)
	r.yoga = yoga.NewNode()
	return r
}

func (r *RenderNinePatch) Paint(screen *ebiten.Image, offset Offset) {
	if r.Source == nil || r.Slice.IsZero() {
		return
	}
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	dstX := float64(offset.X + bounds.X)
	dstY := float64(offset.Y + bounds.Y)
	dstW := float64(bounds.Width)
	dstH := float64(bounds.Height)

	srcBounds := r.Source.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	s := r.Slice

	// 确保切片不超过源图尺寸
	if s.Left+s.Right >= srcW {
		s.Left = srcW / 4
		s.Right = srcW / 4
	}
	if s.Top+s.Bottom >= srcH {
		s.Top = srcH / 4
		s.Bottom = srcH / 4
	}

	// 中心区域在源图中的尺寸
	centerSrcW := srcW - s.Left - s.Right
	centerSrcH := srcH - s.Top - s.Bottom
	if centerSrcW <= 0 || centerSrcH <= 0 {
		return
	}

	// 中心区域在目标中的尺寸（总尺寸减去四角固定区域）
	centerDstW := dstW - float64(s.Left) - float64(s.Right)
	centerDstH := dstH - float64(s.Top) - float64(s.Bottom)
	if centerDstW < 0 {
		centerDstW = 0
	}
	if centerDstH < 0 {
		centerDstH = 0
	}

	// 缩放比例（用于四边和中心）
	scaleX := 1.0
	if centerSrcW > 0 {
		scaleX = centerDstW / float64(centerSrcW)
	}
	scaleY := 1.0
	if centerSrcH > 0 {
		scaleY = centerDstH / float64(centerSrcH)
	}

	// 预计算 9 个源区域
	regions := [9]image.Rectangle{
		// 0: 左上
		{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: s.Left, Y: s.Top}},
		// 1: 上
		{Min: image.Point{X: s.Left, Y: 0}, Max: image.Point{X: srcW - s.Right, Y: s.Top}},
		// 2: 右上
		{Min: image.Point{X: srcW - s.Right, Y: 0}, Max: image.Point{X: srcW, Y: s.Top}},
		// 3: 左
		{Min: image.Point{X: 0, Y: s.Top}, Max: image.Point{X: s.Left, Y: srcH - s.Bottom}},
		// 4: 中心
		{Min: image.Point{X: s.Left, Y: s.Top}, Max: image.Point{X: srcW - s.Right, Y: srcH - s.Bottom}},
		// 5: 右
		{Min: image.Point{X: srcW - s.Right, Y: s.Top}, Max: image.Point{X: srcW, Y: srcH - s.Bottom}},
		// 6: 左下
		{Min: image.Point{X: 0, Y: srcH - s.Bottom}, Max: image.Point{X: s.Left, Y: srcH}},
		// 7: 下
		{Min: image.Point{X: s.Left, Y: srcH - s.Bottom}, Max: image.Point{X: srcW - s.Right, Y: srcH}},
		// 8: 右下
		{Min: image.Point{X: srcW - s.Right, Y: srcH - s.Bottom}, Max: image.Point{X: srcW, Y: srcH}},
	}

	// 预计算 9 个目标位置
	positions := [9][2]float64{
		{dstX, dstY},
		{dstX + float64(s.Left), dstY},
		{dstX + float64(s.Left) + centerDstW, dstY},
		{dstX, dstY + float64(s.Top)},
		{dstX + float64(s.Left), dstY + float64(s.Top)},
		{dstX + float64(s.Left) + centerDstW, dstY + float64(s.Top)},
		{dstX, dstY + float64(s.Top) + centerDstH},
		{dstX + float64(s.Left), dstY + float64(s.Top) + centerDstH},
		{dstX + float64(s.Left) + centerDstW, dstY + float64(s.Top) + centerDstH},
	}

	// 预计算 9 个目标尺寸
	sizes := [9][2]float64{
		{float64(s.Left), float64(s.Top)},
		{centerDstW, float64(s.Top)},
		{float64(s.Right), float64(s.Top)},
		{float64(s.Left), centerDstH},
		{centerDstW, centerDstH},
		{float64(s.Right), centerDstH},
		{float64(s.Left), float64(s.Bottom)},
		{centerDstW, float64(s.Bottom)},
		{float64(s.Right), float64(s.Bottom)},
	}

	// 缩放因子：0/2/6/8 不缩放；1/7 水平缩放；3/5 垂直缩放；4 双向缩放
	scales := [9][2]float64{
		{1, 1}, {scaleX, 1}, {1, 1},
		{1, scaleY}, {scaleX, scaleY}, {1, scaleY},
		{1, 1}, {scaleX, 1}, {1, 1},
	}

	for i := 0; i < 9; i++ {
		sw, sh := sizes[i][0], sizes[i][1]
		if sw <= 0 || sh <= 0 {
			continue
		}
		px, py := positions[i][0], positions[i][1]
		sx, sy := scales[i][0], scales[i][1]

		sub := r.Source.SubImage(regions[i]).(*ebiten.Image)
		if sub == nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(px, py)
		if r.TintColor != nil {
			op.ColorScale.ScaleWithColor(r.TintColor)
		}
		screen.DrawImage(sub, op)
	}
}
