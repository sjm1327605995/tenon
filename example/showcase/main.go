package main

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func App(_ struct{}) *ui.Node {
	th := ui.DarkTheme
	card := func(opts []ui.StyleOpt, kids ...*ui.Node) *ui.Node {
		base := []ui.StyleOpt{ui.Column, ui.Gap(10), ui.Padding(16), ui.Bg(th.Card),
			ui.Border(1, th.Border), ui.Radius(th.Radius + 4)}
		return ui.Div(append([]*ui.Node{ui.Style(append(base, opts...)...)}, kids...)...)
	}
	stat := func(label, value, delta string, up bool) *ui.Node {
		dc := ui.Hex("#22c55e")
		if !up {
			dc = ui.Hex("#ef4444")
		}
		return card([]ui.StyleOpt{ui.Grow(1), ui.Gap(4)},
			shadcn.Muted(label),
			ui.Text(value, ui.FontSize(26), ui.Bold),
			ui.Text(delta, ui.FontSize(13), ui.TextColor(dc)))
	}
	dot := func(c ui.Color, label string) *ui.Node {
		return ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(8)),
			ui.Div(ui.Style(ui.Width(10), ui.Height(10), ui.Radius(5), ui.Bg(c))),
			ui.Text(label, ui.FontSize(13), ui.TextColor(th.MutedForeground)))
	}

	sidebar := shadcn.Sidebar(shadcn.SidebarProps{
		Header: shadcn.H4("◆ Tenon"),
		Groups: []shadcn.SidebarGroup{
			{Label: "概览", Items: []shadcn.SidebarItem{
				{Label: "仪表盘", Icon: ui.Icon(ui.IconSearch, 18), Active: true},
				{Label: "分析", Icon: ui.Icon(ui.IconArrowRight, 18)},
				{Label: "新建", Icon: ui.Icon(ui.IconPlus, 18)},
			}},
			{Label: "系统", Items: []shadcn.SidebarItem{
				{Label: "回收站", Icon: ui.Icon(ui.IconTrash, 18)},
			}},
		},
		Footer: shadcn.Muted("v1.0 · dark"),
	})

	topbar := ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(12)),
		shadcn.H3("仪表盘"),
		ui.Div(ui.Style(ui.Grow(1))),
		ui.Div(ui.Style(ui.Width(240)),
			shadcn.InputGroup(shadcn.InputGroupProps{Leading: ui.Icon(ui.IconSearch, 16),
				Placeholder: "搜索…", Trailing: shadcn.KbdGroup("Ctrl", "K")})),
		shadcn.Badge(shadcn.BadgeProps{}, ui.Text("Pro")),
		shadcn.Button(shadcn.ButtonProps{}, ui.Text("新建报表")),
	)

	stats := ui.Div(ui.Style(ui.Row, ui.Gap(16)),
		stat("总收入", "¥128,430", "▲ 12.5%  较上月", true),
		stat("活跃用户", "8,942", "▲ 4.1%  较上周", true),
		stat("转化率", "3.24%", "▼ 0.8%  较上周", false),
		stat("退款", "¥2,110", "▲ 1.2%  较上月", false),
	)

	series := []float32{18, 24, 20, 32, 28, 40, 36, 48, 44, 52}
	charts := ui.Div(ui.Style(ui.Row, ui.Gap(16)),
		card([]ui.StyleOpt{ui.Grow(1)},
			shadcn.Small("收入趋势"),
			shadcn.LineChart(shadcn.LineChartProps{Data: series, Width: 480, Height: 170, Area: true})),
		card([]ui.StyleOpt{ui.Width(240), ui.ItemsCenter},
			shadcn.Small("渠道占比"),
			shadcn.PieChart(shadcn.PieChartProps{Size: 150, Slices: []shadcn.PieSlice{
				{Value: 45, Color: ui.Hex("#6366f1")}, {Value: 30, Color: ui.Hex("#ec4899")},
				{Value: 15, Color: ui.Hex("#22c55e")}, {Value: 10, Color: ui.Hex("#f59e0b")}}}),
			ui.Div(ui.Style(ui.Column, ui.Gap(4)),
				dot(ui.Hex("#6366f1"), "直接 45%"), dot(ui.Hex("#ec4899"), "搜索 30%"),
				dot(ui.Hex("#22c55e"), "推荐 15%"), dot(ui.Hex("#f59e0b"), "其它 10%"))),
	)

	table := card([]ui.StyleOpt{ui.Gap(12)},
		shadcn.Small("最近订单"),
		shadcn.DataTable(shadcn.DataTableProps{
			Columns: []shadcn.DataColumn{
				{Key: "id", Header: "订单号", Sortable: true, Width: 120},
				{Key: "customer", Header: "客户", Sortable: true},
				{Key: "status", Header: "状态", Width: 110},
				{Key: "amount", Header: "金额", Sortable: true, Width: 110},
			},
			Rows: []map[string]string{
				{"id": "#1024", "customer": "Alice Chen", "status": "已完成", "amount": "¥1,200"},
				{"id": "#1023", "customer": "Bob Li", "status": "处理中", "amount": "¥340"},
				{"id": "#1022", "customer": "Carol Wang", "status": "已完成", "amount": "¥875"},
				{"id": "#1021", "customer": "Dave Zhao", "status": "已退款", "amount": "¥60"},
			}}))

	main := ui.Div(ui.Style(ui.Column, ui.Grow(1), ui.Gap(16), ui.Padding(20), ui.Bg(th.Background)),
		topbar, stats, charts, table)

	return ui.ThemeProvider(th,
		ui.Div(ui.Style(ui.Row, ui.Fill, ui.Bg(th.Background), ui.TextColor(th.Foreground)),
			sidebar, main))
}

func main() {
	ui.WindowSize(1180, 880)
	ui.Run(ui.Use(App, struct{}{}))
}
