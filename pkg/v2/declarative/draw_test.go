package declarative

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/yoga"
)

func TestTextDraw(t *testing.T) {
	if err := fonts.InitDefaultFont(); err != nil {
		t.Skip("skip: no font available")
	}

	// 1. 创建声明式 Text widget 并 mount
	w := Text("Hello").FontSize(24).Color(Black)
	el := w.CreateElement()
	el.Mount(nil, 0)

	// 2. 获取 RenderObject
	ro := el.FindRenderObject()
	if ro == nil {
		t.Fatal("FindRenderObject returned nil")
	}
	rt, ok := ro.(*render.RenderText)
	if !ok {
		t.Fatalf("expected *render.RenderText, got %T", ro)
	}

	// 3. 手动执行 yoga layout
	y := rt.GetYoga()
	if y == nil {
		t.Fatal("RenderText has no yoga node")
	}
	y.StyleSetWidth(200)
	y.StyleSetHeight(50)
	y.CalculateLayout(200, 50, yoga.DirectionLTR)

	// 4. 同步 bounds
	rt.SetBounds(render.Bounds{
		X:      y.LayoutLeft(),
		Y:      y.LayoutTop(),
		Width:  y.LayoutWidth(),
		Height: y.LayoutHeight(),
	})

	t.Logf("RenderText bounds: %+v", rt.GetBounds())
	t.Logf("RenderText Content: %q", rt.Content)

	// 5. 绘制到离屏图像
	img := ebiten.NewImage(400, 400)
	rt.Paint(img, render.Offset{X: 10, Y: 10})

	// 6. 简单检查：图像不应全透明（至少有一些像素被绘制）
	allTransparent := true
	for x := 0; x < 400 && allTransparent; x++ {
		for y := 0; y < 400; y++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			if a > 0 || r > 0 || g > 0 || b > 0 {
				allTransparent = false
				break
			}
		}
	}
	if allTransparent {
		t.Fatal("RenderText painted nothing (image is all transparent)")
	}
	t.Log("RenderText successfully painted to image")
}
