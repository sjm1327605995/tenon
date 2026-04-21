package main

import (
	"fmt"
	"image/color"

	"github.com/sjm1327605995/tenon/internal/renderer"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

func main() {
	// 初始化字体管理器
	if err := fonts.InitDefaultFont(); err != nil {
		panic("Failed to initialize default fonts: " + err.Error())
	}

	// 加载自定义字体
	if err := fonts.LoadFontFromFile(fonts.FontFamilySans, "font/OPPOSans-Medium.ttf"); err != nil {
		fmt.Println("Warning: Failed to load OPPOSans-Medium.ttf:", err.Error())
		fmt.Println("Using default font instead")
	} else {
		fmt.Println("Successfully loaded OPPOSans-Medium.ttf")
		// 设置自定义字体为默认字体
		fonts.SetDefaultFontFamily(fonts.FontFamilySans)
	}

	root := buildUI()
	renderer.Run(root, 800, 600)
}

func buildUI() *components.View {
	return components.NewView().
		SetWidth(800).
		SetHeight(600).
		SetBackgroundColor(color.RGBA{R: 240, G: 240, B: 240, A: 255}).
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetPadding(yoga.EdgeAll, 20).
		Add(
			components.NewText("Tenon UI Framework - 链式编程演示").
				SetFontSize(24).
				SetColor(color.RGBA{R: 51, G: 51, B: 51, A: 255}).
				SetMargin(yoga.EdgeBottom, 30),
			components.NewView().
				SetFlexDirection(yoga.FlexDirectionColumn).
				SetPadding(yoga.EdgeAll, 20).
				SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}).
				SetBorderRadius(12).
				SetBorder(yoga.EdgeAll, 1).
				SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
				Add(
					components.NewText("字体测试").
						SetFontSize(18).
						SetColor(color.RGBA{R: 33, G: 37, B: 41, A: 255}).
						SetMargin(yoga.EdgeBottom, 20),
					components.NewText("12px - 这是 OPPO Sans 字体").
						SetFontSize(12).
						SetColor(color.RGBA{R: 108, G: 117, B: 125, A: 255}).
						SetMargin(yoga.EdgeBottom, 10),
					components.NewText("16px - 这是 OPPO Sans 字体").
						SetFontSize(16).
						SetColor(color.RGBA{R: 73, G: 80, B: 87, A: 255}).
						SetMargin(yoga.EdgeBottom, 10),
					components.NewText("20px - 这是 OPPO Sans 字体").
						SetFontSize(20).
						SetColor(color.RGBA{R: 33, G: 37, B: 41, A: 255}).
						SetMargin(yoga.EdgeBottom, 10),
					components.NewText("24px - 这是 OPPO Sans 字体").
						SetFontSize(24).
						SetColor(color.RGBA{R: 0, G: 123, B: 255, A: 255}),
				),
			components.NewView().
				SetFlexDirection(yoga.FlexDirectionColumn).
				SetPadding(yoga.EdgeAll, 20).
				SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}).
				SetBorderRadius(12).SetBorder(yoga.EdgeAll, 1).
				SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255}).
				SetMargin(yoga.EdgeTop, 20).
				Add(
					components.NewText("按钮组件演示").
						SetFontSize(18).
						SetColor(color.RGBA{R: 33, G: 37, B: 41, A: 255}).
						SetMargin(yoga.EdgeBottom, 20),
					components.NewView().
						SetFlexDirection(yoga.FlexDirectionRow).
						SetJustifyContent(yoga.JustifySpaceAround).
						SetPadding(yoga.EdgeAll, 20).
						Add(
							components.NewButton("普通按钮").
								SetWidth(120).
								SetHeight(40).
								SetOnClick(func() {
									fmt.Println("普通按钮被点击了!")
								}),
							components.NewButton("禁用按钮").
								SetWidth(120).
								SetHeight(40).
								SetDisabled(true),
							components.NewButton("自定义按钮").
								SetWidth(120).
								SetHeight(40).
								SetBackgroundColors(
									color.RGBA{R: 40, G: 167, B: 69, A: 255}, // 正常颜色
									color.RGBA{R: 33, G: 136, B: 56, A: 255}, // 悬停颜色
									color.RGBA{R: 25, G: 105, B: 44, A: 255}, // 按下颜色
								).SetOnClick(func() {
								fmt.Println("自定义按钮被点击了!")
							}),
						)))

}
