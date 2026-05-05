package ui

import (
	"os/exec"
	"runtime"
	"strings"
)

// Clipboard 提供跨平台剪贴板读写能力。
// Linux: xclip/xsel/wl-clipboard
// macOS: pbcopy/pbpaste
// Windows: clip/powershell
type Clipboard struct {
	// 外部注入的读写函数（可用于覆盖平台默认实现）
	ReadFunc  func() string
	WriteFunc func(text string)
}

// NewClipboard 创建剪贴板实例，自动检测平台工具。
func NewClipboard() *Clipboard {
	c := &Clipboard{}
	c.detectPlatform()
	return c
}

// Read 读取剪贴板内容。
func (c *Clipboard) Read() string {
	if c.ReadFunc != nil {
		return c.ReadFunc()
	}
	return ""
}

// Write 写入剪贴板内容。
func (c *Clipboard) Write(text string) {
	if c.WriteFunc != nil {
		c.WriteFunc(text)
	}
}

// HasContent 剪贴板是否有内容。
func (c *Clipboard) HasContent() bool {
	return strings.TrimSpace(c.Read()) != ""
}

// detectPlatform 自动检测当前平台可用的剪贴板工具。
func (c *Clipboard) detectPlatform() {
	switch runtime.GOOS {
	case "darwin":
		c.ReadFunc = func() string {
			out, err := exec.Command("pbpaste").Output()
			if err != nil {
				return ""
			}
			return string(out)
		}
		c.WriteFunc = func(text string) {
			cmd := exec.Command("pbcopy")
			cmd.Stdin = strings.NewReader(text)
			cmd.Run()
		}
	case "linux":
		// 优先 Wayland (wl-copy/wl-paste)，fallback 到 X11 (xclip)
		c.ReadFunc = func() string {
			// Try wl-paste first
			if out, err := exec.Command("wl-paste", "-n").Output(); err == nil {
				return string(out)
			}
			// Try xclip
			if out, err := exec.Command("xclip", "-selection", "clipboard", "-o").Output(); err == nil {
				return string(out)
			}
			// Try xsel
			if out, err := exec.Command("xsel", "--clipboard", "--output").Output(); err == nil {
				return string(out)
			}
			return ""
		}
		c.WriteFunc = func(text string) {
			// Try wl-copy first
			cmd := exec.Command("wl-copy")
			cmd.Stdin = strings.NewReader(text)
			if cmd.Run() == nil {
				return
			}
			// Try xclip
			cmd = exec.Command("xclip", "-selection", "clipboard")
			cmd.Stdin = strings.NewReader(text)
			if cmd.Run() == nil {
				return
			}
			// Try xsel
			cmd = exec.Command("xsel", "--clipboard", "--input")
			cmd.Stdin = strings.NewReader(text)
			cmd.Run()
		}
	case "windows":
		c.ReadFunc = func() string {
			out, err := exec.Command("powershell", "-command", "Get-Clipboard").Output()
			if err != nil {
				return ""
			}
			return strings.TrimRight(string(out), "\r\n")
		}
		c.WriteFunc = func(text string) {
			cmd := exec.Command("cmd", "/c", "clip")
			cmd.Stdin = strings.NewReader(text)
			cmd.Run()
		}
	}
}

// ==================== 全局剪贴板 ====================

var globalClipboard *Clipboard

// GetClipboard 获取全局剪贴板实例。
func GetClipboard() *Clipboard {
	if globalClipboard == nil {
		globalClipboard = NewClipboard()
	}
	return globalClipboard
}

// ClipboardRead 读取剪贴板快捷函数。
func ClipboardRead() string {
	return GetClipboard().Read()
}

// ClipboardWrite 写入剪贴板快捷函数。
func ClipboardWrite(text string) {
	GetClipboard().Write(text)
}
