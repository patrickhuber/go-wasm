package leb128

import "bufio"

func Decode(r *bufio.Reader) (uint32, int, error) {
	var val uint32
	shift := 0
	total := 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		total++
		val |= (uint32(b&0b_0111_1111) << shift)
		if b&0b_1000_0000 == 0 {
			break
		}
		shift += 7
	}
	return val, total, nil
}
