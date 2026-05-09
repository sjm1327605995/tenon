package widgets

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sjm1327605995/tenon/internal/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

func dumpRenderObjectTree(ro render.RenderObject, indent int) string {
	if ro == nil {
		return ""
	}
	var sb strings.Builder
	b := ro.GetBounds()
	name := fmt.Sprintf("%T", ro)
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		name = name[idx+1:]
	}
	sb.WriteString(fmt.Sprintf("%s%s bounds=%+v z=%d visible=%v\n",
		strings.Repeat("  ", indent), name, b, ro.GetZIndex(), ro.IsVisible()))
	for _, child := range ro.GetChildren() {
		sb.WriteString(dumpRenderObjectTree(child, indent+1))
	}
	return sb.String()
}

func TestSelectRenderObjectTree(t *testing.T) {
	selectWidget := Select([]SelectOption{
		{Value: "a", Label: "Option A"},
		{Value: "b", Label: "Option B"},
	})

	env := ui.TestWidget(t, selectWidget, 400, 300)

	t.Logf("=== CLOSED ===\n%s", dumpRenderObjectTree(env.RootElement, 0))

	// 打开
	env.TapAt(100, 20)
	env.Rebuild()

	t.Logf("=== OPEN ===\n%s", dumpRenderObjectTree(env.RootElement, 0))

	// 查找所有 RenderText
	var findTexts func(ro render.RenderObject)
	findTexts = func(ro render.RenderObject) {
		if ro == nil {
			return
		}
		if rt, ok := ro.(*render.RenderText); ok {
			t.Logf("RenderText: %q bounds=%+v", rt.Content, ro.GetBounds())
		}
		for _, child := range ro.GetChildren() {
			findTexts(child)
		}
	}
	findTexts(env.RootElement)
}
