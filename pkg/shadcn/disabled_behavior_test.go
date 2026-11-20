package shadcn

import (
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func TestDisabledInputCannotFocusOrEdit(t *testing.T) {
	changed := 0
	h := ui.MountDefault(ui.Div(ui.Style(ui.Column),
		Input(InputProps{Value: "locked", Placeholder: "disabled", Disabled: true,
			OnChange: func(string) { changed++ }}),
		Input(InputProps{Placeholder: "enabled"})))
	q := h.Root().ByPlaceholder("disabled")
	b := q.Bounds()
	if h.MouseDown(b.X+b.W/2, b.Y+b.H/2).Exists() || q.Focusable() {
		t.Fatal("disabled input can receive mouse focus")
	}
	if got := h.Tab().Placeholder(); got != "enabled" {
		t.Fatalf("Tab focused %q, want enabled", got)
	}
	q.Focus()
	if q.IsFocused() || q.Type("x") || q.Backspace(1) || q.Clear() {
		t.Fatal("disabled input accepted focus or editing")
	}
	if q.Value() != "locked" || changed != 0 {
		t.Fatalf("disabled input changed: value=%q callbacks=%d", q.Value(), changed)
	}
}

func TestDisabledControlsDoNotActivateClickableParent(t *testing.T) {
	for _, tc := range []struct {
		name    string
		control func() *ui.Node
	}{
		{"button", func() *ui.Node { return Button(ButtonProps{Disabled: true}, ui.Text("Disabled")) }},
		{"checkbox", func() *ui.Node { return Checkbox(CheckboxProps{Disabled: true}) }},
		{"switch", func() *ui.Node { return Switch(SwitchProps{Disabled: true}) }},
		{"input", func() *ui.Node { return Input(InputProps{Disabled: true, Value: "locked"}) }},
		{"textarea", func() *ui.Node { return Textarea(TextareaProps{Disabled: true, Value: "locked"}) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			clicked := 0
			h := ui.MountDefault(ui.Div(ui.OnClick(func() { clicked++ }), tc.control()))
			q := h.Root().Child(0)
			if q.Click() || clicked != 0 {
				t.Fatal("clicking disabled control activated its parent")
			}
		})
	}
}
