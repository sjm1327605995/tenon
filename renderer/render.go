package renderer

import (
	"github.com/sjm1327605995/tenon/core/context/node"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/millken/yoga"
)

var whiteImage *ebiten.Image

func init() {
	img := ebiten.NewImage(3, 3)
	img.Fill(color.White)
	whiteImage = img.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
}
func NewEbitenRenderer() *Ebiten {
	return &Ebiten{}
}

type Ebiten struct {
	scaleFactor float32
	startX      float32
	startY      float32
	Screen      *ebiten.Image
}

func (e *Ebiten) UpdateScaleFactor(scaleFactor float32) {
	e.scaleFactor = scaleFactor
}
func (e *Ebiten) SetScreen(screen *ebiten.Image) {
	e.Screen = screen
}
func (e *Ebiten) Image(startX, startY float32, node *yoga.Node, img image.Image) error {
	opts := &ebiten.DrawImageOptions{}
	bounds := img.Bounds()
	w, h := node.StyleGetWidth(), node.StyleGetHeight()

	opts.GeoM.Scale(float64(w/float32(bounds.Dx())), float64(h/float32(bounds.Dy())))
	opts.GeoM.Translate(float64(startX), float64(startY))

	ebitenImage, ok := img.(*ebiten.Image)
	if !ok {
		ebitenImage = ebiten.NewImageFromImage(img)
	}
	e.Screen.DrawImage(ebitenImage, opts)
	return nil
}
func (e *Ebiten) Rectangle(startX, startY float32, node *yoga.Node, radius node.Radius, color color.RGBA) error {
	//包含了border的宽度
	w, h := node.StyleGetWidth(), node.StyleGetHeight()
	if w == 0 || h == 0 {
		return nil
	}
	//leftBorder := node.StyleGetBorder(yoga.EdgeLeft)
	//topBorder := node.StyleGetBorder(yoga.EdgeTop)
	//rightBorder := node.StyleGetBorder(yoga.EdgeRight)
	//bottomBorder := node.StyleGetBorder(yoga.EdgeBottom)
	//allBorder := node.StyleGetBorder(yoga.EdgeAll)
	//if allBorder > 0 {
	//	leftBorder = allBorder
	//	topBorder = allBorder
	//	rightBorder = allBorder
	//	bottomBorder = allBorder
	//}
	//需要要减去border的宽度 且如何绘制border

	if radius.TopLeft > 0 || radius.TopRight > 0 || radius.BottomRight > 0 || radius.BottomLeft > 0 {
		renderFillRoundedRect(e.Screen, w, h, startX, startY, radius, color)
	} else {
		vector.DrawFilledRect(
			e.Screen,
			startX,
			startY,
			w, h,
			color,
			true,
		)
	}
	return nil
}

// renderFillRoundedRect 绘制带圆角的矩形
func renderFillRoundedRect(screen *ebiten.Image, w, h, x, y float32, cornerRadius node.Radius, _color color.RGBA) {
	// 颜色归一化（float32）
	r := float32(_color.R) / 255.0
	g := float32(_color.G) / 255.0
	b := float32(_color.B) / 255.0
	a := float32(_color.A) / 255.0

	// 安全圆角计算：确保半径不会过大
	safeRadius := func(r, maxPossible float32) float32 {
		if r <= 0 {
			return 0
		}
		if r > maxPossible {
			return maxPossible
		}
		return r
	}

	// 计算每个角的最大可能半径
	maxTLR := min(w/2, h/2)
	maxTRR := min(w/2, h/2)
	maxBRR := min(w/2, h/2)
	maxBLR := min(w/2, h/2)

	tlR := safeRadius(cornerRadius.TopLeft, maxTLR)
	trR := safeRadius(cornerRadius.TopRight, maxTRR)
	brR := safeRadius(cornerRadius.BottomRight, maxBRR)
	blR := safeRadius(cornerRadius.BottomLeft, maxBLR)

	// 矩形边界
	x1, y1 := x, y
	x2, y2 := x+w, y+h

	// 顶点管理
	var vertices []ebiten.Vertex
	addVertex := func(dstX, dstY float32) uint16 {
		vertices = append(vertices, ebiten.Vertex{
			DstX:   dstX,
			DstY:   dstY,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		})
		return uint16(len(vertices) - 1)
	}

	// 生成完整轮廓点
	points := generateRoundedRectPoints(x1, y1, x2, y2, tlR, trR, brR, blR, 16)

	// 中心顶点（用于三角化）
	centerX, centerY := (x1+x2)/2, (y1+y2)/2
	centerIdx := addVertex(centerX, centerY)

	// 三角化：中心到轮廓的每个三角形
	var indices []uint16
	for i := 0; i < len(points)-1; i++ {
		p1 := addVertex(points[i][0], points[i][1])
		p2 := addVertex(points[i+1][0], points[i+1][1])
		indices = append(indices, centerIdx, p1, p2)
	}

	// 闭合最后一个三角形
	last := addVertex(points[len(points)-1][0], points[len(points)-1][1])
	first := addVertex(points[0][0], points[0][1])
	indices = append(indices, centerIdx, last, first)

	// 执行绘制
	if len(vertices) > 0 && len(indices) > 0 {
		screen.DrawTriangles(vertices, indices, whiteImage, &ebiten.DrawTrianglesOptions{
			AntiAlias: true,
		})
	}
}

