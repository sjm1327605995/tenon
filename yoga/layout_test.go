package yoga

import "testing"

type layoutRect struct {
	left, top, width, height float32
}

func sizedNode(width, height float32) *Node {
	n := NewNode()
	n.StyleSetWidth(width)
	n.StyleSetHeight(height)
	return n
}

func assertLayoutRect(t *testing.T, n *Node, want layoutRect) {
	t.Helper()
	got := layoutRect{n.LayoutLeft(), n.LayoutTop(), n.LayoutWidth(), n.LayoutHeight()}
	if got != want {
		t.Fatalf("layout=%+v, want %+v", got, want)
	}
}

func TestLayoutFixtures(t *testing.T) {
	tests := []struct {
		name  string
		build func() (*Node, []*Node)
		want  []layoutRect
	}{
		{
			name: "row padding and gap",
			build: func() (*Node, []*Node) {
				root := sizedNode(120, 40)
				root.StyleSetFlexDirection(FlexDirectionRow)
				root.StyleSetPadding(EdgeAll, 10)
				root.StyleSetGap(GutterColumn, 10)
				children := []*Node{sizedNode(20, 20), sizedNode(20, 20)}
				for i, child := range children {
					root.InsertChild(child, uint32(i))
				}
				return root, children
			},
			want: []layoutRect{{10, 10, 20, 20}, {40, 10, 20, 20}},
		},
		{
			name: "flex grow fills remaining width",
			build: func() (*Node, []*Node) {
				root := sizedNode(100, 20)
				root.StyleSetFlexDirection(FlexDirectionRow)
				fixed, flexible := sizedNode(20, 20), NewNode()
				flexible.StyleSetFlexGrow(1)
				root.InsertChild(fixed, 0)
				root.InsertChild(flexible, 1)
				return root, []*Node{fixed, flexible}
			},
			want: []layoutRect{{0, 0, 20, 20}, {20, 0, 80, 20}},
		},
		{
			name: "percentage dimensions",
			build: func() (*Node, []*Node) {
				root := sizedNode(200, 100)
				child := NewNode()
				child.StyleSetWidthPercent(50)
				child.StyleSetHeightPercent(25)
				root.InsertChild(child, 0)
				return root, []*Node{child}
			},
			want: []layoutRect{{0, 0, 100, 25}},
		},
		{
			name: "center alignment",
			build: func() (*Node, []*Node) {
				root := sizedNode(100, 50)
				root.StyleSetFlexDirection(FlexDirectionRow)
				root.StyleSetJustifyContent(JustifyCenter)
				root.StyleSetAlignItems(AlignCenter)
				child := sizedNode(20, 10)
				root.InsertChild(child, 0)
				return root, []*Node{child}
			},
			want: []layoutRect{{40, 20, 20, 10}},
		},
		{
			name: "wrapped row",
			build: func() (*Node, []*Node) {
				root := sizedNode(50, 20)
				root.StyleSetFlexDirection(FlexDirectionRow)
				root.StyleSetFlexWrap(WrapWrap)
				children := []*Node{sizedNode(30, 10), sizedNode(30, 10)}
				for i, child := range children {
					root.InsertChild(child, uint32(i))
				}
				return root, children
			},
			want: []layoutRect{{0, 0, 30, 10}, {0, 10, 30, 10}},
		},
		{
			name: "absolute positioning",
			build: func() (*Node, []*Node) {
				root := sizedNode(100, 100)
				child := sizedNode(30, 40)
				child.StyleSetPositionType(PositionTypeAbsolute)
				child.StyleSetPosition(EdgeLeft, 10)
				child.StyleSetPosition(EdgeTop, 20)
				root.InsertChild(child, 0)
				return root, []*Node{child}
			},
			want: []layoutRect{{10, 20, 30, 40}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, children := test.build()
			root.CalculateLayout(200, 100, DirectionLTR)
			for i, child := range children {
				assertLayoutRect(t, child, test.want[i])
			}
		})
	}
}

