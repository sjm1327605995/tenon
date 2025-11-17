package component

import (
	"bytes"
	"image"
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
	src string
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
func (e *Image) Update(gtx layout.Context) {
	var img image.Image
	var err error
	data, err := os.ReadFile(e.src)
	if err != nil {
		e.gio = &RectGio{}
		return
	}
	// Decode image (SVG is now supported through RegisterFormat)
	img, _, err = image.Decode(bytes.NewReader(data))
	if err != nil {
		e.gio = &RectGio{}
		return
	}
	e.gio = &imageGio{
		Image: widget.Image{
			Src:      paint.NewImageOp(img),
			Fit:      e.fit,
			Scale:    e.scale,
			Position: e.position,
		},
	}
}
