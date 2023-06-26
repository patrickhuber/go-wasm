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
		size, err := t.Size()
		if err != nil {
			return err
		}
		return StoreInt(c, v, ptr, size, false)
	case kind.U8:
		fallthrough
	case kind.U16:
		fallthrough
	case kind.U32:
		fallthrough
	case kind.U64:
		size, err := t.Size()
		if err != nil {
			return err
		}
		return StoreInt(c, val, ptr, size, false)
	case kind.S8:
		fallthrough
	case kind.S16:
		fallthrough
	case kind.S32:
		fallthrough
	case kind.S64:
		size, err := t.Size()
		if err != nil {
			return err
		}
		return StoreInt(c, val, ptr, size, true)
	case kind.Float32:
		fallthrough
	case kind.Float64:
		size, err := t.Size()
		if err != nil {
			return err
		}
		return StoreFloat(c, val, ptr, size)
	case kind.Char:
		size, err := t.Size()
		if err != nil {
			return err
		}
		return StoreInt(c, uint32(val.(rune)), ptr, size, false)
	case kind.String:
		return StoreString(c, val.(string), ptr)
	case kind.List:
		l := t.(*types.List)
		return StoreList(c, val, ptr, l.Type)
	case kind.Record:
		r := t.(*types.Record)
		return StoreRecord(c, val, ptr, r)
	}
	return types.TrapWith("unrecognized kind %s", t.Kind())
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
	size, _ := types.U32{}.Size()
	return StoreInt(c, val, ptr, size, false)
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
	if dstByteLength > types.MaxStringByteLength {
		return 0, 0, types.TrapWith("destination byte length %d is greater than max string byte length %d", dstByteLength, types.MaxStringByteLength)
	}
	ptr, err := cx.Options.Realloc(0, 0, dstAlignment, dstByteLength)
	if err != nil {
		return 0, 0, err
	}
	if ptr != types.AlignTo(ptr, dstAlignment) {
		return 0, 0, types.TrapWith("ptr %d is not aligned to destination %d", ptr, dstAlignment)
	}
	if ptr+dstByteLength > uint32(cx.Options.Memory.Len()) {
		return 0, 0, types.TrapWith("array size %d is greater than the memory size %d", ptr+dstByteLength, cx.Options.Memory.Len())
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

	if worstCaseSize > types.MaxStringByteLength {
		return 0, 0, types.TrapWith("worst case size %d is greater than max string byte length %d", worstCaseSize, types.MaxStringByteLength)
	}

	ptr, err := cx.Options.Realloc(0, 0, 2, worstCaseSize)
	if err != nil {
		return 0, 0, err
	}

	if ptr != types.AlignTo(ptr, 2) {
		return 0, 0, types.TrapWith("ptr %d is not alinged to 2", ptr)
	}

	if ptr+worstCaseSize > uint32(cx.Options.Memory.Len()) {
		return 0, 0, types.TrapWith("worst case size %d is greater than memory size %d", ptr+worstCaseSize, cx.Options.Memory.Len())
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

		if ptr != types.AlignTo(ptr, 2) {
			return 0, 0, types.TrapWith("ptr %d could not be aligned to 2", ptr)
		}

		if hiPtr > uint32(cx.Options.Memory.Len()) {
			return 0, 0, types.TrapWith("ptr %d is greater than memory size %d", hiPtr, cx.Options.Memory.Len())
		}
	}

	codeUnits := uint32(len(encoded) / 2)
	return ptr, codeUnits, nil
}

func StoreList(cx *types.Context, v any, ptr uint32, elementType types.ValType) error {
	begin, length, err := StoreListIntoRange(cx, v, elementType)
	if err != nil {
		return err
	}
	err = StoreUInt32(cx, begin, ptr)
	if err != nil {
		return err
	}
	return StoreUInt32(cx, length, ptr+4)
}

func StoreListIntoRange(cx *types.Context, v any, elementType types.ValType) (uint32, uint32, error) {
	slice, err := ToSlice(v)
	if err != nil {
		return 0, 0, err
	}

	size, err := elementType.Size()
	if err != nil {
		return 0, 0, err
	}

	byteLengthInt := len(slice) * int(size)
	if byteLengthInt >= (1 << 32) {
		return 0, 0, types.TrapWith("byte length %d exceeds max of %d", byteLengthInt, (1 << 32))
	}
	byteLength := uint32(byteLengthInt)

	alignment, err := elementType.Alignment()
	if err != nil {
		return 0, 0, err
	}

	ptr, err := cx.Options.Realloc(0, 0, alignment, byteLength)
	if err != nil {
		return 0, 0, err
	}

	if ptr != types.AlignTo(ptr, alignment) {
		return 0, 0, types.TrapWith("ptr %d not aligned to %d", ptr, alignment)
	}

	if ptr+byteLength > uint32(cx.Options.Memory.Len()) {
		return 0, 0, types.TrapWith("ptr %d exceeds mememory size %d", ptr+byteLength, cx.Options.Memory.Len())
	}

	for i, element := range slice {
		u32Index := uint32(i)
		err = Store(cx, element, elementType, ptr+u32Index*size)
		if err != nil {
			return 0, 0, err
		}
	}

	return ptr, uint32(len(slice)), nil
}

func ToSlice(val any) ([]any, error) {
	switch v := val.(type) {
	case []any:
		return v, nil
	default:
		return nil, types.NewCastError(val, "[]any")
	}
}

func StoreRecord(cx *types.Context, val any, ptr uint32, r *types.Record) error {
	valMap, err := ToMapStringAny(val)
	if err != nil {
		return err
	}

	for _, f := range r.Fields {
		size, err := StoreField(cx, valMap[f.Label], ptr, f)
		if err != nil {
			return err
		}
		ptr += size
	}

	return nil
}

func StoreField(cx *types.Context, val any, ptr uint32, f types.Field) (uint32, error) {
	alignment, err := f.Type.Alignment()

	if err != nil {
		return 0, err
	}

	ptr = types.AlignTo(ptr, alignment)

	err = Store(cx, val, f.Type, ptr)
	if err != nil {
		return 0, err
	}

	return f.Type.Size()
}

func ToMapStringAny(val any) (map[string]any, error) {
	switch v := val.(type) {
	case map[string]any:
		return v, nil
	}
	return nil, types.NewCastError(val, "map[string]any")
}

func PackFlagsIntoInt(v map[string]any, labels []string) (int, error) {
	packed := 0
	shift := 0
	for _, label := range labels {
		val := v[label]
		b, ok := val.(bool)
		if !ok {
			return 0, types.NewCastError(val, "bool")
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
