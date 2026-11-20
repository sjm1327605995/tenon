package ui

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFluentAndClassicProduceSamePaint(t *testing.T) {
	classic := VStack(12, Style(Width(320), Padding(16), Bg(White)),
		Input(Value("Alice"), Placeholder("Name"), OnChange(func(string) {})),
		Button(OnClick(func() {}), Text("Save")))
	compact := VStack(12,
		Input().Bind("Alice", func(string) {}).Placeholder("Name"),
		Button(Text("Save")).OnClick(func() {}),
	).Width(320).Padding(16).Bg(White)
	want := MountDefault(classic).Paint()
	got := MountDefault(compact).Paint()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("compact syntax changed the rendered output:\ngot=%+v\nwant=%+v", got, want)
	}
}

func TestFluentBranchesDoNotModifySharedNode(t *testing.T) {
	clicked := ""
	base := Button(Style(Width(80), Height(30)), Text("Button"), OnClick(func() { clicked = "base" }))
	a := base.Width(120).OnClick(func() { clicked = "a" }).Key("a")
	b := base.Width(160).OnClick(func() { clicked = "b" }).Key("b")
	h := MountDefault(VStack(0, base, a, b))
	for i, want := range []float32{80, 120, 160} {
		q := h.Root().Child(i)
		if q.Bounds().W != want || q.Count() != 1 {
			t.Fatalf("branch %d changed layout or inserted a wrapper: %+v children=%d", i, q.Bounds(), q.Count())
		}
		q.Click()
		if clicked != []string{"base", "a", "b"}[i] {
			t.Fatalf("branch %d changed another branch's callback: %q", i, clicked)
		}
	}
	if base.key != "" {
		t.Fatal("assigning branch keys modified the base node")
	}
}

func TestFluentTextAndIconStylesAndEvents(t *testing.T) {
	clicked := 0
	h := MountDefault(HStack(8,
		Text("Label", FontSize(12)).FontSize(24).TextColor(Red).OnClick(func() { clicked++ }),
		Icon(IconCheck, 16).TextColor(Blue)))
	if !h.Root().ByText("Label").Click() || clicked != 1 {
		t.Fatal("fluent text click did not activate")
	}
	text, icon := false, false
	for _, op := range h.Paint() {
		text = text || op.Kind == "text" && op.Text == "Label" && op.Color == Red
		icon = icon || op.Kind == "strokepath" && op.Color == Blue
	}
	if !text || !icon || h.Root().ByText("Label").Bounds().H < 30 {
		t.Fatal("fluent text/icon styles did not reach the renderer")
	}
}

func TestFluentComponentKeepsStateWhenModifiersChangeOrDisappear(t *testing.T) {
	child := func(_ struct{}) *Node {
		count, set := UseState(0)
		return Button(Style(Width(80), Height(30)), Text(fmt.Sprint(count)), OnClick(func() { set(count + 1) }))
	}
	h := MountDefault(View(func() *Node {
		wide, setWide := UseState(false)
		n := Memo(child, struct{}{})
		if wide {
			n = n.Width(180).Bg(Red)
		}
		return VStack(0, n, Button(Text("Toggle"), OnClick(func() { setWide(!wide) })))
	}))
	h.Root().ByText("0").Click()
	h.Root().ByText("Toggle").Click()
	q := h.Root().Child(0)
	if !q.ByText("1").Exists() || q.Bounds().W != 180 || q.Count() != 1 {
		t.Fatalf("component lost state or gained a host wrapper: text=%v bounds=%+v", q.Texts(), q.Bounds())
	}
	q.ByText("1").Click()
	h.Root().ByText("Toggle").Click()
	if !h.Root().Child(0).ByText("2").Exists() || h.Root().Child(0).Bounds().W != 80 {
		t.Fatal("removing modifiers did not restore style while preserving component state")
	}
}

func TestStableModifiedMemoStillBailsOut(t *testing.T) {
	renders := 0
	child := func(_ struct{}) *Node { renders++; return Text("Stable") }
	stable := Memo(child, struct{}{}).TextColor(Red)
	h := MountDefault(View(func() *Node {
		tick, set := UseState(0)
		return VStack(0, stable, Button(Text("Tick"), OnClick(func() { set(tick + 1) })))
	}))
	h.Root().ByText("Tick").Click()
	if renders != 1 {
		t.Fatalf("stable modified Memo rendered %d times, want 1", renders)
	}
}

func TestFluentComponentUsesLatestCallbackAndBinding(t *testing.T) {
	observed := -1
	field := func(_ struct{}) *Node { return Input() }
	button := func(_ struct{}) *Node { return Button(Text("Record")) }
	h := MountDefault(View(func() *Node {
		tick, setTick := UseState(0)
		value, setValue := UseState("initial")
		return VStack(0,
			Use(field, struct{}{}).Bind(value, setValue).Placeholder("Name"),
			Memo(button, struct{}{}).OnClick(func() { observed = tick }),
			Button(Text("Tick"), OnClick(func() { setTick(tick + 1) })), Text(value))
	}))
	h.Root().ByPlaceholder("Name").Focus().Type("!")
	if !h.Root().ByText("initial!").Exists() {
		t.Fatal("fluent component binding did not update state")
	}
	h.Root().ByText("Tick").Click()
	h.Root().ByText("Record").Click()
	if observed != 1 {
		t.Fatalf("fluent component retained stale callback: %d", observed)
	}
}

func TestFluentSingleRootProviderPreservesContext(t *testing.T) {
	ctx := CreateContext("fallback")
	child := func(_ struct{}) *Node { return Text(UseContext(ctx)) }
	h := MountDefault(ctx.Provider("provided", View(func() *Node { return Use(child, struct{}{}) })).Width(200))
	if !h.Root().ByText("provided").Exists() || h.Root().Bounds().W != 200 {
		t.Fatal("styling through a provider lost context or layout")
	}
}

func TestRemovingLeafModifiersClearsOldStyleAndEvents(t *testing.T) {
	clicks := 0
	h := MountDefault(View(func() *Node {
		modified, set := UseState(true)
		label := Text("Label")
		if modified {
			label = label.Width(180).OnClick(func() { clicks++ })
		}
		return VStack(0, label, Button(Text("Remove")).OnClick(func() { set(false) }))
	}))
	h.Root().ByText("Label").Click()
	h.Root().ByText("Remove").Click()
	q := h.Root().ByText("Label")
	if q.Click() || q.Clickable() || q.Bounds().W == 180 || clicks != 1 {
		t.Fatal("removing leaf modifiers retained old style or click handler")
	}
}

func TestFluentNilAndInvalidRoots(t *testing.T) {
	var empty *Node
	if empty.Width(20).OnClick(func() {}).Key("nil") != nil {
		t.Fatal("fluent methods should preserve optional nil nodes")
	}
	for _, n := range []*Node{Fragment(Text("A"), Text("B")), Div()} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("invalid modifier usage did not fail")
				}
			}()
			if n.typ == typeFragment {
				n.Width(10)
			} else {
				n.With(Text("not an attribute"))
			}
		}()
	}
}
