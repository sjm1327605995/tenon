package main

import (
	"image/color"

	comp "github.com/sjm1327605995/tenon/react/component"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/yoga"
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

func (h *HomePage) ShouldComponentUpdate() {

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
	return comp.NewView().
		WidthPercent(100).
		HeightPercent(100).
		JustifyContent(yoga.JustifyCenter).
		AlignItems(yoga.AlignCenter).
		Background(color.NRGBA{R: 0xff, A: 0xff})
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
