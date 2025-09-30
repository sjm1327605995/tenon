// node.go implements the main game structure and Ebiten integration.
// This code is a modified version of examples/shapes from https://github.com/ebiten/ebiten
// Ebiten is licensed under the following license:
//
// Copyright 2017 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/sjm1327605995/tenon/core/context/node"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/millken/yoga"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game represents the main game structure
type Game struct {
	root *node.Context
}

// NewGame creates a new game instance
func NewGame() *Game {
	image, _, _ := ebitenutil.NewImageFromFile("gopher.png")

	return &Game{
		root: node.NewContext().View(func(v *node.View) {
			v.SetGap(yoga.GutterAll, 10).
				SetPadding(yoga.EdgeAll, 10).
				SetAlignItems(yoga.AlignCenter).
				SetJustifyContent(yoga.JustifyCenter).
				SetBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}).AddChild(
				node.NewView().SetWidth(100).
					SetHeight(100).
					SetRadius(node.RadiusAll, 20).
					SetBackgroundColor(color.RGBA{R: 186, G: 85, B: 211, A: 255}).
					SetOnHover(func(v *node.View) {
						v.SetBackgroundColor(color.RGBA{R: 176, G: 75, B: 201, A: 255})
					}).
					SetOnClick(func(v *node.View) {
						v.SetBackgroundColor(color.RGBA{R: 206, G: 105, B: 231, A: 255})
					}).
					SetMouseover(func(v *node.View) {
						v.SetBackgroundColor(color.RGBA{R: 186, G: 85, B: 211, A: 255})
					}),
				node.NewView().
					SetBorder(yoga.EdgeAll, 20).
					SetPadding(yoga.EdgeAll, 10).
					SetRadius(node.RadiusTopLeft, 50).
					SetRadius(node.RadiusTopRight, 100).
					SetRadius(node.RadiusBottomRight, 50).
					SetRadius(node.RadiusTopRight, 100).
					SetRadius(node.RadiusBottomLeft, 100).
					SetBackgroundColor(color.RGBA{R: 100, G: 149, B: 237, A: 255}).
					AddChild(node.NewImage(image)),
			)

		}),
	}

} // Update handles game logic updates
func (g *Game) Update() error {

	g.root.Update()

	return nil
}

// Draw renders the game screen
func (g *Game) Draw(screen *ebiten.Image) {

	g.root.Render(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()), 0, 0)
}

// Layout returns the game's screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	g.root.SetLayout(float32(outsideWidth), float32(outsideHeight))
	return outsideWidth, outsideHeight
}

// main initializes and starts the game
func main() {
	// Create a new game instance
	game := NewGame()

	// Configure window properties
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Go UI Layout - Ebiten Demo")
	ebiten.SetVsyncEnabled(false)
	// Start the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
