package tenon

import (
	"image/color"

	"github.com/sjm1327605995/tenon/internal/engine"
	"github.com/sjm1327605995/tenon/pkg/v2/declarative"
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
	BorderSlice       = render.BorderSlice
	SpriteSheet       = widgets.SpriteSheet
	FragmentWidget    = engine.FragmentWidget
	Hooks             = engine.Hooks
)

// BaseState 是 BaseStateOf[Widget] 的别名，保持向后兼容。
type BaseState = ui.BaseStateOf[ui.Widget]

// BaseStateOf 是泛型 State 基类，避免 GetWidget() 类型断言。
type BaseStateOf[W ui.Widget] = ui.BaseStateOf[W]

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
	Fragment              = engine.Fragment
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
	Curve               = ui.Curve
	LinearCurve         = ui.LinearCurve
	EaseInOutCurve      = ui.EaseInOutCurve
)

// 注意：Tween[T] 和 Animation[T] 是泛型类型，无法通过别名导出。
// 用户需要直接使用 ui.Tween[T] 和 ui.Animation[T]。

// 动画常量。
const (
	AnimationDismissed = ui.AnimationDismissed
	AnimationForward   = ui.AnimationForward
	AnimationReverse   = ui.AnimationReverse
	AnimationCompleted = ui.AnimationCompleted
)

// Widget 构造函数 — React/SwiftUI 声明式 API（推荐）。
// 用户可直接通过 tenon.XXX 调用，无需再 import 子包。
var (
	// 基础组件
	Text       = declarative.Text
	Button     = declarative.Button
	Input      = declarative.Input
	Image      = declarative.Image
	Icon       = declarative.Icon
	Scroll     = declarative.Scroll
	Container  = declarative.Container
	CardBox    = declarative.Card
	Animated   = declarative.Animated
	Navigator  = declarative.Navigator

	// 布局
	VStack     = declarative.VStack
	HStack     = declarative.HStack
	Spacer     = declarative.Spacer
	Stack      = declarative.Stack
	Positioned = declarative.Positioned
	Grid       = declarative.Grid

	// 游戏 GUI
	NinePatch       = declarative.NinePatch
	Sprite          = declarative.Sprite
	SpriteButton    = declarative.SpriteButton
	GameProgressBar = declarative.GameProgressBar
	Card3D          = declarative.Card3D

	// 交互与特效
	Draggable  = declarative.Draggable
	DropTarget = declarative.DropTarget
	FlipCard   = declarative.FlipCard
	Transform  = declarative.Transform

	// 表单
	Select = declarative.Select

	// 兼容保留 — imperative widgets API（旧代码仍可用）
	Row                  = widgets.Row
	Column               = widgets.Column
	TextField            = widgets.TextField
	EditableText         = widgets.EditableText
	NewAnimatedContainer = widgets.NewAnimatedContainer
	MultiSelect          = widgets.MultiSelect
)

// 图标常量。
const (
	IconArrowUp         = widgets.IconArrowUp
	IconArrowDown       = widgets.IconArrowDown
	IconArrowLeft       = widgets.IconArrowLeft
	IconArrowRight      = widgets.IconArrowRight
	IconChevronRight    = widgets.IconChevronRight
	IconChevronDown     = widgets.IconChevronDown
	IconCheckboxEmpty   = widgets.IconCheckboxEmpty
	IconCheckboxChecked = widgets.IconCheckboxChecked
	IconInfo            = widgets.IconInfo
	IconClose           = widgets.IconClose
	IconCheck           = widgets.IconCheck
	IconMinus           = widgets.IconMinus
)

// 图标配置。
var (
	RegisterIconFallback = widgets.RegisterIconFallback
	NewSpriteSheet       = widgets.NewSpriteSheet
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

// ==================== Hooks 泛型函数包装 ====================
// internal/engine 是 internal 包，泛型函数无法通过 var 转发，
// 以下直接在 tenon 包提供同名泛型函数。

func UseState[T any](h *Hooks, initial T) (func() T, func(T)) {
	return engine.UseState[T](h, initial)
}

func UseMemo[T any](h *Hooks, factory func() T, deps []any) T {
	return engine.UseMemo[T](h, factory, deps)
}

func UseCallback[T any](h *Hooks, fn T, deps []any) T {
	return engine.UseCallback[T](h, fn, deps)
}

// ObjectFit 快捷导出
const (
	ObjectFitCover     = render.ObjectFitCover
	ObjectFitContain   = render.ObjectFitContain
	ObjectFitFill      = render.ObjectFitFill
	ObjectFitNone      = render.ObjectFitNone
	ObjectFitScaleDown = render.ObjectFitScaleDown
)

// RGBA 返回标准库 color.Color，避免 *Color 解引用问题。
func RGBA(r, g, b, a uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: a}
}
