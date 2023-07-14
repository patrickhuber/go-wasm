package io

import "github.com/patrickhuber/go-wasm/abi/types"

func CanonResourceNew(inst *types.ComponentInstance, rt types.ResourceType, rep int) (any, error) {
	h := &types.HandleElem{
		Rep: rep,
		Own: true,
	}
	return inst.Handles.Add(rt, h)
}

func CanonResourceRep(inst *types.ComponentInstance, rt types.ResourceType, rep uint32) (uint32, error) {
	h, err := inst.Handles.Get(rt, rep)
	if err != nil {
		return 0, err
	}
	return uint32(h.Rep), nil
}

func CanonResourceDrop(inst *types.ComponentInstance, rt types.ResourceType, i uint32) error {
	h, err := inst.Handles.Remove(rt, i)
	if err != nil {
		return err
	}
	if !h.Own {
		return nil
	}
	if inst != rt.Impl() {
		return types.TrapWith("ComponentInstance != ResourceType.Impl and ResourceType.Impl.MayEnter == false")
	}
	if rt.DTor() != nil {
		rt.DTor()(h.Rep)
	}
	return nil
}
