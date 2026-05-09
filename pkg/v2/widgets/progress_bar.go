package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// GameProgressBarWidget 是游戏风格的进度条，支持 NinePatch 图片背景和前景。
type GameProgressBarWidget struct {
	ui.BaseWidget
	progress   float32
	trackImage *ebiten.Image
	fillImage  *ebiten.Image
	slice      render.BorderSlice
	width      float32
	height     float32
}

// GameProgressBar 创建游戏风格进度条。
// track 为背景图片，fill 为填充图片；slice 为可选的 NinePatch 切图参数。
func GameProgressBar(progress float32, track, fill *ebiten.Image) GameProgressBarWidget {
	return GameProgressBarWidget{progress: progress, trackImage: track, fillImage: fill}
}

func (p GameProgressBarWidget) Slice(s render.BorderSlice) GameProgressBarWidget {
	p.slice = s
	return p
}

func (p GameProgressBarWidget) W(v float32) GameProgressBarWidget {
	p.width = v
	return p
}

func (p GameProgressBarWidget) H(v float32) GameProgressBarWidget {
	p.height = v
	return p
}

func (p GameProgressBarWidget) CreateElement() ui.Element {
	return ui.NewRenderObjectElement(p)
}

func (p GameProgressBarWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	r := render.NewRenderProgressBar()
	r.SetProgress(p.progress)
	r.TrackImage = p.trackImage
	r.FillImage = p.fillImage
	r.Slice = p.slice
	if p.width > 0 {
		r.StyleSetWidth(p.width)
	}
	if p.height > 0 {
		r.StyleSetHeight(p.height)
	}
	return r
}

func (p GameProgressBarWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	r := ro.(*render.RenderProgressBar)
	old := oldWidget.(GameProgressBarWidget)
	if old.progress != p.progress {
		r.SetProgress(p.progress)
	}
	if old.trackImage != p.trackImage || old.fillImage != p.fillImage || old.slice != p.slice {
		r.SetImages(p.trackImage, p.fillImage, p.slice)
	}
	if old.width != p.width {
		if p.width > 0 {
			r.StyleSetWidth(p.width)
		} else {
			r.StyleSetWidthAuto()
		}
	}
	if old.height != p.height {
		if p.height > 0 {
			r.StyleSetHeight(p.height)
		} else {
			r.StyleSetHeightAuto()
		}
	}
}
