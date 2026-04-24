package main

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/antdesign"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// ==================== 页面标题辅助 ====================

func pageTitle(text string) tenon.Component {
	title := antdesign.NewAntTitle(text).SetLevel(antdesign.AntTitleH1)
	return components.NewView().SetMargin(yoga.EdgeBottom, 24).Add(title)
}

func sectionTitle(text string) tenon.Component {
	title := antdesign.NewAntTitle(text).SetLevel(antdesign.AntTitleH3)
	return components.NewView().SetMargin(yoga.EdgeTop, 16).SetMargin(yoga.EdgeBottom, 12).Add(title)
}

func demoBlock(title string, demos ...tenon.Component) tenon.Component {
	card := antdesign.NewAntCard().SetTitle(title).SetShadow(1)
	for _, d := range demos {
		card.Add(d)
	}
	return components.NewView().SetMargin(yoga.EdgeBottom, 16).Add(card)
}

// ==================== General 示例 ====================

type ButtonDemo struct {
	tenon.BaseWidget
	loadings [5]bool
}

func NewButtonDemo() *ButtonDemo {
	d := &ButtonDemo{}
	d.Init(d)
	return d
}

func (d *ButtonDemo) enterLoading(idx int) {
	if d.loadings[idx] {
		return
	}
	d.loadings[idx] = true
	d.SetState(func() {})
	go func() {
		time.Sleep(3 * time.Second)
		d.loadings[idx] = false
		d.SetState(func() {})
	}()
}

func (d *ButtonDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Button 按钮"),
		demoBlock("按钮类型",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Primary").SetType(antdesign.AntButtonPrimary),
				antdesign.NewAntButton("Default").SetType(antdesign.AntButtonDefault),
				antdesign.NewAntButton("Dashed").SetType(antdesign.AntButtonDashed),
				antdesign.NewAntButton("Text").SetType(antdesign.AntButtonText),
				antdesign.NewAntButton("Link").SetType(antdesign.AntButtonLink),
			),
		),
		demoBlock("危险按钮",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Primary").SetType(antdesign.AntButtonPrimary).SetDanger(true),
				antdesign.NewAntButton("Default").SetType(antdesign.AntButtonDefault).SetDanger(true),
				antdesign.NewAntButton("Text").SetType(antdesign.AntButtonText).SetDanger(true),
				antdesign.NewAntButton("Link").SetType(antdesign.AntButtonLink).SetDanger(true),
			),
		),
		demoBlock("按钮尺寸",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Large").SetType(antdesign.AntButtonPrimary).SetSize(antdesign.AntButtonLarge),
				antdesign.NewAntButton("Middle").SetType(antdesign.AntButtonPrimary).SetSize(antdesign.AntButtonMiddle),
				antdesign.NewAntButton("Small").SetType(antdesign.AntButtonPrimary).SetSize(antdesign.AntButtonSmall),
			),
		),
		demoBlock("图标按钮",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Search").SetType(antdesign.AntButtonPrimary).SetAntIcon(antdesign.AntIconSearch),
				antdesign.NewAntButton("Search").SetType(antdesign.AntButtonPrimary).SetAntIcon(antdesign.AntIconSearch).SetIconPosition(antdesign.AntButtonIconRight),
				antdesign.NewAntButton("").SetType(antdesign.AntButtonPrimary).SetShape(antdesign.AntButtonShapeCircle).SetAntIcon(antdesign.AntIconSearch),
			),
		),
		demoBlock("幽灵按钮",
			components.NewView().SetBackgroundColor(tenon.GetTheme().PrimaryColor).SetPadding(yoga.EdgeAll, 16).SetBorderRadius(8).Add(
				antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
					antdesign.NewAntButton("Primary").SetType(antdesign.AntButtonPrimary).SetGhost(true),
					antdesign.NewAntButton("Default").SetType(antdesign.AntButtonDefault).SetGhost(true),
					antdesign.NewAntButton("Dashed").SetType(antdesign.AntButtonDashed).SetGhost(true),
					antdesign.NewAntButton("Danger").SetType(antdesign.AntButtonPrimary).SetGhost(true).SetDanger(true),
				),
			),
		),
		demoBlock("加载中",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Loading").SetType(antdesign.AntButtonPrimary).SetLoading(true),
				antdesign.NewAntButton("Loading").SetType(antdesign.AntButtonPrimary).SetSize(antdesign.AntButtonSmall).SetLoading(true),
				antdesign.NewAntButton("").SetType(antdesign.AntButtonPrimary).SetShape(antdesign.AntButtonShapeCircle).SetAntIcon(antdesign.AntIconSearch).SetLoading(true),
			),
		),
		demoBlock("点击加载",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Click me").SetType(antdesign.AntButtonPrimary).SetLoading(d.loadings[0]).SetOnClick(func() { d.enterLoading(0) }),
				antdesign.NewAntButton("Icon End").SetType(antdesign.AntButtonPrimary).SetLoading(d.loadings[1]).SetAntIcon(antdesign.AntIconSearch).SetIconPosition(antdesign.AntButtonIconRight).SetOnClick(func() { d.enterLoading(1) }),
				antdesign.NewAntButton("Icon Replace").SetType(antdesign.AntButtonPrimary).SetAntIcon(antdesign.AntIconSearch).SetLoading(d.loadings[2]).SetOnClick(func() { d.enterLoading(2) }),
				antdesign.NewAntButton("").SetType(antdesign.AntButtonPrimary).SetShape(antdesign.AntButtonShapeCircle).SetAntIcon(antdesign.AntIconSearch).SetLoading(d.loadings[3]).SetOnClick(func() { d.enterLoading(3) }),
			),
		),
		demoBlock("Block 按钮",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceVertical).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Primary Block").SetType(antdesign.AntButtonPrimary).SetBlock(true),
				antdesign.NewAntButton("Default Block").SetType(antdesign.AntButtonDefault).SetBlock(true),
				antdesign.NewAntButton("Danger Block").SetType(antdesign.AntButtonPrimary).SetDanger(true).SetBlock(true),
			),
		),
		demoBlock("形状",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Default").SetType(antdesign.AntButtonPrimary),
				antdesign.NewAntButton("Round").SetType(antdesign.AntButtonPrimary).SetShape(antdesign.AntButtonShapeRound),
				antdesign.NewAntButton("").SetType(antdesign.AntButtonPrimary).SetShape(antdesign.AntButtonShapeCircle).SetAntIcon(antdesign.AntIconSearch),
			),
		),
		demoBlock("状态",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Disabled").SetType(antdesign.AntButtonPrimary).SetDisabled(true),
				antdesign.NewAntButton("Disabled").SetType(antdesign.AntButtonDefault).SetDisabled(true),
			),
		),
	)
}

