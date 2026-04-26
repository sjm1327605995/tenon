package tenon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/pkg/v2/debug"
)

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
	Animation    = core.Animation
	Tween        = core.Tween
)

var (
	GetTheme          = core.GetTheme
	SetTheme          = core.SetTheme
	DefaultLightTheme = core.DefaultLightTheme
	DefaultDarkTheme  = core.DefaultDarkTheme
	DefaultAntTheme   = core.DefaultAntTheme
	RegisterStyle     = core.RegisterStyle
	LogDebug          = core.LogDebug
	NewTween          = core.NewTween
	LerpFloat32       = core.LerpFloat32
)

func Run(root core.Widget, width, height int) {
	engine := core.NewEngine(root, width, height)
	engine.Mount()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(engine); err != nil {
		panic(err)
	}
}

func RunWithDebug(root core.Widget, width, height int, debugPort int) *debug.Debugger {
	engine := core.NewEngine(root, width, height)

	d := debug.NewDebugger(engine, debugPort)
	engine.SetDebugger(d)
	if err := d.Start(); err != nil {
		panic("failed to start debugger: " + err.Error())
	}

	engine.Mount()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(engine); err != nil {
		d.Stop()
		panic(err)
	}
	return d
}
