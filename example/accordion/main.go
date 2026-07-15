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

// comp 是侧栏里一个可切换的组件文档条目；group 是一组条目。
type comp struct {
	name, pkg, desc string
	preview         func() *ui.Node
}
type group struct {
	label string
	items []comp
}

// registry 是全部组件的目录。新增一个组件只需在此加一行 + 写一个 pvXxx 预览函数。
func registry() []group {
	c := func(name, pkg, desc string, pv func(struct{}) *ui.Node) comp {
		return comp{name, pkg, desc, func() *ui.Node { return ui.Use(pv, struct{}{}) }}
	}
	return []group{
		{"开始", []comp{
			c("Accordion", "accordion", "一组垂直排列的可交互标题，每个标题都会展开一段内容。", pvAccordion),
		}},
		{"表单", []comp{
			c("Button", "button", "触发一个动作或事件，内置多种变体与尺寸。", pvButton),
			c("Input", "input", "供用户输入文本的单行表单控件。", pvInput),
			c("Textarea", "textarea", "供用户输入多行文本的表单控件。", pvTextarea),
			c("Checkbox", "checkbox", "一个可勾选的复选框，用于开关某个选项。", pvCheckbox),
			c("Switch", "switch", "在开启与关闭两种状态间切换的开关。", pvSwitch),
			c("Radio Group", "radio-group", "一组单选项，同一时刻只能选中其中一个。", pvRadio),
			c("Slider", "slider", "通过拖动在给定范围内选择一个数值。", pvSlider),
		}},
		{"展示", []comp{
			c("Badge", "badge", "用于标注状态或分类的小徽标。", pvBadge),
			c("Avatar", "avatar", "展示用户头像，缺省时回退为首字母。", pvAvatar),
			c("Card", "card", "一个带边框的内容容器，用于分组信息。", pvCard),
			c("Separator", "separator", "用一条细线在视觉与语义上分隔内容。", pvSeparator),
			c("Skeleton", "skeleton", "内容加载时展示的占位骨架。", pvSkeleton),
			c("Progress", "progress", "展示任务完成进度的进度条。", pvProgress),
		}},
		{"反馈", []comp{
			c("Alert", "alert", "向用户显示一条需要注意的提示信息。", pvAlert),
			c("Tabs", "tabs", "一组分段视图，点击标签切换内容。", pvTabs),
			c("Tooltip", "tooltip", "悬停元素时弹出的信息浮层。", pvTooltip),
		}},
	}
}

// app 汇聚一次渲染所需的上下文（主题、当前状态与各 setter）。文档页的各区块与可复用
// 交互都以它的成员方法呈现，从而不必把 th / 状态在函数间层层透传。
type app struct {
	th     ui.Theme
	groups []group
	flat   []comp
	cur    comp
	dark   bool
	sel    int
	show   bool // View Code 展开
	fw     int  // Radix UI / Base UI
	scroll ui.ScrollInfo

	setDark func(bool)
	setSel  func(int)
	setShow func(bool)
	setFw   func(int)
}

func App(_ struct{}) *ui.Node {
	dark, setDark := ui.UseState(false)
	sel, setSel := ui.UseState(0)
	show, setShow := ui.UseState(false)
	fw, setFw := ui.UseState(0)
	scrollRef, scroll := ui.UseScroll()

	groups := registry()
	flat := flatten(groups)
	if sel >= len(flat) {
		sel = 0
	}
	a := &app{
		th: pickTheme(dark), groups: groups, flat: flat, cur: flat[sel],
		dark: dark, sel: sel, show: show, fw: fw, scroll: scroll,
		setDark: setDark, setSel: setSel, setShow: setShow, setFw: setFw,
	}

	content := ui.VStack(28,
		a.reveal(a.header()),
		a.reveal(a.frameTabs()),
		a.reveal(a.preview()),
		vspace(8),
		a.reveal(a.install()),
		vspace(60),
	)
	page := ui.Div(ui.Style(ui.Column, ui.ItemsCenter, ui.WidthPct(100), ui.PaddingXY(32, 44), ui.Bg(a.th.Background)),
		ui.Div(ui.Style(ui.Column, ui.MaxWidth(760), ui.WidthPct(100)), content))

	return ui.ThemeProvider(a.th,
		ui.Div(ui.Style(ui.Row, ui.Fill, ui.Bg(a.th.Background), ui.TextColor(a.th.Foreground)),
			a.sidebar(),
			vline(a.th),
			ui.ScrollView(scrollRef, ui.Style(ui.Grow(1), ui.HeightPct(100)), page)))
}