type FloatButtonDemo struct{ tenon.BaseWidget }

func NewFloatButtonDemo() *FloatButtonDemo {
	d := &FloatButtonDemo{}
	d.Init(d)
	return d
}

func (d *FloatButtonDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("FloatButton 浮动按钮"),
		demoBlock("基本用法",
			components.NewView().SetHeight(300).SetBackgroundColor(tenon.GetTheme().BackgroundColor).Add(
				antdesign.NewAntFloatButton("+").SetOnClick(func() {
					fmt.Println("FloatButton clicked!")
				}),
			),
		),
	)
}

type DividerDemo struct{ tenon.BaseWidget }

func NewDividerDemo() *DividerDemo {
	d := &DividerDemo{}
	d.Init(d)
	return d
}

func (d *DividerDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Divider 分割线"),
		demoBlock("水平分割线",
			components.NewText("Before divider").SetFontSize(tenon.GetTheme().FontSizeBase),
			antdesign.NewAntDivider(),
			components.NewText("After divider").SetFontSize(tenon.GetTheme().FontSizeBase),
		),
		demoBlock("带文字的分割线",
			antdesign.NewAntDivider().SetText("Center"),
			antdesign.NewAntDivider().SetText("Left").SetAlign("left"),
			antdesign.NewAntDivider().SetText("Right").SetAlign("right"),
		),
		demoBlock("竖向分割线",
			components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetAlignItems(yoga.AlignCenter).Add(
				components.NewText("Text").SetFontSize(tenon.GetTheme().FontSizeBase),
				components.NewView().SetHeight(20).Add(antdesign.NewAntDivider().SetVertical(true)),
				components.NewText("Link").SetFontSize(tenon.GetTheme().FontSizeBase),
				components.NewView().SetHeight(20).Add(antdesign.NewAntDivider().SetVertical(true)),
				components.NewText("Button").SetFontSize(tenon.GetTheme().FontSizeBase),
			),
		),
	)
}

