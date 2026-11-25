package ui

import (
	"fmt"
	"sync"
	"time"
)

// Game3DComponent is a backend-neutral scene component with a shared lifecycle.
type Game3DComponent interface {
	Update(time.Duration)
	Close()
}

// Scene3DAdapter creates a native component from a portable node declaration.
type Scene3DAdapter interface {
	Adapt3D(*GameNode3D) (Game3DComponent, error)
}

// Scene3DAdapterFunc adapts Scene3DAdapter to a function.
type Scene3DAdapterFunc func(*GameNode3D) (Game3DComponent, error)

// Adapt3D calls f with the declaration.
func (f Scene3DAdapterFunc) Adapt3D(node *GameNode3D) (Game3DComponent, error) {
	return f(node)
}

// Scene3DValidator validates a plugin declaration before native adaptation.
// Adapters may implement it to reject invalid payloads before allocating resources.
type Scene3DValidator interface {
	Validate3D(*GameNode3D) error
}

type scene3DPluginKey struct {
	name    string
	version int
}

// Scene3DPluginRegistry resolves built-in declarations and third-party plugin nodes.
// Registries are independent: adding a plugin does not affect other scenes or backends.
type Scene3DPluginRegistry struct {
	mu      sync.RWMutex
	adapter Scene3DAdapter
	plugins map[scene3DPluginKey]Scene3DAdapter
}

// NewScene3DPluginRegistry creates an empty registry. A backend normally calls
// RegisterAdapter for its built-in nodes and RegisterPlugin for extensions.
func NewScene3DPluginRegistry() *Scene3DPluginRegistry {
	return &Scene3DPluginRegistry{plugins: make(map[scene3DPluginKey]Scene3DAdapter)}
}

// RegisterAdapter sets the fallback adapter for non-plugin nodes.
func (r *Scene3DPluginRegistry) RegisterAdapter(adapter Scene3DAdapter) {
	if adapter == nil {
		panic("ui: cannot register a nil 3D adapter")
	}
	r.mu.Lock()
	r.adapter = adapter
	r.mu.Unlock()
}

// RegisterPlugin registers an unversioned named adapter used by Plugin3D nodes.
func (r *Scene3DPluginRegistry) RegisterPlugin(name string, adapter Scene3DAdapter) {
	r.registerPlugin(name, 0, adapter)
}

// RegisterPluginVersion registers an adapter for one positive plugin version.
func (r *Scene3DPluginRegistry) RegisterPluginVersion(name string, version int, adapter Scene3DAdapter) {
	if version <= 0 {
		panic("ui: 3D plugin version must be positive")
	}
	r.registerPlugin(name, version, adapter)
}

func (r *Scene3DPluginRegistry) registerPlugin(name string, version int, adapter Scene3DAdapter) {
	if name == "" {
		panic("ui: cannot register an empty 3D plugin name")
	}
	if adapter == nil {
		panic("ui: cannot register a nil 3D adapter")
	}
	r.mu.Lock()
	r.plugins[scene3DPluginKey{name: name, version: version}] = adapter
	r.mu.Unlock()
}

// SupportsPlugin reports whether the exact unversioned (version 0) or positive
// versioned adapter is registered. It does not invoke plugin code.
func (r *Scene3DPluginRegistry) SupportsPlugin(name string, version int) bool {
	if name == "" || version < 0 {
		return false
	}
	r.mu.RLock()
	adapter := r.plugins[scene3DPluginKey{name: name, version: version}]
	r.mu.RUnlock()
	return adapter != nil
}

// Build converts every declaration using the same adapter path and lifecycle.
func (r *Scene3DPluginRegistry) Build(scene *NativeScene3D) (*Scene3DRuntime, error) {
	r.mu.RLock()
	fallback := r.adapter
	plugins := make(map[scene3DPluginKey]Scene3DAdapter, len(r.plugins))
	for key, adapter := range r.plugins {
		plugins[key] = adapter
	}
	r.mu.RUnlock()

	var components []Game3DComponent
	closeComponents := func() {
		for i := len(components) - 1; i >= 0; i-- {
			components[i].Close()
		}
	}
	if scene == nil {
		return NewScene3DRuntime(nil), nil
	}
	for index, node := range scene.Children {
		if node == nil {
			continue
		}
		adapter := fallback
		if node.Kind == Game3DPlugin {
			if node.PluginVersion < 0 {
				closeComponents()
				return nil, fmt.Errorf("ui: 3D plugin %q has invalid version %d at index %d", node.PluginName, node.PluginVersion, index)
			}
			adapter = plugins[scene3DPluginKey{name: node.PluginName, version: node.PluginVersion}]
			if adapter == nil {
				closeComponents()
				if node.PluginVersion == 0 {
					return nil, fmt.Errorf("ui: 3D plugin %q at index %d is not registered", node.PluginName, index)
				}
				return nil, fmt.Errorf("ui: 3D plugin %q version %d at index %d is not registered", node.PluginName, node.PluginVersion, index)
			}
			if validator, ok := adapter.(Scene3DValidator); ok {
				if err := validator.Validate3D(node); err != nil {
					closeComponents()
					return nil, fmt.Errorf("ui: validate 3D plugin %q version %d at index %d: %w", node.PluginName, node.PluginVersion, index, err)
				}
			}
		}
		if adapter == nil {
			closeComponents()
			return nil, fmt.Errorf("ui: no 3D adapter registered for node kind %d at index %d", node.Kind, index)
		}
		component, err := adapter.Adapt3D(node)
		if err != nil {
			closeComponents()
			return nil, fmt.Errorf("ui: adapt 3D node at index %d: %w", index, err)
		}
		if component == nil {
			closeComponents()
			return nil, fmt.Errorf("ui: 3D adapter returned nil for node at index %d", index)
		}
		components = append(components, component)
	}
	return NewScene3DRuntime(components), nil
}

// MustBuild is Build with panic-on-error semantics for application startup.
func (r *Scene3DPluginRegistry) MustBuild(scene *NativeScene3D) *Scene3DRuntime {
	runtime, err := r.Build(scene)
	if err != nil {
		panic(err)
	}
	return runtime
}

// Scene3DRuntime advances and closes a mounted component list. Unlike
// BackendRuntime, it retains its component list after Close and can be reactivated.
type Scene3DRuntime struct {
	components []Game3DComponent
	closed     bool
}

// NewScene3DRuntime wraps native components without taking ownership of their
// construction. Closing the runtime releases them in reverse declaration order.
func NewScene3DRuntime(components []Game3DComponent) *Scene3DRuntime {
	return &Scene3DRuntime{components: components}
}

// Components returns a copy of the component list for native scene mounting.
func (s *Scene3DRuntime) Components() []Game3DComponent {
	if s == nil {
		return nil
	}
	return append([]Game3DComponent(nil), s.components...)
}

// Activate allows updates after Close so the backend can remount the retained
// components. It does not perform native mounting itself.
func (s *Scene3DRuntime) Activate() {
	if s != nil {
		s.closed = false
	}
}

// Update advances every component in declaration order.
func (s *Scene3DRuntime) Update(dt time.Duration) {
	if s == nil || s.closed {
		return
	}
	for _, component := range s.components {
		component.Update(dt)
	}
}

// Close releases components in reverse order and is safe to call repeatedly.
// Call Activate before mounting the same runtime again.
func (s *Scene3DRuntime) Close() {
	if s == nil || s.closed {
		return
	}
	s.closed = true
	for i := len(s.components) - 1; i >= 0; i-- {
		s.components[i].Close()
	}
}
