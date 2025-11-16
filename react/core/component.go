package core

// Component is the basic interface for all UI components
type Component interface {
	Constructor()

	GetDerivedStateFromProps(props, state any)

	ShouldComponentUpdate()
	Render() Node
	GetSnapshotBeforeUpdate()
	ComponentDidMount()
	ComponentDidUpdate()
	ComponentWillUnmount()
}
