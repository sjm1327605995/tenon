package router

import (
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 验证栈式导航：初始屏 -> Push 进详情屏 -> Pop 返回；以及 Replace 不改深度。
func TestStackNavigation(t *testing.T) {
	list := func(_ Params) *ui.Node {
		nav := UseNavigate()
		return ui.Button(ui.OnClick(func() { nav.Push("detail", Params{"id": "7"}) }), ui.Text("open"))
	}
	detail := func(p Params) *ui.Node {
		nav := UseNavigate()
		return ui.Div(
			ui.Text("detail "+p["id"]),
			ui.Button(ui.OnClick(func() { nav.Replace("detail", Params{"id": "8"}) }), ui.Text("next")),
			ui.Button(ui.OnClick(func() { nav.Pop() }), ui.Text("back")),
		)
	}
	h := ui.MountDefault(Router(Props{
		Initial: "list",
		Screens: map[string]Screen{"list": list, "detail": detail},
	}))

	if !h.Root().ByText("open").Exists() {
		t.Fatalf("初始 list 屏未渲染；texts=%v", h.Root().Texts())
	}

	clickText(t, h, "open") // Push -> detail(7)
	if !h.Root().ByText("detail 7").Exists() {
		t.Fatalf("Push 后详情屏未渲染；texts=%v", h.Root().Texts())
	}

	clickText(t, h, "next") // Replace -> detail(8)，深度不变
	if !h.Root().ByText("detail 8").Exists() {
		t.Fatalf("Replace 后未更新参数；texts=%v", h.Root().Texts())
	}

	clickText(t, h, "back") // Pop -> 回到 list（说明 Replace 没有增加深度）
	if !h.Root().ByText("open").Exists() {
		t.Fatalf("Pop 后未回到 list（Replace 误增深度？）；texts=%v", h.Root().Texts())
	}
}

func clickText(t *testing.T, h *ui.Harness, s string) {
	t.Helper()
	q := h.Root().ByText(s)
	if !q.Exists() {
		t.Fatalf("未找到文本 %q", s)
	}
	b := q.Bounds()
	if !h.ClickAt(b.X+b.W/2, b.Y+b.H/2) {
		t.Fatalf("点击 %q 没有命中可点击元素", s)
	}
}
