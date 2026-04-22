package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/internal/renderer"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

type LifecycleDemo struct {
	core.BaseComponent
	counter int
	label   *components.Text
}

func NewLifecycleDemo() *LifecycleDemo {
	l := &LifecycleDemo{
		BaseComponent: core.NewBaseComponent(),
		counter:       0,
	}
	l.Init(l)

	l.label = components.NewText("点击按钮计数: 0")
	l.label.SetFontSize(18)
	l.label.SetColor(color.RGBA{R: 33, G: 37, B: 41, A: 255})
	l.AddChild(l.label)

	return l
}

func (l *LifecycleDemo) ComponentDidMount() {
	l.BaseComponent.ComponentDidMount()
	fmt.Printf("[%s] LifecycleDemo 挂载\n", time.Now().Format("15:04:05"))
}

func (l *LifecycleDemo) ComponentDidUpdate(prevProps, prevState, snapshot interface{}) {
	l.BaseComponent.ComponentDidUpdate(prevProps, prevState, snapshot)
	fmt.Printf("[%s] LifecycleDemo 更新 - 前一个状态: %v\n", time.Now().Format("15:04:05"), prevState)
}

func (l *LifecycleDemo) ComponentWillUnmount() {
	l.BaseComponent.ComponentWillUnmount()
	fmt.Printf("[%s] LifecycleDemo 卸载\n", time.Now().Format("15:04:05"))
}

func (l *LifecycleDemo) Increment() {
	l.counter++
	l.label.SetContent(fmt.Sprintf("点击按钮计数: %d", l.counter))
	l.SetState(l.counter)
}

func (l *LifecycleDemo) Draw(screen *ebiten.Image) {
	l.BaseComponent.Draw(screen)
	for _, child := range l.GetChildren() {
		child.Draw(screen)
	}
}

func (l *LifecycleDemo) Update() error {
	return l.BaseComponent.Update()
}

func (l *LifecycleDemo) DrawOverlay(screen *ebiten.Image) {
	l.BaseComponent.DrawOverlay(screen)
	for _, child := range l.GetChildren() {
		child.DrawOverlay(screen)
	}
}

func (l *LifecycleDemo) HandleInput() bool {
	return l.BaseComponent.HandleInput()
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("Failed to initialize default fonts: " + err.Error())
	}

	if err := fonts.LoadFontFromFile(fonts.FontFamilySans, "font/OPPOSans-Medium.ttf"); err != nil {
		fmt.Println("Warning: Failed to load OPPOSans-Medium.ttf:", err.Error())
		fmt.Println("Using default font instead")
	} else {
		fmt.Println("Successfully loaded OPPOSans-Medium.ttf")
		fonts.SetDefaultFontFamily(fonts.FontFamilySans)
	}

	root := buildUI()
	renderer.Run(root, 800, 600)
}

func buildUI() *components.View {
	root := components.NewView()
	root.SetStyle(core.NewStyle().
		SetWidth(800).
		SetHeight(600).
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetPadding(yoga.EdgeAll, 20),
	)
	root.SetVisualStyle(core.NewVisualStyle().
		SetBackgroundColor(color.RGBA{R: 240, G: 240, B: 240, A: 255}),
	)

	title := components.NewText("Tenon UI Framework - 生命周期演示")
	title.SetFontSize(24)
	title.SetColor(color.RGBA{R: 51, G: 51, B: 51, A: 255})
	title.SetMargin(yoga.EdgeBottom, 30)

	lifecycleDemo := NewLifecycleDemo()
	lifecycleDemo.SetPadding(yoga.EdgeAll, 20)
	lifecycleDemo.SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	lifecycleDemo.SetBorderRadius(12)
	lifecycleDemo.SetBorder(yoga.EdgeAll, 1)
	lifecycleDemo.SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255})
	lifecycleDemo.SetMargin(yoga.EdgeBottom, 20)

	incrementBtn := components.NewButton("点击增加计数")
	incrementBtn.SetWidth(200)
	incrementBtn.SetHeight(44)
	incrementBtn.SetOnClick(func() {
		lifecycleDemo.Increment()
	})

	root.Add(title, lifecycleDemo, incrementBtn)
	return root
}
