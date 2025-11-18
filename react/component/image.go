package component

import (
	"bytes"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

type Image struct {
	src    string
	origin image.Image
	format string // Tracks image format (e.g., "svg", "png")
	data   []byte // Stores original image data for SVG re-rendering
	Base[Image]
	fit      widget.Fit
	position layout.Direction
	scale    float32
}

func (e *Image) Position(pos layout.Direction) *Image {
	e.position = pos
	return e
}
func (e *Image) Scale(scale float32) *Image {
	e.scale = scale
	return e
}

func (e *Image) Src(src string) *Image {
	e.src = src
	return e
}
func (e *Image) Fit(fit widget.Fit) *Image {
	e.fit = fit
	return e
}

func NewImage() *Image {
	img := &Image{}
	img.Base = NewBase[Image](img)
	return img
}

type imageGio struct {
	widget.Image
}

func (i *imageGio) Layout(gtx layout.Context) layout.Dimensions {
	return i.Image.Layout(gtx)
}

func init() {
	// Register SVG format with image package to handle both XML declaration and direct SVG tag cases
	image.RegisterFormat("svg", "<svg", svgDecode, svgDecodeConfig)
}

const ( //96.0 / 25.4
	SvgDPMM = 10
)

func svgDecode(r io.Reader) (image.Image, error) {
	fc, err := canvas.ParseSVG(r)
	if err != nil {
		return nil, err
	}
	img := rasterizer.Draw(fc, canvas.DPMM(SvgDPMM), canvas.DefaultColorSpace)
	return cropTransparent(img), nil
}

// cropTransparent 裁剪图像中多余的透明区域，保留最小有效区域
func cropTransparent(img image.Image) image.Image {
	bounds := img.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	// 遍历所有像素，寻找非透明像素的边界
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 获取像素的 alpha 通道（判断是否透明）
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 { // 非透明像素（alpha > 0）
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	// 如果全是透明像素，返回原图
	if minX > maxX || minY > maxY {
		return img
	}
	// 计算裁剪区域（+1 是因为 image.Rect 是左闭右开区间）
	cropRect := image.Rect(minX, minY, maxX+1, maxY+1)

	// 创建新图像并复制有效区域
	cropped := image.NewRGBA(cropRect)
	draw.Draw(cropped, cropped.Bounds(), img, cropRect.Min, draw.Src)

	return cropped
}

func svgDecodeConfig(r io.Reader) (image.Config, error) {
	fc, err := canvas.ParseSVG(r)
	if err != nil {
		return image.Config{}, err
	}
	w, h := fc.Size()
	return image.Config{
		Width:  int(w * canvas.DPMM(SvgDPMM).DPMM()),
		Height: int(h * canvas.DPMM(SvgDPMM).DPMM()),
	}, nil
}
func (e *Image) Update(gtx layout.Context) {
	if e.origin == nil {
		var img image.Image
		var err error
		data, err := os.ReadFile(e.src)
		if err != nil {
			e.gio = &RectGio{}
			return
		}
		//TODO 如果是svg 当size变化是不通过image 伸缩 而是通过canvas重绘图片
		img, format, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			e.gio = &RectGio{}
			return
		}
		e.origin = img
		e.format = format
		e.data = data // Save data for potential SVG re-rendering
	}

	// 获取图片原始尺寸（像素）
	imgWidth := float32(e.origin.Bounds().Dx())
	imgHeight := float32(e.origin.Bounds().Dy())
	if e.Node.StyleGetWidth() == 0 && e.Node.StyleGetHeight() == 0 {
		e.Node.StyleSetWidth(imgWidth)
		e.Node.StyleSetHeight(imgHeight)
	}
	e.Node.SetMeasureFunc(func(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {

		// 计算最终宽度
		finalWidth := calculateDimension(imgWidth, width, widthMode)
		// 计算最终高度
		finalHeight := calculateDimension(imgHeight, height, heightMode)

		// 如果只约束了一个维度，根据原始比例调整另一个维度（保持宽高比）
		if widthMode == yoga.MeasureModeExactly && heightMode != yoga.MeasureModeExactly {
			// 宽度固定，按比例计算高度
			finalHeight = (finalWidth / imgWidth) * imgHeight
		} else if heightMode == yoga.MeasureModeExactly && widthMode != yoga.MeasureModeExactly {
			// 高度固定，按比例计算宽度
			finalWidth = (finalHeight / imgHeight) * imgWidth
		}
		// 确保尺寸不为负（Yoga 不允许负尺寸）
		finalWidth = max(finalWidth, 0)
		finalHeight = max(finalHeight, 0)
		return yoga.Size{Width: finalWidth, Height: finalHeight}
	})

	e.gio = &imageGio{
		Image: widget.Image{
			Src:      paint.NewImageOp(e.origin),
			Fit:      e.fit,
			Scale:    e.scale,
			Position: e.position,
		},
	}
}

// 辅助函数：根据原始尺寸、约束尺寸和测量模式计算最终维度
func calculateDimension(original, constraint float32, mode yoga.MeasureMode) float32 {
	switch mode {
	case yoga.MeasureModeExactly:
		// 强制使用约束尺寸
		return constraint
	case yoga.MeasureModeAtMost:
		// 不超过约束尺寸，取原始尺寸和约束尺寸的最小值
		return min(original, constraint)
	case yoga.MeasureModeUndefined:
		// 无约束，使用原始尺寸
		return original
	default:
		return original
	}
}
