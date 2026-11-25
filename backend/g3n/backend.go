// Package g3n connects Tenon's declarative UI to the G3N game engine.
// Importing this module selects G3N for ui.Run; the core module remains on Gio.
package g3n

import (
	"image"
	"os"
	"strconv"
	"time"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/window"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

// Game exposes the native G3N scene and camera. Launch and Update can build
// and animate meshes, models, lights, physics, and audio without Tenon wrappers.
type Game struct {
	Scene    *core.Node
	Camera   *camera.Camera
	Window   window.IWindow
	Renderer *renderer.Renderer
	scenes   []*Scene
}

// Options configures native game callbacks.
type Options struct {
	Launch        func(*Game)
	Update        func(*Game, time.Duration)
	ClearColor    *Color
	SceneViewport func(width, height int) (x, y, widthOut, heightOut int)
}

// Backend is a Tenon rendering backend backed by G3N.
type Backend struct {
	options Options
	fonts   map[fontKey]*font
}

var activeWindow *window.GlfwWindow

// Exit requests the active G3N window to close.
func Exit() {
	if activeWindow != nil {
		activeWindow.SetShouldClose(true)
	}
}

// Minimise iconifies the active G3N window.
func Minimise() {
	if activeWindow != nil {
		activeWindow.Iconify()
	}
}

// Use selects G3N for subsequent ui.Run calls.
func Use(options Options) {
	ui.SetBackend(&Backend{options: options, fonts: make(map[fontKey]*font)})
}

func init() { Use(Options{}) }

func smokeFrameLimit() int {
	value, err := strconv.Atoi(os.Getenv("TENON_SMOKE_FRAMES"))
	if err != nil || value <= 0 {
		return 0
	}
	return value
}

func (b *Backend) NewFont(px float32, weight int, italic bool) ui.FontFace {
	key := fontKey{px: px, weight: weight, italic: italic}
	if cached := b.fonts[key]; cached != nil {
		return cached
	}
	face, err := text.NewFontFromData(ui.DefaultFontData())
	if err != nil {
		return nil
	}
	cached := &font{font: face, px: px, fauxBold: weight >= 600, fauxItalic: italic}
	face.SetPointSize(float64(px))
	b.fonts[key] = cached
	return cached
}

func (b *Backend) NewBitmap(img image.Image) ui.Bitmap {
	if img == nil {
		return nil
	}
	return newBitmap(img)
}

func (b *Backend) NewVecPath(svgPath string, scale float32) ui.VecPath {
	return &vectorPath{svg: svgPath, scale: scale}
}

func (b *Backend) Capabilities() ui.BackendCapabilities {
	return ui.BackendCapabilities{SVGPaths: true, LayerTransforms: true}
}

func (b *Backend) Run(root *ui.Node, cfg ui.WindowConfig) {
	w, h := cfg.Size()
	if err := window.Init(w, h, cfg.Title()); err != nil {
		panic(err)
	}
	win := window.Get()
	activeWindow = win.(*window.GlfwWindow)
	defer func() {
		if activeWindow == win {
			activeWindow = nil
		}
		win.Destroy()
	}()
	rend := renderer.NewRenderer(win.Gls())
	rend.AddShader("tenon_vector_vertex", vectorVertexShader)
	rend.AddShader("tenon_vector_fragment", vectorFragmentShader)
	rend.AddProgram("tenon_vector", "tenon_vector_vertex", "tenon_vector_fragment")
	rend.AddShader("tenon_unlit_vertex", unlitVertexShader)
	rend.AddShader("tenon_unlit_fragment", unlitFragmentShader)
	rend.AddProgram("tenon_unlit", "tenon_unlit_vertex", "tenon_unlit_fragment")
	if err := rend.AddDefaultShaders(); err != nil {
		panic(err)
	}
	scene := core.NewNode()
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)
	game := &Game{Scene: scene, Camera: cam, Window: win, Renderer: rend}
	resize := func() {
		fbw, fbh := win.GetFramebufferSize()
		win.Gls().Viewport(0, 0, int32(fbw), int32(fbh))
		cam.SetAspect(float32(fbw) / float32(fbh))
	}
	win.Subscribe(window.OnWindowSize, func(string, interface{}) { resize() })
	resize()
	gui.Manager().Set(scene)
	if b.options.Launch != nil {
		b.options.Launch(game)
	}
	defer game.closeScenes()

	input := newInput(win)
	runtime := ui.NewBackendRuntime(root, cfg, input)
	defer runtime.Close()
	uiRoot := core.NewNode()
	lastPaintVersion := uint64(0)
	paintW, paintH := 0, 0
	painted := false
	last := time.Now()
	smokeRemaining := smokeFrameLimit()
	capture := nativeCaptureFromEnv()
	currentCursor := window.ArrowCursor
	for !win.(*window.GlfwWindow).ShouldClose() {
		now := time.Now()
		dt := now.Sub(last)
		last = now
		if b.options.Update != nil {
			b.options.Update(game, dt)
		}
		fbw, fbh := win.GetFramebufferSize()
		scaleX, _ := win.GetScale()
		runtime.Resize(fbw, fbh, float32(scaleX))
		cursorX, cursorY := input.Cursor()
		uiInteractive := runtime.PointerInteractive(cursorX, cursorY)
		runtime.Frame(float32(dt.Seconds() * 1000))
		if next := nativeCursor(runtime.PointerCursor(cursorX, cursorY)); next != currentCursor {
			win.(*window.GlfwWindow).SetCursor(next)
			currentCursor = next
		}
		viewportX, viewportY, viewportW, viewportH := 0, 0, fbw, fbh
		if b.options.SceneViewport != nil {
			viewportX, viewportY, viewportW, viewportH = b.options.SceneViewport(fbw, fbh)
		}
		cam.SetAspect(float32(viewportW) / float32(viewportH))
		if !uiInteractive {
			game.HandlePointer(pointerEventAt(cursorX, cursorY, viewportX, viewportY, viewportW, viewportH,
				input.MouseJustPressed(ui.MouseLeft), input.MouseJustPressed(ui.MouseRight)))
		}
		paintVersion := runtime.PaintVersion(now)
		if !painted || paintVersion != lastPaintVersion || paintW != fbw || paintH != fbh {
			uiRoot.DisposeChildren(true)
			runtime.Paint(newPainter(uiRoot, win.Gls(), fbw, fbh))
			lastPaintVersion, paintW, paintH, painted = paintVersion, fbw, fbh, true
		}
		if b.options.ClearColor != nil {
			color := b.options.ClearColor
			win.Gls().ClearColor(color.R, color.G, color.B, 1)
		}
		win.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		win.Gls().Viewport(int32(viewportX), int32(fbh-viewportY-viewportH), int32(viewportW), int32(viewportH))
		win.Gls().Enable(gls.SCISSOR_TEST)
		win.Gls().Scissor(int32(viewportX), int32(fbh-viewportY-viewportH), uint32(viewportW), uint32(viewportH))
		if err := rend.Render(scene, cam); err != nil {
			panic(err)
		}
		win.Gls().Disable(gls.SCISSOR_TEST)
		win.Gls().Viewport(0, 0, int32(fbw), int32(fbh))
		cam.SetAspect(float32(fbw) / float32(fbh))
		if err := rend.Render(uiRoot, cam); err != nil {
			panic(err)
		}
		if capture.capture(win.Gls(), fbw, fbh) {
			win.(*window.GlfwWindow).SetShouldClose(true)
		}
		win.(*window.GlfwWindow).SwapBuffers()
		input.endFrame()
		win.(*window.GlfwWindow).PollEvents()
		if smokeRemaining > 0 {
			smokeRemaining--
			if smokeRemaining == 0 {
				win.(*window.GlfwWindow).SetShouldClose(true)
			}
		}
	}
}
