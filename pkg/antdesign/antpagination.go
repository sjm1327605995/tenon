package antdesign

import (
	"fmt"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntPagination is a pagination component.
type AntPagination struct {
	tenon.BaseWidget
	current          int
	total            int
	pageSize         int
	showSizeChanger  bool
	showQuickJumper  bool
	disabled         bool
	onChange         func(page, pageSize int)
	onShowSizeChange func(current, size int)
}

// NewAntPagination creates an AntPagination.
func NewAntPagination() *AntPagination {
	p := &AntPagination{pageSize: 10}
	p.Init(p)
	return p
}

// Render returns the pagination UI.
func (p *AntPagination) Render() tenon.Component {
	theme := NewAntTheme()

	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignCenter).
		SetFlexWrap(yoga.WrapWrap)

	totalPages := (p.total + p.pageSize - 1) / p.pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	// Prev button
	prevDisabled := p.current <= 1 || p.disabled
	prevBtn := components.NewButton("<").
		SetOnClick(func() {
			if !prevDisabled && p.current > 1 {
				p.current--
				if p.onChange != nil {
					p.onChange(p.current, p.pageSize)
				}
			}
		})
	if prevDisabled {
		prevBtn.SetDisabled(true)
	}
	root.AddChild(prevBtn)

	// Page numbers (simplified: show all or ellipsis logic)
	for i := 1; i <= totalPages; i++ {
		isActive := i == p.current
		btn := components.NewButton(fmt.Sprintf("%d", i)).
			SetOnClick(func(page int) func() {
				return func() {
					if !p.disabled && p.current != page {
						p.current = page
						if p.onChange != nil {
							p.onChange(p.current, p.pageSize)
						}
					}
				}
			}(i))
		if isActive {
			btn.SetBackgroundColors(theme.PrimaryColor, theme.PrimaryHoverColor, theme.ButtonPressedColor)
			btn.SetTextColor(theme.ButtonTextColor)
		}
		if p.disabled {
			btn.SetDisabled(true)
		}
		root.AddChild(btn)
	}

	// Next button
	nextDisabled := p.current >= totalPages || p.disabled
	nextBtn := components.NewButton(">").
		SetOnClick(func() {
			if !nextDisabled && p.current < totalPages {
				p.current++
				if p.onChange != nil {
					p.onChange(p.current, p.pageSize)
				}
			}
		})
	if nextDisabled {
		nextBtn.SetDisabled(true)
	}
	root.AddChild(nextBtn)

	// Show total info
	if p.total > 0 {
		root.Add(components.NewText(fmt.Sprintf("Total %d items", p.total)).
			SetFontSize(theme.FontSizeBase).
			SetColor(theme.TextMutedColor).
			SetMargin(yoga.EdgeLeft, 16))
	}

	return root
}

func (p *AntPagination) SetCurrent(c int) *AntPagination          { p.current = c; return p }
func (p *AntPagination) SetTotal(t int) *AntPagination            { p.total = t; return p }
func (p *AntPagination) SetPageSize(s int) *AntPagination         { p.pageSize = s; return p }
func (p *AntPagination) SetShowSizeChanger(v bool) *AntPagination { p.showSizeChanger = v; return p }
func (p *AntPagination) SetShowQuickJumper(v bool) *AntPagination { p.showQuickJumper = v; return p }
func (p *AntPagination) SetDisabled(v bool) *AntPagination        { p.disabled = v; return p }
func (p *AntPagination) SetOnChange(fn func(page, pageSize int)) *AntPagination {
	p.onChange = fn
	return p
}
