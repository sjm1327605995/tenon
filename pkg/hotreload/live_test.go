package hotreload

import (
	"os"
	"path/filepath"
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 解释示例 view 并挂载，确认（含 shadcn 组件）真的渲染出来。
func TestInterpretViewRenders(t *testing.T) {
	node, err := Interpret(filepath.Join("testdata", "view.go"))
	if err != nil {
		t.Fatalf("interpret: %v", err)
	}
	h := ui.MountDefault(node)
	if !h.Root().ByText("Hot Reload 🔥").Exists() {
		t.Fatalf("interpreted view did not render; texts=%v", h.Root().Texts())
	}
	if !h.Root().ByText("Primary").Exists() {
		t.Fatal("compiled shadcn.Button not rendered from interpreted source")
	}
}

// 解释出错（这里故意用泛型触发 yaegi panic）应被 recover 成 error，而不是让宿主崩溃。
func TestInterpretErrorRecovered(t *testing.T) {
	f := filepath.Join(t.TempDir(), "bad.go")
	src := "package view\nfunc Map[T any](x T) T { return x }\nfunc View() *struct{} { return nil }\n"
	if err := os.WriteFile(f, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Interpret(f); err == nil {
		t.Fatal("expected error (generics), got nil — panic not recovered?")
	}
}

func TestInterpretedDisabledInputCannotBeEdited(t *testing.T) {
	f := filepath.Join(t.TempDir(), "disabled.go")
	src := `package view
import ui "github.com/sjm1327605995/tenon/pkg/ui"
func View() *ui.Node {
	return ui.Input(ui.Disabled(true), ui.Value("locked"))
}
`
	if err := os.WriteFile(f, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	node, err := Interpret(f)
	if err != nil {
		t.Fatalf("interpret Disabled attribute: %v", err)
	}
	h := ui.MountDefault(node)
	q := h.Root().ByKind("input")
	if !q.Exists() || q.Focusable() || q.Type("x") || q.Value() != "locked" {
		t.Fatal("interpreted Disabled input remained editable")
	}
}

func TestInterpretedCompactSyntaxRenders(t *testing.T) {
	f := filepath.Join(t.TempDir(), "compact.go")
	src := `package view
import (
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/shadcn"
)
func View() *ui.Node {
	return ui.VStack(12,
		shadcn.TextField("Alice", nil).Placeholder("Name"),
		shadcn.TextButton("Save", nil),
	).Width(320).Padding(16)
}
`
	if err := os.WriteFile(f, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	node, err := Interpret(f)
	if err != nil {
		t.Fatalf("interpret compact syntax: %v", err)
	}
	h := ui.MountDefault(node)
	if h.Root().Bounds().W != 320 || h.Root().ByPlaceholder("Name").Value() != "Alice" || !h.Root().ByText("Save").Exists() {
		t.Fatal("interpreted compact form did not render correctly")
	}
}
