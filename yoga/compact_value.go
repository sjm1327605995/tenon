package yoga

import (
	"math"
)

type CompactValue struct {
	repr uint32
}

const (
	LOWER_BOUND         float32 = 1.08420217e-19
	UPPER_BOUND_POINT   float32 = 36893485948395847680.0
	UPPER_BOUND_PERCENT float32 = 18446742974197923840.0
	PERCENT_BIT         uint32  = 0x40000000
	AUTO_BITS           uint32  = 0x7faaaaaa
	ZERO_BITS_POINT     uint32  = 0x7f8f0f0f
	ZERO_BITS_PERCENT   uint32  = 0x7f80f0f0
	BIAS                uint32  = 0x20000000
)

func (cv CompactValue) IsUndefined() bool {
	return cv.repr != AUTO_BITS && cv.repr != ZERO_BITS_POINT && cv.repr != ZERO_BITS_PERCENT && IsNaN(math.Float32frombits(cv.repr))
}

func (cv CompactValue) IsDefined() bool {
	return !cv.IsUndefined()
}

func (cv CompactValue) IsAuto() bool {
	return cv.repr == AUTO_BITS
}

func (cv CompactValue) Equal(other CompactValue) bool {
	return cv.repr == other.repr
}

func (cv CompactValue) Value() Value {
	if cv.IsUndefined() {
		return ValueUndefined
	}
	switch cv.repr {
	case AUTO_BITS:
		return ValueAuto
	case ZERO_BITS_POINT:
		return Value{value: 0.0, unit: UnitPoint}
	case ZERO_BITS_PERCENT:
		return Value{value: 0.0, unit: UnitPercent}
	}

	if IsNaN(math.Float32frombits(cv.repr)) {
		return ValueUndefined
	}

	data := cv.repr
	data &= ^PERCENT_BIT
	data += BIAS

	return Value{value: math.Float32frombits(data), unit: If(cv.repr&0x40000000 != 0, UnitPercent, UnitPoint)}
}

func CompactValueOf(value float32, unit Unit) CompactValue {
	if IsNaN(value) || IsInf(value, 0) {
		return CompactValueOfUndefined()
	}
	if value == 0.0 || (value < LOWER_BOUND && value > -LOWER_BOUND) {
		zero := ZERO_BITS_POINT
		if unit == UnitPercent {
			zero = ZERO_BITS_PERCENT
		}
		return CompactValue{repr: zero}
	}

	upperBound := UPPER_BOUND_POINT
	if unit == UnitPercent {
		upperBound = UPPER_BOUND_PERCENT
	}
	if value > upperBound || value < -upperBound {
		value = float32(math.Copysign(float64(upperBound), float64(value)))
	}

	unitBit := uint32(0)
	if unit == UnitPercent {
		unitBit = PERCENT_BIT
	}
	data := math.Float32bits(value)
	data -= BIAS
	data |= unitBit
	return CompactValue{repr: data}
}

func CompactValueOfPoint(value float32) CompactValue {
	return CompactValueOf(value, UnitPoint)
}

func CompactValuePercent(value float32) CompactValue {
	return CompactValueOf(value, UnitPercent)
}

func CompactValueOfUndefined() CompactValue {
	return CompactValue{repr: 0x7FC00000}
}

func CompactValueOfAuto() CompactValue {
	return CompactValue{repr: AUTO_BITS}
}
