package runtime

import (
	"bufio"
	"fmt"
	"io"

	"github.com/patrickhuber/go-wasm/leb128"
	"github.com/patrickhuber/go-wasm/opcode"
)

type interpreter struct {
	stack []uint32
}

type Interpreter interface {
	Run(r io.Reader) error
}

func NewInterpreter() Interpreter {
	return &interpreter{}
}

func (i *interpreter) Run(r io.Reader) error {
	reader := bufio.NewReader(r)
	code, _, err := leb128.Decode(reader)
	if err != nil {
		return err
	}
	switch opcode.Opcode(code) {

	case opcode.Unreachable:
		return fmt.Errorf("unreachable")

	case opcode.Call:
		// get the function from the function instance
		// push values onto the stack
		// run the function body instructions
	case opcode.Drop:
		_ = i.pop()

	case opcode.I32Const:
		u32, _, err := leb128.Decode(reader)
		if err != nil {
			return err
		}
		i.push(u32)

	case opcode.I32Add:
		lhs := i.pop()
		rhs := i.pop()
		i.push(lhs + rhs)

	}
	return nil
}

func (i *interpreter) pop() uint32 {
	if len(i.stack) == 0 {
		return uint32(opcode.Unreachable)
	}
	v := i.stack[len(i.stack)-1]
	i.stack = i.stack[:len(i.stack)-1]
	return v
}

func (i *interpreter) push(v uint32) {
	i.stack = append(i.stack, v)
}
