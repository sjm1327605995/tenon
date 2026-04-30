package components

import (
	"fmt"

	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Pagination is a page navigation control.
type Pagination struct {
	core.BaseElement
	totalPages  int
	current     int
	onChange    func(page int)
	prevBtn     *Button
	nextBtn     *Button
	pageButtons []*Button
}

// NewPagination creates pagination controls.
func NewPagination(totalPages int) *Pagination {
	p := &Pagination{totalPages: totalPages, current: 1}
	p.Init(p)
	p.SetFlexDirection(yoga.FlexDirectionRow)
	p.SetAlignItems(yoga.AlignCenter)
	p.SetGap(yoga.GutterAll, 4)
	p.initButtons()
	p.updateButtons()
	return p
}

func (p *Pagination) initButtons() {
	p.prevBtn = NewButton("←").SetVariant(ButtonOutline)
	p.prevBtn.SetOnClick(func() { p.SetPage(p.current - 1) })
	p.AppendChild(p.prevBtn)

	maxVisible := 5
	if p.totalPages < maxVisible {
		maxVisible = p.totalPages
	}
	p.pageButtons = make([]*Button, maxVisible)
	for i := 0; i < maxVisible; i++ {
		btn := NewButton("")
		p.pageButtons[i] = btn
		p.AppendChild(btn)
	}

	p.nextBtn = NewButton("→").SetVariant(ButtonOutline)
	p.nextBtn.SetOnClick(func() { p.SetPage(p.current + 1) })
	p.AppendChild(p.nextBtn)
}

func (p *Pagination) updateButtons() {
	start, end := p.visibleRange()

	p.prevBtn.SetDisabled(p.current <= 1)

	for i, btn := range p.pageButtons {
		page := start + i
		if page > end {
			btn.SetDisplay(yoga.DisplayNone)
		} else {
			btn.SetDisplay(yoga.DisplayFlex)
			btn.SetText(fmt.Sprintf("%d", page))
			if page == p.current {
				btn.SetVariant(ButtonDefault)
			} else {
				btn.SetVariant(ButtonOutline)
			}
			idx := page
			btn.SetOnClick(func() { p.SetPage(idx) })
		}
	}

	p.nextBtn.SetDisabled(p.current >= p.totalPages)
	p.Mark(core.FlagNeedLayout)
}

func (p *Pagination) visibleRange() (int, int) {
	start := p.current - 2
	if start < 1 {
		start = 1
	}
	end := start + 4
	if end > p.totalPages {
		end = p.totalPages
		start = end - 4
		if start < 1 {
			start = 1
		}
	}
	return start, end
}

// SetPage sets the current page.
func (p *Pagination) SetPage(page int) *Pagination {
	if page < 1 || page > p.totalPages {
		return p
	}
	p.current = page
	p.updateButtons()
	if p.onChange != nil {
		p.onChange(page)
	}
	return p
}

// SetOnChange sets the page change callback.
func (p *Pagination) SetOnChange(fn func(page int)) *Pagination {
	p.onChange = fn
	return p
}
