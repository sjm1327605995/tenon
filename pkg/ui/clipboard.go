package ui

// 应用内剪贴板；默认在应用内工作。要接系统剪贴板，用 SetClipboardProvider 注入
// （例如包一层 github.com/atotto/clipboard）。
var (
	clipboardText string
	getClipboard  = func() string { return clipboardText }
	setClipboard  = func(s string) { clipboardText = s }
)

// Clipboard 返回当前剪贴板文本。
func Clipboard() string { return getClipboard() }

// SetClipboardText 设置剪贴板文本。
func SetClipboardText(s string) { setClipboard(s) }

// SetClipboardProvider 用外部实现（如系统剪贴板）替换默认的应用内剪贴板。
func SetClipboardProvider(get func() string, set func(string)) {
	if get != nil {
		getClipboard = get
	}
	if set != nil {
		setClipboard = set
	}
}
