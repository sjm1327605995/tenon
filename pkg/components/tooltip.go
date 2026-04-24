package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Tooltip 是悬停提示框组件，支持自适应内容高度。
type Tooltip struct {
	core.BaseHost
	title       *Text
	desc        *Text
	stats       *Text
	maxWidth    float32
	padding     float32
}

// NewTooltip 创建一个悬停提示框。
func NewTooltip() *Tooltip {
	t := &Tooltip{
		maxWidth: 280,
		padding:  12,
	}
	t.Init(t)
	t.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	t.GetElement().Yoga.StyleSetPadding(yoga.EdgeAll, t.padding)
	t.GetElement().Yoga.StyleSetGap(yoga.GutterRow, 4)
	t.GetElement().Yoga.StyleSetPositionType(yoga.PositionTypeAbsolute)
	t.GetElement().Yoga.StyleSetMaxWidth(t.maxWidth)
	t.GetElement().BackgroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 250}
	t.GetElement().BorderColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	t.GetElement().BorderRadius = core.BorderRadius{TopLeft: 6, TopRight: 6, BottomRight: 6, BottomLeft: 6}
	t.GetElement().ShadowColor = color.RGBA{R: 0, G: 0, B: 0, A: 60}
	t.GetElement().ShadowBlur = 8
	t.GetElement().ShadowOffsetY = 2

	title := NewText("")
	title.SetFontSize(core.GetTheme().FontSizeBase + 2)
	title.SetColor(core.GetTheme().TextColor)
	title.GetElement().PointerEvents = core.PointerEventsNone
	t.title = title
	t.AddChild(title)

	desc := NewText("")
	desc.SetFontSize(core.GetTheme().FontSizeBase)
	desc.SetColor(core.GetTheme().TextMutedColor)
	desc.GetElement().PointerEvents = core.PointerEventsNone
	t.desc = desc
	t.AddChild(desc)

	stats := NewText("")
	stats.SetFontSize(core.GetTheme().FontSizeSM)
	stats.SetColor(core.GetTheme().PrimaryColor)
	stats.GetElement().PointerEvents = core.PointerEventsNone
	t.stats = stats
	t.AddChild(stats)

	t.Hide()
	return t
}

// SetCardInfo 设置游戏王卡片信息。
func (t *Tooltip) SetCardInfo(name, effect string, atk, def, level int, attr, race string) *Tooltip {
	t.title.SetContent(name)
	if effect != "" {
		t.desc.SetContent(effect)
		t.desc.GetElement().Visible = true
	} else {
		t.desc.GetElement().Visible = false
	}
	statsText := ""
	if level > 0 {
		statsText += "LV" + string(rune('0'+level)) + " "
	}
	if attr != "" {
		statsText += attr + " / "
	}
	if race != "" {
		statsText += race
	}
	if atk >= 0 || def >= 0 {
		statsText += "\n"
		if atk >= 0 {
			statsText += "ATK " + formatNumber(atk) + "  "
		}
		if def >= 0 {
			statsText += "DEF " + formatNumber(def)
		}
	}
	if statsText != "" {
		t.stats.SetContent(statsText)
		t.stats.GetElement().Visible = true
	} else {
		t.stats.GetElement().Visible = false
	}
	return t
}

func formatNumber(n int) string {
	if n < 0 {
		return "?"
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf) - 1
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		buf[i] = byte('0' + n%10)
		n /= 10
		i--
	}
	if neg {
		buf[i] = '-'
		i--
	}
	return string(buf[i+1:])
}

// Show 显示提示框。
func (t *Tooltip) Show() *Tooltip {
	t.GetElement().Visible = true
	return t
}

// Hide 隐藏提示框。
func (t *Tooltip) Hide() *Tooltip {
	t.GetElement().Visible = false
	return t
}

// SetPosition 设置提示框位置（相对于父组件的坐标）。
func (t *Tooltip) SetPosition(x, y float32) *Tooltip {
	t.GetElement().Yoga.StyleSetPosition(yoga.EdgeLeft, x)
	t.GetElement().Yoga.StyleSetPosition(yoga.EdgeTop, y)
	return t
}

// SetMaxWidth 设置最大宽度。
func (t *Tooltip) SetMaxWidth(w float32) *Tooltip {
	t.maxWidth = w
	t.GetElement().Yoga.StyleSetMaxWidth(w)
	return t
}

// Draw 绘制提示框背景和边框。
func (t *Tooltip) Draw(screen *ebiten.Image) {
	el := t.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := t.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if el.BackgroundColor != nil {
		if hasRadius(el.BorderRadius) {
			t.drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BorderRadius, el.BackgroundColor)
		} else {
			vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
		}
	}

	if el.BorderColor != nil {
		if hasRadius(el.BorderRadius) {
			t.drawRoundedRectStroke(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BorderRadius, 1, el.BorderColor)
		} else {
			vector.StrokeRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, 1, el.BorderColor, false)
		}
	}
}

// HandleEvent 消费所有事件防止穿透。
func (t *Tooltip) HandleEvent(e *core.Event) bool {
	return t.GetElement().Visible
}

func (t *Tooltip) drawRoundedRectFill(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

func (t *Tooltip) drawRoundedRectStroke(screen *ebiten.Image, x, y, w, h float32, r core.BorderRadius, stroke float32, clr color.Color) {
	var path vector.Path
	buildRoundedRectPath(&path, x, y, w, h, r)
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.AntiAlias = true
	vector.StrokePath(screen, &path, strokeOp, op)
}

// SyncFrom 同步提示框属性。
func (t *Tooltip) SyncFrom(other core.Host) {
	if o, ok := other.(*Tooltip); ok {
		t.maxWidth = o.maxWidth
		t.padding = o.padding
		if t.title != nil && o.title != nil {
			t.title.Content = o.title.Content
			t.title.cachedLayout = nil
		}
		if t.desc != nil && o.desc != nil {
			t.desc.Content = o.desc.Content
			t.desc.cachedLayout = nil
		}
		if t.stats != nil && o.stats != nil {
			t.stats.Content = o.stats.Content
			t.stats.cachedLayout = nil
		}
	}
}
