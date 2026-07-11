// Package shadcn is a shadcn/ui-style component library built on top of pkg/ui.
// Components are the base ui primitives restyled with theme tokens (ui.UseTheme)
// and interaction state (ui.UseInteraction); they mix freely with base ui nodes.
package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// Variant 是按钮视觉变体。
type Variant int

const (
	Default Variant = iota
	Destructive
	Outline
	Secondary
	Ghost
	Link
)

// Size 是按钮尺寸。
type Size int

const (
	SizeDefault Size = iota
	SizeSm
	SizeLg
	SizeIcon
)

// ButtonProps 配置按钮外观与行为。
type ButtonProps struct {
	Variant  Variant
	Size     Size
	OnClick  func()
	Disabled bool
	children []*ui.Node
}

// Button 渲染一个 shadcn 风格按钮。children 作为内容（文本/图标）。
func Button(p ButtonProps, children ...*ui.Node) *ui.Node {
	p.children = children
	return ui.Use(button, p)
}

func button(p ButtonProps) *ui.Node {
	th := ui.UseTheme()
	hovered, pressed, ia := ui.UseInteraction()

	// 基础：rounded-md, text-sm, font-medium, gap-2。
	style := []ui.StyleOpt{
		ui.Row, ui.ItemsCenter, ui.JustifyCenter, ui.Gap(8),
		ui.Radius(radiusMd(th)), ui.FontSize(14), ui.Medium,
	}
	// 尺寸（Tailwind: h-9 px-4 / h-8 px-3 / h-10 px-6 / size-9），高度固定，垂直居中。
	switch p.Size {
	case SizeSm:
		style = append(style, ui.Height(32), ui.PaddingXY(12, 0), ui.Gap(6))
	case SizeLg:
		style = append(style, ui.Height(40), ui.PaddingXY(24, 0))
	case SizeIcon:
		style = append(style, ui.Width(36), ui.Height(36))
	default:
		style = append(style, ui.Height(36), ui.PaddingXY(16, 0))
	}

	// 变体配色 + hover（bg-*/90|80 用 alpha-over-背景 近似）。
	active := hovered && !p.Disabled
	bg, fg, border := th.Primary, th.PrimaryForeground, ui.Transparent
	bordered := false
	var shadow ui.StyleOpt
	switch p.Variant {
	case Destructive:
		bg, fg = th.Destructive, ui.White
		if active {
			bg = over(th.Destructive, th.Background, 0.9)
		}
	case Outline:
		bg, fg, border, bordered = th.Background, th.Foreground, th.Border, true
		shadow = shadowXs()
		if active {
			bg, fg = th.Accent, th.AccentForeground
		}
	case Secondary:
		bg, fg = th.Secondary, th.SecondaryForeground
		if active {
			bg = over(th.Secondary, th.Background, 0.8)
		}
	case Ghost:
		bg, fg = ui.Transparent, th.Foreground
		if active {
			bg, fg = th.Accent, th.AccentForeground
		}
	case Link:
		bg, fg = ui.Transparent, th.Primary // hover:underline 暂无下划线支持
	default: // Default
		if active {
			bg = over(th.Primary, th.Background, 0.9)
		}
	}

	style = append(style, ui.Bg(bg), ui.TextColor(fg))
	if bordered {
		style = append(style, ui.Border(1, border))
	}
	if shadow != nil {
		style = append(style, shadow)
	}
	if pressed && !p.Disabled {
		style = append(style, ui.Scale(0.97))
	}
	if p.Disabled {
		style = append(style, ui.Opacity(0.5))
	}

	attrs := []*ui.Node{ui.Style(style...)}
	if !p.Disabled {
		attrs = append(attrs, ui.OnClick(p.OnClick), ia)
	}
	return ui.Button(append(attrs, p.children...)...)
}
