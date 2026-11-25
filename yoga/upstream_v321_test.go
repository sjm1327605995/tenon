// Expected rectangles follow the corresponding Yoga v3.2.1 upstream fixtures.
package yoga

import "testing"

func TestUpstreamV321FlexGrowLessThanOne(t *testing.T) {
	root := sizedNode(200, 500)
	root.StyleSetPositionType(PositionTypeAbsolute)
	children := []*Node{NewNode(), NewNode(), NewNode()}
	children[0].StyleSetFlexGrow(0.2)
	children[0].StyleSetFlexBasis(40)
	children[1].StyleSetFlexGrow(0.2)
	children[2].StyleSetFlexGrow(0.4)
	for i, child := range children {
		root.InsertChild(child, uint32(i))
	}

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, children[0], layoutRect{0, 0, 200, 132})
	assertLayoutRect(t, children[1], layoutRect{0, 132, 200, 92})
	assertLayoutRect(t, children[2], layoutRect{0, 224, 200, 184})
}

func TestUpstreamV321FlexBasisOverridesMainSize(t *testing.T) {
	root := sizedNode(100, 100)
	root.StyleSetPositionType(PositionTypeAbsolute)
	children := []*Node{NewNode(), NewNode(), NewNode()}
	for i, child := range children {
		child.StyleSetFlexGrow(1)
		root.InsertChild(child, uint32(i))
	}
	children[0].StyleSetFlexBasis(50)
	children[0].StyleSetHeight(20)
	children[1].StyleSetHeight(10)
	children[2].StyleSetHeight(10)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, children[0], layoutRect{0, 0, 100, 60})
	assertLayoutRect(t, children[1], layoutRect{0, 60, 100, 20})
	assertLayoutRect(t, children[2], layoutRect{0, 80, 100, 20})
}

func TestUpstreamV321PercentagePositionAndDirection(t *testing.T) {
	root := sizedNode(400, 400)
	root.StyleSetPositionType(PositionTypeAbsolute)
	root.StyleSetFlexDirection(FlexDirectionRow)
	child := NewNode()
	child.StyleSetPositionPercent(EdgeLeft, 10)
	child.StyleSetPositionPercent(EdgeTop, 20)
	child.StyleSetWidthPercent(45)
	child.StyleSetHeightPercent(55)
	root.InsertChild(child, 0)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{40, 80, 180, 220})
	root.CalculateLayout(Undefined, Undefined, DirectionRTL)
	assertLayoutRect(t, child, layoutRect{260, 80, 180, 220})
}

func TestUpstreamV321MeasureCacheAcrossAvailableHeight(t *testing.T) {
	root := NewNode()
	child := NewNode()
	measureCalls := 0
	child.SetMeasureFunc(func(_ *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		measureCalls++
		return Size{
			Width:  measureMinDimension(width, widthMode),
			Height: measureMinDimension(height, heightMode),
		}
	})
	root.InsertChild(child, 0)

	root.CalculateLayout(100, 100, DirectionLTR)
	root.CalculateLayout(100, 50, DirectionLTR)
	if measureCalls != 1 {
		t.Fatalf("measure calls=%d, want 1", measureCalls)
	}
}

func measureMinDimension(value float32, mode MeasureMode) float32 {
	if mode == MeasureModeUndefined || mode == MeasureModeAtMost && value > 10 {
		return 10
	}
	return value
}

func TestUpstreamV321FlexBasisGrowRowAndDirection(t *testing.T) {
	root := sizedNode(100, 100)
	root.StyleSetPositionType(PositionTypeAbsolute)
	root.StyleSetFlexDirection(FlexDirectionRow)
	first, second := NewNode(), NewNode()
	first.StyleSetFlexGrow(1)
	first.StyleSetFlexBasis(50)
	second.StyleSetFlexGrow(1)
	root.InsertChild(first, 0)
	root.InsertChild(second, 1)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, first, layoutRect{0, 0, 75, 100})
	assertLayoutRect(t, second, layoutRect{75, 0, 25, 100})
	root.CalculateLayout(Undefined, Undefined, DirectionRTL)
	assertLayoutRect(t, first, layoutRect{25, 0, 75, 100})
	assertLayoutRect(t, second, layoutRect{0, 0, 25, 100})
}

func TestUpstreamV321FlexShrinkToZero(t *testing.T) {
	root := NewNode()
	root.StyleSetPositionType(PositionTypeAbsolute)
	root.StyleSetHeight(75)
	children := []*Node{sizedNode(50, 50), sizedNode(50, 50), sizedNode(50, 50)}
	children[1].StyleSetFlexShrink(1)
	for i, child := range children {
		root.InsertChild(child, uint32(i))
	}

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, root, layoutRect{0, 0, 50, 75})
	assertLayoutRect(t, children[0], layoutRect{0, 0, 50, 50})
	assertLayoutRect(t, children[1], layoutRect{0, 50, 50, 0})
	assertLayoutRect(t, children[2], layoutRect{0, 50, 50, 50})
}