type TypographyDemo struct{ tenon.BaseWidget }

func NewTypographyDemo() *TypographyDemo {
	d := &TypographyDemo{}
	d.Init(d)
	return d
}

func (d *TypographyDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Typography 排版"),
		demoBlock("标题",
			antdesign.NewAntTitle("h1. Ant Design").SetLevel(antdesign.AntTitleH1),
			antdesign.NewAntTitle("h2. Ant Design").SetLevel(antdesign.AntTitleH2),
			antdesign.NewAntTitle("h3. Ant Design").SetLevel(antdesign.AntTitleH3),
			antdesign.NewAntTitle("h4. Ant Design").SetLevel(antdesign.AntTitleH4),
		),
		demoBlock("文本样式",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntText("Primary").SetType(antdesign.AntTextPrimary),
				antdesign.NewAntText("Secondary").SetType(antdesign.AntTextSecondary),
				antdesign.NewAntText("Success").SetType(antdesign.AntTextSuccess),
				antdesign.NewAntText("Warning").SetType(antdesign.AntTextWarning),
				antdesign.NewAntText("Danger").SetType(antdesign.AntTextDanger),
			),
		),
		demoBlock("段落",
			antdesign.NewAntParagraph("Tenon is a cross-platform UI framework built with Go and Ebiten. It provides a React-like component model with Flexbox layout."),
		),
		demoBlock("链接",
			antdesign.NewAntLink("Ant Design Official").SetHref("https://ant-design.antgroup.com"),
		),
	)
}

type SpaceDemo struct{ tenon.BaseWidget }

func NewSpaceDemo() *SpaceDemo {
	d := &SpaceDemo{}
	d.Init(d)
	return d
}

func (d *SpaceDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Space 间距"),
		demoBlock("水平间距",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Button 1").SetType(antdesign.AntButtonPrimary),
				antdesign.NewAntButton("Button 2").SetType(antdesign.AntButtonPrimary),
				antdesign.NewAntButton("Button 3").SetType(antdesign.AntButtonPrimary),
			),
		),
		demoBlock("垂直间距",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceVertical).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntButton("Button 1").SetType(antdesign.AntButtonPrimary),
				antdesign.NewAntButton("Button 2").SetType(antdesign.AntButtonPrimary),
				antdesign.NewAntButton("Button 3").SetType(antdesign.AntButtonPrimary),
			),
		),
	)
}

// ==================== Layout 示例 ====================

type LayoutDemo struct{ tenon.BaseWidget }

func NewLayoutDemo() *LayoutDemo {
	d := &LayoutDemo{}
	d.Init(d)
	return d
}

func (d *LayoutDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Layout 布局"),
		demoBlock("经典布局",
			components.NewView().SetHeight(240).Add(
				antdesign.NewAntHeader().SetHeight(48).Add(
					components.NewText("Header").SetFontSize(tenon.GetTheme().FontSizeBase).SetColor(tenon.GetTheme().SurfaceColor),
				),
				components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetFlexGrow(1).Add(
					antdesign.NewAntSider().SetWidth(80).Add(
						components.NewText("Sider").SetFontSize(tenon.GetTheme().FontSizeBase).SetColor(tenon.GetTheme().SurfaceColor),
					),
					antdesign.NewAntContent().Add(
						components.NewText("Content").SetFontSize(tenon.GetTheme().FontSizeBase),
					),
				),
				antdesign.NewAntFooter().SetHeight(40).Add(
					components.NewText("Footer").SetFontSize(tenon.GetTheme().FontSizeBase).SetColor(tenon.GetTheme().SurfaceColor),
				),
			),
		),
		demoBlock("Flex 布局",
			antdesign.NewAntFlex().SetJustify(yoga.JustifyCenter).SetAlign(yoga.AlignCenter).SetGap(8).Add(
				components.NewView().SetWidth(80).SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				components.NewView().SetWidth(80).SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				components.NewView().SetWidth(80).SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
			),
			components.NewView().SetHeight(8),
			antdesign.NewAntFlex().SetJustify(yoga.JustifySpaceBetween).SetGap(8).Add(
				components.NewView().SetWidth(80).SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				components.NewView().SetWidth(80).SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				components.NewView().SetWidth(80).SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
			),
		),
		demoBlock("Grid 网格",
			antdesign.NewAntRow().SetGutter(8).Add(
				antdesign.NewAntCol().SetSpan(6).Add(
					components.NewView().SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				),
				antdesign.NewAntCol().SetSpan(6).Add(
					components.NewView().SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				),
				antdesign.NewAntCol().SetSpan(6).Add(
					components.NewView().SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				),
				antdesign.NewAntCol().SetSpan(6).Add(
					components.NewView().SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				),
			),
			components.NewView().SetHeight(8),
			antdesign.NewAntRow().SetGutter(8).Add(
				antdesign.NewAntCol().SetSpan(8).Add(
					components.NewView().SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				),
				antdesign.NewAntCol().SetSpan(8).Add(
					components.NewView().SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				),
				antdesign.NewAntCol().SetSpan(8).Add(
					components.NewView().SetHeight(40).SetBackgroundColor(tenon.GetTheme().PrimaryColor),
				),
			),
		),
	)
}

