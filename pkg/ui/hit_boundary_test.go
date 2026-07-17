package ui

import "testing"

// 元素边界落在小数坐标上时（uiScale=1.5 之类的非整数缩放下很常见），
// 光标必须能精确命中边界内侧：把光标截断成整数会让边界处最多差 1px，
// 表现为「贴着边框时 hover 时有时无」。
func TestHitAtFractionalBoundary(t *testing.T) {
	rn := &renderNode{bounds: Rect{X: 10.5, Y: 10.5, W: 20, H: 20}}

	cases := []struct {
		x, y float32
		want bool
		why  string
	}{
		{10.5, 10.5, true, "左上角内侧起点"},
		{10.6, 10.6, true, "边界内侧 0.1px"},
		{10.4, 10.6, false, "左边界外侧 0.1px"},
		{30.4, 30.4, true, "右下角内侧"},
		{30.6, 20, false, "右边界外侧"},
	}
	for _, c := range cases {
		if got := hitNode(rn, c.x, c.y) != nil; got != c.want {
			t.Errorf("(%v,%v) 命中=%v want %v —— %s", c.x, c.y, got, c.want, c.why)
		}
	}
}

// 光标坐标必须保留亚像素精度：截断成整数后，(10.6,10.6) 会变成 (10,10)，
// 落到 X=10.5 起的元素之外 —— 明明指针就在元素里，却命中不到。
func TestCursorKeepsSubPixelPrecision(t *testing.T) {
	in := &gioInput{}
	in.setCursor(10.6, 10.6)
	x, y := in.cursor()
	if x != 10.6 || y != 10.6 {
		t.Fatalf("cursor=(%v,%v) want (10.6,10.6) —— 亚像素精度被丢掉了", x, y)
	}

	rn := &renderNode{bounds: Rect{X: 10.5, Y: 10.5, W: 20, H: 20}}
	if hitNode(rn, x, y) == nil {
		t.Fatal("指针在元素内却命中不到（坐标被截断到了元素之外）")
	}
}
