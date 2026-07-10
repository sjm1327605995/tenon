package main

import (
	"fmt"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// ============================================================
// 主题 —— 用 Context 跨层传递
// ============================================================

type Theme struct {
	Name             string
	Bg, Card, Border ui.Color
	Text, Sub        ui.Color
	Accent, AccentFg ui.Color
	Done             ui.Color
}

var light = Theme{
	Name: "light",
	Bg:   ui.Hex("#f3f4f6"), Card: ui.White, Border: ui.Hex("#e5e7eb"),
	Text: ui.Hex("#111827"), Sub: ui.Hex("#6b7280"),
	Accent: ui.Hex("#2563eb"), AccentFg: ui.White, Done: ui.Hex("#9ca3af"),
}

var dark = Theme{
	Name: "dark",
	Bg:   ui.Hex("#0b1120"), Card: ui.Hex("#1e293b"), Border: ui.Hex("#334155"),
	Text: ui.Hex("#f1f5f9"), Sub: ui.Hex("#94a3b8"),
	Accent: ui.Hex("#38bdf8"), AccentFg: ui.Hex("#0b1120"), Done: ui.Hex("#64748b"),
}

var ThemeCtx = ui.CreateContext(light)

// ============================================================
// Todos —— 用 useReducer 管理
// ============================================================

type Todo struct {
	ID   int
	Text string
	Done bool
}

type action struct {
	kind string // "add" | "toggle" | "delete"
	id   int
	text string
}

func todosReducer(s []Todo, a action) []Todo {
	switch a.kind {
	case "add":
		id := 1
		for _, t := range s {
			if t.ID >= id {
				id = t.ID + 1
			}
		}
		return append(append([]Todo{}, s...), Todo{ID: id, Text: a.text, Done: false})
	case "toggle":
		out := append([]Todo{}, s...)
		for i := range out {
			if out[i].ID == a.id {
				out[i].Done = !out[i].Done
			}
		}
		return out
	case "delete":
		out := make([]Todo, 0, len(s))
		for _, t := range s {
			if t.ID != a.id {
				out = append(out, t)
			}
		}
		return out
	}
	return s
}

// ============================================================
// 组件
// ============================================================

func App(_ struct{}) *ui.Node {
	isDark, setDark := ui.UseState(false)
	th := light
	if isDark {
		th = dark
	}
	return ThemeCtx.Provider(th,
		ui.Use(Page, PageProps{ToggleTheme: func() { setDark(!isDark) }}),
	)
}

type PageProps struct{ ToggleTheme func() }

func Page(p PageProps) *ui.Node {
	th := ui.UseContext(ThemeCtx)
	text, setText := ui.UseState("")
	todos, dispatch := ui.UseReducer(todosReducer, []Todo{
		{ID: 1, Text: "学习 hooks 模型", Done: true},
		{ID: 2, Text: "构建声明式 UI", Done: false},
		{ID: 3, Text: "接入 yoga 布局引擎", Done: true},
		{ID: 4, Text: "用 ebiten/vector 绘制", Done: true},
		{ID: 5, Text: "实现 Context 主题系统", Done: false},
		{ID: 6, Text: "受控输入与键盘编辑", Done: false},
		{ID: 7, Text: "带 key 的列表复用", Done: false},
		{ID: 8, Text: "ScrollView 滚动裁剪", Done: false},
	})

	add := func() {
		if text != "" {
			dispatch(action{kind: "add", text: text})
			setText("")
		}
	}

	items := make([]*ui.Node, len(todos))
	for i, td := range todos {
		items[i] = ui.Keyed(fmt.Sprintf("%d", td.ID),
			ui.Memo(TodoItem, TodoItemProps{Todo: td, Dispatch: dispatch}))
	}

	return ui.Div(
		ui.Style(ui.Column, ui.Width(560), ui.Height(600), ui.Bg(th.Bg),
			ui.Padding(28), ui.Gap(20)),

		// 顶栏
		ui.Div(
			ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween),
			ui.Text("Tenon Todos", ui.FontSize(26), ui.TextColor(th.Text)),
			ui.Button(
				ui.Style(ui.PaddingXY(14, 8), ui.Bg(th.Card), ui.Border(1, th.Border), ui.Radius(8)),
				ui.OnClick(p.ToggleTheme),
				ui.Text("主题: "+th.Name, ui.FontSize(13), ui.TextColor(th.Sub)),
			),
		),

		// 输入行
		ui.Div(
			ui.Style(ui.Row, ui.Gap(10), ui.ItemsCenter),
			ui.Input(
				ui.Style(ui.Grow(1), ui.Height(40), ui.PaddingXY(12, 0), ui.Bg(th.Card),
					ui.Border(1, th.Border), ui.Radius(8),
					ui.FontSize(16), ui.TextColor(th.Text)),
				ui.Value(text),
				ui.Placeholder("添加一项待办…"),
				ui.OnChange(setText),
			),
			ui.Button(
				ui.Style(ui.Width(72), ui.Height(40), ui.Bg(th.Accent), ui.Radius(8),
					ui.ItemsCenter, ui.JustifyCenter),
				ui.OnClick(add),
				ui.Text("添加", ui.FontSize(15), ui.TextColor(th.AccentFg)),
			),
		),

		// 列表（可滚动）
		ui.ScrollView(append([]*ui.Node{
			ui.Style(ui.Column, ui.Gap(10), ui.Grow(1), ui.PaddingXY(0, 4)),
		}, items...)...),

		ui.Text(fmt.Sprintf("共 %d 项", len(todos)), ui.FontSize(13), ui.TextColor(th.Sub)),
	)
}

type TodoItemProps struct {
	Todo     Todo
	Dispatch func(action)
}

func TodoItem(p TodoItemProps) *ui.Node {
	th := ui.UseContext(ThemeCtx)
	td := p.Todo

	label := td.Text
	labelColor := th.Text
	if td.Done {
		label = "✓ " + label
		labelColor = th.Done
	}

	return ui.Div(
		ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(12), ui.PaddingXY(14, 10),
			ui.Bg(th.Card), ui.Border(1, th.Border), ui.Radius(10)),

		ui.Button(
			ui.Style(ui.Width(28), ui.Height(28), ui.Radius(6), ui.ItemsCenter, ui.JustifyCenter,
				ui.Bg(boolColor(td.Done, th.Accent, th.Card)), ui.Border(1, th.Border)),
			ui.OnClick(func() { p.Dispatch(action{kind: "toggle", id: td.ID}) }),
			ui.Text(boolStr(td.Done, "✓", ""), ui.FontSize(14), ui.TextColor(th.AccentFg)),
		),

		ui.Div(ui.Style(ui.Grow(1)),
			ui.Text(label, ui.FontSize(16), ui.TextColor(labelColor)),
		),

		ui.Button(
			ui.Style(ui.PaddingXY(10, 6), ui.Radius(6)),
			ui.OnClick(func() { p.Dispatch(action{kind: "delete", id: td.ID}) }),
			ui.Text("删除", ui.FontSize(13), ui.TextColor(ui.Red)),
		),
	)
}

func boolColor(b bool, t, f ui.Color) ui.Color {
	if b {
		return t
	}
	return f
}
func boolStr(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

func main() {
	ui.Run(ui.Use(App, struct{}{}))
}
