package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// ProgressBar creates a linear progress indicator (0.0 ~ 1.0).
func ProgressBar(progress float32) engine.Widget {
	return progressBarWidget{Progress: progress}
}

// internal widget type
type progressBarWidget struct {
	engine.BaseWidget
	Progress float32
}

func (p progressBarWidget) CreateElement() engine.Element {
	return engine.NewRenderObjectElement(p)
}

// CreateRenderObject implements RenderObjectFactory.
func (p progressBarWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	theme := engine.GetTheme()
	r := render.NewRenderProgressBar()
	r.Progress = p.Progress
	r.TrackColor = theme.MutedColor
	r.FillColor = theme.PrimaryColor
	r.BarRadius = theme.BorderRadius / 2
	r.SetBorderRadius(theme.BorderRadius)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (p progressBarWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	r := ro.(*render.RenderProgressBar)
	r.SetProgress(p.Progress)
}
