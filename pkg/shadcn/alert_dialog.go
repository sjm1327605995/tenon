package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type AlertDialogProps struct {
	Open        bool
	Title       string
	Description string
	CancelLabel string // 默认「取消」
	ActionLabel string // 默认「确定」
	OnCancel    func()
	OnAction    func()
	Destructive bool // Action 按钮用危险色
}

// AlertDialog 是需要用户明确确认的模态对话框：标题 + 说明 + 取消/确认按钮。
// 复用 Dialog（Esc / 点遮罩 = 取消）。
func AlertDialog(p AlertDialogProps) *ui.Node {
	cancel := p.CancelLabel
	if cancel == "" {
		cancel = "取消"
	}
	action := p.ActionLabel
	if action == "" {
		action = "确定"
	}
	av := Default
	if p.Destructive {
		av = Destructive
	}
	return Dialog(DialogProps{Open: p.Open, OnClose: p.OnCancel},
		DialogTitle(p.Title),
		DialogDescription(p.Description),
		ui.Div(ui.Style(ui.Row, ui.JustifyEnd, ui.Gap(8)),
			Button(ButtonProps{Variant: Outline, OnClick: p.OnCancel}, ui.Text(cancel)),
			Button(ButtonProps{Variant: av, OnClick: p.OnAction}, ui.Text(action)),
		),
	)
}
