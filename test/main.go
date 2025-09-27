package main

import (
	"fmt"

	"github.com/millken/yoga"
)

func main() {
	c := yoga.ConfigNew()
	c.SetUseWebDefaults(true)
	n := yoga.NewNodeWithConfig(c)
	n.StyleSetWidth(100)
	n.StyleSetBorder(yoga.EdgeAll, 100)
	yoga.CalculateLayout(n, yoga.Undefined, yoga.Undefined, yoga.DirectionInherit)
	fmt.Println(n.StyleGetWidth())
}
