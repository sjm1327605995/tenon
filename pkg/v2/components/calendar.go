package components

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

type Calendar struct {
	core.BaseElement
	month    time.Month
	year     int
	selected time.Time
	onSelect func(date time.Time)
}

func NewCalendar() *Calendar {
	now := time.Now()
	c := &Calendar{
		month:    now.Month(),
		year:     now.Year(),
		selected: now,
	}
	c.Init(c)
	c.SetFlexDirection(yoga.FlexDirectionColumn)
	c.SetGap(yoga.GutterAll, 4)
	c.SetWidth(280)
	c.buildGrid()
	return c
}

func (c *Calendar) ElementType() string { return "Calendar" }

func (c *Calendar) buildGrid() {
	c.ClearChildren()

	theme := core.GetTheme()
	cellFontSize := float32(12)
	cellH := cellFontSize * 2.5

	header := NewView()
	header.SetFlexDirection(yoga.FlexDirectionRow)
	header.SetJustifyContent(yoga.JustifySpaceBetween)
	header.SetAlignItems(yoga.AlignCenter)
	header.SetPadding(yoga.EdgeHorizontal, 4)

	prev := NewButton("←").SetVariant(ButtonGhost)
	prev.SetOnClick(func() { c.prevMonth() })
	title := NewText(fmt.Sprintf("%s %d", c.month.String(), c.year)).SetFontSize(16)
	next := NewButton("→").SetVariant(ButtonGhost)
	next.SetOnClick(func() { c.nextMonth() })
	header.Add(prev, title, next)
	c.Add(header)

	// 固定格子宽度：280 / 7 = 40，刚好整除，无任何舍入误差
	cellW := float32(40)

	daysRow := NewView()
	daysRow.SetFlexDirection(yoga.FlexDirectionRow)
	for _, d := range []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"} {
		cell := NewView().SetWidth(cellW).SetHeight(cellH)
		cell.SetJustifyContent(yoga.JustifyCenter)
		cell.SetAlignItems(yoga.AlignCenter)
		cell.SetFlexShrink(0)
		label := NewText(d).SetFontSize(float64(cellFontSize)).SetColor(theme.MutedForegroundColor)
		cell.Add(label)
		daysRow.Add(cell)
	}
	c.Add(daysRow)

	firstDay := time.Date(c.year, c.month, 1, 0, 0, 0, 0, time.Local)
	startOffset := int(firstDay.Weekday())
	daysInMonth := 32 - time.Date(c.year, c.month, 32, 0, 0, 0, 0, time.Local).Day()

	week := NewView()
	week.SetFlexDirection(yoga.FlexDirectionRow)

	for i := 0; i < startOffset; i++ {
		empty := NewView().SetWidth(cellW).SetHeight(cellH)
		empty.SetFlexShrink(0)
		week.Add(empty)
	}

	for day := 1; day <= daysInMonth; day++ {
		d := day
		cell := NewButton(fmt.Sprintf("%d", d)).SetVariant(ButtonGhost)
		cell.SetWidth(cellW)
		cell.SetHeight(cellH)
		cell.SetFlexShrink(0)
		cell.SetJustifyContent(yoga.JustifyCenter)
		cell.SetAlignItems(yoga.AlignCenter)
		cell.SetPadding(yoga.EdgeAll, 0) // 覆盖 Button 默认 16px padding，让文字在 40px 内居中
		cell.SetOnClick(func() {
			c.selected = time.Date(c.year, c.month, d, 0, 0, 0, 0, time.Local)
			if c.onSelect != nil {
				c.onSelect(c.selected)
			}
			c.buildGrid()
		})
		if d == c.selected.Day() && c.month == c.selected.Month() && c.year == c.selected.Year() {
			cell.SetVariant(ButtonDefault)
		}
		week.Add(cell)
		if (startOffset+day)%7 == 0 && day != daysInMonth {
			c.Add(week)
			week = NewView()
			week.SetFlexDirection(yoga.FlexDirectionRow)
		}
	}
	c.Add(week)

	c.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
}

func (c *Calendar) prevMonth() {
	c.month--
	if c.month < 1 {
		c.month = 12
		c.year--
	}
	c.buildGrid()
}

func (c *Calendar) nextMonth() {
	c.month++
	if c.month > 12 {
		c.month = 1
		c.year++
	}
	c.buildGrid()
}

func (c *Calendar) SetDate(date time.Time) *Calendar {
	c.selected = date
	c.month = date.Month()
	c.year = date.Year()
	c.buildGrid()
	return c
}

func (c *Calendar) SetOnSelect(fn func(date time.Time)) *Calendar {
	c.onSelect = fn
	return c
}
