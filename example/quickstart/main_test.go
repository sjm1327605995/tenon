package main

import (
	"testing"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

func TestProfileEditingAndClear(t *testing.T) {
	profile := &Profile{Name: value.New("Tenon")}
	h := ui.Mount(ui.Render(profile), 480, 320)
	defer h.Close()

	input := h.Root().ByPlaceholder("姓名")
	if !input.Exists() || input.Value() != "Tenon" {
		t.Fatalf("initial input exists=%t value=%q, want Tenon", input.Exists(), input.Value())
	}
	input.SetValue("Kimi")
	if got := profile.Name.Get(); got != "Kimi" {
		t.Fatalf("bound name=%q, want Kimi", got)
	}
	if !h.Root().ByText("Kimi").Exists() {
		t.Fatal("updated name was not rendered")
	}
	h.Root().ByText("清空").Click()
	if got := profile.Name.Get(); got != "" {
		t.Fatalf("cleared name=%q, want empty", got)
	}
	if got := h.Root().ByPlaceholder("姓名").Value(); got != "" {
		t.Fatalf("cleared input=%q, want empty", got)
	}
}
