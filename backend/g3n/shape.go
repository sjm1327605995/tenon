package g3n

import (
	"math"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

func minf(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

const vectorVertexShader = `#include <attributes>
uniform mat4 ModelMatrix;
out vec2 FragTexcoord;
void main() {
    FragTexcoord = VertexTexcoord;
    FragTexcoord.y = 1.0 - FragTexcoord.y;
    gl_Position = ModelMatrix * vec4(VertexPosition, 1.0);
}`

const vectorFragmentShader = `precision highp float;
uniform vec4 Shape;
uniform vec4 FromColor;
uniform vec4 ToColor;
uniform vec2 Direction;
uniform float Opacity;
in vec2 FragTexcoord;
out vec4 FragColor;
void main() {
    float w = Shape.x;
    float h = Shape.y;
    float radius = max(0.0, Shape.z);
    float stroke = Shape.w;
    float cut = max(0.0, -Shape.z);
    radius = min(radius, min(w, h) * 0.5);
    cut = min(cut, min(w, h) * 0.5);
    vec2 point = vec2(FragTexcoord.x * w, FragTexcoord.y * h);
    vec2 q = abs(point - vec2(w, h) * 0.5) - (vec2(w, h) * 0.5 - radius);
    float distance = length(max(q, vec2(0.0))) + min(max(q.x, q.y), 0.0) - radius;
    if (cut > 0.0) {
        distance = max(distance, (point.x + point.y - (w + h - cut)) / sqrt(2.0));
    }
    float coverage = clamp(0.5 - distance / max(fwidth(distance), 0.0001), 0.0, 1.0);
    if (stroke > 0.0) {
        float edge = min(-distance, distance + stroke);
        coverage = clamp(0.5 + edge / max(fwidth(distance), 0.0001), 0.0, 1.0);
    }
    if (coverage <= 0.0) {
        discard;
    }
    float halfSpan = (abs(w * Direction.x) + abs(h * Direction.y)) * 0.5;
    float t = 0.5 + dot(point - vec2(w, h) * 0.5, Direction) / max(2.0 * halfSpan, 0.0001);
    vec4 color = mix(FromColor, ToColor, clamp(t, 0.0, 1.0));
    FragColor = vec4(color.rgb, color.a * coverage * Opacity);
}`

type vectorMaterial struct {
	material.Material
	shape     math32.Vector4
	from, to  math32.Color4
	direction math32.Vector2
	opacity   float32
	uniShape  gls.Uniform
	uniFrom   gls.Uniform
	uniTo     gls.Uniform
	uniDir    gls.Uniform
	uniAlpha  gls.Uniform
}

func newVectorMaterial(w, h, r float32, from, to ui.Color, angle, stroke float32) *vectorMaterial {
	rad := float64(angle) * math.Pi / 180
	mat := &vectorMaterial{
		shape:     math32.Vector4{X: w, Y: h, Z: r, W: stroke},
		from:      *color4(from),
		to:        *color4(to),
		direction: math32.Vector2{X: float32(math.Cos(rad)), Y: float32(math.Sin(rad))},
		opacity:   1,
	}
	mat.Material.Init()
	mat.SetShader("tenon_vector")
	mat.SetShaderUnique(true)
	mat.SetUseLights(material.UseLightNone)
	mat.SetSide(material.SideDouble)
	mat.SetTransparent(true)
	mat.SetDepthTest(false)
	mat.SetDepthMask(false)
	mat.uniShape.Init("Shape")
	mat.uniFrom.Init("FromColor")
	mat.uniTo.Init("ToColor")
	mat.uniDir.Init("Direction")
	mat.uniAlpha.Init("Opacity")
	return mat
}

func (m *vectorMaterial) RenderSetup(gs *gls.GLS) {
	m.Material.RenderSetup(gs)
	gs.Uniform4f(m.uniShape.Location(gs), m.shape.X, m.shape.Y, m.shape.Z, m.shape.W)
	gs.Uniform4f(m.uniFrom.Location(gs), m.from.R, m.from.G, m.from.B, m.from.A)
	gs.Uniform4f(m.uniTo.Location(gs), m.to.R, m.to.G, m.to.B, m.to.A)
	gs.Uniform2f(m.uniDir.Location(gs), m.direction.X, m.direction.Y)
	gs.Uniform1f(m.uniAlpha.Location(gs), m.opacity)
}

type vectorPanel struct {
	*gui.Panel
	material *vectorMaterial
	uniModel gls.Uniform
	clipData
}

func newVectorPanel(x, y, w, h, r float32, from, to ui.Color, angle, stroke float32) *vectorPanel {
	panel := &vectorPanel{Panel: gui.NewPanel(w, h), material: newVectorMaterial(w, h, r, from, to, angle, stroke)}
	gr := graphic.NewGraphic(panel, panel.GetGeometry(), gls.TRIANGLES)
	gr.AddMaterial(panel, panel.material, 0, 0)
	panel.Panel.InitializeGraphic(w, h, gr)
	panel.uniModel.Init("ModelMatrix")
	panel.SetPosition(x, y)
	panel.SetColor4(color4(ui.White))
	return panel
}

func (p *vectorPanel) SetColor4(color *math32.Color4) *gui.Panel {
	p.material.opacity = color.A
	return p.Panel.SetColor4(color)
}

func (p *vectorPanel) RenderSetup(gs *gls.GLS, info *core.RenderInfo) {
	p.applyClip(gs)
	renderPanelTransform(gs, p.Panel, p.Panel, &p.uniModel, info)
}

func gradientPanel(x, y, w, h, r float32, from, to ui.Color, angle float32) *vectorPanel {
	return newVectorPanel(x, y, w, h, r, from, to, angle, 0)
}

func strokePanel(x, y, w, h, r, width float32, c ui.Color) *vectorPanel {
	return newVectorPanel(x, y, w, h, r, c, c, 0, width)
}
