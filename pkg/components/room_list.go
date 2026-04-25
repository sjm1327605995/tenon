package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// RoomInfo 定义房间信息。
type RoomInfo struct {
	ID       string
	Name     string
	Players  int
	MaxPlayers int
	Password bool
	Rule     string // 规则说明
}

// RoomList 是房间列表组件。
type RoomList struct {
	core.BaseHost
	rooms      []RoomInfo
	selectedID string
	onSelect   func(roomID string)
	onJoin     func(roomID string)
	itemHosts  map[string]*View
}

// NewRoomList 创建一个房间列表。
func NewRoomList() *RoomList {
	rl := &RoomList{
		itemHosts: make(map[string]*View),
	}
	rl.Init(rl)
	rl.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	rl.GetElement().Yoga.StyleSetGap(yoga.GutterRow, 2)
	rl.GetElement().Yoga.StyleSetWidth(360)
	rl.GetElement().Yoga.StyleSetHeight(400)
	rl.GetElement().Yoga.StyleSetOverflow(yoga.OverflowScroll)
	rl.GetElement().BackgroundColor = core.GetTheme().BackgroundColor
	return rl
}

// SetRooms 设置房间列表数据。
func (rl *RoomList) SetRooms(rooms []RoomInfo) *RoomList {
	rl.rooms = rooms
	rl.rebuild()
	return rl
}

// GetSelected 获取当前选中的房间ID。
func (rl *RoomList) GetSelected() string {
	return rl.selectedID
}

// SetOnSelect 设置选择回调。
func (rl *RoomList) SetOnSelect(fn func(roomID string)) *RoomList {
	rl.onSelect = fn
	return rl
}

// SetOnJoin 设置加入回调。
func (rl *RoomList) SetOnJoin(fn func(roomID string)) *RoomList {
	rl.onJoin = fn
	return rl
}

func (rl *RoomList) rebuild() {
	rl.ClearChildren()
	rl.itemHosts = make(map[string]*View)

	for _, room := range rl.rooms {
		row := NewView()
		row.Init(row)
		row.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
		row.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
		row.GetElement().Yoga.StyleSetPadding(yoga.EdgeVertical, 10)
		row.GetElement().Yoga.StyleSetPadding(yoga.EdgeHorizontal, 12)
		row.GetElement().Yoga.StyleSetBorderRadius(6)
		row.GetElement().PointerEvents = core.PointerEventsAuto

		// 选中高亮
		if room.ID == rl.selectedID {
			row.GetElement().BackgroundColor = core.GetTheme().MenuItemHoverBg
		}

		// 房间名
		nameText := NewText(room.Name)
		nameText.SetFontSize(core.GetTheme().FontSizeBase)
		nameText.SetColor(core.GetTheme().TextColor)
		nameText.GetElement().Yoga.StyleSetFlexGrow(1)
		nameText.GetElement().PointerEvents = core.PointerEventsNone
		row.AddChild(nameText)

		// 密码标识
		if room.Password {
			lockText := NewText("🔒")
			lockText.SetFontSize(core.GetTheme().FontSizeSM)
			lockText.GetElement().PointerEvents = core.PointerEventsNone
			row.AddChild(lockText)
		}

		// 玩家数
		playerText := NewText(formatPlayers(room.Players, room.MaxPlayers))
		playerText.SetFontSize(core.GetTheme().FontSizeSM)
		playerText.SetColor(core.GetTheme().TextMutedColor)
		playerText.GetElement().Yoga.StyleSetMargin(yoga.EdgeLeft, 8)
		playerText.GetElement().PointerEvents = core.PointerEventsNone
		row.AddChild(playerText)

		// 规则
		if room.Rule != "" {
			ruleText := NewText(room.Rule)
			ruleText.SetFontSize(core.GetTheme().FontSizeSM)
			ruleText.SetColor(core.GetTheme().TextMutedColor)
			ruleText.GetElement().Yoga.StyleSetMargin(yoga.EdgeLeft, 8)
			ruleText.GetElement().PointerEvents = core.PointerEventsNone
			row.AddChild(ruleText)
		}

		// 悬停效果
		roomID := room.ID
		row.SetOnHover(func(hovered bool) {
			if roomID != rl.selectedID {
				if hovered {
					row.GetElement().BackgroundColor = core.GetTheme().MenuItemHoverBg
				} else {
					row.GetElement().BackgroundColor = nil
				}
			}
		})

		// 点击选择
		row.SetOnClick(func() {
			rl.selectedID = roomID
			if rl.onSelect != nil {
				rl.onSelect(roomID)
			}
			// 刷新所有行的高亮
			for _, item := range rl.itemHosts {
				item.GetElement().BackgroundColor = nil
			}
			row.GetElement().BackgroundColor = core.GetTheme().MenuItemHoverBg
		})

		rl.itemHosts[room.ID] = row
		rl.AddChild(row)
	}
}

func formatPlayers(current, max int) string {
	return itoa(current) + "/" + itoa(max)
}

// Draw 绘制列表背景。
func (rl *RoomList) Draw(screen *ebiten.Image) {
	el := rl.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := rl.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}
}

// ==================== 链式 API ====================

func (rl *RoomList) SetWidth(width float32) *RoomList {
	rl.GetElement().Yoga.StyleSetWidth(width)
	return rl
}
func (rl *RoomList) SetHeight(height float32) *RoomList {
	rl.GetElement().Yoga.StyleSetHeight(height)
	return rl
}
func (rl *RoomList) SetMargin(edge yoga.Edge, value float32) *RoomList {
	rl.GetElement().Yoga.StyleSetMargin(edge, value)
	return rl
}

// SyncFrom 同步房间列表属性。
func (rl *RoomList) SyncFrom(other core.Host) {
	if o, ok := other.(*RoomList); ok {
		rl.rooms = o.rooms
		rl.selectedID = o.selectedID
		rl.onSelect = o.onSelect
		rl.onJoin = o.onJoin
		rl.rebuild()
	}
}
