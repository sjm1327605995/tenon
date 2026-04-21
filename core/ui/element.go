package ui

import "github.com/sjm1327605995/tenon/yoga"

type Element interface {
	Render() Element
	Children() []Element
}

type FlexDirection = yoga.FlexDirection
type Justify = yoga.Justify
type Align = yoga.Align
type Wrap = yoga.Wrap
type Overflow = yoga.Overflow
type Display = yoga.Display
type PositionType = yoga.PositionType
type Direction = yoga.Direction
type Edge = yoga.Edge
type Gutter = yoga.Gutter

const (
	FlexDirectionColumn       = yoga.FlexDirectionColumn
	FlexDirectionColumnReverse = yoga.FlexDirectionColumnReverse
	FlexDirectionRow          = yoga.FlexDirectionRow
	FlexDirectionRowReverse   = yoga.FlexDirectionRowReverse

	JustifyFlexStart    = yoga.JustifyFlexStart
	JustifyCenter       = yoga.JustifyCenter
	JustifyFlexEnd      = yoga.JustifyFlexEnd
	JustifySpaceBetween = yoga.JustifySpaceBetween
	JustifySpaceAround  = yoga.JustifySpaceAround
	JustifySpaceEvenly  = yoga.JustifySpaceEvenly

	AlignAuto       = yoga.AlignAuto
	AlignFlexStart  = yoga.AlignFlexStart
	AlignCenter     = yoga.AlignCenter
	AlignFlexEnd    = yoga.AlignFlexEnd
	AlignStretch    = yoga.AlignStretch
	AlignBaseline   = yoga.AlignBaseline
	AlignSpaceBetween = yoga.AlignSpaceBetween
	AlignSpaceAround  = yoga.AlignSpaceAround
	AlignSpaceEvenly  = yoga.AlignSpaceEvenly

	WrapNoWrap   = yoga.WrapNoWrap
	WrapWrap     = yoga.WrapWrap
	WrapWrapReverse = yoga.WrapWrapReverse

	OverflowVisible = yoga.OverflowVisible
	OverflowHidden  = yoga.OverflowHidden
	OverflowScroll  = yoga.OverflowScroll

	DisplayFlex  = yoga.DisplayFlex
	DisplayNone  = yoga.DisplayNone

	PositionTypeStatic   = yoga.PositionTypeStatic
	PositionTypeRelative  = yoga.PositionTypeRelative
	PositionTypeAbsolute = yoga.PositionTypeAbsolute

	DirectionInherit = yoga.DirectionInherit
	DirectionLTR     = yoga.DirectionLTR
	DirectionRTL     = yoga.DirectionRTL

	EdgeLeft   = yoga.EdgeLeft
	EdgeTop    = yoga.EdgeTop
	EdgeRight  = yoga.EdgeRight
	EdgeBottom = yoga.EdgeBottom
	EdgeStart  = yoga.EdgeStart
	EdgeEnd    = yoga.EdgeEnd
	EdgeHorizontal = yoga.EdgeHorizontal
	EdgeVertical   = yoga.EdgeVertical
	EdgeAll    = yoga.EdgeAll

	GutterColumn = yoga.GutterColumn
	GutterRow    = yoga.GutterRow
	GutterAll    = yoga.GutterAll

	UnitPoint    = yoga.UnitPoint
	UnitPercent  = yoga.UnitPercent
	UnitAuto     = yoga.UnitAuto
	UnitUndefined = yoga.UnitUndefined
)

type ViewStyle struct {
	ID        string
	ClassName string
	Tag       string
	Data      map[string]interface{}
	OnClick   func()
	OnLayout  func(LayoutResults)
	Ref       *ElementRef
}

type ElementRef struct {
	current *ViewElement
}

func (r *ElementRef) Get() *ViewElement {
	return r.current
}

func (r *ElementRef) Set(v *ViewElement) {
	r.current = v
}

type ViewElement struct {
	index    int
	children []Element
	style    ViewStyle
	node     *yoga.Node
}

func View(children ...Element) *ViewElement {
	v := &ViewElement{
		index:    0,
		children: children,
		style: ViewStyle{
			Data: make(map[string]interface{}),
		},
		node: yoga.NewNode(),
	}
	return v
}

func (v *ViewElement) Render() Element {
	return v
}

func (v *ViewElement) Children() []Element {
	return v.children
}

func (v *ViewElement) Node() *yoga.Node {
	return v.node
}

func (v *ViewElement) Style() *ViewStyle {
	return &v.style
}

