package g3n

import (
	"fmt"
	"image"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/texture"
	"github.com/sjm1327605995/tenon/pkg/value"
)

// MaterialProps describes a shared standard material. Resources with the same
// CacheKey or binding identity are created once per Scene.
const unlitVertexShader = `#include <attributes>
#include <material>
uniform mat4 MVP;
out vec2 FragTexcoord;
void main() {
    vec2 texcoord = VertexTexcoord;
#if MAT_TEXTURES > 0
    if (MatTexFlipY(0)) {
        texcoord.y = 1.0 - texcoord.y;
    }
#endif
    FragTexcoord = texcoord;
    gl_Position = MVP * vec4(VertexPosition, 1.0);
}`

const unlitFragmentShader = `precision highp float;
#include <material>
in vec2 FragTexcoord;
out vec4 FragColor;
void main() {
    vec4 texColor = vec4(1.0);
#if MAT_TEXTURES > 0
    texColor = texture(MatTexture[0], FragTexcoord * MatTexRepeat(0) + MatTexOffset(0));
#endif
    FragColor = vec4(MatDiffuseColor, MatOpacity) * texColor;
}`

type MaterialProps struct {
	Color    value.Value[Color]
	Opacity  value.Value[float32]
	Texture  *texture.Texture2D
	Image    image.Image
	ImageKey string
	Unlit    bool
	CacheKey string
}

// StandardMaterial creates material properties for MeshProps.
func StandardMaterial(color value.Value[Color]) MaterialProps {
	return MaterialProps{Color: color}
}

func (c *SceneContext) materialFor(props MaterialProps) *material.Standard {
	key := props.CacheKey
	if key == "" {
		key = fmt.Sprintf("%s/%s/%p/%t", bindingKey(props.Color), bindingKey(props.Opacity), props.Texture, props.Unlit)
	}
	if cached := c.materials[key]; cached != nil {
		return cached
	}
	created := material.NewStandard(&Color{R: 1, G: 1, B: 1})
	if props.Texture != nil {
		created.AddTexture(props.Texture.Incref())
	}
	if props.Unlit {
		created.SetShader("tenon_unlit")
		created.SetUseLights(material.UseLightNone)
	}
	c.materials[key] = created
	return created
}

type shapeKind uint8

const (
	shapeSphere shapeKind = iota
	shapeCube
)

// Shape identifies geometry shared by mesh components.
type Shape struct {
	kind           shapeKind
	native         *geometry.Geometry
	borrowed       bool
	radius         float64
	size           float32
	widthSegments  int
	heightSegments int
}

// GeometryShape wraps application-owned native geometry without requiring a Tenon primitive.
func GeometryShape(geometry *geometry.Geometry) Shape {
	return Shape{native: geometry, borrowed: true}
}

// SphereShape describes a UV sphere.
func SphereShape(radius float64, widthSegments, heightSegments int) Shape {
	if radius <= 0 {
		radius = 1
	}
	if widthSegments <= 0 {
		widthSegments = 32
	}
	if heightSegments <= 0 {
		heightSegments = 32
	}
	return Shape{kind: shapeSphere, radius: radius, widthSegments: widthSegments, heightSegments: heightSegments}
}

// CubeShape describes a cube with the specified edge length.
func CubeShape(size float32) Shape {
	if size <= 0 {
		size = 1
	}
	return Shape{kind: shapeCube, size: size}
}

func (s Shape) key() string {
	if s.native != nil {
		return fmt.Sprintf("native/%p", s.native)
	}
	return fmt.Sprintf("%d/%g/%g/%d/%d", s.kind, s.radius, s.size, s.widthSegments, s.heightSegments)
}

func (c *SceneContext) geometryFor(shape Shape, textured ...bool) *geometry.Geometry {
	key := shape.key()
	if len(textured) > 0 && textured[0] {
		key += "/textured"
	}
	if cached := c.geometries[key]; cached != nil {
		return cached
	}
	var created *geometry.Geometry
	if shape.native != nil {
		created = shape.native
	} else {
		switch shape.kind {
		case shapeSphere:
			created = geometry.NewSphere(shape.radius, shape.widthSegments, shape.heightSegments)
		case shapeCube:
			created = geometry.NewCube(shape.size)
		}
	}
	if !shape.borrowed {
		c.ownedGeometries[created] = true
	}
	c.geometries[key] = created
	return created
}

// MeshProps describes a native mesh with shared geometry and material.
type MeshProps struct {
	Shape     Shape
	Material  MaterialProps
	Transform TransformProps
}

// MeshComponent is a declarative native mesh.
type MeshComponent struct {
	nodeComponent
	props  MeshProps
	Native *graphic.Mesh
}

// Mesh creates a mesh component.
func Mesh(props MeshProps) *MeshComponent { return &MeshComponent{props: props} }

// SphereProps describes a sphere component.
type SphereProps struct {
	Radius         float64
	WidthSegments  int
	HeightSegments int
	Material       MaterialProps
	Transform      TransformProps
}

// Sphere creates a sphere mesh component.
func Sphere(props SphereProps) *MeshComponent {
	return Mesh(MeshProps{
		Shape:     SphereShape(props.Radius, props.WidthSegments, props.HeightSegments),
		Material:  props.Material,
		Transform: props.Transform,
	})
}

