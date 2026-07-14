package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type AccordionItemData struct {
	Title   string
	Content []*ui.Node
}

type accordionProps struct{ items []AccordionItemData }

// Accordion 是可折叠面板组（shadcn 单开模式）：同一时刻至多展开一项，默认首项展开，
// 展开/收起带高度滑动动画、chevron 旋转，与 shadcn/ui 一致。
func Accordion(items []AccordionItemData) *ui.Node {
	return ui.Use(accordion, accordionProps{items: items})
}

func accordion(p accordionProps) *ui.Node {
	openIdx, setOpen := ui.UseState(0) // 默认首项展开
	kids := []*ui.Node{ui.Style(ui.Column)}
	for i, it := range p.items {
		idx := i
		kids = append(kids, ui.Use(accordionItem, accordionItemProps{
			title: it.Title, content: it.Content, open: openIdx == idx, last: i == len(p.items)-1,
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
	hovered, _, ia := ui.UseInteraction()
	cref, crect := ui.UseMeasure() // 测量内容自然高度

	target := float32(0)
	if p.open {
		target = crect.H
	}
	h := ui.UseTween(target, 200, ui.EaseOut) // 高度滑动，ease-out 与 shadcn 一致
	rot := float32(0)
	if p.open {
		rot = 180
	}
	chevRot := ui.UseTween(rot, 200, ui.EaseOut) // chevron 旋转动画

	titleColor := th.Foreground
	if hovered {
		titleColor = over(th.Foreground, th.MutedForeground, 0.4) // 悬停轻微变化（替代 text-decoration:underline）
	}

	// 触发头：justify-between，medium 14px，py-4，右侧旋转 chevron
	header := ui.Div(
		ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.PaddingXY(0, 16)),
		ui.OnClick(p.onToggle), ia,
		ui.Text(p.title, ui.FontSize(14), ui.Medium, ui.TextColor(titleColor)),
		ui.Div(ui.Style(ui.Rotate(chevRot)),
			ui.Icon(ui.IconChevronDown, 16, ui.TextColor(th.MutedForeground))),
	)

	// 外层裁剪到动画高度；内层绝对定位，按自然高度布局并被测量（不被动画高度挤压）。
	body := ui.Div(
		ui.Style(ui.Height(h), ui.Clip),
		ui.Div(cref, ui.Style(ui.Absolute, ui.Left(0), ui.Top(0), ui.WidthPct(100), ui.Column),
			ui.Div(append([]*ui.Node{ui.Style(ui.Column, ui.Gap(6),
				ui.FontSize(14), ui.TextColor(th.Foreground))}, p.content...)...),
			ui.Div(ui.Style(ui.Height(16))), // pb-4
		),
	)

	item := ui.Div(ui.Style(ui.Column), header, body)
	if !p.last { // 项间下边框（border-b）
		return ui.Div(ui.Style(ui.Column), item, ui.Div(ui.Style(ui.Height(1), ui.Bg(th.Border))))
	}
	return item
}
