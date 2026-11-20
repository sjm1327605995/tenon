package yoga

import "testing"

func TestSizingModeConversions(t *testing.T) {
	modes := []MeasureMode{MeasureModeUndefined, MeasureModeExactly, MeasureModeAtMost}
	for _, mode := range modes {
		if got := SizingModeFromMeasureMode(mode).MeasureMode(); got != mode {
			t.Fatalf("round trip for %v produced %v", mode, got)
		}
	}
}

func TestDisplayContentsFlattensLayoutChildren(t *testing.T) {
	root := sizedNode(100, 20)
	root.StyleSetFlexDirection(FlexDirectionRow)
	contents := NewNode()
	contents.StyleSetDisplay(DisplayContents)
	contents.StyleSetPadding(EdgeAll, 50)
	first := sizedNode(20, 20)
	nested := NewNode()
	nested.StyleSetDisplay(DisplayContents)
	second := sizedNode(30, 20)
	last := sizedNode(10, 20)
	root.InsertChild(contents, 0)
	contents.InsertChild(first, 0)
	contents.InsertChild(nested, 1)
	nested.InsertChild(second, 0)
	root.InsertChild(last, 1)

	if got := root.getLayoutChildCount(); got != 3 {
		t.Fatalf("layout child count=%d, want 3", got)
	}
	root.CalculateLayout(100, 20, DirectionLTR)
	assertLayoutRect(t, contents, layoutRect{})
	assertLayoutRect(t, nested, layoutRect{})
	assertLayoutRect(t, first, layoutRect{0, 0, 20, 20})
	assertLayoutRect(t, second, layoutRect{20, 0, 30, 20})
	assertLayoutRect(t, last, layoutRect{50, 0, 10, 20})
}

func TestContentBoxDimensionsIncludePaddingAndBorder(t *testing.T) {
	newRoot := func(sizing BoxSizing) *Node {
		root := sizedNode(100, 100)
		root.StyleSetBoxSizing(sizing)
		root.StyleSetPadding(EdgeAll, 10)
		root.StyleSetBorder(EdgeAll, 5)
		return root
	}
	borderBox := newRoot(BoxSizingBorderBox)
	borderBox.CalculateLayout(200, 200, DirectionLTR)
	assertLayoutRect(t, borderBox, layoutRect{0, 0, 100, 100})

	contentBox := newRoot(BoxSizingContentBox)
	contentBox.CalculateLayout(200, 200, DirectionLTR)
	assertLayoutRect(t, contentBox, layoutRect{0, 0, 130, 130})
	if contentBox.StyleGetBoxSizing() != BoxSizingContentBox {
		t.Fatal("box sizing setting was not retained")
	}
}

func TestAbsoluteInsetsHonorPaddingRulesAndAuto(t *testing.T) {
	config := NewConfig(nil)
	root := NewNodeWithConfig(config)
	root.StyleSetWidth(100)
	root.StyleSetHeight(100)
	root.StyleSetPadding(EdgeAll, 10)
	root.StyleSetBorder(EdgeAll, 5)
	child := NewNodeWithConfig(config)
	child.StyleSetWidth(20)
	child.StyleSetHeight(20)
	child.StyleSetPositionType(PositionTypeAbsolute)
	root.InsertChild(child, 0)

	root.CalculateLayout(100, 100, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{15, 15, 20, 20})

	config.AddErrata(ErrataAbsolutePositionWithoutInsetsExcludesPadding)
	root.CalculateLayout(100, 100, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{5, 5, 20, 20})

	child.StyleSetPosition(EdgeLeft, 2)
	child.StyleSetPosition(EdgeTop, 3)
	root.CalculateLayout(100, 100, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{7, 8, 20, 20})

	child.StyleSetPositionAuto(EdgeLeft)
	child.StyleSetPositionAuto(EdgeTop)
	root.CalculateLayout(100, 100, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{5, 5, 20, 20})
}

func TestPercentageGapUsesContainerAxis(t *testing.T) {
	root := sizedNode(200, 20)
	root.StyleSetFlexDirection(FlexDirectionRow)
	root.StyleSetGapPercent(GutterColumn, 10)
	first, second := sizedNode(20, 20), sizedNode(20, 20)
	root.InsertChild(first, 0)
	root.InsertChild(second, 1)
	root.CalculateLayout(200, 20, DirectionLTR)
	assertLayoutRect(t, first, layoutRect{0, 0, 20, 20})
	assertLayoutRect(t, second, layoutRect{40, 0, 20, 20})
}

