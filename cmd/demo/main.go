package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/internal/renderer"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/pkg/hooks"
	"github.com/sjm1327605995/tenon/yoga"
)

// Counter 函数组件 - React 19 风格
func Counter(props map[string]interface{}) core.Component {
	count, setCount := hooks.UseState(0)
	isPending, startTransition := hooks.UseTransition()

	// useEffect 示例
	hooks.UseEffect(func() func() {
		fmt.Printf("Effect: count changed to %d\n", count)
		return func() {
			fmt.Printf("Effect cleanup: count was %d\n", count)
		}
	}, []interface{}{count})

	// useId 示例
	id := hooks.UseId()

	view := components.NewView()
	view.SetPadding(yoga.EdgeAll, 20)
	view.SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	view.SetBorderRadius(12)
	view.SetBorder(yoga.EdgeAll, 1)
	view.SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255})
	view.SetMargin(yoga.EdgeBottom, 20)

	// 标题文本
	title := components.NewText(fmt.Sprintf("Counter %s", id))
	title.SetFontSize(20)
	title.SetColor(color.RGBA{R: 51, G: 51, B: 51, A: 255})
	title.SetMargin(yoga.EdgeBottom, 10)

	// 计数显示
	countText := components.NewText(fmt.Sprintf("Count: %d", count))
	countText.SetFontSize(18)
	countText.SetColor(color.RGBA{R: 33, G: 37, B: 41, A: 255})
	countText.SetMargin(yoga.EdgeBottom, 15)

	// 过渡状态文本
	pendingText := components.NewText("")
	if isPending {
		pendingText.SetContent("Updating...")
		pendingText.SetColor(color.RGBA{R: 255, G: 165, B: 0, A: 255})
	}
	pendingText.SetMargin(yoga.EdgeBottom, 10)

	// 增加按钮
	incrementBtn := components.NewButton("Increment")
	incrementBtn.SetWidth(150)
	incrementBtn.SetHeight(40)
	incrementBtn.SetMargin(yoga.EdgeBottom, 10)
	incrementBtn.SetOnClick(func() {
		setCount(count + 1)
	})

	// 使用 startTransition 的按钮
	transitionBtn := components.NewButton("Slow Increment")
	transitionBtn.SetWidth(150)
	transitionBtn.SetHeight(40)
	transitionBtn.SetOnClick(func() {
		startTransition(func() {
			setCount(count + 1)
		})
	})

	view.Add(title, countText, pendingText, incrementBtn, transitionBtn)
	return view
}

// TodoList 函数组件 - useOptimistic 示例
func TodoList(props map[string]interface{}) core.Component {
	todos, setTodos := hooks.UseState([]string{"Learn Tenon", "Build UI"})

	// useOptimistic 示例
	optimisticTodos, addOptimisticTodo := hooks.UseOptimistic(
		todos,
		func(currentState []string, optimisticValue interface{}) []string {
			newTodo := optimisticValue.(string)
			return append([]string{newTodo + " (adding...)"}, currentState...)
		},
	)

	view := components.NewView()
	view.SetPadding(yoga.EdgeAll, 20)
	view.SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	view.SetBorderRadius(12)
	view.SetBorder(yoga.EdgeAll, 1)
	view.SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255})

	title := components.NewText("Todo List (useOptimistic)")
	title.SetFontSize(20)
	title.SetColor(color.RGBA{R: 51, G: 51, B: 51, A: 255})
	title.SetMargin(yoga.EdgeBottom, 15)
	view.AddChild(title)

	// 渲染列表
	for _, todo := range optimisticTodos {
		todoText := components.NewText("- " + todo)
		todoText.SetFontSize(16)
		todoText.SetColor(color.RGBA{R: 66, G: 66, B: 66, A: 255})
		todoText.SetMargin(yoga.EdgeBottom, 5)
		view.AddChild(todoText)
	}

	// 添加按钮
	addBtn := components.NewButton("Add Todo")
	addBtn.SetWidth(150)
	addBtn.SetHeight(40)
	addBtn.SetMargin(yoga.EdgeTop, 15)
	addBtn.SetOnClick(func() {
		newTodo := fmt.Sprintf("Todo %d", len(todos)+1)
		addOptimisticTodo(newTodo)
		setTodos(append([]string{newTodo}, todos...))
	})
	view.AddChild(addBtn)

	return view
}

// ActionDemo 函数组件 - useActionState 示例
func ActionDemo(props map[string]interface{}) core.Component {
	name, submitAction, isPending, err := hooks.UseActionState(
		func(prevState string, formData map[string]string) (string, error) {
			inputName := formData["name"]
			if inputName == "" {
				return prevState, fmt.Errorf("name is required")
			}
			return fmt.Sprintf("Hello, %s!", inputName), nil
		},
		"",
	)

	view := components.NewView()
	view.SetPadding(yoga.EdgeAll, 20)
	view.SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	view.SetBorderRadius(12)
	view.SetBorder(yoga.EdgeAll, 1)
	view.SetBorderColor(color.RGBA{R: 222, G: 226, B: 230, A: 255})
	view.SetMargin(yoga.EdgeTop, 20)

	title := components.NewText("Action State Demo")
	title.SetFontSize(20)
	title.SetColor(color.RGBA{R: 51, G: 51, B: 51, A: 255})
	title.SetMargin(yoga.EdgeBottom, 15)
	view.AddChild(title)

	resultText := components.NewText(name)
	resultText.SetFontSize(16)
	resultText.SetColor(color.RGBA{R: 33, G: 37, B: 41, A: 255})
	resultText.SetMargin(yoga.EdgeBottom, 10)
	view.AddChild(resultText)

	if isPending {
		pendingText := components.NewText("Submitting...")
		pendingText.SetColor(color.RGBA{R: 255, G: 165, B: 0, A: 255})
		pendingText.SetMargin(yoga.EdgeBottom, 10)
		view.AddChild(pendingText)
	}

	if err != nil {
		errText := components.NewText("Error: "+err.Error())
		errText.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
		errText.SetMargin(yoga.EdgeBottom, 10)
		view.AddChild(errText)
	}

	// 模拟表单提交按钮
	submitBtn := components.NewButton("Submit 'World'")
	submitBtn.SetWidth(150)
	submitBtn.SetHeight(40)
	submitBtn.SetOnClick(func() {
		submitAction(map[string]string{"name": "World"})
	})
	view.AddChild(submitBtn)

	return view
}

func buildUI() core.Component {
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

	title := components.NewText("Tenon UI Framework - React 19 Style")
	title.SetFontSize(24)
	title.SetColor(color.RGBA{R: 51, G: 51, B: 51, A: 255})
	title.SetMargin(yoga.EdgeBottom, 30)

	// 创建函数组件的包装组件
	counterWrapper := components.NewFunctionComponent(Counter, nil)
	todoWrapper := components.NewFunctionComponent(TodoList, nil)
	actionWrapper := components.NewFunctionComponent(ActionDemo, nil)

	root.Add(title, counterWrapper, todoWrapper, actionWrapper)
	return root
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
