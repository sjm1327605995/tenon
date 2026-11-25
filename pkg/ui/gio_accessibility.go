package ui

import (
	"image"

	"gioui.org/io/event"
	"gioui.org/io/semantic"
	"gioui.org/op/clip"
)

func (p *gioPainter) PushAccessibility(node AccessibilityNode) {
	b := node.Bounds
	state := clip.Rect(image.Rect(int(b.X+0.5), int(b.Y+0.5), int(b.X+b.W+0.5), int(b.Y+b.H+0.5))).Push(p.ops)
	p.clips = append(p.clips, state.Pop)
	if node.Enabled && node.Role != AccessibilityStaticText && node.nativeTag != nil {
		event.Op(p.ops, node.nativeTag)
		p.accessibilityTags = append(p.accessibilityTags, node.nativeTag)
	}

	switch node.Role {
	case AccessibilityButton:
		semantic.Button.Add(p.ops)
		if node.Label != "" {
			semantic.LabelOp(node.Label).Add(p.ops)
		}
	case AccessibilityTextInput:
		semantic.Editor.Add(p.ops)
		if node.Value != "" {
			semantic.LabelOp(node.Value).Add(p.ops)
			if node.Label != "" {
				semantic.DescriptionOp(node.Label).Add(p.ops)
			}
		} else if node.Label != "" {
			semantic.LabelOp(node.Label).Add(p.ops)
		}
	case AccessibilityCheckbox:
		semantic.CheckBox.Add(p.ops)
		if node.Label != "" {
			semantic.LabelOp(node.Label).Add(p.ops)
		}
	case AccessibilitySwitch:
		semantic.Switch.Add(p.ops)
		if node.Label != "" {
			semantic.LabelOp(node.Label).Add(p.ops)
		}
	case AccessibilityRadioButton:
		semantic.RadioButton.Add(p.ops)
		if node.Label != "" {
			semantic.LabelOp(node.Label).Add(p.ops)
		}
	case AccessibilityStaticText:
		semantic.LabelOp(node.Label).Add(p.ops)
	}
	if node.HasSelected {
		semantic.SelectedOp(node.Selected).Add(p.ops)
	}
	semantic.EnabledOp(node.Enabled).Add(p.ops)
}

func (p *gioPainter) PopAccessibility() { p.PopClip() }
