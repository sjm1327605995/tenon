package render

import "image/color"

type RenderCommandType int32

const (
	RENDER_COMMAND_TYPE_NONE = RenderCommandType(iota)
	RENDER_COMMAND_TYPE_RECTANGLE
	RENDER_COMMAND_TYPE_BORDER
	RENDER_COMMAND_TYPE_TEXT
	RENDER_COMMAND_TYPE_IMAGE
	RENDER_COMMAND_TYPE_SCISSOR_START
	RENDER_COMMAND_TYPE_SCISSOR_END
	RENDER_COMMAND_TYPE_CUSTOM
)

type BoundingBox struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}
type RenderCommand struct {
	BoundingBox BoundingBox
	RenderData  RenderData

	Id          uint32
	ZIndex      int16
	CommandType RenderCommandType
}
type RenderData struct {
	// union
	Rectangle RectangleRenderData
	Text      TextRenderData
	Image     ImageRenderData
	Custom    CustomRenderData
	Border    BorderRenderData
	Clip      ClipRenderData
}
type TextRenderData struct {
	StringContents string
	TextColor      color.RGBA
	FontId         uint16
	FontSize       uint16
	LetterSpacing  uint16
	LineHeight     uint16
}
type CornerRadius struct {
	TopLeft     float32
	TopRight    float32
	BottomLeft  float32
	BottomRight float32
}
type BorderWidth struct {
	Left            uint16
	Right           uint16
	Top             uint16
	Bottom          uint16
	BetweenChildren uint16
}
type RectangleRenderData struct {
	BackgroundColor color.RGBA
	CornerRadius    CornerRadius
}
type ImageRenderData struct {
	BackgroundColor color.RGBA
	CornerRadius    CornerRadius
	ImageData       any
}
type CustomRenderData struct {
	BackgroundColor color.RGBA
	CornerRadius    CornerRadius
	CustomData      any
}
type ScrollRenderData struct {
	Horizontal bool
	Vertical   bool
}
type ClipRenderData ScrollRenderData
type BorderRenderData struct {
	Color        color.RGBA
	CornerRadius CornerRadius
	Width        BorderWidth
}
