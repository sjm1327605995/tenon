package g3n

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type transformedPanel struct {
	gui.IPanel
	uniModel gls.Uniform
	clipData
}

func withTransforms(panel gui.IPanel) *transformedPanel {
	wrapped := &transformedPanel{IPanel: panel}
	wrapped.uniModel.Init("ModelMatrix")
	panel.GetGraphic().SetIGraphic(wrapped)
	return wrapped
}

func (p *transformedPanel) SetSize(width, height float32) {
	p.IPanel.GetPanel().SetSize(width, height)
}

func (p *transformedPanel) SetRotationZ(angle float32) {
	p.IPanel.GetPanel().SetRotationZ(angle)
}

func (p *transformedPanel) SetColor4(color *math32.Color4) *gui.Panel {
	return p.IPanel.GetPanel().SetColor4(color)
}

func (p *transformedPanel) Color4() math32.Color4 { return p.IPanel.GetPanel().Color4() }

func (p *transformedPanel) Material() *material.Material { return p.IPanel.GetPanel().Material() }

func (p *transformedPanel) RenderSetup(gs *gls.GLS, info *core.RenderInfo) {
	p.applyClip(gs)
	renderPanelTransform(gs, p.IPanel, p.IPanel.GetPanel(), &p.uniModel, info)
}

func renderPanelTransform(gs *gls.GLS, graphic gui.IPanel, panel *gui.Panel, uniform *gls.Uniform, info *core.RenderInfo) {
	graphic.RenderSetup(gs, info)
	angle := panel.Rotation().Z
	if angle == 0 {
		return
	}
	var model math32.Matrix4
	panel.SetModelMatrix(gs, &model)
	center, rotation, back := math32.Matrix4{}, math32.Matrix4{}, math32.Matrix4{}
	center.MakeTranslation(0.5, -0.5, 0)
	rotation.MakeRotationZ(angle)
	back.MakeTranslation(-0.5, 0.5, 0)
	center.Multiply(&rotation).Multiply(&back)
	model.Multiply(&center)
	gs.UniformMatrix4fv(uniform.Location(gs), 1, false, &model[0])
}
