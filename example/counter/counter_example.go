package main

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/yoga"

	"gioui.org/layout"
)

// CounterComponent is a component that maintains a counter state
type CounterComponent struct {
	count        int
	stateChanged bool
}

func (c *CounterComponent) Constructor() {
	fmt.Println("CounterComponent: Constructor called")
	c.count = 0
	c.stateChanged = false
}

func (c *CounterComponent) GetDerivedStateFromProps(props, state any) {
	// Not implemented yet
}

func (c *CounterComponent) ShouldComponentUpdate() bool {
	shouldUpdate := c.stateChanged
	c.stateChanged = false
	return shouldUpdate
}

func (c *CounterComponent) Render() core.Node {
	count := c.count
	fmt.Printf("CounterComponent: Render called, count: %d\n", count)
	return &counterNode{count: count}
}

func (c *CounterComponent) GetSnapshotBeforeUpdate() {
	// Not implemented yet
}

func (c *CounterComponent) ComponentDidMount() {
	fmt.Println("CounterComponent: ComponentDidMount called")
}

func (c *CounterComponent) ComponentDidUpdate() {
	fmt.Printf("CounterComponent: ComponentDidUpdate called, count: %d\n", c.count)
}

func (c *CounterComponent) ComponentWillUnmount() {
	fmt.Println("CounterComponent: ComponentWillUnmount called")
}

// SetState updates the component's state
func (c *CounterComponent) SetState(newCount int) {
	c.count = newCount
	c.stateChanged = true
}

// counterNode is a simplified implementation of core.Node for the counter
type counterNode struct {
	count    int
	yogaNode *yoga.Node
}

func (n *counterNode) Yoga() *yoga.Node {
	if n.yogaNode == nil {
		n.yogaNode = yoga.NewNode()
	}
	return n.yogaNode
}

func (n *counterNode) Children() []core.Node {
	return nil
}

func (n *counterNode) Update(ctx layout.Context) {
	// Not implemented yet
}

func (n *counterNode) Gio() core.Gio {
	return &counterGio{text: fmt.Sprintf("Count: %d", n.count)}
}

// counterGio is a simplified implementation of core.Gio for the counter
type counterGio struct {
	text string
}

func (g *counterGio) Layout(gtx layout.Context) layout.Dimensions {
	// For this example, we'll just print the text and return a minimal size
	fmt.Printf("counterGio.Layout: Rendering %s\n", g.text)
	return layout.Dimensions{Size: gtx.Constraints.Min}
}

func main() {
	// Create a new engine
	engine := core.NewEngine()

	// Create a counter component instance
	counter := &CounterComponent{}

	// Add a route for our counter component
	engine.AddRoute(core.Route{
		Path: "/",
		ComponentFn: func() core.Component {
			return counter
		},
	})

	// Navigate to the route
	engine.Navigate("/")

	// Simulate async updates to the component state
	go func() {
		for i := 1; i <= 3; i++ {
			// Wait for 2 seconds
			time.Sleep(time.Second * 2)

			// Update the component state
			engine.Update(func() {
				fmt.Printf("main: Enqueueing update %d\n", i)
				counter.SetState(i)
			})
		}

		// After 3 updates, exit the program
		time.Sleep(time.Second * 2)
		fmt.Println("main: Exiting program")
		// Note: In a real Gio application, we'd close the window properly
	}()

	// Start the engine
	fmt.Println("main: Starting engine")
	engine.Run()
}