func TestUpstreamV321MeasureOnceSingleFlexibleChild(t *testing.T) {
	root := sizedNode(100, 100)
	root.StyleSetFlexDirection(FlexDirectionRow)
	root.StyleSetAlignItems(AlignFlexStart)
	child := NewNode()
	child.StyleSetFlexGrow(1)
	measureCalls := 0
	child.SetMeasureFunc(func(_ *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		measureCalls++
		return Size{
			Width:  If(widthMode == MeasureModeUndefined, float32(10), width),
			Height: If(heightMode == MeasureModeUndefined, float32(10), height),
		}
	})
	root.InsertChild(child, 0)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	if measureCalls != 1 {
		t.Fatalf("measure calls=%d, want 1", measureCalls)
	}
}

func TestUpstreamV321PercentagePaddingUsesWidth(t *testing.T) {
	root := sizedNode(200, 100)
	root.StyleSetPositionType(PositionTypeAbsolute)
	child := NewNode()
	child.StyleSetFlexGrow(1)
	for _, edge := range []Edge{EdgeLeft, EdgeTop, EdgeRight, EdgeBottom} {
		child.StyleSetPaddingPercent(edge, 10)
	}
	grandchild := sizedNode(10, 10)
	root.InsertChild(child, 0)
	child.InsertChild(grandchild, 0)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{0, 0, 200, 100})
	assertLayoutRect(t, grandchild, layoutRect{20, 20, 10, 10})
	root.CalculateLayout(Undefined, Undefined, DirectionRTL)
	assertLayoutRect(t, grandchild, layoutRect{170, 20, 10, 10})
}

func TestUpstreamV321PercentageAbsolutePosition(t *testing.T) {
	root := sizedNode(200, 100)
	root.StyleSetPositionType(PositionTypeAbsolute)
	child := sizedNode(10, 10)
	child.StyleSetPositionType(PositionTypeAbsolute)
	child.StyleSetPositionPercent(EdgeLeft, 30)
	child.StyleSetPositionPercent(EdgeTop, 10)
	root.InsertChild(child, 0)

	for _, direction := range []Direction{DirectionLTR, DirectionRTL} {
		root.CalculateLayout(Undefined, Undefined, direction)
		assertLayoutRect(t, child, layoutRect{60, 10, 10, 10})
	}
}

func TestUpstreamV321AspectRatioBothDimensionsDefined(t *testing.T) {
	for _, test := range []struct {
		name      string
		direction FlexDirection
		want      layoutRect
	}{
		{name: "row", direction: FlexDirectionRow, want: layoutRect{0, 0, 100, 100}},
		{name: "column", direction: FlexDirectionColumn, want: layoutRect{0, 0, 50, 50}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := sizedNode(100, 100)
			root.StyleSetAlignItems(AlignFlexStart)
			root.StyleSetFlexDirection(test.direction)
			child := sizedNode(100, 50)
			child.StyleSetAspectRatio(1)
			root.InsertChild(child, 0)

			root.CalculateLayout(Undefined, Undefined, DirectionLTR)
			assertLayoutRect(t, child, test.want)
		})
	}
}

func TestUpstreamV321AspectRatioMinMaxConstraints(t *testing.T) {
	t.Run("maximum main dimension", func(t *testing.T) {
		root := sizedNode(100, 100)
		root.StyleSetAlignItems(AlignFlexStart)
		child := NewNode()
		child.StyleSetWidth(50)
		child.StyleSetMaxHeight(40)
		child.StyleSetAspectRatio(1)
		root.InsertChild(child, 0)
		root.CalculateLayout(Undefined, Undefined, DirectionLTR)
		assertLayoutRect(t, child, layoutRect{0, 0, 40, 40})
	})
	t.Run("maximum cross dimension", func(t *testing.T) {
		root := sizedNode(100, 100)
		root.StyleSetAlignItems(AlignFlexStart)
		child := NewNode()
		child.StyleSetHeight(50)
		child.StyleSetMaxWidth(40)
		child.StyleSetAspectRatio(1)
		root.InsertChild(child, 0)
		root.CalculateLayout(Undefined, Undefined, DirectionLTR)
		assertLayoutRect(t, child, layoutRect{0, 0, 40, 50})
	})
}