func TestConfigVersionInvalidatesCachedLayout(t *testing.T) {
	config := NewConfig(nil)
	root := NewNodeWithConfig(config)
	root.StyleSetWidth(10.25)
	root.StyleSetHeight(10)
	root.CalculateLayout(20, 20, DirectionLTR)
	if root.LayoutWidth() != 10 {
		t.Fatalf("initial rounded width=%v, want 10", root.LayoutWidth())
	}

	version := config.GetVersion()
	config.SetPointScaleFactor(2)
	if config.GetVersion() == version {
		t.Fatal("config change did not increment its version")
	}
	root.CalculateLayout(20, 20, DirectionLTR)
	if root.LayoutWidth() != 10.5 {
		t.Fatalf("width after scale change=%v, want 10.5", root.LayoutWidth())
	}
	version = config.GetVersion()
	config.SetPointScaleFactor(2)
	if config.GetVersion() != version {
		t.Fatal("unchanged config value incremented its version")
	}
}

func TestDegenerateAspectRatioActsAsAuto(t *testing.T) {
	n := NewNode()
	n.StyleSetAspectRatio(2)
	n.StyleSetAspectRatio(0)
	if !IsUndefined(n.StyleGetAspectRatio()) {
		t.Fatal("zero aspect ratio should be treated as auto")
	}
}

func TestInvalidMeasureResultIsClamped(t *testing.T) {
	warnings := 0
	config := NewConfig(func(_ *Config, _ *Node, level LogLevel, _ string, _ ...any) int {
		if level == LogLevelWarn {
			warnings++
		}
		return 0
	})
	root := sizedNode(100, 20)
	leaf := NewNodeWithConfig(config)
	leaf.StyleSetAlignSelf(AlignFlexStart)
	leaf.SetMeasureFunc(func(_ *Node, _ float32, _ MeasureMode, _ float32, _ MeasureMode) Size {
		return Size{Width: -1, Height: Undefined}
	})
	root.InsertChild(leaf, 0)

	root.CalculateLayout(100, 20, DirectionLTR)
	assertLayoutRect(t, leaf, layoutRect{})
	if warnings == 0 {
		t.Fatal("invalid measure result did not produce a warning")
	}
}

func TestDisplayContentsWithFlexibleChildren(t *testing.T) {
	root := sizedNode(100, 100)
	root.StyleSetPositionType(PositionTypeAbsolute)
	root.StyleSetFlexDirection(FlexDirectionRow)
	contents := NewNode()
	contents.StyleSetDisplay(DisplayContents)
	root.InsertChild(contents, 0)
	children := []*Node{NewNode(), NewNode()}
	for i, child := range children {
		child.StyleSetFlexGrow(1)
		child.StyleSetFlexShrink(1)
		child.StyleSetFlexBasisPercent(0)
		child.StyleSetHeight(float32(10 * (i + 1)))
		contents.InsertChild(child, uint32(i))
	}

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, contents, layoutRect{})
	assertLayoutRect(t, children[0], layoutRect{0, 0, 50, 10})
	assertLayoutRect(t, children[1], layoutRect{50, 0, 50, 20})

	root.CalculateLayout(Undefined, Undefined, DirectionRTL)
	assertLayoutRect(t, children[0], layoutRect{50, 0, 50, 10})
	assertLayoutRect(t, children[1], layoutRect{0, 0, 50, 20})
}

func TestContentBoxPercentageDimensions(t *testing.T) {
	root := sizedNode(100, 100)
	root.StyleSetPositionType(PositionTypeAbsolute)
	child := NewNode()
	child.StyleSetPadding(EdgeAll, 4)
	child.StyleSetBorder(EdgeAll, 16)
	child.StyleSetWidthPercent(50)
	child.StyleSetHeightPercent(25)
	child.StyleSetBoxSizing(BoxSizingContentBox)
	root.InsertChild(child, 0)

	root.CalculateLayout(Undefined, Undefined, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{0, 0, 90, 65})
}

func TestLegacyStaticPositionErrata(t *testing.T) {
	config := NewConfig(nil)
	config.AddErrata(ErrataPositionStaticBehavesLikeRelative)
	root := NewNodeWithConfig(config)
	root.StyleSetWidth(100)
	root.StyleSetHeight(100)
	child := NewNodeWithConfig(config)
	child.StyleSetWidth(20)
	child.StyleSetHeight(20)
	child.StyleSetPositionType(PositionTypeStatic)
	child.StyleSetPosition(EdgeLeft, 10)
	root.InsertChild(child, 0)

	root.CalculateLayout(100, 100, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{10, 0, 20, 20})
}

func TestDisplayChangeUpdatesContentsCount(t *testing.T) {
	root := NewNode()
	child := NewNode()
	grandchild := sizedNode(20, 20)
	root.InsertChild(child, 0)
	child.InsertChild(grandchild, 0)
	if root.contentsChildrenCount_ != 0 {
		t.Fatal("plain child was counted as display:contents")
	}
	child.StyleSetDisplay(DisplayContents)
	if root.contentsChildrenCount_ != 1 || root.getLayoutChildren()[0] != grandchild {
		t.Fatal("display:contents change did not update flattened children")
	}
	child.StyleSetDisplay(DisplayFlex)
	if root.contentsChildrenCount_ != 0 || root.getLayoutChildren()[0] != child {
		t.Fatal("display:flex change did not restore direct children")
	}
}
