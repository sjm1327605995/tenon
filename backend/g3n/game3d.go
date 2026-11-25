package g3n

import (
	"fmt"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

func vec3Value(source value.Value[ui.Vec3]) value.Value[Vec3] {
	if source == nil {
		return nil
	}
	return value.Map(source, func(v ui.Vec3) Vec3 { return V3(v.X, v.Y, v.Z) })
}

func color3D(c ui.Color) Color {
	return RGB(float32(c.R)/255, float32(c.G)/255, float32(c.B)/255)
}

func color3DValue(source value.Value[ui.Color]) value.Value[Color] {
	if source == nil {
		return nil
	}
	return value.Map(source, color3D)
}

func transform3D(props ui.Transform3DProps) TransformProps {
	return TransformProps{
		Position: vec3Value(props.Position), Rotation: vec3Value(props.Rotation),
		Scale: vec3Value(props.Scale), Spin: vec3Value(props.Spin),
	}
}

func shape3D(shape ui.Shape3D) Shape {
	switch shape.Kind {
	case ui.Shape3DCube:
		return CubeShape(shape.Size)
	case ui.Shape3DSphere:
		return SphereShape(shape.Radius, shape.WidthSegments, shape.HeightSegments)
	case ui.Shape3DCustom:
		if shape.Geometry == nil {
			return Shape{}
		}
		native := geometry.NewGeometry()
		if len(shape.Geometry.Positions) > 0 {
			native.AddVBO(gls.NewVBO(math32.ArrayF32(shape.Geometry.Positions)).AddAttrib(gls.VertexPosition))
		}
		if len(shape.Geometry.Normals) > 0 {
			native.AddVBO(gls.NewVBO(math32.ArrayF32(shape.Geometry.Normals)).AddAttrib(gls.VertexNormal))
		}
		if len(shape.Geometry.UVs) > 0 {
			native.AddVBO(gls.NewVBO(math32.ArrayF32(shape.Geometry.UVs)).AddAttrib(gls.VertexTexcoord))
		}
		return Shape{native: native}
	default:
		return Shape{}
	}
}

func material3D(props ui.Material3DProps) MaterialProps {
	return MaterialProps{
		Color: color3DValue(props.Color), Opacity: props.Opacity,
		Image: props.Texture, ImageKey: props.TextureKey,
		Unlit: props.Unlit, CacheKey: props.CacheKey,
	}
}

// Mesh3D adapts a backend-neutral mesh declaration to a native component.
func Mesh3D(node *ui.GameNode3D) *MeshComponent {
	if node == nil || node.Kind != ui.Game3DMesh || node.Mesh == nil {
		panic("g3n: Mesh3D requires a mesh declaration")
	}
	return Mesh(MeshProps{Shape: shape3D(node.Mesh.Shape), Material: material3D(node.Mesh.Material), Transform: transform3D(node.Mesh.Transform)})
}

// Component3D adapts one Tenon scene declaration.
func Component3D(node *ui.GameNode3D) Component {
	if node == nil {
		return nil
	}
	switch node.Kind {
	case ui.Game3DMesh:
		return Mesh3D(node)
	case ui.Game3DCamera:
		props := node.Camera
		return PerspectiveCamera(CameraProps{FOV: props.FOV, Near: props.Near, Far: props.Far, Transform: transform3D(props.Transform)})
	case ui.Game3DAmbient, ui.Game3DPoint, ui.Game3DDirectional:
		props := node.Light
		lightProps := LightProps{Color: color3DValue(props.Color), Intensity: props.Intensity, Transform: transform3D(props.Transform)}
		switch node.Kind {
		case ui.Game3DAmbient:
			return AmbientLight(lightProps)
		case ui.Game3DPoint:
			return PointLight(lightProps)
		default:
			return DirectionalLight(lightProps)
		}
	case ui.Game3DNative:
		component, ok := node.Native.(Component)
		if !ok || component == nil {
			panic("g3n: Native3D does not contain a G3N component")
		}
		return component
	default:
		panic(fmt.Sprintf("g3n: unsupported 3D node kind %d", node.Kind))
	}
}

// Scene3DAdapter adapts a portable declaration to a native component.
type Scene3DAdapter = ui.Scene3DAdapter

// Scene3DAdapterFunc adapts a function to a Scene3DAdapter.
type Scene3DAdapterFunc = ui.Scene3DAdapterFunc

// NewScene3DPluginRegistry creates a registry with G3N's built-in adapter.
// Third-party nodes can be added without changing Tenon or this backend.
func NewScene3DPluginRegistry() *ui.Scene3DPluginRegistry {
	registry := ui.NewScene3DPluginRegistry()
	registry.RegisterAdapter(ui.Scene3DAdapterFunc(func(node *ui.GameNode3D) (ui.Game3DComponent, error) {
		return Component3D(node), nil
	}))
	return registry
}

// Scene3D builds a native scene entirely from Tenon declarations.
func Scene3D(scene *ui.NativeScene3D) *Scene {
	return GameSceneRuntime(NewScene3DPluginRegistry().MustBuild(scene))
}
