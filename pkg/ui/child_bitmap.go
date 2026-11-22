package ui

import "math/bits"

// childBitmap is a sparse, power-of-two bitmap indexed by a child's position
// in its parent. Empty nodes allocate no words.
type childBitmap struct {
	words []uint64
}

func (b *childBitmap) ensureWord(index uint32) {
	if index < uint32(len(b.words)) {
		return
	}
	size := max(1, len(b.words))
	for size <= int(index) {
		size <<= 1
	}
	b.words = append(b.words, make([]uint64, size-len(b.words))...)
}

func (b *childBitmap) Set(index uint32) {
	word := index / 64
	b.ensureWord(word)
	b.words[word] |= uint64(1) << (index % 64)
}

func (b *childBitmap) Has(index uint32) bool {
	word := index / 64
	return word < uint32(len(b.words)) && b.words[word]&(uint64(1)<<(index%64)) != 0
}

func (b *childBitmap) Delete(index uint32) {
	word := index / 64
	if word < uint32(len(b.words)) {
		b.words[word] &^= uint64(1) << (index % 64)
	}
}

func (b *childBitmap) Reset() {
	for i := range b.words {
		b.words[i] = 0
	}
}

func (b *childBitmap) Any() bool {
	for _, word := range b.words {
		if word != 0 {
			return true
		}
	}
	return false
}

func (b *childBitmap) Range(fn func(index uint32)) {
	for wordIndex, word := range b.words {
		for word != 0 {
			bit := bits.TrailingZeros64(word)
			fn(uint32(wordIndex*64 + bit))
			word &= word - 1
		}
	}
}

type hoverState struct {
	self     [2]bool
	children [2]childBitmap
}

func (rn *renderNode) advanceHover() {
	rn.hover.self = [2]bool{false, rn.hover.self[0]}
	rn.hover.children[0], rn.hover.children[1] = rn.hover.children[1], rn.hover.children[0]
	rn.hover.children[0].Reset()
	for _, child := range rn.children {
		if child.hover.self[0] || child.hover.self[1] ||
			child.hover.children[0].Any() || child.hover.children[1].Any() {
			child.advanceHover()
		}
	}
}

func (rn *renderNode) leavePreviousHover() {
	if rn.hover.self[1] && !rn.hover.self[0] && rn.onHover != nil {
		rn.onHover(false)
	}
	// 命中仍在这条父分支时，也可能已经从子元素退到容器本身，必须继续向下检查。
	rn.hover.children[1].Range(func(index uint32) {
		if index < uint32(len(rn.children)) {
			rn.children[index].leavePreviousHover()
		}
	})
}

func (rn *renderNode) clearDisabledHoverAtLevel(level uint32) {
	if rn.level == level {
		if rn.hover.self[0] && rn.isDisabled() {
			rn.hover.self[0] = false
			if rn.onHover != nil {
				rn.onHover(false)
			}
		}
		return
	}
	if rn.level > level {
		return
	}
	rn.hover.children[0].Range(func(index uint32) {
		if index < uint32(len(rn.children)) {
			rn.children[index].clearDisabledHoverAtLevel(level)
		}
	})
}

func (rn *renderNode) forEachAtLevel(level uint32, fn func(*renderNode)) {
	if rn.level > level {
		return
	}
	if rn.level == level {
		fn(rn)
		return
	}
	for _, child := range rn.children {
		child.forEachAtLevel(level, fn)
	}
}
