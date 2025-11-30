package unit

type Unit interface {
	Type() string
	Value() float32
}
type percent struct {
	value float32
}

func Percent(v float32) Unit {
	return percent{value: v}
}

func (p percent) Type() string {
	return "percent"
}
func (p percent) Value() float32 {
	return p.value
}

type pt struct {
	value float32
}

func Pt(v float32) Unit {
	return pt{value: v}
}

func (p pt) Type() string {
	return "pt"
}
func (p pt) Value() float32 {
	return p.value
}

type auto struct {
}

func (a auto) Type() string {
	return "auto"
}
func (a auto) Value() float32 {
	return 0
}

func Auto() Unit {
	return auto{}
}
