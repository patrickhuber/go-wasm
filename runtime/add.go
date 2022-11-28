package runtime

func I32Add(stack Stack) {
	stack.PushUint32(stack.PopUint32() + stack.PopUint32())
}

func I64Add(stack Stack) {
	stack.PushUint64(stack.PopUint64() + stack.PopUint64())
}

func F32Add(stack Stack) {
	stack.PushFloat32(stack.PopFloat32() + stack.PopFloat32())
}

func F64Add(stack Stack) {
	stack.PushFloat64(stack.PopFloat64() + stack.PopFloat64())
}
