package yoga

var (
	undefinedFloatOptional = FloatOptional{value: NaN}
)

type FloatOptional struct {
	value float32
}

func NewFloatOptional(value float32) FloatOptional {
	return FloatOptional{value: value}
}

func (opt FloatOptional) unwrap() float32 {
	return opt.value
}

func (opt FloatOptional) isUndefined() bool {
	return IsNaN(opt.value)
}

func (opt FloatOptional) isDefined() bool {
	return !IsNaN(opt.value)
}

func (opt FloatOptional) unwrapOrDefault(defaultValue float32) float32 {
	return If(opt.isUndefined(), defaultValue, opt.value)
}

func (opt FloatOptional) equal(other FloatOptional) bool {
	return opt.value == other.value || (opt.isUndefined() && other.isUndefined())
}
