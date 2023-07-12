package io

import "github.com/patrickhuber/go-wasm/abi/types"

func CanonResourceNew(inst *types.ComponentInstance, rt types.ResourceType, rep int) (any, error) {
	h := &types.HandleElem{
		Rep: rep,
		Own: true,
	}
	return inst.Handles.Add(h, rt)
}

func CanonResourceRep(inst *types.ComponentInstance, rt types.ResourceType, rep uint32) (uint32, error) {
	h, err := inst.Handles.Get(rt, rep)
	if err != nil {
		return 0, err
	}
	return uint32(h.Rep), nil
}