func TestAbsolutePositionWithinStaticAncestor(t *testing.T) {
	root := sizedNode(100, 100)
	wrapper := sizedNode(80, 80)
	wrapper.StyleSetPositionType(PositionTypeStatic)
	wrapper.StyleSetMargin(EdgeLeft, 10)
	wrapper.StyleSetMargin(EdgeTop, 20)
	child := sizedNode(20, 10)
	child.StyleSetPositionType(PositionTypeAbsolute)
	child.StyleSetPosition(EdgeLeft, 5)
	child.StyleSetPosition(EdgeTop, 7)
	root.InsertChild(wrapper, 0)
	wrapper.InsertChild(child, 0)

	root.CalculateLayout(100, 100, DirectionLTR)

	left := wrapper.LayoutLeft() + child.LayoutLeft()
	top := wrapper.LayoutTop() + child.LayoutTop()
	if left != 5 || top != 7 {
		t.Fatalf("absolute child global position=(%v,%v), want (5,7); wrapper=(%v,%v), child=(%v,%v)",
			left, top, wrapper.LayoutLeft(), wrapper.LayoutTop(), child.LayoutLeft(), child.LayoutTop())
	}
}

func TestAlwaysFormsContainingBlock(t *testing.T) {
	root := sizedNode(100, 100)
	wrapper := sizedNode(80, 80)
	wrapper.StyleSetPositionType(PositionTypeStatic)
	wrapper.StyleSetMargin(EdgeLeft, 10)
	wrapper.StyleSetMargin(EdgeTop, 20)
	wrapper.SetAlwaysFormsContainingBlock(true)
	child := sizedNode(20, 10)
	child.StyleSetPositionType(PositionTypeAbsolute)
	child.StyleSetPosition(EdgeLeft, 5)
	child.StyleSetPosition(EdgeTop, 7)
	root.InsertChild(wrapper, 0)
	wrapper.InsertChild(child, 0)

	if !wrapper.AlwaysFormsContainingBlock() {
		t.Fatal("AlwaysFormsContainingBlock did not retain its setting")
	}
	root.CalculateLayout(100, 100, DirectionLTR)

	assertLayoutRect(t, child, layoutRect{5, 7, 20, 10})
}

func TestFlexShrinkAndWeightedGrow(t *testing.T) {
	t.Run("shrink", func(t *testing.T) {
		root := sizedNode(100, 20)
		root.StyleSetFlexDirection(FlexDirectionRow)
		children := []*Node{sizedNode(100, 20), sizedNode(100, 20)}
		for i, child := range children {
			child.StyleSetFlexShrink(1)
			root.InsertChild(child, uint32(i))
		}
		root.CalculateLayout(100, 20, DirectionLTR)
		assertLayoutRect(t, children[0], layoutRect{0, 0, 50, 20})
		assertLayoutRect(t, children[1], layoutRect{50, 0, 50, 20})
	})

	t.Run("weighted grow", func(t *testing.T) {
		root := sizedNode(120, 20)
		root.StyleSetFlexDirection(FlexDirectionRow)
		first, second := NewNode(), NewNode()
		first.StyleSetFlexGrow(1)
		second.StyleSetFlexGrow(2)
		root.InsertChild(first, 0)
		root.InsertChild(second, 1)
		root.CalculateLayout(120, 20, DirectionLTR)
		assertLayoutRect(t, first, layoutRect{0, 0, 40, 20})
		assertLayoutRect(t, second, layoutRect{40, 0, 80, 20})
	})
}

func TestMinAndMaxConstraintsDuringFlexing(t *testing.T) {
	t.Run("maximum growth", func(t *testing.T) {
		root := sizedNode(200, 20)
		root.StyleSetFlexDirection(FlexDirectionRow)
		child := NewNode()
		child.StyleSetFlexGrow(1)
		child.StyleSetMaxWidth(80)
		root.InsertChild(child, 0)
		root.CalculateLayout(200, 20, DirectionLTR)
		assertLayoutRect(t, child, layoutRect{0, 0, 80, 20})
	})

	t.Run("minimum shrink", func(t *testing.T) {
		root := sizedNode(50, 20)
		root.StyleSetFlexDirection(FlexDirectionRow)
		child := sizedNode(100, 20)
		child.StyleSetFlexShrink(1)
		child.StyleSetMinWidth(70)
		root.InsertChild(child, 0)
		root.CalculateLayout(50, 20, DirectionLTR)
		assertLayoutRect(t, child, layoutRect{0, 0, 70, 20})
	})
}

