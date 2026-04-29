package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// ProgressBar is a linear progress indicator.
type ProgressBar struct {
	core.BaseElement
	progress     float32 // 0.0 ~ 1.0
	trackColor   color.Color
	fillColor    color.Color
	borderRadius float32
}

// NewProgressBar creates a progress bar.
func NewProgressBar() *ProgressBar {
	theme := core.GetTheme()
	pb := &ProgressBar{
		progress:     0,
		trackColor:   theme.ProgressBarTrackColor,
		fillColor:    theme.ProgressBarFillColor,
		borderRadius: theme.BorderRadius / 2,
	}
	pb.Init(pb)
	pb.SetHeight(8)
	return pb
}

// ElementType returns type identifier.
func (pb *ProgressBar) ElementType() string { return "ProgressBar" }

// Draw renders the track and fill.
func (pb *ProgressBar) Draw(screen *ebiten.Image) {
	bounds := pb.GetBounds()

	// Track
	pb.drawRoundedRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, pb.trackColor)

	// Fill
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

// SetProgress sets the progress value [0, 1].
func (pb *ProgressBar) SetProgress(progress float32) *ProgressBar {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	pb.progress = progress
	pb.Mark(core.FlagNeedDraw)
	return pb
}

// GetProgress returns the current progress value.
func (pb *ProgressBar) GetProgress() float32 {
	return pb.progress
}

// SetTrackColor sets the track background color.
func (pb *ProgressBar) SetTrackColor(clr color.Color) *ProgressBar {
	pb.trackColor = clr
	pb.Mark(core.FlagNeedDraw)
	return pb
}

// SetFillColor sets the fill color.
func (pb *ProgressBar) SetFillColor(clr color.Color) *ProgressBar {
	pb.fillColor = clr
	pb.Mark(core.FlagNeedDraw)
	return pb
}

// SetBorderRadius sets the corner radius.
func (pb *ProgressBar) SetBorderRadius(radius float32) *ProgressBar {
	pb.borderRadius = radius
	pb.Mark(core.FlagNeedDraw)
	return pb
}

// DebugProps returns visual properties for debugger preview.
func (pb *ProgressBar) DebugProps() map[string]interface{} {
	props := make(map[string]interface{})
	if pb.trackColor != nil {
		props["backgroundColor"] = colorToCSS(pb.trackColor)
	}
	if pb.fillColor != nil {
		props["fillColor"] = colorToCSS(pb.fillColor)
	}
	props["progress"] = pb.progress
	if pb.borderRadius > 0 {
		props["borderRadius"] = pb.borderRadius
	}
	return props
}

// SyncFrom 同步新 ProgressBar 的属性到当前 Element（声明式重建）。
func (pb *ProgressBar) SyncFrom(src core.Element) {
	other, ok := src.(*ProgressBar)
	if !ok {
		return
	}
	sb := &SyncBuilder{}
	syncField(sb, &pb.progress, other.progress)
	syncColor(sb, &pb.trackColor, other.trackColor)
	syncColor(sb, &pb.fillColor, other.fillColor)
	syncField(sb, &pb.borderRadius, other.borderRadius)
	sb.MarkDraw(pb)
}

