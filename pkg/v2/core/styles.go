package core

// StyleFunc applies styles to an element.
type StyleFunc func(Element)

// StyleRegistry 持有样式注册表。
type StyleRegistry struct {
	globalStyles map[string]StyleFunc
	typedStyles  map[string]map[string]StyleFunc // elemType -> tag -> func
}

// NewStyleRegistry 创建样式注册表。
func NewStyleRegistry() *StyleRegistry {
	return &StyleRegistry{
		globalStyles: make(map[string]StyleFunc),
		typedStyles:  make(map[string]map[string]StyleFunc),
	}
}

// RegisterStyle registers a global style by tag.
func (r *StyleRegistry) RegisterStyle(tag string, apply StyleFunc) {
	if r.globalStyles == nil {
		r.globalStyles = make(map[string]StyleFunc)
	}
	r.globalStyles[tag] = apply
}

// RegisterStyleForType registers a style for a specific element type.
func (r *StyleRegistry) RegisterStyleForType(elemType, tag string, apply StyleFunc) {
	if r.typedStyles == nil {
		r.typedStyles = make(map[string]map[string]StyleFunc)
	}
	if r.typedStyles[elemType] == nil {
		r.typedStyles[elemType] = make(map[string]StyleFunc)
	}
	r.typedStyles[elemType][tag] = apply
}

// ApplyStyles applies registered styles to an element based on its tags/classes.
func (r *StyleRegistry) ApplyStyles(el Element) {
	styled, isStyled := el.(StyledElement)

	// Apply global styles by tag
	if isStyled && styled.GetTag() != "" {
		if fn, ok := r.globalStyles[styled.GetTag()]; ok {
			fn(el)
		}
	}
	// Apply global styles by class
	if isStyled {
		for _, class := range styled.GetClass() {
			if fn, ok := r.globalStyles[class]; ok {
				fn(el)
			}
		}
	}
	// Apply type-specific styles
	if m, ok := r.typedStyles[el.ElementType()]; ok && isStyled {
		if styled.GetTag() != "" {
			if fn, ok := m[styled.GetTag()]; ok {
				fn(el)
			}
		}
		for _, class := range styled.GetClass() {
			if fn, ok := m[class]; ok {
				fn(el)
			}
		}
	}
	// Recursively apply to children
	for _, child := range el.GetChildren() {
		r.ApplyStyles(child)
	}
}

// Clone 返回样式注册表的深拷贝。
func (r *StyleRegistry) Clone() *StyleRegistry {
	nr := NewStyleRegistry()
	for k, v := range r.globalStyles {
		nr.globalStyles[k] = v
	}
	for t, m := range r.typedStyles {
		nm := make(map[string]StyleFunc)
		for k, v := range m {
			nm[k] = v
		}
		nr.typedStyles[t] = nm
	}
	return nr
}
