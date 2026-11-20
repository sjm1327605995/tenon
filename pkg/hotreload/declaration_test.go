package hotreload

import (
	"os"
	"path/filepath"
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func TestInterpretedDeclarationReturnsWithoutExplicitConversion(t *testing.T) {
	f := filepath.Join(t.TempDir(), "struct.go")
	src := `package view
import d "github.com/sjm1327605995/tenon/pkg/decl"
func View() d.View {
	return d.Column{
		Width: 320, Gap: 12, Padding: 16,
		Children: []d.View{
			d.Input{Value: "Alice", Placeholder: "Name", Disabled: true},
			d.Button{Text: "Save"},
		},
	}
}
`
	if err := os.WriteFile(f, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	node, err := Interpret(f)
	if err != nil {
		t.Fatalf("interpret struct declaration: %v", err)
	}
	h := ui.MountDefault(node)
	q := h.Root().ByPlaceholder("Name")
	if h.Root().Bounds().W != 320 || q.Value() != "Alice" || q.Focusable() || q.Type("x") || !h.Root().ByText("Save").Exists() {
		t.Fatal("interpreted struct declaration lost layout, content, or disabled behavior")
	}
}
