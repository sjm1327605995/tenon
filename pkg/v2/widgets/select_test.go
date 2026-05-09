package widgets

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

func TestSelectOpenFindsOptions(t *testing.T) {
	selectWidget := Select([]SelectOption{
		{Value: "a", Label: "Option A"},
		{Value: "b", Label: "Option B"},
		{Value: "c", Label: "Option C"},
	})

	env := ui.TestWidget(t, selectWidget, 400, 300)

	// 初始状态：应该能找到 placeholder 文本
	if rt := env.FindText("请选择..."); rt == nil {
		t.Fatal("expected placeholder text before opening")
	}

	// 点击 Select trigger（在中心位置）
	env.TapAt(100, 20)
	env.Rebuild()

	// 打开后，应该能找到 dropdown 的选项文本
	for _, label := range []string{"Option A", "Option B", "Option C"} {
		if rt := env.FindText(label); rt == nil {
			t.Fatalf("expected to find '%s' after opening select", label)
		}
	}
}

func TestSelectOpenClose(t *testing.T) {
	selectWidget := Select([]SelectOption{
		{Value: "a", Label: "Option A"},
	})

	env := ui.TestWidget(t, selectWidget, 400, 300)

	// 打开
	env.TapAt(100, 20)
	env.Rebuild()
	if rt := env.FindText("Option A"); rt == nil {
		t.Fatal("expected option after opening")
	}

	// 关闭（再次点击 trigger）
	env.TapAt(100, 20)
	env.Rebuild()
	if rt := env.FindText("Option A"); rt != nil {
		t.Fatal("expected option to disappear after closing")
	}
}
