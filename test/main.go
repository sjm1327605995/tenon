package main

import (
	"fmt"
	"github.com/sjm1327605995/tenon/yoga"
)

func main() {
	root := yoga.NewNode()
	root.StyleSetWidth(300)
	root.StyleSetHeight(200)
	n1 := yoga.NewNode()
	n1.StyleSetWidth(100)
	n1.StyleSetHeight(100)
	n2 := yoga.NewNode()
	n2.StyleSetWidth(100)
	n2.StyleSetHeight(100)
	root.InsertChild(n1, 0)
	root.InsertChild(n2, 1)
	fmt.Println(root.StyleGetWidth())
	root.CalculateLayout(800, 600, yoga.DirectionInherit)
	fmt.Println(root.LayoutWidth())
}
