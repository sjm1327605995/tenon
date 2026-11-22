package main

import (
	"github.com/sjm1327605995/tenon/backend/g3n"
	"github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

func main() {
	spin := value.New(ui.V3(0, 1, 0))
	scene := g3n.Scene3D(ui.GameScene3D(
		ui.PerspectiveCamera3D(ui.Camera3DProps{
			FOV:       60,
			Transform: ui.Transform3DProps{Position: value.Static(ui.V3(0, 0, 3))},
		}),
		ui.Mesh3D(ui.Mesh3DProps{
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
	))
	g3n.Use(g3n.Options{Launch: scene.Launch, Update: scene.Frame})
	ui.WindowTitle("Tenon declarative UI + G3N components")
	ui.WindowSize(960, 640)
	ui.Run(ui.Use(func(_ struct{}) *ui.Node {
		active, setActive := ui.UseState(true)
		label := "暂停旋转"
		if !active {
			label = "继续旋转"
		}
		return ui.VStack(12,
			ui.Text("Tenon + G3N").FontSize(24).TextColor(ui.White),
			ui.Text("声明式 3D 组件 + 声明式界面").TextColor(ui.White),
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
