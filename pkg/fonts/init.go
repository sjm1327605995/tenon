package fonts

import (
	"io"
	"os"

	"golang.org/x/image/font/gofont/goregular"
)

// InitDefaultFont 初始化默认字体
func InitDefaultFont() error {
	fontManager := GetFontManager()

	// 加载默认字体
	err := fontManager.LoadFontFromBytes(FontFamilyDefault, goregular.TTF)
	if err != nil {
		return err
	}

	// 预加载常用尺寸
	sizes := []float32{12, 14, 16, 18, 20, 24, 32, 48}
	return fontManager.PreloadCommonSizes(FontFamilyDefault, sizes)
}

// LoadFontFromFile 从文件加载字体
func LoadFontFromFile(family FontFamily, filePath string) error {
	fontManager := GetFontManager()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return fontManager.LoadFontFromReader(family, file)
}

// LoadFontFromBytes 从字节加载字体（便捷函数）
func LoadFontFromBytes(family FontFamily, fontData []byte) error {
	fontManager := GetFontManager()
	return fontManager.LoadFontFromBytes(family, fontData)
}

// LoadFontFromReader 从读取器加载字体（便捷函数）
func LoadFontFromReader(family FontFamily, reader io.Reader) error {
	fontManager := GetFontManager()
	return fontManager.LoadFontFromReader(family, reader)
}

// SetDefaultFontFamily 设置默认字体族（便捷函数）
func SetDefaultFontFamily(family FontFamily) {
	fontManager := GetFontManager()
	fontManager.SetDefaultFontFamily(family)
}

// GetFontFace 获取字体面（便捷函数）
func GetFontFace(descriptor FontDescriptor) (*FontFace, error) {
	fontManager := GetFontManager()
	return fontManager.GetFontFace(descriptor)
}

// GetDefaultFontFace 获取默认字体面（便捷函数）
func GetDefaultFontFace(size float32) (*FontFace, error) {
	fontManager := GetFontManager()
	return fontManager.GetDefaultFontFace(size)
}

// ListLoadedFonts 列出已加载的字体（便捷函数）
func ListLoadedFonts() []FontFamily {
	fontManager := GetFontManager()
	return fontManager.ListLoadedFonts()
}

// HasFontFamily 检查字体族是否已加载（便捷函数）
func HasFontFamily(family FontFamily) bool {
	fontManager := GetFontManager()
	return fontManager.HasFontFamily(family)
}

// PreloadCommonSizes 预加载常用尺寸（便捷函数）
func PreloadCommonSizes(family FontFamily, sizes []float32) error {
	fontManager := GetFontManager()
	return fontManager.PreloadCommonSizes(family, sizes)
}
