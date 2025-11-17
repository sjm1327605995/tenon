package component

import (
	"bytes"
	"io"
	"os"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Image struct {
	src string
	Base[Image]
	load  bool
	image widget.Image
}

func (e *Image) Position(pos layout.Direction) *Image {
	e.image.Position = pos
	return e
}
func (e *Image) Scale(scale float32) *Image {
	e.image.Scale = scale
	return e
}

func (e *Image) Src(src string) *Image {
	e.src = src
	return e
}
func (e *Image) Fit(fit widget.Fit) *Image {
	e.image.Fit = fit
	return e
}

func NewImage() *Image {
	img := &Image{}
	img.Base = NewBase[Image](img)
	return img
}
func init() {
	// Register SVG format with image package to handle both XML declaration and direct SVG tag cases
	image.RegisterFormat("svg", "<svg", svgDecode, svgDecodeConfig)
}

func svgDecode(r io.Reader) (image.Image, error) {
	fc, err := canvas.ParseSVG(r)
	if err != nil {
		return nil, err
	}
	img := rasterizer.Draw(fc, canvas.DPMM(10.0), canvas.DefaultColorSpace)
	return img, nil
}

func svgDecodeConfig(r io.Reader) (image.Config, error) {
	fc, err := canvas.ParseSVG(r)
	if err != nil {
		return image.Config{}, err
	}
	w, h := fc.Size()
	return image.Config{
		Width:  int(w*canvas.DPMM(10.0).DPMM() + 0.5),
		Height: int(h*canvas.DPMM(10.0).DPMM() + 0.5),
	}, nil
}

func (e *Image) Layout(gtx layout.Context) layout.Dimensions {
	if !e.load {
		var img image.Image
		var err error
		data, err := os.ReadFile(e.src)
		if err != nil {
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}
		// Decode image (SVG is now supported through RegisterFormat)
		img, _, err = image.Decode(bytes.NewReader(data))
		if err != nil {
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}
		if img != nil {
			e.image.Src = paint.NewImageOp(img)
		}
		//e.Node.SetMeasureFunc(func(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
		//
		//})
		// Prevent reloading on every frame
		e.load = true
	}
	return e.image.Layout(gtx)
}
