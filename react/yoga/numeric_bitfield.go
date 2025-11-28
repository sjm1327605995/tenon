package yoga

import "math/bits"

type EnumBitset struct {
	bits uint
}

func NewEnumBitset() EnumBitset {
	return EnumBitset{bits: 0}
}

func (bitset *EnumBitset) Set(index uint) {
	bitset.bits |= 1 << index
}

func (bitset *EnumBitset) Reset(index uint) {
	bitset.bits &= ^(1 << index)
}

func (bitset *EnumBitset) Test(index uint) bool {
	return (bitset.bits & (1 << index)) != 0
}

func (bitset *EnumBitset) Count() int {
	return bits.OnesCount(bitset.bits)
}
