package widgets

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// TestGallerySmoke 模拟 Gallery 的构建和布局流程，验证根节点 bounds 非零。
func TestGallerySmoke(t *testing.T) {
	app := func() engine.Widget {
		return NewAnimatedContainer().
			WithChild(Text("Hello Gallery").FontSize(16)).
			WithSize(400, 300).
			WithBackground(*render.NewColor(200, 200, 200, 255))
	}

	eng := engine.NewEngine(app, 900, 800)
	eng.Mount()

	// Simulate first frame Update
	eng.Update()

	ro := eng.GetRootRenderObject()
	if ro == nil {
		t.Fatal("rootRenderObject is nil")
	}

	bounds := ro.GetBounds()
	t.Logf("root bounds: %+v", bounds)

	if bounds.Width <= 0 || bounds.Height <= 0 {
		t.Errorf("expected root bounds > 0, got %+v", bounds)
	}

	// Walk all children and check bounds
	var check func(ro render.RenderObject, depth int)
	check = func(ro render.RenderObject, depth int) {
		b := ro.GetBounds()
		y := ro.GetYoga()
		t.Logf("%*s%T bounds=%+v yoga=%v", depth*2, "", ro, b, y != nil)
		for _, child := range ro.GetChildren() {
			check(child, depth+1)
		}
	}
	check(ro, 0)
}

// TestGalleryScrollSmoke 模拟 Gallery 的 Scroll + Column + Text 结构
func TestGalleryScrollSmoke(t *testing.T) {
	app := func() engine.Widget {
		return Scroll(
			Column(
				Text("Line 1").FontSize(16),
				Text("Line 2").FontSize(16),
				Text("Line 3").FontSize(16),
			).Gapf(8),
		)
	}

	eng := engine.NewEngine(app, 900, 800)
	eng.Mount()
	eng.Update()

	ro := eng.GetRootRenderObject()
	if ro == nil {
		t.Fatal("rootRenderObject is nil")
	}

	bounds := ro.GetBounds()
	t.Logf("root bounds: %+v", bounds)

	if bounds.Width <= 0 || bounds.Height <= 0 {
		t.Errorf("expected root bounds > 0, got %+v", bounds)
	}

	var check func(ro render.RenderObject, depth int)
	check = func(ro render.RenderObject, depth int) {
		b := ro.GetBounds()
		t.Logf("%*s%T bounds=%+v children=%d", depth*2, "", ro, b, len(ro.GetChildren()))
		for _, child := range ro.GetChildren() {
			check(child, depth+1)
		}
	}
	check(ro, 0)
}
