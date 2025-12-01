package elements

import (
	"bytes"
	"image"
	"image/draw"
	"io"
	"os"

	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/components"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

type Image struct {
	*components.Image
}

func (i *Image) Render() api.Node {
	return i
}

func NewImage() *Image {
	return &Image{Image: components.NewImage()}
}
func (i *Image) SetExtendedStyle(style styles.IExtendedStyle) {
	//TODO implement me
	panic("implement me")
}

func (i *Image) GetChildrenCount() int {
	return 0
}
func (i *Image) Source(path string) *Image {
	data, err := os.ReadFile(path)
	if err != nil {
		return i
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return i
	}
	imgBounds := img.Bounds()
	imgWidth := float32(imgBounds.Dx())
	imgHeight := float32(imgBounds.Dy())
	userWidth := i.Yoga().StyleGetWidth()
	userHeight := i.Yoga().StyleGetHeight()
	// 最终显示尺寸（默认用原始尺寸，后续根据用户设置覆盖）
	finalWidth := imgWidth
	finalHeight := imgHeight

	// 尺寸逻辑处理 + 计算最终显示尺寸
	if userWidth == 0 && userHeight == 0 {
		// 场景1：未设置宽高 → 使用原图尺寸（缩放率=1）
		i.Yoga().StyleSetWidth(imgWidth)
		i.Yoga().StyleSetHeight(imgHeight)
	} else if userWidth > 0 && userHeight == 0 {
		// 场景2：只设宽度 → 按比例算高度
		aspectRatio := imgHeight / imgWidth
		finalHeight = userWidth * aspectRatio
		finalWidth = userWidth // 用户指定的宽度作为最终宽度
		i.Yoga().StyleSetHeight(finalHeight)
	} else if userWidth == 0 && userHeight > 0 {
		// 场景3：只设高度 → 按比例算宽度
		aspectRatio := imgWidth / imgHeight
		finalWidth = userHeight * aspectRatio
		finalHeight = userHeight // 用户指定的高度作为最终高度
		i.Yoga().StyleSetWidth(finalWidth)
	} else {
		// 场景4：用户已设置宽高 → 直接使用用户尺寸（后续计算缩放率）
		finalWidth = userWidth
		finalHeight = userHeight
	}

	// 核心：计算缩放率（取宽/高缩放率的最小值，避免图片拉伸）
	// 缩放率 = 最终显示尺寸 / 原始尺寸
	scaleX := finalWidth / imgWidth     // 宽度缩放比例
	scaleY := finalHeight / imgHeight   // 高度缩放比例
	i.Image.Scale = min(scaleX, scaleY) // 取最小值（保持图片比例，不拉伸）
	i.Origin = img
	return i
}
func (i *Image) GetChildAt(index int) api.Element {
	return nil
}

func (i *Image) Rendering(renderer api.Renderer) {
	renderer.DrawImage(i.Image)
}

func (i *Image) GetChildren() []api.Element {
	return nil
}

func (i *Image) SetStyle(style *styles.Style) {
	style.Apply(i)
}
func (i *Image) Style(option *styles.Style) *Image {
	i.SetStyle(option)
	return i
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
func init() {
	// Register SVG format with image package to handle both XML declaration and direct SVG tag cases
	image.RegisterFormat("svg", "<svg", svgDecode, svgDecodeConfig)
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
