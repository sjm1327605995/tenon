package ui

import (
	"github.com/sjm1327605995/tenon/yoga"
	"testing"
)

type CardComponent struct {
	title    string
	content  string
	children []Element
	node     *yoga.Node
}

func Card(title, content string, children ...Element) *CardComponent {
	return &CardComponent{
		title:    title,
		content:  content,
		children: children,
		node:     yoga.NewNode(),
	}
}

func (c *CardComponent) Render() Element {
	return c
}

func (c *CardComponent) Children() []Element {
	return c.children
}

func (c *CardComponent) Node() *yoga.Node {
	return c.node
}

func (c *CardComponent) Title() string {
	return c.title
}

func (c *CardComponent) Content() string {
	return c.content
}

func TestCardComponentCreation(t *testing.T) {
	card := Card("Title", "Content")
	if card == nil {
		t.Fatal("Card() returned nil")
	}
	if card.Node() == nil {
		t.Fatal("Card().Node() returned nil")
	}
}

func TestCardComponentProperties(t *testing.T) {
	card := Card("Hello", "World")

	if card.Title() != "Hello" {
		t.Errorf("expected title 'Hello', got '%s'", card.Title())
	}
	if card.Content() != "World" {
		t.Errorf("expected content 'World', got '%s'", card.Content())
	}
}

func TestCardComponentWithChildren(t *testing.T) {
	child := View().Width(100).Height(50)
	card := Card("Parent", "Content", child)

	if len(card.Children()) != 1 {
		t.Errorf("expected 1 child, got %d", len(card.Children()))
	}
}

func TestNativeAndCardMixedTree(t *testing.T) {
	childCard := Card("Child Card", "This is child")
	childCard.Node().StyleSetWidth(150)
	childCard.Node().StyleSetHeight(100)

	container := View().
		Width(375).
		Height(812).
		FlexDirectionColumn().
		JustifyContentCenter().
		AlignItemsCenter()

	container.AddChild(childCard)

	container.CalculateLayout(375, 812, DirectionLTR)

	layout := container.GetLayout()
	if layout.Width != 375 {
		t.Errorf("expected container Width 375, got %v", layout.Width)
	}

	if childCard.Node().LayoutWidth() == 0 {
		t.Errorf("expected childCard to have layout after container CalculateLayout")
	}
}

func TestCardRenderConversion(t *testing.T) {
	card := Card("Test", "Rendering")

	element := card.Render()

	if _, ok := element.(*CardComponent); !ok {
		t.Errorf("expected CardComponent after Render(), got %T", element)
	}
}

func TestCardChildrenTracked(t *testing.T) {
	childCard := Card("Child", "Content")
	card := Card("Parent", "Content", childCard)

	if len(card.Children()) != 1 {
		t.Errorf("expected 1 child, got %d", len(card.Children()))
	}
}

func TestCardInHStack(t *testing.T) {
	nativeView := View().Width(100).Height(100)
	card := Card("Card", "In HStack")
	text := Text("Label").Width(75).Height(100)

	hstack := HStack(nativeView, card, text)

	hstack.CalculateLayout(375, 100, DirectionLTR)

	if len(hstack.Children()) != 3 {
		t.Errorf("expected 3 children, got %d", len(hstack.Children()))
	}
}

func TestCardInVStack(t *testing.T) {
	text := Text("Header").Width(300).Height(40)
	card := Card("Main", "Content")
	button := View().Width(300).Height(50)

	vstack := VStack(text, card, button)

	vstack.CalculateLayout(375, 812, DirectionLTR)

	if len(vstack.Children()) != 3 {
		t.Errorf("expected 3 children, got %d", len(vstack.Children()))
	}
}

func TestCardInZStack(t *testing.T) {
	background := Card("Background", "")
	overlay := Card("Overlay", "On top")

	zstack := ZStack()
	zstack.AddChild(background)
	zstack.AddChild(overlay)

	zstack.CalculateLayout(375, 812, DirectionLTR)

	if len(zstack.Children()) != 2 {
		t.Errorf("expected 2 children, got %d", len(zstack.Children()))
	}
}

func TestBuildTreeWithCard(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)

	card := Card("Tree Card", "Testing tree build")

	root := mapper.BuildTree(card)
	if root == nil {
		t.Fatal("BuildTree returned nil")
	}

	node := mapper.GetNode(card)
	if node == nil {
		t.Fatal("GetNode returned nil for CardComponent")
	}
}

func TestUpdateCard(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)

	card := Card("Original", "Content")

	mapper.BuildTree(card)

	updatedNode := mapper.UpdateElement(card)

	if updatedNode == nil {
		t.Fatal("UpdateElement returned nil")
	}

	if !updatedNode.IsDirty() {
		t.Errorf("updated node should be dirty")
	}
}

func TestRemoveCard(t *testing.T) {
	config := DefaultRendererConfig()
	mapper := NewElementMapper(config)

	card := Card("To Remove", "Content")

	container := View().Width(375).Height(812)
	container.AddChild(card)

	mapper.BuildTree(container)

	mapper.RemoveElement(card)

	removedNode := mapper.GetNode(card)
	if removedNode != nil {
		t.Errorf("expected nil after RemoveElement")
	}
}

func TestCardWidthFromStyle(t *testing.T) {
	card := Card("Test", "Content")

	card.Node().StyleSetWidth(200)
	card.Node().StyleSetHeight(100)

	card.Node().CalculateLayout(375, 812, yoga.DirectionLTR)

	if card.Node().LayoutWidth() != 200 {
		t.Errorf("expected Width 200, got %v", card.Node().LayoutWidth())
	}
	if card.Node().LayoutHeight() != 100 {
		t.Errorf("expected Height 100, got %v", card.Node().LayoutHeight())
	}
}

func TestNestedCard(t *testing.T) {
	innerCard := Card("Inner", "Nested card")
	outerCard := Card("Outer", "Contains inner card", innerCard)

	outerCard.Node().CalculateLayout(375, 812, yoga.DirectionLTR)

	if outerCard.Node().LayoutWidth() == 0 {
		t.Errorf("expected outerCard to have layout")
	}
}

func TestMixedCardAndView(t *testing.T) {
	card := Card("Header", "Welcome")
	view := View().Width(100).Height(50)

	container := View().
		Width(375).
		Height(812).
		FlexDirectionColumn()

	container.AddChild(card)
	container.AddChild(view)

	container.CalculateLayout(375, 812, DirectionLTR)

	if len(container.Children()) != 2 {
		t.Errorf("expected 2 children, got %d", len(container.Children()))
	}
}


