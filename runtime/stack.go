package runtime

import (
	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/instance"
	"github.com/patrickhuber/go-wasm/values"
)

// Stack holds the runtime state of execution
// the stack contains Values, Labels and Activations
type Stack struct {
	Values      []values.Value
	Labels      []api.Label
	Activations []Frame
}

type Frame struct {
	FrameState  *FrameState
	InnerFrames []Frame
}

type FrameState struct {
	Locals []values.Value
	Module instance.Module
}
