package core

import (
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// RenderEngine 负责绘制管理。
type RenderEngine struct {
	engine *Engine
}

func newRenderEngine(e *Engine) *RenderEngine {
	return &RenderEngine{engine: e}
}

func (r *RenderEngine) drawScreen(screen *ebiten.Image) {
	drawStart := time.Now()
	screen.Fill(color.RGBA{R: 245, G: 245, B: 245, A: 255})
	if r.engine.rootElement != nil {
		r.drawElement(screen, r.engine.rootElement, 0, 0)
	}
	for _, overlay := range r.engine.overlays {
		if overlay != nil && overlay.IsVisible() {
			r.drawElement(screen, overlay, 0, 0)
		}
	}
	r.engine.perf.LastDrawTime = time.Since(drawStart)
}

func (r *RenderEngine) drawElement(screen *ebiten.Image, el Element, offsetX, offsetY float32) {
	if el == nil || !el.IsVisible() {
		return
	}
	bounds := el.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 若标记了 ClipChildren，子节点绘制在裁剪后的子图上
	var childScreen *ebiten.Image = screen
	childOffsetX := offsetX
	childOffsetY := offsetY
	if el.HasFlag(FlagClipChildren) {
		sub := screen.SubImage(image.Rect(
			int(bounds.X), int(bounds.Y),
			int(bounds.X+bounds.Width), int(bounds.Y+bounds.Height),
		))
		if subImg, ok := sub.(*ebiten.Image); ok {
			childScreen = subImg
			// SubImage uses the same coordinate system as the original image;
			// it only acts as a clip mask. Children should use screen coordinates.
			childOffsetX = 0
			childOffsetY = 0
		}
	}

	// Use screen coordinates for drawing. SubImage will clip automatically.
	relBounds := LayoutBounds{
		X:      bounds.X - childOffsetX,
		Y:      bounds.Y - childOffsetY,
		Width:  bounds.Width,
		Height: bounds.Height,
	}
	el.SetBounds(relBounds)
	el.Draw(childScreen)
	el.SetBounds(bounds)

	for _, child := range el.GetChildren() {
		if r.engine.IsOverlay(child) {
			continue
		}
		r.drawElement(childScreen, child, childOffsetX, childOffsetY)
	}
}