func TestSpaceBetweenAndAlignSelf(t *testing.T) {
	root := sizedNode(120, 50)
	root.StyleSetFlexDirection(FlexDirectionRow)
	root.StyleSetJustifyContent(JustifySpaceBetween)
	root.StyleSetAlignItems(AlignFlexStart)
	children := []*Node{sizedNode(20, 10), sizedNode(20, 10), sizedNode(20, 10)}
	children[1].StyleSetAlignSelf(AlignFlexEnd)
	for i, child := range children {
		root.InsertChild(child, uint32(i))
	}
	root.CalculateLayout(120, 50, DirectionLTR)
	assertLayoutRect(t, children[0], layoutRect{0, 0, 20, 10})
	assertLayoutRect(t, children[1], layoutRect{50, 40, 20, 10})
	assertLayoutRect(t, children[2], layoutRect{100, 0, 20, 10})
}

func TestAspectRatioAndDisplayNone(t *testing.T) {
	root := sizedNode(100, 50)
	root.StyleSetFlexDirection(FlexDirectionRow)
	withRatio := NewNode()
	withRatio.StyleSetWidth(40)
	withRatio.StyleSetAspectRatio(2)
	hidden := sizedNode(20, 20)
	hidden.StyleSetDisplay(DisplayNone)
	last := sizedNode(20, 20)
	root.InsertChild(withRatio, 0)
	root.InsertChild(hidden, 1)
	root.InsertChild(last, 2)

	root.CalculateLayout(100, 50, DirectionLTR)
	assertLayoutRect(t, withRatio, layoutRect{0, 0, 40, 20})
	assertLayoutRect(t, hidden, layoutRect{})
	assertLayoutRect(t, last, layoutRect{40, 0, 20, 20})
}

func TestRightToLeftRow(t *testing.T) {
	root := sizedNode(100, 20)
	root.StyleSetDirection(DirectionRTL)
	root.StyleSetFlexDirection(FlexDirectionRow)
	first, second := sizedNode(20, 20), sizedNode(30, 20)
	root.InsertChild(first, 0)
	root.InsertChild(second, 1)

	root.CalculateLayout(100, 20, DirectionRTL)
	assertLayoutRect(t, first, layoutRect{80, 0, 20, 20})
	assertLayoutRect(t, second, layoutRect{50, 0, 30, 20})
}

func TestMeasureCacheAndInvalidation(t *testing.T) {
	root := sizedNode(100, 20)
	leaf := NewNode()
	leaf.StyleSetAlignSelf(AlignFlexStart)
	measureCalls := 0
	leaf.SetMeasureFunc(func(_ *Node, _ float32, _ MeasureMode, _ float32, _ MeasureMode) Size {
		measureCalls++
		return Size{Width: 40, Height: 10}
	})
	root.InsertChild(leaf, 0)

	root.CalculateLayout(100, 20, DirectionLTR)
	assertLayoutRect(t, leaf, layoutRect{0, 0, 40, 10})
	initialCalls := measureCalls
	if initialCalls == 0 {
		t.Fatal("measure function was not called")
	}
	root.CalculateLayout(100, 20, DirectionLTR)
	if measureCalls != initialCalls {
		t.Fatalf("clean layout called measure %d additional times", measureCalls-initialCalls)
	}

	leaf.MarkDirty()
	root.CalculateLayout(100, 20, DirectionLTR)
	if measureCalls <= initialCalls {
		t.Fatal("MarkDirty did not invalidate the measurement cache")
	}
}
