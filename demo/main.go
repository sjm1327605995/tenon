package main

import (
	"fmt"
	"os"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/yoga"
)

type CounterProps struct {
	InitialCount int
}

func Counter(props CounterProps) tenon.UI {
	return tenon.View(&tenon.ViewProps{
		Style: &tenon.ViewStyle{
			Width:          tenon.Px(200),
			Height:         tenon.Px(100),
			Background:     "#ffffff",
			FlexDirection:  yoga.FlexDirectionColumn,
			JustifyContent: yoga.JustifySpaceAround,
			AlignItems:     yoga.AlignCenter,
			Padding:        tenon.Px(20),
		},
	},
		tenon.Text(&tenon.TextProps{
			Content: "Counter Component",
			Style: &tenon.TextStyle{
				FontSize: tenon.Px(16),
				Color:    "#333333",
			},
		}),

		tenon.View(&tenon.ViewProps{
			Style: &tenon.ViewStyle{
				Width:          tenon.Px(120),
				Height:         tenon.Px(40),
				Background:     "#f0f0f0",
				FlexDirection:  yoga.FlexDirectionRow,
				JustifyContent: yoga.JustifyCenter,
				AlignItems:     yoga.AlignCenter,
			},
		},
			tenon.Text(&tenon.TextProps{
				Content: fmt.Sprintf("Count: %d", props.InitialCount),
				Style: &tenon.TextStyle{
					FontSize: tenon.Px(14),
					Color:    "#666666",
				},
			}),
		),
	)
}

func Dashboard() tenon.UI {
	return tenon.View(&tenon.ViewProps{
		Style: &tenon.ViewStyle{
			Width:          tenon.Px(400),
			Height:         tenon.Px(300),
			Background:     "#ffffff",
			FlexDirection:  yoga.FlexDirectionColumn,
			Padding:        tenon.Px(20),
		},
	},
		tenon.View(&tenon.ViewProps{
			Style: &tenon.ViewStyle{
				Height:         tenon.Px(50),
				FlexDirection:  yoga.FlexDirectionRow,
				JustifyContent: yoga.JustifySpaceBetween,
				AlignItems:     yoga.AlignCenter,
				MarginBottom:   tenon.Px(20),
			},
		},
			tenon.Text(&tenon.TextProps{
				Content: "Dashboard",
				Style: &tenon.TextStyle{
					FontSize: tenon.Px(20),
					Color:    "#495057",
				},
			}),

			tenon.Text(&tenon.TextProps{
				Content: "Welcome back!",
				Style: &tenon.TextStyle{
					FontSize: tenon.Px(14),
					Color:    "#868e96",
				},
			}),
		),

		tenon.View(&tenon.ViewProps{
			Style: &tenon.ViewStyle{
				FlexDirection: yoga.FlexDirectionRow,
				FlexGrow:      1,
			},
		},
			tenon.View(&tenon.ViewProps{
				Style: &tenon.ViewStyle{
					Width:       tenon.Px(200),
					MarginRight: tenon.Px(20),
				},
			},
				tenon.View(&tenon.ViewProps{
					Style: &tenon.ViewStyle{
						Width:          tenon.Px(300),
						Height:         tenon.Px(150),
						Background:     "#f8f9fa",
						FlexDirection:  yoga.FlexDirectionRow,
						Padding:        tenon.Px(15),
					},
				},
					tenon.View(&tenon.ViewProps{
						Style: &tenon.ViewStyle{
							Width:      tenon.Px(80),
							Height:     tenon.Px(80),
							Background: "#e9ecef",
							Margin:     tenon.Px(10),
						},
					},
						tenon.Text(&tenon.TextProps{
							Content: "Avatar",
							Style: &tenon.TextStyle{
								FontSize: tenon.Px(12),
								Color:    "#6c757d",
							},
						}),
					),

					tenon.View(&tenon.ViewProps{
						Style: &tenon.ViewStyle{
							FlexDirection:  yoga.FlexDirectionColumn,
							JustifyContent: yoga.JustifySpaceAround,
							FlexGrow:       1,
						},
					},
						tenon.Text(&tenon.TextProps{
							Content: "John Doe",
							Style: &tenon.TextStyle{
								FontSize: tenon.Px(18),
								Color:    "#212529",
							},
						}),

						tenon.Text(&tenon.TextProps{
							Content: "Software Engineer",
							Style: &tenon.TextStyle{
								FontSize: tenon.Px(14),
								Color:    "#6c757d",
							},
						}),

						tenon.Text(&tenon.TextProps{
							Content: "San Francisco, CA",
							Style: &tenon.TextStyle{
								FontSize: tenon.Px(12),
								Color:    "#adb5bd",
							},
						}),
					),
				),
			),

			tenon.View(&tenon.ViewProps{
				Style: &tenon.ViewStyle{
					FlexGrow:      1,
					FlexDirection: yoga.FlexDirectionColumn,
					JustifyContent: yoga.JustifyCenter,
				},
			},
				Counter(CounterProps{InitialCount: 42}),
			),
		),
	)
}

func main() {
	fmt.Println("=== Tenon React19 Style Demo ===")
	fmt.Println()

	reconciler := tenon.NewReconciler(Dashboard())
	rootElement := reconciler.GetRootElement()

	html := tenon.RenderToHTML(rootElement)
	os.WriteFile("debug.html", []byte(html), 0644)
	fmt.Println("HTML debug file saved to debug.html")

	fmt.Println("\nDemo completed successfully!")
}