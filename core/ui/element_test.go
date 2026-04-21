package ui

import (
	"testing"
)

func TestViewCreation(t *testing.T) {
	view := View()
	if view == nil {
		t.Fatal("View() returned nil")
	}
	if view.Node() == nil {
		t.Fatal("View().Node() returned nil")
	}
}

func TestViewWithChildren(t *testing.T) {
	child1 := View()
	child2 := View()

	view := View(child1, child2)
	if len(view.Children()) != 2 {
		t.Errorf("expected 2 children, got %d", len(view.Children()))
	}
}

func TestViewChainAPISize(t *testing.T) {
	view := View().
		Width(100).
		Height(200)

	view.CalculateLayout(375, 812, DirectionLTR)

	layout := view.GetLayout()
	if layout.Width != 100 {
		t.Errorf("expected Width 100, got %v", layout.Width)
	}
	if layout.Height != 200 {
		t.Errorf("expected Height 200, got %v", layout.Height)
	}
}

func TestViewChainAPIFlexDirection(t *testing.T) {
	view := View().
		FlexDirectionColumn()

	if view.Node().StyleGetFlexDirection() != FlexDirectionColumn {
		t.Errorf("expected FlexDirectionColumn")
	}

	view2 := View().FlexDirectionRow()
	if view2.Node().StyleGetFlexDirection() != FlexDirectionRow {
		t.Errorf("expected FlexDirectionRow")
	}
}

func TestViewChainAPIJustifyContent(t *testing.T) {
	view := View().
		JustifyContentCenter().
		AlignItemsCenter()

	if view.Node().StyleGetJustifyContent() != JustifyCenter {
		t.Errorf("expected JustifyCenter")
	}
	if view.Node().StyleGetAlignItems() != AlignCenter {
		t.Errorf("expected AlignCenter")
	}
}

func TestViewChainAPIMargin(t *testing.T) {
	view := View().
		MarginAll(10).
		MarginLeft(5).
		MarginTop(15).
		MarginRight(20).
		MarginBottom(25)

	_ = view.Node().StyleGetMargin(EdgeAll)
}

func TestViewChainAPIPadding(t *testing.T) {
	view := View().
		PaddingAll(16).
		PaddingHorizontal(8).
		PaddingVertical(12)

	_ = view.Node().StyleGetPadding(EdgeAll)
}

func TestViewChainAPIBorder(t *testing.T) {
	view := View().
		BorderAll(2).
		BorderLeft(1).
		BorderTop(1)

	if view.Node().StyleGetBorder(EdgeAll) != 2 {
		t.Errorf("expected BorderAll 2")
	}
}

func TestViewChainAPIPosition(t *testing.T) {
	view := View().
		PositionAbsolute().
		PositionLeft(10).
		PositionTop(20).
		PositionRight(30).
		PositionBottom(40)

	if view.Node().StyleGetPositionType() != PositionTypeAbsolute {
		t.Errorf("expected PositionTypeAbsolute")
	}
}

func TestViewChainAPIFlex(t *testing.T) {
	view := View().
		Flex(1).
		FlexGrow(1).
		FlexShrink(0).
		FlexBasis(100)

	if view.Node().StyleGetFlex() != 1 {
		t.Errorf("expected Flex 1")
	}
	if view.Node().StyleGetFlexGrow() != 1 {
		t.Errorf("expected FlexGrow 1")
	}
}

func TestViewChainAPIGap(t *testing.T) {
	view := View().
		GapColumn(10).
		GapRow(20)

	if view.Node().StyleGetGap(GutterColumn) != 10 {
		t.Errorf("expected GapColumn 10")
	}
	if view.Node().StyleGetGap(GutterRow) != 20 {
		t.Errorf("expected GapRow 20")
	}
}

func TestViewChainAPIDirection(t *testing.T) {
	view := View().
		DirectionLTR()

	if view.Node().StyleGetDirection() != DirectionLTR {
		t.Errorf("expected DirectionLTR")
	}
}

