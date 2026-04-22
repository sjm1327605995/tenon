package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ProgressBar 是进度条组件。
type ProgressBar struct {
	core.BaseHost
	progress      float32 // 0.0 ~ 1.0
	trackColor    color.Color
	fillColor     color.Color
	borderRadius  float32
}

// NewProgressBar 创建一个进度条。
func NewProgressBar() *ProgressBar {
	pb := &ProgressBar{
		progress:     0,
		trackColor:   color.RGBA{R: 224, G: 224, B: 224, A: 255},
		fillColor:    color.RGBA{R: 0, G: 123, B: 255, A: 255},
		borderRadius: 4,
	}
	pb.Init(pb)
	pb.GetElement().Yoga.StyleSetHeight(8)
	return pb
}

// Draw 绘制进度条背景和填充。
func (pb *ProgressBar) Draw(screen *ebiten.Image) {
	el := pb.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := pb.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 绘制轨道
	pb.drawRoundedRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, pb.trackColor)

	// 绘制填充
	if pb.progress > 0 {
		fillW := bounds.Width * pb.progress
		if fillW > bounds.Width {
			fillW = bounds.Width
		}
		pb.drawRoundedRect(screen, bounds.X, bounds.Y, fillW, bounds.Height, pb.fillColor)
	}
}

func (pb *ProgressBar) drawRoundedRect(screen *ebiten.Image, x, y, w, h float32, clr color.Color) {
	if pb.borderRadius > 0 {
		var path vector.Path
		r := pb.borderRadius
		if r > w/2 {
			r = w / 2
		}
		if r > h/2 {
			r = h / 2
		}
		buildRoundedRectPath(&path, x, y, w, h, core.BorderRadius{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r})
		op := &vector.DrawPathOptions{}
		op.ColorScale.ScaleWithColor(clr)
		op.AntiAlias = true
		vector.FillPath(screen, &path, &vector.FillOptions{}, op)
	} else {
		vector.FillRect(screen, x, y, w, h, clr, false)
	}
}

// ==================== 链式 API ====================

func (pb *ProgressBar) SetProgress(progress float32) *ProgressBar {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	pb.progress = progress
	return pb
}
func (pb *ProgressBar) SetTrackColor(clr color.Color) *ProgressBar {
	pb.trackColor = clr
	return pb
}
func (pb *ProgressBar) SetFillColor(clr color.Color) *ProgressBar {
	pb.fillColor = clr
	return pb
}
func (pb *ProgressBar) SetWidth(width float32) *ProgressBar {
	pb.GetElement().Yoga.StyleSetWidth(width)
	return pb
}
func (pb *ProgressBar) SetHeight(height float32) *ProgressBar {
	pb.GetElement().Yoga.StyleSetHeight(height)
	return pb
}
func (pb *ProgressBar) SetBorderRadius(radius float32) *ProgressBar {
	pb.borderRadius = radius
	return pb
}
func (pb *ProgressBar) SetMargin(edge yoga.Edge, value float32) *ProgressBar {
	pb.GetElement().Yoga.StyleSetMargin(edge, value)
	return pb
}
func (pb *ProgressBar) SetFlexGrow(grow float32) *ProgressBar {
	pb.GetElement().Yoga.StyleSetFlexGrow(grow)
	return pb
}

// SyncFrom 同步进度条属性。
func (pb *ProgressBar) SyncFrom(other core.Host) {
	if o, ok := other.(*ProgressBar); ok {
		pb.progress = o.progress
		pb.trackColor = o.trackColor
		pb.fillColor = o.fillColor
		pb.borderRadius = o.borderRadius
	}
}
