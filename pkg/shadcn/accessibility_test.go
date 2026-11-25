package shadcn

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/ui"
)

func TestCompositeControlsExposeRolesAndSelection(t *testing.T) {
	h := ui.MountDefault(ui.VStack(4,
		Checkbox(CheckboxProps{Checked: true}),
		Switch(SwitchProps{Checked: false}),
		RadioGroup(RadioGroupProps{Value: "A", Options: []string{"A", "B"}}),
	))
	defer h.Close()
	nodes := h.AccessibilitySnapshot()

	checkbox := accessibilityRole(nodes, ui.AccessibilityCheckbox, "")
	if checkbox == nil || !checkbox.HasSelected || !checkbox.Selected {
		t.Fatalf("checkbox semantic=%+v", checkbox)
	}
	switchNode := accessibilityRole(nodes, ui.AccessibilitySwitch, "")
	if switchNode == nil || !switchNode.HasSelected || switchNode.Selected {
		t.Fatalf("switch semantic=%+v", switchNode)
	}
	radioA := accessibilityRole(nodes, ui.AccessibilityRadioButton, "A")
	radioB := accessibilityRole(nodes, ui.AccessibilityRadioButton, "B")
	if radioA == nil || !radioA.Selected || radioB == nil || radioB.Selected {
		t.Fatalf("radio semantics A=%+v B=%+v", radioA, radioB)
	}
}

func TestTabsToggleAndToggleGroupSemantics(t *testing.T) {
	t.Run("tabs", func(t *testing.T) {
		h := ui.MountDefault(Tabs(TabsProps{Tabs: []string{"A", "B"}, Active: 1}))
		defer h.Close()
		nodes := h.AccessibilitySnapshot()
		a := accessibilityRole(nodes, ui.AccessibilityRadioButton, "A")
		b := accessibilityRole(nodes, ui.AccessibilityRadioButton, "B")
		if a == nil || a.Selected || b == nil || !b.Selected {
			t.Fatalf("tab semantics A=%+v B=%+v", a, b)
		}
	})
	t.Run("toggle", func(t *testing.T) {
		h := ui.MountDefault(Toggle(ToggleProps{Pressed: true}, ui.Text("Bold")))
		defer h.Close()
		node := accessibilityRole(h.AccessibilitySnapshot(), ui.AccessibilityCheckbox, "Bold")
		if node == nil || !node.Selected {
			t.Fatalf("toggle semantic=%+v", node)
		}
	})
	t.Run("toggle group", func(t *testing.T) {
		h := ui.MountDefault(ToggleGroup(ToggleGroupProps{Value: "B", Options: []string{"A", "B"}}))
		defer h.Close()
		nodes := h.AccessibilitySnapshot()
		a := accessibilityRole(nodes, ui.AccessibilityRadioButton, "A")
		b := accessibilityRole(nodes, ui.AccessibilityRadioButton, "B")
		if a == nil || a.Selected || b == nil || !b.Selected {
			t.Fatalf("toggle group semantics A=%+v B=%+v", a, b)
		}
	})
}

func TestDisabledCheckboxRemainsInAccessibilityTree(t *testing.T) {
	h := ui.MountDefault(Checkbox(CheckboxProps{Disabled: true}))
	defer h.Close()
	nodes := h.AccessibilitySnapshot()
	if len(nodes) != 1 || nodes[0].Role != ui.AccessibilityCheckbox || nodes[0].Enabled {
		t.Fatalf("disabled checkbox snapshot=%+v", nodes)
	}
}

func accessibilityRole(nodes []ui.AccessibilityNode, role ui.AccessibilityRole, label string) *ui.AccessibilityNode {
	for i := range nodes {
		if nodes[i].Role == role && (label == "" || nodes[i].Label == label) {
			return &nodes[i]
		}
	}
	return nil
}
