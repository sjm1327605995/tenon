package hotreload

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

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

func TestLivePreviewReloadsAndRecovers(t *testing.T) {
	file := filepath.Join(t.TempDir(), "view.go")
	var h *ui.Harness
	writeText := func(text string) {
		t.Helper()
		src := fmt.Sprintf("package view\nimport ui \"github.com/sjm1327605995/tenon/pkg/ui\"\nfunc View() *ui.Node { return ui.Text(%q) }\n", text)
		if err := os.WriteFile(file, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	waitFor := func(text string) {
		t.Helper()
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if h := h; h != nil && h.Root().ByText(text).Exists() {
				return
			}
			time.Sleep(50 * time.Millisecond)
			if h != nil {
				h.Step(0)
			}
		}
		t.Fatalf("text %q did not appear; texts=%v", text, h.Root().Texts())
	}

	writeText("Before")
	h = ui.MountDefault(LivePreview(file))
	defer h.Close()
	if !h.Root().ByText("Before").Exists() {
		t.Fatal("initial live preview did not render")
	}
	time.Sleep(20 * time.Millisecond)
	writeText("After")
	waitFor("After")

	if err := os.WriteFile(file, []byte("package view\nfunc Map[T any](x T) T { return x }\nfunc View() int { return 0 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	waitFor("编译/解释错误")
	if !h.Root().ByText("After").Exists() {
		t.Fatal("interpret error removed the last good view")
	}
	time.Sleep(20 * time.Millisecond)
	writeText("Recovered")
	waitFor("Recovered")
	if h.Root().ByText("编译/解释错误").Exists() {
		t.Fatal("successful reload retained the error overlay")
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

func TestInterpretedRecentUIAPIs(t *testing.T) {
	f := filepath.Join(t.TempDir(), "recent.go")
	src := `package view
import (
	"image"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)
type sink struct{ calls float32 }
func (s *sink) MoveTo(float32, float32) { s.calls++ }
func (s *sink) LineTo(float32, float32) { s.calls++ }
func (s *sink) CubeTo(float32, float32, float32, float32, float32, float32) { s.calls++ }
func (s *sink) QuadTo(float32, float32, float32, float32) { s.calls++ }
func (s *sink) ClosePath() {}
func View() *ui.Node {
	s := &sink{}
	ui.ParseSVGPath(ui.IconCheck, float32(1), s)
	if s.calls == 0 {
		return ui.Text("missing")
	}
	warped, _ := ui.WarpImageToQuad(image.NewNRGBA(image.Rect(0, 0, 1, 1)),
		[4]ui.Point{{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 0, Y: 10}, {X: 10, Y: 10}})
	if warped == nil {
		return ui.Text("missing")
	}
	return ui.VStack(8,
		ui.Text("recent"),
		ui.Icon(ui.IconCheck, 16),
		ui.Img(ui.SrcImage("hot-recent", warped),
			ui.Style(ui.ImageBlur(1), ui.ImageGrayscale(1))),
	)
}
`
	if err := os.WriteFile(f, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	node, err := Interpret(f)
	if err != nil {
		t.Fatalf("interpret recent APIs: %v", err)
	}
	h := ui.MountDefault(node)
	defer h.Close()
	if !h.Root().ByText("recent").Exists() || h.Root().ByText("missing").Exists() {
		t.Fatalf("ParseSVGPath was not callable: texts=%v", h.Root().Texts())
	}
	if !h.Root().ByKind("image").Exists() {
		t.Fatal("recent filtered image did not render")
	}
	foundVector := false
	for _, op := range h.Paint() {
		if op.Kind == "strokepath" || op.Kind == "fillpath" {
			foundVector = true
		}
	}
	if !foundVector {
		t.Fatal("recent SVG icon did not render")
	}
}
