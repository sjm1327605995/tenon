package widgets

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

func TestButtonSmoke(t *testing.T) {
	app := func() ui.Widget {
		return Button("Click Me").OnTap(func() {})
	}

	eng := ui.NewEngine(app, 900, 800)
	eng.Mount()
	eng.Update()

	ro := eng.GetRootRenderObject()
	if ro == nil {
		t.Fatal("rootRenderObject is nil")
	}

	bounds := ro.GetBounds()
	t.Logf("root bounds: %+v", bounds)

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
