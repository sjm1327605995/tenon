package decl

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
	values "github.com/sjm1327605995/tenon/pkg/value"
)

// Text describes plain text. A zero FontSize or Color preserves inherited styling.
type Text struct {
	Key      string
	Value    string
	FontSize float32
	Color    ui.Color
	Bold     bool
	Style    []ui.StyleOpt
	Bind     values.Value[string]
}

// Build creates a text node with the requested text styles.
func (t Text) Build() *ui.Node {
	opts := make([]ui.StyleOpt, 0, 3+len(t.Style))
	if t.FontSize != 0 {
		opts = append(opts, ui.FontSize(t.FontSize))
	}
	if t.Color != (ui.Color{}) {
		opts = append(opts, ui.TextColor(t.Color))
	}
	if t.Bold {
		opts = append(opts, ui.Bold)
	}
	var content any = t.Value
	if t.Bind != nil {
		content = t.Bind
	}
	return finish(ui.Text(content), t.Key, append(opts, t.Style...))
}

// Button describes a themed button. Text is prepended to optional custom Children.
type Button struct {
	Key      string
	Text     string
	Variant  shadcn.Variant
	Size     shadcn.Size
	Disabled bool
	OnClick  func()
	Style    []ui.StyleOpt
	Children []View
}

// Build delegates button appearance, interaction, and lifecycle to pkg/shadcn.
func (b Button) Build() *ui.Node {
	content := make([]*ui.Node, 0, len(b.Children)+1)
	if b.Text != "" {
		content = append(content, ui.Text(b.Text))
	}
	content = append(content, children(b.Children)...)
	return finish(shadcn.Button(shadcn.ButtonProps{
		Variant: b.Variant, Size: b.Size, Disabled: b.Disabled, OnClick: b.OnClick,
	}, content...), b.Key, b.Style)
}

// Input describes a controlled, themed single-line input.
type Input struct {
	Key         string
	Value       string
	Placeholder string
	Disabled    bool
	OnChange    func(string)
	OnSubmit    func(string)
	Style       []ui.StyleOpt
	Bind        values.Bind[string]
}

// Build creates an input whose value is updated through OnChange.
func (i Input) Build() *ui.Node {
	n := shadcn.Input(shadcn.InputProps{
		Value: i.Value, Placeholder: i.Placeholder, Disabled: i.Disabled, OnChange: i.OnChange, Bind: i.Bind,
	})
	if i.OnSubmit != nil {
		n = n.OnSubmit(i.OnSubmit)
	}
	return finish(n, i.Key, i.Style)
}

// Textarea describes a controlled, themed multiline input.
type Textarea struct {
	Key         string
	Value       string
	Placeholder string
	Rows        int
	Disabled    bool
	OnChange    func(string)
	Style       []ui.StyleOpt
	Bind        values.Bind[string]
}

// Build creates a multiline input using the existing Textarea behavior.
func (t Textarea) Build() *ui.Node {
	return finish(shadcn.Textarea(shadcn.TextareaProps{
		Value: t.Value, Placeholder: t.Placeholder, Rows: t.Rows, Disabled: t.Disabled, OnChange: t.OnChange, Bind: t.Bind,
	}), t.Key, t.Style)
}

// Checkbox describes a controlled, themed checkbox.
type Checkbox struct {
	Key      string
	Checked  bool
	Disabled bool
	OnChange func(bool)
	Style    []ui.StyleOpt
}

// Build creates a checkbox with the existing checked and disabled appearance.
func (c Checkbox) Build() *ui.Node {
	return finish(shadcn.Checkbox(shadcn.CheckboxProps{
		Checked: c.Checked, Disabled: c.Disabled, OnChange: c.OnChange,
	}), c.Key, c.Style)
}

// Switch describes a controlled, themed switch using the same fields as Checkbox.
type Switch Checkbox

// Build creates an animated switch using the existing component lifecycle.
func (s Switch) Build() *ui.Node {
	return finish(shadcn.Switch(shadcn.SwitchProps{
		Checked: s.Checked, Disabled: s.Disabled, OnChange: s.OnChange,
	}), s.Key, s.Style)
}

var (
	_ View = Text{}
	_ View = Button{}
	_ View = Input{}
	_ View = Textarea{}
	_ View = Checkbox{}
	_ View = Switch{}
)