// goTo 切换到第 i 个组件（并收起 View Code）。
func (a *app) goTo(i int) {
	if i >= 0 && i < len(a.flat) {
		a.setSel(i)
		a.setShow(false)
	}
}

// reveal 把区块包进滚动淡入容器。
func (a *app) reveal(kids ...*ui.Node) *ui.Node {
	return ui.Use(revealBox, revealProps{off: a.scroll.Offset, vp: a.scroll.Viewport, kids: kids})
}

func (a *app) sidebar() *ui.Node {
	idx := 0
	var sgroups []shadcn.SidebarGroup
	for _, g := range a.groups {
		var items []shadcn.SidebarItem
		for _, c := range g.items {
			i := idx // 快照跨组递增的下标供事件闭包捕获
			items = append(items, shadcn.SidebarItem{Label: c.name, Active: i == a.sel, OnClick: func() { a.goTo(i) }})
			idx++
		}
		sgroups = append(sgroups, shadcn.SidebarGroup{Label: g.label, Items: items})
	}
	return shadcn.Sidebar(shadcn.SidebarProps{
		Header: ui.HStack(8, ui.Text("◆", ui.FontSize(16), ui.TextColor(a.th.Primary)),
			ui.Text("Tenon UI", ui.FontSize(15), ui.Semibold)),
		Groups: sgroups,
		Footer: ui.HStack(8, shadcn.Switch(shadcn.SwitchProps{Checked: a.dark, OnChange: a.setDark}),
			shadcn.Label("深色主题")),
	})
}

func (a *app) header() *ui.Node {
	actions := ui.HStack(8,
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline},
			ui.Text("复制当前页", ui.FontSize(13)),
			ui.Icon(ui.IconChevronDown, 14, ui.TextColor(a.th.MutedForeground))),
		a.arrow(ui.IconChevronLeft, a.sel-1),
		a.arrow(ui.IconChevronRight, a.sel+1))
	return ui.VStack(12,
		ui.HStack(6, muted(a.th, "组件", 13), muted(a.th, "/", 13), ui.Text(a.cur.name, ui.FontSize(13))),
		ui.HStack(16, ui.Text(a.cur.name, ui.FontSize(30), ui.Bold), ui.Spacer(), actions),
		muted(a.th, a.cur.desc, 16))
}

// arrow 是标题栏的上/下一个组件按钮。
func (a *app) arrow(path string, to int) *ui.Node {
	return shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline, Size: shadcn.SizeIcon, OnClick: func() { a.goTo(to) }},
		ui.Icon(path, 16, ui.TextColor(a.th.MutedForeground)))
}

func (a *app) frameTabs() *ui.Node {
	return a.underlineTabs([]string{"Radix UI", "Base UI"}, a.fw, a.setFw)
}

func (a *app) preview() *ui.Node {
	th, c := a.th, a.cur
	sym := symbol(c.name)
	kids := []*ui.Node{ui.Style(ui.Column, ui.Border(1, th.Border), ui.Radius(th.Radius+4), ui.Clip, ui.Bg(th.Card)),
		ui.Center(ui.Style(ui.MinHeight(320), ui.PaddingXY(40, 40)), c.preview()),
		hline(th),
		ui.Div(ui.Style(ui.Column, ui.ItemsCenter, ui.JustifyCenter, ui.Bg(th.Muted), ui.PaddingXY(24, 26)),
			ui.Div(ui.Style(ui.Column, ui.Gap(3), ui.ItemsCenter, ui.Opacity(0.55)),
				muted(th, fmt.Sprintf("shadcn.%s(…)", sym), 12.5)),
			vspace(14),
			shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline, OnClick: func() { a.setShow(!a.show) }},
				ui.Text(codeLabel(a.show), ui.FontSize(13)))),
	}
	if a.show {
		kids = append(kids, hline(th),
			ui.Div(ui.Style(ui.Column, ui.Gap(3), ui.Bg(th.Muted), ui.PaddingXY(20, 18)),
				code(th, "import \"github.com/sjm1327605995/tenon/pkg/shadcn\""),
				vspace(8),
				code(th, "func Demo() *ui.Node {"),
				code(th, fmt.Sprintf("  return shadcn.%s(…)", sym)),
				code(th, "}")))
	}
	return ui.Div(kids...)
}

