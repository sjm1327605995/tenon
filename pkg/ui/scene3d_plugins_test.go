package ui

import (
	"fmt"
	"testing"
	"time"
)

type test3DComponent struct {
	name  string
	log   *[]string
	close int
}

func (c *test3DComponent) Update(time.Duration) {
	*c.log = append(*c.log, "update:"+c.name)
}

func (c *test3DComponent) Close() {
	c.close++
	*c.log = append(*c.log, "close:"+c.name)
}

func test3DAdapter(log *[]string) Scene3DAdapterFunc {
	return func(node *GameNode3D) (Game3DComponent, error) {
		if node.Kind == Game3DPlugin {
			return &test3DComponent{name: node.PluginName, log: log}, nil
		}
		return &test3DComponent{name: "builtin", log: log}, nil
	}
}

func TestScene3DPluginLifecycle(t *testing.T) {
	var log []string
	registry := NewScene3DPluginRegistry()
	registry.RegisterAdapter(test3DAdapter(&log))
	registry.RegisterPlugin("particles", test3DAdapter(&log))
	runtime := registry.MustBuild(GameScene3D(
		Mesh3D(Mesh3DProps{}),
		Plugin3D("particles", struct{ Count int }{Count: 100}),
	))
	if got := len(runtime.Components()); got != 2 {
		t.Fatalf("component count=%d, want 2", got)
	}

	runtime.Update(time.Second)
	runtime.Close()
	runtime.Close()
	want := []string{"update:builtin", "update:particles", "close:particles", "close:builtin"}
	if len(log) != len(want) {
		t.Fatalf("lifecycle log=%v, want %v", log, want)
	}
	for i := range want {
		if log[i] != want[i] {
			t.Fatalf("lifecycle log=%v, want %v", log, want)
		}
	}
}

func TestUnknownPluginClosesEarlierComponents(t *testing.T) {
	var log []string
	registry := NewScene3DPluginRegistry()
	registry.RegisterAdapter(test3DAdapter(&log))
	_, err := registry.Build(GameScene3D(
		Mesh3D(Mesh3DProps{}),
		Plugin3D("missing", nil),
	))
	if err == nil {
		t.Fatal("unknown plugin did not return an error")
	}
	if len(log) != 1 || log[0] != "close:builtin" {
		t.Fatalf("cleanup log=%v, want close:builtin", log)
	}
}

func TestPluginRegistriesAreIndependent(t *testing.T) {
	var log []string
	first := NewScene3DPluginRegistry()
	first.RegisterAdapter(test3DAdapter(&log))
	first.RegisterPlugin("custom", test3DAdapter(&log))
	second := NewScene3DPluginRegistry()
	second.RegisterAdapter(test3DAdapter(&log))
	if _, err := second.Build(GameScene3D(Plugin3D("custom", nil))); err == nil {
		t.Fatal("plugin leaked into an independent registry")
	}
	if _, err := first.Build(GameScene3D(Plugin3D("custom", nil))); err != nil {
		t.Fatalf("plugin was not available in its own registry: %v", err)
	}
}

func TestEmptyPluginScene(t *testing.T) {
	runtime := NewScene3DPluginRegistry().MustBuild(nil)
	if len(runtime.Components()) != 0 {
		t.Fatal("nil scene produced components")
	}
}

func TestScene3DRuntimeCanReactivate(t *testing.T) {
	var log []string
	component := &test3DComponent{name: "component", log: &log}
	runtime := NewScene3DRuntime([]Game3DComponent{component})

	runtime.Update(time.Second)
	runtime.Close()
	runtime.Update(time.Second)
	runtime.Activate()
	runtime.Update(time.Second)
	runtime.Close()

	want := []string{"update:component", "close:component", "update:component", "close:component"}
	if len(log) != len(want) {
		t.Fatalf("reactivation log=%v, want %v", log, want)
	}
	for i := range want {
		if log[i] != want[i] {
			t.Fatalf("reactivation log=%v, want %v", log, want)
		}
	}
	if len(runtime.Components()) != 1 {
		t.Fatal("Close discarded components needed for remounting")
	}
}