func (v *ViewElement) Index() int {
	return v.index
}

func (v *ViewElement) AddChild(child Element) *ViewElement {
	v.children = append(v.children, child)
	if ve, ok := child.(*ViewElement); ok {
		v.node.InsertChild(ve.node, v.node.GetChildCount())
	} else if te, ok := child.(*TextElement); ok {
		v.node.InsertChild(te.node, v.node.GetChildCount())
	} else {
		if nodeGetter, ok := child.(interface{ Node() *yoga.Node }); ok {
			v.node.InsertChild(nodeGetter.Node(), v.node.GetChildCount())
		}
	}
	return v
}

func (v *ViewElement) RemoveChild(child Element) *ViewElement {
	for i, c := range v.children {
		if c == child {
			v.children = append(v.children[:i], v.children[i+1:]...)
			if ve, ok := child.(*ViewElement); ok {
				v.node.RemoveChild(ve.node)
			}
			break
		}
	}
	return v
}

func (v *ViewElement) GetChildCount() int {
	return len(v.children)
}

func (v *ViewElement) Width(width float64) *ViewElement {
	v.node.StyleSetWidth(float32(width))
	return v
}

func (v *ViewElement) WidthPercent(percent float64) *ViewElement {
	v.node.StyleSetWidthPercent(float32(percent))
	return v
}

func (v *ViewElement) WidthAuto() *ViewElement {
	v.node.StyleSetWidthAuto()
	return v
}

func (v *ViewElement) Height(height float64) *ViewElement {
	v.node.StyleSetHeight(float32(height))
	return v
}

func (v *ViewElement) HeightPercent(percent float64) *ViewElement {
	v.node.StyleSetHeightPercent(float32(percent))
	return v
}

func (v *ViewElement) HeightAuto() *ViewElement {
	v.node.StyleSetHeightAuto()
	return v
}

func (v *ViewElement) FlexDirection(dir FlexDirection) *ViewElement {
	v.node.StyleSetFlexDirection(dir)
	return v
}

func (v *ViewElement) FlexDirectionColumn() *ViewElement {
	return v.FlexDirection(FlexDirectionColumn)
}

func (v *ViewElement) FlexDirectionRow() *ViewElement {
	return v.FlexDirection(FlexDirectionRow)
}

func (v *ViewElement) FlexDirectionColumnReverse() *ViewElement {
	return v.FlexDirection(FlexDirectionColumnReverse)
}

func (v *ViewElement) FlexDirectionRowReverse() *ViewElement {
	return v.FlexDirection(FlexDirectionRowReverse)
}

func (v *ViewElement) JustifyContent(justify Justify) *ViewElement {
	v.node.StyleSetJustifyContent(justify)
	return v
}

func (v *ViewElement) JustifyContentCenter() *ViewElement {
	return v.JustifyContent(JustifyCenter)
}

func (v *ViewElement) JustifyContentFlexStart() *ViewElement {
	return v.JustifyContent(JustifyFlexStart)
}

func (v *ViewElement) JustifyContentFlexEnd() *ViewElement {
	return v.JustifyContent(JustifyFlexEnd)
}

func (v *ViewElement) JustifyContentSpaceBetween() *ViewElement {
	return v.JustifyContent(JustifySpaceBetween)
}

func (v *ViewElement) JustifyContentSpaceAround() *ViewElement {
	return v.JustifyContent(JustifySpaceAround)
}

func (v *ViewElement) JustifyContentSpaceEvenly() *ViewElement {
	return v.JustifyContent(JustifySpaceEvenly)
}

func (v *ViewElement) AlignItems(align Align) *ViewElement {
	v.node.StyleSetAlignItems(align)
	return v
}

func (v *ViewElement) AlignItemsCenter() *ViewElement {
	return v.AlignItems(AlignCenter)
}

func (v *ViewElement) AlignItemsFlexStart() *ViewElement {
	return v.AlignItems(AlignFlexStart)
}

func (v *ViewElement) AlignItemsFlexEnd() *ViewElement {
	return v.AlignItems(AlignFlexEnd)
}

func (v *ViewElement) AlignItemsStretch() *ViewElement {
	return v.AlignItems(AlignStretch)
}

func (v *ViewElement) AlignItemsBaseline() *ViewElement {
	return v.AlignItems(AlignBaseline)
}

func (v *ViewElement) AlignContent(align Align) *ViewElement {
	v.node.StyleSetAlignContent(align)
	return v
}

