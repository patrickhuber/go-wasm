package io

import (
	"encoding/binary"
	"math"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

func Store(c *types.Context, val any, t types.ValType, ptr uint32) error {
	switch t.Kind() {
	case kind.Bool:
		var v byte = 0
		if val.(bool) {
			v = 1
		}
		return StoreInt(c, v, ptr, t.Size(), false)
	case kind.U8:
		fallthrough
	case kind.U16:
		fallthrough
	case kind.U32:
		fallthrough
	case kind.U64:
		return StoreInt(c, val, ptr, t.Size(), false)
	case kind.S8:
		fallthrough
	case kind.S16:
		fallthrough
	case kind.S32:
		fallthrough
	case kind.S64:
		return StoreInt(c, val, ptr, t.Size(), true)
	case kind.Float32:
		fallthrough
	case kind.Float64:
		return StoreFloat(c, val, ptr, t.Size())
	case kind.Char:
		return StoreInt(c, uint32(val.(rune)), ptr, t.Size(), false)
	case kind.String:
		return StoreString(c, val.(string), ptr)
	case kind.List:
		l := t.(*types.List)
		return StoreList(c, val, ptr, l.Type)
	}
	return types.Trap()
}

func StoreFloat(c *types.Context, val any, ptr uint32, nbytes uint32) error {
	if nbytes == 4 {
		f := val.(float32)
		i := math.Float32bits(f)
		StoreInt(c, i, ptr, nbytes, false)
	} else {
		f := val.(float64)
		i := math.Float64bits(f)
		StoreInt(c, i, ptr, nbytes, false)
	}
	return nil
}

func StoreUInt32(c *types.Context, val uint32, ptr uint32) error {
	return StoreInt(c, val, ptr, 4, false)
}

func StoreInt(c *types.Context, val any, ptr uint32, nbytes uint32, signed bool) error {
	buf := c.Options.Memory[ptr : ptr+nbytes]
	switch nbytes {
	case 1:
		var b byte
		if signed {
			b = byte(val.(int8))
		} else {
			b = val.(byte)
		}
		buf[0] = b

	case 2:
		var u16 uint16
		if signed {
			u16 = uint16(val.(int16))
		} else {
			u16 = val.(uint16)
		}
		binary.LittleEndian.PutUint16(buf, u16)

	case 4:
		var u32 uint32
		if signed {
			u32 = uint32(val.(int32))
		} else {
			u32 = val.(uint32)
		}
		binary.LittleEndian.PutUint32(buf, u32)

	case 8:
		var u64 uint64
		if signed {
			u64 = uint64(val.(int64))
		} else {
			u64 = val.(uint64)
		}
		binary.LittleEndian.PutUint64(buf, u64)
	}
	return nil
}

func StoreString(c *types.Context, str string, ptr uint32) error {
	begin, taggedCodeUnits, err := StoreStringIntoRange(c, str)
	if err != nil {
		return err
	}

	err = StoreUInt32(c, begin, ptr)
	if err != nil {
		return err
	}
	return StoreUInt32(c, taggedCodeUnits, ptr+4)
}

func StoreStringIntoRange(cx *types.Context, str string) (uint32, uint32, error) {
	var encoder *encoding.Encoder

	switch cx.Options.StringEncoding {
	case types.Utf8:
		encoder = unicode.UTF8.NewEncoder()

	case types.Utf16:
		encoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()

	case types.Latin1Utf16:
		encoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
		// 	encoder = charmap.ISO8859_1.NewEncoder()
	default:
		return 0, 0, types.Trap()
	}
	if encoder == nil {
		return 0, 0, nil
	}
	return 0, 0, nil
}

func StoreStringCopy(cx *types.Context, src string, srcCodeUnits uint32, dstCodeUnitSize uint32, dstAlignment uint32, dstEncoding encoding.Encoder) error {

	dstByteLength := dstCodeUnitSize * srcCodeUnits
	err := types.TrapIf(dstByteLength > types.MaxStringByteLength)
	if err != nil {
		return err
	}

	ptr := cx.Options.Realloc(0, 0, dstAlignment, dstByteLength)
	err = types.TrapIf(ptr != types.AlignTo(ptr, dstAlignment))
	if err != nil {
		return err
	}

	err = types.TrapIf(ptr+dstByteLength > uint32(len(cx.Options.Memory)))
	if err != nil {
		return err
	}

	buf := cx.Options.Memory[ptr : ptr+dstByteLength]
	_, _, err = dstEncoding.Transformer.Transform(buf, []byte(src), false)
	if err != nil {
		return err
	}
	return nil

}

func StoreList(c *types.Context, v any, ptr uint32, elementType types.ValType) error {
	return nil
}