// install：Tenon 是一个 Go 库，没有单独的组件安装步骤，导入 pkg/shadcn 即可使用。
// 因此这里只展示 Go 的导入与用法，而不是 shadcn 那样的包管理器安装命令。
func (a *app) install() *ui.Node {
	th, c := a.th, a.cur
	box := func(s string) *ui.Node {
		return ui.Div(ui.Style(ui.Column, ui.Border(1, th.Border), ui.Radius(th.Radius+4), ui.Clip, ui.Bg(th.Muted)),
			ui.Div(ui.Style(ui.Row, ui.ItemsCenter, ui.PaddingXY(16, 16)), code(th, s)))
	}
	return ui.VStack(16,
		ui.Text("安装", ui.FontSize(22), ui.Semibold),
		muted(th, "Tenon 是一个 Go 库，无需单独安装组件——直接导入 pkg/shadcn 即可使用。", 15),
		ui.Text("导入", ui.FontSize(16), ui.Semibold),
		box("import \"github.com/sjm1327605995/tenon/pkg/shadcn\""),
		ui.Text("用法", ui.FontSize(16), ui.Semibold),
		box(fmt.Sprintf("shadcn.%s(…)", symbol(c.name))))
}

// underlineTabs 是下划线式标签栏（选中项文字加深、底部 2px 下划线），下方带一条整宽分隔线。
// 框架标签（Radix/Base）与安装标签（命令/手动）共用它。
func (a *app) underlineTabs(labels []string, active int, onSel func(int)) *ui.Node {
	tabs := make([]*ui.Node, len(labels))
	for i, label := range labels {
		col, under := a.th.MutedForeground, ui.Transparent
		if i == active {
			col, under = a.th.Foreground, a.th.Foreground
		}
		tabs[i] = ui.Div(ui.Style(ui.Column, ui.Gap(10)), ui.OnClick(func() { onSel(i) }),
			ui.Text(label, ui.FontSize(14), ui.Medium, ui.TextColor(col)),
			ui.Div(ui.Style(ui.Height(2), ui.Radius(1), ui.Bg(under))))
	}
	return ui.VStack(0,
		ui.Div(append([]*ui.Node{ui.Style(ui.Row, ui.Gap(20))}, tabs...)...),
		hline(a.th))
}

// ---------- 纯助手（无状态，预览函数也在用，故保持自由函数） ----------

func pickTheme(dark bool) ui.Theme {
	if dark {
		return ui.DarkTheme
	}
	return ui.LightTheme
}
func flatten(gs []group) []comp {
	var f []comp
	for _, g := range gs {
		f = append(f, g.items...)
	}
	return f
}
func symbol(name string) string  { return strings.ReplaceAll(name, " ", "") }
func hline(th ui.Theme) *ui.Node { return ui.Div(ui.Style(ui.Height(1), ui.Bg(th.Border))) }
func vline(th ui.Theme) *ui.Node {
	return ui.Div(ui.Style(ui.Width(1), ui.HeightPct(100), ui.Bg(th.Border)))
}
func vspace(h float32) *ui.Node { return ui.Div(ui.Style(ui.Height(h))) }
func code(th ui.Theme, s string) *ui.Node {
	return ui.Text(s, ui.FontSize(12.5), ui.TextColor(th.Foreground))
}
func muted(th ui.Theme, s string, size float32) *ui.Node {
	return ui.Text(s, ui.FontSize(size), ui.TextColor(th.MutedForeground))
}
func codeLabel(show bool) string {
	if show {
		return "Hide Code"
	}
	return "View Code"
}

// ---------- 各组件预览（独立子组件） ----------

func pvAccordion(_ struct{}) *ui.Node {
	th := ui.UseTheme()
	m := func(s string) *ui.Node { return muted(th, s, 14) }
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
	return ui.HStack(12,
		shadcn.Button(shadcn.ButtonProps{}, ui.Text("Default")),
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Secondary}, ui.Text("Secondary")),
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline}, ui.Text("Outline")),
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Ghost}, ui.Text("Ghost")),
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Destructive}, ui.Text("Destructive")))
}

func pvInput(_ struct{}) *ui.Node {
	v, set := ui.UseState("")
	return ui.VStack(8, ui.Style(ui.Width(300)),
		shadcn.Label("邮箱"),
		shadcn.Input(shadcn.InputProps{Value: v, OnChange: set, Placeholder: "you@example.com"}))
}

func pvTextarea(_ struct{}) *ui.Node {
	v, set := ui.UseState("")
	return ui.VStack(8,
		shadcn.Label("留言"),
		shadcn.Textarea(shadcn.TextareaProps{Value: v, OnChange: set, Placeholder: "输入你的留言…", Rows: 4}))
}

func pvCheckbox(_ struct{}) *ui.Node {
	c, set := ui.UseState(true)
	return ui.HStack(8, shadcn.Checkbox(shadcn.CheckboxProps{Checked: c, OnChange: set}),
		shadcn.Label("接受条款与条件"))
}

