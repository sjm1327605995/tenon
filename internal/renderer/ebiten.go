package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
)

// Game 实现了 ebiten.Game 接口，是 Tenon 与 Ebiten 的桥接层。
type Game struct {
	engine *core.Engine
}

// NewGame 创建一个新的 Game 实例。
// root 可以是 Host 或 Widget；如果是 Widget，会先 Render 一次拿到根 Host。
func NewGame(root core.Component, width, height int) *Game {
	var rootHost core.Host

	if h, ok := root.(core.Host); ok {
		rootHost = h
	} else if w, ok := root.(core.Widget); ok {
		rendered := w.Render()
		if h, ok := rendered.(core.Host); ok {
			rootHost = h
		}
	}

	if rootHost == nil {
		panic("tenon: root component must be a Host or render to a Host")
	}

	engine := core.NewEngine(rootHost, width, height)
	engine.Mount()

	return &Game{engine: engine}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.engine.SetScreenSize(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	return g.engine.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.engine.Draw(screen)
}

// Run 启动 Tenon 应用。
func Run(root core.Component, width, height int) {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Tenon")

	game := NewGame(root, width, height)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
