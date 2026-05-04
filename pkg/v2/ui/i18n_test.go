package ui

import (
	"testing"
)

func TestLocalizationTranslate(t *testing.T) {
	translations := map[string]map[string]string{
		"en-US": {
			"greeting": "Hello",
			"farewell": "Goodbye",
		},
		"zh-CN": {
			"greeting": "你好",
			"farewell": "再见",
		},
	}

	loc := NewLocalization("en-US", translations, nil)
	if loc.Translate("greeting") != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", loc.Translate("greeting"))
	}
	if loc.Translate("farewell") != "Goodbye" {
		t.Errorf("expected 'Goodbye', got '%s'", loc.Translate("farewell"))
	}
}

func TestLocalizationMissingKey(t *testing.T) {
	loc := NewLocalization("en-US", map[string]map[string]string{
		"en-US": {"greeting": "Hello"},
	}, nil)
	if loc.Translate("nonexistent") != "nonexistent" {
		t.Errorf("expected key fallback, got '%s'", loc.Translate("nonexistent"))
	}
}

func TestLocalizationMissingLocale(t *testing.T) {
	loc := NewLocalization("fr-FR", map[string]map[string]string{
		"en-US": {"greeting": "Hello"},
	}, nil)
	if loc.Translate("greeting") != "greeting" {
		t.Errorf("expected key fallback for missing locale, got '%s'", loc.Translate("greeting"))
	}
}

func TestLocalizationGetLocale(t *testing.T) {
	loc := NewLocalization("zh-CN", nil, nil)
	if loc.GetLocale() != "zh-CN" {
		t.Errorf("expected 'zh-CN', got '%s'", loc.GetLocale())
	}
}

func TestLocalizationUpdateShouldNotify(t *testing.T) {
	old := NewLocalization("en-US", nil, nil)
	newSame := NewLocalization("en-US", nil, nil)
	newDiff := NewLocalization("zh-CN", nil, nil)

	if newSame.UpdateShouldNotify(old) {
		t.Error("same locale should not notify")
	}
	if !newDiff.UpdateShouldNotify(old) {
		t.Error("different locale should notify")
	}
}

func TestLFromContext(t *testing.T) {
	translations := map[string]map[string]string{
		"en-US": {"title": "My App"},
		"zh-CN": {"title": "我的应用"},
	}

	var title string
	eng := NewEngine(func() Widget {
		return NewLocalization("zh-CN", translations,
			NewBuilder(func(ctx BuildContext) Widget {
				title = L(ctx, "title")
				return nil
			}),
		)
	}, 800, 600)
	eng.Mount()

	if title != "我的应用" {
		t.Errorf("expected '我的应用', got '%s'", title)
	}
}

func TestLFromContextEnglish(t *testing.T) {
	translations := map[string]map[string]string{
		"en-US": {"title": "My App"},
		"zh-CN": {"title": "我的应用"},
	}

	var title string
	eng := NewEngine(func() Widget {
		return NewLocalization("en-US", translations,
			NewBuilder(func(ctx BuildContext) Widget {
				title = L(ctx, "title")
				return nil
			}),
		)
	}, 800, 600)
	eng.Mount()

	if title != "My App" {
		t.Errorf("expected 'My App', got '%s'", title)
	}
}

func TestLNilContext(t *testing.T) {
	result := L(nil, "key")
	if result != "key" {
		t.Errorf("expected 'key' fallback, got '%s'", result)
	}
}
