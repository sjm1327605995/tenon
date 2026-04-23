package tenon

import (
	"github.com/sjm1327605995/tenon/internal/renderer"
	"github.com/sjm1327605995/tenon/pkg/core"
)

// 导出公共类型，方便用户直接使用 tenon.XXX
type (
	Component  = core.Component
	Widget     = core.Widget
	Host       = core.Host
	BaseWidget = core.BaseWidget
	BaseHost   = core.BaseHost
	Element    = core.Element
	LayoutBounds = core.LayoutBounds
	Event      = core.Event
	EventType  = core.EventType
	BorderRadius = core.BorderRadius
	PointerEvents = core.PointerEvents
	Context    = core.Context
	Ref        = core.Ref
	Theme      = core.Theme
)

// 导出便捷变量/函数
var (
	NewElement      = core.NewElement
	NewContext      = core.NewContext
	PointerEventsAuto  = core.PointerEventsAuto
	PointerEventsNone  = core.PointerEventsNone
	GetTheme        = core.GetTheme
	SetTheme        = core.SetTheme
	DefaultLightTheme = core.DefaultLightTheme
	DefaultDarkTheme  = core.DefaultDarkTheme
	DefaultAntTheme   = core.DefaultAntTheme
)

// Run 启动应用。root 可以是 Host 或 Widget。
func Run(root core.Component, width, height int) {
	renderer.Run(root, width, height)
}
