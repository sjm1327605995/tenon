package g3n

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/sjm1327605995/tenon/pkg/value"
)

// Vec3 is a native G3N vector used by transforms and cameras.
type Vec3 = math32.Vector3

// Color is a native G3N linear RGB color.
type Color = math32.Color

// V3 constructs a vector.
func V3(x, y, z float32) Vec3 { return Vec3{X: x, Y: y, Z: z} }

// RGB constructs an opaque linear color.
func RGB(r, g, b float32) Color { return Color{R: r, G: g, B: b} }

// Component is a mountable game object. Mount must create native resources once;
// Update may animate them, and Close must release subscriptions.
type Component interface {
	Mount(*SceneContext) error
	Update(time.Duration)
	Close()
}

// SceneContext provides the native scene root and per-scene resource caches.
type SceneContext struct {
	Game            *Game
	parent          *core.Node
	geometries      map[string]*geometry.Geometry
	ownedGeometries map[*geometry.Geometry]bool
	materials       map[string]*material.Standard
	textures        map[string]*texture.Texture2D
}

// Add attaches a native node to the declarative scene root.
func (c *SceneContext) Add(node core.INode) { c.parent.Add(node) }

// Scene owns a declarative game tree and its native resources.
type Scene struct {
	children []Component
	game     *Game
	root     *core.Node
	context  *SceneContext
}

// GameScene creates a scene from native 3D components.
func GameScene(children ...Component) *Scene {
	return &Scene{children: append([]Component(nil), children...)}
}

// Mount creates the scene once and caches shared geometry and materials.
func (s *Scene) Mount(game *Game) error {
	if game == nil || game.Scene == nil {
		return errors.New("g3n: scene requires a native game")
	}
	if s.game != nil {
		return errors.New("g3n: scene is already mounted")
	}
	s.game = game
	s.root = core.NewNode()
	s.context = &SceneContext{
		Game:            game,
		parent:          s.root,
		geometries:      make(map[string]*geometry.Geometry),
		ownedGeometries: make(map[*geometry.Geometry]bool),
		materials:       make(map[string]*material.Standard),
		textures:        make(map[string]*texture.Texture2D),
	}
	game.Scene.Add(s.root)
	game.scenes = append(game.scenes, s)
	for _, child := range s.children {
		if err := child.Mount(s.context); err != nil {
			s.Close()
			return err
		}
	}
	return nil
}

// Launch mounts the scene and panics on initialization failure, matching ui.Run.
func (s *Scene) Launch(game *Game) {
	if err := s.Mount(game); err != nil {
		panic(err)
	}
}

// Update advances component animation and dynamic properties.
func (s *Scene) Update(dt time.Duration) {
	if s.game == nil {
		return
	}
	for _, child := range s.children {
		child.Update(dt)
	}
}

// Frame adapts Scene.Update to Options.Update.
func (s *Scene) Frame(_ *Game, dt time.Duration) { s.Update(dt) }

// Close releases subscriptions and detaches the scene without touching game objects
// owned by other scenes or by native G3N code.
func (s *Scene) Close() {
	if s.game == nil {
		return
	}
	for i := len(s.children) - 1; i >= 0; i-- {
		s.children[i].Close()
	}
	for _, cached := range s.context.materials {
		cached.Dispose()
	}
	for _, cached := range s.context.textures {
		cached.Dispose()
	}
	for cached := range s.context.ownedGeometries {
		cached.Dispose()
	}
	s.game.Scene.Remove(s.root)
	for i, scene := range s.game.scenes {
		if scene == s {
			s.game.scenes = append(s.game.scenes[:i], s.game.scenes[i+1:]...)
			break
		}
	}
	s.game, s.root, s.context = nil, nil, nil
}

// Root returns the native node containing this declarative scene.
func (s *Scene) Root() *core.Node { return s.root }

func (g *Game) closeScenes() {
	for i := len(g.scenes) - 1; i >= 0; i-- {
		g.scenes[i].Close()
	}
}

func bindingKey(v any) string {
	if v == nil {
		return "nil"
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		return fmt.Sprintf("%x", rv.Pointer())
	}
	return fmt.Sprintf("%v", v)
}

func watch[T any](unsubs *[]func(), source value.Value[T], apply func(T)) {
	if source == nil {
		return
	}
	apply(source.Get())
	*unsubs = append(*unsubs, source.Subscribe(apply))
}
