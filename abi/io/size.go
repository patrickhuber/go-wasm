package io

import "github.com/patrickhuber/go-wasm/abi/types"

func Size(vt types.ValType) (uint32, error) {
	switch t := vt.(type) {
	case types.Bool:
		return 1, nil
	case types.UInt8:
		return 1, nil
	case types.UInt16:
		return 2, nil
	case types.UInt32:
		return 4, nil
	case types.Int8:
		return 1, nil
	case types.Int16:
		return 2, nil
	case types.Int32:
		return 4, nil
	case types.Record:
		return SizeRecord(t)
	}
	return 0, nil
}

func SizeRecord(r types.Record) (uint32, error) {
	var s uint32 = 0
	for _, f :=range r.Fields(){
		
	}
	return 0, nil
}
