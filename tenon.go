package tenon

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/shadcn"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// 核心类型别名，用户可直接使用。
type (
	Widget            = ui.Widget
	Element           = ui.Element
	StatefulWidget    = ui.StatefulWidget
	State             = ui.State
	BaseState         = ui.BaseState
	BaseWidget        = ui.BaseWidget
	BuildContext      = ui.BuildContext
	Builder           = ui.Builder
	StatefulBuilder   = ui.StatefulBuilder
	GlobalKey         = ui.GlobalKey
	Engine            = ui.Engine
	Theme             = ui.Theme
	EdgeInsets        = ui.EdgeInsets
	Color             = render.Color
	FlexDirection     = ui.FlexDirection
	Justify           = ui.Justify
	Align             = ui.Align
	Wrap              = ui.Wrap
	PositionType      = ui.PositionType
	Display           = ui.Display
	Overflow          = ui.Overflow
	Edge              = ui.Edge
	Gutter            = ui.Gutter
	AnimatedContainer = widgets.AnimatedContainer
	SelectOption      = widgets.SelectOption
)

// 核心函数。
var (
	NewEngine             = ui.NewEngine
	NewStatefulElement    = ui.NewStatefulElement
	Rebuild               = ui.RebuildDefault
	SetTheme              = ui.SetTheme
	GetTheme              = ui.GetTheme
	ThemeOf               = ui.ThemeOf
	NewBuilder            = ui.NewBuilder
	NewStatefulBuilder    = ui.NewStatefulBuilder
	NewGlobalKey          = ui.NewGlobalKey
	DefaultLightTheme     = ui.DefaultLightTheme
	DefaultDarkTheme      = ui.DefaultDarkTheme
	EdgeInsetsAll         = ui.EdgeInsetsAll
	EdgeInsetsSymmetric   = ui.EdgeInsetsSymmetric
	EdgeInsetsOnly        = ui.EdgeInsetsOnly
	NewColor              = render.NewColor
	NewColorFrom          = render.NewColorFrom
)

// 动画类型别名。
type (
	AnimationController = ui.AnimationController
	AnimationStatus     = ui.AnimationStatus
	Tween[T any]        = ui.Tween[T]
	Animation[T any]    = ui.Animation[T]
	Curve               = ui.Curve
	LinearCurve         = ui.LinearCurve
	EaseInOutCurve      = ui.EaseInOutCurve
)

// 动画常量。
const (
	AnimationDismissed = ui.AnimationDismissed
	AnimationForward   = ui.AnimationForward
	AnimationReverse   = ui.AnimationReverse
	AnimationCompleted = ui.AnimationCompleted
)

// Widget 构造函数，SwiftUI-like API。
var (
	Text               = widgets.Text
	Row                = widgets.Row
	Column             = widgets.Column
	Container          = widgets.Container
	Button             = widgets.Button
	Image              = widgets.Image
	Stack              = widgets.Stack
	Positioned         = widgets.Positioned
	Scroll             = widgets.Scroll
	TextField          = widgets.TextField
	EditableText       = widgets.EditableText
	NewAnimatedContainer = widgets.NewAnimatedContainer
	Select               = widgets.Select
	MultiSelect          = widgets.MultiSelect
)

// Shadcn 组件。
var (
	Badge       = shadcn.Badge
	DotBadge    = shadcn.DotBadge
	CountBadge  = shadcn.CountBadge
	Separator   = shadcn.ShadcnSeparator
	Avatar      = shadcn.Avatar
	Card        = shadcn.ShadcnCard
	Alert       = shadcn.ShadcnAlert
	Label       = shadcn.Label
	Skeleton    = shadcn.Skeleton
	ProgressBar = shadcn.ProgressBar
	Checkbox    = shadcn.Checkbox
	Switch      = shadcn.Switch
	Radio       = shadcn.Radio
	Slider      = shadcn.Slider
	Tabs        = shadcn.ShadcnTabs
	Breadcrumb  = shadcn.ShadcnBreadcrumb
	Pagination  = shadcn.ShadcnPagination
	Table       = shadcn.ShadcnTable
	Accordion   = shadcn.ShadcnAccordion
	Toggle      = shadcn.ShadcnToggle
	Textarea    = shadcn.ShadcnTextarea
	Calendar    = shadcn.ShadcnCalendar
)