func (v *ViewElement) AlignSelf(align Align) *ViewElement {
	v.node.StyleSetAlignSelf(align)
	return v
}

func (v *ViewElement) FlexWrap(wrap Wrap) *ViewElement {
	v.node.StyleSetFlexWrap(wrap)
	return v
}

func (v *ViewElement) FlexWrapNowrap() *ViewElement {
	return v.FlexWrap(WrapNoWrap)
}

func (v *ViewElement) FlexWrapWrap() *ViewElement {
	return v.FlexWrap(WrapWrap)
}

func (v *ViewElement) FlexWrapWrapReverse() *ViewElement {
	return v.FlexWrap(WrapWrapReverse)
}

func (v *ViewElement) Overflow(overflow Overflow) *ViewElement {
	v.node.StyleSetOverflow(overflow)
	return v
}

func (v *ViewElement) OverflowVisible() *ViewElement {
	return v.Overflow(OverflowVisible)
}

func (v *ViewElement) OverflowHidden() *ViewElement {
	return v.Overflow(OverflowHidden)
}

func (v *ViewElement) OverflowScroll() *ViewElement {
	return v.Overflow(OverflowScroll)
}

func (v *ViewElement) Display(display Display) *ViewElement {
	v.node.StyleSetDisplay(display)
	return v
}

func (v *ViewElement) DisplayFlex() *ViewElement {
	return v.Display(DisplayFlex)
}

func (v *ViewElement) DisplayNone() *ViewElement {
	return v.Display(DisplayNone)
}

func (v *ViewElement) Position(typ PositionType) *ViewElement {
	v.node.StyleSetPositionType(typ)
	return v
}

func (v *ViewElement) PositionRelative() *ViewElement {
	return v.Position(PositionTypeRelative)
}

func (v *ViewElement) PositionAbsolute() *ViewElement {
	return v.Position(PositionTypeAbsolute)
}

func (v *ViewElement) PositionStatic() *ViewElement {
	return v.Position(PositionTypeStatic)
}

func (v *ViewElement) PositionLeft(value float64) *ViewElement {
	v.node.StyleSetPosition(EdgeLeft, float32(value))
	return v
}

func (v *ViewElement) PositionTop(value float64) *ViewElement {
	v.node.StyleSetPosition(EdgeTop, float32(value))
	return v
}

func (v *ViewElement) PositionRight(value float64) *ViewElement {
	v.node.StyleSetPosition(EdgeRight, float32(value))
	return v
}

func (v *ViewElement) PositionBottom(value float64) *ViewElement {
	v.node.StyleSetPosition(EdgeBottom, float32(value))
	return v
}

func (v *ViewElement) PositionLeftPercent(percent float64) *ViewElement {
	v.node.StyleSetPositionPercent(EdgeLeft, float32(percent))
	return v
}

func (v *ViewElement) PositionTopPercent(percent float64) *ViewElement {
	v.node.StyleSetPositionPercent(EdgeTop, float32(percent))
	return v
}

func (v *ViewElement) Flex(flex float64) *ViewElement {
	v.node.StyleSetFlex(float32(flex))
	return v
}

func (v *ViewElement) FlexGrow(grow float64) *ViewElement {
	v.node.StyleSetFlexGrow(float32(grow))
	return v
}

func (v *ViewElement) FlexShrink(shrink float64) *ViewElement {
	v.node.StyleSetFlexShrink(float32(shrink))
	return v
}

func (v *ViewElement) FlexBasis(value float64) *ViewElement {
	v.node.StyleSetFlexBasis(float32(value))
	return v
}

func (v *ViewElement) FlexBasisPercent(percent float64) *ViewElement {
	v.node.StyleSetFlexBasisPercent(float32(percent))
	return v
}

func (v *ViewElement) FlexBasisAuto() *ViewElement {
	v.node.StyleSetFlexBasisAuto()
	return v
}

func (v *ViewElement) Margin(edge Edge, value float64) *ViewElement {
	v.node.StyleSetMargin(edge, float32(value))
	return v
}

func (v *ViewElement) MarginAll(value float64) *ViewElement {
	v.node.StyleSetMargin(EdgeAll, float32(value))
	return v
}

func (v *ViewElement) MarginHorizontal(value float64) *ViewElement {
	v.node.StyleSetMargin(EdgeHorizontal, float32(value))
	return v
}

func (v *ViewElement) MarginVertical(value float64) *ViewElement {
	v.node.StyleSetMargin(EdgeVertical, float32(value))
	return v
}

