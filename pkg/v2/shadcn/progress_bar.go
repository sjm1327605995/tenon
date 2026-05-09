package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// ProgressBar creates a linear progress indicator (0.0 ~ 1.0).
func ProgressBar(progress float32) ui.Widget {
	return progressBarWidget{Progress: progress}
}

// internal widget type
type progressBarWidget struct {
	ui.BaseWidget
	Progress float32
}

func (p progressBarWidget) CreateElement() ui.Element {
	return ui.NewRenderObjectElement(p)
}

// CreateRenderObject implements RenderObjectFactory.
func (p progressBarWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	theme := ui.GetTheme()
	r := render.NewRenderProgressBar()
	r.Progress = p.Progress
	r.TrackColor = theme.MutedColor
	r.FillColor = theme.PrimaryColor
	r.BarRadius = theme.BorderRadius / 2
	r.SetBorderRadius(theme.BorderRadius)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (p progressBarWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	r := ro.(*render.RenderProgressBar)
	r.SetProgress(p.Progress)
}
