package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// DialogButton 定义对话框按钮。
type DialogButton struct {
	Label  string
	Type   string // "primary" | "default" | "danger"
	OnClick func()
}

// Dialog 是确认/取消对话框组件。
type Dialog struct {
	core.BaseHost
	title      string
	content    string
	buttons    []DialogButton
	overlay    *View
	panel      *View
	titleText  *Text
	contentText *Text
	buttonRow  *View
}

// NewDialog 创建一个对话框。
func NewDialog() *Dialog {
	d := &Dialog{}
	d.Init(d)
	d.GetElement().Yoga.StyleSetPositionType(yoga.PositionTypeAbsolute)
	d.GetElement().Yoga.StyleSetPosition(yoga.EdgeLeft, 0)
	d.GetElement().Yoga.StyleSetPosition(yoga.EdgeTop, 0)
	d.GetElement().Yoga.StyleSetWidthPercent(100)
	d.GetElement().Yoga.StyleSetHeightPercent(100)
	d.GetElement().Yoga.StyleSetJustifyContent(yoga.JustifyCenter)
	d.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
	d.GetElement().BackgroundColor = color.RGBA{R: 0, G: 0, B: 0, A: 100}
	d.GetElement().PointerEvents = core.PointerEventsAuto

	// 遮罩层点击关闭
	d.overlay = NewView()
	d.overlay.Init(d.overlay)
	d.overlay.GetElement().Yoga.StyleSetPositionType(yoga.PositionTypeAbsolute)
	d.overlay.GetElement().Yoga.StyleSetPosition(yoga.EdgeLeft, 0)
	d.overlay.GetElement().Yoga.StyleSetPosition(yoga.EdgeTop, 0)
	d.overlay.GetElement().Yoga.StyleSetWidthPercent(100)
	d.overlay.GetElement().Yoga.StyleSetHeightPercent(100)
	d.overlay.GetElement().PointerEvents = core.PointerEventsAuto
	d.AddChild(d.overlay)

	// 面板
	panel := NewView()
	panel.Init(panel)
	panel.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	panel.GetElement().Yoga.StyleSetWidth(360)
	panel.GetElement().Yoga.StyleSetMaxWidthPercent(90)
	panel.GetElement().Yoga.StyleSetPadding(yoga.EdgeAll, 24)
	panel.GetElement().BackgroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	panel.GetElement().BorderRadius = core.BorderRadius{TopLeft: 8, TopRight: 8, BottomRight: 8, BottomLeft: 8}
	panel.GetElement().ShadowColor = color.RGBA{R: 0, G: 0, B: 0, A: 80}
	panel.GetElement().ShadowBlur = 16
	panel.GetElement().ShadowOffsetY = 4
	panel.GetElement().PointerEvents = core.PointerEventsAuto
	d.panel = panel
	d.AddChild(panel)

	// 标题
	titleText := NewText("")
	titleText.SetFontSize(core.GetTheme().FontSizeLG)
	titleText.SetColor(core.GetTheme().TextColor)
	titleText.GetElement().PointerEvents = core.PointerEventsNone
	d.titleText = titleText
	panel.AddChild(titleText)

	// 内容
	contentText := NewText("")
	contentText.SetFontSize(core.GetTheme().FontSizeBase)
	contentText.SetColor(core.GetTheme().TextMutedColor)
	contentText.GetElement().Yoga.StyleSetMargin(yoga.EdgeTop, 12)
	contentText.GetElement().PointerEvents = core.PointerEventsNone
	d.contentText = contentText
	panel.AddChild(contentText)

	// 按钮行
	buttonRow := NewView()
	buttonRow.Init(buttonRow)
	buttonRow.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	buttonRow.GetElement().Yoga.StyleSetJustifyContent(yoga.JustifyFlexEnd)
	buttonRow.GetElement().Yoga.StyleSetGap(yoga.GutterColumn, 8)
	buttonRow.GetElement().Yoga.StyleSetMargin(yoga.EdgeTop, 20)
	buttonRow.GetElement().PointerEvents = core.PointerEventsAuto
	d.buttonRow = buttonRow
	panel.AddChild(buttonRow)

	d.Hide()
	return d
}

// SetTitle 设置标题。
func (d *Dialog) SetTitle(title string) *Dialog {
	d.title = title
	d.titleText.SetContent(title)
	return d
}

// SetContent 设置内容文本。
func (d *Dialog) SetContent(content string) *Dialog {
	d.content = content
	d.contentText.SetContent(content)
	return d
}

// SetButtons 设置按钮。
func (d *Dialog) SetButtons(buttons []DialogButton) *Dialog {
	d.buttons = buttons
	d.buttonRow.ClearChildren()
	for _, btn := range buttons {
		b := NewButton(btn.Label)
		switch btn.Type {
		case "primary":
			b.SetType(ButtonTypePrimary)
		case "danger":
			b.SetType(ButtonTypeDanger)
		default:
			b.SetType(ButtonTypeDefault)
		}
		if btn.OnClick != nil {
			b.SetOnClick(btn.OnClick)
		}
		d.buttonRow.AddChild(b)
	}
	return d
}

// Show 显示对话框。
func (d *Dialog) Show() *Dialog {
	d.GetElement().Visible = true
	return d
}

// Hide 隐藏对话框。
func (d *Dialog) Hide() *Dialog {
	d.GetElement().Visible = false
	return d
}

// HandleEvent 消费所有事件防止穿透。
func (d *Dialog) HandleEvent(e *core.Event) bool {
	if e.Type == core.EventClick {
		// 点击遮罩层关闭
		if d.overlay != nil {
			bounds := d.overlay.GetLayoutBounds()
			if e.X >= bounds.X && e.X < bounds.X+bounds.Width &&
				e.Y >= bounds.Y && e.Y < bounds.Y+bounds.Height {
				d.Hide()
				return true
			}
		}
	}
	return d.GetElement().Visible
}

// SyncFrom 同步对话框属性。
func (d *Dialog) SyncFrom(other core.Host) {
	if o, ok := other.(*Dialog); ok {
		d.title = o.title
		d.content = o.content
		d.buttons = o.buttons
		d.titleText.SetContent(d.title)
		d.contentText.SetContent(d.content)
		d.SetButtons(d.buttons)
	}
}
