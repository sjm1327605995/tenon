package main

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/fonts"
)

type galleryApp struct {
	counter int
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
			// CounterButton: StatefulWidget demo, state encapsulated inside component
			tenon.Row(
				NewCounterButton(0),
			).Gapf(8),
			// Traditional manual Rebuild mode still works
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
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				checked1, checked2 := false, false
				return tenon.Column(
					tenon.Checkbox("Accept terms and conditions", checked1, func(v bool) {
						setState(func() { checked1 = v })
					}),
					tenon.Checkbox("Subscribe to newsletter", checked2, func(v bool) {
						setState(func() { checked2 = v })
					}),
				).Gapf(4)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Switch
			sectionTitle("Switch"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				s1, s2 := false, false
				return tenon.Column(
					tenon.Switch(s1, func(v bool) {
						setState(func() { s1 = v })
					}),
					tenon.Switch(s2, func(v bool) {
						setState(func() { s2 = v })
					}),
				).Gapf(4)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Radio
			sectionTitle("Radio"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				radioIdx := 0
				return tenon.Column(
					tenon.Radio("Option A", radioIdx == 0, func(v bool) {
						if v { setState(func() { radioIdx = 0 }) }
					}),
					tenon.Radio("Option B", radioIdx == 1, func(v bool) {
						if v { setState(func() { radioIdx = 1 }) }
					}),
					tenon.Radio("Option C", radioIdx == 2, func(v bool) {
						if v { setState(func() { radioIdx = 2 }) }
					}),
				).Gapf(4)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Slider
			sectionTitle("Slider"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				sliderVal := float32(50)
				return tenon.Column(
					tenon.Text(fmt.Sprintf("Value: %.0f", sliderVal)).FontSize(14).Color(muted),
					tenon.Slider(0, 100, sliderVal, func(v float32) {
						setState(func() { sliderVal = v })
					}),
				).Gapf(4)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Tabs
			sectionTitle("Tabs"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				activeTab := 0
				return tenon.Tabs(
					[]tenon.TabItem{
						{Label: "Account", Content: tenon.Text("Manage your account settings.").FontSize(14)},
						{Label: "Password", Content: tenon.Text("Change your password.").FontSize(14)},
						{Label: "Notifications", Content: tenon.Text("Set notification preferences.").FontSize(14)},
					},
					activeTab,
					func(idx int) { setState(func() { activeTab = idx }) },
				)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Breadcrumb
			sectionTitle("Breadcrumb"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				idx := 2
				return tenon.Breadcrumb(
					[]tenon.BreadcrumbItem{
						{Label: "Home", IsActive: idx == 0, OnTap: func() { setState(func() { idx = 0 }) }},
						{Label: "Products", IsActive: idx == 1, OnTap: func() { setState(func() { idx = 1 }) }},
						{Label: "Electronics", IsActive: idx == 2, OnTap: func() { setState(func() { idx = 2 }) }},
					},
					" / ",
				)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Pagination
			sectionTitle("Pagination"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				page := 1
				return tenon.Column(
					tenon.Text(fmt.Sprintf("Current page: %d", page)).FontSize(14).Color(muted),
					tenon.Pagination(page, 10, func(p int) {
						setState(func() { page = p })
					}),
				).Gapf(4)
			}),

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
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				expanded := map[int]bool{0: true}
				return tenon.Accordion(
					[]tenon.AccordionPane{
						{Title: "Is it accessible?", Content: tenon.Text("Yes. It adheres to the WAI-ARIA design pattern.").FontSize(14)},
						{Title: "Is it styled?", Content: tenon.Text("Yes. It comes with default styles that match the other components.").FontSize(14)},
						{Title: "Is it animated?", Content: tenon.Text("Yes. It's animated by default, but you can disable it.").FontSize(14)},
					},
					expanded,
					func(idx int) {
						setState(func() {
							expanded[idx] = !expanded[idx]
						})
					},
				)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Toggle
			sectionTitle("Toggle"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				t1, t2 := false, false
				return tenon.Row(
					tenon.Toggle("Bold", t1, func(v bool) {
						setState(func() { t1 = v })
					}),
					tenon.Toggle("Italic", t2, func(v bool) {
						setState(func() { t2 = v })
					}),
				).Gapf(8)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Textarea
			sectionTitle("Textarea"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				text := ""
				return tenon.Column(
					tenon.Textarea(text, func(v string) {
						setState(func() { text = v })
					}),
					tenon.Text(fmt.Sprintf("Length: %d", len(text))).FontSize(12).Color(muted),
				).Gapf(4)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// Calendar
			sectionTitle("Calendar"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				selected := time.Now()
				return tenon.Column(
					tenon.Text(fmt.Sprintf("Selected: %s", selected.Format("2006-01-02"))).FontSize(14).Color(muted),
					tenon.Calendar(selected.Year(), selected.Month(), selected, func(d time.Time) {
						setState(func() { selected = d })
					}),
				).Gapf(4)
			}),

			tenon.Separator(tenon.SeparatorHorizontal),

			// AnimatedContainer
			sectionTitle("AnimatedContainer"),
			tenon.NewStatefulBuilder(func(ctx tenon.BuildContext, setState func(fn func())) tenon.Widget {
				expanded := false
				var w, h float32 = 200, 50
				bg := *tenon.NewColor(59, 130, 246, 255)
				r := float32(8)
				if expanded {
					w, h = 300, 80
					bg = *tenon.NewColor(239, 68, 68, 255)
					r = 16
				}
				return tenon.Column(
					tenon.Row(
						tenon.NewAnimatedContainer().
							WithChild(tenon.Text("Tap Toggle").FontSize(14).Color(tenon.NewColor(255, 255, 255, 255))).
							WithSize(w, h).
							WithBackground(bg).
							WithRadius(r).
							WithDuration(300 * time.Millisecond).
							WithCurve(tenon.EaseInOutCurve{}),
					).AlignItems(tenon.AlignFlexStart),
					tenon.Button("Toggle Size & Color").OnTap(func() {
						setState(func() { expanded = !expanded })
					}),
				).Gapf(12)
			}),

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
