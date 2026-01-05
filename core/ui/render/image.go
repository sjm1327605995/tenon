package render

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/sjm1327605995/tenon/core/ui/style"
)

type ImageStyle struct {
	Src string `json:"src"`
	Fit style.Fit
}

func (i ImageStyle) ToRender() Render {
	data, err := os.ReadFile(i.Src)
	if err != nil {
		return nil
	}
	reader := bytes.NewReader(data)
	if bytes.HasPrefix(data, []byte("<svg")) {
		return NewSvg(reader, i)
	}
	img, _, err := image.Decode(reader)

	if err != nil {
		return &Image{
			Image: widget.Image{
				Src: paint.NewImageOp(img),
				Fit: widget.Fit(i.Fit),
			},
			style:   i,
			maxSize: img.Bounds().Max,
		}
	}
	return nil
}

type Image struct {
	widget.Image
	style   ImageStyle
	maxSize image.Point
	BaseRender
}

func (i *Image) DefaultSize() image.Point {
	return i.maxSize
}

func (i *Image) HasDefault() bool {
	return true
}

func (i *Image) Layout(ctx layout.Context) layout.Dimensions {
	return i.Image.Layout(ctx)
}
