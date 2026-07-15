// Command accordion 是一个用 Tenon 复刻的 shadcn/ui 文档站：左侧是分组组件菜单，
// 点击即可切换；右侧是该组件的文档页（面包屑、标题栏、框架标签、实时预览、安装区），
// 整页放在 ScrollView 里、各区块随滚动逐个淡入。左上角 ‹ › 切换上/下一个组件。
//
//	go run ./example/accordion
package main

import (
	"fmt"
	"strings"

	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func main() {
	ui.WindowSize(1200, 860)
	ui.Run(ui.Use(App, struct{}{}))
}

// comp 是侧栏里一个可切换的组件文档条目。
type comp struct {
	name, pkg, desc string
	preview         func() *ui.Node
}

type group struct {
	label string
	items []comp
}

func registry() []group {
	u := func(fn func(struct{}) *ui.Node) func() *ui.Node {
		return func() *ui.Node { return ui.Use(fn, struct{}{}) }
	}
	return []group{
		{"开始", []comp{
			{"Accordion", "accordion", "一组垂直排列的可交互标题，每个标题都会展开一段内容。", u(pvAccordion)},
		}},
		{"表单", []comp{
			{"Button", "button", "触发一个动作或事件，内置多种变体与尺寸。", u(pvButton)},
			{"Input", "input", "供用户输入文本的单行表单控件。", u(pvInput)},
			{"Textarea", "textarea", "供用户输入多行文本的表单控件。", u(pvTextarea)},
			{"Checkbox", "checkbox", "一个可勾选的复选框，用于开关某个选项。", u(pvCheckbox)},
			{"Switch", "switch", "在开启与关闭两种状态间切换的开关。", u(pvSwitch)},
			{"Radio Group", "radio-group", "一组单选项，同一时刻只能选中其中一个。", u(pvRadio)},
			{"Slider", "slider", "通过拖动在给定范围内选择一个数值。", u(pvSlider)},
		}},
		{"展示", []comp{
			{"Badge", "badge", "用于标注状态或分类的小徽标。", u(pvBadge)},
			{"Avatar", "avatar", "展示用户头像，缺省时回退为首字母。", u(pvAvatar)},
			{"Card", "card", "一个带边框的内容容器，用于分组信息。", u(pvCard)},
			{"Separator", "separator", "用一条细线在视觉与语义上分隔内容。", u(pvSeparator)},
			{"Skeleton", "skeleton", "内容加载时展示的占位骨架。", u(pvSkeleton)},
			{"Progress", "progress", "展示任务完成进度的进度条。", u(pvProgress)},
		}},
		{"反馈", []comp{
			{"Alert", "alert", "向用户显示一条需要注意的提示信息。", u(pvAlert)},
			{"Tabs", "tabs", "一组分段视图，点击标签切换内容。", u(pvTabs)},
			{"Tooltip", "tooltip", "悬停元素时弹出的信息浮层。", u(pvTooltip)},
		}},
	}
}

func App(_ struct{}) *ui.Node {
	dark, setDark := ui.UseState(false)
	sel, setSel := ui.UseState(0)
	showCode, setShowCode := ui.UseState(false)
	fw, setFw := ui.UseState(0)
	inst, setInst := ui.UseState(0)
	pm, setPm := ui.UseState(0)

	th := ui.LightTheme
	if dark {
		th = ui.DarkTheme
	}
	strip := th.Muted

	groups := registry()
	var flat []comp
	for _, g := range groups {
		flat = append(flat, g.items...)
	}
	if sel >= len(flat) {
		sel = 0
	}
	cur := flat[sel]

	scrollRef, sc := ui.UseScroll()
	off, vp := sc.Offset, sc.Viewport
	rev := func(kids ...*ui.Node) *ui.Node {
		return ui.Use(revealBox, revealProps{off: off, vp: vp, kids: kids})
	}
	divider := func() *ui.Node { return ui.Div(ui.Style(ui.Height(1), ui.Bg(th.Border))) }

	// ---------- 侧栏 ----------
	idx := 0
	var sgroups []shadcn.SidebarGroup
	for _, g := range groups {
		var items []shadcn.SidebarItem
		for _, c := range g.items {
			i := idx
			items = append(items, shadcn.SidebarItem{Label: c.name, Active: i == sel,
				OnClick: func() { setSel(i); setShowCode(false) }})
			idx++
		}
		sgroups = append(sgroups, shadcn.SidebarGroup{Label: g.label, Items: items})
	}
	sidebar := shadcn.Sidebar(shadcn.SidebarProps{
		Header: ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(8)),
			ui.Text("◆", ui.FontSize(16), ui.TextColor(th.Primary)),
			ui.Text("Tenon UI", ui.FontSize(15), ui.Semibold)),
		Groups: sgroups,
		Footer: ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(8)),
			shadcn.Switch(shadcn.SwitchProps{Checked: dark, OnChange: setDark}),
			shadcn.Label("深色主题")),
	})

	// ---------- 头部 ----------
	prev := func() {
		if sel > 0 {
			setSel(sel - 1)
			setShowCode(false)
		}
	}
	next := func() {
		if sel < len(flat)-1 {
			setSel(sel + 1)
			setShowCode(false)
		}
	}
	iconBtn := func(path string, onClick func()) *ui.Node {
		return shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline, Size: shadcn.SizeIcon, OnClick: onClick},
			ui.Icon(path, 16, ui.TextColor(th.MutedForeground)))
	}
	header := rev(ui.Div(ui.Style(ui.Column, ui.Gap(12)),
		ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(6)),
			ui.Text("组件", ui.FontSize(13), ui.TextColor(th.MutedForeground)),
			ui.Text("/", ui.FontSize(13), ui.TextColor(th.MutedForeground)),
			ui.Text(cur.name, ui.FontSize(13), ui.TextColor(th.Foreground))),
		ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.JustifyBetween, ui.Gap(16)),
			ui.Text(cur.name, ui.FontSize(30), ui.Bold),
			ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(8)),
				shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline},
					ui.Text("复制当前页", ui.FontSize(13)),
					ui.Icon(ui.IconChevronDown, 14, ui.TextColor(th.MutedForeground))),
				iconBtn(ui.IconChevronLeft, prev),
				iconBtn(ui.IconChevronRight, next))),
		ui.Text(cur.desc, ui.FontSize(16), ui.TextColor(th.MutedForeground)),
	))

	// ---------- 框架标签 ----------
	fwTab := func(label string, i int) *ui.Node {
		active := fw == i
		col, under := th.MutedForeground, ui.Transparent
		if active {
			col, under = th.Foreground, th.Foreground
		}
		return ui.Div(ui.Style(ui.Column, ui.Gap(10)), ui.OnClick(func() { setFw(i) }),
			ui.Text(label, ui.FontSize(14), ui.Medium, ui.TextColor(col)),
			ui.Div(ui.Style(ui.Height(2), ui.Radius(1), ui.Bg(under))))
	}
	fwTabs := rev(ui.Div(ui.Style(ui.Column),
		ui.Div(ui.Style(ui.Row, ui.Gap(20)), fwTab("Radix UI", 0), fwTab("Base UI", 1)),
		divider()))

	// ---------- 预览卡 ----------
	sym := strings.ReplaceAll(cur.name, " ", "")
	previewInner := []*ui.Node{ui.Style(ui.Column, ui.Border(1, th.Border), ui.Radius(th.Radius+4), ui.Clip, ui.Bg(th.Card)),
		ui.Div(ui.Style(ui.Column, ui.ItemsCenter, ui.JustifyCenter, ui.MinHeight(320), ui.PaddingXY(40, 40)),
			cur.preview()),
		divider(),
		ui.Div(ui.Style(ui.Column, ui.ItemsCenter, ui.JustifyCenter, ui.Bg(strip), ui.PaddingXY(24, 26)),
			ui.Div(ui.Style(ui.Column, ui.Gap(3), ui.Opacity(0.55), ui.ItemsCenter),
				ui.Text(fmt.Sprintf("import { %s } from", sym), ui.FontSize(12.5), ui.TextColor(th.MutedForeground)),
				ui.Text(fmt.Sprintf("  \"@/components/ui/%s\"", cur.pkg), ui.FontSize(12.5), ui.TextColor(th.MutedForeground))),
			ui.Div(ui.Style(ui.Height(14))),
			shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline, OnClick: func() { setShowCode(!showCode) }},
				ui.Text(codeLabel(showCode), ui.FontSize(13)))),
	}
	if showCode {
		previewInner = append(previewInner, divider(),
			ui.Div(ui.Style(ui.Column, ui.Gap(3), ui.Bg(strip), ui.PaddingXY(20, 18)),
				codeLine(th, fmt.Sprintf("import { %s } from \"@/components/ui/%s\"", sym, cur.pkg)),
				ui.Div(ui.Style(ui.Height(8))),
				codeLine(th, "export function Demo() {"),
				codeLine(th, fmt.Sprintf("  return <%s />", sym)),
				codeLine(th, "}")))
	}
	preview := rev(ui.Div(previewInner...))

	// ---------- 安装区 ----------
	miniTab := func(label string, i int) *ui.Node {
		active := inst == i
		col, under := th.MutedForeground, ui.Transparent
		if active {
			col, under = th.Foreground, th.Foreground
		}
		return ui.Div(ui.Style(ui.Column, ui.Gap(9)), ui.OnClick(func() { setInst(i) }),
			ui.Text(label, ui.FontSize(13.5), ui.Medium, ui.TextColor(col)),
			ui.Div(ui.Style(ui.Height(2), ui.Radius(1), ui.Bg(under))))
	}
	pmNames := []string{"pnpm", "npm", "yarn", "bun"}
	pmPre := []string{"pnpm dlx", "npx", "yarn dlx", "bunx"}
	pmTab := func(i int) *ui.Node {
		active := pm == i
		st := []ui.StyleOpt{ui.PaddingXY(10, 5), ui.Radius(6)}
		col := th.MutedForeground
		if active {
			st = append(st, ui.Bg(th.Background))
			col = th.Foreground
		}
		return ui.Button(ui.Style(st...), ui.OnClick(func() { setPm(i) }),
			ui.Text(pmNames[i], ui.FontSize(12.5), ui.TextColor(col)))
	}
	addCmd := "--bun shadcn@latest add " + cur.pkg
	if pm != 3 {
		addCmd = "shadcn@latest add " + cur.pkg
	}
	terminal := ui.Div(ui.Style(ui.Column, ui.Border(1, th.Border), ui.Radius(th.Radius+4), ui.Clip, ui.Bg(strip)),
		ui.Div(ui.Style(ui.Row, ui.Gap(2), ui.Padding(6), ui.Bg(strip)),
			pmTab(0), pmTab(1), pmTab(2), pmTab(3)),
		divider(),
		ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(12), ui.PaddingXY(16, 14)),
			ui.Text("›", ui.FontSize(13), ui.TextColor(th.MutedForeground)),
			ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.Gap(6), ui.Grow(1)),
				ui.Text(pmPre[pm], ui.FontSize(13), ui.TextColor(th.MutedForeground)),
				ui.Text(addCmd, ui.FontSize(13), ui.TextColor(th.Foreground))),
			shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Ghost, Size: shadcn.SizeSm},
				ui.Text("复制", ui.FontSize(12.5), ui.TextColor(th.MutedForeground)))))
	importBox := ui.Div(ui.Style(ui.Column, ui.Border(1, th.Border), ui.Radius(th.Radius+4), ui.Clip, ui.Bg(strip)),
		ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.PaddingXY(16, 16)),
			codeLine(th, fmt.Sprintf("import { %s } from \"@/components/ui/%s\"", sym, cur.pkg))))
	install := rev(ui.Div(ui.Style(ui.Column, ui.Gap(16)),
		ui.Text("安装", ui.FontSize(22), ui.Semibold),
		ui.Div(ui.Style(ui.Column),
			ui.Div(ui.Style(ui.Row, ui.Gap(18)), miniTab("命令", 0), miniTab("手动", 1)),
			divider()),
		terminal,
		ui.Text("导入", ui.FontSize(16), ui.Semibold),
		importBox))

	// ---------- 组装 ----------
	content := ui.Div(ui.Style(ui.Column, ui.MaxWidth(760), ui.WidthPct(100), ui.Gap(28)),
		header, fwTabs, preview,
		ui.Div(ui.Style(ui.Height(8))),
		install,
		ui.Div(ui.Style(ui.Height(60))))
	page := ui.Div(ui.Style(ui.Column, ui.ItemsCenter, ui.WidthPct(100), ui.PaddingXY(32, 44), ui.Bg(th.Background)),
		content)

	return ui.ThemeProvider(th,
		ui.Div(ui.Style(ui.Row, ui.Fill, ui.Bg(th.Background), ui.TextColor(th.Foreground)),
			sidebar,
			ui.Div(ui.Style(ui.Width(1), ui.HeightPct(100), ui.Bg(th.Border))),
			ui.ScrollView(scrollRef, ui.Style(ui.Grow(1), ui.HeightPct(100)), page)))
}

