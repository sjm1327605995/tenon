// Command router 演示 pkg/router 的栈式导航：一个收件箱列表，点邮件 Push 进详情屏，
// 详情屏可「返回」(Pop)、看「下一封」(Replace，深度不变)。左上角标题栏根据栈深度
// 显示返回按钮。
//
//	go run ./example/router
package main

import (
	"fmt"

	"github.com/sjm1327605995/tenon/pkg/router"
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func main() {
	ui.WindowSize(880, 720)
	ui.Run(ui.Use(App, struct{}{}))
}

type email struct{ id, from, initials, subject, body string }

var inbox = []email{
	{"1", "Alice Chen", "AC", "欢迎加入 Tenon", "很高兴见到你！这是一个用纯 Go 写的 React 风格 GUI 工具包。随便点点，感受一下导航栈。"},
	{"2", "构建机器人", "CI", "构建通过 ✓", "main 分支 #1421 构建成功，全部 128 个测试通过，用时 42s。"},
	{"3", "Bob Li", "BL", "关于路由设计", "我觉得桌面端不该照搬 URL 路由，栈式导航（push/pop）更贴切，你怎么看？"},
	{"4", "Newsletter", "NL", "本周前端要闻", "本期：签名式响应、局部渲染、以及为什么 flexbox 依然好用。"},
	{"5", "Carol Wang", "CW", "午饭？", "楼下新开了一家，12:30 一起？"},
}

func App(_ struct{}) *ui.Node {
	th := ui.LightTheme
	return ui.ThemeProvider(th,
		ui.Div(ui.Style(ui.Fill, ui.Bg(th.Background), ui.TextColor(th.Foreground)),
			router.Router(router.Props{
				Initial: "list",
				Screens: map[string]router.Screen{
					"list":   listScreen,
					"detail": detailScreen,
				},
			})))
}

// ---------- 列表屏 ----------

func listScreen(_ router.Params) *ui.Node {
	nav := router.UseNavigate()
	rows := make([]*ui.Node, len(inbox))
	for i, e := range inbox {
		rows[i] = ui.Use(emailRow, rowProps{e: e, onOpen: func() {
			nav.Push("detail", router.Params{"id": e.id})
		}})
	}
	return screen(
		topBar(ui.Text("收件箱", ui.FontSize(20), ui.Semibold), ui.Spacer(),
			shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeSecondary}, ui.Text(fmt.Sprintf("%d 封", len(inbox))))),
		ui.VStack(0, rows...))
}

type rowProps struct {
	e      email
	onOpen func()
}

func emailRow(p rowProps) *ui.Node {
	hovered, _, ia := ui.UseInteraction()
	bg := ui.Transparent
	if hovered {
		bg = ui.UseTheme().Muted
	}
	return ui.VStack(0,
		ui.HStack(12, ui.Style(ui.PaddingXY(20, 14), ui.Bg(bg)), ui.OnClick(p.onOpen), ia,
			shadcn.Avatar(p.e.initials, 40),
			ui.VStack(2, ui.Text(p.e.from, ui.FontSize(14), ui.Medium), muted(p.e.subject, 13)),
			ui.Spacer(),
			ui.Icon(ui.IconChevronRight, 16, dim())),
		shadcn.Separator(shadcn.SeparatorProps{}))
}

// ---------- 详情屏 ----------

func detailScreen(p router.Params) *ui.Node {
	nav := router.UseNavigate()
	e, idx := findEmail(p["id"])
	next := inbox[(idx+1)%len(inbox)]

	var back *ui.Node // 仅当能返回时才显示返回按钮
	if nav.CanPop() {
		back = shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Ghost, Size: shadcn.SizeSm, OnClick: nav.Pop},
			ui.Icon(ui.IconChevronLeft, 16, dim()), ui.Text("收件箱", ui.FontSize(13)))
	}

	body := ui.VStack(18, ui.Style(ui.PaddingXY(24, 20)),
		ui.Text(e.subject, ui.FontSize(22), ui.Semibold),
		ui.HStack(12, shadcn.Avatar(e.initials, 40),
			ui.VStack(2, ui.Text(e.from, ui.FontSize(14), ui.Medium), muted("发送至：我", 12))),
		shadcn.Separator(shadcn.SeparatorProps{}),
		ui.Text(e.body, ui.FontSize(14)),
		ui.HStack(8,
			shadcn.Button(shadcn.ButtonProps{}, ui.Text("回复")),
			// Replace：原地换成下一封，栈深度不变（返回仍回到列表）。
			shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline,
				OnClick: func() { nav.Replace("detail", router.Params{"id": next.id}) }},
				ui.Text("下一封"))))

	return screen(
		topBar(back, ui.Spacer(), muted(fmt.Sprintf("第 %d 层", nav.Depth()), 12)),
		body)
}

func findEmail(id string) (email, int) {
	for i, e := range inbox {
		if e.id == id {
			return e, i
		}
	}
	return inbox[0], 0
}

// ---------- 屏脚手架与小助手 ----------

// screen 是标准屏布局：固定顶栏 + 分隔线 + 可滚动内容。
func screen(bar, content *ui.Node) *ui.Node {
	return ui.Div(ui.Style(ui.Column, ui.Fill),
		bar,
		shadcn.Separator(shadcn.SeparatorProps{}),
		ui.ScrollView(ui.Style(ui.Grow(1)), content))
}

// topBar 是高 56 的顶栏，横向排入给定内容。
func topBar(items ...*ui.Node) *ui.Node {
	return ui.HStack(10, append([]*ui.Node{ui.Style(ui.Height(56), ui.PaddingXY(16, 0))}, items...)...)
}

// muted 是次要色文本；dim 是次要色（用于图标 currentColor）。二者自取当前主题。
func muted(s string, size float32) *ui.Node { return ui.Text(s, ui.FontSize(size), dim()) }
func dim() ui.StyleOpt                      { return ui.TextColor(ui.UseTheme().MutedForeground) }
