package components

import (
	"fmt"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Pagination is a page navigation control.
type Pagination struct {
	core.BaseElement
	totalPages int
	current    int
	onChange   func(page int)
}

// NewPagination creates pagination controls.
func NewPagination(totalPages int) *Pagination {
	p := &Pagination{
		totalPages: totalPages,
		current:    1,
	}
	p.Init(p)
	p.SetFlexDirection(yoga.FlexDirectionRow)
	p.SetAlignItems(yoga.AlignCenter)
	p.SetGap(yoga.GutterAll, 4)
	p.buildButtons()
	return p
}

// ElementType returns type identifier.
func (p *Pagination) ElementType() string { return "Pagination" }

func (p *Pagination) buildButtons() {
	p.ClearChildren()

	prev := NewButton("←").SetVariant(ButtonOutline)
	prev.SetOnClick(func() { p.SetPage(p.current - 1) })
	if p.current <= 1 {
		prev.SetDisabled(true)
	}
	p.AppendChild(prev)

	start, end := p.visibleRange()
	for i := start; i <= end; i++ {
		idx := i
		btn := NewButton(fmt.Sprintf("%d", i))
		if i == p.current {
			btn.SetVariant(ButtonDefault)
		} else {
			btn.SetVariant(ButtonOutline)
		}
		btn.SetOnClick(func() { p.SetPage(idx) })
		p.AppendChild(btn)
	}

	next := NewButton("→").SetVariant(ButtonOutline)
	next.SetOnClick(func() { p.SetPage(p.current + 1) })
	if p.current >= p.totalPages {
		next.SetDisabled(true)
	}
	p.AppendChild(next)
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
	p.buildButtons()
	p.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
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
