package ui

import (
	"testing"

	"gioui.org/font/opentype"
	"gioui.org/text"
)

// 内置字体必须解析得出来。它是 //go:embed 的资源，解析不了就是资源坏了/换错了。
func TestEmbeddedFontParses(t *testing.T) {
	if len(cjkFont) == 0 {
		t.Fatal("cjkFont 是空的 —— //go:embed 没生效")
	}
	if _, err := opentype.Parse(cjkFont); err != nil {
		t.Fatalf("内置字体解析失败: %v", err)
	}
	t.Logf("assets/OPPOSans-Medium.ttf 解析成功（%d 字节）", len(cjkFont))
}

// shaper 必须真的在用内置字体，而不是悄悄回落到系统字体。
//
// 判据是「字形 ID 与系统回落不同」。不能用「有没有 notdef 豆腐」来判 ——
// gio 在字体集为空时会回落到系统字体，而开发机上通常装着中文字体，
// 于是空字体集下中文照样没有豆腐（实测 notdef=0），那种测试证明不了任何事。
func TestShaperUsesEmbeddedFontNotSystemFallback(t *testing.T) {
	glyphs := func(sh *text.Shaper, s string) []uint32 {
		f := gioNewFont(20, 400, false).(*gioFont)
		sh.LayoutString(f.params(), s)
		var ids []uint32
		for {
			g, ok := sh.NextGlyph()
			if !ok {
				return ids
			}
			ids = append(ids, uint32(g.ID))
		}
	}
	const cjk = "中文字体"
	ours := glyphs(gioShaper(), cjk)
	fallback := glyphs(text.NewShaper(text.WithCollection(nil)), cjk) // 系统字体

	t.Logf("内置字体   %q -> %v", cjk, ours)
	t.Logf("系统回落   %q -> %v", cjk, fallback)

	if len(ours) != len([]rune(cjk)) {
		t.Fatalf("内置字体给出 %d 个字形，want %d", len(ours), len([]rune(cjk)))
	}
	for _, id := range ours {
		if id == 0 {
			t.Fatalf("内置字体下 %q 出现 notdef（豆腐）：%v", cjk, ours)
		}
	}
	same := len(ours) == len(fallback)
	if same {
		for i := range ours {
			if ours[i] != fallback[i] {
				same = false
				break
			}
		}
	}
	if same {
		t.Errorf("字形 ID 与系统回落完全一致（%v）—— 内置字体很可能没生效", ours)
	}
}
