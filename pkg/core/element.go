package core

import (
	"image/color"

	"github.com/sjm1327605995/tenon/yoga"
)

// Element 是宿主组件的样式与布局节点。
// 每个 Host 拥有一个 Element，包含 Yoga 节点和视觉属性。
type Element struct {
	Yoga            *yoga.Node
	Visible         bool
	BackgroundColor color.Color
	BorderColor     color.Color
	ShadowColor     color.Color
	ShadowBlur      float32
	ShadowOffsetX   float32
	ShadowOffsetY   float32
	PointerEvents   PointerEvents
	BorderRadius    BorderRadius
	Key             string
}

// LayoutBounds 描述组件在屏幕上的位置和尺寸。
type LayoutBounds struct {
	X, Y, Width, Height float32
}

// BorderRadius 描述四个角的圆角半径。
type BorderRadius struct {
	TopLeft, TopRight, BottomRight, BottomLeft float32
}

// PointerEvents 控制组件是否响应指针事件。
type PointerEvents int

const (
	PointerEventsAuto PointerEvents = iota
	PointerEventsNone
)

// NewElement 创建一个带 Yoga 节点的默认 Element。
func NewElement() *Element {
	return &Element{
		Yoga:          yoga.NewNode(),
		Visible:       true,
		PointerEvents: PointerEventsAuto,
	}
}
