package leb128

import (
	"bufio"
	"fmt"
	"io"
)

func DecodeReader(r io.Reader) (uint32, int, error) {
	var val uint32
	shift := 0
	total := 0
	buf := make([]byte, 1)

	for {
		n, err := r.Read(buf)
		if n == 0 {
			return 0, 0, fmt.Errorf("expected 1 byte read but read 0")
		}
		if err != nil {
			return 0, 0, err
		}
		b := buf[0]
		total++
		val |= (uint32(b&0b_0111_1111) << shift)
		if b&0b_1000_0000 == 0 {
			break
		}
		shift += 7
	}
	return val, total, nil
}

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

func DecodeSlice(s []byte) (uint32, int, error) {
	var val uint32
	shift := 0
	total := 0
	for index := 0; index < len(s); index++ {
		b := s[index]
		total++
		val |= (uint32(b&0b_0111_1111) << shift)
		if b&0b_1000_0000 == 0 {
			break
		}
		shift += 7
	}
	return val, total, nil
}
