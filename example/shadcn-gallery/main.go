package main

import (
	"image/color"
	"time"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/pkg/v2/components"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// demoCard wraps content in a styled card for the gallery.
func demoCard(title string, content core.Element) *components.View {
	card := components.NewCard()
	card.SetTitle(title)
	if content != nil {
		card.AddContent(content)
	}
	return &card.View
}

// demoSection creates a labeled section.
func demoSection(label string, children ...core.Element) core.Element {
	v := components.NewView()
	v.SetFlexDirection(yoga.FlexDirectionColumn)
	v.SetGap(yoga.GutterAll, 8)
	v.SetMargin(yoga.EdgeBottom, 24)
	title := components.NewText(label).SetFontSize(20).SetColor(core.GetTheme().TextColor)
	v.Add(title)
	for _, c := range children {
		v.Add(c)
	}
	return v
}

// ========== Pages ==========

type LayoutPage struct{ core.BaseWidget }

func NewLayoutPage() *LayoutPage {
	p := &LayoutPage{}
	p.Init(p)
	return p
}

func (p *LayoutPage) Render() core.Element {
	// Card demo
	cardDemo := components.NewCard().
		SetTitle("Payment Method").
		SetDescription("Add a new payment method to your account.").
		AddContent(
			components.NewTextInput().SetPlaceholder("Card number"),
			components.NewTextInput().SetPlaceholder("Name on card"),
		).
		SetFooter(
			components.NewButton("Cancel").SetVariant(components.ButtonOutline),
			components.NewButton("Save Card"),
		)

	// Accordion demo
	acc := components.NewAccordion([]components.AccordionItem{
		{Title: "Is it accessible?", Content: func() core.Element {
			return components.NewText("Yes. It adheres to the WAI-ARIA design pattern.").SetColor(core.GetTheme().TextColor)
		}},
		{Title: "Is it styled?", Content: func() core.Element {
			return components.NewText("Yes. It comes with default styles that matches the other components.").SetColor(core.GetTheme().TextColor)
		}},
		{Title: "Is it animated?", Content: func() core.Element {
			return components.NewText("Yes. It's animated by default, but you can disable it.").SetColor(core.GetTheme().TextColor)
		}},
	})

	// Collapsible demo
	coll := components.NewCollapsible(
		"What are the terms and conditions?",
		components.NewText("By using this service, you agree to our terms and conditions.").SetColor(core.GetTheme().TextColor),
	)

	// AspectRatio demo
	ar := components.NewAspectRatio(16.0 / 9.0)
	arBox := components.NewView()
	arBox.SetBackgroundColor(core.GetTheme().MutedColor)
	arBox.SetFlexDirection(yoga.FlexDirectionColumn)
	arBox.SetJustifyContent(yoga.JustifyCenter)
	arBox.SetAlignItems(yoga.AlignCenter)
	arBox.Add(components.NewText("16:9 Aspect Ratio").SetColor(core.GetTheme().MutedForegroundColor))
	ar.SetChild(arBox)

	return components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 16).
		SetPadding(yoga.EdgeAll, 24).
		Add(
			demoSection("Card", demoCard("Card Component", &cardDemo.View)),
			demoSection("Accordion", demoCard("Accordion Component", acc)),
			demoSection("Collapsible", demoCard("Collapsible Component", coll)),
			demoSection("Aspect Ratio", demoCard("Aspect Ratio Component", ar)),
		)
}

type FormsPage struct{ core.BaseWidget }

func NewFormsPage() *FormsPage {
	p := &FormsPage{}
	p.Init(p)
	return p
}

