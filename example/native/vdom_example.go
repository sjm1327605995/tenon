package main

import (
	"fmt"
	"log"

	"tenon/core/vdom"
	"tenon/react"
)

// 创建一个简单的vdom树示例
func createVDOMTree() *vdom.Element {
	// 创建根容器
	container := vdom.CreateElement("View", vdom.Props{
		"style": map[string]interface{}{
			"flexDirection": "column",
			"padding":       20,
		},
	})

	// 创建标题
	title := vdom.CreateElement("Text", vdom.Props{
		"style": map[string]interface{}{
			"fontSize":   24,
			"fontWeight": "bold",
			"marginBottom": 10,
		},
		"children": "React Native Style VDOM Example",
	})

	// 创建内容文本
	content := vdom.CreateElement("Text", vdom.Props{
		"style": map[string]interface{}{
			"fontSize": 16,
			"marginBottom": 20,
		},
		"children": "This example demonstrates how VDOM integrates with native components.",
	})

	// 创建一个可点击的文本元素
	button := vdom.CreateElement("Text", vdom.Props{
		"style": map[string]interface{}{
			"fontSize":   18,
			"color":      "blue",
			"padding":    10,
			"borderWidth": 1,
			"borderColor": "blue",
			"textAlign":  "center",
		},
		"children": "Click Me",
		"onPress": func(event map[string]interface{}) {
			fmt.Println("Button clicked!")
		},
	})

	// 添加子元素到容器
	container.AppendChild(title)
	container.AppendChild(content)
	container.AppendChild(button)

	return container
}

// 更新vdom树的示例函数
func updateVDOMTree(root *vdom.Element) {
	// 找到按钮元素并更新其文本
	if len(root.Children) > 2 {
		button := root.Children[2]
		if textElement, ok := button.(*vdom.Element); ok {
			textElement.Props["children"] = "Button Clicked!"
			// 更新样式以反映点击状态
			if style, ok := textElement.Props["style"].(map[string]interface{}); ok {
				style["color"] = "green"
				style["borderColor"] = "green"
			}
		}
	}
}

func main() {
	// 初始化渲染器
	nativeRenderer := react.NewNativeRenderer()
	react.SetNativeRenderer(nativeRenderer)

	// 初始化原生组件系统
	react.InitializeNativeComponents()

	// 创建vdom树
	rootElement := createVDOMTree()

	// 使用渲染器将vdom渲染到原生组件
	fmt.Println("Rendering VDOM to native components...")
	nativeComponent := nativeRenderer.Render(rootElement)

	// 模拟用户交互并更新vdom
	fmt.Println("\nSimulating user interaction...")
	updateVDOMTree(rootElement)

	// 重新渲染更新后的vdom
	fmt.Println("\nRe-rendering with updated VDOM...")
	nativeRenderer.Update(rootElement, nativeComponent)

	// 模拟组件卸载
	fmt.Println("\nSimulating component unmount...")
	nativeComponent.UnmountNativeView()
	nativeComponent.WillUnmount()

	fmt.Println("\nExample completed successfully!")
}