// Package main demonstrates how to build a game-specific Card component
// on top of the Tenon UI framework. This is NOT part of Tenon itself.
package main

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/components"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Card renders a playing card with front/back faces, flip animation,
// and status indicator overlays. This is a game-specific extension,
// not part of the Tenon framework core.
type Card struct {
	core.BaseElement
	frontImage *ebiten.Image
	backImage  *ebiten.Image

	// State
	isFaceUp      bool
	flipAngle     float32 // 0=front, 180=back, animated
	isSelectable  bool
	isSelected    bool
	isHighlighted bool

	// Status overlays
	showNegate bool

	// Animation
	flipTween *core.Tween
}

// NewCard creates a Card element.
func NewCard() *Card {
	c := &Card{}
	c.Init(c)
	return c
}

// ElementType returns type identifier.
func (c *Card) ElementType() string { return "Card" }

// Draw renders the card with flip and overlays.
func (c *Card) Draw(screen *ebiten.Image) {
	if !c.IsVisible() {
		return
	}
	bounds := c.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	img := c.pickImage()
	if img == nil {
		return
	}

	// Compute flip scaleX: cos(angle) simulates Y-axis rotation
	rad := float64(c.flipAngle) * math.Pi / 180
	flipScaleX := float32(math.Cos(rad))
	if math.Abs(float64(flipScaleX)) < 0.001 {
		return // edge-on, nothing visible
	}

	// Combine user transform with card flip
	tr := c.GetTransform()
	combined := tr
	combined.ScaleX *= flipScaleX

	// Render full transform: Scale to bounds → Apply combined
	sw := float64(bounds.Width) / float64(img.Bounds().Dx())
	sh := float64(bounds.Height) / float64(img.Bounds().Dy())
	var geoM ebiten.GeoM
	geoM.Scale(sw, sh)
	geoM.Concat(core.BuildTransformGeoM(bounds, combined))

	op := &ebiten.DrawImageOptions{GeoM: geoM}
	core.ApplyColorScaleAlpha(&op.ColorScale, combined.Alpha)
	screen.DrawImage(img, op)

	// Draw selection / highlight border
	c.drawOverlays(screen, bounds, flipScaleX)
}

func (c *Card) pickImage() *ebiten.Image {
	if c.flipAngle < 90 || c.flipAngle > 270 {
		return c.frontImage
	}
	return c.backImage
}

func (c *Card) drawOverlays(screen *ebiten.Image, bounds core.LayoutBounds, flipScaleX float32) {
	if !c.isSelectable && !c.isHighlighted && !c.showNegate {
		return
	}

	tr := c.GetTransform()
	combined := tr
	combined.ScaleX *= flipScaleX

	sw := float64(bounds.Width) / 177
	sh := float64(bounds.Height) / 254
	var geoM ebiten.GeoM
	geoM.Scale(sw, sh)
	geoM.Concat(core.BuildTransformGeoM(bounds, combined))

	corners := [4][2]float32{{0, 0}, {177, 0}, {0, 254}, {177, 254}}
	screenCorners := make([][2]float32, 4)
	for i, p := range corners {
		sx, sy := geoM.Apply(float64(p[0]), float64(p[1]))
		screenCorners[i] = [2]float32{float32(sx), float32(sy)}
	}

	if c.showNegate {
		red := color.RGBA{255, 0, 0, 180}
		vector.StrokeLine(screen, screenCorners[0][0], screenCorners[0][1], screenCorners[3][0], screenCorners[3][1], 3, red, false)
		vector.StrokeLine(screen, screenCorners[1][0], screenCorners[1][1], screenCorners[2][0], screenCorners[2][1], 3, red, false)
	}
	if c.isSelectable || c.isHighlighted {
		var clr color.Color
		if c.isHighlighted {
			clr = color.RGBA{0, 255, 255, 255}
		} else if c.isSelected {
			clr = color.RGBA{255, 255, 0, 255}
		} else {
			clr = color.RGBA{255, 255, 0, 200}
		}
		for i := 0; i < 4; i++ {
			a, b := screenCorners[i], screenCorners[(i+1)%4]
			vector.StrokeLine(screen, a[0], a[1], b[0], b[1], 2, clr, false)
		}
	}
}

