package io_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/stretchr/testify/require"
)

func TestCodeUnits(t *testing.T) {
	type test struct {
		integer uint32
		tc      io.TaggedCodeUnits
	}
	tests := []test{
		{
			integer: io.UTF16Tag ^ 10,
			tc: io.TaggedCodeUnits{
				CodeUnits: 10,
				UTF16:     true,
			},
		},
		{
			integer: 10,
			tc: io.TaggedCodeUnits{
				CodeUnits: 10,
				UTF16:     false,
			},
		},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%d", test.integer)
		t.Run(name, func(t *testing.T) {
			tc := io.UInt32ToTaggedCodeUnits(test.integer)
			require.Equal(t, tc, test.tc)
		})

	}
}