// CubeProps describes a cube component.
type CubeProps struct {
	Size      float32
	Material  MaterialProps
	Transform TransformProps
}

// Cube creates a cube mesh component.
func Cube(props CubeProps) *MeshComponent {
	return Mesh(MeshProps{Shape: CubeShape(props.Size), Material: props.Material, Transform: props.Transform})
}

func (c *SceneContext) textureFor(img image.Image, key string) *texture.Texture2D {
	if img == nil {
		return nil
	}
	if key == "" {
		return TextureFromImage(img)
	}
	if cached := c.textures[key]; cached != nil {
		return cached
	}
	created := TextureFromImage(img)
	c.textures[key] = created
	return created
}

func (c *MeshComponent) Mount(ctx *SceneContext) error {
	if c.Native != nil {
		return fmt.Errorf("g3n: mesh is already mounted")
	}
	materialProps := c.props.Material
	if materialProps.Texture == nil && materialProps.Image != nil {
		key := materialProps.ImageKey
		if key == "" {
			key = materialProps.CacheKey
		}
		materialProps.Texture = ctx.textureFor(materialProps.Image, key)
	}
	mat := ctx.materialFor(materialProps)
	mat.Incref()
	geom := ctx.geometryFor(c.props.Shape, materialProps.Texture != nil)
	geom.Incref()
	c.Native = graphic.NewMesh(geom, mat)
	c.node = c.Native
	ctx.Add(c.Native)
	c.bindTransform(c.props.Transform)
	watch(&c.unsubs, c.props.Material.Color, func(v Color) { mat.SetColor(&v) })
	watch(&c.unsubs, c.props.Material.Opacity, func(v float32) {
		mat.SetOpacity(v)
		mat.SetTransparent(v < 1)
	})
	return nil
}

func (c *MeshComponent) Close() {
	c.nodeComponent.Close()
	if c.Native != nil {
		c.Native.Dispose()
		c.Native = nil
	}
}

// CameraProps configures the shared perspective camera.
type CameraProps struct {
	FOV       float32
	Near      float32
	Far       float32
	Transform TransformProps
}

// CameraComponent configures the game's existing camera rather than replacing it.
type CameraComponent struct {
	nodeComponent
	props  CameraProps
	Native *camera.Camera
}

// PerspectiveCamera creates a camera component.
func PerspectiveCamera(props CameraProps) *CameraComponent {
	return &CameraComponent{props: props}
}

func (c *CameraComponent) Mount(ctx *SceneContext) error {
	if ctx.Game.Camera == nil {
		return fmt.Errorf("g3n: game has no perspective camera")
	}
	c.Native = ctx.Game.Camera
	c.node = c.Native
	if c.props.FOV > 0 {
		c.Native.SetFov(c.props.FOV)
	}
	if c.props.Near > 0 {
		c.Native.SetNear(c.props.Near)
	}
	if c.props.Far > 0 {
		c.Native.SetFar(c.props.Far)
	}
	c.bindTransform(c.props.Transform)
	return nil
}

func (c *CameraComponent) Close() {
	for _, unsubscribe := range c.unsubs {
		unsubscribe()
	}
	c.unsubs = nil
	c.spin = nil
	c.node = nil
	c.Native = nil
}

// LightProps configures ambient, point, or directional lighting.
type LightProps struct {
	Color     value.Value[Color]
	Intensity value.Value[float32]
	Transform TransformProps
}

type nativeLight interface {
	core.INode
	SetColor(*Color)
	SetIntensity(float32)
}

// LightComponent is a declarative light node.
type LightComponent struct {
	nodeComponent
	kind   string
	props  LightProps
	Native core.INode
}

// AmbientLight creates uniform scene lighting.
func AmbientLight(props LightProps) *LightComponent {
	return &LightComponent{kind: "ambient", props: props}
}

// PointLight creates an omnidirectional light at Transform.Position.
func PointLight(props LightProps) *LightComponent {
	return &LightComponent{kind: "point", props: props}
}

// DirectionalLight creates a direction-only light.
func DirectionalLight(props LightProps) *LightComponent {
	return &LightComponent{kind: "directional", props: props}
}

func (c *LightComponent) Mount(ctx *SceneContext) error {
	if c.Native != nil {
		return fmt.Errorf("g3n: light is already mounted")
	}
	var native nativeLight
	switch c.kind {
	case "ambient":
		native = light.NewAmbient(&Color{R: 1, G: 1, B: 1}, 1)
	case "point":
		native = light.NewPoint(&Color{R: 1, G: 1, B: 1}, 1)
	case "directional":
		native = light.NewDirectional(&Color{R: 1, G: 1, B: 1}, 1)
	default:
		return fmt.Errorf("g3n: unknown light kind %q", c.kind)
	}
	c.Native = native
	c.node = native
	ctx.Add(native)
	c.bindTransform(c.props.Transform)
	watch(&c.unsubs, c.props.Color, func(v Color) { native.SetColor(&v) })
	watch(&c.unsubs, c.props.Intensity, native.SetIntensity)
	return nil
}

func (c *LightComponent) Close() {
	c.nodeComponent.Close()
	if c.Native != nil {
		c.Native.Dispose()
		c.Native = nil
	}
}
