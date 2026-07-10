package shadcn

import (
	"math"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// ---- Separator ----

type SeparatorProps struct{ Vertical bool }

func Separator(p SeparatorProps) *ui.Node { return ui.Use(separator, p) }

func separator(p SeparatorProps) *ui.Node {
	th := ui.UseTheme()
	if p.Vertical {
		return ui.Div(ui.Style(ui.Width(1), ui.Bg(th.Border)))
	}
	return ui.Div(ui.Style(ui.Height(1), ui.Bg(th.Border)))
}

// ---- Skeleton ----

type SkeletonProps struct{ Width, Height float32 }

func Skeleton(width, height float32) *ui.Node {
	return ui.Use(skeleton, SkeletonProps{Width: width, Height: height})
}

func skeleton(p SkeletonProps) *ui.Node {
	th := ui.UseTheme()
	t := ui.UseElapsed()
	o := 0.55 + 0.35*float32(0.5+0.5*math.Sin(float64(t)*3)) // 呼吸式脉冲
	return ui.Div(ui.Style(ui.Width(p.Width), ui.Height(p.Height), ui.Radius(6),
		ui.Bg(th.Muted), ui.Opacity(o)))
}

// ---- Avatar ----

type avatarProps struct {
	initials string
	size     float32
}

func Avatar(initials string, size float32) *ui.Node {
	return ui.Use(avatar, avatarProps{initials: initials, size: size})
}

func avatar(p avatarProps) *ui.Node {
	th := ui.UseTheme()
	return ui.Div(
		ui.Style(ui.Width(p.size), ui.Height(p.size), ui.Radius(p.size/2), ui.Bg(th.Muted),
			ui.ItemsCenter, ui.JustifyCenter),
		ui.Text(p.initials, ui.FontSize(p.size*0.4), ui.TextColor(th.MutedForeground)),
	)
}