func codeLabel(show bool) string {
	if show {
		return "Hide Code"
	}
	return "View Code"
}

func codeLine(th ui.Theme, s string) *ui.Node {
	return ui.Text(s, ui.FontSize(12.5), ui.TextColor(th.Foreground))
}

// ================= 各组件预览 =================

func pvAccordion(_ struct{}) *ui.Node {
	th := ui.UseTheme()
	m := func(s string) *ui.Node { return ui.Text(s, ui.FontSize(14), ui.TextColor(th.MutedForeground)) }
	items := []shadcn.AccordionItemData{
		{Title: "What are your shipping options?", Content: []*ui.Node{
			m("We offer standard (5–7 days), express (2–3 days), and overnight shipping. Free shipping on international orders.")}},
		{Title: "What is your return policy?", Content: []*ui.Node{
			m("Items can be returned within 30 days of delivery in their original condition.")}},
		{Title: "How can I contact customer support?", Content: []*ui.Node{
			m("Reach us 24/7 by email at support@example.com, or use the live chat on any page.")}},
	}
	return ui.Div(ui.Style(ui.Width(430)), shadcn.Accordion(items))
}

func pvButton(_ struct{}) *ui.Node {
	return ui.Div(ui.Style(ui.Row, ui.Gap(12), ui.ItemsCenter, ui.JustifyCenter),
		shadcn.Button(shadcn.ButtonProps{}, ui.Text("Default")),
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Secondary}, ui.Text("Secondary")),
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline}, ui.Text("Outline")),
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Ghost}, ui.Text("Ghost")),
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Destructive}, ui.Text("Destructive")))
}

