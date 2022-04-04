package leb128

import "bufio"

func Encode(w *bufio.Writer, value uint32) (int, error) {
	total := 0
	for {
		b := byte(value & 0b_0111_1111)
		value >>= 7
		if value != 0 {
			b |= 0b_1000_0000
		}
		err := w.WriteByte(b)
		if err != nil {
			return 0, err
		}
		total++
		if value == 0 {
			break
		}
	}
	return total, nil
}
