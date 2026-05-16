package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/canvas"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/widget"
)

// App implements ebiten.Game and hosts a widget tree.
type App struct {
	root       widget.Widget
	ctx        *widget.ContextImpl
	screenSize geometry.Size
}

// New creates a new App with the given root widget.
func New(root widget.Widget) *App {
	ctx := widget.NewContext()
	return &App{
		root: root,
		ctx:  ctx,
	}
}

// SetRoot sets a new root widget.
func (a *App) SetRoot(root widget.Widget) {
	a.root = root
}

// Layout implements ebiten.Game.
func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	a.screenSize = geometry.Sz(float32(outsideWidth), float32(outsideHeight))
	a.ctx.SetWindowSize(a.screenSize)
	return outsideWidth, outsideHeight
}

// Update implements ebiten.Game.
func (a *App) Update() error {
	if a.root == nil {
		return nil
	}

	// Run layout pass
	constraints := geometry.Constraints{
		MinWidth:  0,
		MinHeight: 0,
		MaxWidth:  a.screenSize.Width,
		MaxHeight: a.screenSize.Height,
	}
	a.root.Layout(a.ctx, constraints)

	// Ensure root bounds fill the screen
	if sb, ok := a.root.(interface{ SetBounds(geometry.Rect) }); ok {
		sb.SetBounds(geometry.FromPointSize(geometry.Pt(0, 0), a.screenSize))
	}

	return nil
}

// Draw implements ebiten.Game.
func (a *App) Draw(screen *ebiten.Image) {
	if a.root == nil {
		return
	}

	c := canvas.New(screen)
	c.Clear(widget.ColorWhite)
	widget.DrawTree(a.root, a.ctx, c)
}

// Run starts the ebiten game loop.
func (a *App) Run() error {
	ebiten.SetWindowTitle("Tenon")
	ebiten.SetWindowSize(800, 600)
	return ebiten.RunGame(a)
}
