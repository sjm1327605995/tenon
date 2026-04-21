package tenon

import (
	"fmt"
	"os"
	"testing"

	"github.com/sjm1327605995/tenon/internal/renderer"
	"github.com/sjm1327605995/tenon/internal/reconciler"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/types"
	"github.com/sjm1327605995/tenon/yoga"
)

func saveHTMLFile(filename, html string) {
	os.MkdirAll("test_output", 0755)
	os.WriteFile(filename, []byte(html), 0644)
}

func contains(html, substr string) bool {
	return len(html) > 0 && len(substr) > 0 && len(html) >= len(substr) &&
		(findSubstring(html, substr) >= 0)
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestFlexDirection(t *testing.T) {
	flexDirectionCases := []struct {
		name   string
		dir    yoga.FlexDirection
		layout string
	}{
		{"Row", yoga.FlexDirectionRow, "flex-direction: row;"},
		{"Column", yoga.FlexDirectionColumn, "flex-direction: column;"},
		{"RowReverse", yoga.FlexDirectionRowReverse, "flex-direction: row-reverse;"},
		{"ColumnReverse", yoga.FlexDirectionColumnReverse, "flex-direction: column-reverse;"},
	}

	for _, tc := range flexDirectionCases {
		t.Run(tc.name, func(t *testing.T) {
			component := components.View(&types.ViewProps{
				Style: &types.ViewStyle{
					Width:         types.Px(300),
					Height:        types.Px(200),
					FlexDirection: tc.dir,
				},
			},
				components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(50), Background: "#ff0000"}}),
				components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(50), Background: "#00ff00"}}),
				components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(50), Background: "#0000ff"}}),
			)

			r := reconciler.NewReconciler(component)
			rootElement := r.GetRootElement()

			html := renderer.RenderToHTML(rootElement)

			if !contains(html, tc.layout) {
				t.Errorf("Expected HTML to contain %s, got:\n%s", tc.layout, html)
			}

			saveHTMLFile(fmt.Sprintf("test_output/flex_direction_%s.html", tc.name), html)
		})
	}
}

func TestJustifyContent(t *testing.T) {
	justifyCases := []struct {
		name   string
		justify yoga.Justify
		layout  string
	}{
		{"FlexStart", yoga.JustifyFlexStart, "justify-content: flex-start;"},
		{"FlexEnd", yoga.JustifyFlexEnd, "justify-content: flex-end;"},
		{"Center", yoga.JustifyCenter, "justify-content: center;"},
		{"SpaceBetween", yoga.JustifySpaceBetween, "justify-content: space-between;"},
		{"SpaceAround", yoga.JustifySpaceAround, "justify-content: space-around;"},
		{"SpaceEvenly", yoga.JustifySpaceEvenly, "justify-content: space-evenly;"},
	}

	for _, tc := range justifyCases {
		t.Run(tc.name, func(t *testing.T) {
			component := components.View(&types.ViewProps{
				Style: &types.ViewStyle{
					Width:          types.Px(300),
					Height:         types.Px(100),
					JustifyContent: tc.justify,
				},
			},
				components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(50), Background: "#ff0000"}}),
				components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(50), Background: "#00ff00"}}),
				components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(50), Background: "#0000ff"}}),
			)

			r := reconciler.NewReconciler(component)
			rootElement := r.GetRootElement()

			html := renderer.RenderToHTML(rootElement)

			if !contains(html, tc.layout) {
				t.Errorf("Expected HTML to contain %s, got:\n%s", tc.layout, html)
			}

			saveHTMLFile(fmt.Sprintf("test_output/justify_%s.html", tc.name), html)
		})
	}
}

