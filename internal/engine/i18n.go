package engine

import (
	"fmt"
	"reflect"
)

// ==================== Localization ====================

var localizationType = reflect.TypeOf(Localization{})

// Localization 是一个 InheritedWidget，向下传递本地化数据。
type Localization struct {
	BaseWidget
	Locale        string            // 当前语言，如 "zh-CN", "en-US"
	Translations  map[string]map[string]string // locale -> key -> value
	child         Widget
}

// NewLocalization 创建本地化 Widget。
func NewLocalization(locale string, translations map[string]map[string]string, child Widget) Localization {
	return Localization{
		Locale:       locale,
		Translations: translations,
		child:        child,
	}
}

func (l Localization) CreateElement() Element {
	return NewInheritedElement(l)
}

func (l Localization) UpdateShouldNotify(oldWidget InheritedWidget) bool {
	old := oldWidget.(Localization)
	return l.Locale != old.Locale
}

func (l Localization) BuildChild(ctx BuildContext) Widget {
	return l.child
}

// L 从 BuildContext 获取本地化文本。
// key 是翻译键，如 "app.title"。
func L(ctx BuildContext, key string) string {
	if ctx == nil {
		return key
	}
	iw, ok := ctx.DependOnInheritedWidgetOfExactType(localizationType)
	if !ok || iw == nil {
		return key
	}
	loc, ok := iw.(Localization)
	if !ok {
		return key
	}
	return loc.Translate(key)
}

// Lf 从 BuildContext 获取本地化文本并格式化。
func Lf(ctx BuildContext, key string, args ...any) string {
	tmpl := L(ctx, key)
	if len(args) == 0 {
		return tmpl
	}
	return fmt.Sprintf(tmpl, args...)
}

// Translate 翻译指定 key。
func (l Localization) Translate(key string) string {
	if l.Translations == nil {
		return key
	}
	localeData, ok := l.Translations[l.Locale]
	if !ok {
		return key
	}
	if val, ok := localeData[key]; ok {
		return val
	}
	return key
}

// GetLocale 返回当前语言。
func (l Localization) GetLocale() string {
	return l.Locale
}

// ==================== Convenience ====================

// GetLocalization 从 BuildContext 获取 Localization。
func GetLocalization(ctx BuildContext) *Localization {
	if ctx == nil {
		return nil
	}
	iw, ok := ctx.DependOnInheritedWidgetOfExactType(localizationType)
	if !ok || iw == nil {
		return nil
	}
	loc, ok := iw.(Localization)
	if !ok {
		return nil
	}
	// 注意：返回的是 InheritedWidget 中存储的副本的指针。
	// 修改 *loc 不会影响树中的实际 widget。
	cp := loc
	return &cp
}

// GetLocale 从 BuildContext 获取当前语言。
func GetLocale(ctx BuildContext) string {
	if loc := GetLocalization(ctx); loc != nil {
		return loc.GetLocale()
	}
	return ""
}
