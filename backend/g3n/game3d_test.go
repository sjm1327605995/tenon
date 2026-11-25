package g3n

import (
	"image"
	"testing"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

func TestPortableScene3DMountsAndUpdates(t *testing.T) {
	position := value.New(ui.V3(0, 0, 0))
	spin := value.New(ui.V3(0, 1, 0))
	color := value.New(ui.Hex("#ff6600"))
	imageValue := image.NewRGBA(image.Rect(0, 0, 1, 1))
	scene := Scene3D(ui.GameScene3D(
		ui.Mesh3D(ui.Mesh3DProps{
			Shape: ui.Sphere3D(1, 12, 8),
			Material: ui.Material3DProps{
				Color: color, Texture: imageValue, TextureKey: "portable-texture", Unlit: true,
			},
			Transform: ui.Transform3DProps{Position: position, Spin: spin},
		}),
		ui.PerspectiveCamera3D(ui.Camera3DProps{FOV: 60}),
	))
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	defer scene.Close()
	mesh := scene.Components()[0].(*MeshComponent)
	position.Set(ui.V3(2, 3, 4))
	if got := mesh.Native.Position(); got.X != 2 || got.Y != 3 || got.Z != 4 {
		t.Fatalf("portable position = %+v", got)
	}
	if len(scene.context.textures) != 1 || mesh.Native.GetGraphic().Materials()[0].IMaterial().GetMaterial().TextureCount() != 1 {
		t.Fatal("portable image was not converted and cached as a texture")
	}
	if got := game.Camera.Fov(); got != 60 {
		t.Fatalf("portable camera FOV = %v", got)
	}
}

func TestPortableCustomGeometry(t *testing.T) {
	node := ui.Mesh3D(ui.Mesh3DProps{Shape: ui.Custom3D(&ui.Geometry3D{
		Positions: []float32{0, 0, 0, 1, 0, 0, 0, 1, 0},
		Normals:   []float32{0, 0, 1, 0, 0, 1, 0, 0, 1},
		UVs:       []float32{0, 0, 1, 0, 0, 1},
	})})
	mesh := Mesh3D(node)
	scene := GameScene(mesh)
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	defer scene.Close()
	geometry := mesh.Native.GetGeometry()
	if geometry.VBO(gls.VertexPosition) == nil || geometry.VBO(gls.VertexNormal) == nil || geometry.VBO(gls.VertexTexcoord) == nil {
		t.Fatal("portable geometry did not preserve its vertex buffers")
	}
}

func TestNativeSceneRemountRecreatesLightAndKeepsSharedCamera(t *testing.T) {
	sharedCamera := camera.New(1)
	root := core.NewNode()
	root.Add(sharedCamera)
	scene := GameScene(
		PerspectiveCamera(CameraProps{FOV: 60}),
		AmbientLight(LightProps{}),
	)
	game := &Game{Scene: root, Camera: sharedCamera}
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	defer scene.Close()
	light := scene.Components()[1].(*LightComponent)
	firstLight := light.Native
	if firstLight == nil {
		t.Fatal("light did not mount")
	}

	scene.Close()
	if sharedCamera.Parent() == nil {
		t.Fatal("closing a declarative scene detached the shared game camera")
	}
	if light.Native != nil {
		t.Fatal("closed light retained its native object")
	}
	if err := scene.Mount(game); err != nil {
		t.Fatalf("native scene remount failed: %v", err)
	}
	if light.Native == nil || light.Native == firstLight {
		t.Fatal("remount did not create a fresh light object")
	}
	scene.Close()
	if sharedCamera.Parent() == nil {
		t.Fatal("second close detached the shared game camera")
	}
}
