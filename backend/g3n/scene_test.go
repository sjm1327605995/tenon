package g3n

import (
	"image"
	"testing"
	"time"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/texture"
	"github.com/sjm1327605995/tenon/pkg/value"
)

func TestSceneSharesResourcesAndUpdatesNativeObjects(t *testing.T) {
	position := value.New(V3(0, 0, 0))
	spin := value.New(V3(0, 1, 0))
	color := value.New(RGB(0.2, 0.4, 0.6))
	material := StandardMaterial(color)
	first := Sphere(SphereProps{Radius: 1, Material: material, Transform: TransformProps{Position: position, Spin: spin}})
	second := Sphere(SphereProps{Radius: 1, Material: material})
	cam := camera.New(1)
	game := &Game{Scene: core.NewNode(), Camera: cam}
	scene := GameScene(first, second, PerspectiveCamera(CameraProps{FOV: 75}))
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	if first.Native.GetGeometry() != second.Native.GetGeometry() {
		t.Fatal("identical spheres did not share geometry")
	}
	if len(scene.context.materials) != 1 {
		t.Fatalf("material cache has %d entries, want 1", len(scene.context.materials))
	}
	native := first.Native
	position.Set(V3(1, 2, 3))
	if got := native.Position(); got.X != 1 || got.Y != 2 || got.Z != 3 {
		t.Fatalf("position = %+v, want (1,2,3)", got)
	}
	scene.Update(time.Second)
	if first.Native != native {
		t.Fatal("scene update recreated the mesh")
	}
	if got := native.Rotation().Y; got != 1 {
		t.Fatalf("rotation Y = %v, want 1", got)
	}
	if cam.Fov() != 75 {
		t.Fatalf("camera FOV = %v, want 75", cam.Fov())
	}
	scene.Close()
	if scene.Root() != nil || len(game.Scene.Children()) != 0 {
		t.Fatal("scene close did not detach only its own root")
	}
	position.Set(V3(9, 9, 9))
	if got := native.Position(); got.X != 1 || got.Y != 2 || got.Z != 3 {
		t.Fatalf("closed subscription still updated position: %+v", got)
	}
}

func TestSceneClosePreservesNativeObjects(t *testing.T) {
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	native := core.NewNode()
	game.Scene.Add(native)
	scene := GameScene(Sphere(SphereProps{}))
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	scene.Close()
	if len(game.Scene.Children()) != 1 || game.Scene.ChildAt(0) != native {
		t.Fatal("scene close removed a native object it does not own")
	}
}

func TestTexturedGeometryUsesSeparateVAOCacheEntry(t *testing.T) {
	tex := texture.NewTexture2DFromRGBA(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	plain := Cube(CubeProps{Size: 1})
	first := Cube(CubeProps{Size: 1, Material: MaterialProps{Texture: tex, Unlit: true}})
	second := Cube(CubeProps{Size: 1, Material: MaterialProps{Texture: tex, Unlit: true}})
	scene := GameScene(plain, first, second)
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	defer scene.Close()
	if plain.Native.GetGeometry() == first.Native.GetGeometry() {
		t.Fatal("textured and untextured materials share a VAO cache entry")
	}
	if first.Native.GetGeometry() != second.Native.GetGeometry() {
		t.Fatal("compatible textured meshes did not share geometry")
	}
}

func TestSharedMeshCloseKeepsCachedResourcesAlive(t *testing.T) {
	material := StandardMaterial(value.Static(RGB(1, 1, 1)))
	first := Cube(CubeProps{Size: 1, Material: material})
	second := Cube(CubeProps{Size: 1, Material: material})
	scene := GameScene(first, second)
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	first.Close()
	if second.Native.GetGeometry().VBO(gls.VertexPosition) == nil || len(second.Native.GetGraphic().Materials()) != 1 {
		t.Fatal("closing one shared mesh destroyed its sibling's resources")
	}
	scene.Close()
}

func TestGeometryShapeMountsApplicationOwnedGeometry(t *testing.T) {
	native := geometry.NewCube(1)
	mesh := Mesh(MeshProps{Shape: GeometryShape(native)})
	scene := GameScene(mesh)
	game := &Game{Scene: core.NewNode(), Camera: camera.New(1)}
	if err := scene.Mount(game); err != nil {
		t.Fatal(err)
	}
	if mesh.Native.GetGeometry() != native {
		t.Fatal("native geometry was replaced during mount")
	}
	scene.Close()
	if native.VBO(gls.VertexPosition) == nil {
		t.Fatal("scene close disposed application-owned geometry")
	}
}