func TestViewChainAPISpecialProps(t *testing.T) {
	ref := &ElementRef{}

	view := View().
		Id("test-id").
		ClassName("test-class").
		Tag("div").
		Data("key1", "value1").
		Data("key2", 123).
		OnClick(func() {}).
		Ref(ref).
		BackgroundColor(0xFF0000FF).
		Opacity(0.5).
		ZIndex(10)

	style := view.Style()
	if style.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", style.ID)
	}
	if style.ClassName != "test-class" {
		t.Errorf("expected ClassName 'test-class', got '%s'", style.ClassName)
	}
	if style.Tag != "div" {
		t.Errorf("expected Tag 'div', got '%s'", style.Tag)
	}
	if style.Data["key1"] != "value1" {
		t.Errorf("expected Data['key1'] = 'value1'")
	}
	if style.Data["key2"] != 123 {
		t.Errorf("expected Data['key2'] = 123")
	}
	if style.OnClick == nil {
		t.Errorf("expected OnClick handler")
	}
	if style.Ref != ref {
		t.Errorf("expected Ref")
	}
	if style.Data["backgroundColor"] != uint32(0xFF0000FF) {
		t.Errorf("expected backgroundColor")
	}
	if style.Data["opacity"] != 0.5 {
		t.Errorf("expected opacity 0.5")
	}
	if style.Data["zIndex"] != 10 {
		t.Errorf("expected zIndex 10")
	}
}

func TestViewAddChild(t *testing.T) {
	parent := View()
	child := View()

	parent.AddChild(child)

	if len(parent.Children()) != 1 {
		t.Errorf("expected 1 child, got %d", len(parent.Children()))
	}
	if parent.GetChildCount() != 1 {
		t.Errorf("expected GetChildCount() = 1, got %d", parent.GetChildCount())
	}
}

func TestViewRemoveChild(t *testing.T) {
	parent := View()
	child := View()

	parent.AddChild(child)
	parent.RemoveChild(child)

	if len(parent.Children()) != 0 {
		t.Errorf("expected 0 children, got %d", len(parent.Children()))
	}
}

func TestHStack(t *testing.T) {
	child1 := View()
	child2 := View()

	hstack := HStack(child1, child2)

	if hstack.Node().StyleGetFlexDirection() != FlexDirectionRow {
		t.Errorf("expected HStack to have FlexDirectionRow")
	}
}

func TestVStack(t *testing.T) {
	child1 := View()
	child2 := View()

	vstack := VStack(child1, child2)

	if vstack.Node().StyleGetFlexDirection() != FlexDirectionColumn {
		t.Errorf("expected VStack to have FlexDirectionColumn")
	}
}

func TestZStack(t *testing.T) {
	zstack := ZStack()

	if zstack.Node().StyleGetPositionType() != PositionTypeAbsolute {
		t.Errorf("expected ZStack to have PositionTypeAbsolute")
	}
}

func TestTextCreation(t *testing.T) {
	text := Text("Hello World")
	if text == nil {
		t.Fatal("Text() returned nil")
	}
	if text.Text() != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", text.Text())
	}
	if text.Node() == nil {
		t.Fatal("Text().Node() returned nil")
	}
}

func TestTextChainAPI(t *testing.T) {
	text := Text("Hello").
		FontSize(24).
		Color(0xFF0000FF).
		Width(100).
		Height(50)

	style := text.Style()
	if style.Data["fontSize"] != 24.0 {
		t.Errorf("expected fontSize 24")
	}
	if style.Data["color"] != uint32(0xFF0000FF) {
		t.Errorf("expected color")
	}
}

func TestTextSetText(t *testing.T) {
	text := Text("Hello")
	text.SetText("World")

	if text.Text() != "World" {
		t.Errorf("expected 'World', got '%s'", text.Text())
	}
}


