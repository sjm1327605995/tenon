package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/types"
	"image/color"
)

type Drawable interface {
	Draw(screen *ebiten.Image, x, y int)
}

var drawableRegistry = make(map[string]func(element types.Element) Drawable)

func RegisterDrawable(elementType string, creator func(element types.Element) Drawable) {
	drawableRegistry[elementType] = creator
}

type EbitenRenderer struct {
	screenWidth  int
	screenHeight int
}

func NewEbitenRenderer(width, height int) *EbitenRenderer {
	return &EbitenRenderer{
		screenWidth:  width,
		screenHeight: height,
	}
}

func (r *EbitenRenderer) GetScreenSize() (int, int) {
	return r.screenWidth, r.screenHeight
}

func (r *EbitenRenderer) Render(screen *ebiten.Image, element types.Element) {
	if element == nil {
		return
	}

	screen.Fill(backgroundColor)
	r.renderElement(screen, element, 0, 0)
}

func (r *EbitenRenderer) renderElement(screen *ebiten.Image, element types.Element, parentX, parentY int) {
	if element == nil {
		return
	}

	layout := element.GetLayout()
	x := parentX + int(layout.X)
	y := parentY + int(layout.Y)

	creator, ok := drawableRegistry[element.Type()]
	if !ok {
		return
	}

	drawable := creator(element)
	drawable.Draw(screen, x, y)

	for _, child := range element.GetChildren() {
		if child == nil {
			continue
		}
		r.renderElement(screen, child, x, y)
	}
}

type Game struct {
	renderer *EbitenRenderer
	root     types.Element
}

func NewGame(root types.Element, width, height int) *Game {
	RegisterDrawables()
	return &Game{
		renderer: NewEbitenRenderer(width, height),
		root:     root,
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.renderer.screenWidth, g.renderer.screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Render(screen, g.root)
}

func (g *Game) Update() error {
	return nil
}

func Run(root types.Element, width, height int) {
	ebiten.SetWindowSize(width, height)
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