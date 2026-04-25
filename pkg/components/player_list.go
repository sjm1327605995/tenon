package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// PlayerInfo 定义玩家信息。
type PlayerInfo struct {
	ID       string
	Name     string
	Ready    bool
	Host     bool // 房主
	Duelist  bool // 决斗者（非观众）
}

// PlayerList 是玩家列表组件。
type PlayerList struct {
	core.BaseHost
	players    []PlayerInfo
	itemHosts  map[string]*View
}

// NewPlayerList 创建一个玩家列表。
func NewPlayerList() *PlayerList {
	pl := &PlayerList{
		itemHosts: make(map[string]*View),
	}
	pl.Init(pl)
	pl.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	pl.GetElement().Yoga.StyleSetGap(yoga.GutterRow, 2)
	pl.GetElement().Yoga.StyleSetWidth(200)
	pl.GetElement().BackgroundColor = core.GetTheme().BackgroundColor
	return pl
}

// SetPlayers 设置玩家列表数据。
func (pl *PlayerList) SetPlayers(players []PlayerInfo) *PlayerList {
	pl.players = players
	pl.rebuild()
	return pl
}

func (pl *PlayerList) rebuild() {
	pl.ClearChildren()
	pl.itemHosts = make(map[string]*View)

	for _, player := range pl.players {
		row := NewView()
		row.Init(row)
		row.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
		row.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
		row.GetElement().Yoga.StyleSetPadding(yoga.EdgeVertical, 6)
		row.GetElement().Yoga.StyleSetPadding(yoga.EdgeHorizontal, 10)
		row.GetElement().PointerEvents = core.PointerEventsNone

		// 状态指示器
		status := NewView()
		status.Init(status)
		status.GetElement().Yoga.StyleSetWidth(8)
		status.GetElement().Yoga.StyleSetHeight(8)
		status.GetElement().Yoga.StyleSetMargin(yoga.EdgeRight, 8)
		status.GetElement().Yoga.StyleSetBorderRadius(4)
		status.GetElement().PointerEvents = core.PointerEventsNone
		if player.Ready {
			status.GetElement().BackgroundColor = color.RGBA{R: 56, G: 158, B: 13, A: 255}
		} else {
			status.GetElement().BackgroundColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
		}
		row.AddChild(status)

		// 房主标识
		if player.Host {
			hostBadge := NewText("[房主]")
			hostBadge.SetFontSize(core.GetTheme().FontSizeSM)
			hostBadge.SetColor(core.GetTheme().PrimaryColor)
			hostBadge.GetElement().Yoga.StyleSetMargin(yoga.EdgeRight, 4)
			hostBadge.GetElement().PointerEvents = core.PointerEventsNone
			row.AddChild(hostBadge)
		}

		// 玩家名
		nameText := NewText(player.Name)
		nameText.SetFontSize(core.GetTheme().FontSizeBase)
		if player.Ready {
			nameText.SetColor(core.GetTheme().TextColor)
		} else {
			nameText.SetColor(core.GetTheme().TextMutedColor)
		}
		nameText.GetElement().Yoga.StyleSetFlexGrow(1)
		nameText.GetElement().PointerEvents = core.PointerEventsNone
		row.AddChild(nameText)

		// 准备状态
		if player.Ready {
			readyText := NewText("✓")
			readyText.SetFontSize(core.GetTheme().FontSizeSM)
			readyText.SetColor(color.RGBA{R: 56, G: 158, B: 13, A: 255})
			readyText.GetElement().PointerEvents = core.PointerEventsNone
			row.AddChild(readyText)
		}

		pl.itemHosts[player.ID] = row
		pl.AddChild(row)
	}
}

// Draw 绘制列表背景。
func (pl *PlayerList) Draw(screen *ebiten.Image) {
	el := pl.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := pl.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}
}

// ==================== 链式 API ====================

func (pl *PlayerList) SetWidth(width float32) *PlayerList {
	pl.GetElement().Yoga.StyleSetWidth(width)
	return pl
}
func (pl *PlayerList) SetMargin(edge yoga.Edge, value float32) *PlayerList {
	pl.GetElement().Yoga.StyleSetMargin(edge, value)
	return pl
}

// SyncFrom 同步玩家列表属性。
func (pl *PlayerList) SyncFrom(other core.Host) {
	if o, ok := other.(*PlayerList); ok {
		pl.players = o.players
		pl.rebuild()
	}
}