// Flip animates the card from front to back or back to front.
func (c *Card) Flip() {
	if c.flipTween != nil && c.flipTween.IsRunning() {
		c.flipTween.Stop()
	}
	from := c.flipAngle
	to := float32(180)
	if c.flipAngle >= 90 {
		to = 0
	}
	c.flipTween = core.NewTween(300*time.Millisecond, core.EaseInOutQuad)
	c.flipTween.OnUpdate(func(p float32) {
		c.flipAngle = from + (to-from)*p
		c.Mark(core.FlagNeedDraw)
	}).OnComplete(func() {
		c.isFaceUp = (c.flipAngle < 90)
		c.flipAngle = 0
		if !c.isFaceUp {
			c.flipAngle = 180
		}
	})
	c.flipTween.Start()
	if eng := c.GetEngine(); eng != nil {
		eng.AddAnimation(c.flipTween)
	}
}

// Setters
func (c *Card) SetFrontImage(img *ebiten.Image) *Card {
	c.frontImage = img
	c.Mark(core.FlagNeedDraw)
	return c
}
func (c *Card) SetBackImage(img *ebiten.Image) *Card {
	c.backImage = img
	c.Mark(core.FlagNeedDraw)
	return c
}
func (c *Card) SetFaceUp(v bool) *Card {
	if c.isFaceUp != v {
		c.isFaceUp = v
		if v {
			c.flipAngle = 0
		} else {
			c.flipAngle = 180
		}
		c.Mark(core.FlagNeedDraw)
	}
	return c
}
func (c *Card) SetSelectable(v bool) *Card {
	c.isSelectable = v
	c.Mark(core.FlagNeedDraw)
	return c
}
func (c *Card) SetSelected(v bool) *Card {
	c.isSelected = v
	c.Mark(core.FlagNeedDraw)
	return c
}
func (c *Card) SetHighlighted(v bool) *Card {
	c.isHighlighted = v
	c.Mark(core.FlagNeedDraw)
	return c
}
func (c *Card) SetShowNegate(v bool) *Card {
	c.showNegate = v
	c.Mark(core.FlagNeedDraw)
	return c
}

// ==================== Demo Widget ====================

type cardDemo struct {
	core.BaseWidget
	card *Card
}

func newCardDemo() *cardDemo {
	return &cardDemo{}
}

func (d *cardDemo) Render() core.Element {
	root := components.NewView()
	root.SetWidthPercent(100)
	root.SetHeightPercent(100)
	root.SetBackgroundColor(color.RGBA{30, 30, 50, 255})

	// Create placeholder card images (colored rectangles)
	front := ebiten.NewImage(177, 254)
	front.Fill(color.RGBA{200, 150, 100, 255})
	back := ebiten.NewImage(177, 254)
	back.Fill(color.RGBA{50, 80, 150, 255})

	d.card = NewCard()
	d.card.SetFrontImage(front)
	d.card.SetBackImage(back)
	d.card.SetFaceUp(true)
	d.card.SetSelectable(true)
	d.card.SetOrigin(0.5, 0.5)

	// Position card in center using absolute positioning
	d.card.SetPositionType(yoga.PositionTypeAbsolute)
	d.card.SetPosition(yoga.EdgeLeft, 200)
	d.card.SetPosition(yoga.EdgeTop, 100)
	d.card.SetWidth(177)
	d.card.SetHeight(254)

	// Button to flip
	btn := components.NewButton("Flip").SetOnClick(func() {
		d.card.Flip()
	})
	btn.SetPositionType(yoga.PositionTypeAbsolute)
	btn.SetPosition(yoga.EdgeLeft, 20)
	btn.SetPosition(yoga.EdgeTop, 20)

	// Button to rotate
	rotBtn := components.NewButton("Rotate").SetOnClick(func() {
		d.card.SetRotation(d.card.GetTransform().Rotation + 15)
	})
	rotBtn.SetPositionType(yoga.PositionTypeAbsolute)
	rotBtn.SetPosition(yoga.EdgeLeft, 20)
	rotBtn.SetPosition(yoga.EdgeTop, 60)

	root.Add(btn, rotBtn, d.card)
	return root
}