func TestAdapterFailureAndNilComponentCleanup(t *testing.T) {
	for _, test := range []struct {
		name    string
		adapter Scene3DAdapterFunc
	}{
		{
			name: "adapter error",
			adapter: func(*GameNode3D) (Game3DComponent, error) {
				return nil, fmt.Errorf("adapter failed")
			},
		},
		{
			name: "nil component",
			adapter: func(*GameNode3D) (Game3DComponent, error) {
				return nil, nil
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var log []string
			registry := NewScene3DPluginRegistry()
			registry.RegisterAdapter(test3DAdapter(&log))
			registry.RegisterPlugin("bad", test.adapter)
			_, err := registry.Build(GameScene3D(
				Mesh3D(Mesh3DProps{}),
				Plugin3D("bad", nil),
			))
			if err == nil {
				t.Fatal("adapter failure did not return an error")
			}
			if len(log) != 1 || log[0] != "close:builtin" {
				t.Fatalf("cleanup log=%v, want close:builtin", log)
			}
		})
	}
}

func TestVersionedPluginsResolveExactly(t *testing.T) {
	var log []string
	registry := NewScene3DPluginRegistry()
	for _, version := range []int{1, 2} {
		version := version
		registry.RegisterPluginVersion("particles", version, Scene3DAdapterFunc(func(*GameNode3D) (Game3DComponent, error) {
			return &test3DComponent{name: fmt.Sprintf("v%d", version), log: &log}, nil
		}))
	}
	if !registry.SupportsPlugin("particles", 2) || registry.SupportsPlugin("particles", 3) || registry.SupportsPlugin("particles", 0) {
		t.Fatal("plugin capability query does not match registered versions")
	}
	runtime := registry.MustBuild(GameScene3D(Plugin3DVersion("particles", 2, nil)))
	if got := len(runtime.Components()); got != 1 {
		t.Fatalf("component count=%d, want 1", got)
	}
	runtime.Update(0)
	if len(log) != 1 || log[0] != "update:v2" {
		t.Fatalf("versioned update log=%v, want update:v2", log)
	}
	if _, err := registry.Build(GameScene3D(Plugin3DVersion("particles", 3, nil))); err == nil {
		t.Fatal("unregistered plugin version was accepted")
	}
	if _, err := registry.Build(GameScene3D(Plugin3D("particles", nil))); err == nil {
		t.Fatal("versioned adapters received an unversioned declaration")
	}
}

type test3DValidatedAdapter struct {
	validations, adaptations int
	log                      *[]string
}

func (a *test3DValidatedAdapter) Validate3D(node *GameNode3D) error {
	a.validations++
	payload, ok := node.Native.(struct{ Count int })
	if !ok || payload.Count <= 0 {
		return fmt.Errorf("payload must contain a positive Count")
	}
	return nil
}

func (a *test3DValidatedAdapter) Adapt3D(*GameNode3D) (Game3DComponent, error) {
	a.adaptations++
	return &test3DComponent{name: "validated", log: a.log}, nil
}

func TestPluginValidationRunsBeforeAdaptation(t *testing.T) {
	var log []string
	adapter := &test3DValidatedAdapter{log: &log}
	registry := NewScene3DPluginRegistry()
	registry.RegisterAdapter(test3DAdapter(&log))
	registry.RegisterPluginVersion("particles", 1, adapter)

	_, err := registry.Build(GameScene3D(
		Mesh3D(Mesh3DProps{}),
		Plugin3DVersion("particles", 1, nil),
	))
	if err == nil {
		t.Fatal("invalid plugin payload was accepted")
	}
	if adapter.validations != 1 || adapter.adaptations != 0 {
		t.Fatalf("validation=%d adaptation=%d, want 1/0", adapter.validations, adapter.adaptations)
	}
	if len(log) != 1 || log[0] != "close:builtin" {
		t.Fatalf("cleanup log=%v, want close:builtin", log)
	}

	if _, err := registry.Build(GameScene3D(Plugin3DVersion("particles", 1, struct{ Count int }{Count: 2}))); err != nil {
		t.Fatalf("valid plugin payload was rejected: %v", err)
	}
	if adapter.validations != 2 || adapter.adaptations != 1 {
		t.Fatalf("validation=%d adaptation=%d, want 2/1", adapter.validations, adapter.adaptations)
	}
}