func pvInput(_ struct{}) *ui.Node {
	v, set := ui.UseState("")
	return ui.Div(ui.Style(ui.Column, ui.Gap(8), ui.Width(300)),
		shadcn.Label("邮箱"),
		shadcn.Input(shadcn.InputProps{Value: v, OnChange: set, Placeholder: "you@example.com"}))
}

func pvTextarea(_ struct{}) *ui.Node {
	v, set := ui.UseState("")
	return ui.Div(ui.Style(ui.Column, ui.Gap(8)),
		shadcn.Label("留言"),
		shadcn.Textarea(shadcn.TextareaProps{Value: v, OnChange: set, Placeholder: "输入你的留言…", Rows: 4}))
}

func pvCheckbox(_ struct{}) *ui.Node {
	c, set := ui.UseState(true)
	return ui.Div(ui.Style(ui.Row, ui.Gap(8), ui.ItemsCenter),
		shadcn.Checkbox(shadcn.CheckboxProps{Checked: c, OnChange: set}),
		shadcn.Label("接受条款与条件"))
}

func pvSwitch(_ struct{}) *ui.Node {
	a, setA := ui.UseState(true)
	b, setB := ui.UseState(false)
	row := func(sw, lbl *ui.Node) *ui.Node {
		return ui.Div(ui.Style(ui.Row, ui.Gap(10), ui.ItemsCenter), sw, lbl)
	}
	return ui.Div(ui.Style(ui.Column, ui.Gap(14)),
		row(shadcn.Switch(shadcn.SwitchProps{Checked: a, OnChange: setA}), shadcn.Label("推送通知")),
		row(shadcn.Switch(shadcn.SwitchProps{Checked: b, OnChange: setB}), shadcn.Label("营销邮件")))
}

