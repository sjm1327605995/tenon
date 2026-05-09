package tenon

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/declarative"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	"github.com/sjm1327605995/tenon/pkg/widgets"
)

// 核心类型别名，用户可直接使用。
type (
	Widget            = engine.Widget
	Element           = engine.Element
	StatefulWidget    = engine.StatefulWidget
	BaseWidget        = engine.BaseWidget
	BuildContext      = engine.BuildContext
	Builder           = engine.Builder
	StatefulBuilder   = engine.StatefulBuilder
	GlobalKey         = engine.GlobalKey
	Engine            = engine.Engine
	Theme             = engine.Theme
	EdgeInsets        = engine.EdgeInsets
	Color             = render.Color
	FlexDirection     = engine.FlexDirection
	Justify           = engine.Justify
	Align             = engine.Align
	Wrap              = engine.Wrap
	PositionType      = engine.PositionType
	Display           = engine.Display
	Overflow          = engine.Overflow
	Edge              = engine.Edge
	Gutter            = engine.Gutter
	AnimatedContainer = widgets.AnimatedContainer
	SelectOption      = widgets.SelectOption
	BorderSlice       = render.BorderSlice
	SpriteSheet       = widgets.SpriteSheet
	FragmentWidget    = engine.FragmentWidget
	Hooks             = engine.Hooks
)

// BaseState 是 BaseStateOf[Widget] 的别名，保持向后兼容。
type BaseState = engine.BaseStateOf[engine.Widget]

// BaseStateOf 是泛型 State 基类，避免 GetWidget() 类型断言。
type BaseStateOf[W engine.Widget] = engine.BaseStateOf[W]

// 核心函数。
var (
	NewEngine             = engine.NewEngine
	NewStatefulElement    = engine.NewStatefulElement
	Rebuild               = engine.RebuildDefault
	SetTheme              = engine.SetTheme
	GetTheme              = engine.GetTheme
	ThemeOf               = engine.ThemeOf
	NewBuilder            = engine.NewBuilder
	NewStatefulBuilder    = engine.NewStatefulBuilder
	NewGlobalKey          = engine.NewGlobalKey
	DefaultLightTheme     = engine.DefaultLightTheme
	Fragment              = engine.Fragment
	DefaultDarkTheme      = engine.DefaultDarkTheme
	EdgeInsetsAll         = engine.EdgeInsetsAll
	EdgeInsetsSymmetric   = engine.EdgeInsetsSymmetric
	EdgeInsetsOnly        = engine.EdgeInsetsOnly
	NewColor              = render.NewColor
	NewColorFrom          = render.NewColorFrom
)

// 动画类型别名。
type (
	AnimationController = engine.AnimationController
	AnimationStatus     = engine.AnimationStatus
	Curve               = engine.Curve
	LinearCurve         = engine.LinearCurve
	EaseInOutCurve      = engine.EaseInOutCurve
)

// 注意：Tween[T] 和 Animation[T] 是泛型类型，无法通过别名导出。
// 用户需要直接使用 engine.Tween[T] 和 engine.Animation[T]。

// 动画常量。
const (
	AnimationDismissed = engine.AnimationDismissed
	AnimationForward   = engine.AnimationForward
	AnimationReverse   = engine.AnimationReverse
	AnimationCompleted = engine.AnimationCompleted
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
	FlexDirectionColumn  = engine.FlexDirectionColumn
	FlexDirectionRow     = engine.FlexDirectionRow
	JustifyFlexStart     = engine.JustifyFlexStart
	JustifyCenter        = engine.JustifyCenter
	JustifyFlexEnd       = engine.JustifyFlexEnd
	JustifySpaceBetween  = engine.JustifySpaceBetween
	JustifySpaceAround   = engine.JustifySpaceAround
	JustifySpaceEvenly   = engine.JustifySpaceEvenly
	AlignAuto            = engine.AlignAuto
	AlignFlexStart       = engine.AlignFlexStart
	AlignCenter          = engine.AlignCenter
	AlignFlexEnd         = engine.AlignFlexEnd
	AlignStretch         = engine.AlignStretch
	AlignBaseline        = engine.AlignBaseline
	WrapNoWrap           = engine.WrapNoWrap
	WrapWrap             = engine.WrapWrap
	PositionTypeRelative = engine.PositionTypeRelative
	PositionTypeAbsolute = engine.PositionTypeAbsolute
	DisplayFlex          = engine.DisplayFlex
	DisplayNone          = engine.DisplayNone
	OverflowVisible      = engine.OverflowVisible
	OverflowHidden       = engine.OverflowHidden
	OverflowScroll       = engine.OverflowScroll
	EdgeLeft             = engine.EdgeLeft
	EdgeTop              = engine.EdgeTop
	EdgeRight            = engine.EdgeRight
	EdgeBottom           = engine.EdgeBottom
	EdgeStart            = engine.EdgeStart
	EdgeEnd              = engine.EdgeEnd
	EdgeHorizontal       = engine.EdgeHorizontal
	EdgeVertical         = engine.EdgeVertical
	EdgeAll              = engine.EdgeAll
	GutterColumn         = engine.GutterColumn
	GutterRow            = engine.GutterRow
	GutterAll            = engine.GutterAll
)

// Run 启动 ebiten 窗口并运行应用。
// 根 Widget 会自动包裹在 ThemeProvider 中。
func Run(buildFunc engine.BuildFunc, width, height int) {
	wrappedBuildFunc := func() engine.Widget {
		return engine.ThemeProvider(engine.GetTheme(), buildFunc())
	}
	eng := engine.NewEngine(wrappedBuildFunc, width, height)
	eng.Run()
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
