package g3n

import (
	"errors"
	"testing"
	"time"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

type testPluginComponent struct {
	mounts, updates, closes int
}

func (c *testPluginComponent) Mount(*SceneContext) error {
	c.mounts++
	return nil
}

func (c *testPluginComponent) Update(time.Duration) { c.updates++ }
func (c *testPluginComponent) Close()               { c.closes++ }

func TestThirdPartyPluginUsesSceneLifecycle(t *testing.T) {
	component := &testPluginComponent{}
	registry := NewScene3DPluginRegistry()
	registry.RegisterPlugin("test-plugin", Scene3DAdapterFunc(func(node *ui.GameNode3D) (ui.Game3DComponent, error) {
		if node.Native != "payload" {
			t.Fatalf("plugin payload=%v, want payload", node.Native)
		}
		return component, nil
	}))
	scene := GameSceneRuntime(registry.MustBuild(ui.GameScene3D(
		ui.Plugin3D("test-plugin", "payload"),
	)))
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	if component.mounts != 1 {
		t.Fatalf("plugin mounts=%d, want 1", component.mounts)
	}
	scene.Update(time.Second)
	if component.updates != 1 {
		t.Fatalf("plugin updates=%d, want 1", component.updates)
	}
	scene.Close()
	if component.closes != 1 {
		t.Fatalf("plugin closes=%d, want 1", component.closes)
	}
	if len(game.Scene.Children()) != 0 {
		t.Fatal("plugin scene was not detached")
	}
}

func TestSceneCanCloseAndRemount(t *testing.T) {
	component := &testPluginComponent{}
	registry := NewScene3DPluginRegistry()
	registry.RegisterPlugin("remount", Scene3DAdapterFunc(func(*ui.GameNode3D) (ui.Game3DComponent, error) {
		return component, nil
	}))
	scene := GameSceneRuntime(registry.MustBuild(ui.GameScene3D(ui.Plugin3D("remount", nil))))
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	scene.Close()
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	if component.mounts != 2 {
		t.Fatalf("mounts=%d, want 2 after remount", component.mounts)
	}
	scene.Update(time.Second)
	if component.updates != 1 {
		t.Fatalf("updates=%d, want 1 after remount", component.updates)
	}
	scene.Close()
	if component.closes != 2 {
		t.Fatalf("closes=%d, want 2", component.closes)
	}
}

type retryPluginComponent struct {
	mounts int
}

func (c *retryPluginComponent) Mount(*SceneContext) error {
	c.mounts++
	if c.mounts == 1 {
		return errors.New("first mount failed")
	}
	return nil
}
func (c *retryPluginComponent) Update(time.Duration) {}
func (c *retryPluginComponent) Close()               {}

func TestSceneMountFailureCanRetry(t *testing.T) {
	component := &retryPluginComponent{}
	registry := NewScene3DPluginRegistry()
	registry.RegisterPlugin("retry", Scene3DAdapterFunc(func(*ui.GameNode3D) (ui.Game3DComponent, error) {
		return component, nil
	}))
	scene := GameSceneRuntime(registry.MustBuild(ui.GameScene3D(ui.Plugin3D("retry", nil))))
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := scene.Mount(game); err == nil {
		t.Fatal("first mount unexpectedly succeeded")
	}
	if len(game.Scene.Children()) != 0 {
		t.Fatal("failed mount left the scene root attached")
	}
	if err := scene.Mount(game); err != nil {
		t.Fatalf("retry failed: %v", err)
	}
	if component.mounts != 2 {
		t.Fatalf("mount attempts=%d, want 2", component.mounts)
	}
	scene.Close()
}

type updateOnlyPluginComponent struct {
	updates, closes int
}

func (c *updateOnlyPluginComponent) Update(time.Duration) { c.updates++ }
func (c *updateOnlyPluginComponent) Close()               { c.closes++ }

func TestPluginWithoutNativeMountUsesSceneLifecycle(t *testing.T) {
	component := &updateOnlyPluginComponent{}
	registry := NewScene3DPluginRegistry()
	registry.RegisterPlugin("update-only", Scene3DAdapterFunc(func(*ui.GameNode3D) (ui.Game3DComponent, error) {
		return component, nil
	}))
	scene := GameSceneRuntime(registry.MustBuild(ui.GameScene3D(ui.Plugin3D("update-only", nil))))
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := scene.Mount(game); err != nil {
		t.Fatalf("update-only plugin mount failed: %v", err)
	}
	scene.Update(time.Second)
	scene.Close()
	if component.updates != 1 || component.closes != 1 {
		t.Fatalf("update-only lifecycle = %d updates/%d closes, want 1/1", component.updates, component.closes)
	}
}

func TestIndependentScenesCloseSeparately(t *testing.T) {
	firstComponent := &testPluginComponent{}
	secondComponent := &testPluginComponent{}
	first := GameScene(firstComponent)
	second := GameScene(secondComponent)
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := first.Mount(game); err != nil {
		t.Fatal(err)
	}
	if err := second.Mount(game); err != nil {
		t.Fatal(err)
	}
	if first.Root() == nil || second.Root() == nil || first.Root() == second.Root() {
		t.Fatal("independent scenes did not receive distinct native roots")
	}
	if len(game.Scene.Children()) != 2 {
		t.Fatalf("game root children = %d, want 2", len(game.Scene.Children()))
	}

	first.Close()
	if len(game.scenes) != 1 || game.scenes[0] != second {
		t.Fatal("closing the first scene removed the wrong scene registration")
	}
	if len(game.Scene.Children()) != 1 || second.Root().Parent() != game.Scene {
		t.Fatal("closing the first scene detached the second scene root")
	}
	second.Update(time.Second)
	if secondComponent.updates != 1 || firstComponent.updates != 0 {
		t.Fatalf("unexpected scene updates: first=%d second=%d", firstComponent.updates, secondComponent.updates)
	}
	second.Close()
	if len(game.Scene.Children()) != 0 || firstComponent.closes != 1 || secondComponent.closes != 1 {
		t.Fatalf("scene cleanup = %d roots, first closes=%d second closes=%d", len(game.Scene.Children()), firstComponent.closes, secondComponent.closes)
	}
}
