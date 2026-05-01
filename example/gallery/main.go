package main

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/fonts"
)

type galleryApp struct {
	counter int

	// Checkbox
	checked1, checked2 bool

	// Switch
	switch1, switch2 bool

	// Radio
	radioIdx int

	// Slider
	sliderVal float32

	// Tabs
	activeTab int

	// Breadcrumb
	breadcrumbIdx int

	// Pagination
	page int

	// Accordion
	expanded map[int]bool

	// Toggle
	toggle1, toggle2 bool

	// Textarea
	textarea string

	// Calendar
	selectedDate time.Time

	// AnimatedContainer
	acExpanded bool
}

func (g *galleryApp) Build() tenon.Widget {
	t := tenon.GetTheme()
	muted := t.TextMutedColor

	if g.expanded == nil {
		g.expanded = map[int]bool{0: true}
	}
	if g.selectedDate.IsZero() {
		g.selectedDate = time.Now()
	}

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
				NewCounterButton(0),
			).Gapf(8),
			tenon.Row(
				tenon.Button(fmt.Sprintf("Global Counter: %d", g.counter)).OnTap(func() {
					g.counter++
					tenon.Rebuild()
				}),
				tenon.Button("Reset Global").Variantf(tenon.ButtonOutline).OnTap(func() {
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
			tenon.Column(
				tenon.Checkbox("Accept terms and conditions", g.checked1, func(v bool) {
					g.checked1 = v
					tenon.Rebuild()
				}),
				tenon.Checkbox("Subscribe to newsletter", g.checked2, func(v bool) {
					g.checked2 = v
					tenon.Rebuild()
				}),
			).Gapf(4),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Switch
			sectionTitle("Switch"),
			tenon.Column(
				tenon.Switch(g.switch1, func(v bool) {
					g.switch1 = v
					tenon.Rebuild()
				}),
				tenon.Switch(g.switch2, func(v bool) {
					g.switch2 = v
					tenon.Rebuild()
				}),
			).Gapf(4),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Radio
			sectionTitle("Radio"),
			tenon.Column(
				tenon.Radio("Option A", g.radioIdx == 0, func(v bool) {
					if v { g.radioIdx = 0; tenon.Rebuild() }
				}),
				tenon.Radio("Option B", g.radioIdx == 1, func(v bool) {
					if v { g.radioIdx = 1; tenon.Rebuild() }
				}),
				tenon.Radio("Option C", g.radioIdx == 2, func(v bool) {
					if v { g.radioIdx = 2; tenon.Rebuild() }
				}),
			).Gapf(4),

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
				tenon.Text(fmt.Sprintf("Current page: %d", g.page)).FontSize(14).Color(muted),
				tenon.Pagination(g.page, 10, func(p int) {
					g.page = p
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
				tenon.Textarea(g.textarea, func(v string) {
					g.textarea = v
					tenon.Rebuild()
				}),
				tenon.Text(fmt.Sprintf("Length: %d", len(g.textarea))).FontSize(12).Color(muted),
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

			// AnimatedContainer
			sectionTitle("AnimatedContainer"),
			tenon.Column(
				tenon.Row(
					tenon.NewAnimatedContainer().
						WithChild(tenon.Text("Tap Toggle").FontSize(14).Color(tenon.NewColor(255, 255, 255, 255))).
						WithSize(func() (float32, float32) {
							if g.acExpanded {
								return 300, 80
							}
							return 200, 50
						}()).
						WithBackground(func() tenon.Color {
							if g.acExpanded {
								return *tenon.NewColor(239, 68, 68, 255)
							}
							return *tenon.NewColor(59, 130, 246, 255)
						}()).
						WithRadius(func() float32 {
							if g.acExpanded {
								return 16
							}
							return 8
						}()).
						WithDuration(300 * time.Millisecond).
						WithCurve(tenon.EaseInOutCurve{}),
				).AlignItems(tenon.AlignFlexStart),
				tenon.Button("Toggle Size & Color").OnTap(func() {
					g.acExpanded = !g.acExpanded
					tenon.Rebuild()
				}),
			).Gapf(12),

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

// ========== CounterButton: StatefulWidget Demo ==========

type CounterButton struct {
	tenon.BaseWidget
	initial int
}

func NewCounterButton(initial int) CounterButton {
	return CounterButton{initial: initial}
}

func (c CounterButton) CreateElement() tenon.Element {
	return tenon.NewStatefulElement(c)
}

func (c CounterButton) CreateState() tenon.State {
	s := &counterState{}
	s.Init(s)
	return s
}

type counterState struct {
	tenon.BaseState
	count int
}

func (s *counterState) InitState() {
	s.count = s.GetWidget().(CounterButton).initial
}

func (s *counterState) Build(ctx tenon.BuildContext) tenon.Widget {
	return tenon.Row(
		tenon.Button(fmt.Sprintf("Stateful Counter: %d", s.count)).OnTap(func() {
			s.SetState(func() {
				s.count++
			})
		}),
		tenon.Button("Reset").Variantf(tenon.ButtonOutline).OnTap(func() {
			s.SetState(func() {
				s.count = s.GetWidget().(CounterButton).initial
			})
		}),
	).Gapf(8)
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("failed to init font: " + err.Error())
	}
	app := &galleryApp{counter: 0}
	tenon.SetTheme(tenon.DefaultLightTheme())
	tenon.Run(app.Build, 900, 800)
}
