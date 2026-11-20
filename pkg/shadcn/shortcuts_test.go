package shadcn

import (
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

func TestShortcutsComposeAndKeepControlledState(t *testing.T) {
	var saved string
	h := ui.MountDefault(ui.View(func() *ui.Node {
		name, setName := ui.UseState("")
		return ui.VStack(12,
			TextField(name, setName).Placeholder("Name").OnSubmit(func(s string) { saved = s }),
			TextButton("Save", func() { saved = name }),
		).Width(320).Padding(24)
	}))
	if h.Root().Bounds().W != 320 || h.Root().Count() != 2 {
		t.Fatal("shortcut composition changed layout structure")
	}
	h.Root().ByPlaceholder("Name").Focus().Type("Alice")
	h.Enter()
	if saved != "Alice" {
		t.Fatalf("shortcut field did not submit controlled value: %q", saved)
	}
	h.Root().ByPlaceholder("Name").SetValue("Bob")
	h.Root().ByText("Save").Click()
	if saved != "Bob" {
		t.Fatalf("shortcut button retained stale state: %q", saved)
	}
}

func TestTextFieldAcceptsBindDirectly(t *testing.T) {
	name := value.New("Alice")
	caption := value.New("Save")
	var observed string
	h := ui.MountDefault(ui.VStack(8,
		TextField(name, func(next string) { observed = next }).Placeholder("Name"),
		Label(name),
		TextButton(caption, func() {})))
	h.Root().ByPlaceholder("Name").Focus().SetValue("Bob")
	if observed != "Bob" || name.Get() != "Bob" || !h.Root().ByText("Bob").Exists() {
		t.Fatal("themed Bind input and label did not synchronize")
	}
	name.Set("Carol")
	caption.Set("Submit")
	h.Step(0)
	if h.Root().ByPlaceholder("Name").Value() != "Carol" || !h.Root().ByText("Submit").Exists() || observed != "Bob" {
		t.Fatal("programmatic Bind update did not synchronize or incorrectly fired the input observer")
	}
}
