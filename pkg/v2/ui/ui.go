// Package ui 提供 UI 引擎核心类型。
//
// 已迁移至 internal/engine，此包仅为兼容性保留。
// 请直接使用 tenon 包的统一入口。
package ui

import (
	internal "github.com/sjm1327605995/tenon/internal/engine"
	"github.com/sjm1327605995/tenon/internal/render"
)

// 核心类型别名。
type (
	Widget            = internal.Widget
	Element           = internal.Element
	Key               = internal.Key
	NilKey            = internal.NilKey
	GlobalKey         = internal.GlobalKey
	BaseWidget        = internal.BaseWidget
	BaseElement       = internal.BaseElement
	ComponentElement  = internal.ComponentElement
	StatefulWidget    = internal.StatefulWidget
	State             = internal.State
	BaseState         = internal.BaseState
	StatefulElement   = internal.StatefulElement
	SingleChildComponentElement = internal.SingleChildComponentElement
	RenderObjectElement = internal.RenderObjectElement
	SingleChildRenderObjectElement = internal.SingleChildRenderObjectElement
	MultiChildRenderObjectElement = internal.MultiChildRenderObjectElement
	BuildContext      = internal.BuildContext
	BuildFunc         = internal.BuildFunc
	Engine            = internal.Engine
	Theme             = internal.Theme
	InheritedWidget   = internal.InheritedWidget
	InheritedElement  = internal.InheritedElement
	InheritedTheme    = internal.InheritedTheme
	Localization      = internal.Localization
	NavigatorWidget   = internal.NavigatorWidget
	NavigatorState    = internal.NavigatorState
	NavigatorContext  = internal.NavigatorContext
	Page              = internal.Page
	PageTransition    = internal.PageTransition
	RouteBuilder      = internal.RouteBuilder
	RouteParams       = internal.RouteParams
	FocusNode         = internal.FocusNode
	FocusManager      = internal.FocusManager
	Focusabler        = internal.Focusabler
	Scrollabler       = internal.Scrollabler
	Clipper           = internal.Clipper
	InputHandler      = internal.InputHandler
	TextInputHandler  = internal.TextInputHandler
	Blinker           = internal.Blinker
	ClipboardHandler  = internal.ClipboardHandler
	Clipboard         = internal.Clipboard
	ShortcutManager   = internal.ShortcutManager
	Shortcut          = internal.Shortcut
	Inspector         = internal.Inspector
	PerformanceOverlay = internal.PerformanceOverlay
	AnimationController = internal.AnimationController
	AnimationStatus   = internal.AnimationStatus

	Curve             = internal.Curve
	LinearCurve       = internal.LinearCurve
	EaseInOutCurve    = internal.EaseInOutCurve
	AnimationGroup    = internal.AnimationGroup
	AnimationSequence = internal.AnimationSequence
	CurvedAnimation   = internal.CurvedAnimation
	DelayedAnimation  = internal.DelayedAnimation
	PathAnimation     = internal.PathAnimation
	Path              = internal.Path
	LinePath          = internal.LinePath
	QuadBezierPath    = internal.QuadBezierPath
	CubicBezierPath   = internal.CubicBezierPath
	ArcPath           = internal.ArcPath
	EdgeInsets        = internal.EdgeInsets
	Builder           = internal.Builder
	StatefulBuilder   = internal.StatefulBuilder
	OpacityWidget     = internal.OpacityWidget
	SlideOffsetWidget = internal.SlideOffsetWidget
	TooltipWidget     = internal.TooltipWidget
	TooltipOverlay    = internal.TooltipOverlay
	ContextMenuWidget = internal.ContextMenuWidget
	MenuBarWidget     = internal.MenuBarWidget
	MenuItem          = internal.MenuItem
	MenuBarItem       = internal.MenuBarItem
	OffsetTween       = internal.OffsetTween
	SizeTween         = internal.SizeTween
	ColorTween        = internal.ColorTween
	TransformTween    = internal.TransformTween
	Float32Tween      = internal.Float32Tween
	TestEnvironment   = internal.TestEnvironment
)

