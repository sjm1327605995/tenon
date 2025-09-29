package main

import (
	"fmt"
	"github.com/millken/yoga"
)

func main() {
	//c := yoga.ConfigNew()

	//c.SetUseWebDefaults(true)
	root := yoga.NewNode()
	root.StyleSetWidth(500)
	root.StyleSetHeight(500)
	root.SetDirtiedFunc(func(node *yoga.Node) {
		fmt.Println("root")
	})
	n := yoga.NewNode()

	n.SetDirtiedFunc(func(node *yoga.Node) {
		fmt.Println("n")
	})
	n.StyleSetAlignItems(yoga.AlignCenter)
	n.StyleSetJustifyContent(yoga.JustifyCenter)
	root.InsertChild(n, 0)
	child := yoga.NewNode()
	child.StyleSetWidth(100)
	child.StyleSetHeight(100)
	child.SetDirtiedFunc(func(node *yoga.Node) {
		fmt.Println("child")
	})
	n.InsertChild(child, 0)

	yoga.CalculateLayout(n, yoga.Undefined, yoga.Undefined, yoga.DirectionInherit)
	fmt.Println(child.LayoutLeft())
	child.StyleSetWidth(200)
	yoga.CalculateLayout(n, yoga.Undefined, yoga.Undefined, yoga.DirectionInherit)
	fmt.Println(child.LayoutLeft())
}