// ==================== Navigation 示例 ====================

type BreadcrumbDemo struct{ tenon.BaseWidget }

func NewBreadcrumbDemo() *BreadcrumbDemo {
	d := &BreadcrumbDemo{}
	d.Init(d)
	return d
}

func (d *BreadcrumbDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Breadcrumb 面包屑"),
		demoBlock("基本用法",
			antdesign.NewAntBreadcrumb().SetItems([]antdesign.AntBreadcrumbItem{
				{Label: "Home"},
				{Label: "Application Center"},
				{Label: "Application List"},
				{Label: "An Application"},
			}),
		),
		demoBlock("带链接",
			antdesign.NewAntBreadcrumb().SetItems([]antdesign.AntBreadcrumbItem{
				{Label: "Home", Href: "/"},
				{Label: "Application Center", Href: "/app"},
				{Label: "An Application"},
			}),
		),
		demoBlock("自定义分隔符",
			antdesign.NewAntBreadcrumb().SetSeparator(">").SetItems([]antdesign.AntBreadcrumbItem{
				{Label: "Home"},
				{Label: "Application Center"},
				{Label: "An Application"},
			}),
		),
	)
}

type PaginationDemo struct{ tenon.BaseWidget }

func NewPaginationDemo() *PaginationDemo {
	d := &PaginationDemo{}
	d.Init(d)
	return d
}

func (d *PaginationDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Pagination 分页"),
		demoBlock("基本分页",
			antdesign.NewAntPagination().SetCurrent(1).SetTotal(85).SetPageSize(10),
		),
		demoBlock("更多分页",
			antdesign.NewAntPagination().SetCurrent(6).SetTotal(500).SetPageSize(10),
		),
	)
}

type StepsDemo struct{ tenon.BaseWidget }

func NewStepsDemo() *StepsDemo {
	d := &StepsDemo{}
	d.Init(d)
	return d
}

func (d *StepsDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Steps 步骤条"),
		demoBlock("基本用法",
			antdesign.NewAntSteps().SetCurrent(1).SetItems([]antdesign.AntStepItem{
				{Title: "Finished", Description: "This is a description."},
				{Title: "In Progress", Description: "This is a description."},
				{Title: "Waiting", Description: "This is a description."},
			}),
		),
		demoBlock("小型步骤条",
			antdesign.NewAntSteps().SetCurrent(1).SetSize("small").SetItems([]antdesign.AntStepItem{
				{Title: "Finished"},
				{Title: "In Progress"},
				{Title: "Waiting"},
			}),
		),
		demoBlock("竖直方向",
			antdesign.NewAntSteps().SetCurrent(1).SetDirection("vertical").SetItems([]antdesign.AntStepItem{
				{Title: "Finished", Description: "This is a description."},
				{Title: "In Progress", Description: "This is a description."},
				{Title: "Waiting", Description: "This is a description."},
			}),
		),
	)
}

type TabsDemo struct{ tenon.BaseWidget }