func (p *FormsPage) Render() core.Element {
	// Button variants
	btnRow := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetAlignItems(yoga.AlignCenter)
	btnRow.Add(
		components.NewButton("Default"),
		components.NewButton("Secondary").SetVariant(components.ButtonSecondary),
		components.NewButton("Outline").SetVariant(components.ButtonOutline),
		components.NewButton("Ghost").SetVariant(components.ButtonGhost),
		components.NewButton("Destructive").SetVariant(components.ButtonDestructive),
	)

	// Input & Textarea
	inputGroup := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 8)
	inputGroup.Add(
		components.NewLabel("Email"),
		components.NewTextInput().SetPlaceholder("Enter your email"),
		components.NewLabel("Message"),
		components.NewTextarea().SetRows(4).SetPlaceholder("Type your message here..."),
	)

	// Checkbox & Radio & Switch
	controls := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 8)
	controls.Add(
		components.NewCheckbox("Accept terms and conditions"),
		components.NewRadio("Select this option"),
		components.NewSwitch(),
	)

	// Select
	selectEl := components.NewSelect()
	selectEl.SetOptions([]string{"Apple", "Banana", "Cherry", "Date"})
	selectEl.SetPlaceholder("Select a fruit...")

	// Toggle & ToggleGroup
	tg := components.NewToggleGroup([]string{"Bold", "Italic", "Underline"})
	tg.SetSingle(false)

	// Slider
	slider := components.NewSlider(0, 100)

	return components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 16).
		SetPadding(yoga.EdgeAll, 24).
		Add(
			demoSection("Button", demoCard("Button Variants", btnRow)),
			demoSection("Input & Textarea", demoCard("Form Inputs", inputGroup)),
			demoSection("Checkbox / Radio / Switch", demoCard("Form Controls", controls)),
			demoSection("Select", demoCard("Select Component", selectEl)),
			demoSection("Toggle Group", demoCard("Toggle Group", tg)),
			demoSection("Slider", demoCard("Slider Component", slider)),
		)
}

type NavigationPage struct{ core.BaseWidget }

func NewNavigationPage() *NavigationPage {
	p := &NavigationPage{}
	p.Init(p)
	return p
}

func (p *NavigationPage) Render() core.Element {
	// Breadcrumb
	bc := components.NewBreadcrumb([]components.BreadcrumbItem{
		{Label: "Home", IsActive: false},
		{Label: "Products", IsActive: false},
		{Label: "Laptops", IsActive: true},
	})

	// Menubar
	mb := components.NewMenubar([]components.MenubarItem{
		{Label: "File"},
		{Label: "Edit"},
		{Label: "View"},
		{Label: "Help"},
	})

	// Pagination
	pg := components.NewPagination(20)

	// Tabs
	tabs := components.NewTab([]components.TabItem{
		{Label: "Account", Content: func() core.Element {
			return components.NewText("Account settings content.").SetColor(core.GetTheme().TextColor)
		}},
		{Label: "Password", Content: func() core.Element {
			return components.NewText("Password settings content.").SetColor(core.GetTheme().TextColor)
		}},
		{Label: "Notifications", Content: func() core.Element {
			return components.NewText("Notification settings content.").SetColor(core.GetTheme().TextColor)
		}},
	})

	return components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 16).
		SetPadding(yoga.EdgeAll, 24).
		Add(
			demoSection("Breadcrumb", demoCard("Breadcrumb Component", bc)),
			demoSection("Menubar", demoCard("Menubar Component", mb)),
			demoSection("Pagination", demoCard("Pagination Component", pg)),
			demoSection("Tabs", demoCard("Tabs Component", tabs)),
		)
}

type OverlaysPage struct {
	core.BaseWidget
	alertDialog *components.AlertDialog
	drawer      *components.Drawer
	sheet       *components.Sheet
}

func NewOverlaysPage() *OverlaysPage {
	p := &OverlaysPage{}
	p.Init(p)
	return p
}

