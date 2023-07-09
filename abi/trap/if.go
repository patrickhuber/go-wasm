package trap

import "github.com/patrickhuber/go-wasm/abi/types"

func If(condition bool) {
	if !condition {
		return
	}
	panic(types.Trap())
}

func Iff(condition bool, format string, args ...any) {
	if !condition {
		return
	}
	panic(types.TrapWith(format, args...))
}