func NewTabsDemo() *TabsDemo {
	d := &TabsDemo{}
	d.Init(d)
	return d
}

func (d *TabsDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Tabs 标签页"),
		demoBlock("线型标签页",
			antdesign.NewAntTabs().SetActiveKey("1").SetItems([]antdesign.AntTabItem{
				{Key: "1", Label: "Tab 1", Children: components.NewText("Content of Tab Pane 1").SetFontSize(tenon.GetTheme().FontSizeBase)},
				{Key: "2", Label: "Tab 2", Children: components.NewText("Content of Tab Pane 2").SetFontSize(tenon.GetTheme().FontSizeBase)},
				{Key: "3", Label: "Tab 3", Children: components.NewText("Content of Tab Pane 3").SetFontSize(tenon.GetTheme().FontSizeBase)},
			}),
		),
		demoBlock("卡片型标签页",
			antdesign.NewAntTabs().SetActiveKey("1").SetType(antdesign.AntTabCard).SetItems([]antdesign.AntTabItem{
				{Key: "1", Label: "Tab 1", Children: components.NewText("Content of Card Tab 1").SetFontSize(tenon.GetTheme().FontSizeBase)},
				{Key: "2", Label: "Tab 2", Children: components.NewText("Content of Card Tab 2").SetFontSize(tenon.GetTheme().FontSizeBase)},
			}),
		),
	)
}

// ==================== Data Entry 示例 ====================

type InputDemo struct{ tenon.BaseWidget }

func NewInputDemo() *InputDemo {
	d := &InputDemo{}
	d.Init(d)
	return d
}

func (d *InputDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Input 输入框"),
		demoBlock("基本用法",
			antdesign.NewAntInput().SetPlaceholder("Basic usage").SetWidth(300),
		),
		demoBlock("前缀和后缀",
			antdesign.NewAntInput().SetPrefix("￥").SetSuffix("RMB").SetWidth(300),
		),
		demoBlock("搜索框",
			antdesign.NewAntInput().SetPlaceholder("input search text").SetSearch(true).SetWidth(300),
		),
	)
}

// ==================== Data Display 示例 ====================

type AvatarDemo struct{ tenon.BaseWidget }

func NewAvatarDemo() *AvatarDemo {
	d := &AvatarDemo{}
	d.Init(d)
	return d
}

func (d *AvatarDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Avatar 头像"),
		demoBlock("基本用法",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntAvatar().SetSize(40).SetChildren("U"),
				antdesign.NewAntAvatar().SetSize(40).SetIcon("👤"),
				antdesign.NewAntAvatar().SetSize(40).SetShape(antdesign.AntAvatarSquare).SetChildren("U"),
			),
		),
		demoBlock("尺寸",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntAvatar().SetSize(24).SetChildren("S"),
				antdesign.NewAntAvatar().SetSize(32).SetChildren("M"),
				antdesign.NewAntAvatar().SetSize(40).SetChildren("L"),
				antdesign.NewAntAvatar().SetSize(64).SetChildren("XL"),
			),
		),
		demoBlock("头像组",
			antdesign.NewAntAvatarGroup().SetMaxCount(3).Add(
				antdesign.NewAntAvatar().SetChildren("A"),
				antdesign.NewAntAvatar().SetChildren("B"),
				antdesign.NewAntAvatar().SetChildren("C"),
				antdesign.NewAntAvatar().SetChildren("D"),
				antdesign.NewAntAvatar().SetChildren("E"),
			),
		),
	)
}

type BadgeDemo struct{ tenon.BaseWidget }

func NewBadgeDemo() *BadgeDemo {
	d := &BadgeDemo{}
	d.Init(d)
	return d
}

func (d *BadgeDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Badge 徽标数"),
		demoBlock("基本用法",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntBadge(
					antdesign.NewAntButton("Messages").SetType(antdesign.AntButtonDefault),
				).SetCount(5),
				antdesign.NewAntBadge(
					antdesign.NewAntButton("Notifications").SetType(antdesign.AntButtonDefault),
				).SetCount(0).SetDot(true),
				antdesign.NewAntBadge(
					antdesign.NewAntButton("Overflow").SetType(antdesign.AntButtonDefault),
				).SetCount(120),
			),
		),
	)
}

