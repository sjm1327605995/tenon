package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ShortcutKey 修饰键。
type ShortcutKey int

const (
	ShortcutCtrl  ShortcutKey = iota // Ctrl (Windows/Linux) 或 Cmd (macOS)
	ShortcutShift                    // Shift
	ShortcutAlt                      // Alt
)

// Shortcut 定义一个快捷键组合。
type Shortcut struct {
	Modifiers []ShortcutKey // 修饰键组合
	Key       ebiten.Key    // 主键
	Handler   func()        // 触发时的回调
	Priority  int           // 优先级，高优先级先处理
	Global    bool          // 是否全局（不受焦点限制）
}

// ShortcutManager 管理快捷键注册和分发。
type ShortcutManager struct {
	shortcuts []Shortcut
	enabled   bool
}

// NewShortcutManager 创建快捷键管理器。
func NewShortcutManager() *ShortcutManager {
	return &ShortcutManager{enabled: true}
}

// Register 注册快捷键。
func (sm *ShortcutManager) Register(s Shortcut) {
	sm.shortcuts = append(sm.shortcuts, s)
}

// Unregister 移除匹配的快捷键。
func (sm *ShortcutManager) Unregister(key ebiten.Key, modifiers ...ShortcutKey) {
	for i, s := range sm.shortcuts {
		if s.Key == key && matchModifiers(s.Modifiers, modifiers) {
			sm.shortcuts = append(sm.shortcuts[:i], sm.shortcuts[i+1:]...)
			return
		}
	}
}

// SetEnabled 启用/禁用快捷键处理。
func (sm *ShortcutManager) SetEnabled(v bool) {
	sm.enabled = v
}

// IsEnabled 返回是否启用。
func (sm *ShortcutManager) IsEnabled() bool {
	return sm.enabled
}

// Handle 处理快捷键，返回是否有快捷键被触发。
func (sm *ShortcutManager) Handle() bool {
	if !sm.enabled {
		return false
	}

	for _, s := range sm.shortcuts {
		if sm.matchShortcut(s) {
			if s.Handler != nil {
				s.Handler()
			}
			return true
		}
	}
	return false
}

func (sm *ShortcutManager) matchShortcut(s Shortcut) bool {
	// 检查主键
	if !inpututil.IsKeyJustPressed(s.Key) {
		return false
	}
	// 检查修饰键
	for _, mod := range s.Modifiers {
		switch mod {
		case ShortcutCtrl:
			if !ebiten.IsKeyPressed(ebiten.KeyControl) {
				return false
			}
		case ShortcutShift:
			if !ebiten.IsKeyPressed(ebiten.KeyShift) {
				return false
			}
		case ShortcutAlt:
			if !ebiten.IsKeyPressed(ebiten.KeyAlt) {
				return false
			}
		}
	}
	return true
}

func matchModifiers(a, b []ShortcutKey) bool {
	if len(a) != len(b) {
		return false
	}
	for i, m := range a {
		if m != b[i] {
			return false
		}
	}
	return true
}

// ==================== 快捷键注册辅助 ====================

// Ctrl 快捷键组合。
func Ctrl(key ebiten.Key, handler func()) Shortcut {
	return Shortcut{
		Modifiers: []ShortcutKey{ShortcutCtrl},
		Key:       key,
		Handler:   handler,
	}
}

// CtrlShift 快捷键组合。
func CtrlShift(key ebiten.Key, handler func()) Shortcut {
	return Shortcut{
		Modifiers: []ShortcutKey{ShortcutCtrl, ShortcutShift},
		Key:       key,
		Handler:   handler,
	}
}

// Alt 快捷键组合。
func Alt(key ebiten.Key, handler func()) Shortcut {
	return Shortcut{
		Modifiers: []ShortcutKey{ShortcutAlt},
		Key:       key,
		Handler:   handler,
	}
}
