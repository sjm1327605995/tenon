package shadcn

import (
	"fmt"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
	values "github.com/sjm1327605995/tenon/pkg/value"
)

// TextButton 是默认变体的纯文字按钮。需要变体、尺寸、禁用外观或图标时使用 ButtonProps。
func TextButton(text any, onClick func()) *ui.Node {
	return Button(ButtonProps{OnClick: onClick}, ui.Text(text))
}

// TextField accepts a Bind[string] directly, or a legacy string plus callback.
// With a Bind, UI edits and programmatic Set calls synchronize automatically.
func TextField(source any, onChange ...func(string)) *ui.Node {
	if len(onChange) > 1 {
		panic("shadcn.TextField: at most one change callback is allowed")
	}
	var observer func(string)
	if len(onChange) == 1 {
		observer = onChange[0]
	}
	switch value := source.(type) {
	case values.Bind[string]:
		return Input(InputProps{Bind: value, OnChange: observer})
	case string:
		return Input(InputProps{Value: value, OnChange: observer})
	default:
		panic(fmt.Sprintf("shadcn.TextField: expected string or value.Bind[string], got %T", source))
	}
}