func TestAlignItems(t *testing.T) {
	alignCases := []struct {
		name   string
		align  yoga.Align
		layout string
	}{
		{"Stretch", yoga.AlignStretch, "align-items: stretch;"},
		{"FlexStart", yoga.AlignFlexStart, "align-items: flex-start;"},
		{"FlexEnd", yoga.AlignFlexEnd, "align-items: flex-end;"},
		{"Center", yoga.AlignCenter, "align-items: center;"},
		{"Baseline", yoga.AlignBaseline, "align-items: baseline;"},
	}

	for _, tc := range alignCases {
		t.Run(tc.name, func(t *testing.T) {
			component := components.View(&types.ViewProps{
				Style: &types.ViewStyle{
					Width:       types.Px(300),
					Height:      types.Px(100),
					AlignItems: tc.align,
				},
			},
				components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(50), Background: "#ff0000"}}),
				components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(50), Background: "#00ff00"}}),
				components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(50), Background: "#0000ff"}}),
			)

			r := reconciler.NewReconciler(component)
			rootElement := r.GetRootElement()

			html := renderer.RenderToHTML(rootElement)

			if !contains(html, tc.layout) {
				t.Errorf("Expected HTML to contain %s, got:\n%s", tc.layout, html)
			}

			saveHTMLFile(fmt.Sprintf("test_output/align_%s.html", tc.name), html)
		})
	}
}

func TestNestedLayout(t *testing.T) {
	component := components.View(&types.ViewProps{
		Style: &types.ViewStyle{
			Width:         types.Px(400),
			Height:        types.Px(300),
			FlexDirection: yoga.FlexDirectionColumn,
		},
	},
		components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(400), Height: types.Px(100), Background: "#ff0000", MarginBottom: types.Px(10)}}),
		components.View(&types.ViewProps{Style: &types.ViewStyle{
			Width:         types.Px(400),
			Height:        types.Px(100),
			FlexDirection: yoga.FlexDirectionRow,
		},
		},
			components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(190), Height: types.Px(100), Background: "#00ff00", MarginRight: types.Px(10)}}),
			components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(200), Height: types.Px(100), Background: "#0000ff"}}),
		),
	)

	r := reconciler.NewReconciler(component)
	rootElement := r.GetRootElement()

	html := renderer.RenderToHTML(rootElement)
	saveHTMLFile("test_output/nested_layout.html", html)
}

func TestFlexGrow(t *testing.T) {
	component := components.View(&types.ViewProps{
		Style: &types.ViewStyle{
			Width:         types.Px(400),
			Height:        types.Px(100),
			FlexDirection: yoga.FlexDirectionRow,
		},
	},
		components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(100), Background: "#ff0000"}}),
		components.View(&types.ViewProps{Style: &types.ViewStyle{FlexGrow: 1, Height: types.Px(100), Background: "#00ff00"}}),
		components.View(&types.ViewProps{Style: &types.ViewStyle{Width: types.Px(50), Height: types.Px(100), Background: "#0000ff"}}),
	)

	r := reconciler.NewReconciler(component)
	rootElement := r.GetRootElement()

	html := renderer.RenderToHTML(rootElement)

	if !contains(html, "flex-grow: 1") {
		t.Errorf("Expected HTML to contain flex-grow: 1, got:\n%s", html)
	}

	saveHTMLFile("test_output/flex_grow.html", html)
}

func TestMargin(t *testing.T) {
	component := components.View(&types.ViewProps{
		Style: &types.ViewStyle{
			Width:  types.Px(300),
			Height: types.Px(200),
			Margin: types.Px(20),
		},
	})

	r := reconciler.NewReconciler(component)
	rootElement := r.GetRootElement()

	html := renderer.RenderToHTML(rootElement)

	if !contains(html, "margin: 20px") {
		t.Errorf("Expected HTML to contain margin: 20px, got:\n%s", html)
	}

	saveHTMLFile("test_output/margin.html", html)
}

func TestPadding(t *testing.T) {
	component := components.View(&types.ViewProps{
		Style: &types.ViewStyle{
			Width:    types.Px(300),
			Height:   types.Px(200),
			Padding:  types.Px(20),
		},
	})

	r := reconciler.NewReconciler(component)
	rootElement := r.GetRootElement()

	html := renderer.RenderToHTML(rootElement)

	if !contains(html, "padding: 20px") {
		t.Errorf("Expected HTML to contain padding: 20px, got:\n%s", html)
	}

	saveHTMLFile("test_output/padding.html", html)
}

