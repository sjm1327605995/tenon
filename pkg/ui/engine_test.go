package ui

import (
	"fmt"
	"testing"

	"github.com/sjm1327605995/tenon/pkg/font"
	"github.com/sjm1327605995/tenon/yoga"
)

var (
	parentRenders int
	childRenders  int
)

func testChild(_ struct{}) *Node {
	childRenders++
	count, setCount := UseState(0)
	return Button(
		Style(Width(60), Height(30)),
		OnClick(func() { setCount(count + 1) }),
		Text(fmt.Sprintf("%d", count)),
	)
}

func testParent(_ struct{}) *Node {
	parentRenders++
	return Div(
		Style(Column, Gap(8), Padding(10), Width(240), Height(200)),
		Text("static-parent"),
		Use(testChild, struct{}{}),
	)
}

func layoutAll(g *game) {
	g.rootRN = rootRenderNode(g.rootFiber)
	relink(g.rootFiber)
	if g.rootRN != nil {
		resolveInherited(g.rootRN, Black, false, 16, false)
		g.rootRN.yn.CalculateLayout(float32(g.w), float32(g.h), yoga.DirectionLTR)
		computeBounds(g.rootRN, 0, 0)
	}
	g.layoutPortals(true)
	if g.rootRN != nil {
		syncMeasures(g.rootRN)
	}
}

func collectTexts(rn *renderNode, out *[]string) {
	if rn.kind == rnText {
		*out = append(*out, rn.text)
	}
	for _, c := range rn.children {
		collectTexts(c, out)
	}
}

func findClickable(rn *renderNode) *renderNode {
	if rn.onClick != nil {
		return rn
	}
	for _, c := range rn.children {
		if r := findClickable(c); r != nil {
			return r
		}
	}
	return nil
}

// TestLocalRerender 验证：子组件 setState 只触发子组件重渲染，父组件不重跑。
func TestLocalRerender(t *testing.T) {
	if err := font.InitDefaultFont(); err != nil {
		t.Fatal(err)
	}
	parentRenders, childRenders = 0, 0
	g := &game{w: 400, h: 400}
	activeGame = g

	g.rootFiber = reconcile(nil, nil, Use(testParent, struct{}{}))
	layoutAll(g)

	if parentRenders != 1 || childRenders != 1 {
		t.Fatalf("initial renders: parent=%d child=%d, want 1/1", parentRenders, childRenders)
	}

	var texts []string
	collectTexts(g.rootRN, &texts)
	if !contains(texts, "0") || !contains(texts, "static-parent") {
		t.Fatalf("initial texts = %v", texts)
	}

	// 命中 + 按钮并点击：用 hitTest 走真实命中路径
	btn := findClickable(g.rootRN)
	if btn == nil {
		t.Fatal("no clickable found")
	}
	cx := btn.bounds.X + btn.bounds.W/2
	cy := btn.bounds.Y + btn.bounds.H/2
	handler := hitTest(g.rootRN, cx, cy)
	if handler == nil {
		t.Fatalf("hitTest missed button at (%.1f,%.1f)", cx, cy)
	}
	handler() // setCount(1) -> markDirty(childFiber)

	if len(g.dirty) != 1 {
		t.Fatalf("dirty queue = %d, want 1 (only child)", len(g.dirty))
	}
	g.flushDirty()
	layoutAll(g)

	if parentRenders != 1 {
		t.Fatalf("parent re-rendered (=%d); expected local re-render only", parentRenders)
	}
	if childRenders != 2 {
		t.Fatalf("child renders=%d, want 2", childRenders)
	}

	texts = nil
	collectTexts(g.rootRN, &texts)
	if !contains(texts, "1") {
		t.Fatalf("after click texts = %v, want a \"1\"", texts)
	}
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}