// Shadcn 类型别名。
type (
	TabItem         = shadcn.TabItem
	BreadcrumbItem  = shadcn.BreadcrumbItem
	AccordionPane   = shadcn.AccordionPane
)

// Yoga/Flex 常量快捷导出。
const (
	FlexDirectionColumn  = ui.FlexDirectionColumn
	FlexDirectionRow     = ui.FlexDirectionRow
	JustifyFlexStart     = ui.JustifyFlexStart
	JustifyCenter        = ui.JustifyCenter
	JustifyFlexEnd       = ui.JustifyFlexEnd
	JustifySpaceBetween  = ui.JustifySpaceBetween
	JustifySpaceAround   = ui.JustifySpaceAround
	JustifySpaceEvenly   = ui.JustifySpaceEvenly
	AlignAuto            = ui.AlignAuto
	AlignFlexStart       = ui.AlignFlexStart
	AlignCenter          = ui.AlignCenter
	AlignFlexEnd         = ui.AlignFlexEnd
	AlignStretch         = ui.AlignStretch
	AlignBaseline        = ui.AlignBaseline
	WrapNoWrap           = ui.WrapNoWrap
	WrapWrap             = ui.WrapWrap
	PositionTypeRelative = ui.PositionTypeRelative
	PositionTypeAbsolute = ui.PositionTypeAbsolute
	DisplayFlex          = ui.DisplayFlex
	DisplayNone          = ui.DisplayNone
	OverflowVisible      = ui.OverflowVisible
	OverflowHidden       = ui.OverflowHidden
	OverflowScroll       = ui.OverflowScroll
	EdgeLeft             = ui.EdgeLeft
	EdgeTop              = ui.EdgeTop
	EdgeRight            = ui.EdgeRight
	EdgeBottom           = ui.EdgeBottom
	EdgeStart            = ui.EdgeStart
	EdgeEnd              = ui.EdgeEnd
	EdgeHorizontal       = ui.EdgeHorizontal
	EdgeVertical         = ui.EdgeVertical
	EdgeAll              = ui.EdgeAll
	GutterColumn         = ui.GutterColumn
	GutterRow            = ui.GutterRow
	GutterAll            = ui.GutterAll
)

// Run 启动 ebiten 窗口并运行应用。
// 根 Widget 会自动包裹在 ThemeProvider 中。
func Run(buildFunc ui.BuildFunc, width, height int) {
	wrappedBuildFunc := func() ui.Widget {
		return ui.ThemeProvider(ui.GetTheme(), buildFunc())
	}
	engine := ui.NewEngine(wrappedBuildFunc, width, height)
	engine.Run()
}

// Button 变体
const (
	ButtonDefault     = widgets.ButtonDefault
	ButtonSecondary   = widgets.ButtonSecondary
	ButtonOutline     = widgets.ButtonOutline
	ButtonGhost       = widgets.ButtonGhost
	ButtonDestructive = widgets.ButtonDestructive
	ButtonLink        = widgets.ButtonLink
)

// Badge 变体
const (
	BadgeDefault     = shadcn.BadgeDefault
	BadgeSecondary   = shadcn.BadgeSecondary
	BadgeOutline     = shadcn.BadgeOutline
	BadgeDestructive = shadcn.BadgeDestructive
)

// Alert 变体
const (
	AlertDefault     = shadcn.AlertDefault
	AlertDestructive = shadcn.AlertDestructive
)

// Separator 方向
const (
	SeparatorHorizontal = shadcn.SeparatorHorizontal
	SeparatorVertical   = shadcn.SeparatorVertical
)

// ObjectFit 快捷导出
const (
	ObjectFitCover     = render.ObjectFitCover
	ObjectFitContain   = render.ObjectFitContain
	ObjectFitFill      = render.ObjectFitFill
	ObjectFitNone      = render.ObjectFitNone
	ObjectFitScaleDown = render.ObjectFitScaleDown
)
