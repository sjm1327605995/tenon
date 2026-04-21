package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
)

type Game struct {
	root         core.Component
	screenWidth  int
	screenHeight int
}

func NewGame(root core.Component, width, height int) *Game {
	return &Game{
		root:         root,
		screenWidth:  width,
		screenHeight: height,
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// 当窗口大小改变时，更新屏幕尺寸并重新计算布局
	g.screenWidth = outsideWidth
	g.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	// 每次绘制时都重新计算布局，确保布局与当前窗口大小匹配
	core.CalculateLayout(g.root, float32(g.screenWidth), float32(g.screenHeight))
	g.root.Draw(screen)
	g.root.DrawOverlay(screen)
}

func (g *Game) Update() error {
	return g.root.Update()
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