func (v *ViewElement) MarginLeft(value float64) *ViewElement {
	return v.Margin(EdgeLeft, value)
}

func (v *ViewElement) MarginTop(value float64) *ViewElement {
	return v.Margin(EdgeTop, value)
}

func (v *ViewElement) MarginRight(value float64) *ViewElement {
	return v.Margin(EdgeRight, value)
}

func (v *ViewElement) MarginBottom(value float64) *ViewElement {
	return v.Margin(EdgeBottom, value)
}

func (v *ViewElement) MarginAuto(edge Edge) *ViewElement {
	v.node.StyleSetMarginAuto(edge)
	return v
}

func (v *ViewElement) MarginLeftAuto() *ViewElement {
	return v.MarginAuto(EdgeLeft)
}

func (v *ViewElement) MarginRightAuto() *ViewElement {
	return v.MarginAuto(EdgeRight)
}

func (v *ViewElement) MarginPercent(edge Edge, percent float64) *ViewElement {
	v.node.StyleSetMarginPercent(edge, float32(percent))
	return v
}

func (v *ViewElement) Padding(edge Edge, value float64) *ViewElement {
	v.node.StyleSetPadding(edge, float32(value))
	return v
}

func (v *ViewElement) PaddingAll(value float64) *ViewElement {
	v.node.StyleSetPadding(EdgeAll, float32(value))
	return v
}

func (v *ViewElement) PaddingHorizontal(value float64) *ViewElement {
	v.node.StyleSetPadding(EdgeHorizontal, float32(value))
	return v
}

func (v *ViewElement) PaddingVertical(value float64) *ViewElement {
	v.node.StyleSetPadding(EdgeVertical, float32(value))
	return v
}

func (v *ViewElement) PaddingLeft(value float64) *ViewElement {
	return v.Padding(EdgeLeft, value)
}

func (v *ViewElement) PaddingTop(value float64) *ViewElement {
	return v.Padding(EdgeTop, value)
}

func (v *ViewElement) PaddingRight(value float64) *ViewElement {
	return v.Padding(EdgeRight, value)
}

func (v *ViewElement) PaddingBottom(value float64) *ViewElement {
	return v.Padding(EdgeBottom, value)
}

func (v *ViewElement) PaddingPercent(edge Edge, percent float64) *ViewElement {
	v.node.StyleSetPaddingPercent(edge, float32(percent))
	return v
}

func (v *ViewElement) Border(edge Edge, width float64) *ViewElement {
	v.node.StyleSetBorder(edge, float32(width))
	return v
}

func (v *ViewElement) BorderAll(width float64) *ViewElement {
	return v.Border(EdgeAll, width)
}

func (v *ViewElement) BorderLeft(width float64) *ViewElement {
	return v.Border(EdgeLeft, width)
}

func (v *ViewElement) BorderTop(width float64) *ViewElement {
	return v.Border(EdgeTop, width)
}

func (v *ViewElement) BorderRight(width float64) *ViewElement {
	return v.Border(EdgeRight, width)
}

func (v *ViewElement) BorderBottom(width float64) *ViewElement {
	return v.Border(EdgeBottom, width)
}

func (v *ViewElement) Gap(gutter Gutter, value float64) *ViewElement {
	v.node.StyleSetGap(gutter, float32(value))
	return v
}

func (v *ViewElement) GapColumn(value float64) *ViewElement {
	return v.Gap(GutterColumn, value)
}

func (v *ViewElement) GapRow(value float64) *ViewElement {
	return v.Gap(GutterRow, value)
}

func (v *ViewElement) GapAll(value float64) *ViewElement {
	return v.Gap(GutterAll, value)
}

func (v *ViewElement) Direction(dir Direction) *ViewElement {
	v.node.StyleSetDirection(dir)
	return v
}

func (v *ViewElement) DirectionLTR() *ViewElement {
	return v.Direction(DirectionLTR)
}

func (v *ViewElement) DirectionRTL() *ViewElement {
	return v.Direction(DirectionRTL)
}

func (v *ViewElement) Id(id string) *ViewElement {
	v.style.ID = id
	return v
}

func (v *ViewElement) ClassName(className string) *ViewElement {
	v.style.ClassName = className
	return v
}

func (v *ViewElement) Tag(tag string) *ViewElement {
	v.style.Tag = tag
	return v
}

func (v *ViewElement) Data(key string, value interface{}) *ViewElement {
	v.style.Data[key] = value
	return v
}

