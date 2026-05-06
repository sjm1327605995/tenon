// Package render 提供渲染树节点。
//
// 已迁移至 internal/render，此包仅为兼容性保留。
// 请直接使用 tenon 包的统一入口。
package render

import internal "github.com/sjm1327605995/tenon/internal/render"

// Re-export 所有公共类型和函数。
type (
	RenderObject  = internal.RenderObject
	BaseRenderObject = internal.BaseRenderObject
	RenderBox     = internal.RenderBox
	RenderFlex    = internal.RenderFlex
	RenderText    = internal.RenderText
	RenderImage   = internal.RenderImage
	RenderButton  = internal.RenderButton
	RenderScroll  = internal.RenderScroll
	RenderStack   = internal.RenderStack
	RenderCheckbox = internal.RenderCheckbox
	RenderRadio   = internal.RenderRadio
	RenderSwitch  = internal.RenderSwitch
	RenderSlider  = internal.RenderSlider
	RenderProgressBar = internal.RenderProgressBar
	RenderEditableText = internal.RenderEditableText
	ParticleSystemRenderObject = internal.ParticleSystemRenderObject
	PipelineOwner = internal.PipelineOwner
	Bounds        = internal.Bounds
	Offset        = internal.Offset
	Size          = internal.Size
	BoxConstraints = internal.BoxConstraints
	Color         = internal.Color
	BorderRadius  = internal.BorderRadius
	Transform     = internal.Transform
	ObjectFit     = internal.ObjectFit
	Particle      = internal.Particle
	ParticleEmitter = internal.ParticleEmitter
	PointEmitter  = internal.PointEmitter
	RectEmitter   = internal.RectEmitter
	CircleEmitter = internal.CircleEmitter
	ParticleBehavior = internal.ParticleBehavior
)

var (
	NewRenderBox     = internal.NewRenderBox
	NewRenderFlex    = internal.NewRenderFlex
	NewRenderText    = internal.NewRenderText
	NewRenderImage   = internal.NewRenderImage
	NewRenderButton  = internal.NewRenderButton
	NewRenderScroll  = internal.NewRenderScroll
	NewRenderStack   = internal.NewRenderStack
	NewRenderCheckbox = internal.NewRenderCheckbox
	NewRenderRadio   = internal.NewRenderRadio
	NewRenderSwitch  = internal.NewRenderSwitch
	NewRenderSlider  = internal.NewRenderSlider
	NewRenderProgressBar = internal.NewRenderProgressBar
	NewRenderEditableText = internal.NewRenderEditableText
	NewParticleSystemRenderObject = internal.NewParticleSystemRenderObject
	NewPipelineOwner = internal.NewPipelineOwner
	NewColor         = internal.NewColor
	NewColorFrom     = internal.NewColorFrom
	UniformBorderRadius = internal.UniformBorderRadius
	BuildRoundedRectPath = internal.BuildRoundedRectPath
	DrawRoundedRectFill = internal.DrawRoundedRectFill
	DrawRoundedRectStroke = internal.DrawRoundedRectStroke
	DrawRect         = internal.DrawRect
	DrawBorder       = internal.DrawBorder
	DrawShadow       = internal.DrawShadow
	DrawFilledCircle = internal.DrawFilledCircle
	DrawFilledCirclePath = internal.DrawFilledCirclePath
	DrawStrokedCirclePath = internal.DrawStrokedCirclePath
	StrokeCircle     = internal.StrokeCircle
	SubImage         = internal.SubImage
	Darken           = internal.Darken
	Lighten          = internal.Lighten
	DefaultTransform = internal.DefaultTransform
	BuildTransformGeoM = internal.BuildTransformGeoM
	ApplyColorScaleAlpha = internal.ApplyColorScaleAlpha
	ColorEquals      = internal.ColorEquals
	ColorPtrEquals   = internal.ColorPtrEquals
)

const (
	ObjectFitCover     = internal.ObjectFitCover
	ObjectFitContain   = internal.ObjectFitContain
	ObjectFitFill      = internal.ObjectFitFill
	ObjectFitNone      = internal.ObjectFitNone
	ObjectFitScaleDown = internal.ObjectFitScaleDown
)

// ButtonState 常量。
const (
	ButtonStateNormal   = internal.ButtonStateNormal
	ButtonStateHover    = internal.ButtonStateHover
	ButtonStatePressed  = internal.ButtonStatePressed
	ButtonStateDisabled = internal.ButtonStateDisabled
)

// ButtonState 类型。
type ButtonState = internal.ButtonState
