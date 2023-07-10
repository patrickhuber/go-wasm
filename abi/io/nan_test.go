package io_test

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/values"
)

func TestNan32(t *testing.T) {
	type test struct {
		inbits  uint32
		outbits uint32
	}
	tests := []test{
		{0x7fc00000, values.Float32Nan},
		{0x7fc00001, values.Float32Nan},
		{0x7fe00000, values.Float32Nan},
		{0x7fffffff, values.Float32Nan},
		{0xffffffff, values.Float32Nan},
		{0x7f800000, 0x7f800000},
		{0x3fc00000, 0x3fc00000},
	}
	for _, test := range tests {
		t.Run(strconv.Itoa(int(test.inbits)), func(t *testing.T) {
			f := math.Float32frombits(test.inbits)
			a, err := io.LiftFlat(Context(), values.NewIterator(values.Float32(f)), Float32())
			if err != nil {
				t.Fatal(err)
			}
			f, ok := a.(float32)
			if !ok {
				t.Fatalf("expected type float32 but found %T", a)
			}
			if math.Float32bits(f) != test.outbits {
				t.Fatalf("expected %d to equal %d", math.Float32bits(f), test.outbits)
			}

			buf := binary.LittleEndian.AppendUint32([]byte{}, test.inbits)
			cx := Context(CanonicalOptions(Memory(bytes.NewBuffer(buf))))

			a, err = io.Load(cx, Float32(), 0)
			if err != nil {
				t.Fatal(err)
			}
			f, ok = a.(float32)
			if !ok {
				t.Fatalf("expected type float32 but found %T", a)
			}
			if math.Float32bits(f) != test.outbits {
				t.Fatalf("expected %d to equal %d", math.Float32bits(f), test.outbits)
			}
		})
	}
}

func TestNan64(t *testing.T) {
	type test struct {
		inbits  uint64
		outbits uint64
	}
	tests := []test{
		{0x7ff8000000000000, values.Float64Nan},
		{0x7ff8000000000001, values.Float64Nan},
		{0x7ffc000000000000, values.Float64Nan},
		{0x7fffffffffffffff, values.Float64Nan},
		{0xffffffffffffffff, values.Float64Nan},
		{0x7ff0000000000000, 0x7ff0000000000000},
		{0x3ff0000000000000, 0x3ff0000000000000},
	}
	for _, test := range tests {
		t.Run(strconv.Itoa(int(test.inbits)), func(t *testing.T) {
			f := math.Float64frombits(test.inbits)
			a, err := io.LiftFlat(Context(), values.NewIterator(values.Float64(f)), Float64())
			if err != nil {
				t.Fatal(err)
			}
			f, ok := a.(float64)
			if !ok {
				t.Fatalf("expected type float64 but found %T", a)
			}
			if math.Float64bits(f) != test.outbits {
				t.Fatalf("expected %d to equal %d", math.Float64bits(f), test.outbits)
			}

			buf := binary.LittleEndian.AppendUint64([]byte{}, test.inbits)
			cx := Context(CanonicalOptions(Memory(bytes.NewBuffer(buf))))

			a, err = io.Load(cx, Float64(), 0)
			if err != nil {
				t.Fatal(err)
			}
			f, ok = a.(float64)
			if !ok {
				t.Fatalf("expected type float64 but found %T", a)
			}
			if math.Float64bits(f) != test.outbits {
				t.Fatalf("expected %d to equal %d", math.Float64bits(f), test.outbits)
			}
		})
	}
}
