package core

// Route defines a route with path and component type
// The component type should be a constructor function that returns a Component

type Route struct {
	Path        string
	ComponentFn func() Component
}
