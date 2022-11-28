package wasm

import (
	"bufio"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/patrickhuber/go-wasm/leb128"
)

func ReadUtf8String(reader *bufio.Reader) (string, int, error) {
	size, read, err := leb128.Decode(reader)
	if err != nil {
		return "", 0, err
	}
	buf := make([]byte, size)
	if _, err = io.ReadFull(reader, buf); err != nil {
		return "", 0, fmt.Errorf("failed to read string %w", err)
	}
	if !utf8.Valid(buf) {
		return "", 0, fmt.Errorf("invalid utf8 string")
	}
	return string(buf), read + int(size), nil
}

func DecodeUtf8String(buf []byte) (string, int, error) {
	size, read, err := leb128.DecodeSlice(buf)
	if err != nil {
		return "", 0, err
	}
	limit := read + int(size)
	if !utf8.Valid(buf[read:limit]) {
		return "", 0, fmt.Errorf("invalid utf8 string")
	}
	return string(buf[read:limit]), limit, nil
}