type CardDemo struct{ tenon.BaseWidget }

func NewCardDemo() *CardDemo {
	d := &CardDemo{}
	d.Init(d)
	return d
}

func (d *CardDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Card 卡片"),
		demoBlock("基本卡片",
			antdesign.NewAntCard().SetTitle("Default size card").Add(
				components.NewText("Card content").SetFontSize(tenon.GetTheme().FontSizeBase),
				components.NewText("Card content").SetFontSize(tenon.GetTheme().FontSizeBase),
				components.NewText("Card content").SetFontSize(tenon.GetTheme().FontSizeBase),
			),
		),
		demoBlock("无标题卡片",
			antdesign.NewAntCard().Add(
				components.NewText("This card has no title.").SetFontSize(tenon.GetTheme().FontSizeBase),
			),
		),
		demoBlock("不同阴影",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceMiddle).Add(
				antdesign.NewAntCard().SetTitle("Shadow: 1").SetShadow(1).Add(components.NewText("Small shadow").SetFontSize(tenon.GetTheme().FontSizeBase)),
				antdesign.NewAntCard().SetTitle("Shadow: 2").SetShadow(2).Add(components.NewText("Medium shadow").SetFontSize(tenon.GetTheme().FontSizeBase)),
				antdesign.NewAntCard().SetTitle("Shadow: 3").SetShadow(3).Add(components.NewText("Large shadow").SetFontSize(tenon.GetTheme().FontSizeBase)),
			),
		),
	)
}

type EmptyDemo struct{ tenon.BaseWidget }

func NewEmptyDemo() *EmptyDemo {
	d := &EmptyDemo{}
	d.Init(d)
	return d
}

func (d *EmptyDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Empty 空状态"),
		demoBlock("基本用法",
			antdesign.NewAntEmpty().SetDescription("No Data"),
		),
		demoBlock("自定义描述",
			antdesign.NewAntEmpty().SetDescription("No search results"),
		),
	)
}

type TableDemo struct{ tenon.BaseWidget }

func NewTableDemo() *TableDemo {
	d := &TableDemo{}
	d.Init(d)
	return d
}

func (d *TableDemo) Render() tenon.Component {
	columns := []antdesign.AntColumn{
		{Title: "Name", Key: "name"},
		{Title: "Age", Key: "age", Width: 80},
		{Title: "Address", Key: "address"},
	}
	data := []map[string]any{
		{"name": "John Brown", "age": "32", "address": "New York No. 1 Lake Park"},
		{"name": "Jim Green", "age": "42", "address": "London No. 1 Lake Park"},
		{"name": "Joe Black", "age": "32", "address": "Sidney No. 1 Lake Park"},
		{"name": "Disabled User", "age": "99", "address": "Sidney No. 1 Lake Park"},
	}
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Table 表格"),
		demoBlock("基本用法",
			antdesign.NewAntTable(columns).SetData(data),
		),
		demoBlock("斑马纹表格",
			antdesign.NewAntTable(columns).SetData(data).SetStripe(true),
		),
	)
}

type TagDemo struct{ tenon.BaseWidget }

func NewTagDemo() *TagDemo {
	d := &TagDemo{}
	d.Init(d)
	return d
}

func (d *TagDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Tag 标签"),
		demoBlock("基本用法",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceSmall).Add(
				antdesign.NewAntTag("Tag 1"),
				antdesign.NewAntTag("Tag 2").SetColor(antdesign.AntTagRed),
				antdesign.NewAntTag("Tag 3").SetColor(antdesign.AntTagOrange),
				antdesign.NewAntTag("Tag 4").SetColor(antdesign.AntTagGreen),
				antdesign.NewAntTag("Tag 5").SetColor(antdesign.AntTagBlue),
			),
		),
		demoBlock("可关闭标签",
			antdesign.NewAntSpace().SetDirection(antdesign.AntSpaceHorizontal).SetSize(antdesign.AntSpaceSmall).Add(
				antdesign.NewAntTag("Closable").SetClosable(true),
				antdesign.NewAntTag("Red").SetColor(antdesign.AntTagRed).SetClosable(true),
				antdesign.NewAntTag("Blue").SetColor(antdesign.AntTagBlue).SetClosable(true),
			),
		),
	)
}

