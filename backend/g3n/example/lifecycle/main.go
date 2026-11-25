package main

import (
	"fmt"
	"time"

	"github.com/sjm1327605995/tenon/backend/g3n"
	"github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

type retryMesh struct {
	delegate *g3n.MeshComponent
	mounts   int
}

func (c *retryMesh) Mount(ctx *g3n.SceneContext) error {
	c.mounts++
	if c.mounts == 1 {
		return fmt.Errorf("lifecycle: intentional first mount failure")
	}
	return c.delegate.Mount(ctx)
}

func (c *retryMesh) Update(dt time.Duration) { c.delegate.Update(dt) }
func (c *retryMesh) Close()                  { c.delegate.Close() }

type sphereAdapter struct {
	component *retryMesh
}

func (*sphereAdapter) Validate3D(node *ui.GameNode3D) error {
	props, ok := node.Native.(ui.Mesh3DProps)
	if !ok || props.Shape.Kind != ui.Shape3DSphere || props.Shape.Radius <= 0 {
		return fmt.Errorf("lifecycle: expected a sphere Mesh3DProps payload")
	}
	return nil
}

func (a *sphereAdapter) Adapt3D(node *ui.GameNode3D) (ui.Game3DComponent, error) {
	a.component = &retryMesh{delegate: g3n.Mesh3D(ui.Mesh3D(node.Native.(ui.Mesh3DProps)))}
	return a.component, nil
}

func main() {
	registry := g3n.NewScene3DPluginRegistry()
	adapter := &sphereAdapter{}
	registry.RegisterPluginVersion("lifecycle-sphere", 1, adapter)
	scene := g3n.GameSceneRuntime(registry.MustBuild(ui.GameScene3D(
		ui.PerspectiveCamera3D(ui.Camera3DProps{
			FOV:       60,
			Transform: ui.Transform3DProps{Position: value.Static(ui.V3(0, 0, 3))},
		}),
		ui.Plugin3DVersion("lifecycle-sphere", 1, ui.Mesh3DProps{
			Shape:    ui.Sphere3D(1, 32, 32),
			Material: ui.Material3DProps{Color: value.Static(ui.Color{R: 41, G: 158, B: 235, A: 255})},
		}),
		ui.AmbientLight3D(ui.Light3DProps{Intensity: value.Static(float32(0.8))}),
	)))
	frame := 0
	g3n.Use(g3n.Options{
		Launch: func(game *g3n.Game) {
			if err := scene.Mount(game); err == nil {
				panic("lifecycle: first mount unexpectedly succeeded")
			}
		},
		Update: func(game *g3n.Game, dt time.Duration) {
			frame++
			switch frame {
			case 2, 6:
				scene.Close()
			case 4, 8:
				if err := scene.Mount(game); err != nil {
					panic(err)
				}
			}
			scene.Update(dt)
		},
	})
	ui.WindowTitle("Tenon — native G3N lifecycle")
	ui.WindowSize(640, 360)
	ui.Run(ui.Text("Native G3N Close/Mount lifecycle"))
	attempts := 0
	if adapter.component != nil {
		attempts = adapter.component.mounts
	}
	if attempts != 3 {
		panic(fmt.Sprintf("lifecycle: mount attempts=%d, want 3", attempts))
	}
}