func pvRadio(_ struct{}) *ui.Node {
	v, set := ui.UseState("标准配送")
	return shadcn.RadioGroup(shadcn.RadioGroupProps{Value: v,
		Options: []string{"标准配送", "加急配送", "次日达"}, OnChange: set})
}

func pvSlider(_ struct{}) *ui.Node {
	th := ui.UseTheme()
	v, set := ui.UseState(float32(40))
	return ui.Div(ui.Style(ui.Column, ui.Gap(12), ui.Width(240)),
		ui.Div(ui.Style(ui.Row, ui.JustifyBetween, ui.ItemsCenter),
			shadcn.Label("音量"),
			ui.Text(fmt.Sprintf("%.0f%%", v), ui.FontSize(13), ui.TextColor(th.MutedForeground))),
		shadcn.Slider(shadcn.SliderProps{Value: v, Min: 0, Max: 100, OnChange: set}))
}

func pvBadge(_ struct{}) *ui.Node {
	return ui.Div(ui.Style(ui.Row, ui.Gap(8), ui.ItemsCenter),
		shadcn.Badge(shadcn.BadgeProps{}, ui.Text("Default")),
		shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeSecondary}, ui.Text("Secondary")),
		shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeOutline}, ui.Text("Outline")),
		shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeDestructive}, ui.Text("Destructive")))
}

func pvAvatar(_ struct{}) *ui.Node {
	return ui.Div(ui.Style(ui.Row, ui.Gap(12), ui.ItemsCenter),
		shadcn.Avatar("SJ", 44), shadcn.Avatar("KM", 44),
		shadcn.Avatar("TN", 44), shadcn.Avatar("+5", 44))
}