func TestUpstreamV321AlignBaselineCustomFunc(t *testing.T) {
	root := sizedNode(100, 100)
	root.StyleSetFlexDirection(FlexDirectionRow)
	root.StyleSetAlignItems(AlignBaseline)
	first := sizedNode(50, 50)
	second := sizedNode(50, 20)
	baselineChild := sizedNode(50, 20)
	baselineChild.SetBaselineFunc(func(_ *Node, _, _ float32) float32 { return 10 })
	root.InsertChild(first, 0)
	root.InsertChild(second, 1)
	second.InsertChild(baselineChild, 0)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, first, layoutRect{0, 0, 50, 50})
	assertLayoutRect(t, second, layoutRect{50, 40, 50, 20})
	assertLayoutRect(t, baselineChild, layoutRect{0, 0, 50, 20})
}

func TestUpstreamV321DisplayNoneZeroesSubtree(t *testing.T) {
	root := sizedNode(100, 100)
	root.StyleSetPositionType(PositionTypeAbsolute)
	root.StyleSetFlexDirection(FlexDirectionRow)
	first, hidden, last := NewNode(), NewNode(), NewNode()
	for _, child := range []*Node{first, hidden, last} {
		child.StyleSetFlexGrow(1)
		child.StyleSetFlexShrink(1)
		child.StyleSetFlexBasisPercent(0)
	}
	hidden.StyleSetDisplay(DisplayNone)
	hiddenChild := NewNode()
	hiddenChild.StyleSetWidth(20)
	hidden.InsertChild(hiddenChild, 0)
	root.InsertChild(first, 0)
	root.InsertChild(hidden, 1)
	root.InsertChild(last, 2)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, first, layoutRect{0, 0, 50, 100})
	assertLayoutRect(t, hidden, layoutRect{})
	assertLayoutRect(t, hiddenChild, layoutRect{})
	assertLayoutRect(t, last, layoutRect{50, 0, 50, 100})
	root.CalculateLayout(Undefined, Undefined, DirectionRTL)
	assertLayoutRect(t, first, layoutRect{50, 0, 50, 100})
	assertLayoutRect(t, last, layoutRect{0, 0, 50, 100})
}

func TestUpstreamV321DisplayContentsWithSiblings(t *testing.T) {
	root := sizedNode(100, 100)
	root.StyleSetPositionType(PositionTypeAbsolute)
	root.StyleSetFlexDirection(FlexDirectionRow)
	first, contents, last := NewNode(), NewNode(), NewNode()
	first.StyleSetHeight(30)
	last.StyleSetHeight(30)
	contents.StyleSetDisplay(DisplayContents)
	children := []*Node{first, NewNode(), NewNode(), last}
	for _, child := range children {
		child.StyleSetFlexGrow(1)
		child.StyleSetFlexShrink(1)
		child.StyleSetFlexBasisPercent(0)
	}
	children[1].StyleSetHeight(10)
	children[2].StyleSetHeight(20)
	root.InsertChild(first, 0)
	root.InsertChild(contents, 1)
	contents.InsertChild(children[1], 0)
	contents.InsertChild(children[2], 1)
	root.InsertChild(last, 2)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	for i, want := range []layoutRect{{0, 0, 25, 30}, {25, 0, 25, 10}, {50, 0, 25, 20}, {75, 0, 25, 30}} {
		assertLayoutRect(t, children[i], want)
	}
	assertLayoutRect(t, contents, layoutRect{})
	root.CalculateLayout(Undefined, Undefined, DirectionRTL)
	for i, want := range []layoutRect{{75, 0, 25, 30}, {50, 0, 25, 10}, {25, 0, 25, 20}, {0, 0, 25, 30}} {
		assertLayoutRect(t, children[i], want)
	}
}

func TestUpstreamV321DirtiedPropagationRunsOnce(t *testing.T) {
	root := sizedNode(100, 100)
	root.StyleSetAlignItems(AlignFlexStart)
	first, second := sizedNode(50, 20), sizedNode(50, 20)
	root.InsertChild(first, 0)
	root.InsertChild(second, 1)
	root.CalculateLayout(Undefined, Undefined, DirectionLTR)

	dirtiedCount := 0
	root.SetDirtiedFunc(func(_ *Node) { dirtiedCount++ })
	first.MarkDirty()
	if dirtiedCount != 1 {
		t.Fatalf("dirtied count=%d, want 1", dirtiedCount)
	}
	first.MarkDirty()
	if dirtiedCount != 1 {
		t.Fatalf("repeated dirty propagation count=%d, want 1", dirtiedCount)
	}
}

func TestUpstreamV321SetChildrenReordersAndDetaches(t *testing.T) {
	root := NewNode()
	children := []*Node{NewNode(), NewNode(), NewNode()}
	root.SetChildren(children)
	added := NewNode()
	root.SetChildren([]*Node{children[2], children[1], added})

	if children[0].GetOwner() != nil {
		t.Fatal("removed child retained its owner")
	}
	for i, child := range []*Node{children[2], children[1], added} {
		if root.GetChild(uint32(i)) != child || child.GetOwner() != root {
			t.Fatalf("child %d has inconsistent order or ownership", i)
		}
	}
}

