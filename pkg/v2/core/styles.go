package core

// StyleFunc applies styles to an element.
type StyleFunc func(Element)

var (
	globalStyles     = make(map[string]StyleFunc)
	typedStyles      = make(map[string]map[string]StyleFunc) // elemType -> tag -> func
)

// RegisterStyle registers a global style by tag.
// Any element with SetTag(tag) or SetClass(tag) will have this style applied.
func RegisterStyle(tag string, apply StyleFunc) {
	globalStyles[tag] = apply
}

// RegisterStyleForType registers a style for a specific element type.
func RegisterStyleForType(elemType, tag string, apply StyleFunc) {
	if typedStyles[elemType] == nil {
		typedStyles[elemType] = make(map[string]StyleFunc)
	}
	typedStyles[elemType][tag] = apply
}

// applyStyles applies registered styles to an element based on its tags/classes.
func applyStyles(el Element) {
	// Apply global styles by tag
	if el.GetTag() != "" {
		if fn, ok := globalStyles[el.GetTag()]; ok {
			fn(el)
		}
	}
	// Apply global styles by class
	for _, class := range el.GetClass() {
		if fn, ok := globalStyles[class]; ok {
			fn(el)
		}
	}
	// Apply type-specific styles
	if m, ok := typedStyles[el.ElementType()]; ok {
		if el.GetTag() != "" {
			if fn, ok := m[el.GetTag()]; ok {
				fn(el)
			}
		}
		for _, class := range el.GetClass() {
			if fn, ok := m[class]; ok {
				fn(el)
			}
		}
	}
	// Recursively apply to children
	for _, child := range el.GetChildren() {
		applyStyles(child)
	}
}
