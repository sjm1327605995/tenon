// Run from the repository root: go run _debug_snapshot/test_v3.go
// Exercises the current Node/Fiber engine without opening a window.
package main

import (
	"fmt"
	"log"
	"math"

	d "github.com/sjm1327605995/tenon/pkg/decl"
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

var boundsChecks int

func require(ok bool, message string) {
	if !ok {
		log.Fatal(message)
	}
}

func bounds(q *ui.Query) ui.Rect {
	require(q.Exists(), "expected UI node is missing")
	b := q.Bounds()
	require(b.W > 0 && b.H > 0, fmt.Sprintf("invalid bounds: %+v", b))
	for _, v := range []float32{b.X, b.Y, b.W, b.H} {
		require(!math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0), "non-finite bounds")
	}
	boundsChecks++
	return b
}

func verifyForms() {
	h := ui.Mount(ui.Use(func(_ struct{}) *ui.Node {
		disabled, setDisabled := ui.UseState(false)
		value, setValue := ui.UseState("initial")
		return ui.VStack(12, ui.Style(ui.Width(320), ui.Padding(16)),
			shadcn.Button(shadcn.ButtonProps{OnClick: func() { setDisabled(!disabled) }}, ui.Text("Toggle disabled")),
			shadcn.Input(shadcn.InputProps{Value: value, OnChange: setValue, Placeholder: "Entry", Disabled: disabled}),
			shadcn.Textarea(shadcn.TextareaProps{Value: value, OnChange: setValue, Placeholder: "Notes", Disabled: disabled}))
	}, struct{}{}), 640, 480)
	entry := func() *ui.Query { return h.Root().ByPlaceholder("Entry") }
	notes := func() *ui.Query { return h.Root().ByPlaceholder("Notes") }
	bounds(entry())
	bounds(notes())
	entry().Focus().Type("!")
	require(entry().Value() == "initial!", "enabled input did not update")
	bounds(entry())
	bounds(notes())
	require(h.Root().ByText("Toggle disabled").Click(), "disable action did not fire")
	require(!entry().Focusable() && !notes().Focusable() && !h.Focused().Exists(), "disabled fields retained focus")
	require(!entry().Type("x") && !notes().Type("x"), "disabled fields accepted editing")
	bounds(entry())
	bounds(notes())
	require(h.Root().ByText("Toggle disabled").Click(), "reenable action did not fire")
	require(entry().Focusable() && notes().Focusable(), "fields did not reenable")
	bounds(entry())
	bounds(notes())
}

func verifyFloatingPanel() {
	h := ui.Mount(ui.Div(ui.Style(ui.Fill),
		ui.Div(ui.Style(ui.Absolute, ui.Right(4), ui.Bottom(4)),
			shadcn.Popover(ui.Text("Open panel"), ui.Text("Floating content")))), 400, 240)
	bounds(h.Root().ByText("Open panel"))
	require(h.Root().ByText("Open panel").Click(), "popover did not open")
	check := func(w, height float32) {
		require(len(h.Overlays()) == 1, "popover overlay is missing")
		b := bounds(h.Overlays()[0].ByKind("scroll"))
		require(b.X >= 0 && b.Y >= 0 && b.X+b.W <= w+0.1 && b.Y+b.H <= height+0.1, "popover escaped viewport")
		bounds(h.Overlays()[0].ByText("Floating content"))
	}
	check(400, 240)
	h.Resize(220, 160)
	check(220, 160)
	h.Escape()
	require(len(h.Overlays()) == 0, "popover did not close")
}

func verifyKeyboardScroll() {
	h := ui.Mount(ui.Div(ui.Style(ui.Fill), ui.Portal(
		ui.ScrollView(ui.Style(ui.Width(200), ui.Height(90)),
			ui.Button(ui.Style(ui.Height(30)), ui.OnClick(func() {}), ui.Text("First")),
			ui.Div(ui.Style(ui.Height(200))),
			ui.Button(ui.Style(ui.Height(30)), ui.OnClick(func() {}), ui.Text("Last"))))), 400, 240)
	require(h.Tab().AllText() == "First", "first keyboard target is missing")
	bounds(h.Focused())
	require(h.Tab().AllText() == "Last", "second keyboard target is missing")
	b := bounds(h.Focused())
	v := bounds(h.Overlays()[0].ByKind("scroll"))
	require(b.Y >= v.Y && b.Y+b.H <= v.Y+v.H+0.1, "keyboard focus was not scrolled into view")
}

