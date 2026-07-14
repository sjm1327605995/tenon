package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type FieldProps struct {
	Label       string
	Description string // 辅助说明（次要色）
	Error       string // 非空时以错误色显示，覆盖 Description
}

// Field 是表单字段容器：标签 + 控件 + 说明/错误。控件（如 shadcn.Input）作为子节点传入。
//
//	shadcn.Field(shadcn.FieldProps{Label: "邮箱", Description: "我们不会公开你的邮箱。"},
//	    shadcn.Input(shadcn.InputProps{Placeholder: "you@example.com"}))
func Field(p FieldProps, control *ui.Node) *ui.Node {
	return ui.Use(field, fieldProps{p: p, control: control})
}

type fieldProps struct {
	p       FieldProps
	control *ui.Node
}

func field(fp fieldProps) *ui.Node {
	th := ui.UseTheme()
	p := fp.p
	kids := []*ui.Node{ui.Style(ui.Column, ui.Gap(6))}
	if p.Label != "" {
		kids = append(kids, ui.Text(p.Label, ui.FontSize(14), ui.Medium, ui.TextColor(th.Foreground)))
	}
	if fp.control != nil {
		kids = append(kids, fp.control)
	}
	if p.Error != "" {
		kids = append(kids, ui.Text(p.Error, ui.FontSize(13), ui.TextColor(th.Destructive)))
	} else if p.Description != "" {
		kids = append(kids, ui.Text(p.Description, ui.FontSize(13), ui.TextColor(th.MutedForeground)))
	}
	return ui.Div(kids...)
}