func pvSwitch(_ struct{}) *ui.Node {
	a, setA := ui.UseState(true)
	b, setB := ui.UseState(false)
	return ui.VStack(14,
		ui.HStack(10, shadcn.Switch(shadcn.SwitchProps{Checked: a, OnChange: setA}), shadcn.Label("推送通知")),
		ui.HStack(10, shadcn.Switch(shadcn.SwitchProps{Checked: b, OnChange: setB}), shadcn.Label("营销邮件")))
}

func pvRadio(_ struct{}) *ui.Node {
	v, set := ui.UseState("标准配送")
	return shadcn.RadioGroup(shadcn.RadioGroupProps{Value: v,
		Options: []string{"标准配送", "加急配送", "次日达"}, OnChange: set})
}

func pvSlider(_ struct{}) *ui.Node {
	th := ui.UseTheme()
	v, set := ui.UseState(float32(40))
	return ui.VStack(12,
		ui.HStack(0, ui.Style(ui.Width(240), ui.JustifyBetween),
			shadcn.Label("音量"), muted(th, fmt.Sprintf("%.0f%%", v), 13)),
		shadcn.Slider(shadcn.SliderProps{Value: v, Min: 0, Max: 100, OnChange: set}))
}

func pvBadge(_ struct{}) *ui.Node {
	return ui.HStack(8,
		shadcn.Badge(shadcn.BadgeProps{}, ui.Text("Default")),
		shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeSecondary}, ui.Text("Secondary")),
		shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeOutline}, ui.Text("Outline")),
		shadcn.Badge(shadcn.BadgeProps{Variant: shadcn.BadgeDestructive}, ui.Text("Destructive")))
}

func pvAvatar(_ struct{}) *ui.Node {
	return ui.HStack(12, shadcn.Avatar("SJ", 44), shadcn.Avatar("KM", 44),
		shadcn.Avatar("TN", 44), shadcn.Avatar("+5", 44))
}

func pvCard(_ struct{}) *ui.Node {
	return ui.Div(ui.Style(ui.Width(340)),
		shadcn.Card(
			shadcn.CardHeader(
				shadcn.CardTitle("创建项目"),
				shadcn.CardDescription("部署你的新项目，仅需几步即可上线。")),
			shadcn.CardContent(
				shadcn.Label("项目名称"),
				vspace(8),
				shadcn.Input(shadcn.InputProps{Placeholder: "my-app"})),
			shadcn.CardFooter(
				ui.Spacer(),
				shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline, Size: shadcn.SizeSm}, ui.Text("取消")),
				shadcn.Button(shadcn.ButtonProps{Size: shadcn.SizeSm}, ui.Text("部署")))))
}

func pvSeparator(_ struct{}) *ui.Node {
	th := ui.UseTheme()
	return ui.VStack(12,
		ui.VStack(2, ui.Text("Tenon UI", ui.FontSize(15), ui.Semibold), muted(th, "一套可复制粘贴的组件。", 13)),
		shadcn.Separator(shadcn.SeparatorProps{}),
		ui.HStack(16, ui.Text("博客", ui.FontSize(14)), ui.Text("文档", ui.FontSize(14)), ui.Text("源码", ui.FontSize(14))))
}

func pvSkeleton(_ struct{}) *ui.Node {
	return ui.HStack(14, shadcn.Skeleton(48, 48),
		ui.VStack(10, shadcn.Skeleton(220, 14), shadcn.Skeleton(170, 14)))
}

func pvProgress(_ struct{}) *ui.Node {
	return ui.VStack(16, shadcn.Progress(0.3), shadcn.Progress(0.6), shadcn.Progress(0.9))
}

func pvAlert(_ struct{}) *ui.Node {
	return ui.VStack(12, ui.Style(ui.Width(420)),
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
	return ui.VStack(16, ui.Style(ui.Width(380)),
		shadcn.Tabs(shadcn.TabsProps{Tabs: []string{"账户", "密码", "通知"}, Active: i, OnChange: set}),
		ui.Div(ui.Style(ui.Border(1, th.Border), ui.Radius(th.Radius+2), ui.Padding(16), ui.Bg(th.Card)),
			muted(th, body[i], 14)))
}

func pvTooltip(_ struct{}) *ui.Node {
	return shadcn.Tooltip("这是一个提示浮层",
		shadcn.Button(shadcn.ButtonProps{Variant: shadcn.Outline}, ui.Text("悬停我")))
}

// ---------- 滚动淡入容器 ----------

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
