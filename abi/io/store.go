package io

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/encoding"
	"golang.org/x/text/encoding/charmap"
)

func Store(c *types.CallContext, val any, t types.ValType, ptr uint32) error {
	if err := StoreValidate(c, t, ptr); err != nil {
		return err
	}

	t = t.Despecialize()
	switch t.Kind() {
	case kind.Bool:
		var v byte = 0
		b, ok := val.(bool)
		if !ok {
			return types.NewCastError(val, "bool")
		}
		if b {
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
		r, ok := val.(rune)
		if !ok {
			return types.NewCastError(val, "rune")
		}
		return StoreInt(c, uint32(r), ptr, size, false)
	case kind.String:
		s, ok := val.(string)
		if !ok {
			return types.NewCastError(val, "string")
		}
		return StoreString(c, s, ptr)
	case kind.List:
		l := t.(*types.List)
		return StoreList(c, val, ptr, l.Type)
	case kind.Record:
		r := t.(*types.Record)
		return StoreRecord(c, val, ptr, r)
	case kind.Variant:
		v := t.(*types.Variant)
		return StoreVariant(c, val, ptr, v)
	}
	return types.TrapWith("Store: unrecognized kind.%s", t.Kind())
}

func StoreValidate(c *types.CallContext, t types.ValType, ptr uint32) error {
	alignment, err := t.Alignment()
	if err != nil {
		return err
	}

	if ptr != types.AlignTo(ptr, alignment) {
		return fmt.Errorf("Store: ptr %d is not aligned to %d", ptr, alignment)
	}

	size, err := t.Size()
	if err != nil {
		return err
	}

	if ptr+size > uint32(c.Options.Memory.Len()) {
		return fmt.Errorf("Store: %d exceeds memory length %d", ptr+size, c.Options.Memory.Len())
	}

	return err
}

func StoreFloat(c *types.CallContext, val any, ptr uint32, nbytes uint32) error {
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

func StoreUInt32(c *types.CallContext, val uint32, ptr uint32) error {
	buf := c.Options.Memory.Bytes()[ptr : ptr+types.SizeOfU32]
	binary.LittleEndian.PutUint32(buf, val)
	return nil
}

func StoreInt(c *types.CallContext, val any, ptr uint32, nbytes uint32, signed bool) error {
	var u64 uint64
	var max uint64
	var min int64 = 0
	sign := false
	switch t := val.(type) {
	case uint8:
		u64 = uint64(t)
		max = math.MaxUint8
	case uint16:
		u64 = uint64(t)
		max = math.MaxUint16
	case uint32:
		u64 = uint64(t)
		max = math.MaxUint32
	case uint64:
		u64 = t
		max = math.MaxUint64
	case int8:
		u64 = uint64(t)
		sign = true
		max = math.MaxInt8
		min = math.MinInt8
	case int16:
		u64 = uint64(t)
		sign = true
		max = math.MaxInt16
		min = math.MinInt16
	case int32:
		u64 = uint64(t)
		sign = true
		max = math.MaxInt32
		min = math.MinInt32
	case int64:
		u64 = uint64(t)
		sign = true
		max = math.MaxInt64
		min = math.MinInt64
	}

	if sign != signed {
		signCh := "+"
		if sign {
			signCh = "-"
		}
		return fmt.Errorf("invalid integer sign %v for value %v", signCh, val)
	}

	if !signed && u64 > max {
		return fmt.Errorf("invalid integer %d exceeds max value %d", val, max)
	}
	if signed {
		if int64(u64) > int64(max) {
			return fmt.Errorf("invalid integer %d exceeds max value %d", val, int64(max))
		}
		if int64(u64) < min {
			return fmt.Errorf("invalid integer %d exceeds min value %d", val, min)
		}
	}

	buf := c.Options.Memory.Bytes()[ptr : ptr+nbytes]
	switch nbytes {
	case types.SizeOfS8:
		buf[0] = uint8(u64)
	case types.SizeOfS16:
		binary.LittleEndian.PutUint16(buf, uint16(u64))
	case types.SizeOfS32:
		binary.LittleEndian.PutUint32(buf, uint32(u64))
	case types.SizeOfS64:
		binary.LittleEndian.PutUint64(buf, u64)
	}
	return nil
}

// StoreString stores the string to linear memory using the context encoding
// All strings in go are assumed to be utf8 encoded
func StoreString(c *types.CallContext, str string, ptr uint32) error {

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

func StoreStringIntoRange(cx *types.CallContext, str string) (uint32, uint32, error) {

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
	cx *types.CallContext,
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

func StoreStringCopy(cx *types.CallContext, src string, srcCodeUnits uint32, dstCodeUnitSize uint32, dstAlignment uint32, dstEncoding encoding.Encoder) (uint32, uint32, error) {

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

func StoreUtf8ToUtf16(cx *types.CallContext, src string, srcCodeUnits uint32) (uint32, uint32, error) {

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

func StoreList(cx *types.CallContext, v any, ptr uint32, elementType types.ValType) error {
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

func StoreListIntoRange(cx *types.CallContext, v any, elementType types.ValType) (uint32, uint32, error) {
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

func StoreRecord(cx *types.CallContext, val any, ptr uint32, r *types.Record) error {
	valMap, err := ToMapStringAny(val)
	if err != nil {
		return err
	}
	for _, f := range r.Fields {
		alignment, err := f.Type.Alignment()
		if err != nil {
			return err
		}

		ptr = types.AlignTo(ptr, alignment)

		err = Store(cx, valMap[f.Label], f.Type, ptr)
		if err != nil {
			return err
		}

		size, err := f.Type.Size()
		if err != nil {
			return err
		}

		ptr += size
	}
	return nil
}

func StoreVariant(cx *types.CallContext, val any, ptr uint32, v *types.Variant) error {
	caseIndex, caseValue, err := MatchCase(val, v.Cases)
	if err != nil {
		return err
	}
	dt, err := v.DiscriminantType()
	if err != nil {
		return err
	}
	size, err := dt.Size()
	if err != nil {
		return err
	}
	err = StoreInt(cx, caseIndex, ptr, size, false)
	if err != nil {
		return err
	}
	ptr += size
	alignment, err := v.MaxCaseAlignment()
	if err != nil {
		return err
	}
	ptr = types.AlignTo(ptr, alignment)
	c := v.Cases[caseIndex]
	if c.Type == nil {
		return nil
	}
	return Store(cx, caseValue, c.Type, ptr)
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
