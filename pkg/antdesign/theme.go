package antdesign

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/core"
)

// AntTheme 扩展 Tenon 主题，补充 Ant Design 语义化颜色。
// 继承 core.Theme 的所有字段，同时增加 AntD 特有的配色维度。
type AntTheme struct {
	*core.Theme

	// 语义色板
	ErrorColor   color.Color // 错误 / 危险
	WarningColor color.Color // 警告
	SuccessColor color.Color // 成功
	InfoColor    color.Color // 信息

	// 危险按钮
	DangerNormalColor  color.Color
	DangerHoverColor   color.Color
	DangerPressedColor color.Color
	DangerTextColor    color.Color
	DangerBgColor      color.Color // 危险背景（用于 Alert、Tag 等）

	// 警告按钮
	WarningNormalColor  color.Color
	WarningHoverColor   color.Color
	WarningPressedColor color.Color
	WarningTextColor    color.Color
	WarningBgColor      color.Color

	// 成功按钮
	SuccessNormalColor  color.Color
	SuccessHoverColor   color.Color
	SuccessPressedColor color.Color
	SuccessTextColor    color.Color
	SuccessBgColor      color.Color

	// 文字按钮 / 链接按钮
	TextButtonColor      color.Color
	TextButtonHoverColor color.Color
	LinkColor            color.Color
	LinkHoverColor       color.Color

	// 幽灵按钮边框
	GhostBorderColor      color.Color
	GhostBorderHoverColor color.Color
	GhostTextColor        color.Color

	// 禁用状态（更精细）
	DisabledBgColor   color.Color
	DisabledTextColor color.Color
	DisabledBorderColor color.Color

	// 斑马纹、悬停背景
	TableRowHoverBg  color.Color
	TableStripeBg    color.Color
	TableHeaderBg    color.Color
	TableHeaderColor color.Color

	// Alert 背景色（带透明度）
	AlertErrorBg   color.Color
	AlertWarningBg color.Color
	AlertSuccessBg color.Color
	AlertInfoBg    color.Color

	// Tag 颜色映射
	TagRedBg    color.Color
	TagRedText  color.Color
	TagOrangeBg color.Color
	TagOrangeText color.Color
	TagGreenBg  color.Color
	TagGreenText color.Color
	TagBlueBg   color.Color
	TagBlueText color.Color
}

// NewAntTheme 从当前全局主题创建 Ant Design 扩展主题。
func NewAntTheme() *AntTheme {
	return ExtendTheme(core.GetTheme())
}

// ExtendTheme 将 core.Theme 扩展为 AntTheme。
func ExtendTheme(base *core.Theme) *AntTheme {
	return &AntTheme{
		Theme: base,

		// 语义色板（Ant Design 官方色值）
		ErrorColor:   color.RGBA{R: 255, G: 77, B: 79, A: 255},   // #ff4d4f
		WarningColor: color.RGBA{R: 250, G: 173, B: 20, A: 255},  // #faad14
		SuccessColor: color.RGBA{R: 82, G: 196, B: 26, A: 255},  // #52c41a
		InfoColor:    color.RGBA{R: 22, G: 119, B: 255, A: 255},  // #1677ff

		// 危险按钮
		DangerNormalColor:  color.RGBA{R: 255, G: 77, B: 79, A: 255},
		DangerHoverColor:   color.RGBA{R: 255, G: 120, B: 117, A: 255},
		DangerPressedColor: color.RGBA{R: 207, G: 35, B: 35, A: 255},
		DangerTextColor:    color.White,
		DangerBgColor:      color.RGBA{R: 255, G: 241, B: 240, A: 255}, // #fff1f0

		// 警告按钮
		WarningNormalColor:  color.RGBA{R: 250, G: 173, B: 20, A: 255},
		WarningHoverColor:   color.RGBA{R: 255, G: 197, B: 61, A: 255},
		WarningPressedColor: color.RGBA{R: 212, G: 136, B: 6, A: 255},
		WarningTextColor:    color.White,
		WarningBgColor:      color.RGBA{R: 255, G: 251, B: 230, A: 255}, // #fffbe6

		// 成功按钮
		SuccessNormalColor:  color.RGBA{R: 82, G: 196, B: 26, A: 255},
		SuccessHoverColor:    color.RGBA{R: 115, G: 209, B: 61, A: 255},
		SuccessPressedColor:  color.RGBA{R: 56, G: 158, B: 13, A: 255},
		SuccessTextColor:     color.White,
		SuccessBgColor:       color.RGBA{R: 246, G: 255, B: 237, A: 255}, // #f6ffed

		// 文字 / 链接按钮
		TextButtonColor:      base.PrimaryColor,
		TextButtonHoverColor: base.PrimaryHoverColor,
		LinkColor:            base.PrimaryColor,
		LinkHoverColor:       base.PrimaryHoverColor,

		// 幽灵按钮
		GhostBorderColor:      color.RGBA{R: 217, G: 217, B: 217, A: 255},
		GhostBorderHoverColor: base.PrimaryColor,
		GhostTextColor:        color.RGBA{R: 0, G: 0, B: 0, A: 224},

		// 禁用状态
		DisabledBgColor:     color.RGBA{R: 245, G: 245, B: 245, A: 255}, // #f5f5f5
		DisabledTextColor:   color.RGBA{R: 0, G: 0, B: 0, A: 64},       // rgba(0,0,0,0.25)
		DisabledBorderColor: color.RGBA{R: 217, G: 217, B: 217, A: 255}, // #d9d9d9

		// 表格
		TableRowHoverBg:  color.RGBA{R: 250, G: 250, B: 250, A: 255}, // #fafafa
		TableStripeBg:    color.RGBA{R: 250, G: 250, B: 250, A: 255},
		TableHeaderBg:    color.RGBA{R: 250, G: 250, B: 250, A: 255},
		TableHeaderColor: color.RGBA{R: 0, G: 0, B: 0, A: 224},

		// Alert 背景
		AlertErrorBg:   color.RGBA{R: 255, G: 241, B: 240, A: 255},
		AlertWarningBg: color.RGBA{R: 255, G: 251, B: 230, A: 255},
		AlertSuccessBg: color.RGBA{R: 246, G: 255, B: 237, A: 255},
		AlertInfoBg:    color.RGBA{R: 230, G: 244, B: 255, A: 255},

		// Tag 颜色
		TagRedBg:      color.RGBA{R: 255, G: 241, B: 240, A: 255},
		TagRedText:    color.RGBA{R: 207, G: 19, B: 34, A: 255},
		TagOrangeBg:   color.RGBA{R: 255, G: 247, B: 230, A: 255},
		TagOrangeText: color.RGBA{R: 212, G: 107, B: 8, A: 255},
		TagGreenBg:    color.RGBA{R: 246, G: 255, B: 237, A: 255},
		TagGreenText:  color.RGBA{R: 56, G: 158, B: 13, A: 255},
		TagBlueBg:     color.RGBA{R: 230, G: 244, B: 255, A: 255},
		TagBlueText:   color.RGBA{R: 9, G: 88, B: 217, A: 255},
	}
}

// SetAntTheme 将 AntTheme 基础色设为全局主题。
func SetAntTheme() {
	core.SetTheme(ExtendTheme(core.DefaultAntTheme()).Theme)
}
