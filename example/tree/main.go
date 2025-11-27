package main

import (
	"fmt"
	"github.com/sjm1327605995/tenon/core/vdom"
)

func main() {
	t1 := &vdom.Tree{}
	t1.Children = append(t1.Children, &vdom.Element{
		Name: "div",
		Attrs: []vdom.Attr{
			{Name: "id", Value: "1"},
		},
	})
	t2 := &vdom.Tree{}
	t2.Children = append(t2.Children, &vdom.Element{
		Name: "div",
		Attrs: []vdom.Attr{
			{Name: "id", Value: "1"},
			{Name: "class", Value: "text"},
		},
	})

	list, err := vdom.Diff(t1, t2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(list)
}