func verifyCompactSyntax() {
	saved := ""
	h := ui.Mount(ui.View(func() *ui.Node {
		name, setName := ui.UseState("")
		return ui.VStack(12,
			shadcn.TextField(name, setName).Placeholder("Compact entry"),
			shadcn.TextButton("Save compact", func() { saved = name }),
			ui.For([]string{"A", "B"}, func(label string) *ui.Node { return ui.Text(label).Key(label) }),
		).Width(320).Padding(16)
	}), 640, 480)
	bounds(h.Root())
	bounds(h.Root().ByPlaceholder("Compact entry"))
	h.Root().ByPlaceholder("Compact entry").Focus().Type("compact")
	require(h.Root().ByText("Save compact").Click() && saved == "compact", "compact form did not preserve controlled state")
	bounds(h.Root().ByPlaceholder("Compact entry"))
	bounds(h.Root().ByText("A"))
	bounds(h.Root().ByText("B"))
}

func main() {
	verifyForms()
	verifyFloatingPanel()
	verifyKeyboardScroll()
	verifyCompactSyntax()
	verifyStructSyntax()
	verifyValueBinding()
	fmt.Printf("UI smoke check passed (%d bounds checks)\n", boundsChecks)
}

func verifyValueBinding() {
	name := value.New("Alice")
	h := ui.MountDefault(ui.VStack(8,
		ui.Text(value.Static("Bound profile")),
		shadcn.TextField(name).Placeholder("Value entry"),
		ui.Text(name),
		ui.Text(value.Map(name, func(s string) string { return "Hello " + s })),
	))
	bounds(h.Root().ByPlaceholder("Value entry"))
	bounds(h.Root().ByText("Alice"))
	name.Set("Bob")
	h.Step(0)
	require(h.Root().ByPlaceholder("Value entry").Value() == "Bob" && h.Root().ByText("Hello Bob").Exists(), "Set did not synchronize Value consumers")
	bounds(h.Root().ByPlaceholder("Value entry"))
	bounds(h.Root().ByText("Bob"))
	h.Root().ByPlaceholder("Value entry").SetValue("Carol")
	require(name.Get() == "Carol" && h.Root().ByText("Carol").Exists(), "Value input did not write back to Bind")
	bounds(h.Root().ByText("Hello Carol"))
}

func verifyStructSyntax() {
	var saved string
	h := d.Mount(func() d.View {
		name, setName := ui.UseState("")
		disabled, setDisabled := ui.UseState(false)
		checked, setChecked := ui.UseState(false)
		on, setOn := ui.UseState(false)
		notes, setNotes := ui.UseState("")
		return d.Column{Width: 320, Gap: 8, Padding: 16, Children: []d.View{
			d.Text{Value: "Struct form", FontSize: 20},
			d.Input{Value: name, OnChange: setName, Placeholder: "Struct entry", Disabled: disabled},
			d.Row{Gap: 8, Children: []d.View{
				d.Checkbox{Checked: checked, OnChange: setChecked},
				d.Switch{Checked: on, OnChange: setOn},
			}},
			d.Textarea{Value: notes, OnChange: setNotes, Placeholder: "Struct notes"},
			d.Button{Text: "Save struct", OnClick: func() { saved = name }, Disabled: disabled},
			d.Button{Text: "Toggle struct", OnClick: func() { setDisabled(!disabled) }},
		}}
	}, 640, 480)
	bounds(h.Root())
	bounds(h.Root().ByPlaceholder("Struct entry"))
	h.Root().ByPlaceholder("Struct entry").Focus().Type("Alice")
	require(h.Root().ByText("Save struct").Click() && saved == "Alice", "struct form did not update")
	h.Root().Child(2).Child(0).Click()
	h.Root().Child(2).Child(1).Click()
	h.Root().ByPlaceholder("Struct notes").Focus().Type("Notes")
	h.Step(160)
	for _, q := range []*ui.Query{h.Root().Child(2), h.Root().Child(2).Child(0), h.Root().Child(2).Child(1), h.Root().ByPlaceholder("Struct notes")} {
		bounds(q)
	}
	h.Root().ByText("Toggle struct").Click()
	require(!h.Root().ByPlaceholder("Struct entry").Type("x"), "disabled struct input accepted editing")
	bounds(h.Root().ByPlaceholder("Struct entry"))
	h.Root().ByText("Toggle struct").Click()
	bounds(h.Root().ByPlaceholder("Struct entry"))
	bounds(h.Root().ByText("Save struct"))
}
