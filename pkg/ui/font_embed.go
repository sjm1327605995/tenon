package ui

import _ "embed"

// cjkFont 是内置的中日韩字体（OPPOSans），作为默认字体以支持 CJK 文本。
//
//go:embed assets/OPPOSans-Medium.ttf
var cjkFont []byte

// DefaultFontData 返回内置 OPPOSans 字体的只读字节，供可选渲染后端复用。
func DefaultFontData() []byte { return cjkFont }