func TestUpstreamV321RoundingValues(t *testing.T) {
	tests := []struct {
		value       float64
		ceil, floor bool
		want        float32
	}{
		{value: 6.000001, want: 6},
		{value: 6.000001, ceil: true, want: 6},
		{value: 5.999999, floor: true, want: 6},
		{value: 6.01, want: 6},
		{value: 6.01, ceil: true, want: 6.5},
		{value: 5.99, floor: true, want: 5.5},
		{value: -6.000001, want: -6},
		{value: -6.01, floor: true, want: -6.5},
		{value: -5.99, ceil: true, want: -5.5},
	}
	for _, test := range tests {
		if got := RoundValueToPixelGrid(test.value, 2, test.ceil, test.floor); got != test.want {
			t.Fatalf("RoundValueToPixelGrid(%v, 2, %t, %t)=%v, want %v", test.value, test.ceil, test.floor, got, test.want)
		}
	}
}

func TestUpstreamV321PerNodePointScaleFactor(t *testing.T) {
	rootConfig := NewConfig(nil)
	rootConfig.SetPointScaleFactor(2)
	childConfig := NewConfig(nil)
	childConfig.SetPointScaleFactor(1)
	grandchildConfig := NewConfig(nil)
	grandchildConfig.SetPointScaleFactor(0.5)
	root := sizedNode(11.5, 11.5)
	root.SetConfig(rootConfig)
	child := sizedNode(9.5, 9.5)
	child.SetConfig(childConfig)
	grandchild := sizedNode(7, 7)
	grandchild.SetConfig(grandchildConfig)
	root.InsertChild(child, 0)
	child.InsertChild(grandchild, 0)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, root, layoutRect{0, 0, 11.5, 11.5})
	assertLayoutRect(t, child, layoutRect{0, 0, 10, 10})
	assertLayoutRect(t, grandchild, layoutRect{0, 0, 8, 8})
}

func TestUpstreamV321ConsistentRoundingDuringRepeatedLayouts(t *testing.T) {
	config := NewConfig(nil)
	config.SetPointScaleFactor(2)
	root := NewNodeWithConfig(config)
	root.StyleSetMargin(EdgeTop, -1.49)
	root.StyleSetWidth(500)
	root.StyleSetHeight(500)
	container := NewNodeWithConfig(config)
	leaf := NewNodeWithConfig(config)
	leaf.SetMeasureFunc(func(_ *Node, _ float32, _ MeasureMode, _ float32, _ MeasureMode) Size {
		return Size{Width: 10, Height: 10}
	})
	root.InsertChild(container, 0)
	container.InsertChild(leaf, 0)

	for i := 0; i < 5; i++ {
		root.StyleSetMargin(EdgeLeft, float32(i+1))
		root.CalculateLayout(Undefined, Undefined, DirectionLTR)
		if got := leaf.LayoutHeight(); got != 10 {
			t.Fatalf("layout %d height=%v, want 10", i, got)
		}
	}
}

func TestUpstreamV321ScaleChangeInvalidatesChildPosition(t *testing.T) {
	config := NewConfig(nil)
	config.SetPointScaleFactor(1)
	root := NewNodeWithConfig(config)
	root.StyleSetFlexDirection(FlexDirectionRow)
	root.StyleSetWidth(50)
	root.StyleSetHeight(50)
	first, second := NewNodeWithConfig(config), NewNodeWithConfig(config)
	first.StyleSetFlexGrow(1)
	second.StyleSetFlexGrow(1)
	root.InsertChild(first, 0)
	root.InsertChild(second, 1)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, second, layoutRect{25, 0, 25, 50})
	config.SetPointScaleFactor(1.5)
	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	if got := second.LayoutLeft(); got != 25.333334 {
		t.Fatalf("left after scale change=%v, want 25.333334", got)
	}
}

func TestUpstreamV321LayoutEventCacheStats(t *testing.T) {
	ResetEventSubscribers()
	t.Cleanup(ResetEventSubscribers)
	root := NewNode()
	var passes []LayoutData
	Subscribe(EventLayoutPassEnd, func(event Event) {
		passes = append(passes, event.LayoutData)
	})
	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	root.CalculateLayout(Undefined, Undefined, DirectionLTR)

	if len(passes) != 2 {
		t.Fatalf("layout pass count=%d, want 2", len(passes))
	}
	cached := passes[1]
	if cached.Layouts() != 0 || cached.Measures() != 0 || cached.CachedLayouts() != 1 || cached.CachedMeasures() != 0 {
		t.Fatalf("cached layout stats=%+v, want layouts=0 measures=0 cachedLayouts=1 cachedMeasures=0", cached)
	}
}
