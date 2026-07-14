package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

type InputOTPProps struct {
	Length   int // 位数（默认 6）
	Value    string
	OnChange func(string)
}

// InputOTP 是分段一次性验证码输入：N 个格子显示数字，透明输入层捕获键入，
// 当前格高亮。仅接受数字，超长截断。
func InputOTP(p InputOTPProps) *ui.Node { return ui.Use(inputOTP, p) }

func inputOTP(p InputOTPProps) *ui.Node {
	th := ui.UseTheme()
	n := p.Length
	if n <= 0 {
		n = 6
	}
	val := digitsOnly(p.Value)
	if len(val) > n {
		val = val[:n]
	}
	active := len(val) // 下一个待输入的格

	slots := []*ui.Node{ui.Style(ui.Row, ui.Gap(6))}
	for i := 0; i < n; i++ {
		ch := ""
		if i < len(val) {
			ch = string(val[i])
		}
		st := []ui.StyleOpt{ui.Width(40), ui.Height(46), ui.ItemsCenter, ui.JustifyCenter,
			ui.Radius(radiusMd(th)), ui.Bg(th.Background), ui.Border(1, th.Input)}
		if i == active && active < n {
			st = append(st, ui.Border(2, th.Ring)) // 当前格高亮
		}
		slots = append(slots, ui.Div(ui.Style(st...),
			ui.Text(ch, ui.FontSize(18), ui.Medium, ui.TextColor(th.Foreground))))
	}

	// 透明输入覆盖层：捕获键入、点击聚焦，文字/光标全透明（显示交给格子）
	overlay := ui.Input(
		ui.Style(ui.Absolute, ui.Left(0), ui.Top(0), ui.WidthPct(100), ui.HeightPct(100),
			ui.Bg(ui.Color{}), ui.TextColor(ui.Color{})),
		ui.Value(val),
		ui.OnChange(func(v string) {
			out := digitsOnly(v)
			if len(out) > n {
				out = out[:n]
			}
			if p.OnChange != nil {
				p.OnChange(out)
			}
		}),
	)
	return ui.Div(ui.Style(ui.Column), ui.Div(slots...), overlay)
}

func digitsOnly(s string) string {
	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			b = append(b, s[i])
		}
	}
	return string(b)
}
