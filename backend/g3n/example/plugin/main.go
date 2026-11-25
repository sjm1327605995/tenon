package main

import (
	"fmt"

	"github.com/sjm1327605995/tenon/backend/g3n"
	"github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

type spinningSphereAdapter struct{}

func (spinningSphereAdapter) Validate3D(node *ui.GameNode3D) error {
	props, ok := node.Native.(ui.Mesh3DProps)
	if !ok {
		return fmt.Errorf("spinning-sphere: payload must be ui.Mesh3DProps")
	}
	if props.Shape.Kind != ui.Shape3DSphere || props.Shape.Radius <= 0 {
		return fmt.Errorf("spinning-sphere: payload must contain a positive sphere radius")
	}
	return nil
}

func (spinningSphereAdapter) Adapt3D(node *ui.GameNode3D) (ui.Game3DComponent, error) {
	return g3n.Mesh3D(ui.Mesh3D(node.Native.(ui.Mesh3DProps))), nil
}

func main() {
	spin := value.New(ui.V3(0, 1, 0))
	registry := g3n.NewScene3DPluginRegistry()
	registry.RegisterPluginVersion("spinning-sphere", 1, spinningSphereAdapter{})
	if !registry.SupportsPlugin("spinning-sphere", 1) {
		panic("spinning-sphere v1 is unavailable")
	}
	scene := g3n.GameSceneRuntime(registry.MustBuild(ui.GameScene3D(
		ui.PerspectiveCamera3D(ui.Camera3DProps{
			FOV:       60,
			Transform: ui.Transform3DProps{Position: value.Static(ui.V3(0, 0, 3))},
		}),
		ui.Plugin3DVersion("spinning-sphere", 1, ui.Mesh3DProps{
			Shape:    ui.Sphere3D(1, 32, 32),
			Material: ui.Material3DProps{Color: value.Static(ui.Color{R: 41, G: 158, B: 235, A: 255})},
			Transform: ui.Transform3DProps{
				Spin: spin,
			},
		}),
		ui.AmbientLight3D(ui.Light3DProps{Intensity: value.Static(float32(0.8))}),
		ui.PointLight3D(ui.Light3DProps{
			Intensity: value.Static(float32(5)),
			Transform: ui.Transform3DProps{Position: value.Static(ui.V3(1, 0, 2))},
		}),
	)))
	g3n.Use(g3n.Options{Launch: scene.Launch, Update: scene.Frame})
	ui.WindowTitle("Tenon declarative UI + third-party 3D plugin")
	ui.WindowSize(960, 640)
	ui.Run(ui.Use(func(_ struct{}) *ui.Node {
		active, setActive := ui.UseState(true)
		label := "暂停旋转"
		if !active {
			label = "继续旋转"
		}
		return ui.VStack(12,
			ui.Text("Tenon + third-party 3D plugin").FontSize(24).TextColor(ui.White),
			ui.Text("统一场景格式与插件生命周期").TextColor(ui.White),
			ui.Button(
				ui.Style(ui.Width(120), ui.Height(36), ui.Bg(ui.Blue)),
				ui.OnClick(func() {
					if active {
						spin.Set(ui.V3(0, 0, 0))
					} else {
						spin.Set(ui.V3(0, 1, 0))
					}
					setActive(!active)
				}),
				ui.Text(label).TextColor(ui.White),
			),
		).Width(520).Padding(20).Bg(ui.Color{R: 15, G: 23, B: 42, A: 220}).Radius(12)
	}, struct{}{}))
}
