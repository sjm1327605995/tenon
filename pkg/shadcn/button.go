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

	// 基础 + 尺寸
	style := []ui.StyleOpt{
		ui.Row, ui.ItemsCenter, ui.JustifyCenter, ui.Gap(8), ui.Radius(th.Radius),
	}
	switch p.Size {
	case SizeSm:
		style = append(style, ui.PaddingXY(12, 6), ui.FontSize(13), ui.Height(32))
	case SizeLg:
		style = append(style, ui.PaddingXY(24, 10), ui.FontSize(16), ui.Height(44))
	case SizeIcon:
		style = append(style, ui.Width(36), ui.Height(36))
	default:
		style = append(style, ui.PaddingXY(16, 8), ui.FontSize(14), ui.Height(38))
	}

	// 变体配色
	bg, fg, border := th.Primary, th.PrimaryForeground, ui.Transparent
	bordered := false
	switch p.Variant {
	case Destructive:
		bg, fg = th.Destructive, th.DestructiveForeground
	case Outline:
		bg, fg, border, bordered = th.Background, th.Foreground, th.Border, true
	case Secondary:
		bg, fg = th.Secondary, th.SecondaryForeground
	case Ghost:
		bg, fg = ui.Transparent, th.Foreground
	case Link:
		bg, fg = ui.Transparent, th.Primary
	}

	// hover 态
	active := hovered && !p.Disabled
	switch p.Variant {
	case Ghost, Outline:
		if active {
			bg, fg = th.Accent, th.AccentForeground
		}
	case Link:
		// 链接变体保持透明
	default:
		if active {
			bg = ui.Mix(bg, th.Background, 0.12)
		}
	}

	style = append(style, ui.Bg(bg), ui.TextColor(fg))
	if bordered {
		style = append(style, ui.Border(1, border))
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
