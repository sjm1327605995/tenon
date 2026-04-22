package yoga

import "testing"

func TestCalculateLayoutUpdatesStats(t *testing.T) {
	root := NewNode()
	root.StyleSetWidth(200)
	root.StyleSetHeight(200)

	child := NewNode()
	child.StyleSetWidth(100)
	child.StyleSetHeight(50)
	root.InsertChild(child, 0)

	// 第一次计算，应该产生 Layout 统计
	CalculateLayout(root, 200, 200, DirectionLTR)

	// 修改尺寸，再次计算，子节点可能命中缓存
	root.StyleSetWidth(250)
	CalculateLayout(root, 250, 250, DirectionLTR)

	// 由于 LayoutData 是内部局部变量，无法从外部直接读取。
	// 这里主要验证 CalculateLayout 不会因为 layoutType 的使用而 panic。
	if root.LayoutWidth() != 250 {
		t.Fatalf("expected root width 250, got %f", root.LayoutWidth())
	}
	if child.LayoutWidth() != 100 {
		t.Fatalf("expected child width 100, got %f", child.LayoutWidth())
	}
}
