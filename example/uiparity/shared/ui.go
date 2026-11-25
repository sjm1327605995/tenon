package shared

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	"github.com/sjm1327605995/tenon/pkg/value"
)

type State struct {
	Text   value.Bind[string]
	Status value.Bind[string]
}

func NewState() *State {
	return &State{Text: value.New(""), Status: value.New("等待交互")}
}

func Root(state *State) *ui.Node {
	if state == nil {
		state = NewState()
	}
	return ui.VStack(18,
		ui.Text("G3N 控件一致性", ui.FontSize(24), ui.Bold),
		ui.HStack(12,
			ui.Text("Regular", ui.FontSize(15)),
			ui.Text("Semibold", ui.FontSize(15), ui.Semibold),
			ui.Text("Italic", ui.FontSize(15), ui.Italic),
		),
		ui.HStack(12,
			ui.Icon(ui.IconCheck, 24, ui.TextColor(ui.Hex("#16a34a"))),
			ui.Icon(ui.IconChevronDown, 24, ui.TextColor(ui.Hex("#475569"))),
			ui.Icon(ui.IconSearch, 24, ui.TextColor(ui.Hex("#2563eb"))),
			ui.Icon(ui.IconTrash, 24, ui.TextColor(ui.Hex("#dc2626"))),
		).With(ui.Style(ui.Scale(1.05), ui.Rotate(3))),
		shadcn.Input(shadcn.InputProps{Placeholder: "输入文本", Bind: state.Text}),
		shadcn.Button(shadcn.ButtonProps{OnClick: func() { state.Status.Set("按钮已点击") }},
			ui.Icon(ui.IconCheck, 16), ui.Text("保存")),
		ui.Text(state.Status, ui.FontSize(13), ui.TextColor(ui.Hex("#64748b"))),
	).Padding(28).Bg(ui.Hex("#f8fafc"))
}
