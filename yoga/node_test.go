package yoga

import "testing"

func TestNodeResetRestoresDefaultsAndKeepsConfig(t *testing.T) {
	config := NewConfig(nil)
	config.SetPointScaleFactor(2)
	n := sizedNode(100, 50)
	n.SetConfig(config)
	n.SetContext("context")
	n.SetIsReferenceBaseline(true)

	n.Reset()

	if n.GetConfig() != config {
		t.Fatal("Reset replaced the node's config")
	}
	if n.GetContext() != nil {
		t.Fatal("Reset retained the old context")
	}
	if !IsUndefined(n.StyleGetWidth()) || !IsUndefined(n.StyleGetHeight()) {
		t.Fatal("Reset retained the old dimensions")
	}
	if n.IsReferenceBaseline() {
		t.Fatal("Reset retained the reference baseline flag")
	}
}

func TestSetIsReferenceBaselineChangesValueAndDirtiesNode(t *testing.T) {
	n := NewNode()
	n.SetIsReferenceBaseline(true)
	if !n.IsReferenceBaseline() || !n.IsDirty() {
		t.Fatal("enabling reference baseline did not update and dirty the node")
	}

	n.setDirty(false)
	n.SetIsReferenceBaseline(true)
	if n.IsDirty() {
		t.Fatal("setting the same baseline value dirtied the node")
	}
	n.SetIsReferenceBaseline(false)
	if n.IsReferenceBaseline() || !n.IsDirty() {
		t.Fatal("disabling reference baseline did not update and dirty the node")
	}
}

func TestCloneChildrenIfNeededInstallsClone(t *testing.T) {
	originalParent := NewNode()
	shared := sizedNode(20, 20)
	originalParent.InsertChild(shared, 0)

	newParent := NewNode()
	newParent.children_ = []*Node{shared}
	newParent.cloneChildrenIfNeeded()

	clone := newParent.GetChild(0)
	if clone == shared {
		t.Fatal("shared child was not replaced by a clone")
	}
	if clone.GetOwner() != newParent {
		t.Fatal("clone does not belong to the new parent")
	}
	if shared.GetOwner() != originalParent {
		t.Fatal("cloning changed the original child's owner")
	}
}

func TestCloneHasIndependentChildList(t *testing.T) {
	child := sizedNode(20, 20)
	original := NewNode()
	original.InsertChild(child, 0)
	clone := original.Clone()

	if clone.GetOwner() != nil || clone.GetChild(0) != child {
		t.Fatal("Clone did not produce a detached copy sharing the original child")
	}
	clone.RemoveChild(child)
	if clone.GetChildCount() != 0 || original.GetChildCount() != 1 {
		t.Fatal("editing the clone's child list changed the original")
	}
	if child.GetOwner() != original {
		t.Fatal("removing a shared child from the clone detached the original child")
	}
}

func TestStaticPositionIgnoresInsets(t *testing.T) {
	root := sizedNode(100, 100)
	child := sizedNode(20, 20)
	child.StyleSetPositionType(PositionTypeStatic)
	child.StyleSetPosition(EdgeLeft, 10)
	child.StyleSetPosition(EdgeTop, 20)
	root.InsertChild(child, 0)

	root.CalculateLayout(100, 100, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{0, 0, 20, 20})

	child.StyleSetPositionType(PositionTypeRelative)
	root.CalculateLayout(100, 100, DirectionLTR)
	assertLayoutRect(t, child, layoutRect{10, 20, 20, 20})
}

func TestRemoveChildInvalidatesLayout(t *testing.T) {
	root := sizedNode(100, 20)
	root.StyleSetFlexDirection(FlexDirectionRow)
	first, second := sizedNode(20, 20), sizedNode(20, 20)
	root.InsertChild(first, 0)
	root.InsertChild(second, 1)
	root.CalculateLayout(100, 20, DirectionLTR)
	assertLayoutRect(t, second, layoutRect{20, 0, 20, 20})

	root.RemoveChild(first)
	if !root.IsDirty() {
		t.Fatal("removing a child did not dirty its parent")
	}
	if first.GetOwner() != nil || first.LayoutWidth() != 0 {
		t.Fatal("removed child retained its owner or stale layout")
	}
	root.CalculateLayout(100, 20, DirectionLTR)
	assertLayoutRect(t, second, layoutRect{0, 0, 20, 20})
}

func TestSwapChildUpdatesOwnershipAndLayout(t *testing.T) {
	root := sizedNode(100, 20)
	oldChild := sizedNode(20, 20)
	root.InsertChild(oldChild, 0)
	root.CalculateLayout(100, 20, DirectionLTR)
	newChild := sizedNode(30, 20)

	root.SwapChild(newChild, 0)

	if oldChild.GetOwner() != nil || newChild.GetOwner() != root {
		t.Fatal("SwapChild left inconsistent ownership")
	}
	if !root.IsDirty() {
		t.Fatal("SwapChild did not dirty its parent")
	}
	root.CalculateLayout(100, 20, DirectionLTR)
	assertLayoutRect(t, newChild, layoutRect{0, 0, 30, 20})
}

func TestSetChildrenReplacesAndReordersNodes(t *testing.T) {
	root := sizedNode(100, 20)
	root.StyleSetFlexDirection(FlexDirectionRow)
	oldChild := sizedNode(20, 20)
	root.InsertChild(oldChild, 0)
	first, second := sizedNode(20, 20), sizedNode(30, 20)

	children := []*Node{first, second}
	root.SetChildren(children)
	children[0] = nil
	if oldChild.GetOwner() != nil || root.GetChild(0) != first {
		t.Fatal("SetChildren did not replace the child list safely")
	}

	root.SetChildren([]*Node{second, first})
	root.CalculateLayout(100, 20, DirectionLTR)
	assertLayoutRect(t, second, layoutRect{0, 0, 30, 20})
	assertLayoutRect(t, first, layoutRect{30, 0, 20, 20})
}

func TestDefaultConfigUsesStandardPointScaleFactor(t *testing.T) {
	if got := NewNode().GetConfig().GetPointScaleFactor(); got != 1 {
		t.Fatalf("default point scale factor=%v, want 1", got)
	}
}
