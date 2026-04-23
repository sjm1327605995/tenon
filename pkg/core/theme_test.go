package core

import (
	"image/color"
	"testing"
)

func TestDefaultLightTheme(t *testing.T) {
	th := DefaultLightTheme()
	if th == nil {
		t.Fatal("DefaultLightTheme should not be nil")
	}
	if th.PrimaryColor == nil {
		t.Error("PrimaryColor should not be nil")
	}
	if th.BorderRadius <= 0 {
		t.Error("BorderRadius should be > 0")
	}
	if th.FontSizeBase <= 0 {
		t.Error("FontSizeBase should be > 0")
	}
}

func TestDefaultDarkTheme(t *testing.T) {
	th := DefaultDarkTheme()
	if th == nil {
		t.Fatal("DefaultDarkTheme should not be nil")
	}
	if th.PrimaryColor == nil {
		t.Error("PrimaryColor should not be nil")
	}
}

func TestGetThemeReturnsDefault(t *testing.T) {
	// 重置默认主题
	defaultTheme = nil
	th := GetTheme()
	if th == nil {
		t.Fatal("GetTheme should return default theme when none is set")
	}
}

func TestSetTheme(t *testing.T) {
	custom := &Theme{
		PrimaryColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		BorderRadius: 16,
		FontSizeBase: 20,
	}
	SetTheme(custom)
	if GetTheme() != custom {
		t.Fatal("GetTheme should return the theme set by SetTheme")
	}
	if GetTheme().BorderRadius != 16 {
		t.Errorf("expected BorderRadius=16, got %f", GetTheme().BorderRadius)
	}

	// 恢复默认，避免影响其他测试
	SetTheme(DefaultLightTheme())
}

func TestThemeUsedByComponents(t *testing.T) {
	// 使用自定义主题创建组件，验证组件读取了主题值
	custom := &Theme{
		PrimaryColor:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
		ButtonNormalColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		SwitchOnColor:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
		SliderFillColor:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
		TextColor:        color.RGBA{R: 0, G: 255, B: 0, A: 255},
		BorderRadius:     16,
		FontSizeBase:     20,
	}
	SetTheme(custom)
	defer SetTheme(DefaultLightTheme())

	// 由于组件测试在 components 包，这里只验证 Theme 本身的设置/获取
	if GetTheme().PrimaryColor != custom.PrimaryColor {
		t.Error("PrimaryColor mismatch")
	}
	if GetTheme().BorderRadius != 16 {
		t.Errorf("expected BorderRadius=16, got %f", GetTheme().BorderRadius)
	}
}
