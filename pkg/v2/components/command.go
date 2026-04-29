package components

import (
	"strings"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// CommandItem represents a single command entry.
type CommandItem struct {
	Label    string
	Shortcut string
	Action   func()
}

// Command is a command palette (search + selectable list).
type Command struct {
	View
	items    []CommandItem
	filtered []CommandItem
	input    *TextInput
	list     *View
	isOpen   bool
}

// NewCommand creates a command palette.
func NewCommand() *Command {
	c := &Command{}
	c.Init(c)
	c.SetFlexDirection(yoga.FlexDirectionColumn)
	c.SetBackgroundColor(core.GetTheme().CardColor)
	c.SetBorderColor(core.GetTheme().BorderColor)
	c.SetBorderRadius(core.GetTheme().BorderRadius)
	c.SetPadding(yoga.EdgeAll, 16)
	c.SetShadow(core.GetTheme().ShadowColor, 8, 0, 4)

	c.input = NewTextInput()
	c.input.SetPlaceholder("Type a command or search...")
	c.input.SetOnChange(func(text string) { c.filter(text) })
	c.Add(c.input)

	c.list = NewView()
	c.list.SetFlexDirection(yoga.FlexDirectionColumn)
	c.list.SetMaxHeight(300)
	c.Add(c.list)
	return c
}

// ElementType returns type identifier.
func (c *Command) ElementType() string { return "Command" }

// SetItems sets the command list.
func (c *Command) SetItems(items []CommandItem) *Command {
	c.items = items
	c.filtered = append([]CommandItem{}, items...)
	c.renderList()
	return c
}

func (c *Command) filter(query string) {
	c.filtered = c.filtered[:0]
	q := strings.ToLower(query)
	for _, item := range c.items {
		if strings.Contains(strings.ToLower(item.Label), q) {
			c.filtered = append(c.filtered, item)
		}
	}
	c.renderList()
}

func (c *Command) renderList() {
	c.list.ClearChildren()
	for _, item := range c.filtered {
		it := item
		row := NewView()
		row.SetFlexDirection(yoga.FlexDirectionRow)
		row.SetJustifyContent(yoga.JustifySpaceBetween)
		row.SetAlignItems(yoga.AlignCenter)
		row.SetPadding(yoga.EdgeVertical, 8)
		row.SetPadding(yoga.EdgeHorizontal, 12)

		label := NewText(it.Label).SetFontSize(14).SetColor(core.GetTheme().TextColor)
		row.Add(label)
		if it.Shortcut != "" {
			shortcut := NewKbd(it.Shortcut)
			row.Add(shortcut)
		}
		c.list.Add(row)
	}
	c.Mark(core.FlagNeedLayout)
}

// Open shows the command palette.
func (c *Command) Open() *Command {
	c.isOpen = true
	c.SetVisible(true)
	c.Mark(core.FlagNeedLayout)
	return c
}

// Close hides the command palette.
func (c *Command) Close() *Command {
	c.isOpen = false
	c.SetVisible(false)
	return c
}
