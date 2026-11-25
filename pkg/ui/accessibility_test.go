package ui

import (
	"testing"

	"gioui.org/f32"
	gioinput "gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/op"
)

func accessibilityTestApp() *Node {
	return Div(
		Button(OnClick(func() {}), Text("Save")),
		Input(Placeholder("Name"), Value("Tenon")),
		Text("Ready"),
		Button(OnClick(func() {}), AccessibilityLabel("Close"), Disabled(true), Text("X")),
	)
}

func TestAccessibilitySnapshot(t *testing.T) {
	h := MountDefault(accessibilityTestApp())
	defer h.Close()
	nodes := h.AccessibilitySnapshot()
	if len(nodes) != 4 {
		t.Fatalf("semantic node count=%d, want 4: %+v", len(nodes), nodes)
	}
	if nodes[0].Role != AccessibilityButton || nodes[0].Label != "Save" || nodes[0].Parent != -1 {
		t.Fatalf("button semantic=%+v", nodes[0])
	}
	if nodes[1].Role != AccessibilityTextInput || nodes[1].Label != "Name" || nodes[1].Value != "Tenon" {
		t.Fatalf("input semantic=%+v", nodes[1])
	}
	if nodes[2].Role != AccessibilityStaticText || nodes[2].Label != "Ready" || nodes[2].Parent != -1 {
		t.Fatalf("static text semantic=%+v", nodes[2])
	}
	if nodes[3].Label != "Close" || nodes[3].Enabled {
		t.Fatalf("disabled button semantic=%+v", nodes[3])
	}

	h.Root().ByPlaceholder("Name").Focus()
	nodes = h.AccessibilitySnapshot()
	if !nodes[1].Focused || nodes[0].Focused {
		t.Fatalf("focus semantic is incorrect: %+v", nodes)
	}
	h.Root().ByPlaceholder("Name").SetValue("Updated")
	if got := h.AccessibilitySnapshot()[1].Value; got != "Updated" {
		t.Fatalf("input semantic value=%q, want Updated", got)
	}
	h.Close()
	if h.AccessibilitySnapshot() != nil {
		t.Fatal("closed harness retained accessibility nodes")
	}
}

func TestGioAccessibilityTree(t *testing.T) {
	h := MountDefault(accessibilityTestApp())
	defer h.Close()
	ops := new(op.Ops)
	painter := newGioPainter(ops, h.g.w, h.g.h)
	paint(painter, h.g.rootRN)
	var router gioinput.Router
	router.Event(pointer.Filter{Target: gioTag, Kinds: pointer.Press | pointer.Release})
	for _, tag := range painter.accessibilityTags {
		router.Event(pointer.Filter{Target: tag, Kinds: pointer.Press | pointer.Release})
	}
	router.Frame(ops)
	nodes := router.AppendSemantics(nil)

	button := findSemanticByClass(nodes, semantic.Button)
	if button == nil || button.Desc.Label != "Save" || button.Desc.Disabled || button.Desc.Gestures&gioinput.ClickGesture == 0 {
		t.Fatalf("Gio button semantic=%+v", button)
	}
	editor := findSemanticByClass(nodes, semantic.Editor)
	if editor == nil || editor.Desc.Label != "Tenon" || editor.Desc.Description != "Name" || editor.Desc.Gestures&gioinput.ClickGesture == 0 {
		t.Fatalf("Gio editor semantic=%+v", editor)
	}
	if label := findSemanticByLabel(nodes, "Ready"); label == nil || label.Desc.Gestures != 0 {
		t.Fatalf("Gio static text semantic=%+v, want no click gesture", label)
	}
	if close := findSemanticByLabel(nodes, "Close"); close == nil || !close.Desc.Disabled || close.Desc.Gestures != 0 {
		t.Fatalf("Gio disabled button semantic=%+v, want disabled without click gesture", close)
	}
}

func TestCompositeGioAccessibilityRoles(t *testing.T) {
	for _, test := range []struct {
		role     AccessibilityRole
		class    semantic.ClassOp
		selected bool
	}{
		{role: AccessibilityCheckbox, class: semantic.CheckBox, selected: true},
		{role: AccessibilitySwitch, class: semantic.Switch, selected: false},
		{role: AccessibilityRadioButton, class: semantic.RadioButton, selected: true},
	} {
		t.Run(test.class.String(), func(t *testing.T) {
			h := MountDefault(Div(
				AccessibilityRoleOf(test.role), AccessibilitySelected(test.selected), AccessibilityLabel("choice"),
			))
			defer h.Close()
			snapshot := h.AccessibilitySnapshot()
			if len(snapshot) != 1 || snapshot[0].Role != test.role || !snapshot[0].HasSelected || snapshot[0].Selected != test.selected {
				t.Fatalf("composite snapshot=%+v", snapshot)
			}
			ops := new(op.Ops)
			painter := newGioPainter(ops, h.g.w, h.g.h)
			paint(painter, h.g.rootRN)
			var router gioinput.Router
			registerSemanticPointerFilters(&router, painter)
			router.Frame(ops)
			semanticNode := findSemanticByClass(router.AppendSemantics(nil), test.class)
			if semanticNode == nil || semanticNode.Desc.Label != "choice" || semanticNode.Desc.Selected != test.selected {
				t.Fatalf("Gio composite semantic=%+v", semanticNode)
			}
		})
	}
}

