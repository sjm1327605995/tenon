package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// InputOTP is a one-time password input with multiple boxes.
type InputOTP struct {
	core.BaseElement
	length     int
	boxes      []*TextInput
	onComplete func(value string)
}

// NewInputOTP creates an OTP input.
func NewInputOTP(length int) *InputOTP {
	if length < 1 {
		length = 6
	}
	i := &InputOTP{length: length}
	i.Init(i)
	i.SetFlexDirection(yoga.FlexDirectionRow)
	i.SetGap(yoga.GutterAll, 8)

	for idx := 0; idx < length; idx++ {
		box := NewTextInput()
		box.SetWidth(40)
		box.SetHeight(40)
		i.boxes = append(i.boxes, box)
		i.Add(box)
	}
	return i
}

// ElementType returns type identifier.
func (i *InputOTP) ElementType() string { return "InputOTP" }

// SetOnComplete sets the completion callback.
func (i *InputOTP) SetOnComplete(fn func(value string)) *InputOTP {
	i.onComplete = fn
	return i
}
