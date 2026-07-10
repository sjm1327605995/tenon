package shadcn

import (
	"fmt"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// ---- Breadcrumb ----

type breadcrumbProps struct {
	items    []string
	onSelect func(int)
}

// Breadcrumb 是面包屑导航（末项为当前页）。
func Breadcrumb(items []string, onSelect func(int)) *ui.Node {
	return ui.Use(breadcrumb, breadcrumbProps{items: items, onSelect: onSelect})
}

func breadcrumb(p breadcrumbProps) *ui.Node {
	th := ui.UseTheme()
	kids := []*ui.Node{ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(8))}
	for i, it := range p.items {
		idx := i
		last := i == len(p.items)-1
		col := th.MutedForeground
		if last {
			col = th.Foreground
		}
		kids = append(kids, ui.Div(
			ui.OnClick(func() {
				if p.onSelect != nil {
					p.onSelect(idx)
				}
			}),
			ui.Text(it, ui.FontSize(14), ui.TextColor(col)),
		))
		if !last {
			kids = append(kids, ui.Text("/", ui.FontSize(14), ui.TextColor(th.MutedForeground)))
		}
	}
	return ui.Div(kids...)
}

// ---- Pagination ----

type PaginationProps struct {
	Page, Total int
	OnChange    func(int)
}

// Pagination 是分页器（1..Total，带上一页/下一页）。
func Pagination(p PaginationProps) *ui.Node { return ui.Use(pagination, p) }

func pagination(p PaginationProps) *ui.Node {
	th := ui.UseTheme()
	cell := func(label string, target int, disabled, active bool) *ui.Node {
		st := []ui.StyleOpt{ui.Width(34), ui.Height(34), ui.Radius(th.Radius),
			ui.ItemsCenter, ui.JustifyCenter, ui.Border(1, th.Border)}
		fg := th.Foreground
		if active {
			st = append(st, ui.Bg(th.Primary))
			fg = th.PrimaryForeground
		}
		if disabled {
			st = append(st, ui.Opacity(0.5))
		}
		attrs := []*ui.Node{ui.Style(st...)}
		if !disabled {
			attrs = append(attrs, ui.OnClick(func() {
				if p.OnChange != nil {
					p.OnChange(target)
				}
			}))
		}
		return ui.Div(append(attrs, ui.Text(label, ui.FontSize(13), ui.TextColor(fg)))...)
	}
	kids := []*ui.Node{ui.Style(ui.Row, ui.Gap(6), ui.ItemsCenter)}
	kids = append(kids, cell("‹", p.Page-1, p.Page <= 1, false))
	for i := 1; i <= p.Total; i++ {
		kids = append(kids, cell(fmt.Sprintf("%d", i), i, false, i == p.Page))
	}
	kids = append(kids, cell("›", p.Page+1, p.Page >= p.Total, false))
	return ui.Div(kids...)
}

// ---- AspectRatio ----

// AspectRatio 按给定宽度与比例（宽/高）裁剪出固定纵横比的容器。
func AspectRatio(width, ratio float32, children ...*ui.Node) *ui.Node {
	h := width
	if ratio > 0 {
		h = width / ratio
	}
	base := ui.Style(ui.Width(width), ui.Height(h), ui.Clip)
	return ui.Div(append([]*ui.Node{base}, children...)...)
}

// ---- ScrollArea ----

// ScrollArea 是固定高度的可滚动区域。
func ScrollArea(height float32, children ...*ui.Node) *ui.Node {
	base := ui.Style(ui.Height(height), ui.Column)
	return ui.ScrollView(append([]*ui.Node{base}, children...)...)
}

// ---- Collapsible ----

type collapsibleProps struct {
	open    bool
	trigger *ui.Node
	content []*ui.Node
}

// Collapsible 是单个可折叠区块（受控 open），带高度动画。
func Collapsible(open bool, onToggle func(), trigger *ui.Node, content ...*ui.Node) *ui.Node {
	// 把 onToggle 挂到 trigger 外层，内部用组件跑动画
	wrapped := ui.Div(ui.OnClick(onToggle), trigger)
	return ui.Use(collapsible, collapsibleProps{open: open, trigger: wrapped, content: content})
}

func collapsible(p collapsibleProps) *ui.Node {
	cref, crect := ui.UseMeasure()
	target := float32(0)
	if p.open {
		target = crect.H
	}
	h := ui.UseTween(target, 180, ui.EaseInOut)
	return ui.Div(ui.Style(ui.Column, ui.Gap(6)),
		p.trigger,
		ui.Div(ui.Style(ui.Height(h), ui.Clip),
			ui.Div(append([]*ui.Node{cref, ui.Style(ui.Column, ui.Gap(6))}, p.content...)...),
		),
	)
}
