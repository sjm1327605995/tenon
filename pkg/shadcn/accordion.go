package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type AccordionItemData struct {
	Title   string
	Content []*ui.Node
}

type accordionProps struct{ items []AccordionItemData }

// Accordion 是可折叠面板组，同一时刻至多展开一项，展开/收起带高度动画。
func Accordion(items []AccordionItemData) *ui.Node {
	return ui.Use(accordion, accordionProps{items: items})
}

func accordion(p accordionProps) *ui.Node {
	th := ui.UseTheme()
	openIdx, setOpen := ui.UseState(-1)
	kids := []*ui.Node{ui.Style(ui.Column, ui.Border(1, th.Border), ui.Radius(th.Radius), ui.Clip)}
	for i, it := range p.items {
		idx := i
		open := openIdx == idx
		kids = append(kids, ui.Use(accordionItem, accordionItemProps{
			title: it.Title, content: it.Content, open: open, last: i == len(p.items)-1,
			onToggle: func() {
				if openIdx == idx {
					setOpen(-1)
				} else {
					setOpen(idx)
				}
			},
		}))
	}
	return ui.Div(kids...)
}

type accordionItemProps struct {
	title    string
	content  []*ui.Node
	open     bool
	last     bool
	onToggle func()
}

func accordionItem(p accordionItemProps) *ui.Node {
	th := ui.UseTheme()
	cref, crect := ui.UseMeasure() // 测量内容自然高度

	target := float32(0)
	if p.open {
		target = crect.H
	}
	h := ui.UseTween(target, 200, ui.EaseInOut)

	chevron := "▸"
	if p.open {
		chevron = "▾"
	}

	header := ui.Div(
		ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.PaddingXY(16, 14)),
		ui.OnClick(p.onToggle),
		ui.Text(p.title, ui.FontSize(15)),
		ui.Text(chevron, ui.FontSize(12), ui.TextColor(th.MutedForeground)),
	)

	// 外层裁剪到动画高度；内层始终按自然高度布局并被测量
	body := ui.Div(
		ui.Style(ui.Height(h), ui.Clip),
		ui.Div(cref, ui.Style(ui.Column, ui.PaddingXY(16, 0), ui.Gap(6)),
			ui.Div(ui.Style(ui.Height(2))),
			ui.Div(append([]*ui.Node{ui.Style(ui.Column, ui.Gap(6))}, p.content...)...),
			ui.Div(ui.Style(ui.Height(14))),
		),
	)

	item := ui.Div(ui.Style(ui.Column), header, body)
	if !p.last {
		return ui.Div(ui.Style(ui.Column), item, ui.Div(ui.Style(ui.Height(1), ui.Bg(th.Border))))
	}
	return item
}