// ==================== Feedback 示例 ====================

type AlertDemo struct{ tenon.BaseWidget }

func NewAlertDemo() *AlertDemo {
	d := &AlertDemo{}
	d.Init(d)
	return d
}

func (d *AlertDemo) Render() tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle("Alert 警告提示"),
		demoBlock("四种样式",
			antdesign.NewAntAlert("Success Text").SetType(antdesign.AntAlertSuccess),
			antdesign.NewAntAlert("Info Text").SetType(antdesign.AntAlertInfo),
			antdesign.NewAntAlert("Warning Text").SetType(antdesign.AntAlertWarning),
			antdesign.NewAntAlert("Error Text").SetType(antdesign.AntAlertError),
		),
		demoBlock("含有辅助性文字",
			antdesign.NewAntAlert("Success Text").SetType(antdesign.AntAlertSuccess).SetDesc("Detailed description and advice about successful copywriting."),
			antdesign.NewAntAlert("Info Text").SetType(antdesign.AntAlertInfo).SetDesc("Additional description and information about copywriting."),
		),
		demoBlock("可关闭",
			antdesign.NewAntAlert("Warning Text").SetType(antdesign.AntAlertWarning).SetClosable(true),
			antdesign.NewAntAlert("Error Text").SetType(antdesign.AntAlertError).SetClosable(true),
		),
	)
}

// ==================== 根应用：左侧菜单 + 右侧内容 ====================

type App struct {
	tenon.BaseWidget
	currentPage string
}

func NewApp() *App {
	a := &App{currentPage: "button"}
	a.Init(a)
	return a
}

func placeholderPage(title string) tenon.Component {
	return components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle(title),
		components.NewText("该组件示例待实现").SetFontSize(tenon.GetTheme().FontSizeBase).SetColor(tenon.GetTheme().TextMutedColor),
	)
}

