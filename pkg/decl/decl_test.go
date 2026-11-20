package decl_test

import (
	"fmt"
	"reflect"
	"testing"

	d "github.com/sjm1327605995/tenon/pkg/decl"
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func TestDeclarationMatchesExistingLayoutAndPaint(t *testing.T) {
	classic := ui.VStack(12, ui.Style(ui.Width(320), ui.Padding(24), ui.Bg(ui.White)),
		ui.Text("Profile", ui.FontSize(22), ui.Bold),
		shadcn.Input(shadcn.InputProps{Value: "Alice", Placeholder: "Name"}),
		shadcn.Button(shadcn.ButtonProps{}, ui.Text("Save")))
	want := ui.MountDefault(classic).Paint()
	h := d.Mount(func() d.View {
		return d.Column{Width: 320, Gap: 12, Padding: 24, Background: ui.White, Children: []d.View{
			d.Text{Value: "Profile", FontSize: 22, Bold: true},
			d.Input{Value: "Alice", Placeholder: "Name"},
			d.Button{Text: "Save"},
		}}
	}, 800, 600)
	if !reflect.DeepEqual(h.Paint(), want) || h.Root().Count() != 3 {
		t.Fatal("declaration changed rendering or added host wrappers")
	}
}

func TestControlledFormUpdatesAndDisablesWithoutExplicitConversion(t *testing.T) {
	var saved string
	h := d.Mount(func() d.View {
		name, setName := ui.UseState("")
		disabled, setDisabled := ui.UseState(false)
		return d.Column{Width: 320, Gap: 8, Children: []d.View{
			d.Input{Value: name, OnChange: setName, Placeholder: "Name", Disabled: disabled,
				OnSubmit: func(v string) { saved = v }},
			d.Button{Text: "Save", Disabled: disabled, OnClick: func() { saved = name }},
			d.Button{Text: "Toggle", OnClick: func() { setDisabled(!disabled) }},
		}}
	}, 400, 240)
	h.Root().ByPlaceholder("Name").Focus().Type("Alice")
	h.Enter()
	if saved != "Alice" {
		t.Fatalf("declaration input did not submit: %q", saved)
	}
	h.Root().ByText("Toggle").Click()
	if h.Root().ByPlaceholder("Name").Focusable() || h.Root().ByPlaceholder("Name").Type("x") || h.Root().ByText("Save").Click() {
		t.Fatal("disabled declarations remained interactive")
	}
	h.Root().ByText("Toggle").Click()
	h.Root().ByPlaceholder("Name").SetValue("Bob")
	h.Root().ByText("Save").Click()
	if saved != "Bob" || h.Root().ByPlaceholder("Name").Bounds().H != 36 {
		t.Fatal("reenabled declaration lost state, layout, or callback freshness")
	}
}

func firstPage() d.View {
	count, set := ui.UseState(1)
	return d.Button{Text: fmt.Sprint(count), OnClick: func() { set(count + 1) }}
}

func secondPage() d.View {
	count, _ := ui.UseState(99)
	return d.Text{Value: fmt.Sprint(count)}
}

func TestNestedComponentsKeepOriginalFunctionIdentity(t *testing.T) {
	h := d.Mount(func() d.View {
		second, setSecond := ui.UseState(false)
		page := d.Func(firstPage)
		if second {
			page = d.Func(secondPage)
		}
		return d.Column{Children: []d.View{
			page, d.Button{Text: "Switch", OnClick: func() { setSecond(!second) }},
		}}
	}, 400, 240)
	h.Root().ByText("1").Click()
	if !h.Root().ByText("2").Exists() {
		t.Fatal("component local state did not update")
	}
	h.Root().ByText("Switch").Click()
	if !h.Root().ByText("99").Exists() {
		t.Fatal("distinct declaration functions shared component state")
	}
}

func TestTypedPropsAndKeysPreserveStateDuringReorder(t *testing.T) {
	row := func(id string) d.View {
		count, set := ui.UseState(0)
		return d.Button{Text: fmt.Sprintf("%s:%d", id, count), OnClick: func() { set(count + 1) }}
	}
	h := d.Mount(func() d.View {
		reversed, setReversed := ui.UseState(false)
		ids := []string{"a", "b"}
		if reversed {
			ids = []string{"b", "a"}
		}
		rows := make([]d.View, 0, len(ids))
		for _, id := range ids {
			rows = append(rows, d.Use(row, id).Key(id))
		}
		return d.Column{Children: []d.View{
			d.Fragment{Children: rows},
			d.Button{Text: "Reverse", OnClick: func() { setReversed(!reversed) }},
		}}
	}, 400, 240)
	h.Root().ByText("a:0").Click()
	h.Root().ByText("Reverse").Click()
	if got := h.Root().Texts(); !reflect.DeepEqual(got, []string{"b:0", "a:1", "Reverse"}) || h.Root().Count() != 3 {
		t.Fatalf("declaration reorder lost keyed state or introduced wrappers: %v", got)
	}
}

func TestMemoDeclarationKeepsBailout(t *testing.T) {
	renders := 0
	type props struct{ Label string }
	child := func(p props) d.View { renders++; return d.Text{Value: p.Label} }
	h := d.Mount(func() d.View {
		tick, set := ui.UseState(0)
		return d.Column{Children: []d.View{
			d.Memo(child, props{Label: "Stable"}),
			d.Button{Text: "Tick", OnClick: func() { set(tick + 1) }},
		}}
	}, 400, 240)
	h.Root().ByText("Tick").Click()
	if renders != 1 {
		t.Fatalf("unchanged declaration props rendered %d times", renders)
	}
}

func TestMixedChildrenAndTypedNil(t *testing.T) {
	var oldNode *ui.Node
	var optional *d.Text
	h := d.Mount(func() d.View {
		return d.Column{Children: []d.View{nil, oldNode, optional, d.Text{Value: "Struct"}, ui.Text("Node")}}
	}, 400, 240)
	if h.Root().Count() != 2 || !reflect.DeepEqual(h.Root().Texts(), []string{"Struct", "Node"}) {
		t.Fatal("mixed declaration/node children did not ignore nils")
	}
	if d.Mount(func() d.View { return optional }, 400, 240).Root().Exists() {
		t.Fatal("typed nil page should render nothing")
	}
}

func TestStyleOverridesCanExpressExplicitZero(t *testing.T) {
	h := d.Mount(func() d.View {
		return d.Box{Width: 120, Height: 50, Style: []ui.StyleOpt{ui.Width(0)}}
	}, 400, 240)
	if h.Root().Bounds().W != 0 || h.Root().Bounds().H != 50 {
		t.Fatalf("explicit zero override was lost: %+v", h.Root().Bounds())
	}
}

func TestStructScrollUsesKeyboardVisibility(t *testing.T) {
	h := d.Mount(func() d.View {
		return d.ScrollView{Width: 200, Height: 90, Children: []d.View{
			d.Button{Text: "First", OnClick: func() {}},
			d.Box{Height: 200},
			d.Button{Text: "Last", OnClick: func() {}},
		}}
	}, 400, 240)
	h.Tab()
	h.Tab()
	b, v := h.Focused().Bounds(), h.Root().Bounds()
	if h.Focused().AllText() != "Last" || b.Y < v.Y || b.Y+b.H > v.Y+v.H+0.1 {
		t.Fatal("keyboard focus did not reveal scrolled declaration")
	}
}

func TestCheckboxSwitchAndTextareaRemainControlled(t *testing.T) {
	h := d.Mount(func() d.View {
		checked, setChecked := ui.UseState(false)
		on, setOn := ui.UseState(false)
		text, setText := ui.UseState("")
		return d.Column{Width: 300, Gap: 8, Children: []d.View{
			d.Row{Gap: 8, Children: []d.View{
				d.Checkbox{Checked: checked, OnChange: setChecked},
				d.Switch{Checked: on, OnChange: setOn},
			}},
			d.Textarea{Value: text, OnChange: setText, Placeholder: "Notes"},
			d.Text{Value: fmt.Sprintf("%v/%v/%s", checked, on, text)},
		}}
	}, 400, 300)
	h.Root().Child(0).Child(0).Click()
	h.Root().Child(0).Child(1).Click()
	h.Root().ByPlaceholder("Notes").Focus().Type("Hello")
	h.Step(160)
	if !h.Root().ByText("true/true/Hello").Exists() || h.Root().ByPlaceholder("Notes").Bounds().H <= 0 {
		t.Fatal("declaration controls did not update their parent state")
	}
}
