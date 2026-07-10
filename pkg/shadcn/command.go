package shadcn

import (
	"strings"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

type CommandItem struct {
	Label    string
	OnSelect func()
}

type CommandProps struct {
	Open        bool
	OnClose     func()
	Items       []CommandItem
	Placeholder string
}

// Command 是命令面板：顶部搜索框 + 可过滤的命令列表（Portal 居中弹出）。
func Command(p CommandProps) *ui.Node { return ui.Use(command, p) }

func command(p CommandProps) *ui.Node {
	th := ui.UseTheme()
	mounted, prog := ui.UseTransition(p.Open, 160)
	query, setQuery := ui.UseState("")
	ui.UseEffect(func() ui.Cleanup {
		if p.Open {
			setQuery("")
		}
		return nil
	}, p.Open)
	ui.UseEscape(p.Open, func() {
		if p.OnClose != nil {
			p.OnClose()
		}
	})
	if !mounted {
		return nil
	}

	placeholder := p.Placeholder
	if placeholder == "" {
		placeholder = "输入命令搜索…"
	}

	rows := []*ui.Node{ui.Style(ui.Column)}
	for _, it := range p.Items {
		if query != "" && !strings.Contains(strings.ToLower(it.Label), strings.ToLower(query)) {
			continue
		}
		item := it
		rows = append(rows, menuRow(item.Label, func() {
			if item.OnSelect != nil {
				item.OnSelect()
			}
			if p.OnClose != nil {
				p.OnClose()
			}
		}))
	}

	box := []*ui.Node{
		ui.Style(ui.Column, ui.Width(460), ui.Bg(th.Popover), ui.TextColor(th.PopoverForeground),
			ui.Border(1, th.Border), ui.Radius(th.Radius+2), ui.Padding(6),
			ui.Scale(0.98+0.02*prog)),
		ui.OnClick(func() {}),
		ui.Input(
			ui.Style(ui.Height(40), ui.PaddingXY(10, 0), ui.Bg(th.Popover), ui.TextColor(th.Foreground),
				ui.FontSize(15), ui.ItemsCenter),
			ui.Value(query), ui.OnChange(setQuery), ui.Placeholder(placeholder),
		),
		ui.Div(ui.Style(ui.Height(1), ui.Bg(th.Border))),
		ui.ScrollView(append([]*ui.Node{ui.Style(ui.Column, ui.MaxHeight(300), ui.PaddingXY(0, 4))}, rows...)...),
	}

	return ui.Portal(
		ui.Div(
			ui.Style(ui.Grow(1), ui.Column, ui.ItemsCenter, ui.Bg(ui.Color{R: 0, G: 0, B: 0, A: uint8(120 * prog)})),
			ui.OnClick(func() {
				if p.OnClose != nil {
					p.OnClose()
				}
			}),
			ui.Div(ui.Style(ui.Height(90))), // 顶部留白
			ui.Div(box...),
		),
	)
}
