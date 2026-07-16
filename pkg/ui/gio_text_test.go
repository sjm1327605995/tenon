package ui

import "testing"

// gioFont.Measure 应随字符数近似线性增长（横向排布）。若 shaper 因 MaxWidth=0 逐字竖排，
// 单字与多字的推进宽度会退化、不成比例，本测试即可捕获。
func TestGioTextHorizontalAdvance(t *testing.T) {
	f, ok := gioNewFont(16, 400, false).(*gioFont)
	if !ok || f == nil {
		t.Fatal("gioNewFont 应返回 *gioFont")
	}
	w1 := f.Measure("A", 16*1.3)
	w3 := f.Measure("AAA", 16*1.3)
	if w1 <= 0 || w3 <= 0 {
		t.Fatalf("测得宽度应为正: w1=%v w3=%v", w1, w3)
	}
	// 三个字应约等于单字的 3 倍（容一定字距误差），至少要显著大于单字，证明是横向累加而非竖排。
	if w3 < w1*2.2 || w3 > w1*3.8 {
		t.Fatalf("横向推进异常: w1=%v w3=%v (期望 w3≈3*w1)", w1, w3)
	}
	if a, _, _ := f.Metrics(); a <= 0 {
		t.Fatalf("ascent 应为正: %v", a)
	}
}
