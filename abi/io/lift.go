package io

import "github.com/patrickhuber/go-wasm/abi/types"

func LiftOwn(cx *types.Context, i uint32, t *types.Own) (*types.Handle, error) {
	return cx.Instance.Handles.Transfer(i, t)
}

func LiftBorrow(cx *types.Context, i uint32, t *types.Borrow) (*types.Handle, error) {
	h, err := cx.Instance.Handles.Get(i, t.ResourceType)
	if err != nil {
		return nil, err
	}
	h.LendCount += 1
	cx.Lenders = append(cx.Lenders, h)
	return &types.Handle{
		Rep:       h.Rep,
		LendCount: 0,
		Context:   nil,
	}, nil
}
