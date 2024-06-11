package common

const (
	MaxUserSize  = 1 << 25
	BitSlotSize  = 1 << 10
	BytesForSlot = BitSlotSize >> 3
)

type BitSlot [BytesForSlot]byte

func (bs *BitSlot) Add(pos uint16) {
	ByteIdx := pos / 8
	BitIdx := pos % 8
	bs[ByteIdx] |= 1 << BitIdx
}
func (bs *BitSlot) Clear(pos uint16) {
	ByteIdx := pos / 8
	BitIdx := pos % 8
	bs[ByteIdx] &^= 1 << BitIdx
}

func (bs *BitSlot) Has(pos uint16) bool {
	ByteIdx := pos / 8
	BitIdx := pos % 8
	return bs[ByteIdx]&(1<<BitIdx) != 0
}

func (bs *BitSlot) Union(bs2 BitSlot) BitSlot {
	var result BitSlot
	for i := range bs {
		result[i] = bs[i] | bs2[i]
	}
	return result
}

func (bs *BitSlot) Intersection(bs2 BitSlot) BitSlot {
	var result BitSlot
	for i := range bs {
		result[i] = bs[i] & bs2[i]
	}
	return result
}

func (bs *BitSlot) GetBitsAsArray() []uint16 {
	var result []uint16
	for i, byteVal := range bs {
		for j := 0; j < 8; j++ {
			if byteVal&(1<<j) != 0 {
				bitPos := uint16(i*8 + j)
				result = append(result, bitPos)
			}
		}
	}
	return result
}
