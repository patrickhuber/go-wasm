package io

import (
	"encoding/binary"
	"math"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
)

func Store(c *types.Context, val any, t types.ValType, ptr uint32) {
	switch t.Kind() {
	case kind.Bool:
		var v byte = 0
		if val.(bool) {
			v = 1
		}
		StoreInt(c, v, ptr, t.Size())
	case kind.U8:
		fallthrough
	case kind.U16:
		fallthrough
	case kind.U32:
		fallthrough
	case kind.U64:
		fallthrough
	case kind.S8:
		fallthrough
	case kind.S16:
		fallthrough
	case kind.S32:
		fallthrough
	case kind.S64:
		StoreInt(c, val, ptr, t.Size())
	case kind.Float32:
		fallthrough
	case kind.Float64:
		StoreFloat(c, val, ptr, t.Size())
	case kind.Char:
		StoreInt(c, uint32(val.(rune)), ptr, t.Size())
	case kind.String:
		StoreString(c, val.(string), ptr)
	case kind.List:
		l := t.(*types.List)
		storeList(c, val, ptr, l.Type)
	}
}

func StoreFloat(c *types.Context, val any, ptr uint32, nbytes uint32) {
	if nbytes == 4 {
		f := val.(float32)
		i := math.Float32bits(f)
		StoreInt(c, i, ptr, nbytes)
	} else {
		f := val.(float64)
		i := math.Float64bits(f)
		StoreInt(c, i, ptr, nbytes)
	}
}

func StoreInt(c *types.Context, val any, ptr uint32, nbytes uint32) {
	buf := c.Options.Memory[ptr : ptr+nbytes]
	switch nbytes {
	case 1:
		v := val.(byte)
		buf[0] = byte(v)
	case 2:
		v := val.(uint16)
		binary.LittleEndian.PutUint16(buf, v)
	case 4:
		v := val.(uint32)
		binary.LittleEndian.PutUint32(buf, v)
	case 8:
		v := val.(uint64)
		binary.LittleEndian.PutUint64(buf, v)
	}
}

func StoreString(c *types.Context, str string, ptr uint32) {
	begin, taggedCodeUnits := storeStringIntoRange(c, str)
	StoreInt(c, begin, ptr, 4)
	StoreInt(c, taggedCodeUnits, ptr+4, 4)
}

func storeStringIntoRange(c *types.Context, str string) (uint32, uint32) {
	return 0, 0
}

func storeList(c *types.Context, v any, ptr uint32, elementType types.ValType) {

}