func (p *OverlaysPage) Render() core.Element {
	// AlertDialog
	p.alertDialog = components.NewAlertDialog()
	p.alertDialog.SetTitle("Are you absolutely sure?").
		SetDescription("This action cannot be undone. This will permanently delete your account and remove your data from our servers.").
		SetActionText("Continue").
		SetCancelText("Cancel")

	alertBtn := components.NewButton("Show Alert Dialog")
	alertBtn.SetOnClick(func() { p.alertDialog.Open() })

	// Popover
	pop := components.NewPopover("Open popover")
	pop.AddContent(
		components.NewText("This is a popover content.").SetColor(core.GetTheme().TextColor),
		components.NewButton("Close").SetVariant(components.ButtonGhost).SetOnClick(func() { pop.Close() }),
	)

	// HoverCard
	hc := components.NewHoverCard(components.NewButton("Hover me").SetVariant(components.ButtonLink))
	hc.AddContent(
		components.NewText("@nextjs").SetFontSize(14).SetColor(core.GetTheme().TextColor),
		components.NewText("The React Framework for the Web.").SetFontSize(12).SetColor(core.GetTheme().MutedForegroundColor),
	)

	// Drawer
	p.drawer = components.NewDrawer(components.DrawerRight)
	p.drawer.Panel().Add(
		components.NewText("Drawer Content").SetFontSize(18).SetColor(core.GetTheme().TextColor),
		components.NewText("You can put any content here.").SetFontSize(14).SetColor(core.GetTheme().MutedForegroundColor),
	)
	drawerBtn := components.NewButton("Open Drawer").SetVariant(components.ButtonOutline)
	drawerBtn.SetOnClick(func() { p.drawer.Open() })

	// Sheet
	p.sheet = components.NewSheet(components.SheetBottom)
	p.sheet.Panel().Add(
		components.NewText("Sheet Content").SetFontSize(18).SetColor(core.GetTheme().TextColor),
		components.NewText("This is a sheet panel.").SetFontSize(14).SetColor(core.GetTheme().MutedForegroundColor),
	)
	sheetBtn := components.NewButton("Open Sheet").SetVariant(components.ButtonOutline)
	sheetBtn.SetOnClick(func() { p.sheet.Open() })

	// Tooltip
	ttBtn := components.NewButton("Hover for Tooltip")

	return components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 16).
		SetPadding(yoga.EdgeAll, 24).
		Add(
			demoSection("Alert Dialog", demoCard("Alert Dialog", alertBtn)),
			demoSection("Popover", demoCard("Popover", pop)),
			demoSection("Hover Card", demoCard("Hover Card", hc)),
			demoSection("Drawer", demoCard("Drawer", drawerBtn)),
			demoSection("Sheet", demoCard("Sheet", sheetBtn)),
			demoSection("Tooltip", demoCard("Tooltip", ttBtn)),
			p.alertDialog,
			p.drawer,
			p.sheet,
		)
}

type FeedbackPage struct{ core.BaseWidget }

func NewFeedbackPage() *FeedbackPage {
	p := &FeedbackPage{}
	p.Init(p)
	return p
}

func (p *FeedbackPage) Render() core.Element {
	// Alert variants
	alertDefault := components.NewAlert("Heads up!").
		SetDescription("You can add components to your app using the CLI.")
	alertDestructive := components.NewAlert("Error").
		SetDescription("Your session has expired. Please log in again.").
		SetVariant(components.AlertDestructive)

	// Progress
	pb := components.NewProgressBar().SetProgress(0.6)

	// Skeleton
	skel := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 8)
	skel.Add(
		components.NewSkeleton().SetWidth(200),
		components.NewSkeleton().SetWidth(150),
		components.NewSkeleton().SetWidth(180),
	)

	// Toast
	toastManager := components.NewToastManager()
	toastBtn := components.NewButton("Show Toast").SetVariant(components.ButtonOutline)
	toastBtn.SetOnClick(func() {
		toastManager.Show("Event has been created.", 3000)
	})

	return components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 16).
		SetPadding(yoga.EdgeAll, 24).
		Add(
			func() core.Element {
				alertWrap := components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).SetGap(yoga.GutterAll, 8)
				alertWrap.Add(alertDefault, alertDestructive)
				return demoSection("Alert", demoCard("Alert Variants", alertWrap))
			}(),
			demoSection("Progress", demoCard("Progress Bar", pb)),
			demoSection("Skeleton", demoCard("Skeleton Placeholder", skel)),
			func() core.Element {
				toastWrap := components.NewView().SetFlexDirection(yoga.FlexDirectionColumn).SetGap(yoga.GutterAll, 8)
				toastWrap.Add(toastBtn, toastManager)
				return demoSection("Toast", demoCard("Toast Notification", toastWrap))
			}(),
		)
}

type DataDisplayPage struct{ core.BaseWidget }

func NewDataDisplayPage() *DataDisplayPage {
	p := &DataDisplayPage{}
	p.Init(p)
	return p
}

