package ui

import (
	"testing"
	"time"
)

func TestRendererConfig(t *testing.T) {
	config := DefaultRendererConfig()
	if config.Width != 375 {
		t.Errorf("expected Width 375, got %v", config.Width)
	}
	if config.Height != 812 {
		t.Errorf("expected Height 812, got %v", config.Height)
	}
	if config.Direction != DirectionLTR {
		t.Errorf("expected DirectionLTR")
	}
}

func TestRendererCreation(t *testing.T) {
	config := DefaultRendererConfig()
	renderer := NewRenderer(config)
	if renderer == nil {
		t.Fatal("NewRenderer returned nil")
	}
	defer renderer.Stop()
}

func TestElementMapperCreation(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)
	if mapper == nil {
		t.Fatal("NewElementMapper returned nil")
	}
}

func TestBuildTree(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)

	element := View().
		Width(375).
		Height(812).
		FlexDirectionColumn()

	rootNode := mapper.BuildTree(element)
	if rootNode == nil {
		t.Fatal("BuildTree returned nil rootNode")
	}
	if rootNode.yogaNode == nil {
		t.Fatal("rootNode.yogaNode is nil")
	}
}

func TestBuildTreeWithChildren(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)

	element := VStack(
		Text("Header"),
		View().
			Width(100).
			Height(100),
		Text("Footer"),
	)

	rootNode := mapper.BuildTree(element)
	if rootNode == nil {
		t.Fatal("BuildTree returned nil")
	}
	if len(rootNode.children) != 3 {
		t.Errorf("expected 3 children, got %d", len(rootNode.children))
	}
}

func TestCalculateLayout(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)

	element := View().
		Width(375).
		Height(812).
		FlexDirectionColumn()

	mapper.BuildTree(element)
	mapper.CalculateLayout()

	layout := mapper.GetRootNode().GetLayout()
	if layout.Width != 375 {
		t.Errorf("expected Width 375, got %v", layout.Width)
	}
	if layout.Height != 812 {
		t.Errorf("expected Height 812, got %v", layout.Height)
	}
}

func TestRenderNodeDirty(t *testing.T) {
	node := NewRenderNode(View())
	if node.IsDirty() {
		t.Errorf("new node should not be dirty")
	}

	node.MarkDirty()
	if !node.IsDirty() {
		t.Errorf("node should be dirty after MarkDirty()")
	}

	node.ClearDirty()
	if node.IsDirty() {
		t.Errorf("node should not be dirty after ClearDirty()")
	}
}

func TestGetNode(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)

	element := View()
	mapper.BuildTree(element)

	node := mapper.GetNode(element)
	if node == nil {
		t.Fatal("GetNode returned nil")
	}
}

func TestUpdateElement(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)

	element := View()
	mapper.BuildTree(element)

	node := mapper.UpdateElement(element)
	if node == nil {
		t.Fatal("UpdateElement returned nil")
	}
	if !node.IsDirty() {
		t.Errorf("updated node should be dirty")
	}
}

func TestStartStopRenderer(t *testing.T) {
	config := DefaultRendererConfig()
	renderer := NewRenderer(config)

	element := View().Width(100).Height(100)
	renderer.SetRoot(element)

	renderer.Start()
	time.Sleep(50 * time.Millisecond)

	if !renderer.isRunning.Load() {
		t.Errorf("renderer should be running after Start()")
	}

	renderer.Stop()
	time.Sleep(50 * time.Millisecond)

	if renderer.isRunning.Load() {
		t.Errorf("renderer should not be running after Stop()")
	}
}

func TestSetRoot(t *testing.T) {
	config := DefaultRendererConfig()
	renderer := NewRenderer(config)

	element := View().Width(375).Height(812)
	renderer.SetRoot(element)
	renderer.mapper.CalculateLayout()

	layout := renderer.GetLayout()
	if layout.Width != 375 {
		t.Errorf("expected Width 375, got %v", layout.Width)
	}
}

func TestRegisterHandler(t *testing.T) {
	config := DefaultRendererConfig()
	renderer := NewRenderer(config)

	handler := func(e RenderEventData) {}

	renderer.RegisterHandler(RenderEventLayoutChange, handler)
	if len(renderer.handlers[RenderEventLayoutChange]) != 1 {
		t.Errorf("expected 1 handler")
	}

	renderer.Stop()
}

func TestLayoutChangeDetection(t *testing.T) {
	config := DefaultRendererConfig()
	renderer := NewRenderer(config)

	element := View().Width(375).Height(812)
	renderer.SetRoot(element)
	renderer.Start()
	renderer.mapper.CalculateLayout()

	time.Sleep(50 * time.Millisecond)

	renderer.Stop()
}

func TestPrintTree(t *testing.T) {
	config := DefaultRendererConfig()
	renderer := NewRenderer(config)

	element := VStack(
		Text("Header"),
		View().
			Width(300).
			Height(200),
	)

	renderer.SetRoot(element)
	renderer.PrintTree()
	renderer.Stop()
}

func TestStartLayoutMonitorLoop(t *testing.T) {
	element := VStack(
		View().Width(100).Height(100),
		View().Width(100).Height(100),
	)

	renderer := StartLayoutMonitorLoop(element, 50*time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	renderer.Stop()
}

func TestCollectLayouts(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)

	element := VStack(
		Text("Header"),
		View().Width(100).Height(100),
	)

	mapper.BuildTree(element)
	mapper.CalculateLayout()

	rootNode := mapper.GetRootNode()
	layouts := collectLayouts(rootNode)

	if len(layouts) == 0 {
		t.Errorf("expected layouts to be collected")
	}
}

func TestCompareLayouts(t *testing.T) {
	element := View().
		Width(375).
		Height(200)

	element.CalculateLayout(375, 812, DirectionLTR)
	oldLayout := element.GetLayout()

	element.Width(400)
	element.Height(300)
	element.CalculateLayout(375, 812, DirectionLTR)
	newLayout := element.GetLayout()

	if oldLayout.Width == newLayout.Width && oldLayout.Height == newLayout.Height {
		t.Errorf("expected layout changes to be detected")
	}
}

func TestElementTypeName(t *testing.T) {
	viewNode := NewRenderNode(View())
	if elementTypeName(viewNode) != "View" {
		t.Errorf("expected 'View'")
	}

	textNode := NewRenderNode(Text("Hello"))
	if elementTypeName(textNode) != "Text(\"Hello\")" {
		t.Errorf("expected 'Text(\"Hello\")'")
	}
}


