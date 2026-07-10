package ui

import _ "embed"

// cjkFont 是内置的中日韩字体（OPPOSans），作为默认字体以支持 CJK 文本。
//
//go:embed assets/OPPOSans-Medium.ttf
var cjkFont []byte
