package font

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// FontFamily 字体族定义
type FontFamily string

const (
	FontFamilyDefault FontFamily = "default"
	FontFamilySerif   FontFamily = "serif"
	FontFamilySans    FontFamily = "sans"
	FontFamilyMono    FontFamily = "mono"
)

// FontWeight 字体粗细
type FontWeight int

const (
	FontWeightNormal FontWeight = 400
	FontWeightBold   FontWeight = 700
	FontWeightLight  FontWeight = 300
)

// FontStyle 字体样式
type FontStyle int

const (
	FontStyleNormal FontStyle = iota
	FontStyleItalic
)

// FontDescriptor 字体描述符
type FontDescriptor struct {
	Family FontFamily
	Weight FontWeight
	Style  FontStyle
	Size   float32
}

// FontFace 字体面。当没有对应字重/斜体的真实字体源时，FauxBold/FauxItalic 为真，
// 由渲染层用合成粗体（重叠描边）/斜体（错切）近似。
type FontFace struct {
	Descriptor FontDescriptor
	Face       *text.GoTextFace
	FauxBold   bool
	FauxItalic bool
}

type variantKey struct {
	family       FontFamily
	bold, italic bool
}

// FontManager 字体管理器
type FontManager struct {
	mu             sync.RWMutex
	fontSources    map[FontFamily]*text.GoTextFaceSource
	variantSources map[variantKey]*text.GoTextFaceSource // 已注册的粗体/斜体变体
	fontCache      map[FontDescriptor]*FontFace
	defaultFamily  FontFamily
}

var (
	instance *FontManager
	once     sync.Once
)

// GetFontManager 获取字体管理器单例
func GetFontManager() *FontManager {
	once.Do(func() {
		instance = &FontManager{
			fontSources:    make(map[FontFamily]*text.GoTextFaceSource),
			variantSources: make(map[variantKey]*text.GoTextFaceSource),
			fontCache:      make(map[FontDescriptor]*FontFace),
			defaultFamily:  FontFamilyDefault,
		}
	})
	return instance
}

// LoadFontFromBytes 从字节加载字体
func (fm *FontManager) LoadFontFromBytes(family FontFamily, fontData []byte) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if _, exists := fm.fontSources[family]; exists {
		return fmt.Errorf("font family '%s' already loaded", family)
	}

	return fm.loadFontLocked(family, fontData)
}

// ReloadFontFromBytes 重新加载字体（覆盖已存在的字体族）
func (fm *FontManager) ReloadFontFromBytes(family FontFamily, fontData []byte) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	return fm.loadFontLocked(family, fontData)
}

func (fm *FontManager) loadFontLocked(family FontFamily, fontData []byte) error {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		return fmt.Errorf("failed to load font: %v", err)
	}

	fm.fontSources[family] = source
	fm.clearFontCache() // 清除缓存以重新生成字体面
	return nil
}

// LoadFontFromReader 从读取器加载字体
func (fm *FontManager) LoadFontFromReader(family FontFamily, reader io.Reader) error {
	fontData, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read font data: %v", err)
	}
	return fm.LoadFontFromBytes(family, fontData)
}

// SetDefaultFontFamily 设置默认字体族
func (fm *FontManager) SetDefaultFontFamily(family FontFamily) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.defaultFamily = family
}

// GetFontFace 获取字体面
func (fm *FontManager) GetFontFace(descriptor FontDescriptor) (*FontFace, error) {
	fm.mu.RLock()

	// 检查缓存
	if face, exists := fm.fontCache[descriptor]; exists {
		fm.mu.RUnlock()
		return face, nil
	}
	fm.mu.RUnlock()

	fm.mu.Lock()
	defer fm.mu.Unlock()

	// 双重检查
	if face, exists := fm.fontCache[descriptor]; exists {
		return face, nil
	}

	// 选择字体源：优先注册过的粗体/斜体变体，否则用基础字体并标记合成
	bold := descriptor.Weight >= 600
	italic := descriptor.Style == FontStyleItalic
	source, real := fm.pickSource(descriptor.Family, bold, italic)
	if source == nil {
		return nil, fmt.Errorf("font family '%s' not found", descriptor.Family)
	}

	face := &FontFace{
		Descriptor: descriptor,
		Face:       &text.GoTextFace{Source: source, Size: float64(descriptor.Size)},
	}
	if !real {
		face.FauxBold = bold
		face.FauxItalic = italic
	}

	fm.fontCache[descriptor] = face
	return face, nil
}

// pickSource 选择最匹配的字体源；real 表示命中真实变体（而非需合成）。
func (fm *FontManager) pickSource(family FontFamily, bold, italic bool) (*text.GoTextFaceSource, bool) {
	if bold || italic {
		if s, ok := fm.variantSources[variantKey{family, bold, italic}]; ok {
			return s, true
		}
	}
	return fm.getFontSource(family), false
}

// LoadFontVariant 注册某字体族的粗体/斜体变体源（可选，用于替代合成、获得更高质量）。
func (fm *FontManager) LoadFontVariant(family FontFamily, bold, italic bool, data []byte) error {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return err
	}
	fm.mu.Lock()
	fm.variantSources[variantKey{family, bold, italic}] = src
	fm.fontCache = make(map[FontDescriptor]*FontFace)
	fm.mu.Unlock()
	return nil
}

// GetDefaultFontFace 获取默认字体面（常规字重）
func (fm *FontManager) GetDefaultFontFace(size float32) (*FontFace, error) {
	return fm.GetDefaultFace(size, FontWeightNormal, false)
}

// GetDefaultFace 获取默认字体族在指定字重/斜体下的字体面。
func (fm *FontManager) GetDefaultFace(size float32, weight FontWeight, italic bool) (*FontFace, error) {
	style := FontStyleNormal
	if italic {
		style = FontStyleItalic
	}
	return fm.GetFontFace(FontDescriptor{Family: fm.defaultFamily, Weight: weight, Style: style, Size: size})
}

// getFontSource 获取字体源
func (fm *FontManager) getFontSource(family FontFamily) *text.GoTextFaceSource {
	if source, exists := fm.fontSources[family]; exists {
		return source
	}

	// 回退到默认字体
	if source, exists := fm.fontSources[FontFamilyDefault]; exists {
		return source
	}

	return nil
}

// clearFontCache 清除字体缓存
func (fm *FontManager) clearFontCache() {
	fm.fontCache = make(map[FontDescriptor]*FontFace)
}

// ListLoadedFonts 列出已加载的字体
func (fm *FontManager) ListLoadedFonts() []FontFamily {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	families := make([]FontFamily, 0, len(fm.fontSources))
	for family := range fm.fontSources {
		families = append(families, family)
	}
	return families
}

// HasFontFamily 检查字体族是否已加载
func (fm *FontManager) HasFontFamily(family FontFamily) bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	_, exists := fm.fontSources[family]
	return exists
}

// PreloadCommonSizes 预加载常用尺寸的字体
func (fm *FontManager) PreloadCommonSizes(family FontFamily, sizes []float32) error {
	for _, size := range sizes {
		_, err := fm.GetFontFace(FontDescriptor{
			Family: family,
			Weight: FontWeightNormal,
			Style:  FontStyleNormal,
			Size:   size,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
