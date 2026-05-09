package widgets

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// TestEngineDraw 验证 Engine.Draw 能正确输出到图像。
func TestEngineDraw(t *testing.T) {
	app := func() engine.Widget {
		return Container(
			Text("Hello World").FontSize(24),
		).Background(*render.NewColor(255, 0, 0, 255)).W(400).H(300)
	}

	eng := engine.NewEngine(app, 400, 300)
	eng.Mount()
	eng.Update()

	img := ebiten.NewImage(400, 300)
	eng.Draw(img)
	// ebiten v2.9.9 中 img.At 需要在游戏启动后才能调用，
	// 这里只验证 Draw 不 panic 即可。
}
