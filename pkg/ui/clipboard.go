package ui

// 剪贴板。默认实现会同时更新应用内缓存并经后端写入系统剪贴板，所以复制/剪切默认就能
// 和其它程序互通。clipboardText 是最近一次读写的缓存。
var (
	clipboardText string
	getClipboard  = defaultGetClipboard
	setClipboard  = defaultSetClipboard
)

// backendWriteClipboard 由当前后端登记：把文本写进系统剪贴板。
// 注意要放在默认 provider 里，而不是在后端 init 里直接改写 setClipboard ——
// 否则任何「把 setClipboard 还原成默认」的代码都会把系统剪贴板接线悄悄弄丢。
var backendWriteClipboard func(string)

func defaultGetClipboard() string { return clipboardText }

func defaultSetClipboard(s string) {
	clipboardText = s
	if backendWriteClipboard != nil {
		backendWriteClipboard(s)
	}
}

// Clipboard 返回最近一次读写到的剪贴板文本（缓存）。
//
// 注意它不会去读系统剪贴板：系统剪贴板的读取是异步的，粘贴（Ctrl+V）由后端接管并直接
// 落到聚焦输入框，不经过这里。所以在别的程序里复制后、本应用尚未粘贴过时，这里可能是旧值。
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