func TestTextElement(t *testing.T) {
	component := components.Text(&types.TextProps{
		Content: "Hello, Tenon!",
		Style: &types.TextStyle{
			FontSize: types.Px(16),
			Color:    "#333333",
		},
	})

	r := reconciler.NewReconciler(component)
	rootElement := r.GetRootElement()

	html := renderer.RenderToHTML(rootElement)

	if !contains(html, "Hello, Tenon!") {
		t.Errorf("Expected HTML to contain Hello, Tenon!, got:\n%s", html)
	}

	saveHTMLFile("test_output/text_element.html", html)
}

func TestImageElement(t *testing.T) {
	component := components.Image(&types.ImageProps{
		Source: "test.png",
		Style: &types.ImageStyle{
			Width:  types.Px(200),
			Height: types.Px(150),
		},
	})

	r := reconciler.NewReconciler(component)
	rootElement := r.GetRootElement()

	html := renderer.RenderToHTML(rootElement)

	if !contains(html, `src="test.png"`) {
		t.Errorf("Expected HTML to contain src=\"test.png\", got:\n%s", html)
	}

	saveHTMLFile("test_output/image_element.html", html)
}

func TestComplexDashboard(t *testing.T) {
	component := components.View(&types.ViewProps{
		Style: &types.ViewStyle{
			Width:         types.Px(800),
			Height:        types.Px(600),
			Padding:       types.Px(20),
			FlexDirection: yoga.FlexDirectionColumn,
		},
	},
		components.View(&types.ViewProps{
			Style: &types.ViewStyle{
				Width:         types.Px(760),
				Height:        types.Px(100),
				Background:    "#2196F3",
				MarginBottom:  types.Px(20),
				JustifyContent: yoga.JustifyCenter,
				AlignItems:     yoga.AlignCenter,
			},
		}, components.Text(&types.TextProps{Content: "Dashboard Header", Style: &types.TextStyle{FontSize: types.Px(24), Color: "#ffffff"}})),
		components.View(&types.ViewProps{
			Style: &types.ViewStyle{
				Width:         types.Px(760),
				Height:        types.Px(400),
				FlexDirection: yoga.FlexDirectionRow,
			},
		},
			components.View(&types.ViewProps{
				Style: &types.ViewStyle{
					Width:         types.Px(370),
					Height:        types.Px(380),
					Background:    "#4CAF50",
					MarginRight:   types.Px(20),
					FlexDirection: yoga.FlexDirectionColumn,
					JustifyContent: yoga.JustifySpaceBetween,
				},
			},
				components.Text(&types.TextProps{Content: "Card 1", Style: &types.TextStyle{FontSize: types.Px(18), Color: "#ffffff"}}),
				components.Text(&types.TextProps{Content: "Description 1", Style: &types.TextStyle{FontSize: types.Px(14), Color: "#e0e0e0"}}),
			),
			components.View(&types.ViewProps{
				Style: &types.ViewStyle{
					Width:         types.Px(370),
					Height:        types.Px(380),
					Background:    "#FF9800",
					FlexDirection: yoga.FlexDirectionColumn,
					JustifyContent: yoga.JustifySpaceBetween,
				},
			},
				components.Text(&types.TextProps{Content: "Card 2", Style: &types.TextStyle{FontSize: types.Px(18), Color: "#ffffff"}}),
				components.Text(&types.TextProps{Content: "Description 2", Style: &types.TextStyle{FontSize: types.Px(14), Color: "#e0e0e0"}}),
			),
		),
	)

	r := reconciler.NewReconciler(component)
	rootElement := r.GetRootElement()

	html := renderer.RenderToHTML(rootElement)
	saveHTMLFile("test_output/complex_dashboard.html", html)
}