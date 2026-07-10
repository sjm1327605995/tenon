package ui

// Theme 是语义化设计令牌，供各扩展风格库（shadcn 等）按语义取色，随明暗主题切换。
// 命名对齐 shadcn 的 CSS 变量。
type Theme struct {
	Background, Foreground             Color
	Card, CardForeground               Color
	Popover, PopoverForeground         Color
	Primary, PrimaryForeground         Color
	Secondary, SecondaryForeground     Color
	Muted, MutedForeground             Color
	Accent, AccentForeground           Color
	Destructive, DestructiveForeground Color
	Border, Input, Ring                Color
	Radius                             float32
}

// LightTheme 是默认浅色主题（shadcn zinc 风格）。
var LightTheme = Theme{
	Background: Hex("#ffffff"), Foreground: Hex("#0a0a0a"),
	Card: Hex("#ffffff"), CardForeground: Hex("#0a0a0a"),
	Popover: Hex("#ffffff"), PopoverForeground: Hex("#0a0a0a"),
	Primary: Hex("#18181b"), PrimaryForeground: Hex("#fafafa"),
	Secondary: Hex("#f4f4f5"), SecondaryForeground: Hex("#18181b"),
	Muted: Hex("#f4f4f5"), MutedForeground: Hex("#71717a"),
	Accent: Hex("#f4f4f5"), AccentForeground: Hex("#18181b"),
	Destructive: Hex("#ef4444"), DestructiveForeground: Hex("#fafafa"),
	Border: Hex("#e4e4e7"), Input: Hex("#e4e4e7"), Ring: Hex("#a1a1aa"),
	Radius: 10,
}

// DarkTheme 是默认深色主题（shadcn zinc 风格）。
var DarkTheme = Theme{
	Background: Hex("#0a0a0a"), Foreground: Hex("#fafafa"),
	Card: Hex("#0a0a0a"), CardForeground: Hex("#fafafa"),
	Popover: Hex("#0a0a0a"), PopoverForeground: Hex("#fafafa"),
	Primary: Hex("#fafafa"), PrimaryForeground: Hex("#18181b"),
	Secondary: Hex("#27272a"), SecondaryForeground: Hex("#fafafa"),
	Muted: Hex("#27272a"), MutedForeground: Hex("#a1a1aa"),
	Accent: Hex("#27272a"), AccentForeground: Hex("#fafafa"),
	Destructive: Hex("#7f1d1d"), DestructiveForeground: Hex("#fafafa"),
	Border: Hex("#27272a"), Input: Hex("#27272a"), Ring: Hex("#d4d4d8"),
	Radius: 10,
}

var themeContext = CreateContext(LightTheme)

// ThemeProvider 向子树提供主题。
func ThemeProvider(t Theme, children ...*Node) *Node {
	return themeContext.Provider(t, children...)
}

// UseTheme 读取当前主题（未包裹 ThemeProvider 时返回 LightTheme）。
func UseTheme() Theme { return UseContext(themeContext) }