func (a *App) Render() tenon.Component {
	menu := components.NewMenu().SetItems([]components.MenuItemData{
		{Key: "", Label: "General"},
		{Key: "button", Label: "Button"},
		{Key: "floatbutton", Label: "FloatButton"},
		{Key: "icon", Label: "Icon"},
		{Key: "typography", Label: "Typography"},
		{Key: "", Label: "Layout"},
		{Key: "divider", Label: "Divider"},
		{Key: "flex", Label: "Flex"},
		{Key: "grid", Label: "Grid"},
		{Key: "layout", Label: "Layout"},
		{Key: "masonry", Label: "Masonry"},
		{Key: "space", Label: "Space"},
		{Key: "splitter", Label: "Splitter"},
		{Key: "", Label: "Navigation"},
		{Key: "anchor", Label: "Anchor"},
		{Key: "breadcrumb", Label: "Breadcrumb"},
		{Key: "dropdown", Label: "Dropdown"},
		{Key: "menu", Label: "Menu"},
		{Key: "pagination", Label: "Pagination"},
		{Key: "steps", Label: "Steps"},
		{Key: "tabs", Label: "Tabs"},
		{Key: "", Label: "Data Entry"},
		{Key: "autocomplete", Label: "AutoComplete"},
		{Key: "cascader", Label: "Cascader"},
		{Key: "checkbox", Label: "Checkbox"},
		{Key: "colorpicker", Label: "ColorPicker"},
		{Key: "datepicker", Label: "DatePicker"},
		{Key: "form", Label: "Form"},
		{Key: "input", Label: "Input"},
		{Key: "inputnumber", Label: "InputNumber"},
		{Key: "mentions", Label: "Mentions"},
		{Key: "radio", Label: "Radio"},
		{Key: "rate", Label: "Rate"},
		{Key: "select", Label: "Select"},
		{Key: "slider", Label: "Slider"},
		{Key: "switch", Label: "Switch"},
		{Key: "timepicker", Label: "TimePicker"},
		{Key: "transfer", Label: "Transfer"},
		{Key: "treeselect", Label: "TreeSelect"},
		{Key: "upload", Label: "Upload"},
		{Key: "", Label: "Other"},
		{Key: "affix", Label: "Affix"},
		{Key: "app", Label: "App"},
	}).SetSelectedKey(a.currentPage).SetOnSelect(func(key string) {
		a.SetState(func() { a.currentPage = key })
	})

	content := components.NewView().
		SetFlexGrow(1).
		SetPadding(yoga.EdgeAll, 24).
		SetBackgroundColor(tenon.GetTheme().BackgroundColor)

	switch a.currentPage {
	// General
	case "button":
		content.AddChild(NewButtonDemo())
	case "floatbutton":
		content.AddChild(NewFloatButtonDemo())
	case "icon":
		content.AddChild(placeholderPage("Icon 图标"))
	case "typography":
		content.AddChild(NewTypographyDemo())
	// Layout
	case "divider":
		content.AddChild(NewDividerDemo())
	case "flex":
		content.AddChild(placeholderPage("Flex 弹性布局"))
	case "grid":
		content.AddChild(placeholderPage("Grid 栅格"))
	case "layout":
		content.AddChild(NewLayoutDemo())
	case "masonry":
		content.AddChild(placeholderPage("Masonry 瀑布流"))
	case "space":
		content.AddChild(NewSpaceDemo())
	case "splitter":
		content.AddChild(placeholderPage("Splitter 分割面板"))
	// Navigation
	case "anchor":
		content.AddChild(placeholderPage("Anchor 锚点"))
	case "breadcrumb":
		content.AddChild(NewBreadcrumbDemo())
	case "dropdown":
		content.AddChild(placeholderPage("Dropdown 下拉菜单"))
	case "menu":
		content.AddChild(placeholderPage("Menu 导航菜单"))
	case "pagination":
		content.AddChild(NewPaginationDemo())
	case "steps":
		content.AddChild(NewStepsDemo())
	case "tabs":
		content.AddChild(NewTabsDemo())
	// Data Entry
	case "autocomplete":
		content.AddChild(placeholderPage("AutoComplete 自动完成"))
	case "cascader":
		content.AddChild(placeholderPage("Cascader 级联选择"))
	case "checkbox":
		content.AddChild(placeholderPage("Checkbox 多选框"))
	case "colorpicker":
		content.AddChild(placeholderPage("ColorPicker 颜色选择器"))
	case "datepicker":
		content.AddChild(placeholderPage("DatePicker 日期选择框"))
	case "form":
		content.AddChild(placeholderPage("Form 表单"))
	case "input":
		content.AddChild(NewInputDemo())
	case "inputnumber":
		content.AddChild(placeholderPage("InputNumber 数字输入框"))
	case "mentions":
		content.AddChild(placeholderPage("Mentions 提及"))
	case "radio":
		content.AddChild(placeholderPage("Radio 单选框"))
	case "rate":
		content.AddChild(placeholderPage("Rate 评分"))
	case "select":
		content.AddChild(placeholderPage("Select 选择器"))
	case "slider":
		content.AddChild(placeholderPage("Slider 滑动输入条"))
	case "switch":
		content.AddChild(placeholderPage("Switch 开关"))
	case "timepicker":
		content.AddChild(placeholderPage("TimePicker 时间选择框"))
	case "transfer":
		content.AddChild(placeholderPage("Transfer 穿梭框"))
	case "treeselect":
		content.AddChild(placeholderPage("TreeSelect 树选择"))
	case "upload":
		content.AddChild(placeholderPage("Upload 上传"))
	// Other
	case "affix":
		content.AddChild(placeholderPage("Affix 固钉"))
	case "app":
		content.AddChild(placeholderPage("App 包裹组件"))
	}

	scrollView := components.NewScrollView().
		SetHeightPercent(100)
	scrollView.GetElement().Yoga.StyleSetFlexGrow(1)
	scrollView.Content().AddChild(content)

	return components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetWidthPercent(100).
		SetHeightPercent(100).
		Add(menu, scrollView)
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("Failed to init default font: " + err.Error())
	}
	if err := fonts.LoadFontFromFile(fonts.FontFamilySans, "font/OPPOSans-Medium.ttf"); err == nil {
		fonts.SetDefaultFontFamily(fonts.FontFamilySans)
	}

	antdesign.SetAntTheme()

	fmt.Println("Ant Design Component Showcase")
	fmt.Println("Open: http://ant-design.antgroup.com/components/overview-cn")

	tenon.Run(NewApp(), 1100, 750)
}