func (p *DataDisplayPage) Render() core.Element {
	// Avatar
	avatars := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetAlignItems(yoga.AlignCenter)
	avatars.Add(
		components.NewAvatar("JD").SetSize(40),
		components.NewAvatar("AB").SetSize(48),
		components.NewAvatar("CD").SetSize(56),
	)

	// Badge variants
	badges := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetAlignItems(yoga.AlignCenter)
	badges.Add(
		components.NewBadge("default").SetVariant(components.BadgeDefault),
		components.NewBadge("secondary").SetVariant(components.BadgeSecondary),
		components.NewBadge("outline").SetVariant(components.BadgeOutline),
		components.NewBadge("destructive").SetVariant(components.BadgeDestructive),
	)

	// Calendar
	cal := components.NewCalendar().SetDate(time.Now())

	// Table
	tbl := components.NewTable()
	tbl.SetHeaders([]string{"Invoice", "Status", "Method", "Amount"})
	tbl.AddRow([]string{"INV001", "Paid", "Credit Card", "$250.00"})
	tbl.AddRow([]string{"INV002", "Pending", "PayPal", "$150.00"})
	tbl.AddRow([]string{"INV003", "Unpaid", "Bank Transfer", "$350.00"})
	tbl.AddRow([]string{"INV004", "Paid", "Credit Card", "$450.00"})

	// Carousel
	items := []core.Element{
		components.NewView().SetBackgroundColor(color.RGBA{R: 255, G: 200, B: 200, A: 255}).SetHeight(150).SetWidthPercent(100),
		components.NewView().SetBackgroundColor(color.RGBA{R: 200, G: 255, B: 200, A: 255}).SetHeight(150).SetWidthPercent(100),
		components.NewView().SetBackgroundColor(color.RGBA{R: 200, G: 200, B: 255, A: 255}).SetHeight(150).SetWidthPercent(100),
	}
	carousel := components.NewCarousel(items)

	return components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 16).
		SetPadding(yoga.EdgeAll, 24).
		Add(
			demoSection("Avatar", demoCard("Avatar Component", avatars)),
			demoSection("Badge", demoCard("Badge Variants", badges)),
			demoSection("Calendar", demoCard("Calendar Component", cal)),
			demoSection("Table", demoCard("Table Component", tbl)),
			demoSection("Carousel", demoCard("Carousel Component", carousel)),
		)
}

// ========== App ==========

type GalleryApp struct {
	core.BaseWidget
	page    *core.State[string]
	sidebar *components.Sidebar
}

func NewGalleryApp() *GalleryApp {
	a := &GalleryApp{page: core.NewState("layout")}
	a.Init(a)
	return a
}

func (a *GalleryApp) Render() core.Element {
	// Sidebar
	a.sidebar = components.NewSidebar([]components.SidebarNavItem{
		{Label: "Components", Children: []components.SidebarNavItem{
			{Label: "Layout", Page: "layout"},
			{Label: "Forms", Page: "forms"},
			{Label: "Navigation", Page: "navigation"},
			{Label: "Overlays", Page: "overlays"},
			{Label: "Feedback", Page: "feedback"},
			{Label: "Data Display", Page: "data-display"},
		}},
	})
	a.sidebar.SetSelected(a.page.Get())
	a.sidebar.SetOnSelect(func(page string) {
		a.page.Set(page)
	})

	// Main content
	content := core.NewSwitch(a.page).
		Case("layout", func() core.Element { return NewLayoutPage().Render() }).
		Case("forms", func() core.Element { return NewFormsPage().Render() }).
		Case("navigation", func() core.Element { return NewNavigationPage().Render() }).
		Case("overlays", func() core.Element { return NewOverlaysPage().Render() }).
		Case("feedback", func() core.Element { return NewFeedbackPage().Render() }).
		Case("data-display", func() core.Element { return NewDataDisplayPage().Render() }).
		Default(func() core.Element {
			return components.NewText("Select a page from the sidebar").SetColor(core.GetTheme().TextColor)
		})

	// Scrollable main area - wrap in a flex container to prevent overflow
	scroll := components.NewScrollView()
	scroll.SetFlexGrow(1)
	scroll.SetHeightPercent(100)
	scroll.Content().Add(content)

	// Root layout
	root := components.NewView()
	root.SetWidthPercent(100)
	root.SetHeightPercent(100)
	root.SetFlexDirection(yoga.FlexDirectionRow)
	root.SetBackgroundColor(core.GetTheme().BackgroundColor)
	root.Add(a.sidebar, scroll)

	return root
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("failed to init font: " + err.Error())
	}
	if err := fonts.LoadFontFromFile(fonts.FontFamilySans, "font/OPPOSans-Medium.ttf"); err == nil {
		fonts.SetDefaultFontFamily(fonts.FontFamilySans)
	}

	core.SetTheme(core.DefaultShadcnLightTheme())

	app := NewGalleryApp()
	tenon.RunWithDebug(app, 1200, 800, 8765)
}
