package engine

import (
	"testing"
)

func TestClipboardReadWrite(t *testing.T) {
	clip := NewClipboard()
	// 注入读写函数（避免依赖系统剪贴板）
	var stored string
	clip.WriteFunc = func(text string) { stored = text }
	clip.ReadFunc = func() string { return stored }

	clip.Write("hello")
	if clip.Read() != "hello" {
		t.Errorf("expected 'hello', got '%s'", clip.Read())
	}

	clip.Write("world")
	if clip.Read() != "world" {
		t.Errorf("expected 'world', got '%s'", clip.Read())
	}
}

func TestClipboardHasContent(t *testing.T) {
	clip := NewClipboard()
	var stored string
	clip.WriteFunc = func(text string) { stored = text }
	clip.ReadFunc = func() string { return stored }

	if clip.HasContent() {
		t.Error("expected empty clipboard")
	}

	clip.Write("test")
	if !clip.HasContent() {
		t.Error("expected non-empty clipboard")
	}
}

func TestShortcutManagerRegister(t *testing.T) {
	sm := NewShortcutManager()
	called := false
	_ = called
	sm.Register(Ctrl(0, func() { called = true }))

	if len(sm.shortcuts) != 1 {
		t.Errorf("expected 1 shortcut, got %d", len(sm.shortcuts))
	}
}

func TestShortcutManagerUnregister(t *testing.T) {
	sm := NewShortcutManager()
	sm.Register(Ctrl(0, func() {}))
	sm.Unregister(0, ShortcutCtrl)

	if len(sm.shortcuts) != 0 {
		t.Errorf("expected 0 shortcuts, got %d", len(sm.shortcuts))
	}
}

func TestShortcutManagerEnabled(t *testing.T) {
	sm := NewShortcutManager()
	if !sm.IsEnabled() {
		t.Error("expected enabled by default")
	}

	sm.SetEnabled(false)
	if sm.IsEnabled() {
		t.Error("expected disabled")
	}
}

func TestMatchModifiers(t *testing.T) {
	a := []ShortcutKey{ShortcutCtrl, ShortcutShift}
	b := []ShortcutKey{ShortcutCtrl, ShortcutShift}
	if !matchModifiers(a, b) {
		t.Error("expected match")
	}

	c := []ShortcutKey{ShortcutCtrl}
	if matchModifiers(a, c) {
		t.Error("expected no match")
	}
}
