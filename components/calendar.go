package components

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

const calendarMaxWeeks = 6

// Calendar is a month calendar grid.
type Calendar struct {
	core.BaseElement
	month     time.Month
	year      int
	selected  time.Time
	onSelect  func(date time.Time)
	title     *native.Text
	prevBtn   *Button
	nextBtn   *Button
	weekRows  []*native.View
	dayCells  [][]*Button // [week][day]
}

// NewCalendar creates a calendar for the current month.
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
	c.initGrid()
	c.updateGrid()
	return c
}

func (c *Calendar) initGrid() {
	theme := core.GetTheme()

	header := native.NewView()
	header.SetFlexDirection(yoga.FlexDirectionRow)
	header.SetJustifyContent(yoga.JustifySpaceBetween)
	header.SetAlignItems(yoga.AlignCenter)
	header.SetPadding(yoga.EdgeHorizontal, 4)

	c.prevBtn = NewButton("←").SetVariant(ButtonGhost)
	c.prevBtn.SetOnClick(func() { c.prevMonth() })
	c.title = native.NewText("").SetFontSize(16)
	c.nextBtn = NewButton("→").SetVariant(ButtonGhost)
	c.nextBtn.SetOnClick(func() { c.nextMonth() })
	header.Add(c.prevBtn, c.title, c.nextBtn)
	c.Add(header)

	cellFontSize := float32(12)
	cellH := cellFontSize * 2.5
	cellW := float32(40)

	daysRow := native.NewView()
	daysRow.SetFlexDirection(yoga.FlexDirectionRow)
	for _, d := range []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"} {
		cell := native.NewView().SetWidth(cellW).SetHeight(cellH)
		cell.SetJustifyContent(yoga.JustifyCenter)
		cell.SetAlignItems(yoga.AlignCenter)
		cell.SetFlexShrink(0)
		label := native.NewText(d).SetFontSize(float64(cellFontSize)).SetColor(theme.MutedForegroundColor)
		cell.Add(label)
		daysRow.Add(cell)
	}
	c.Add(daysRow)

	c.weekRows = make([]*native.View, calendarMaxWeeks)
	c.dayCells = make([][]*Button, calendarMaxWeeks)
	for w := 0; w < calendarMaxWeeks; w++ {
		row := native.NewView()
		row.SetFlexDirection(yoga.FlexDirectionRow)
		c.weekRows[w] = row
		c.dayCells[w] = make([]*Button, 7)
		for d := 0; d < 7; d++ {
			btn := NewButton("").SetVariant(ButtonGhost)
			btn.SetWidth(cellW)
			btn.SetHeight(cellH)
			btn.SetFlexShrink(0)
			btn.SetJustifyContent(yoga.JustifyCenter)
			btn.SetAlignItems(yoga.AlignCenter)
			btn.SetPadding(yoga.EdgeAll, 0)
			c.dayCells[w][d] = btn
			row.Add(btn)
		}
		c.Add(row)
	}
}

func (c *Calendar) updateGrid() {
	c.title.SetContent(fmt.Sprintf("%s %d", c.month.String(), c.year))

	firstDay := time.Date(c.year, c.month, 1, 0, 0, 0, 0, time.Local)
	startOffset := int(firstDay.Weekday())
	daysInMonth := 32 - time.Date(c.year, c.month, 32, 0, 0, 0, 0, time.Local).Day()

	for w := 0; w < calendarMaxWeeks; w++ {
		rowVisible := false
		for d := 0; d < 7; d++ {
			btn := c.dayCells[w][d]
			cellDay := w*7 + d - startOffset + 1
			if cellDay >= 1 && cellDay <= daysInMonth {
				rowVisible = true
				btn.SetDisplay(yoga.DisplayFlex)
				btn.SetText(fmt.Sprintf("%d", cellDay))
				if cellDay == c.selected.Day() && c.month == c.selected.Month() && c.year == c.selected.Year() {
					btn.SetVariant(ButtonDefault)
				} else {
					btn.SetVariant(ButtonGhost)
				}
				idx := cellDay
				btn.SetOnClick(func() {
					c.selected = time.Date(c.year, c.month, idx, 0, 0, 0, 0, time.Local)
					if c.onSelect != nil {
						c.onSelect(c.selected)
					}
					c.updateGrid()
				})
			} else {
				btn.SetDisplay(yoga.DisplayNone)
			}
		}
		if rowVisible {
			c.weekRows[w].SetDisplay(yoga.DisplayFlex)
		} else {
			c.weekRows[w].SetDisplay(yoga.DisplayNone)
		}
	}

	c.Mark(core.FlagNeedLayout)
}

func (c *Calendar) prevMonth() {
	c.month--
	if c.month < 1 {
		c.month = 12
		c.year--
	}
	c.updateGrid()
}

func (c *Calendar) nextMonth() {
	c.month++
	if c.month > 12 {
		c.month = 1
		c.year++
	}
	c.updateGrid()
}

// SetDate sets the selected date and updates the displayed month.
func (c *Calendar) SetDate(date time.Time) *Calendar {
	c.selected = date
	c.month = date.Month()
	c.year = date.Year()
	c.updateGrid()
	return c
}

// SetOnSelect sets the date selection callback.
func (c *Calendar) SetOnSelect(fn func(date time.Time)) *Calendar {
	c.onSelect = fn
	return c
}
