package main

import (
	"testing"

	"github.com/sjm1327605995/tenon/example/uiparity/shared"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func TestSharedButtonAndInputInteraction(t *testing.T) {
	state := shared.NewState()
	h := ui.Mount(shared.Root(state), 480, 320)
	defer h.Close()
	if !h.Root().ByText("保存").Click() || state.Status.Get() != "按钮已点击" {
		t.Fatalf("button status = %q, want 按钮已点击", state.Status.Get())
	}
	h.Root().ByPlaceholder("输入文本").SetValue("你好")
	if got := state.Text.Get(); got != "你好" {
		t.Fatalf("input binding = %q, want 你好", got)
	}
}
