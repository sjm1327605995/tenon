// Package fonts 提供字体管理功能。
//
// 已迁移至 internal/font，此包仅为兼容性保留。
// 请直接使用 tenon 包的统一入口。
package fonts

import internal "github.com/sjm1327605995/tenon/internal/font"

type (
	FontFamily   = internal.FontFamily
	FontWeight   = internal.FontWeight
	FontStyle    = internal.FontStyle
	FontDescriptor = internal.FontDescriptor
	FontFace     = internal.FontFace
	FontManager  = internal.FontManager
)

var (
	GetFontManager     = internal.GetFontManager
	LoadFontFromBytes  = internal.LoadFontFromBytes
	LoadFontFromReader = internal.LoadFontFromReader
	LoadFontFromFile   = internal.LoadFontFromFile
	ReloadFontFromFile = internal.ReloadFontFromFile
	SetDefaultFontFamily = internal.SetDefaultFontFamily
	GetFontFace        = internal.GetFontFace
	GetDefaultFontFace = internal.GetDefaultFontFace
	ListLoadedFonts    = internal.ListLoadedFonts
	HasFontFamily      = internal.HasFontFamily
	PreloadCommonSizes = internal.PreloadCommonSizes
	InitDefaultFont    = internal.InitDefaultFont
)

const (
	FontFamilyDefault = internal.FontFamilyDefault
	FontFamilySerif   = internal.FontFamilySerif
	FontFamilySans    = internal.FontFamilySans
	FontFamilyMono    = internal.FontFamilyMono
	FontWeightNormal  = internal.FontWeightNormal
	FontWeightBold    = internal.FontWeightBold
	FontWeightLight   = internal.FontWeightLight
	FontStyleNormal   = internal.FontStyleNormal
	FontStyleItalic   = internal.FontStyleItalic
)