func (v *ViewElement) GetData(key string) interface{} {
	return v.style.Data[key]
}

func (v *ViewElement) OnClick(handler func()) *ViewElement {
	v.style.OnClick = handler
	return v
}

func (v *ViewElement) OnLayout(handler func(LayoutResults)) *ViewElement {
	v.style.OnLayout = handler
	return v
}

func (v *ViewElement) Ref(ref *ElementRef) *ViewElement {
	v.style.Ref = ref
	if ref != nil {
		ref.Set(v)
	}
	return v
}

func (v *ViewElement) GetRef() *ElementRef {
	return v.style.Ref
}

func (v *ViewElement) BackgroundColor(color uint32) *ViewElement {
	v.style.Data["backgroundColor"] = color
	return v
}

func (v *ViewElement) Opacity(opacity float64) *ViewElement {
	v.style.Data["opacity"] = opacity
	return v
}

func (v *ViewElement) ZIndex(index int) *ViewElement {
	v.style.Data["zIndex"] = index
	return v
}

func (v *ViewElement) Transform(transform string) *ViewElement {
	v.style.Data["transform"] = transform
	return v
}

func (v *ViewElement) GetLayout() LayoutResults {
	return getNodeLayout(v.node)
}

func (v *ViewElement) CalculateLayout(width, height float64, direction Direction) {
	v.node.CalculateLayout(float32(width), float32(height), direction)
	if v.style.OnLayout != nil {
		v.style.OnLayout(getNodeLayout(v.node))
	}
}

func (v *ViewElement) MarkDirty() {
	v.node.MarkDirty()
}

func (v *ViewElement) MarkDirtyAndPropagate() {
	v.node.MarkDirty()
}

type TextElement struct {
	text   string
	style  ViewStyle
	node   *yoga.Node
}

func Text(text string) *TextElement {
	t := &TextElement{
		text:  text,
		style: ViewStyle{Data: make(map[string]interface{})},
		node:  yoga.NewNode(),
	}
	return t
}

func (t *TextElement) Render() Element {
	return t
}

func (t *TextElement) Children() []Element {
	return nil
}

func (t *TextElement) Text() string {
	return t.text
}

func (t *TextElement) SetText(text string) *TextElement {
	t.text = text
	return t
}

func (t *TextElement) Style() *ViewStyle {
	return &t.style
}

func (t *TextElement) Node() *yoga.Node {
	return t.node
}

func (t *TextElement) Width(width float64) *TextElement {
	t.node.StyleSetWidth(float32(width))
	return t
}

func (t *TextElement) Height(height float64) *TextElement {
	t.node.StyleSetHeight(float32(height))
	return t
}

func (t *TextElement) FlexGrow(grow float64) *TextElement {
	t.node.StyleSetFlexGrow(float32(grow))
	return t
}

func (t *TextElement) FlexShrink(shrink float64) *TextElement {
	t.node.StyleSetFlexShrink(float32(shrink))
	return t
}

func (t *TextElement) Margin(edge Edge, value float64) *TextElement {
	t.node.StyleSetMargin(edge, float32(value))
	return t
}

func (t *TextElement) Padding(edge Edge, value float64) *TextElement {
	t.node.StyleSetPadding(edge, float32(value))
	return t
}

func (t *TextElement) Id(id string) *TextElement {
	t.style.ID = id
	return t
}

func (t *TextElement) ClassName(className string) *TextElement {
	t.style.ClassName = className
	return t
}

func (t *TextElement) Data(key string, value interface{}) *TextElement {
	t.style.Data[key] = value
	return t
}

func (t *TextElement) Color(color uint32) *TextElement {
	t.style.Data["color"] = color
	return t
}

func (t *TextElement) FontSize(size float64) *TextElement {
	t.style.Data["fontSize"] = size
	return t
}

func (t *TextElement) GetLayout() LayoutResults {
	return getNodeLayout(t.node)
}

func (t *TextElement) CalculateLayout(width, height float64, direction Direction) {
	t.node.CalculateLayout(float32(width), float32(height), direction)
	if t.style.OnLayout != nil {
		t.style.OnLayout(getNodeLayout(t.node))
	}
}

func HStack(children ...Element) *ViewElement {
	view := View(children...)
	view.FlexDirectionRow()
	return view
}

func VStack(children ...Element) *ViewElement {
	view := View(children...)
	view.FlexDirectionColumn()
	return view
}

func ZStack(children ...Element) *ViewElement {
	view := View(children...)
	view.PositionAbsolute()
	return view
}