// 核心函数。
var (
	NewRenderObjectElement         = internal.NewRenderObjectElement
	NewSingleChildRenderObjectElement = internal.NewSingleChildRenderObjectElement
	NewMultiChildRenderObjectElement  = internal.NewMultiChildRenderObjectElement
	NewComponentElement            = internal.NewComponentElement
	NewSingleChildComponentElement = internal.NewSingleChildComponentElement

	NewEngine             = internal.NewEngine
	NewStatefulElement    = internal.NewStatefulElement
	RebuildDefault        = internal.RebuildDefault
	DefaultEngine         = internal.DefaultEngine
	SetTheme              = internal.SetTheme
	GetTheme              = internal.GetTheme
	ThemeOf               = internal.ThemeOf
	ThemeProvider         = internal.ThemeProvider
	DefaultLightTheme     = internal.DefaultLightTheme
	DefaultDarkTheme      = internal.DefaultDarkTheme
	NewBuilder            = internal.NewBuilder
	NewStatefulBuilder    = internal.NewStatefulBuilder
	NewGlobalKey          = internal.NewGlobalKey
	NewLocalKey           = internal.NewLocalKey
	NewStringKey          = internal.NewStringKey
	IsNilKey              = internal.IsNilKey
	CanUpdate             = internal.CanUpdate
	UpdateChild           = internal.UpdateChild
	UpdateChildren        = internal.UpdateChildren
	EdgeInsetsAll         = internal.EdgeInsetsAll
	EdgeInsetsSymmetric   = internal.EdgeInsetsSymmetric
	EdgeInsetsOnly        = internal.EdgeInsetsOnly
	NewInheritedElement   = internal.NewInheritedElement
	NewLocalization       = internal.NewLocalization
	L                     = internal.L
	Lf                    = internal.Lf
	GetLocalization       = internal.GetLocalization
	GetLocale             = internal.GetLocale
	NewNavigatorContext   = internal.NewNavigatorContext
	Navigator             = internal.Navigator
	NavigatorWithParams   = internal.NavigatorWithParams
	GetNavigator          = internal.GetNavigator
	NavPush               = internal.NavPush
	NavPop                = internal.NavPop
	NavPushReplacement    = internal.NavPushReplacement
	NavPopToRoot          = internal.NavPopToRoot
	NewFocusManager       = internal.NewFocusManager
	RegisterPopupDismisser   = internal.RegisterPopupDismisser
	UnregisterPopupDismisser = internal.UnregisterPopupDismisser
	DismissAllPopups         = internal.DismissAllPopups
	NewShortcutManager    = internal.NewShortcutManager
	NewClipboard          = internal.NewClipboard
	GetClipboard          = internal.GetClipboard
	NewTooltipOverlay     = internal.NewTooltipOverlay
	TestWidget            = internal.TestWidget
	Opacity               = internal.Opacity
	SlideOffset           = internal.SlideOffset
	Tooltip               = internal.Tooltip
	ContextMenu           = internal.ContextMenu
	MenuBar               = internal.MenuBar
)

// 动画常量。
const (
	AnimationDismissed = internal.AnimationDismissed
	AnimationForward   = internal.AnimationForward
	AnimationReverse   = internal.AnimationReverse
	AnimationCompleted = internal.AnimationCompleted
)

// 导航常量。
const (
	TransitionNone    = internal.TransitionNone
	TransitionSlide   = internal.TransitionSlide
	TransitionFade    = internal.TransitionFade
	TransitionSlideUp = internal.TransitionSlideUp
)

// 快捷键常量。
const (
	ShortcutCtrl  = internal.ShortcutCtrl
	ShortcutShift = internal.ShortcutShift
	ShortcutAlt   = internal.ShortcutAlt
)

// Tooltip 位置常量。
const (
	TooltipTop    = internal.TooltipTop
	TooltipBottom = internal.TooltipBottom
	TooltipLeft   = internal.TooltipLeft
	TooltipRight  = internal.TooltipRight
)

// Flex/Yoga 类型别名。
type (
	FlexDirection = internal.FlexDirection
	Justify       = internal.Justify
	Align         = internal.Align
	Wrap          = internal.Wrap
	PositionType  = internal.PositionType
	Display       = internal.Display
	Overflow      = internal.Overflow
	Edge          = internal.Edge
	Gutter        = internal.Gutter
	Direction     = internal.Direction
)

// Render 类型别名（供 ui 包用户使用）。
type RenderObject = render.RenderObject

// Flex/Yoga 常量。
const (
	FlexDirectionColumn    = internal.FlexDirectionColumn
	FlexDirectionRow       = internal.FlexDirectionRow
	JustifyFlexStart       = internal.JustifyFlexStart
	JustifyCenter          = internal.JustifyCenter
	JustifyFlexEnd         = internal.JustifyFlexEnd
	JustifySpaceBetween    = internal.JustifySpaceBetween
	JustifySpaceAround     = internal.JustifySpaceAround
	JustifySpaceEvenly     = internal.JustifySpaceEvenly
	AlignAuto              = internal.AlignAuto
	AlignFlexStart         = internal.AlignFlexStart
	AlignCenter            = internal.AlignCenter
	AlignFlexEnd           = internal.AlignFlexEnd
	AlignStretch           = internal.AlignStretch
	AlignBaseline          = internal.AlignBaseline
	WrapNoWrap             = internal.WrapNoWrap
	WrapWrap               = internal.WrapWrap
	PositionTypeRelative   = internal.PositionTypeRelative
	PositionTypeAbsolute   = internal.PositionTypeAbsolute
	DisplayFlex            = internal.DisplayFlex
	DisplayNone            = internal.DisplayNone
	OverflowVisible        = internal.OverflowVisible
	OverflowHidden         = internal.OverflowHidden
	OverflowScroll         = internal.OverflowScroll
	EdgeLeft               = internal.EdgeLeft
	EdgeTop                = internal.EdgeTop
	EdgeRight              = internal.EdgeRight
	EdgeBottom             = internal.EdgeBottom
	EdgeStart              = internal.EdgeStart
	EdgeEnd                = internal.EdgeEnd
	EdgeHorizontal         = internal.EdgeHorizontal
	EdgeVertical           = internal.EdgeVertical
	EdgeAll                = internal.EdgeAll
	GutterColumn           = internal.GutterColumn
	GutterRow              = internal.GutterRow
	GutterAll              = internal.GutterAll
	DirectionLTR           = internal.DirectionLTR
	DirectionRTL           = internal.DirectionRTL
)
