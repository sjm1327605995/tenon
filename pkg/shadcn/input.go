package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type InputProps struct {
	Value       string
	OnChange    func(string)
	Placeholder string
	Disabled    bool
}

// Input 是主题化的单行输入框。
func Input(p InputProps) *ui.Node { return ui.Use(input, p) }

func input(p InputProps) *ui.Node {
	th := ui.UseTheme()
	// shadcn v4: h-9 rounded-md border border-input px-3 py-1 text-sm shadow-xs。
	st := []ui.StyleOpt{
		ui.Height(36), ui.PaddingXY(12, 0), ui.Radius(radiusMd(th)), ui.ItemsCenter,
		ui.Bg(th.Background), ui.Border(1, th.Input), ui.TextColor(th.Foreground),
		ui.FontSize(14), shadowXs(),
	}
	if p.Disabled {
		st = append(st, ui.Opacity(0.5))
	}
	attrs := []*ui.Node{ui.Style(st...), ui.Value(p.Value), ui.Placeholder(p.Placeholder)}
	if !p.Disabled {
		attrs = append(attrs, ui.OnChange(p.OnChange))
	}
	return ui.Input(attrs...)
}

// ---- Textarea ----

type TextareaProps struct {
	Value       string
	OnChange    func(string)
	Placeholder string
	Rows        int
}

// Textarea 是多行文本域：自动折行、Enter 换行、高度随内容增长。
func Textarea(p TextareaProps) *ui.Node { return ui.Use(textarea, p) }

func textarea(p TextareaProps) *ui.Node {
	th := ui.UseTheme()
	rows := p.Rows
	if rows <= 0 {
		rows = 3
	}
	// shadcn v4: min-h-16 rounded-md border px-3 py-2 text-sm shadow-xs。
	minH := float32(rows)*22 + 16
	if minH < 64 {
		minH = 64
	}
	st := []ui.StyleOpt{
		ui.Width(280), ui.MinHeight(minH), ui.PaddingXY(12, 8),
		ui.Radius(radiusMd(th)), ui.Bg(th.Background), ui.Border(1, th.Input),
		ui.TextColor(th.Foreground), ui.FontSize(14), shadowXs(),
	}
	return ui.Input(ui.Style(st...), ui.Multiline(),
		ui.Value(p.Value), ui.Placeholder(p.Placeholder), ui.OnChange(p.OnChange))
}

// Label 是表单标签。
func Label(text string) *ui.Node { return ui.Use(labelC, text) }

func labelC(text string) *ui.Node {
	th := ui.UseTheme()
	return ui.Text(text, ui.FontSize(14), ui.TextColor(th.Foreground))
}