func TestDisabledCompositeRoleDoesNotRequireClickHandler(t *testing.T) {
	h := MountDefault(Div(AccessibilityRoleOf(AccessibilityCheckbox), AccessibilityLabel("disabled choice"), Disabled(true)))
	defer h.Close()
	nodes := h.AccessibilitySnapshot()
	if len(nodes) != 1 || nodes[0].Enabled {
		t.Fatalf("disabled composite snapshot=%+v", nodes)
	}
}

func findSemanticByClass(nodes []gioinput.SemanticNode, class semantic.ClassOp) *gioinput.SemanticNode {
	for i := range nodes {
		if nodes[i].Desc.Class == class {
			return &nodes[i]
		}
		if found := findSemanticByClass(nodes[i].Children, class); found != nil {
			return found
		}
	}
	return nil
}

func findSemanticByLabel(nodes []gioinput.SemanticNode, label string) *gioinput.SemanticNode {
	for i := range nodes {
		if nodes[i].Desc.Label == label {
			return &nodes[i]
		}
		if found := findSemanticByLabel(nodes[i].Children, label); found != nil {
			return found
		}
	}
	return nil
}

func registerSemanticPointerFilters(router *gioinput.Router, painter *gioPainter) {
	for _, tag := range painter.accessibilityTags {
		router.Event(pointer.Filter{Target: tag, Kinds: pointer.Press | pointer.Release})
	}
}

func TestGioSemanticHitTesting(t *testing.T) {
	h := Mount(Div(
		Button(Style(Width(100), Height(30)), OnClick(func() {}), Text("Save")),
		Input(Style(Width(100), Height(30)), Placeholder("Name")),
	), 200, 100)
	defer h.Close()
	ops := new(op.Ops)
	painter := newGioPainter(ops, h.g.w, h.g.h)
	paint(painter, h.g.rootRN)
	var router gioinput.Router
	registerSemanticPointerFilters(&router, painter)
	router.Frame(ops)
	nodes := router.AppendSemantics(nil)

	for _, test := range []struct {
		point f32.Point
		class semantic.ClassOp
	}{
		{point: f32.Pt(20, 15), class: semantic.Button},
		{point: f32.Pt(20, 45), class: semantic.Editor},
	} {
		id, ok := router.SemanticAt(test.point)
		if !ok {
			t.Fatalf("no semantic node at %v", test.point)
		}
		node := findSemanticByID(nodes, id)
		if node == nil || node.Desc.Class != test.class {
			t.Fatalf("semantic at %v = %+v, want %v", test.point, node, test.class)
		}
	}
}

func findSemanticByID(nodes []gioinput.SemanticNode, id gioinput.SemanticID) *gioinput.SemanticNode {
	for i := range nodes {
		if nodes[i].ID == id {
			return &nodes[i]
		}
		if found := findSemanticByID(nodes[i].Children, id); found != nil {
			return found
		}
	}
	return nil
}

func TestGioSemanticHitTestingUsesTransformedBounds(t *testing.T) {
	h := Mount(Button(
		Style(Width(100), Height(30), TranslateXY(40, 20)),
		OnClick(func() {}), Text("Move"),
	), 300, 100)
	defer h.Close()
	if got := h.Root().VisualBounds().X; got != 40 {
		t.Fatalf("visual button x=%v, want 40", got)
	}
	ops := new(op.Ops)
	painter := newGioPainter(ops, h.g.w, h.g.h)
	paint(painter, h.g.rootRN)
	var router gioinput.Router
	registerSemanticPointerFilters(&router, painter)
	router.Frame(ops)
	nodes := router.AppendSemantics(nil)

	id, ok := router.SemanticAt(f32.Pt(70, 35))
	if !ok {
		t.Fatal("no semantic node at the transformed visual position")
	}
	if node := findSemanticByID(nodes, id); node == nil || node.Desc.Class != semantic.Button {
		t.Fatalf("transformed semantic node=%+v, want Button", node)
	}
	id, ok = router.SemanticAt(f32.Pt(20, 15))
	if ok {
		if node := findSemanticByID(nodes, id); node != nil && node.Desc.Class == semantic.Button {
			t.Fatal("button semantic remained at its untransformed layout position")
		}
	}
}

func TestGioSemanticHitTestingUsesTopmostPortal(t *testing.T) {
	h := MountDefault(Div(
		Button(Style(Width(100), Height(40)), OnClick(func() {}), Text("Background")),
		Portal(Button(Style(Width(100), Height(40)), OnClick(func() {}), Text("Overlay"))),
	))
	defer h.Close()
	ops := new(op.Ops)
	painter := newGioPainter(ops, h.g.w, h.g.h)
	paint(painter, h.g.rootRN)
	for _, portal := range h.g.portals {
		paint(painter, portal.overlayRoot)
	}
	var router gioinput.Router
	registerSemanticPointerFilters(&router, painter)
	router.Frame(ops)
	nodes := router.AppendSemantics(nil)
	id, ok := router.SemanticAt(f32.Pt(20, 20))
	if !ok {
		t.Fatal("no semantic node at the portal button")
	}
	node := findSemanticByID(nodes, id)
	if node == nil || node.Desc.Class != semantic.Button || node.Desc.Label != "Overlay" {
		t.Fatalf("portal semantic=%+v, want Overlay button", node)
	}
}