// generateRoundedRectPoints 生成圆角矩形的连续轮廓点
func generateRoundedRectPoints(x1, y1, x2, y2, tlR, trR, brR, blR float32, segments int) [][2]float32 {
	var points [][2]float32

	// 左上角圆角
	if tlR > 0 {
		centerX, centerY := x1+tlR, y1+tlR
		startAngle := math.Pi       // 180度
		endAngle := math.Pi * 3 / 2 // 270度
		// 生成圆弧点（已处理类型转换）
		arcPoints := generateArcPoints(centerX, centerY, tlR, startAngle, endAngle, segments)
		points = append(points, arcPoints...)
	} else {
		points = append(points, [2]float32{x1, y1})
	}

	// 上边缘直线段
	topEdgeEndX := x2 - trR
	points = append(points, [2]float32{topEdgeEndX, y1})

	// 右上角圆角
	if trR > 0 {
		centerX, centerY := x2-trR, y1+trR
		startAngle := math.Pi * 3 / 2 // 270度
		endAngle := math.Pi * 2       // 360度
		arcPoints := generateArcPoints(centerX, centerY, trR, startAngle, endAngle, segments)
		// 跳过第一个点（与上边缘终点重合）
		if len(arcPoints) > 1 {
			points = append(points, arcPoints[1:]...)
		}
	} else {
		points = append(points, [2]float32{x2, y1})
	}

	// 右边缘直线段
	rightEdgeEndY := y2 - brR
	points = append(points, [2]float32{x2, rightEdgeEndY})

	// 右下角圆角
	if brR > 0 {
		centerX, centerY := x2-brR, y2-brR
		startAngle := 0         // 0度
		endAngle := math.Pi / 2 // 90度
		arcPoints := generateArcPoints(centerX, centerY, brR, float64(startAngle), endAngle, segments)
		if len(arcPoints) > 1 {
			points = append(points, arcPoints[1:]...)
		}
	} else {
		points = append(points, [2]float32{x2, y2})
	}

	// 下边缘直线段
	bottomEdgeEndX := x1 + blR
	points = append(points, [2]float32{bottomEdgeEndX, y2})

	// 左下角圆角
	if blR > 0 {
		centerX, centerY := x1+blR, y2-blR
		startAngle := math.Pi / 2 // 90度
		endAngle := math.Pi       // 180度
		arcPoints := generateArcPoints(centerX, centerY, blR, startAngle, endAngle, segments)
		if len(arcPoints) > 1 {
			points = append(points, arcPoints[1:]...)
		}
	} else {
		points = append(points, [2]float32{x1, y2})
	}

	// 左边缘直线段（闭合到起点）
	leftEdgeEndY := y1 + tlR
	points = append(points, [2]float32{x1, leftEdgeEndY})

	return points
}

// generateArcPoints 生成圆弧上的点（修复类型转换错误）
func generateArcPoints(centerX, centerY, radius float32, startAngle, endAngle float64, segments int) [][2]float32 {
	var points [][2]float32
	if segments <= 0 {
		return points
	}

	// 角度计算使用float64确保精度
	totalAngle := endAngle - startAngle
	angleStep := totalAngle / float64(segments)

	for i := 0; i <= segments; i++ {
		// 角度计算全程使用float64
		currentAngle := startAngle + float64(i)*angleStep
		// 三角函数计算后显式转换为float32
		x := centerX + float32(math.Cos(currentAngle))*radius
		y := centerY + float32(math.Sin(currentAngle))*radius
		points = append(points, [2]float32{x, y})
	}

	return points
}
