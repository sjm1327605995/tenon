package yoga

import "testing"

func yogaLayoutBenchTree(rows int) *Node {
	root := NewNode()
	root.StyleSetWidth(400)
	root.StyleSetHeight(400)
	for i := 0; i < rows; i++ {
		row := sizedNode(40, 20)
		root.InsertChild(row, uint32(i))
	}
	return root
}

func BenchmarkYogaLayoutDirty(b *testing.B) {
	root := yogaLayoutBenchTree(200)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.StyleSetGap(GutterAll, float32(i%2))
		root.CalculateLayout(400, 400, DirectionLTR)
	}
}

func BenchmarkYogaLayoutCached(b *testing.B) {
	root := yogaLayoutBenchTree(200)
	root.CalculateLayout(400, 400, DirectionLTR)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.CalculateLayout(400, 400, DirectionLTR)
	}
}
