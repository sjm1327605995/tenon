package tenon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// 导出公共类型，方便用户直接使用 tenon.XXX
type (
	Widget       = core.Widget
	BaseWidget   = core.BaseWidget
	Element      = core.Element
	BaseElement  = core.BaseElement
	LayoutBounds = core.LayoutBounds
	Event        = core.Event
	EventType    = core.EventType
	BorderRadius = core.BorderRadius
	Theme        = core.Theme
	Engine       = core.Engine
	Animation = core.Animation
	Tween     = core.Tween
)

// 导出便捷变量/函数
var (
	GetTheme          = core.GetTheme
	SetTheme          = core.SetTheme
	DefaultLightTheme = core.DefaultLightTheme
	DefaultDarkTheme  = core.DefaultDarkTheme
	DefaultAntTheme   = core.DefaultAntTheme
	RegisterStyle     = core.RegisterStyle
	LogDebug          = core.LogDebug
	NewTween    = core.NewTween
	LerpFloat32 = core.LerpFloat32
)

// Run 启动 Tenon 应用。
func Run(root core.Widget, width, height int) {
	engine := core.NewEngine(root, width, height)
	engine.Mount()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(engine); err != nil {
		panic(err)
	}
}
