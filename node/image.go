package node

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type Image struct {
	Src image.Image
	*View
}

func (i *Image) Measure() {
	bounds := i.Src.Bounds()
	originalWidth := float32(bounds.Dx())
	originalHeight := float32(bounds.Dy())
	w, h := i.Node.StyleGetWidth(), i.Node.StyleGetHeight()

	// 如果未指定宽度和高度，则使用原始尺寸
	if w == 0 && h == 0 {
		i.Node.StyleSetWidth(originalWidth)
		i.Node.StyleSetHeight(originalHeight)
	} else if w == 0 { // 仅指定高度，按比例计算宽度
		scale := h / originalHeight
		i.Node.StyleSetWidth(originalWidth * scale)
	} else if h == 0 { // 仅指定宽度，按比例计算高度
		scale := w / originalWidth
		i.Node.StyleSetHeight(originalHeight * scale)
	}
}

func (i *Image) OnDraw(r Renderer) {
	_ = r.Image(i.X, i.Y, i.Node.Node, i.Src)
	i.Node.OnDraw(r)
}
func NewImage(image image.Image) *Image {
	return &Image{
		View: NewView(),
		Src:  image,
	}

}
func NewFileImage(file string) (image.Image, error) {
	f, err := os.Open("gopher.png")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}
