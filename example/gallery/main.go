package main

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/fonts"
)

type galleryApp struct {
	counter      int
	checkbox1    bool
	checkbox2    bool
	switch1      bool
	switch2      bool
	radioIdx     int
	sliderVal    float32
	activeTab    int
	currentPage  int
	toggle1       bool
	toggle2       bool
	breadcrumbIdx int
	expanded      map[int]bool
	textareaVal  string
	selectedDate time.Time
}

func newGalleryApp() *galleryApp {
	return &galleryApp{
		expanded:      map[int]bool{0: true},
		selectedDate:  time.Now(),
		sliderVal:     50,
		currentPage:   1,
		breadcrumbIdx: 2,
	}
}

func (g *galleryApp) Build() tenon.Widget {
	t := tenon.GetTheme()
	muted := t.TextMutedColor

	return tenon.Scroll(
		tenon.Column(
			// Header
			tenon.Text("Tenon UI Gallery").FontSize(28).Color(t.TextColor),
			tenon.Text("Shadcn-style components for Go").FontSize(14).Color(muted),
			tenon.Separator(tenon.SeparatorHorizontal),

			// Badge
			sectionTitle("Badge"),
			tenon.Row(
				tenon.Badge("Default", tenon.BadgeDefault),
				tenon.Badge("Secondary", tenon.BadgeSecondary),
				tenon.Badge("Outline", tenon.BadgeOutline),
				tenon.Badge("Destructive", tenon.BadgeDestructive),
				tenon.DotBadge(),
				tenon.CountBadge(42),
			).Gapf(8),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Button
			sectionTitle("Button"),
			tenon.Row(
				tenon.Button("Default").OnTap(func() {}),
				tenon.Button("Secondary").Variantf(tenon.ButtonSecondary).OnTap(func() {}),
				tenon.Button("Outline").Variantf(tenon.ButtonOutline).OnTap(func() {}),
				tenon.Button("Ghost").Variantf(tenon.ButtonGhost).OnTap(func() {}),
				tenon.Button("Destructive").Variantf(tenon.ButtonDestructive).OnTap(func() {}),
				tenon.Button("Link").Variantf(tenon.ButtonLink).OnTap(func() {}),
			).Gapf(8),
			tenon.Row(
				tenon.Button(fmt.Sprintf("Clicked %d", g.counter)).OnTap(func() {
					g.counter++
						tenon.Rebuild()
					}),
				tenon.Button("Reset").Variantf(tenon.ButtonOutline).OnTap(func() {
					g.counter = 0
						tenon.Rebuild()
					}),
			).Gapf(8),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Avatar
			sectionTitle("Avatar"),
			tenon.Row(
				tenon.Avatar("AB", 40),
				tenon.Avatar("CD", 32),
				tenon.Avatar("E", 48),
			).Gapf(8).AlignItems(tenon.AlignCenter),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Card
			sectionTitle("Card"),
			tenon.Card(
				"Beautiful Card",
				"A card component with title, description and actions.",
				[]tenon.Widget{
					tenon.Text("This is the card body content. You can put any widgets here.").FontSize(14).Color(t.TextColor),
				},
				[]tenon.Widget{
					tenon.Button("Cancel").Variantf(tenon.ButtonOutline).OnTap(func() {}),
					tenon.Button("Save").OnTap(func() {}),
				},
			),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Alert
			sectionTitle("Alert"),
			tenon.Alert("Heads up!", "You can add components to your app using the cli.", tenon.AlertDefault),
			tenon.Alert("Error", "Something went wrong. Please try again.", tenon.AlertDestructive),

			tenon.Separator(tenon.SeparatorHorizontal),

			// ProgressBar
			sectionTitle("ProgressBar"),
			tenon.ProgressBar(0.25),
			tenon.ProgressBar(0.60),
			tenon.ProgressBar(1.0),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Skeleton
			sectionTitle("Skeleton"),
			tenon.Row(
				tenon.Column(
					tenon.Skeleton(200, 20),
					tenon.Skeleton(160, 14),
					tenon.Skeleton(120, 14),
				).Gapf(8),
			),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Checkbox
			sectionTitle("Checkbox"),
			tenon.Checkbox("Accept terms and conditions", g.checkbox1, func(v bool) {
				g.checkbox1 = v
				tenon.Rebuild()
			}),
			tenon.Checkbox("Subscribe to newsletter", g.checkbox2, func(v bool) {
				g.checkbox2 = v
				tenon.Rebuild()
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Switch
			sectionTitle("Switch"),
			tenon.Switch(g.switch1, func(v bool) {
				g.switch1 = v
				tenon.Rebuild()
			}),
			tenon.Switch(g.switch2, func(v bool) {
				g.switch2 = v
				tenon.Rebuild()
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Radio
			sectionTitle("Radio"),
			tenon.Radio("Option A", g.radioIdx == 0, func(v bool) { if v { g.radioIdx = 0 } else { g.radioIdx = -1 }; tenon.Rebuild() }),
			tenon.Radio("Option B", g.radioIdx == 1, func(v bool) { if v { g.radioIdx = 1 } else { g.radioIdx = -1 }; tenon.Rebuild() }),
			tenon.Radio("Option C", g.radioIdx == 2, func(v bool) { if v { g.radioIdx = 2 } else { g.radioIdx = -1 }; tenon.Rebuild() }),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Slider
			sectionTitle("Slider"),
			tenon.Column(
				tenon.Text(fmt.Sprintf("Value: %.0f", g.sliderVal)).FontSize(14).Color(muted),
				tenon.Slider(0, 100, g.sliderVal, func(v float32) {
					g.sliderVal = v
						tenon.Rebuild()
					}),
			).Gapf(4),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Tabs
			sectionTitle("Tabs"),
			tenon.Tabs(
				[]tenon.TabItem{
					{Label: "Account", Content: tenon.Text("Manage your account settings.").FontSize(14)},
					{Label: "Password", Content: tenon.Text("Change your password.").FontSize(14)},
					{Label: "Notifications", Content: tenon.Text("Set notification preferences.").FontSize(14)},
				},
				g.activeTab,
				func(idx int) { g.activeTab = idx; tenon.Rebuild() },
			),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Breadcrumb
			sectionTitle("Breadcrumb"),
			tenon.Breadcrumb(
				[]tenon.BreadcrumbItem{
					{Label: "Home", IsActive: g.breadcrumbIdx == 0, OnTap: func() { g.breadcrumbIdx = 0; tenon.Rebuild() }},
					{Label: "Products", IsActive: g.breadcrumbIdx == 1, OnTap: func() { g.breadcrumbIdx = 1; tenon.Rebuild() }},
					{Label: "Electronics", IsActive: g.breadcrumbIdx == 2, OnTap: func() { g.breadcrumbIdx = 2; tenon.Rebuild() }},
				},
				" / ",
			),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Pagination
			sectionTitle("Pagination"),
			tenon.Column(
				tenon.Text(fmt.Sprintf("Current page: %d", g.currentPage)).FontSize(14).Color(muted),
				tenon.Pagination(g.currentPage, 10, func(p int) {
					g.currentPage = p
						tenon.Rebuild()
					}),
			).Gapf(4),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Table
			sectionTitle("Table"),
			tenon.Table(
				[]string{"Name", "Role", "Status"},
				[][]string{
					{"Alice", "Admin", "Active"},
					{"Bob", "Editor", "Active"},
					{"Charlie", "Viewer", "Inactive"},
				},
			),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Accordion
			sectionTitle("Accordion"),
			tenon.Accordion(
				[]tenon.AccordionPane{
					{Title: "Is it accessible?", Content: tenon.Text("Yes. It adheres to the WAI-ARIA design pattern.").FontSize(14)},
					{Title: "Is it styled?", Content: tenon.Text("Yes. It comes with default styles that match the other components.").FontSize(14)},
					{Title: "Is it animated?", Content: tenon.Text("Yes. It's animated by default, but you can disable it.").FontSize(14)},
				},
				g.expanded,
				func(idx int) {
					g.expanded[idx] = !g.expanded[idx]
						tenon.Rebuild()
					},
			),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Toggle
			sectionTitle("Toggle"),
			tenon.Row(
				tenon.Toggle("Bold", g.toggle1, func(v bool) {
					g.toggle1 = v
						tenon.Rebuild()
					}),
				tenon.Toggle("Italic", g.toggle2, func(v bool) {
					g.toggle2 = v
						tenon.Rebuild()
					}),
			).Gapf(8),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Textarea
			sectionTitle("Textarea"),
			tenon.Column(
				tenon.Textarea(g.textareaVal, func(v string) {
					g.textareaVal = v
						tenon.Rebuild()
					}),
				tenon.Text(fmt.Sprintf("Length: %d", len(g.textareaVal))).FontSize(12).Color(muted),
			).Gapf(4),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Calendar
			sectionTitle("Calendar"),
			tenon.Column(
				tenon.Text(fmt.Sprintf("Selected: %s", g.selectedDate.Format("2006-01-02"))).FontSize(14).Color(muted),
				tenon.Calendar(g.selectedDate.Year(), g.selectedDate.Month(), g.selectedDate, func(d time.Time) {
					g.selectedDate = d
						tenon.Rebuild()
					}),
			).Gapf(4),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Separator demo
			sectionTitle("Separator"),
			tenon.Row(
				tenon.Text("Left").FontSize(14),
				tenon.Separator(tenon.SeparatorVertical),
				tenon.Text("Right").FontSize(14),
			).Gapf(8),

		).Gapf(20).Paddingf(tenon.EdgeInsetsAll(24)),
	)
}

func sectionTitle(text string) tenon.Widget {
	return tenon.Label(text)
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("failed to init font: " + err.Error())
	}
	app := newGalleryApp()
	tenon.SetTheme(tenon.DefaultLightTheme())
	tenon.Run(app.Build, 900, 800)
}


