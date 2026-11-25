package ui

import (
	"image"

	"github.com/sjm1327605995/tenon/pkg/value"
)

// Vec3 is a backend-neutral 3D vector.
type Vec3 struct{ X, Y, Z float32 }

func V3(x, y, z float32) Vec3 { return Vec3{X: x, Y: y, Z: z} }

// Geometry3D contains triangle vertices in backend-neutral arrays.
type Geometry3D struct {
	Positions, Normals, UVs []float32
}

type Shape3DKind uint8

const (
	Shape3DCustom Shape3DKind = iota
	Shape3DCube
	Shape3DSphere
)

// Shape3D describes shared geometry without exposing a native engine object.
type Shape3D struct {
	Kind                          Shape3DKind
	Size                          float32
	Radius                        float64
	WidthSegments, HeightSegments int
	Geometry                      *Geometry3D
}

func Cube3D(size float32) Shape3D { return Shape3D{Kind: Shape3DCube, Size: size} }

func Sphere3D(radius float64, widthSegments, heightSegments int) Shape3D {
	return Shape3D{Kind: Shape3DSphere, Radius: radius, WidthSegments: widthSegments, HeightSegments: heightSegments}
}

func Custom3D(geometry *Geometry3D) Shape3D { return Shape3D{Geometry: geometry} }

// Material3DProps is a portable material declaration. Texture is decoded by the
// active backend and cached under TextureKey within each scene.
type Material3DProps struct {
	Color      value.Value[Color]
	Opacity    value.Value[float32]
	Texture    image.Image
	TextureKey string
	Unlit      bool
	CacheKey   string
}

// Transform3DProps controls position, Euler rotation, scale and spin.
type Transform3DProps struct {
	Position value.Value[Vec3]
	Rotation value.Value[Vec3]
	Scale    value.Value[Vec3]
	Spin     value.Value[Vec3]
}

type Mesh3DProps struct {
	Shape     Shape3D
	Material  Material3DProps
	Transform Transform3DProps
}

type Camera3DProps struct {
	FOV, Near, Far float32
	Transform      Transform3DProps
}

type Light3DProps struct {
	Color     value.Value[Color]
	Intensity value.Value[float32]
	Transform Transform3DProps
}

type Game3DNodeKind uint8

const (
	Game3DMesh Game3DNodeKind = iota
	Game3DCamera
	Game3DAmbient
	Game3DPoint
	Game3DDirectional
	Game3DNative
	Game3DPlugin
)

// GameNode3D is an immutable scene declaration. Native is an advanced escape
// hatch for backend-specific components; ordinary game code uses the Props above.
type GameNode3D struct {
	Kind          Game3DNodeKind
	Mesh          *Mesh3DProps
	Camera        *Camera3DProps
	Light         *Light3DProps
	Native        any
	PluginName    string
	PluginVersion int
}

func Mesh3D(props Mesh3DProps) *GameNode3D {
	return &GameNode3D{Kind: Game3DMesh, Mesh: &props}
}

func PerspectiveCamera3D(props Camera3DProps) *GameNode3D {
	return &GameNode3D{Kind: Game3DCamera, Camera: &props}
}

func AmbientLight3D(props Light3DProps) *GameNode3D {
	return &GameNode3D{Kind: Game3DAmbient, Light: &props}
}

func PointLight3D(props Light3DProps) *GameNode3D {
	return &GameNode3D{Kind: Game3DPoint, Light: &props}
}

func DirectionalLight3D(props Light3DProps) *GameNode3D {
	return &GameNode3D{Kind: Game3DDirectional, Light: &props}
}

func Native3D(component any) *GameNode3D {
	return &GameNode3D{Kind: Game3DNative, Native: component}
}

// Plugin3D declares an unversioned backend extension resolved by name through a
// Scene3DPluginRegistry. The payload is interpreted only by that plugin.
func Plugin3D(name string, payload any) *GameNode3D {
	return &GameNode3D{Kind: Game3DPlugin, PluginName: name, Native: payload}
}

// Plugin3DVersion declares an extension that requires the same positive adapter
// version at build time. Versioned nodes do not fall back to unversioned plugins.
func Plugin3DVersion(name string, version int, payload any) *GameNode3D {
	if version <= 0 {
		panic("ui: 3D plugin version must be positive")
	}
	return &GameNode3D{Kind: Game3DPlugin, PluginName: name, PluginVersion: version, Native: payload}
}

// NativeScene3D groups declarations for a native backend scene.
type NativeScene3D struct{ Children []*GameNode3D }

func GameScene3D(children ...*GameNode3D) *NativeScene3D {
	return &NativeScene3D{Children: append([]*GameNode3D(nil), children...)}
}
