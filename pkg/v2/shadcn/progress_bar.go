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
	e := &progressBarElement{}
	e.RenderObjectElement.BaseElement.Init(e, p)
	return e
}

type progressBarElement struct {
	ui.RenderObjectElement
}

func (e *progressBarElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(progressBarWidget)
	theme := ui.GetTheme()
	r := render.NewRenderProgressBar()
	r.Progress = w.Progress
	r.TrackColor = theme.MutedColor
	r.FillColor = theme.PrimaryColor
	r.BarRadius = theme.BorderRadius / 2
	r.SetBorderRadius(theme.BorderRadius)
	return r
}

func (e *progressBarElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(progressBarWidget)
	r := e.GetRenderObject().(*render.RenderProgressBar)
	r.SetProgress(w.Progress)
}
