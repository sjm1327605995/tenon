package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ChatMessage 定义一条聊天消息。
type ChatMessage struct {
	Sender  string
	Content string
	System  bool // 系统消息
}

// ChatPanel 是聊天面板组件。
type ChatPanel struct {
	core.BaseHost
	messages     []ChatMessage
	scrollView   *ScrollView
	messageViews []*View
	maxMessages  int
}

// NewChatPanel 创建一个聊天面板。
func NewChatPanel() *ChatPanel {
	cp := &ChatPanel{
		maxMessages: 100,
	}
	cp.Init(cp)
	cp.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	cp.GetElement().Yoga.StyleSetWidth(280)
	cp.GetElement().Yoga.StyleSetHeight(200)
	cp.GetElement().BackgroundColor = color.RGBA{R: 30, G: 30, B: 35, A: 230}
	cp.GetElement().BorderRadius = core.BorderRadius{TopLeft: 6, TopRight: 6, BottomRight: 6, BottomLeft: 6}

	scroll := NewScrollView()
	scroll.GetElement().Yoga.StyleSetFlexGrow(1)
	scroll.GetElement().Yoga.StyleSetWidthPercent(100)
	cp.scrollView = scroll
	cp.AddChild(scroll)

	return cp
}

// AddMessage 添加一条消息。
func (cp *ChatPanel) AddMessage(msg ChatMessage) *ChatPanel {
	cp.messages = append(cp.messages, msg)
	if len(cp.messages) > cp.maxMessages {
		cp.messages = cp.messages[len(cp.messages)-cp.maxMessages:]
	}
	cp.rebuildMessages()
	return cp
}

// ClearMessages 清空所有消息。
func (cp *ChatPanel) ClearMessages() *ChatPanel {
	cp.messages = nil
	cp.rebuildMessages()
	return cp
}

// SetMaxMessages 设置最大消息数。
func (cp *ChatPanel) SetMaxMessages(max int) *ChatPanel {
	cp.maxMessages = max
	if len(cp.messages) > max {
		cp.messages = cp.messages[len(cp.messages)-max:]
		cp.rebuildMessages()
	}
	return cp
}

func (cp *ChatPanel) rebuildMessages() {
	if cp.scrollView == nil {
		return
	}
	cp.scrollView.ClearChildren()
	cp.messageViews = nil

	for _, msg := range cp.messages {
		row := NewView()
		row.Init(row)
		row.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
		row.GetElement().Yoga.StyleSetPadding(yoga.EdgeVertical, 2)
		row.GetElement().Yoga.StyleSetPadding(yoga.EdgeHorizontal, 8)
		row.GetElement().PointerEvents = core.PointerEventsNone

		if msg.System {
			// 系统消息居中灰色
			sysText := NewText(msg.Content)
			sysText.SetFontSize(core.GetTheme().FontSizeSM)
			sysText.SetColor(core.GetTheme().TextMutedColor)
			sysText.GetElement().PointerEvents = core.PointerEventsNone
			row.AddChild(sysText)
		} else {
			// 普通消息：发送者名字 + 内容
			senderText := NewText(msg.Sender + ": ")
			senderText.SetFontSize(core.GetTheme().FontSizeSM)
			senderText.SetColor(core.GetTheme().PrimaryColor)
			senderText.GetElement().PointerEvents = core.PointerEventsNone
			row.AddChild(senderText)

			contentText := NewText(msg.Content)
			contentText.SetFontSize(core.GetTheme().FontSizeSM)
			contentText.SetColor(core.GetTheme().TextColor)
			contentText.GetElement().PointerEvents = core.PointerEventsNone
			row.AddChild(contentText)
		}

		cp.messageViews = append(cp.messageViews, row)
		cp.scrollView.AddChild(row)
	}

	// 滚动到底部
	if cp.GetEngine() != nil {
		cp.GetEngine().InvalidateAll()
	}
}

// HandleEvent 消费所有事件防止穿透。
func (cp *ChatPanel) HandleEvent(e *core.Event) bool {
	return cp.GetElement().Visible
}

// ==================== 链式 API ====================

func (cp *ChatPanel) SetWidth(width float32) *ChatPanel {
	cp.GetElement().Yoga.StyleSetWidth(width)
	return cp
}
func (cp *ChatPanel) SetHeight(height float32) *ChatPanel {
	cp.GetElement().Yoga.StyleSetHeight(height)
	return cp
}
func (cp *ChatPanel) SetMargin(edge yoga.Edge, value float32) *ChatPanel {
	cp.GetElement().Yoga.StyleSetMargin(edge, value)
	return cp
}

// SyncFrom 同步聊天面板属性。
func (cp *ChatPanel) SyncFrom(other core.Host) {
	if o, ok := other.(*ChatPanel); ok {
		cp.messages = o.messages
		cp.maxMessages = o.maxMessages
		cp.rebuildMessages()
	}
}
