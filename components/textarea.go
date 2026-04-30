package components

// Textarea is a multi-line text input.
type Textarea struct {
	*TextInput
}

// NewTextarea creates a multi-line text input.
func NewTextarea() *Textarea {
	t := NewTextInput()
	t.SetMultiline(true)
	t.SetHeight(80)
	return &Textarea{TextInput: t}
}

// ElementType returns type identifier.
func (t *Textarea) ElementType() string { return "Textarea" }

// SetRows sets a visual hint for minimum height (approximate).
func (t *Textarea) SetRows(rows int) *Textarea {
	t.SetHeight(float32(rows) * 24)
	return t
}