func pvCard(_ struct{}) *ui.Node {
	// 用 Card 的正规分区（CardHeader/Content/Footer 各带 px-6 水平内边距）。
	return ui.Div(ui.Style(ui.Width(340)),
		shadcn.Card(
			shadcn.CardHeader(
				shadcn.CardTitle("创建项目"),
				shadcn.CardDescription("部署你的新项目，仅需几步即可上线。")),
			shadcn.CardContent(
				shadcn.Label("项目名称"),
				ui.Div(ui.Style(ui.Height(8))),
				shadcn.Input(shadcn.InputProps{Placeholder: "my-app"})),
			shadcn.CardFooter(
				ui.Div(ui.Style(ui.Row, ui.Grow(1), ui.JustifyEnd, ui.Gap(8)),
					shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline, Size: shadcn.SizeSm}, ui.Text("取消")),
					shadcn.Button(shadcn.ButtonProps{Size: shadcn.SizeSm}, ui.Text("部署"))))))
}

func pvSeparator(_ struct{}) *ui.Node {
	th := ui.UseTheme()
	return ui.Div(ui.Style(ui.Column, ui.Gap(12), ui.Width(300)),
		ui.Div(ui.Style(ui.Column, ui.Gap(2)),
			ui.Text("Tenon UI", ui.FontSize(15), ui.Semibold),
			ui.Text("一套可复制粘贴的组件。", ui.FontSize(13), ui.TextColor(th.MutedForeground))),
		shadcn.Separator(shadcn.SeparatorProps{}),
		ui.Div(ui.Style(ui.Row, ui.Gap(16)),
			ui.Text("博客", ui.FontSize(14)), ui.Text("文档", ui.FontSize(14)), ui.Text("源码", ui.FontSize(14))))
}

func pvSkeleton(_ struct{}) *ui.Node {
	return ui.Div(ui.Style(ui.Row, ui.Gap(14), ui.ItemsCenter),
		shadcn.Skeleton(48, 48),
		ui.Div(ui.Style(ui.Column, ui.Gap(10)),
			shadcn.Skeleton(220, 14), shadcn.Skeleton(170, 14)))
}

func pvProgress(_ struct{}) *ui.Node {
	return ui.Div(ui.Style(ui.Column, ui.Gap(16), ui.ItemsCenter),
		shadcn.Progress(0.3), shadcn.Progress(0.6), shadcn.Progress(0.9))
}

func pvAlert(_ struct{}) *ui.Node {
	return ui.Div(ui.Style(ui.Column, ui.Gap(12), ui.Width(420)),
		shadcn.Alert(shadcn.AlertProps{},
			shadcn.AlertTitle("提示"),
			shadcn.AlertDescription("这是一条默认提示信息，用于向用户传达状态。")),
		shadcn.Alert(shadcn.AlertProps{Variant: shadcn.AlertDestructive},
			shadcn.AlertTitle("出错了"),
			shadcn.AlertDescription("你的会话已过期，请重新登录。")))
}

func pvTabs(_ struct{}) *ui.Node {
	th := ui.UseTheme()
	i, set := ui.UseState(0)
	body := []string{"在这里管理你的账户信息。", "修改密码以保护账户安全。", "设置你希望接收的通知类型。"}
	return ui.Div(ui.Style(ui.Column, ui.Gap(16), ui.Width(380)),
		shadcn.Tabs(shadcn.TabsProps{Tabs: []string{"账户", "密码", "通知"}, Active: i, OnChange: set}),
		ui.Div(ui.Style(ui.Border(1, th.Border), ui.Radius(th.Radius+2), ui.Padding(16), ui.Bg(th.Card)),
			ui.Text(body[i], ui.FontSize(14), ui.TextColor(th.MutedForeground))))
}

func pvTooltip(_ struct{}) *ui.Node {
	return shadcn.Tooltip("这是一个提示浮层",
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline}, ui.Text("悬停我")))
}

// ================= 滚动淡入容器 =================

type revealProps struct {
	off, vp float32
	kids    []*ui.Node
}

func revealBox(p revealProps) *ui.Node {
	cref, r := ui.UseMeasure()
	seen := ui.UseRef(false)
	if r.H > 0 && r.Y <= p.off+p.vp*0.9 {
		*seen = true
	}
	target := float32(0)
	if *seen {
		target = 1
	}
	op := ui.UseTween(target, 380, ui.EaseOut)
	ty := ui.UseTween((1-target)*16, 380, ui.EaseOut)
	return ui.Div(append([]*ui.Node{cref, ui.Style(ui.Column, ui.Opacity(op), ui.TranslateXY(0, ty))}, p.kids...)...)
}
