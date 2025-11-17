package main

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/yoga"

	"gioui.org/layout"
)

// TestComponent is a simple component for testing async updates
type TestComponent struct {
	count int
	core.Component
}

func (t *TestComponent) Constructor() {
	fmt.Println("Constructor called")
	t.count = 0
}

func (t *TestComponent) GetDerivedStateFromProps(props, state any) {
	// Not implemented yet
}

func (t *TestComponent) ShouldComponentUpdate() bool {
	// Always return true for testing
	return true
}

func (t *TestComponent) Render() core.Node {
	fmt.Println("Render called, count:", t.count)

	// Create a simple view with some text
	// This is a simplified example - in real usage, you'd use actual components
	return &simpleNode{text: fmt.Sprintf("Count: %d", t.count)}
}

func (t *TestComponent) GetSnapshotBeforeUpdate() {
	// Not implemented yet
}

func (t *TestComponent) ComponentDidMount() {
	fmt.Println("ComponentDidMount called")
}

func (t *TestComponent) ComponentDidUpdate() {
	fmt.Println("ComponentDidUpdate called")
}

func (t *TestComponent) ComponentWillUnmount() {
	fmt.Println("ComponentWillUnmount called")
}

// simpleNode is a simplified implementation of core.Node for testing
type simpleNode struct {
	text     string
	children []core.Node
	yogaNode *yoga.Node
}

func (n *simpleNode) Yoga() *yoga.Node {
	if n.yogaNode == nil {
		n.yogaNode = yoga.NewNode()
	}
	return n.yogaNode
}

func (n *simpleNode) Children() []core.Node {
	return n.children
}

func (n *simpleNode) Update(ctx layout.Context) {
	// Not implemented yet
}

func (n *simpleNode) Gio() core.Gio {
	// Return a simple implementation of core.Gio
	return &simpleGio{text: n.text}
}

// simpleGio is a simplified implementation of core.Gio for testing
type simpleGio struct {
	text string
}

func (s *simpleGio) Layout(gtx layout.Context) layout.Dimensions {
	// For testing purposes, just return a minimal size
	return layout.Dimensions{Size: gtx.Constraints.Min}
}

func main() {
	// Create a new engine
	engine := core.NewEngine()

	// Add a route for our test component
	engine.AddRoute(core.Route{
		Path: "/",
		ComponentFn: func() core.Component {
			return &TestComponent{}
		},
	})

	// Navigate to the route
	engine.Navigate("/")

	// Simulate an async update
	go func() {
		// Wait for a short time
		time.Sleep(time.Second * 2)

		// Update the component state
		engine.Update(func() {
			fmt.Println("Async update executed")
			// Note: We can't directly access engine.CurrentPage because it's not public
			// In a real implementation, you'd have a way to access the current component
		})

		// Wait for another short time and update again
		time.Sleep(time.Second * 2)
		engine.Update(func() {
			fmt.Println("Second async update executed")
		})
	}()

	// Start the engine
	fmt.Println("Starting engine...")
	engine.Run()
}
