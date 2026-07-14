package shadcn

import (
	"sort"
	"strconv"
	"strings"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

type DataColumn struct {
	Key      string  // 行数据的键
	Header   string  // 列标题
	Sortable bool    // 可点击表头排序
	Width    float32 // 固定列宽（0 = 平分剩余宽度）
}

type DataTableProps struct {
	Columns  []DataColumn
	Rows     []map[string]string
	Search   bool // 顶部全局搜索框
	PageSize int  // 每页行数（0 = 不分页）
}

// DataTable 是带搜索、排序、分页的数据表格。数据以 []map[string]string 提供，列引用键名。
func DataTable(p DataTableProps) *ui.Node { return ui.Use(dataTable, p) }

func dataTable(p DataTableProps) *ui.Node {
	th := ui.UseTheme()
	query, setQuery := ui.UseState("")
	sortKey, setSortKey := ui.UseState("")
	sortDesc, setSortDesc := ui.UseState(false)
	page, setPage := ui.UseState(0)

	// 过滤（拷贝，不修改 prop）
	q := strings.ToLower(strings.TrimSpace(query))
	rows := make([]map[string]string, 0, len(p.Rows))
	for _, r := range p.Rows {
		if q == "" || rowMatches(r, q) {
			rows = append(rows, r)
		}
	}
	// 排序
	if sortKey != "" {
		sort.SliceStable(rows, func(i, j int) bool {
			a, b := rows[i][sortKey], rows[j][sortKey]
			if sortDesc {
				a, b = b, a
			}
			return lessValue(a, b)
		})
	}
	// 分页
	total := len(rows)
	ps := p.PageSize
	if ps <= 0 {
		ps = maxInt(total, 1)
	}
	pages := (total + ps - 1) / ps
	if pages == 0 {
		pages = 1
	}
	cur := clampInt(page, 0, pages-1)
	start := cur * ps
	end := minInt(start+ps, total)
	var pageRows []map[string]string
	if start < total {
		pageRows = rows[start:end]
	}

	toggleSort := func(key string) func() {
		return func() {
			if sortKey == key {
				setSortDesc(!sortDesc)
			} else {
				setSortKey(key)
				setSortDesc(false)
			}
		}
	}

	var out []*ui.Node
	out = append(out, ui.Style(ui.Column, ui.Gap(12)))

	if p.Search {
		out = append(out, InputGroup(InputGroupProps{
			Leading: ui.Icon(ui.IconSearch, 16), Placeholder: "搜索…",
			Value: query, OnChange: func(v string) { setQuery(v); setPage(0) }}))
	}

	// 表格容器
	table := []*ui.Node{ui.Style(ui.Column, ui.Border(1, th.Border), ui.Radius(radiusMd(th)), ui.Clip)}
	// 表头
	head := []*ui.Node{ui.Style(ui.Row, ui.ItemsCenter, ui.Height(42), ui.PaddingXY(4, 0), ui.Bg(th.Muted))}
	for _, c := range p.Columns {
		col := c
		head = append(head, ui.Use(dtHeadCell, dtHeadProps{
			col: col, active: sortKey == col.Key, desc: sortDesc, onSort: toggleSort(col.Key)}))
	}
	table = append(table, ui.Div(head...))
	// 表身
	if len(pageRows) == 0 {
		table = append(table, ui.Div(ui.Style(ui.Row, ui.JustifyCenter, ui.PaddingXY(0, 24)),
			ui.Text("无数据", ui.FontSize(14), ui.TextColor(th.MutedForeground))))
	}
	for i, r := range pageRows {
		table = append(table, ui.Use(dtRow, dtRowProps{cols: p.Columns, row: r, last: i == len(pageRows)-1}))
	}
	out = append(out, ui.Div(table...))

	// 分页
	if p.PageSize > 0 && pages > 1 {
		out = append(out, ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween),
			ui.Text(pageInfo(total, start, end), ui.FontSize(13), ui.TextColor(th.MutedForeground)),
			ui.Div(ui.Style(ui.Row, ui.Gap(8)),
				Button(ButtonProps{Variant: Outline, Size: SizeSm, Disabled: cur == 0,
					OnClick: func() { setPage(cur - 1) }}, ui.Text("上一页")),
				Button(ButtonProps{Variant: Outline, Size: SizeSm, Disabled: cur >= pages-1,
					OnClick: func() { setPage(cur + 1) }}, ui.Text("下一页")),
			),
		))
	}
	return ui.Div(out...)
}

type dtHeadProps struct {
	col    DataColumn
	active bool
	desc   bool
	onSort func()
}

func dtHeadCell(p dtHeadProps) *ui.Node {
	th := ui.UseTheme()
	kids := []*ui.Node{cellStyle(p.col, ui.PaddingXY(12, 0)),
		ui.Text(p.col.Header, ui.FontSize(13), ui.Medium, ui.TextColor(th.MutedForeground))}
	if p.col.Sortable {
		icon := ui.IconChevronDown
		if p.active && !p.desc {
			icon = ui.IconChevronUp
		}
		color := th.Border
		if p.active {
			color = th.Foreground
		}
		kids = append(kids, ui.Icon(icon, 14, ui.TextColor(color)))
	}
	attrs := []*ui.Node{}
	if p.col.Sortable {
		attrs = append(attrs, ui.OnClick(p.onSort))
	}
	return ui.Div(append(attrs, kids...)...)
}

type dtRowProps struct {
	cols []DataColumn
	row  map[string]string
	last bool
}

func dtRow(p dtRowProps) *ui.Node {
	th := ui.UseTheme()
	hovered, _, ia := ui.UseInteraction()
	st := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.Height(44), ui.PaddingXY(4, 0)}
	if hovered {
		st = append(st, ui.Bg(over(th.Muted, th.Background, 0.5)))
	}
	cells := []*ui.Node{ui.Style(st...), ia}
	for _, c := range p.cols {
		cells = append(cells, ui.Div(cellStyle(c, ui.PaddingXY(12, 0)),
			ui.Text(p.row[c.Key], ui.FontSize(14), ui.TextColor(th.Foreground))))
	}
	row := ui.Div(cells...)
	if !p.last {
		return ui.Div(ui.Style(ui.Column), row, ui.Div(ui.Style(ui.Height(1), ui.Bg(th.Border))))
	}
	return row
}

// cellStyle 给单元格一致的宽度策略：固定宽，或 flex:1 1 0 等宽（Width(0)+Grow 保证
// 表头与表身列宽一致——否则 auto-basis 会让列宽随内容变化而错位）。Clip 截断超长内容。
func cellStyle(c DataColumn, extra ...ui.StyleOpt) *ui.Node {
	base := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.Gap(4), ui.Clip}
	if c.Width > 0 {
		base = append(base, ui.Width(c.Width))
	} else {
		base = append(base, ui.Width(0), ui.Grow(1))
	}
	return ui.Style(append(base, extra...)...)
}

func rowMatches(r map[string]string, q string) bool {
	for _, v := range r {
		if strings.Contains(strings.ToLower(v), q) {
			return true
		}
	}
	return false
}

func lessValue(a, b string) bool {
	af, ae := strconv.ParseFloat(a, 64)
	bf, be := strconv.ParseFloat(b, 64)
	if ae == nil && be == nil {
		return af < bf
	}
	return a < b
}

func pageInfo(total, start, end int) string {
	if total == 0 {
		return "共 0 条"
	}
	return strconv.Itoa(start+1) + "–" + strconv.Itoa(end) + " / 共 " + strconv.Itoa(total) + " 条"
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
