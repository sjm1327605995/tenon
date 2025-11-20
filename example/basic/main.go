package main

import (
	comp "github.com/sjm1327605995/tenon/react/component"
	"github.com/sjm1327605995/tenon/react/core"
)

// HomePage is a simple page component that uses the View component
type HomePage struct {
}

// Constructor is called when the component is created
func (h *HomePage) Constructor() {
}

// GetDerivedStateFromProps is called when props are updated
func (h *HomePage) GetDerivedStateFromProps(props, state any) {
}

func (h *HomePage) ShouldComponentUpdate() bool {
	// Default implementation returns true to always update
	return true
}

// Render returns the UI tree for the component
func (h *HomePage) Render() core.Node {
	// Create Image component
	//imageComponent := comp.NewImage().
	//	Src("react.svg")
	//
	//// Create inner View component
	//innerView := comp.NewView().
	//	WidthPercent(80).
	//	HeightPercent(80).
	//	JustifyContent(yoga.JustifyCenter).
	//	AlignItems(yoga.AlignCenter).
	//	Background(color.NRGBA{R: 189, G: 193, B: 193, A: 0xff}).
	//	Body(imageComponent)

	// Create outer View component and return it
	//return comp.NewView().
	//	WidthPercent(50).
	//	HeightPercent(50).
	//	BorderWidth(10).
	//	BorderColor(color.NRGBA{
	//		G: 0xff,
	//		B: 0,
	//		A: 0xff,
	//	}).
	//	JustifyContent(yoga.JustifyCenter).
	//	AlignItems(yoga.AlignCenter).
	//	Background(color.NRGBA{R: 0xff, A: 0xff}).
	//	Body(comp.NewImage().
	//		Width(100).Height(100).Src("react.svg"), comp.NewText("test"))
	return comp.NewText("test")

}

// GetSnapshotBeforeUpdate is called before updating the DOM
func (h *HomePage) GetSnapshotBeforeUpdate() {
}

// ComponentDidMount is called after the component is mounted
func (h *HomePage) ComponentDidMount() {
}

// ComponentDidUpdate is called after the component is updated
func (h *HomePage) ComponentDidUpdate() {
}

// ComponentWillUnmount is called before the component is unmounted
func (h *HomePage) ComponentWillUnmount() {
}

func main() {
	// Create a new engine instance
	e := core.NewEngine()

	// Add a route for the home page
	e.AddRoute(core.Route{
		Path:        "/",
		ComponentFn: func() core.Component { return &HomePage{} },
	})

	// Navigate to the home page
	e.Navigate("/")

	// Run the engine
	e.Run()
}
