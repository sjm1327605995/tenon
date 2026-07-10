package shadcn

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// Shared design tokens translating shadcn/ui (Tailwind v4, new-york) into Tenon
// styles. The Tailwind spacing scale is 1 unit = 4px (px-3 = 12, h-9 = 36, …);
// components below reference these helpers so spacing/radius/shadow stay in one
// place and track the theme.

// ---- Radius steps ----
// shadcn v4 derives radii from --radius (Theme.Radius, default 10px):
//   sm = base-4, md = base-2, lg = base, xl = base+4  (rounded-full ~= 9999).

func radiusSm(th ui.Theme) float32 { return maxf(th.Radius-4, 0) }
func radiusMd(th ui.Theme) float32 { return maxf(th.Radius-2, 0) }
func radiusLg(th ui.Theme) float32 { return th.Radius }
func radiusXl(th ui.Theme) float32 { return th.Radius + 4 }

// radiusFull is an effectively pill/circular radius (rounded-full).
const radiusFull float32 = 9999

func maxf(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// ---- Elevation (box-shadow) ----
// Alpha values approximate Tailwind's shadow opacities (0.05–0.12 black).

func shadowXs() ui.StyleOpt { return ui.Shadow(ui.Color{A: 13}, 0, 1, 2, 0) }
func shadowSm() ui.StyleOpt { return ui.Shadow(ui.Color{A: 25}, 0, 1, 3, 0) }
func shadowMd() ui.StyleOpt { return ui.Shadow(ui.Color{A: 30}, 0, 4, 8, -1) }
func shadowLg() ui.StyleOpt { return ui.Shadow(ui.Color{A: 36}, 0, 10, 18, -3) }

// ---- Color helpers ----

// over returns the solid color equivalent to fg painted at the given alpha over
// bg — the translation of Tailwind's "bg-primary/90" style opacity modifiers.
// e.g. over(th.Primary, th.Background, 0.9) == bg-primary/90.
func over(fg, bg ui.Color, alpha float32) ui.Color { return ui.Mix(bg, fg, alpha) }
