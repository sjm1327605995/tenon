package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/scheduler"
)

type Game struct {
	root         core.Component
	screenWidth  int
	screenHeight int
	isFirstFrame bool
}

func NewGame(root core.Component, width, height int) *Game {
	return &Game{
		root:         root,
		screenWidth:  width,
		screenHeight: height,
		isFirstFrame: true,
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenWidth = outsideWidth
	g.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	core.CalculateLayout(g.root, float32(g.screenWidth), float32(g.screenHeight))

	if g.isFirstFrame {
		scheduler.GetInstance().ProcessMount(g.root)
		g.isFirstFrame = false
	}

	g.root.Draw(screen)
	g.root.DrawOverlay(screen)
}

func (g *Game) Update() error {
	scheduler.GetInstance().ProcessUpdates()
	return g.root.Update()
}

func (g *Game) Close() error {
	scheduler.GetInstance().ProcessUnmount(g.root)
	return nil
}

func Run(root core.Component, width, height int) {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Tenon - UI Framework")

	game := NewGame(root, width, height)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

var backgroundColor = color.RGBA{R: 245, G: 245, B: 245, A: 255}

func ParseColor(hex string) color.RGBA {
	var r, g, b, a uint8 = 255, 255, 255, 255

	if len(hex) >= 7 && hex[0] == '#' {
		r = parseHexChannel(hex[1:3])
		g = parseHexChannel(hex[3:5])
		b = parseHexChannel(hex[5:7])
		if len(hex) >= 9 {
			a = parseHexChannel(hex[7:9])
		}
	}

	return color.RGBA{R: r, G: g, B: b, A: a}
}

func parseHexChannel(s string) uint8 {
	var v uint8
	for _, c := range s {
		var val uint8
		if c >= '0' && c <= '9' {
			val = uint8(c - '0')
		} else if c >= 'a' && c <= 'f' {
			val = uint8(c - 'a' + 10)
		} else if c >= 'A' && c <= 'F' {
			val = uint8(c - 'A' + 10)
		} else {
			val = 0
		}
		v = v*16 + val
	}
	return v
}
