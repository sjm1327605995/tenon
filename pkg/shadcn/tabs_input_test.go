package shadcn

import (
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func center(q *ui.Query) (float32, float32) {
	b := q.Bounds()
	return b.X + b.W/2, b.Y + b.H/2
}

// An Input inside a Tabs panel must be clickable and typable: a real mousedown
// focuses the input itself (not a neighboring tab trigger), and typing updates
// the controlled value.
func TestTabsPanelInputInteractive(t *testing.T) {
	form := func(_ struct{}) *ui.Node {
		active, setActive := ui.UseState(0)
		v, setV := ui.UseState("")
		var panel *ui.Node
		if active == 0 {
			panel = ui.VStack(8, Label("名称"),
				Input(InputProps{Value: v, OnChange: setV, Placeholder: "name"}))
		} else {
			panel = Label("其他")
		}
		return ui.VStack(16, ui.Style(ui.Width(380)),
			Tabs(TabsProps{Tabs: []string{"账户", "密码"}, Active: active, OnChange: setActive}),
			ui.Div(ui.Style(ui.Padding(16)), panel))
	}

	h := ui.Mount(ui.Use(form, struct{}{}), 600, 400)
	if h.MouseDown(center(h.Root().ByPlaceholder("name"))).Kind() != "input" {
		t.Fatal("click in tab panel did not focus the input")
	}
	h.Root().ByPlaceholder("name").Type("hello")
	if got := h.Root().ByPlaceholder("name").Value(); got != "hello" {
		t.Fatalf("tab-panel input not typable: value=%q", got)
	}
}

// Each tab panel carries its own input; after switching tabs, the newly shown
// panel's input must still be focusable and typable. Guards the reported bug
// where only the first tab accepted input.
func TestTabsSwitchThenInput(t *testing.T) {
	form := func(_ struct{}) *ui.Node {
		active, setActive := ui.UseState(0)
		a, setA := ui.UseState("")
		p, setP := ui.UseState("")
		var panel *ui.Node
		if active == 0 {
			panel = ui.VStack(8, Label("名称"),
				Input(InputProps{Value: a, OnChange: setA, Placeholder: "acct"}))
		} else {
			panel = ui.VStack(8, Label("当前密码"),
				Input(InputProps{Value: p, OnChange: setP, Placeholder: "pwd"}))
		}
		return ui.VStack(16, ui.Style(ui.Width(380)),
			Tabs(TabsProps{Tabs: []string{"账户", "密码"}, Active: active, OnChange: setActive}),
			ui.Div(ui.Style(ui.Padding(16)), panel))
	}

	h := ui.Mount(ui.Use(form, struct{}{}), 600, 400)

	// Tab 0 input works.
	if h.MouseDown(center(h.Root().ByPlaceholder("acct"))).Kind() != "input" {
		t.Fatal("tab0: click did not focus input")
	}
	h.Root().ByPlaceholder("acct").Type("alice")
	if got := h.Root().ByPlaceholder("acct").Value(); got != "alice" {
		t.Fatalf("tab0 input not typable: %q", got)
	}

	// Switch to the 密码 tab with a real mousedown (clears focus AND fires OnChange).
	h.MouseDown(center(h.Root().ByText("密码")))

	in1 := h.Root().ByPlaceholder("pwd")
	if !in1.Exists() {
		t.Fatalf("tab1 input missing after switch; texts=%v", h.Root().Texts())
	}
	foc := h.MouseDown(center(in1))
	if foc.Kind() != "input" || foc.Placeholder() != "pwd" {
		t.Fatalf("BUG: after switch, focused kind=%q placeholder=%q; want input/pwd", foc.Kind(), foc.Placeholder())
	}
	h.Root().ByPlaceholder("pwd").Type("secret")
	if got := h.Root().ByPlaceholder("pwd").Value(); got != "secret" {
		t.Fatalf("BUG: after switch, password input not typable: value=%q", got)
	}
}
