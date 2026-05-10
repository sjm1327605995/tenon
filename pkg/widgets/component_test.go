package widgets

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/render"
)

// verify RenderComponent can be embedded and delegates correctly via Impl
type testComp struct {
	RenderComponent
	created bool
}

func (t *testComp) RenderObject(element engine.Element) render.RenderObject {
	t.created = true
	return render.NewRenderBox()
}

func TestRenderComponent(t *testing.T) {
	c := &testComp{}
	c.RenderComponent.Impl = c

	el := c.CreateElement()
	if el == nil {
		t.Fatal("CreateElement returned nil")
	}

	// Mount triggers CreateRenderObject via RenderObjectFactory assertion
	el.Mount(nil, 0)
	if !c.created {
		t.Fatal("CreateRenderObject was not called")
	}
}
