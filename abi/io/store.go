package io

import (
	"encoding/binary"
	"math"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/encoding"
	"golang.org/x/text/encoding/charmap"
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
		return StoreInt(c, i, ptr, nbytes, false)
	} else {
		f := val.(float64)
		i := math.Float64bits(f)
		return StoreInt(c, i, ptr, nbytes, false)
	}
}

func StoreUInt32(c *types.Context, val uint32, ptr uint32) error {
	return StoreInt(c, val, ptr, 4, false)
}

func StoreInt(c *types.Context, val any, ptr uint32, nbytes uint32, signed bool) error {
	buf := c.Options.Memory.Bytes()[ptr : ptr+nbytes]
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

// StoreString stores the string to linear memory using the context encoding
// All strings in go are assumed to be utf8 encoded
func StoreString(c *types.Context, str string, ptr uint32) error {

	// string storage in wasm components stores the string first
	begin, taggedCodeUnits, err := StoreStringIntoRange(c, str)
	if err != nil {
		return err
	}

	// once the string is stored the pointer and code units are stored next
	err = StoreUInt32(c, begin, ptr)
	if err != nil {
		return err
	}
	return StoreUInt32(c, taggedCodeUnits, ptr+4)
}

func StoreStringIntoRange(cx *types.Context, str string) (uint32, uint32, error) {

	codec, err := encoding.DefaultFactory().Get(cx.Options.StringEncoding)
	if err != nil {
		return 0, 0, err
	}

	if cx.Options.StringEncoding != encoding.Latin1Utf16 {
		return StoreStringDynamic(cx, str, codec)
	}

	// set the default encoding
	enc := encoding.Latin1
	if !isLatin1(str) {
		enc = encoding.UTF16LE
	}

	codec, err = encoding.DefaultFactory().Get(enc)
	if err != nil {
		return 0, 0, err
	}

	ptr, size, err := StoreStringDynamic(cx, str, codec)
	if err != nil {
		return 0, 0, err
	}

	tcu := TaggedCodeUnits{
		CodeUnits: size,
		UTF16:     enc == encoding.UTF16LE,
	}

	return ptr, tcu.ToUInt32(), nil
}

func isLatin1(str string) bool {
	latin1 := charmap.ISO8859_1

	// scan the string looking for invalid latin1 characters
	for _, r := range str {
		_, ok := latin1.EncodeRune(r)
		if !ok {
			return false
		}
	}
	return true
}

// StoreStringDynamic assumes the incoming string is in utf8 and stores the string to the given codec's encoding at the end of the context memory
func StoreStringDynamic(
	cx *types.Context,
	str string,
	codec encoding.Codec) (uint32, uint32, error) {

	encoded, err := encoding.EncodeString(codec, str)
	if err != nil {
		return 0, 0, err
	}

	dstAlignment := uint32(codec.Alignment())
	lenEncoded := uint32(len(encoded))

	ptr, err := cx.Options.Realloc(0, 0, dstAlignment, lenEncoded)
	if err != nil {
		return 0, 0, err
	}

	buf := cx.Options.Memory.Bytes()[ptr : ptr+lenEncoded]
	copy(buf, encoded)

	// return the pointer and the adjusted length (in runes)
	return ptr, lenEncoded / uint32(codec.RuneSize()), nil
}

func StoreStringCopy(cx *types.Context, src string, srcCodeUnits uint32, dstCodeUnitSize uint32, dstAlignment uint32, dstEncoding encoding.Encoder) (uint32, uint32, error) {

	dstByteLength := dstCodeUnitSize * srcCodeUnits
	err := types.TrapIf(dstByteLength > types.MaxStringByteLength)
	if err != nil {
		return 0, 0, err
	}

	ptr, err := cx.Options.Realloc(0, 0, dstAlignment, dstByteLength)
	if err != nil {
		return 0, 0, err
	}
	err = types.TrapIf(ptr != types.AlignTo(ptr, dstAlignment))
	if err != nil {
		return 0, 0, err
	}

	err = types.TrapIf(ptr+dstByteLength > uint32(cx.Options.Memory.Len()))
	if err != nil {
		return 0, 0, err
	}

	encoded, err := encoding.EncodeString(dstEncoding, src)
	if err != nil {
		return 0, 0, err
	}

	buf := cx.Options.Memory.Bytes()[ptr : ptr+dstByteLength]
	copy(buf, encoded)

	return ptr, srcCodeUnits, err
}

func StoreUtf8ToUtf16(cx *types.Context, src string, srcCodeUnits uint32) (uint32, uint32, error) {

	worstCaseSize := 2 * srcCodeUnits

	err := types.TrapIf(worstCaseSize > types.MaxStringByteLength)
	if err != nil {
		return 0, 0, err
	}

	ptr, err := cx.Options.Realloc(0, 0, 2, worstCaseSize)
	if err != nil {
		return 0, 0, err
	}

	err = types.TrapIf(ptr != types.AlignTo(ptr, 2))
	if err != nil {
		return 0, 0, err
	}

	err = types.TrapIf(ptr+worstCaseSize > uint32(cx.Options.Memory.Len()))
	if err != nil {
		return 0, 0, err
	}

	encoded, err := encoding.EncodeString(encoding.NewUTF16(), src)
	if err != nil {
		return 0, 0, err
	}

	hiPtr := ptr + uint32(len(encoded))
	buf := cx.Options.Memory.Bytes()[ptr:hiPtr]
	copy(buf, encoded)

	if len(encoded) < int(worstCaseSize) {

		ptr, err = cx.Options.Realloc(ptr, worstCaseSize, 2, uint32(len(encoded)))
		if err != nil {
			return 0, 0, err
		}

		err = types.TrapIf(ptr != types.AlignTo(ptr, 2))
		if err != nil {
			return 0, 0, err
		}

		err = types.TrapIf(hiPtr > uint32(cx.Options.Memory.Len()))
		if err != nil {
			return 0, 0, err
		}
	}

	codeUnits := uint32(len(encoded) / 2)
	return ptr, codeUnits, nil
}

func StoreUtf8ToLatin1OrUtf16(cx *types.Context, src string, srcCodeUnits uint32) (uint32, uint32, error) {
	return 0, 0, nil
}

func StoreList(c *types.Context, v any, ptr uint32, elementType types.ValType) error {
	return nil
}

func PackFlagsIntoInt(v map[string]any, labels []string) (int, error) {
	packed := 0
	shift := 0
	for _, label := range labels {
		val := v[label]
		b, ok := val.(bool)
		if !ok {
			return 0, NewCastError(val, "bool")
		}
		i := 0
		if b {
			i = 1
		}
		packed |= i << shift
		shift += 1
	}

	return packed, nil
}
